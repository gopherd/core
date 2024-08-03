package erron_test

import (
	"testing"

	. "github.com/gopherd/core/erron"
)

func testingDo() error {
	return Throwf("testing error")
}

func TestThrow(t *testing.T) {
	err := testingDo()
	t.Log(err)
}

func testingTry() error {
	panic("try panic")
}

func testingTry2() error {
	return testingTry()
}

func TestTry(t *testing.T) {
	err := Try(testingTry2)
	t.Log(err)
}
