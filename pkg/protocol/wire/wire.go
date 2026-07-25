package wire

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"time"
	"unicode/utf8"

	"github.com/KonjacBot/go-mc/chat"
	"github.com/KonjacBot/go-mc/chat/sign"
	"github.com/KonjacBot/go-mc/nbt"
	pk "github.com/KonjacBot/go-mc/net/packet"
	"github.com/KonjacBot/go-mc/yggdrasil/user"
	"github.com/google/uuid"
)

type (
	Field         = pk.Field
	FieldEncoder  = pk.FieldEncoder
	FieldDecoder  = pk.FieldDecoder
	Boolean       = pk.Boolean
	Byte          = pk.Byte
	UnsignedByte  = pk.UnsignedByte
	Short         = pk.Short
	UnsignedShort = pk.UnsignedShort
	Int           = pk.Int
	Long          = pk.Long
	Float         = pk.Float
	Double        = pk.Double
	VarInt        = pk.VarInt
	VarLong       = pk.VarLong
	UUID          = pk.UUID
	FixedBitSet   = pk.FixedBitSet

	String     string
	Identifier = String
	ByteArray  []byte
	BitSet     []int64
)

const (
	MaxStringChars       = 32767
	MaxStringBytes       = MaxStringChars * 3
	MaxCollectionEntries = 32767
	maxNBTBytes          = pk.MaxDataLength
	maxNBTDepth          = 512
	maxNBTEntries        = 65536
)

func (s String) WriteTo(w io.Writer) (n int64, err error) {
	value := string(s)
	if !utf8.ValidString(value) {
		return 0, errors.New("string contains invalid UTF-8")
	}
	if len(value) > MaxStringBytes || utf16CodeUnits(value) > MaxStringChars {
		return 0, fmt.Errorf("string exceeds %d characters", MaxStringChars)
	}
	n, err = pk.VarInt(len(value)).WriteTo(w)
	if err != nil {
		return n, err
	}
	written, err := io.WriteString(w, value)
	return n + int64(written), err
}

func (s *String) ReadFrom(r io.Reader) (n int64, err error) {
	var length pk.VarInt
	n, err = length.ReadFrom(r)
	if err != nil {
		return n, err
	}
	if length < 0 {
		return n, errors.New("string length less than zero")
	}
	if length > MaxStringBytes {
		return n, fmt.Errorf("string length %d exceeds %d bytes", length, MaxStringBytes)
	}
	if remaining, ok := r.(interface{ Len() int }); ok && int(length) > remaining.Len() {
		return n, io.ErrUnexpectedEOF
	}

	data := make([]byte, int(length))
	read, err := io.ReadFull(r, data)
	n += int64(read)
	if err != nil {
		return n, err
	}
	if !utf8.Valid(data) {
		return n, errors.New("string contains invalid UTF-8")
	}
	value := string(data)
	if utf16CodeUnits(value) > MaxStringChars {
		return n, fmt.Errorf("string exceeds %d characters", MaxStringChars)
	}
	*s = String(value)
	return n, nil
}

func utf16CodeUnits(value string) int {
	count := 0
	for _, r := range value {
		if r >= 0x10000 {
			count += 2
		} else {
			count++
		}
	}
	return count
}

func (b ByteArray) WriteTo(w io.Writer) (n int64, err error) {
	if len(b) > pk.MaxDataLength {
		return 0, fmt.Errorf("byte array exceeds %d bytes", pk.MaxDataLength)
	}
	n, err = pk.VarInt(len(b)).WriteTo(w)
	if err != nil {
		return n, err
	}
	written, err := w.Write(b)
	return n + int64(written), err
}

func (b *ByteArray) ReadFrom(r io.Reader) (n int64, err error) {
	var length pk.VarInt
	n, err = length.ReadFrom(r)
	if err != nil {
		return n, err
	}
	if length < 0 {
		return n, errors.New("byte array length less than zero")
	}
	if length > pk.MaxDataLength {
		return n, fmt.Errorf("byte array length %d exceeds %d", length, pk.MaxDataLength)
	}
	if remaining, ok := r.(interface{ Len() int }); ok && int(length) > remaining.Len() {
		return n, io.ErrUnexpectedEOF
	}
	data := make([]byte, int(length))
	read, err := io.ReadFull(r, data)
	n += int64(read)
	if err != nil {
		return n, err
	}
	*b = data
	return n, nil
}

func (b BitSet) WriteTo(w io.Writer) (n int64, err error) {
	if len(b) > MaxCollectionEntries {
		return 0, fmt.Errorf("bit set exceeds %d longs", MaxCollectionEntries)
	}
	n, err = pk.VarInt(len(b)).WriteTo(w)
	if err != nil {
		return n, err
	}
	for _, value := range b {
		written, writeErr := pk.Long(value).WriteTo(w)
		n += written
		if writeErr != nil {
			return n, writeErr
		}
	}
	return n, nil
}

func (b *BitSet) ReadFrom(r io.Reader) (n int64, err error) {
	var length pk.VarInt
	n, err = length.ReadFrom(r)
	if err != nil {
		return n, err
	}
	if length < 0 {
		return n, errors.New("bit set length less than zero")
	}
	if length > MaxCollectionEntries {
		return n, fmt.Errorf("bit set exceeds %d longs", MaxCollectionEntries)
	}
	if remaining, ok := r.(interface{ Len() int }); ok && int64(length)*8 > int64(remaining.Len()) {
		return n, io.ErrUnexpectedEOF
	}
	values := make(BitSet, int(length))
	for i := range values {
		written, readErr := (*pk.Long)(&values[i]).ReadFrom(r)
		n += written
		if readErr != nil {
			return n, readErr
		}
	}
	*b = values
	return n, nil
}

func NewFixedBitSet(bits int64) FixedBitSet {
	return pk.NewFixedBitSet(bits)
}

type arrayField struct {
	value any
}

func Array(value any) pk.Field {
	return arrayField{value: value}
}

func (a arrayField) WriteTo(w io.Writer) (n int64, err error) {
	value, err := sliceValue(a.value, false)
	if err != nil {
		return 0, err
	}
	if value.Len() > MaxCollectionEntries {
		return 0, fmt.Errorf("array length %d exceeds %d", value.Len(), MaxCollectionEntries)
	}
	n, err = pk.VarInt(value.Len()).WriteTo(w)
	if err != nil {
		return n, err
	}
	for i := 0; i < value.Len(); i++ {
		field, ok := fieldEncoder(value.Index(i))
		if !ok {
			return n, fmt.Errorf("array element %d of type %s is not encodable", i, value.Type().Elem())
		}
		written, writeErr := field.WriteTo(w)
		n += written
		if writeErr != nil {
			return n, writeErr
		}
	}
	return n, nil
}

func (a arrayField) ReadFrom(r io.Reader) (n int64, err error) {
	value, err := sliceValue(a.value, true)
	if err != nil {
		return 0, err
	}
	var length pk.VarInt
	n, err = length.ReadFrom(r)
	if err != nil {
		return n, err
	}
	if length < 0 {
		return n, errors.New("array length less than zero")
	}
	if length > MaxCollectionEntries {
		return n, fmt.Errorf("array length %d exceeds %d", length, MaxCollectionEntries)
	}
	if length == 0 {
		value.Set(reflect.Zero(value.Type()))
		return n, nil
	}

	capacity := int(length)
	if capacity > 64 {
		capacity = 64
	}
	decoded := reflect.MakeSlice(value.Type(), 0, capacity)
	elementType := value.Type().Elem()
	for i := 0; i < int(length); i++ {
		var element reflect.Value
		var appendValue reflect.Value
		if elementType.Kind() == reflect.Pointer {
			element = reflect.New(elementType.Elem())
			appendValue = element
		} else {
			element = reflect.New(elementType)
			appendValue = element.Elem()
		}
		field, ok := element.Interface().(pk.FieldDecoder)
		if !ok {
			return n, fmt.Errorf("array element %d of type %s is not decodable", i, elementType)
		}
		read, readErr := field.ReadFrom(r)
		n += read
		if readErr != nil {
			value.Set(decoded)
			return n, readErr
		}
		decoded = reflect.Append(decoded, appendValue)
	}
	value.Set(decoded)
	return n, nil
}

func sliceValue(value any, writable bool) (reflect.Value, error) {
	if value == nil {
		return reflect.Value{}, errors.New("array value is nil")
	}
	result := reflect.ValueOf(value)
	if writable {
		if result.Kind() != reflect.Pointer || result.IsNil() {
			return reflect.Value{}, errors.New("array decode target must be a non-nil pointer to a slice")
		}
		result = result.Elem()
	} else {
		for result.Kind() == reflect.Pointer {
			if result.IsNil() {
				return reflect.Value{}, errors.New("array value is nil")
			}
			result = result.Elem()
		}
	}
	if result.Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("array value has type %s, want slice", result.Type())
	}
	return result, nil
}

func fieldEncoder(value reflect.Value) (pk.FieldEncoder, bool) {
	if value.CanInterface() {
		if field, ok := value.Interface().(pk.FieldEncoder); ok {
			return field, true
		}
	}
	if value.CanAddr() && value.Addr().CanInterface() {
		field, ok := value.Addr().Interface().(pk.FieldEncoder)
		return field, ok
	}
	return nil, false
}

type IDSet struct {
	TagName Identifier
	IDs     []int32
}

func (i IDSet) WriteTo(w io.Writer) (n int64, err error) {
	if i.TagName != "" {
		n, err = pk.VarInt(0).WriteTo(w)
		if err != nil {
			return n, err
		}
		written, err := i.TagName.WriteTo(w)
		return n + written, err
	}
	if len(i.IDs) > MaxCollectionEntries {
		return 0, fmt.Errorf("holder set exceeds %d ids", MaxCollectionEntries)
	}
	n, err = pk.VarInt(len(i.IDs) + 1).WriteTo(w)
	if err != nil {
		return n, err
	}
	for _, id := range i.IDs {
		written, writeErr := pk.VarInt(id).WriteTo(w)
		n += written
		if writeErr != nil {
			return n, writeErr
		}
	}
	return n, nil
}

func (i *IDSet) ReadFrom(r io.Reader) (n int64, err error) {
	*i = IDSet{}
	var encodedLength pk.VarInt
	n, err = encodedLength.ReadFrom(r)
	if err != nil {
		return n, err
	}
	if encodedLength < 0 {
		return n, errors.New("holder set length less than zero")
	}
	if encodedLength == 0 {
		read, err := i.TagName.ReadFrom(r)
		return n + read, err
	}
	length := encodedLength - 1
	if length > MaxCollectionEntries {
		return n, fmt.Errorf("holder set exceeds %d ids", MaxCollectionEntries)
	}
	i.IDs = make([]int32, int(length))
	for index := range i.IDs {
		var id pk.VarInt
		read, readErr := id.ReadFrom(r)
		n += read
		if readErr != nil {
			return n, readErr
		}
		i.IDs[index] = int32(id)
	}
	return n, nil
}

type NBTField struct {
	V                     any
	DisallowUnknownFields bool
}

// Message applies the bounded NBT decoder to Minecraft text components while
// preserving chat.Message's public methods through embedding.
type Message struct {
	chat.Message
}

func (m Message) WriteTo(w io.Writer) (int64, error) {
	return NBT(&m.Message).WriteTo(w)
}

func (m *Message) ReadFrom(r io.Reader) (int64, error) {
	m.Message = chat.Message{}
	return NBT(&m.Message).ReadFrom(r)
}

func NewMessage(message chat.Message) Message {
	return Message{Message: message}
}

type JsonMessage struct {
	chat.Message
}

type Property struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Signature string `json:"signature,omitempty"`
}

type Session struct {
	sign.Session
}

func (s Session) WriteTo(w io.Writer) (int64, error) {
	return s.Session.WriteTo(w)
}

func (s *Session) ReadFrom(r io.Reader) (n int64, err error) {
	var sessionID uuid.UUID
	read, err := (*pk.UUID)(&sessionID).ReadFrom(r)
	n += read
	if err != nil {
		return n, err
	}
	var expiresAt pk.Long
	read, err = expiresAt.ReadFrom(r)
	n += read
	if err != nil {
		return n, err
	}
	var publicKey, signature ByteArray
	read, err = publicKey.ReadFrom(r)
	n += read
	if err != nil {
		return n, err
	}
	read, err = signature.ReadFrom(r)
	n += read
	if err != nil {
		return n, err
	}
	decoded, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return n, err
	}
	rsaKey, ok := decoded.(*rsa.PublicKey)
	if !ok {
		return n, fmt.Errorf("profile public key is %T, want RSA", decoded)
	}
	s.Session = sign.Session{
		SessionID: sessionID,
		PublicKey: user.PublicKey{
			ExpiresAt: time.UnixMilli(int64(expiresAt)),
			PubKey:    rsaKey,
			Signature: bytes.Clone(signature),
		},
	}
	return n, nil
}

func (p Property) WriteTo(w io.Writer) (n int64, err error) {
	for _, value := range []string{p.Name, p.Value} {
		written, writeErr := String(value).WriteTo(w)
		n += written
		if writeErr != nil {
			return n, writeErr
		}
	}
	written, err := pk.Boolean(p.Signature != "").WriteTo(w)
	n += written
	if err != nil || p.Signature == "" {
		return n, err
	}
	written, err = String(p.Signature).WriteTo(w)
	return n + written, err
}

func (p *Property) ReadFrom(r io.Reader) (n int64, err error) {
	*p = Property{}
	for _, target := range []*string{&p.Name, &p.Value} {
		var value String
		read, readErr := value.ReadFrom(r)
		n += read
		if readErr != nil {
			return n, readErr
		}
		*target = string(value)
	}
	var hasSignature pk.Boolean
	read, err := hasSignature.ReadFrom(r)
	n += read
	if err != nil || !hasSignature {
		return n, err
	}
	var signature String
	read, err = signature.ReadFrom(r)
	n += read
	p.Signature = string(signature)
	return n, err
}

func (m JsonMessage) WriteTo(w io.Writer) (n int64, err error) {
	data, err := json.Marshal(m.Message)
	if err != nil {
		return 0, err
	}
	return String(data).WriteTo(w)
}

func (m *JsonMessage) ReadFrom(r io.Reader) (n int64, err error) {
	var data String
	n, err = data.ReadFrom(r)
	if err != nil {
		return n, err
	}
	m.Message = chat.Message{}
	return n, json.Unmarshal([]byte(data), &m.Message)
}

func NBT(value any) pk.Field {
	return NBTField{V: value}
}

func (n NBTField) WriteTo(w io.Writer) (written int64, err error) {
	var payload bytes.Buffer
	limited := &limitedWriter{writer: &payload, remaining: maxNBTBytes}
	_, err = (pk.NBTField{V: n.V, DisallowUnknownFields: n.DisallowUnknownFields}).WriteTo(limited)
	if err != nil {
		return 0, err
	}
	count, err := payload.WriteTo(w)
	return count, err
}

func (n NBTField) ReadFrom(r io.Reader) (read int64, err error) {
	raw, validationErr := validatedNBT(r)
	read = int64(len(raw))
	if validationErr != nil {
		return read, validationErr
	}
	decoder := nbt.NewDecoder(bytes.NewReader(raw))
	decoder.NetworkFormat(true)
	if n.DisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	_, err = decoder.Decode(n.V)
	if errors.Is(err, nbt.ErrEND) {
		err = nil
	}
	return read, err
}

type limitedWriter struct {
	writer    io.Writer
	remaining int
}

func (w *limitedWriter) Write(data []byte) (int, error) {
	if len(data) > w.remaining {
		return 0, fmt.Errorf("NBT payload exceeds %d bytes", maxNBTBytes)
	}
	written, err := w.writer.Write(data)
	w.remaining -= written
	return written, err
}

func validatedNBT(r io.Reader) ([]byte, error) {
	var recorded bytes.Buffer
	limited := &io.LimitedReader{R: r, N: maxNBTBytes + 1}
	validator := nbtValidator{r: io.TeeReader(limited, &recorded)}
	tag, err := validator.byte()
	if err == nil {
		err = validator.payload(tag, 0)
	}
	if recorded.Len() > maxNBTBytes {
		err = fmt.Errorf("NBT payload exceeds %d bytes", maxNBTBytes)
	}
	return bytes.Clone(recorded.Bytes()), err
}

type nbtValidator struct {
	r io.Reader
}

func (v *nbtValidator) payload(tag byte, depth int) error {
	if depth > maxNBTDepth {
		return fmt.Errorf("NBT nesting exceeds %d", maxNBTDepth)
	}
	switch tag {
	case nbt.TagEnd:
		return nil
	case nbt.TagByte:
		return v.skip(1)
	case nbt.TagShort:
		return v.skip(2)
	case nbt.TagInt, nbt.TagFloat:
		return v.skip(4)
	case nbt.TagLong, nbt.TagDouble:
		return v.skip(8)
	case nbt.TagByteArray:
		length, err := v.length()
		if err != nil {
			return err
		}
		return v.skip(int64(length))
	case nbt.TagString:
		return v.string()
	case nbt.TagList:
		elementTag, err := v.byte()
		if err != nil {
			return err
		}
		length, err := v.length()
		if err != nil {
			return err
		}
		if length > maxNBTEntries {
			return fmt.Errorf("NBT list length %d exceeds %d", length, maxNBTEntries)
		}
		if elementTag == nbt.TagEnd && length != 0 {
			return errors.New("NBT list with TAG_End has non-zero length")
		}
		for i := int32(0); i < length; i++ {
			if err := v.payload(elementTag, depth+1); err != nil {
				return err
			}
		}
		return nil
	case nbt.TagCompound:
		for entries := 0; ; entries++ {
			if entries > maxNBTEntries {
				return fmt.Errorf("NBT compound exceeds %d entries", maxNBTEntries)
			}
			childTag, err := v.byte()
			if err != nil {
				return err
			}
			if childTag == nbt.TagEnd {
				return nil
			}
			if err := v.string(); err != nil {
				return err
			}
			if err := v.payload(childTag, depth+1); err != nil {
				return err
			}
		}
	case nbt.TagIntArray:
		length, err := v.length()
		if err != nil {
			return err
		}
		return v.skip(int64(length) * 4)
	case nbt.TagLongArray:
		length, err := v.length()
		if err != nil {
			return err
		}
		return v.skip(int64(length) * 8)
	default:
		return fmt.Errorf("unknown NBT tag %#x", tag)
	}
}

func (v *nbtValidator) byte() (byte, error) {
	var value [1]byte
	_, err := io.ReadFull(v.r, value[:])
	return value[0], err
}

func (v *nbtValidator) length() (int32, error) {
	var data [4]byte
	if _, err := io.ReadFull(v.r, data[:]); err != nil {
		return 0, err
	}
	length := int32(binary.BigEndian.Uint32(data[:]))
	if length < 0 {
		return 0, errors.New("NBT collection length less than zero")
	}
	return length, nil
}

func (v *nbtValidator) string() error {
	var data [2]byte
	if _, err := io.ReadFull(v.r, data[:]); err != nil {
		return err
	}
	return v.skip(int64(binary.BigEndian.Uint16(data[:])))
}

func (v *nbtValidator) skip(length int64) error {
	if length < 0 || length > maxNBTBytes {
		return fmt.Errorf("invalid NBT payload length %d", length)
	}
	_, err := io.CopyN(io.Discard, v.r, length)
	return err
}
