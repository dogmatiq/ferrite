package variable

// Value is the value of an environment variable.
type Value interface {
	// Verbatim returns the string representation of the variable as it appears
	// in the environment.
	Verbatim() Literal

	// Canonical returns the canonical string representation of the variable.
	Canonical() Literal

	// IsDefault returns true if the value is the default, as opposed to being
	// set explicitly in the environment.
	IsDefault() bool
}

// valueOf is a value of an environment variable depicted by type T.
type valueOf[T any] struct {
	verbatim  Literal
	canonical Literal
	native    T
	isDefault bool
}

// Verbatim returns the string representation of the variable as it appears
// in the environment.
func (v valueOf[T]) Verbatim() Literal {
	return v.verbatim
}

// Canonical returns the canonical string representation of the variable.
func (v valueOf[T]) Canonical() Literal {
	return v.canonical
}

// IsDefault returns true if the value is the default, as opposed to being
// set explicitly in the environment.
func (v valueOf[T]) IsDefault() bool {
	return v.isDefault
}
