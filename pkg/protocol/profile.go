package protocol

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
	"github.com/KonjacBot/minego/pkg/protocol/wire"
	"github.com/google/uuid"
)

//codec:gen
type GameProfile struct {
	UUID       uuid.UUID       `mc:"UUID" json:"id"`
	Name       string          `json:"name"`
	Properties []wire.Property `json:"properties"`
}

//codec:gen
type PartialProfile struct {
	Username   pk.Option[wire.String, *wire.String]
	UUID       pk.Option[pk.UUID, *pk.UUID]
	Properties []wire.Property
}

//codec:gen
type ResolvableProfile struct {
	Type int32 `mc:"VarInt"`
	//opt:enum:Type:0
	Partial *PartialProfile
	//opt:enum:Type:1
	GameProfile *GameProfile

	Body   pk.Option[wire.Identifier, *wire.Identifier]
	Cape   pk.Option[wire.Identifier, *wire.Identifier]
	Elytra pk.Option[wire.Identifier, *wire.Identifier]
	Model  pk.Option[pk.VarInt, *pk.VarInt]
}
