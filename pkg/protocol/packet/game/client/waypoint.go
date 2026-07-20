package client

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
	"github.com/google/uuid"
)

//codec:gen
type WaypointColor struct {
	R, G, B uint8
}

type WaypointIcon struct {
	Style string `mc:"Identifier"`
	Color pk.Option[WaypointColor, *WaypointColor]
}

//codec:gen
type WaypointVec3i struct {
	X, Y, Z int32 `mc:"VarInt"`
}

//codec:gen
type WaypointChunkPos struct {
	X, Z int32 `mc:"VarInt"`
}

//codec:gen
type WaypointAzimuth struct {
	Angle float32
}

type Waypoint struct {
	Operation        int32 `mc:"VarInt"`
	IsUUIDIdentifier bool
	//opt:enum:IsUUIDIdentifier:true
	UUID uuid.UUID `mc:"UUID"`
	//opt:enum:IsUUIDIdentifier:false
	Name string
	Icon WaypointIcon
	// 0 = empty, 1 = vec3i, 2 = chunk, 3 = azimuth.
	WaypointType int32 `mc:"VarInt"`
	//opt:enum:WaypointType:1
	WaypointPlayerPos WaypointVec3i
	//opt:enum:WaypointType:2
	WaypointChunkPos WaypointChunkPos
	//opt:enum:WaypointType:3
	WaypointAzimuth WaypointAzimuth
}
