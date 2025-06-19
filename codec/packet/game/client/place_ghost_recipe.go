package client

import "git.konjactw.dev/patyhank/minego/codec/slot/display/recipe"

//codec:gen
type PlaceGhostRecipe struct {
	WindowID      int32 `mc:"VarInt"`
	RecipeDisplay recipe.Display
}
