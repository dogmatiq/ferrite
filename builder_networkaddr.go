package ferrite

import (
	"errors"
	"net"

	"github.com/dogmatiq/ferrite/internal/variable"
)

// NetworkAddr is a network address in host:port form.
type NetworkAddr struct {
	Host string
	Port string
}

// String returns the network address in host:port form.
func (a NetworkAddr) String() string {
	return net.JoinHostPort(a.Host, a.Port)
}

// NetworkAddress configures an environment variable as a network address in
// host:port form, as accepted by net.SplitHostPort.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func NetworkAddress(name, desc string) *NetworkAddressBuilder {
	b := &NetworkAddressBuilder{
		schema: variable.TypedOther[NetworkAddr]{
			Marshaler: networkAddrMarshaler{},
		},
	}

	b.builder.Name(name)
	b.builder.Description(desc)
	b.builder.BuiltInConstraint(
		"**MUST** be a valid network address",
		func(_ variable.ConstraintContext, v NetworkAddr) variable.ConstraintError {
			if v.Host == "" {
				return errors.New("host must not be empty")
			}
			return validatePort(v.Port)
		},
	)
	b.builder.NonNormativeExample(
		mustParseNetworkAddr("192.168.0.1:8080"),
		"an IPv4 address with a port",
	)
	b.builder.NonNormativeExample(
		mustParseNetworkAddr("[::1]:8080"),
		"an IPv6 address with a port",
	)
	b.builder.NonNormativeExample(
		mustParseNetworkAddr("host.example.org:https"),
		"a named host with an IANA service name",
	)
	buildNetworkAddrSyntaxDocumentation(b.builder.Documentation())

	return b
}

// NetworkAddressBuilder builds a specification for a network address variable.
type NetworkAddressBuilder struct {
	schema  variable.TypedOther[NetworkAddr]
	builder variable.TypedSpecBuilder[NetworkAddr]
}

var _ isBuilderOf[
	NetworkAddr,
	string,
	*NetworkAddressBuilder,
]

// WithDefault sets the default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *NetworkAddressBuilder) WithDefault(v string) *NetworkAddressBuilder {
	b.builder.Default(mustParseNetworkAddr(v))
	return b
}

// WithExample adds an example value to the variable's documentation.
func (b *NetworkAddressBuilder) WithExample(v string, desc string) *NetworkAddressBuilder {
	b.builder.NormativeExample(mustParseNetworkAddr(v), desc)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *NetworkAddressBuilder) Required(options ...RequiredOption) Required[NetworkAddr] {
	return required(b.schema, &b.builder, options...)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *NetworkAddressBuilder) Optional(options ...OptionalOption) Optional[NetworkAddr] {
	return optional(b.schema, &b.builder, options...)
}

// Deprecated completes the build process and registers a deprecated variable
// with Ferrite's validation system.
func (b *NetworkAddressBuilder) Deprecated(options ...DeprecatedOption) Deprecated[NetworkAddr] {
	return deprecated(b.schema, &b.builder, options...)
}

type networkAddrMarshaler struct{}

func (networkAddrMarshaler) Marshal(v NetworkAddr) (variable.Literal, error) {
	return variable.Literal{
		String: v.String(),
	}, nil
}

func (networkAddrMarshaler) Unmarshal(v variable.Literal) (NetworkAddr, error) {
	host, port, err := net.SplitHostPort(v.String)
	if err != nil {
		return NetworkAddr{}, err
	}
	return NetworkAddr{Host: host, Port: port}, nil
}

func mustParseNetworkAddr(v string) NetworkAddr {
	host, port, err := net.SplitHostPort(v)
	if err != nil {
		panic(err)
	}
	return NetworkAddr{Host: host, Port: port}
}

func buildNetworkAddrSyntaxDocumentation(d variable.DocumentationBuilder) {
	d.
		Summary("Network address syntax").
		Paragraph(
			"Addresses may be specified as `<host>:<port>`, where `<host>` is a",
			"hostname or IP address and `<port>` is a numeric port number or an",
			"IANA service name.",
			"IPv6 addresses must be enclosed in square brackets, e.g. `[::1]:8080`.",
		).
		Format().
		Done()
}
