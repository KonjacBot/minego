package client

import (
	"github.com/KonjacBot/go-mc/nbt"
)

//codec:gen
type ShowDialog struct {
	DialogID int32 `mc:"VarInt"`
	//opt:id:DialogID
	DialogData nbt.RawMessage `mc:"NBT"`
}
