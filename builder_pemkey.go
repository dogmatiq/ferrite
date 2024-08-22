package ferrite

import (
	"crypto"
	"crypto/rsa"

	"github.com/dogmatiq/ferrite/internal/variable"
)

// PrivateKey configures an environment variable as a filename referring to a
// file that contains a PEM-encoded private key.
//
// That key is parsed to produce a value of type K.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func PrivateKey[K crypto.PrivateKey](
	name,
	desc string,
) *PrivateKeyBuilder[K] {
	b := &PrivateKeyBuilder[K]{
		schema: variable.TypedOther[PrivateKeyFile[K]]{
			Marshaler: privateKeyMarshaler[K]{},
		},
	}
	b.builder.Name(name)
	b.builder.Description(desc)
	b.builder.NonNormativeExample(PrivateKeyFile[K]{Name: "/path/to/privatekey.pem"}, "an absolute file path")
	b.builder.NonNormativeExample(PrivateKeyFile[K]{Name: "./path/to/privatekey.pem"}, "a relative file path")
	return b
}

// PrivateKeyBuilder builds a specification for a string value that is the name
// of a PEM file containing a private key of type T.
type PrivateKeyBuilder[K crypto.PrivateKey] struct {
	schema  variable.TypedOther[PrivateKeyFile[K]]
	builder variable.TypedSpecBuilder[PrivateKeyFile[K]]
}

var _ isBuilderOf[PrivateKeyFile[*rsa.PrivateKey], *PrivateKeyBuilder[*rsa.PrivateKey]]

// WithDefault sets the default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *PrivateKeyBuilder[K]) WithDefault(v string) *PrivateKeyBuilder[K] {
	b.builder.Default(
		PrivateKeyFile[K]{Name: v},
	)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *PrivateKeyBuilder[T]) Required(options ...RequiredOption) Required[PrivateKeyFile[T]] {
	return required(b.schema, &b.builder, options...)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *PrivateKeyBuilder[T]) Optional(options ...OptionalOption) Optional[PrivateKeyFile[T]] {
	return optional(b.schema, &b.builder, options...)
}

// Deprecated completes the build process and registers a deprecated variable
// with Ferrite's validation system.
func (b *PrivateKeyBuilder[T]) Deprecated(options ...DeprecatedOption) Deprecated[PrivateKeyFile[T]] {
	return deprecated(b.schema, &b.builder, options...)
}

type PrivateKeyFile[K crypto.PrivateKey] struct {
	Name string
	Key  K
}

type privateKeyMarshaler[K crypto.PrivateKey] struct{}

func (privateKeyMarshaler[K]) Marshal(PrivateKeyFile[K]) (variable.Literal, error) {
	panic("not implemented")
}

func (privateKeyMarshaler[K]) Unmarshal(variable.Literal) (PrivateKeyFile[K], error) {
	panic("not implemented")
}
