package variable

import (
	"os"
	"strings"
	"sync"
)

// Environment is an interface for interacting with the environment.
type Environment interface {
	// Get returns the value of an environment variable.
	Get(string) Literal

	// Set sets the value of an environment variable.
	Set(string, Literal)

	// Unset removes an environment variable.
	Unset(string)

	// Range calls fn for each environment variable.
	//
	// It stops iterating if fn returns false.
	Range(func(string, Literal) bool)
}

// OSEnvironment is the operating system's actual environment.
var OSEnvironment osEnvironment

// OperatingSystem is an implementation of Environment that uses the operating
// system's actual environment.
type osEnvironment struct{}

// Get returns the value of an environment variable.
func (osEnvironment) Get(n string) Literal {
	return Literal{
		String: os.Getenv(string(n)),
	}
}

// Set sets the value of an environment variable.
func (osEnvironment) Set(n string, v Literal) {
	if err := os.Setenv(string(n), v.String); err != nil {
		panic(err)
	}
}

// Unset removes an environment variable.
func (osEnvironment) Unset(n string) {
	if err := os.Unsetenv(string(n)); err != nil {
		panic(err)
	}
}

// Range calls fn for each environment variable.
//
// It stops iterating if fn returns false. It returns true if iteration
// reached the env.
func (osEnvironment) Range(fn func(string, Literal) bool) {
	for _, env := range os.Environ() {
		i := strings.IndexByte(env, '=')
		n := env[:i]
		v := Literal{
			String: env[i+1:],
		}

		if !fn(n, v) {
			return
		}
	}
}

// MemoryEnvironment is an Environment that stores environment variables in this
// processes's memory, as opposed to using the operating system environment.
type MemoryEnvironment struct {
	m sync.Map // map[string]Literal
}

// Get returns the value of an environment variable.
func (e *MemoryEnvironment) Get(n string) Literal {
	if v, ok := e.m.Load(n); ok {
		return v.(Literal)
	}

	return Literal{}
}

// Set sets the value of an environment variable.
func (e *MemoryEnvironment) Set(n string, v Literal) {
	e.m.Store(n, v)
}

// Unset removes an environment variable.
func (e *MemoryEnvironment) Unset(n string) {
	e.m.Delete(n)
}

// Range calls fn for each environment variable.
//
// It stops iterating if fn returns false. It returns true if iteration
// reached the env.
func (e *MemoryEnvironment) Range(fn func(string, Literal) bool) {
	e.m.Range(
		func(k, v any) bool {
			return fn(
				k.(string),
				v.(Literal),
			)
		},
	)
}

// Snapshot is a snapshot of an environment.
type Snapshot struct {
	values map[string]Literal
}

// TakeSnapshot takes a snapshot of the variables within an environment.
func TakeSnapshot(env Environment) Snapshot {
	sn := Snapshot{
		values: map[string]Literal{},
	}

	env.Range(func(n string, v Literal) bool {
		sn.values[n] = v
		return true
	})

	return sn
}

// RestoreSnapshot restores an environment to the state it was in when the
// given snapshot was taken.
func RestoreSnapshot(env Environment, sn Snapshot) {
	env.Range(func(n string, v Literal) bool {
		if _, ok := sn.values[n]; !ok {
			env.Unset(n)
		}
		return true
	})

	for n, v := range sn.values {
		env.Set(n, v)
	}
}
