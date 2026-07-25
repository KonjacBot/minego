package msa

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
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

func TestDoJSONRejectsOversizedResponse(t *testing.T) {
	auth := NewAuth("", nil)
	auth.Client = &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       io.NopCloser(strings.NewReader(strings.Repeat("x", (1<<20)+1))),
		}, nil
	})}
	req, err := http.NewRequest(http.MethodGet, "https://example.invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := auth.doJSON(req, nil); err == nil {
		t.Fatal("doJSON() accepted an oversized response")
	}
}

func TestDoJSONPropagatesBodyReadError(t *testing.T) {
	auth := NewAuth("", nil)
	auth.Client = &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       io.NopCloser(errorReader{}),
		}, nil
	})}
	req, err := http.NewRequest(http.MethodGet, "https://example.invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := auth.doJSON(req, nil); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("doJSON() error = %v, want io.ErrUnexpectedEOF", err)
	}
}

func TestCallbackHandlerIgnoresUnrelatedInvalidState(t *testing.T) {
	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)
	handler := callbackHandler("/callback", "expected", codeCh, errCh)

	invalid := httptest.NewRequest(http.MethodGet, "/callback?state=wrong&code=attacker", nil)
	invalidResponse := httptest.NewRecorder()
	handler.ServeHTTP(invalidResponse, invalid)
	if invalidResponse.Code != http.StatusBadRequest {
		t.Fatalf("invalid-state status = %d, want %d", invalidResponse.Code, http.StatusBadRequest)
	}
	select {
	case err := <-errCh:
		t.Fatalf("invalid-state request aborted login: %v", err)
	default:
	}

	valid := httptest.NewRequest(http.MethodGet, "/callback?state=expected&code=valid", nil)
	validResponse := httptest.NewRecorder()
	handler.ServeHTTP(validResponse, valid)
	select {
	case code := <-codeCh:
		if code != "valid" {
			t.Fatalf("code = %q, want valid", code)
		}
	default:
		t.Fatal("valid callback did not publish its code")
	}
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }
