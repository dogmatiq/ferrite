package ferrite

// Required is the application-facing representation of an environment variable
// that must have a value.
type Required[T any] interface {
	// Value returns the parsed and validated value of the environment variable.
	//
	// It panics if the environment variable is undefined or invalid.
	Value() T
}

// Optional is the application-facing representation of an environment variable
// that may optionally be undefined.
type Optional[T any] interface {
	// Value returns the parsed and validated value of the environment variable,
	// if it is defined.
	//
	// If the environment variable is not defined (and there is no default
	// value), ok is false; otherwise, ok is true and v is the value.
	//
	// It panics if the environment variable is defined but invalid.
	Value() (v T, ok bool)
}
