package client

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type SetCursorItem struct {
	CarriedItem slot.Slot
}
