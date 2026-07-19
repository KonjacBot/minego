package world

import (
	"sync"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"

	"github.com/KonjacBot/go-mc/data/entity"

	"github.com/KonjacBot/minego/pkg/protocol/metadata"
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

type Entity struct {
	mu         sync.RWMutex
	id         int32
	entityUUID uuid.UUID
	entityType entity.ID
	pos        mgl64.Vec3
	rot        mgl64.Vec2
	metadata   map[uint8]metadata.Metadata
	equipment  map[int8]slot.Slot
}

func (e *Entity) ID() int32 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.id
}

func (e *Entity) UUID() uuid.UUID {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.entityUUID
}

func (e *Entity) Type() entity.ID {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.entityType
}

func (e *Entity) Position() mgl64.Vec3 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.pos
}

func (e *Entity) Rotation() mgl64.Vec2 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.rot
}

func (e *Entity) Metadata() map[uint8]metadata.Metadata {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.metadata == nil {
		return nil
	}
	result := make(map[uint8]metadata.Metadata, len(e.metadata))
	for index, value := range e.metadata {
		result[index] = value
	}
	return result
}

func (e *Entity) Equipment() map[int8]slot.Slot {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.equipment == nil {
		return nil
	}
	result := make(map[int8]slot.Slot, len(e.equipment))
	for index, value := range e.equipment {
		result[index] = value.Clone()
	}
	return result
}

func (e *Entity) SetPosition(pos mgl64.Vec3) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.pos = pos
}

func (e *Entity) SetRotation(rot mgl64.Vec2) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rot = rot
}

func (e *Entity) SetID(id int32) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.id = id
}

func (e *Entity) addPosition(delta mgl64.Vec3) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.pos = e.pos.Add(delta)
}

func (e *Entity) updateMetadata(values map[uint8]metadata.Metadata) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.metadata == nil {
		e.metadata = make(map[uint8]metadata.Metadata)
	}
	for index, value := range values {
		e.metadata[index] = value
	}
}
