package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type Repairable struct {
	Items pk.IDSet
}

func (*Repairable) ID() string {
	return "minecraft:repairable"
}
