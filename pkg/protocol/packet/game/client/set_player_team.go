package client

import (
	"git.konjactw.dev/falloutBot/go-mc/chat"
)

//codec:gen
type UpdateTeams struct {
	TeamName string
	Type     int8
	//opt:enum:Type:0
	CreateTeam UpdateTeamsCreateTeam
	//opt:enum:Type:1
	RemoveTeam UpdateTeamsRemoveTeam
	//opt:enum:Type:2
	UpdateTeam UpdateTeamsUpdateTeam
	//opt:enum:Type:3
	AddEntities UpdateTeamsAddEntities
	//opt:enum:Type:4
	RemoveEntities UpdateTeamsRemoveEntities
}

//codec:gen
type UpdateTeamsCreateTeam struct {
	TeamDisplayName   chat.Message
	FriendlyFlags     int8
	NameTagVisibility int32 `mc:"VarInt"`
	CollisionRule     int32 `mc:"VarInt"`
	TeamColor         int32 `mc:"VarInt"`
	TeamPrefix        chat.Message
	TeamSuffix        chat.Message
	Entities          []string `mc:"String"`
}

//codec:gen
type UpdateTeamsRemoveTeam struct {
}

//codec:gen
type UpdateTeamsUpdateTeam struct {
	DisplayName       chat.Message
	FriendlyFlags     int8
	NameTagVisibility int32 `mc:"VarInt"`
	CollisionRule     int32 `mc:"VarInt"`
	TeamColor         int32 `mc:"VarInt"`
	TeamPrefix        chat.Message
	TeamSuffix        chat.Message
}

//codec:gen
type UpdateTeamsAddEntities struct {
	Entities []string `mc:"String"`
}

//codec:gen
type UpdateTeamsRemoveEntities struct {
	Entities []string `mc:"String"`
}
