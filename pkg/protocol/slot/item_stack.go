package slot

import (
	"io"

	"git.konjactw.dev/falloutBot/go-mc/level/item"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
)

type Slot struct {
	Count           int32
	ItemID          item.ID
	AddComponent    []Component
	RemoveComponent []ComponentID
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
	for _, c := range s.AddComponent {
		temp, err = pk.VarInt(c.Type()).WriteTo(w)
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

	temp, err = pk.VarInt(len(s.RemoveComponent)).WriteTo(w)
	n += temp
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
	for i := int32(0); i < addLens; i++ {
		temp, err = (*pk.VarInt)(&id).ReadFrom(r)
		n += temp
		if err != nil {
			return temp, err
		}
		c := ComponentFromID(ComponentID(id))

		temp, err = c.ReadFrom(r)
		n += temp
		if err != nil {
			return temp, err
		}
		s.AddComponent = append(s.AddComponent, c)
	}

	for i := int32(0); i < removeLens; i++ {
		temp, err = (*pk.VarInt)(&id).ReadFrom(r)
		n += temp
		if err != nil {
			return temp, err
		}
		s.RemoveComponent = append(s.RemoveComponent, ComponentID(id))
	}
	return n, nil
}
