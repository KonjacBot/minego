package client

import "git.konjactw.dev/falloutBot/go-mc/chat"

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
