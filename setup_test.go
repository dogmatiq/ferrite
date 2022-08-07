package ferrite_test

import (
	"os"
	"strings"

	"github.com/dogmatiq/ferrite"
	"github.com/dogmatiq/ferrite/spec"
)

// setUp configures the validation exit behavior such that it prints to stdout
// and does NOT exit the process on failure. This is useful inside testable
// examples.
func setUp() {
	ferrite.SetExitBehavior(
		os.Stdout,
		func(code int) {},
	)
}

// tearDown resets the global state after a test.
func tearDown() {
	ferrite.SetExitBehavior(os.Stderr, os.Exit)
	ferrite.ClearValidators()
	spec.ResetRegistry()

	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "FERRITE_") {
			i := strings.Index(env, "=")
			os.Unsetenv(env[:i])
		}
	}
}
