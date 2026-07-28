package slot

import (
	"bytes"
	"testing"
)

func TestEmptyAndAnyFuelDisplayRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		want DisplayType
		data Display
	}{
		{name: "nil defaults to empty", want: DisplayEmpty, data: Display{}},
		{name: "empty", want: DisplayEmpty, data: Display{SlotDisplay: &Empty{}}},
		{name: "fuel", want: DisplayAnyFuel, data: Display{SlotDisplay: &AnyFuel{}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var encoded bytes.Buffer
			written, err := test.data.WriteTo(&encoded)
			if err != nil {
				t.Fatal(err)
			}
			var got Display
			read, err := got.ReadFrom(&encoded)
			if err != nil {
				t.Fatal(err)
			}
			if written != read ||
				got.SlotDisplay == nil ||
				got.SlotDisplay.SlotDisplayType() != test.want {
				t.Fatalf(
					"display = %#v, written/read = %d/%d, want type %d",
					got,
					written,
					read,
					test.want,
				)
			}
		})
	}
}
