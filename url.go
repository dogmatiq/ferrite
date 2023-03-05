package ferrite

import (
	"errors"
	"net/url"

	"github.com/dogmatiq/ferrite/maybe"
	"github.com/dogmatiq/ferrite/variable"
)

// URL configures an environment variable as a URL.
//
// The URL must be fully-qualified (i.e. it must have a scheme and hostname).
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func URL(name, desc string) URLBuilder {
	return URLBuilder{
		name: name,
		desc: desc,
		options: []variable.SpecOption[*url.URL]{
			variable.WithConstraint(
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
			),
			variable.WithNonNormativeExample(
				mustParseURL("https://example.org/path"),
				"a typical URL for a web page",
			),
			variable.WithDocumentation[*url.URL](
				variable.Documentation{
					Summary: "URL syntax",
					Paragraphs: []string{
						"A fully-qualified URL includes both a scheme (protocol) and a hostname. " +
							"URLs are not necessarily web addresses; `https://example.org` and " +
							"`mailto:contact@example.org` are both examples of fully-qualified URLs.",
					},
				},
			),
		},
	}
}

// URLBuilder builds a specification for a URL variable.
type URLBuilder struct {
	name, desc string
	def        maybe.Value[*url.URL]
	options    []variable.SpecOption[*url.URL]
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b URLBuilder) WithDefault(v string) URLBuilder {
	u := mustParseURL(v)
	b.def = maybe.Some(u)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b URLBuilder) Required(options ...variable.RegisterOption) Required[*url.URL] {
	return registerRequired(b.spec(true), options)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b URLBuilder) Optional(options ...variable.RegisterOption) Optional[*url.URL] {
	return registerOptional(b.spec(false), options)
}

func (b URLBuilder) spec(req bool) variable.TypedSpec[*url.URL] {
	s, err := variable.NewSpec(
		b.name,
		b.desc,
		b.def,
		req,
		variable.TypedOther[*url.URL]{
			Marshaler: urlMarshaler{},
		},
		b.options...,
	)
	if err != nil {
		panic(err.Error())
	}

	return s
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
