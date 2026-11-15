package bot

import (
	"context"
	"fmt"

	"golang.org/x/net/proxy"

	mcnet "github.com/KonjacBot/go-mc/net"
)

// Dialer returns a Minecraft dialer that connects through this SOCKS5 proxy.
// It lives next to ProxyConfig so the client and the keeper dial the same way.
func (p *ProxyConfig) Dialer() (mcnet.MCDialer, error) {
	if p.Type != "" && p.Type != "socks5" {
		return nil, fmt.Errorf("unsupported proxy type %q", p.Type)
	}
	var auth *proxy.Auth
	if p.Username != "" || p.Password != "" {
		auth = &proxy.Auth{
			User:     p.Username,
			Password: p.Password,
		}
	}

	dialer, err := proxy.SOCKS5("tcp", p.Host, auth, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("failed to create SOCKS5 proxy dialer: %w", err)
	}

	return &socks5MCDialer{
		dialer: dialer,
	}, nil
}

// socks5MCDialer 實作 mcnet.MCDialer 介面，使用 SOCKS5 proxy
type socks5MCDialer struct {
	dialer proxy.Dialer
}

func (d *socks5MCDialer) DialMCContext(ctx context.Context, address string) (*mcnet.Conn, error) {
	contextDialer, ok := d.dialer.(proxy.ContextDialer)
	if !ok {
		return nil, fmt.Errorf("SOCKS5 dialer does not support contexts")
	}
	conn, err := contextDialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to dial through SOCKS5 proxy: %w", err)
	}

	// 將 net.Conn 包裝成 mcnet.Conn
	return mcnet.WrapConn(conn), nil
}
