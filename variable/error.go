package variable

// Error is an error that indicates a problem parsing or validating an
// environment variable.
type Error interface {
	error

	// Name returns the name of the environment variable.
	Name() Name
}
