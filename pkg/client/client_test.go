package client

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	mcnet "github.com/KonjacBot/go-mc/net"
)

func TestNewClientAllowsNilOptions(t *testing.T) {
	c := NewClient(nil)
	if c == nil {
		t.Fatal("NewClient(nil) returned nil")
	}
}

func TestCloseBeforeConnect(t *testing.T) {
	c := NewClient(nil)
	if err := c.Close(context.Background()); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestConfigurationReturnsReadError(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	b := &botClient{conn: mcnet.WrapConn(clientConn)}
	_ = serverConn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := b.configuration(ctx); err == nil {
		t.Fatal("configuration() returned nil after peer closed")
	}
}

func TestHandleGameStopsPromptlyWhenContextIsCanceled(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer serverConn.Close()
	b := &botClient{
		conn:          mcnet.WrapConn(clientConn),
		packetHandler: newPacketHandler(),
		eventHandler:  NewEventHandler(),
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- b.HandleGame(ctx) }()
	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("HandleGame() error = %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("HandleGame did not stop after context cancellation")
	}
}
