package world

import (
	"context"

	"github.com/KonjacBot/go-mc/level"
	"github.com/KonjacBot/go-mc/level/block"
	"github.com/KonjacBot/go-mc/nbt"

	"github.com/KonjacBot/minego/v2/pkg/bot"
	"github.com/KonjacBot/minego/v2/pkg/protocol"
	cp "github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/client"
)

// 方塊實體(告示牌、箱子、旗幟…)存在區塊自己的 BlockEntity 清單裡:區塊封包本來就帶著,
// 卸載、重生時跟著區塊一起丟掉;之後的 BlockEntityData 封包(例如告示牌改字)就地更新。

// handleBlockEntities 註冊 BlockEntityData 的處理。
func (w *World) handleBlockEntities(c bot.Client) {
	c.Bus().On(func(ctx context.Context, p *cp.BlockEntityData) {
		if w.noChunks {
			return
		}
		pos := protocol.Position{int32(p.Position.X), int32(p.Position.Y), int32(p.Position.Z)}
		w.chunkLock.Lock()
		defer w.chunkLock.Unlock()
		chunk, ok := w.columns[level.ChunkPos{pos[0] >> 4, pos[2] >> 4}]
		if !ok {
			return
		}
		var be level.BlockEntity
		if !be.PackXZ(int(pos[0]&15), int(pos[2]&15)) {
			return
		}
		be.Y, be.Type, be.Data = int16(pos[1]), block.EntityType(p.Type), p.Data
		for i := range chunk.BlockEntity {
			if chunk.BlockEntity[i].XZ == be.XZ && chunk.BlockEntity[i].Y == be.Y {
				chunk.BlockEntity[i] = be
				return
			}
		}
		chunk.BlockEntity = append(chunk.BlockEntity, be)
	})
}

// BlockEntity 回傳 pos 的方塊實體種類(go-mc block.EntityType)與 NBT;
// 區塊沒載入或那裡沒有方塊實體時 ok 為 false。
//
// 回傳值只用 protocol 與 go-mc 的型別,呼叫端可以用介面斷言選用這個能力,
// 不必跟著升級 minego 才編得過。
func (w *World) BlockEntity(pos protocol.Position) (typ int32, data nbt.RawMessage, ok bool) {
	w.chunkLock.RLock()
	defer w.chunkLock.RUnlock()
	chunk, found := w.columns[level.ChunkPos{pos[0] >> 4, pos[2] >> 4}]
	if !found {
		return 0, nbt.RawMessage{}, false
	}
	for _, be := range chunk.BlockEntity {
		x, z := be.UnpackXZ()
		if int32(x) == pos[0]&15 && int32(z) == pos[2]&15 && int32(be.Y) == pos[1] {
			return int32(be.Type), be.Data, true
		}
	}
	return 0, nbt.RawMessage{}, false
}

// BlockEntitiesNear 回傳 center 周圍(每一軸相差不超過 radius 格)已載入的方塊實體位置。
func (w *World) BlockEntitiesNear(center protocol.Position, radius int32) []protocol.Position {
	w.chunkLock.RLock()
	defer w.chunkLock.RUnlock()
	var out []protocol.Position
	for cx := (center[0] - radius) >> 4; cx <= (center[0]+radius)>>4; cx++ {
		for cz := (center[2] - radius) >> 4; cz <= (center[2]+radius)>>4; cz++ {
			chunk, ok := w.columns[level.ChunkPos{cx, cz}]
			if !ok {
				continue
			}
			for _, be := range chunk.BlockEntity {
				x, z := be.UnpackXZ()
				p := protocol.Position{cx<<4 + int32(x), int32(be.Y), cz<<4 + int32(z)}
				if abs(p[0]-center[0]) <= radius && abs(p[1]-center[1]) <= radius && abs(p[2]-center[2]) <= radius {
					out = append(out, p)
				}
			}
		}
	}
	return out
}
