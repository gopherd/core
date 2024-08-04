/*
Package errkit provides utilities for handling multiple errors and a flexible error code mechanism for Go applications.

This package offers two main features:
1. A List type for managing multiple errors
2. An error code mechanism for structured error handling

Key features:
  - Collect and manage multiple errors using the List type
  - Associate integer error codes with errors
  - Wrap existing errors with error codes
  - Add context information to errors
  - Check error types using error codes

List Usage:

The List type allows you to collect and manage multiple errors:

	var errList errkit.List
	errList.Append(err1)
	errList.Append(err2)
	errList.Append(err3)
	if err := errList.Err(); err != nil {
		// handle combined error
	}

List methods:
  - Append(err error): Adds an error to the list if it's non-nil
  - Len() int: Returns the number of errors in the list
  - First() error: Returns the first error in the list, or nil if empty
  - Last() error: Returns the last error in the list, or nil if empty
  - Err() error: Returns all errors as a single error, or nil if no errors

Error Code Mechanism Usage:

1. Define your error code type and constants:

	type MyErrno int

	const (
		EUnknown MyErrno = errkit.EUnknown
		EOK      MyErrno = errkit.EOK
		ENotFound MyErrno = iota + 1
		EInvalidInput
		// Define more error codes as needed
	)

2. Create errors with error codes:

	func FindUser(id int) error {
		// Simulate a database lookup
		if id < 0 {
			return errkit.New(EInvalidInput, fmt.Errorf("invalid user id: %d", id))
		}
		// User not found scenario
		return errkit.New(ENotFound, fmt.Errorf("user with id %d not found", id))
	}

3. Add context to errors:

	func ProcessUser(id int) error {
		if err := FindUser(id); err != nil {
			return errkit.NewWithContext(errkit.Errno(err), err, "failed to process user")
		}
		// Process user...
		return nil
	}

4. Check error types:

	err := ProcessUser(42)
	switch errkit.Errno(err) {
	case ENotFound:
		fmt.Println("User not found")
	case EInvalidInput:
		fmt.Println("Invalid input provided")
	default:
		fmt.Println("An error occurred:", err)
	}

5. Use the Is function for more idiomatic error checking:

	if errkit.Is(err, ENotFound) {
		fmt.Println("User not found")
	}

By using errkit, you can create more structured and easily identifiable errors
in your Go applications, improving error handling and debugging. The List type
provides a convenient way to collect and manage multiple errors, while the
error code mechanism allows for more detailed and type-safe error handling.
*/
package errkit
