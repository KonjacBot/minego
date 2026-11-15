package keeper

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KonjacBot/go-mc/data/packetid"
	mcnet "github.com/KonjacBot/go-mc/net"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
)

// ControlChannel is the custom payload channel a core uses to talk to the
// keeper inside the Minecraft stream. The payload "close" ends the session
// (the player leaves); closing the socket without it only detaches.
const ControlChannel = "minego:keeper"

// SnapshotDone is the payload the keeper sends on ControlChannel right after
// the snapshot: everything before it is the state the core resumes from, so a
// core waits for it before trusting its inventory and world.
const SnapshotDone = "ready"

// errResync ends a live attachment because the server reconfigured; the core
// is taken through a new configuration phase on the same socket.
var errResync = errors.New("server reconfigured")

// errNotEntered means the configuration is done but the server's Login has not
// arrived, so there is no world to snapshot yet.
var errNotEntered = errors.New("server has not sent Login yet")

// serveCore speaks the server side of the Minecraft protocol to a core over
// the local socket: login (no encryption, no compression), a replay of the
// recorded configuration, then the snapshot and the live stream.
func (s *session) serveCore(ctx context.Context, conn *mcnet.Conn) error {
	if err := s.coreLogin(conn); err != nil {
		return err
	}
	for {
		gen, err := s.replayConfig(ctx, conn)
		if err != nil {
			return err
		}
		a, snapshot, err := s.goLive(conn, gen)
		for errors.Is(err, errNotEntered) {
			// Configured but the server's Login is not here yet.
			s.mu.Lock()
			entered := s.entered
			s.mu.Unlock()
			select {
			case <-entered:
			case <-s.ended:
				return s.endReason()
			case <-ctx.Done():
				return ctx.Err()
			}
			a, snapshot, err = s.goLive(conn, gen)
		}
		if errors.Is(err, errResync) {
			// The server reconfigured while the core was configuring: take it
			// through the new configuration before it sees any play packet.
			if cfgErr := s.restartConfig(conn); cfgErr != nil {
				return cfgErr
			}
			continue
		}
		if err != nil {
			return err
		}
		err = s.pump(ctx, a, snapshot)
		if !errors.Is(err, errResync) {
			return err
		}
	}
}

// coreLogin accepts the core's handshake and login and answers with the
// player's real profile from the server's own LoginFinished.
func (s *session) coreLogin(conn *mcnet.Conn) error {
	_ = conn.Socket.SetDeadline(time.Now().Add(s.k.readIdle))
	defer func() { _ = conn.Socket.SetDeadline(time.Time{}) }()
	var p pk.Packet
	if err := conn.ReadPacket(&p); err != nil {
		return fmt.Errorf("讀核心的握手: %w", err)
	}
	var (
		version pk.VarInt
		host    pk.String
		port    pk.UnsignedShort
		intent  pk.VarInt
	)
	if err := p.Scan(&version, &host, &port, &intent); err != nil {
		return fmt.Errorf("讀核心的握手: %w", err)
	}
	if int32(version) != protocol.ProtocolVersion || intent != 2 {
		return fmt.Errorf("核心的協定版本 %d(意圖 %d)和保持者的 %d 不同", version, intent, protocol.ProtocolVersion)
	}
	if err := conn.ReadPacket(&p); err != nil || packetid.ServerboundPacketID(p.ID) != packetid.ServerboundLoginHello {
		return fmt.Errorf("讀核心的登入: %w (封包 %#x)", err, p.ID)
	}
	if err := conn.WritePacket(pk.Packet{ID: int32(packetid.ClientboundLoginLoginFinished), Data: s.finished}); err != nil {
		return fmt.Errorf("回核心登入完成: %w", err)
	}
	if err := conn.ReadPacket(&p); err != nil || packetid.ServerboundPacketID(p.ID) != packetid.ServerboundLoginLoginAcknowledged {
		return fmt.Errorf("等核心確認登入: %w (封包 %#x)", err, p.ID)
	}
	return nil
}

// replayConfig waits for the server's configuration to finish, sends the
// recorded configuration to the core and waits for it to finish too. The
// core's answers (client information, known packs, keepalives) were already
// given to the server by the keeper and are dropped.
func (s *session) replayConfig(ctx context.Context, conn *mcnet.Conn) (gen int, err error) {
	s.mu.Lock()
	configured := s.configured
	s.mu.Unlock()
	select {
	case <-configured:
	case <-s.ended:
		return 0, s.endReason()
	case <-ctx.Done():
		return 0, ctx.Err()
	}
	s.mu.Lock()
	gen = s.configGen
	config := s.st.config
	s.mu.Unlock()
	for _, p := range config {
		if err := writeCore(conn, p); err != nil {
			return 0, err
		}
	}
	if err := writeCore(conn, pk.Marshal(packetid.ClientboundConfigFinishConfiguration)); err != nil {
		return 0, err
	}
	for {
		var p pk.Packet
		_ = conn.Socket.SetReadDeadline(time.Now().Add(s.k.readIdle))
		if err := conn.ReadPacket(&p); err != nil {
			return 0, fmt.Errorf("等核心完成設定: %w", err)
		}
		if packetid.ServerboundPacketID(p.ID) == packetid.ServerboundConfigFinishConfiguration {
			_ = conn.Socket.SetReadDeadline(time.Time{})
			return gen, nil
		}
	}
}

// restartConfig moves a core that is in the play state back to configuration.
func (s *session) restartConfig(conn *mcnet.Conn) error {
	if err := writeCore(conn, pk.Marshal(packetid.ClientboundStartConfiguration)); err != nil {
		return err
	}
	for {
		var p pk.Packet
		_ = conn.Socket.SetReadDeadline(time.Now().Add(s.k.readIdle))
		if err := conn.ReadPacket(&p); err != nil {
			return fmt.Errorf("等核心回到設定階段: %w", err)
		}
		if packetid.ServerboundPacketID(p.ID) == packetid.ServerboundConfigurationAcknowledged {
			_ = conn.Socket.SetReadDeadline(time.Time{})
			return nil
		}
	}
}

// goLive atomically takes the snapshot and makes the core live, so every
// server packet after the snapshot reaches the core exactly once. A previous
// core is dropped.
func (s *session) goLive(conn *mcnet.Conn, gen int) (*attachment, []pk.Packet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.alive() {
		return nil, nil, s.endErr
	}
	if gen != s.configGen || !s.playing {
		return nil, nil, errResync
	}
	snapshot, _, ok := s.st.snapshot()
	if !ok {
		return nil, nil, errNotEntered
	}
	if old := s.core; old != nil && old.conn != conn {
		s.log.Info("新的核心接上,讓舊的離開")
		s.dropCore(old)
	}
	a := &attachment{conn: conn, queue: make(chan pk.Packet, queueSize), done: make(chan struct{}), live: true}
	s.core = a
	if s.idle != nil {
		s.idle.Stop()
		s.idle = nil
	}
	s.log.Info("核心接上", "snapshot", len(snapshot))
	return a, snapshot, nil
}

// pump writes the snapshot and the live stream to the core and forwards what
// the core sends, until the core leaves, the server ends, or the server
// reconfigures (errResync, the core already told to reconfigure).
func (s *session) pump(ctx context.Context, a *attachment, snapshot []pk.Packet) error {
	writer := make(chan error, 1)
	go func() { writer <- s.writeLive(a, snapshot) }()

	err := s.readCore(ctx, a)
	if errors.Is(err, errResync) {
		if werr := <-writer; !errors.Is(werr, errResync) {
			return werr
		}
		return errResync
	}
	s.mu.Lock()
	s.dropCore(a)
	s.mu.Unlock()
	a.stop()
	<-writer
	return err
}

func (s *session) writeLive(a *attachment, snapshot []pk.Packet) error {
	for _, p := range snapshot {
		if err := writeCore(a.conn, p); err != nil {
			a.stop()
			return err
		}
	}
	done := client.CustomPayload{Channel: ControlChannel, Data: []byte(SnapshotDone)}
	if err := writeCore(a.conn, marshal(&done)); err != nil {
		a.stop()
		return err
	}
	for {
		select {
		case <-a.done:
			return nil
		case p := <-a.queue:
			switch p.ID {
			case markerResync:
				a.resyncSent.Store(true)
				if err := writeCore(a.conn, pk.Marshal(packetid.ClientboundStartConfiguration)); err != nil {
					a.stop()
					return err
				}
				return errResync
			case markerEnd:
				a.stop()
				return s.endReason()
			}
			if err := writeCore(a.conn, p); err != nil {
				a.stop()
				return err
			}
		}
	}
}

// readCore forwards the core's packets to the server until the core leaves or
// acknowledges a reconfiguration.
func (s *session) readCore(ctx context.Context, a *attachment) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		var p pk.Packet
		if err := a.conn.ReadPacket(&p); err != nil {
			return fmt.Errorf("核心離開: %w", err)
		}
		if packetid.ServerboundPacketID(p.ID) == packetid.ServerboundConfigurationAcknowledged {
			if a.resyncSent.Load() {
				return errResync
			}
			continue
		}
		if err := s.fromCore(a, p); err != nil {
			return err
		}
	}
}

// fromCore handles one packet from a live core: answers the keeper already
// gave are swallowed, keeper control messages are handled, the rest is
// recorded and forwarded.
func (s *session) fromCore(a *attachment, p pk.Packet) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.core != a || !a.live || !s.playing {
		return nil
	}
	switch packetid.ServerboundPacketID(p.ID) {
	case packetid.ServerboundKeepAlive:
		var v pk.Long
		if err := p.Scan(&v); err != nil || !s.pending.keepAlive[int64(v)] {
			return nil
		}
		delete(s.pending.keepAlive, int64(v))
	case packetid.ServerboundPong:
		var v pk.Int
		if err := p.Scan(&v); err != nil || !s.pending.ping[int32(v)] {
			return nil
		}
		delete(s.pending.ping, int32(v))
	case packetid.ServerboundAcceptTeleportation:
		var v server.AcceptTeleportation
		if err := decode(p, &v); err != nil || !s.pending.teleport[v.TeleportID] {
			return nil
		}
		delete(s.pending.teleport, v.TeleportID)
	case packetid.ServerboundChunkBatchReceived:
		if s.pending.batches == 0 {
			return nil
		}
		s.pending.batches--
	case packetid.ServerboundResourcePack:
		var v server.ResourcePack
		if err := decode(p, &v); err != nil || !s.pending.packs[v.UUID] {
			return nil
		}
		if v.Result != 3 && v.Result != 4 {
			// Accepted and downloaded are intermediate; the next answer is final.
			delete(s.pending.packs, v.UUID)
		}
	case packetid.ServerboundPlayerLoaded:
		if s.st.loaded {
			return nil
		}
	case packetid.ServerboundCustomPayload:
		var v server.CustomPayload
		if err := decode(p, &v); err == nil && v.Channel == ControlChannel {
			if string(v.Data) == "close" {
				go s.close(errors.New("核心要求結束連線"))
			}
			return nil
		}
	}
	if err := s.st.fromCore(p); err != nil {
		s.log.Debug("保持者看不懂核心的封包,狀態可能少記一筆", "packet", packetid.ServerboundPacketID(p.ID), "error", err)
	}
	return s.writeUp(p)
}

func writeCore(conn *mcnet.Conn, p pk.Packet) error {
	_ = conn.Socket.SetWriteDeadline(time.Now().Add(writeTimeout))
	if err := conn.WritePacket(p); err != nil {
		return fmt.Errorf("送封包給核心: %w", err)
	}
	return nil
}

func (s *session) endReason() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.endErr != nil {
		return s.endErr
	}
	return errors.New("連線已結束")
}
