package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
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
