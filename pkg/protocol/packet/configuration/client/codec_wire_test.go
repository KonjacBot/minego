package client

import (
	"bytes"
	"testing"
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
