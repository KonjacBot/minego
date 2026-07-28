package client

import (
	"bytes"
	"testing"
)

func TestLevelChunksLoadStartGameEventWireLayout(t *testing.T) {
	event := GameEvent{
		Event: GameEventLevelChunksLoadStart,
		Param: 0,
	}
	var encoded bytes.Buffer
	if _, err := event.WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	want := []byte{13, 0, 0, 0, 0}
	if !bytes.Equal(encoded.Bytes(), want) {
		t.Fatalf("encoded event = %v, want %v", encoded.Bytes(), want)
	}

	var decoded GameEvent
	if _, err := decoded.ReadFrom(bytes.NewReader(want)); err != nil {
		t.Fatal(err)
	}
	if decoded != event {
		t.Fatalf("decoded event = %+v, want %+v", decoded, event)
	}
}
