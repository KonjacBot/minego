package slot

import (
	"bytes"
	"errors"
	"io"
	"testing"

	pk "github.com/KonjacBot/go-mc/net/packet"
)

func TestTradeSlotRejectsUnknownComponent(t *testing.T) {
	data := encodeVarInts(t, 1, 1, 1, 1<<20)
	var value TradeSlot
	if _, err := value.ReadFrom(bytes.NewReader(data)); err == nil {
		t.Fatal("ReadFrom() accepted an unknown component")
	}
}

func TestTradeSlotRejectsOversizedComponentCount(t *testing.T) {
	data := encodeVarInts(t, 1, 1, maxComponentPatchEntries+1)
	var value TradeSlot
	if _, err := value.ReadFrom(bytes.NewReader(data)); err == nil {
		t.Fatal("ReadFrom() accepted an oversized component count")
	}
}

func TestTradeSlotPropagatesTruncatedInput(t *testing.T) {
	data := encodeVarInts(t, 1, 1, 1)
	var value TradeSlot
	if _, err := value.ReadFrom(bytes.NewReader(data)); !errors.Is(err, io.EOF) {
		t.Fatalf("ReadFrom() error = %v, want io.EOF", err)
	}
}

func TestTradeSlotPropagatesWriteError(t *testing.T) {
	value := TradeSlot{ID: 1, Count: 1}
	if _, err := value.WriteTo(errorWriter{}); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("WriteTo() error = %v, want io.ErrClosedPipe", err)
	}
}

func encodeVarInts(t *testing.T, values ...int) []byte {
	t.Helper()
	var data bytes.Buffer
	for _, value := range values {
		if _, err := pk.VarInt(value).WriteTo(&data); err != nil {
			t.Fatal(err)
		}
	}
	return data.Bytes()
}

type errorWriter struct{}

func (errorWriter) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }
