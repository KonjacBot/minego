package client

import "github.com/Tnze/go-mc/chat"

//codec:gen
type TestInstanceBlockStatus struct {
	Status  chat.Message
	HasSize bool
	//opt:optional:HasSize
	SizeX float64
	//opt:optional:HasSize
	SizeY float64
	//opt:optional:HasSize
	SizeZ float64
}
