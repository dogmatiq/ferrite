package ferrite

import (
	"fmt"

	"github.com/dogmatiq/ferrite/maybe"
	"github.com/dogmatiq/ferrite/variable"
)

// Bool configures an environment variable as a boolean.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Bool(name, desc string) BoolBuilder[bool] {
	return BoolAs[bool](name, desc)
}

// BoolAs configures an environment variable as a boolean using a user-defined
// type.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func BoolAs[T ~bool](name, desc string) BoolBuilder[T] {
	return BoolBuilder[T]{
		name: name,
		desc: desc,
		t:    fmt.Sprint(T(true)),
		f:    fmt.Sprint(T(false)),
	}
}

// BoolBuilder builds a specification for a boolean value.
type BoolBuilder[T ~bool] struct {
	name, desc string
	t, f       string
	def        maybe.Value[T]
}

// WithLiterals overrides the default literals used to represent true and false.
//
// The default literals "true" and "false" are no longer valid values when using
// custom literals.
func (b BoolBuilder[T]) WithLiterals(t, f string) BoolBuilder[T] {
	b.t = t
	b.f = f
	return b
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b BoolBuilder[T]) WithDefault(v T) BoolBuilder[T] {
	b.def = maybe.Some(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b BoolBuilder[T]) Required(options ...variable.RegisterOption) Required[T] {
	return registerRequired(b.spec(true), options)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b BoolBuilder[T]) Optional(options ...variable.RegisterOption) Optional[T] {
	return registerOptional(b.spec(false), options)
}

func (b BoolBuilder[T]) spec(req bool) variable.TypedSpec[T] {
	s, err := variable.NewSpec(
		b.name,
		b.desc,
		b.def,
		req,
		variable.TypedSet[T]{
			Members: []T{true, false},
			ToLiteral: func(v T) variable.Literal {
				s := b.f
				if v {
					s = b.t
				}

				return variable.Literal{
					String: s,
				}
			},
		},
		variable.WithExample(T(true), ""),
		variable.WithExample(T(false), ""),
	)
	if err != nil {
		panic(err.Error())
	}

	return s
}
