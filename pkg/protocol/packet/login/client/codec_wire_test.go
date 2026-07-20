package client

import (
	"bytes"
	"testing"
)

func TestLoginCustomQueryRejectsOversizeData(t *testing.T) {
	packet := LoginCustomQuery{MessageID: 1, Channel: "x", Data: make([]byte, (1<<20)+1)}
	if _, err := packet.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted an oversized login custom query")
	}

	wire := append([]byte{1, 1, 'x'}, make([]byte, (1<<20)+1)...)
	if _, err := new(LoginCustomQuery).ReadFrom(bytes.NewReader(wire)); err == nil {
		t.Fatal("ReadFrom() accepted an oversized login custom query")
	}
}
