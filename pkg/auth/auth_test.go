package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"testing"

	mcnet "github.com/KonjacBot/go-mc/net"
	"github.com/KonjacBot/minego/pkg/protocol/packet/login/client"
)

func TestProfileDecodesYggdrasilSelectedProfile(t *testing.T) {
	var response struct {
		SelectedProfile Profile `json:"selectedProfile"`
	}
	err := json.Unmarshal([]byte(`{"selectedProfile":{"id":"12345678123456789abcdef012345678","name":"Steve"}}`), &response)
	if err != nil {
		t.Fatal(err)
	}
	if response.SelectedProfile.Name != "Steve" || response.SelectedProfile.UUID.String() != "12345678-1234-5678-9abc-def012345678" {
		t.Fatalf("decoded profile = %#v", response.SelectedProfile)
	}
}

type nilProfileProvider struct{}

func (nilProfileProvider) Authenticate(context.Context, *mcnet.Conn, client.LoginHello) error {
	return nil
}
func (nilProfileProvider) FetchProfile(context.Context) *Profile { return nil }

func TestHandleLoginRejectsNilProfile(t *testing.T) {
	a := Auth{Conn: &mcnet.Conn{}, Provider: nilProfileProvider{}}
	if err := a.HandleLogin(context.Background()); err == nil {
		t.Fatal("HandleLogin() returned nil for a nil profile")
	}
}

func TestEncryptionResponseRejectsNonRSAKey(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := genEncryptionKeyResponse(make([]byte, 16), publicKey, []byte("token")); err == nil {
		t.Fatal("genEncryptionKeyResponse() accepted a non-RSA key")
	}
}
