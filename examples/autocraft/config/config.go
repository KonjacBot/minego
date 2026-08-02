package config

import (
	"os"

	"github.com/KonjacBot/minego/pkg/bot"
	"github.com/KonjacBot/minego/pkg/protocol"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Address     string              `toml:"address"`
	Server      string              `toml:"server"`
	Proxy       *bot.ProxyConfig    `toml:"proxy,omitempty"`
	UserCode    string              `toml:"user_code"`
	Owners      []string            `toml:"owners"`
	TakePos     protocol.Position   `toml:"take_pos"`
	TakeButton  protocol.Position   `toml:"take_button"`
	PlacePos    []protocol.Position `toml:"place_pos"`
	PlaceButton []protocol.Position `toml:"place_button"`
}

func ReadConfig() (c Config, err error) {
	_, err = os.Stat("config.toml")
	if err != nil {
		if os.IsNotExist(err) {
			data, err := toml.Marshal(Config{
				Address:     "mcfallout.net",
				Server:      "server71",
				Proxy:       &bot.ProxyConfig{},
				UserCode:    "artif",
				Owners:      []string{},
				TakePos:     protocol.Position{0, 1, 2},
				PlacePos:    []protocol.Position{{0, 1, 2}},
				PlaceButton: []protocol.Position{{0, 1, 2}},
				TakeButton:  protocol.Position{0, 1, 2},
			})
			if err != nil {
				return Config{}, err
			}
			err = os.WriteFile("config.toml", data, 0644)
			return Config{}, err
		}
		return Config{}, err
	}

	data, err := os.ReadFile("config.toml")
	if err != nil {
		return Config{}, err
	}

	// Keep the configured target deterministic for older config files that do
	// not yet have a top-level server key. A top-level value still overrides it.
	c.Server = "server71"
	err = toml.Unmarshal(data, &c)
	return c, err
}
