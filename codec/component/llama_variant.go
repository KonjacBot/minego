package component

import "git.konjactw.dev/patyhank/minego/codec/data/slot"

//codec:gen
type LlamaVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*LlamaVariant) Type() slot.ComponentID {
	return 90
}

func (*LlamaVariant) ID() string {
	return "minecraft:llama/variant"
}
