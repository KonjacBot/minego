package client

//codec:gen
type SetExperience struct {
	ExperienceBar   float32
	Level           int32 `mc:"VarInt"`
	TotalExperience int32 `mc:"VarInt"`
}
