package metadata

import (
	"git.konjactw.dev/falloutBot/go-mc/chat"
	"git.konjactw.dev/falloutBot/go-mc/nbt"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/component"
	"git.konjactw.dev/patyhank/minego/pkg/protocol/particle"
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

type Byte struct {
	pk.Byte
}

func (b Byte) EntityMetadataType() MetadataType {
	return MetadataByte
}

type VarInt struct {
	pk.VarInt
}

func (b VarInt) EntityMetadataType() MetadataType {
	return MetadataVarInt
}

type VarLong struct {
	pk.VarLong
}

func (b VarLong) EntityMetadataType() MetadataType {
	return MetadataVarLong
}

type Float struct {
	pk.Float
}

func (b Float) EntityMetadataType() MetadataType {
	return MetadataFloat
}

type String struct {
	pk.String
}

func (b String) EntityMetadataType() MetadataType {
	return MetadataString
}

type Chat struct {
	chat.Message
}

func (b Chat) EntityMetadataType() MetadataType {
	return MetadataChat
}

type OptChat struct {
	pk.Option[chat.Message, *chat.Message]
}

func (b OptChat) EntityMetadataType() MetadataType {
	return MetadataChat
}

type Slot struct {
	slot.Slot
}

func (b Slot) EntityMetadataType() MetadataType {
	return MetadataSlot
}

type Boolean struct {
	pk.Boolean
}

func (b Boolean) EntityMetadataType() MetadataType {
	return MetadataBoolean
}

//codec:gen
type Rotation struct {
	X, Y, Z float32
}

func (b Rotation) EntityMetadataType() MetadataType {
	return MetadataRotation
}

type Position struct {
	pk.Position
}

func (b Position) EntityMetadataType() MetadataType {
	return MetadataPosition
}

type OptPosition struct {
	pk.Option[pk.Position, *pk.Position]
}

func (b OptPosition) EntityMetadataType() MetadataType {
	return MetadataOptPosition
}

type Direction struct {
	pk.VarInt
}

func (b Direction) EntityMetadataType() MetadataType {
	return MetadataDirection
}

type OptLivingEntity struct {
	pk.Option[pk.UUID, *pk.UUID]
}

func (b OptLivingEntity) EntityMetadataType() MetadataType {
	return MetadataOptLivingEntity
}

type BlockState struct {
	pk.VarInt
}

func (b BlockState) EntityMetadataType() MetadataType {
	return MetadataBlockState
}

type OptBlockState struct {
	pk.VarInt
}

func (b OptBlockState) EntityMetadataType() MetadataType {
	return MetadataOptBlockState
}

//codec:gen
type NBT struct {
	Data nbt.RawMessage `mc:"NBT"`
}

func (b NBT) EntityMetadataType() MetadataType {
	return MetadataNBT
}

type Particle struct {
	particle.Particle
}

func (b Particle) EntityMetadataType() MetadataType {
	return MetadataParticle
}

//codec:gen
type Particles struct {
	Particles []particle.Particle
}

func (b Particles) EntityMetadataType() MetadataType {
	return MetadataParticles
}

//codec:gen
type VillagerData struct {
	Type, Profession, Level int32 `mc:"VarInt"`
}

func (b VillagerData) EntityMetadataType() MetadataType {
	return MetadataVillagerData
}

type OptVarInt struct {
	pk.Option[pk.VarInt, *pk.VarInt]
}

func (b OptVarInt) EntityMetadataType() MetadataType {
	return MetadataOptVarInt
}

type Pose struct {
	pk.VarInt
}

func (b Pose) EntityMetadataType() MetadataType {
	return MetadataPose
}

type CatVariant struct {
	pk.VarInt
}

func (b CatVariant) EntityMetadataType() MetadataType {
	return MetadataCatVariant
}

type CowVariant struct {
	pk.VarInt
}

func (CowVariant) EntityMetadataType() MetadataType {
	return MetadataCowVariant
}

type WolfVariant struct {
	pk.VarInt
}

func (WolfVariant) EntityMetadataType() MetadataType {
	return MetadataWolfVariant
}

type WolfSoundVariant struct {
	pk.VarInt
}

func (WolfSoundVariant) EntityMetadataType() MetadataType {
	return MetadataWolfSoundVariant
}

type FrogVariant struct {
	pk.VarInt
}

func (FrogVariant) EntityMetadataType() MetadataType {
	return MetadataFrogVariant
}

type PigVariant struct {
	pk.VarInt
}

func (PigVariant) EntityMetadataType() MetadataType {
	return MetadataPigVariant
}

type ChickenVariant struct {
	pk.VarInt
}

func (ChickenVariant) EntityMetadataType() MetadataType {
	return MetadataChickenVariant
}

//codec:gen
type GlobalPosition struct {
	Dimension pk.Identifier
	Position  pk.Position
}

//codec:gen
type OptGlobalPosition struct {
	GlobalPosition pk.Option[GlobalPosition, *GlobalPosition]
}

func (OptGlobalPosition) EntityMetadataType() MetadataType {
	return MetadataOptGlobalPosition
}

type PaintingVariant struct {
	pk.OptID[component.PaintingVariant, *component.PaintingVariant]
}

func (PaintingVariant) EntityMetadataType() MetadataType {
	return MetadataPaintingVariant
}

type SnifferVariant struct {
	pk.VarInt
}

func (SnifferVariant) EntityMetadataType() MetadataType {
	return MetadataSnifferVariant
}

type ArmadilloState struct {
	pk.VarInt
}

func (ArmadilloState) EntityMetadataType() MetadataType {
	return MetadataArmadilloState
}

//codec:gen
type Vector3 struct {
	X, Y, Z float32
}

func (Vector3) EntityMetadataType() MetadataType {
	return MetadataVector3
}

//codec:gen
type Quaternion struct {
	X, Y, Z, W float32
}

func (Quaternion) EntityMetadataType() MetadataType {
	return MetadataQuaternion
}
