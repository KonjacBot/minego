package client

type Commands struct {
	Nodes     []CommandNode
	RootIndex int32 `mc:"VarInt"`
}

type CommandNode struct {
	Flags          int8
	Children       []int32
	Redirect       int32 `mc:"VarInt"`
	Name           string
	Parser         CommandParser
	SuggestionType string `mc:"Identifier"`
}

type CommandParser struct {
	ID         int32 `mc:"VarInt"`
	Properties []byte
}
