package component

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/KonjacBot/go-mc/level/item"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

func TestLodestoneTrackerOfficialWire(t *testing.T) {
	want := LodestoneTracker{
		Target: pk.Option[GlobalPosition, *GlobalPosition]{
			Has: true,
			Val: GlobalPosition{
				Dimension: "minecraft:overworld",
				Position:  pk.Position{X: 1, Y: 64, Z: -3},
			},
		},
		Tracked: true,
	}

	wire := []byte{
		1,
		19, 'm', 'i', 'n', 'e', 'c', 'r', 'a', 'f', 't', ':', 'o', 'v', 'e', 'r', 'w', 'o', 'r', 'l', 'd',
		0, 0, 0, 127, 255, 255, 208, 64,
		1,
	}

	assertComponentGolden(t, wire, &want, new(LodestoneTracker))
}

func TestEquippableOfficialWire(t *testing.T) {
	want := Equippable{
		Slot:            3,
		EquipSound:      pk.OptID[SoundEvent, *SoundEvent]{Has: true, ID: 2},
		Dispensable:     true,
		Swappable:       false,
		DamageOnHurt:    true,
		EquipOnInteract: false,
		CanBeSheared:    true,
		ShearingSound:   pk.OptID[SoundEvent, *SoundEvent]{Has: true, ID: 5},
	}
	wire := []byte{3, 3, 0, 0, 0, 1, 0, 1, 0, 1, 6}

	assertComponentGolden(t, wire, &want, new(Equippable))
}

func TestEquippableReadClearsReusedValue(t *testing.T) {
	full := buildEquippableWire(t,
		func(c *Equippable) {
			c.Slot = 3
			c.EquipSound = pk.OptID[SoundEvent, *SoundEvent]{Has: true, ID: 2}
			c.AssetID = pk.Option[pk.Identifier, *pk.Identifier]{Has: true, Val: "minecraft:test_asset"}
			c.CameraOverlay = pk.Option[pk.Identifier, *pk.Identifier]{Has: true, Val: "minecraft:test_overlay"}
			c.AllowedEntities = pk.Option[pk.IDSet, *pk.IDSet]{Has: true, Val: pk.IDSet{TagName: "minecraft:test", IDs: []int32{1}}}
			c.Dispensable = true
			c.DamageOnHurt = true
			c.EquipOnInteract = true
			c.CanBeSheared = true
			c.ShearingSound = pk.OptID[SoundEvent, *SoundEvent]{Has: true, ID: 5}
		},
	)
	minimal := []byte{3, 3, 0, 0, 0, 1, 0, 1, 0, 1, 6}

	value := Equippable{AssetID: pk.Option[pk.Identifier, *pk.Identifier]{Has: true, Val: "stale"}}
	if _, err := value.ReadFrom(bytes.NewReader(full)); err != nil {
		t.Fatal(err)
	}
	if _, err := value.ReadFrom(bytes.NewReader(minimal)); err != nil {
		t.Fatal(err)
	}
	if value.AssetID.Has || value.CameraOverlay.Has || value.AllowedEntities.Has {
		t.Fatalf("stale optional state retained: %#v", value)
	}
	if value.EquipOnInteract || !value.CanBeSheared || !value.Dispensable || value.Swappable || !value.DamageOnHurt {
		t.Fatalf("decoded booleans wrong after reuse: %#v", value)
	}
}

func TestPotDecorationsOfficialWire(t *testing.T) {
	one, two, three, four := item.ID(1), item.ID(2), item.ID(3), item.ID(4)
	want := PotDecorations{Back: &one, Left: &two, Right: &three, Front: &four}
	wire := []byte{4, 1, 2, 3, 4}

	assertComponentGolden(t, wire, &want, new(PotDecorations))
}

func TestPotDecorationsReadClearsReusedValue(t *testing.T) {
	one, two, three, four := item.ID(1), item.ID(2), item.ID(3), item.ID(4)
	value := PotDecorations{Back: &one, Left: &two, Right: &three, Front: &four}

	if _, err := value.ReadFrom(bytes.NewReader([]byte{4, 1, 2, 3, 4})); err != nil {
		t.Fatal(err)
	}
	if _, err := value.ReadFrom(bytes.NewReader([]byte{2, 1, 2})); err != nil {
		t.Fatal(err)
	}
	if value.Back == nil || *value.Back != 1 || value.Left == nil || *value.Left != 2 || value.Right != nil || value.Front != nil {
		t.Fatalf("reused pot decorations retained stale state: %#v", value)
	}
}

func TestPotDecorationsRejectsOverCap(t *testing.T) {
	if _, err := new(PotDecorations).ReadFrom(bytes.NewReader([]byte{5})); err == nil {
		t.Fatal("ReadFrom() accepted pot decorations length > 4")
	}
}

func TestPotionContentsOfficialWireCustomNameOptional(t *testing.T) {
	absent := PotionContents{}
	assertComponentGolden(t, []byte{0, 0, 0, 0}, &absent, new(PotionContents))

	present := PotionContents{CustomName: pk.Option[pk.String, *pk.String]{Has: true, Val: "custom"}}
	assertComponentGolden(t, []byte{0, 0, 0, 1, 6, 'c', 'u', 's', 't', 'o', 'm'}, &present, new(PotionContents))
}

func TestPotionContentsReadClearsReusedValue(t *testing.T) {
	full := buildPotionContentsWire(t, func(c *PotionContents) {
		c.PotionID = pk.Option[pk.VarInt, *pk.VarInt]{Has: true, Val: 7}
		c.CustomColor = pk.Option[pk.Int, *pk.Int]{Has: true, Val: 0x112233}
		c.CustomEffects = []PotionEffect{{TypeID: 5, Details: PotionEffectDetails{HasHiddenEffect: true, HiddenEffect: &PotionEffect{TypeID: 9, Details: PotionEffectDetails{ShowParticles: true}}}}}
		c.CustomName = pk.Option[pk.String, *pk.String]{Has: true, Val: "custom"}
	})
	minimal := []byte{0, 0, 0, 0}

	value := PotionContents{CustomName: pk.Option[pk.String, *pk.String]{Has: true, Val: "stale"}}
	if _, err := value.ReadFrom(bytes.NewReader(full)); err != nil {
		t.Fatal(err)
	}
	if _, err := value.ReadFrom(bytes.NewReader(minimal)); err != nil {
		t.Fatal(err)
	}
	if value.PotionID.Has || value.CustomColor.Has || value.CustomName.Has || len(value.CustomEffects) != 0 {
		t.Fatalf("reused potion contents retained stale state: %#v", value)
	}
}

func TestPotionEffectDetailsReadClearsHiddenEffect(t *testing.T) {
	withHidden := []byte{2, 120, 1, 1, 0, 1, 9, 0, 20, 0, 1, 0, 0}
	withoutHidden := []byte{2, 120, 1, 1, 0, 0}

	value := PotionEffectDetails{HasHiddenEffect: true, HiddenEffect: &PotionEffect{TypeID: 99}}
	if _, err := value.ReadFrom(bytes.NewReader(withHidden)); err != nil {
		t.Fatal(err)
	}
	if value.HiddenEffect == nil {
		t.Fatal("expected hidden effect after first decode")
	}
	if _, err := value.ReadFrom(bytes.NewReader(withoutHidden)); err != nil {
		t.Fatal(err)
	}
	if value.HasHiddenEffect || value.HiddenEffect != nil {
		t.Fatalf("hidden effect retained after absent decode: %#v", value)
	}
}

func TestRecipesOfficialWireUsesRootList(t *testing.T) {
	want := Recipes{RecipeIDs: []string{"minecraft:test", "minecraft:other"}}
	wire := appendNBTStringList(nil, []string{"minecraft:test", "minecraft:other"})
	assertComponentGolden(t, wire, &want, new(Recipes))
	if wire[0] != 9 {
		t.Fatalf("root tag = %d, want TAG_List (9)", wire[0])
	}
}

func TestRecipesReadClearsReusedValue(t *testing.T) {
	value := Recipes{RecipeIDs: []string{"stale", "entries"}}
	if _, err := value.ReadFrom(bytes.NewReader(appendNBTStringList(nil, []string{"minecraft:test", "minecraft:other"}))); err != nil {
		t.Fatal(err)
	}
	if _, err := value.ReadFrom(bytes.NewReader(appendNBTStringList(nil, []string{"minecraft:one"}))); err != nil {
		t.Fatal(err)
	}
	if len(value.RecipeIDs) != 1 || value.RecipeIDs[0] != "minecraft:one" {
		t.Fatalf("reused recipes retained stale state: %#v", value)
	}
}

func TestContainerLootReadClearsReusedValue(t *testing.T) {
	withSeed := buildContainerLootWire(t, ContainerLoot{LootTable: "minecraft:chests/test", Seed: 12})
	withoutSeed := buildContainerLootWire(t, ContainerLoot{LootTable: "minecraft:chests/empty"})

	value := ContainerLoot{LootTable: "stale", Seed: 99}
	if _, err := value.ReadFrom(bytes.NewReader(withSeed)); err != nil {
		t.Fatal(err)
	}
	if _, err := value.ReadFrom(bytes.NewReader(withoutSeed)); err != nil {
		t.Fatal(err)
	}
	if value.LootTable != "minecraft:chests/empty" || value.Seed != 0 {
		t.Fatalf("reused container loot retained stale state: %#v", value)
	}
}

func TestDebugStickStateReadClearsReusedValue(t *testing.T) {
	full := buildDebugStickStateWire(t, DebugStickState{Properties: map[string]string{"minecraft:stone": "facing"}})
	empty := buildDebugStickStateWire(t, DebugStickState{Properties: map[string]string{}})

	value := DebugStickState{Properties: map[string]string{"stale": "value"}}
	if _, err := value.ReadFrom(bytes.NewReader(full)); err != nil {
		t.Fatal(err)
	}
	if _, err := value.ReadFrom(bytes.NewReader(empty)); err != nil {
		t.Fatal(err)
	}
	if len(value.Properties) != 0 {
		t.Fatalf("reused debug stick state retained stale map: %#v", value)
	}
}

func TestMapDecorationsReadClearsReusedValue(t *testing.T) {
	full := buildMapDecorationsWire(t, MapDecorations{Decorations: map[string]MapDecorationEntry{"home": {Type: "minecraft:player", X: 1.5, Z: -2.5, Rotation: 0.25}}})
	empty := buildMapDecorationsWire(t, MapDecorations{Decorations: map[string]MapDecorationEntry{}})

	value := MapDecorations{Decorations: map[string]MapDecorationEntry{"stale": {Type: "minecraft:player"}}}
	if _, err := value.ReadFrom(bytes.NewReader(full)); err != nil {
		t.Fatal(err)
	}
	if _, err := value.ReadFrom(bytes.NewReader(empty)); err != nil {
		t.Fatal(err)
	}
	if len(value.Decorations) != 0 {
		t.Fatalf("reused map decorations retained stale map: %#v", value)
	}
}

func assertComponentGolden[T comparableComponent](t *testing.T, wire []byte, want T, got T) {
	t.Helper()

	read, err := got.ReadFrom(bytes.NewReader(wire))
	if err != nil {
		t.Fatal(err)
	}
	if read != int64(len(wire)) {
		t.Fatalf("ReadFrom() = %d, want %d", read, len(wire))
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("decoded = %#v, want %#v", got, want)
	}

	var encoded bytes.Buffer
	written, err := want.WriteTo(&encoded)
	if err != nil {
		t.Fatal(err)
	}
	if written != int64(len(wire)) || !bytes.Equal(encoded.Bytes(), wire) {
		t.Fatalf("WriteTo() = %x (%d), want %x (%d)", encoded.Bytes(), written, wire, len(wire))
	}
}

type comparableComponent interface {
	ReadFrom(io.Reader) (int64, error)
	WriteTo(io.Writer) (int64, error)
}

func appendNBTStringList(dst []byte, values []string) []byte {
	dst = append(dst, 9, 8)
	dst = append(dst, 0, 0, 0, byte(len(values)))
	for _, value := range values {
		dst = append(dst, byte(len(value)>>8), byte(len(value)))
		dst = append(dst, value...)
	}
	return dst
}

func buildEquippableWire(t *testing.T, apply func(*Equippable)) []byte {
	t.Helper()
	value := Equippable{}
	apply(&value)
	var encoded bytes.Buffer
	if _, err := value.WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	return encoded.Bytes()
}

func buildPotionContentsWire(t *testing.T, apply func(*PotionContents)) []byte {
	t.Helper()
	value := PotionContents{}
	apply(&value)
	var encoded bytes.Buffer
	if _, err := value.WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	return encoded.Bytes()
}

func buildContainerLootWire(t *testing.T, value ContainerLoot) []byte {
	t.Helper()
	var encoded bytes.Buffer
	if _, err := value.WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	return encoded.Bytes()
}

func buildDebugStickStateWire(t *testing.T, value DebugStickState) []byte {
	t.Helper()
	var encoded bytes.Buffer
	if _, err := value.WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	return encoded.Bytes()
}

func buildMapDecorationsWire(t *testing.T, value MapDecorations) []byte {
	t.Helper()
	var encoded bytes.Buffer
	if _, err := value.WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	return encoded.Bytes()
}
