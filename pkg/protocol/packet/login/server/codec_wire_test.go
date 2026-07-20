package server

import (
	"bytes"
	"testing"
)

func TestLoginCookieResponseRejectsOversizePayload(t *testing.T) {
	packet := LoginCookieResponse{Key: "x", HasPayload: true, Payload: make([]byte, 5121)}
	if _, err := packet.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted an oversized login cookie payload")
	}
}

func TestLoginCustomQueryAnswerRejectsOversizeData(t *testing.T) {
	packet := LoginCustomQueryAnswer{MessageID: 1, HasData: true, Data: make([]byte, (1<<20)+1)}
	if _, err := packet.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted an oversized login custom query answer")
	}
}
