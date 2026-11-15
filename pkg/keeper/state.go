package keeper

import (
	"bytes"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/google/uuid"

	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/v2/pkg/protocol/metadata"
	"github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/server"
)

// state is what the keeper remembers about the player's world, built from the
// packets that pass through it in both directions. A core that attaches gets a
// snapshot of it (see snapshot.go) and ends up where a client that never left
// would be.
//
// Packets the keeper does not understand are still forwarded; they are simply
// not part of the snapshot. Decoding failures only lose that piece of state.
type state struct {
	// config is every configuration-phase packet the server sent (registries,
	// tags, known packs, features), replayed to the core's configuration phase.
	config []pk.Packet

	// login is the last Login with the spawn info of the latest Respawn, so the
	// snapshot enters the current dimension directly.
	login    *client.Login
	playerID int32
	// loaded is set once PlayerLoaded went to the server for the current world;
	// the core's duplicate is swallowed.
	loaded bool

	// singles keeps the last packet of each kind whose latest copy is the whole
	// state (time, health, experience, world border, recipes...).
	singles map[packetid.ClientboundPacketID]pk.Packet
	// events are the persistent game events (game mode, weather, waiting for
	// chunks), by event id in arrival order.
	events []pk.Packet
	// recipes is the recipe book additions and removals in order.
	recipes []pk.Packet
	// effects are the player's own mob effects by effect id.
	effects map[int32]pk.Packet
	players map[uuid.UUID]map[reflect.Type]client.PlayerInfo

	world    world
	inv      inventory
	self     selfState
	nextPort int32
}

// selfState is the player state the server never echoes back: the client
// predicts it and the keeper follows what the core sent.
type selfState struct {
	known      bool
	x, y, z    float64
	yaw, pitch float32
	moveFlags  int8
	held       int32
	// abilities is the server's last PlayerAbilities; flying is what the core
	// reported since.
	flying, flyingKnown bool
	input               uint8
	sprinting           bool
	usingItem           bool
	sequence            int32
}

// syntheticTeleportBase starts the teleport ids the keeper makes up for
// snapshots, far from the server's counter, so their confirmations can be
// told apart and swallowed.
const syntheticTeleportBase = 1 << 30

func newState() *state {
	s := &state{nextPort: syntheticTeleportBase}
	s.resetPlay()
	return s
}

// resetPlay forgets everything tied to one play session (a new Login).
func (s *state) resetPlay() {
	s.login = nil
	s.loaded = false
	s.singles = make(map[packetid.ClientboundPacketID]pk.Packet)
	s.events = nil
	s.recipes = nil
	s.effects = make(map[int32]pk.Packet)
	s.players = make(map[uuid.UUID]map[reflect.Type]client.PlayerInfo)
	s.world = newWorld()
	s.inv = newInventory()
	s.self = selfState{held: -1}
}

// singleKinds are the packets whose last copy is all the state they carry.
var singleKinds = map[packetid.ClientboundPacketID]bool{
	packetid.ClientboundChangeDifficulty:         true,
	packetid.ClientboundUpdateRecipes:            true,
	packetid.ClientboundRecipeBookSettings:       true,
	packetid.ClientboundCommands:                 true,
	packetid.ClientboundServerData:               true,
	packetid.ClientboundServerLinks:              true,
	packetid.ClientboundUpdateTags:               true,
	packetid.ClientboundTabList:                  true,
	packetid.ClientboundInitializeBorder:         true,
	packetid.ClientboundSetBorderCenter:          true,
	packetid.ClientboundSetBorderLerpSize:        true,
	packetid.ClientboundSetBorderSize:            true,
	packetid.ClientboundSetBorderWarningDelay:    true,
	packetid.ClientboundSetBorderWarningDistance: true,
	packetid.ClientboundSetTime:                  true,
	packetid.ClientboundSetDefaultSpawnPosition:  true,
	packetid.ClientboundGameRuleValues:           true,
	packetid.ClientboundTickingState:             true,
	packetid.ClientboundTickingStep:              true,
	packetid.ClientboundSetChunkCacheRadius:      true,
	packetid.ClientboundSetSimulationDistance:    true,
	packetid.ClientboundSetChunkCacheCenter:      true,
	packetid.ClientboundSetExperience:            true,
	packetid.ClientboundSetHealth:                true,
}

// persistentEvents are the GameEvent ids that describe lasting state.
var persistentEvents = map[uint8]bool{
	1:  true, // stop raining
	2:  true, // start raining
	3:  true, // change game mode
	7:  true, // rain level
	8:  true, // thunder level
	11: true, // immediate respawn
	12: true, // limited crafting
	13: true, // level chunks load start
}

// fromServer records a play-phase packet the server sent. It returns an error
// only for packets it should have understood; the caller logs and moves on.
func (s *state) fromServer(p pk.Packet) error {
	id := packetid.ClientboundPacketID(p.ID)
	if singleKinds[id] {
		s.singles[id] = p
		return nil
	}
	switch id {
	case packetid.ClientboundLogin:
		var login client.Login
		if err := decode(p, &login); err != nil {
			return err
		}
		s.resetPlay()
		s.login = &login
		s.playerID = login.PlayerID
		s.world.dimension = login.CommonPlayerSpawnInfo.Dimension
	case packetid.ClientboundRespawn:
		var respawn client.Respawn
		if err := decode(p, &respawn); err != nil {
			return err
		}
		if s.login != nil {
			s.login.CommonPlayerSpawnInfo = respawn.CommonPlayerSpawnInfo
		}
		if respawn.CommonPlayerSpawnInfo.Dimension != s.world.dimension {
			// Like the vanilla client (and minego): only a new dimension replaces the
			// level; respawning in the same one keeps chunks and entities.
			s.world = newWorld()
			s.world.dimension = respawn.CommonPlayerSpawnInfo.Dimension
		}
		s.loaded = false
		s.effects = make(map[int32]pk.Packet)
		s.inv.open = nil
		s.self.sprinting, s.self.usingItem = false, false
	case packetid.ClientboundGameEvent:
		var event client.GameEvent
		if err := decode(p, &event); err != nil {
			return err
		}
		if !persistentEvents[event.Event] {
			return nil
		}
		s.events = slices.DeleteFunc(s.events, func(old pk.Packet) bool {
			return p.Data[0] == old.Data[0] || (event.Event == 1 || event.Event == 2) && (old.Data[0] == 1 || old.Data[0] == 2)
		})
		s.events = append(s.events, p)
	case packetid.ClientboundSetHeldSlot:
		var held client.SetHeldItem
		if err := decode(p, &held); err != nil {
			return err
		}
		s.self.held = held.Slot
	case packetid.ClientboundPlayerAbilities:
		s.singles[id] = p
		s.self.flyingKnown = false
	case packetid.ClientboundRecipeBookAdd, packetid.ClientboundRecipeBookRemove:
		s.recipes = append(s.recipes, p)
	case packetid.ClientboundUpdateMobEffect:
		var effect client.EntityEffect
		if err := decode(p, &effect); err != nil {
			return err
		}
		if effect.EntityID == s.playerID {
			s.effects[effect.EffectID] = p
		}
	case packetid.ClientboundRemoveMobEffect:
		var effect client.RemoveMobEffect
		if err := decode(p, &effect); err != nil {
			return err
		}
		if effect.EntityID == s.playerID {
			delete(s.effects, effect.EffectID)
		}
	case packetid.ClientboundPlayerInfoUpdate:
		var update client.PlayerInfoUpdate
		if err := decode(p, &update); err != nil {
			return err
		}
		for id, infos := range update.Players {
			known := s.players[id]
			if known == nil {
				known = make(map[reflect.Type]client.PlayerInfo)
				s.players[id] = known
			}
			for _, info := range infos {
				known[reflect.TypeOf(info)] = info
			}
		}
	case packetid.ClientboundPlayerInfoRemove:
		var remove client.PlayerInfoRemove
		if err := decode(p, &remove); err != nil {
			return err
		}
		for _, id := range remove.UUIDs {
			delete(s.players, id)
		}
	case packetid.ClientboundPlayerPosition:
		var pos client.PlayerPosition
		if err := decode(p, &pos); err != nil {
			return err
		}
		s.teleport(pos)
	default:
		if handled, err := s.world.fromServer(id, p, s.playerID); handled {
			return err
		}
		if handled, err := s.inv.fromServer(id, p); handled {
			return err
		}
	}
	return nil
}

// teleport applies a server position correction the way the client does.
func (s *state) teleport(p client.PlayerPosition) {
	rel := func(bit int32, old, v float64) float64 {
		if p.Flags&bit != 0 {
			return old + v
		}
		return v
	}
	s.self.x = rel(0x01, s.self.x, p.X)
	s.self.y = rel(0x02, s.self.y, p.Y)
	s.self.z = rel(0x04, s.self.z, p.Z)
	s.self.yaw = float32(rel(0x08, float64(s.self.yaw), float64(p.XRot)))
	s.self.pitch = float32(rel(0x10, float64(s.self.pitch), float64(p.YRot)))
	s.self.known = true
}

// fromCore records a play-phase packet the core sent to the server.
func (s *state) fromCore(p pk.Packet) error {
	id := packetid.ServerboundPacketID(p.ID)
	switch id {
	case packetid.ServerboundMovePlayerPos:
		var m server.MovePlayerPos
		if err := decode(p, &m); err != nil {
			return err
		}
		s.self.x, s.self.y, s.self.z, s.self.moveFlags, s.self.known = m.X, m.FeetY, m.Z, m.Flags, true
	case packetid.ServerboundMovePlayerPosRot:
		var m server.MovePlayerPosRot
		if err := decode(p, &m); err != nil {
			return err
		}
		s.self.x, s.self.y, s.self.z, s.self.moveFlags, s.self.known = m.X, m.FeetY, m.Z, m.Flags, true
		s.self.yaw, s.self.pitch = m.XRot, m.YRot
	case packetid.ServerboundMovePlayerRot:
		var m server.MovePlayerRot
		if err := decode(p, &m); err != nil {
			return err
		}
		s.self.yaw, s.self.pitch, s.self.moveFlags = m.XRot, m.YRot, m.Flags
	case packetid.ServerboundMovePlayerStatusOnly:
		var m server.MovePlayerStatusOnly
		if err := decode(p, &m); err != nil {
			return err
		}
		s.self.moveFlags = m.Flags
	case packetid.ServerboundPlayerAbilities:
		var a server.PlayerAbilities
		if err := decode(p, &a); err != nil {
			return err
		}
		s.self.flying, s.self.flyingKnown = a.Flags&0x02 != 0, true
	case packetid.ServerboundSetCarriedItem:
		var held server.SetCarriedItem
		if err := decode(p, &held); err != nil {
			return err
		}
		s.self.held = int32(held.Slot)
		s.self.usingItem = false
	case packetid.ServerboundPlayerInput:
		var in server.PlayerInput
		if err := decode(p, &in); err != nil {
			return err
		}
		s.self.input = in.Flags
	case packetid.ServerboundPlayerCommand:
		var cmd server.PlayerCommand
		if err := decode(p, &cmd); err != nil {
			return err
		}
		switch cmd.ActionID {
		case commandStartSprint:
			s.self.sprinting = true
		case commandStopSprint:
			s.self.sprinting = false
		}
	case packetid.ServerboundUseItem:
		s.self.usingItem = true
	case packetid.ServerboundPlayerAction:
		var action server.PlayerAction
		if err := decode(p, &action); err != nil {
			return err
		}
		s.self.sequence = max(s.self.sequence, action.Sequence)
		if action.Status == actionReleaseUseItem {
			s.self.usingItem = false
		}
	case packetid.ServerboundPlayerLoaded:
		s.loaded = true
	default:
		return s.inv.fromCore(id, p)
	}
	return nil
}

// PlayerCommand and PlayerAction values the keeper reads or sends.
const (
	commandStartSprint   = 1
	commandStopSprint    = 2
	actionReleaseUseItem = 5
)

// neutral returns the packets that leave the player standing still when the
// core goes away mid-action: release keys and sprint, stop using an item and
// close an open container (the core reopens what it needs after checkpoint).
func (s *state) neutral() []server.ServerboundPacket {
	var out []server.ServerboundPacket
	if s.self.input != 0 {
		out = append(out, &server.PlayerInput{})
		s.self.input = 0
	}
	if s.self.sprinting {
		out = append(out, &server.PlayerCommand{EntityID: s.playerID, ActionID: commandStopSprint})
		s.self.sprinting = false
	}
	if s.self.usingItem {
		out = append(out, &server.PlayerAction{Status: actionReleaseUseItem, Sequence: s.self.sequence})
		s.self.usingItem = false
	}
	if s.inv.open != nil {
		out = append(out, &server.ContainerClose{WindowID: s.inv.open.id})
		s.inv.open = nil
	}
	return out
}

func decode(p pk.Packet, v pk.FieldDecoder) error {
	if _, err := v.ReadFrom(bytes.NewReader(p.Data)); err != nil {
		return fmt.Errorf("decode packet %#x: %w", p.ID, err)
	}
	return nil
}

func marshal(p client.ClientboundPacket) pk.Packet {
	return pk.Marshal(p.PacketID(), p)
}

// playerGroups splits the tab list into PlayerInfoUpdate packets whose players
// all carry the same actions, which is what the packet format requires.
func (s *state) playerGroups() []pk.Packet {
	groups := map[string]map[uuid.UUID][]client.PlayerInfo{}
	var keys []string
	for id, infos := range s.players {
		names := make([]string, 0, len(infos))
		list := make([]client.PlayerInfo, 0, len(infos))
		for t, info := range infos {
			names = append(names, t.String())
			list = append(list, info)
		}
		slices.Sort(names)
		key := strings.Join(names, ",")
		if groups[key] == nil {
			groups[key] = map[uuid.UUID][]client.PlayerInfo{}
			keys = append(keys, key)
		}
		groups[key][id] = list
	}
	slices.Sort(keys)
	out := make([]pk.Packet, 0, len(keys))
	for _, key := range keys {
		out = append(out, marshal(&client.PlayerInfoUpdate{Players: groups[key]}))
	}
	return out
}

// metadataCopy keeps entity metadata merged by index.
func metadataCopy(dst map[uint8]metadata.Metadata, src metadata.EntityMetadata) {
	for index, value := range src.Data {
		dst[index] = value
	}
}
