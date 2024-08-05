package history_test

import (
	"testing"

	"github.com/gopherd/core/container/history"
)

func TestDefaultRecorder_AddRecord(t *testing.T) {
	recorder := &history.BaseRecorder{}

	t.Run("Add single record", func(t *testing.T) {
		value := 5
		action := history.ValueUndoAction(&value, 0)
		recorder.PushAction(action)
	})

	t.Run("Add multiple records", func(t *testing.T) {
		value1, value2 := 10, "test"
		action1 := history.ValueUndoAction(&value1, 0)
		action2 := history.ValueUndoAction(&value2, "")

		recorder.PushAction(action1)
		recorder.PushAction(action2)
	})
}

func TestDefaultRecorder_Undo(t *testing.T) {
	t.Run("Undo single record", func(t *testing.T) {
		recorder := &history.BaseRecorder{}
		value := 5
		originalValue := 0
		action := history.ValueUndoAction(&value, originalValue)
		recorder.PushAction(action)

		recorder.Undo()

		if value != originalValue {
			t.Errorf("Expected value to be %d after undo, got %d", originalValue, value)
		}
	})

	t.Run("Undo multiple records", func(t *testing.T) {
		recorder := &history.BaseRecorder{}
		value1, value2 := 10, "test"
		originalValue1, originalValue2 := 0, ""
		action1 := history.ValueUndoAction(&value1, originalValue1)
		action2 := history.ValueUndoAction(&value2, originalValue2)

		recorder.PushAction(action1)
		recorder.PushAction(action2)

		recorder.UndoAll()

		if value1 != originalValue1 || value2 != originalValue2 {
			t.Errorf("Expected values to be reset after undo, got %d and %s", value1, value2)
		}
	})

	t.Run("Undo with no records", func(t *testing.T) {
		recorder := &history.BaseRecorder{}
		ok := recorder.Undo() // Should not panic
		if ok {
			t.Error("Expected Undo to return false with no records")
		}
	})
}

func TestValueRecord(t *testing.T) {
	t.Run("Integer value", func(t *testing.T) {
		value := 5
		originalValue := 0
		action := history.ValueUndoAction(&value, originalValue)
		action.Undo()
		if value != originalValue {
			t.Errorf("Expected value to be %d after undo, got %d", originalValue, value)
		}
	})

	t.Run("String value", func(t *testing.T) {
		value := "new"
		originalValue := "old"
		record := history.ValueUndoAction(&value, originalValue)
		record.Undo()
		if value != originalValue {
			t.Errorf("Expected value to be '%s' after undo, got '%s'", originalValue, value)
		}
	})

	t.Run("Struct value", func(t *testing.T) {
		type testStruct struct {
			field int
		}
		value := testStruct{field: 10}
		originalValue := testStruct{field: 5}
		record := history.ValueUndoAction(&value, originalValue)
		record.Undo()
		if value != originalValue {
			t.Errorf("Expected struct to be %+v after undo, got %+v", originalValue, value)
		}
	})
}

func TestRecorder_Interface(t *testing.T) {
	var _ history.Recorder = &history.BaseRecorder{}
}

func TestUndoAction_Interface(t *testing.T) {
	var _ history.Action = history.ValueUndoAction(&struct{}{}, struct{}{})
}

// Benchmarks
func BenchmarkDefaultRecorder_AddRecord(b *testing.B) {
	recorder := &history.BaseRecorder{}
	value := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		recorder.PushAction(history.ValueUndoAction(&value, i))
	}
}

func BenchmarkDefaultRecorder_Undo(b *testing.B) {
	recorder := &history.BaseRecorder{}
	value := 0
	for i := 0; i < b.N; i++ {
		recorder.PushAction(history.ValueUndoAction(&value, i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		recorder.Undo()
	}
}
