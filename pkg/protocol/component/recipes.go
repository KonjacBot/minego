package component

type Recipes struct {
	RecipeIDs []string
}

func (*Recipes) ID() string {
	return "minecraft:recipes"
}
