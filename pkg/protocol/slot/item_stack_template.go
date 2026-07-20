package slot

import (
	"errors"
	"fmt"
	"io"

	"github.com/KonjacBot/go-mc/level/item"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

// ItemStackTemplate is a non-empty item descriptor used inside data components
// (introduced in 26.2). Unlike [Slot], it is always present (no leading count
// short-circuit) and its wire order is: VarInt itemID, VarInt count, then a
// DataComponentPatch (added count, removed count, added components, removed ids).
type ItemStackTemplate struct {
	ItemID          item.ID
	Count           int32
	AddComponent    map[int32]Component
	RemoveComponent []int32
}

func (s ItemStackTemplate) WriteTo(w io.Writer) (n int64, err error) {
	if s.ItemID == 0 {
		return 0, errors.New("item stack template must be non-empty")
	}
	if s.Count == 0 {
		return 0, errors.New("item stack template must be non-empty")
	}

	temp, err := pk.VarInt(s.ItemID).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = pk.VarInt(s.Count).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = pk.VarInt(len(s.AddComponent)).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = pk.VarInt(len(s.RemoveComponent)).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	for id, c := range s.AddComponent {
		temp, err = pk.VarInt(id).WriteTo(w)
		n += temp
		if err != nil {
			return n, err
		}
		temp, err = c.WriteTo(w)
		n += temp
		if err != nil {
			return n, err
		}
	}
	for _, id := range s.RemoveComponent {
		temp, err = pk.VarInt(id).WriteTo(w)
		n += temp
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (s *ItemStackTemplate) ReadFrom(r io.Reader) (n int64, err error) {
	*s = ItemStackTemplate{}

	var itemID int32
	temp, err := (*pk.VarInt)(&itemID).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	if itemID == 0 {
		return n, errors.New("item stack template must be non-empty")
	}
	s.ItemID = item.ID(itemID)

	temp, err = (*pk.VarInt)(&s.Count).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	if s.Count == 0 {
		return n, errors.New("item stack template must be non-empty")
	}

	addLens, removeLens := int32(0), int32(0)
	temp, err = (*pk.VarInt)(&addLens).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = (*pk.VarInt)(&removeLens).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}

	if addLens > 0 {
		s.AddComponent = make(map[int32]Component)
	}
	var id int32
	for i := int32(0); i < addLens; i++ {
		temp, err = (*pk.VarInt)(&id).ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
		c := ComponentFromID(int(id))
		if c == nil {
			return n, fmt.Errorf("unknown data component id %d", id)
		}
		temp, err = c.ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
		s.AddComponent[id] = c
	}
	for i := int32(0); i < removeLens; i++ {
		temp, err = (*pk.VarInt)(&id).ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
		s.RemoveComponent = append(s.RemoveComponent, id)
	}
	return n, nil
}
