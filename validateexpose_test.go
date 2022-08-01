package ferrite

import (
	"io"
	"os"
	"strings"
)

// SetExitBehavior sets the behavior of the ValidateEnvironment() function when
// the environment is invalid.
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

// Setup configures the validation exit behavior such that it prints to stdout
// and does NOT exit the process on failure. This is useful inside testable
// examples.
func Setup() {
	SetExitBehavior(os.Stdout, func(code int) {})
}

// Teardown resets the global state after a test.
func Teardown() {
	SetExitBehavior(os.Stderr, os.Exit)
	registry = nil

	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "FERRITE_") {
			i := strings.Index(env, "=")
			os.Unsetenv(env[:i])
		}
	}
}
