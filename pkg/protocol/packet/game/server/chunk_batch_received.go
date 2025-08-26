package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type ChunkBatchReceived struct {
	ChunksPerTick float32
}

func (*ChunkBatchReceived) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundChunkBatchReceived
}

func init() {
	registerPacket(packetid.ServerboundChunkBatchReceived, func() ServerboundPacket {
		return &ChunkBatchReceived{}
	})
}
