package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
)

var _ ClientboundPacket = (*Login)(nil)

//codec:gen
type GlobalPos struct {
	Dimension string `mc:"Identifier"`
	Pos       pk.Position
}

//codec:gen
type CommonPlayerSpawnInfo struct {
	DimensionType     int32  `mc:"VarInt"`
	Dimension         string `mc:"Identifier"`
	Seed              int64
	GameType          uint8
	PreviousGameType  int8
	IsDebug           bool
	IsFlat            bool
	LastDeathLocation pk.Option[GlobalPos, *GlobalPos]
	PortalCooldown    int32 `mc:"VarInt"`
	SeaLevel          int32 `mc:"VarInt"`
}

//codec:gen
type Login struct {
	PlayerID              int32 `mc:"VarInt"`
	Hardcore              bool
	Levels                []string `mc:"Identifier"`
	MaxPlayers            int32    `mc:"VarInt"`
	ChunkRadius           int32    `mc:"VarInt"`
	SimulationDistance    int32    `mc:"VarInt"`
	ReducedDebugInfo      bool
	ShowDeathScreen       bool
	DoLimitedCrafting     bool
	CommonPlayerSpawnInfo CommonPlayerSpawnInfo
	EnforcesSecureChat    bool
}

func (Login) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundLogin
}
