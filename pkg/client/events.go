package client

import "github.com/KonjacBot/minego/v2/pkg/protocol/packet"

// ConnectionStateChangeEvent is emitted when the connection moves between
// login, configuration and play.
type ConnectionStateChangeEvent struct {
	From packet.State
	To   packet.State
}
