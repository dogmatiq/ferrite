package ferrite

import (
	"fmt"
	"strconv"

	"github.com/dogmatiq/ferrite/internal/reflectx"
	"github.com/dogmatiq/ferrite/maybe"
	"github.com/dogmatiq/ferrite/variable"
	"golang.org/x/exp/constraints"
)

// Signed configures an environment variable as a signed integer.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Signed[T constraints.Signed](name, desc string) *SignedBuilder[T] {
	b := &SignedBuilder[T]{
		schema: variable.TypedNumeric[T]{
			Marshaler: signedMarshaler[T]{},
		},
	}

	b.spec.Init(name, desc)
	b.spec.
		Documentation("Signed integer syntax").
		Paragraph(
			"Signed integers can only be specified using decimal notation.",
			"A leading positive sign (`+`) is **OPTIONAL**.",
			"A leading negative sign (`-`) is **REQUIRED** in order to specify a negative value.",
		).
		Format().
		Paragraph(
			"Internally, the `%s` variable is represented using a signed %d-bit integer type (`%s`);",
			"any value that overflows this data-type is invalid.",
		).
		Format(
			name,
			reflectx.BitSize[T](),
			reflectx.KindOf[T](),
		).
		Done()

	return b
}

// SignedBuilder builds a specification for a signed integer value.
type SignedBuilder[T constraints.Signed] struct {
	schema variable.TypedNumeric[T]
	spec   variable.SpecBuilder[T]
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *SignedBuilder[T]) WithDefault(v T) *SignedBuilder[T] {
	b.spec.Default(v)
	return b
}

// WithMinimum sets the minimum acceptable value of the variable.
func (b *SignedBuilder[T]) WithMinimum(v T) *SignedBuilder[T] {
	b.schema.NativeMin = maybe.Some(v)
	return b
}

// WithMaximum sets the maximum acceptable value of the variable.
func (b *SignedBuilder[T]) WithMaximum(v T) *SignedBuilder[T] {
	b.schema.NativeMax = maybe.Some(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *SignedBuilder[T]) Required(options ...Option) Required[T] {
	b.spec.MarkRequired()
	v := b.spec.Done(b.schema, options)
	return requiredOne(v)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *SignedBuilder[T]) Optional(options ...Option) Optional[T] {
	v := b.spec.Done(b.schema, options)
	return optionalOne(v)

}

type signedMarshaler[T constraints.Signed] struct{}

func (signedMarshaler[T]) Marshal(v T) (variable.Literal, error) {
	return variable.Literal{
		String: fmt.Sprintf("%+d", v),
	}, nil
}

func (signedMarshaler[T]) Unmarshal(v variable.Literal) (T, error) {
	n, err := strconv.ParseInt(v.String, 10, reflectx.BitSize[T]())
	if err != nil {
		return 0, err
	}

	return T(n), nil
}
