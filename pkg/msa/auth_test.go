package msa

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return fn(request)
}

func TestTransientMicrosoftRefreshFailureDoesNotStartInteractiveLogin(t *testing.T) {
	openedLogin := false
	auth := NewAuth("", nil)
	auth.MSToken = &MSToken{RefreshToken: "refresh", Expiry: time.Now().Add(-time.Minute)}
	auth.Client = &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("temporary network failure")
	})}
	auth.OnAuthURL = func(context.Context, string) error {
		openedLogin = true
		return nil
	}

	if err := auth.ensureMSToken(context.Background()); err == nil {
		t.Fatal("temporary refresh failure returned nil")
	}
	if openedLogin {
		t.Fatal("temporary refresh failure opened an interactive login")
	}
	if auth.MSToken == nil || auth.MSToken.RefreshToken != "refresh" {
		t.Fatal("temporary refresh failure discarded the reusable refresh token")
	}
}

func TestOnlyPermanentMicrosoftOAuthErrorsRequireInteractiveLogin(t *testing.T) {
	if !requiresInteractiveLogin(&oauthError{Code: "invalid_grant"}) {
		t.Fatal("invalid_grant should require interactive login")
	}
	if requiresInteractiveLogin(&oauthError{Code: "temporarily_unavailable"}) {
		t.Fatal("temporary OAuth failure should be retried without interactive login")
	}
}

func TestSessionAuthenticationClassifiesDefinitiveCredentialRejection(t *testing.T) {
	auth := NewAuth("", nil)
	auth.MSToken = &MSToken{RefreshToken: "rejected", Expiry: time.Now().Add(-time.Minute)}
	auth.Client = &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Status:     "400 Bad Request",
			Body:       io.NopCloser(strings.NewReader(`{"error":"invalid_grant","error_description":"refresh token revoked"}`)),
		}, nil
	})}

	err := auth.prepareSessionAuthentication(context.Background())
	if !errors.Is(err, ErrAuthRequired) {
		t.Fatalf("session authentication error = %v, want ErrAuthRequired", err)
	}
}

func TestSessionAuthenticationDoesNotClassifyNetworkFailure(t *testing.T) {
	auth := NewAuth("", nil)
	auth.MSToken = &MSToken{RefreshToken: "refresh", Expiry: time.Now().Add(-time.Minute)}
	networkErr := errors.New("temporary network failure")
	auth.Client = &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, networkErr
	})}

	err := auth.prepareSessionAuthentication(context.Background())
	if !errors.Is(err, networkErr) || errors.Is(err, ErrAuthRequired) {
		t.Fatalf("session authentication error = %v, want unclassified network failure", err)
	}
}
