package metadata

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
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

type entityMetadata interface {
	EntityMetadataType() MetadataType
	pk.Field
}

type EntityMetadata struct {
	Data map[uint8]entityMetadata
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

type metadataCreator func() entityMetadata

var metadataType = map[MetadataType]metadataCreator{}

func init() {
	metadataType[MetadataByte] = func() entityMetadata { return &Byte{} }
	metadataType[MetadataVarInt] = func() entityMetadata { return &VarInt{} }
	metadataType[MetadataVarLong] = func() entityMetadata { return &VarLong{} }
	metadataType[MetadataFloat] = func() entityMetadata { return &Float{} }
	metadataType[MetadataString] = func() entityMetadata { return &String{} }
	metadataType[MetadataChat] = func() entityMetadata { return &Chat{} }
	metadataType[MetadataOptChat] = func() entityMetadata { return &OptChat{} }
	metadataType[MetadataSlot] = func() entityMetadata { return &Slot{} }
	metadataType[MetadataBoolean] = func() entityMetadata { return &Boolean{} }
	metadataType[MetadataRotation] = func() entityMetadata { return &Rotation{} }
	metadataType[MetadataPosition] = func() entityMetadata { return &Position{} }
	metadataType[MetadataOptPosition] = func() entityMetadata { return &OptPosition{} }
	metadataType[MetadataDirection] = func() entityMetadata { return &Direction{} }
	metadataType[MetadataOptLivingEntity] = func() entityMetadata { return &OptLivingEntity{} }
	metadataType[MetadataBlockState] = func() entityMetadata { return &BlockState{} }
	metadataType[MetadataOptBlockState] = func() entityMetadata { return &OptBlockState{} }
	metadataType[MetadataNBT] = func() entityMetadata { return &NBT{} }
	metadataType[MetadataParticle] = func() entityMetadata { return &Particle{} }
	metadataType[MetadataParticles] = func() entityMetadata { return &Particles{} }
	metadataType[MetadataVillagerData] = func() entityMetadata { return &VillagerData{} }
	metadataType[MetadataOptVarInt] = func() entityMetadata { return &OptVarInt{} }
	metadataType[MetadataPose] = func() entityMetadata { return &Pose{} }
	metadataType[MetadataCatVariant] = func() entityMetadata { return &CatVariant{} }
	metadataType[MetadataCowVariant] = func() entityMetadata { return &CowVariant{} }
	metadataType[MetadataWolfVariant] = func() entityMetadata { return &WolfVariant{} }
	metadataType[MetadataWolfSoundVariant] = func() entityMetadata { return &WolfSoundVariant{} }
	metadataType[MetadataFrogVariant] = func() entityMetadata { return &FrogVariant{} }
	metadataType[MetadataPigVariant] = func() entityMetadata { return &PigVariant{} }
	metadataType[MetadataChickenVariant] = func() entityMetadata { return &ChickenVariant{} }
	metadataType[MetadataOptGlobalPosition] = func() entityMetadata { return &OptGlobalPosition{} }
	metadataType[MetadataPaintingVariant] = func() entityMetadata { return &PaintingVariant{} }
	metadataType[MetadataSnifferVariant] = func() entityMetadata { return &SnifferVariant{} }
	metadataType[MetadataArmadilloState] = func() entityMetadata { return &ArmadilloState{} }
	metadataType[MetadataVector3] = func() entityMetadata { return &Vector3{} }
	metadataType[MetadataQuaternion] = func() entityMetadata { return &Quaternion{} }
}
