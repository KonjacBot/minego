package client

import "git.konjactw.dev/patyhank/minego/codec/metadata"

//codec:gen
type SetHealth struct {
	Health         float32
	Food           metadata.VarInt
	FoodSaturation float32
}
