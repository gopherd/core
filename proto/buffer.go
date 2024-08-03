package proto

import (
	"encoding/json"
	"strconv"
	"sync"

	"github.com/gopherd/core/text/resp"
)

var bufferp = sync.Pool{
	New: func() any {
		return new(Buffer)
	},
}

// AllocBuffer gets buffer from pool
func AllocBuffer() *Buffer {
	buf := bufferp.Get().(*Buffer)
	buf.Reset()
	return buf
}

// FreeBuffer puts buffer to pool if cap of buffer less than 64k
func FreeBuffer(b *Buffer) {
	if b.Cap() < 1<<16 {
		bufferp.Put(b)
	}
}

type Buffer struct {
	buf []byte
	off int
}

func (b *Buffer) Cap() int {
	return cap(b.buf)
}

func (b *Buffer) Len() int {
	return len(b.buf) - b.off
}

func (b *Buffer) Bytes() []byte {
	return b.buf[b.off:]
}

func (b *Buffer) Reset() {
	b.off = 0
	b.buf = b.buf[:0]
}

func (b *Buffer) Reserve(n int) {
	if cap(b.buf) < n {
		buf := make([]byte, len(b.buf)-b.off, n)
		copy(buf, b.buf[b.off:])
		b.buf = buf
		b.off = 0
	}
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	b.buf = append(b.buf, p...)
	return len(p), nil
}

func (b *Buffer) Marshal(m Message) error {
	buf, err := m.MarshalAppend(b.buf, false)
	if err == nil {
		b.buf = buf
	}
	return err
}

func (b *Buffer) Unmarshal(m Message) error {
	return Unmarshal(b.Bytes(), m)
}

func (b *Buffer) Encode(m Message, contentType ContentType) error {
	b.Reset()
	var err error
	switch contentType {
	case ContentTypeText:
		b.buf = append(b.buf, resp.StringType.Byte())
		b.buf = strconv.AppendInt(b.buf, int64(m.Typeof()), 10)
		b.buf = append(b.buf, ' ')
		err = json.NewEncoder(b).Encode(m)
		if err == nil {
			if n := len(b.buf); n > 0 && b.buf[n-1] == '\n' {
				b.buf[n-1] = '\r'
				b.buf = append(b.buf, '\n')
			} else {
				b.buf = append(b.buf, '\r', '\n')
			}
		}
	case ContentTypeProtobuf:
		b.buf, err = EncodeAppend(b.buf, m)
	default:
		err = ErrUnsupportedContentType
	}
	return err
}
