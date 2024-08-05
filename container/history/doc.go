// Package history provides data structures with undo functionality.
//
// This package includes implementations of Map, Set, and Slice data structures
// that support undo operations. It also provides a Recorder interface for
// managing undo actions.
//
// The main components of this package are:
//
//   - Map: A generic map that supports undo operations.
//   - Set: A generic set that supports undo operations.
//   - Slice: A generic slice that supports undo operations.
//   - Recorder: An interface for managing undo actions.
//   - BaseRecorder: A basic implementation of the Recorder interface.
//
// Each data structure (Map, Set, and Slice) is designed to work with a Recorder,
// which keeps track of changes and allows for undoing operations. The BaseRecorder
// provides a simple implementation of the Recorder interface that can be used
// with any of the data structures.
//
// Example usage:
//
//	recorder := &history.BaseRecorder{}
//	myMap := history.NewMap[string, int](recorder, 10)
//	myMap.Set("key", 42)
//	value, _ := myMap.Get("key")
//	fmt.Println(value) // Output: 42
//	recorder.Undo()
//	_, ok := myMap.Get("key")
//	fmt.Println(ok) // Output: false
//
// This package is useful for scenarios where you need to maintain a history of
// changes and potentially revert them, such as in text editors, game state
// management, or any application where an undo feature is desired.
package history
