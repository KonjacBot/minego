package slot

import (
	"bytes"
	"hash/crc32"
	"io"
	"reflect"
	"testing"

	"github.com/KonjacBot/go-mc/level/item"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

type hashedTestComponent struct {
	Value int32
}

func (hashedTestComponent) ID() string { return "test" }
func (c *hashedTestComponent) ReadFrom(r io.Reader) (int64, error) {
	return (*pk.Int)(&c.Value).ReadFrom(r)
}
func (c hashedTestComponent) WriteTo(w io.Writer) (int64, error) {
	return pk.Int(c.Value).WriteTo(w)
}

func TestHashedSlotCodecSupportsComponentArrays(t *testing.T) {
	want := HashedSlot{
		HasItem: true, ItemID: 7, ItemCount: 2,
		AddComponents:     []AddedHashedComponent{{Type: 3, DataHash: 11}, {Type: 8, DataHash: 12}},
		RemovedComponents: []int32{4, 9},
	}
	var encoded bytes.Buffer
	if _, err := want.WriteTo(&encoded); err != nil {
		t.Fatal(err)
	}
	var got HashedSlot
	if _, err := got.ReadFrom(&encoded); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("decoded hashed slot = %#v, want %#v", got, want)
	}
}

func TestHashSlotUsesCRC32CAndStableComponentOrder(t *testing.T) {
	first := &hashedTestComponent{Value: 17}
	second := &hashedTestComponent{Value: 23}
	got, err := HashSlot(Slot{
		Count: 2, ItemID: item.ID(7),
		AddComponent:    map[int32]Component{8: second, 3: first},
		RemoveComponent: []int32{4},
	})
	if err != nil {
		t.Fatal(err)
	}
	componentBytes := []byte{0, 0, 0, 17}
	wantHash := int32(crc32.Checksum(componentBytes, crc32.MakeTable(crc32.Castagnoli)))
	if len(got.AddComponents) != 2 || got.AddComponents[0].Type != 3 || got.AddComponents[0].DataHash != wantHash {
		t.Fatalf("added component hashes = %#v", got.AddComponents)
	}
	if !got.HasItem || got.ItemID != 7 || got.ItemCount != 2 || !reflect.DeepEqual(got.RemovedComponents, []int32{4}) {
		t.Fatalf("hashed slot = %#v", got)
	}
}
