package client

//codec:gen
type ChunkBiomes struct {
	ChunkX int32
	ChunkZ int32
	Data   []byte `mc:"ByteArray"`
}
