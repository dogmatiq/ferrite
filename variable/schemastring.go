package variable

import (
	"fmt"
	"reflect"

	"github.com/dogmatiq/ferrite/maybe"
)

// String is a schema that allows arbitrary string input.
type String interface {
	Schema

	// MinLength returns the minimum permitted length of the string.
	MinLength() maybe.Value[int]

	// MaxLength returns the maximum permitted length of the string.
	MaxLength() maybe.Value[int]
}

// TypedString is a string value depicted by type T.
type TypedString[T ~string] struct {
	MinLen, MaxLen maybe.Value[int]
}

// MinLength returns the minimum permitted length of the string.
func (s TypedString[T]) MinLength() maybe.Value[int] {
	return s.MinLen
}

// MaxLength returns the maximum permitted length of the string.
func (s TypedString[T]) MaxLength() maybe.Value[int] {
	return s.MaxLen
}

// Type returns the type of the native value.
func (s TypedString[T]) Type() reflect.Type {
	return typeOf[T]()
}

// Finalize prepares the schema for use.
//
// It returns an error if schema is invalid.
func (s TypedString[T]) Finalize() error {
	min := 1

	if v, ok := s.MinLen.Get(); ok {
		if v < min {
			return fmt.Errorf("minimum length: must be at least %d", min)
		}
		min = v
	}

	if v, ok := s.MaxLen.Get(); ok {
		if v < min {
			return fmt.Errorf("maximum length: must be at least %d", min)
		}
	}

	return nil
}

// AcceptVisitor passes s to the appropriate method of v.
func (s TypedString[T]) AcceptVisitor(v SchemaVisitor) {
	v.VisitString(s)
}

// Marshal converts a value to its literal representation.
func (s TypedString[T]) Marshal(v T) (Literal, error) {
	return Literal(v), s.validate(v)
}

// Unmarshal converts a literal value to it's native representation.
func (s TypedString[T]) Unmarshal(v Literal) (T, error) {
	n := T(v)
	return n, s.validate(n)
}

// validate returns an error if v is invalid.
func (s TypedString[T]) validate(v T) error {
	if min, ok := s.MinLen.Get(); ok && len(v) < min {
		return MinLengthError{s}
	}

	if max, ok := s.MaxLen.Get(); ok && len(v) > max {
		return MinLengthError{s}
	}

	return nil
}

// MinLengthError indicates that a string value was shorter than the minimum
// permitted length.
type MinLengthError struct {
	String String
}

var _ SchemaError = MinLengthError{}

// Schema returns the schema that was violated.
func (e MinLengthError) Schema() Schema {
	return e.String
}

// AcceptVisitor passes the error to the appropriate method of v.
func (e MinLengthError) AcceptVisitor(v SchemaErrorVisitor) {
	v.VisitMinLengthError(e)
}

func (e MinLengthError) Error() string {
	return fmt.Sprintf("too short, %s", explainLengthError(e.String))
}

// MaxLengthError indicates that a string value was greater than the maximum
// permitted value.
type MaxLengthError struct {
	String String
}

var _ SchemaError = MaxLengthError{}

// Schema returns the schema that was violated.
func (e MaxLengthError) Schema() Schema {
	return e.String
}

// AcceptVisitor passes the error to the appropriate method of v.
func (e MaxLengthError) AcceptVisitor(v SchemaErrorVisitor) {
	v.VisitMaxLengthError(e)
}

func (e MaxLengthError) Error() string {
	return fmt.Sprintf("too long, %s", explainLengthError(e.String))
}

func explainLengthError(s String) string {
	min, hasMin := s.MinLength().Get()
	max, hasMax := s.MaxLength().Get()

	if !hasMin {
		return fmt.Sprintf("expected length to be %d bytes or fewer", max)
	}

	if !hasMax {
		return fmt.Sprintf("expected length to be %d bytes or more", max)
	}

	if min == max {
		return fmt.Sprintf("expected length to be exactly %d bytes", min)
	}

	return fmt.Sprintf("expected length to be between %d and %d bytes", min, max)
}
