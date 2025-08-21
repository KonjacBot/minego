package client

//codec:gen
type SetEntityLink struct {
	AttachedEntityID int32
	HoldingEntityID  int32 // leader, -1 to detach
}
