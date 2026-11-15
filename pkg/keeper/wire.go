package keeper

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/KonjacBot/minego/pkg/bot"
)

// Version is the local interface version. A core and a keeper with different
// versions refuse each other; the launcher replaces the keeper when it sees a
// mismatch.
const Version = 1

// maxFrame bounds a request or reply so a stray client cannot make the keeper
// allocate without limit before it is authenticated.
const maxFrame = 1 << 20

// Op is what a request asks the keeper to do.
type Op string

const (
	// OpConnect attaches to the named session, opening the server connection
	// first when there is none. The socket then carries the Minecraft protocol.
	OpConnect Op = "connect"
	// OpClose ends the named session (the player leaves the server).
	OpClose Op = "close"
	// OpStatus lists the sessions, for the launcher to decide whether the keeper
	// can be restarted.
	OpStatus Op = "status"
	// OpShutdown ends every session and stops the keeper (the launcher's way to
	// stop it on every platform when the user quits).
	OpShutdown Op = "shutdown"
)

// Request is the first frame a client sends on the local socket.
type Request struct {
	Keeper  int    `json:"keeper"`
	Token   string `json:"token"`
	Op      Op     `json:"op"`
	Session string `json:"session,omitempty"`
	// Protocol is the Minecraft protocol version the client speaks.
	Protocol int32 `json:"protocol,omitempty"`
	// Server is how to open the connection when the session does not exist.
	Server *Server `json:"server,omitempty"`
}

// Server is the server connection a session opens.
type Server struct {
	Address  string           `json:"address"`
	FakeHost string           `json:"fakeHost,omitempty"`
	Proxy    *bot.ProxyConfig `json:"proxy,omitempty"`
	Name     string           `json:"name"`
	UUID     uuid.UUID        `json:"uuid"`
	// AccessToken is the Minecraft session token; empty logs in offline. The
	// keeper uses it once for the session join and does not store it.
	AccessToken string `json:"accessToken,omitempty"`
}

// Reply answers a Request.
type Reply struct {
	Keeper   int    `json:"keeper"`
	Protocol int32  `json:"protocol"`
	OK       bool   `json:"ok"`
	Error    string `json:"error,omitempty"`
	// Resumed reports that the session already existed: the snapshot that
	// follows describes a player that never left the server.
	Resumed bool `json:"resumed,omitempty"`
	// Previous is why the session's previous connection ended while no core was
	// attached (for example a kick), so the core can report it.
	Previous string        `json:"previous,omitempty"`
	Sessions []SessionInfo `json:"sessions,omitempty"`
}

// SessionInfo describes one kept connection in a status reply.
type SessionInfo struct {
	Session  string    `json:"session"`
	Address  string    `json:"address"`
	Attached bool      `json:"attached"`
	Since    time.Time `json:"since"`
	// Detached is when the last core left; zero while one is attached.
	Detached time.Time `json:"detached,omitzero"`
}

// writeFrame writes v as a length-prefixed JSON frame. Frames are length
// prefixed rather than newline delimited so reading one never buffers the
// Minecraft bytes that follow it.
func writeFrame(w io.Writer, v any) error {
	body, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("encode frame: %w", err)
	}
	if len(body) > maxFrame {
		return fmt.Errorf("frame of %d bytes exceeds %d", len(body), maxFrame)
	}
	buf := make([]byte, 4+len(body))
	binary.BigEndian.PutUint32(buf, uint32(len(body)))
	copy(buf[4:], body)
	if _, err := w.Write(buf); err != nil {
		return fmt.Errorf("write frame: %w", err)
	}
	return nil
}

func readFrame(r io.Reader, v any) error {
	var size [4]byte
	if _, err := io.ReadFull(r, size[:]); err != nil {
		return fmt.Errorf("read frame size: %w", err)
	}
	n := binary.BigEndian.Uint32(size[:])
	if n > maxFrame {
		return fmt.Errorf("frame of %d bytes exceeds %d", n, maxFrame)
	}
	body := make([]byte, n)
	if _, err := io.ReadFull(r, body); err != nil {
		return fmt.Errorf("read frame: %w", err)
	}
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("decode frame: %w", err)
	}
	return nil
}

// LoadOrCreateToken returns the token stored at path, creating a random one
// readable only by the current user when the file is missing or empty. The
// keeper and the core both call it, so whichever starts first writes it.
func LoadOrCreateToken(path string) (string, error) {
	if token, err := ReadToken(path); err == nil {
		return token, nil
	} else if !errors.Is(err, os.ErrNotExist) && !errors.Is(err, errEmptyToken) {
		return "", err
	}
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return "", fmt.Errorf("generate keeper token: %w", err)
	}
	token := hex.EncodeToString(secret)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return "", fmt.Errorf("create keeper token directory: %w", err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".keeper-token-*")
	if err != nil {
		return "", fmt.Errorf("create keeper token: %w", err)
	}
	defer os.Remove(tmp.Name())
	if err := tmp.Chmod(0o600); err != nil && !errors.Is(err, errors.ErrUnsupported) {
		_ = tmp.Close()
		return "", fmt.Errorf("protect keeper token: %w", err)
	}
	if _, err := tmp.WriteString(token); err != nil {
		_ = tmp.Close()
		return "", fmt.Errorf("write keeper token: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return "", fmt.Errorf("write keeper token: %w", err)
	}
	// Link instead of rename: when two processes race, the first token wins and
	// the second reads it back instead of overwriting it.
	if err := os.Link(tmp.Name(), path); err != nil {
		if token, readErr := ReadToken(path); readErr == nil {
			return token, nil
		}
		if err := os.Rename(tmp.Name(), path); err != nil {
			return "", fmt.Errorf("install keeper token: %w", err)
		}
	}
	return token, nil
}

var errEmptyToken = errors.New("keeper token file is empty")

// ReadToken reads the token the keeper wrote at path.
func ReadToken(path string) (string, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read keeper token: %w", err)
	}
	token := strings.TrimSpace(string(body))
	if token == "" {
		return "", errEmptyToken
	}
	return token, nil
}

func tokenMatches(want, got string) bool {
	return want != "" && subtle.ConstantTimeCompare([]byte(want), []byte(got)) == 1
}
