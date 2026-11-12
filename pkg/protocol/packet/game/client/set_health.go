package client

import "github.com/KonjacBot/minego/v2/pkg/protocol/metadata"

//codec:gen
type SetHealth struct {
	Health         float32
	Food           metadata.VarInt
	FoodSaturation float32
}
