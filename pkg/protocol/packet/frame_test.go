package packet

import (
	"bytes"
	"compress/zlib"
	"testing"

	pk "github.com/KonjacBot/go-mc/net/packet"
)

func TestFrameRoundTrip(t *testing.T) {
	for _, threshold := range []int{-1, 0, 64} {
		t.Run(string(rune(threshold+2)), func(t *testing.T) {
			want := pk.Packet{ID: 42, Data: bytes.Repeat([]byte{0xab}, 128)}
			var encoded bytes.Buffer
			if err := WriteFrame(&encoded, threshold, want); err != nil {
				t.Fatal(err)
			}
			var got pk.Packet
			if err := ReadFrame(&encoded, threshold, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != want.ID || !bytes.Equal(got.Data, want.Data) {
				t.Fatalf("decoded packet = %#v, want %#v", got, want)
			}
		})
	}
}

func TestReadFrameRejectsExtraDecompressedData(t *testing.T) {
	frame := compressedTestFrame(t, 1, []byte{0, 1}, nil)
	var packet pk.Packet
	if err := ReadFrame(bytes.NewReader(frame), 1, &packet); err == nil {
		t.Fatal("ReadFrame() accepted extra decompressed data")
	}
}

func TestReadFrameRejectsTrailingCompressedFrameData(t *testing.T) {
	frame := compressedTestFrame(t, 1, []byte{0}, []byte{0xaa})
	var packet pk.Packet
	if err := ReadFrame(bytes.NewReader(frame), 1, &packet); err == nil {
		t.Fatal("ReadFrame() accepted trailing compressed-frame data")
	}
}

func TestReadFrameRejectsUncompressedPacketAtThreshold(t *testing.T) {
	frame := []byte{4, 0, 0, 1, 2}
	var packet pk.Packet
	if err := ReadFrame(bytes.NewReader(frame), 2, &packet); err == nil {
		t.Fatal("ReadFrame() accepted an uncompressed packet at the compression threshold")
	}
}

func compressedTestFrame(t *testing.T, declaredLength int, raw, trailing []byte) []byte {
	t.Helper()
	var compressed bytes.Buffer
	zw := zlib.NewWriter(&compressed)
	if _, err := zw.Write(raw); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}

	var body bytes.Buffer
	if _, err := pk.VarInt(declaredLength).WriteTo(&body); err != nil {
		t.Fatal(err)
	}
	_, _ = body.Write(compressed.Bytes())
	_, _ = body.Write(trailing)
	var frame bytes.Buffer
	if _, err := pk.VarInt(body.Len()).WriteTo(&frame); err != nil {
		t.Fatal(err)
	}
	_, _ = body.WriteTo(&frame)
	return frame.Bytes()
}
