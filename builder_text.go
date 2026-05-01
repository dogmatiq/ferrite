package ferrite

import (
	"encoding"

	"github.com/dogmatiq/ferrite/internal/variable"
)

// TextAs configures an environment variable as a value of type T, which
// implements [encoding.TextMarshaler] directly (via a value receiver) and
// [encoding.TextUnmarshaler] via its pointer type P.
//
// The result type of the variable is T (the value type). Use [TextAsP] if T
// only implements [encoding.TextMarshaler] via a pointer receiver, or if the
// result type should be a pointer.
func TextAs[
	T encoding.TextMarshaler,
	P interface {
		*T
		encoding.TextUnmarshaler
	},
](name, desc string) *TextAsBuilder[T] {
	b := &TextAsBuilder[T]{
		schema: variable.TypedOther[T]{
			Marshaler: textMarshaler[T]{
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

// TextAsP configures an environment variable as a value of type P (a pointer
// to E), which implements [encoding.TextMarshaler] and
// [encoding.TextUnmarshaler].
//
// The result type of the variable is P (the pointer type). Use [TextAs] if the
// result type should be a value type.
func TextAsP[
	P interface {
		*E
		encoding.TextMarshaler
		encoding.TextUnmarshaler
	},
	E any,
](name, desc string) *TextAsBuilder[P] {
	b := &TextAsBuilder[P]{
		schema: variable.TypedOther[P]{
			Marshaler: textMarshaler[P]{
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

// TextAsBuilder builds a specification for a variable that uses
// [encoding.TextMarshaler] and [encoding.TextUnmarshaler] for marshaling.
//
// The interface constraints are enforced by the [TextAs] and [TextAsP]
// constructors, not by this type. This allows both constructors to share a
// single builder type without propagating the pointer type parameter.
type TextAsBuilder[T any] struct {
	schema  variable.TypedOther[T]
	builder variable.TypedSpecBuilder[T]
}

var _ isBuilderOf[
	any,
	any,
	*TextAsBuilder[any],
]

// WithDefault sets the default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *TextAsBuilder[T]) WithDefault(v T) *TextAsBuilder[T] {
	b.builder.Default(v)
	return b
}

// WithExample adds an example value to the variable's documentation.
func (b *TextAsBuilder[T]) WithExample(v T, desc string) *TextAsBuilder[T] {
	b.builder.NormativeExample(v, desc)
	return b
}

// WithConstraint adds a constraint to the variable.
//
// fn is called with the environment variable value after it is parsed. If fn
// returns false the value is considered invalid.
func (b *TextAsBuilder[T]) WithConstraint(
	desc string,
	fn func(T) bool,
) *TextAsBuilder[T] {
	b.builder.UserConstraint(desc, fn)
	return b
}

// WithSensitiveContent marks the variable as containing sensitive content.
//
// Values of sensitive variables are not printed to the console or included in
// generated documentation.
func (b *TextAsBuilder[T]) WithSensitiveContent() *TextAsBuilder[T] {
	b.builder.MarkSensitive()
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *TextAsBuilder[T]) Required(options ...RequiredOption) Required[T] {
	return required(b.schema, &b.builder, options...)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *TextAsBuilder[T]) Optional(options ...OptionalOption) Optional[T] {
	return optional(b.schema, &b.builder, options...)
}

// Deprecated completes the build process and registers a deprecated variable
// with Ferrite's validation system.
func (b *TextAsBuilder[T]) Deprecated(options ...DeprecatedOption) Deprecated[T] {
	return deprecated(b.schema, &b.builder, options...)
}

// textMarshaler implements [variable.Marshaler] using closures that invoke the
// [encoding.TextMarshaler] and [encoding.TextUnmarshaler] interfaces.
type textMarshaler[T any] struct {
	marshal   func(T) ([]byte, error)
	unmarshal func([]byte) (T, error)
}

func (m textMarshaler[T]) Marshal(v T) (variable.Literal, error) {
	data, err := m.marshal(v)
	if err != nil {
		return variable.Literal{}, err
	}
	return variable.Literal{String: string(data)}, nil
}

func (m textMarshaler[T]) Unmarshal(v variable.Literal) (T, error) {
	return m.unmarshal([]byte(v.String))
}
