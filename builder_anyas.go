package ferrite

import (
	"github.com/dogmatiq/ferrite/internal/variable"
)

// AnyAs configures an environment variable as a value of type T, using
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
func AnyAs[T any](
	name, desc string,
	unmarshal func(string) (T, error),
	marshal func(T) (string, error),
) *AnyAsBuilder[T] {
	if unmarshal == nil {
		panic("AnyAs: unmarshal function must not be nil")
	}
	if marshal == nil {
		panic("AnyAs: marshal function must not be nil")
	}

	b := &AnyAsBuilder[T]{
		schema: variable.TypedOther[T]{
			Marshaler: anyAsMarshaler[T]{unmarshal, marshal},
		},
	}

	b.builder.Name(name)
	b.builder.Description(desc)

	return b
}

// AnyAsBuilder builds a specification for a variable of an arbitrary type,
// using caller-supplied functions to marshal and unmarshal the value.
type AnyAsBuilder[T any] struct {
	schema  variable.TypedOther[T]
	builder variable.TypedSpecBuilder[T]
}

var _ isBuilderOf[
	any,
	any,
	*AnyAsBuilder[any],
]

// WithDefault sets the default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *AnyAsBuilder[T]) WithDefault(v T) *AnyAsBuilder[T] {
	b.builder.Default(v)
	return b
}

// WithExample adds an example value to the variable's documentation.
func (b *AnyAsBuilder[T]) WithExample(v T, desc string) *AnyAsBuilder[T] {
	b.builder.NormativeExample(v, desc)
	return b
}

// WithConstraint adds a constraint to the variable.
//
// fn is called with the environment variable value after it is parsed. If fn
// returns false the value is considered invalid.
func (b *AnyAsBuilder[T]) WithConstraint(
	desc string,
	fn func(T) bool,
) *AnyAsBuilder[T] {
	b.builder.UserConstraint(desc, fn)
	return b
}

// WithSensitiveContent marks the variable as containing sensitive content.
//
// Values of sensitive variables are not printed to the console or included in
// generated documentation.
func (b *AnyAsBuilder[T]) WithSensitiveContent() *AnyAsBuilder[T] {
	b.builder.MarkSensitive()
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *AnyAsBuilder[T]) Required(options ...RequiredOption) Required[T] {
	return required(b.schema, &b.builder, options...)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *AnyAsBuilder[T]) Optional(options ...OptionalOption) Optional[T] {
	return optional(b.schema, &b.builder, options...)
}

// Deprecated completes the build process and registers a deprecated variable
// with Ferrite's validation system.
func (b *AnyAsBuilder[T]) Deprecated(options ...DeprecatedOption) Deprecated[T] {
	return deprecated(b.schema, &b.builder, options...)
}

type anyAsMarshaler[T any] struct {
	unmarshal func(string) (T, error)
	marshal   func(T) (string, error)
}

func (m anyAsMarshaler[T]) Marshal(v T) (variable.Literal, error) {
	s, err := m.marshal(v)
	if err != nil {
		return variable.Literal{}, err
	}
	return variable.Literal{String: s}, nil
}

func (m anyAsMarshaler[T]) Unmarshal(v variable.Literal) (T, error) {
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
