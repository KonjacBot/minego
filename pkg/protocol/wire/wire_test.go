package wire

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/KonjacBot/go-mc/nbt"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

func TestStringRejectsNegativeLength(t *testing.T) {
	var value String
	if _, err := value.ReadFrom(bytes.NewReader([]byte{0xff, 0xff, 0xff, 0xff, 0x0f})); err == nil {
		t.Fatal("ReadFrom() accepted a negative string length")
	}
}

func TestStringRejectsTruncatedValueBeforeAllocation(t *testing.T) {
	var value String
	data := appendVarInt(nil, MaxStringBytes)
	if _, err := value.ReadFrom(bytes.NewReader(data)); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("ReadFrom() error = %v, want io.ErrUnexpectedEOF", err)
	}
}

func TestArrayRejectsOversizedCount(t *testing.T) {
	var values []pk.VarInt
	data := appendVarInt(nil, MaxCollectionEntries+1)
	if _, err := Array(&values).ReadFrom(bytes.NewReader(data)); err == nil {
		t.Fatal("ReadFrom() accepted an oversized array")
	}
	if values != nil {
		t.Fatalf("oversized array allocated %d values", len(values))
	}
}

func TestArrayZeroLengthClearsReusedValue(t *testing.T) {
	values := []pk.VarInt{1, 2, 3}
	if _, err := Array(&values).ReadFrom(bytes.NewReader([]byte{0})); err != nil {
		t.Fatal(err)
	}
	if values != nil {
		t.Fatalf("zero-length array = %#v, want nil", values)
	}
}

func TestIDSetDecodesIDsRatherThanByteCounts(t *testing.T) {
	var value IDSet
	if _, err := value.ReadFrom(bytes.NewReader([]byte{3, 5, 7})); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value.IDs, []int32{5, 7}) {
		t.Fatalf("IDs = %v, want [5 7]", value.IDs)
	}
}

func TestNBTRejectsNegativeArrayLength(t *testing.T) {
	var value any
	data := []byte{nbt.TagByteArray, 0xff, 0xff, 0xff, 0xff}
	if _, err := NBT(&value).ReadFrom(bytes.NewReader(data)); err == nil {
		t.Fatal("ReadFrom() accepted a negative NBT array length")
	}
}

func TestNBTRejectsExcessiveNesting(t *testing.T) {
	data := []byte{nbt.TagCompound}
	for range maxNBTDepth + 1 {
		data = append(data, nbt.TagCompound, 0, 0)
	}
	for range maxNBTDepth + 2 {
		data = append(data, nbt.TagEnd)
	}

	var value any
	if _, err := NBT(&value).ReadFrom(bytes.NewReader(data)); err == nil {
		t.Fatal("ReadFrom() accepted excessively nested NBT")
	}
}

func appendVarInt(dst []byte, value int) []byte {
	var encoded [pk.MaxVarIntLen]byte
	length := pk.VarInt(value).WriteToBytes(encoded[:])
	return append(dst, encoded[:length]...)
}
