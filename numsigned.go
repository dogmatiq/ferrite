package ferrite

import (
	"fmt"
	"strconv"

	"github.com/dogmatiq/ferrite/maybe"
	"github.com/dogmatiq/ferrite/variable"
	"golang.org/x/exp/constraints"
)

// Signed configures an environment variable as a signed integer.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Signed[T constraints.Signed](name, desc string) SignedBuilder[T] {
	return SignedBuilder[T]{
		name: name,
		desc: desc,
	}
}

// SignedBuilder builds a specification for a signed integer value.
type SignedBuilder[T constraints.Signed] struct {
	name, desc string
	def        maybe.Value[T]
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b SignedBuilder[T]) WithDefault(v T) SignedBuilder[T] {
	b.def = maybe.Some(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b SignedBuilder[T]) Required(options ...variable.RegisterOption) Required[T] {
	return registerRequired(b.spec(true), options)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b SignedBuilder[T]) Optional(options ...variable.RegisterOption) Optional[T] {
	return registerOptional(b.spec(false), options)
}

func (b SignedBuilder[T]) spec(req bool) variable.TypedSpec[T] {
	s, err := variable.NewSpec(
		b.name,
		b.desc,
		b.def,
		req,
		variable.TypedNumeric[T]{
			Marshaler: signedMarshaler[T]{},
		},
	)
	if err != nil {
		panic(err.Error())
	}

	return s
}

type signedMarshaler[T constraints.Signed] struct{}

func (signedMarshaler[T]) Marshal(v T) (variable.Literal, error) {
	return variable.Literal{
		String: fmt.Sprintf("%+d", v),
	}, nil
}

func (signedMarshaler[T]) Unmarshal(v variable.Literal) (T, error) {
	n, err := strconv.ParseInt(v.String, 10, bitSize[T]())
	if err != nil {
		return 0, err
	}

	return T(n), nil
}
