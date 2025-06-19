package client

//codec:gen
type StoreCookie struct {
	Key     string `mc:"Identifier"`
	Payload []int8 `mc:"Byte"`
}
