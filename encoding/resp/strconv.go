package resp

import (
	"bytes"
	"errors"
	"strconv"
)

var ErrNumberSyntax = errors.New("invalid number syntax")
var ErrNumberRange = errors.New("invalid number range")

const maxUint64 = 1<<64 - 1

func lower(c byte) byte {
	return c | ('x' - 'X')
}

func btoi64(s []byte, bitSize int) (int64, error) {
	sLen := len(s)
	if sLen == 0 {
		return 0, ErrNumberSyntax
	}
	if strconv.IntSize == 32 && (0 < sLen && sLen < 10) ||
		strconv.IntSize == 64 && (0 < sLen && sLen < 19) {
		// Fast path for small integers that fit int type.
		s0 := s
		if s[0] == '-' || s[0] == '+' {
			s = s[1:]
			if len(s) < 1 {
				return 0, ErrNumberSyntax
			}
		}

		n := int64(0)
		for _, ch := range []byte(s) {
			ch -= '0'
			if ch > 9 {
				return 0, ErrNumberSyntax
			}
			n = n*10 + int64(ch)
		}
		if s0[0] == '-' {
			n = -n
		}
		return n, nil
	}

	// Pick off leading sign.
	neg := false
	if s[0] == '+' {
		s = s[1:]
	} else if s[0] == '-' {
		neg = true
		s = s[1:]
	}

	// Convert unsigned and check range.
	un, err := btou64(s, bitSize)
	if err != nil {
		return 0, err
	}

	if bitSize == 0 {
		bitSize = int(strconv.IntSize)
	}

	cutoff := uint64(1 << uint(bitSize-1))
	if !neg && un >= cutoff {
		return int64(cutoff - 1), ErrNumberRange
	}
	if neg && un > cutoff {
		return -int64(cutoff), ErrNumberRange
	}
	n := int64(un)
	if neg {
		n = -n
	}
	return n, nil
}

func btou64(s []byte, bitSize int) (uint64, error) {
	if len(s) == 0 {
		return 0, ErrNumberSyntax
	}

	if bitSize == 0 {
		bitSize = int(strconv.IntSize)
	}

	// Cutoff is the smallest number such that cutoff*base > maxUint64.
	// Use compile-time constants for common cases.
	var cutoff = uint64(maxUint64/10 + 1)

	maxVal := uint64(1)<<uint(bitSize) - 1

	var n uint64
	for _, c := range []byte(s) {
		var d byte
		switch {
		case '0' <= c && c <= '9':
			d = c - '0'
		case 'a' <= lower(c) && lower(c) <= 'z':
			d = lower(c) - 'a' + 10
		default:
			return 0, ErrNumberSyntax
		}

		if d >= byte(10) {
			return 0, ErrNumberSyntax
		}

		if n >= cutoff {
			// n*base overflows
			return maxVal, ErrNumberRange
		}
		n *= 10

		n1 := n + uint64(d)
		if n1 < n || n1 > maxVal {
			// n+v overflows
			return maxVal, ErrNumberRange
		}
		n = n1
	}

	return n, nil
}

const (
	minItoa = -128
	maxItoa = 32768
)

var (
	itobOffset [maxItoa - minItoa + 1]uint32
	itobBuffer []byte
)

func init() {
	var b bytes.Buffer
	for i := range itobOffset {
		itobOffset[i] = uint32(b.Len())
		b.WriteString(strconv.Itoa(i + minItoa))
	}
	itobBuffer = b.Bytes()
}

func i64tob(i int64) []byte {
	if i >= minItoa && i <= maxItoa {
		beg := itobOffset[i-minItoa]
		if i == maxItoa {
			return itobBuffer[beg:]
		}
		end := itobOffset[i-minItoa+1]
		return itobBuffer[beg:end]
	}
	return strconv.AppendInt(nil, i, 10)
}

func itob(i int) []byte {
	return i64tob(int64(i))
}
