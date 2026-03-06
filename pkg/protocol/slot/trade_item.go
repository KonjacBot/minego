package slot

import (
	"io"

	pk "github.com/KonjacBot/go-mc/net/packet"
)

type TradeSlot struct {
	ID         int32
	Count      int32
	Components map[int32]Component
}

func (t TradeSlot) WriteTo(w io.Writer) (n int64, err error) {
	pk.VarInt(t.ID).WriteTo(w)
	pk.VarInt(t.Count).WriteTo(w)
	pk.VarInt(len(t.Components)).WriteTo(w)
	for id, component := range t.Components {
		pk.VarInt(id).WriteTo(w)
		component.WriteTo(w)
	}
	return
}

func (t *TradeSlot) ReadFrom(r io.Reader) (n int64, err error) {
	(*pk.VarInt)(&t.ID).ReadFrom(r)
	(*pk.VarInt)(&t.Count).ReadFrom(r)
	var lens pk.VarInt
	lens.ReadFrom(r)
	t.Components = make(map[int32]Component)
	for i := int32(0); i < int32(lens); i++ {
		var id pk.VarInt
		id.ReadFrom(r)
		c := ComponentFromID(int(id))
		c.ReadFrom(r)
		t.Components[int32(id)] = c
	}

	return
}
