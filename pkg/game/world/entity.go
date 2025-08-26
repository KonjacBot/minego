package world

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"

	"git.konjactw.dev/falloutBot/go-mc/data/entity"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/metadata"
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

type Entity struct {
	id         int32
	entityUUID uuid.UUID
	entityType entity.ID
	pos        mgl64.Vec3
	rot        mgl64.Vec2
	metadata   map[uint8]metadata.Metadata
	equipment  map[int8]slot.Slot
}

func (e *Entity) ID() int32 {
	return e.id
}

func (e *Entity) UUID() uuid.UUID {
	return e.entityUUID
}

func (e *Entity) Type() entity.ID {
	return e.entityType
}

func (e *Entity) Position() mgl64.Vec3 {
	return e.pos
}

func (e *Entity) Rotation() mgl64.Vec2 {
	return e.rot
}

func (e *Entity) Metadata() map[uint8]metadata.Metadata {
	return e.metadata
}

func (e *Entity) Equipment() map[int8]slot.Slot {
	return e.equipment
}

func (e *Entity) SetPosition(pos mgl64.Vec3) {
	e.pos = pos
}

func (e *Entity) SetRotation(rot mgl64.Vec2) {
	e.rot = rot
}

func (e *Entity) SetID(id int32) {
	e.id = id
}
