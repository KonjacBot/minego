package bot

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol"
	"git.konjactw.dev/patyhank/minego/pkg/protocol/metadata"
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
	"github.com/Tnze/go-mc/data/entity"
	"github.com/Tnze/go-mc/level/block"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
)

type World interface {
	GetBlock(pos protocol.Position) (block.Block, error)
	SetBlock(pos protocol.Position, b block.Block) error

	GetNearbyBlocks(pos protocol.Position, radius int32) ([]block.Block, error)
	FindNearbyBlock(pos protocol.Position, radius int32, blk block.Block) (protocol.Position, error)

	Entities() []Entity
	GetEntity(id int32) Entity
	GetNearbyEntities(radius int32) []Entity
	GetEntitiesByType(entityType entity.ID) []Entity
}

type Entity interface {
	ID() int32
	UUID() uuid.UUID
	Type() entity.ID
	Position() mgl64.Vec3
	Rotation() mgl64.Vec2

	Metadata() map[uint8]metadata.Metadata
	Equipment() map[int8]slot.Slot
}
