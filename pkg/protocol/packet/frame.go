package packet

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"

	pk "github.com/KonjacBot/go-mc/net/packet"
)

// ReadFrame decodes one Minecraft packet frame with strict frame and
// decompression boundaries. It intentionally does not rely on dependency
// replace directives, which are ignored when minego is consumed as a module.
func ReadFrame(r io.Reader, threshold int, result *pk.Packet) error {
	if result == nil {
		return fmt.Errorf("packet destination is nil")
	}
	if threshold < 0 {
		return readUncompressedFrame(r, result)
	}
	return readCompressedFrame(r, threshold, result)
}

func readUncompressedFrame(r io.Reader, result *pk.Packet) error {
	var frameLength pk.VarInt
	if _, err := frameLength.ReadFrom(r); err != nil {
		return err
	}
	if frameLength < 1 || frameLength > pk.MaxDataLength+pk.MaxVarIntLen {
		return fmt.Errorf("invalid packet frame length %d", frameLength)
	}
	frame := &io.LimitedReader{R: r, N: int64(frameLength)}
	var packetID pk.VarInt
	idBytes, err := packetID.ReadFrom(frame)
	if err != nil {
		return err
	}
	payloadLength := int(frameLength) - int(idBytes)
	if payloadLength < 0 || payloadLength > pk.MaxDataLength {
		return fmt.Errorf("invalid packet payload length %d", payloadLength)
	}
	if err := readPayload(frame, result, int32(packetID), payloadLength); err != nil {
		return err
	}
	if frame.N != 0 {
		return fmt.Errorf("packet frame has %d unread bytes", frame.N)
	}
	return nil
}

func readCompressedFrame(r io.Reader, threshold int, result *pk.Packet) error {
	var frameLength pk.VarInt
	if _, err := frameLength.ReadFrom(r); err != nil {
		return err
	}
	if frameLength < 1 || frameLength > pk.MaxDataLength+pk.MaxVarIntLen*2 {
		return fmt.Errorf("invalid compressed frame length %d", frameLength)
	}
	frame := &io.LimitedReader{R: r, N: int64(frameLength)}
	var dataLength pk.VarInt
	dataLengthBytes, err := dataLength.ReadFrom(frame)
	if err != nil {
		return err
	}

	if dataLength == 0 {
		var packetID pk.VarInt
		idBytes, err := packetID.ReadFrom(frame)
		if err != nil {
			return err
		}
		payloadLength := int(frameLength) - int(dataLengthBytes) - int(idBytes)
		if payloadLength < 0 || payloadLength > pk.MaxDataLength {
			return fmt.Errorf("invalid uncompressed payload length %d", payloadLength)
		}
		if int(idBytes)+payloadLength >= threshold {
			return fmt.Errorf("uncompressed packet size %d is not below threshold %d", int(idBytes)+payloadLength, threshold)
		}
		if err := readPayload(frame, result, int32(packetID), payloadLength); err != nil {
			return err
		}
		if frame.N != 0 {
			return fmt.Errorf("compressed frame has %d unread bytes", frame.N)
		}
		return nil
	}

	if dataLength < 0 || dataLength > pk.MaxDataLength {
		return fmt.Errorf("invalid decompressed packet length %d", dataLength)
	}
	if int(dataLength) < threshold {
		return fmt.Errorf("compressed packet size %d is below threshold %d", dataLength, threshold)
	}

	compressed, err := zlib.NewReader(singleByteReader{Reader: frame})
	if err != nil {
		return err
	}
	defer compressed.Close()
	var packetID pk.VarInt
	idBytes, err := packetID.ReadFrom(compressed)
	if err != nil {
		return err
	}
	payloadLength := int(dataLength) - int(idBytes)
	if payloadLength < 0 || payloadLength > pk.MaxDataLength {
		return fmt.Errorf("invalid decompressed payload length %d", payloadLength)
	}
	if err := readPayload(compressed, result, int32(packetID), payloadLength); err != nil {
		return err
	}
	var extra [1]byte
	read, extraErr := compressed.Read(extra[:])
	if read != 0 {
		return fmt.Errorf("compressed packet has extra decompressed data")
	}
	if extraErr != io.EOF {
		return extraErr
	}
	if frame.N != 0 {
		return fmt.Errorf("compressed frame has %d trailing bytes", frame.N)
	}
	return nil
}

func readPayload(r io.Reader, result *pk.Packet, id int32, length int) error {
	if cap(result.Data) < length {
		result.Data = make([]byte, length)
	} else {
		result.Data = result.Data[:length]
	}
	result.ID = id
	_, err := io.ReadFull(r, result.Data)
	return err
}

// WriteFrame writes one bounded Minecraft packet frame.
func WriteFrame(w io.Writer, threshold int, value pk.Packet) error {
	if len(value.Data) > pk.MaxDataLength {
		return fmt.Errorf("packet payload exceeds %d bytes", pk.MaxDataLength)
	}
	var packetData bytes.Buffer
	if _, err := pk.VarInt(value.ID).WriteTo(&packetData); err != nil {
		return err
	}
	if _, err := packetData.Write(value.Data); err != nil {
		return err
	}

	var frame bytes.Buffer
	if threshold >= 0 && packetData.Len() >= threshold {
		var compressed bytes.Buffer
		zw := zlib.NewWriter(&compressed)
		if _, err := zw.Write(packetData.Bytes()); err != nil {
			return err
		}
		if err := zw.Close(); err != nil {
			return err
		}
		frameLength := pk.VarInt(pk.VarInt(packetData.Len()).Len() + compressed.Len())
		if _, err := frameLength.WriteTo(&frame); err != nil {
			return err
		}
		if _, err := pk.VarInt(packetData.Len()).WriteTo(&frame); err != nil {
			return err
		}
		if _, err := compressed.WriteTo(&frame); err != nil {
			return err
		}
	} else {
		dataLengthBytes := 0
		if threshold >= 0 {
			dataLengthBytes = 1
		}
		if _, err := pk.VarInt(dataLengthBytes + packetData.Len()).WriteTo(&frame); err != nil {
			return err
		}
		if threshold >= 0 {
			if _, err := pk.VarInt(0).WriteTo(&frame); err != nil {
				return err
			}
		}
		if _, err := packetData.WriteTo(&frame); err != nil {
			return err
		}
	}
	_, err := frame.WriteTo(w)
	return err
}

type singleByteReader struct {
	io.Reader
}

func (r singleByteReader) ReadByte() (byte, error) {
	var value [1]byte
	_, err := io.ReadFull(r.Reader, value[:])
	return value[0], err
}
