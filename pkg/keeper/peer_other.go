//go:build !linux

package keeper

import "net"

// checkPeer has no portable peer credential check outside Linux; the socket
// file's permissions and the token decide who may connect.
func checkPeer(net.Conn) error { return nil }
