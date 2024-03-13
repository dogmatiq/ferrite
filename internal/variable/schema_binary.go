package variable

import (
	"fmt"
	"math/rand"
	"reflect"

	"github.com/dogmatiq/ferrite/internal/maybe"
	"github.com/dogmatiq/ferrite/internal/reflectx"
)

// Binary is a schema that allows input of binary data.
type Binary interface {
	LengthLimited

	// EncodingDescription returns a short (one word, ideally) human-readable
	// description of the encoding scheme.
	EncodingDescription() string
}

// TypedBinary is a string value depicted by type T.
type TypedBinary[T ~[]B, B ~byte] struct {
	Marshaler      Marshaler[T]
	MinLen, MaxLen maybe.Value[int]
	EncodingDesc   string
}

// MinLength returns the minimum permitted length of the native value.
func (s TypedBinary[T, B]) MinLength() (int, bool) {
	return s.MinLen.Get()
}

// MaxLength returns the maximum permitted length of the native value.
func (s TypedBinary[T, B]) MaxLength() (int, bool) {
	return s.MaxLen.Get()
}

// ExplainLengthError returns a human-readable description of the length
// constraints, for use in an error message.
func (s TypedBinary[T, B]) ExplainLengthError() string {
	return explainLengthError(s, "(unencoded) length")
}

// EncodingDescription returns a short (one word, ideally) human-readable
// description of the encoding scheme.
func (s TypedBinary[T, B]) EncodingDescription() string {
	return s.EncodingDesc
}

// Type returns the type of the native value.
func (s TypedBinary[T, B]) Type() reflect.Type {
	return reflectx.TypeOf[T]()
}

// Finalize prepares the schema for use.
//
// It returns an error if schema is invalid.
func (s TypedBinary[T, B]) Finalize() error {
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
func (s TypedBinary[T, B]) AcceptVisitor(v SchemaVisitor) {
	v.VisitBinary(s)
}

// Marshal converts a value to its literal representation.
func (s TypedBinary[T, B]) Marshal(v T) (Literal, error) {
	if err := s.validate(v); err != nil {
		return Literal{}, err
	}

	return s.Marshaler.Marshal(v)
}

// Unmarshal converts a literal value to it's native representation.
func (s TypedBinary[T, B]) Unmarshal(v Literal) (T, error) {
	n, err := s.Marshaler.Unmarshal(v)
	if err != nil {
		return nil, err
	}

	return n, s.validate(n)
}

// Examples returns a (possibly empty) set of examples of valid values.
func (s TypedBinary[T, B]) Examples(conservative bool) []TypedExample[T] {
	if conservative {
		return nil
	}

	size, ok := s.MinLen.Get()
	if !ok {
		size = 16

		if max, ok := s.MaxLen.Get(); ok {
			if size > max {
				size = max
			}
		}
	}

	example := make(T, size)

	// Use a *deterministic* random value to as the example.
	src := rand.NewSource(int64(size))
	for i := range example {
		example[i] = B(src.Int63())
	}

	return []TypedExample[T]{
		{
			Native: example,
		},
	}
}

// validate returns an error if v is invalid.
func (s TypedBinary[T, B]) validate(v T) error {
	if min, ok := s.MinLen.Get(); ok && len(v) < min {
		return MinLengthError{s}
	}

	if max, ok := s.MaxLen.Get(); ok && len(v) > max {
		return MaxLengthError{s}
	}

	return nil
}
