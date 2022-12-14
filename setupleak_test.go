package ferrite

// This file is in the "ferrite" package (not "ferrite_test") it is used to
// manipulate some of Ferrite's global state (*gasp*) during tests.

import "io"

// SetExitBehavior sets the behavior of the Init() function when it wants to
// exit the process.
//
// w is the writer to which the validation result is written. fn is a function
// that is called with the process's exit code.
//
// By default these values are os.Stderr and os.Exit, respectively.
//
// Use ferrite_test.Teardown() to undo changes made by this function.
func SetExitBehavior(w io.Writer, fn func(code int)) {
	output = w
	exit = fn
}
