package history

// Recorder defines the interface for adding undo actions and performing undo operations.
type Recorder interface {
	// PushAction adds a new undo action to the recorder.
	PushAction(a Action)

	// UndoAll reverts all actions in reverse order and clears the recorder.
	// It returns the number of actions undone.
	UndoAll() (n int)

	// Undo reverts the last action in the recorder.
	// It returns true if an action was undone.
	Undo() bool
}

// Action represents a single undoable action.
type Action interface {
	// Undo reverts the changes made by this action.
	Undo()
}

// BaseRecorder implements the Recorder interface using a slice of Actions.
type BaseRecorder struct {
	actions []Action
}

// Len returns the number of actions in the recorder.
func (r *BaseRecorder) Len() int {
	return len(r.actions)
}

// PushAction appends a new Action to the BaseRecorder.
func (r *BaseRecorder) PushAction(a Action) {
	r.actions = append(r.actions, a)
}

// UndoAll reverts all actions in reverse order and clears the recorder.
// It returns the number of actions undone.
func (r *BaseRecorder) UndoAll() (n int) {
	n = len(r.actions)
	for i := len(r.actions) - 1; i >= 0; i-- {
		r.actions[i].Undo()
		r.actions[i] = nil // Allow garbage collection
	}
	r.actions = r.actions[:0]
	return
}

// Undo reverts the last action in the recorder.
// It returns true if an action was undone.
func (r *BaseRecorder) Undo() bool {
	if len(r.actions) == 0 {
		return false
	}
	i := len(r.actions) - 1
	r.actions[i].Undo()
	r.actions[i] = nil // Allow garbage collection
	r.actions = r.actions[:i]
	return true
}

// ValueUndoAction creates a new Action for undoing changes to a value of any type.
func ValueUndoAction[T any](ptr *T, old T) Action {
	return &valueUndoAction[T]{
		ptr: ptr,
		old: old,
	}
}

type valueUndoAction[T any] struct {
	ptr *T
	old T
}

// Undo restores the original value.
func (a *valueUndoAction[T]) Undo() {
	*a.ptr = a.old
}
