// resp implements redis RESP
// @see https://redis.io/topics/protocol
//
package resp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

var (
	ErrInvalidType            = errors.New("resp: invalid type")
	ErrNumberBuffer           = errors.New("resp: invalid number buffer")
	ErrUnexpectedNumberPrefix = errors.New("resp: unexpected number prefix")
	ErrNumberOfArguments      = errors.New("resp: invalid number of arguments")
	ErrLengthOfArgument       = errors.New("resp: invalid length of argument")
)

const (
	maxNumberOfArguments = 65536
	maxLengthOfArgument  = 1024 * 1024 // 1M
)

const (
	cr    = '\r'
	lf    = '\n'
	space = ' '
)

var (
	crlf       = []byte{cr, lf}
	nilBytes   = []byte{'-', '1'}
	spaceBytes = []byte{space}
	okBytes    = []byte{'O', 'K'}
)

// Reader used to read resp content
type Reader interface {
	io.Reader
	io.ByteReader

	// ReadLine tries to return a single line, not including the end-of-line bytes.
	ReadLine() (line []byte, isPrefix bool, err error)
}

// Type of content
type Type byte

const (
	StringType  Type = '+'
	ErrorType   Type = '-'
	IntegerType Type = ':'
	BytesType   Type = '$'
	ArrayType   Type = '*'

	None Type = ' '
)

// Byte of type
func (t Type) Byte() byte {
	return byte(t)
}

func (t Type) String() string {
	switch t {
	case StringType:
		return "<string>"
	case ErrorType:
		return "<error>"
	case IntegerType:
		return "<integer>"
	case BytesType:
		return "<blob>"
	case ArrayType:
		return "<array>"
	default:
		return fmt.Sprintf("<%d>", int(t))
	}
}

var type2bytes = map[Type][]byte{
	StringType:  []byte{byte(StringType)},
	ErrorType:   []byte{byte(ErrorType)},
	IntegerType: []byte{byte(IntegerType)},
	BytesType:   []byte{byte(BytesType)},
	ArrayType:   []byte{byte(ArrayType)},
}

func (t Type) bytes() []byte {
	if b, ok := type2bytes[t]; ok {
		return b
	}
	return nil
}

// Command represents a command that using resp
type Command struct {
	Request *Value
}

// NewCommand creates a command
func NewCommand() *Command {
	return &Command{
		Request: NewValue(),
	}
}

// Is reports whether the command name is `name`
func (cmd *Command) Is(name string) bool {
	s := cmd.Request.Elements()[0].Value()
	if len(s) != len(name) {
		return false
	}
	for i := range s {
		x := s[i]
		y := name[i]
		if x == y {
			continue
		}
		if x >= 'A' && x <= 'Z' {
			x += 'a' - 'A'
		}
		if y >= 'A' && y <= 'Z' {
			y += 'a' - 'Z'
		}
		if x != y {
			return false
		}
	}
	return true
}

// Name returns name of command
func (cmd *Command) Name() string {
	return string(cmd.Request.Elements()[0].Value())
}

// NArg returns number of arguments
func (cmd *Command) NArg() int {
	return len(cmd.Request.Elements()) - 1
}

// Arg returns ith argument
func (cmd *Command) Arg(i int) string {
	return string(cmd.Request.Elements()[i+1].Value())
}

// Value of resp
type Value struct {
	Type       Type
	raw        *bytes.Buffer
	root       bool
	isNil      bool
	begin, end int
	elements   *slice
}

// NewValue creates a value
func NewValue() *Value {
	return &Value{
		root:     true,
		raw:      &bytes.Buffer{},
		elements: &slice{},
	}
}

// Reset resets the value
func (v *Value) Reset() *Value {
	v.raw.Reset()
	v.clear()
	return v
}

// SetNil sets the value as nil
func (v *Value) SetNil() *Value {
	v.raw.Reset()
	return v.reset(BytesType, nil, true)
}

// SetError sets the value as an error
func (v *Value) SetError(err error) *Value {
	v.raw.Reset()
	v.reset(ErrorType, []byte(err.Error()), false)
	return v
}

// SetArray sets the value as an array
func (v *Value) SetArray(len int) *Value {
	v.raw.Reset()
	return v.reset(ArrayType, itob(len), false)
}

// SetBytes sets the value as a Bytes
func (v *Value) SetBytes(value []byte) *Value {
	v.raw.Reset()
	return v.reset(BytesType, value, false)
}

// SetString sets the value as a string
func (v *Value) SetString(value string) *Value {
	v.raw.Reset()
	v.clear()
	v.isNil = false
	v.Type = StringType
	v.writeBegin(nil)
	v.raw.WriteString(value)
	v.writeEnd()
	return v
}

// SetInteger sets the value as an integer
func (v *Value) SetInteger(i int64) *Value {
	v.raw.Reset()
	return v.reset(IntegerType, i64tob(i), false)
}

// Set sets the value from interface
func (v *Value) Set(x any) error {
	v.raw.Reset()
	return reflect(v, x)
}

func (v *Value) reset(typ Type, value []byte, isNil bool) *Value {
	v.clear()
	v.Type = typ
	v.isNil = isNil
	if isNil {
		v.raw.WriteByte(byte(BytesType))
		v.raw.Write(nilBytes)
		v.writeEnd()
	} else {
		if typ == ArrayType {
			v.writeBegin(value)
		} else {
			if typ == BytesType {
				v.writeBegin(itob(len(value)))
			} else {
				v.writeBegin(nil)
			}
			v.writeValue(value)
			v.writeEnd()
		}
	}
	return v
}

func (v *Value) clear() {
	v.Type = None
	v.begin = 0
	v.end = 0
	v.isNil = false
	if v.elements != nil {
		v.elements.reset()
	}
}

// Append appends value to elements of array
func (v *Value) Append(typ Type, value []byte) (appended *Value) {
	elem := v.getElements().new(v.raw)
	elem.reset(typ, value, false)
	v.elements.append(elem)
	return elem
}

// AppendNil appends nil to elements of array
func (v *Value) AppendNil() {
	elem := v.getElements().new(v.raw)
	elem.reset(BytesType, nil, true)
	v.elements.append(elem)
}

func (v *Value) writeBegin(lenBuf []byte) {
	if v.isNil && v.Type != ArrayType {
		v.Type = BytesType
	}
	v.raw.WriteByte(byte(v.Type))
	if len(lenBuf) > 0 {
		v.raw.Write(lenBuf)
		v.raw.Write(crlf)
	}
}

func (v *Value) writeEnd() {
	v.raw.Write(crlf)
}

func (v *Value) writeValue(value []byte) {
	v.begin = v.raw.Len()
	v.raw.Write(value)
	v.end = v.raw.Len()
}

func (v *Value) writeStringValue(value string) {
	v.begin = v.raw.Len()
	v.raw.WriteString(value)
	v.end = v.raw.Len()
}

func reflectArray(v *Value, x []any) error {
	elements := v.getElements()
	v.Type = ArrayType
	v.writeBegin(itob(len(x)))
	for i := range x {
		elem := elements.new(v.raw)
		if err := reflect(elem, x[i]); err != nil {
			return err
		}
		elements.append(elem)
	}
	return nil
}

func reflect(to *Value, x any) error {
	to.clear()
	if x == nil {
		to.reset(BytesType, nil, true)
		return nil
	}
	switch v := x.(type) {
	case []any:
		return reflectArray(to, v)
	case string:
		to.Type = BytesType
		to.writeBegin(itob(len(v)))
		to.begin = to.raw.Len()
		to.raw.WriteString(v)
		to.end = to.raw.Len()
		to.writeEnd()
	case []byte:
		to.Type = BytesType
		to.writeBegin(itob(len(v)))
		to.writeValue(v)
		to.writeEnd()
	case int64:
		to.Type = IntegerType
		to.writeBegin(nil)
		to.writeValue(i64tob(v))
		to.writeEnd()
	case int:
		to.Type = IntegerType
		to.writeBegin(nil)
		to.writeValue(itob(v))
		to.writeEnd()
	case int32:
		to.Type = IntegerType
		to.writeBegin(nil)
		to.writeValue(i64tob(int64(v)))
		to.writeEnd()
	case uint64:
		to.Type = IntegerType
		to.writeBegin(nil)
		to.writeValue(i64tob(int64(v)))
		to.writeEnd()
	case uint:
		to.Type = IntegerType
		to.writeBegin(nil)
		to.writeValue(i64tob(int64(v)))
		to.writeEnd()
	case uint32:
		to.Type = IntegerType
		to.writeBegin(nil)
		to.writeValue(i64tob(int64(v)))
		to.writeEnd()
	default:
		return ErrInvalidType
	}
	return nil
}

// Bool returns parsed boolean value
func (v *Value) Bool() (bool, error) {
	i, err := v.Int64()
	return i > 0, err
}

// Int64 returns parsed integer value
func (v *Value) Int64() (int64, error) {
	return btoi64(v.Value(), 64)
}

// Float returns parsed float value
func (v *Value) Float() (float64, error) {
	return strconv.ParseFloat(string(v.Value()), 64)
}

// OK reports whether the value is OK
func (v *Value) OK() bool {
	return (v.Type == StringType || v.Type == BytesType) && bytes.Equal(v.Value(), okBytes)
}

func (v *Value) getElements() *slice {
	if v.elements == nil {
		v.elements = &slice{}
	}
	return v.elements
}

// Value returns content of value
func (v *Value) Value() []byte {
	if v.raw == nil {
		return nil
	}
	return v.raw.Bytes()[v.begin:v.end]
}

// Elements returns elements of array
func (v *Value) Elements() []*Value {
	return v.getElements().values()
}

// IsNil reports whether ther value is nil
func (v *Value) IsNil() bool { return v.isNil }

// String returns marshaled string
func (v *Value) String() string {
	if v.root {
		return v.raw.String()
	}
	var w bytes.Buffer
	v.writeTo(&w)
	return w.String()
}

// WriteTo writes marshaled content to w
func (v *Value) WriteTo(w io.Writer) (int, error) {
	if v.root {
		return writeBytes(None, w, v.raw.Bytes(), false)
	}
	return v.writeTo(w)
}

func (v *Value) writeTo(w io.Writer) (int, error) {
	if v.isNil {
		if v.Type != ArrayType {
			v.Type = BytesType
		}
		return writeBytes(v.Type, w, nilBytes, true)
	}
	switch v.Type {
	case ArrayType:
		wrote := 0
		elements := v.getElements()
		if n, err := writeBytes(v.Type, w, itob(elements.len), true); err != nil {
			wrote += n
			return wrote, err
		} else {
			wrote += n
			for _, value := range elements.values() {
				if n, err := value.writeTo(w); err != nil {
					wrote += n
					return wrote, err
				} else {
					wrote += n
				}
			}
			return wrote, err
		}
	case BytesType:
		wrote := 0
		if n, err := writeBytes(v.Type, w, itob(len(v.Value())), true); err != nil {
			wrote += n
			return wrote, err
		} else {
			wrote += n
			if n, err := writeBytes(None, w, v.Value(), true); err != nil {
				wrote += n
				return wrote, err
			} else {
				wrote += n
			}
		}
		return wrote, nil
	default:
		return writeBytes(v.Type, w, v.Value(), true)
	}
}

// ReadFrom reads value from reader
func (v *Value) ReadFrom(r Reader) error {
	v.raw.Reset()
	v.clear()
	return v.readFrom(r, v.raw)
}

func (v *Value) readFrom(r Reader, raw *bytes.Buffer) error {
	for {
		t, err := r.ReadByte()
		if err != nil {
			return err
		}
		if t == cr || t == lf || t == ' ' {
			continue
		}
		v.Type = Type(t)
		raw.WriteByte(t)
		break
	}
	switch v.Type {
	case StringType, ErrorType, IntegerType:
		v.begin = raw.Len()
		if l, err := readLine(r, raw); err != nil {
			return err
		} else {
			v.end = v.begin + len(l)
		}
	case BytesType:
		n, err := readNumber(r, raw)
		if err != nil {
			return err
		}
		if n > maxLengthOfArgument {
			return ErrLengthOfArgument
		}
		v.begin = raw.Len()
		if n >= 0 {
			if data, err := readBytes(r, raw, n); err != nil {
				return err
			} else {
				v.end = v.begin + len(data)
			}
		} else if n < 0 {
			v.end = v.begin
			// RESP 的 nil 值是用 n == -1 表示的
			v.isNil = true
		}
	case ArrayType:
		n, err := readNumber(r, raw)
		if err != nil {
			return err
		}
		if n < 0 || n > maxNumberOfArguments {
			return ErrNumberOfArguments
		} else if n > 0 {
			if v.elements == nil {
				v.elements = &slice{
					data: make([]*Value, 0, n),
				}
			}
			for i := int64(0); i < n; i++ {
				elem := v.elements.new(raw)
				if err := elem.readFrom(r, raw); err != nil {
					return err
				}
				v.elements.append(elem)
			}
		}
	default:
		// inline command 单行且以空格分隔参数
		begin := raw.Len() - 1
		if l, err := readLine(r, raw); err != nil {
			return err
		} else {
			v.Type = ArrayType
			elements := v.getElements()
			for i := begin; i < len(l); i++ {
				if l[i] == space {
					if i > begin {
						elem := elements.new(raw)
						elem.Type = BytesType
						elem.begin = begin
						elem.end = i
						elements.append(elem)
					} else {
						begin = i + 1
					}
				}
			}
			if begin < len(l) {
				elem := elements.new(raw)
				elem.Type = BytesType
				elem.begin = begin
				elem.end = len(l)
				elements.append(elem)
			}
		}
	}
	return nil
}

//---------------------------------------------------------
// helper functions

type slice struct {
	data []*Value
	len  int // <= len(data)
}

func (s *slice) reset() {
	const maxSliceCacheLen = 128
	s.len = 0
	if n := len(s.data); n > maxSliceCacheLen {
		for i := maxSliceCacheLen; i < n; i++ {
			s.data[i] = nil
		}
		s.data = s.data[:maxSliceCacheLen]
	}
}

func (s slice) values() []*Value {
	return s.data[:s.len]
}

func (s *slice) append(v *Value) {
	if s.len == len(s.data) {
		s.data = append(s.data, v)
	} else {
		s.data[s.len] = v
	}
	s.len++
}

func (s *slice) new(raw *bytes.Buffer) *Value {
	if s.len == len(s.data) {
		v := &Value{raw: raw}
		s.data = append(s.data, v)
		return v
	}
	s.data[s.len].clear()
	s.data[s.len].raw = raw
	return s.data[s.len]
}

func writeBytes(t Type, w io.Writer, v []byte, end bool) (int, error) {
	wrote := 0
	if t != None {
		if b := t.bytes(); b != nil {
			if n, err := writeFull(w, b); err != nil {
				wrote += n
				return wrote, err
			} else {
				wrote += n
			}
		}
	}
	if n, err := writeFull(w, v); err != nil {
		wrote += n
		return wrote, err
	} else {
		wrote += n
	}
	if end {
		if n, err := writeFull(w, crlf); err != nil {
			wrote += n
			return wrote, err
		} else {
			wrote += n
		}
	}
	return wrote, nil
}

func writeFull(w io.Writer, buf []byte) (int, error) {
	wrote := 0
	offset := 0
	for {
		n, err := w.Write(buf[offset:])
		if err != nil {
			wrote += n
			return wrote, err
		}
		wrote += n
		offset += n
		if offset == len(buf) {
			break
		}
	}
	return wrote, nil
}

func readLine(r Reader, to *bytes.Buffer) ([]byte, error) {
	begin := to.Len()
	for {
		l, more, err := r.ReadLine()
		if err != nil {
			return nil, err
		}
		to.Write(l)
		if !more {
			break
		}
	}
	end := to.Len()
	to.Write(crlf)
	return to.Bytes()[begin:end], nil
}

func readNumber(r Reader, to *bytes.Buffer) (int64, error) {
	l, err := readLine(r, to)
	if err != nil {
		return 0, err
	}
	return btoi64(l, 32)
}

func readBytes(r Reader, to *bytes.Buffer, n int64) ([]byte, error) {
	extra := len(crlf)
	to.Grow(int(n) + extra)
	data := to.Bytes()
	begin := len(data)
	if _, err := io.CopyN(to, r, n+int64(extra)); err != nil {
		return nil, err
	}
	return data[begin : begin+int(n)], nil
}
