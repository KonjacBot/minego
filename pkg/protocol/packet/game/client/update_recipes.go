package client

import (
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot/display/slot"
)

//codec:gen
type PropertySet struct {
	Id    string  `mc:"Identifier"`
	Items []int32 `mc:"VarInt"`
}

//codec:gen
type StonecutterRecipe struct {
	Ingredient  pk.IDSet
	SlotDisplay slot.Display
}

//codec:gen
type UpdateRecipes struct {
	PropertySets       []PropertySet
	StonecutterRecipes []StonecutterRecipe
}
