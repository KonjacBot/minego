package player

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"sync"
	"time"

	"github.com/go-gl/mathgl/mgl64"

	"github.com/KonjacBot/go-mc/level/block"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/bot"
	"github.com/KonjacBot/minego/pkg/game/world"
	"github.com/KonjacBot/minego/pkg/protocol"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
)

type Player struct {
	c  bot.Client
	mu sync.RWMutex

	abilities         int8
	entity            *world.Entity
	stateID, sequence int32

	lastReceivedPacketTime time.Time
}

// New 創建新的 Player 實例
func New(c bot.Client) *Player {
	pl := &Player{
		c:       c,
		entity:  &world.Entity{},
		stateID: 1,
	}

	startup := sync.OnceFunc(func() {
		//go func() {
		//	ticker := time.NewTicker(50 * time.Millisecond)
		//	for range ticker.C {
		//		_ = c.WritePacket(context.Background(), &server.ClientTickEnd{})
		//	}
		//}()
	})

	c.PacketHandler().AddGenericPacketHandler(func(ctx context.Context, pk client.ClientboundPacket) {
		pl.mu.Lock()
		pl.lastReceivedPacketTime = time.Now()
		pl.mu.Unlock()
	})

	bot.AddHandler(c, func(ctx context.Context, p *client.KeepAlive) {
		_ = c.WritePacket(ctx, &server.KeepAlive{
			ID: p.ID,
		})
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.PlayerAbilities) {
		pl.mu.Lock()
		pl.abilities = p.Flags
		pl.mu.Unlock()
	})

	bot.AddHandler(c, func(ctx context.Context, p *client.Login) {
		startup()
		c.WritePacket(ctx, &server.ClientInformation{
			Location:            "zh_TW",
			ViewDistance:        16,
			ChatMode:            0,
			ChatColor:           true,
			DisplayedSkinParts:  127,
			MainHand:            0,
			EnableTextFiltering: false,
			AllowListing:        true,
			ParticleStatus:      0,
		})
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.Ping) {
		_ = c.WritePacket(ctx, &server.Pong{ID: p.ID})
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.SystemChatMessage) {
		if !p.Overlay {
			bot.PublishEvent(c, MessageEvent{Message: p.Content})
		}
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.PlayerPosition) {
		position := pl.entity.Position()
		if p.Flags&0x01 != 0 {
			position[0] += p.X
		} else {
			position[0] = p.X
		}

		if p.Flags&0x02 != 0 {
			position[1] += p.Y
		} else {
			position[1] = p.Y
		}

		if p.Flags&0x04 != 0 {
			position[2] += p.Z
		} else {
			position[2] = p.Z
		}

		pl.entity.SetPosition(position)

		rot := pl.entity.Rotation()
		if p.Flags&0x08 != 0 {
			rot[0] += float64(p.XRot)
		} else {
			rot[0] = float64(p.XRot)
		}

		if p.Flags&0x10 != 0 {
			rot[1] += float64(p.YRot)
		} else {
			rot[1] = float64(p.YRot)
		}
		pl.entity.SetRotation(rot)

		_ = c.WritePacket(ctx, &server.AcceptTeleportation{TeleportID: p.ID})
		_ = c.WritePacket(ctx, &server.MovePlayerPosRot{
			X:     position[0],
			FeetY: position[1],
			Z:     position[2],
			XRot:  float32(rot[0]),
			YRot:  float32(rot[1]),
			Flags: 0x00,
		})
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.PlayerRotation) {
		pl.entity.SetRotation(mgl64.Vec2{float64(p.Yaw), float64(p.Pitch)})
	})

	return pl
}

func (p *Player) CheckServer() {
	deadline := time.Now().Add(5 * time.Second)
	for p.c.IsConnected() && time.Now().Before(deadline) {
		p.mu.RLock()
		lastReceived := p.lastReceivedPacketTime
		p.mu.RUnlock()
		if time.Since(lastReceived) <= 50*time.Millisecond {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// StateID 返回當前狀態 ID
func (p *Player) StateID() int32 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.stateID
}

// UpdateStateID 更新狀態 ID
func (p *Player) UpdateStateID(id int32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stateID = id
}

// Sequence 返回當前互動狀態 ID
func (p *Player) Sequence() int32 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.sequence
}

// UpdateSequence 更新互動狀態 ID
func (p *Player) UpdateSequence(id int32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if id > p.sequence {
		p.sequence = id
	}
}

func (p *Player) nextSequence() int32 {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.sequence++
	return p.sequence
}

// Entity 返回玩家實體
func (p *Player) Entity() bot.Entity {
	return p.entity
}

// FlyTo 直線飛行到指定位置，每5格飛行一段
func (p *Player) FlyTo(pos mgl64.Vec3) error {
	if p.c == nil {
		return fmt.Errorf("client is not initialized")
	}

	if p.entity == nil {
		return fmt.Errorf("player entity is not initialized")
	}
	p.mu.RLock()
	canFly := p.abilities&0x04 != 0
	p.mu.RUnlock()
	if !canFly {
		return fmt.Errorf("player abilities not requirements")
	}

	currentPos := p.entity.Position()
	direction := pos.Sub(currentPos)
	distance := direction.Len()

	if distance == 0 {
		return nil // 已經在目標位置
	}

	const segmentLength = 5.0

	slog.Info("flyto begin", "from", vecString(currentPos), "target", vecString(pos), "distance", distance)

	for {
		currentPos = p.entity.Position()

		direction = pos.Sub(currentPos)
		distance = direction.Len()

		if distance < 0.05 {
			slog.Info("flyto done", "at", vecString(currentPos))
			return nil
		}

		// 正規化方向向量
		direction = direction.Normalize()

		moveDistance := math.Min(segmentLength, distance)

		target := currentPos.Add(direction.Mul(moveDistance))

		slog.Info("flyto step", "target", vecString(target), "remaining", distance)
		p.entity.SetPosition(target)

		if err := p.c.WritePacket(context.Background(), &server.MovePlayerPos{
			X:     target.X(),
			FeetY: target.Y(),
			Z:     target.Z(),
			Flags: 0x00,
		}); err != nil {
			return fmt.Errorf("failed to move player: %w", err)
		}
		time.Sleep(50 * time.Millisecond)

		p.mu.RLock()
		canFly = p.abilities&0x04 != 0
		p.mu.RUnlock()
		if !canFly {
			return fmt.Errorf("player abilities not requirements")
		}
		if !p.entity.Position().ApproxEqualThreshold(target, 0.5) {
			return fmt.Errorf("failed to move player: position updated by server")
		}
	}
}

func vecString(v mgl64.Vec3) string {
	return fmt.Sprintf("%.2f,%.2f,%.2f", v.X(), v.Y(), v.Z())
}

// WalkTo 使用 A* 演算法步行到指定位置
func (p *Player) WalkTo(pos mgl64.Vec3) error {
	if p.c == nil {
		return fmt.Errorf("client is not initialized")
	}

	if p.entity == nil {
		return fmt.Errorf("player entity is not initialized")
	}

	currentPos := p.entity.Position()

	path, err := AStar(p.c.World(), currentPos, pos, 4096)
	if err != nil {
		return fmt.Errorf("failed to find path: %w", err)
	}

	if len(path) == 0 {
		return fmt.Errorf("no path found to target position")
	}

	// 沿著路徑移動
	for _, waypoint := range path {
		if err := p.c.WritePacket(context.Background(), &server.MovePlayerPos{
			X:     waypoint.X() + 0.5,
			FeetY: waypoint.Y(),
			Z:     waypoint.Z() + 0.5,
			Flags: 0x0,
		}); err != nil {
			return fmt.Errorf("failed to move to waypoint: %w", err)
		}

		time.Sleep(10 * time.Millisecond)
	}

	return nil
}

func (p *Player) UpdateLocation() {
	_ = p.c.WritePacket(context.Background(), &server.MovePlayerPosRot{
		X:     p.entity.Position().X(),
		FeetY: p.entity.Position().Y(),
		Z:     p.entity.Position().Z(),
		XRot:  float32(p.entity.Rotation().X()),
		YRot:  float32(p.entity.Rotation().Y()),
		Flags: 0x00,
	})
}

// LookAt 看向指定位置
func (p *Player) LookAt(target mgl64.Vec3) error {
	if p.c == nil {
		return fmt.Errorf("client is not initialized")
	}

	if p.entity == nil {
		return fmt.Errorf("player entity is not initialized")
	}

	// 計算視角
	playerPos := p.entity.Position()
	direction := target.Sub(playerPos).Normalize()

	// 計算 yaw 和 pitch
	yaw := float32(math.Atan2(-direction.X(), direction.Z()) * 180 / math.Pi)
	pitch := float32(math.Asin(-direction.Y()) * 180 / math.Pi)

	p.entity.SetRotation(mgl64.Vec2{float64(yaw), float64(pitch)})

	return p.c.WritePacket(context.Background(), &server.MovePlayerRot{
		XRot:  yaw,
		YRot:  pitch,
		Flags: 0x00,
	})
}

// BreakBlock 破壞指定位置的方塊
func (p *Player) BreakBlock(pos protocol.Position) error {
	if p.c == nil {
		return fmt.Errorf("client is not initialized")
	}

	// 發送開始挖掘封包
	startPacket := &server.PlayerAction{
		Status:   0,
		Sequence: p.nextSequence(),
		Location: pk.Position{X: int(pos[0]), Y: int(pos[1]), Z: int(pos[2])},
		Face:     1,
	}

	if err := p.c.WritePacket(context.Background(), startPacket); err != nil {
		return fmt.Errorf("failed to send start destroy packet: %w", err)
	}

	// 發送完成挖掘封包
	finishPacket := &server.PlayerAction{
		Status:   2,
		Sequence: p.nextSequence(),
		Location: pk.Position{X: int(pos[0]), Y: int(pos[1]), Z: int(pos[2])},
		Face:     1,
	}

	return p.c.WritePacket(context.Background(), finishPacket)
}

// PlaceBlock 在指定位置放置方塊
func (p *Player) PlaceBlock(pos protocol.Position) error {
	if p.c == nil {
		return fmt.Errorf("client is not initialized")
	}

	packet := &server.UseItemOn{
		Hand:        0,
		Location:    pk.Position{X: int(pos[0]), Y: int(pos[1]), Z: int(pos[2])},
		Face:        1,
		CursorX:     0.5,
		CursorY:     0.5,
		CursorZ:     0.5,
		InsideBlock: false,
		Sequence:    p.nextSequence(),
	}

	return p.c.WritePacket(context.Background(), packet)
}

// PlaceBlock 在指定位置放置方塊
func (p *Player) PlaceBlockWithArgs(pos protocol.Position, face int32, cursor mgl64.Vec3) error {
	if p.c == nil {
		return fmt.Errorf("client is not initialized")
	}

	packet := &server.UseItemOn{
		Hand:        0,
		Location:    pk.Position{X: int(pos[0]), Y: int(pos[1]), Z: int(pos[2])},
		Face:        face,
		CursorX:     float32(cursor[0]),
		CursorY:     float32(cursor[1]),
		CursorZ:     float32(cursor[2]),
		InsideBlock: false,
		Sequence:    p.nextSequence(),
	}

	return p.c.WritePacket(context.Background(), packet)
}

// OpenContainer 打開指定位置的容器
func (p *Player) OpenContainer(pos protocol.Position, hand int32) (bot.Container, error) {
	if p.c == nil {
		return nil, fmt.Errorf("client is not initialized")
	}
	w := p.c.World()
	if w == nil {
		return nil, fmt.Errorf("world is not initialized")
	}
	blk, err := w.GetBlock(pos)
	if err != nil {
		return nil, fmt.Errorf("failed to open container: %w", err)
	}
	if block.IsAirBlock(blk) {
		return nil, fmt.Errorf("failed to open container: block at %v is air", pos)
	}

	previousContainerID := p.c.Inventory().CurrentContainerID()
	// 發送使用物品封包來打開容器
	packet := &server.UseItemOn{
		Hand:           hand,
		Location:       pk.Position{X: int(pos[0]), Y: int(pos[1]), Z: int(pos[2])},
		Face:           1,
		CursorX:        0.5,
		CursorY:        0.5,
		CursorZ:        0.5,
		InsideBlock:    false,
		WorldBorderHit: false,
		Sequence:       p.nextSequence(),
	}

	if err := p.c.WritePacket(context.Background(), packet); err != nil {
		return nil, fmt.Errorf("failed to open container: %w", err)
	}

	return p.waitForContainer(previousContainerID)
}

// UseItem 使用指定手中的物品
func (p *Player) UseItem(hand int8) error {
	if p.c == nil {
		return fmt.Errorf("client is not initialized")
	}

	return p.c.WritePacket(context.Background(), &server.UseItem{
		Hand:     int32(hand),
		Sequence: p.nextSequence(),
		Yaw:      0,
		Pitch:    0,
	})
}

// OpenMenu 打開指定命令的選單
func (p *Player) OpenMenu(command string) (bot.Container, error) {
	if p.c == nil {
		return nil, fmt.Errorf("client is not initialized")
	}

	previousContainerID := p.c.Inventory().CurrentContainerID()
	if err := p.c.WritePacket(context.Background(), &server.ChatCommand{
		Command: command,
	}); err != nil {
		return nil, fmt.Errorf("failed to open menu with command '%s': %w", command, err)
	}

	return p.waitForContainer(previousContainerID)
}

func (p *Player) waitForContainer(previousID int32) (bot.Container, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		currentID := p.c.Inventory().CurrentContainerID()
		container := p.c.Inventory().Container()
		if currentID > 0 && currentID != previousID && container != nil && container.SlotCount() > 0 {
			return container, nil
		}
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("failed to open container: %w", ctx.Err())
		case <-ticker.C:
		}
	}
}

func (p *Player) Command(msg string) error {
	return p.c.WritePacket(context.Background(), &server.ChatCommand{
		Command: msg,
	})
}

func (p *Player) Chat(msg string) error {
	return p.c.WritePacket(context.Background(), &server.Chat{
		Message: msg,
	})
}
