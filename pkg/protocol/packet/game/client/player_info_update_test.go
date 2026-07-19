package client

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"testing"
	"time"

	"github.com/KonjacBot/go-mc/chat/sign"
	pk "github.com/KonjacBot/go-mc/net/packet"
	"github.com/KonjacBot/go-mc/yggdrasil/user"
	"github.com/google/uuid"
)

// Wire layout for Initialize Chat (wiki.vg / Java ClientboundPlayerInfoUpdatePacket):
//
//	bool present
//	if present:
//	  UUID chatSessionId
//	  long expiresAt
//	  byte[] publicKey (PKIX)
//	  byte[] keySignature
//
// sign.Session already encodes SessionID + PublicKey, so PlayerInfoChatData must
// not prepend another UUID (that misaligned PublicKey into ASN.1 parse).
func TestPlayerInfoUpdateInitializeChatRoundTrip(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	playerID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	sessionID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	sig := make([]byte, 32)
	_, _ = rand.Read(sig)

	want := &PlayerInfoUpdate{
		Players: map[uuid.UUID][]PlayerInfo{
			playerID: {
				&PlayerInfoInitializeChat{
					Data: pk.Option[PlayerInfoChatData, *PlayerInfoChatData]{
						Has: true,
						Val: PlayerInfoChatData{
							Session: sign.Session{
								SessionID: sessionID,
								PublicKey: user.PublicKey{
									ExpiresAt: time.UnixMilli(1_700_000_000_000),
									PubKey:    &priv.PublicKey,
									Signature: sig,
								},
							},
						},
					},
				},
			},
		},
	}

	var encoded bytes.Buffer
	if _, err := want.WriteTo(&encoded); err != nil {
		t.Fatalf("WriteTo: %v", err)
	}

	var got PlayerInfoUpdate
	if _, err := got.ReadFrom(bytes.NewReader(encoded.Bytes())); err != nil {
		t.Fatalf("ReadFrom: %v", err)
	}

	infos, ok := got.Players[playerID]
	if !ok || len(infos) != 1 {
		t.Fatalf("players = %#v", got.Players)
	}
	initChat, ok := infos[0].(*PlayerInfoInitializeChat)
	if !ok || !bool(initChat.Data.Has) {
		t.Fatalf("info = %#v", infos[0])
	}
	if initChat.Data.Val.Session.SessionID != sessionID {
		t.Fatalf("session id = %v, want %v", initChat.Data.Val.Session.SessionID, sessionID)
	}
	if !initChat.Data.Val.Session.PublicKey.PubKey.Equal(&priv.PublicKey) {
		t.Fatal("public key mismatch after round-trip")
	}
}

func TestPlayerInfoUpdateInitializeChatWireDecode(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	pubDER, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	playerID := uuid.MustParse("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	sessionID := uuid.MustParse("ffffffff-0000-1111-2222-333333333333")
	sig := bytes.Repeat([]byte{0xab}, 16)

	var buf bytes.Buffer
	// action bitset: only INITIALIZE_CHAT (bit 1 / mask 0x02)
	bitset := pk.NewFixedBitSet(8)
	bitset.Set(1, true)
	if _, err := bitset.WriteTo(&buf); err != nil {
		t.Fatal(err)
	}
	if _, err := pk.VarInt(1).WriteTo(&buf); err != nil {
		t.Fatal(err)
	}
	if _, err := pk.UUID(playerID).WriteTo(&buf); err != nil {
		t.Fatal(err)
	}
	// present = true
	if _, err := pk.Boolean(true).WriteTo(&buf); err != nil {
		t.Fatal(err)
	}
	if _, err := pk.UUID(sessionID).WriteTo(&buf); err != nil {
		t.Fatal(err)
	}
	if _, err := pk.Long(1_700_000_000_000).WriteTo(&buf); err != nil {
		t.Fatal(err)
	}
	if _, err := pk.ByteArray(pubDER).WriteTo(&buf); err != nil {
		t.Fatal(err)
	}
	if _, err := pk.ByteArray(sig).WriteTo(&buf); err != nil {
		t.Fatal(err)
	}

	var got PlayerInfoUpdate
	n, err := got.ReadFrom(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("ReadFrom: %v", err)
	}
	if n != int64(buf.Len()) {
		t.Fatalf("consumed %d of %d bytes", n, buf.Len())
	}
	infos := got.Players[playerID]
	initChat := infos[0].(*PlayerInfoInitializeChat)
	if initChat.Data.Val.Session.SessionID != sessionID {
		t.Fatalf("session id = %v", initChat.Data.Val.Session.SessionID)
	}
}
