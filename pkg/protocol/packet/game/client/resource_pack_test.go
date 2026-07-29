package client

import (
	"testing"

	"github.com/KonjacBot/go-mc/data/packetid"
)

func TestResourcePackPacketIDs(t *testing.T) {
	if got := (&AddResourcePack{}).PacketID(); got != packetid.ClientboundResourcePackPush {
		t.Fatalf("AddResourcePack packet ID = %v, want push", got)
	}
	if got := (&RemoveResourcePack{}).PacketID(); got != packetid.ClientboundResourcePackPop {
		t.Fatalf("RemoveResourcePack packet ID = %v, want pop", got)
	}
}
