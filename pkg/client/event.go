package client

import (
	"fmt"
	"sync"
)

// EventHandler 是一個泛型事件總線
type EventHandler struct {
	mu       sync.RWMutex
	handlers map[string][]func(event any) error
}

func (e *EventHandler) PublishEvent(event string, data any) error {
	e.mu.RLock()
	hs := append([]func(event any) error(nil), e.handlers[event]...)
	e.mu.RUnlock()

	for i, h := range hs {
		var callbackErr error
		if err := invokeCallback(fmt.Sprintf("event %q handler[%d]", event, i), func() {
			callbackErr = h(data)
		}); err != nil {
			return err
		}
		if callbackErr != nil {
			return callbackErr
		}
	}
	return nil
}

func (e *EventHandler) SubscribeEvent(event string, handler func(data any) error) {
	if handler == nil {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers[event] = append(e.handlers[event], handler)
}

func NewEventHandler() *EventHandler {
	return &EventHandler{
		handlers: make(map[string][]func(any) error),
	}
}
