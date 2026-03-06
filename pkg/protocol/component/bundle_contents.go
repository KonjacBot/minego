package component

import (
	slot2 "github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type BundleContents struct {
	Items []slot2.Slot
}

func (*BundleContents) ID() string {
	return "minecraft:bundle_contents"
}
