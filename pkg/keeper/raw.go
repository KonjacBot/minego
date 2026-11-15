package keeper

import (
	"bytes"

	pk "github.com/KonjacBot/go-mc/net/packet"
)

// The snapshot reuses the server's own bytes wherever it can: decoding a field
// only to find where it ends, then copying that span. Re-encoding decoded
// values would depend on every minego type writing exactly what it reads, and
// a single asymmetric codec would make the core's decoder drop the connection.

// span decodes f from r (positioned inside data) and returns the bytes it used.
func span(data []byte, r *bytes.Reader, f pk.FieldDecoder) ([]byte, error) {
	start := len(data) - r.Len()
	if _, err := f.ReadFrom(r); err != nil {
		return nil, err
	}
	return data[start : len(data)-r.Len()], nil
}

func varint(v int32) []byte {
	var b bytes.Buffer
	_, _ = pk.VarInt(v).WriteTo(&b)
	return b.Bytes()
}
