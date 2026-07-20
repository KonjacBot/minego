package client

import (
	"fmt"
	"io"

	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/packet/codecutil"
)

const (
	maxCustomReportDetailCount = 32
	maxCustomReportKeyChars    = 128
	maxCustomReportValueChars  = 4096
	maxResourcePackHashChars   = 40
	maxCookiePayloadBytes      = 5120
	maxRemainingPayloadBytes   = 1 << 20
)

type boundedReportDetails map[string]string

func (d boundedReportDetails) WriteTo(w io.Writer) (n int64, err error) {
	if len(d) > maxCustomReportDetailCount {
		return 0, fmt.Errorf("detail count exceeds %d", maxCustomReportDetailCount)
	}

	n, err = pk.VarInt(len(d)).WriteTo(w)
	if err != nil {
		return n, err
	}
	for key, value := range d {
		var temp int64
		temp, err = codecutil.BoundedString{Value: &key, MaxChars: maxCustomReportKeyChars}.WriteTo(w)
		n += temp
		if err != nil {
			return n, err
		}
		temp, err = codecutil.BoundedString{Value: &value, MaxChars: maxCustomReportValueChars}.WriteTo(w)
		n += temp
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (d *boundedReportDetails) ReadFrom(r io.Reader) (n int64, err error) {
	var size pk.VarInt
	n, err = (&size).ReadFrom(r)
	if err != nil {
		return n, err
	}
	if size < 0 {
		return n, fmt.Errorf("detail count less than zero")
	}
	if int(size) > maxCustomReportDetailCount {
		return n, fmt.Errorf("detail count exceeds %d", maxCustomReportDetailCount)
	}

	values := make(map[string]string, int(size))
	for i := 0; i < int(size); i++ {
		var key, value string
		var temp int64
		temp, err = codecutil.BoundedString{Value: &key, MaxChars: maxCustomReportKeyChars}.ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
		temp, err = codecutil.BoundedString{Value: &value, MaxChars: maxCustomReportValueChars}.ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
		values[key] = value
	}
	*d = values
	return n, nil
}
