package client

import (
	"io"

	"github.com/KonjacBot/go-mc/chat"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type AdvancementDisplay struct {
	Title       chat.Message
	Description chat.Message
	Icon        slot.Slot
	FrameType   int32 `mc:"VarInt"`
	Flag        AdvancementDisplayFlag
	X, Y        float32
}

type AdvancementDisplayFlag struct {
	Flag     int32
	AssetsId string
}

func (a AdvancementDisplayFlag) WriteTo(w io.Writer) (n int64, err error) {
	if a.AssetsId != "" {
		a.Flag |= 1
	}

	n, err = pk.Int(a.Flag).WriteTo(w)
	if err != nil {
		return n, err
	}
	var nn int64
	if a.AssetsId != "" {
		nn, err = pk.Identifier(a.AssetsId).WriteTo(w)
		if err != nil {
			return n + nn, err
		}
	}
	return n + nn, err
}

func (a *AdvancementDisplayFlag) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = (*pk.Int)(&a.Flag).ReadFrom(r)
	if err != nil {
		return n, err
	}
	var nn int64
	if a.Flag&1 != 0 {
		nn, err = (*pk.Identifier)(&a.AssetsId).ReadFrom(r)
		if err != nil {
			return n + nn, err
		}
	}

	return n + nn, err
}

//codec:gen
type AdvancementRequirements struct {
	OR []string
}

//codec:gen
type Advancement struct {
	ID             string `mc:"Identifier"`
	ParentID       pk.Option[pk.Identifier, *pk.Identifier]
	HasDisplayData bool
	//opt:optional:HasDisplayData
	DisplayData       AdvancementDisplay
	Requirements      []AdvancementRequirements
	SendTelemetryData bool
}

//codec:gen
type AdvancementProgress struct {
	ID       string `mc:"Identifier"`
	Criteria []AdvancementProgressCriteria
}

//codec:gen
type AdvancementProgressCriteria struct {
	CriterionId string
	HasAchieved bool
	//opt:optional:HasAchieved
	AchievingDate int64
}

//codec:gen
type UpdateAdvancements struct {
	Clear                 bool
	Advancements          []Advancement
	RemovedIds            []string `mc:"Identifier"`
	Progress              []AdvancementProgress
	ShowAdvancementsToast bool
}
