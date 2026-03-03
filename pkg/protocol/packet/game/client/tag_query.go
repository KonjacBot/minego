package client

import "github.com/KonjacBot/go-mc/nbt"

//codec:gen
type TagQueryResponse struct {
	TransactionID int32          `mc:"VarInt"`
	NBT           nbt.RawMessage `mc:"NBT"`
}
