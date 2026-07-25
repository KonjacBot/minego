package component

import (
	"bytes"
	"testing"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

func TestTrimRegistryReferencesDoNotConsumeNextComponent(t *testing.T) {
	wire := []byte{6, 13, 0}
	reader := bytes.NewReader(wire)
	var trim Trim

	n, err := trim.ReadFrom(reader)
	if err != nil {
		t.Fatalf("ReadFrom() error = %v", err)
	}
	if n != 2 {
		t.Fatalf("ReadFrom() consumed %d bytes, want 2", n)
	}
	if !trim.TrimMaterial.Has || trim.TrimMaterial.ID != 5 {
		t.Fatalf("unexpected trim material holder: %+v", trim.TrimMaterial)
	}
	if !trim.TrimPattern.Has || trim.TrimPattern.ID != 12 {
		t.Fatalf("unexpected trim pattern holder: %+v", trim.TrimPattern)
	}

	nextComponentID, err := reader.ReadByte()
	if err != nil {
		t.Fatalf("read next component id: %v", err)
	}
	if nextComponentID != 0 {
		t.Fatalf("next component id = %d, want 0", nextComponentID)
	}

	var encoded bytes.Buffer
	n, err = trim.WriteTo(&encoded)
	if err != nil {
		t.Fatalf("WriteTo() error = %v", err)
	}
	if n != 2 {
		t.Fatalf("WriteTo() wrote %d bytes, want 2", n)
	}
	if want := []byte{6, 13}; !bytes.Equal(encoded.Bytes(), want) {
		t.Fatalf("WriteTo() = %v, want %v", encoded.Bytes(), want)
	}
}

func TestSlotTrimRegistryReferencesPreserveFollowingComponent(t *testing.T) {
	wire := []byte{
		1, 63, 2, 0,
		56, 6, 13,
		0, 10, 0,
	}
	var got slot.Slot

	n, err := got.ReadFrom(bytes.NewReader(wire))
	if err != nil {
		t.Fatalf("ReadFrom() error = %v", err)
	}
	if n != int64(len(wire)) {
		t.Fatalf("ReadFrom() consumed %d bytes, want %d", n, len(wire))
	}

	trim, ok := got.AddComponent[56].(*Trim)
	if !ok {
		t.Fatalf("component 56 = %T, want *component.Trim", got.AddComponent[56])
	}
	if !trim.TrimMaterial.Has || trim.TrimMaterial.ID != 5 {
		t.Fatalf("unexpected trim material holder: %+v", trim.TrimMaterial)
	}
	if !trim.TrimPattern.Has || trim.TrimPattern.ID != 12 {
		t.Fatalf("unexpected trim pattern holder: %+v", trim.TrimPattern)
	}
	if _, ok := got.AddComponent[0].(*CustomData); !ok {
		t.Fatalf("component 0 = %T, want *component.CustomData", got.AddComponent[0])
	}
}
