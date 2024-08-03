package resp

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"
	"testing"
)

func newArray(a []*Value) *Value {
	return &Value{Type: ArrayType, elements: &slice{data: a, len: len(a)}}
}

func newString(s string) *Value {
	return &Value{Type: StringType, raw: bytes.NewBufferString(s), end: len(s)}
}

func newBytes(b []byte) *Value {
	return &Value{Type: BytesType, raw: bytes.NewBuffer(b), end: len(b)}
}

func newError(e string) *Value {
	return &Value{Type: ErrorType, raw: bytes.NewBufferString(e), end: len(e)}
}

func newInteger(i int64) *Value {
	v := &Value{Type: IntegerType, raw: bytes.NewBufferString(strconv.FormatInt(i, 10))}
	v.end = v.raw.Len()
	return v
}

func newNil() *Value {
	return &Value{Type: BytesType, isNil: true}
}

func equalsValue(v1, v2 *Value) bool {
	if (v1.IsNil() && !v2.IsNil()) || (!v1.IsNil() && v2.IsNil()) {
		return false
	}
	if v1.IsNil() && v2.IsNil() {
		return true
	}
	if v1.Type != v2.Type {
		return false
	}
	if v1.Type == ArrayType {
		if v1.getElements().len != v2.getElements().len {
			return false
		}
		for i := range v1.elements.values() {
			if !equalsValue(v1.elements.values()[i], v2.elements.values()[i]) {
				return false
			}
		}
		return true
	}
	if len(v1.Value()) != len(v2.Value()) {
		return false
	}
	for i := range v1.Value() {
		if v1.Value()[i] != v2.Value()[i] {
			return false
		}
	}
	return true
}

func TestRead(t *testing.T) {
	// *5\r\n           = > Array(3)
	// *4\r\n             = > Array(4)
	// +ok\r\n              = > String
	// -ERR timeout\r\n     = > Error
	// :123\r\n             = > Integer
	// $7\r\n               = > Blob
	// he\r\nllo\r\n
	// $1\r\n           = > Blob
	// \n\r\n
	// $2\r\n           = > Blob
	// \r\n\r\n
	// $-1\r\n          = > Blob
	// $0\r\n           = > Blob
	// \r\n
	const source = "*5\r\n*4\r\n+ok\r\n-ERR timeout\r\n:123\r\n$7\r\nhe\r\nllo\r\n$1\r\n\n\r\n$2\r\n\r\n\r\n$-1\r\n$0\r\n\r\n"
	var buf = bytes.NewBuffer([]byte(source))
	var expected = newArray([]*Value{
		newArray([]*Value{
			newString("ok"),
			newError("ERR timeout"),
			newInteger(123),
			newBytes([]byte("he\r\nllo")),
		}),
		newBytes([]byte{'\n'}),
		newBytes([]byte{'\r', '\n'}),
		newNil(),
		newBytes(nil),
	})

	var value = NewValue()
	if err := value.ReadFrom(bufio.NewReader(buf)); err != nil {
		t.Errorf("read value error: %v, value: %v, remain bytes: %q", err, value, buf.String())
		return
	} else {
		if !equalsValue(value, expected) {
			t.Errorf("value mismatch: %q vs %q", value, expected)
			return
		}
	}

	encoded := value.String()
	if encoded != source {
		t.Errorf("encodes mismatch with original: %q vs %q", encoded, source)
	}
}

func Test_Reset(t *testing.T) {
	var v = NewValue()
	type testCase struct {
		typ    Type
		value  string
		isNil  bool
		expect *Value
	}
	for i, c := range []testCase{
		testCase{StringType, "ok", false, newString("ok")},
		testCase{StringType, "ok", true, newNil()},
		testCase{BytesType, "hello", false, newBytes([]byte("hello"))},
		testCase{BytesType, "", true, newNil()},
		testCase{IntegerType, "1", false, newInteger(1)},
		testCase{IntegerType, "-1", false, newInteger(-1)},
		testCase{IntegerType, "100", false, newInteger(100)},
		testCase{ErrorType, "Error", false, newError("Error")},
	} {
		v.raw.Reset()
		v.reset(c.typ, []byte(c.value), c.isNil)
		if !equalsValue(c.expect, v) {
			t.Errorf("case %d error: %q vs %q", i, c.expect.String(), v.String())
		}
	}

	type testArrayCase struct {
		isNil  bool
		array  [][]byte
		expect *Value
	}
	for i, c := range []testArrayCase{
		testArrayCase{true, nil, newNil()},
		testArrayCase{false, [][]byte{}, newArray(nil)},
		testArrayCase{
			false,
			[][]byte{append(BytesType.bytes(), []byte("set")...), append(BytesType.bytes(), []byte("key")...)},
			newArray([]*Value{newBytes([]byte("set")), newBytes([]byte("key"))}),
		},
	} {
		if c.isNil {
			v.SetNil()
		} else {
			v.SetArray(len(c.array))
			for _, x := range c.array {
				v.Append(Type(x[0]), x[1:])
			}
		}
		if !equalsValue(c.expect, v) {
			t.Errorf("case %d error: %q vs %q", i, c.expect.String(), v.String())
		}
	}
}

func Benchmark_Read(b *testing.B) {
	const source = "*3\r\n$3\r\nSET\r\n$3\r\nget\r\n$10\r\nabcdefghij\r\n"
	var (
		reader    = strings.NewReader(source)
		bufReader = bufio.NewReader(reader)
		value     = NewValue()
	)
	for i := 0; i < b.N; i++ {
		if err := value.ReadFrom(bufReader); err != nil {
			b.Errorf("read error: %v", err)
		}
		reader.Seek(0, io.SeekStart)
		bufReader.Reset(reader)
	}
}
