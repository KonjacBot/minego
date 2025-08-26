package client

import (
	"context"
	"fmt"

	"golang.org/x/net/proxy"

	mcnet "git.konjactw.dev/falloutBot/go-mc/net"

	"git.konjactw.dev/patyhank/minego/pkg/bot"
)

// createSOCKS5Dialer 建立 SOCKS5 dialer
func socks5(proxyConfig *bot.ProxyConfig) (mcnet.MCDialer, error) {
	var auth *proxy.Auth
	if proxyConfig.Username != "" || proxyConfig.Password != "" {
		auth = &proxy.Auth{
			User:     proxyConfig.Username,
			Password: proxyConfig.Password,
		}
	}

	dialer, err := proxy.SOCKS5("tcp", proxyConfig.Host, auth, proxy.Direct)
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
	// 使用 SOCKS5 proxy 建立連線
	conn, err := d.dialer.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to dial through SOCKS5 proxy: %w", err)
	}

	// 將 net.Conn 包裝成 mcnet.Conn
	return mcnet.WrapConn(conn), nil
}
