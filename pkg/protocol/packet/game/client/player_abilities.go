package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
)

var _ ClientboundPacket = (*PlayerAbilities)(nil)

//codec:gen
type PlayerAbilities struct {
	Flags        int8
	FlyingSpeed  float32
	WalkingSpeed float32
}

func (PlayerAbilities) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundPlayerAbilities
}
