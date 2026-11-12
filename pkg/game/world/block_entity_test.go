package world

import (
	"context"
	"testing"

	"github.com/KonjacBot/go-mc/level"
	"github.com/KonjacBot/go-mc/level/block"
	"github.com/KonjacBot/go-mc/nbt"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/v2/pkg/protocol"
	gameclient "github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/client"
)

func rawNBT(t *testing.T, v any) nbt.RawMessage {
	t.Helper()
	b, err := nbt.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var m nbt.RawMessage
	if err := nbt.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	return m
}

type signNBT struct {
	FrontText struct {
		Messages []string `nbt:"messages"`
	} `nbt:"front_text"`
}

func sign(t *testing.T, line string) nbt.RawMessage {
	var s signNBT
	s.FrontText.Messages = []string{line, "", "", ""}
	return rawNBT(t, s)
}

func firstLine(t *testing.T, data nbt.RawMessage) string {
	t.Helper()
	var s signNBT
	if err := data.Unmarshal(&s); err != nil {
		t.Fatal(err)
	}
	return s.FrontText.Messages[0]
}

func TestBlockEntities(t *testing.T) {
	c := newWorldTestClient()
	w := NewWorld(c)
	c.world = w
	ctx := context.Background()

	// 區塊 (-1, 2):方塊實體在 x=-13 (=-16+3), y=64, z=37 (=32+5)。
	chunk := level.EmptyChunk(24)
	var be level.BlockEntity
	be.PackXZ(3, 5)
	be.Y, be.Type, be.Data = 64, block.EntityType(7), sign(t, "總和16")
	chunk.BlockEntity = append(chunk.BlockEntity, be)
	c.handler.HandlePacket(ctx, &gameclient.LevelChunkWithLight{Pos: level.ChunkPos{-1, 2}, Data: chunk})

	pos := protocol.Position{-13, 64, 37}
	typ, data, ok := w.BlockEntity(pos)
	if !ok || typ != 7 || firstLine(t, data) != "總和16" {
		t.Fatalf("BlockEntity = %d %v %v", typ, data, ok)
	}
	if _, _, ok := w.BlockEntity(protocol.Position{-13, 65, 37}); ok {
		t.Fatal("y 不同不該有方塊實體")
	}

	// 改字:就地更新;新放的告示牌:加進清單。
	c.handler.HandlePacket(ctx, &gameclient.BlockEntityData{Position: pk.Position{X: -13, Y: 64, Z: 37}, Type: 7, Data: sign(t, "總和17")})
	c.handler.HandlePacket(ctx, &gameclient.BlockEntityData{Position: pk.Position{X: -10, Y: 64, Z: 40}, Type: 7, Data: sign(t, "大")})
	if _, data, _ := w.BlockEntity(pos); firstLine(t, data) != "總和17" {
		t.Fatalf("更新後 %q", firstLine(t, data))
	}
	if _, data, ok := w.BlockEntity(protocol.Position{-10, 64, 40}); !ok || firstLine(t, data) != "大" {
		t.Fatal("新的方塊實體沒有加入")
	}

	near := w.BlockEntitiesNear(protocol.Position{-12, 64, 38}, 3)
	if len(near) != 2 {
		t.Fatalf("附近 %v", near)
	}
	if near := w.BlockEntitiesNear(protocol.Position{-12, 64, 38}, 1); len(near) != 1 || near[0] != pos {
		t.Fatalf("半徑 1 內 %v", near)
	}

	// 沒載入的區塊:BlockEntityData 被忽略。
	c.handler.HandlePacket(ctx, &gameclient.BlockEntityData{Position: pk.Position{X: 100, Y: 64, Z: 100}, Type: 7, Data: sign(t, "x")})
	if _, _, ok := w.BlockEntity(protocol.Position{100, 64, 100}); ok {
		t.Fatal("沒載入的區塊不該有方塊實體")
	}

	c.handler.HandlePacket(ctx, &gameclient.ForgetLevelChunk{Pos: gameclient.PackedChunkPos{-1, 2}})
	if _, _, ok := w.BlockEntity(pos); ok {
		t.Fatal("卸載區塊後方塊實體還在")
	}
}
