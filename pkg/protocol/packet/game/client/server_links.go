package client

import chat "github.com/KonjacBot/minego/pkg/protocol/wire"

//codec:gen
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
	Links []ServerLinkData
}
