package client

import (
	"io"

	"github.com/KonjacBot/go-mc/data/packetid"
	"github.com/KonjacBot/go-mc/level"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

var _ ClientboundPacket = (*ForgetLevelChunk)(nil)

//codec:gen
type ForgetLevelChunk struct {
	Pos PackedChunkPos
}

func (ForgetLevelChunk) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundForgetLevelChunk
}

// PackedChunkPos 是原版 FriendlyByteBuf.writeChunkPos 的區塊座標:整個是一個 long,Z 在高 32 位元,
// 所以線上先 Z 後 X。照 level.ChunkPos(先 X 後 Z 兩個 int)讀會把 X、Z 讀反。
type PackedChunkPos level.ChunkPos

func (c *PackedChunkPos) ReadFrom(r io.Reader) (int64, error) {
	var v pk.Long
	n, err := v.ReadFrom(r)
	if err != nil {
		return n, err
	}
	*c = PackedChunkPos{int32(v), int32(v >> 32)}
	return n, nil
}

func (c PackedChunkPos) WriteTo(w io.Writer) (int64, error) {
	return pk.Long(int64(uint32(c[0])) | int64(c[1])<<32).WriteTo(w)
}
