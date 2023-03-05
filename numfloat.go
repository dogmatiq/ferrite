package ferrite

import (
	"strconv"

	"github.com/dogmatiq/ferrite/maybe"
	"github.com/dogmatiq/ferrite/variable"
	"golang.org/x/exp/constraints"
)

// Float configures an environment variable as a floating-point number.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Float[T constraints.Float](name, desc string) FloatBuilder[T] {
	return FloatBuilder[T]{
		name: name,
		desc: desc,
	}
}

// FloatBuilder builds a specification for a floating-point number.
type FloatBuilder[T constraints.Float] struct {
	name, desc    string
	def, min, max maybe.Value[T]
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b FloatBuilder[T]) WithDefault(v T) FloatBuilder[T] {
	b.def = maybe.Some(v)
	return b
}

// WithMinimum sets the minimum acceptable value of the variable.
func (b FloatBuilder[T]) WithMinimum(v T) FloatBuilder[T] {
	b.min = maybe.Some(v)
	return b
}

// WithMaximum sets the maximum acceptable value of the variable.
func (b FloatBuilder[T]) WithMaximum(v T) FloatBuilder[T] {
	b.max = maybe.Some(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b FloatBuilder[T]) Required(options ...variable.RegisterOption) Required[T] {
	return registerRequired(b.spec(true), options)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b FloatBuilder[T]) Optional(options ...variable.RegisterOption) Optional[T] {
	return registerOptional(b.spec(false), options)
}

func (b FloatBuilder[T]) spec(req bool) variable.TypedSpec[T] {
	s, err := variable.NewSpec(
		b.name,
		b.desc,
		b.def,
		req,
		variable.TypedNumeric[T]{
			Marshaler: floatMarshaler[T]{},
			NativeMin: b.min,
			NativeMax: b.max,
		},
	)
	if err != nil {
		panic(err.Error())
	}

	return s
}

type floatMarshaler[T constraints.Float] struct{}

func (floatMarshaler[T]) Marshal(v T) (variable.Literal, error) {
	s := strconv.FormatFloat(
		float64(v),
		'g',
		-1,
		bitSize[T](),
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
	n, err := strconv.ParseFloat(v.String, bitSize[T]())
	if err != nil {
		return 0, err
	}

	return T(n), nil
}
