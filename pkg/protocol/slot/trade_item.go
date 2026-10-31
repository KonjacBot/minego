package slot

import (
	"fmt"
	"io"

	pk "github.com/KonjacBot/go-mc/net/packet"
)

// TradeSlot 是村民交易要的物品(26.2 的 ItemCost):物品編號、數量,再一串「元件要完全一樣」的條件,
// 每個條件是元件種類加該元件自己的網路格式。
type TradeSlot struct {
	ID         int32
	Count      int32
	Components map[int32]Component
}

func (t TradeSlot) WriteTo(w io.Writer) (n int64, err error) {
	fields := []pk.FieldEncoder{pk.VarInt(t.ID), pk.VarInt(t.Count), pk.VarInt(len(t.Components))}
	for id, component := range t.Components {
		fields = append(fields, pk.VarInt(id), component)
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

func (t *TradeSlot) ReadFrom(r io.Reader) (n int64, err error) {
	var lens pk.VarInt
	for _, f := range []pk.FieldDecoder{(*pk.VarInt)(&t.ID), (*pk.VarInt)(&t.Count), &lens} {
		temp, err := f.ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
	}
	t.Components = make(map[int32]Component)
	for i := int32(0); i < int32(lens); i++ {
		var id pk.VarInt
		temp, err := id.ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
		// 不認得的元件不知道有多長,後面讀不下去;回錯誤,不要 panic 讓整條連線掛掉。
		c := ComponentFromID(int(id))
		if c == nil {
			return n, fmt.Errorf("unknown trade item component id %d", id)
		}
		temp, err = c.ReadFrom(r)
		n += temp
		if err != nil {
			return n, fmt.Errorf("read trade item %d component %d (%s): %w", t.ID, id, c.ID(), err)
		}
		t.Components[int32(id)] = c
	}
	return n, nil
}
