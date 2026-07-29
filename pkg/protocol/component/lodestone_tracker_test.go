package component

import (
	"bytes"
	"reflect"
	"testing"

	pk "github.com/KonjacBot/go-mc/net/packet"
)

func TestLodestoneTrackerWireLayout(t *testing.T) {
	tests := []struct {
		name string
		wire []byte
		want LodestoneTracker
	}{
		{
			name: "with target",
			wire: []byte{
				0x01,
				0x14, 'm', 'i', 'n', 'e', 'c', 'r', 'a', 'f', 't', ':',
				'd', 'x', 'l', '_', 'e', 'd', 'i', 't', '_', '2',
				0x00, 0x00, 0x03, 0x80, 0x00, 0x0b, 0xd0, 0x44,
				0x01,
			},
			want: LodestoneTracker{
				HasGlobalPosition: true,
				Dimension: pk.Option[pk.Identifier, *pk.Identifier]{
					Has: true,
					Val: "minecraft:dxl_edit_2",
				},
				Position: pk.Option[pk.Position, *pk.Position]{
					Has: true,
					Val: pk.Position{X: 14, Y: 68, Z: 189},
				},
				Tracked: true,
			},
		},
		{
			name: "without target",
			wire: []byte{0x00, 0x01},
			want: LodestoneTracker{Tracked: true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var payload bytes.Buffer
			payload.Write(test.wire)
			if _, err := pk.VarInt(11).WriteTo(&payload); err != nil {
				t.Fatal(err)
			}

			var got LodestoneTracker
			read, err := got.ReadFrom(&payload)
			if err != nil {
				t.Fatal(err)
			}
			if read != int64(len(test.wire)) {
				t.Fatalf("decoder consumed %d bytes, want %d", read, len(test.wire))
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("decoded tracker = %#v, want %#v", got, test.want)
			}

			var nextComponentID int32
			if _, err := (*pk.VarInt)(&nextComponentID).ReadFrom(&payload); err != nil {
				t.Fatal(err)
			}
			if nextComponentID != 11 || payload.Len() != 0 {
				t.Fatalf("next component id = %d, remaining bytes = %d", nextComponentID, payload.Len())
			}

			var encoded bytes.Buffer
			if _, err := got.WriteTo(&encoded); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(encoded.Bytes(), test.wire) {
				t.Fatalf("encoded tracker = %x, want %x", encoded.Bytes(), test.wire)
			}
		})
	}
}
