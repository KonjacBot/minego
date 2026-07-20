package client

import "github.com/KonjacBot/go-mc/chat"

//codec:gen
type ObjectivesData struct {
	Value           chat.Message
	RenderType      int32 `mc:"VarInt"`
	HasNumberFormat bool
	//opt:optional:HasNumberFormat
	NumberFormat *ScoreNumberFormat
}

type ObjectivesCreateData struct {
	ObjectivesData
}

type ObjectivesUpdateData struct {
	ObjectivesData
}

//codec:gen
type UpdateObjectives struct {
	ObjectiveName string
	Mode          int8
	//opt:enum:Mode:0
	Create ObjectivesCreateData
	//opt:enum:Mode:2
	Update ObjectivesUpdateData
}
