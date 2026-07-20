package codecutil

import (
	"bytes"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/KonjacBot/go-mc/net/packet"
)

type BoundedString struct {
	Value    *string
	MaxChars int
}

func (s BoundedString) WriteTo(w io.Writer) (n int64, err error) {
	if !utf8.ValidString(*s.Value) {
		return 0, fmt.Errorf("string contains invalid utf-8")
	}
	if utf16CodeUnits(*s.Value) > s.MaxChars {
		return 0, fmt.Errorf("string exceeds %d characters", s.MaxChars)
	}
	if len([]byte(*s.Value)) > s.MaxChars*3 {
		return 0, fmt.Errorf("string exceeds %d UTF-8 bytes", s.MaxChars*3)
	}
	return packet.String(*s.Value).WriteTo(w)
}

func (s BoundedString) ReadFrom(r io.Reader) (n int64, err error) {
	var length packet.VarInt
	n, err = (&length).ReadFrom(r)
	if err != nil {
		return n, err
	}
	if length < 0 {
		return n, fmt.Errorf("string length less than zero")
	}
	if int(length) > s.MaxChars*3 {
		return n, fmt.Errorf("string exceeds %d characters", s.MaxChars)
	}

	data := make([]byte, int(length))
	read, err := io.ReadFull(r, data)
	n += int64(read)
	if err != nil {
		return n, err
	}

	if !utf8.Valid(data) {
		return n, fmt.Errorf("string contains invalid utf-8")
	}
	value := string(data)
	if utf16CodeUnits(value) > s.MaxChars {
		return n, fmt.Errorf("string exceeds %d characters", s.MaxChars)
	}
	*s.Value = value
	return n, nil
}

func utf16CodeUnits(value string) int {
	count := 0
	for _, r := range value {
		if r >= 0x10000 {
			count += 2
			continue
		}
		count++
	}
	return count
}

type BoundedByteArray struct {
	Value  *[]byte
	MaxLen int
}

func (b BoundedByteArray) WriteTo(w io.Writer) (n int64, err error) {
	if len(*b.Value) > b.MaxLen {
		return 0, fmt.Errorf("byte array exceeds %d bytes", b.MaxLen)
	}

	n, err = packet.VarInt(len(*b.Value)).WriteTo(w)
	if err != nil {
		return n, err
	}

	read, err := w.Write(*b.Value)
	return n + int64(read), err
}

func (b BoundedByteArray) ReadFrom(r io.Reader) (n int64, err error) {
	var length packet.VarInt
	n, err = (&length).ReadFrom(r)
	if err != nil {
		return n, err
	}
	if length < 0 {
		return n, fmt.Errorf("byte array length less than zero")
	}
	if int(length) > b.MaxLen {
		return n, fmt.Errorf("byte array exceeds %d bytes", b.MaxLen)
	}

	data := make([]byte, int(length))
	read, err := io.ReadFull(r, data)
	n += int64(read)
	if err != nil {
		return n, err
	}
	*b.Value = data
	return n, nil
}

type RemainingBytes struct {
	Value  *[]byte
	MaxLen int
}

func (b RemainingBytes) WriteTo(w io.Writer) (n int64, err error) {
	if len(*b.Value) > b.MaxLen {
		return 0, fmt.Errorf("payload exceeds %d bytes", b.MaxLen)
	}
	read, err := w.Write(*b.Value)
	return int64(read), err
}

func (b RemainingBytes) ReadFrom(r io.Reader) (n int64, err error) {
	data, err := io.ReadAll(io.LimitReader(r, int64(b.MaxLen+1)))
	if err != nil {
		return 0, err
	}
	if len(data) > b.MaxLen {
		return int64(len(data)), fmt.Errorf("payload exceeds %d bytes", b.MaxLen)
	}
	*b.Value = data
	return int64(len(data)), nil
}

type LengthPrefixedNBT struct {
	Value  any
	MaxLen int
}

func (n LengthPrefixedNBT) WriteTo(w io.Writer) (written int64, err error) {
	var payload bytes.Buffer
	_, err = packet.NBT(n.Value).WriteTo(&payload)
	if err != nil {
		return 0, err
	}
	if payload.Len() > n.MaxLen {
		return 0, fmt.Errorf("nbt payload exceeds %d bytes", n.MaxLen)
	}

	written, err = packet.VarInt(payload.Len()).WriteTo(w)
	if err != nil {
		return written, err
	}
	read, err := payload.WriteTo(w)
	return written + read, err
}

func (n LengthPrefixedNBT) ReadFrom(r io.Reader) (read int64, err error) {
	var length packet.VarInt
	read, err = (&length).ReadFrom(r)
	if err != nil {
		return read, err
	}
	if length < 0 {
		return read, fmt.Errorf("nbt payload length less than zero")
	}
	if int(length) > n.MaxLen {
		return read, fmt.Errorf("nbt payload exceeds %d bytes", n.MaxLen)
	}

	data := make([]byte, int(length))
	count, err := io.ReadFull(r, data)
	read += int64(count)
	if err != nil {
		return read, err
	}
	inner := bytes.NewReader(data)
	_, err = packet.NBT(n.Value).ReadFrom(inner)
	if err != nil {
		return read, err
	}
	if inner.Len() != 0 {
		return read, fmt.Errorf("nbt payload has %d trailing bytes", inner.Len())
	}
	return read, err
}
