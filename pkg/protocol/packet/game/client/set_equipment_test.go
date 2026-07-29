package client

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/KonjacBot/go-mc/level/item"
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

func TestEquipmentRoundTripMultipleEntries(t *testing.T) {
	want := Equipment{
		{Slot: 0, Item: slot.Slot{ItemID: item.Stone{}.ID(), Count: 1}},
		{Slot: 5, Item: slot.Slot{}},
	}
	var encoded bytes.Buffer
	if _, err := want.WriteTo(&encoded); err != nil {
		t.Fatalf("WriteTo(): %v", err)
	}
	var got Equipment
	if _, err := got.ReadFrom(&encoded); err != nil {
		t.Fatalf("ReadFrom(): %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("decoded equipment count = %d, want %d", len(got), len(want))
	}
	for index := range want {
		if got[index].Slot != want[index].Slot ||
			got[index].Item.ItemID != want[index].Item.ItemID ||
			got[index].Item.Count != want[index].Item.Count {
			t.Fatalf("decoded equipment[%d] = %#v, want %#v", index, got[index], want[index])
		}
	}
}

func TestEquipmentContinuationBit(t *testing.T) {
	equipment := Equipment{
		{Slot: 0, Item: slot.Slot{}},
		{Slot: 5, Item: slot.Slot{}},
	}

	var encoded bytes.Buffer
	if _, err := equipment.WriteTo(&encoded); err != nil {
		t.Fatalf("WriteTo(): %v", err)
	}

	// The high bit is set on entries that have a following entry, not on
	// the final entry. Each empty item is encoded as a zero count byte.
	wantWire := []byte{0x80, 0x00, 0x05, 0x00}
	if got := encoded.Bytes(); !bytes.Equal(got, wantWire) {
		t.Fatalf("encoded equipment = % x, want % x", got, wantWire)
	}

	var decoded Equipment
	if _, err := decoded.ReadFrom(bytes.NewReader(wantWire)); err != nil {
		t.Fatalf("ReadFrom(): %v", err)
	}
	if !reflect.DeepEqual(decoded, equipment) {
		t.Fatalf("decoded equipment = %#v, want %#v", decoded, equipment)
	}
}
