package slot_test

import (
	"testing"

	_ "github.com/KonjacBot/minego/pkg/protocol/component"
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

func TestComponentIDMatchesRegistryOrder(t *testing.T) {
	tests := []struct {
		identifier string
		want       int
	}{
		{identifier: "minecraft:custom_data", want: 0},
		{identifier: "minecraft:max_stack_size", want: 1},
		{identifier: "minecraft:use_remainder", want: 25},
	}
	for _, test := range tests {
		got, ok := slot.ComponentID(test.identifier)
		if !ok || got != test.want {
			t.Fatalf(
				"ComponentID(%q) = %d,%v, want %d,true",
				test.identifier,
				got,
				ok,
				test.want,
			)
		}
		component := slot.ComponentFromID(got)
		if component == nil || component.ID() != test.identifier {
			t.Fatalf(
				"ComponentFromID(%d) = %#v, want %q",
				got,
				component,
				test.identifier,
			)
		}
	}
}
