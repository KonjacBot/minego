package client

type Tag struct {
	Name    string  `mc:"Identifier"`
	Entries []int32 `mc:"VarInt"`
}

type RegistryTag struct {
	Registry string `mc:"Identifier"`
	Tags     []Tag
}

type UpdateTags struct {
	Data []RegistryTag
}
