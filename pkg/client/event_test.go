package client

import (
	"errors"
	"testing"
	"time"
)

func TestEventHandlerAllowsSubscriptionFromHandler(t *testing.T) {
	e := NewEventHandler()
	e.SubscribeEvent("test", func(any) error {
		e.SubscribeEvent("test", func(any) error { return nil })
		return nil
	})

	done := make(chan error, 1)
	go func() { done <- e.PublishEvent("test", nil) }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("PublishEvent deadlocked during reentrant subscription")
	}
}

func TestEventHandlerContainsSubscriberPanic(t *testing.T) {
	e := NewEventHandler()
	e.SubscribeEvent("test", func(any) error {
		panic("subscriber failure")
	})

	err := e.PublishEvent("test", nil)
	var panicErr *CallbackPanicError
	if !errors.As(err, &panicErr) {
		t.Fatalf("PublishEvent() error = %v, want CallbackPanicError", err)
	}
}
