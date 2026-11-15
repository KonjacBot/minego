package keeper

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/KonjacBot/minego/v2/pkg/auth"
	"github.com/KonjacBot/minego/v2/pkg/bot"
	"github.com/KonjacBot/minego/v2/pkg/protocol"
)

// Dial asks the keeper for the session in opts and returns the socket ready
// for the Minecraft handshake. provider supplies the profile and, for online
// accounts, the access token the keeper logs in with when it has to open the
// connection.
func Dial(ctx context.Context, opts bot.KeeperOptions, addr string, connect *bot.ConnectOptions, provider auth.Provider) (net.Conn, Reply, error) {
	if opts.Socket == "" || opts.Session == "" {
		return nil, Reply{}, errors.New("保持者設定少了 socket 或連線名稱")
	}
	profile := provider.FetchProfile(ctx)
	if profile == nil {
		return nil, Reply{}, errors.New("登入提供者沒有回傳玩家資料")
	}
	srv := &Server{Address: addr, Name: profile.Name, UUID: profile.UUID}
	if connect != nil {
		srv.FakeHost, srv.Proxy = connect.FakeHost, connect.Proxy
	}
	switch p := provider.(type) {
	case *auth.OfflineAuth:
	case auth.TokenSource:
		token, err := p.SessionToken(ctx)
		if err != nil {
			return nil, Reply{}, fmt.Errorf("取得登入權杖: %w", err)
		}
		srv.AccessToken = token
	default:
		return nil, Reply{}, fmt.Errorf("登入方式 %T 不能經由保持者連線", provider)
	}
	req := Request{Keeper: Version, Token: opts.Token, Op: OpConnect, Session: opts.Session, Protocol: protocol.ProtocolVersion, Server: srv}
	conn, reply, err := roundTrip(ctx, opts.Socket, req)
	if err != nil {
		return nil, reply, err
	}
	return conn, reply, nil
}

// CloseSession ends the named session: the player leaves the server.
func CloseSession(ctx context.Context, opts bot.KeeperOptions) error {
	conn, _, err := roundTrip(ctx, opts.Socket, Request{Keeper: Version, Token: opts.Token, Op: OpClose, Session: opts.Session})
	if err != nil {
		return err
	}
	return conn.Close()
}

// Status lists the keeper's sessions.
func Status(ctx context.Context, socket, token string) (Reply, error) {
	conn, reply, err := roundTrip(ctx, socket, Request{Keeper: Version, Token: token, Op: OpStatus})
	if err != nil {
		return reply, err
	}
	_ = conn.Close()
	return reply, nil
}

func roundTrip(ctx context.Context, socket string, req Request) (net.Conn, Reply, error) {
	conn, err := (&net.Dialer{}).DialContext(ctx, "unix", socket)
	if err != nil {
		return nil, Reply{}, fmt.Errorf("連到保持者: %w", err)
	}
	stop := context.AfterFunc(ctx, func() { _ = conn.SetDeadline(time.Now()) })
	defer stop()
	var reply Reply
	err = writeFrame(conn, req)
	if err == nil {
		err = readFrame(conn, &reply)
	}
	if err == nil && !reply.OK {
		err = &RefusedError{Reason: reply.Error}
	}
	if err != nil {
		_ = conn.Close()
		if ctx.Err() != nil {
			err = fmt.Errorf("%w: %w", err, ctx.Err())
		}
		return nil, reply, fmt.Errorf("保持者: %w", err)
	}
	_ = conn.SetDeadline(time.Time{})
	return conn, reply, nil
}

// RefusedError is the keeper declining a request, with its reason.
type RefusedError struct{ Reason string }

func (e *RefusedError) Error() string { return "保持者拒絕: " + e.Reason }
