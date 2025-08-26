package client

import (
	"git.konjactw.dev/falloutBot/go-mc/chat"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type AdvancementDisplay struct {
	Title                chat.Message
	Description          chat.Message
	Icon                 slot.Slot
	FrameType            int32 `mc:"VarInt"`
	Flags                int32
	HasBackgroundTexture bool
	//opt:optional:HasBackgroundTexture
	BackgroundTexture string `mc:"Identifier"`
	X, Y              float32
}

//codec:gen
type AdvancementRequirements struct {
	OR []string
}

//codec:gen
type Advancement struct {
	ID          string `mc:"Identifier"`
	HasParentID bool
	//opt:optional:HasParentID
	ParentID       string
	HasDisplayData bool
	//opt:optional:HasDisplayData
	DisplayData       AdvancementDisplay
	Requirements      []AdvancementRequirements
	SendTelemetryData bool
}

//codec:gen
type AdvancementProgress struct {
	ID          string `mc:"Identifier"`
	CriterionId string `mc:"Identifier"`
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
