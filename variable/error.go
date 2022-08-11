package variable

import "fmt"

// Error is an error that indicates a problem parsing or validating an
// environment variable.
type Error interface {
	error
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

// ValidationError indicates a problem parsing or validating an environment
// variable value.
type ValidationError interface {
	Error

	// Verbatim returns the offending value.
	Verbatim() String

	// Reason returns a human-readable explanation of why the value is invalid.
	Reason() string
}

// formatValidationError formats the message of a validation error.
func formatValidationError(err ValidationError) string {
	return fmt.Sprintf(
		"%s (%s) is invalid: %s",
		err.Name(),
		err.Verbatim(),
		err.Reason(),
	)
}
