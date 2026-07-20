package client

import (
	"bytes"
	"fmt"
	"io"

	"github.com/KonjacBot/go-mc/data/registryid"
	"github.com/KonjacBot/go-mc/net/packet"
)

const (
	commandParserFloat            = 1
	commandParserDouble           = 2
	commandParserInteger          = 3
	commandParserLong             = 4
	commandParserString           = 5
	commandParserEntity           = 6
	commandParserScoreHolder      = 31
	commandParserTime             = 43
	commandParserResourceOrTag    = 44
	commandParserResourceOrTagKey = 45
	commandParserResource         = 46
	commandParserResourceKey      = 47
	commandParserResourceSelector = 48
)

func (c *ChunkBiomeData) ReadFrom(r io.Reader) (n int64, err error) {
	var temp int64
	temp, err = (&c.Pos).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	c.Data = nil
	temp, err = (*packet.ByteArray)(&c.Data).ReadFrom(r)
	n += temp
	return n, err
}

func (c ChunkBiomeData) WriteTo(w io.Writer) (n int64, err error) {
	var temp int64
	temp, err = (&c.Pos).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = (*packet.ByteArray)(&c.Data).WriteTo(w)
	n += temp
	return n, err
}

func (c *CommandNode) ReadFrom(r io.Reader) (n int64, err error) {
	var temp int64
	temp, err = (*packet.Byte)(&c.Flags).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	c.Children = nil
	temp, err = (*Int32VarIntVarIntArray)(&c.Children).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	flags := byte(c.Flags)
	if flags&0x08 != 0 {
		temp, err = (*packet.VarInt)(&c.Redirect).ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
	} else {
		c.Redirect = 0
	}

	switch flags & 0x03 {
	case 0:
		c.Name = ""
		c.Parser = CommandParser{}
		c.SuggestionType = ""
	case 1:
		temp, err = (*packet.String)(&c.Name).ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
		c.Parser = CommandParser{}
		c.SuggestionType = ""
	case 2:
		temp, err = (*packet.String)(&c.Name).ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
		temp, err = (&c.Parser).ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
		if flags&0x10 != 0 {
			temp, err = (*packet.Identifier)(&c.SuggestionType).ReadFrom(r)
			n += temp
			if err != nil {
				return n, err
			}
		} else {
			c.SuggestionType = ""
		}
	default:
		return n, fmt.Errorf("unsupported command node type %#x", flags&0x03)
	}
	return n, nil
}

func (c CommandNode) WriteTo(w io.Writer) (n int64, err error) {
	var temp int64
	flags := byte(c.Flags)
	temp, err = (*packet.Byte)(&c.Flags).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = (*Int32VarIntVarIntArray)(&c.Children).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	if flags&0x08 != 0 {
		temp, err = (*packet.VarInt)(&c.Redirect).WriteTo(w)
		n += temp
		if err != nil {
			return n, err
		}
	}

	switch flags & 0x03 {
	case 0:
		return n, nil
	case 1:
		temp, err = (*packet.String)(&c.Name).WriteTo(w)
		n += temp
		return n, err
	case 2:
		temp, err = (*packet.String)(&c.Name).WriteTo(w)
		n += temp
		if err != nil {
			return n, err
		}
		temp, err = (&c.Parser).WriteTo(w)
		n += temp
		if err != nil {
			return n, err
		}
		if flags&0x10 != 0 {
			temp, err = (*packet.Identifier)(&c.SuggestionType).WriteTo(w)
			n += temp
		}
		return n, err
	default:
		return n, fmt.Errorf("unsupported command node type %#x", flags&0x03)
	}
}

func (c *CommandParser) ReadFrom(r io.Reader) (n int64, err error) {
	var temp int64
	temp, err = (*packet.VarInt)(&c.ID).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	if c.ID < 0 || c.ID >= int32(len(registryid.CommandArgumentType)) {
		return n, fmt.Errorf("unknown command parser id %d", c.ID)
	}
	c.Properties, temp, err = readCommandParserProperties(c.ID, r)
	n += temp
	return n, err
}

func (c CommandParser) WriteTo(w io.Writer) (n int64, err error) {
	var temp int64
	temp, err = (*packet.VarInt)(&c.ID).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = writeRawBytes(w, c.Properties)
	n += temp
	return n, err
}

func (c *WaypointIcon) ReadFrom(r io.Reader) (n int64, err error) {
	var temp int64
	temp, err = (*packet.Identifier)(&c.Style).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = (&c.Color).ReadFrom(r)
	n += temp
	return n, err
}

func (c WaypointIcon) WriteTo(w io.Writer) (n int64, err error) {
	var temp int64
	temp, err = (*packet.Identifier)(&c.Style).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = (&c.Color).WriteTo(w)
	n += temp
	return n, err
}

func readCommandParserProperties(id int32, r io.Reader) ([]byte, int64, error) {
	var raw bytes.Buffer
	var n int64
	writeField := func(field packet.Field) error {
		_, err := field.WriteTo(&raw)
		return err
	}

	switch id {
	case commandParserFloat:
		return readNumericCommandParserProperties(r, &raw, &n, 4)
	case commandParserDouble:
		return readNumericCommandParserProperties(r, &raw, &n, 8)
	case commandParserInteger:
		return readNumericCommandParserProperties(r, &raw, &n, 4)
	case commandParserLong:
		return readNumericCommandParserProperties(r, &raw, &n, 8)
	case commandParserString:
		var stringType int32
		temp, err := (*packet.VarInt)(&stringType).ReadFrom(r)
		n += temp
		if err != nil {
			return nil, n, err
		}
		if err := writeField((*packet.VarInt)(&stringType)); err != nil {
			return nil, n, err
		}
	case commandParserEntity, commandParserScoreHolder:
		var flags int8
		temp, err := (*packet.Byte)(&flags).ReadFrom(r)
		n += temp
		if err != nil {
			return nil, n, err
		}
		if err := writeField((*packet.Byte)(&flags)); err != nil {
			return nil, n, err
		}
	case commandParserTime:
		var min int32
		temp, err := (*packet.Int)(&min).ReadFrom(r)
		n += temp
		if err != nil {
			return nil, n, err
		}
		if err := writeField((*packet.Int)(&min)); err != nil {
			return nil, n, err
		}
	case commandParserResourceOrTag, commandParserResourceOrTagKey, commandParserResource, commandParserResourceKey, commandParserResourceSelector:
		var registry string
		temp, err := (*packet.Identifier)(&registry).ReadFrom(r)
		n += temp
		if err != nil {
			return nil, n, err
		}
		if err := writeField((*packet.Identifier)(&registry)); err != nil {
			return nil, n, err
		}
	}

	return raw.Bytes(), n, nil
}

func readNumericCommandParserProperties(r io.Reader, raw *bytes.Buffer, n *int64, boundSize int) ([]byte, int64, error) {
	var flags int8
	temp, err := (*packet.Byte)(&flags).ReadFrom(r)
	*n += temp
	if err != nil {
		return nil, *n, err
	}
	if _, err := (*packet.Byte)(&flags).WriteTo(raw); err != nil {
		return nil, *n, err
	}
	if flags&0x01 != 0 {
		temp, err = readCommandParserBound(r, raw, boundSize)
		*n += temp
		if err != nil {
			return nil, *n, err
		}
	}
	if flags&0x02 != 0 {
		temp, err = readCommandParserBound(r, raw, boundSize)
		*n += temp
		if err != nil {
			return nil, *n, err
		}
	}
	return raw.Bytes(), *n, nil
}

func readCommandParserBound(r io.Reader, raw *bytes.Buffer, size int) (int64, error) {
	switch size {
	case 4:
		var value int32
		n, err := (*packet.Int)(&value).ReadFrom(r)
		if err != nil {
			return n, err
		}
		_, err = (*packet.Int)(&value).WriteTo(raw)
		return n, err
	case 8:
		var value int64
		n, err := (*packet.Long)(&value).ReadFrom(r)
		if err != nil {
			return n, err
		}
		_, err = (*packet.Long)(&value).WriteTo(raw)
		return n, err
	default:
		return 0, fmt.Errorf("unsupported command parser bound size %d", size)
	}
}

func writeRawBytes(w io.Writer, data []byte) (int64, error) {
	if len(data) == 0 {
		return 0, nil
	}
	n, err := w.Write(data)
	return int64(n), err
}
