package client

import (
	"bytes"
	"testing"

	pk "github.com/KonjacBot/go-mc/net/packet"
)

func TestLevelChunkWithLightReadsBoundedEmptyPayload(t *testing.T) {
	data := make([]byte, 8) // chunk position
	data = append(data,
		0,          // height maps
		0,          // section data
		0,          // block entities
		0, 0, 0, 0, // light masks
		0, 0, // light arrays
	)
	var value LevelChunkWithLight
	read, err := value.ReadFrom(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if read != int64(len(data)) {
		t.Fatalf("ReadFrom() count = %d, want %d", read, len(data))
	}
}

func TestLevelChunkWithLightRejectsTruncatedHugeSectionData(t *testing.T) {
	var encodedLength [pk.MaxVarIntLen]byte
	length := pk.VarInt(pk.MaxDataLength).WriteToBytes(encodedLength[:])
	data := make([]byte, 8)
	data = append(data, 0) // height maps
	data = append(data, encodedLength[:length]...)

	var value LevelChunkWithLight
	if _, err := value.ReadFrom(bytes.NewReader(data)); err == nil {
		t.Fatal("ReadFrom() accepted truncated section data")
	}
}

func TestLevelChunkWithLightRejectsOversizedLightMask(t *testing.T) {
	data := make([]byte, 8)
	data = append(data,
		0, // height maps
		0, // section data
		0, // block entities
		2, // first mask has two longs; an empty chunk permits one
	)
	var value LevelChunkWithLight
	if _, err := value.ReadFrom(bytes.NewReader(data)); err == nil {
		t.Fatal("ReadFrom() accepted an oversized light mask")
	}
}
