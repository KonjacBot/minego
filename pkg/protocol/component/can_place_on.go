package component

import (
	"io"

	"github.com/KonjacBot/go-mc/nbt"
	pk "github.com/KonjacBot/go-mc/net/packet"
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type CanPlaceOn struct {
	BlockPredicates []BlockPredicate
}

//codec:gen
type BlockPredicate struct {
	Blocks                         pk.Option[wire.IDSet, *wire.IDSet]
	Properties                     pk.Option[Properties, *Properties]
	NBT                            pk.Option[wire.NBTField, *wire.NBTField]
	DataComponents                 []ExactDataComponentMatcher
	PartialDataComponentPredicates []PartialDataComponentMatcher
}

type Properties []Property

func (p Properties) WriteTo(w io.Writer) (n int64, err error) {
	return wire.Array(p).WriteTo(w)
}

func (p *Properties) ReadFrom(r io.Reader) (n int64, err error) {
	return wire.Array(p).ReadFrom(r)
}

//codec:gen
type Property struct {
	Name         string
	IsExactMatch bool
	ExactValue   pk.Option[wire.String, *wire.String]
	MinValue     pk.Option[wire.String, *wire.String]
	MaxValue     pk.Option[wire.String, *wire.String]
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

func (*CanPlaceOn) ID() string {
	return "minecraft:can_place_on"
}
