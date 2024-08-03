package operator_test

import (
	"testing"

	"github.com/gopherd/core/operator"
)

func TestOr(t *testing.T) {
	if v := operator.Or("x", "y"); v != "x" {
		t.Fatalf("want %q, but got %q", "x", v)
	}
	if v := operator.Or("", "y"); v != "y" {
		t.Fatalf("want %q, but got %q", "y", v)
	}
	if v := operator.Or(1, 2); v != 1 {
		t.Fatalf("want 1, but got %d", v)
	}
	if v := operator.Or(0, 2); v != 2 {
		t.Fatalf("want 1, but got %d", v)
	}
	if v := operator.Or(true, true); v != true {
		t.Fatalf("want true, but got %v", v)
	}
	if v := operator.Or(true, false); v != true {
		t.Fatalf("want true, but got %v", v)
	}
	if v := operator.Or(false, true); v != true {
		t.Fatalf("want true, but got %v", v)
	}
	if v := operator.Or(false, false); v != false {
		t.Fatalf("want false, but got %v", v)
	}
}
