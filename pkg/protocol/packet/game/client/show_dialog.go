package client

import (
	"fmt"
	"io"
	"math"

	"github.com/KonjacBot/go-mc/nbt"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

// ShowDialog matches the official 26.2 holder wire shape used by
// ClientboundShowDialogPacket: VarInt 0 means an inline direct Dialog follows,
// otherwise the value is a one-based registry holder ID.
type ShowDialog struct {
	HasRegistryID bool
	RegistryID    int32
	DialogData    nbt.RawMessage `mc:"NBT"`
}

func (c *ShowDialog) ReadFrom(r io.Reader) (n int64, err error) {
	var wireID pk.VarInt
	var temp int64
	temp, err = (&wireID).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}

	if wireID == 0 {
		c.HasRegistryID = false
		c.RegistryID = 0
		c.DialogData = nbt.RawMessage{}
		temp, err = pk.NBT(&c.DialogData).ReadFrom(r)
		n += temp
		return n, err
	}

	c.HasRegistryID = true
	c.RegistryID = int32(wireID - 1)
	c.DialogData = nbt.RawMessage{}
	return n, nil
}

func (c ShowDialog) WriteTo(w io.Writer) (n int64, err error) {
	var wireID pk.VarInt
	var temp int64
	if c.HasRegistryID {
		if c.RegistryID < 0 {
			return 0, fmt.Errorf("registry ID must be non-negative")
		}
		if c.RegistryID >= math.MaxInt32 {
			return 0, fmt.Errorf("registry ID exceeds maximum encodable value")
		}
		wireID = pk.VarInt(c.RegistryID + 1)
	}
	temp, err = (&wireID).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	if c.HasRegistryID {
		return n, nil
	}
	temp, err = pk.NBT(&c.DialogData).WriteTo(w)
	n += temp
	return n, err
}
