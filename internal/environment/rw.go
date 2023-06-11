package environment

import (
	"os"
	"strings"
)

// Get returns the value of the environment variable with the given name.
func Get(n string) string {
	return os.Getenv(n)
}

// Set sets the value of the environment variable with the given name, or panics
// if unable to do so.
func Set(n, v string) {
	if err := os.Setenv(n, v); err != nil {
		panic(err)
	}
}

// Unset unsets the environment variable with the given name, or panics if
// unable to do so.
func Unset(n string) {
	if err := os.Unsetenv(n); err != nil {
		panic(err)
	}
}

// Range calls fn for each environment variable.
//
// It stops iterating if fn returns false. It returns true if iteration
// reached the end.
func Range(fn func(n, v string) bool) {
	for _, env := range os.Environ() {
		i := strings.IndexByte(env, '=')
		n := env[:i]
		v := env[i+1:]

		if !fn(n, v) {
			return
		}
	}
}
