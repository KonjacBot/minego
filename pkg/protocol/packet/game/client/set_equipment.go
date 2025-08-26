package client

import (
	"io"

	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

type EquipmentData struct {
	Slot int8
	Item slot.Slot
}

type Equipment []EquipmentData

func (e Equipment) WriteTo(w io.Writer) (n int64, err error) {
	for i, equipment := range e {
		b := equipment.Slot
		if len(e)-1 == i {
			b |= -128
		}

		n1, err := pk.Byte(b).WriteTo(w)
		n += n1
		if err != nil {
			return n, err
		}
		n2, err := equipment.Item.WriteTo(w)
		n += n2
		if err != nil {
			return n, err
		}
	}
	return
}

func (e *Equipment) ReadFrom(r io.Reader) (n int64, err error) {
	for {
		var b pk.Byte
		n1, err := b.ReadFrom(r)
		n += n1
		if err != nil {
			return n, err
		}

		var equipment EquipmentData
		equipment.Slot = int8(b & 127)
		n2, err := equipment.Item.ReadFrom(r)
		n += n2
		if err != nil {
			return n, err
		}

		*e = append(*e, equipment)
		if n&-128 == 0 {
			break
		}
	}

	return
}

//codec:gen
type SetEquipment struct {
	EntityID int32 `mc:"VarInt"`
	Equipment
}
