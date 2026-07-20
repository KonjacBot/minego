package client

//codec:gen
type StoreCookie struct {
	Key     string `mc:"Identifier"`
	Payload []byte `mc:"ByteArray"`
}
