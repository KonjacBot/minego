package player

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl64"

	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/bot"
	"git.konjactw.dev/patyhank/minego/pkg/game/world"
	"git.konjactw.dev/patyhank/minego/pkg/protocol"
	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/client"
	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/server"
)

type Player struct {
	c bot.Client

	entity  *world.Entity
	stateID int32

	lastReceivedPacketTime time.Time
}

// New 創建新的 Player 實例
func New(c bot.Client) *Player {
	pl := &Player{
		c:       c,
		entity:  &world.Entity{},
		stateID: 1,
	}

	c.PacketHandler().AddGenericPacketHandler(func(ctx context.Context, pk client.ClientboundPacket) {
		pl.lastReceivedPacketTime = time.Now()
	})

	bot.AddHandler(c, func(ctx context.Context, p *client.KeepAlive) {
		_ = c.WritePacket(ctx, &server.KeepAlive{
			ID: p.ID,
		})
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.Disconnect) {
		fmt.Println(p.Reason.String())
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.SystemChatMessage) {
		if !p.Overlay {
			bot.PublishEvent(c, MessageEvent{Message: p.Content})
		}
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.PlayerPosition) {
		fmt.Println(p)
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

		c.WritePacket(context.Background(), &server.AcceptTeleportation{TeleportID: p.ID})
		c.WritePacket(context.Background(), &server.MovePlayerPosRot{
			X:     p.X,
			FeetY: p.Y,
			Z:     p.Z,
			Yaw:   p.XRot,
			Pitch: p.YRot,
			Flags: 0x00,
		})
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.PlayerRotation) {
		pl.entity.SetRotation(mgl64.Vec2{float64(p.Yaw), float64(p.Pitch)})
	})

	return pl
}

func (p *Player) CheckServer() {
	for time.Since(p.lastReceivedPacketTime) > 50*time.Millisecond && p.c.IsConnected() {
		time.Sleep(50 * time.Millisecond)
	}
}

// StateID 返回當前狀態 ID
func (p *Player) StateID() int32 {
	return p.stateID
}

// UpdateStateID 更新狀態 ID
func (p *Player) UpdateStateID(id int32) {
	p.stateID = id
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

	currentPos := p.entity.Position()
	direction := pos.Sub(currentPos)
	distance := direction.Len()

	if distance == 0 {
		return nil // 已經在目標位置
	}

	segmentLength := 8.0

	for {
		currentPos = p.entity.Position()

		direction = pos.Sub(currentPos)
		distance = direction.Len()

		if distance == 0 {
			return nil
		}

		// 正規化方向向量
		direction = direction.Normalize()

		moveDistance := math.Min(segmentLength, distance)

		target := currentPos.Add(direction.Mul(moveDistance))

		if err := p.c.WritePacket(context.Background(), &server.MovePlayerPos{
			X:     target.X(),
			FeetY: target.Y(),
			Z:     target.Z(),
			Flags: 0x00,
		}); err != nil {
			return fmt.Errorf("failed to move player: %w", err)
		}

		time.Sleep(50 * time.Millisecond)
	}
	return nil
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

	// 使用 A* 演算法尋找路徑
	path, err := AStar(p.c.World(), currentPos, pos)
	if err != nil {
		return fmt.Errorf("failed to find path: %w", err)
	}

	if len(path) == 0 {
		return fmt.Errorf("no path found to target position")
	}

	// 沿著路徑移動
	for _, waypoint := range path {
		if err := p.c.WritePacket(context.Background(), &server.MovePlayerPos{
			X:     waypoint.X(),
			FeetY: waypoint.Y(),
			Z:     waypoint.Z(),
			Flags: 0x0,
		}); err != nil {
			return fmt.Errorf("failed to move to waypoint: %w", err)
		}

		// 短暫延遲以模擬真實移動
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func (p *Player) UpdateLocation() {
	_ = p.c.WritePacket(context.Background(), &server.MovePlayerPosRot{
		X:     p.entity.Position().X(),
		FeetY: p.entity.Position().Y(),
		Z:     p.entity.Position().Z(),
		Yaw:   float32(p.entity.Rotation().X()),
		Pitch: float32(p.entity.Rotation().Y()),
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
		Yaw:   yaw,
		Pitch: pitch,
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
		Sequence: p.stateID,
		Location: pk.Position{X: int(pos[0]), Y: int(pos[1]), Z: int(pos[2])},
		Face:     1,
	}

	if err := p.c.WritePacket(context.Background(), startPacket); err != nil {
		return fmt.Errorf("failed to send start destroy packet: %w", err)
	}

	// 發送完成挖掘封包
	finishPacket := &server.PlayerAction{
		Status:   2,
		Sequence: p.stateID,
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
		Sequence:    p.stateID,
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
		Sequence:    p.stateID,
	}

	return p.c.WritePacket(context.Background(), packet)
}

// OpenContainer 打開指定位置的容器
func (p *Player) OpenContainer(pos protocol.Position) (bot.Container, error) {
	if p.c == nil {
		return nil, fmt.Errorf("client is not initialized")
	}

	// 發送使用物品封包來打開容器
	packet := &server.UseItemOn{
		Hand:           1,
		Location:       pk.Position{X: int(pos[0]), Y: int(pos[1]), Z: int(pos[2])},
		Face:           1,
		CursorX:        0.5,
		CursorY:        0.5,
		CursorZ:        0.5,
		InsideBlock:    false,
		WorldBorderHit: false,
		Sequence:       p.stateID,
	}

	if err := p.c.WritePacket(context.Background(), packet); err != nil {
		return nil, fmt.Errorf("failed to open container: %w", err)
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	for ctx.Err() == nil && p.c.Inventory().CurrentContainerID() <= 0 {
		time.Sleep(time.Millisecond * 50)
	}

	for ctx.Err() == nil && p.c.Inventory().Container().SlotCount() == 0 {
		time.Sleep(time.Millisecond * 50)
	}

	if ctx.Err() != nil {
		return nil, fmt.Errorf("failed to open container: %w", ctx.Err())
	}
	if p.c.Inventory().CurrentContainerID() <= 0 {
		return nil, fmt.Errorf("failed to open container: no container opened")
	}

	return p.c.Inventory().Container(), nil
}

// UseItem 使用指定手中的物品
func (p *Player) UseItem(hand int8) error {
	if p.c == nil {
		return fmt.Errorf("client is not initialized")
	}

	return p.c.WritePacket(context.Background(), &server.UseItem{
		Hand:     int32(hand),
		Sequence: p.stateID,
		Yaw:      0,
		Pitch:    0,
	})
}

// OpenMenu 打開指定命令的選單
func (p *Player) OpenMenu(command string) (bot.Container, error) {
	if p.c == nil {
		return nil, fmt.Errorf("client is not initialized")
	}

	if err := p.c.WritePacket(context.Background(), &server.ChatCommand{
		Command: command,
	}); err != nil {
		return nil, fmt.Errorf("failed to open menu with command '%s': %w", command, err)
	}

	// 返回客戶端的容器處理器
	return p.c.Inventory().Container(), nil
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
