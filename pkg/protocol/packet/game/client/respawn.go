package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
)

var _ ClientboundPacket = (*Respawn)(nil)

//codec:gen
type Respawn struct {
	CommonPlayerSpawnInfo CommonPlayerSpawnInfo
	DataToKeep            uint8
}

func (Respawn) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundRespawn
}
