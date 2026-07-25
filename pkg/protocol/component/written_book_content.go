package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type WrittenBookContent struct {
	RawTitle      string `mc:"String"`
	FilteredTitle pk.Option[wire.String, *wire.String]
	Author        string `mc:"String"`
	Generation    int32  `mc:"VarInt"`
	Pages         []WrittenBookPage
	Resolved      bool
}

//codec:gen
type WrittenBookPage struct {
	RawContent      wire.Message
	FilteredContent pk.Option[wire.Message, *wire.Message]
}

func (*WrittenBookContent) ID() string {
	return "minecraft:written_book_content"
}
