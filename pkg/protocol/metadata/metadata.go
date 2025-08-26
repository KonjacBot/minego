package metadata

import (
	"io"

	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
)

type MetadataType int32

const (
	MetadataByte MetadataType = iota
	MetadataVarInt
	MetadataVarLong
	MetadataFloat
	MetadataString
	MetadataChat
	MetadataOptChat
	MetadataSlot
	MetadataBoolean
	MetadataRotation
	MetadataPosition
	MetadataOptPosition
	MetadataDirection
	MetadataOptLivingEntity
	MetadataBlockState
	MetadataOptBlockState
	MetadataNBT
	MetadataParticle
	MetadataParticles
	MetadataVillagerData
	MetadataOptVarInt
	MetadataPose
	MetadataCatVariant
	MetadataCowVariant
	MetadataWolfVariant
	MetadataWolfSoundVariant
	MetadataFrogVariant
	MetadataPigVariant
	MetadataChickenVariant
	MetadataOptGlobalPosition
	MetadataPaintingVariant
	MetadataSnifferVariant
	MetadataArmadilloState
	MetadataVector3
	MetadataQuaternion
)

type Metadata interface {
	EntityMetadataType() MetadataType
	pk.Field
}

type EntityMetadata struct {
	Data map[uint8]Metadata
}

func (m EntityMetadata) WriteTo(w io.Writer) (int64, error) {
	n := int64(0)
	for u, metadata := range m.Data {
		n1, err := pk.UnsignedByte(u).WriteTo(w)
		n += n1
		if err != nil {
			return n, err
		}
		n2, err := pk.VarInt(metadata.EntityMetadataType()).WriteTo(w)
		n += n2
		if err != nil {
			return n, err
		}

		n3, err := metadata.WriteTo(w)
		n += n3
		if err != nil {
			return n, err
		}
	}
	n4, err := pk.UnsignedByte(0xff).WriteTo(w)
	n += n4
	if err != nil {
		return n, err
	}
	return n, nil
}

func (m *EntityMetadata) ReadFrom(r io.Reader) (int64, error) {
	m.Data = make(map[uint8]Metadata)
	var index uint8
	n, err := (*pk.UnsignedByte)(&index).ReadFrom(r)
	if err != nil {
		return n, err
	}
	for index != 0xff {
		var typeId MetadataType
		n1, err := (*pk.VarInt)(&typeId).ReadFrom(r)
		n += n1
		if err != nil {
			return n, err
		}

		metadata := metadataType[typeId]()
		n2, err := metadata.ReadFrom(r)
		n += n2
		if err != nil {
			return n, err
		}
		m.Data[index] = metadata

		n3, err := (*pk.UnsignedByte)(&index).ReadFrom(r)
		n += n3
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

type metadataCreator func() Metadata

var metadataType = map[MetadataType]metadataCreator{}

func init() {
	metadataType[MetadataByte] = func() Metadata { return &Byte{} }
	metadataType[MetadataVarInt] = func() Metadata { return &VarInt{} }
	metadataType[MetadataVarLong] = func() Metadata { return &VarLong{} }
	metadataType[MetadataFloat] = func() Metadata { return &Float{} }
	metadataType[MetadataString] = func() Metadata { return &String{} }
	metadataType[MetadataChat] = func() Metadata { return &Chat{} }
	metadataType[MetadataOptChat] = func() Metadata { return &OptChat{} }
	metadataType[MetadataSlot] = func() Metadata { return &Slot{} }
	metadataType[MetadataBoolean] = func() Metadata { return &Boolean{} }
	metadataType[MetadataRotation] = func() Metadata { return &Rotation{} }
	metadataType[MetadataPosition] = func() Metadata { return &Position{} }
	metadataType[MetadataOptPosition] = func() Metadata { return &OptPosition{} }
	metadataType[MetadataDirection] = func() Metadata { return &Direction{} }
	metadataType[MetadataOptLivingEntity] = func() Metadata { return &OptLivingEntity{} }
	metadataType[MetadataBlockState] = func() Metadata { return &BlockState{} }
	metadataType[MetadataOptBlockState] = func() Metadata { return &OptBlockState{} }
	metadataType[MetadataNBT] = func() Metadata { return &NBT{} }
	metadataType[MetadataParticle] = func() Metadata { return &Particle{} }
	metadataType[MetadataParticles] = func() Metadata { return &Particles{} }
	metadataType[MetadataVillagerData] = func() Metadata { return &VillagerData{} }
	metadataType[MetadataOptVarInt] = func() Metadata { return &OptVarInt{} }
	metadataType[MetadataPose] = func() Metadata { return &Pose{} }
	metadataType[MetadataCatVariant] = func() Metadata { return &CatVariant{} }
	metadataType[MetadataCowVariant] = func() Metadata { return &CowVariant{} }
	metadataType[MetadataWolfVariant] = func() Metadata { return &WolfVariant{} }
	metadataType[MetadataWolfSoundVariant] = func() Metadata { return &WolfSoundVariant{} }
	metadataType[MetadataFrogVariant] = func() Metadata { return &FrogVariant{} }
	metadataType[MetadataPigVariant] = func() Metadata { return &PigVariant{} }
	metadataType[MetadataChickenVariant] = func() Metadata { return &ChickenVariant{} }
	metadataType[MetadataOptGlobalPosition] = func() Metadata { return &OptGlobalPosition{} }
	metadataType[MetadataPaintingVariant] = func() Metadata { return &PaintingVariant{} }
	metadataType[MetadataSnifferVariant] = func() Metadata { return &SnifferVariant{} }
	metadataType[MetadataArmadilloState] = func() Metadata { return &ArmadilloState{} }
	metadataType[MetadataVector3] = func() Metadata { return &Vector3{} }
	metadataType[MetadataQuaternion] = func() Metadata { return &Quaternion{} }
}
