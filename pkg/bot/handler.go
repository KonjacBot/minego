package bot

import (
	"context"

	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	"git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/client"
)

type PacketHandler interface {
	AddPacketHandler(id packetid.ClientboundPacketID, handler func(ctx context.Context, p client.ClientboundPacket))
	AddRawPacketHandler(id packetid.ClientboundPacketID, handler func(ctx context.Context, p packet.Packet))
	AddGenericPacketHandler(handler func(ctx context.Context, p client.ClientboundPacket))
	HandlePacket(ctx context.Context, p client.ClientboundPacket)
}

type HandlerFunc[T client.ClientboundPacket] func(ctx context.Context, p T)

func AddHandler[T client.ClientboundPacket](c Client, f HandlerFunc[T]) {
	var t T
	handler := c.PacketHandler()
	handler.AddPacketHandler(t.PacketID(), func(ctx context.Context, p client.ClientboundPacket) {
		f(ctx, p.(T))
	})
}
