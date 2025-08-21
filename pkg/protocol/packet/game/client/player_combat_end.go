package client

//codec:gen
type EndCombat struct {
	Duration int32 `mc:"VarInt"`
}
