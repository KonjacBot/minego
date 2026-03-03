package client

import (
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/slot/display/recipe"
)

//codec:gen
type RecipeIngredients struct {
	Data []pk.IDSet
}

//codec:gen
type Recipe struct {
	RecipeID       int32 `mc:"VarInt"`
	Display        recipe.Display
	GroupID        int32 `mc:"VarInt"`
	CategoryID     int32 `mc:"VarInt"`
	HasIngredients bool
	//opt:optional:HasIngredients
	Ingredients []pk.IDSet
	Flags       int8
}

//codec:gen
type RecipeBookAdd struct {
	Recipes []Recipe
	Replace bool
}
