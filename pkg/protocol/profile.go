package protocol

import (
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
	"git.konjactw.dev/falloutBot/go-mc/yggdrasil/user"
	"github.com/google/uuid"
)

//codec:gen
type GameProfile struct {
	UUID       uuid.UUID `mc:"UUID"`
	Name       string
	Properties []user.Property
}

//codec:gen
type PartialProfile struct {
	Username   pk.Option[pk.String, *pk.String]
	UUID       pk.Option[pk.UUID, *pk.UUID]
	Properties []user.Property
}

//codec:gen
type ResolvableProfile struct {
	Type int32 `mc:"VarInt"`
	//opt:enum:Type:0
	Partial *PartialProfile
	//opt:enum:Type:1
	GameProfile *ResolvableProfile

	Body   pk.Option[pk.Identifier, *pk.Identifier]
	Cape   pk.Option[pk.Identifier, *pk.Identifier]
	Elytra pk.Option[pk.Identifier, *pk.Identifier]
	Model  pk.Option[pk.VarInt, *pk.VarInt]
}
