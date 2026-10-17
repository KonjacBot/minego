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

// setEquipmentWire 是實機格式的 SetEquipment 封包內容:
// EntityID=300(VarInt),主手(slot 0,帶 continuation bit)放 1 個 ID 5 的物品，
// 頭盔(slot 5,最後一筆)放 1 個 ID 1112 的物品，物品皆無 component 變更。
var setEquipmentWire = []byte{
	0xac, 0x02, // EntityID 300
	0x80, 0x01, 0x05, 0x00, 0x00, // slot 0|0x80, count 1, item 5, +0 -0 components
	0x05, 0x01, 0xd8, 0x08, 0x00, 0x00, // slot 5, count 1, item 1112, +0 -0 components
}

func TestSetEquipmentReadFromWire(t *testing.T) {
	var p SetEquipment
	r := bytes.NewReader(setEquipmentWire)
	n, err := p.ReadFrom(r)
	if err != nil {
		t.Fatalf("ReadFrom(): %v", err)
	}
	if n != int64(len(setEquipmentWire)) || r.Len() != 0 {
		t.Fatalf("read %d bytes, %d left; want all %d consumed", n, r.Len(), len(setEquipmentWire))
	}
	if p.EntityID != 300 {
		t.Fatalf("EntityID = %d, want 300", p.EntityID)
	}
	if len(p.Equipment) != 2 {
		t.Fatalf("len(Equipment) = %d, want 2: %#v", len(p.Equipment), p.Equipment)
	}
	want := []struct {
		slot  int8
		item  item.ID
		count int32
	}{{0, 5, 1}, {5, 1112, 1}}
	for i, w := range want {
		got := p.Equipment[i]
		if got.Slot != w.slot || got.Item.ItemID != w.item || got.Item.Count != w.count {
			t.Fatalf("Equipment[%d] = {Slot:%d ItemID:%d Count:%d}, want {%d %d %d}",
				i, got.Slot, got.Item.ItemID, got.Item.Count, w.slot, w.item, w.count)
		}
	}

	var encoded bytes.Buffer
	if _, err := p.WriteTo(&encoded); err != nil {
		t.Fatalf("WriteTo(): %v", err)
	}
	if !bytes.Equal(encoded.Bytes(), setEquipmentWire) {
		t.Fatalf("re-encoded = % x, want % x", encoded.Bytes(), setEquipmentWire)
	}
}
