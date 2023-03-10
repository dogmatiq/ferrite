package ferrite

import (
	"strconv"

	"github.com/dogmatiq/ferrite/internal/reflectx"
	"github.com/dogmatiq/ferrite/maybe"
	"github.com/dogmatiq/ferrite/variable"
	"golang.org/x/exp/constraints"
)

// Float configures an environment variable as a floating-point number.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Float[T constraints.Float](name, desc string) *FloatBuilder[T] {
	b := &FloatBuilder[T]{
		schema: variable.TypedNumeric[T]{
			Marshaler: floatMarshaler[T]{},
		},
	}

	b.spec.Name(name)
	b.spec.Description(desc)
	b.spec.Documentation().
		Summary("Floating-point syntax").
		Paragraph(
			"Floating-point values can be specified using decimal (base-10) or hexadecimal (base-16) notation, and may use scientific notation.",
			"A leading positive sign (`+`) is **OPTIONAL**.",
			"A leading negative sign (`-`) is **REQUIRED** in order to specify a negative value.",
		).
		Format().
		Paragraph(
			"Internally, the `%s` variable is represented using a %d-bit floating point type (`%s`);",
			"any value that overflows this data-type is invalid.",
			"Values are rounded to the nearest floating-point number using IEEE 754 unbiased rounding.",
		).
		Format(
			name,
			reflectx.BitSize[T](),
			reflectx.KindOf[T](),
		).
		Done()

	return b
}

// FloatBuilder builds a specification for a floating-point number.
type FloatBuilder[T constraints.Float] struct {
	schema variable.TypedNumeric[T]
	spec   variable.TypedSpecBuilder[T]
}

var _ isBuilderOf[float32, *FloatBuilder[float32]]

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *FloatBuilder[T]) WithDefault(v T) *FloatBuilder[T] {
	b.spec.Default(v)
	return b
}

// WithMinimum sets the minimum acceptable value of the variable.
func (b *FloatBuilder[T]) WithMinimum(v T) *FloatBuilder[T] {
	b.schema.NativeMin = maybe.Some(v)
	return b
}

// WithMaximum sets the maximum acceptable value of the variable.
func (b *FloatBuilder[T]) WithMaximum(v T) *FloatBuilder[T] {
	b.schema.NativeMax = maybe.Some(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *FloatBuilder[T]) Required(options ...VariableOption) Required[T] {
	return req(b.schema, &b.spec, options)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *FloatBuilder[T]) Optional(options ...VariableOption) Optional[T] {
	return opt(b.schema, &b.spec, options)
}

// Deprecated completes the build process and registers a deprecated variable
// with Ferrite's validation system.
func (b *FloatBuilder[T]) Deprecated(reason string, options ...VariableOption) Deprecated[T] {
	return dep(b.schema, &b.spec, reason, options)
}

type floatMarshaler[T constraints.Float] struct{}

func (floatMarshaler[T]) Marshal(v T) (variable.Literal, error) {
	s := strconv.FormatFloat(
		float64(v),
		'g',
		-1,
		reflectx.BitSize[T](),
	)

	switch s[0] {
	case '+', '-':
	default:
		s = "+" + s
	}

	return variable.Literal{
		String: s,
	}, nil
}

func (floatMarshaler[T]) Unmarshal(v variable.Literal) (T, error) {
	n, err := strconv.ParseFloat(v.String, reflectx.BitSize[T]())
	if err != nil {
		return 0, err
	}

	return T(n), nil
}
