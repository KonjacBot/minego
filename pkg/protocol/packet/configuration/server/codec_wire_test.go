package server

import (
	"bytes"
	"testing"

	pk "github.com/KonjacBot/go-mc/net/packet"
)

func TestConfigCustomClickActionUsesLengthPrefixedOptionalNBT(t *testing.T) {
	var encoded bytes.Buffer
	if _, err := (ConfigCustomClickAction{Action: "x"}).WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	if got := encoded.Bytes(); !bytes.Equal(got, []byte{1, 'x', 1, 0}) {
		t.Fatalf("encoded config custom click action = %v", got)
	}
}

func TestConfigCustomClickActionRejectsOversizedPayload(t *testing.T) {
	var wire bytes.Buffer
	if _, err := pk.Identifier("x").WriteTo(&wire); err != nil {
		t.Fatal(err)
	}
	if _, err := pk.VarInt(65537).WriteTo(&wire); err != nil {
		t.Fatal(err)
	}
	if _, err := new(ConfigCustomClickAction).ReadFrom(bytes.NewReader(wire.Bytes())); err == nil {
		t.Fatal("ReadFrom() accepted an oversized config custom click payload")
	}
}

func TestConfigCustomClickActionRejectsTrailingGarbageInsideLengthPrefix(t *testing.T) {
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
	if _, err := new(ConfigCustomClickAction).ReadFrom(bytes.NewReader(wire.Bytes())); err == nil {
		t.Fatal("ReadFrom() accepted trailing garbage inside a length-prefixed config custom click payload")
	}
}
