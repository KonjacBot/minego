package client

import (
	"bytes"
	"strings"
	"testing"

	"github.com/KonjacBot/go-mc/nbt"
)

func TestConfigCustomPayloadUsesRemainingPacketData(t *testing.T) {
	want := ConfigCustomPayload{Channel: "x", Data: []byte{0x80, 0x01}}
	var encoded bytes.Buffer
	written, err := want.WriteTo(&encoded)
	if err != nil {
		t.Fatal(err)
	}
	if got := encoded.Bytes(); !bytes.Equal(got, []byte{1, 'x', 0x80, 0x01}) {
		t.Fatalf("encoded custom payload = %v, want no data length prefix", got)
	}

	var decoded ConfigCustomPayload
	read, err := decoded.ReadFrom(bytes.NewReader(encoded.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	if read != written || decoded.Channel != want.Channel || !bytes.Equal(decoded.Data, want.Data) {
		t.Fatalf("decoded = %#v after %d bytes, want %#v after %d", decoded, read, want, written)
	}
}

func TestConfigStoreCookieUsesByteArray(t *testing.T) {
	want := ConfigStoreCookie{Key: "x", Payload: []byte{1, 2}}
	var encoded bytes.Buffer
	if _, err := want.WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	if got := encoded.Bytes(); !bytes.Equal(got, []byte{1, 'x', 2, 1, 2}) {
		t.Fatalf("encoded cookie = %v", got)
	}
}

func TestConfigCustomPayloadRejectsOversizeData(t *testing.T) {
	packet := ConfigCustomPayload{Channel: "x", Data: make([]byte, (1<<20)+1)}
	if _, err := packet.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted an oversized config custom payload")
	}

	wire := append([]byte{1, 'x'}, make([]byte, (1<<20)+1)...)
	if _, err := new(ConfigCustomPayload).ReadFrom(bytes.NewReader(wire)); err == nil {
		t.Fatal("ReadFrom() accepted an oversized config custom payload")
	}
}

func TestConfigStoreCookieRejectsOversizePayload(t *testing.T) {
	packet := ConfigStoreCookie{Key: "x", Payload: make([]byte, 5121)}
	if _, err := packet.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted an oversized config cookie")
	}
}

func TestConfigShowDialogUsesSingleNBTValue(t *testing.T) {
	want := ConfigShowDialog{DialogData: nbt.RawMessage{Type: nbt.TagString, Data: []byte{0, 0}}}
	var encoded bytes.Buffer
	if _, err := want.WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded.Bytes(), []byte{nbt.TagString, 0, 0}) {
		t.Fatalf("encoded dialog = %v", encoded.Bytes())
	}

	var decoded ConfigShowDialog
	if _, err := decoded.ReadFrom(bytes.NewReader(encoded.Bytes())); err != nil {
		t.Fatal(err)
	}
	if decoded.DialogData.Type != want.DialogData.Type || !bytes.Equal(decoded.DialogData.Data, want.DialogData.Data) {
		t.Fatalf("decoded dialog = %v", decoded.DialogData)
	}
}

func TestConfigCodeOfConductRejectsOversizeText(t *testing.T) {
	packet := ConfigCodeOfConduct{CodeOfConduct: strings.Repeat("a", 32768)}
	if _, err := packet.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted an oversized code of conduct")
	}
}
