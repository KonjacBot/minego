package component

import (
	"bytes"
	"testing"
)

func TestWritableBookContentWireFormat(t *testing.T) {
	wire := []byte{1, 1, 'p', 0}
	var content WritableBookContent

	if n, err := content.ReadFrom(bytes.NewReader(wire)); err != nil || n != int64(len(wire)) {
		t.Fatalf("ReadFrom() = (%d, %v), want (%d, nil)", n, err, len(wire))
	}
	if content.Pages[0].RawContent != "p" || content.Pages[0].FilteredContent.Has {
		t.Fatalf("unexpected decoded content: %+v", content)
	}
}

func TestWrittenBookContentWireFormat(t *testing.T) {
	wire := []byte{1, 't', 0, 1, 'a', 0, 0, 1}
	var content WrittenBookContent

	if n, err := content.ReadFrom(bytes.NewReader(wire)); err != nil || n != int64(len(wire)) {
		t.Fatalf("ReadFrom() = (%d, %v), want (%d, nil)", n, err, len(wire))
	}
	if content.RawTitle != "t" || content.Author != "a" || !content.Resolved {
		t.Fatalf("unexpected decoded content: %+v", content)
	}
}
