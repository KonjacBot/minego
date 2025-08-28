package world

import (
	"container/list"
	"context"
	"errors"
	"sync"

	"github.com/go-gl/mathgl/mgl64"
	"golang.org/x/exp/constraints"

	"git.konjactw.dev/falloutBot/go-mc/data/entity"
	"git.konjactw.dev/falloutBot/go-mc/level"
	"git.konjactw.dev/falloutBot/go-mc/level/block"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/bot"
	"git.konjactw.dev/patyhank/minego/pkg/protocol"
	"git.konjactw.dev/patyhank/minego/pkg/protocol/metadata"
	cp "git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/client"
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

type World struct {
	c bot.Client

	columns  map[level.ChunkPos]*level.Chunk
	entities map[int32]*Entity

	entityLock sync.Mutex
	chunkLock  sync.Mutex
}

func NewWorld(c bot.Client) *World {
	w := &World{
		c:        c,
		columns:  make(map[level.ChunkPos]*level.Chunk),
		entities: make(map[int32]*Entity),
	}

	bot.AddHandler(c, func(ctx context.Context, p *cp.LevelChunkWithLight) {
		w.chunkLock.Lock()
		defer w.chunkLock.Unlock()

		w.columns[p.Pos] = p.Data
	})

	bot.AddHandler(c, func(ctx context.Context, p *cp.ForgetLevelChunk) {
		w.chunkLock.Lock()
		defer w.chunkLock.Unlock()

		delete(w.columns, p.Pos)
	})
	bot.AddHandler(c, func(ctx context.Context, p *cp.Respawn) {
		w.chunkLock.Lock()
		defer w.chunkLock.Unlock()

		w.columns = make(map[level.ChunkPos]*level.Chunk)
	})

	bot.AddHandler(c, func(ctx context.Context, p *cp.AddEntity) {
		w.entityLock.Lock()
		defer w.entityLock.Unlock()

		w.entities[p.ID] = &Entity{
			id:         p.ID,
			entityUUID: p.UUID,
			entityType: entity.ID(p.Type),
			pos:        mgl64.Vec3{p.X, p.Y, p.Z},
			rot:        mgl64.Vec2{pk.Angle(p.XRot).ToDeg(), pk.Angle(p.YRot).ToDeg()},
			metadata:   nil,
			equipment:  nil,
		}
	})
	bot.AddHandler(c, func(ctx context.Context, p *cp.RemoveEntities) {
		w.entityLock.Lock()
		defer w.entityLock.Unlock()
		for _, d := range p.EntityIDs {
			e, ok := w.entities[d]
			if ok {
				bot.PublishEvent(c, EntityRemoveEvent{Entity: e})
				delete(w.entities, d)
			}
		}
	})
	bot.AddHandler(c, func(ctx context.Context, p *cp.SetEntityMetadata) {
		w.entityLock.Lock()
		defer w.entityLock.Unlock()
		e, ok := w.entities[p.EntityID]
		if ok {
			if e.metadata == nil {
				e.metadata = make(map[uint8]metadata.Metadata)
			}
			for u, entityMetadata := range p.Metadata.Data {
				e.metadata[u] = entityMetadata
			}
		}
	})
	bot.AddHandler(c, func(ctx context.Context, p *cp.SetEquipment) {
		w.entityLock.Lock()
		defer w.entityLock.Unlock()
		e, ok := w.entities[p.EntityID]
		if ok {
			if e.equipment == nil {
				e.equipment = make(map[int8]slot.Slot)
			}
			for _, equipment := range p.Equipment {
				e.equipment[equipment.Slot] = equipment.Item
			}
		}
	})
	bot.AddHandler(c, func(ctx context.Context, p *cp.UpdateEntityPosition) {
		w.entityLock.Lock()
		defer w.entityLock.Unlock()
		if e, ok := w.entities[p.EntityID]; ok {
			e.pos = e.pos.Add(mgl64.Vec3{float64(p.DeltaX) / 4096.0, float64(p.DeltaY) / 4096.0, float64(p.DeltaZ) / 4096.0})
		}
	})

	bot.AddHandler(c, func(ctx context.Context, p *cp.UpdateEntityRotation) {
		w.entityLock.Lock()
		defer w.entityLock.Unlock()
		if e, ok := w.entities[p.EntityID]; ok {
			e.rot = mgl64.Vec2{float64(p.Yaw), float64(p.Pitch)}
		}
	})

	bot.AddHandler(c, func(ctx context.Context, p *cp.UpdateEntityPositionAndRotation) {
		w.entityLock.Lock()
		defer w.entityLock.Unlock()
		if e, ok := w.entities[p.EntityID]; ok {
			e.pos = e.pos.Add(mgl64.Vec3{float64(p.DeltaX) / 4096.0, float64(p.DeltaY) / 4096.0, float64(p.DeltaZ) / 4096.0})
		}
	})

	bot.AddHandler(c, func(ctx context.Context, p *cp.BlockUpdate) {
		w.chunkLock.Lock()
		defer w.chunkLock.Unlock()

		pos := protocol.Position{int32(p.Position.X), int32(p.Position.Y), int32(p.Position.Z)}
		chunkX := pos[0] >> 4
		chunkZ := pos[2] >> 4
		pos2d := level.ChunkPos{chunkX, chunkZ}

		chunk, ok := w.columns[pos2d]
		if !ok {
			return // chunk not loaded, ignore update
		}

		blockX := pos[0] & 15
		blockZ := pos[2] & 15
		sectionY := pos[1] >> 4
		blockY := pos[1] & 15

		if sectionY < 0 || int(sectionY) >= len(chunk.Sections) {
			return // invalid section Y coordinate
		}

		section := chunk.Sections[sectionY]
		blockIdx := (blockY << 8) | (blockZ << 4) | blockX
		section.SetBlock(int(blockIdx), level.BlocksState(p.BlockState))
	})

	bot.AddHandler(c, func(ctx context.Context, p *cp.UpdateSectionsBlocks) {
		w.chunkLock.Lock()
		defer w.chunkLock.Unlock()

		sectionX, sectionY, sectionZ := p.ToSectionPos()
		chunkX := sectionX
		chunkZ := sectionZ
		pos2d := level.ChunkPos{chunkX, chunkZ}

		chunk, ok := w.columns[pos2d]
		if !ok {
			return // chunk not loaded, ignore update
		}

		if sectionY < 0 || int(sectionY) >= len(chunk.Sections) {
			return // invalid section Y coordinate
		}

		section := chunk.Sections[sectionY]
		blocks := p.ParseBlocks()

		for localPos, stateID := range blocks {
			blockX := localPos[0]
			blockY := localPos[1]
			blockZ := localPos[2]
			blockIdx := (blockY << 8) | (blockZ << 4) | blockX
			section.SetBlock(int(blockIdx), level.BlocksState(stateID))
		}
	})

	return w
}

func (w *World) GetBlock(pos protocol.Position) (block.Block, error) {
	w.chunkLock.Lock()
	defer w.chunkLock.Unlock()
	chunkX := pos[0] >> 4
	chunkZ := pos[2] >> 4
	pos2d := level.ChunkPos{chunkX, chunkZ}

	chunk, ok := w.columns[pos2d]
	if !ok {
		return nil, errors.New("chunk not loaded")
	}

	blockX := pos[0] & 15
	blockZ := pos[2] & 15
	blockY := pos[1] & 15
	blockIdx := (blockY << 8) | (blockZ << 4) | blockX
	sectionY := pos[1] >> 4
	if sectionY < 0 || int(sectionY) >= len(chunk.Sections) {
		return nil, errors.New("invalid section Y coordinate")
	}
	blockStateId := chunk.Sections[sectionY].GetBlock(int(blockIdx))
	return block.StateList[blockStateId], nil
}

func (w *World) SetBlock(pos protocol.Position, blk block.Block) error {
	w.chunkLock.Lock()
	defer w.chunkLock.Unlock()

	chunkX := pos[0] >> 4
	chunkZ := pos[2] >> 4
	pos2d := level.ChunkPos{chunkX, chunkZ}

	chunk, ok := w.columns[pos2d]
	if !ok {
		return errors.New("chunk not loaded")
	}

	blockX := pos[0] & 15
	blockZ := pos[2] & 15
	sectionY := pos[1] >> 4
	blockY := pos[1] & 15

	if sectionY < 0 || int(sectionY) >= len(chunk.Sections) {
		return errors.New("invalid section Y coordinate")
	}

	section := chunk.Sections[sectionY]

	blockIdx := (blockY << 8) | (blockZ << 4) | blockX
	section.SetBlock(int(blockIdx), block.ToStateID[blk])
	return nil
}

func (w *World) GetNearbyBlocks(pos protocol.Position, radius int32) ([]block.Block, error) {
	var blocks []block.Block

	for dx := -radius; dx <= radius; dx++ {
		for dy := -radius; dy <= radius; dy++ {
			for dz := -radius; dz <= radius; dz++ {
				blk, err := w.GetBlock(protocol.Position{pos[0] + dx, pos[1] + dy, pos[2] + dz})
				if err != nil {
					continue
				}
				blocks = append(blocks, blk)
			}
		}
	}

	return blocks, nil
}

func (w *World) FindNearbyBlock(pos protocol.Position, radius int32, blk block.Block) (protocol.Position, error) {
	visited := make(map[protocol.Position]bool)
	queue := list.New()
	start := pos
	queue.PushBack(start)
	visited[start] = true

	// Direction vectors for 6-way adjacent blocks
	dirs := []protocol.Position{
		{1, 0, 0}, {-1, 0, 0},
		{0, 1, 0}, {0, -1, 0},
		{0, 0, 1}, {0, 0, -1},
	}
	for queue.Len() > 0 {
		current := queue.Remove(queue.Front()).(protocol.Position)

		// Skip if beyond the radius
		if abs(current[0]-pos[0]) > radius || abs(current[1]-pos[1]) > radius || abs(current[2]-pos[2]) > radius {
			continue
		}

		// Check if current block matches target
		if currentBlock, err := w.GetBlock(current); err == nil {
			if currentBlock == blk {
				return current, nil
			}
		}

		// Check all 6 adjacent blocks
		for _, dir := range dirs {
			next := protocol.Position{
				current[0] + dir[0],
				current[1] + dir[1],
				current[2] + dir[2],
			}

			if !visited[next] {
				visited[next] = true
				queue.PushBack(next)
			}
		}
	}

	return protocol.Position{}, errors.New("block not found")
}

func (w *World) Entities() []bot.Entity {
	w.entityLock.Lock()
	defer w.entityLock.Unlock()
	var entities []bot.Entity
	for _, e := range w.entities {
		entities = append(entities, e)
	}
	return entities
}

func (w *World) GetEntity(id int32) bot.Entity {
	w.entityLock.Lock()
	defer w.entityLock.Unlock()
	return w.entities[id]
}

func (w *World) GetNearbyEntities(radius int32) []bot.Entity {
	w.entityLock.Lock()
	defer w.entityLock.Unlock()

	selfPos := w.c.Player().Entity().Position()
	var entities []bot.Entity

	for _, e := range w.entities {
		sqr := e.pos.Sub(selfPos).LenSqr()
		if sqr <= float64(radius*radius) {
			entities = append(entities, e)
		}
	}
	return entities
}

func (w *World) GetEntitiesByType(entityType entity.ID) []bot.Entity {
	w.entityLock.Lock()
	defer w.entityLock.Unlock()

	var entities []bot.Entity
	for _, e := range w.entities {
		if e.entityType == entityType {
			entities = append(entities, e)
		}
	}
	return entities
}

func abs[T constraints.Signed | constraints.Float](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
