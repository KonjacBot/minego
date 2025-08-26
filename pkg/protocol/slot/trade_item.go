package slot

import (
	"io"

	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
)

type TradeSlot struct {
	ID         int32
	Count      int32
	Components []Component
}

func (t TradeSlot) WriteTo(w io.Writer) (n int64, err error) {
	pk.VarInt(t.ID).WriteTo(w)
	pk.VarInt(t.Count).WriteTo(w)
	pk.VarInt(len(t.Components)).WriteTo(w)
	for _, component := range t.Components {
		pk.VarInt(component.Type()).WriteTo(w)
		component.WriteTo(w)
	}
	return
}

func (t *TradeSlot) ReadFrom(r io.Reader) (n int64, err error) {
	(*pk.VarInt)(&t.ID).ReadFrom(r)
	(*pk.VarInt)(&t.Count).ReadFrom(r)
	var lens pk.VarInt
	lens.ReadFrom(r)
	t.Components = make([]Component, lens)
	for i := range t.Components {
		var id pk.VarInt
		id.ReadFrom(r)
		c := ComponentFromID(ComponentID(id))
		c.ReadFrom(r)
		t.Components[i] = c
	}

	return
}
