package protocol

import (
	"bytes"
	"testing"
)

func TestResolvableProfileFullProfileRoundTrip(t *testing.T) {
	want := ResolvableProfile{
		Type:        1,
		GameProfile: &GameProfile{Name: "Alex"},
	}
	var encoded bytes.Buffer
	written, err := want.WriteTo(&encoded)
	if err != nil {
		t.Fatal(err)
	}

	var got ResolvableProfile
	read, err := got.ReadFrom(bytes.NewReader(encoded.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	if read != written {
		t.Fatalf("ReadFrom() = %d bytes, WriteTo() = %d bytes", read, written)
	}
	if got.GameProfile == nil || got.GameProfile.Name != want.GameProfile.Name {
		t.Fatalf("decoded full profile = %#v, want %#v", got.GameProfile, want.GameProfile)
	}
}

func TestResolvableProfileRejectsUnknownType(t *testing.T) {
	if _, err := new(ResolvableProfile).ReadFrom(bytes.NewReader([]byte{2})); err == nil {
		t.Fatal("ReadFrom() accepted an unknown profile type")
	}
	if _, err := (ResolvableProfile{Type: 2}).WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted an unknown profile type")
	}
}
