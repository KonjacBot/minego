package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type PaintingVariant struct {
	Width   int32
	Height  int32
	AssetID wire.Identifier
	Title   pk.Option[wire.Message, *wire.Message]
	Author  pk.Option[wire.Message, *wire.Message]
}

func (*PaintingVariant) ID() string {
	return "minecraft:painting/variant"
}
