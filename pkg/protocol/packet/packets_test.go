package packet

import (
	"testing"

	"github.com/KonjacBot/go-mc/data/packetid"
	configclient "github.com/KonjacBot/minego/pkg/protocol/packet/configuration/client"
	configserver "github.com/KonjacBot/minego/pkg/protocol/packet/configuration/server"
)

func TestGetClientPacketReturnsNilForUnknownID(t *testing.T) {
	for _, state := range []State{StateLogin, StateConfig, StatePlay, State(99)} {
		if packet := GetClientPacket(state, 1<<30); packet != nil {
			t.Fatalf("GetClientPacket(%d, unknown) = %T, want nil", state, packet)
		}
	}
}

func TestConfigurationCodeOfConductPacketsAreRegistered(t *testing.T) {
	if packet, ok := GetClientPacket(StateConfig, int32(packetid.ClientboundConfigCodeOfConduct)).(*configclient.ConfigCodeOfConduct); !ok || packet == nil {
		t.Fatalf("GetClientPacket(StateConfig, CodeOfConduct) = %T", packet)
	}
	if packet, ok := GetServerPacket(StateConfig, int32(packetid.ServerboundConfigAcceptCodeOfConduct)).(*configserver.ConfigAcceptCodeOfConduct); !ok || packet == nil {
		t.Fatalf("GetServerPacket(StateConfig, AcceptCodeOfConduct) = %T", packet)
	}
}
