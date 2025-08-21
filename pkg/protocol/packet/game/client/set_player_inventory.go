package client

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type SetPlayerInventory struct {
	Slot int32 `mc:"VarInt"`
	Data slot.Slot
}
