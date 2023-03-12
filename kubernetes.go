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

// KubernetesService reads Kubernetes service discovery environment variables for a
// service's port.
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

	b.hostSpec.Name(
		fmt.Sprintf(
			"%s_SERVICE_HOST",
			kubernetesNameToEnv(svc),
		),
	)
	b.hostSpec.Description(
		fmt.Sprintf(
			"kubernetes %q service host",
			b.service,
		),
	)
	b.hostSpec.BuiltInConstraint(
		"**MUST** be a valid hostname",
		func(h string) variable.ConstraintError {
			return validateHost(h)
		},
	)
	b.hostSpec.Documentation().
		Paragraph(docs).
		Format().
		Important().
		Done()

	b.portSpec.Name(
		fmt.Sprintf(
			"%s_SERVICE_PORT",
			kubernetesNameToEnv(svc),
		),
	)
	b.portSpec.Description(
		fmt.Sprintf(
			"kubernetes %q service port",
			b.service,
		),
	)
	b.portSpec.BuiltInConstraint(
		"**MUST** be a valid network port",
		func(p string) variable.ConstraintError {
			return validatePort(p)
		},
	)
	b.portSpec.Documentation().
		Paragraph(docs).
		Format().
		Important().
		Done()

	buildNetworkPortSyntaxDocumentation(b.portSpec.Documentation())

	seeAlso(b.hostSpec.Peek(), b.portSpec.Peek())
	seeAlso(b.portSpec.Peek(), b.hostSpec.Peek())

	return b
}

// KubernetesServiceBuilder is the specification for a Kubernetes service.
type KubernetesServiceBuilder struct {
	service    string
	hostSchema variable.TypedString[string]
	portSchema variable.TypedString[string]
	hostSpec   variable.TypedSpecBuilder[string]
	portSpec   variable.TypedSpecBuilder[string]
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

	b.portSpec.Name(
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
	b.hostSpec.Default(host)
	b.portSpec.Default(port)
	return b
}

// SeeAlso creates a relationship between this variable and those used by i.
func (b *KubernetesServiceBuilder) SeeAlso(i Input, options ...SeeAlsoOption) *KubernetesServiceBuilder {
	seeAlsoInput(&b.hostSpec, i, options...)
	seeAlsoInput(&b.portSpec, i, options...)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *KubernetesServiceBuilder) Required(options ...RequiredOption) Required[KubernetesAddress] {
	b.hostSpec.MarkRequired()
	b.portSpec.MarkRequired()

	hostConfig := buildRequiredConfig(&b.hostSpec, options...)
	portConfig := buildRequiredConfig(&b.portSpec, options...)

	host := variable.Register(
		hostConfig.Registry,
		b.hostSpec.Done(b.hostSchema),
	)

	port := variable.Register(
		portConfig.Registry,
		b.portSpec.Done(b.portSchema),
	)

	return requiredFunc[KubernetesAddress]{
		[]variable.Any{host, port},
		func() (KubernetesAddress, error) {
			h, ok, err := host.NativeValue()
			if err != nil {
				return KubernetesAddress{}, err
			}

			if !ok {
				return KubernetesAddress{}, undefinedError(host)
			}

			p, ok, err := port.NativeValue()
			if err != nil {
				return KubernetesAddress{}, err
			}

			if !ok {
				return KubernetesAddress{}, undefinedError(port)
			}

			return KubernetesAddress{h, p}, nil
		},
	}
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *KubernetesServiceBuilder) Optional(options ...OptionalOption) Optional[KubernetesAddress] {
	hostConfig := buildOptionalConfig(&b.hostSpec, options...)
	portConfig := buildOptionalConfig(&b.portSpec, options...)

	host := variable.Register(
		hostConfig.Registry,
		b.hostSpec.Done(b.hostSchema),
	)

	port := variable.Register(
		portConfig.Registry,
		b.portSpec.Done(b.portSchema),
	)

	return optionalFunc[KubernetesAddress]{
		[]variable.Any{host, port},
		b.optionalResolver(host, port),
	}
}

// Deprecated completes the build process and registers a deprecated variable
// with Ferrite's validation system.
func (b *KubernetesServiceBuilder) Deprecated(options ...DeprecatedOption) Deprecated[KubernetesAddress] {
	b.hostSpec.MarkDeprecated()
	b.portSpec.MarkDeprecated()

	hostConfig := buildDeprecatedConfig(&b.hostSpec, options...)
	portConfig := buildDeprecatedConfig(&b.portSpec, options...)

	host := variable.Register(
		hostConfig.Registry,
		b.hostSpec.Done(b.hostSchema),
	)

	port := variable.Register(
		portConfig.Registry,
		b.portSpec.Done(b.portSchema),
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
		h, hostOk, err := host.NativeValue()
		if err != nil {
			return KubernetesAddress{}, false, err
		}

		p, portOk, err := port.NativeValue()
		if err != nil {
			return KubernetesAddress{}, false, err
		}

		if hostOk && portOk {
			return KubernetesAddress{h, p}, true, nil
		}

		if hostOk {
			return KubernetesAddress{}, false, fmt.Errorf(
				"%s is defined but %s is not, define both or neither",
				host.Spec().Name(),
				port.Spec().Name(),
			)
		}

		if portOk {
			return KubernetesAddress{}, false, fmt.Errorf(
				"%s is defined but %s is not, define both or neither",
				port.Spec().Name(),
				host.Spec().Name(),
			)
		}

		return KubernetesAddress{}, false, nil
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
