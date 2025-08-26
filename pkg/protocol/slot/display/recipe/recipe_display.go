package recipe

import (
	"io"

	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot/display/slot"
)

type DisplayType int32

const (
	DisplayCraftingShapeless DisplayType = 0 + iota
	DisplayCraftingShaped
	DisplayFurnace
	DisplayStonecutter
	DisplaySmithing
)

type Display struct {
	Display RecipeDisplay
}

func (d Display) WriteTo(w io.Writer) (n int64, err error) {
	pk.VarInt(d.Display.RecipeType()).WriteTo(w)
	return d.Display.WriteTo(w)
}

func (d *Display) ReadFrom(r io.Reader) (n int64, err error) {
	var displayType DisplayType
	_, err = (*pk.VarInt)(&displayType).ReadFrom(r)
	if err != nil {
		return
	}
	switch displayType {
	case DisplayCraftingShapeless:
		d.Display = new(Shapeless)
	case DisplayCraftingShaped:
		d.Display = new(Shaped)
	case DisplayFurnace:
		d.Display = new(Furnace)
	case DisplayStonecutter:
		d.Display = new(Stonecutter)
	case DisplaySmithing:
		d.Display = new(Smithing)
	}
	if d.Display != nil {
		return d.Display.ReadFrom(r)
	}
	return
}

type RecipeDisplay interface {
	RecipeType() DisplayType
	pk.Field
}

//codec:gen
type Shapeless struct {
	Ingredients     []slot.Display
	Result          slot.Display
	CraftingStation slot.Display
}

func (i Shapeless) RecipeType() DisplayType {
	return DisplayCraftingShapeless
}

//codec:gen
type Shaped struct {
	Width           int32 `mc:"VarInt"`
	Height          int32 `mc:"VarInt"`
	Ingredients     []slot.Display
	Result          slot.Display
	CraftingStation slot.Display
}

func (i Shaped) RecipeType() DisplayType {
	return DisplayCraftingShaped
}

//codec:gen
type Furnace struct {
	Ingredient      slot.Display
	Fuel            slot.Display
	Result          slot.Display
	CraftingStation slot.Display
	CookingTime     int32 `mc:"VarInt"`
	Experience      float32
}

func (i Furnace) RecipeType() DisplayType {
	return DisplayFurnace
}

//codec:gen
type Stonecutter struct {
	Ingredient      slot.Display
	Result          slot.Display
	CraftingStation slot.Display
}

func (i Stonecutter) RecipeType() DisplayType {
	return DisplayStonecutter
}

//codec:gen
type Smithing struct {
	Template        slot.Display
	Base            slot.Display
	Addition        slot.Display
	Result          slot.Display
	CraftingStation slot.Display
}

func (i Smithing) RecipeType() DisplayType {
	return DisplaySmithing
}
