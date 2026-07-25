package client

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot/display/slot"
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type PropertySet struct {
	Id    string  `mc:"Identifier"`
	Items []int32 `mc:"VarInt"`
}

//codec:gen
type StonecutterRecipe struct {
	Ingredient  wire.IDSet
	SlotDisplay slot.Display
}

//codec:gen
type UpdateRecipes struct {
	PropertySets       []PropertySet
	StonecutterRecipes []StonecutterRecipe
}
