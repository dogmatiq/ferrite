package ferrite_test

import (
	"os"
	"strings"

	"github.com/dogmatiq/ferrite"
	"github.com/dogmatiq/ferrite/variable"
)

// example is a helper function that sets up the global state for a testable
// example. It returns a function that resets the global state after the test.
func example() func() {
	ferrite.XDefaultInitOptions.Err = os.Stdout
	ferrite.XDefaultInitOptions.Exit = func(code int) {}

	return tearDown
}

// tearDown resets the environemnt and Ferrite global state after a test.
func tearDown() {
	ferrite.XDefaultInitOptions.Err = os.Stderr
	ferrite.XDefaultInitOptions.Exit = os.Exit

	variable.DefaultRegistry.Reset()

	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "FERRITE_") {
			i := strings.Index(env, "=")
			os.Unsetenv(env[:i])
		}
	}
}
