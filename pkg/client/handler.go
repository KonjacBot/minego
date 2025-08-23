package client

import (
	"context"

	"git.konjactw.dev/patyhank/minego/pkg/bot"
	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/client"
	"github.com/Tnze/go-mc/data/packetid"
)

func newPacketHandler() bot.PacketHandler {
	return &packetHandler{
		handlerMap: make(map[packetid.ClientboundPacketID][]func(ctx context.Context, p client.ClientboundPacket)),
	}
}

type packetHandler struct {
	handlerMap map[packetid.ClientboundPacketID][]func(ctx context.Context, p client.ClientboundPacket)
	genericMap []func(ctx context.Context, p client.ClientboundPacket)
}

func (ph *packetHandler) AddPacketHandler(id packetid.ClientboundPacketID, handler func(ctx context.Context, p client.ClientboundPacket)) {
	f := ph.handlerMap[id]
	f = append(f, handler)
	ph.handlerMap[id] = f
}

func (ph *packetHandler) AddGenericPacketHandler(handler func(ctx context.Context, p client.ClientboundPacket)) {
	ph.genericMap = append(ph.genericMap, handler)
}

func (ph *packetHandler) HandlePacket(ctx context.Context, p client.ClientboundPacket) {
	f := ph.handlerMap[p.PacketID()]
	if f != nil {
		for _, handler := range f {
			handler(ctx, p)
		}
	}
	for _, handler := range ph.genericMap {
		handler(ctx, p)
	}
}
