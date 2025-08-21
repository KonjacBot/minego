package server

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/server"
	"github.com/Tnze/go-mc/data/packetid"
)

type ConfigCustomPayload struct {
	server.CustomPayload
}

func (ConfigCustomPayload) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundConfigCustomPayload
}

func init() {
	registerPacket(packetid.ServerboundConfigCustomPayload, func() ServerboundPacket {
		return &ConfigCustomPayload{}
	})
}
