package inventory

import "github.com/Tnze/go-mc/chat"

type ContainerOpenEvent struct {
	WindowID int32
	Type     int32
	Title    chat.Message
}

func (c ContainerOpenEvent) EventID() string {
	return "inventory:container_open"
}
