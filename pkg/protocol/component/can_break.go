package component

//codec:gen
type CanBreak struct {
	BlockPredicates []BlockPredicate
}

func (*CanBreak) ID() string {
	return "minecraft:can_break"
}
