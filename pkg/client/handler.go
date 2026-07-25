package client

import (
	"context"
	"fmt"
	"runtime/debug"
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

// CallbackPanicError identifies a panic raised by a consumer callback. The
// client converts it to a connection error instead of letting it kill the
// process.
type CallbackPanicError struct {
	Callback string
	Value    any
	Stack    []byte
}

func (e *CallbackPanicError) Error() string {
	return fmt.Sprintf("%s panicked: %v", e.Callback, e.Value)
}

func invokeCallback(name string, callback func()) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = &CallbackPanicError{Callback: name, Value: recovered, Stack: debug.Stack()}
		}
	}()
	callback()
	return nil
}

func (ph *packetHandler) AddPacketHandler(id packetid.ClientboundPacketID, handler func(ctx context.Context, p client.ClientboundPacket)) {
	if handler == nil {
		return
	}
	ph.mu.Lock()
	defer ph.mu.Unlock()
	f := ph.handlerMap[id]
	f = append(f, handler)
	ph.handlerMap[id] = f
}

func (ph *packetHandler) AddGenericPacketHandler(handler func(ctx context.Context, p client.ClientboundPacket)) {
	if handler == nil {
		return
	}
	ph.mu.Lock()
	defer ph.mu.Unlock()
	ph.genericMap = append(ph.genericMap, handler)
}

func (ph *packetHandler) AddRawPacketHandler(id packetid.ClientboundPacketID, handler func(ctx context.Context, p pk.Packet)) {
	if handler == nil {
		return
	}
	ph.mu.Lock()
	defer ph.mu.Unlock()
	ph.rawMap[id] = append(ph.rawMap[id], handler)
}

func (ph *packetHandler) HandlePacket(ctx context.Context, p client.ClientboundPacket) {
	if p == nil {
		return
	}
	_ = ph.handlePacket(ctx, p.PacketID(), p)
}

func (ph *packetHandler) handlePacket(ctx context.Context, id packetid.ClientboundPacketID, p client.ClientboundPacket) error {
	ph.mu.RLock()
	genericHandlers := append([]func(context.Context, client.ClientboundPacket){}, ph.genericMap...)
	handlers := append([]func(context.Context, client.ClientboundPacket){}, ph.handlerMap[id]...)
	ph.mu.RUnlock()

	for i, handler := range genericHandlers {
		if err := invokeCallback(fmt.Sprintf("generic packet handler[%d]", i), func() {
			handler(ctx, p)
		}); err != nil {
			return err
		}
	}

	for i, handler := range handlers {
		if err := invokeCallback(fmt.Sprintf("packet %d handler[%d]", id, i), func() {
			handler(ctx, p)
		}); err != nil {
			return err
		}
	}
	return nil
}

func (ph *packetHandler) rawHandlers(id packetid.ClientboundPacketID) []func(context.Context, pk.Packet) {
	ph.mu.RLock()
	defer ph.mu.RUnlock()
	return append([]func(context.Context, pk.Packet){}, ph.rawMap[id]...)
}
