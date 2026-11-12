package client

import (
	"github.com/KonjacBot/minego/v2/pkg/protocol/slot"
)

//codec:gen
type SetCursorItem struct {
	CarriedItem slot.Slot
}
