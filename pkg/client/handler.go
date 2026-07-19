package client

import (
	"context"
	"sync"

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
	mu         sync.RWMutex
	handlerMap map[packetid.ClientboundPacketID][]func(ctx context.Context, p client.ClientboundPacket)
	genericMap []func(ctx context.Context, p client.ClientboundPacket)
	rawMap     map[packetid.ClientboundPacketID][]func(ctx context.Context, p pk.Packet)
}

func (ph *packetHandler) AddPacketHandler(id packetid.ClientboundPacketID, handler func(ctx context.Context, p client.ClientboundPacket)) {
	ph.mu.Lock()
	defer ph.mu.Unlock()
	f := ph.handlerMap[id]
	f = append(f, handler)
	ph.handlerMap[id] = f
}

func (ph *packetHandler) AddGenericPacketHandler(handler func(ctx context.Context, p client.ClientboundPacket)) {
	ph.mu.Lock()
	defer ph.mu.Unlock()
	ph.genericMap = append(ph.genericMap, handler)
}

func (ph *packetHandler) AddRawPacketHandler(id packetid.ClientboundPacketID, handler func(ctx context.Context, p pk.Packet)) {
	ph.mu.Lock()
	defer ph.mu.Unlock()
	ph.rawMap[id] = append(ph.rawMap[id], handler)
}

func (ph *packetHandler) HandlePacket(ctx context.Context, p client.ClientboundPacket) {
	ph.mu.RLock()
	genericHandlers := append([]func(context.Context, client.ClientboundPacket){}, ph.genericMap...)
	handlers := append([]func(context.Context, client.ClientboundPacket){}, ph.handlerMap[p.PacketID()]...)
	ph.mu.RUnlock()

	for _, handler := range genericHandlers {
		handler(ctx, p)
	}

	for _, handler := range handlers {
		handler(ctx, p)
	}
}

func (ph *packetHandler) rawHandlers(id packetid.ClientboundPacketID) []func(context.Context, pk.Packet) {
	ph.mu.RLock()
	defer ph.mu.RUnlock()
	return append([]func(context.Context, pk.Packet){}, ph.rawMap[id]...)
}
