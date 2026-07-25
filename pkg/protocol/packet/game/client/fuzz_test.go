package client

import (
	"bytes"
	"testing"

	"github.com/KonjacBot/go-mc/data/packetid"
)

func FuzzClientboundPacketDecoders(f *testing.F) {
	seeds := [][]byte{
		{},
		{0xff, 0xff, 0xff, 0xff, 0x0f},
		bytes.Repeat([]byte{0xff}, 64),
		make([]byte, 256),
	}
	for id := range clientboundPackets {
		for _, data := range seeds {
			f.Add(uint8(id), data)
		}
	}

	f.Fuzz(func(t *testing.T, packetID uint8, data []byte) {
		if len(data) > 4096 {
			t.Skip()
		}
		id := packetid.ClientboundPacketID(packetID)
		creator, ok := clientboundPackets[id]
		if !ok {
			t.Skip()
		}
		defer func() {
			if recovered := recover(); recovered != nil {
				t.Fatalf("packet %d decoder panicked: %v", id, recovered)
			}
		}()
		packet := creator()
		_, _ = packet.ReadFrom(bytes.NewReader(data))
	})
}
