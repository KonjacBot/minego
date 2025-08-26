package server

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/server"
)

type ConfigCustomPayload struct {
	server.CustomPayload
}

func (*ConfigCustomPayload) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundConfigCustomPayload
}

func init() {
	registerPacket(packetid.ServerboundConfigCustomPayload, func() ServerboundPacket {
		return &ConfigCustomPayload{}
	})
}
