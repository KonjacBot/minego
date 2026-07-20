package client

import (
	"github.com/KonjacBot/go-mc/chat"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

type TestInstanceBlockStatus struct {
	Status chat.Message
	Size   pk.Option[WaypointVec3i, *WaypointVec3i]
}
