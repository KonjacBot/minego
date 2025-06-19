package component

import (
	"git.konjactw.dev/patyhank/minego/codec/data/slot"
	"github.com/Tnze/go-mc/nbt"
)

//codec:gen
type BucketEntityData struct {
	Data nbt.RawMessage `mc:"NBT"`
}

func (*BucketEntityData) Type() slot.ComponentID {
	return 50
}

func (*BucketEntityData) ID() string {
	return "minecraft:bucket_entity_data"
}
