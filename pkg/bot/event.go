package bot

import "fmt"

type EventHandler interface {
	PublishEvent(event string, data any) error
	SubscribeEvent(event string, handler func(data any) error)
}

type Event interface {
	EventID() string
}

func PublishEvent(client Client, event Event) error {
	return client.EventHandler().PublishEvent(event.EventID(), event)
}

func SubscribeEvent[T Event](client Client, handler func(event T) error) {
	if client == nil || handler == nil {
		return
	}
	var t T
	client.EventHandler().SubscribeEvent(t.EventID(), func(data any) error {
		typed, ok := data.(T)
		if !ok {
			return fmt.Errorf("event %q payload is %T, want %T", t.EventID(), data, t)
		}
		return handler(typed)
	})
}
