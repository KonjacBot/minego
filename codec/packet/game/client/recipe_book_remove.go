package client

//codec:gen
type RecipeBookRemove struct {
	Recipes []int32 `mc:"VarInt"`
}
