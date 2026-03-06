package component

//codec:gen
type RepairCost struct {
	Cost int32 `mc:"VarInt"`
}

func (*RepairCost) ID() string {
	return "minecraft:repair_cost"
}
