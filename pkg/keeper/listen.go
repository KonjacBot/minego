package keeper

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

// ErrRunning is returned by Listen when another keeper already serves the socket.
var ErrRunning = errors.New("another keeper is already listening on this socket")

// Listen opens the keeper's local socket at path, readable only by the current
// user. A socket file left by a keeper that died is removed; one that still
// answers belongs to a running keeper and Listen reports ErrRunning.
//
// Windows 10 1803 and later support Unix domain sockets too, so every platform
// uses the same transport; the token (and on Linux the peer's uid) decides who
// may talk to the keeper.
func Listen(path string) (net.Listener, error) {
	if _, err := os.Stat(path); err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		conn, dialErr := (&net.Dialer{}).DialContext(ctx, "unix", path)
		cancel()
		if dialErr == nil {
			_ = conn.Close()
			return nil, ErrRunning
		}
		if err := os.Remove(path); err != nil {
			return nil, fmt.Errorf("remove stale keeper socket: %w", err)
		}
	}
	ln, err := net.Listen("unix", path)
	if err != nil {
		return nil, fmt.Errorf("listen on keeper socket: %w", err)
	}
	if err := os.Chmod(path, 0o600); err != nil && !errors.Is(err, errors.ErrUnsupported) {
		_ = ln.Close()
		return nil, fmt.Errorf("protect keeper socket: %w", err)
	}
	return ln, nil
}
