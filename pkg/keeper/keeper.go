// Package keeper keeps Minecraft connections alive while the program driving
// them restarts.
//
// A keeper is a small long-running process that owns the server connections
// (TCP, encryption, compression, login). The driving program — the core —
// attaches over a local socket and speaks the ordinary Minecraft protocol to
// it, as if the keeper were the server. While no core is attached the keeper
// answers what the server requires (keepalives, pings, teleport confirmations,
// chunk batch acknowledgements, resource packs, configuration phases) and
// keeps recording the world. A core that attaches gets a snapshot of that
// state first, then the live stream, and ends up in the state of a client that
// never left. This is the approach of proxies like ZenithProxy.
//
// The core side is minego's client with bot.ConnectOptions.Keeper set.
package keeper

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sort"
	"sync"
	"time"

	mcnet "github.com/KonjacBot/go-mc/net"

	"github.com/KonjacBot/minego/v2/pkg/protocol"
)

// Options configures a Keeper.
type Options struct {
	// Token is the shared secret cores must present (see LoadOrCreateToken).
	Token string
	// DetachTimeout closes a session no core attached to for this long; zero
	// means 10 minutes, negative never.
	DetachTimeout time.Duration
	// ReadIdle is how long the server or a core may stay silent; zero means 30s.
	ReadIdle time.Duration
	// OpenTimeout bounds connecting and logging in to a server; zero means 60s.
	OpenTimeout time.Duration
	Logger      *slog.Logger
}

// Keeper holds the sessions and serves cores on a local socket.
type Keeper struct {
	token         string
	detachTimeout time.Duration
	readIdle      time.Duration
	openTimeout   time.Duration
	log           *slog.Logger

	mu       sync.Mutex
	sessions map[string]*session
	// previous is why a session ended while no core was attached, told to the
	// next core that connects with that name.
	previous map[string]string

	shutdown     chan struct{}
	shutdownOnce sync.Once
}

// New returns a Keeper; call Serve to accept cores.
func New(opts Options) *Keeper {
	k := &Keeper{
		token:         opts.Token,
		detachTimeout: opts.DetachTimeout,
		readIdle:      opts.ReadIdle,
		openTimeout:   opts.OpenTimeout,
		log:           opts.Logger,
		sessions:      map[string]*session{},
		previous:      map[string]string{},
		shutdown:      make(chan struct{}),
	}
	if k.detachTimeout == 0 {
		k.detachTimeout = 10 * time.Minute
	}
	if k.readIdle <= 0 {
		k.readIdle = 30 * time.Second
	}
	if k.openTimeout <= 0 {
		k.openTimeout = time.Minute
	}
	if k.log == nil {
		k.log = slog.Default()
	}
	return k
}

// Serve accepts cores on ln until ctx ends or ln fails. The sessions keep
// running after Serve returns; Close ends them.
func (k *Keeper) Serve(ctx context.Context, ln net.Listener) error {
	stop := context.AfterFunc(ctx, func() { _ = ln.Close() })
	defer stop()
	for {
		conn, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return fmt.Errorf("接受核心連線: %w", err)
		}
		go k.handle(ctx, conn)
	}
}

// Shutdown is closed when a client asked the keeper to stop (OpShutdown); the
// program should then Close it and exit.
func (k *Keeper) Shutdown() <-chan struct{} { return k.shutdown }

// Close ends every session: the players leave their servers.
func (k *Keeper) Close() {
	k.mu.Lock()
	sessions := make([]*session, 0, len(k.sessions))
	for _, s := range k.sessions {
		sessions = append(sessions, s)
	}
	k.mu.Unlock()
	for _, s := range sessions {
		s.close(errors.New("保持者關閉"))
	}
	for _, s := range sessions {
		<-s.ended
	}
}

// Sessions describes the kept connections.
func (k *Keeper) Sessions() []SessionInfo {
	k.mu.Lock()
	sessions := make([]*session, 0, len(k.sessions))
	for _, s := range k.sessions {
		sessions = append(sessions, s)
	}
	k.mu.Unlock()
	out := make([]SessionInfo, 0, len(sessions))
	for _, s := range sessions {
		out = append(out, s.info())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Session < out[j].Session })
	return out
}

func (k *Keeper) handle(ctx context.Context, conn net.Conn) {
	defer func() { _ = conn.Close() }()
	if err := checkPeer(conn); err != nil {
		k.log.Warn("拒絕本機連線", "error", err)
		return
	}
	_ = conn.SetDeadline(time.Now().Add(10 * time.Second))
	var req Request
	if err := readFrame(conn, &req); err != nil {
		k.log.Debug("讀不到核心的請求", "error", err)
		return
	}
	reply := Reply{Keeper: Version, Protocol: protocol.ProtocolVersion}
	fail := func(format string, args ...any) {
		reply.Error = fmt.Sprintf(format, args...)
		_ = writeFrame(conn, reply)
	}
	if !tokenMatches(k.token, req.Token) {
		fail("保持者的通行碼不對")
		return
	}
	if req.Keeper != Version {
		fail("本機介面版本不同:核心 %d,保持者 %d", req.Keeper, Version)
		return
	}
	switch req.Op {
	case OpStatus:
		reply.OK, reply.Sessions = true, k.Sessions()
		_ = writeFrame(conn, reply)
	case OpShutdown:
		reply.OK = true
		_ = writeFrame(conn, reply)
		k.shutdownOnce.Do(func() { close(k.shutdown) })
	case OpClose:
		k.mu.Lock()
		s := k.sessions[req.Session]
		k.mu.Unlock()
		if s != nil {
			s.close(errors.New("核心要求結束連線"))
		}
		reply.OK = true
		_ = writeFrame(conn, reply)
	case OpConnect:
		if req.Protocol != protocol.ProtocolVersion {
			fail("Minecraft 協定版本不同:核心 %d,保持者 %d", req.Protocol, protocol.ProtocolVersion)
			return
		}
		s, resumed, previous, err := k.session(ctx, req)
		if err != nil {
			fail("%v", err)
			return
		}
		reply.OK, reply.Resumed, reply.Previous = true, resumed, previous
		if err := writeFrame(conn, reply); err != nil {
			return
		}
		_ = conn.SetDeadline(time.Time{})
		if err := s.serveCore(ctx, mcnet.WrapConn(conn)); err != nil {
			s.log.Info("核心離開", "reason", err)
		}
	default:
		fail("不認得的請求 %q", req.Op)
	}
}

// session returns the live session named in req, opening the server
// connection when there is none (or it points at another server).
func (k *Keeper) session(ctx context.Context, req Request) (s *session, resumed bool, previous string, err error) {
	if req.Session == "" || req.Server == nil || req.Server.Address == "" || req.Server.Name == "" {
		return nil, false, "", errors.New("請求少了連線名稱或伺服器")
	}
	k.mu.Lock()
	s = k.sessions[req.Session]
	if s != nil && s.srv.Address != req.Server.Address {
		k.mu.Unlock()
		s.close(fmt.Errorf("改連到 %s", req.Server.Address))
		<-s.ended
		k.mu.Lock()
		s = k.sessions[req.Session]
	}
	created := s == nil
	if created {
		s = newSession(k, req.Session, *req.Server)
		k.sessions[req.Session] = s
		previous = k.previous[req.Session]
		delete(k.previous, req.Session)
	}
	k.mu.Unlock()

	if created {
		k.log.Info("開新連線", "session", req.Session, "address", req.Server.Address)
		openCtx, cancel := context.WithTimeout(ctx, k.openTimeout)
		s.open(openCtx)
		cancel()
	}
	select {
	case <-s.ready:
	case <-ctx.Done():
		return nil, false, "", ctx.Err()
	}
	if s.openErr != nil {
		return nil, false, previous, s.openErr
	}
	if !s.alive() {
		return nil, false, previous, fmt.Errorf("連線已結束: %w", s.endReason())
	}
	return s, !created, previous, nil
}

// forget removes an ended session.
func (k *Keeper) forget(s *session, reason error, attached bool) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.sessions[s.name] != s {
		return
	}
	delete(k.sessions, s.name)
	if !attached && reason != nil && s.openErr == nil {
		k.previous[s.name] = reason.Error()
	}
}
