package ferrite

import (
	"errors"
	"net/url"

	"github.com/dogmatiq/ferrite/variable"
)

// URL configures an environment variable as a URL.
//
// The URL must be fully-qualified (i.e. it must have a scheme and hostname).
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func URL(name, desc string) *URLBuilder {
	b := &URLBuilder{
		schema: variable.TypedOther[*url.URL]{
			Marshaler: urlMarshaler{},
		},
	}

	b.spec.Name(name)
	b.spec.Description(desc)
	b.spec.BuiltInConstraint(
		"**MUST** be a fully-qualified URL",
		func(v *url.URL) variable.ConstraintError {
			if v.Scheme == "" {
				return errors.New("URL must have a scheme")
			}

			if v.Host == "" {
				return errors.New("URL must have a hostname")
			}

			return nil
		},
	)
	b.spec.NonNormativeExample(
		mustParseURL("https://example.org/path"),
		"a typical URL for a web page",
	)
	b.spec.Documentation().
		Summary("URL syntax").
		Paragraph(
			"A fully-qualified URL includes both a scheme (protocol) and a hostname.",
			"URLs are not necessarily web addresses;",
			"`https://example.org` and `mailto:contact@example.org` are both examples of fully-qualified URLs.",
		).
		Format().
		Done()

	return b
}

// URLBuilder builds a specification for a URL variable.
type URLBuilder struct {
	schema variable.TypedOther[*url.URL]
	spec   variable.TypedSpecBuilder[*url.URL]
}

var _ isBuilderOf[*url.URL, *URLBuilder]

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *URLBuilder) WithDefault(v string) *URLBuilder {
	b.spec.Default(mustParseURL(v))
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *URLBuilder) Required(options ...RequiredOption) Required[*url.URL] {
	return required(b.schema, &b.spec, options)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *URLBuilder) Optional(options ...OptionalOption) Optional[*url.URL] {
	return optional(b.schema, &b.spec, options)
}

// Deprecated completes the build process and registers a deprecated variable
// with Ferrite's validation system.
func (b *URLBuilder) Deprecated(reason string, options ...DeprecatedOption) Deprecated[*url.URL] {
	return deprecated(b.schema, &b.spec, reason, options)
}

type urlMarshaler struct{}

func (urlMarshaler) Marshal(v *url.URL) (variable.Literal, error) {
	return variable.Literal{
		String: v.String(),
	}, nil
}

func (urlMarshaler) Unmarshal(v variable.Literal) (*url.URL, error) {
	return url.Parse(v.String)
}

func mustParseURL(v string) *url.URL {
	u, err := url.Parse(v)
	if err != nil {
		panic(err)
	}
	return u
}
