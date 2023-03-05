package ferrite

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/dogmatiq/ferrite/maybe"
	"github.com/dogmatiq/ferrite/variable"
	"golang.org/x/exp/constraints"
)

// Unsigned configures an environment variable as a unsigned integer.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Unsigned[T constraints.Unsigned](name, desc string) UnsignedBuilder[T] {
	return UnsignedBuilder[T]{
		name: name,
		desc: desc,
	}
}

// UnsignedBuilder builds a specification for an unsigned integer value.
type UnsignedBuilder[T constraints.Unsigned] struct {
	name, desc    string
	def, min, max maybe.Value[T]
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b UnsignedBuilder[T]) WithDefault(v T) UnsignedBuilder[T] {
	b.def = maybe.Some(v)
	return b
}

// WithMinimum sets the minimum acceptable value of the variable.
func (b UnsignedBuilder[T]) WithMinimum(v T) UnsignedBuilder[T] {
	b.min = maybe.Some(v)
	return b
}

// WithMaximum sets the maximum acceptable value of the variable.
func (b UnsignedBuilder[T]) WithMaximum(v T) UnsignedBuilder[T] {
	b.max = maybe.Some(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b UnsignedBuilder[T]) Required(options ...variable.RegisterOption) Required[T] {
	return registerRequired(b.spec(true), options)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b UnsignedBuilder[T]) Optional(options ...variable.RegisterOption) Optional[T] {
	return registerOptional(b.spec(false), options)
}

func (b UnsignedBuilder[T]) spec(req bool) variable.TypedSpec[T] {
	s, err := variable.NewSpec(
		b.name,
		b.desc,
		b.def,
		req,
		variable.TypedNumeric[T]{
			Marshaler: unsignedMarshaler[T]{},
			NativeMin: b.min,
			NativeMax: b.max,
		},
		variable.WithDocumentation[T]().
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
				b.name,
				bitSize[T](),
				reflect.TypeOf(T(0)).Kind(),
			).
			Done(),
	)
	if err != nil {
		panic(err.Error())
	}

	return s
}

type unsignedMarshaler[T constraints.Unsigned] struct{}

func (unsignedMarshaler[T]) Marshal(v T) (variable.Literal, error) {
	return variable.Literal{
		String: fmt.Sprintf("%d", v),
	}, nil
}

func (unsignedMarshaler[T]) Unmarshal(v variable.Literal) (T, error) {
	n, err := strconv.ParseUint(v.String, 10, bitSize[T]())
	if err != nil {
		return 0, err
	}

	return T(n), nil
}
