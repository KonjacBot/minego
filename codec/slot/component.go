package slot

import pk "github.com/Tnze/go-mc/net/packet"

type Component interface {
	Type() ComponentID
	ID() string

	pk.Field
}

type ComponentID int32

var components = make(map[ComponentID]Component)

func ComponentFromID(id ComponentID) Component {
	return components[id]
}

func RegisterComponent(c Component) {
	components[c.Type()] = c
}
