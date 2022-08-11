package variable

import (
	"os"
	"strings"
)

// Environment is an interface for interacting with the environment.
type Environment interface {
	// Get returns the value of an environment variable.
	Get(Name) Literal

	// Set sets the value of an environment variable.
	Set(Name, Literal)

	// Unset removes an environment variable.
	Unset(Name)

	// ForEach calls fn for each environment variable.
	//
	// It stops iterating if fn returns false. It returns true if iteration
	// reached the env.
	ForEach(func(Name, Literal) bool) bool
}

// OSEnvironment is the operating system's actual environment.
var OSEnvironment osEnvironment

// OperatingSystem is an implementation of Environment that uses the operating
// system's actual environment.
type osEnvironment struct{}

// Get returns the value of an environment variable.
func (osEnvironment) Get(n Name) Literal {
	return Literal(os.Getenv(string(n)))
}

// Set sets the value of an environment variable.
func (osEnvironment) Set(n Name, v Literal) {
	if err := os.Setenv(string(n), string(v)); err != nil {
		panic(err)
	}
}

// Unset removes an environment variable.
func (osEnvironment) Unset(n Name) {
	if err := os.Unsetenv(string(n)); err != nil {
		panic(err)
	}
}

// ForEach calls fn for each environment variable.
//
// It stops iterating if fn returns false. It returns true if iteration
// reached the env.
func (osEnvironment) ForEach(fn func(Name, Literal) bool) bool {
	for _, env := range os.Environ() {
		i := strings.IndexByte(env, '=')
		n := Name(env[:i])
		v := Literal(env[i+1:])

		if !fn(n, v) {
			return false
		}
	}

	return true
}

// Snapshot is a snapshot of an environment.
type Snapshot struct {
	values map[Name]Literal
}

// TakeSnapshot takes a snapshot of the variables within an environment.
func TakeSnapshot(env Environment) Snapshot {
	sn := Snapshot{
		values: map[Name]Literal{},
	}

	env.ForEach(func(n Name, v Literal) bool {
		sn.values[n] = v
		return true
	})

	return sn
}

// RestoreSnapshot restores an environment to the state it was in when the
// given snapshot was taken.
func RestoreSnapshot(env Environment, sn Snapshot) {
	env.ForEach(func(n Name, v Literal) bool {
		if _, ok := sn.values[n]; !ok {
			env.Unset(n)
		}
		return true
	})

	for n, v := range sn.values {
		env.Set(n, v)
	}
}
