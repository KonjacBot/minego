package recipe

import (
	"bytes"
	"io"
	"testing"
)

type testRecipeDisplay struct{}

func (*testRecipeDisplay) RecipeType() DisplayType { return DisplayCraftingShapeless }

func (*testRecipeDisplay) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte{42})
	return int64(n), err
}

func (*testRecipeDisplay) ReadFrom(io.Reader) (int64, error) { return 0, nil }

func TestDisplayWriteReportsAllBytes(t *testing.T) {
	var encoded bytes.Buffer
	written, err := (Display{Display: &testRecipeDisplay{}}).WriteTo(&encoded)
	if err != nil {
		t.Fatal(err)
	}
	if written != int64(encoded.Len()) || !bytes.Equal(encoded.Bytes(), []byte{0, 42}) {
		t.Fatalf("WriteTo() = %d bytes, encoded %v", written, encoded.Bytes())
	}
}

func TestDisplayRejectsInvalidValue(t *testing.T) {
	if _, err := (Display{}).WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("WriteTo() accepted a nil recipe display")
	}

	display := Display{Display: &testRecipeDisplay{}}
	read, err := display.ReadFrom(bytes.NewReader([]byte{99}))
	if err == nil {
		t.Fatal("ReadFrom() accepted an unknown recipe display type")
	}
	if read != 1 || display.Display != nil {
		t.Fatalf("ReadFrom() = (%d, %v), display = %T", read, err, display.Display)
	}
}
