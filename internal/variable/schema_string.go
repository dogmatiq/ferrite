package variable

import (
	"fmt"
	"math"
	"reflect"

	"github.com/dogmatiq/ferrite/internal/maybe"
	"github.com/dogmatiq/ferrite/internal/reflectx"
)

// String is a schema that allows arbitrary string input.
type String interface {
	LengthLimited
}

// TypedString is a string value depicted by type T.
type TypedString[T ~string] struct {
	MinLen, MaxLen maybe.Value[int]
}

// MinLength returns the minimum permitted length of the native value.
func (s TypedString[T]) MinLength() (int, bool) {
	return s.MinLen.Get()
}

// MaxLength returns the maximum permitted length of the native value.
func (s TypedString[T]) MaxLength() (int, bool) {
	return s.MaxLen.Get()
}

// ExplainLengthError returns a human-readable description of the length
// constraints, for use in an error message.
func (s TypedString[T]) ExplainLengthError() string {
	return explainLengthError(s, "length")
}

// Type returns the type of the native value.
func (s TypedString[T]) Type() reflect.Type {
	return reflectx.TypeOf[T]()
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
func (s TypedString[T]) Examples(conservative bool) []TypedExample[T] {
	if conservative {
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
		if example != "" {
			example += " "
		}

		example += words[word%len(words)]
		word++
	}

	// If the variable is too long truncate it to the max length, and ensure it
	// doesn't end in a space.
	if len(example) > max {
		example = example[:max-1] + "x"
	}

	return []TypedExample[T]{
		{
			Native: example,
		},
	}
}

// validate returns an error if v is invalid.
func (s TypedString[T]) validate(v T) error {
	if min, ok := s.MinLen.Get(); ok && len(v) < min {
		return MinLengthError{s}
	}

	if max, ok := s.MaxLen.Get(); ok && len(v) > max {
		return MaxLengthError{s}
	}

	return nil
}
