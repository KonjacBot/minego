package client

import (
	"bytes"
	"testing"
)

func TestForgetLevelChunkWireLayout(t *testing.T) {
	tests := []struct {
		name string
		pos  PackedChunkPos
		wire []byte
	}{
		{"正的", PackedChunkPos{1, 2}, []byte{0, 0, 0, 2, 0, 0, 0, 1}},
		{"負的", PackedChunkPos{-213, -218}, []byte{0xff, 0xff, 0xff, 0x26, 0xff, 0xff, 0xff, 0x2b}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var decoded ForgetLevelChunk
			if _, err := decoded.ReadFrom(bytes.NewReader(tc.wire)); err != nil {
				t.Fatal(err)
			}
			if decoded.Pos != tc.pos {
				t.Fatalf("decoded pos = %v, want %v", decoded.Pos, tc.pos)
			}
			var encoded bytes.Buffer
			if _, err := (ForgetLevelChunk{Pos: tc.pos}).WriteTo(&encoded); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(encoded.Bytes(), tc.wire) {
				t.Fatalf("encoded = %v, want %v", encoded.Bytes(), tc.wire)
			}
		})
	}
}
