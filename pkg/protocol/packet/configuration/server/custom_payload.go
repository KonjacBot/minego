package server

import (
	"github.com/KonjacBot/go-mc/data/packetid"

	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
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
