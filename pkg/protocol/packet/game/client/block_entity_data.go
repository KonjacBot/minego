package client

import (
	"git.konjactw.dev/falloutBot/go-mc/nbt"
	"git.konjactw.dev/falloutBot/go-mc/net/packet"
)

// codec:gen
type BlockEntityData struct {
	Position packet.Position
	Type     int32          `mc:"VarInt"`
	Data     nbt.RawMessage `mc:"NBT"`
}
