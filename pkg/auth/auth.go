package auth

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"git.konjactw.dev/falloutBot/go-mc/chat"
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	"git.konjactw.dev/falloutBot/go-mc/net"
	"git.konjactw.dev/falloutBot/go-mc/net/CFB8"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/login/client"
	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/login/server"
)

var (
	ErrLogin   = errors.New("login")
	ErrKick    = errors.New("login.kicked")
	ErrEncrypt = errors.New("login.encrypt")
)

type Profile struct {
	Name string
	UUID uuid.UUID
}

type Provider interface {
	Authenticate(ctx context.Context, conn *net.Conn, content client.LoginHello) error
	FetchProfile(ctx context.Context) *Profile
}

type Auth struct {
	*net.Conn
	Provider
}

func (a *Auth) HandleLogin(ctx context.Context) error {
	profile := a.FetchProfile(ctx)

	err := a.WritePacket(pk.Marshal(packetid.ServerboundLoginHello, server.LoginHello{
		Name: profile.Name,
		UUID: profile.UUID,
	}))
	if err != nil {
		return errors.Join(ErrLogin, fmt.Errorf("write login hello fail: %w", err))
	}

	var p pk.Packet
	for {
		err = a.ReadPacket(&p)
		if err != nil {
			return errors.Join(ErrLogin, fmt.Errorf("read packet fail: %w", err))
		}

		switch packetid.ClientboundPacketID(p.ID) {
		case packetid.ClientboundLoginLoginDisconnect:
			var reason chat.JsonMessage
			err = p.Scan(&reason)

			return errors.Join(ErrKick, fmt.Errorf("kicked by server: %s", chat.Message(reason).ClearString()))
		case packetid.ClientboundLoginHello:
			var hello client.LoginHello
			err = p.Scan(&hello)
			if err != nil {
				return errors.Join(ErrLogin, fmt.Errorf("read login hello fail: %w", err))
			}

			err = a.Authenticate(ctx, a.Conn, hello)
			if err != nil {
				return errors.Join(ErrLogin, fmt.Errorf("authenticate fail: %w", err))
			}
		case packetid.ClientboundLoginLoginFinished:
			err = a.WritePacket(pk.Marshal(packetid.ServerboundLoginLoginAcknowledged))
			if err != nil {
				return errors.Join(ErrLogin, fmt.Errorf("write login ack fail: %w", err))
			}
			return nil
		case packetid.ClientboundLoginLoginCompression:
			var threshold int32

			err = p.Scan((*pk.VarInt)(&threshold))
			if err != nil {
				return errors.Join(ErrLogin, fmt.Errorf("read login compression fail: %w", err))
			}
			a.Conn.SetThreshold(int(threshold))
		case packetid.ClientboundLoginCustomQuery:
			var query client.LoginCustomQuery

			err = p.Scan(&query)
			if err != nil {
				return errors.Join(ErrLogin, fmt.Errorf("read login custom query fail: %w", err))
			}

			err = a.WritePacket(pk.Marshal(
				packetid.ServerboundLoginCustomQueryAnswer,
				&server.LoginCustomQueryAnswer{MessageID: query.MessageID},
			))
			if err != nil {
				return errors.Join(ErrLogin, fmt.Errorf("read login custom query fail: %w", err))
			}
		case packetid.ClientboundLoginCookieRequest:
			var cookie client.LoginCookieRequest
			err = p.Scan(&cookie)
			if err != nil {
				return errors.Join(ErrLogin, fmt.Errorf("read login cookie request fail: %w", err))
			}
			err = a.WritePacket(pk.Marshal(
				packetid.ServerboundLoginCookieResponse,
				&server.LoginCookieResponse{Key: cookie.Key},
			))
			if err != nil {
				return errors.Join(ErrLogin, fmt.Errorf("read login cookie request fail: %w", err))
			}
		}
	}
}

type OnlineAuthServer struct {
	SessionServer string
	AuthServer    string
}

type OnlineAuth struct {
	AccessToken string
	Profile     Profile
}

func (o *OnlineAuth) Authenticate(ctx context.Context, conn *net.Conn, content client.LoginHello) error {
	key, encodeStream, decodeStream := newSymmetricEncryption()

	err := o.LoginAuth(ctx, content, key)
	if err != nil {
		return errors.Join(ErrEncrypt, fmt.Errorf("login auth fail: %w", err))
	}

	// Response with Encryption Key
	var pkt pk.Packet
	pkt, err = genEncryptionKeyResponse(key, content.PublicKey, content.VerifyToken)
	if err != nil {
		return fmt.Errorf("gen encryption key response fail: %v", err)
	}

	err = conn.WritePacket(pkt)
	if err != nil {
		return err
	}

	// Set Connection Encryption
	conn.SetCipher(encodeStream, decodeStream)
	return nil
}

func genEncryptionKeyResponse(shareSecret, publicKey, verifyToken []byte) (erp pk.Packet, err error) {
	iPK, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		err = fmt.Errorf("decode public key fail: %v", err)
		return
	}
	rsaKey := iPK.(*rsa.PublicKey)
	cryptPK, err := rsa.EncryptPKCS1v15(rand.Reader, rsaKey, shareSecret)
	if err != nil {
		err = fmt.Errorf("encryption share secret fail: %v", err)
		return
	}

	verifyT, err := rsa.EncryptPKCS1v15(rand.Reader, rsaKey, verifyToken)
	if err != nil {
		err = fmt.Errorf("encryption verfy tokenfail: %v", err)
		return erp, err
	}
	return pk.Marshal(
		packetid.ServerboundLoginKey,
		pk.ByteArray(cryptPK),
		pk.ByteArray(verifyT),
	), nil
}

func (o *OnlineAuth) LoginAuth(ctx context.Context, content client.LoginHello, key []byte) error {
	digest := authDigest(content.ServerID, key, content.PublicKey)

	request, err := json.Marshal(struct {
		AccessToken     string `json:"accessToken"`
		SelectedProfile string `json:"selectedProfile"`
		ServerID        string `json:"serverId"`
	}{
		AccessToken:     o.AccessToken,
		SelectedProfile: hex.EncodeToString(o.Profile.UUID[:]),
		ServerID:        digest,
	})

	PostRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://sessionserver.mojang.com/session/minecraft/join",
		bytes.NewReader(request))
	if err != nil {
		return errors.Join(ErrEncrypt, fmt.Errorf("make request fail: %w", err))
	}
	PostRequest.Header.Set("User-agent", "go-mc")
	PostRequest.Header.Set("Connection", "keep-alive")
	PostRequest.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(PostRequest)
	if err != nil {
		return errors.Join(ErrEncrypt, fmt.Errorf("session mojang fail: %w", err))
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusNoContent {
		return errors.Join(ErrEncrypt, fmt.Errorf("session join fail: %s", string(body)))
	}
	return nil
}

func newSymmetricEncryption() (key []byte, encoStream, decoStream cipher.Stream) {
	key = make([]byte, 16)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}

	b, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	decoStream = CFB8.NewCFB8Decrypt(b, key)
	encoStream = CFB8.NewCFB8Encrypt(b, key)
	return
}

func (o *OnlineAuth) FetchProfile(ctx context.Context) *Profile {
	return &o.Profile
}

func authDigest(serverID string, sharedSecret, publicKey []byte) string {
	h := sha1.New()
	h.Write([]byte(serverID))
	h.Write(sharedSecret)
	h.Write(publicKey)
	hash := h.Sum(nil)

	// Check for negative hashes
	negative := (hash[0] & 0x80) == 0x80
	if negative {
		hash = twosComplement(hash)
	}

	// Trim away zeroes
	res := strings.TrimLeft(hex.EncodeToString(hash), "0")
	if negative {
		res = "-" + res
	}

	return res
}

// little endian
func twosComplement(p []byte) []byte {
	carry := true
	for i := len(p) - 1; i >= 0; i-- {
		p[i] = ^p[i]
		if carry {
			carry = p[i] == 0xff
			p[i]++
		}
	}
	return p
}

type OfflineAuth struct {
	Username string
}

// NameToUUID return the UUID from player name in offline mode
func NameToUUID(name string) uuid.UUID {
	version := 3
	h := md5.New()
	h.Write([]byte("OfflinePlayer:"))
	h.Write([]byte(name))
	var id uuid.UUID
	h.Sum(id[:0])
	id[6] = (id[6] & 0x0f) | uint8((version&0xf)<<4)
	id[8] = (id[8] & 0x3f) | 0x80 // RFC 4122 variant
	return id
}

func (o *OfflineAuth) FetchProfile(ctx context.Context) *Profile {
	return &Profile{
		Name: o.Username,
		UUID: NameToUUID(o.Username),
	}
}

func (o *OfflineAuth) Authenticate(ctx context.Context, conn *net.Conn, content client.LoginHello) error {
	return nil
}

type KonjacAuth struct {
	UserCode string
}

func (k *KonjacAuth) Authenticate(ctx context.Context, conn *net.Conn, content client.LoginHello) error {
	key, encodeStream, decodeStream := newSymmetricEncryption()

	err := k.LoginAuth(ctx, content, key)
	if err != nil {
		return errors.Join(ErrEncrypt, fmt.Errorf("login auth fail: %w", err))
	}

	// Response with Encryption Key
	var pkt pk.Packet
	pkt, err = genEncryptionKeyResponse(key, content.PublicKey, content.VerifyToken)
	if err != nil {
		return fmt.Errorf("gen encryption key response fail: %v", err)
	}

	err = conn.WritePacket(pkt)
	if err != nil {
		return err
	}

	// Set Connection Encryption
	conn.SetCipher(encodeStream, decodeStream)
	return nil
}

func (k *KonjacAuth) LoginAuth(ctx context.Context, content client.LoginHello, key []byte) error {
	digest := authDigest(content.ServerID, key, content.PublicKey)

	request, err := json.Marshal(struct {
		AccessToken     string `json:"accessToken"`
		SelectedProfile string `json:"selectedProfile"`
		ServerID        string `json:"serverId"`
	}{
		AccessToken:     k.UserCode,
		SelectedProfile: "-",
		ServerID:        digest,
	})

	PostRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://127.0.0.1:37565/ss/session/minecraft/join",
		bytes.NewReader(request))
	if err != nil {
		return errors.Join(ErrEncrypt, fmt.Errorf("make request fail: %w", err))
	}
	PostRequest.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(PostRequest)
	if err != nil {
		return errors.Join(ErrEncrypt, fmt.Errorf("session mojang fail: %w", err))
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusNoContent {
		return errors.Join(ErrEncrypt, fmt.Errorf("session join fail: %s", string(body)))
	}
	return nil
}

func (k *KonjacAuth) FetchProfile(ctx context.Context) *Profile {
	data, err := json.Marshal(map[string]any{
		"agent": map[string]any{
			"name":    "minego",
			"version": 0,
		},
		"username":    k.UserCode,
		"password":    "",
		"clientToken": "",
		"requestUser": "",
	})
	if err != nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://127.0.0.1:37565/as/authenticate", bytes.NewReader(data))
	if err != nil {
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil
	}

	var profile struct {
		SelectedProfile Profile `json:"selectedProfile"`
	}

	err = json.Unmarshal(data, &profile)
	if err != nil {
		return nil
	}

	return &profile.SelectedProfile
}
