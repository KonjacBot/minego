package slot

import (
	"errors"
	"fmt"
	"io"
	"sort"
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

func (s Slot) Clone() Slot {
	cloned := s
	if s.AddComponent != nil {
		cloned.AddComponent = make(map[int32]Component, len(s.AddComponent))
		for id, component := range s.AddComponent {
			cloned.AddComponent[id] = component
		}
	}
	cloned.RemoveComponent = append([]int32(nil), s.RemoveComponent...)
	return cloned
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
	if err := validateComponentPatch(s.AddComponent, s.RemoveComponent); err != nil {
		return 0, err
	}
	temp, err := pk.VarInt(s.Count).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	if s.Count < 0 {
		return n, errors.New("slot count less than zero")
	}
	if s.Count == 0 {
		return n, nil
	}
	if s.ItemID < 0 {
		return n, errors.New("slot item id less than zero")
	}
	temp, err = pk.VarInt(s.ItemID).WriteTo(w)
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

	componentIDs := make([]int, 0, len(s.AddComponent))
	for id := range s.AddComponent {
		componentIDs = append(componentIDs, int(id))
	}
	sort.Ints(componentIDs)
	for _, componentID := range componentIDs {
		id := int32(componentID)
		c := s.AddComponent[id]
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

func (s *Slot) ReadFrom(r io.Reader) (n int64, err error) {
	reader, leave, err := enterComponentDecode(r)
	if err != nil {
		return 0, err
	}
	defer leave()
	r = reader

	*s = Slot{}
	temp, err := (*pk.VarInt)(&s.Count).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	if s.Count < 0 {
		return n, errors.New("slot count less than zero")
	}
	if s.Count == 0 {
		return n, nil
	}

	var itemID int32
	temp, err = (*pk.VarInt)(&itemID).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}

	s.ItemID = item.ID(itemID)
	if itemID < 0 {
		return n, errors.New("slot item id less than zero")
	}

	addLens := int32(0)
	temp, err = (*pk.VarInt)(&addLens).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}

	removeLens := int32(0)
	temp, err = (*pk.VarInt)(&removeLens).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	if addLens < 0 || removeLens < 0 {
		return n, errors.New("slot component count less than zero")
	}
	if addLens > maxComponentPatchEntries || removeLens > maxComponentPatchEntries {
		return n, fmt.Errorf("slot component count greater than %d", maxComponentPatchEntries)
	}

	var id int32
	if addLens > 0 {
		s.AddComponent = make(map[int32]Component)
	}
	for i := int32(0); i < addLens; i++ {
		temp, err = (*pk.VarInt)(&id).ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
		if _, exists := s.AddComponent[id]; exists {
			return n, fmt.Errorf("duplicate slot component id %d", id)
		}
		c := ComponentFromID(int(id))
		if c == nil {
			return n, fmt.Errorf("unknown slot component id %d", id)
		}

		temp, err = c.ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
		s.AddComponent[int32(id)] = c
	}
	for i := int32(0); i < removeLens; i++ {
		temp, err = (*pk.VarInt)(&id).ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
		s.RemoveComponent = append(s.RemoveComponent, id)
	}
	if err := validateComponentPatch(s.AddComponent, s.RemoveComponent); err != nil {
		return n, err
	}
	return n, nil
}
