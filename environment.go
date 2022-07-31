package ferrite

import "os"

// Environment is an interface for reading environment variables.
type Environment interface {
	// Lookup retrieves the value of an environment variable.
	//
	// If the variable is present in the environment, value (which may be empty)
	// is returned and ok is true. Otherwise, the value is empty and ok is false.
	Lookup(name string) (value string, ok bool)
}

// ShellEnvironment is the process's shell environment.
var ShellEnvironment Environment = shellEnvironment{}

// shellEnvironment reads environment variables from the process's shell
// environment.
type shellEnvironment struct{}

// Lookup retrieves the value of an environment variable.
//
// If the variable is present in the environment, value (which may be empty)
// is returned and ok is true. Otherwise, the value is empty and ok is false.
func (shellEnvironment) Lookup(name string) (value string, ok bool) {
	return os.LookupEnv(name)
}

// MemoryEnvironment is an in-memory map of environment variables
type MemoryEnvironment map[string]string

// Lookup retrieves the value of an environment variable.
//
// If the variable is present in the environment, value (which may be empty)
// is returned and ok is true. Otherwise, the value is empty and ok is false.
func (e MemoryEnvironment) Lookup(name string) (value string, ok bool) {
	value, ok = e[name]
	return value, ok
}
