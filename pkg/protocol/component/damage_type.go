package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type DamageTypeData struct {
	MessageId  string
	Scaling    string
	Exhaustion float32

	Effects          pk.Option[wire.String, *wire.String]
	DeathMessageType pk.Option[wire.String, *wire.String]
}

//codec:gen
type DamageType struct {
	IsHolder bool
	//opt:enum:IsHolder:true
	HolderID int32 `mc:"VarInt"`
	//opt:enum:IsHolder:false
	HolderType string `mc:"Identifier"`
}

func (*DamageType) ID() string {
	return "minecraft:damage_type"
}
