package msa

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	mcnet "github.com/KonjacBot/go-mc/net"
	"github.com/KonjacBot/minego/pkg/auth"
	"github.com/KonjacBot/minego/pkg/protocol/packet/login/client"
	"github.com/google/uuid"
)

const (
	msAuthorizeURL  = "https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize"
	msDeviceCodeURL = "https://login.microsoftonline.com/consumers/oauth2/v2.0/devicecode"
	msTokenURL      = "https://login.microsoftonline.com/consumers/oauth2/v2.0/token"

	xblAuthURL     = "https://user.auth.xboxlive.com/user/authenticate"
	xstsAuthURL    = "https://xsts.auth.xboxlive.com/xsts/authorize"
	mcXboxLoginURL = "https://api.minecraftservices.com/authentication/login_with_xbox"
	mcProfileURL   = "https://api.minecraftservices.com/minecraft/profile"
	mojangJoinURL  = "https://sessionserver.mojang.com/session/minecraft/join"

	defaultMSScope  = "XboxLive.signin offline_access"
	defaultClientId = "e22c5ce3-fd95-4c29-997b-cfc256b5a05a"
)

var ErrAuthRequired = errors.New("microsoft authentication required")

type LoginMode string

const (
	LoginModeAuthCode   LoginMode = "auth_code"
	LoginModeDeviceCode LoginMode = "device_code"
)

type TokenStore interface {
	LoadToken(ctx context.Context) (*TokenState, error)
	SaveToken(ctx context.Context, state *TokenState) error
}

type DeviceCodeHandler func(ctx context.Context, code MSDeviceCode) error
type AuthURLHandler func(ctx context.Context, authURL string) error

type Auth struct {
	ClientID string
	Scope    string
	Client   *http.Client
	Store    TokenStore
	Mode     LoginMode

	RedirectHost     string
	CallbackPath     string
	LocalhostTimeout time.Duration
	OnAuthURL        AuthURLHandler

	OnCode DeviceCodeHandler

	MSToken   *MSToken
	XBLToken  *XBToken
	XSTSToken *XBToken
	MCToken   *MCToken
	MCProfile *MCProfile

	loaded  bool
	profile *auth.Profile
}

type TokenState struct {
	MSToken   *MSToken   `json:"ms_token,omitempty" toml:"ms_token"`
	XBLToken  *XBToken   `json:"xbl_token,omitempty" toml:"xbl_token"`
	XSTSToken *XBToken   `json:"xsts_token,omitempty" toml:"xsts_token"`
	MCToken   *MCToken   `json:"mc_token,omitempty" toml:"mc_token"`
	MCProfile *MCProfile `json:"mc_profile,omitempty" toml:"mc_profile"`
}

type MSDeviceCode struct {
	UserCode                string `json:"user_code"`
	DeviceCode              string `json:"device_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURL         string `json:"verification_url"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
	Message                 string `json:"message"`
}

type MSToken struct {
	TokenType    string    `json:"token_type" toml:"token_type"`
	Scope        string    `json:"scope" toml:"scope"`
	AccessToken  string    `json:"access_token" toml:"access_token"`
	RefreshToken string    `json:"refresh_token" toml:"refresh_token"`
	ExpiresIn    int       `json:"expires_in" toml:"expires_in"`
	Expiry       time.Time `json:"expiry" toml:"expiry"`
}

type XBToken struct {
	IssueInstant  time.Time `json:"IssueInstant" toml:"issue_instant"`
	NotAfter      time.Time `json:"NotAfter" toml:"not_after"`
	Token         string    `json:"Token" toml:"token"`
	DisplayClaims struct {
		Xui []struct {
			Uhs string `json:"uhs" toml:"uhs"`
		} `json:"xui" toml:"xui"`
	} `json:"DisplayClaims" toml:"display_claims"`
}

type MCToken struct {
	Username    string        `json:"username" toml:"username"`
	Roles       []interface{} `json:"roles" toml:"roles"`
	AccessToken string        `json:"access_token" toml:"access_token"`
	TokenType   string        `json:"token_type" toml:"token_type"`
	ExpiresIn   int           `json:"expires_in" toml:"expires_in"`
	Expiry      time.Time     `json:"expiry" toml:"expiry"`
}

type MCProfile struct {
	ID   string `json:"id" toml:"id"`
	Name string `json:"name" toml:"name"`
}

func NewAuth(clientID string, store TokenStore) *Auth {
	if clientID == "" {
		clientID = defaultClientId
	}
	return &Auth{
		ClientID:         clientID,
		Scope:            defaultMSScope,
		Client:           http.DefaultClient,
		Store:            store,
		Mode:             LoginModeAuthCode,
		RedirectHost:     "localhost",
		CallbackPath:     "/callback",
		LocalhostTimeout: 3 * time.Minute,
	}
}

func (m *Auth) Login(ctx context.Context) error {
	m.normalize()
	if err := m.load(ctx); err != nil {
		return err
	}
	if m.hasUsableMinecraftToken() {
		return nil
	}
	if err := m.ensureMSToken(ctx); err != nil {
		return err
	}
	if err := m.authXBL(ctx); err != nil {
		return err
	}
	if err := m.authXSTS(ctx); err != nil {
		return err
	}
	if err := m.authMinecraft(ctx); err != nil {
		return err
	}
	if err := m.fetchMCProfile(ctx); err != nil {
		return err
	}
	return m.save(ctx)
}

func (m *Auth) Authenticate(ctx context.Context, conn *mcnet.Conn, content client.LoginHello) error {
	if !m.hasUsableMinecraftToken() {
		if err := m.prepareSessionAuthentication(ctx); err != nil {
			return errors.Join(auth.ErrEncrypt, err)
		}
	}

	return (&auth.OnlineAuth{
		AccessToken: m.MCToken.AccessToken,
		Profile:     *m.profile,
	}).Authenticate(ctx, conn, content)
}

func (m *Auth) CachedProfile(ctx context.Context) *auth.Profile {
	m.normalize()
	if err := m.load(ctx); err != nil {
		return nil
	}
	if m.profile != nil {
		return m.profile
	}
	_ = m.rebuildProfile()
	return m.profile
}

func (m *Auth) FetchProfile(ctx context.Context) *auth.Profile {
	if m.profile == nil {
		_ = m.Login(ctx)
	}
	return m.profile
}

func (m *Auth) prepareSessionAuthentication(ctx context.Context) error {
	m.normalize()
	if err := m.load(ctx); err != nil {
		return err
	}
	if m.hasUsableMinecraftToken() {
		return nil
	}
	if err := m.ensureSessionMSToken(ctx); err != nil {
		return err
	}
	for _, step := range []struct {
		name string
		run  func(context.Context) error
	}{
		{name: "Xbox Live", run: m.authXBL},
		{name: "XSTS", run: m.authXSTS},
		{name: "Minecraft", run: m.authMinecraft},
		{name: "Minecraft profile", run: m.fetchMCProfile},
	} {
		if err := step.run(ctx); err != nil {
			if definitiveCredentialRejection(err) {
				return errors.Join(ErrAuthRequired, fmt.Errorf("%s authentication: %w", step.name, err))
			}
			return err
		}
	}
	return m.save(ctx)
}

func (m *Auth) ensureSessionMSToken(ctx context.Context) error {
	if m.ClientID == "" {
		return errors.New("missing Microsoft OAuth client id")
	}
	if m.MSToken != nil && m.MSToken.AccessToken != "" && m.MSToken.Expiry.After(time.Now().Add(30*time.Second)) {
		return nil
	}
	if m.MSToken == nil || m.MSToken.RefreshToken == "" {
		return ErrAuthRequired
	}
	if err := m.refreshMSToken(ctx); err != nil {
		if definitiveCredentialRejection(err) {
			return errors.Join(ErrAuthRequired, err)
		}
		return err
	}
	return m.save(ctx)
}

func definitiveCredentialRejection(err error) bool {
	var oauthErr *oauthError
	if errors.As(err, &oauthErr) {
		switch oauthErr.Code {
		case "invalid_grant", "interaction_required", "consent_required", "invalid_token":
			return true
		}
	}
	var statusErr *httpStatusError
	return errors.As(err, &statusErr) && statusErr.StatusCode == http.StatusUnauthorized
}

func (m *Auth) normalize() {
	if m.Client == nil {
		m.Client = http.DefaultClient
	}
	if m.Scope == "" {
		m.Scope = defaultMSScope
	}
	if m.Mode == "" {
		m.Mode = LoginModeAuthCode
	}
	if m.RedirectHost == "" {
		m.RedirectHost = "localhost"
	}
	if m.CallbackPath == "" {
		m.CallbackPath = "/callback"
	}
	if !strings.HasPrefix(m.CallbackPath, "/") {
		m.CallbackPath = "/" + m.CallbackPath
	}
	if m.LocalhostTimeout <= 0 {
		m.LocalhostTimeout = 3 * time.Minute
	}
}

func (m *Auth) load(ctx context.Context) error {
	if m.loaded || m.Store == nil {
		m.loaded = true
		return nil
	}
	state, err := m.Store.LoadToken(ctx)
	if err != nil {
		return fmt.Errorf("load Microsoft auth token state: %w", err)
	}
	m.loaded = true
	if state == nil {
		return nil
	}
	m.MSToken = state.MSToken
	m.XBLToken = state.XBLToken
	m.XSTSToken = state.XSTSToken
	m.MCToken = state.MCToken
	m.MCProfile = state.MCProfile
	return m.rebuildProfile()
}

func (m *Auth) save(ctx context.Context) error {
	if m.Store == nil {
		return nil
	}
	return m.Store.SaveToken(ctx, &TokenState{
		MSToken:   m.MSToken,
		XBLToken:  m.XBLToken,
		XSTSToken: m.XSTSToken,
		MCToken:   m.MCToken,
		MCProfile: m.MCProfile,
	})
}

func (m *Auth) rebuildProfile() error {
	if m.MCProfile == nil || m.MCProfile.ID == "" || m.MCProfile.Name == "" {
		return nil
	}
	id, err := uuidFromUndashed(m.MCProfile.ID)
	if err != nil {
		return err
	}
	m.profile = &auth.Profile{Name: m.MCProfile.Name, UUID: id}
	return nil
}

func (m *Auth) hasUsableMinecraftToken() bool {
	if m.MCToken == nil || m.MCToken.AccessToken == "" || m.MCProfile == nil {
		return false
	}
	if !m.MCToken.Expiry.After(time.Now().Add(30 * time.Second)) {
		return false
	}
	return m.rebuildProfile() == nil && m.profile != nil
}

func (m *Auth) ensureMSToken(ctx context.Context) error {
	if m.ClientID == "" {
		return errors.New("missing Microsoft OAuth client id")
	}
	if m.MSToken != nil && m.MSToken.AccessToken != "" && m.MSToken.Expiry.After(time.Now().Add(30*time.Second)) {
		return nil
	}
	if m.MSToken != nil && m.MSToken.RefreshToken != "" {
		if err := m.refreshMSToken(ctx); err == nil {
			return m.save(ctx)
		} else if !requiresInteractiveLogin(err) {
			return err
		}
		m.MSToken = nil
	}
	if m.Mode == LoginModeDeviceCode {
		return m.loginByDeviceCode(ctx)
	}
	if err := m.loginByAuthCode(ctx); err != nil {
		if m.OnCode != nil {
			return m.loginByDeviceCode(ctx)
		}
		return err
	}
	return nil
}

func requiresInteractiveLogin(err error) bool {
	var oauthErr *oauthError
	if !errors.As(err, &oauthErr) {
		return false
	}
	switch oauthErr.Code {
	case "invalid_grant", "interaction_required", "consent_required":
		return true
	default:
		return false
	}
}

func (m *Auth) loginByAuthCode(ctx context.Context) error {
	ln, err := net.Listen("tcp", m.RedirectHost+":0")
	if err != nil {
		return fmt.Errorf("listen localhost OAuth callback: %w", err)
	}
	defer ln.Close()

	addr := ln.Addr().(*net.TCPAddr)
	redirectURI := fmt.Sprintf("http://%s:%d%s", m.RedirectHost, addr.Port, m.CallbackPath)
	state, err := randomURLSafe(32)
	if err != nil {
		return err
	}
	verifier, err := randomURLSafe(64)
	if err != nil {
		return err
	}
	challenge := pkceChallenge(verifier)

	authValues := url.Values{}
	authValues.Set("client_id", m.ClientID)
	authValues.Set("response_type", "code")
	authValues.Set("redirect_uri", redirectURI)
	authValues.Set("scope", m.Scope)
	authValues.Set("state", state)
	authValues.Set("prompt", "select_account")
	authValues.Set("code_challenge", challenge)
	authValues.Set("code_challenge_method", "S256")
	authURL := msAuthorizeURL + "?" + authValues.Encode()

	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)
	server := &http.Server{Handler: callbackHandler(m.CallbackPath, state, codeCh, errCh)}
	go func() {
		if err := server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()
	defer server.Shutdown(context.Background())

	if m.OnAuthURL != nil {
		if err := m.OnAuthURL(ctx, authURL); err != nil {
			return err
		}
	} else {
		fmt.Println("Open this URL to sign in:", authURL)
	}

	waitCtx, cancel := context.WithTimeout(ctx, m.LocalhostTimeout)
	defer cancel()
	var code string
	select {
	case <-waitCtx.Done():
		return waitCtx.Err()
	case err := <-errCh:
		return err
	case code = <-codeCh:
	}

	form := url.Values{}
	form.Set("client_id", m.ClientID)
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", redirectURI)
	form.Set("code_verifier", verifier)
	form.Set("scope", m.Scope)

	var tok MSToken
	if err := m.postForm(ctx, msTokenURL, form, &tok); err != nil {
		return fmt.Errorf("exchange Microsoft authorization code: %w", err)
	}
	tok.Expiry = time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second)
	m.MSToken = &tok
	return m.save(ctx)
}

func callbackHandler(path, expectedState string, codeCh chan<- string, errCh chan<- error) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		gotState := q.Get("state")
		if subtle.ConstantTimeCompare([]byte(gotState), []byte(expectedState)) != 1 {
			http.Error(w, "invalid state", http.StatusBadRequest)
			select {
			case errCh <- errors.New("invalid OAuth state"):
			default:
			}
			return
		}
		if e := q.Get("error"); e != "" {
			desc := q.Get("error_description")
			http.Error(w, html.EscapeString(e), http.StatusBadRequest)
			select {
			case errCh <- fmt.Errorf("Microsoft OAuth error: %s: %s", e, desc):
			default:
			}
			return
		}
		code := q.Get("code")
		if code == "" {
			http.Error(w, "missing code", http.StatusBadRequest)
			select {
			case errCh <- errors.New("missing OAuth code"):
			default:
			}
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, "<html><body>登入完成，您現在可以關閉此視窗。</body></html>")
		select {
		case codeCh <- code:
		default:
		}
	})
	return mux
}

func (m *Auth) loginByDeviceCode(ctx context.Context) error {
	form := url.Values{}
	form.Set("client_id", m.ClientID)
	form.Set("scope", m.Scope)

	var dc MSDeviceCode
	if err := m.postForm(ctx, msDeviceCodeURL, form, &dc); err != nil {
		return fmt.Errorf("request Microsoft device code: %w", err)
	}
	if dc.VerificationURI == "" {
		dc.VerificationURI = dc.VerificationURL
	}
	if m.OnCode != nil {
		if err := m.OnCode(ctx, dc); err != nil {
			return err
		}
	} else {
		fmt.Println(dc.Message)
	}

	interval := dc.Interval
	if interval <= 0 {
		interval = 5
	}
	deadline := time.Now().Add(time.Duration(dc.ExpiresIn) * time.Second)
	for {
		if !deadline.IsZero() && time.Now().After(deadline) {
			return errors.New("device code expired")
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(interval) * time.Second):
		}

		form = url.Values{}
		form.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
		form.Set("client_id", m.ClientID)
		form.Set("device_code", dc.DeviceCode)

		var tok MSToken
		err := m.postForm(ctx, msTokenURL, form, &tok)
		if err == nil {
			tok.Expiry = time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second)
			m.MSToken = &tok
			return m.save(ctx)
		}
		var oauthErr *oauthError
		if !errors.As(err, &oauthErr) {
			return err
		}
		switch oauthErr.Code {
		case "authorization_pending":
			continue
		case "slow_down":
			interval++
			continue
		default:
			return oauthErr
		}
	}
}

func (m *Auth) refreshMSToken(ctx context.Context) error {
	form := url.Values{}
	form.Set("client_id", m.ClientID)
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", m.MSToken.RefreshToken)
	form.Set("scope", m.Scope)

	oldRefresh := m.MSToken.RefreshToken
	var tok MSToken
	if err := m.postForm(ctx, msTokenURL, form, &tok); err != nil {
		return fmt.Errorf("refresh Microsoft token: %w", err)
	}
	if tok.RefreshToken == "" {
		tok.RefreshToken = oldRefresh
	}
	tok.Expiry = time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second)
	m.MSToken = &tok
	return nil
}

func (m *Auth) authXBL(ctx context.Context) error {
	payload := map[string]any{
		"Properties": map[string]any{
			"AuthMethod": "RPS",
			"SiteName":   "user.auth.xboxlive.com",
			"RpsTicket":  "d=" + m.MSToken.AccessToken,
		},
		"RelyingParty": "http://auth.xboxlive.com",
		"TokenType":    "JWT",
	}
	var token XBToken
	if err := m.postJSON(ctx, xblAuthURL, payload, &token); err != nil {
		return fmt.Errorf("obtain XBL token: %w", err)
	}
	m.XBLToken = &token
	return m.save(ctx)
}

func (m *Auth) authXSTS(ctx context.Context) error {
	payload := map[string]any{
		"Properties": map[string]any{
			"SandboxId":  "RETAIL",
			"UserTokens": []string{m.XBLToken.Token},
		},
		"RelyingParty": "rp://api.minecraftservices.com/",
		"TokenType":    "JWT",
	}
	var token XBToken
	if err := m.postJSON(ctx, xstsAuthURL, payload, &token); err != nil {
		return fmt.Errorf("obtain XSTS token: %w", err)
	}
	m.XSTSToken = &token
	return m.save(ctx)
}

func (m *Auth) authMinecraft(ctx context.Context) error {
	if m.XSTSToken == nil || len(m.XSTSToken.DisplayClaims.Xui) == 0 {
		return errors.New("missing XSTS xui claim")
	}
	identityToken := "XBL3.0 x=" + m.XSTSToken.DisplayClaims.Xui[0].Uhs + ";" + m.XSTSToken.Token
	payload := map[string]string{"identityToken": identityToken}

	var token MCToken
	if err := m.postJSON(ctx, mcXboxLoginURL, payload, &token); err != nil {
		return fmt.Errorf("obtain Minecraft token: %w", err)
	}
	token.Expiry = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	m.MCToken = &token
	return m.save(ctx)
}

func (m *Auth) fetchMCProfile(ctx context.Context) error {
	var profile MCProfile
	if err := m.getJSON(ctx, mcProfileURL, "Bearer "+m.MCToken.AccessToken, &profile); err != nil {
		return fmt.Errorf("fetch Minecraft profile: %w", err)
	}
	id, err := uuidFromUndashed(profile.ID)
	if err != nil {
		return err
	}
	m.MCProfile = &profile
	m.profile = &auth.Profile{Name: profile.Name, UUID: id}
	return m.save(ctx)
}

func (m *Auth) postJSON(ctx context.Context, endpoint string, payload any, out any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return m.doJSON(req, out)
}

func (m *Auth) postForm(ctx context.Context, endpoint string, form url.Values, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return m.doJSON(req, out)
}

func (m *Auth) getJSON(ctx context.Context, endpoint, authorization string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	return m.doJSON(req, out)
}

func (m *Auth) doJSON(req *http.Request, out any) error {
	resp, err := m.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var oe oauthError
		if json.Unmarshal(body, &oe) == nil && oe.Code != "" {
			return &oe
		}
		return &httpStatusError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(body)}
	}
	if out == nil || len(body) == 0 {
		return nil
	}
	return json.Unmarshal(body, out)
}

type httpStatusError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *httpStatusError) Error() string {
	return e.Status + ": " + e.Body
}

type oauthError struct {
	Code        string `json:"error"`
	Description string `json:"error_description"`
	Codes       []int  `json:"error_codes"`
	Timestamp   string `json:"timestamp"`
	TraceID     string `json:"trace_id"`
	Correlation string `json:"correlation_id"`
}

func (e *oauthError) Error() string {
	if e.Description != "" {
		return e.Code + ": " + e.Description
	}
	return e.Code
}

func randomURLSafe(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func pkceChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func uuidFromUndashed(s string) (uuid.UUID, error) {
	if strings.Contains(s, "-") {
		return uuid.Parse(s)
	}
	if len(s) != 32 {
		return uuid.Nil, fmt.Errorf("invalid Minecraft UUID length %d", len(s))
	}
	b := make([]byte, 16)
	for i := range 16 {
		v, err := strconv.ParseUint(s[i*2:i*2+2], 16, 8)
		if err != nil {
			return uuid.Nil, err
		}
		b[i] = byte(v)
	}
	var id uuid.UUID
	copy(id[:], b)
	return id, nil
}

var _ auth.Provider = (*Auth)(nil)
