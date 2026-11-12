package client

import "github.com/KonjacBot/minego/v2/pkg/protocol/slot/display/recipe"

//codec:gen
type PlaceGhostRecipe struct {
	WindowID      int32 `mc:"VarInt"`
	RecipeDisplay recipe.Display
}
