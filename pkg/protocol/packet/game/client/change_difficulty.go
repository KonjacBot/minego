package client

// codec:gen
type ChangeDifficulty struct {
	Difficulty int32 `mc:"VarInt"`
	Locked     bool
}
