package component

import "git.konjactw.dev/patyhank/minego/codec/slot"

//codec:gen
type BundleContents struct {
	Items []slot.Slot
}

func (*BundleContents) Type() slot.ComponentID {
	return 41
}

func (*BundleContents) ID() string {
	return "minecraft:bundle_contents"
}
