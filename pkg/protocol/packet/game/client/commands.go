package client

import "github.com/Tnze/go-mc/server/command"

//codec:gen
type Commands struct {
	Nodes     []command.Node
	RootIndex int32 `mc:"VarInt"`
}
