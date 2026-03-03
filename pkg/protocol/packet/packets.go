package packet

import (
	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"
)
import loginclient "github.com/KonjacBot/minego/pkg/protocol/packet/login/client"
import gameclient "github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
import configclient "github.com/KonjacBot/minego/pkg/protocol/packet/configuration/client"
import loginserver "github.com/KonjacBot/minego/pkg/protocol/packet/login/server"
import gameserver "github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
import configserver "github.com/KonjacBot/minego/pkg/protocol/packet/configuration/server"

type ServerboundPacket interface {
	pk.Field
	PacketID() packetid.ServerboundPacketID
}

type ClientboundPacket interface {
	pk.Field
	PacketID() packetid.ClientboundPacketID
}

type State int32

const (
	StateLogin State = iota
	StateConfig
	StatePlay
)

func GetClientPacket(state State, id int32) ClientboundPacket {
	switch state {
	case StateLogin:
		return loginclient.ClientboundPackets[packetid.ClientboundPacketID(id)]()
	case StateConfig:
		return configclient.ClientboundPackets[packetid.ClientboundPacketID(id)]()
	case StatePlay:
		return gameclient.ClientboundPackets[packetid.ClientboundPacketID(id)]()
	}
	return nil
}

func GetServerPacket(state State, id int32) ServerboundPacket {
	switch state {
	case StateLogin:
		creator := loginserver.ServerboundPackets[packetid.ServerboundPacketID(id)]
		if creator == nil {
			return nil
		}
		return creator()
	case StateConfig:
		creator := configserver.ServerboundPackets[packetid.ServerboundPacketID(id)]
		if creator == nil {
			return nil
		}
		return creator()
	case StatePlay:
		creator := gameserver.ServerboundPackets[packetid.ServerboundPacketID(id)]
		if creator == nil {
			return nil
		}
		return creator()
	}
	return nil
}
