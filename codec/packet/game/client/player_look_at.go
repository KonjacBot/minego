package client

//codec:gen
type LookedAtEntity struct {
	EntityID int32 `mc:"VarInt"`
	LookType int32 `mc:"VarInt"`
}

//codec:gen
type LookAt struct {
	LookType                  int32 `mc:"VarInt"` // Feet = 0 Eyes = 1
	TargetX, TargetY, TargetZ float64
	HasLookedAtEntity         bool
	//opt:optional:HasLookedAtEntity
	LookedAtEntity *LookedAtEntity
}
