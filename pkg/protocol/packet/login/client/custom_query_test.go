package client

import (
	"bytes"
	"testing"
)

func TestLoginCustomQueryUsesRemainingPacketData(t *testing.T) {
	want := LoginCustomQuery{MessageID: 1, Channel: "x", Data: []byte{0x80, 0x01}}
	var encoded bytes.Buffer
	written, err := want.WriteTo(&encoded)
	if err != nil {
		t.Fatal(err)
	}
	if got := encoded.Bytes(); !bytes.Equal(got, []byte{1, 1, 'x', 0x80, 0x01}) {
		t.Fatalf("encoded custom query = %v, want no data length prefix", got)
	}

	var decoded LoginCustomQuery
	read, err := decoded.ReadFrom(bytes.NewReader(encoded.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	if read != written || decoded.MessageID != want.MessageID || decoded.Channel != want.Channel || !bytes.Equal(decoded.Data, want.Data) {
		t.Fatalf("decoded = %#v after %d bytes, want %#v after %d", decoded, read, want, written)
	}
}
