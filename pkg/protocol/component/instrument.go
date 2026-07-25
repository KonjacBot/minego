package component

import (
	"github.com/KonjacBot/go-mc/net/packet"
	"github.com/KonjacBot/minego/pkg/protocol/wire"
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
	Description wire.Message
}

func (*Instrument) ID() string {
	return "minecraft:instrument"
}
