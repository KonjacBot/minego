package packet

import (
	"bytes"
	"fmt"
	"reflect"

	pk "github.com/KonjacBot/go-mc/net/packet"
)

// CodecPanicError reports a panic raised by a packet field implementation.
// Network data must never be able to propagate such a panic to the process.
type CodecPanicError struct {
	Operation string
	Value     any
}

func (e *CodecPanicError) Error() string {
	return fmt.Sprintf("packet %s panicked: %v", e.Operation, e.Value)
}

// Marshal encodes fields without using go-mc's panic-based packet builder.
func Marshal[ID ~int32 | ~int](id ID, fields ...pk.FieldEncoder) (result pk.Packet, err error) {
	defer func() {
		if value := recover(); value != nil {
			result = pk.Packet{}
			err = &CodecPanicError{Operation: "encode", Value: value}
		}
	}()

	var data bytes.Buffer
	for i, field := range fields {
		if isNilField(field) {
			return pk.Packet{}, fmt.Errorf("encode packet field[%d]: field is nil", i)
		}
		if _, err := field.WriteTo(&data); err != nil {
			return pk.Packet{}, fmt.Errorf("encode packet field[%d]: %w", i, err)
		}
		if data.Len() > pk.MaxDataLength {
			return pk.Packet{}, fmt.Errorf("encoded packet exceeds %d bytes", pk.MaxDataLength)
		}
	}

	return pk.Packet{ID: int32(id), Data: data.Bytes()}, nil
}

// Scan decodes every field and rejects unread trailing payload bytes.
func Scan(p pk.Packet, fields ...pk.FieldDecoder) (err error) {
	defer func() {
		if value := recover(); value != nil {
			err = &CodecPanicError{Operation: "decode", Value: value}
		}
	}()

	reader := bytes.NewReader(p.Data)
	for i, field := range fields {
		if isNilField(field) {
			return fmt.Errorf("decode packet field[%d]: field is nil", i)
		}
		if _, err := field.ReadFrom(reader); err != nil {
			return fmt.Errorf("decode packet field[%d]: %w", i, err)
		}
	}
	if reader.Len() != 0 {
		return fmt.Errorf("packet has %d trailing payload bytes", reader.Len())
	}
	return nil
}

func isNilField(field any) bool {
	if field == nil {
		return true
	}
	value := reflect.ValueOf(field)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
