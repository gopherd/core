package history

type Recorder interface {
	AddRecord(r Record)
	Undo()
}

type Record interface {
	Undo()
}

type DefaultRecorder struct {
	records []Record
}

func (recorder *DefaultRecorder) AddRecord(r Record) {
	recorder.records = append(recorder.records, r)
}

func (recorder *DefaultRecorder) Undo() {
	for i := len(recorder.records) - 1; i >= 0; i-- {
		recorder.records[i].Undo()
		recorder.records[i] = nil
	}
	recorder.records = recorder.records[:0]
}

func ValueRecord[T any](ptr *T, old T) Record {
	return &valueRecord[T]{
		ptr: ptr,
		old: old,
	}
}

type valueRecord[T any] struct {
	ptr *T
	old T
}

func (vc *valueRecord[T]) Undo() {
	*(vc.ptr) = vc.old
}
