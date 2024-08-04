package proto

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"

	"github.com/gopherd/core/event"
)

const (
	// max size of content: 1G
	MaxSize = 1 << 30
	// max message type
	MaxType = 1 << 31
)

var (
	ErrVarintOverflow         = errors.New("proto: varint overflow")
	ErrSizeOverflow           = errors.New("proto: size overflow")
	ErrTypeOverflow           = errors.New("proto: type overflow")
	ErrOutOfRange             = errors.New("proto: out of range")
	ErrUnsupportedContentType = errors.New("proto: unsupported content type")
)

type UnrecognizedTypeError struct {
	Type Type
}

func ErrUnrecognizedType(typ Type) error {
	return &UnrecognizedTypeError{Type: typ}
}

func (err *UnrecognizedTypeError) Error() string {
	return "proto: unrecognized message type " + strconv.FormatUint(uint64(err.Type), 10)
}

// ContentType represents encoding type of content
type ContentType int

const (
	ContentTypeProtobuf ContentType = iota
	ContentTypeText
)

// IsTextproto reports whether the contentType is a textproto type
func IsTextproto(contentType ContentType) bool {
	return contentType > ContentTypeProtobuf
}

// Type represents message type
type Type = uint32

// Body represents message body
type Body interface {
	io.Reader
	io.ByteReader

	// Len returns remain length of body
	Len() int

	// Peek returns the next n bytes without advancing the reader. The bytes stop
	// being valid at the next read call. If Peek returns fewer than n bytes, it
	// also returns an error explaining why the read is short.
	Peek(n int) ([]byte, error)

	// Discard skips the next n bytes, returning the number of bytes discarded.
	// If Discard skips fewer than n bytes, it also returns an error.
	Discard(n int) (discarded int, err error)
}

// textReader implements Body interface
type textReader struct {
	*bufio.Reader
	content *bytes.Reader
}

// Len implements Body Len method
func (tr *textReader) Len() int { return tr.content.Len() }

// Text creates a textproto body
func Text(b []byte) Body {
	content := bytes.NewReader(b)
	return &textReader{
		Reader:  bufio.NewReader(content),
		content: content,
	}
}

// Message represents a message interface
type Message interface {
	// Type of message
	Typeof() Type
	// Size of message
	Sizeof() int
	// Name of message
	Nameof() string

	// MarshalAppend marshals message to buf
	MarshalAppend(buf []byte, useCachedSize bool) ([]byte, error)
	// Unmarshal unmarshals message from buf
	Unmarshal(buf []byte) error
}

type MessageInfo struct {
	Type   Type
	Name   string
	Module string
}

var (
	creators  = make(map[Type]func() Message)
	mods      = make(map[string][]Type)
	type2mods = make(map[Type]string)
	messages  []MessageInfo
)

// Register registers a message creator by type. Register is not
// concurrent-safe, it is recommended to call in `init` function.
//
// e.g.
//
//	package foo
//
//	import "github.com/gopherd/core/proto"
//
//	func init() {
//		proto.Register("foo", BarType, func() proto.Message { return new(Bar) })
//	}
func Register(mod string, typ Type, creator func() Message) {
	if typ > MaxType {
		panic(fmt.Sprintf("proto: Register type %d out of range [0, %d]", typ, MaxType))
	}
	if creator == nil {
		panic(fmt.Sprintf("proto: Register creator is nil for type %d", typ))
	}
	if _, dup := creators[typ]; dup {
		panic(fmt.Sprintf("proto: Register called twice for type %d", typ))
	}
	creators[typ] = creator
	mods[mod] = append(mods[mod], typ)
	type2mods[typ] = mod
	messages = append(messages, MessageInfo{
		Type:   typ,
		Name:   creator().Nameof(),
		Module: mod,
	})
}

// Messages returns all registered message informations
func Messages() []MessageInfo {
	return messages
}

// New creates a message by type, nil returned if type not found
func New(typ Type) Message {
	if creator, ok := creators[typ]; ok {
		return creator()
	}
	return nil
}

// Arena represents a message factory
type Arena interface {
	Get(typ Type) Message
	Put(m Message)
}

// ArenaFunc wraps function as an Arena
type ArenaFunc func(Type) Message

// Get implements Arena Get method
func (fn ArenaFunc) Get(typ Type) Message { return fn(typ) }

// Put implements Arena Put method
func (fn ArenaFunc) Put(_ Message) {}

// Pool implements Arena interface to reuse message objects
type Pool struct {
	mu    sync.RWMutex
	pools map[Type]*sync.Pool
}

func (pp *Pool) getp(typ Type) *sync.Pool {
	if p := pp.findp(typ); p != nil {
		return p
	}
	f, ok := creators[typ]
	if !ok {
		return nil
	}
	pp.mu.Lock()
	defer pp.mu.Unlock()
	if pp.pools == nil {
		pp.pools = make(map[Type]*sync.Pool)
	}
	if p, ok := pp.pools[typ]; ok {
		return p
	}
	p := &sync.Pool{New: func() any { return f() }}
	pp.pools[typ] = p
	return p
}

func (pp *Pool) findp(typ Type) *sync.Pool {
	pp.mu.RLock()
	defer pp.mu.RUnlock()
	if pp.pools == nil {
		return nil
	}
	return pp.pools[typ]
}

// Get selects an message object from the Pool by type, removes it from the
// Pool, and returns it to the caller.
func (pp *Pool) Get(typ Type) Message {
	p := pp.getp(typ)
	if p == nil {
		return nil
	}
	if x := p.Get(); x == nil {
		println("proto: get a nil message from pool")
		return nil
	} else if m, ok := x.(Message); !ok {
		println("proto: get an unexpected message type from pool")
		return nil
	} else {
		return m
	}
}

// Put adds x to the pool.
func (pp *Pool) Put(m Message) {
	if m != nil {
		if p := pp.findp(m.Typeof()); p != nil {
			p.Put(m)
		}
	}
}

// Lookup lookups all registered types by module
func Lookup(module string) []Type {
	return mods[module]
}

// Moduleof returns the module name of message by type
func Moduleof(typ Type) string {
	return type2mods[typ]
}

// Peeker peeks n bytes
type Peeker interface {
	Peek(n int) ([]byte, error)
}

func peekUvarint(peeker Peeker) (int, uint64, error) {
	var x uint64
	var s uint
	var n int
	for i := 0; i < binary.MaxVarintLen64; i++ {
		n++
		buf, err := peeker.Peek(n)
		if err != nil {
			return n, x, err
		}
		b := buf[i]
		if b < 0x80 {
			if i == binary.MaxVarintLen64-1 && b > 1 {
				return n, x, ErrTypeOverflow
			}
			return n, x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return n, x, ErrTypeOverflow
}

func convertType(x uint64, err error) (Type, error) {
	if err != nil {
		return 0, err
	}
	if x > MaxType {
		return 0, ErrTypeOverflow
	}
	return Type(x), nil
}

// ParseType parses type from string
func ParseType(s string) (Type, error) {
	typ, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return 0, err
	}
	if typ < 0 || typ >= MaxType {
		return 0, ErrTypeOverflow
	}
	return Type(typ), nil
}

// ReadType reads message type from reader
func ReadType(r io.ByteReader) (typ Type, err error) {
	return convertType(binary.ReadUvarint(r))
}

// PeekType reads message type without advancing underlying reader offset
func PeekType(peeker Peeker) (n int, typ Type, err error) {
	var x uint64
	n, x, err = peekUvarint(peeker)
	typ, err = convertType(x, err)
	return
}

// EncodeType encodes type as varint to buf and returns number of bytes written.
func EncodeType(buf []byte, typ Type) int {
	return binary.PutUvarint(buf, uint64(typ))
}

func convertSize(x uint64, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	if x > MaxSize {
		return 0, ErrSizeOverflow
	}
	return int(x), nil
}

// ReadSize reads message size from reader
func ReadSize(r io.ByteReader) (size int, err error) {
	return convertSize(binary.ReadUvarint(r))
}

// PeekSize reads message size without advancing underlying reader offset
func PeekSize(peeker Peeker) (n int, size int, err error) {
	var x uint64
	n, x, err = peekUvarint(peeker)
	size, err = convertSize(x, err)
	return
}

// EncodeSize encodes type as varint to buf and returns number of bytes written.
func EncodeSize(buf []byte, size int) int {
	return binary.PutUvarint(buf, uint64(size))
}

// sizeofUvarint returns the number of unsigned-varint encoding-bytes.
func sizeofUvarint(x uint64) int {
	i := 0
	for x >= 0x80 {
		x >>= 7
		i++
	}
	return i + 1
}

// Encode returns the wire-format encoding of m with type and size.
//
//	|type|body.size|body|
func Encode(m Message, reservedHeadLen int) ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	size := m.Sizeof()
	if size > MaxSize {
		return nil, ErrSizeOverflow
	}
	off := reservedHeadLen
	tsize := sizeofUvarint(uint64(m.Typeof()))
	ssize := sizeofUvarint(uint64(size))
	buf := make([]byte, off+tsize+ssize+size)
	_, err := encodeAppend(buf[off:], m, size)
	return buf, err
}

// EncodeAppend encodes m to buf
func EncodeAppend(buf []byte, m Message) ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	size := m.Sizeof()
	if size > MaxSize {
		return nil, ErrSizeOverflow
	}
	off := len(buf)
	tsize := sizeofUvarint(uint64(m.Typeof()))
	ssize := sizeofUvarint(uint64(size))
	n := tsize + ssize + size
	if cap(buf)-off < n {
		newbuf := make([]byte, off+tsize+ssize, off+n)
		copy(newbuf, buf)
		buf = newbuf
	} else {
		buf = buf[:off+tsize+ssize]
	}
	return encodeAppend(buf[off:], m, size)
}

func encodeAppend(buf []byte, m Message, size int) ([]byte, error) {
	off := 0
	off += binary.PutUvarint(buf[off:], uint64(m.Typeof()))
	off += binary.PutUvarint(buf[off:], uint64(size))
	_, err := m.MarshalAppend(buf[off:off], true)
	return buf, err
}

// Marshal returns the wire-format encoding of m without type or size.
func Marshal(m Message) ([]byte, error) {
	return m.MarshalAppend(nil, false)
}

// Decode decodes one message with type and size from buf and
// returns number of bytes read and unmarshaled message.
func Decode(buf []byte, arena Arena) (int, Message, error) {
	off := 0
	// decode type
	typ, n := binary.Uvarint(buf[off:])
	if n == 0 {
		return off, nil, io.ErrShortBuffer
	} else if n < 0 {
		return off, nil, ErrTypeOverflow
	}
	off += n
	if typ > MaxType {
		return off, nil, ErrTypeOverflow
	}
	var m Message
	if arena == nil {
		m = New(Type(typ))
	} else {
		m = arena.Get(Type(typ))
	}
	if m == nil {
		return off, nil, &UnrecognizedTypeError{Type: Type(typ)}
	}
	// decode size
	size, n := binary.Uvarint(buf[off:])
	if n == 0 {
		return 0, nil, io.ErrShortBuffer
	} else if n < 0 {
		return -n, nil, ErrSizeOverflow
	}
	off += n
	if size > MaxSize {
		return n, nil, ErrSizeOverflow
	}
	end := off + int(size)
	if end > len(buf) {
		return n, nil, ErrOutOfRange
	}
	// decode body
	err := Unmarshal(buf[off:end], m)
	off += int(size)
	return off, m, err
}

// Unmarshal parses the wire-format message in b and places the result in m.
// The provided message must be mutable (e.g., a non-nil pointer to a message).
func Unmarshal(buf []byte, m Message) error {
	return m.Unmarshal(buf)
}

// Listener aliases event.Listener for Type
type Listener = event.Listener[Type]

// Dispatcher aliases event.Dispatcher for Type
type Dispatcher = event.Dispatcher[Type]

// Listen listens message handler
func Listen[H ~func(context.Context, M), M Message](h H) Listener {
	var m M
	return event.Listen[Type, M](m.Typeof(), h)
}
