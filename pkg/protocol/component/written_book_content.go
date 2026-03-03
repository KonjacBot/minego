package component

import (
	"github.com/KonjacBot/go-mc/chat"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type WrittenBookContent struct {
	RawTitle         string `mc:"String"`
	HasFilteredTitle bool
	FilteredTitle    pk.Option[pk.String, *pk.String]
	Author           string `mc:"String"`
	Generation       int32  `mc:"VarInt"`
	Pages            []WrittenBookPage
}

//codec:gen
type WrittenBookPage struct {
	RawContent         chat.Message
	HasFilteredContent bool
	FilteredContent    pk.Option[chat.Message, *chat.Message]
}

func (*WrittenBookContent) Type() slot.ComponentID {
	return 46
}

func (*WrittenBookContent) ID() string {
	return "minecraft:written_book_content"
}
