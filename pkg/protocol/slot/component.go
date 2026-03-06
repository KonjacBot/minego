package slot

import (
	"fmt"

	pk "github.com/KonjacBot/go-mc/net/packet"
)

type Component interface {
	ID() string

	pk.Field
}

type componentCreator func() Component

var components = make(map[int]componentCreator)

func ComponentFromID(id int) Component {
	fmt.Println(id)
	if components[id] == nil {
		fmt.Println(components[id], id)
		return nil
	}
	return components[id]()
}

func RegisterComponent(c componentCreator) {
	components[len(components)] = c
}
