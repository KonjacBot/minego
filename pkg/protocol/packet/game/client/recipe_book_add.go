package client

import (
	"git.konjactw.dev/patyhank/minego/codec/slot/display/recipe"
	pk "github.com/Tnze/go-mc/net/packet"
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
	Ingredients RecipeIngredients
	Flags       int8
}

//codec:gen
type RecipeBookAdd struct {
	Recipes []Recipe
	Replace bool
}
