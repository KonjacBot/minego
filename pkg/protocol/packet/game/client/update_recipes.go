package client

import (
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/slot/display/slot"
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
