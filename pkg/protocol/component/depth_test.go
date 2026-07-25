package component

import (
	"bytes"
	"testing"

	pk "github.com/KonjacBot/go-mc/net/packet"
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

func TestItemStackTemplateRejectsExcessiveComponentNesting(t *testing.T) {
	useRemainderID := -1
	for id := 0; ; id++ {
		value := slot.ComponentFromID(id)
		if value == nil {
			break
		}
		if value.ID() == "minecraft:use_remainder" {
			useRemainderID = id
			break
		}
	}
	if useRemainderID < 0 {
		t.Fatal("use_remainder component is not registered")
	}

	var data bytes.Buffer
	write := func(value int) {
		t.Helper()
		if _, err := pk.VarInt(value).WriteTo(&data); err != nil {
			t.Fatal(err)
		}
	}
	for range 65 {
		write(1) // item ID
		write(1) // count
		write(1) // one added component
		write(0) // no removed components
		write(useRemainderID)
	}
	write(1)
	write(1)
	write(0)
	write(0)

	var value slot.ItemStackTemplate
	if _, err := value.ReadFrom(bytes.NewReader(data.Bytes())); err == nil {
		t.Fatal("ReadFrom() accepted excessively nested item components")
	}
}
