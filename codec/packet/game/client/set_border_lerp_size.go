package client

//codec:gen
type SetBorderLerpSize struct {
	OldDiameter, NewDiameter float64
	Speed                    int64 `mc:"VarLong"`
}
