package client

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type SetCursorItem struct {
	CarriedItem slot.Slot
}
