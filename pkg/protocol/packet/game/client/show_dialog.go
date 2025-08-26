package client

import (
	"git.konjactw.dev/falloutBot/go-mc/nbt"
)

//codec:gen
type ShowDialog struct {
	DialogID int32 `mc:"VarInt"`
	//opt:id:DialogID
	DialogData nbt.RawMessage `mc:"NBT"`
}
