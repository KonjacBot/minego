package slot

import (
	"fmt"
	"io"

	pk "github.com/KonjacBot/go-mc/net/packet"
)

type Component interface {
	ID() string

	pk.Field
}

type componentCreator func() Component

const maxComponentPatchEntries = 32767

const maxComponentDecodeDepth = 64

type componentDecodeReader struct {
	io.Reader
	depth int
}

func enterComponentDecode(r io.Reader) (*componentDecodeReader, func(), error) {
	reader, ok := r.(*componentDecodeReader)
	if !ok {
		reader = &componentDecodeReader{Reader: r}
	}
	if reader.depth >= maxComponentDecodeDepth {
		return nil, nil, fmt.Errorf("item component nesting exceeds %d", maxComponentDecodeDepth)
	}
	reader.depth++
	return reader, func() { reader.depth-- }, nil
}

var components = make(map[int]componentCreator)

func ComponentFromID(id int) Component {
	if components[id] == nil {
		return nil
	}
	return components[id]()
}

func validateComponentPatch(added map[int32]Component, removed []int32) error {
	if len(added) > maxComponentPatchEntries || len(removed) > maxComponentPatchEntries {
		return fmt.Errorf("component patch count greater than %d", maxComponentPatchEntries)
	}
	for id, component := range added {
		if id < 0 {
			return fmt.Errorf("data component id %d is negative", id)
		}
		if component == nil {
			return fmt.Errorf("data component %d is nil", id)
		}
	}
	seen := make(map[int32]struct{}, len(removed))
	for _, id := range removed {
		if id < 0 {
			return fmt.Errorf("removed data component id %d is negative", id)
		}
		if _, exists := added[id]; exists {
			return fmt.Errorf("data component id %d is both added and removed", id)
		}
		if _, exists := seen[id]; exists {
			return fmt.Errorf("duplicate removed data component id %d", id)
		}
		seen[id] = struct{}{}
	}
	return nil
}

func RegisterComponent(c componentCreator) {
	components[len(components)] = c
}
