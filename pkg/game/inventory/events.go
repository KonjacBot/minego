package inventory

import "github.com/KonjacBot/go-mc/chat"

type ContainerOpenEvent struct {
	WindowID int32
	Type     int32
	Title    chat.Message
}
