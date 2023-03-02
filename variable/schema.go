package variable

import (
	"reflect"
)

// Schema describes the valid values of an environment value.
type Schema interface {
	// Type returns the type of the native value.
	Type() reflect.Type

	// Finalize prepares the schema for use.
	//
	// It returns an error if schema is invalid.
	Finalize() error

	// AcceptVisitor passes the schema to the appropriate method of v.
	AcceptVisitor(SchemaVisitor)
}

// SchemaError indicates that a value is invalid because it violates its schema.
type SchemaError interface {
	error

	// Schema returns the schema that was violated.
	Schema() Schema

	// AcceptVisitor passes the error to the appropriate method of v.
	AcceptVisitor(SchemaErrorVisitor)
}

// SchemaVisitor dispatches based on a variable's schema.
type SchemaVisitor interface {
	VisitNumeric(Numeric)
	VisitSet(Set)
	VisitString(String)
	VisitOther(Other)
}

// SchemaErrorVisitor dispatches based on the type of a SchemaError.
type SchemaErrorVisitor interface {
	// Numeric errors ...
	VisitMinError(MinError)
	VisitMaxError(MaxError)

	// Set errors...
	VisitSetMembershipError(SetMembershipError)

	// String errors ...
	VisitMinLengthError(MinLengthError)
	VisitMaxLengthError(MaxLengthError)
}

// TypedSchema describes the valid values of an environment varible value
// depicted by type T.
type TypedSchema[T any] interface {
	Schema
	Marshaler[T]

	// Examples returns a (possibly empty) set of examples of valid values.
	Examples(hasOtherExamples bool) []TypedExample[T]
}

// typeOf returns the type of T.
func typeOf[T any]() reflect.Type {
	return reflect.TypeOf([...]T{}).Elem()
}
