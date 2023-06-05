package ferrite

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/dogmatiq/ferrite/variable"
)

// Binary configures an environment variable as a raw binary value, represented
// as a byte-slice.
//
// Binary values are represented in environment variables using a suitable
// encoding schema. The default encoding is the standard "base64" encoding
// described by RFC 4648.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Binary(name, desc string) *BinaryBuilder[[]byte, byte] {
	return BinaryAs[[]byte](name, desc)
}

// BinaryAs configures an environment variable as a raw binary value,
// represented as a user-defined byte-slice type.
//
// Binary values are represented in environment variables using a suitable
// encoding schema. The default encoding is the standard "base64" encoding
// described by RFC 4648.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func BinaryAs[T ~[]B, B ~byte](name, desc string) *BinaryBuilder[T, B] {
	b := &BinaryBuilder[T, B]{
		schema: variable.TypedBinary[T, B]{},
	}

	b.builder.Name(name)
	b.builder.Description(desc)
	b.WithBase64Encoding(base64.StdEncoding)

	return b
}

// BinaryBuilder builds a specification for a binary variable.
type BinaryBuilder[T ~[]B, B ~byte] struct {
	schema  variable.TypedBinary[T, B]
	builder variable.TypedSpecBuilder[T]
}

var _ isBuilderOf[[]byte, *BinaryBuilder[[]byte, byte]]

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *BinaryBuilder[T, B]) WithDefault(v T) *BinaryBuilder[T, B] {
	b.builder.Default(v)
	return b
}

// WithConstraint adds a constraint to the variable.
//
// fn is called with the environment variable value after it is parsed. If fn
// returns false the value is considered invalid.
func (b *BinaryBuilder[T, B]) WithConstraint(
	desc string,
	fn func(T) bool,
) *BinaryBuilder[T, B] {
	b.builder.UserConstraint(desc, fn)
	return b
}

// WithSensitiveContent marks the variable as containing sensitive content.
//
// Values of sensitive variables are not printed to the console or included in
// generated documentation.
func (b *BinaryBuilder[T, B]) WithSensitiveContent() *BinaryBuilder[T, B] {
	b.builder.MarkSensitive()
	return b
}

// WithBase64Encoding configures the variable to use base64 encoding to
// represent the binary value within the environment.
func (b *BinaryBuilder[T, B]) WithBase64Encoding(enc *base64.Encoding) *BinaryBuilder[T, B] {
	b.schema.Marshaler = base64BinaryMarshaler[T, B]{Encoding: enc}
	return b
}

// WithHexEncoding configures the variable to use hexadecimal encoding to
// represent the binary value within the environment.
func (b *BinaryBuilder[T, B]) WithHexEncoding() *BinaryBuilder[T, B] {
	b.schema.Marshaler = hexBinaryMarshaler[T, B]{}
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *BinaryBuilder[T, B]) Required(options ...RequiredOption) Required[T] {
	return required(b.schema, &b.builder, options...)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *BinaryBuilder[T, B]) Optional(options ...OptionalOption) Optional[T] {
	return optional(b.schema, &b.builder, options...)
}

// Deprecated completes the build process and registers a deprecated variable
// with Ferrite's validation system.
func (b *BinaryBuilder[T, B]) Deprecated(options ...DeprecatedOption) Deprecated[T] {
	return deprecated(b.schema, &b.builder, options...)
}

type base64BinaryMarshaler[T ~[]B, B ~byte] struct {
	Encoding *base64.Encoding
}

func (m base64BinaryMarshaler[T, B]) Marshal(v T) (variable.Literal, error) {
	return variable.Literal{
		String: m.Encoding.EncodeToString(toByteSlice(v)),
	}, nil
}

func (m base64BinaryMarshaler[T, B]) Unmarshal(v variable.Literal) (T, error) {
	data, err := m.Encoding.DecodeString(v.String)
	return fromByteSlice[T](data), err
}

func (m base64BinaryMarshaler[T, B]) EncodingDescription() string {
	switch *m.Encoding {
	case *base64.StdEncoding:
		return "base64"
	case *base64.RawStdEncoding:
		return "base64-nopad"
	case *base64.URLEncoding:
		return "base64url"
	case *base64.RawURLEncoding:
		return "base64url-nopad"
	default:
		name := "base64-noncanonical"
		if m.Encoding.EncodedLen(1) != base64.StdEncoding.EncodedLen(1) {
			name += "-nopad"
		}
		return name
	}
}

func (m base64BinaryMarshaler[T, B]) EncodedLen(n int) (min, max int) {
	n = m.Encoding.EncodedLen(n)
	return n, n
}

type hexBinaryMarshaler[T ~[]B, B ~byte] struct{}

func (m hexBinaryMarshaler[T, B]) EncodingDescription() string {
	return "hex"
}

func (m hexBinaryMarshaler[T, B]) Marshal(v T) (variable.Literal, error) {
	return variable.Literal{
		String: hex.EncodeToString(toByteSlice(v)),
	}, nil
}

func (m hexBinaryMarshaler[T, B]) Unmarshal(v variable.Literal) (T, error) {
	data, err := hex.DecodeString(v.String)
	return fromByteSlice[T](data), err
}

func (m hexBinaryMarshaler[T, B]) EncodedLen(n int) (min, max int) {
	n = hex.EncodedLen(n)
	return n, n
}

// toByteSlice converts a slice of user-defined-byte type to []byte.
func toByteSlice[T ~[]B, B ~byte](in T) []byte {
	out := make([]byte, len(in))
	for i, o := range in {
		out[i] = byte(o)
	}
	return out
}

// fromByteSlice converts a []byte to a slice of user-defined-byte types.
func fromByteSlice[T ~[]B, B ~byte](in []byte) T {
	out := make(T, len(in))
	for i, o := range in {
		out[i] = B(o)
	}
	return out
}