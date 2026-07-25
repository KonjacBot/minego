package slot

import (
	"errors"
	"fmt"
	"io"
	"sort"

	pk "github.com/KonjacBot/go-mc/net/packet"
)

type TradeSlot struct {
	ID         int32
	Count      int32
	Components map[int32]Component
}

func (t TradeSlot) WriteTo(w io.Writer) (n int64, err error) {
	if t.ID < 0 {
		return 0, errors.New("trade slot item id less than zero")
	}
	if t.Count < 0 {
		return 0, errors.New("trade slot count less than zero")
	}
	if len(t.Components) > maxComponentPatchEntries {
		return 0, fmt.Errorf("trade slot component count greater than %d", maxComponentPatchEntries)
	}

	temp, err := pk.VarInt(t.ID).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = pk.VarInt(t.Count).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = pk.VarInt(len(t.Components)).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}

	ids := make([]int, 0, len(t.Components))
	for id := range t.Components {
		ids = append(ids, int(id))
	}
	sort.Ints(ids)
	for _, componentID := range ids {
		id := int32(componentID)
		component := t.Components[id]
		if component == nil {
			return n, fmt.Errorf("trade slot component %d is nil", id)
		}
		temp, err = pk.VarInt(id).WriteTo(w)
		n += temp
		if err != nil {
			return n, err
		}
		temp, err = component.WriteTo(w)
		n += temp
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (t *TradeSlot) ReadFrom(r io.Reader) (n int64, err error) {
	reader, leave, err := enterComponentDecode(r)
	if err != nil {
		return 0, err
	}
	defer leave()
	r = reader

	*t = TradeSlot{}
	temp, err := (*pk.VarInt)(&t.ID).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	if t.ID < 0 {
		return n, errors.New("trade slot item id less than zero")
	}
	temp, err = (*pk.VarInt)(&t.Count).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	if t.Count < 0 {
		return n, errors.New("trade slot count less than zero")
	}

	var count pk.VarInt
	temp, err = count.ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	if count < 0 {
		return n, errors.New("trade slot component count less than zero")
	}
	if count > maxComponentPatchEntries {
		return n, fmt.Errorf("trade slot component count greater than %d", maxComponentPatchEntries)
	}
	if count > 0 {
		t.Components = make(map[int32]Component, min(int(count), len(components)))
	}
	for i := int32(0); i < int32(count); i++ {
		var id pk.VarInt
		temp, err = id.ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
		if _, exists := t.Components[int32(id)]; exists {
			return n, fmt.Errorf("duplicate trade slot component id %d", id)
		}
		component := ComponentFromID(int(id))
		if component == nil {
			return n, fmt.Errorf("unknown trade slot component id %d", id)
		}
		temp, err = component.ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
		t.Components[int32(id)] = component
	}

	return n, nil
}
