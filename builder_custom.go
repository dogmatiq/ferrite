package ferrite

import (
	"github.com/dogmatiq/ferrite/internal/variable"
)

// Custom configures an environment variable as a value of type T, using
// caller-supplied functions to marshal and unmarshal the value.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
//
// unmarshal is called to convert the environment variable's string value into T.
// If it returns an error, the value is considered invalid.
//
// marshal is called to convert a T value back to its string representation.
// It is used to render defaults and examples in validation output.
func Custom[T any](
	name, desc string,
	unmarshal func(string) (T, error),
	marshal func(T) (string, error),
) *CustomBuilder[T] {
	if unmarshal == nil {
		panic("Custom: unmarshal function must not be nil")
	}
	if marshal == nil {
		panic("Custom: marshal function must not be nil")
	}

	b := &CustomBuilder[T]{
		schema: variable.TypedOther[T]{
			Marshaler: customMarshaler[T]{unmarshal, marshal},
		},
	}

	b.builder.Name(name)
	b.builder.Description(desc)

	return b
}

// CustomBuilder builds a specification for a variable of an arbitrary type,
// using caller-supplied functions to marshal and unmarshal the value.
type CustomBuilder[T any] struct {
	schema  variable.TypedOther[T]
	builder variable.TypedSpecBuilder[T]
}

var _ isBuilderOf[
	any,
	any,
	*CustomBuilder[any],
]

// WithDefault sets the default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *CustomBuilder[T]) WithDefault(v T) *CustomBuilder[T] {
	b.builder.Default(v)
	return b
}

// WithExample adds an example value to the variable's documentation.
func (b *CustomBuilder[T]) WithExample(v T, desc string) *CustomBuilder[T] {
	b.builder.NormativeExample(v, desc)
	return b
}

// WithConstraint adds a constraint to the variable.
//
// fn is called with the environment variable value after it is parsed. If fn
// returns false the value is considered invalid.
func (b *CustomBuilder[T]) WithConstraint(
	desc string,
	fn func(T) bool,
) *CustomBuilder[T] {
	b.builder.UserConstraint(desc, fn)
	return b
}

// WithSensitiveContent marks the variable as containing sensitive content.
//
// Values of sensitive variables are not printed to the console or included in
// generated documentation.
func (b *CustomBuilder[T]) WithSensitiveContent() *CustomBuilder[T] {
	b.builder.MarkSensitive()
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *CustomBuilder[T]) Required(options ...RequiredOption) Required[T] {
	return required(b.schema, &b.builder, options...)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *CustomBuilder[T]) Optional(options ...OptionalOption) Optional[T] {
	return optional(b.schema, &b.builder, options...)
}

// Deprecated completes the build process and registers a deprecated variable
// with Ferrite's validation system.
func (b *CustomBuilder[T]) Deprecated(options ...DeprecatedOption) Deprecated[T] {
	return deprecated(b.schema, &b.builder, options...)
}

type customMarshaler[T any] struct {
	unmarshal func(string) (T, error)
	marshal   func(T) (string, error)
}

func (m customMarshaler[T]) Marshal(v T) (variable.Literal, error) {
	s, err := m.marshal(v)
	if err != nil {
		return variable.Literal{}, err
	}
	return variable.Literal{String: s}, nil
}

func (m customMarshaler[T]) Unmarshal(v variable.Literal) (T, error) {
	n, err := m.unmarshal(v.String)
	if err != nil {
		return n, err
	}

	// Verify that the unmarshaled value can be round-tripped through the
	// marshal function. [variable.TypedSpec] re-marshals every valid value to
	// compute its canonical literal representation, and panics if marshaling
	// fails. Since these are user-supplied functions, it's plausible that
	// unmarshal accepts values that marshal can't represent, so we check here
	// to surface that as a validation error rather than a panic.
	if _, err := m.marshal(n); err != nil {
		return n, err
	}

	return n, nil
}
