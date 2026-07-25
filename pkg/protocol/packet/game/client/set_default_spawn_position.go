package client

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type SetDefaultSpawnPosition struct {
	DimensionName wire.Identifier
	Location      pk.Position
	Yaw           float32
	Pitch         float32
}
