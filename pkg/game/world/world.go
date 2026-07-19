package world

import (
	"container/list"
	"context"
	"errors"
	"sync"

	"github.com/go-gl/mathgl/mgl64"
	"golang.org/x/exp/constraints"

	"github.com/KonjacBot/go-mc/data/entity"
	"github.com/KonjacBot/go-mc/level"
	"github.com/KonjacBot/go-mc/level/block"

	"github.com/KonjacBot/minego/pkg/bot"
	"github.com/KonjacBot/minego/pkg/protocol"
	cp "github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

type World struct {
	c bot.Client

	chunkHeight uint32
	columns     map[level.ChunkPos]*level.Chunk
	entities    map[int32]*Entity

	entityLock sync.RWMutex
	chunkLock  sync.RWMutex
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

	bot.AddHandler(c, func(ctx context.Context, p *cp.Login) {
		w.chunkLock.Lock()
		switch p.CommonPlayerSpawnInfo.DimensionType {
		case 0:
			w.chunkHeight = 384
		case 1, 2:
			fallthrough
		default:
			w.chunkHeight = 256
		}
		w.columns = make(map[level.ChunkPos]*level.Chunk)
		w.chunkLock.Unlock()

		w.entityLock.Lock()
		w.entities = make(map[int32]*Entity)
		w.entityLock.Unlock()
	})

	bot.AddHandler(c, func(ctx context.Context, p *cp.ForgetLevelChunk) {
		w.chunkLock.Lock()
		defer w.chunkLock.Unlock()

		delete(w.columns, p.Pos)
	})

	bot.AddHandler(c, func(ctx context.Context, p *cp.Respawn) {
		w.chunkLock.Lock()
		switch p.CommonPlayerSpawnInfo.DimensionType {
		case 0:
			w.chunkHeight = 384
		case 1, 2:
			fallthrough
		default:
			w.chunkHeight = 256
		}

		w.columns = make(map[level.ChunkPos]*level.Chunk)
		w.chunkLock.Unlock()

		w.entityLock.Lock()
		w.entities = make(map[int32]*Entity)
		w.entityLock.Unlock()
	})

	bot.AddHandler(c, func(ctx context.Context, p *cp.AddEntity) {
		entity := &Entity{
			id:         p.ID,
			entityUUID: p.UUID,
			entityType: entity.ID(p.Type),
			pos:        mgl64.Vec3{p.X, p.Y, p.Z},
			rot:        mgl64.Vec2{p.Yaw.ToDeg(), p.Pitch.ToDeg()},
			metadata:   nil,
			equipment:  nil,
		}
		w.entityLock.Lock()
		w.entities[p.ID] = entity
		w.entityLock.Unlock()
		_ = bot.PublishEvent(c, EntityAddEvent{EntityID: p.ID})
	})
	bot.AddHandler(c, func(ctx context.Context, p *cp.RemoveEntities) {
		var removed []*Entity
		w.entityLock.Lock()
		for _, d := range p.EntityIDs {
			e, ok := w.entities[d]
			if ok {
				delete(w.entities, d)
				removed = append(removed, e)
			}
		}
		w.entityLock.Unlock()
		for _, e := range removed {
			_ = bot.PublishEvent(c, EntityRemoveEvent{Entity: e})
		}
	})
	bot.AddHandler(c, func(ctx context.Context, p *cp.SetEntityMetadata) {
		w.entityLock.RLock()
		e, ok := w.entities[p.EntityID]
		w.entityLock.RUnlock()
		if ok {
			e.updateMetadata(p.Metadata.Data)
		}
	})
	bot.AddHandler(c, func(ctx context.Context, p *cp.SetEquipment) {
		w.entityLock.RLock()
		e, ok := w.entities[p.EntityID]
		w.entityLock.RUnlock()
		if ok {
			e.mu.Lock()
			if e.equipment == nil {
				e.equipment = make(map[int8]slot.Slot)
			}
			for _, equipment := range p.Equipment {
				e.equipment[equipment.Slot] = equipment.Item.Clone()
			}
			e.mu.Unlock()
		}
	})
	bot.AddHandler(c, func(ctx context.Context, p *cp.UpdateEntityPosition) {
		w.entityLock.RLock()
		defer w.entityLock.RUnlock()
		if e, ok := w.entities[p.EntityID]; ok {
			e.addPosition(mgl64.Vec3{float64(p.DeltaX) / 4096.0, float64(p.DeltaY) / 4096.0, float64(p.DeltaZ) / 4096.0})
		}
	})

	bot.AddHandler(c, func(ctx context.Context, p *cp.UpdateEntityRotation) {
		w.entityLock.RLock()
		e, ok := w.entities[p.EntityID]
		w.entityLock.RUnlock()
		if ok {
			e.SetRotation(mgl64.Vec2{p.Yaw.ToDeg(), p.Pitch.ToDeg()})
		}
	})

	bot.AddHandler(c, func(ctx context.Context, p *cp.UpdateEntityPositionAndRotation) {
		w.entityLock.RLock()
		e, ok := w.entities[p.EntityID]
		w.entityLock.RUnlock()
		if ok {
			e.addPosition(mgl64.Vec3{float64(p.DeltaX) / 4096.0, float64(p.DeltaY) / 4096.0, float64(p.DeltaZ) / 4096.0})
			e.SetRotation(mgl64.Vec2{p.Yaw.ToDeg(), p.Pitch.ToDeg()})
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

		if len(chunk.Sections) > 16 {
			sectionY += 4
		}

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

		if len(chunk.Sections) > 16 {
			sectionY += 4
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
	w.chunkLock.RLock()
	defer w.chunkLock.RUnlock()
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
	if len(chunk.Sections) > 16 {
		sectionY += 4
	}
	if sectionY < 0 || int(sectionY) >= len(chunk.Sections) {
		return nil, errors.New("invalid section Y coordinate")
	}
	blockStateId := chunk.Sections[sectionY].GetBlock(int(blockIdx))
	if int(blockStateId) < 0 || int(blockStateId) >= len(block.StateList) {
		return nil, errors.New("unknown block state")
	}
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
	if len(chunk.Sections) > 16 {
		sectionY += 4
	}
	if sectionY < 0 || int(sectionY) >= len(chunk.Sections) {
		return errors.New("invalid section Y coordinate")
	}

	section := chunk.Sections[sectionY]

	blockIdx := (blockY << 8) | (blockZ << 4) | blockX
	stateID, ok := block.ToStateID[blk]
	if !ok {
		return errors.New("unknown block")
	}
	section.SetBlock(int(blockIdx), stateID)
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
			if currentBlock.ID() == blk.ID() {
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
	w.entityLock.RLock()
	defer w.entityLock.RUnlock()
	var entities []bot.Entity
	for _, e := range w.entities {
		entities = append(entities, e)
	}
	return entities
}

func (w *World) GetEntity(id int32) bot.Entity {
	w.entityLock.RLock()
	defer w.entityLock.RUnlock()
	entity, ok := w.entities[id]
	if !ok {
		return nil
	}
	return entity
}

func (w *World) GetNearbyEntities(radius int32) []bot.Entity {
	w.entityLock.RLock()
	defer w.entityLock.RUnlock()

	selfPos := w.c.Player().Entity().Position()
	var entities []bot.Entity

	for _, e := range w.entities {
		sqr := e.Position().Sub(selfPos).LenSqr()
		if sqr <= float64(radius)*float64(radius) {
			entities = append(entities, e)
		}
	}
	return entities
}

func (w *World) GetEntitiesByType(entityType entity.ID) []bot.Entity {
	w.entityLock.RLock()
	defer w.entityLock.RUnlock()

	var entities []bot.Entity
	for _, e := range w.entities {
		if e.Type() == entityType {
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
