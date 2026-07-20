package slot

import (
	"bytes"
	"testing"

	"github.com/KonjacBot/go-mc/level/item"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

func TestHashedSlotOfficialEmptyWire(t *testing.T) {
	wire := []byte{0}

	var got HashedSlot
	read, err := got.ReadFrom(bytes.NewReader(wire))
	if err != nil {
		t.Fatal(err)
	}
	if read != int64(len(wire)) {
		t.Fatalf("ReadFrom() = %d, want %d", read, len(wire))
	}
	if got.HasItem {
		t.Fatalf("decoded = %#v, want empty hashed slot", got)
	}

	var encoded bytes.Buffer
	written, err := got.WriteTo(&encoded)
	if err != nil {
		t.Fatal(err)
	}
	if written != int64(len(wire)) || !bytes.Equal(encoded.Bytes(), wire) {
		t.Fatalf("WriteTo() = %x (%d), want %x (%d)", encoded.Bytes(), written, wire, len(wire))
	}
}

func TestHashedSlotOfficialWire(t *testing.T) {
	want := HashedSlot{
		HasItem:           true,
		ItemID:            5,
		ItemCount:         2,
		AddComponents:     []AddedHashedComponent{{Type: 7, DataHash: 0x01020304}},
		RemovedComponents: []int32{9, 10},
	}
	wire := []byte{1, 5, 2, 1, 7, 1, 2, 3, 4, 2, 9, 10}

	var got HashedSlot
	read, err := got.ReadFrom(bytes.NewReader(wire))
	if err != nil {
		t.Fatal(err)
	}
	if read != int64(len(wire)) {
		t.Fatalf("ReadFrom() = %d, want %d", read, len(wire))
	}
	if got.HasItem != want.HasItem || got.ItemID != want.ItemID || got.ItemCount != want.ItemCount {
		t.Fatalf("decoded header = %#v, want %#v", got, want)
	}
	if len(got.AddComponents) != 1 || got.AddComponents[0] != want.AddComponents[0] {
		t.Fatalf("decoded added = %#v, want %#v", got.AddComponents, want.AddComponents)
	}
	if len(got.RemovedComponents) != 2 || got.RemovedComponents[0] != 9 || got.RemovedComponents[1] != 10 {
		t.Fatalf("decoded removed = %#v, want %#v", got.RemovedComponents, want.RemovedComponents)
	}

	var encoded bytes.Buffer
	written, err := want.WriteTo(&encoded)
	if err != nil {
		t.Fatal(err)
	}
	if written != int64(len(wire)) || !bytes.Equal(encoded.Bytes(), wire) {
		t.Fatalf("WriteTo() = %x (%d), want %x (%d)", encoded.Bytes(), written, wire, len(wire))
	}
}

func TestHashedSlotReadClearsReusedValue(t *testing.T) {
	value := HashedSlot{HasItem: true, ItemID: 99, ItemCount: 5, AddComponents: []AddedHashedComponent{{Type: 1, DataHash: 2}}, RemovedComponents: []int32{3}}
	if _, err := value.ReadFrom(bytes.NewReader([]byte{0})); err != nil {
		t.Fatal(err)
	}
	if value.HasItem || value.ItemID != 0 || value.ItemCount != 0 || value.AddComponents != nil || value.RemovedComponents != nil {
		t.Fatalf("reused hashed slot retained stale state: %#v", value)
	}
}

func TestHashedSlotRejectsOfficialCollectionCaps(t *testing.T) {
	t.Run("decode added components over cap", func(t *testing.T) {
		var wire bytes.Buffer
		mustWriteVarInt(t, &wire, 1)
		mustWriteVarInt(t, &wire, 1)
		mustWriteVarInt(t, &wire, 1)
		mustWriteVarInt(t, &wire, 257)
		mustWriteVarInt(t, &wire, 0)

		if _, err := new(HashedSlot).ReadFrom(bytes.NewReader(wire.Bytes())); err == nil {
			t.Fatal("ReadFrom() accepted added component count > 256")
		}
	})

	t.Run("decode removed components over cap", func(t *testing.T) {
		var wire bytes.Buffer
		mustWriteVarInt(t, &wire, 1)
		mustWriteVarInt(t, &wire, 1)
		mustWriteVarInt(t, &wire, 1)
		mustWriteVarInt(t, &wire, 0)
		mustWriteVarInt(t, &wire, 257)

		if _, err := new(HashedSlot).ReadFrom(bytes.NewReader(wire.Bytes())); err == nil {
			t.Fatal("ReadFrom() accepted removed component count > 256")
		}
	})

	t.Run("write added components over cap", func(t *testing.T) {
		value := HashedSlot{HasItem: true, ItemID: 1, ItemCount: 1, AddComponents: make([]AddedHashedComponent, 257)}
		if _, err := value.WriteTo(&bytes.Buffer{}); err == nil {
			t.Fatal("WriteTo() accepted added component count > 256")
		}
	})

	t.Run("write removed components over cap", func(t *testing.T) {
		value := HashedSlot{HasItem: true, ItemID: 1, ItemCount: 1, RemovedComponents: make([]int32, 257)}
		if _, err := value.WriteTo(&bytes.Buffer{}); err == nil {
			t.Fatal("WriteTo() accepted removed component count > 256")
		}
	})
}

func TestItemStackTemplateOfficialWire(t *testing.T) {
	want := ItemStackTemplate{ItemID: item.ID(1), Count: 2}
	wire := []byte{1, 2, 0, 0}

	var got ItemStackTemplate
	read, err := got.ReadFrom(bytes.NewReader(wire))
	if err != nil {
		t.Fatal(err)
	}
	if read != int64(len(wire)) {
		t.Fatalf("ReadFrom() = %d, want %d", read, len(wire))
	}
	if got.ItemID != want.ItemID || got.Count != want.Count || got.AddComponent != nil || got.RemoveComponent != nil {
		t.Fatalf("decoded = %#v, want %#v", got, want)
	}

	var encoded bytes.Buffer
	written, err := want.WriteTo(&encoded)
	if err != nil {
		t.Fatal(err)
	}
	if written != int64(len(wire)) || !bytes.Equal(encoded.Bytes(), wire) {
		t.Fatalf("WriteTo() = %x (%d), want %x (%d)", encoded.Bytes(), written, wire, len(wire))
	}
}

func TestItemStackTemplateRejectsEmptyOfficialWire(t *testing.T) {
	t.Run("read empty item id", func(t *testing.T) {
		wire := []byte{0, 1, 0, 0}
		if _, err := new(ItemStackTemplate).ReadFrom(bytes.NewReader(wire)); err == nil {
			t.Fatal("ReadFrom() accepted empty item stack template item id")
		}
	})

	t.Run("read zero count", func(t *testing.T) {
		wire := []byte{1, 0, 0, 0}
		if _, err := new(ItemStackTemplate).ReadFrom(bytes.NewReader(wire)); err == nil {
			t.Fatal("ReadFrom() accepted zero-count item stack template")
		}
	})

	t.Run("write empty item id", func(t *testing.T) {
		value := ItemStackTemplate{ItemID: 0, Count: 1}
		if _, err := value.WriteTo(&bytes.Buffer{}); err == nil {
			t.Fatal("WriteTo() accepted empty item stack template item id")
		}
	})

	t.Run("write zero count", func(t *testing.T) {
		value := ItemStackTemplate{ItemID: 1, Count: 0}
		if _, err := value.WriteTo(&bytes.Buffer{}); err == nil {
			t.Fatal("WriteTo() accepted zero-count item stack template")
		}
	})
}

func TestItemStackTemplateReadClearsReusedValue(t *testing.T) {
	value := ItemStackTemplate{ItemID: item.ID(99), Count: 5, AddComponent: map[int32]Component{1: nil}, RemoveComponent: []int32{2}}
	if _, err := value.ReadFrom(bytes.NewReader([]byte{1, 2, 0, 0})); err != nil {
		t.Fatal(err)
	}
	if value.ItemID != 1 || value.Count != 2 || value.AddComponent != nil || value.RemoveComponent != nil {
		t.Fatalf("reused item stack template retained stale state: %#v", value)
	}
}

func mustWriteVarInt(t *testing.T, w *bytes.Buffer, value int32) {
	t.Helper()
	if _, err := pk.VarInt(value).WriteTo(w); err != nil {
		t.Fatal(err)
	}
}
