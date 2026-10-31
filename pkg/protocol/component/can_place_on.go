package component

import (
	"errors"
	"fmt"
	"io"

	"github.com/KonjacBot/go-mc/nbt"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type CanPlaceOn struct {
	BlockPredicates []BlockPredicate
}

// BlockPredicate 是冒險模式能放、能挖哪些方塊的條件(26.2 的 BlockPredicate.STREAM_CODEC)。
//
//codec:gen
type BlockPredicate struct {
	Blocks                         pk.Option[pk.IDSet, *pk.IDSet]
	Properties                     pk.Option[Properties, *Properties]
	NBT                            pk.Option[NBTValue, *NBTValue]
	DataComponents                 []ExactDataComponentMatcher
	PartialDataComponentPredicates []PartialDataComponentMatcher
}

// NBTValue 是原樣保留的一段網路 NBT(根標籤沒有名字)。pk.NBTField 沒給目標時解不出來,這裡只留原始位元組。
type NBTValue struct {
	Data nbt.RawMessage
}

func (v NBTValue) WriteTo(w io.Writer) (int64, error) {
	return pk.NBT(&v.Data).WriteTo(w)
}

func (v *NBTValue) ReadFrom(r io.Reader) (int64, error) {
	return pk.NBT(&v.Data).ReadFrom(r)
}

type Properties []Property

func (p Properties) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Array(p).WriteTo(w)
}

func (p *Properties) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Array(p).ReadFrom(r)
}

// Property 是一個方塊狀態條件:名稱後面是 Either,true 接一個要相等的值(不是 optional),
// false 接 optional 的最小值與最大值。
type Property struct {
	Name         string
	IsExactMatch bool
	ExactValue   string
	MinValue     pk.Option[pk.String, *pk.String]
	MaxValue     pk.Option[pk.String, *pk.String]
}

func (p *Property) ReadFrom(r io.Reader) (n int64, err error) {
	fields := []pk.FieldDecoder{(*pk.String)(&p.Name), (*pk.Boolean)(&p.IsExactMatch)}
	if n, err = readAll(r, fields); err != nil {
		return n, err
	}
	if p.IsExactMatch {
		fields = []pk.FieldDecoder{(*pk.String)(&p.ExactValue)}
	} else {
		fields = []pk.FieldDecoder{&p.MinValue, &p.MaxValue}
	}
	temp, err := readAll(r, fields)
	return n + temp, err
}

// readAll 依序讀完 fields。
func readAll(r io.Reader, fields []pk.FieldDecoder) (n int64, err error) {
	for _, f := range fields {
		temp, err := f.ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (p Property) WriteTo(w io.Writer) (n int64, err error) {
	fields := []pk.FieldEncoder{pk.String(p.Name), pk.Boolean(p.IsExactMatch)}
	if p.IsExactMatch {
		fields = append(fields, pk.String(p.ExactValue))
	} else {
		fields = append(fields, p.MinValue, p.MaxValue)
	}
	for _, f := range fields {
		temp, err := f.WriteTo(w)
		n += temp
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

// ExactDataComponentMatcher 是「元件要完全一樣」的條件:元件種類後面接的是該元件自己的網路格式,不是 NBT。
type ExactDataComponentMatcher struct {
	Type  int32
	Value slot.Component
}

// errUnknownComponent 是條件裡的元件種類 minego 不認得:後面的長度不知道,沒辦法跳過。
var errUnknownComponent = errors.New("unknown data component")

func (m *ExactDataComponentMatcher) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = (*pk.VarInt)(&m.Type).ReadFrom(r)
	if err != nil {
		return n, err
	}
	m.Value = slot.ComponentFromID(int(m.Type))
	if m.Value == nil {
		return n, fmt.Errorf("exact component predicate %d: %w", m.Type, errUnknownComponent)
	}
	temp, err := m.Value.ReadFrom(r)
	return n + temp, err
}

func (m ExactDataComponentMatcher) WriteTo(w io.Writer) (n int64, err error) {
	n, err = pk.VarInt(m.Type).WriteTo(w)
	if err != nil || m.Value == nil {
		return n, err
	}
	temp, err := m.Value.WriteTo(w)
	return n + temp, err
}

// PartialDataComponentMatcher 是「元件符合某個條件」:先一個 bool,true 時 Type 是
// minecraft:data_component_predicate_type 的編號,false 時是元件種類的編號;後面接 NBT 條件。
//
//codec:gen
type PartialDataComponentMatcher struct {
	IsPredicateType bool
	Type            int32          `mc:"VarInt"`
	Predicate       nbt.RawMessage `mc:"NBT"`
}

func (*CanPlaceOn) ID() string {
	return "minecraft:can_place_on"
}
