//codec:ignore
package client

import (
	"github.com/Tnze/go-mc/data/packetid"
	"github.com/Tnze/go-mc/net/packet"
)

type ClientboundPacket interface {
	packet.Field
}

var ClientboundPackets = make(map[packetid.ClientboundPacketID]ClientboundPacket)

func init() {
	ClientboundPackets[packetid.ClientboundAddEntity] = &AddEntity{}
	ClientboundPackets[packetid.ClientboundAnimate] = &Animate{}
	ClientboundPackets[packetid.ClientboundAwardStats] = &AwardStats{}
	ClientboundPackets[packetid.ClientboundBlockChangedAck] = &BlockChangedAck{}
	ClientboundPackets[packetid.ClientboundBlockDestruction] = &BlockDestruction{}
	ClientboundPackets[packetid.ClientboundBlockEntityData] = &BlockEntityData{}
	ClientboundPackets[packetid.ClientboundBlockEvent] = &BlockEvent{}
	ClientboundPackets[packetid.ClientboundBlockUpdate] = &BlockUpdate{}
	ClientboundPackets[packetid.ClientboundBossEvent] = &BossEvent{}
	ClientboundPackets[packetid.ClientboundChangeDifficulty] = &ChangeDifficulty{}
	ClientboundPackets[packetid.ClientboundChunkBatchFinished] = &ChunkBatchFinished{}
	ClientboundPackets[packetid.ClientboundChunkBatchStart] = &ChunkBatchStart{}
	ClientboundPackets[packetid.ClientboundChunksBiomes] = &ChunkBiomes{}
	ClientboundPackets[packetid.ClientboundClearTitles] = &ClearTitles{}
	ClientboundPackets[packetid.ClientboundContainerClose] = &CloseContainer{}
	ClientboundPackets[packetid.ClientboundCommandSuggestions] = &CommandSuggestions{}
	ClientboundPackets[packetid.ClientboundCommands] = &Commands{}
	ClientboundPackets[packetid.ClientboundContainerSetData] = &ContainerSetData{}
	ClientboundPackets[packetid.ClientboundContainerSetSlot] = &ContainerSetSlot{}
	ClientboundPackets[packetid.ClientboundCooldown] = &Cooldown{}
	ClientboundPackets[packetid.ClientboundCustomChatCompletions] = &CustomChatCompletions{}
	ClientboundPackets[packetid.ClientboundDamageEvent] = &DamageEvent{}
	ClientboundPackets[packetid.ClientboundDebugSample] = &DebugSample{}
	ClientboundPackets[packetid.ClientboundDeleteChat] = &DeleteChat{}
	ClientboundPackets[packetid.ClientboundDisguisedChat] = &DisguisedChat{}
	ClientboundPackets[packetid.ClientboundEntityEvent] = &EntityEvent{}
}
