//codec:ignore
package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	"git.konjactw.dev/falloutBot/go-mc/net/packet"
)

type ClientboundPacket interface {
	packet.Field
	PacketID() packetid.ClientboundPacketID
}

type clientBoundPacketCreator func() ClientboundPacket

var ClientboundPackets = make(map[packetid.ClientboundPacketID]clientBoundPacketCreator)

func init() {
	registerPacket(func() ClientboundPacket {
		return &AddEntity{}
	})
	registerPacket(func() ClientboundPacket {
		return &Animate{}
	})
	registerPacket(func() ClientboundPacket {
		return &AwardStats{}
	})
	registerPacket(func() ClientboundPacket {
		return &BlockChangedAck{}
	})
	registerPacket(func() ClientboundPacket {
		return &BlockDestruction{}
	})
	registerPacket(func() ClientboundPacket {
		return &BlockEntityData{}
	})
	registerPacket(func() ClientboundPacket {
		return &BlockEvent{}
	})
	registerPacket(func() ClientboundPacket {
		return &BlockUpdate{}
	})
	registerPacket(func() ClientboundPacket {
		return &BossEvent{}
	})
	registerPacket(func() ClientboundPacket {
		return &ChangeDifficulty{}
	})
	registerPacket(func() ClientboundPacket {
		return &ChunkBatchFinished{}
	})
	registerPacket(func() ClientboundPacket {
		return &ChunkBatchStart{}
	})
	registerPacket(func() ClientboundPacket {
		return &ChunkBiomes{}
	})
	registerPacket(func() ClientboundPacket {
		return &ClearTitles{}
	})
	registerPacket(func() ClientboundPacket {
		return &CommandSuggestions{}
	})
	registerPacket(func() ClientboundPacket {
		return &Commands{}
	})
	registerPacket(func() ClientboundPacket {
		return &CloseContainer{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetContainerContent{}
	})
	registerPacket(func() ClientboundPacket {
		return &ContainerSetData{}
	})
	registerPacket(func() ClientboundPacket {
		return &ContainerSetSlot{}
	})
	registerPacket(func() ClientboundPacket {
		return &CookieRequest{}
	})
	registerPacket(func() ClientboundPacket {
		return &Cooldown{}
	})
	registerPacket(func() ClientboundPacket {
		return &CustomChatCompletions{}
	})
	registerPacket(func() ClientboundPacket {
		return &CustomPayload{}
	})
	registerPacket(func() ClientboundPacket {
		return &DamageEvent{}
	})
	registerPacket(func() ClientboundPacket {
		return &DebugSample{}
	})
	registerPacket(func() ClientboundPacket {
		return &DeleteChat{}
	})
	registerPacket(func() ClientboundPacket {
		return &Disconnect{}
	})
	registerPacket(func() ClientboundPacket {
		return &DisguisedChat{}
	})
	registerPacket(func() ClientboundPacket {
		return &EntityEvent{}
	})
	registerPacket(func() ClientboundPacket {
		return &TeleportEntity{}
	})
	registerPacket(func() ClientboundPacket {
		return &Explode{}
	})
	registerPacket(func() ClientboundPacket {
		return &ForgetLevelChunk{}
	})
	registerPacket(func() ClientboundPacket {
		return &GameEvent{}
	})
	registerPacket(func() ClientboundPacket {
		return &OpenHorseScreen{}
	})
	registerPacket(func() ClientboundPacket {
		return &HurtAnimation{}
	})
	registerPacket(func() ClientboundPacket {
		return &InitializeWorldBorder{}
	})
	registerPacket(func() ClientboundPacket {
		return &KeepAlive{}
	})
	registerPacket(func() ClientboundPacket {
		return &LevelChunkWithLight{}
	})
	registerPacket(func() ClientboundPacket {
		return &LevelEvent{}
	})
	registerPacket(func() ClientboundPacket {
		return &Particle{}
	})
	registerPacket(func() ClientboundPacket {
		return &UpdateLight{}
	})
	registerPacket(func() ClientboundPacket {
		return &Login{}
	})
	registerPacket(func() ClientboundPacket {
		return &MapData{}
	})
	registerPacket(func() ClientboundPacket {
		return &MerchantOffers{}
	})
	registerPacket(func() ClientboundPacket {
		return &UpdateEntityPosition{}
	})
	registerPacket(func() ClientboundPacket {
		return &UpdateEntityPositionAndRotation{}
	})
	registerPacket(func() ClientboundPacket {
		return &MoveMinecartAlongTrack{}
	})
	registerPacket(func() ClientboundPacket {
		return &UpdateEntityRotation{}
	})
	registerPacket(func() ClientboundPacket {
		return &MoveVehicle{}
	})
	registerPacket(func() ClientboundPacket {
		return &OpenBook{}
	})
	registerPacket(func() ClientboundPacket {
		return &OpenScreen{}
	})
	registerPacket(func() ClientboundPacket {
		return &OpenSignEditor{}
	})
	registerPacket(func() ClientboundPacket {
		return &Ping{}
	})
	registerPacket(func() ClientboundPacket {
		return &PingResponse{}
	})
	registerPacket(func() ClientboundPacket {
		return &PlaceGhostRecipe{}
	})
	registerPacket(func() ClientboundPacket {
		return &PlayerAbilities{}
	})
	registerPacket(func() ClientboundPacket {
		return &EndCombat{}
	})
	registerPacket(func() ClientboundPacket {
		return &EnterCombat{}
	})
	registerPacket(func() ClientboundPacket {
		return &CombatDeath{}
	})
	registerPacket(func() ClientboundPacket {
		return &PlayerInfoRemove{}
	})
	registerPacket(func() ClientboundPacket {
		return &PlayerInfoUpdate{}
	})
	registerPacket(func() ClientboundPacket {
		return &LookAt{}
	})
	registerPacket(func() ClientboundPacket {
		return &PlayerPosition{}
	})
	registerPacket(func() ClientboundPacket {
		return &PlayerRotation{}
	})
	registerPacket(func() ClientboundPacket {
		return &RecipeBookAdd{}
	})
	registerPacket(func() ClientboundPacket {
		return &RecipeBookRemove{}
	})
	registerPacket(func() ClientboundPacket {
		return &RecipeBookSettings{}
	})
	registerPacket(func() ClientboundPacket {
		return &RemoveEntities{}
	})
	registerPacket(func() ClientboundPacket {
		return &RemoveMobEffect{}
	})
	registerPacket(func() ClientboundPacket {
		return &ResetScore{}
	})
	registerPacket(func() ClientboundPacket {
		return &AddResourcePack{}
	})
	registerPacket(func() ClientboundPacket {
		return &RemoveResourcePack{}
	})
	registerPacket(func() ClientboundPacket {
		return &Respawn{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetHeadRotation{}
	})
	registerPacket(func() ClientboundPacket {
		return &UpdateSectionsBlocks{}
	})
	registerPacket(func() ClientboundPacket {
		return &SelectAdvancementsTab{}
	})
	registerPacket(func() ClientboundPacket {
		return &ServerData{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetActionBarText{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetBorderCenter{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetBorderLerpSize{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetBorderSize{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetBorderWarningDelay{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetBorderWarningDistance{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetCamera{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetCenterChunk{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetRenderDistance{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetCursorItem{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetDefaultSpawnPosition{}
	})
	registerPacket(func() ClientboundPacket {
		return &DisplayObjective{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetEntityMetadata{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetEntityLink{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetEntityVelocity{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetEquipment{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetExperience{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetHealth{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetHeldItem{}
	})
	registerPacket(func() ClientboundPacket {
		return &UpdateObjectives{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetPassengers{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetPlayerInventory{}
	})
	registerPacket(func() ClientboundPacket {
		return &UpdateTeams{}
	})
	registerPacket(func() ClientboundPacket {
		return &UpdateScore{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetSimulationDistance{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetSubtitleText{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetTime{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetTitleText{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetTitleAnimationTimes{}
	})
	registerPacket(func() ClientboundPacket {
		return &EntitySoundEffect{}
	})
	registerPacket(func() ClientboundPacket {
		return &SoundEffect{}
	})
	registerPacket(func() ClientboundPacket {
		return &StartConfiguration{}
	})
	registerPacket(func() ClientboundPacket {
		return &StopSound{}
	})
	registerPacket(func() ClientboundPacket {
		return &StoreCookie{}
	})
	registerPacket(func() ClientboundPacket {
		return &SystemChatMessage{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetTabListHeaderAndFooter{}
	})
	registerPacket(func() ClientboundPacket {
		return &TagQueryResponse{}
	})
	registerPacket(func() ClientboundPacket {
		return &PickupItem{}
	})
	registerPacket(func() ClientboundPacket {
		return &SynchronizeVehiclePosition{}
	})
	registerPacket(func() ClientboundPacket {
		return &TestInstanceBlockStatus{}
	})
	registerPacket(func() ClientboundPacket {
		return &SetTickingState{}
	})
	registerPacket(func() ClientboundPacket {
		return &StepTick{}
	})
	registerPacket(func() ClientboundPacket {
		return &Transfer{}
	})
	registerPacket(func() ClientboundPacket {
		return &UpdateAdvancements{}
	})
	registerPacket(func() ClientboundPacket {
		return &UpdateAttributes{}
	})
	registerPacket(func() ClientboundPacket {
		return &EntityEffect{}
	})
	registerPacket(func() ClientboundPacket {
		return &UpdateRecipes{}
	})
	registerPacket(func() ClientboundPacket {
		return &UpdateTags{}
	})
	registerPacket(func() ClientboundPacket {
		return &ProjectilePower{}
	})
	registerPacket(func() ClientboundPacket {
		return &CustomReportDetails{}
	})
	registerPacket(func() ClientboundPacket {
		return &ServerLinks{}
	})
	registerPacket(func() ClientboundPacket {
		return &Waypoint{}
	})
	registerPacket(func() ClientboundPacket {
		return &ClearDialog{}
	})
	registerPacket(func() ClientboundPacket {
		return &ShowDialog{}
	})
}

func registerPacket(creator clientBoundPacketCreator) {
	ClientboundPackets[creator().PacketID()] = creator
}

func (*AddEntity) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundAddEntity
}
func (*Animate) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundAnimate
}
func (*AwardStats) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundAwardStats
}
func (*BlockChangedAck) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundBlockChangedAck
}
func (*BlockDestruction) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundBlockDestruction
}
func (*BlockEntityData) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundBlockEntityData
}
func (*BlockEvent) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundBlockEvent
}
func (*BlockUpdate) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundBlockUpdate
}
func (*BossEvent) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundBossEvent
}
func (*ChangeDifficulty) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundChangeDifficulty
}
func (*ChunkBatchFinished) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundChunkBatchFinished
}
func (*ChunkBatchStart) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundChunkBatchStart
}
func (*ChunkBiomes) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundChunksBiomes
}
func (*ClearTitles) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundClearTitles
}
func (*CommandSuggestions) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundCommandSuggestions
}
func (*Commands) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundCommands
}
func (*CloseContainer) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundContainerClose
}
func (*SetContainerContent) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundContainerSetContent
}
func (*ContainerSetData) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundContainerSetData
}
func (*ContainerSetSlot) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundContainerSetSlot
}
func (*CookieRequest) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundCookieRequest
}
func (*Cooldown) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundCooldown
}
func (*CustomChatCompletions) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundCustomChatCompletions
}
func (*CustomPayload) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundCustomPayload
}
func (*DamageEvent) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundDamageEvent
}
func (*DebugSample) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundDebugSample
}
func (*DeleteChat) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundDeleteChat
}
func (*Disconnect) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundDisconnect
}
func (*DisguisedChat) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundDisguisedChat
}
func (*EntityEvent) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundEntityEvent
}
func (*TeleportEntity) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundEntityPositionSync
}
func (*Explode) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundExplode
}
func (*ForgetLevelChunk) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundForgetLevelChunk
}
func (*GameEvent) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundGameEvent
}
func (*OpenHorseScreen) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundHorseScreenOpen
}
func (*HurtAnimation) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundHurtAnimation
}
func (*InitializeWorldBorder) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundInitializeBorder
}
func (*KeepAlive) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundKeepAlive
}
func (*LevelChunkWithLight) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundLevelChunkWithLight
}
func (*LevelEvent) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundLevelEvent
}
func (*Particle) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundLevelParticles
}
func (*UpdateLight) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundLightUpdate
}
func (*Login) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundLogin
}
func (*MapData) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundMapItemData
}
func (*MerchantOffers) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundMerchantOffers
}
func (*UpdateEntityPosition) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundMoveEntityPos
}
func (*UpdateEntityPositionAndRotation) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundMoveEntityPosRot
}
func (*MoveMinecartAlongTrack) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundMoveMinecartAlongTrack
}
func (*UpdateEntityRotation) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundMoveEntityRot
}
func (*MoveVehicle) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundMoveVehicle
}
func (*OpenBook) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundOpenBook
}
func (*OpenScreen) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundOpenScreen
}
func (*OpenSignEditor) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundOpenSignEditor
}
func (*Ping) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundPing
}
func (*PingResponse) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundPongResponse
}
func (*PlaceGhostRecipe) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundPlaceGhostRecipe
}
func (*PlayerAbilities) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundPlayerAbilities
}
func (*EndCombat) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundPlayerCombatEnd
}
func (*EnterCombat) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundPlayerCombatEnter
}
func (*CombatDeath) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundPlayerCombatKill
}
func (*PlayerInfoRemove) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundPlayerInfoRemove
}
func (*PlayerInfoUpdate) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundPlayerInfoUpdate
}
func (*LookAt) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundPlayerLookAt
}
func (*PlayerPosition) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundPlayerPosition
}
func (*PlayerRotation) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundPlayerRotation
}
func (*RecipeBookAdd) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundRecipeBookAdd
}
func (*RecipeBookRemove) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundRecipeBookRemove
}
func (*RecipeBookSettings) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundRecipeBookSettings
}
func (*RemoveEntities) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundRemoveEntities
}
func (*RemoveMobEffect) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundRemoveMobEffect
}
func (*ResetScore) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundResetScore
}
func (*AddResourcePack) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundResourcePackPop
}
func (*RemoveResourcePack) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundResourcePackPush
}
func (*Respawn) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundRespawn
}
func (*SetHeadRotation) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundRotateHead
}
func (*UpdateSectionsBlocks) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSectionBlocksUpdate
}
func (*SelectAdvancementsTab) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSelectAdvancementsTab
}
func (*ServerData) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundServerData
}
func (*SetActionBarText) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetActionBarText
}
func (*SetBorderCenter) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetBorderCenter
}
func (*SetBorderLerpSize) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetBorderLerpSize
}
func (*SetBorderSize) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetBorderSize
}
func (*SetBorderWarningDelay) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetBorderWarningDelay
}
func (*SetBorderWarningDistance) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetBorderWarningDistance
}
func (*SetCamera) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetCamera
}
func (*SetCenterChunk) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetChunkCacheCenter
}
func (*SetRenderDistance) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetChunkCacheRadius
}
func (*SetCursorItem) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetCursorItem
}
func (*SetDefaultSpawnPosition) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetDefaultSpawnPosition
}
func (*DisplayObjective) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetDisplayObjective
}
func (*SetEntityMetadata) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetEntityData
}
func (*SetEntityLink) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetEntityLink
}
func (*SetEntityVelocity) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetEntityMotion
}
func (*SetEquipment) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetEquipment
}
func (*SetExperience) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetExperience
}
func (*SetHealth) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetHealth
}
func (*SetHeldItem) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetHeldSlot
}
func (*UpdateObjectives) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetObjective
}
func (*SetPassengers) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetPassengers
}
func (*SetPlayerInventory) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetPlayerInventory
}
func (*UpdateTeams) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetPlayerTeam
}
func (*UpdateScore) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetScore
}
func (*SetSimulationDistance) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetSimulationDistance
}
func (*SetSubtitleText) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetSubtitleText
}
func (*SetTime) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetTime
}
func (*SetTitleText) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetTitleText
}
func (*SetTitleAnimationTimes) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSetTitlesAnimation
}
func (*EntitySoundEffect) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSoundEntity
}
func (*SoundEffect) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSound
}
func (*StartConfiguration) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundStartConfiguration
}
func (*StopSound) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundStopSound
}
func (*StoreCookie) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundStoreCookie
}
func (*SystemChatMessage) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundSystemChat
}
func (*SetTabListHeaderAndFooter) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundTabList
}
func (*TagQueryResponse) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundTagQuery
}
func (*PickupItem) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundTakeItemEntity
}
func (*SynchronizeVehiclePosition) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundTeleportEntity
}
func (*TestInstanceBlockStatus) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundTestInstanceBlockStatus
}
func (*SetTickingState) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundTickingState
}
func (*StepTick) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundTickingStep
}
func (*Transfer) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundTransfer
}
func (*UpdateAdvancements) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundUpdateAdvancements
}
func (*UpdateAttributes) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundUpdateAttributes
}
func (*EntityEffect) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundUpdateMobEffect
}
func (*UpdateRecipes) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundUpdateRecipes
}
func (*UpdateTags) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundUpdateTags
}
func (*ProjectilePower) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundProjectilePower
}
func (*CustomReportDetails) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundCustomReportDetails
}
func (*ServerLinks) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundServerLinks
}
func (*Waypoint) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundWaypoint
}
func (*ClearDialog) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundClearDialog
}
func (*ShowDialog) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundShowDialog
}
