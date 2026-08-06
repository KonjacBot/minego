package bot

import (
	"context"

	"github.com/go-gl/mathgl/mgl64"

	"github.com/KonjacBot/minego/pkg/protocol"
)

type Player interface {
	StateID() int32
	UpdateStateID(id int32)
	Sequence() int32
	UpdateSequence(id int32)
	Entity() Entity

	FlyTo(pos mgl64.Vec3) error
	WalkTo(pos mgl64.Vec3) error
	LookAt(vec3 mgl64.Vec3) error
	UpdateLocation()

	BreakBlock(pos protocol.Position) error
	PlaceBlock(pos protocol.Position) error
	PlaceBlockWithArgs(pos protocol.Position, face int32, cursor mgl64.Vec3) error
	OpenContainer(pos protocol.Position, hand int32) (Container, error)

	UseItem(hand int8) error

	OpenMenu(command string) (Container, error)
	Command(command string) error
	Chat(message string) error
	CheckServer()
}

// BlockActionPlayer exposes context-aware, sequenced block actions. The
// returned sequence is the exact value the server will acknowledge with a
// clientbound BlockChangedAck packet.
type BlockActionPlayer interface {
	AcknowledgedSequence() int32
	WaitForSequence(ctx context.Context, id int32) error
	StartBreakingBlock(ctx context.Context, pos protocol.Position, face int8) (int32, error)
	FinishBreakingBlock(ctx context.Context, pos protocol.Position, face int8) (int32, error)
	CancelBreakingBlock(ctx context.Context, pos protocol.Position, face int8) (int32, error)
	PlaceBlockWithArgsContext(ctx context.Context, pos protocol.Position, face int32, cursor mgl64.Vec3) (int32, error)
}
