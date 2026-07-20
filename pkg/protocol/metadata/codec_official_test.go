package metadata

import (
	"bytes"
	"io"
	"testing"

	"github.com/KonjacBot/go-mc/chat"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

func TestOptChatUsesOptionalMetadataType(t *testing.T) {
	value := OptChat{}
	if value.EntityMetadataType() != MetadataOptChat {
		t.Fatalf("EntityMetadataType() = %d, want %d", value.EntityMetadataType(), MetadataOptChat)
	}
}

func TestOptBlockStateOfficialWire(t *testing.T) {
	if wrote, got := writeMetadataValue(t, OptBlockState{}), []byte{0}; !bytes.Equal(wrote, got) {
		t.Fatalf("empty WriteTo() = %x, want %x", wrote, got)
	}

	value := OptBlockState{Value: 42}
	if wrote, got := writeMetadataValue(t, value), []byte{42}; !bytes.Equal(wrote, got) {
		t.Fatalf("present WriteTo() = %x, want %x", wrote, got)
	}

	var decoded OptBlockState
	if _, err := decoded.ReadFrom(bytes.NewReader([]byte{42})); err != nil {
		t.Fatal(err)
	}
	if decoded.Value != 42 {
		t.Fatalf("decoded = %#v, want value 42", decoded)
	}

	decoded = OptBlockState{Value: 99}
	if _, err := decoded.ReadFrom(bytes.NewReader([]byte{0})); err != nil {
		t.Fatal(err)
	}
	if decoded.Value != 0 {
		t.Fatalf("reused opt block state retained stale value: %#v", decoded)
	}
}

func TestOptVarIntOfficialWire(t *testing.T) {
	if wrote, got := writeMetadataValue(t, OptVarInt{}), []byte{0}; !bytes.Equal(wrote, got) {
		t.Fatalf("empty WriteTo() = %x, want %x", wrote, got)
	}

	value := OptVarInt{Has: true, Value: 0}
	if wrote, got := writeMetadataValue(t, value), []byte{1}; !bytes.Equal(wrote, got) {
		t.Fatalf("zero WriteTo() = %x, want %x", wrote, got)
	}

	value = OptVarInt{Has: true, Value: 41}
	if wrote, got := writeMetadataValue(t, value), []byte{42}; !bytes.Equal(wrote, got) {
		t.Fatalf("present WriteTo() = %x, want %x", wrote, got)
	}

	var decoded OptVarInt
	if _, err := decoded.ReadFrom(bytes.NewReader([]byte{42})); err != nil {
		t.Fatal(err)
	}
	if !decoded.Has || decoded.Value != 41 {
		t.Fatalf("decoded = %#v, want has=true value=41", decoded)
	}

	decoded = OptVarInt{Has: true, Value: 99}
	if _, err := decoded.ReadFrom(bytes.NewReader([]byte{0})); err != nil {
		t.Fatal(err)
	}
	if decoded.Has || decoded.Value != 0 {
		t.Fatalf("reused opt varint retained stale value: %#v", decoded)
	}
}

func TestEntityMetadataWritesOptionalChatTypeID(t *testing.T) {
	meta := EntityMetadata{Data: map[uint8]Metadata{
		1: &OptChat{Option: pk.Option[chat.Message, *chat.Message]{Has: true, Val: chat.Text("hello")}},
	}}
	encoded := writeMetadataValue(t, meta)
	if len(encoded) < 2 || encoded[0] != 1 || encoded[1] != byte(MetadataOptChat) {
		t.Fatalf("metadata header = %x, want index 1 type %d", encoded, MetadataOptChat)
	}
}

func writeMetadataValue[T interface {
	WriteTo(io.Writer) (int64, error)
}](t *testing.T, value T) []byte {
	t.Helper()
	var encoded bytes.Buffer
	if _, err := value.WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	return encoded.Bytes()
}
