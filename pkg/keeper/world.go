package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"maps"
	"math"
	"slices"

	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/v2/pkg/protocol/metadata"
	"github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/v2/pkg/protocol/slot"
)

// world is the loaded part of the current dimension: chunks with the block
// changes since they were sent, entities and map contents.
//
// Chunks stay as the server's raw packet plus a map of changed blocks instead
// of being decoded, so the keeper does not have to understand the chunk format
// and memory is bounded by the view distance.
type world struct {
	dimension string
	chunks    map[chunkPos]*chunk
	entities  map[int32]*entity
	// maps survive dimension changes: map ids are global to the server.
	maps map[int32]*mapState
}

type chunkPos struct{ x, z int32 }

type blockPos struct{ x, y, z int32 }

type chunk struct {
	data          []byte
	light         []byte
	blocks        map[blockPos]int32
	blockEntities map[blockPos][]byte
}

// entity keeps the server's bytes: the spawn packet as sent, metadata entries
// and equipment slots by index, and the latest of the other packets. Only the
// position is tracked as numbers, replayed as a position sync after the spawn.
type entity struct {
	// spawn is nil for the player itself, which has no spawn packet.
	spawn      []byte
	x, y, z    float64
	yaw, pitch float32
	moved      bool
	head       []byte
	meta       map[uint8][]byte
	equipment  map[int8][]byte
	attributes []byte
	passengers []byte
	link       []byte
}

func newEntity(spawn []byte) *entity {
	return &entity{spawn: spawn, meta: map[uint8][]byte{}, equipment: map[int8][]byte{}}
}

type mapState struct {
	scale     int8
	locked    bool
	icons     []byte // the icon array as the server sent it
	hasIcons  bool
	colors    [128 * 128]byte
	hasColors bool
}

func newWorld() world {
	return world{
		chunks:   make(map[chunkPos]*chunk),
		entities: make(map[int32]*entity),
		maps:     make(map[int32]*mapState),
	}
}

// fromServer records the packets that describe the world. handled is false for
// packets that are not about the world.
func (w *world) fromServer(id packetid.ClientboundPacketID, p pk.Packet, self int32) (handled bool, err error) {
	switch id {
	case packetid.ClientboundLevelChunkWithLight:
		if len(p.Data) < 8 {
			return true, fmt.Errorf("chunk packet of %d bytes", len(p.Data))
		}
		pos := chunkPos{int32(binary.BigEndian.Uint32(p.Data)), int32(binary.BigEndian.Uint32(p.Data[4:]))}
		w.chunks[pos] = &chunk{data: p.Data, blocks: map[blockPos]int32{}, blockEntities: map[blockPos][]byte{}}
	case packetid.ClientboundForgetLevelChunk:
		var v pk.Long
		if _, err := v.ReadFrom(bytes.NewReader(p.Data)); err != nil {
			return true, err
		}
		delete(w.chunks, chunkPos{int32(v), int32(v >> 32)})
	case packetid.ClientboundLightUpdate:
		r := bytes.NewReader(p.Data)
		var x, z pk.VarInt
		if _, err := x.ReadFrom(r); err != nil {
			return true, err
		}
		if _, err := z.ReadFrom(r); err != nil {
			return true, err
		}
		if c := w.chunks[chunkPos{int32(x), int32(z)}]; c != nil {
			c.light = p.Data
		}
	case packetid.ClientboundBlockUpdate:
		var u client.BlockUpdate
		if err := decode(p, &u); err != nil {
			return true, err
		}
		w.setBlock(blockPos{int32(u.Position.X), int32(u.Position.Y), int32(u.Position.Z)}, u.BlockState)
	case packetid.ClientboundSectionBlocksUpdate:
		var u client.UpdateSectionsBlocks
		if err := decode(p, &u); err != nil {
			return true, err
		}
		sx := int32(u.ChunkSectionPosition >> 42)
		sy := int32(u.ChunkSectionPosition << 44 >> 44)
		sz := int32(u.ChunkSectionPosition << 22 >> 42)
		for _, entry := range u.Blocks {
			local := int32(entry & 0xfff)
			w.setBlock(blockPos{sx<<4 | local>>8&0xf, sy<<4 | local&0xf, sz<<4 | local>>4&0xf}, int32(entry>>12))
		}
	case packetid.ClientboundBlockEntityData:
		var pos pk.Position
		if _, err := pos.ReadFrom(bytes.NewReader(p.Data)); err != nil {
			return true, err
		}
		b := blockPos{int32(pos.X), int32(pos.Y), int32(pos.Z)}
		if c := w.chunks[chunkOf(b)]; c != nil {
			c.blockEntities[b] = p.Data
		}
	case packetid.ClientboundAddEntity:
		var spawn client.AddEntity
		if err := decode(p, &spawn); err != nil {
			return true, err
		}
		e := newEntity(p.Data)
		e.x, e.y, e.z, e.yaw, e.pitch = spawn.X, spawn.Y, spawn.Z, float32(spawn.Yaw.ToDeg()), float32(spawn.Pitch.ToDeg())
		w.entities[spawn.ID] = e
	case packetid.ClientboundRemoveEntities:
		var remove client.RemoveEntities
		if err := decode(p, &remove); err != nil {
			return true, err
		}
		for _, id := range remove.EntityIDs {
			delete(w.entities, id)
		}
	case packetid.ClientboundMoveEntityPos:
		var m client.UpdateEntityPosition
		if err := decode(p, &m); err != nil {
			return true, err
		}
		w.move(m.EntityID, m.DeltaX, m.DeltaY, m.DeltaZ, nil)
	case packetid.ClientboundMoveEntityPosRot:
		var m client.UpdateEntityPositionAndRotation
		if err := decode(p, &m); err != nil {
			return true, err
		}
		w.move(m.EntityID, m.DeltaX, m.DeltaY, m.DeltaZ, &[2]pk.Angle{m.Yaw, m.Pitch})
	case packetid.ClientboundMoveEntityRot:
		var m client.UpdateEntityRotation
		if err := decode(p, &m); err != nil {
			return true, err
		}
		w.move(m.EntityID, 0, 0, 0, &[2]pk.Angle{m.Yaw, m.Pitch})
	case packetid.ClientboundEntityPositionSync:
		var m client.TeleportEntity
		if err := decode(p, &m); err != nil {
			return true, err
		}
		if e := w.entities[m.EntityID]; e != nil && e.spawn != nil {
			e.x, e.y, e.z, e.yaw, e.pitch, e.moved = m.X, m.Y, m.Z, m.Yaw, m.Pitch, true
		}
	case packetid.ClientboundTeleportEntity:
		var m client.SynchronizeVehiclePosition
		if err := decode(p, &m); err != nil {
			return true, err
		}
		if e := w.entities[m.EntityID]; e != nil && e.spawn != nil {
			rel := func(bit int32, old, v float64) float64 {
				if m.Flags&bit != 0 {
					return old + v
				}
				return v
			}
			e.x = rel(0x01, e.x, m.X)
			e.y = rel(0x02, e.y, m.Y)
			e.z = rel(0x04, e.z, m.Z)
			e.yaw = float32(rel(0x08, float64(e.yaw), float64(m.Yaw)))
			e.pitch = float32(rel(0x10, float64(e.pitch), float64(m.Pitch)))
			e.moved = true
		}
	case packetid.ClientboundRotateHead:
		if e := w.entities[leadingVarInt(p.Data)]; e != nil && e.spawn != nil {
			e.head = p.Data
		}
	case packetid.ClientboundSetEntityData:
		r := bytes.NewReader(p.Data)
		var id pk.VarInt
		if _, err := id.ReadFrom(r); err != nil {
			return true, err
		}
		entries, _, err := metadata.ReadRaw(p.Data[len(p.Data)-r.Len():])
		if err != nil {
			return true, err
		}
		if e := w.entityOrSelf(int32(id), self); e != nil {
			maps.Copy(e.meta, entries)
		}
	case packetid.ClientboundSetEquipment:
		r := bytes.NewReader(p.Data)
		var id pk.VarInt
		if _, err := id.ReadFrom(r); err != nil {
			return true, err
		}
		items := map[int8][]byte{}
		for {
			var b pk.Byte
			if _, err := b.ReadFrom(r); err != nil {
				return true, err
			}
			raw, err := span(p.Data, r, &slot.Slot{})
			if err != nil {
				return true, err
			}
			items[int8(b)&0x7f] = raw
			if int8(b) >= 0 {
				break
			}
		}
		if e := w.entityOrSelf(int32(id), self); e != nil {
			maps.Copy(e.equipment, items)
		}
	case packetid.ClientboundUpdateAttributes:
		if e := w.entityOrSelf(leadingVarInt(p.Data), self); e != nil {
			e.attributes = p.Data
		}
	case packetid.ClientboundSetPassengers:
		if e := w.entityOrSelf(leadingVarInt(p.Data), self); e != nil {
			e.passengers = p.Data
		}
	case packetid.ClientboundSetEntityLink:
		if len(p.Data) >= 4 {
			if e := w.entityOrSelf(int32(binary.BigEndian.Uint32(p.Data)), self); e != nil {
				e.link = p.Data
			}
		}
	case packetid.ClientboundMapItemData:
		return true, w.mapData(p)
	default:
		return false, nil
	}
	return true, nil
}

func chunkOf(b blockPos) chunkPos { return chunkPos{b.x >> 4, b.z >> 4} }

func (w *world) setBlock(b blockPos, state int32) {
	c := w.chunks[chunkOf(b)]
	if c == nil {
		return
	}
	c.blocks[b] = state
	// The server sends new block entity data after a change when the new block
	// has one; the old data must not be replayed onto a different block.
	delete(c.blockEntities, b)
}

func (w *world) move(id int32, dx, dy, dz int16, rot *[2]pk.Angle) {
	e := w.entities[id]
	if e == nil || e.spawn == nil {
		return
	}
	e.x += float64(dx) / 4096
	e.y += float64(dy) / 4096
	e.z += float64(dz) / 4096
	if rot != nil {
		e.yaw, e.pitch = float32(rot[0].ToDeg()), float32(rot[1].ToDeg())
	}
	e.moved = true
}

// entityOrSelf finds a tracked entity; data for the player itself is kept in
// an entity without a spawn packet.
func (w *world) entityOrSelf(id, self int32) *entity {
	if e := w.entities[id]; e != nil {
		return e
	}
	if id != self {
		return nil
	}
	e := newEntity(nil)
	w.entities[id] = e
	return e
}

func leadingVarInt(data []byte) int32 {
	var v pk.VarInt
	if _, err := v.ReadFrom(bytes.NewReader(data)); err != nil {
		return math.MinInt32
	}
	return int32(v)
}

// mapData merges a map update (header, optional icons, optional color patch)
// into the full picture, so the snapshot sends each map once in full.
func (w *world) mapData(p pk.Packet) error {
	r := bytes.NewReader(p.Data)
	var head client.MapData
	if _, err := head.ReadFrom(r); err != nil {
		return fmt.Errorf("decode map header: %w", err)
	}
	m := w.maps[head.MapID]
	if m == nil {
		m = &mapState{}
		w.maps[head.MapID] = m
	}
	m.scale, m.locked = head.Scale, head.Locked
	var hasIcons pk.Boolean
	if _, err := hasIcons.ReadFrom(r); err != nil {
		return fmt.Errorf("decode map icons: %w", err)
	}
	if hasIcons {
		var icons []client.MapIcon
		raw, err := span(p.Data, r, pk.Array(&icons))
		if err != nil {
			return fmt.Errorf("decode map icons: %w", err)
		}
		m.icons, m.hasIcons = raw, true
	}
	var patch client.MapColorPatch
	if _, err := patch.ReadFrom(r); err != nil && err != io.EOF {
		return fmt.Errorf("decode map colors: %w", err)
	}
	if patch.Columns == 0 {
		return nil
	}
	for j := 0; j < int(patch.Rows); j++ {
		for i := 0; i < int(patch.Columns); i++ {
			x, z, k := int(patch.X)+i, int(patch.Z)+j, i+j*int(patch.Columns)
			if x < 128 && z < 128 && k < len(patch.Data) {
				m.colors[x+z*128] = byte(patch.Data[k])
			}
		}
	}
	m.hasColors = true
	return nil
}

// packets returns the world as packets in an order the client accepts: chunks
// and their changes, then entities, then passengers (which name other
// entities), then maps.
func (w *world) packets(self int32) []pk.Packet {
	var out []pk.Packet
	positions := make([]chunkPos, 0, len(w.chunks))
	for pos := range w.chunks {
		positions = append(positions, pos)
	}
	slices.SortFunc(positions, func(a, b chunkPos) int {
		if a.x != b.x {
			return int(a.x - b.x)
		}
		return int(a.z - b.z)
	})
	for _, pos := range positions {
		c := w.chunks[pos]
		out = append(out, pk.Packet{ID: int32(packetid.ClientboundLevelChunkWithLight), Data: c.data})
		out = append(out, sectionUpdates(c.blocks)...)
		for _, data := range c.blockEntities {
			out = append(out, pk.Packet{ID: int32(packetid.ClientboundBlockEntityData), Data: data})
		}
		if c.light != nil {
			out = append(out, pk.Packet{ID: int32(packetid.ClientboundLightUpdate), Data: c.light})
		}
	}

	ids := make([]int32, 0, len(w.entities))
	for id := range w.entities {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	var later []pk.Packet
	for _, id := range ids {
		e := w.entities[id]
		if e.spawn != nil {
			out = append(out, pk.Packet{ID: int32(packetid.ClientboundAddEntity), Data: e.spawn})
			if e.moved {
				out = append(out, marshal(&client.TeleportEntity{EntityID: id, X: e.x, Y: e.y, Z: e.z, Yaw: e.yaw, Pitch: e.pitch}))
			}
			if e.head != nil {
				out = append(out, pk.Packet{ID: int32(packetid.ClientboundRotateHead), Data: e.head})
			}
		} else if id != self {
			continue
		}
		if len(e.meta) > 0 {
			data := varint(id)
			for _, index := range slices.Sorted(maps.Keys(e.meta)) {
				data = append(data, e.meta[index]...)
			}
			out = append(out, pk.Packet{ID: int32(packetid.ClientboundSetEntityData), Data: append(data, 0xff)})
		}
		if len(e.equipment) > 0 {
			data := varint(id)
			slots := slices.Sorted(maps.Keys(e.equipment))
			for i, s := range slots {
				b := byte(s)
				if i < len(slots)-1 {
					b |= 0x80
				}
				data = append(data, b)
				data = append(data, e.equipment[s]...)
			}
			out = append(out, pk.Packet{ID: int32(packetid.ClientboundSetEquipment), Data: data})
		}
		if e.attributes != nil {
			out = append(out, pk.Packet{ID: int32(packetid.ClientboundUpdateAttributes), Data: e.attributes})
		}
		if e.passengers != nil {
			later = append(later, pk.Packet{ID: int32(packetid.ClientboundSetPassengers), Data: e.passengers})
		}
		if e.link != nil {
			later = append(later, pk.Packet{ID: int32(packetid.ClientboundSetEntityLink), Data: e.link})
		}
	}
	out = append(out, later...)

	mapIDs := make([]int32, 0, len(w.maps))
	for id := range w.maps {
		mapIDs = append(mapIDs, id)
	}
	slices.Sort(mapIDs)
	for _, id := range mapIDs {
		out = append(out, w.maps[id].packet(id))
	}
	return out
}

// sectionUpdates encodes changed blocks as one SectionBlocksUpdate per section.
func sectionUpdates(blocks map[blockPos]int32) []pk.Packet {
	if len(blocks) == 0 {
		return nil
	}
	type section struct{ x, y, z int32 }
	bySection := map[section][]int64{}
	var order []section
	for b, state := range blocks {
		s := section{b.x >> 4, b.y >> 4, b.z >> 4}
		if bySection[s] == nil {
			order = append(order, s)
		}
		local := int64(b.x&0xf)<<8 | int64(b.z&0xf)<<4 | int64(b.y&0xf)
		bySection[s] = append(bySection[s], int64(state)<<12|local)
	}
	slices.SortFunc(order, func(a, b section) int {
		switch {
		case a.x != b.x:
			return int(a.x - b.x)
		case a.z != b.z:
			return int(a.z - b.z)
		default:
			return int(a.y - b.y)
		}
	})
	out := make([]pk.Packet, 0, len(order))
	for _, s := range order {
		entries := bySection[s]
		slices.Sort(entries)
		pos := int64(s.x&0x3fffff)<<42 | int64(s.z&0x3fffff)<<20 | int64(s.y&0xfffff)
		out = append(out, marshal(&client.UpdateSectionsBlocks{ChunkSectionPosition: pos, Blocks: entries}))
	}
	return out
}

func (m *mapState) packet(id int32) pk.Packet {
	var buf bytes.Buffer
	_, _ = client.MapData{MapID: id, Scale: m.scale, Locked: m.locked}.WriteTo(&buf)
	_, _ = pk.Boolean(m.hasIcons).WriteTo(&buf)
	if m.hasIcons {
		buf.Write(m.icons)
	}
	patch := client.MapColorPatch{}
	if m.hasColors {
		patch = client.MapColorPatch{Columns: 128, Rows: 128, Data: make([]pk.UnsignedByte, len(m.colors))}
		for i, c := range m.colors {
			patch.Data[i] = pk.UnsignedByte(c)
		}
	}
	_, _ = patch.WriteTo(&buf)
	return pk.Packet{ID: int32(packetid.ClientboundMapItemData), Data: buf.Bytes()}
}
