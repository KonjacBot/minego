package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	mcnet "github.com/KonjacBot/go-mc/net"
	"github.com/KonjacBot/minego/pkg/protocol/packet/login/client"
	"github.com/google/uuid"
)

type profileTransport func(*http.Request) (*http.Response, error)

func (f profileTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestKonjacAuthFetchProfileYggdrasilID(t *testing.T) {
	original := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = original })
	http.DefaultClient = &http.Client{Transport: profileTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "http://127.0.0.1:37565/as/authenticate" {
			t.Fatalf("unexpected authentication URL: %s", r.URL)
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(
			`{"selectedProfile":{"name":"ProxyAccount","id":"12345678123456781234567812345678"}}`)), Header: make(http.Header)}, nil
	})}
	profile := (&KonjacAuth{UserCode: "test-code"}).FetchProfile(context.Background())
	if profile == nil || profile.Name == "" || profile.UUID == uuid.Nil {
		t.Fatal("invalid upstream authentication profile: selectedProfile.id was not decoded")
	}
	if profile.Name != "ProxyAccount" || profile.UUID.String() != "12345678-1234-5678-1234-567812345678" {
		t.Fatal("incorrect authenticated account profile")
	}
}

func TestKonjacAuthProfileValidation(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		valid      bool
	}{
		{"hyphenated id", `{"selectedProfile":{"name":"ProxyAccount","id":"12345678-1234-5678-1234-567812345678"}}`, true},
		{"legacy UUID", `{"selectedProfile":{"Name":"ProxyAccount","UUID":"12345678-1234-5678-1234-567812345678"}}`, true},
		{"missing profile", `{}`, false},
		{"null profile", `{"selectedProfile":null}`, false},
		{"missing UUID", `{"selectedProfile":{"name":"ProxyAccount"}}`, false},
		{"missing name", `{"selectedProfile":{"id":"12345678123456781234567812345678"}}`, false},
		{"malformed UUID", `{"selectedProfile":{"name":"ProxyAccount","id":"invalid"}}`, false},
		{"conflicting UUID", `{"selectedProfile":{"name":"ProxyAccount","id":"12345678123456781234567812345678","uuid":"87654321876543218765432187654321"}}`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			original := http.DefaultClient
			t.Cleanup(func() { http.DefaultClient = original })
			http.DefaultClient = &http.Client{Transport: profileTransport(func(*http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(tc.body)), Header: make(http.Header)}, nil
			})}
			profile := (&KonjacAuth{UserCode: "test-code"}).FetchProfile(context.Background())
			if (profile != nil) != tc.valid {
				t.Fatalf("profile valid=%v, want %v", profile != nil, tc.valid)
			}
		})
	}
}

// Opt-in: uses the local authentication service; never logs credentials or profile data.
func TestKonjacAuthLocalProfileLive(t *testing.T) {
	if os.Getenv("RNG_LIVE_AUTH") != "1" {
		t.Skip("local KonjacAuth check is opt-in")
	}
	code := os.Getenv("MINEGO_USER_CODE")
	if code == "" {
		t.Fatal("MINEGO_USER_CODE is unset")
	}
	profile := (&KonjacAuth{UserCode: code}).FetchProfile(context.Background())
	if profile == nil || profile.Name == "" || profile.UUID == uuid.Nil {
		t.Fatal("invalid upstream authentication profile from local KonjacAuth")
	}
	t.Log("local KonjacAuth returned a nonempty name and nonzero UUID")
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
