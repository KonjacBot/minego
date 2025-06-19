package component

import "git.konjactw.dev/patyhank/minego/codec/data/slot"

//codec:gen
type RepairCost struct {
	Cost int32 `mc:"VarInt"`
}

func (*RepairCost) Type() slot.ComponentID {
	return 16
}

func (*RepairCost) ID() string {
	return "minecraft:repair_cost"
}
