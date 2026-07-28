package client

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestRecipeBookSettingsWireLayout(t *testing.T) {
	want := RecipeBookSettings{
		CraftingRecipeBookOpen:                 true,
		CraftingRecipeBookFilterEnabled:        false,
		SmeltingRecipeBookOpen:                 true,
		SmeltingRecipeBookFilterEnabled:        false,
		BlastingFurnaceRecipeBookOpen:          true,
		BlastingFurnaceRecipeBookFilterEnabled: false,
		SmokingRecipeBookOpen:                  true,
		SmokingRecipeBookFilterEnabled:         false,
	}

	var encoded bytes.Buffer
	written, err := want.WriteTo(&encoded)
	if err != nil {
		t.Fatalf("WriteTo() error = %v", err)
	}
	if written != 8 || encoded.Len() != 8 {
		t.Fatalf(
			"RecipeBookSettings encoded size = %d/%d, want exactly 8 booleans",
			written,
			encoded.Len(),
		)
	}
	if got := encoded.Bytes(); !bytes.Equal(
		got,
		[]byte{1, 0, 1, 0, 1, 0, 1, 0},
	) {
		t.Fatalf("RecipeBookSettings bytes = %v", got)
	}

	var decoded RecipeBookSettings
	read, err := decoded.ReadFrom(&encoded)
	if err != nil {
		t.Fatalf("ReadFrom() error = %v", err)
	}
	if read != 8 || !reflect.DeepEqual(decoded, want) {
		t.Fatalf("ReadFrom() = (%+v, %d), want (%+v, 8)", decoded, read, want)
	}
	if _, err := decoded.ReadFrom(
		bytes.NewReader([]byte{1, 0, 1, 0, 1, 0, 1}),
	); err == nil || err != io.EOF {
		t.Fatalf("ReadFrom() truncated error = %v, want io.EOF", err)
	}
}
