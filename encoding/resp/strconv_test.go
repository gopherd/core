package resp

import (
	"math"
	"strconv"
	"testing"
)

func Test_Btoi(t *testing.T) {
	var (
		MaxUint64 = []byte("18446744073709551615")
		MaxInt64  = []byte("9223372036854775807")
		MinInt64  = []byte("-9223372036854775808")
	)
	if u, err := btou64(MaxUint64, 64); err != nil {
		t.Errorf("btou64 %q error: %v", MaxUint64, err)
	} else if u != math.MaxUint64 {
		t.Errorf("btou64 %q mismatch: %d vs %d", MaxUint64, u, uint64(math.MaxUint64))
	}

	if u, err := btoi64(MaxInt64, 64); err != nil {
		t.Errorf("btoi64 %q error: %v", MaxInt64, err)
	} else if u != math.MaxInt64 {
		t.Errorf("btoi64 %q mismatch: %d vs %d", MaxInt64, u, uint64(math.MaxInt64))
	}

	if u, err := btoi64(MinInt64, 64); err != nil {
		t.Errorf("btoi64 %q error: %v", MinInt64, err)
	} else if u != math.MinInt64 {
		t.Errorf("btoi64 %q mismatch: %d vs %d", MinInt64, u, int64(math.MinInt64))
	}

	var int64s = map[int64]string{}
	var min = -int64(1 << 20)
	var max = int64(1 << 20)
	for i := min; i < max; i++ {
		int64s[i] = strconv.FormatInt(i, 10)
	}
	for i, s := range int64s {
		i64, err := btoi64([]byte(s), 64)
		if err != nil {
			t.Fatalf("btoi64 error: %v", err)
		}
		if i != i64 {
			t.Fatalf("btoi64 failed: s=%s, want=%d, got=%d", s, i, i64)
		}
	}
}

func Test_Itob(t *testing.T) {
	var int64s = map[int64]string{}
	var min = -int64(1 << 20)
	var max = int64(1 << 20)
	for i := min; i < max; i++ {
		int64s[i] = strconv.FormatInt(i, 10)
	}
	for i, s := range int64s {
		b := i64tob(i)
		if s != string(b) {
			t.Fatalf("i64tob failed: i=%d, want=%s, got=%s", i, s, string(b))
		}
	}
}
