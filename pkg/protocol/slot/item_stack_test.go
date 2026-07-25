package slot

import (
	"bytes"
	"io"
	"testing"

	"github.com/KonjacBot/go-mc/level/item"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

type testPatchComponent struct{}

func (testPatchComponent) ID() string                        { return "test" }
func (testPatchComponent) WriteTo(io.Writer) (int64, error)  { return 0, nil }
func (testPatchComponent) ReadFrom(io.Reader) (int64, error) { return 0, nil }

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

func TestComponentPatchRejectsInvalidRemovedIDs(t *testing.T) {
	tests := []struct {
		name    string
		added   map[int32]Component
		removed []int32
	}{
		{name: "duplicate", removed: []int32{0, 0}},
		{name: "unknown", removed: []int32{-1}},
		{name: "added and removed", added: map[int32]Component{0: testPatchComponent{}}, removed: []int32{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateComponentPatch(tt.added, tt.removed); err == nil {
				t.Fatal("validateComponentPatch() accepted an invalid patch")
			}
		})
	}
}

func TestComponentPatchCodecsRejectDuplicateRemovedIDs(t *testing.T) {
	wire := []byte{1, 1, 0, 2, 0, 0}
	t.Run("slot read", func(t *testing.T) {
		if _, err := new(Slot).ReadFrom(bytes.NewReader(wire)); err == nil {
			t.Fatal("ReadFrom() accepted duplicate removed component ids")
		}
	})
	t.Run("template read", func(t *testing.T) {
		if _, err := new(ItemStackTemplate).ReadFrom(bytes.NewReader(wire)); err == nil {
			t.Fatal("ReadFrom() accepted duplicate removed component ids")
		}
	})
	t.Run("slot write", func(t *testing.T) {
		value := Slot{Count: 1, ItemID: item.ID(1), RemoveComponent: []int32{0, 0}}
		if _, err := value.WriteTo(&bytes.Buffer{}); err == nil {
			t.Fatal("WriteTo() accepted duplicate removed component ids")
		}
	})
	t.Run("template write", func(t *testing.T) {
		value := ItemStackTemplate{Count: 1, ItemID: item.ID(1), RemoveComponent: []int32{0, 0}}
		if _, err := value.WriteTo(&bytes.Buffer{}); err == nil {
			t.Fatal("WriteTo() accepted duplicate removed component ids")
		}
	})
}

func TestGeneratedSlotArrayRejectsOversizedWrite(t *testing.T) {
	value := make(Int32VarIntVarIntArray, maxComponentPatchEntries+1)
	if _, err := value.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted an oversized array")
	}
}
