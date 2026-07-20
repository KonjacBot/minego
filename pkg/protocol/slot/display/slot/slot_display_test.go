package slot

import (
	"bytes"
	"io"
	"testing"
)

func TestDisplayReadItemUsesCurrentRegistryID(t *testing.T) {
	var display Display
	n, err := display.ReadFrom(bytes.NewReader([]byte{4, 5}))
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("bytes read = %d, want 2", n)
	}

	item, ok := display.SlotDisplay.(*Item)
	if !ok {
		t.Fatalf("display id 4 decoded as %T, want *Item", display.SlotDisplay)
	}
	if item.ID != 5 {
		t.Fatalf("item id = %d, want 5", item.ID)
	}
}

func TestDisplayEmptyClearsReusedState(t *testing.T) {
	display := Display{SlotDisplay: &Item{ID: 5}}
	if n, err := display.ReadFrom(bytes.NewReader([]byte{0})); err != nil || n != 1 {
		t.Fatalf("ReadFrom() = (%d, %v), want (1, nil)", n, err)
	}
	if _, ok := display.SlotDisplay.(Empty); !ok {
		t.Fatalf("empty display decoded as %T", display.SlotDisplay)
	}
}

func TestDisplayWriteReportsAllBytes(t *testing.T) {
	var buf bytes.Buffer
	n, err := (Display{SlotDisplay: &Item{ID: 5}}).WriteTo(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if n != int64(buf.Len()) || !bytes.Equal(buf.Bytes(), []byte{4, 5}) {
		t.Fatalf("WriteTo() = (%d, %v), bytes %v", n, err, buf.Bytes())
	}
}

func TestDisplayRejectsUnknownType(t *testing.T) {
	var display Display
	if _, err := display.ReadFrom(bytes.NewReader([]byte{99})); err == nil {
		t.Fatal("ReadFrom() accepted unknown display type")
	}
}

var _ io.ReaderFrom = (*Display)(nil)
