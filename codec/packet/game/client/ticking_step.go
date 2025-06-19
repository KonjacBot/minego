package client

//codec:gen
type StepTick struct {
	TickSteps int32 `mc:"VarInt"`
}
