package client

import "github.com/KonjacBot/minego/pkg/protocol/slot/display/recipe"

//codec:gen
type PlaceGhostRecipe struct {
	WindowID      int32 `mc:"VarInt"`
	RecipeDisplay recipe.Display
}
