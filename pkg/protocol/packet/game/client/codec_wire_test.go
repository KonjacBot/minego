package client

import (
	"bytes"
	"math"
	"strings"
	"testing"

	"github.com/KonjacBot/go-mc/chat/sign"
	"github.com/KonjacBot/go-mc/nbt"
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

func TestCustomPayloadRejectsOversizeData(t *testing.T) {
	packet := CustomPayload{Channel: "x", Data: make([]byte, maxRemainingPayloadBytes+1)}
	if _, err := packet.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted an oversized custom payload")
	}

	wire := append([]byte{1, 'x'}, make([]byte, maxRemainingPayloadBytes+1)...)
	if _, err := new(CustomPayload).ReadFrom(bytes.NewReader(wire)); err == nil {
		t.Fatal("ReadFrom() accepted an oversized custom payload")
	}
}

func TestCustomReportDetailsUsesBoundedMap(t *testing.T) {
	want := CustomReportDetails{Details: map[string]string{"k": "v"}}
	var encoded bytes.Buffer
	if _, err := want.WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	if got := encoded.Bytes(); !bytes.Equal(got, []byte{1, 1, 'k', 1, 'v'}) {
		t.Fatalf("encoded report details = %v", got)
	}

	var decoded CustomReportDetails
	if _, err := decoded.ReadFrom(bytes.NewReader(encoded.Bytes())); err != nil {
		t.Fatal(err)
	}
	if decoded.Details["k"] != "v" || len(decoded.Details) != 1 {
		t.Fatalf("decoded report details = %#v", decoded.Details)
	}
}

func TestCustomReportDetailsRejectsTooManyEntries(t *testing.T) {
	details := make(map[string]string, maxCustomReportDetailCount+1)
	for i := 0; i <= maxCustomReportDetailCount; i++ {
		details[strings.Repeat("a", 1)+string(rune('a'+i))] = "v"
	}
	if _, err := (CustomReportDetails{Details: details}).WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted too many custom report details")
	}
}

func TestAddResourcePackRejectsLongHash(t *testing.T) {
	packet := AddResourcePack{Hash: strings.Repeat("a", 41)}
	if _, err := packet.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted a resource pack hash longer than 40 characters")
	}
}

func TestStoreCookieRejectsOversizedPayload(t *testing.T) {
	packet := StoreCookie{Key: "x", Payload: make([]byte, maxCookiePayloadBytes+1)}
	if _, err := packet.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted a cookie payload larger than 5120 bytes")
	}
}

func TestShowDialogUsesInlineDirectHolder(t *testing.T) {
	want := ShowDialog{DialogData: nbt.RawMessage{Type: nbt.TagString, Data: []byte{0, 0}}}
	var encoded bytes.Buffer
	if _, err := want.WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	if got := encoded.Bytes(); !bytes.Equal(got, []byte{0, nbt.TagString, 0, 0}) {
		t.Fatalf("encoded show dialog = %v", got)
	}

	var decoded ShowDialog
	if _, err := decoded.ReadFrom(bytes.NewReader(encoded.Bytes())); err != nil {
		t.Fatal(err)
	}
	if decoded.HasRegistryID || decoded.RegistryID != 0 || decoded.DialogData.Type != want.DialogData.Type || !bytes.Equal(decoded.DialogData.Data, want.DialogData.Data) {
		t.Fatalf("decoded show dialog = %#v", decoded)
	}
}

func TestShowDialogUsesRegistryHolderReference(t *testing.T) {
	want := ShowDialog{HasRegistryID: true, RegistryID: 1}
	var encoded bytes.Buffer
	if _, err := want.WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	if got := encoded.Bytes(); !bytes.Equal(got, []byte{2}) {
		t.Fatalf("encoded show dialog holder = %v", got)
	}

	var decoded ShowDialog
	if _, err := decoded.ReadFrom(bytes.NewReader(encoded.Bytes())); err != nil {
		t.Fatal(err)
	}
	if !decoded.HasRegistryID || decoded.RegistryID != 1 || decoded.DialogData.Type != 0 || len(decoded.DialogData.Data) != 0 {
		t.Fatalf("decoded show dialog holder = %#v", decoded)
	}
}

func TestShowDialogClearsInlineStateWhenReusedForHolderReference(t *testing.T) {
	var decoded ShowDialog
	if _, err := decoded.ReadFrom(bytes.NewReader([]byte{0, nbt.TagString, 0, 0})); err != nil {
		t.Fatal(err)
	}
	if _, err := decoded.ReadFrom(bytes.NewReader([]byte{2})); err != nil {
		t.Fatal(err)
	}
	if !decoded.HasRegistryID || decoded.RegistryID != 1 || decoded.DialogData.Type != 0 || len(decoded.DialogData.Data) != 0 {
		t.Fatalf("reused decoded show dialog = %#v", decoded)
	}
}

func TestShowDialogClearsRegistryStateWhenReusedForInlineValue(t *testing.T) {
	decoded := ShowDialog{HasRegistryID: true, RegistryID: 7}
	if _, err := decoded.ReadFrom(bytes.NewReader([]byte{0, nbt.TagString, 0, 0})); err != nil {
		t.Fatal(err)
	}
	if decoded.HasRegistryID || decoded.RegistryID != 0 || decoded.DialogData.Type != nbt.TagString || !bytes.Equal(decoded.DialogData.Data, []byte{0, 0}) {
		t.Fatalf("reused decoded show dialog = %#v", decoded)
	}
}

func TestShowDialogRejectsNegativeRegistryIDOnWrite(t *testing.T) {
	packet := ShowDialog{HasRegistryID: true, RegistryID: -1}
	if _, err := packet.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted a negative registry ID")
	}
}

func TestShowDialogRejectsOverflowRegistryIDOnWrite(t *testing.T) {
	packet := ShowDialog{HasRegistryID: true, RegistryID: math.MaxInt32}
	if _, err := packet.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted an overflowing registry ID")
	}
}
