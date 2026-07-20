package slot

import (
	"fmt"
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
		return 0, fmt.Errorf("slot display is nil")
	}
	temp, err := pk.VarInt(s.SlotDisplay.SlotDisplayType()).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = s.SlotDisplay.WriteTo(w)
	n += temp
	return n, err
}

func (s *Display) ReadFrom(r io.Reader) (n int64, err error) {
	s.SlotDisplay = nil
	var displayType DisplayType
	temp, err := (*pk.VarInt)(&displayType).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	read := func(display SlotDisplay) error {
		read, readErr := display.ReadFrom(r)
		n += read
		if readErr == nil {
			s.SlotDisplay = display
		}
		return readErr
	}
	switch displayType {
	case DisplayEmpty:
		s.SlotDisplay = Empty{}
	case DisplayAnyFuel:
		s.SlotDisplay = AnyFuel{}
	case DisplayWithAnyPotion:
		err = read(&WithAnyPotion{})
	case DisplayOnlyWithComponent:
		err = read(&OnlyWithComponent{})
	case DisplayItem:
		err = read(&Item{})
	case DisplayItemStack:
		err = read(&ItemStack{})
	case DisplayTag:
		err = read(&Tag{})
	case DisplayDyed:
		err = read(&Dyed{})
	case DisplaySmithingTrim:
		err = read(&SmithingTrim{})
	case DisplayWithRemainder:
		err = read(&WithRemainder{})
	case DisplayComposite:
		err = read(&Composite{})
	default:
		err = fmt.Errorf("unknown slot display type %d", displayType)
	}
	return n, err
}

type SlotDisplay interface {
	SlotDisplayType() DisplayType
	pk.Field
}

type Empty struct{}

func (Empty) SlotDisplayType() DisplayType      { return DisplayEmpty }
func (Empty) ReadFrom(io.Reader) (int64, error) { return 0, nil }
func (Empty) WriteTo(io.Writer) (int64, error)  { return 0, nil }

type AnyFuel struct{}

func (AnyFuel) SlotDisplayType() DisplayType      { return DisplayAnyFuel }
func (AnyFuel) ReadFrom(io.Reader) (int64, error) { return 0, nil }
func (AnyFuel) WriteTo(io.Writer) (int64, error)  { return 0, nil }

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
