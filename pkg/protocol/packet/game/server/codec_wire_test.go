package server

import (
	"bytes"
	"testing"

	"github.com/KonjacBot/go-mc/chat/sign"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

func TestChangeDifficultyUsesVarInt(t *testing.T) {
	var encoded bytes.Buffer
	if _, err := (ChangeDifficulty{Difficulty: 128}).WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	if got := encoded.Bytes(); !bytes.Equal(got, []byte{0x80, 0x01}) {
		t.Fatalf("encoded difficulty = %v, want VarInt", got)
	}
}

func TestChatSignatureHasNoLengthPrefix(t *testing.T) {
	var signature sign.Signature
	signature[0] = 0xaa
	packet := Chat{
		HasSignature: true,
		Signature:    signature,
		Acknowledged: pk.NewFixedBitSet(20),
	}
	var encoded bytes.Buffer
	written, err := packet.WriteTo(&encoded)
	if err != nil {
		t.Fatal(err)
	}
	if written != 279 || encoded.Len() != 279 {
		t.Fatalf("encoded chat size = (%d, %d), want 279", written, encoded.Len())
	}
	if got := encoded.Bytes()[18]; got != 0xaa {
		t.Fatalf("first signature byte = %x, want aa", got)
	}
}

func TestSignedArgumentSignatureHasNoLengthPrefix(t *testing.T) {
	var signature sign.Signature
	signature[0] = 0xbb
	var encoded bytes.Buffer
	written, err := (SignedSignatures{ArgumentName: "x", Signature: signature}).WriteTo(&encoded)
	if err != nil {
		t.Fatal(err)
	}
	if written != 258 || encoded.Len() != 258 {
		t.Fatalf("encoded argument signature size = (%d, %d), want 258", written, encoded.Len())
	}
	if got := encoded.Bytes()[2]; got != 0xbb {
		t.Fatalf("first signature byte = %x, want bb", got)
	}
}

func TestCustomPayloadReadsOnlyAvailableData(t *testing.T) {
	wire := []byte{1, 'x', 1, 2, 3}
	var payload CustomPayload
	read, err := payload.ReadFrom(bytes.NewReader(wire))
	if err != nil {
		t.Fatal(err)
	}
	if read != int64(len(wire)) || payload.Channel != "x" || !bytes.Equal(payload.Data, []byte{1, 2, 3}) {
		t.Fatalf("decoded payload = %#v after %d bytes", payload, read)
	}
}

func TestCustomPayloadRejectsOversizeData(t *testing.T) {
	packet := CustomPayload{Channel: "x", Data: make([]byte, 32768)}
	if _, err := packet.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted an oversized custom payload")
	}

	wire := append([]byte{1, 'x'}, make([]byte, 32768)...)
	if _, err := new(CustomPayload).ReadFrom(bytes.NewReader(wire)); err == nil {
		t.Fatal("ReadFrom() accepted an oversized custom payload")
	}
}
