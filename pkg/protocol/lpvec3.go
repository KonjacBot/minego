package protocol

import (
	"encoding/binary"
	"io"
	"math"

	pk "github.com/KonjacBot/go-mc/net/packet"
)

type LpVec3 struct {
	X, Y, Z float64
}

const (
	lpMaxVal = 1.7179869183e10
	lpMinVal = 3.051944088384301e-05
)

// ReadFrom 實作 io.ReaderFrom 介面
func (v *LpVec3) ReadFrom(r io.Reader) (n int64, err error) {
	// 1. 讀取第一個 byte (lowest)
	var lowestBuf [1]byte
	if _, err := io.ReadFull(r, lowestBuf[:]); err != nil {
		return 0, err
	}
	n += 1
	lowest := uint32(lowestBuf[0])

	if lowest == 0 {
		v.X, v.Y, v.Z = 0, 0, 0
		return n, nil
	}

	// 2. 讀取 middle (1 byte) 與 highest (4 bytes)
	var remain [5]byte
	if _, err := io.ReadFull(r, remain[:]); err != nil {
		return n, err
	}
	n += 5

	middle := uint32(remain[0])
	highest := binary.BigEndian.Uint32(remain[1:])

	// 重組 64-bit buffer
	buffer := uint64(highest)<<16 | uint64(middle)<<8 | uint64(lowest)

	// 3. 處理 Scale 與 Continuation Bit
	scaleVal := uint64(lowest & 3)
	if (lowest & 4) == 4 {
		var vInt pk.VarInt
		vn, err := vInt.ReadFrom(r)
		n += vn
		if err != nil {
			return n, err
		}
		// Java: scale |= (VarInt.read(input) & 4294967295L) << 2
		scaleVal |= (uint64(int32(vInt)) & 0xFFFFFFFF) << 2
	}
	scale := float64(scaleVal)

	// 4. 解包座標
	v.X = lpUnpack(int64(buffer>>3)) * scale
	v.Y = lpUnpack(int64(buffer>>18)) * scale
	v.Z = lpUnpack(int64(buffer>>33)) * scale

	return n, nil
}

// WriteTo 實作 io.WriterTo 介面
func (v *LpVec3) WriteTo(w io.Writer) (n int64, err error) {
	x, y, z := lpSanitize(v.X), lpSanitize(v.Y), lpSanitize(v.Z)
	chessboard := math.Max(math.Abs(x), math.Max(math.Abs(y), math.Abs(z)))

	if chessboard < lpMinVal {
		nn, err := w.Write([]byte{0})
		return int64(nn), err
	}

	scale := int64(math.Ceil(chessboard))
	isPartial := (scale & 3) != scale

	var markers int64
	if isPartial {
		markers = (scale & 3) | 4
	} else {
		markers = scale
	}

	// 封裝位元組
	xn := lpPack(x/float64(scale)) << 3
	yn := lpPack(y/float64(scale)) << 18
	zn := lpPack(z/float64(scale)) << 33
	buffer := uint64(markers | xn | yn | zn)

	// 寫入 6 bytes (lowest, middle, highest[4])
	out := make([]byte, 6)
	out[0] = byte(buffer)
	out[1] = byte(buffer >> 8)
	binary.BigEndian.PutUint32(out[2:], uint32(buffer>>16))

	nn, err := w.Write(out)
	n += int64(nn)
	if err != nil {
		return n, err
	}

	// 處理 Partial Scale
	if isPartial {
		vn, err := pk.VarInt(scale >> 2).WriteTo(w)
		n += vn
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

// 內部轉換工具
func lpSanitize(v float64) float64 {
	if math.IsNaN(v) {
		return 0
	}
	if v > lpMaxVal {
		return lpMaxVal
	}
	if v < -lpMaxVal {
		return -lpMaxVal
	}
	return v
}

func lpPack(v float64) int64 {
	return int64(math.Round((v*0.5 + 0.5) * 32766.0))
}

func lpUnpack(v int64) float64 {
	val := float64(v & 32767)
	if val > 32766.0 {
		val = 32766.0
	}
	return val*2.0/32766.0 - 1.0
}
