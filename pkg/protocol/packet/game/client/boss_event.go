package client

import (
	"io"

	"github.com/google/uuid"

	"git.konjactw.dev/falloutBot/go-mc/chat"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
)

// codec:gen
type BossEvent struct {
	UUID        uuid.UUID `mc:"UUID"`
	BossEventOp BossEventOperation
}

// codec:gen
type BossEventOpAdd struct {
	Message  chat.Message
	Progress float32
	Color    int32 `mc:"VarInt"`
	Overlay  int32 `mc:"VarInt"`
	Flags    uint8
}

// codec:gen
type BossEventOpRemove struct {
}

// codec:gen
type BossEventOpUpdate struct {
	Progress float32
}

// codec:gen
type BossEventOpUpdateName struct {
	Message chat.Message
}

// codec:gen
type BossEventOpUpdateStyle struct {
	Color   int32 `mc:"VarInt"`
	Overlay int32 `mc:"VarInt"`
}

// codec:gen
type BossEventOpUpdateProperties struct {
	Flags uint8
}

type BossEventOperation struct {
	Operation int32 `mc:"VarInt"`
	OpData    any
}

func (c *BossEventOperation) ReadFrom(r io.Reader) (int64, error) {
	o, err := (*pk.VarInt)(&c.Operation).ReadFrom(r)
	if err != nil {
		return o, err
	}
	var i int64

	switch c.Operation {
	case 0:
		var op BossEventOpAdd
		i, err = op.ReadFrom(r)
		if err != nil {
			return i, err
		}
		c.OpData = &op
	case 1:
		var op BossEventOpRemove
		i, err = op.ReadFrom(r)
		if err != nil {
			return i, err
		}
		c.OpData = &op
	case 2:
		var op BossEventOpUpdate
		i, err = op.ReadFrom(r)
		if err != nil {
			return i, err
		}
		c.OpData = &op
	case 3:
		var op BossEventOpUpdateName
		i, err = op.ReadFrom(r)
		if err != nil {
			return i, err
		}
		c.OpData = &op
	case 4:
		var op BossEventOpUpdateStyle
		i, err = op.ReadFrom(r)
		if err != nil {
			return i, err
		}
		c.OpData = &op
	case 5:
		var op BossEventOpUpdateProperties
		i, err = op.ReadFrom(r)
		if err != nil {
			return i, err
		}
		c.OpData = &op
	}

	return o + i, nil
}

func (c *BossEventOperation) WriteTo(w io.Writer) (int64, error) {
	o, err := (*pk.VarInt)(&c.Operation).WriteTo(w)
	if err != nil {
		return o, err
	}

	switch c.Operation {
	case 0:
		op := c.OpData.(*BossEventOpAdd)
		i, err := op.WriteTo(w)
		if err != nil {
			return o + i, err
		}
		return o + i, nil
	case 1:
		op := c.OpData.(*BossEventOpRemove)
		i, err := op.WriteTo(w)
		if err != nil {
			return o + i, err
		}
		return o + i, nil
	case 2:
		op := c.OpData.(*BossEventOpUpdate)
		i, err := op.WriteTo(w)
		if err != nil {
			return o + i, err
		}
		return o + i, nil
	case 3:
		op := c.OpData.(*BossEventOpUpdateName)
		i, err := op.WriteTo(w)
		if err != nil {
			return o + i, err
		}
		return o + i, nil
	case 4:
		op := c.OpData.(*BossEventOpUpdateStyle)
		i, err := op.WriteTo(w)
		if err != nil {
			return o + i, err
		}
		return o + i, nil
	case 5:
		op := c.OpData.(*BossEventOpUpdateProperties)
		i, err := op.WriteTo(w)
		if err != nil {
			return o + i, err
		}
		return o + i, nil
	}

	return o, nil
}
