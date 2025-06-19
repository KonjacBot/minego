package client

//codec:gen
type SetBorderWarningDelay struct {
	WarningTime int32 `mc:"VarInt"`
}
