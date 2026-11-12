package client

import "testing"

func TestRegistryEntry(t *testing.T) {
	b := &Client{registryEntries: map[string]map[int32]string{
		"minecraft:damage_type": {0: "minecraft:arrow", 26: "minecraft:mace_smash"},
	}}
	tests := []struct {
		registry string
		id       int32
		name     string
		ok       bool
	}{
		{"minecraft:damage_type", 26, "minecraft:mace_smash", true},
		{"minecraft:damage_type", 99, "", false},
		{"minecraft:biome", 0, "", false},
	}
	for _, tt := range tests {
		name, ok := b.RegistryEntry(tt.registry, tt.id)
		if name != tt.name || ok != tt.ok {
			t.Errorf("RegistryEntry(%q, %d) = %q, %v; want %q, %v", tt.registry, tt.id, name, ok, tt.name, tt.ok)
		}
	}
	if _, ok := (&Client{}).RegistryEntry("minecraft:damage_type", 0); ok {
		t.Error("RegistryEntry before configuration should report false")
	}
}
