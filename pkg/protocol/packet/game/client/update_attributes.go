package client

//codec:gen
type AttributeModifier struct {
	Id        string `mc:"Identifier"`
	Amount    float64
	Operation int8
}

//codec:gen
type Attribute struct {
	Id        int32 `mc:"VarInt"`
	Value     float64
	Modifiers []AttributeModifier
}

//codec:gen
type UpdateAttributes struct {
	EntityID   int32 `mc:"VarInt"`
	Attributes []Attribute
}
