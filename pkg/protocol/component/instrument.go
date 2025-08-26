package component

import (
	"git.konjactw.dev/falloutBot/go-mc/chat"
	"git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type Instrument struct {
	Instrument packet.OptID[InstrumentData, *InstrumentData]
}

//codec:gen
type InstrumentData struct {
	SoundEvent  packet.OptID[SoundEvent, *SoundEvent]
	SoundRange  float32
	Range       float32
	Description chat.Message
}

func (*Instrument) Type() slot.ComponentID {
	return 52
}

func (*Instrument) ID() string {
	return "minecraft:instrument"
}
