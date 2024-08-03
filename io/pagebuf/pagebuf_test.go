package pagebuf_test

import (
	"io"
	"testing"

	. "github.com/gopherd/core/io/pagebuf"
)

func TestPageBuffer(t *testing.T) {
	const testdata = "hello,pagebuf"

	for _, read := range []func(io.Reader, []byte) (int, error){
		io.ReadFull,
		func(r io.Reader, p []byte) (int, error) {
			return r.Read(p)
		},
	} {
		p := NewPageBufferSize(2, 8)
		p.Write([]byte(testdata))
		var buf = make([]byte, len(testdata))
		if n, err := read(p, buf); err != nil {
			t.Fatalf("read error: %v", err)
		} else {
			if got := string(buf[:n]); got != testdata {
				t.Fatalf("result mismatched: %q vs %q, size=%d", got, testdata, p.Len())
			}
			t.Log(buf)
		}
	}
}
