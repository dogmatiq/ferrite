package ferrite

import (
	"github.com/dogmatiq/ferrite/internal/maybe"
	"github.com/dogmatiq/ferrite/internal/variable"
)

// String configures an environment variable as a string.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func String(name, desc string) *StringBuilder[string] {
	return StringAs[string](name, desc)
}

// StringAs configures an environment variable as a string using a user-defined
// type.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func StringAs[T ~string](name, desc string) *StringBuilder[T] {
	b := &StringBuilder[T]{}

	b.builder.Name(name)
	b.builder.Description(desc)

	return b
}

// StringBuilder builds a specification for a string variable.
type StringBuilder[T ~string] struct {
	schema  variable.TypedString[T]
	builder variable.TypedSpecBuilder[T]
}

var _ isBuilderOf[
	string,
	string,
	*StringBuilder[string],
]

// WithDefault sets the default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *StringBuilder[T]) WithDefault(v T) *StringBuilder[T] {
	b.builder.Default(variable.ConstDefault(v))
	return b
}

// WithExample adds an example value to the variable's documentation.
func (b *StringBuilder[T]) WithExample(v T, desc string) *StringBuilder[T] {
	b.builder.NormativeExample(v, desc)
	return b
}

// WithMinimumLength sets the minimum permitted length of the variable, in
// bytes (not runes).
func (b *StringBuilder[T]) WithMinimumLength(min int) *StringBuilder[T] {
	b.schema.MinLen = maybe.Some(min)
	return b
}

// WithMaximumLength sets the maximum permitted length of the variable, in
// bytes (not runes).
func (b *StringBuilder[T]) WithMaximumLength(max int) *StringBuilder[T] {
	b.schema.MaxLen = maybe.Some(max)
	return b
}

// WithLength sets the exact permitted length of the variable, in bytes (not
// runes).
func (b *StringBuilder[T]) WithLength(n int) *StringBuilder[T] {
	l := maybe.Some(n)
	b.schema.MinLen = l
	b.schema.MaxLen = l
	return b
}

// WithConstraint adds a constraint to the variable.
//
// fn is called with the environment variable value after it is parsed. If fn
// returns false the value is considered invalid.
func (b *StringBuilder[T]) WithConstraint(
	desc string,
	fn func(T) bool,
) *StringBuilder[T] {
	b.builder.UserConstraint(desc, fn)
	return b
}

// WithSensitiveContent marks the variable as containing sensitive content.
//
// Values of sensitive variables are not printed to the console or included in
// generated documentation.
func (b *StringBuilder[T]) WithSensitiveContent() *StringBuilder[T] {
	b.builder.MarkSensitive()
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *StringBuilder[T]) Required(options ...RequiredOption) Required[T] {
	return required(b.schema, &b.builder, options...)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *StringBuilder[T]) Optional(options ...OptionalOption) Optional[T] {
	return optional(b.schema, &b.builder, options...)
}

// Deprecated completes the build process and registers a deprecated variable
// with Ferrite's validation system.
func (b *StringBuilder[T]) Deprecated(options ...DeprecatedOption) Deprecated[T] {
	return deprecated(b.schema, &b.builder, options...)
}
