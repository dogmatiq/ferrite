package ferrite

import (
	"fmt"
	"strconv"

	"github.com/dogmatiq/ferrite/internal/maybe"
	"github.com/dogmatiq/ferrite/internal/reflectx"
	"github.com/dogmatiq/ferrite/internal/variable"
	"golang.org/x/exp/constraints"
)

// Unsigned configures an environment variable as a unsigned integer.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Unsigned[T constraints.Unsigned](name, desc string) *UnsignedBuilder[T] {
	b := &UnsignedBuilder[T]{
		schema: variable.TypedNumeric[T]{
			Marshaler: unsignedMarshaler[T]{},
		},
	}

	b.builder.Name(name)
	b.builder.Description(desc)
	b.builder.Documentation().
		Summary("Unsigned integer syntax").
		Paragraph(
			"Unsigned integers can only be specified using decimal (base-10) notation.",
			"A leading sign (`+` or `-`) is not supported and **MUST NOT** be specified.",
		).
		Format().
		Paragraph(
			"Internally, the `%s` variable is represented using an unsigned %d-bit integer type (`%s`);",
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

// UnsignedBuilder builds a specification for an unsigned integer value.
type UnsignedBuilder[T constraints.Unsigned] struct {
	schema  variable.TypedNumeric[T]
	builder variable.TypedSpecBuilder[T]
}

var _ isBuilderOf[uint, *UnsignedBuilder[uint]]

// WithDefault sets the default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *UnsignedBuilder[T]) WithDefault(v T) *UnsignedBuilder[T] {
	b.builder.Default(v)
	return b
}

// WithMinimum sets the minimum acceptable value of the variable.
func (b *UnsignedBuilder[T]) WithMinimum(v T) *UnsignedBuilder[T] {
	b.schema.NativeMin = maybe.Some(v)
	return b
}

// WithMaximum sets the maximum acceptable value of the variable.
func (b *UnsignedBuilder[T]) WithMaximum(v T) *UnsignedBuilder[T] {
	b.schema.NativeMax = maybe.Some(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *UnsignedBuilder[T]) Required(options ...RequiredOption) Required[T] {
	return required(b.schema, &b.builder, options...)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *UnsignedBuilder[T]) Optional(options ...OptionalOption) Optional[T] {
	return optional(b.schema, &b.builder, options...)
}

// Deprecated completes the build process and registers a deprecated variable
// with Ferrite's validation system.
func (b *UnsignedBuilder[T]) Deprecated(options ...DeprecatedOption) Deprecated[T] {
	return deprecated(b.schema, &b.builder, options...)
}

type unsignedMarshaler[T constraints.Unsigned] struct{}

func (unsignedMarshaler[T]) Marshal(v T) (variable.Literal, error) {
	return variable.Literal{
		String: formatUnsigned(v),
	}, nil
}

func (unsignedMarshaler[T]) Unmarshal(v variable.Literal) (T, error) {
	n, err := strconv.ParseUint(v.String, 10, reflectx.BitSize[T]())
	return T(n), variable.UnwrapNumericParseError(err, formatUnsigned[T])
}

func formatUnsigned[T constraints.Unsigned](v T) string {
	return fmt.Sprintf("%d", v)
}
