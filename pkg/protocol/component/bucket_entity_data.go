package component

import (
	"git.konjactw.dev/falloutBot/go-mc/nbt"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
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
