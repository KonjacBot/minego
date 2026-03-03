package client

import (
	"github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/component"
)

//codec:gen
type EntitySoundEffect struct {
	SoundEvent    packet.OptID[component.SoundEvent, *component.SoundEvent]
	SoundCategory int32 `mc:"VarInt"`
	EntityID      int32 `mc:"VarInt"`
	Volume        float32
	Pitch         float32
	Seed          int32
}
