package client

import "git.konjactw.dev/patyhank/minego/codec/slot"

//codec:gen
type SetCursorItem struct {
	CarriedItem slot.Slot
}
