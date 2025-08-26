package component

import (
	"io"

	"git.konjactw.dev/falloutBot/go-mc/nbt"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type CanPlaceOn struct {
	BlockPredicates []BlockPredicate
}

//codec:gen
type BlockPredicate struct {
	Blocks                         pk.Option[pk.IDSet, *pk.IDSet]
	Properties                     pk.Option[Properties, *Properties]
	NBT                            pk.Option[pk.NBTField, *pk.NBTField]
	DataComponents                 []ExactDataComponentMatcher
	PartialDataComponentPredicates []PartialDataComponentMatcher
}

type Properties []Property

func (p Properties) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Array(p).WriteTo(w)
}

func (p *Properties) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Array(p).ReadFrom(r)
}

//codec:gen
type Property struct {
	Name         string
	IsExactMatch bool
	ExactValue   pk.Option[pk.String, *pk.String]
	MinValue     pk.Option[pk.String, *pk.String]
	MaxValue     pk.Option[pk.String, *pk.String]
}

//codec:gen
type ExactDataComponentMatcher struct {
	Type  int32          `mc:"VarInt"`
	Value nbt.RawMessage `mc:"NBT"`
}

//codec:gen
type PartialDataComponentMatcher struct {
	Type      int32          `mc:"VarInt"`
	Predicate nbt.RawMessage `mc:"NBT"`
}

func (*CanPlaceOn) Type() slot.ComponentID {
	return 11
}

func (*CanPlaceOn) ID() string {
	return "minecraft:can_place_on"
}
