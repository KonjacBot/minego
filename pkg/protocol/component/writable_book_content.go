package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
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

func (*WritableBookContent) ID() string {
	return "minecraft:writable_book_content"
}
