package history_test

import (
	"testing"

	"github.com/gopherd/core/container/history"
)

func TestSlice(t *testing.T) {
	var recorder history.DefaultRecorder
	slice := history.NewSlice[int](&recorder, 0, 0)
	slice.Append(1, 2, 3)
	t.Logf("slice=%v", slice)
	slice.Insert(1, 10, 11)
	t.Logf("slice=%v", slice)
	recorder.Undo()
	t.Logf("slice=%v", slice)
}
