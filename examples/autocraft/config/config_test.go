package config

import (
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func TestServerDefaultsForLegacyConfig(t *testing.T) {
	c := Config{Server: "server71"}
	if err := toml.Unmarshal([]byte("address = 'mcfallout.net'\n"), &c); err != nil {
		t.Fatal(err)
	}
	if c.Server != "server71" {
		t.Fatalf("Server = %q, want server71", c.Server)
	}
}

func TestTopLevelServerOverridesDefault(t *testing.T) {
	c := Config{Server: "server71"}
	if err := toml.Unmarshal([]byte("server = 'server42'\n"), &c); err != nil {
		t.Fatal(err)
	}
	if c.Server != "server42" {
		t.Fatalf("Server = %q, want server42", c.Server)
	}
}
