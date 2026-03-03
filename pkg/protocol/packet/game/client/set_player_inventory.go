package client

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type SetPlayerInventory struct {
	Slot int32 `mc:"VarInt"`
	Data slot.Slot
}
