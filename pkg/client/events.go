package client

import "github.com/KonjacBot/minego/pkg/protocol/packet"

const EventConnectionStateChange = "client:connection_state_change"

type ConnectionStateChangeEvent struct {
	From packet.State
	To   packet.State
}

func (c ConnectionStateChangeEvent) EventID() string {
	return EventConnectionStateChange
}
