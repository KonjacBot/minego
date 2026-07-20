package server

import (
	"fmt"
	"io"

	"github.com/KonjacBot/go-mc/chat"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/packet/codecutil"
)

const (
	maxArgumentSignatureCount = 8
	maxArgumentNameChars      = 16
	maxChatMessageChars       = 256
	maxClientLanguageChars    = 16
	maxCommandSuggestionChars = 32500
	maxContainerChangedSlots  = 128
	maxBookPages              = 100
	maxBookPageChars          = 1024
	maxBookTitleChars         = 32
	maxSignLineChars          = 384
	maxCustomClickNBTBytes    = 65536
)

type boundedArgumentSignatures []SignedSignatures

func (a boundedArgumentSignatures) WriteTo(w io.Writer) (n int64, err error) {
	if len(a) > maxArgumentSignatureCount {
		return 0, fmt.Errorf("argument signature count exceeds %d", maxArgumentSignatureCount)
	}

	n, err = pk.VarInt(len(a)).WriteTo(w)
	if err != nil {
		return n, err
	}
	for i := range a {
		var temp int64
		temp, err = a[i].WriteTo(w)
		n += temp
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (a *boundedArgumentSignatures) ReadFrom(r io.Reader) (n int64, err error) {
	var size pk.VarInt
	n, err = (&size).ReadFrom(r)
	if err != nil {
		return n, err
	}
	if size < 0 {
		return n, fmt.Errorf("argument signature count less than zero")
	}
	if int(size) > maxArgumentSignatureCount {
		return n, fmt.Errorf("argument signature count exceeds %d", maxArgumentSignatureCount)
	}

	values := make([]SignedSignatures, int(size))
	for i := range values {
		var temp int64
		temp, err = values[i].ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
	}
	*a = values
	return n, nil
}

type boundedChangedSlots []ChangedSlot

func (s boundedChangedSlots) WriteTo(w io.Writer) (n int64, err error) {
	if len(s) > maxContainerChangedSlots {
		return 0, fmt.Errorf("changed slot count exceeds %d", maxContainerChangedSlots)
	}

	n, err = pk.VarInt(len(s)).WriteTo(w)
	if err != nil {
		return n, err
	}
	for i := range s {
		var temp int64
		temp, err = s[i].WriteTo(w)
		n += temp
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (s *boundedChangedSlots) ReadFrom(r io.Reader) (n int64, err error) {
	var size pk.VarInt
	n, err = (&size).ReadFrom(r)
	if err != nil {
		return n, err
	}
	if size < 0 {
		return n, fmt.Errorf("changed slot count less than zero")
	}
	if int(size) > maxContainerChangedSlots {
		return n, fmt.Errorf("changed slot count exceeds %d", maxContainerChangedSlots)
	}

	values := make([]ChangedSlot, int(size))
	for i := range values {
		var temp int64
		temp, err = values[i].ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
	}
	*s = values
	return n, nil
}

type limitedBookEntries []string

func (e limitedBookEntries) WriteTo(w io.Writer) (n int64, err error) {
	if len(e) > maxBookPages {
		return 0, fmt.Errorf("book page count exceeds %d", maxBookPages)
	}

	n, err = pk.VarInt(len(e)).WriteTo(w)
	if err != nil {
		return n, err
	}
	for i := range e {
		var temp int64
		temp, err = codecutil.BoundedString{Value: &e[i], MaxChars: maxBookPageChars}.WriteTo(w)
		n += temp
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (e *limitedBookEntries) ReadFrom(r io.Reader) (n int64, err error) {
	var size pk.VarInt
	n, err = (&size).ReadFrom(r)
	if err != nil {
		return n, err
	}
	if size < 0 {
		return n, fmt.Errorf("book page count less than zero")
	}
	if int(size) > maxBookPages {
		return n, fmt.Errorf("book page count exceeds %d", maxBookPages)
	}

	values := make([]string, int(size))
	for i := range values {
		var temp int64
		temp, err = codecutil.BoundedString{Value: &values[i], MaxChars: maxBookPageChars}.ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
	}
	*e = values
	return n, nil
}

type TestInstanceBlockVec3i struct {
	X, Y, Z int32 `mc:"VarInt"`
}

func (v *TestInstanceBlockVec3i) ReadFrom(r io.Reader) (n int64, err error) {
	var temp int64
	temp, err = (*pk.VarInt)(&v.X).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = (*pk.VarInt)(&v.Y).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = (*pk.VarInt)(&v.Z).ReadFrom(r)
	n += temp
	return n, err
}

func (v TestInstanceBlockVec3i) WriteTo(w io.Writer) (n int64, err error) {
	var temp int64
	temp, err = (*pk.VarInt)(&v.X).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = (*pk.VarInt)(&v.Y).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = (*pk.VarInt)(&v.Z).WriteTo(w)
	n += temp
	return n, err
}

type TestInstanceBlockData struct {
	Test           pk.Option[pk.Identifier, *pk.Identifier]
	Size           TestInstanceBlockVec3i
	Rotation       int32 `mc:"VarInt"`
	IgnoreEntities bool
	Status         int32 `mc:"VarInt"`
	ErrorMessage   pk.Option[chat.Message, *chat.Message]
}

func (d *TestInstanceBlockData) ReadFrom(r io.Reader) (n int64, err error) {
	var temp int64
	temp, err = (&d.Test).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = (&d.Size).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = (*pk.VarInt)(&d.Rotation).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = (*pk.Boolean)(&d.IgnoreEntities).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = (*pk.VarInt)(&d.Status).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = (&d.ErrorMessage).ReadFrom(r)
	n += temp
	return n, err
}

func (d TestInstanceBlockData) WriteTo(w io.Writer) (n int64, err error) {
	var temp int64
	temp, err = d.Test.WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = d.Size.WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = (*pk.VarInt)(&d.Rotation).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = (*pk.Boolean)(&d.IgnoreEntities).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = (*pk.VarInt)(&d.Status).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = d.ErrorMessage.WriteTo(w)
	n += temp
	return n, err
}
