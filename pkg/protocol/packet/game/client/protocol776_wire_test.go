package client

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/KonjacBot/go-mc/data/packetid"
	"github.com/KonjacBot/go-mc/level"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/component"
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

type writerTo interface {
	WriteTo(io.Writer) (int64, error)
}

func mustWire(t *testing.T, fields ...writerTo) []byte {
	t.Helper()
	var buf bytes.Buffer
	for _, field := range fields {
		if _, err := field.WriteTo(&buf); err != nil {
			t.Fatal(err)
		}
	}
	return buf.Bytes()
}

// Derived from the protocol 776 server jar via:
// javap -classpath <server-26.2.jar> -p -c net.minecraft.commands.synchronization.ArgumentTypeInfos
// The static initializer registers command argument types in wire-registry order.
var official776CommandParserNames = []string{
	"brigadier:bool",
	"brigadier:float",
	"brigadier:double",
	"brigadier:integer",
	"brigadier:long",
	"brigadier:string",
	"minecraft:entity",
	"minecraft:game_profile",
	"minecraft:block_pos",
	"minecraft:column_pos",
	"minecraft:vec3",
	"minecraft:vec2",
	"minecraft:block_state",
	"minecraft:block_predicate",
	"minecraft:item_stack",
	"minecraft:item_predicate",
	"minecraft:color",
	"minecraft:hex_color",
	"minecraft:component",
	"minecraft:style",
	"minecraft:message",
	"minecraft:nbt_compound_tag",
	"minecraft:nbt_tag",
	"minecraft:nbt_path",
	"minecraft:objective",
	"minecraft:objective_criteria",
	"minecraft:operation",
	"minecraft:particle",
	"minecraft:angle",
	"minecraft:rotation",
	"minecraft:scoreboard_slot",
	"minecraft:score_holder",
	"minecraft:swizzle",
	"minecraft:team",
	"minecraft:item_slot",
	"minecraft:item_slots",
	"minecraft:resource_location",
	"minecraft:function",
	"minecraft:entity_anchor",
	"minecraft:int_range",
	"minecraft:float_range",
	"minecraft:dimension",
	"minecraft:gamemode",
	"minecraft:time",
	"minecraft:resource_or_tag",
	"minecraft:resource_or_tag_key",
	"minecraft:resource",
	"minecraft:resource_key",
	"minecraft:resource_selector",
	"minecraft:template_mirror",
	"minecraft:template_rotation",
	"minecraft:heightmap",
	"minecraft:loot_table",
	"minecraft:loot_predicate",
	"minecraft:loot_modifier",
	"minecraft:dialog",
	"minecraft:uuid",
}

func TestOfficial776CommandParserRegistryIDs(t *testing.T) {
	checks := map[int32]string{
		commandParserFloat:            "brigadier:float",
		commandParserDouble:           "brigadier:double",
		commandParserInteger:          "brigadier:integer",
		commandParserLong:             "brigadier:long",
		commandParserString:           "brigadier:string",
		commandParserEntity:           "minecraft:entity",
		commandParserScoreHolder:      "minecraft:score_holder",
		commandParserTime:             "minecraft:time",
		commandParserResourceOrTag:    "minecraft:resource_or_tag",
		commandParserResourceOrTagKey: "minecraft:resource_or_tag_key",
		commandParserResource:         "minecraft:resource",
		commandParserResourceKey:      "minecraft:resource_key",
		commandParserResourceSelector: "minecraft:resource_selector",
	}

	for id, want := range checks {
		if int(id) >= len(official776CommandParserNames) {
			t.Fatalf("parser id %d out of bounds for official 776 table", id)
		}
		if got := official776CommandParserNames[id]; got != want {
			t.Fatalf("official parser id %d = %q, want %q", id, got, want)
		}
	}
}

func TestCommandParserPropertyWireCoverage(t *testing.T) {
	tests := []struct {
		name  string
		id    int32
		props []byte
	}{
		{name: "float", id: 1, props: []byte{0x03, 0x3f, 0x80, 0x00, 0x00, 0x40, 0x40, 0x00, 0x00}},
		{name: "double", id: 2, props: []byte{0x01, 0x3f, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{name: "integer", id: 3, props: []byte{0x03, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x40}},
		{name: "long", id: 4, props: []byte{0x02, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}},
		{name: "string", id: 5, props: []byte{0x02}},
		{name: "entity", id: 6, props: []byte{0x03}},
		{name: "score holder", id: 31, props: []byte{0x01}},
		{name: "time", id: 43, props: []byte{0x00, 0x00, 0x00, 0x2a}},
		{name: "resource or tag", id: 44, props: mustWire(t, pk.Identifier("minecraft:block"))},
		{name: "resource or tag key", id: 45, props: mustWire(t, pk.Identifier("minecraft:item"))},
		{name: "resource", id: 46, props: mustWire(t, pk.Identifier("minecraft:loot_table"))},
		{name: "resource key", id: 47, props: mustWire(t, pk.Identifier("minecraft:biome"))},
		{name: "resource selector", id: 48, props: mustWire(t, pk.Identifier("minecraft:dimension"))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := CommandParser{ID: tt.id, Properties: tt.props}
			encoded := mustWire(t, want)
			var got CommandParser
			read, err := got.ReadFrom(bytes.NewReader(encoded))
			if err != nil {
				t.Fatal(err)
			}
			if read != int64(len(encoded)) {
				t.Fatalf("ReadFrom() count = %d, encoded length = %d", read, len(encoded))
			}
			if got.ID != want.ID || !bytes.Equal(got.Properties, want.Properties) {
				t.Fatalf("decoded parser = %#v, want %#v", got, want)
			}
		})
	}
}

func TestCommandsPacketWireRoundTrip(t *testing.T) {
	want := Commands{
		Nodes: []CommandNode{
			{Flags: 0x00, Children: []int32{1, 2}},
			{Flags: 0x21, Children: []int32{2}, Name: "say"},
			{
				Flags:          0x1e,
				Children:       nil,
				Redirect:       1,
				Name:           "amount",
				Parser:         CommandParser{ID: 3, Properties: []byte{0x03, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x40}},
				SuggestionType: "minecraft:ask_server",
			},
		},
		RootIndex: 0,
	}

	expected := mustWire(t,
		pk.VarInt(3),
		pk.Byte(0x00), pk.VarInt(2), pk.VarInt(1), pk.VarInt(2),
		pk.Byte(0x21), pk.VarInt(1), pk.VarInt(2), pk.String("say"),
		pk.Byte(0x1e), pk.VarInt(0), pk.VarInt(1), pk.String("amount"), pk.VarInt(3),
	)
	expected = append(expected, want.Nodes[2].Parser.Properties...)
	expected = append(expected, mustWire(t, pk.Identifier("minecraft:ask_server"), pk.VarInt(0))...)

	var encoded bytes.Buffer
	written, err := want.WriteTo(&encoded)
	if err != nil {
		t.Fatal(err)
	}
	if written != int64(encoded.Len()) {
		t.Fatalf("WriteTo() count = %d, encoded length = %d", written, encoded.Len())
	}
	if !bytes.Equal(encoded.Bytes(), expected) {
		t.Fatalf("encoded commands = %v, want %v", encoded.Bytes(), expected)
	}

	var got Commands
	read, err := got.ReadFrom(bytes.NewReader(encoded.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	if read != written {
		t.Fatalf("ReadFrom() count = %d, want %d", read, written)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("decoded commands = %#v, want %#v", got, want)
	}
}

func TestSetEquipmentUsesContinuationBitAndResetsSlice(t *testing.T) {
	want := SetEquipment{
		EntityID: 7,
		Equipment: Equipment{
			{Slot: 1, Item: slot.Slot{}},
			{Slot: 2, Item: slot.Slot{}},
		},
	}
	expected := []byte{0x07, 0x81, 0x00, 0x02, 0x00}

	var encoded bytes.Buffer
	if _, err := want.WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded.Bytes(), expected) {
		t.Fatalf("encoded equipment = %v, want %v", encoded.Bytes(), expected)
	}

	var got SetEquipment
	got.Equipment = Equipment{{Slot: 99}, {Slot: 98}, {Slot: 97}}
	if _, err := got.ReadFrom(bytes.NewReader(expected)); err != nil {
		t.Fatal(err)
	}
	if got.EntityID != 7 || len(got.Equipment) != 2 || got.Equipment[0].Slot != 1 || got.Equipment[1].Slot != 2 {
		t.Fatalf("decoded equipment = %#v", got)
	}

	shorter := []byte{0x07, 0x04, 0x00}
	if _, err := got.ReadFrom(bytes.NewReader(shorter)); err != nil {
		t.Fatal(err)
	}
	if len(got.Equipment) != 1 || got.Equipment[0].Slot != 4 {
		t.Fatalf("decoded shorter equipment = %#v", got.Equipment)
	}
}

func TestChunkBiomesUsesOuterListWire(t *testing.T) {
	want := ChunkBiomes{Chunks: []ChunkBiomeData{{Pos: level.ChunkPos{0x11223344, 0x55667788}, Data: []byte{0xaa, 0xbb}}}}
	expected := []byte{0x01, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x02, 0xaa, 0xbb}

	var encoded bytes.Buffer
	if _, err := want.WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded.Bytes(), expected) {
		t.Fatalf("encoded chunk biomes = %v, want %v", encoded.Bytes(), expected)
	}

	var got ChunkBiomes
	if _, err := got.ReadFrom(bytes.NewReader(expected)); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("decoded chunk biomes = %#v, want %#v", got, want)
	}
}

func TestChunkBiomesReadFromClearsReusedChunksSlice(t *testing.T) {
	full := []byte{0x02, 0, 0, 0, 1, 0, 0, 0, 2, 0x01, 0xaa, 0, 0, 0, 3, 0, 0, 0, 4, 0x01, 0xbb}
	shorter := []byte{0x01, 0, 0, 0, 5, 0, 0, 0, 6, 0x01, 0xcc}
	empty := []byte{0x00}

	var got ChunkBiomes
	if _, err := got.ReadFrom(bytes.NewReader(full)); err != nil {
		t.Fatal(err)
	}
	if len(got.Chunks) != 2 {
		t.Fatalf("decoded chunks len = %d, want 2", len(got.Chunks))
	}
	if _, err := got.ReadFrom(bytes.NewReader(shorter)); err != nil {
		t.Fatal(err)
	}
	if len(got.Chunks) != 1 || got.Chunks[0].Pos != (level.ChunkPos{5, 6}) || !bytes.Equal(got.Chunks[0].Data, []byte{0xcc}) {
		t.Fatalf("decoded shorter chunks = %#v", got.Chunks)
	}
	if _, err := got.ReadFrom(bytes.NewReader(empty)); err != nil {
		t.Fatal(err)
	}
	if len(got.Chunks) != 0 {
		t.Fatalf("decoded empty chunks len = %d, want 0", len(got.Chunks))
	}
}

func TestCommandsReadFromClearsReusedNodesSlice(t *testing.T) {
	full := mustWire(t,
		pk.VarInt(2),
		CommandNode{Flags: 0x00, Children: []int32{1}},
		CommandNode{Flags: 0x00},
		pk.VarInt(0),
	)
	shorter := mustWire(t,
		pk.VarInt(1),
		CommandNode{Flags: 0x00},
		pk.VarInt(0),
	)
	empty := mustWire(t, pk.VarInt(0), pk.VarInt(0))

	var got Commands
	if _, err := got.ReadFrom(bytes.NewReader(full)); err != nil {
		t.Fatal(err)
	}
	if len(got.Nodes) != 2 {
		t.Fatalf("decoded nodes len = %d, want 2", len(got.Nodes))
	}
	if _, err := got.ReadFrom(bytes.NewReader(shorter)); err != nil {
		t.Fatal(err)
	}
	if len(got.Nodes) != 1 || len(got.Nodes[0].Children) != 0 {
		t.Fatalf("decoded shorter nodes = %#v", got.Nodes)
	}
	if _, err := got.ReadFrom(bytes.NewReader(empty)); err != nil {
		t.Fatal(err)
	}
	if len(got.Nodes) != 0 {
		t.Fatalf("decoded empty nodes len = %d, want 0", len(got.Nodes))
	}
}

func TestGameResourcePackPacketIDsAreNotSwapped(t *testing.T) {
	if got := (&AddResourcePack{}).PacketID(); got != packetid.ClientboundResourcePackPush {
		t.Fatalf("AddResourcePack PacketID = %v, want %v", got, packetid.ClientboundResourcePackPush)
	}
	if got := (&RemoveResourcePack{}).PacketID(); got != packetid.ClientboundResourcePackPop {
		t.Fatalf("RemoveResourcePack PacketID = %v, want %v", got, packetid.ClientboundResourcePackPop)
	}
}

func TestSoundEffectIncludesSource(t *testing.T) {
	want := SoundEffect{SoundID: 1, Source: 2, EffectPositionX: 3, EffectPositionY: 4, EffectPositionZ: 5, Volume: 1.25, Pitch: 0.5, Seed: 9}
	expected := mustWire(t, pk.VarInt(1), pk.VarInt(2), pk.Int(3), pk.Int(4), pk.Int(5), pk.Float(1.25), pk.Float(0.5), pk.Long(9))

	var encoded bytes.Buffer
	if _, err := want.WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded.Bytes(), expected) {
		t.Fatalf("encoded sound = %v, want %v", encoded.Bytes(), expected)
	}

	var got SoundEffect
	if _, err := got.ReadFrom(bytes.NewReader(expected)); err != nil {
		t.Fatal(err)
	}
	if got.Source != want.Source || got.Seed != want.Seed {
		t.Fatalf("decoded sound = %#v, want %#v", got, want)
	}
}

func TestEntitySoundEffectUsesLongSeed(t *testing.T) {
	want := EntitySoundEffect{
		SoundEvent:    pk.OptID[component.SoundEvent, *component.SoundEvent]{Has: true, ID: 0},
		SoundCategory: 2,
		EntityID:      3,
		Volume:        1.25,
		Pitch:         0.5,
		Seed:          0x0102030405060708,
	}
	expected := mustWire(t, want)

	var got EntitySoundEffect
	read, err := got.ReadFrom(bytes.NewReader(expected))
	if err != nil {
		t.Fatal(err)
	}
	if read != int64(len(expected)) || got.Seed != want.Seed || got.SoundCategory != want.SoundCategory {
		t.Fatalf("decoded entity sound = %#v after %d bytes, want %#v after %d", got, read, want, len(expected))
	}
}

func TestMapDataOptionalDecorationsAndPatchRoundTrip(t *testing.T) {
	wire := mustWire(t,
		pk.VarInt(12),
		pk.Byte(3),
		pk.Boolean(true),
		pk.Boolean(true),
		pk.VarInt(1),
		MapIcon{Type: 4, X: 5, Z: 6, Direction: 7},
		MapColorPatch{Columns: 1, Rows: 2, X: 3, Z: 4, Data: []pk.UnsignedByte{5, 6}},
	)

	var got MapData
	if _, err := got.ReadFrom(bytes.NewReader(wire)); err != nil {
		t.Fatal(err)
	}
	if got.MapID != 12 || got.Scale != 3 || !got.Locked || !got.HasDecorations || len(got.Decorations) != 1 {
		t.Fatalf("decoded map data = %#v", got)
	}
	if got.Decorations[0].Type != 4 || got.Decorations[0].X != 5 || got.Decorations[0].Z != 6 || got.Decorations[0].Direction != 7 || got.Decorations[0].DisplayName.Has {
		t.Fatalf("decoded decoration = %#v", got.Decorations[0])
	}
	if got.ColorPatch.Columns != 1 || got.ColorPatch.Rows != 2 || got.ColorPatch.X != 3 || got.ColorPatch.Z != 4 || !reflect.DeepEqual(got.ColorPatch.Data, []pk.UnsignedByte{5, 6}) {
		t.Fatalf("decoded patch = %#v", got.ColorPatch)
	}

	var encoded bytes.Buffer
	if _, err := got.WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded.Bytes(), wire) {
		t.Fatalf("re-encoded map data = %v, want %v", encoded.Bytes(), wire)
	}
}

func TestMapDataOptionalDecorationsAndPatchAbsent(t *testing.T) {
	wire := mustWire(t, pk.VarInt(9), pk.Byte(1), pk.Boolean(false), pk.Boolean(false), MapColorPatch{})
	var got MapData
	if _, err := got.ReadFrom(bytes.NewReader(wire)); err != nil {
		t.Fatal(err)
	}
	if got.HasDecorations || len(got.Decorations) != 0 || got.ColorPatch.Columns != 0 || len(got.ColorPatch.Data) != 0 {
		t.Fatalf("decoded absent map data = %#v", got)
	}
}

func TestMapDataReadFromClearsReusedDecorationsSlice(t *testing.T) {
	full := mustWire(t,
		pk.VarInt(1), pk.Byte(1), pk.Boolean(false), pk.Boolean(true), pk.VarInt(2),
		MapIcon{Type: 1, X: 2, Z: 3, Direction: 4},
		MapIcon{Type: 5, X: 6, Z: 7, Direction: 8},
		MapColorPatch{},
	)
	shorter := mustWire(t,
		pk.VarInt(1), pk.Byte(1), pk.Boolean(false), pk.Boolean(true), pk.VarInt(1),
		MapIcon{Type: 9, X: 10, Z: 11, Direction: 12},
		MapColorPatch{},
	)
	absent := mustWire(t, pk.VarInt(1), pk.Byte(1), pk.Boolean(false), pk.Boolean(false), MapColorPatch{})

	var got MapData
	if _, err := got.ReadFrom(bytes.NewReader(full)); err != nil {
		t.Fatal(err)
	}
	if len(got.Decorations) != 2 {
		t.Fatalf("decoded decorations len = %d, want 2", len(got.Decorations))
	}
	if _, err := got.ReadFrom(bytes.NewReader(shorter)); err != nil {
		t.Fatal(err)
	}
	if len(got.Decorations) != 1 || got.Decorations[0].Type != 9 {
		t.Fatalf("decoded shorter decorations = %#v", got.Decorations)
	}
	if _, err := got.ReadFrom(bytes.NewReader(absent)); err != nil {
		t.Fatal(err)
	}
	if len(got.Decorations) != 0 {
		t.Fatalf("decoded absent decorations len = %d, want 0", len(got.Decorations))
	}
}

func TestMapColorPatchReadFromResetsZeroColumnsAndShorterData(t *testing.T) {
	full := mustWire(t, pk.UnsignedByte(1), pk.UnsignedByte(2), pk.UnsignedByte(3), pk.UnsignedByte(4), pk.VarInt(2), pk.UnsignedByte(5), pk.UnsignedByte(6))
	shorter := mustWire(t, pk.UnsignedByte(1), pk.UnsignedByte(7), pk.UnsignedByte(8), pk.UnsignedByte(9), pk.VarInt(1), pk.UnsignedByte(10))
	empty := []byte{0x00}

	var got MapColorPatch
	if _, err := got.ReadFrom(bytes.NewReader(full)); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Data, []pk.UnsignedByte{5, 6}) {
		t.Fatalf("decoded full patch = %#v", got)
	}
	if _, err := got.ReadFrom(bytes.NewReader(shorter)); err != nil {
		t.Fatal(err)
	}
	if got.Rows != 7 || got.X != 8 || got.Z != 9 || !reflect.DeepEqual(got.Data, []pk.UnsignedByte{10}) {
		t.Fatalf("decoded shorter patch = %#v", got)
	}
	if _, err := got.ReadFrom(bytes.NewReader(empty)); err != nil {
		t.Fatal(err)
	}
	if got.Columns != 0 || got.Rows != 0 || got.X != 0 || got.Z != 0 || len(got.Data) != 0 {
		t.Fatalf("decoded zero-columns patch = %#v", got)
	}
}

func TestTestInstanceBlockStatusUsesOptionalVarIntVec3i(t *testing.T) {
	wire := appendNBTString(nil, "ok")
	wire = append(wire, 0x01, 0x01, 0x02, 0x03)

	var got TestInstanceBlockStatus
	if _, err := got.ReadFrom(bytes.NewReader(wire)); err != nil {
		t.Fatal(err)
	}
	if got.Status.Text != "ok" || !got.Size.Has || got.Size.Val != (WaypointVec3i{X: 1, Y: 2, Z: 3}) {
		t.Fatalf("decoded test block status = %#v", got)
	}
}

func TestWaypointUsesOfficialIconAndEmptyType(t *testing.T) {
	waypointID := uuid.MustParse("00112233-4455-6677-8899-aabbccddeeff")
	want := Waypoint{
		Operation:        1,
		IsUUIDIdentifier: true,
		UUID:             waypointID,
		Icon: WaypointIcon{
			Style: "minecraft:default",
			Color: pk.Option[WaypointColor, *WaypointColor]{Has: true, Val: WaypointColor{R: 1, G: 2, B: 3}},
		},
		WaypointType: 0,
	}
	expected := mustWire(t, pk.VarInt(1), pk.Boolean(true), pk.UUID(waypointID), want.Icon, pk.VarInt(0))

	var encoded bytes.Buffer
	if _, err := want.WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded.Bytes(), expected) {
		t.Fatalf("encoded waypoint = %v, want %v", encoded.Bytes(), expected)
	}

	var got Waypoint
	if _, err := got.ReadFrom(bytes.NewReader(expected)); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("decoded waypoint = %#v, want %#v", got, want)
	}
}

func TestWaypointReadFromClearsInactiveIdentifierAndUnionFields(t *testing.T) {
	first := mustWire(t,
		pk.VarInt(1),
		pk.Boolean(false),
		pk.String("spawn"),
		WaypointIcon{Style: "minecraft:default"},
		pk.VarInt(1),
		WaypointVec3i{X: 1, Y: 2, Z: 3},
	)
	waypointID := uuid.MustParse("00112233-4455-6677-8899-aabbccddeeff")
	second := mustWire(t,
		pk.VarInt(2),
		pk.Boolean(true),
		pk.UUID(waypointID),
		WaypointIcon{Style: "minecraft:default"},
		pk.VarInt(2),
		WaypointChunkPos{X: 4, Z: 5},
	)
	third := mustWire(t,
		pk.VarInt(0),
		pk.Boolean(true),
		pk.UUID(uuid.UUID{}),
		WaypointIcon{Style: "minecraft:default"},
		pk.VarInt(0),
	)

	var got Waypoint
	if _, err := got.ReadFrom(bytes.NewReader(first)); err != nil {
		t.Fatal(err)
	}
	if got.Name != "spawn" || got.UUID != (uuid.UUID{}) || got.WaypointPlayerPos != (WaypointVec3i{X: 1, Y: 2, Z: 3}) {
		t.Fatalf("decoded first waypoint = %#v", got)
	}
	if _, err := got.ReadFrom(bytes.NewReader(second)); err != nil {
		t.Fatal(err)
	}
	if got.Name != "" || got.UUID != waypointID || got.WaypointPlayerPos != (WaypointVec3i{}) || got.WaypointChunkPos != (WaypointChunkPos{X: 4, Z: 5}) || got.WaypointAzimuth != (WaypointAzimuth{}) {
		t.Fatalf("decoded second waypoint = %#v", got)
	}
	if _, err := got.ReadFrom(bytes.NewReader(third)); err != nil {
		t.Fatal(err)
	}
	if got.Name != "" || got.UUID != (uuid.UUID{}) || got.WaypointPlayerPos != (WaypointVec3i{}) || got.WaypointChunkPos != (WaypointChunkPos{}) || got.WaypointAzimuth != (WaypointAzimuth{}) {
		t.Fatalf("decoded third waypoint = %#v", got)
	}
}

func TestStopSoundCountsAndWire(t *testing.T) {
	want := StopSound{Flags: 0x03, Source: 1, Sound: "minecraft:test"}
	expected := mustWire(t, pk.Byte(0x03), pk.VarInt(1), pk.Identifier("minecraft:test"))

	var encoded bytes.Buffer
	written, err := want.WriteTo(&encoded)
	if err != nil {
		t.Fatal(err)
	}
	if written != int64(len(expected)) || !bytes.Equal(encoded.Bytes(), expected) {
		t.Fatalf("encoded stop sound = (%d, %v), want (%d, %v)", written, encoded.Bytes(), len(expected), expected)
	}

	var got StopSound
	read, err := got.ReadFrom(bytes.NewReader(expected))
	if err != nil {
		t.Fatal(err)
	}
	if read != written || got != want {
		t.Fatalf("decoded stop sound = %#v after %d bytes, want %#v after %d", got, read, want, written)
	}
}

func TestStopSoundReadFromClearsOmittedFieldsOnReuse(t *testing.T) {
	full := mustWire(t, pk.Byte(0x03), pk.VarInt(7), pk.Identifier("minecraft:bell"))
	sourceOnly := mustWire(t, pk.Byte(0x01), pk.VarInt(9))
	soundOnly := mustWire(t, pk.Byte(0x02), pk.Identifier("minecraft:note_block"))
	none := mustWire(t, pk.Byte(0x00))

	var got StopSound
	if _, err := got.ReadFrom(bytes.NewReader(full)); err != nil {
		t.Fatal(err)
	}
	if got.Source != 7 || got.Sound != "minecraft:bell" {
		t.Fatalf("decoded full stop sound = %#v", got)
	}
	if _, err := got.ReadFrom(bytes.NewReader(sourceOnly)); err != nil {
		t.Fatal(err)
	}
	if got.Source != 9 || got.Sound != "" {
		t.Fatalf("decoded source-only stop sound = %#v", got)
	}
	if _, err := got.ReadFrom(bytes.NewReader(soundOnly)); err != nil {
		t.Fatal(err)
	}
	if got.Source != 0 || got.Sound != "minecraft:note_block" {
		t.Fatalf("decoded sound-only stop sound = %#v", got)
	}
	if _, err := got.ReadFrom(bytes.NewReader(none)); err != nil {
		t.Fatal(err)
	}
	if got.Source != 0 || got.Sound != "" {
		t.Fatalf("decoded empty stop sound = %#v", got)
	}
}
