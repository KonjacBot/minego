package packet

import "testing"

func TestGetClientPacketReturnsNilForUnknownID(t *testing.T) {
	for _, state := range []State{StateLogin, StateConfig, StatePlay, State(99)} {
		if packet := GetClientPacket(state, 1<<30); packet != nil {
			t.Fatalf("GetClientPacket(%d, unknown) = %T, want nil", state, packet)
		}
	}
}
