package client

import (
	"context"

	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
)

func newPacketHandler() *packetHandler {
	return &packetHandler{
		handlerMap: make(map[packetid.ClientboundPacketID][]func(ctx context.Context, p client.ClientboundPacket)),
		rawMap:     make(map[packetid.ClientboundPacketID][]func(ctx context.Context, p pk.Packet)),
	}
}

type packetHandler struct {
	handlerMap map[packetid.ClientboundPacketID][]func(ctx context.Context, p client.ClientboundPacket)
	genericMap []func(ctx context.Context, p client.ClientboundPacket)
	rawMap     map[packetid.ClientboundPacketID][]func(ctx context.Context, p pk.Packet)
}

func (ph *packetHandler) AddPacketHandler(id packetid.ClientboundPacketID, handler func(ctx context.Context, p client.ClientboundPacket)) {
	f := ph.handlerMap[id]
	f = append(f, handler)
	ph.handlerMap[id] = f
}

func (ph *packetHandler) AddGenericPacketHandler(handler func(ctx context.Context, p client.ClientboundPacket)) {
	ph.genericMap = append(ph.genericMap, handler)
}

func (ph *packetHandler) AddRawPacketHandler(id packetid.ClientboundPacketID, handler func(ctx context.Context, p pk.Packet)) {
	ph.rawMap[id] = append(ph.rawMap[id], handler)
}

func (ph *packetHandler) HandlePacket(ctx context.Context, p client.ClientboundPacket) {
	for _, handler := range ph.genericMap {
		handler(ctx, p)
	}

	f := ph.handlerMap[p.PacketID()]
	for _, handler := range f {
		handler(ctx, p)
	}
}
