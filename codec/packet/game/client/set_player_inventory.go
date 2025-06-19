package client

import "git.konjactw.dev/patyhank/minego/codec/slot"

//codec:gen
type SetPlayerInventory struct {
	Slot int32 `mc:"VarInt"`
	Data slot.Slot
}
