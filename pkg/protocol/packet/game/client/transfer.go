package client

//codec:gen
type Transfer struct {
	Host string
	Port int32 `mc:"VarInt"`
}
