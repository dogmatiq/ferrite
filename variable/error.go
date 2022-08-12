package variable

import "fmt"

// Error is an error that indicates a problem parsing or validating an
// environment variable.
type Error interface {
	error

	// Name returns the name of the environment variable.
	Name() Name
}

// SpecError represents a problem with a variable specification itself, rather
// than the variable's value.
type SpecError struct {
	name  Name
	cause error
}

// Name returns the name of the environment variable.
func (e SpecError) Name() Name {
	return e.name
}

func (e SpecError) Unwrap() error {
	return e.cause
}

func (e SpecError) Error() string {
	if e.name == "" {
		return fmt.Sprintf("invalid specification: %s", e.cause)
	}

	return fmt.Sprintf("invalid specification for %s: %s", e.name, e.cause)
}

// ValidationError indicates that a well-formed variable value has other issues.
type ValidationError interface {
	Error

	// Literal returns the invalid value.
	Literal() Literal

	// Unwrap returns the underlying error.
	Unwrap() error
}

type validationError struct {
	name    Name
	literal Literal
	reason  error
}

func (e validationError) Name() Name {
	return e.name
}

func (e validationError) Literal() Literal {
	return e.literal
}

func (e validationError) Unwrap() error {
	return e.reason
}

func (e validationError) Error() string {
	return fmt.Sprintf(
		"%s (%s) is invalid: %s",
		e.name,
		e.literal,
		e.reason,
	)
}
