package client

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type SetContainerContent struct {
	WindowID    int32 `mc:"VarInt"`
	StateID     int32 `mc:"VarInt"`
	Slots       []slot.Slot
	CarriedItem slot.Slot
}
