package component

import (
	"git.konjactw.dev/patyhank/minego/codec/data/slot"
	pk "github.com/Tnze/go-mc/net/packet"
)

//codec:gen
type WritableBookContent struct {
	Pages []WritableBookPage
}

//codec:gen
type WritableBookPage struct {
	RawContent         string `mc:"String"`
	HasFilteredContent bool
	FilteredContent    pk.Option[pk.String, *pk.String]
}

func (*WritableBookContent) Type() slot.ComponentID {
	return 45
}

func (*WritableBookContent) ID() string {
	return "minecraft:writable_book_content"
}
