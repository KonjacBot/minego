package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type WritableBookContent struct {
	Pages []WritableBookPage
}

//codec:gen
type WritableBookPage struct {
	RawContent      string `mc:"String"`
	FilteredContent pk.Option[wire.String, *wire.String]
}

func (*WritableBookContent) ID() string {
	return "minecraft:writable_book_content"
}
