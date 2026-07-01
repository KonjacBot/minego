package slot

import (
	"io"

	"github.com/KonjacBot/go-mc/chat"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

type DisplayType int32

const (
	DisplayEmpty DisplayType = iota
	DisplayAnyFuel
	DisplayWithAnyPotion
	DisplayOnlyWithComponent
	DisplayItem
	DisplayItemStack
	DisplayTag
	DisplayDyed
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
	case DisplayAnyFuel:
		return
	case DisplayWithAnyPotion:
		var potion WithAnyPotion
		if _, err = potion.ReadFrom(r); err != nil {
			return
		}
		s.SlotDisplay = &potion
	case DisplayOnlyWithComponent:
		var only OnlyWithComponent
		if _, err = only.ReadFrom(r); err != nil {
			return
		}
		s.SlotDisplay = &only
	case DisplayItem:
		var item Item
		if _, err = item.ReadFrom(r); err != nil {
			return
		}
		s.SlotDisplay = &item
	case DisplayItemStack:
		var itemStack ItemStack
		if _, err = itemStack.ReadFrom(r); err != nil {
			return
		}
		s.SlotDisplay = &itemStack
	case DisplayTag:
		var tag Tag
		if _, err = tag.ReadFrom(r); err != nil {
			return
		}
		s.SlotDisplay = &tag
	case DisplayDyed:
		var dyed Dyed
		if _, err = dyed.ReadFrom(r); err != nil {
			return
		}
		s.SlotDisplay = &dyed
	case DisplaySmithingTrim:
		var trim SmithingTrim
		if _, err = trim.ReadFrom(r); err != nil {
			return
		}
		s.SlotDisplay = &trim
	case DisplayWithRemainder:
		var remainder WithRemainder
		if _, err = remainder.ReadFrom(r); err != nil {
			return
		}
		s.SlotDisplay = &remainder
	case DisplayComposite:
		var composite Composite
		if _, err = composite.ReadFrom(r); err != nil {
			return
		}
		s.SlotDisplay = &composite
	}
	return
}

type SlotDisplay interface {
	SlotDisplayType() DisplayType
	pk.Field
}

//codec:gen
type WithAnyPotion struct {
	Display Display
}

func (i WithAnyPotion) SlotDisplayType() DisplayType {
	return DisplayWithAnyPotion
}

//codec:gen
type OnlyWithComponent struct {
	Source    Display
	Component int32 `mc:"VarInt"`
}

func (i OnlyWithComponent) SlotDisplayType() DisplayType {
	return DisplayOnlyWithComponent
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
	ItemStack slot.ItemStackTemplate
}

func (i ItemStack) SlotDisplayType() DisplayType {
	return DisplayItemStack
}

//codec:gen
type Tag struct {
	Tag string `mc:"Identifier"`
}

func (i Tag) SlotDisplayType() DisplayType {
	return DisplayTag
}

//codec:gen
type Dyed struct {
	Dye    Display
	Target Display
}

func (i Dyed) SlotDisplayType() DisplayType {
	return DisplayDyed
}

//codec:gen
type SmithingTrimData struct {
	AssetId     string `mc:"Identifier"`
	Description chat.Message
	Decal       bool
}

//codec:gen
type SmithingTrim struct {
	Base     Display
	Material Display
	Pattern  pk.OptID[SmithingTrimData, *SmithingTrimData]
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
