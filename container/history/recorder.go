// Package history provides a simple undo mechanism for recording and reverting changes.
package history

// Recorder defines the interface for adding records and performing undo operations.
type Recorder interface {
	// AddRecord adds a new record to the recorder.
	AddRecord(r Record)

	// Undo reverts all recorded changes in reverse order.
	Undo()
}

// Record represents a single undoable action.
type Record interface {
	// Undo reverts the changes made by this record.
	Undo()
}

// DefaultRecorder implements the Recorder interface using a slice of Records.
type DefaultRecorder struct {
	records []Record
}

// AddRecord appends a new Record to the DefaultRecorder.
func (recorder *DefaultRecorder) AddRecord(r Record) {
	recorder.records = append(recorder.records, r)
}

// Undo reverts all records in reverse order and clears the recorder.
func (recorder *DefaultRecorder) Undo() {
	for i := len(recorder.records) - 1; i >= 0; i-- {
		recorder.records[i].Undo()
		recorder.records[i] = nil // Allow garbage collection
	}
	recorder.records = recorder.records[:0]
}

// ValueRecord creates a new Record for undoing changes to a value of any type.
func ValueRecord[T any](ptr *T, old T) Record {
	return &valueRecord[T]{
		ptr: ptr,
		old: old,
	}
}

// valueRecord is a generic implementation of the Record interface for value types.
type valueRecord[T any] struct {
	ptr *T
	old T
}

// Undo restores the original value.
func (vr *valueRecord[T]) Undo() {
	*vr.ptr = vr.old
}
