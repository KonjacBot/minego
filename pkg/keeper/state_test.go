package keeper

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/KonjacBot/go-mc/chat"
	"github.com/KonjacBot/go-mc/data/packetid"
	"github.com/KonjacBot/go-mc/level"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/v2/pkg/protocol/component"
	"github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/server"
	"github.com/KonjacBot/minego/v2/pkg/protocol/slot"
)

func loggedIn(t *testing.T) *state {
	t.Helper()
	s := newState()
	login := client.Login{PlayerID: 7, CommonPlayerSpawnInfo: client.CommonPlayerSpawnInfo{Dimension: "minecraft:overworld"}}
	if err := s.fromServer(marshal(&login)); err != nil {
		t.Fatal(err)
	}
	return s
}

func chunkPacket(x, z int32) pk.Packet {
	return marshal(&client.LevelChunkWithLight{Pos: level.ChunkPos{x, z}, Data: level.EmptyChunk(24)})
}

// The snapshot carries a loaded chunk with every block change since, and
// forgets chunks the server unloaded.
func TestWorldKeepsBlockChanges(t *testing.T) {
	s := loggedIn(t)
	for _, p := range []pk.Packet{
		chunkPacket(0, 0),
		chunkPacket(-1, 0),
		marshal(&client.BlockUpdate{Position: pk.Position{X: 1, Y: 64, Z: 2}, BlockState: 10}),
		marshal(&client.BlockUpdate{Position: pk.Position{X: 1, Y: 64, Z: 2}, BlockState: 11}),
		sectionUpdate(-1, 3, 0, map[[3]int]int32{{15, 1, 0}: 20, {0, 2, 15}: 21}),
		marshal(&client.BlockUpdate{Position: pk.Position{X: 100, Y: 64, Z: 100}, BlockState: 99}), // not loaded
	} {
		if err := s.fromServer(p); err != nil {
			t.Fatal(err)
		}
	}
	got := changedBlocks(t, s.world.packets(s.playerID))
	want := map[[3]int32]int32{{1, 64, 2}: 11, {-1, 49, 0}: 20, {-16, 50, 15}: 21}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("changed blocks = %v, want %v", got, want)
	}

	var forget pk.Packet
	forget.ID = int32(packetid.ClientboundForgetLevelChunk)
	forget.Data = pk.Marshal(0, pk.Long(int64(0)<<32|int64(uint32(-1&0xffffffff)))).Data
	if err := s.fromServer(forget); err != nil {
		t.Fatal(err)
	}
	if _, ok := s.world.chunks[chunkPos{-1, 0}]; ok {
		t.Fatal("forgotten chunk still in the snapshot")
	}
	got = changedBlocks(t, s.world.packets(s.playerID))
	if !reflect.DeepEqual(got, map[[3]int32]int32{{1, 64, 2}: 11}) {
		t.Fatalf("after forget, changed blocks = %v", got)
	}
}

func sectionUpdate(sx, sy, sz int32, blocks map[[3]int]int32) pk.Packet {
	u := client.UpdateSectionsBlocks{ChunkSectionPosition: int64(sx&0x3fffff)<<42 | int64(sz&0x3fffff)<<20 | int64(sy&0xfffff)}
	for local, state := range blocks {
		u.AddBlock(local[0], local[1], local[2], state)
	}
	return marshal(&u)
}

func changedBlocks(t *testing.T, packets []pk.Packet) map[[3]int32]int32 {
	t.Helper()
	out := map[[3]int32]int32{}
	for _, p := range packets {
		if packetid.ClientboundPacketID(p.ID) != packetid.ClientboundSectionBlocksUpdate {
			continue
		}
		var u client.UpdateSectionsBlocks
		if err := decode(p, &u); err != nil {
			t.Fatal(err)
		}
		sx, sy, sz := u.ToSectionPos()
		for local, state := range u.ParseBlocks() {
			out[[3]int32{sx<<4 + local[0], sy<<4 + local[1], sz<<4 + local[2]}] = state
		}
	}
	return out
}

// A click the server confirms silently is applied from the core's hashed
// prediction, recovering the full stack (with its components) it moved.
func TestInventoryResolvesHashedClicks(t *testing.T) {
	nameID, ok := slot.ComponentID("minecraft:custom_name")
	if !ok {
		t.Fatal("custom_name component not registered")
	}
	named := slot.Slot{ItemID: 900, Count: 5, AddComponent: map[int32]slot.Component{int32(nameID): &component.CustomName{Name: chat.Text("編號1")}}}
	plain := slot.Slot{ItemID: 900, Count: 5}
	hash := func(s slot.Slot) slot.HashedSlot {
		h, err := slot.HashSlot(s)
		if err != nil {
			t.Fatal(err)
		}
		return h
	}
	split := named
	split.Count = 2

	tests := []struct {
		name    string
		click   server.ContainerClick
		slot36  slot.Slot
		slot9   slot.Slot
		cursor  slot.Slot
		missing int
	}{
		{
			name: "move the named stack to the main inventory",
			click: server.ContainerClick{ChangedSlots: []server.ChangedSlot{
				{Slot: 36, SlotData: slot.HashedSlot{}}, {Slot: 9, SlotData: hash(named)},
			}},
			slot9: named,
		},
		{
			name:   "pick up part of it keeps the name on both halves",
			click:  server.ContainerClick{ChangedSlots: []server.ChangedSlot{{Slot: 36, SlotData: hash(func() slot.Slot { s := named; s.Count = 3; return s }())}}, CarriedSlot: hash(split)},
			slot36: func() slot.Slot { s := named; s.Count = 3; return s }(),
			cursor: split,
		},
		{
			name: "an item the keeper never saw keeps id and count",
			click: server.ContainerClick{ChangedSlots: []server.ChangedSlot{
				{Slot: 9, SlotData: hash(slot.Slot{ItemID: 901, Count: 1, AddComponent: named.AddComponent})},
			}},
			slot36:  named,
			slot9:   slot.Slot{ItemID: 901, Count: 1},
			missing: 1,
		},
		{
			name:   "a click on another window is ignored",
			click:  server.ContainerClick{WindowID: 3, ChangedSlots: []server.ChangedSlot{{Slot: 36, SlotData: slot.HashedSlot{}}}},
			slot36: named,
			cursor: plain,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := loggedIn(t)
			content := client.SetContainerContent{Slots: make([]slot.Slot, playerSlots), CarriedItem: tt.cursor}
			content.Slots[36] = named
			if tt.name == "a click on another window is ignored" {
				content.CarriedItem = plain
			}
			if err := s.fromServer(marshal(&content)); err != nil {
				t.Fatal(err)
			}
			if err := s.fromCore(pk.Marshal(packetid.ServerboundContainerClick, &tt.click)); err != nil {
				t.Fatal(err)
			}
			if !sameStack(s.inv.player[36], tt.slot36) || !sameStack(s.inv.player[9], tt.slot9) || !sameStack(s.inv.cursor, tt.cursor) {
				t.Fatalf("slot36=%+v slot9=%+v cursor=%+v", s.inv.player[36], s.inv.player[9], s.inv.cursor)
			}
			if s.inv.unresolved != tt.missing {
				t.Fatalf("unresolved = %d, want %d", s.inv.unresolved, tt.missing)
			}
		})
	}
}

func sameStack(raw []byte, b slot.Slot) bool {
	var a slot.Slot
	if _, err := a.ReadFrom(bytes.NewReader(raw)); err != nil {
		return false
	}
	if a.Count <= 0 && b.Count <= 0 {
		return true
	}
	ha, errA := slot.HashSlot(a)
	hb, errB := slot.HashSlot(b)
	return errA == nil && errB == nil && a.Count == b.Count && sameItem(ha, hb)
}

// Entities come back where they moved to, with merged metadata; a new
// dimension forgets them, respawning in the same one does not.
func TestEntitiesAndDimensions(t *testing.T) {
	s := loggedIn(t)
	for _, p := range []pk.Packet{
		chunkPacket(0, 0),
		marshal(&client.AddEntity{ID: 3, Type: 5, X: 1, Y: 64, Z: 1}),
		marshal(&client.UpdateEntityPosition{EntityID: 3, DeltaX: 4096 * 2, DeltaZ: -4096}),
		marshal(&client.TeleportEntity{EntityID: 3, X: 10, Y: 65, Z: 10, Yaw: 90}),
		marshal(&client.UpdateEntityPosition{EntityID: 3, DeltaY: 2048}),
	} {
		if err := s.fromServer(p); err != nil {
			t.Fatal(err)
		}
	}
	e := s.world.entities[3]
	if e == nil || e.x != 10 || e.y != 65.5 || e.z != 10 || e.yaw != 90 {
		t.Fatalf("entity = %+v", e)
	}

	respawn := func(dim string) pk.Packet {
		return marshal(&client.Respawn{CommonPlayerSpawnInfo: client.CommonPlayerSpawnInfo{Dimension: dim}})
	}
	if err := s.fromServer(respawn("minecraft:overworld")); err != nil {
		t.Fatal(err)
	}
	if len(s.world.entities) != 1 || len(s.world.chunks) != 1 {
		t.Fatal("respawning in the same dimension forgot the world")
	}
	if err := s.fromServer(respawn("minecraft:the_nether")); err != nil {
		t.Fatal(err)
	}
	if len(s.world.entities) != 0 || len(s.world.chunks) != 0 {
		t.Fatal("a new dimension kept the old world")
	}
	packets, _, ok := s.snapshot()
	if !ok {
		t.Fatal("no snapshot after login")
	}
	// The server's own Login then its latest Respawn: the core ends up in the nether.
	var last client.Respawn
	if len(packets) < 2 || packetid.ClientboundPacketID(packets[1].ID) != packetid.ClientboundRespawn {
		t.Fatal("snapshot does not replay the respawn after the login")
	}
	if err := decode(packets[1], &last); err != nil {
		t.Fatal(err)
	}
	if last.CommonPlayerSpawnInfo.Dimension != "minecraft:the_nether" {
		t.Fatalf("snapshot enters %q, want the nether", last.CommonPlayerSpawnInfo.Dimension)
	}
}

// The core's own position and flying state, which the server never echoes,
// come back in the snapshot; the synthetic teleport ids never repeat.
func TestSnapshotRestoresPredictedState(t *testing.T) {
	s := loggedIn(t)
	if err := s.fromServer(marshal(&client.PlayerAbilities{Flags: 0x04})); err != nil {
		t.Fatal(err)
	}
	for _, p := range []server.ServerboundPacket{
		&server.MovePlayerPosRot{X: 1.5, FeetY: 70, Z: -2.5, XRot: 45, YRot: 10},
		&server.PlayerAbilities{Flags: 0x02},
		&server.SetCarriedItem{Slot: 4},
	} {
		if err := s.fromCore(marshalServer(p)); err != nil {
			t.Fatal(err)
		}
	}
	first, teleport, _ := s.snapshot()
	_, second, _ := s.snapshot()
	if teleport == second || teleport < syntheticTeleportBase {
		t.Fatalf("teleport ids %d then %d", teleport, second)
	}
	var sawPos, sawFlying, sawHeld bool
	for _, p := range first {
		switch packetid.ClientboundPacketID(p.ID) {
		case packetid.ClientboundPlayerPosition:
			var pos client.PlayerPosition
			_ = decode(p, &pos)
			sawPos = pos.X == 1.5 && pos.Y == 70 && pos.Z == -2.5 && pos.XRot == 45 && pos.YRot == 10
		case packetid.ClientboundPlayerAbilities:
			var a client.PlayerAbilities
			_ = decode(p, &a)
			sawFlying = a.Flags == 0x06
		case packetid.ClientboundSetHeldSlot:
			var h client.SetHeldItem
			_ = decode(p, &h)
			sawHeld = h.Slot == 4
		}
	}
	if !sawPos || !sawFlying || !sawHeld {
		t.Fatalf("position %v flying %v held %v", sawPos, sawFlying, sawHeld)
	}
}

// Every packet of a snapshot must decode with minego's own decoder: the core
// drops the connection on the first packet it cannot read.
func TestSnapshotDecodes(t *testing.T) {
	s := loggedIn(t)
	nameID, _ := slot.ComponentID("minecraft:custom_name")
	named := slot.Slot{ItemID: 900, Count: 5, AddComponent: map[int32]slot.Component{int32(nameID): &component.CustomName{Name: chat.Text("編號1")}}}
	content := client.SetContainerContent{Slots: make([]slot.Slot, playerSlots)}
	content.Slots[36] = named
	for _, p := range []pk.Packet{
		chunkPacket(0, 0),
		marshal(&client.BlockUpdate{Position: pk.Position{X: 1, Y: 64, Z: 2}, BlockState: 10}),
		marshal(&client.AddEntity{ID: 3, Type: 5, X: 1, Y: 64, Z: 1}),
		marshal(&client.UpdateEntityPosition{EntityID: 3, DeltaX: 4096}),
		marshal(&client.SetEquipment{EntityID: 3, Equipment: client.Equipment{{Slot: 0, Item: named}, {Slot: 5, Item: named}}}),
		marshal(&content),
		marshal(&client.PlayerAbilities{Flags: 0x04}),
		marshal(&client.PlayerInfoUpdate{Players: map[uuid.UUID][]client.PlayerInfo{
			uuid.New(): {&client.PlayerInfoAddPlayer{Name: "a"}, &client.PlayerInfoUpdateLatency{Ping: 5}},
		}}),
		marshal(&client.PlayerInfoUpdate{Players: map[uuid.UUID][]client.PlayerInfo{
			uuid.New(): {&client.PlayerInfoUpdateLatency{Ping: 9}},
		}}),
		marshal(&client.Respawn{CommonPlayerSpawnInfo: client.CommonPlayerSpawnInfo{Dimension: "minecraft:overworld"}}),
	} {
		if err := s.fromServer(p); err != nil {
			t.Fatalf("record %#x: %v", p.ID, err)
		}
	}
	if err := s.fromCore(marshalServer(&server.MovePlayerPos{X: 1, FeetY: 64, Z: 1})); err != nil {
		t.Fatal(err)
	}
	packets, _, _ := s.snapshot()
	for _, p := range packets {
		creator := client.ClientboundPackets[packetid.ClientboundPacketID(p.ID)]
		if creator == nil {
			continue
		}
		r := bytes.NewReader(p.Data)
		if _, err := creator().ReadFrom(r); err != nil {
			t.Errorf("snapshot packet %v does not decode: %v", packetid.ClientboundPacketID(p.ID), err)
		}
	}
}
