package ferrite

import (
	"encoding"

	"github.com/dogmatiq/ferrite/internal/variable"
)

// TextEncoded configures an environment variable as a value of type T, which
// implements [encoding.TextMarshaler] directly (via a value receiver) and
// [encoding.TextUnmarshaler] via its pointer type P.
//
// The result type of the variable is T (the value type). Use [TextEncodedP] if
// T only implements [encoding.TextMarshaler] via a pointer receiver, or if the
// result type should be a pointer.
func TextEncoded[
	T encoding.TextMarshaler,
	P interface {
		*T
		encoding.TextUnmarshaler
	},
](name, desc string) *TextEncodedBuilder[T] {
	b := &TextEncodedBuilder[T]{
		schema: variable.TypedOther[T]{
			Marshaler: textEncodedMarshaler[T]{
				marshal: T.MarshalText,
				unmarshal: func(data []byte) (T, error) {
					var v T
					if err := P(&v).UnmarshalText(data); err != nil {
						return v, err
					}
					return v, nil
				},
			},
		},
	}

	b.builder.Name(name)
	b.builder.Description(desc)

	return b
}

// TextEncodedP configures an environment variable as a value of type P (a
// pointer to E), which implements [encoding.TextMarshaler] and
// [encoding.TextUnmarshaler].
//
// The result type of the variable is P (the pointer type). Use [TextEncoded] if
// the result type should be a value type.
func TextEncodedP[
	P interface {
		*E
		encoding.TextMarshaler
		encoding.TextUnmarshaler
	},
	E any,
](name, desc string) *TextEncodedBuilder[P] {
	b := &TextEncodedBuilder[P]{
		schema: variable.TypedOther[P]{
			Marshaler: textEncodedMarshaler[P]{
				marshal: P.MarshalText,
				unmarshal: func(data []byte) (P, error) {
					v := P(new(E))
					if err := v.UnmarshalText(data); err != nil {
						return nil, err
					}
					return v, nil
				},
			},
		},
	}

	b.builder.Name(name)
	b.builder.Description(desc)

	return b
}

// TextEncodedBuilder builds a specification for a variable that uses
// [encoding.TextMarshaler] and [encoding.TextUnmarshaler] for marshaling.
//
// The interface constraints are enforced by the [TextEncoded] and
// [TextEncodedP] constructors, not by this type. This allows both constructors
// to share a single builder type without propagating the pointer type
// parameter.
type TextEncodedBuilder[T any] struct {
	schema  variable.TypedOther[T]
	builder variable.TypedSpecBuilder[T]
}

var _ isBuilderOf[
	any,
	any,
	*TextEncodedBuilder[any],
]

// WithDefault sets the default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *TextEncodedBuilder[T]) WithDefault(v T) *TextEncodedBuilder[T] {
	b.builder.Default(v)
	return b
}

// WithExample adds an example value to the variable's documentation.
func (b *TextEncodedBuilder[T]) WithExample(v T, desc string) *TextEncodedBuilder[T] {
	b.builder.NormativeExample(v, desc)
	return b
}

// WithConstraint adds a constraint to the variable.
//
// fn is called with the environment variable value after it is parsed. If fn
// returns false the value is considered invalid.
func (b *TextEncodedBuilder[T]) WithConstraint(
	desc string,
	fn func(T) bool,
) *TextEncodedBuilder[T] {
	b.builder.UserConstraint(desc, fn)
	return b
}

// WithSensitiveContent marks the variable as containing sensitive content.
//
// Values of sensitive variables are not printed to the console or included in
// generated documentation.
func (b *TextEncodedBuilder[T]) WithSensitiveContent() *TextEncodedBuilder[T] {
	b.builder.MarkSensitive()
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *TextEncodedBuilder[T]) Required(options ...RequiredOption) Required[T] {
	return required(b.schema, &b.builder, options...)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *TextEncodedBuilder[T]) Optional(options ...OptionalOption) Optional[T] {
	return optional(b.schema, &b.builder, options...)
}

// Deprecated completes the build process and registers a deprecated variable
// with Ferrite's validation system.
func (b *TextEncodedBuilder[T]) Deprecated(options ...DeprecatedOption) Deprecated[T] {
	return deprecated(b.schema, &b.builder, options...)
}

// textEncodedMarshaler implements [variable.Marshaler] using closures that
// invoke the [encoding.TextMarshaler] and [encoding.TextUnmarshaler]
// interfaces.
type textEncodedMarshaler[T any] struct {
	marshal   func(T) ([]byte, error)
	unmarshal func([]byte) (T, error)
}

func (m textEncodedMarshaler[T]) Marshal(v T) (variable.Literal, error) {
	data, err := m.marshal(v)
	if err != nil {
		return variable.Literal{}, err
	}
	return variable.Literal{String: string(data)}, nil
}

func (m textEncodedMarshaler[T]) Unmarshal(v variable.Literal) (T, error) {
	return m.unmarshal([]byte(v.String))
}
