package client

import (
	"bytes"
	"testing"

	pk "github.com/KonjacBot/go-mc/net/packet"
	"github.com/KonjacBot/minego/pkg/protocol/component"
	"github.com/google/uuid"
)

func TestSetContainerContentConsumesResolvedProfileBeforeNextSlot(t *testing.T) {
	profileUUID := uuid.MustParse("b0f05e1a-0594-4c08-ba51-d57dde9602a7")
	var wire bytes.Buffer
	for _, value := range []pk.VarInt{
		1,    // window ID
		1,    // state ID
		2,    // slots
		1,    // first slot count
		1265, // player head item ID
		1,    // added components
		0,    // removed components
		70,   // minecraft:profile
		1,    // resolved profile variant
	} {
		if _, err := value.WriteTo(&wire); err != nil {
			t.Fatal(err)
		}
	}
	wire.Write(profileUUID[:])
	if _, err := pk.String("PowRu").WriteTo(&wire); err != nil {
		t.Fatal(err)
	}
	if _, err := pk.VarInt(0).WriteTo(&wire); err != nil { // profile properties
		t.Fatal(err)
	}
	wire.Write([]byte{0, 0, 0, 0}) // body, cape, elytra, and model overrides
	wire.WriteByte(0)              // second slot is empty
	wire.WriteByte(0)              // carried item is empty
	wireLength := wire.Len()

	var content SetContainerContent
	read, err := content.ReadFrom(&wire)
	if err != nil {
		t.Fatal(err)
	}
	if read != int64(wireLength) || wire.Len() != 0 {
		t.Fatalf("container decoder consumed %d/%d bytes", read, wireLength)
	}
	if len(content.Slots) != 2 || content.Slots[1].Count != 0 || content.CarriedItem.Count != 0 {
		t.Fatalf("decoded slots = %#v, carried = %#v", content.Slots, content.CarriedItem)
	}
	profile, ok := content.Slots[0].AddComponent[70].(*component.Profile)
	if !ok || profile.Profile.GameProfile == nil {
		t.Fatalf("decoded profile component = %#v", content.Slots[0].AddComponent[70])
	}
	if got := profile.Profile.GameProfile; got.UUID != profileUUID || got.Name != "PowRu" || len(got.Properties) != 0 {
		t.Fatalf("decoded game profile = %#v", got)
	}
}
