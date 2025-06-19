package client

type SetBorderLerpSize struct {
	OldDiameter, NewDiameter float64
	Speed                    int64 `mc:"VarLong"`
}
