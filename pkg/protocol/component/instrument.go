package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/net/packet"
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
