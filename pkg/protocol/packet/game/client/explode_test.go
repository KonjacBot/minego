package client

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// TestExplodeDecodesServerPacket decodes a 26.2 explosion captured from a live server (a TNT next to the bot):
// it used to stop after the knockback field and fail with unexpected EOF, ending the connection.
func TestExplodeDecodesServerPacket(t *testing.T) {
	raw, err := hex.DecodeString("c09b6570f264824940" + "50a519c7eba82e40a2ce09b737fb2440600000000000000021ed0d00")
	if err != nil {
		t.Fatal(err)
	}
	var p Explode
	n, err := p.ReadFrom(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("ReadFrom(): %v", err)
	}
	if n != int64(len(raw)) {
		t.Fatalf("read %d bytes, packet has %d", n, len(raw))
	}
	if p.Radius != 3.5 || p.BlockCount != 0 || p.PlayerKnockbackVelocity.Has || p.ExplosionParticle.ID != 33 || len(p.BlockParticles) != 0 {
		t.Errorf("decoded %+v", p)
	}
	if int(p.CenterX) != -1753 || int(p.CenterZ) != 2407 {
		t.Errorf("center = %v, %v, %v", p.CenterX, p.CenterY, p.CenterZ)
	}

	var encoded bytes.Buffer
	if _, err := p.WriteTo(&encoded); err != nil {
		t.Fatalf("WriteTo(): %v", err)
	}
	if !bytes.Equal(encoded.Bytes(), raw) {
		t.Errorf("round trip = %x, want %x", encoded.Bytes(), raw)
	}
}
