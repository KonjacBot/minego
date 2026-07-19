package client

import (
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

	for _, h := range hs {
		if err := h(data); err != nil {
			return err
		}
	}
	return nil
}

func (e *EventHandler) SubscribeEvent(event string, handler func(data any) error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers[event] = append(e.handlers[event], handler)
}

func NewEventHandler() *EventHandler {
	return &EventHandler{
		handlers: make(map[string][]func(any) error),
	}
}
