package slot

import (
	"fmt"
	"io"
	"strings"

	"github.com/KonjacBot/go-mc/level/item"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

type Slot struct {
	Count           int32
	ItemID          item.ID
	AddComponent    map[int32]Component
	RemoveComponent []int32
}

func (s Slot) String() string {
	baseItem := fmt.Sprintf("{%d: %d} ", s.ItemID, s.Count)
	if s.AddComponent != nil {
		var comStrings []string
		for i, component := range s.AddComponent {
			comStrings = append(comStrings, fmt.Sprintf("%d: %#v", i, component))
		}

		baseItem += " [" + strings.Join(comStrings, ", ") + "] "
	}
	if s.RemoveComponent != nil {
		var comStrings []string
		for _, component := range s.RemoveComponent {
			comStrings = append(comStrings, fmt.Sprintf("%d", component))
		}

		baseItem += " (" + strings.Join(comStrings, ", ") + ") "
	}

	return baseItem
}

func (s *Slot) WriteTo(w io.Writer) (n int64, err error) {
	temp, err := pk.VarInt(s.Count).WriteTo(w)
	if s.Count <= 0 || err != nil {
		return temp, err
	}
	n += temp
	temp, err = pk.VarInt(s.ItemID).WriteTo(w)
	n += temp
	if err != nil {
		return temp, err
	}

	temp, err = pk.VarInt(len(s.AddComponent)).WriteTo(w)
	n += temp
	if err != nil {
		return temp, err
	}

	temp, err = pk.VarInt(len(s.RemoveComponent)).WriteTo(w)
	n += temp

	for id, c := range s.AddComponent {
		temp, err = pk.VarInt(id).WriteTo(w)
		n += temp
		if err != nil {
			return temp, err
		}
		temp, err = c.WriteTo(w)
		n += temp
		if err != nil {
			return 0, err
		}
	}

	if err != nil {
		return temp, err
	}
	for _, id := range s.RemoveComponent {
		temp, err = pk.VarInt(id).WriteTo(w)
		n += temp
		if err != nil {
			return temp, err
		}
	}
	return temp, nil
}

func (s *Slot) ReadFrom(r io.Reader) (n int64, err error) {
	temp, err := (*pk.VarInt)(&s.Count).ReadFrom(r)
	if s.Count <= 0 || err != nil {
		return temp, err
	}
	n += temp

	var itemID int32
	temp, err = (*pk.VarInt)(&itemID).ReadFrom(r)
	n += temp
	if err != nil {
		return temp, err
	}

	s.ItemID = item.ID(itemID)

	addLens := int32(0)
	temp, err = (*pk.VarInt)(&addLens).ReadFrom(r)
	n += temp
	if err != nil {
		return temp, err
	}

	removeLens := int32(0)
	temp, err = (*pk.VarInt)(&removeLens).ReadFrom(r)
	n += temp
	if err != nil {
		return temp, err
	}

	var id int32
	if addLens > 0 {
		s.AddComponent = make(map[int32]Component)
	}
	for i := int32(0); i < addLens; i++ {
		temp, err = (*pk.VarInt)(&id).ReadFrom(r)
		n += temp
		if err != nil {
			return temp, err
		}
		c := ComponentFromID(int(id))
		if c == nil {
			return temp, err
		}

		temp, err = c.ReadFrom(r)
		n += temp
		if err != nil {
			return temp, err
		}
		s.AddComponent[int32(id)] = c
	}
	for i := int32(0); i < removeLens; i++ {
		temp, err = (*pk.VarInt)(&id).ReadFrom(r)
		n += temp
		if err != nil {
			return temp, err
		}
		s.RemoveComponent = append(s.RemoveComponent, id)
	}
	return n, nil
}
