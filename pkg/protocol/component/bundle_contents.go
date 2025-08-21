package component

import (
	slot2 "git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type BundleContents struct {
	Items []slot2.Slot
}

func (*BundleContents) Type() slot2.ComponentID {
	return 41
}

func (*BundleContents) ID() string {
	return "minecraft:bundle_contents"
}
