package client

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
	chat "github.com/KonjacBot/minego/pkg/protocol/wire"
)

type TestInstanceBlockStatus struct {
	Status chat.Message
	Size   pk.Option[WaypointVec3i, *WaypointVec3i]
}
