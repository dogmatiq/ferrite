package variable

import (
	"fmt"
	"math/rand"
	"reflect"

	"github.com/dogmatiq/ferrite/internal/reflectx"
	"github.com/dogmatiq/ferrite/maybe"
)

// Binary is a schema that allows input of binary data.
type Binary interface {
	LengthLimited

	// EncodingDescription returns a short (one word, ideally) human-readable
	// description of the encoding scheme.
	EncodingDescription() string
}

// BinaryMarshaler is a Marshaler that encodes and decodes binary data.
type BinaryMarshaler[T ~[]B, B ~byte] interface {
	Marshaler[T]

	// EncodingDescription returns a short (one word, ideally) human-readable
	// description of the encoding scheme.
	EncodingDescription() string

	// EncodedLen returns the minimum and maximum length of the string that
	// encodes n bytes of binary data.
	EncodedLen(n int) (min, max int)
}

// TypedBinary is a string value depicted by type T.
type TypedBinary[T ~[]B, B ~byte] struct {
	Marshaler      BinaryMarshaler[T, B]
	MinLen, MaxLen maybe.Value[int]
}

// MinLengthLiteral returns the minimum permitted length of the literal
// environment variable value, in bytes.
func (s TypedBinary[T, B]) MinLengthLiteral() (int, bool) {
	if min, ok := s.MinLengthNative(); ok {
		min, _ = s.Marshaler.EncodedLen(min)
		return min, true
	}
	return 0, false
}

// MaxLengthLiteral returns the maximum permitted length of the literal
// environment variable value, in bytes.
func (s TypedBinary[T, B]) MaxLengthLiteral() (int, bool) {
	if max, ok := s.MaxLengthNative(); ok {
		_, max = s.Marshaler.EncodedLen(max)
		return max, true
	}
	return 0, false
}

// MinLengthNative returns the minimum permitted length of the native value.
func (s TypedBinary[T, B]) MinLengthNative() (int, bool) {
	return s.MinLen.Get()
}

// MaxLengthNative returns the maximum permitted length of the native value.
func (s TypedBinary[T, B]) MaxLengthNative() (int, bool) {
	return s.MaxLen.Get()
}

// MinLengthEncoded returns the minimum permitted length of the encoded
// data, in bytes.
func (s TypedBinary[T, B]) MinLengthEncoded() (int, bool) {
	if min, ok := s.MinLen.Get(); ok {
		min, _ = s.Marshaler.EncodedLen(min)
		return min, true
	}
	return 0, false
}

// MaxLengthEncoded returns the maximum permitted length of the encoded
// data, in bytes.
func (s TypedBinary[T, B]) MaxLengthEncoded() (int, bool) {
	if max, ok := s.MaxLen.Get(); ok {
		max, _ = s.Marshaler.EncodedLen(max)
		return max, true
	}
	return 0, false
}

// EncodingDescription returns a short (one word, ideally) human-readable
// description of the encoding scheme.
func (s TypedBinary[T, B]) EncodingDescription() string {
	return s.Marshaler.EncodingDescription()
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
func (s TypedBinary[T, B]) Examples(hasOtherExamples bool) []TypedExample[T] {
	if hasOtherExamples {
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
		return MinLengthError{s}
	}

	return nil
}
