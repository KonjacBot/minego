package client

//codec:gen
type InitializeWorldBorder struct {
	X, Z                     float64
	OldDiameter, NewDiameter float64
	Speed                    int64 `mc:"VarLong"`
	PortalTeleportBoundary   int32 `mc:"VarInt"`
	WarningBlocks            int32 `mc:"VarInt"`
	WarningTime              int32 `mc:"VarInt"`
}
