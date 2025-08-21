package client

import "github.com/google/uuid"

//codec:gen
type WaypointColor struct {
	R, G, B uint8
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

//codec:gen
type Waypoint struct {
	Operation        int32 `mc:"VarInt"`
	IsUUIDIdentifier bool
	//opt:enum:IsUUIDIdentifier:true
	UUID uuid.UUID `mc:"UUID"`
	//opt:enum:IsUUIDIdentifier:false
	Name     string
	HasColor bool
	//opt:optional:HasColor
	Color        WaypointColor
	WaypointType int32 `mc:"VarInt"`
	//opt:enum:WaypointType:1
	WaypointPlayerPos WaypointVec3i
	//opt:enum:WaypointType:2
	WaypointChunkPos WaypointChunkPos
	//opt:enum:WaypointType:3
	WaypointAzimuth WaypointAzimuth
}
