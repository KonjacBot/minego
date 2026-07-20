package client

import (
	"bytes"
	"testing"

	"github.com/KonjacBot/go-mc/chat/sign"
)

func TestChangeDifficultyUsesVarInt(t *testing.T) {
	want := ChangeDifficulty{Difficulty: 128, Locked: true}
	var encoded bytes.Buffer
	written, err := want.WriteTo(&encoded)
	if err != nil {
		t.Fatal(err)
	}
	if got := encoded.Bytes(); !bytes.Equal(got, []byte{0x80, 0x01, 0x01}) {
		t.Fatalf("encoded difficulty = %v, want VarInt followed by bool", got)
	}
	if written != int64(encoded.Len()) {
		t.Fatalf("WriteTo() count = %d, encoded length = %d", written, encoded.Len())
	}

	var decoded ChangeDifficulty
	read, err := decoded.ReadFrom(bytes.NewReader(encoded.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	if read != written || decoded != want {
		t.Fatalf("decoded = %#v after %d bytes, want %#v after %d", decoded, read, want, written)
	}
}

func TestDisguisedChatReadsBoundChatType(t *testing.T) {
	wire := appendNBTString(nil, "message")
	wire = append(wire, 2) // One-based chat_type holder ID.
	wire = appendNBTString(wire, "sender")
	wire = append(wire, 0) // No target name.
	var got DisguisedChat
	read, err := got.ReadFrom(bytes.NewReader(wire))
	if err != nil {
		t.Fatal(err)
	}
	if read != int64(len(wire)) || got.Message.Text != "message" || got.ChatType.HolderID != 2 || got.ChatType.Name.Text != "sender" || got.ChatType.TargetName.Has {
		t.Fatalf("decoded = %#v after %d of %d bytes", got, read, len(wire))
	}
}

func appendNBTString(dst []byte, value string) []byte {
	dst = append(dst, 8, byte(len(value)>>8), byte(len(value)))
	return append(dst, value...)
}

func TestPackedMessageSignatureWireFormats(t *testing.T) {
	var cached bytes.Buffer
	if _, err := (PackedMessageSignature{ID: 4}).WriteTo(&cached); err != nil {
		t.Fatal(err)
	}
	if got := cached.Bytes(); !bytes.Equal(got, []byte{5}) {
		t.Fatalf("cached signature = %v, want [5]", got)
	}
	var cachedDecoded PackedMessageSignature
	if _, err := cachedDecoded.ReadFrom(bytes.NewReader(cached.Bytes())); err != nil {
		t.Fatal(err)
	}
	if cachedDecoded.ID != 4 || cachedDecoded.Signature != nil {
		t.Fatalf("decoded cached signature = %#v", cachedDecoded)
	}

	signature := new(sign.Signature)
	for i := range signature {
		signature[i] = byte(i)
	}
	var inline bytes.Buffer
	written, err := (PackedMessageSignature{Signature: signature}).WriteTo(&inline)
	if err != nil {
		t.Fatal(err)
	}
	if written != 257 || inline.Len() != 257 || inline.Bytes()[0] != 0 {
		t.Fatalf("inline signature size = (%d, %d), first byte = %d", written, inline.Len(), inline.Bytes()[0])
	}
	var inlineDecoded PackedMessageSignature
	if _, err := inlineDecoded.ReadFrom(bytes.NewReader(inline.Bytes())); err != nil {
		t.Fatal(err)
	}
	if inlineDecoded.ID != -1 || inlineDecoded.Signature == nil || *inlineDecoded.Signature != *signature {
		t.Fatal("inline signature did not round-trip")
	}

	if _, err := inlineDecoded.ReadFrom(bytes.NewReader(inline.Bytes()[:100])); err == nil {
		t.Fatal("ReadFrom() accepted a truncated inline signature")
	}
}
