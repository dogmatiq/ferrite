package ferrite

import (
	"errors"
	"strconv"

	"github.com/dogmatiq/ferrite/variable"
)

// NetworkPort configures an environment variable as a network port.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func NetworkPort(name, desc string) *NetworkPortBuilder {
	b := &NetworkPortBuilder{}

	b.spec.Name(name)
	b.spec.Description(desc)
	b.spec.BuiltInConstraint(
		"**MUST** be a valid network port",
		func(v string) variable.ConstraintError {
			return validatePort(v)
		},
	)
	b.spec.NonNormativeExample("8000", "a port commonly used for private web servers")
	b.spec.NonNormativeExample("https", "the IANA service name that maps to port 443")
	buildNetworkPortSyntaxDocumentation(b.spec.Documentation())

	return b
}

// NetworkPortBuilder builds a specification for a network port variable.
type NetworkPortBuilder struct {
	schema variable.TypedString[string]
	spec   variable.TypedSpecBuilder[string]
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *NetworkPortBuilder) WithDefault(v string) *NetworkPortBuilder {
	b.spec.Default(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *NetworkPortBuilder) Required(options ...VariableOption) Required[string] {
	return req(b.schema, &b.spec, options)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *NetworkPortBuilder) Optional(options ...VariableOption) Optional[string] {
	return opt(b.schema, &b.spec, options)
}

// validateHost returns an error of port is not a valid numeric port or IANA
// service name.
func validatePort(port string) error {
	if port == "" {
		return errors.New("port must not be empty")
	}

	n, err := strconv.ParseUint(port, 10, 16)

	if errors.Is(err, strconv.ErrSyntax) {
		return validateIANAServiceName(port)
	}

	if err != nil || n == 0 {
		return errors.New("numeric ports must be between 1 and 65535")
	}

	return nil
}

// validateIANAServiceName returns an error if name is not a valid IANA service
// name.
//
// See https://www.rfc-editor.org/rfc/rfc6335.html#section-5.1.
func validateIANAServiceName(name string) error {
	n := len(name)

	// RFC-6335: MUST be at least 1 character and no more than 15 characters
	// long.
	if n == 0 || n > 15 {
		return errors.New("IANA service name must be between 1 and 15 characters")
	}

	// RFC-6335: MUST NOT begin or end with a hyphen.
	if name[0] == '-' || name[n-1] == '-' {
		return errors.New("IANA service name must not begin or end with a hyphen")
	}

	hasLetter := false

	for i := range name {
		ch := name[i] // iterate by byte (not rune)

		// RFC-6335: MUST contain only US-ASCII letters 'A' - 'Z' and 'a' - 'z',
		// digits '0' - '9', and hyphens ('-', ASCII 0x2D or decimal 45).
		switch {
		case ch >= 'A' && ch <= 'Z':
			hasLetter = true
		case ch >= 'a' && ch <= 'z':
			hasLetter = true
		case ch >= '0' && ch <= '9':
			// digit ok!
		case ch == '-':
			// RFC-6335: hyphens MUST NOT be adjacent to other hyphens.
			if name[i-1] == '-' {
				return errors.New("IANA service name must not contain adjacent hyphens")
			}
		default:
			return errors.New("IANA service name must contain only ASCII letters, digits and hyphen")
		}
	}

	// RFC-6335: MUST contain at least one letter ('A' - 'Z' or 'a' - 'z').
	if !hasLetter {
		return errors.New("IANA service name must contain at least one letter")
	}

	return nil
}

func buildNetworkPortSyntaxDocumentation(d variable.DocumentationBuilder) {
	d.
		Summary("Network port syntax").
		Paragraph(
			"Ports may be specified as a numeric value no greater than `65535`.",
			"Alternatively, a service name can be used.",
			"Service names are resolved against the system's service database,",
			"typically located in the `/etc/service` file on UNIX-like systems.",
			"Standard service names are published by IANA.",
		).
		Format().
		Done()
}
