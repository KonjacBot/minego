package world

import "github.com/KonjacBot/minego/v2/pkg/bot"

type EntityRemoveEvent struct {
	Entity bot.Entity
}

type EntityAddEvent struct {
	EntityID int32
}
