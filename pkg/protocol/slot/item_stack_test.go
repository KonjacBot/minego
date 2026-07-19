package slot

import (
	"bytes"
	"testing"

	"github.com/KonjacBot/go-mc/level/item"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

func TestSlotRoundTripReportsCompleteSize(t *testing.T) {
	want := Slot{Count: 2, ItemID: item.ID(3), RemoveComponent: []int32{4, 5}}
	var encoded bytes.Buffer
	written, err := want.WriteTo(&encoded)
	if err != nil {
		t.Fatal(err)
	}
	if written != int64(encoded.Len()) {
		t.Fatalf("WriteTo() count = %d, encoded length = %d", written, encoded.Len())
	}

	var got Slot
	read, err := got.ReadFrom(bytes.NewReader(encoded.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	if read != int64(encoded.Len()) {
		t.Fatalf("ReadFrom() count = %d, encoded length = %d", read, encoded.Len())
	}
	if got.Count != want.Count || got.ItemID != want.ItemID || len(got.RemoveComponent) != 2 {
		t.Fatalf("decoded slot = %#v, want %#v", got, want)
	}
}

func TestSlotRejectsUnknownComponent(t *testing.T) {
	var encoded bytes.Buffer
	fields := []pk.VarInt{
		pk.VarInt(1),
		pk.VarInt(1),
		pk.VarInt(1),
		pk.VarInt(0),
		pk.VarInt(1 << 20),
	}
	for i := range fields {
		if _, err := fields[i].WriteTo(&encoded); err != nil {
			t.Fatal(err)
		}
	}

	var got Slot
	if _, err := got.ReadFrom(bytes.NewReader(encoded.Bytes())); err == nil {
		t.Fatal("ReadFrom() accepted an unknown component")
	}
}

func TestEmptySlotClearsReusedValue(t *testing.T) {
	got := Slot{Count: 1, ItemID: item.ID(9), RemoveComponent: []int32{1}}
	if _, err := got.ReadFrom(bytes.NewReader([]byte{0})); err != nil {
		t.Fatal(err)
	}
	if got.Count != 0 || got.ItemID != 0 || got.RemoveComponent != nil {
		t.Fatalf("reused slot retained stale state: %#v", got)
	}
}
