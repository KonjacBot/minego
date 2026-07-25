package client

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
	chat "github.com/KonjacBot/minego/pkg/protocol/wire"
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
	TeamPrefix        chat.Message
	TeamSuffix        chat.Message
	NameTagVisibility int32 `mc:"VarInt"`
	CollisionRule     int32 `mc:"VarInt"`
	TeamColor         pk.Option[pk.VarInt, *pk.VarInt]
	FriendlyFlags     int8
	Entities          []string `mc:"String"`
}

//codec:gen
type UpdateTeamsRemoveTeam struct {
}

//codec:gen
type UpdateTeamsUpdateTeam struct {
	DisplayName       chat.Message
	TeamPrefix        chat.Message
	TeamSuffix        chat.Message
	NameTagVisibility int32 `mc:"VarInt"`
	CollisionRule     int32 `mc:"VarInt"`
	TeamColor         pk.Option[pk.VarInt, *pk.VarInt]
	FriendlyFlags     int8
}

//codec:gen
type UpdateTeamsAddEntities struct {
	Entities []string `mc:"String"`
}

//codec:gen
type UpdateTeamsRemoveEntities struct {
	Entities []string `mc:"String"`
}
