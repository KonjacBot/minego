package packet

import (
	"errors"
	"io"
	"testing"

	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"
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

func TestMarshalReturnsFieldErrors(t *testing.T) {
	_, err := Marshal(1, failingField{})
	if !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("Marshal() error = %v, want io.ErrClosedPipe", err)
	}
}

func TestMarshalContainsFieldPanic(t *testing.T) {
	_, err := Marshal(1, panicField{})
	var panicErr *CodecPanicError
	if !errors.As(err, &panicErr) {
		t.Fatalf("Marshal() error = %v, want CodecPanicError", err)
	}
}

func TestScanContainsFieldPanic(t *testing.T) {
	err := Scan(pk.Packet{Data: []byte{1}}, panicField{})
	var panicErr *CodecPanicError
	if !errors.As(err, &panicErr) {
		t.Fatalf("Scan() error = %v, want CodecPanicError", err)
	}
}

type failingField struct{}

func (failingField) WriteTo(io.Writer) (int64, error)  { return 0, io.ErrClosedPipe }
func (failingField) ReadFrom(io.Reader) (int64, error) { return 0, io.ErrClosedPipe }

type panicField struct{}

func (panicField) WriteTo(io.Writer) (int64, error)  { panic("encode") }
func (panicField) ReadFrom(io.Reader) (int64, error) { panic("decode") }

func TestConfigurationCodeOfConductPacketsAreRegistered(t *testing.T) {
	if packet, ok := GetClientPacket(StateConfig, int32(packetid.ClientboundConfigCodeOfConduct)).(*configclient.ConfigCodeOfConduct); !ok || packet == nil {
		t.Fatalf("GetClientPacket(StateConfig, CodeOfConduct) = %T", packet)
	}
	if packet, ok := GetServerPacket(StateConfig, int32(packetid.ServerboundConfigAcceptCodeOfConduct)).(*configserver.ConfigAcceptCodeOfConduct); !ok || packet == nil {
		t.Fatalf("GetServerPacket(StateConfig, AcceptCodeOfConduct) = %T", packet)
	}
}
