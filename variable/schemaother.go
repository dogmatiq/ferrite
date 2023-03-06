package variable

import (
	"io"
	"reflect"
)

// Other is a schema for representing values of arbitrary types.
//
// It should be used as a last resort when no other schema offers a better
// explanation of the value.
type Other interface {
	Schema
}

// TypedOther is a schema for representing values of arbitrary types.
//
// It should be used as a last resort when no other schema offers a better
// explanation of the value.
type TypedOther[T any] struct {
	Marshaler Marshaler[T]

	SchemaRenderer     func(io.Writer, TypedOther[T])
	ValueErrorRenderer func(io.Writer, TypedOther[T], ValueError)
}

// Type returns the type of the native value.
func (s TypedOther[T]) Type() reflect.Type {
	return typeOf[T]()
}

// Finalize prepares the schema for use.
//
// It returns an error if schema is invalid.
func (s TypedOther[T]) Finalize() error {
	return nil
}

// AcceptVisitor passes s to the appropriate method of v.
func (s TypedOther[T]) AcceptVisitor(v SchemaVisitor) {
	v.VisitOther(s)
}

// Marshal converts a value to its literal representation.
func (s TypedOther[T]) Marshal(v T) (Literal, error) {
	return s.Marshaler.Marshal(v)
}

// Unmarshal converts a literal value to it's native representation.
func (s TypedOther[T]) Unmarshal(v Literal) (T, error) {
	return s.Marshaler.Unmarshal(v)
}

// Examples returns a (possibly empty) set of examples of valid values.
func (s TypedOther[T]) Examples(hasOtherExamples bool) []TypedExample[T] {
	return nil
}
