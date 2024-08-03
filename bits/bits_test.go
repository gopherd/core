package bits

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCount(t *testing.T) {
	if c := Count32(0xFFFFFFFF); c != 32 {
		t.Errorf("bit count of 0xFFFFFFFF != 32, got %d", c)
	}
	if c := Count64(0xFFFFFFFFFFFFFFFF); c != 64 {
		t.Errorf("bit count of 0xFFFFFFFFFFFFFFFF != 64, got %d", c)
	}

	// generate normalizerArray
	var buf bytes.Buffer
	var buf2 bytes.Buffer
	for i := uint16(0); i <= 255; i++ {
		n := Count16(i)
		x := 1<<uint(n) - 1
		fmt.Fprintf(&buf, "%d,", x)
		if (i+1)%8 == 0 {
			buf.WriteString("\n")
		}

		na := Count16(i & 0xF)
		nb := Count16(i & 0xF0)
		xa := 1<<uint(na) - 1
		xb := 1<<uint(nb) - 1
		x2 := xa | (xb << 4)
		fmt.Fprintf(&buf2, "%d,", x2)
		if (i+1)%8 == 0 {
			buf2.WriteString("\n")
		}
	}
	t.Logf(buf.String())
	t.Logf("\n")
	t.Logf(buf2.String())
}
