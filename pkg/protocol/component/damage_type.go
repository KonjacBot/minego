package component

import (
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

// DamageType 是武器造成的傷害種類:26.2 是 holderRegistry,只有一個 minecraft:damage_type 的登錄表編號,不能內嵌。
//
//codec:gen
type DamageType struct {
	Type int32 `mc:"VarInt"`
}

func (*DamageType) ID() string {
	return "minecraft:damage_type"
}
