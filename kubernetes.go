package ferrite

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"unicode"

	"github.com/dogmatiq/ferrite/maybe"
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
func KubernetesService(svc string) KubernetesServiceBuilder {
	if err := validateKubernetesName(svc); err != nil {
		panic(fmt.Sprintf(
			"kubernetes service name is invalid: %s",
			err,
		))
	}

	return KubernetesServiceBuilder{
		service: svc,
		hostVar: fmt.Sprintf(
			"%s_SERVICE_HOST",
			kubernetesNameToEnv(svc),
		),
		portVar: fmt.Sprintf(
			"%s_SERVICE_PORT",
			kubernetesNameToEnv(svc),
		),
	}
}

// KubernetesServiceBuilder is the specification for a Kubernetes service.
type KubernetesServiceBuilder struct {
	service          string
	hostVar, portVar string
	def              maybe.Value[KubernetesAddress]
}

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
func (b KubernetesServiceBuilder) WithNamedPort(port string) KubernetesServiceBuilder {
	if err := validateKubernetesName(port); err != nil {
		panic(fmt.Sprintf(
			"specification of kubernetes %q service is invalid: invalid named port: %s",
			b.service,
			err,
		))
	}

	b.portVar = fmt.Sprintf(
		"%s_SERVICE_PORT_%s",
		kubernetesNameToEnv(b.service),
		kubernetesNameToEnv(port),
	)

	return b
}

// WithDefault sets a default value to use when the environment variables are
// undefined.
//
// The port may be a numeric value between 1 and 65535, or an IANA registered
// service name (such as "https"). The IANA name is not to be confused with the
// Kubernetes servcice name or port name.
func (b KubernetesServiceBuilder) WithDefault(host, port string) KubernetesServiceBuilder {
	b.def = maybe.Some(KubernetesAddress{host, port})
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b KubernetesServiceBuilder) Required(options ...variable.RegisterOption) Required[KubernetesAddress] {
	host := variable.Register(b.hostSpec(true), options)
	port := variable.Register(b.portSpec(true), options)

	return Required[KubernetesAddress]{
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
func (b KubernetesServiceBuilder) Optional(options ...variable.RegisterOption) Optional[KubernetesAddress] {
	host := variable.Register(b.hostSpec(false), options)
	port := variable.Register(b.portSpec(false), options)

	return Optional[KubernetesAddress]{
		func() (KubernetesAddress, bool, error) {
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
		},
	}
}

func (b KubernetesServiceBuilder) hostSpec(req bool) variable.TypedSpec[string] {
	s, err := variable.NewSpec(
		b.hostVar,
		fmt.Sprintf(
			"kubernetes %q service host",
			b.service,
		),
		maybe.Map(b.def, func(addr KubernetesAddress) string {
			return addr.Host
		}),
		req,
		variable.TypedString[string]{},
		variable.WithConstraint(
			"must be a valid hostname",
			func(h string) variable.ConstraintError {
				return validateHost(h)
			},
		),
	)
	if err != nil {
		panic(err.Error())
	}

	return s
}

func (b KubernetesServiceBuilder) portSpec(req bool) variable.TypedSpec[string] {
	s, err := variable.NewSpec(
		b.portVar,
		fmt.Sprintf(
			"kubernetes %q service port",
			b.service,
		),
		maybe.Map(b.def, func(addr KubernetesAddress) string {
			return addr.Port
		}),
		req,
		variable.TypedString[string]{},
		variable.WithConstraint(
			"must be a valid port",
			func(p string) variable.ConstraintError {
				return validatePort(p)
			},
		),
	)
	if err != nil {
		panic(err.Error())
	}

	return s
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
