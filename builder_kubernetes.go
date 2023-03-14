package ferrite

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"unicode"

	"github.com/dogmatiq/ferrite/variable"
)

// KubernetesAddress is the address of a Kubernetes service.
type KubernetesAddress struct {
	Host string
	Port string
}

func (a KubernetesAddress) String() string {
	return net.JoinHostPort(a.Host, a.Port)
}

// KubernetesService configures environment variables used to obtain the network
// address of a specific Kubernetes service.
//
// svc is the name of the Kubernetes service, NOT the environment variable.
//
// The environment variables "<svc>_SERVICE_HOST" and "<svc>_SERVICE_PORT" are
// used to construct an address for the service.
func KubernetesService(svc string) *KubernetesServiceBuilder {
	if err := validateKubernetesName(svc); err != nil {
		panic(fmt.Sprintf(
			"kubernetes service name is invalid: %s",
			err,
		))
	}

	b := &KubernetesServiceBuilder{
		service: svc,
	}

	const docs = "It is expected that this variable will be implicitly defined by Kubernetes; " +
		"it typically does not need to be specified in the pod manifest."

	b.hostBuilder.Name(
		fmt.Sprintf(
			"%s_SERVICE_HOST",
			kubernetesNameToEnv(svc),
		),
	)
	b.hostBuilder.Description(
		fmt.Sprintf(
			"kubernetes %q service host",
			b.service,
		),
	)
	b.hostBuilder.BuiltInConstraint(
		"**MUST** be a valid hostname",
		func(h string) variable.ConstraintError {
			return validateHost(h)
		},
	)
	b.hostBuilder.Documentation().
		Paragraph(docs).
		Format().
		Important().
		Done()

	b.portBuilder.Name(
		fmt.Sprintf(
			"%s_SERVICE_PORT",
			kubernetesNameToEnv(svc),
		),
	)
	b.portBuilder.Description(
		fmt.Sprintf(
			"kubernetes %q service port",
			b.service,
		),
	)
	b.portBuilder.BuiltInConstraint(
		"**MUST** be a valid network port",
		func(p string) variable.ConstraintError {
			return validatePort(p)
		},
	)
	b.portBuilder.Documentation().
		Paragraph(docs).
		Format().
		Important().
		Done()

	buildNetworkPortSyntaxDocumentation(b.portBuilder.Documentation())

	variable.EstablishRelationships(
		variable.RefersTo{
			Subject:  b.hostBuilder.Peek(),
			RefersTo: b.portBuilder.Peek(),
		},
		variable.RefersTo{
			Subject:  b.portBuilder.Peek(),
			RefersTo: b.hostBuilder.Peek(),
		},
	)

	return b
}

// KubernetesServiceBuilder is the specification for a Kubernetes service.
type KubernetesServiceBuilder struct {
	service     string
	hostSchema  variable.TypedString[string]
	portSchema  variable.TypedString[string]
	hostBuilder variable.TypedSpecBuilder[string]
	portBuilder variable.TypedSpecBuilder[string]
}

var _ isBuilderOf[KubernetesAddress, *KubernetesServiceBuilder]

// WithNamedPort uses a Kubernetes named port instead of the default service
// port.
//
// The "<service>_SERVICE_PORT_<port>" environment variable is used instead of
// "<service>_SERVICE_PORT".
//
// The Kubernetes port name is the name configured in the service manifest. It
// is not to be confused with an IANA registered service name (e.g. "https"),
// although the two may use the same names.
//
// See https://kubernetes.io/docs/concepts/services-networking/service/#multi-port-services
func (b *KubernetesServiceBuilder) WithNamedPort(port string) *KubernetesServiceBuilder {
	if err := validateKubernetesName(port); err != nil {
		panic(fmt.Sprintf(
			"specification of kubernetes %q service is invalid: invalid named port: %s",
			b.service,
			err,
		))
	}

	b.portBuilder.Name(
		fmt.Sprintf(
			"%s_SERVICE_PORT_%s",
			kubernetesNameToEnv(b.service),
			kubernetesNameToEnv(port),
		),
	)

	return b
}

// WithDefault sets a default value to use when the environment variables are
// undefined.
//
// The port may be a numeric value between 1 and 65535, or an IANA registered
// service name (such as "https"). The IANA name is not to be confused with the
// Kubernetes servcice name or port name.
func (b *KubernetesServiceBuilder) WithDefault(host, port string) *KubernetesServiceBuilder {
	b.hostBuilder.Default(host)
	b.portBuilder.Default(port)
	return b
}

// Required completes the build process and registers required variables with
// Ferrite's validation system.
func (b *KubernetesServiceBuilder) Required(options ...RequiredOption) Required[KubernetesAddress] {
	b.hostBuilder.MarkRequired()
	b.portBuilder.MarkRequired()

	var cfg variableSetConfig
	for _, opt := range options {
		opt.applyRequiredOptionToConfig(&cfg)
		opt.applyRequiredOptionToSpec(&b.hostBuilder)
		opt.applyRequiredOptionToSpec(&b.portBuilder)
	}

	host := variable.Register(
		cfg.Registry,
		b.hostBuilder.Done(b.hostSchema),
	)

	port := variable.Register(
		cfg.Registry,
		b.portBuilder.Done(b.portSchema),
	)

	return requiredFunc[KubernetesAddress]{
		[]variable.Any{host, port},
		func() (KubernetesAddress, error) {
			if err := host.Error(); err != nil {
				return KubernetesAddress{}, err
			}

			if err := port.Error(); err != nil {
				return KubernetesAddress{}, err
			}

			return KubernetesAddress{
				host.NativeValue(),
				port.NativeValue(),
			}, nil
		},
	}
}

// Optional completes the build process and registers optional variables with
// Ferrite's validation system.
func (b *KubernetesServiceBuilder) Optional(options ...OptionalOption) Optional[KubernetesAddress] {
	var cfg variableSetConfig
	for _, opt := range options {
		opt.applyOptionalOptionToConfig(&cfg)
		opt.applyOptionalOptionToSpec(&b.hostBuilder)
		opt.applyOptionalOptionToSpec(&b.portBuilder)
	}

	host := variable.Register(
		cfg.Registry,
		b.hostBuilder.Done(b.hostSchema),
	)

	port := variable.Register(
		cfg.Registry,
		b.portBuilder.Done(b.portSchema),
	)

	return optionalFunc[KubernetesAddress]{
		[]variable.Any{host, port},
		b.optionalResolver(host, port),
	}
}

// Deprecated completes the build process and registers deprecated variables
// with Ferrite's validation system.
func (b *KubernetesServiceBuilder) Deprecated(options ...DeprecatedOption) Deprecated[KubernetesAddress] {
	b.hostBuilder.MarkDeprecated()
	b.portBuilder.MarkDeprecated()

	var cfg variableSetConfig
	for _, opt := range options {
		opt.applyDeprecatedOptionToConfig(&cfg)
		opt.applyDeprecatedOptionToSpec(&b.hostBuilder)
		opt.applyDeprecatedOptionToSpec(&b.portBuilder)
	}

	host := variable.Register(
		cfg.Registry,
		b.hostBuilder.Done(b.hostSchema),
	)

	port := variable.Register(
		cfg.Registry,
		b.portBuilder.Done(b.portSchema),
	)

	return deprecatedFunc[KubernetesAddress]{
		[]variable.Any{host, port},
		b.optionalResolver(host, port),
	}
}

func (b *KubernetesServiceBuilder) optionalResolver(
	host, port *variable.OfType[string],
) func() (KubernetesAddress, bool, error) {
	return func() (KubernetesAddress, bool, error) {
		if err := host.Error(); err != nil {
			return KubernetesAddress{}, false, err
		}
		if err := port.Error(); err != nil {
			return KubernetesAddress{}, false, err
		}

		availability := host.Availability()

		if port.Availability() != availability {
			def, undef := host, port
			if availability != variable.AvailabilityOK {
				def, undef = undef, def
			}

			return KubernetesAddress{}, false, fmt.Errorf(
				"%s is defined but %s is not, define both or neither",
				def.Spec().Name(),
				undef.Spec().Name(),
			)
		}

		return KubernetesAddress{
			host.NativeValue(),
			port.NativeValue(),
		}, availability == variable.AvailabilityOK, nil
	}
}

// kubernetesNameToEnv converts a kubernetes resource name to an environment variable
// name, as per Kubernetes own behavior.
func kubernetesNameToEnv(s string) string {
	return strings.ToUpper(
		strings.ReplaceAll(s, "-", "_"),
	)
}

// validateKubernetesName returns an error if name is not a valid Kubernetes
// resource name.
func validateKubernetesName(name string) error {
	if name == "" {
		return errors.New("name must not be empty")
	}

	n := len(name)

	if name[0] == '-' || name[n-1] == '-' {
		return errors.New("name must not begin or end with a hyphen")
	}

	for i := range name {
		ch := name[i] // iterate by byte (not rune)

		switch {
		case ch >= 'a' && ch <= 'z':
		case ch >= '0' && ch <= '9':
		case ch == '-':
		default:
			return errors.New("name must contain only lowercase ASCII letters, digits and hyphen")
		}
	}

	return nil
}

// validateHost returns an error of host is not a valid hostname.
func validateHost(host string) error {
	if host == "" {
		return errors.New("host must not be empty")
	}

	if net.ParseIP(host) != nil {
		return nil
	}

	n := len(host)
	if host[0] == '.' || host[n-1] == '.' {
		return errors.New("host must not begin or end with a dot")
	}

	for _, r := range host {
		if unicode.IsSpace(r) {
			return errors.New("host must not contain whitespace")
		}
	}

	return nil
}
