package client

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
	slot2 "git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type TradeOption struct {
	Input          slot.TradeSlot
	Output         slot2.Slot
	HasSecondInput bool
	//opt:optional:HasSecondInput
	SecondInput     slot.TradeSlot
	TradeDisabled   bool
	TradeUses       int32
	MaxTradeUses    int32
	Experience      int32
	SpecialPrice    int32
	PriceMultiplier float32
	Demand          int32
}

//codec:gen
type MerchantOffers struct {
	WindowID   int32 `mc:"VarInt"`
	Offers     []TradeOption
	Level      int32 `mc:"VarInt"`
	Experience int32 `mc:"VarInt"`
	IsRegular  bool
	CanRestock bool
}
