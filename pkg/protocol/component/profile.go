package component

import (
	"github.com/KonjacBot/minego/v2/pkg/protocol"
)

//codec:gen
type Profile struct {
	Profile protocol.ResolvableProfile
}

func (*Profile) ID() string {
	return "minecraft:profile"
}
