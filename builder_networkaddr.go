package ferrite

import (
	"net"

	"github.com/dogmatiq/ferrite/internal/variable"
)

// NetworkAddress is a network address in host:port form.
type NetworkAddress struct {
	Host string
	Port string
}

// String returns the network address in host:port form.
func (a NetworkAddress) String() string {
	return net.JoinHostPort(a.Host, a.Port)
}

// NetworkAddr configures an environment variable as a network address in
// host:port form, as accepted by net.SplitHostPort.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func NetworkAddr(name, desc string) *NetworkAddrBuilder {
	b := &NetworkAddrBuilder{
		schema: variable.TypedOther[NetworkAddress]{
			Marshaler: networkAddrMarshaler{},
		},
	}

	b.builder.Name(name)
	b.builder.Description(desc)
	b.builder.BuiltInConstraint(
		"**MUST** be a valid network address",
		func(_ variable.ConstraintContext, v NetworkAddress) variable.ConstraintError {
			return nil
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

// NetworkAddrBuilder builds a specification for a network address variable.
type NetworkAddrBuilder struct {
	schema  variable.TypedOther[NetworkAddress]
	builder variable.TypedSpecBuilder[NetworkAddress]
}

var _ isBuilderOf[
	NetworkAddress,
	string,
	*NetworkAddrBuilder,
]

// WithDefault sets the default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *NetworkAddrBuilder) WithDefault(v string) *NetworkAddrBuilder {
	b.builder.Default(mustParseNetworkAddr(v))
	return b
}

// WithExample adds an example value to the variable's documentation.
func (b *NetworkAddrBuilder) WithExample(v string, desc string) *NetworkAddrBuilder {
	b.builder.NormativeExample(mustParseNetworkAddr(v), desc)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *NetworkAddrBuilder) Required(options ...RequiredOption) Required[NetworkAddress] {
	return required(b.schema, &b.builder, options...)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *NetworkAddrBuilder) Optional(options ...OptionalOption) Optional[NetworkAddress] {
	return optional(b.schema, &b.builder, options...)
}

// Deprecated completes the build process and registers a deprecated variable
// with Ferrite's validation system.
func (b *NetworkAddrBuilder) Deprecated(options ...DeprecatedOption) Deprecated[NetworkAddress] {
	return deprecated(b.schema, &b.builder, options...)
}

type networkAddrMarshaler struct{}

func (networkAddrMarshaler) Marshal(v NetworkAddress) (variable.Literal, error) {
	return variable.Literal{
		String: v.String(),
	}, nil
}

func (networkAddrMarshaler) Unmarshal(v variable.Literal) (NetworkAddress, error) {
	host, port, err := net.SplitHostPort(v.String)
	if err != nil {
		return NetworkAddress{}, err
	}
	return NetworkAddress{Host: host, Port: port}, nil
}

func mustParseNetworkAddr(v string) NetworkAddress {
	host, port, err := net.SplitHostPort(v)
	if err != nil {
		panic(err)
	}
	return NetworkAddress{Host: host, Port: port}
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
