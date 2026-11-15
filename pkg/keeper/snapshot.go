package keeper

import (
	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/client"
)

// snapshot returns the play-phase packets that bring a freshly configured
// client to the current state, in roughly the order a server sends them on
// join: world entry, player settings, position, tab list, world settings,
// chunks, entities, inventory, vitals. teleport is the id of the synthetic
// position packet; the core's confirmation of it must not reach the server.
//
// It returns nil before the server's Login arrived.
func (s *state) snapshot() (packets []pk.Packet, teleport int32, ok bool) {
	if s.login == nil {
		return nil, 0, false
	}
	var out []pk.Packet
	single := func(ids ...packetid.ClientboundPacketID) {
		for _, id := range ids {
			if p, ok := s.singles[id]; ok {
				out = append(out, p)
			}
		}
	}

	out = append(out, pk.Packet{ID: int32(packetid.ClientboundLogin), Data: s.login})
	if s.respawn != nil {
		out = append(out, pk.Packet{ID: int32(packetid.ClientboundRespawn), Data: s.respawn})
	}
	single(packetid.ClientboundChangeDifficulty)
	if p, ok := s.abilities(); ok {
		out = append(out, p)
	}
	if s.self.held >= 0 {
		out = append(out, marshal(&client.SetHeldItem{Slot: s.self.held}))
	}
	single(packetid.ClientboundUpdateRecipes, packetid.ClientboundRecipeBookSettings)
	out = append(out, s.recipes...)
	single(packetid.ClientboundCommands, packetid.ClientboundServerData, packetid.ClientboundServerLinks,
		packetid.ClientboundUpdateTags)

	if s.self.known {
		teleport = s.nextPort
		s.nextPort++
		out = append(out, marshal(&client.PlayerPosition{
			ID: teleport, X: s.self.x, Y: s.self.y, Z: s.self.z, XRot: s.self.yaw, YRot: s.self.pitch,
		}))
	}

	out = append(out, s.playerGroups()...)
	single(packetid.ClientboundTabList,
		packetid.ClientboundInitializeBorder, packetid.ClientboundSetBorderCenter, packetid.ClientboundSetBorderLerpSize,
		packetid.ClientboundSetBorderSize, packetid.ClientboundSetBorderWarningDelay, packetid.ClientboundSetBorderWarningDistance,
		packetid.ClientboundSetTime, packetid.ClientboundSetDefaultSpawnPosition, packetid.ClientboundGameRuleValues,
		packetid.ClientboundTickingState, packetid.ClientboundTickingStep)
	out = append(out, s.events...)
	single(packetid.ClientboundSetChunkCacheRadius, packetid.ClientboundSetSimulationDistance,
		packetid.ClientboundSetChunkCacheCenter)

	out = append(out, s.world.packets(s.playerID)...)
	out = append(out, s.inv.packets()...)
	single(packetid.ClientboundSetExperience, packetid.ClientboundSetHealth)
	for _, p := range s.effects {
		out = append(out, p)
	}
	if s.inBundle {
		// The snapshot falls inside a bundle: open one so the closing delimiter
		// that arrives live pairs up instead of starting a bundle of its own.
		out = append(out, pk.Packet{ID: int32(packetid.BundleDelimiter)})
	}
	return out, teleport, true
}

// abilities is the server's last PlayerAbilities with the flying bit the core
// reported since, so a bot that was flying keeps flying.
func (s *state) abilities() (pk.Packet, bool) {
	p, ok := s.singles[packetid.ClientboundPlayerAbilities]
	if !ok {
		return pk.Packet{}, false
	}
	if !s.self.flyingKnown {
		return p, true
	}
	var a client.PlayerAbilities
	if err := decode(p, &a); err != nil {
		return p, true
	}
	if s.self.flying {
		a.Flags |= 0x02
	} else {
		a.Flags &^= 0x02
	}
	return marshal(&a), true
}
