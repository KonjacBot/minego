package world

import "github.com/KonjacBot/minego/pkg/bot"

type EntityRemoveEvent struct {
	Entity bot.Entity
}

func (e EntityRemoveEvent) EventID() string {
	return "world:entity_remove"
}

type EntityAddEvent struct {
	EntityID int32
}

func (e EntityAddEvent) EventID() string {
	return "world:entity_add"
}
