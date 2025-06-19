package server

import (
	"github.com/Tnze/go-mc/data/packetid"
	pk "github.com/Tnze/go-mc/net/packet"
)

//codec:gen
type SetStructureBlock struct {
	Location                  pk.Position
	Action                    int32 `mc:"VarInt"`
	Mode                      int32 `mc:"VarInt"`
	Name                      string
	OffsetX, OffsetY, OffsetZ int8
	SizeX, SizeY, SizeZ       int8
	Mirror                    int32 `mc:"VarInt"`
	Rotation                  int32 `mc:"VarInt"`
	Metadata                  string
	Integrity                 float32
	Seed                      int32 `mc:"VarLong"`
	Flags                     int8
}

func (SetStructureBlock) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundSetStructureBlock
}

func init() {
	registerPacket(packetid.ServerboundSetStructureBlock, func() ServerboundPacket {
		return &SetStructureBlock{}
	})
}
