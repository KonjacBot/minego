package client

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot/display/slot"
	pk "github.com/Tnze/go-mc/net/packet"
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
