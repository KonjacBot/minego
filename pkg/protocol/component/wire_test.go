package component

import (
	"bytes"
	"testing"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

// 26.2 的格式(對照伺服器 jar 的 StreamCodec):雞的種類只有一個 VarInt;畫、唱片、武器音效是 holder(大於 0 是編號加 1,0 接內嵌);
// 旗幟圖案是 holderSet;藥水的隱藏效果只有 Details 沒有效果編號;盔甲紋飾圖案沒有 template_item;
// 傷害種類只有編號;模型顏色是 INT;工具與穿戴有新欄位;方塊條件的狀態是 Either、元件條件接元件自己的格式。
// 每個元件都要剛好讀完自己的位元組,不能吃到下一個元件。
func TestComponentWireFormat(t *testing.T) {
	tests := []struct {
		name string
		c    slot.Component
		wire []byte
	}{
		{"雞:編號 0", &ChickenVariant{}, []byte{0}},
		{"雞:編號 200", &ChickenVariant{}, []byte{0xc8, 0x01}},
		{"畫:登錄表編號", &PaintingVariant{}, []byte{3}},
		{"畫:內嵌", &PaintingVariant{}, append(append([]byte{0, 4, 3, 11}, "minecraft:a"...), 0, 0)},
		{"唱片:登錄表編號", &JukeboxPlayable{}, []byte{3}},
		{"旗幟圖案:標籤", &ProvidesBannerPatterns{}, append([]byte{0, 11}, "minecraft:a"...)},
		{"旗幟圖案:兩個編號", &ProvidesBannerPatterns{}, []byte{3, 5, 6}},
		{"刺擊武器:音效用登錄表編號", &PiercingWeapon{}, []byte{1, 0, 1, 5, 0}},
		{"衝刺武器:音效用登錄表編號", &KineticWeapon{}, []byte{10, 0, 0, 0, 0, 0x3f, 0x80, 0, 0, 0x3f, 0x80, 0, 0, 1, 7, 1, 8}},
		{"藥水:隱藏效果", &PotionContents{}, []byte{0, 0, 1, 2, 1, 100, 0, 1, 1, 1, 0, 50, 0, 1, 1, 0, 0}},
		{"盔甲紋飾:內嵌圖案", &Trim{}, append(append([]byte{1, 0}, append([]byte{11}, "minecraft:a"...)...), 0x0a, 8, 0, 4, 't', 'e', 'x', 't', 0, 1, 'x', 0, 0)},
		{"傷害種類:只有編號", &DamageType{}, []byte{5}},
		{"自訂模型資料:顏色是 4 bytes", &CustomModelData{}, []byte{0, 0, 0, 1, 0, 0, 0, 7}},
		{"工具:規則與創造模式能不能破壞", &Tool{}, []byte{1, 2, 9, 1, 0x40, 0, 0, 0, 0, 0x3f, 0x80, 0, 0, 1, 1}},
		{"穿戴:剪的欄位", &Equippable{}, []byte{4, 3, 0, 0, 0, 1, 1, 0, 1, 0, 5}},
		{"食用效果:播放音效是 holder", &Consumable{}, []byte{0x3f, 0xc0, 0, 0, 1, 3, 1, 1, 4, 2}},
		{"冒險模式條件:狀態、NBT、元件", &CanBreak{}, []byte{1, 1, 2, 9, 1, 2, 1, 'a', 1, 1, 'b', 1, 'c', 0, 1, 1, '1', 0, 1, 0x0a, 0, 1, 3, 5, 1, 1, 2, 0x0a, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(append(append([]byte{}, tt.wire...), 0x7f))
			if _, err := tt.c.ReadFrom(r); err != nil {
				t.Fatalf("ReadFrom: %v", err)
			}
			if r.Len() != 1 {
				t.Fatalf("讀了 %d bytes,要 %d", len(tt.wire)+1-r.Len(), len(tt.wire))
			}
			var out bytes.Buffer
			if _, err := tt.c.WriteTo(&out); err != nil {
				t.Fatalf("WriteTo: %v", err)
			}
			if !bytes.Equal(out.Bytes(), tt.wire) {
				t.Errorf("寫出 % x,要 % x", out.Bytes(), tt.wire)
			}
		})
	}
	var p PaintingVariant
	if _, err := p.ReadFrom(bytes.NewReader([]byte{0, 4, 3, 11, 'm', 'i', 'n', 'e', 'c', 'r', 'a', 'f', 't', ':', 'a', 0, 0})); err != nil {
		t.Fatal(err)
	}
	if p.Variant.Val.Width != 4 || p.Variant.Val.Height != 3 || p.Variant.Val.AssetID != "minecraft:a" {
		t.Errorf("內嵌的畫 = %+v", p.Variant.Val)
	}
}
