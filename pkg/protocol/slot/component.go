package slot

import (
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
)

type Component interface {
	Type() ComponentID
	ID() string

	pk.Field
}

type ComponentID int32
type componentCreator func() Component

var components = make(map[ComponentID]componentCreator)

func ComponentFromID(id ComponentID) Component {
	return components[id]()
}

func RegisterComponent(c componentCreator) {
	components[c().Type()] = c
}
