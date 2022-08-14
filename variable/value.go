package variable

import (
	"fmt"
	"regexp"
	"strings"
)

// Literal the string representation of an environment variable value.
type Literal struct {
	String string
}

// needsQuotes is a pattern that matches values that require shell escaping.
var needsQuotes = regexp.MustCompile(`[^\w@%+=:,./-]`)

// Quote returns an escaped string representation of the value that can be used
// directly in the shell.
func (l Literal) Quote() string {
	if !needsQuotes.MatchString(l.String) {
		return l.String
	}

	return `'` + strings.ReplaceAll(l.String, `'`, `'"'"'`) + `'`
}

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

// ValueError indicates that there is a problem with a variable's value.
type ValueError interface {
	Error

	// Literal returns the invalid value.
	Literal() Literal

	// Unwrap returns the underlying cause of the value error.
	Unwrap() error

	// AcceptVisitor passes the error to the appropriate method of v.
	AcceptVisitor(v ValueErrorVisitor)
}

// ValueErrorVisitor dispatches based on the the cause of a value error.
type ValueErrorVisitor interface {
	SchemaErrorVisitor

	VisitGenericError(error)
}

// valueError indicates that there is a problem with a variable's value.
type valueError struct {
	name    string
	literal Literal
	cause   error
}

func (e valueError) Name() string {
	return e.name
}

func (e valueError) Literal() Literal {
	return e.literal
}

func (e valueError) Unwrap() error {
	return e.cause
}

func (e valueError) AcceptVisitor(v ValueErrorVisitor) {
	switch err := e.cause.(type) {
	case SchemaError:
		err.AcceptVisitor(v)
	default:
		v.VisitGenericError(err)
	}
}

func (e valueError) Error() string {
	return fmt.Sprintf(
		"value of %s (%s) is invalid: %s",
		e.name,
		e.literal.Quote(),
		e.cause,
	)
}
