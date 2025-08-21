package client

//codec:gen
type Tag struct {
	Name    string  `mc:"Identifier"`
	Entries []int32 `mc:"VarInt"`
}

//codec:gen
type RegistryTag struct {
	Registry string `mc:"Identifier"`
	Tags     []Tag
}

//codec:gen
type UpdateTags struct {
	Data []RegistryTag
}
