package slot

import (
	"git.konjactw.dev/patyhank/minego/codec/slot"
	pk "github.com/Tnze/go-mc/net/packet"
	"io"
)

type DisplayType int32

const (
	DisplayEmpty DisplayType = iota
	DisplayItem
	DisplayItemStack
	DisplayTag
	DisplaySmithingTrim
	DisplayWithRemainder
	DisplayComposite
)

type Display struct {
	SlotDisplay
}

func (s Display) WriteTo(w io.Writer) (n int64, err error) {
	pk.VarInt(s.SlotDisplay.SlotDisplayType()).WriteTo(w)
	s.SlotDisplay.WriteTo(w)
	return
}

func (s *Display) ReadFrom(r io.Reader) (n int64, err error) {
	var displayType DisplayType
	_, err = (*pk.VarInt)(&displayType).ReadFrom(r)
	if err != nil {
		return
	}
	switch displayType {
	case DisplayEmpty:
		return
	case DisplayItem:
		var item Item
		if _, err = item.ReadFrom(r); err != nil {
			return
		}
	case DisplayItemStack:
		var itemStack ItemStack
		if _, err = itemStack.ReadFrom(r); err != nil {
			return
		}
	case DisplayTag:
		var tag Tag
		if _, err = tag.ReadFrom(r); err != nil {
			return
		}
	case DisplaySmithingTrim:
		var trim SmithingTrim
		if _, err = trim.ReadFrom(r); err != nil {
			return
		}
	case DisplayWithRemainder:
		var remainder WithRemainder
		if _, err = remainder.ReadFrom(r); err != nil {
			return
		}
	case DisplayComposite:
		var composite Composite
		if _, err = composite.ReadFrom(r); err != nil {
			return
		}
	}
	return
}

type SlotDisplay interface {
	SlotDisplayType() DisplayType
	pk.Field
}

//codec:gen
type Item struct {
	ID int32 `mc:"VarInt"`
}

func (i Item) SlotDisplayType() DisplayType {
	return DisplayItem
}

//codec:gen
type ItemStack struct {
	ItemStack slot.Slot
}

func (i ItemStack) SlotDisplayType() DisplayType {
	return DisplayItemStack
}

//codec:gen
type Tag struct {
	Tag pk.Identifier
}

func (i Tag) SlotDisplayType() DisplayType {
	return DisplayTag
}

//codec:gen
type SmithingTrim struct {
	Base      Display
	Trim      Display
	Remainder Display
}

func (i SmithingTrim) SlotDisplayType() DisplayType {
	return DisplaySmithingTrim
}

//codec:gen
type WithRemainder struct {
	Ingredient Display
	Remainder  Display
}

func (i WithRemainder) SlotDisplayType() DisplayType {
	return DisplayWithRemainder
}

//codec:gen
type Composite struct {
	Displays []Display
}

func (i Composite) SlotDisplayType() DisplayType {
	return DisplayComposite
}
