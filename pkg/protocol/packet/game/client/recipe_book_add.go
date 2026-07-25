package client

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot/display/recipe"
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type RecipeIngredients struct {
	Data []wire.IDSet
}

//codec:gen
type Recipe struct {
	RecipeID       int32 `mc:"VarInt"`
	Display        recipe.Display
	GroupID        int32 `mc:"VarInt"`
	CategoryID     int32 `mc:"VarInt"`
	HasIngredients bool
	//opt:optional:HasIngredients
	Ingredients []wire.IDSet
	Flags       int8
}

//codec:gen
type RecipeBookAdd struct {
	Recipes []Recipe
	Replace bool
}
