package server

import (
	"bytes"
	"strings"
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

func TestCustomClickActionUsesLengthPrefixedOptionalNBT(t *testing.T) {
	var encoded bytes.Buffer
	if _, err := (CustomClickAction{ID: "x"}).WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	if got := encoded.Bytes(); !bytes.Equal(got, []byte{1, 'x', 1, 0}) {
		t.Fatalf("encoded custom click action = %v, want identifier plus length-prefixed empty NBT", got)
	}
}

func TestCustomClickActionRejectsOversizedPayload(t *testing.T) {
	var wire bytes.Buffer
	if _, err := pk.Identifier("x").WriteTo(&wire); err != nil {
		t.Fatal(err)
	}
	if _, err := pk.VarInt(65537).WriteTo(&wire); err != nil {
		t.Fatal(err)
	}
	if _, err := new(CustomClickAction).ReadFrom(bytes.NewReader(wire.Bytes())); err == nil {
		t.Fatal("ReadFrom() accepted an oversized custom click payload")
	}
}

func TestCustomClickActionRejectsTrailingGarbageInsideLengthPrefix(t *testing.T) {
	var wire bytes.Buffer
	if _, err := pk.Identifier("x").WriteTo(&wire); err != nil {
		t.Fatal(err)
	}
	if _, err := pk.VarInt(2).WriteTo(&wire); err != nil {
		t.Fatal(err)
	}
	if _, err := wire.Write([]byte{0, 1}); err != nil {
		t.Fatal(err)
	}
	if _, err := new(CustomClickAction).ReadFrom(bytes.NewReader(wire.Bytes())); err == nil {
		t.Fatal("ReadFrom() accepted trailing garbage inside a length-prefixed custom click payload")
	}
}

func TestChatRejectsTooLongMessage(t *testing.T) {
	if _, err := (Chat{Message: strings.Repeat("a", 257), Acknowledged: pk.NewFixedBitSet(20)}).WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted a chat message longer than 256 characters")
	}
}

func TestChatCommandSignedRejectsTooManyArgumentSignatures(t *testing.T) {
	packet := ChatCommandSigned{ArgumentSignatures: make([]SignedSignatures, 9), Acknowledged: pk.NewFixedBitSet(20)}
	if _, err := packet.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted more than 8 signed arguments")
	}
}

func TestSignedSignaturesRejectsLongArgumentName(t *testing.T) {
	if _, err := (SignedSignatures{ArgumentName: strings.Repeat("a", 17)}).WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted an argument name longer than 16 characters")
	}
}

func TestContainerClickRejectsTooManyChangedSlots(t *testing.T) {
	packet := ContainerClick{ChangedSlots: make([]ChangedSlot, 129)}
	if _, err := packet.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted more than 128 changed slots")
	}
}

func TestEditBookRejectsLongTitle(t *testing.T) {
	packet := EditBook{HasTitle: true, Title: strings.Repeat("a", 33)}
	if _, err := packet.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted a title longer than 32 characters")
	}
}

func TestClientInformationRejectsLongLanguage(t *testing.T) {
	packet := ClientInformation{Location: strings.Repeat("a", 17)}
	if _, err := packet.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted a language longer than 16 characters")
	}
}

func TestCommandSuggestionRejectsLongText(t *testing.T) {
	packet := CommandSuggestion{Text: strings.Repeat("a", 32501)}
	if _, err := packet.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted a command suggestion longer than 32500 characters")
	}
}

func TestCookieResponseRejectsOversizedPayload(t *testing.T) {
	packet := CookieResponse{Key: "x", HasPayload: true, Payload: make([]byte, 5121)}
	if _, err := packet.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted a cookie payload larger than 5120 bytes")
	}
}

func TestSignUpdateRejectsLongLine(t *testing.T) {
	packet := SignUpdate{Line1: strings.Repeat("a", 385)}
	if _, err := packet.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted a sign line longer than 384 characters")
	}
}

func TestTestInstanceBlockActionUsesDataRecord(t *testing.T) {
	want := TestInstanceBlockAction{
		Position: pk.Position{X: 1, Y: 2, Z: 3},
		Action:   2,
		Data: TestInstanceBlockData{
			Test:           pk.Option[pk.Identifier, *pk.Identifier]{Has: true, Val: pk.Identifier("minecraft:test")},
			Size:           TestInstanceBlockVec3i{X: 4, Y: 5, Z: 6},
			Rotation:       1,
			IgnoreEntities: true,
			Status:         2,
		},
	}
	var encoded bytes.Buffer
	if _, err := want.WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	var got TestInstanceBlockAction
	if _, err := got.ReadFrom(bytes.NewReader(encoded.Bytes())); err != nil {
		t.Fatal(err)
	}
	if !bool(got.Data.Test.Has) || string(got.Data.Test.Val) != "minecraft:test" || got.Data.Size != want.Data.Size || got.Data.Rotation != want.Data.Rotation || got.Data.Status != want.Data.Status || !got.Data.IgnoreEntities {
		t.Fatalf("decoded test instance action = %#v", got)
	}
}
