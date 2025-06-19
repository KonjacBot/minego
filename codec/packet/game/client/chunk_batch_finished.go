package client

// codec:gen
type ChunkBatchFinished struct {
	BatchSize int32 `mc:"VarInt"`
}
