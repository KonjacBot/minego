package slot

import (
	"errors"
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
	if s.SlotDisplay == nil {
		s.SlotDisplay = &Empty{}
	}
	n1, err := pk.VarInt(s.SlotDisplay.SlotDisplayType()).WriteTo(w)
	n += n1
	if err != nil {
		return n, err
	}
	n2, err := s.SlotDisplay.WriteTo(w)
	return n + n2, err
}

func (s *Display) ReadFrom(r io.Reader) (n int64, err error) {
	var displayType DisplayType
	n1, err := (*pk.VarInt)(&displayType).ReadFrom(r)
	n += n1
	if err != nil {
		return n, err
	}
	switch displayType {
	case DisplayEmpty:
		s.SlotDisplay = &Empty{}
	case DisplayAnyFuel:
		s.SlotDisplay = &AnyFuel{}
	case DisplayWithAnyPotion:
		s.SlotDisplay = &WithAnyPotion{}
	case DisplayOnlyWithComponent:
		s.SlotDisplay = &OnlyWithComponent{}
	case DisplayItem:
		s.SlotDisplay = &Item{}
	case DisplayItemStack:
		s.SlotDisplay = &ItemStack{}
	case DisplayTag:
		s.SlotDisplay = &Tag{}
	case DisplayDyed:
		s.SlotDisplay = &Dyed{}
	case DisplaySmithingTrim:
		s.SlotDisplay = &SmithingTrim{}
	case DisplayWithRemainder:
		s.SlotDisplay = &WithRemainder{}
	case DisplayComposite:
		s.SlotDisplay = &Composite{}
	default:
		return n, errors.New("unknown slot display type")
	}
	n2, err := s.SlotDisplay.ReadFrom(r)
	return n + n2, err
}

type SlotDisplay interface {
	SlotDisplayType() DisplayType
	pk.Field
}

type Empty struct{}

func (Empty) SlotDisplayType() DisplayType {
	return DisplayEmpty
}

func (Empty) WriteTo(io.Writer) (int64, error) {
	return 0, nil
}

func (*Empty) ReadFrom(io.Reader) (int64, error) {
	return 0, nil
}

type AnyFuel struct{}

func (AnyFuel) SlotDisplayType() DisplayType {
	return DisplayAnyFuel
}

func (AnyFuel) WriteTo(io.Writer) (int64, error) {
	return 0, nil
}

func (*AnyFuel) ReadFrom(io.Reader) (int64, error) {
	return 0, nil
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
