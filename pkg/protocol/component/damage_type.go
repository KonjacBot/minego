package component

import (
	"git.konjactw.dev/falloutBot/go-mc/net/packet"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type DamageTypeData struct {
	MessageId  string
	Scaling    string
	Exhaustion float32

	Effects          pk.Option[pk.String, *pk.String]
	DeathMessageType pk.Option[pk.String, *pk.String]
}

//codec:gen
type DamageType struct {
	Data packet.OptID[DamageTypeData, *DamageTypeData]
}

func (*DamageType) ID() string {
	return "minecraft:damage_type"
}
