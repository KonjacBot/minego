package slot

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
)

type Component interface {
	ID() string

	pk.Field
}

type componentCreator func() Component

var components = make(map[int]componentCreator)

func ComponentFromID(id int) Component {
	if components[id] == nil {
		return nil
	}
	return components[id]()
}

func RegisterComponent(c componentCreator) {
	components[len(components)] = c
}
