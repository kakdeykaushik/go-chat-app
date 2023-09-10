package utils

import "log"

// RemoveIndex removes an element at the specified index from a slice of any type.
// The original slice is not modified; a new slice is returned without the element.
//
// Parameters:
//   - s: The input slice.
//   - i: The index of the element to be removed.
//
// Returns:
//   - A new slice with the specified element removed.
//
// Example usage:
//
//	s := []int{1, 2, 3, 4, 5}
//	modifiedSlice := RemoveIndex(s, 2) // Removes the element at index 2 (value 3)
func RemoveIndex[T any](s []T, i int) []T {
	return append(s[:i], s[i+1:]...)
}

// HandleError checks if an error occurred. If so, it logs a formatted error message
// and exits the program with a fatal error.
//
// Parameters:
//   - err: The error to check.
//   - msg: A custom error message to include in the log output.
//
// Example usage:
//
//	err := someFunction()
//	HandleError(err, "Error while performing someFunction")
func HandleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s - %v\n", msg, err)
	}
}
