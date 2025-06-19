package client

import "github.com/Tnze/go-mc/chat"

type ServerLinkData struct {
	IsBuiltin bool
	//opt:enum:IsBuiltin:true
	Type int32 `mc:"VarInt"`
	//opt:enum:IsBuiltin:false
	Name chat.Message
	URL  string
}

//codec:gen
type ServerLinks struct {
}
