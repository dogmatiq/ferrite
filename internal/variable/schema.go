package variable

import (
	"reflect"
)

// Schema describes the valid values of an environment variable value.
type Schema interface {
	// Type returns the type of the native value.
	Type() reflect.Type

	// Finalize prepares the schema for use.
	//
	// It returns an error if schema is invalid.
	Finalize() error

	// AcceptVisitor passes the schema to the appropriate method of v.
	AcceptVisitor(v SchemaVisitor)
}

// LengthLimited is an interface for a [Schema] that impose a minimum
// and/or maximum length on an environment variable value.
type LengthLimited interface {
	Schema

	// MinLength returns the minimum permitted length of the native value.
	MinLength() (int, bool)

	// MaxLength returns the maximum permitted length of the native value.
	MaxLength() (int, bool)

	// LengthDescription returns a human-readable description of the length that
	// the limit applies to.
	//
	// It must produce a gramatical sentence of the form:
	//  "The value must have <desc> of exactly 5 bytes.
	LengthDescription() string
}

// SchemaError indicates that a value is invalid because it violates its schema.
type SchemaError interface {
	error

	// Schema returns the schema that was violated.
	Schema() Schema

	// AcceptVisitor passes the error to the appropriate method of v.
	AcceptVisitor(v SchemaErrorVisitor)
}

// SchemaVisitor dispatches based on a variable's schema.
type SchemaVisitor interface {
	VisitBinary(Binary)
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
	//
	// If conservative is true, the schema should only return examples that
	// are fairly likely to be meaningful to the application.
	//
	// If conservative is false, the schema should return as many examples
	// as possible, even if they are not very likely to be meaningful.
	Examples(conservative bool) []TypedExample[T]
}
