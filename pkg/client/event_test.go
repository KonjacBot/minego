package client

import (
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
