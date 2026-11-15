package keeper

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"github.com/KonjacBot/go-mc/chat"
	"github.com/KonjacBot/go-mc/data/packetid"
	mcnet "github.com/KonjacBot/go-mc/net"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
)

// writeTimeout bounds a write to the server or the core: a peer that stops
// reading for this long is treated as gone.
const writeTimeout = 10 * time.Second

// queueSize is how many live packets may wait for a slow core before the
// keeper drops it (the session keeps running; the core can attach again).
const queueSize = 16384

// Markers travel through an attachment's queue with the live packets so they
// keep their order relative to them.
const (
	markerResync = -1 // the server started a new configuration phase
	markerEnd    = -2 // the server connection ended; close after the packets before it
)

// session is one kept server connection.
type session struct {
	k     *Keeper
	name  string
	srv   Server
	since time.Time
	log   *slog.Logger

	up       *mcnet.Conn
	upMu     sync.Mutex
	finished []byte
	cookieMu sync.Mutex
	cookies  map[string][]int8

	// ready is closed when the first login and configuration finished; openErr
	// is why it failed.
	ready   chan struct{}
	openErr error

	mu sync.Mutex
	// playing is false while the server runs a configuration phase.
	playing bool
	// configGen counts configuration phases; configured is closed when the
	// current one finished.
	configGen  int
	configured chan struct{}
	// entered is closed when the server's Login of the current play phase
	// arrived; a core waits for it before taking a snapshot.
	entered    chan struct{}
	st         *state
	core       *attachment
	detachedAt time.Time
	idle       *time.Timer
	pending    pending
	ended      chan struct{}
	endErr     error
}

// pending are the server's requests that were forwarded to the core and not
// answered yet. When the core goes away the keeper answers them; an answer for
// anything not pending was already given by the keeper and is swallowed.
type pending struct {
	keepAlive map[int64]bool
	ping      map[int32]bool
	teleport  map[int32]bool
	packs     map[uuid.UUID]bool
	batches   int
}

func (p *pending) reset() {
	*p = pending{keepAlive: map[int64]bool{}, ping: map[int32]bool{}, teleport: map[int32]bool{}, packs: map[uuid.UUID]bool{}}
}

// attachment is one core connected to the session.
type attachment struct {
	conn  *mcnet.Conn
	queue chan pk.Packet
	done  chan struct{}
	once  sync.Once
	// live is set while the core gets the server's packets and the keeper
	// leaves the answers to it.
	live bool
	// resyncSent is set once the core was told to reconfigure.
	resyncSent atomic.Bool
}

func (a *attachment) stop() {
	a.once.Do(func() {
		close(a.done)
		_ = a.conn.Close()
	})
}

// push queues a packet for the core; false means the core fell too far behind.
func (a *attachment) push(p pk.Packet) bool {
	select {
	case a.queue <- p:
		return true
	default:
		return false
	}
}

func newSession(k *Keeper, name string, srv Server) *session {
	s := &session{
		k: k, name: name, srv: srv, since: time.Now(),
		log:        k.log.With("session", name),
		cookies:    map[string][]int8{},
		ready:      make(chan struct{}),
		configured: make(chan struct{}),
		entered:    make(chan struct{}),
		st:         newState(),
		ended:      make(chan struct{}),
		detachedAt: time.Now(),
	}
	s.pending.reset()
	return s
}

// open logs in and configures, then keeps reading the server in the
// background until the connection ends.
func (s *session) open(ctx context.Context) {
	err := s.dial(ctx)
	if err != nil {
		s.openErr = err
		close(s.ready)
		s.end(err)
		return
	}
	s.mu.Lock()
	s.playing = true
	close(s.configured)
	s.armIdle()
	s.mu.Unlock()
	close(s.ready)
	go func() { s.end(s.loop()) }()
}

// loop reads the server until the connection ends.
func (s *session) loop() error {
	ctx := context.Background()
	for {
		p, err := s.readUp(ctx)
		if err != nil {
			return err
		}
		if packetid.ClientboundPacketID(p.ID) == packetid.ClientboundStartConfiguration {
			if err := s.reconfigure(ctx); err != nil {
				return err
			}
			continue
		}
		if err := s.fromServer(p); err != nil {
			return err
		}
	}
}

// fromServer handles one play packet: record it, then forward it to a live
// core or answer it when no core will.
func (s *session) fromServer(p pk.Packet) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.st.fromServer(p); err != nil {
		s.log.Debug("保持者看不懂封包,快照會缺這一筆", "packet", packetid.ClientboundPacketID(p.ID), "error", err)
	}
	if packetid.ClientboundPacketID(p.ID) == packetid.ClientboundLogin && s.st.login != nil {
		select {
		case <-s.entered:
		default:
			close(s.entered)
		}
	}
	a := s.core
	live := a != nil && a.live
	id := packetid.ClientboundPacketID(p.ID)
	switch id {
	case packetid.ClientboundKeepAlive:
		var v pk.Long
		if err := p.Scan(&v); err != nil {
			return fmt.Errorf("讀 keepalive: %w", err)
		}
		if !live {
			return s.writeUp(pk.Marshal(packetid.ServerboundKeepAlive, v))
		}
		s.pending.keepAlive[int64(v)] = true
	case packetid.ClientboundPing:
		var v pk.Int
		if err := p.Scan(&v); err != nil {
			return fmt.Errorf("讀 ping: %w", err)
		}
		if !live {
			return s.writeUp(pk.Marshal(packetid.ServerboundPong, v))
		}
		s.pending.ping[int32(v)] = true
	case packetid.ClientboundPlayerPosition:
		var v client.PlayerPosition
		if err := decode(p, &v); err != nil {
			return err
		}
		if !live {
			return s.confirmTeleport(v.ID)
		}
		s.pending.teleport[v.ID] = true
	case packetid.ClientboundChunkBatchFinished:
		if !live {
			return s.writeUp(marshalServer(&server.ChunkBatchReceived{ChunksPerTick: chunksPerTick}))
		}
		s.pending.batches++
	case packetid.ClientboundResourcePackPush:
		var v client.AddResourcePack
		if err := decode(p, &v); err != nil {
			return err
		}
		if !live {
			return s.acceptPack(v.UUID)
		}
		s.pending.packs[v.UUID] = true
	case packetid.ClientboundCookieRequest:
		if !live {
			var v client.CookieRequest
			if err := decode(p, &v); err != nil {
				return err
			}
			response := s.cookieResponse(string(v.Key))
			return s.writeUp(marshalServer(&response))
		}
	case packetid.ClientboundStoreCookie:
		var v client.StoreCookie
		if err := decode(p, &v); err == nil {
			s.storeCookie(v.Key, v.Payload)
		}
	case packetid.ClientboundLogin, packetid.ClientboundRespawn:
		if !live {
			// What minego's player sends on entering a world; the server ignores
			// movement until it hears PlayerLoaded.
			s.st.loaded = true
			if err := s.writeUp(marshalServer(&server.PlayerLoaded{})); err != nil {
				return err
			}
			return s.writeUp(marshalServer(&server.ClientTickEnd{}))
		}
	case packetid.ClientboundDisconnect:
		var reason chat.Message
		_ = p.Scan(&reason)
		if live {
			s.forward(a, p)
		}
		return &KickError{Reason: reason.ClearString()}
	case packetid.ClientboundTransfer:
		if live {
			s.forward(a, p)
		}
		return errors.New("伺服器要求轉到別的伺服器,保持者不跟著轉")
	}
	if live {
		s.forward(a, p)
	}
	return nil
}

// chunksPerTick is what minego reports after each chunk batch.
const chunksPerTick = 64

func (s *session) forward(a *attachment, p pk.Packet) {
	if !a.push(p) {
		s.log.Warn("核心收封包太慢,先讓它離開,連線由保持者接手")
		s.dropCore(a)
	}
}

// confirmTeleport answers a position correction like the client does: confirm
// it and send the corrected position back.
func (s *session) confirmTeleport(id int32) error {
	if err := s.writeUp(marshalServer(&server.AcceptTeleportation{TeleportID: id})); err != nil {
		return err
	}
	self := s.st.self
	return s.writeUp(marshalServer(&server.MovePlayerPosRot{
		X: self.x, FeetY: self.y, Z: self.z, XRot: self.yaw, YRot: self.pitch, Flags: self.moveFlags,
	}))
}

func (s *session) acceptPack(id uuid.UUID) error {
	for _, result := range []int32{3, 0} {
		if err := s.writeUp(marshalServer(&server.ResourcePack{UUID: id, Result: result})); err != nil {
			return err
		}
	}
	return nil
}

// answerPending answers whatever the leaving core did not, then leaves the
// player standing still. Called with s.mu held when a core stops being live.
func (s *session) answerPending() error {
	var errs []error
	for id := range s.pending.keepAlive {
		errs = append(errs, s.writeUp(pk.Marshal(packetid.ServerboundKeepAlive, pk.Long(id))))
	}
	for id := range s.pending.ping {
		errs = append(errs, s.writeUp(pk.Marshal(packetid.ServerboundPong, pk.Int(id))))
	}
	for id := range s.pending.teleport {
		errs = append(errs, s.confirmTeleport(id))
	}
	for i := 0; i < s.pending.batches; i++ {
		errs = append(errs, s.writeUp(marshalServer(&server.ChunkBatchReceived{ChunksPerTick: chunksPerTick})))
	}
	for id := range s.pending.packs {
		errs = append(errs, s.acceptPack(id))
	}
	s.pending.reset()
	if s.playing {
		for _, p := range s.st.neutral() {
			errs = append(errs, s.writeUp(marshalServer(p)))
		}
	}
	return errors.Join(errs...)
}

// reconfigure handles a configuration phase the server started during play.
// A live core is taken out of the stream first (the keeper answers for it),
// then told to reconfigure; it replays the new configuration once the keeper
// finished it (see core.go).
func (s *session) reconfigure(ctx context.Context) error {
	s.mu.Lock()
	s.playing = false
	s.configGen++
	s.configured = make(chan struct{})
	s.entered = make(chan struct{})
	// Like the vanilla client, a configuration phase forgets the world; the
	// server sends a new Login and everything after it.
	s.st.resetPlay()
	if a := s.core; a != nil && a.live {
		a.live = false
		if err := s.answerPending(); err != nil {
			s.mu.Unlock()
			return err
		}
		if !a.push(pk.Packet{ID: markerResync}) {
			s.dropCore(a)
		}
	}
	// Under s.mu so no core packet in the play format slips out after the ack.
	err := s.writeUp(pk.Marshal(packetid.ServerboundConfigurationAcknowledged))
	s.mu.Unlock()
	if err != nil {
		return err
	}
	if err := s.configure(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	s.playing = true
	close(s.configured)
	s.mu.Unlock()
	return nil
}

// dropCore detaches a core from the session (s.mu held).
func (s *session) dropCore(a *attachment) {
	if s.core != a {
		return
	}
	s.core = nil
	if a.live {
		a.live = false
		if err := s.answerPending(); err != nil {
			s.log.Warn("核心離開時代回封包失敗", "error", err)
		}
	}
	s.detachedAt = time.Now()
	s.armIdle()
	a.stop()
}

// armIdle starts the countdown that closes a session no core comes back to.
func (s *session) armIdle() {
	if s.idle != nil {
		s.idle.Stop()
		s.idle = nil
	}
	if s.core != nil || s.k.detachTimeout <= 0 {
		return
	}
	s.idle = time.AfterFunc(s.k.detachTimeout, func() {
		s.mu.Lock()
		detached := s.core == nil
		s.mu.Unlock()
		if detached {
			s.close(fmt.Errorf("%s 內沒有核心接回", s.k.detachTimeout))
		}
	})
}

// close ends the server connection on purpose.
func (s *session) close(reason error) {
	s.log.Info("結束連線", "reason", reason)
	s.mu.Lock()
	if s.endErr == nil {
		s.endErr = reason
	}
	s.mu.Unlock()
	if s.up != nil {
		_ = s.up.Close()
	}
}

// end records why the connection ended and lets the core know.
func (s *session) end(err error) {
	s.mu.Lock()
	select {
	case <-s.ended:
		s.mu.Unlock()
		return
	default:
	}
	if s.endErr == nil {
		s.endErr = err
	}
	if s.idle != nil {
		s.idle.Stop()
	}
	attached := s.core != nil
	if a := s.core; a != nil {
		s.core = nil
		if !a.push(pk.Packet{ID: markerEnd}) {
			a.stop()
		}
	}
	close(s.ended)
	reason := s.endErr
	s.mu.Unlock()
	if s.up != nil {
		_ = s.up.Close()
	}
	s.log.Info("連線結束", "reason", reason)
	s.k.forget(s, reason, attached)
}

// alive reports whether the server connection is still up.
func (s *session) alive() bool {
	select {
	case <-s.ended:
		return false
	default:
		return true
	}
}

func (s *session) info() SessionInfo {
	s.mu.Lock()
	defer s.mu.Unlock()
	info := SessionInfo{Session: s.name, Address: s.srv.Address, Attached: s.core != nil, Since: s.since}
	if s.core == nil {
		info.Detached = s.detachedAt
	}
	return info
}

func marshalServer(p server.ServerboundPacket) pk.Packet {
	return pk.Marshal(p.PacketID(), p)
}
