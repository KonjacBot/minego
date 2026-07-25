package client

import (
	"bytes"
	"testing"
)

func FuzzLoginPacketDecoders(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0xff, 0xff, 0xff, 0xff, 0x0f})
	f.Add(bytes.Repeat([]byte{0xff}, 64))

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 4096 {
			t.Skip()
		}
		for id, creator := range ClientboundPackets {
			func() {
				defer func() {
					if recovered := recover(); recovered != nil {
						t.Fatalf("login packet %d decoder panicked: %v", id, recovered)
					}
				}()
				_, _ = creator().ReadFrom(bytes.NewReader(data))
			}()
		}
	})
}
