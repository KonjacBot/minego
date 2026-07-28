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
var componentIDs = make(map[string]int)

func ComponentFromID(id int) Component {
	if components[id] == nil {
		return nil
	}
	return components[id]()
}

// ComponentID returns the protocol registry ID for the named data component.
func ComponentID(identifier string) (int, bool) {
	id, ok := componentIDs[identifier]
	return id, ok
}

func RegisterComponent(c componentCreator) {
	component := c()
	identifier := component.ID()
	if _, exists := componentIDs[identifier]; exists {
		panic("duplicate data component registration: " + identifier)
	}
	id := len(components)
	components[id] = c
	componentIDs[identifier] = id
}
