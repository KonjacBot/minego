package component

import (
	"git.konjactw.dev/patyhank/minego/codec/slot"
	pk "github.com/Tnze/go-mc/net/packet"
)

//codec:gen
type Profile struct {
	HasName     bool
	Name        pk.Option[pk.String, *pk.String]
	HasUniqueID bool
	UniqueID    pk.Option[pk.UUID, *pk.UUID]
	Properties  []ProfileProperty
}

//codec:gen
type ProfileProperty struct {
	Name         string `mc:"String"`
	Value        string `mc:"String"`
	HasSignature bool
	Signature    pk.Option[pk.String, *pk.String]
}

func (*Profile) Type() slot.ComponentID {
	return 61
}

func (*Profile) ID() string {
	return "minecraft:profile"
}
