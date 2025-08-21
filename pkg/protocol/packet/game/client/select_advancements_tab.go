package client

//codec:gen
type SelectAdvancementsTab struct {
	HasIdentifier bool
	//opt:optional:HasIdentifier
	Identifier string `mc:"Identifier"`
}
