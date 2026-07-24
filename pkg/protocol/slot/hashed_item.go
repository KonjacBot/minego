package slot

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"sort"
)

//codec:gen
type AddedHashedComponent struct {
	Type     int32 `mc:"VarInt"`
	DataHash int32
}

//codec:gen
type HashedSlot struct {
	HasItem bool
	//opt:optional:HasItem
	ItemID int32 `mc:"VarInt"`
	//opt:optional:HasItem
	ItemCount int32 `mc:"VarInt"`
	//opt:optional:HasItem
	AddComponents []AddedHashedComponent
	//opt:optional:HasItem
	RemovedComponents []int32 `mc:"VarInt"`
}

var componentHashTable = crc32.MakeTable(crc32.Castagnoli)

// HashSlot converts an authoritative item stack to the compact prediction
// format used by the serverbound container-click packet.
func HashSlot(value Slot) (HashedSlot, error) {
	if value.Count <= 0 {
		return HashedSlot{}, nil
	}
	hashed := HashedSlot{
		HasItem:           true,
		ItemID:            int32(value.ItemID),
		ItemCount:         value.Count,
		RemovedComponents: append([]int32(nil), value.RemoveComponent...),
	}
	componentIDs := make([]int, 0, len(value.AddComponent))
	for id := range value.AddComponent {
		componentIDs = append(componentIDs, int(id))
	}
	sort.Ints(componentIDs)
	for _, componentID := range componentIDs {
		component := value.AddComponent[int32(componentID)]
		if component == nil {
			return HashedSlot{}, fmt.Errorf("hash slot component %d: nil component", componentID)
		}
		var encoded bytes.Buffer
		if _, err := component.WriteTo(&encoded); err != nil {
			return HashedSlot{}, fmt.Errorf("hash slot component %d: %w", componentID, err)
		}
		hashed.AddComponents = append(hashed.AddComponents, AddedHashedComponent{
			Type:     int32(componentID),
			DataHash: int32(crc32.Checksum(encoded.Bytes(), componentHashTable)),
		})
	}
	return hashed, nil
}
