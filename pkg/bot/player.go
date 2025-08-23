package bot

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol"
	"github.com/go-gl/mathgl/mgl64"
)

type Player interface {
	StateID() int32
	UpdateStateID(id int32)
	Entity() Entity

	FlyTo(pos mgl64.Vec3) error
	WalkTo(pos mgl64.Vec3) error
	LookAt(vec3 mgl64.Vec3) error

	BreakBlock(pos protocol.Position) error
	PlaceBlock(pos protocol.Position) error
	PlaceBlockWithArgs(pos protocol.Position, face int32, cursor mgl64.Vec3) error
	OpenContainer(pos protocol.Position) (Container, error)

	UseItem(hand int8) error

	OpenMenu(command string) (Container, error)
}
