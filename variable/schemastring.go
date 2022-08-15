package variable

import (
	"fmt"
	"math"
	"reflect"

	"github.com/dogmatiq/ferrite/maybe"
)

// String is a schema that allows arbitrary string input.
type String interface {
	Schema

	// MinLength returns the minimum permitted length of the string.
	MinLength() (int, bool)

	// MaxLength returns the maximum permitted length of the string.
	MaxLength() (int, bool)
}

// TypedString is a string value depicted by type T.
type TypedString[T ~string] struct {
	MinLen, MaxLen maybe.Value[int]
}

// MinLength returns the minimum permitted length of the string.
func (s TypedString[T]) MinLength() (int, bool) {
	return s.MinLen.Get()
}

// MaxLength returns the maximum permitted length of the string.
func (s TypedString[T]) MaxLength() (int, bool) {
	return s.MaxLen.Get()
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
	if err := s.validate(v); err != nil {
		return Literal{}, err
	}

	return Literal{
		String: string(v),
	}, nil
}

// Unmarshal converts a literal value to it's native representation.
func (s TypedString[T]) Unmarshal(v Literal) (T, error) {
	n := T(v.String)
	return n, s.validate(n)
}

// Examples returns a (possibly empty) set of examples of valid values.
func (s TypedString[T]) Examples(hasOtherExamples bool) []TypedExample[T] {
	if hasOtherExamples {
		return nil
	}

	min, ok := s.MinLen.Get()
	if !ok {
		min = 1
	}

	max, ok := s.MaxLen.Get()
	if !ok {
		max = math.MaxInt
	}

	words := [...]T{
		"foo", "bar", "baz",
		"qux", "quux", "corge",
		"grault", "garply", "waldo",
		"fred", "plugh", "xyzzy",
	}

	var example T
	word := 0

	// Add enough words to meet the minimum requirement.
	for len(example) < min {
		example += words[word%len(words)]
		example += " "
		word++
	}

	// If the variable is too long truncate it to the max length, and ensure it
	// doesn't end in a space.
	if len(example) > max {
		example = example[:max-1] + "x"
	}

	return []TypedExample[T]{
		{
			Native:      example,
			Description: "randomly generated example",
		},
	}
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
	min, hasMin := s.MinLength()
	max, hasMax := s.MaxLength()

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
