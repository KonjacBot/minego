package client

//codec:gen
type ResetScore struct {
	EntityName       string
	HasObjectiveName bool
	//opt:optional:HasObjectiveName
	ObjectiveName string
}
