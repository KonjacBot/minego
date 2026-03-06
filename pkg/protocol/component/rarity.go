package component

//codec:gen
type Rarity struct {
	Rarity int32 `mc:"VarInt"` // 0=Common, 1=Uncommon, 2=Rare, 3=Epic
}

func (*Rarity) ID() string {
	return "minecraft:rarity"
}
