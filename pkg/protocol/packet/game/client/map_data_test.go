package client

import (
	"bytes"
	"testing"

	pk "github.com/KonjacBot/go-mc/net/packet"
)

func TestMapColorPatchRoundTrip(t *testing.T) {
	want := MapColorPatch{
		Columns: 1,
		Rows:    2,
		X:       3,
		Z:       4,
		Data:    []pk.UnsignedByte{5, 6},
	}
	var encoded bytes.Buffer
	written, err := want.WriteTo(&encoded)
	if err != nil {
		t.Fatal(err)
	}
	if written != int64(encoded.Len()) {
		t.Fatalf("WriteTo() count = %d, encoded length = %d", written, encoded.Len())
	}

	var got MapColorPatch
	read, err := got.ReadFrom(bytes.NewReader(encoded.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	if read != int64(encoded.Len()) {
		t.Fatalf("ReadFrom() count = %d, encoded length = %d", read, encoded.Len())
	}
	if got.Columns != want.Columns || got.Rows != want.Rows || got.X != want.X || got.Z != want.Z || len(got.Data) != 2 {
		t.Fatalf("decoded patch = %#v, want %#v", got, want)
	}
}

func TestMapColorPatchReturnsTruncatedInputError(t *testing.T) {
	var patch MapColorPatch
	if _, err := patch.ReadFrom(bytes.NewReader([]byte{1, 2})); err == nil {
		t.Fatal("ReadFrom() accepted truncated patch")
	}
}
