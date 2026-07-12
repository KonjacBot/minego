package component

import (
	"github.com/KonjacBot/go-mc/chat"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type WrittenBookContent struct {
	RawTitle      string `mc:"String"`
	FilteredTitle pk.Option[pk.String, *pk.String]
	Author        string `mc:"String"`
	Generation    int32  `mc:"VarInt"`
	Pages         []WrittenBookPage
	Resolved      bool
}

//codec:gen
type WrittenBookPage struct {
	RawContent      chat.Message
	FilteredContent pk.Option[chat.Message, *chat.Message]
}

func (*WrittenBookContent) ID() string {
	return "minecraft:written_book_content"
}
