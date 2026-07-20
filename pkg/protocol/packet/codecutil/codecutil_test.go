package codecutil

import (
	"bytes"
	"testing"
)

func TestBoundedStringCountsUTF16CodeUnits(t *testing.T) {
	value := "😀"
	if _, err := (BoundedString{Value: &value, MaxChars: 1}).WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted an astral character in a 1-code-unit field")
	}

	if _, err := (BoundedString{Value: &value, MaxChars: 2}).WriteTo(&bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
}

func TestBoundedStringRejectsReadLengthOverMaxCharsTimesThree(t *testing.T) {
	var decoded string
	wire := []byte{7, 'a', 'a', 'a', 'a', 'a', 'a', 'a'}
	if _, err := (BoundedString{Value: &decoded, MaxChars: 2}).ReadFrom(bytes.NewReader(wire)); err == nil {
		t.Fatal("ReadFrom() accepted UTF-8 longer than maxChars*3 bytes")
	}
}

func TestBoundedStringRejectsAstralReadWhenUTF16UnitsOverflow(t *testing.T) {
	var decoded string
	wire := []byte{4, 0xF0, 0x9F, 0x98, 0x80}
	if _, err := (BoundedString{Value: &decoded, MaxChars: 1}).ReadFrom(bytes.NewReader(wire)); err == nil {
		t.Fatal("ReadFrom() accepted an astral character in a 1-code-unit field")
	}
}

func TestBoundedStringAllowsAstralReadAtTwoCodeUnits(t *testing.T) {
	var decoded string
	wire := []byte{4, 0xF0, 0x9F, 0x98, 0x80}
	if _, err := (BoundedString{Value: &decoded, MaxChars: 2}).ReadFrom(bytes.NewReader(wire)); err != nil {
		t.Fatal(err)
	}
	if decoded != "😀" {
		t.Fatalf("decoded = %q", decoded)
	}
}
