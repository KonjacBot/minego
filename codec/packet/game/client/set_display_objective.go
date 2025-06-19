package client

//codec:gen
type DisplayObjective struct {
	Position  int32 `mc:"VarInt"`
	ScoreName string
}
