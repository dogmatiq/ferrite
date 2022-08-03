package ferrite

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// KubeService reads Kubernetes service discovery environment variables for a
// service's port.
//
// svc is the name of the Kubernetes service, NOT the environment variable.
//
// The environment variables "<svc>_SERVICE_HOST" and "<svc>_SERVICE_PORT" are
// used to construct an address for the service.
func KubeService(svc string) *KubeServiceSpec {
	s := &KubeServiceSpec{
		service: svc,
		hostResult: ValidationResult{
			Name:        fmt.Sprintf("%s_SERVICE_HOST", kubeToEnv(svc)),
			Description: fmt.Sprintf(`Hostname or IP address of the "%s" service.`, svc),
			ValidInput:  "[string]",
		},
		portResult: ValidationResult{
			Name:        fmt.Sprintf("%s_SERVICE_PORT", kubeToEnv(svc)),
			Description: fmt.Sprintf(`Network port of the "%s" service.`, svc),
			ValidInput:  "[string]|(1..65535)",
		},
	}

	Register(s)

	return s
}

// KubeServiceSpec is the specification for a Kubernetes service.
type KubeServiceSpec struct {
	service string

	seal       seal
	host, port string
	hostResult ValidationResult
	portResult ValidationResult
}

// WithNamedPort uses a Kubernetes named port instead of the default service
// port.
//
//
// The "<service>_SERVICE_PORT_<port>" environment variable is used instead of
// "<service>_SERVICE_PORT".
//
// The Kubernetes port name is the name configured in the service manifest. It
// is not to be confused with an IANA registered service name (e.g. "https"),
// although the two may use the same names.
//
// See https://kubernetes.io/docs/concepts/services-networking/service/#multi-port-services
func (s *KubeServiceSpec) WithNamedPort(port string) *KubeServiceSpec {
	// TODO: panic if port is not As with Kubernetes names in general, names for
	// ports must only contain lowercase alphanumeric characters and -. Port
	// names must also start and end with an alphanumeric character.

	return s.with(func() {
		s.portResult.Name = fmt.Sprintf(
			"%s_SERVICE_PORT_%s",
			kubeToEnv(s.service),
			kubeToEnv(port),
		)

		s.portResult.Description = fmt.Sprintf(
			`Network port of the "%s" service's "%s" port.`,
			s.service,
			port,
		)
	})
}

// WithDefault sets a default value to use when the environment variables are
// undefined.
//
// The port may be a numeric value between 1 and 65535, or an IANA registered
// service name (such as "https"). The IANA name is not to be confused with the
// Kubernetes servcice name or port name.
func (s *KubeServiceSpec) WithDefault(host, port string) *KubeServiceSpec {
	if err := validateHost(host); err != nil {
		panic(fmt.Sprintf(
			"default value of %s is invalid: %s",
			s.hostResult.Name,
			err,
		))
	}

	if err := validatePort(port); err != nil {
		panic(fmt.Sprintf(
			"default value of %s is invalid: %s",
			s.portResult.Name,
			err,
		))
	}

	return s.with(func() {
		// TODO: https://github.com/dogmatiq/ferrite/issues/1

		s.host = host
		s.port = port

		s.hostResult.DefaultValue = host
		s.portResult.DefaultValue = port
	})
}

// Address returns the address (host:port) of the Kubernetes service.
func (s *KubeServiceSpec) Address() string {
	return net.JoinHostPort(s.Host(), s.Port())
}

// Host returns the hostname or IP address of the Kubernetes service.
func (s *KubeServiceSpec) Host() string {
	s.resolve()

	if s.hostResult.Error != nil {
		panic(fmt.Sprintf(
			"%s is invalid: %s",
			s.hostResult.Name,
			s.hostResult.Error,
		))
	}

	return s.host
}

// Port returns the port of the Kubernetes service.
//
// It may be a port number of an IANA registered service name (e.g. "https").
func (s *KubeServiceSpec) Port() string {
	s.resolve()

	if s.portResult.Error != nil {
		panic(fmt.Sprintf(
			"%s is invalid: %s",
			s.portResult.Name,
			s.portResult.Error,
		))
	}

	return s.port
}

// Validate validates the environment variables.
func (s *KubeServiceSpec) Validate() []ValidationResult {
	s.resolve()

	return []ValidationResult{
		s.hostResult,
		s.portResult,
	}
}

// resolve populates s.host, s.port and the validation results, or returns
// immediately if they are already populated.
func (s *KubeServiceSpec) resolve() {
	// TODO: https://github.com/dogmatiq/ferrite/issues/1

	s.seal.Close(func() {
		host := os.Getenv(s.hostResult.Name)
		port := os.Getenv(s.portResult.Name)

		if host != "" {
			s.hostResult.ExplicitValue = host
			s.host = host
		} else if s.host != "" {
			s.hostResult.UsingDefault = true
		} else {
			s.hostResult.Error = errUndefined
		}

		if port != "" {
			s.portResult.ExplicitValue = port
			s.port = port
		} else if s.port != "" {
			s.portResult.UsingDefault = true
		} else {
			s.portResult.Error = errUndefined
		}
	})
}

// with calls fn while holding a lock on s.
//
// It panics if the value has already been resolved.
func (s *KubeServiceSpec) with(fn func()) *KubeServiceSpec {
	s.seal.Do(fn)
	return s
}

// kubeToEnv converts a kubernetes resource name to an environment variable
// name, as per Kubernetes own behavior.
func kubeToEnv(s string) string {
	return strings.ToUpper(
		strings.ReplaceAll(s, "-", "_"),
	)
}

// validateHost reutrns an error if host is not a valid hostname or IP address.
//
// The hostname validation is intentionally very permissive as it's not unheard
// of to encounter functioning services in the wild that have DNS names that are
// technically invalid. This includes hostnames that start with hyphens!
func validateHost(host string) error {
	if host == "" {
		return errUndefined
	}

	if net.ParseIP(host) != nil {
		return nil
	}

	n := len(host)

	if host[0] == '.' || host[n-1] == '.' {
		return errors.New("hostname must not begin or end with a dot")
	}

	for _, r := range host {
		if unicode.IsSpace(r) {
			return errors.New("hostname must not contain whitespace")
		}
	}

	return nil
}

// validatePort returns an error if port is neither a numeric port number, nor a
// valid IANA registered service name.
func validatePort(port string) error {
	if port == "" {
		return errUndefined
	}

	n, err := strconv.ParseUint(port, 10, 16)

	if errors.Is(err, strconv.ErrSyntax) {
		if err := validateIANAServiceName(port); err != nil {
			return fmt.Errorf("%q is not a valid IANA service name (%s)", port, err)
		}

		return nil
	}

	if err == nil && n != 0 {
		return nil
	}

	return errors.New("numeric ports must be between 1 and 65535")
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
		return errors.New("must be between 1 and 15 characters")
	}

	// RFC-6335: MUST NOT begin or end with a hyphen.
	if name[0] == '-' || name[n-1] == '-' {
		return errors.New("must not begin or end with a hyphen")
	}

	hasLetter := false

	for i := range name {
		ch := name[i] // iterate by byte (not rune)

		// RFC-6335: MUST contain only US-ASCII [ANSI.X3.4-1986] letters 'A' -
		// 'Z' and 'a' - 'z', digits '0' - '9', and hyphens ('-', ASCII 0x2D or
		// decimal 45).
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
				return errors.New("must not contain adjacent hyphens")
			}
		default:
			return errors.New("must contain only ASCII letters, digits and hyphen")
		}
	}

	//RFC-6335: MUST contain at least one letter ('A' - 'Z' or 'a' - 'z').
	if !hasLetter {
		return errors.New("must contain at least one letter")
	}

	return nil
}
