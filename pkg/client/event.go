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
	defer e.mu.RUnlock()
	if hs, ok := e.handlers[event]; ok {
		for _, h := range hs {
			if err := h(data); err != nil {
				return err
			}
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

//func SubscribeEvent[T bot.Event[T]](eb *EventHandler, t T, f func(event T) error) {
//	eb.mu.Lock()
//	defer eb.mu.Unlock()
//
//	eb.handlers[t.ID()] = append(eb.handlers[t.ID()], func(d any) error {
//		t2 := d.(T)
//		return f(t2)
//	})
//}
//
//func PublishEvent[T bot.Event[T]](eb *EventHandler, t T) error {
//	eb.mu.RLock()
//	defer eb.mu.RUnlock()
//	if hs, ok := eb.handlers[t.ID()]; ok {
//		for _, h := range hs {
//			if err := h(t); err != nil {
//				return err
//			}
//		}
//	}
//
//	return nil
//}
//
//func GoPublishEvent[T bot.Event[T]](eb *EventHandler, t T) {
//	eb.mu.RLock()
//	defer eb.mu.RUnlock()
//	if hs, ok := eb.handlers[t.ID()]; ok {
//		for _, h := range hs {
//			go h(t)
//		}
//	}
//}
