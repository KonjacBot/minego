package keeper

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

// checkPeer refuses connections from other users: the socket's permissions
// already stop them, this guards against a socket file with loose permissions.
func checkPeer(conn net.Conn) error {
	unix, ok := conn.(*net.UnixConn)
	if !ok {
		return nil
	}
	raw, err := unix.SyscallConn()
	if err != nil {
		return fmt.Errorf("inspect keeper peer: %w", err)
	}
	var cred *syscall.Ucred
	var credErr error
	if err := raw.Control(func(fd uintptr) {
		cred, credErr = syscall.GetsockoptUcred(int(fd), syscall.SOL_SOCKET, syscall.SO_PEERCRED)
	}); err != nil {
		return fmt.Errorf("inspect keeper peer: %w", err)
	}
	if credErr != nil {
		return fmt.Errorf("inspect keeper peer: %w", credErr)
	}
	if int(cred.Uid) != os.Getuid() {
		return fmt.Errorf("keeper peer runs as uid %d, not %d", cred.Uid, os.Getuid())
	}
	return nil
}
