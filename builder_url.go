package ferrite

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/dogmatiq/ferrite/internal/variable"
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

	b.builder.Name(name)
	b.builder.Description(desc)
	b.builder.BuiltInConstraint(
		"**MUST** be a fully-qualified URL",
		func(_ variable.ConstraintContext, v *url.URL) variable.ConstraintError {
			if v.Scheme == "" {
				return errors.New("URL must have a scheme")
			}

			if v.Host == "" {
				return errors.New("URL must have a hostname")
			}

			return nil
		},
	)
	b.builder.NonNormativeExample(
		mustParseURL("https://example.org/path"),
		"a typical URL for a web page",
	)
	return b
}

// URLBuilder builds a specification for a URL variable.
type URLBuilder struct {
	schema  variable.TypedOther[*url.URL]
	builder variable.TypedSpecBuilder[*url.URL]
	schemes []allowedScheme
}

type allowedScheme struct {
	Name        string
	Description string
}

var _ isBuilderOf[
	*url.URL,
	string,
	*URLBuilder,
]

// WithDefault sets the default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *URLBuilder) WithDefault(v string) *URLBuilder {
	b.builder.Default(mustParseURL(v))
	return b
}

// WithExample adds an example value to the variable's documentation.
func (b *URLBuilder) WithExample(v string, desc string) *URLBuilder {
	b.builder.NormativeExample(mustParseURL(v), desc)
	return b
}

// WithConstraint adds a constraint to the variable.
//
// fn is called with the environment variable value after it is parsed. If fn
// returns false the value is considered invalid.
func (b *URLBuilder) WithConstraint(
	desc string,
	fn func(*url.URL) bool,
) *URLBuilder {
	b.builder.UserConstraint(desc, fn)
	return b
}

// WithScheme restricts the variable to URLs that use the given scheme.
//
// Each call adds another permitted scheme, so call WithScheme() multiple times
// to list every protocol you want to allow. desc is an optional human-readable
// description of when the scheme should be used. scheme must not be empty or
// whitespace-only.
func (b *URLBuilder) WithScheme(scheme, desc string) *URLBuilder {
	s := strings.TrimSpace(scheme)
	if s == "" {
		panic("URL schemes must not be empty")
	}

	b.schemes = append(
		b.schemes,
		allowedScheme{
			Name:        strings.ToLower(s),
			Description: desc,
		},
	)

	return b
}

func (b *URLBuilder) applySchemeConstraint() {
	allowed := uniqueSchemes(b.schemes)
	if len(allowed) == 0 {
		return
	}

	allowedSet := map[string]struct{}{}
	for _, scheme := range allowed {
		allowedSet[scheme.Name] = struct{}{}
	}

	names := schemeNames(allowed)
	constraintDesc := schemeConstraintDescription(names)
	constraintErr := schemeConstraintError(names)

	b.builder.BuiltInConstraint(
		constraintDesc,
		func(_ variable.ConstraintContext, v *url.URL) variable.ConstraintError {
			scheme := strings.ToLower(v.Scheme)
			if _, ok := allowedSet[scheme]; ok {
				return nil
			}

			return errors.New(constraintErr)
		},
	)
}

func (b *URLBuilder) applyDocumentation() {
	builder := b.builder.Documentation().
		Summary("URL syntax").
		Paragraph(
			"A fully-qualified URL includes both a scheme (protocol) and a hostname.",
			"URLs are not necessarily web addresses;",
			"`https://example.org` and `mailto:contact@example.org` are both examples of fully-qualified URLs.",
		).
		Format()

	if allowed := uniqueSchemes(b.schemes); len(allowed) != 0 {
		builder = builder.
			Paragraph(allowedSchemeParagraph(allowed)).
			Format()

		for _, p := range allowedSchemeDescriptions(allowed) {
			builder = builder.Paragraph(p).Format()
		}
	}

	builder.Done()
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *URLBuilder) Required(options ...RequiredOption) Required[*url.URL] {
	b.applySchemeConstraint()
	b.applyDocumentation()
	return required(b.schema, &b.builder, options...)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *URLBuilder) Optional(options ...OptionalOption) Optional[*url.URL] {
	b.applySchemeConstraint()
	b.applyDocumentation()
	return optional(b.schema, &b.builder, options...)
}

// Deprecated completes the build process and registers a deprecated variable
// with Ferrite's validation system.
func (b *URLBuilder) Deprecated(options ...DeprecatedOption) Deprecated[*url.URL] {
	b.applySchemeConstraint()
	b.applyDocumentation()
	return deprecated(b.schema, &b.builder, options...)
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

func formatSchemeList(schemes []string) string {
	quoted := make([]string, len(schemes))
	for i, scheme := range schemes {
		quoted[i] = fmt.Sprintf("`%s`", scheme)
	}

	switch len(quoted) {
	case 1:
		return quoted[0]
	case 2:
		return quoted[0] + " or " + quoted[1]
	default:
		return strings.Join(quoted[:len(quoted)-1], ", ") + ", or " + quoted[len(quoted)-1]
	}
}

func allowedSchemeParagraph(schemes []allowedScheme) string {
	names := schemeNames(schemes)
	if len(names) == 1 {
		return fmt.Sprintf("Only the %s scheme is permitted.", formatSchemeList(names))
	}

	return fmt.Sprintf("Only the %s schemes are permitted.", formatSchemeList(names))
}

func allowedSchemeDescriptions(schemes []allowedScheme) []string {
	var descs []string
	for _, scheme := range schemes {
		desc := strings.TrimSpace(scheme.Description)
		if desc == "" {
			continue
		}

		descs = append(
			descs,
			fmt.Sprintf("- `%s` — %s", scheme.Name, desc),
		)
	}

	return descs
}

func schemeConstraintDescription(schemes []string) string {
	if len(schemes) == 1 {
		return fmt.Sprintf("**MUST** use the %s scheme", formatSchemeList(schemes))
	}

	return fmt.Sprintf("**MUST** use one of the %s schemes", formatSchemeList(schemes))
}

func schemeConstraintError(schemes []string) string {
	if len(schemes) == 1 {
		return fmt.Sprintf("expected the scheme to be %s", formatSchemeList(schemes))
	}

	return fmt.Sprintf("expected the scheme to be one of %s", formatSchemeList(schemes))
}

func uniqueSchemes(values []allowedScheme) []allowedScheme {
	seen := map[string]struct{}{}
	uniq := make([]allowedScheme, 0, len(values))
	for _, v := range values {
		if _, ok := seen[v.Name]; ok {
			continue
		}
		seen[v.Name] = struct{}{}
		uniq = append(uniq, v)
	}
	return uniq
}

func schemeNames(schemes []allowedScheme) []string {
	names := make([]string, len(schemes))
	for i, scheme := range schemes {
		names[i] = scheme.Name
	}
	return names
}
