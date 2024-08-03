package pagebuf

import (
	"io"
	"math/bits"
)

const (
	DefaultMaxIdlePages = 4
	DefaultPageSize     = 1024
)

// PageBuffer ...
type PageBuffer struct {
	bits         int
	mask         int
	maxIdlePages int

	pages [][]byte
	size  int
	off   struct {
		page int
		from int
	}
}

// NewPageBuffer creates a PageBuffer with default options
func NewPageBuffer() *PageBuffer {
	return NewPageBufferSize(DefaultMaxIdlePages, DefaultPageSize)
}

// NewPageBufferSize creates a PageBuffer with specified options
func NewPageBufferSize(maxIdlePages, pageSize int) *PageBuffer {
	mask := pageSize - 1
	if mask&pageSize != 0 {
		panic("pipe: pageSize is not a power of two")
	}
	if maxIdlePages <= 0 {
		maxIdlePages = 2
	}
	return &PageBuffer{
		bits:         bits.OnesCount(uint(mask)),
		mask:         mask,
		maxIdlePages: maxIdlePages,
	}
}

// Reset clears the buffer
func (p *PageBuffer) Reset() {
	p.off.page = 0
	p.off.from = 0
	p.size = 0
	if len(p.pages) > p.maxIdlePages {
		pages := make([][]byte, p.maxIdlePages)
		copy(pages, p.pages[:p.maxIdlePages])
		p.pages = pages
	}
}

// Len returns the length of buffer
func (p *PageBuffer) Len() int {
	return p.size
}

// PageSize returns size per page
func (p *PageBuffer) PageSize() int {
	return p.mask + 1
}

// Format formats the buffer as a string
func (p *PageBuffer) Format() string {
	const hex = "0123456789abcdef"
	if p.size == 0 {
		return "[]"
	}
	var (
		buf  = make([]byte, 0, p.size*3+1)
		page = p.off.page
		from = p.off.from
	)
	buf = append(buf, '[')
	for i := 0; i < p.size; i++ {
		b := p.pages[page][from]
		if i > 0 {
			buf = append(buf, ' ')
		}
		buf = append(buf, hex[b>>4], hex[b&0xF])
		if from == p.mask {
			from = 0
			page++
		} else {
			from++
		}
	}
	buf = append(buf, ']')
	return string(buf)
}

func (p *PageBuffer) tail() (page, from, size int) {
	end := p.off.from + p.size
	page = p.off.page + end>>p.bits
	from = end & p.mask
	size = (len(p.pages)-page)<<p.bits - from
	return
}

func (p *PageBuffer) grow(n int) (endPage, endFrom int) {
	var remainSize int
	endPage, endFrom, remainSize = p.tail()
	if remainSize < n {
		remainSize += p.off.page << p.bits
		if endFrom > 0 {
			endPage++
		}
		if remainSize < n {
			// allocate new pages
			need := (n - remainSize + p.mask) >> p.bits
			if need < (len(p.pages) >> 1) {
				need = len(p.pages) >> 1
			}
			pages := make([][]byte, need+len(p.pages))
			for i := range pages {
				pages[i] = make([]byte, p.mask+1)
			}
			copy(pages, p.pages[p.off.page:endPage])
			p.pages = pages
		} else {
			// shift content pages to front
			maxIdleSize := p.maxIdlePages << p.bits
			numPage := len(p.pages)
			if numPage > p.maxIdlePages &&
				maxIdleSize > (p.size+n)<<1 &&
				maxIdleSize > p.size+p.off.from+n {
				pages := make([][]byte, p.maxIdlePages)
				if numPage >= p.off.page+p.maxIdlePages {
					copy(pages, p.pages[p.off.page:p.off.page+p.maxIdlePages])
				} else {
					copy(pages, p.pages[p.off.page:])
					copy(pages[numPage-p.off.page:], p.pages[:p.maxIdlePages-numPage+p.off.page])
				}
				p.pages = pages
			} else {
				copy(p.pages, p.pages[p.off.page:endPage])
			}
		}
		p.off.page = 0
		endPage, endFrom, _ = p.tail()
	}
	return
}

// Write implements io.Writer Write method
func (p *PageBuffer) Write(data []byte) (n int, err error) {
	size := len(data)
	if size == 0 {
		return
	}
	endPage, endFrom := p.grow(size)
	for n < size {
		remain := p.mask + 1 - endFrom
		if size-n > remain {
			copy(p.pages[endPage][endFrom:], data[n:n+remain])
			n += remain
			endPage++
			endFrom = 0
		} else {
			copy(p.pages[endPage][endFrom:], data[n:])
			n = size
		}
	}
	p.size += n
	return
}

// Read implements io.Reader Read method
func (p *PageBuffer) Read(data []byte) (n int, err error) {
	size := len(data)
	for n < size {
		if p.size == 0 {
			err = io.EOF
			return
		}
		nn := p.mask + 1 - p.off.from
		if nn > p.size {
			nn = p.size
		}
		if size-n < nn {
			nn = size - n
			copy(data[n:], p.pages[p.off.page][p.off.from:p.off.from+nn])
			p.off.from += nn
		} else {
			copy(data[n:], p.pages[p.off.page][p.off.from:])
			p.off.page++
			p.off.from = 0
		}
		p.size -= nn
		n += nn
	}
	if p.size == 0 {
		p.Reset()
	}
	return
}
