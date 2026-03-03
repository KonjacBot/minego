package client

import (
	"github.com/KonjacBot/go-mc/nbt"
	"github.com/KonjacBot/go-mc/net/packet"
)

// codec:gen
type BlockEntityData struct {
	Position packet.Position
	Type     int32          `mc:"VarInt"`
	Data     nbt.RawMessage `mc:"NBT"`
}
