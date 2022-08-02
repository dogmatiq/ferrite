package ferrite

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// KubeService reads Kubernetes service discovery environment variables for a
// service's port.
//
// svc is the name of the Kubernetes service, NOT the environment variable.
//
// The environment variables "<svc>_SERVICE_HOST" and "<svc>_SERVICE_PORT" are
// used to construct an address for the service.
func KubeService(svc, desc string) *KubeServiceSpec {
	s := &KubeServiceSpec{
		service: svc,
		hostResult: ValidationResult{
			Name: fmt.Sprintf("%s_SERVICE_HOST", kubeToEnv(svc)),
		},
		portResult: ValidationResult{
			Name: fmt.Sprintf("%s_SERVICE_PORT", kubeToEnv(svc)),
		},
	}

	Register(s)

	return s
}

// KubeServiceSpec is the specification for a Kubernetes service.
type KubeServiceSpec struct {
	service    string
	defaulted  bool
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
// is not to be confused with an IANA registered port name (e.g. "https"),
// although the two may use the same names.
//
// See https://kubernetes.io/docs/concepts/services-networking/service/#multi-port-services
func (s *KubeServiceSpec) WithNamedPort(port string) {
	s.portResult.Name = fmt.Sprintf(
		"%s_SERVICE_PORT_%s",
		kubeToEnv(s.service),
		kubeToEnv(port),
	)
}

// Address returns the address (host:port) of the Kubernetes service.
func (s *KubeServiceSpec) Address() string {
	host := os.Getenv(s.hostResult.Name)
	port := os.Getenv(s.portResult.Name)

	if s.defaulted {
		if host == "" {
			host = s.host
		}

		if port == "" {
			port = s.port
		}
	} else if host == "" {
		panic(fmt.Sprintf(
			"%s is invalid: %s",
			s.hostResult.Name,
			errUndefined,
		))
	} else if port == "" {
		panic(fmt.Sprintf(
			"%s is invalid: %s",
			s.portResult.Name,
			errUndefined,
		))
	}

	return net.JoinHostPort(host, port)
}

// Port returns the port of the Kubernetes service.
//
// It may be a port number of an IANA registered port name (e.g. "https").
func (s *KubeServiceSpec) Port() string {
	port := os.Getenv(s.portResult.Name)

	if s.defaulted {
		if port == "" {
			port = s.port
		}
	} else if port == "" {
		panic(fmt.Sprintf(
			"%s is invalid: %s",
			s.portResult.Name,
			errUndefined,
		))
	}

	return port
}

// Validate validates the environment variables.
func (s *KubeServiceSpec) Validate() []ValidationResult {
	return []ValidationResult{
		s.hostResult,
		s.portResult,
	}
}

// WithDefault sets a default value to use when the environment variables are
// undefined.
//
// port may be a port number or an IANA registered port name (such as "https").
// The IANA name is not to be confused with the Kubernetes port name.
func (s *KubeServiceSpec) WithDefault(host, port string) *KubeServiceSpec {
	// TODO: validate host/port
	s.defaulted = true
	s.host = host
	s.port = port

	return s
}

// kubeToEnv converts a kubernetes resource name to an environment variable
// name, as per Kubernetes own behavior.
func kubeToEnv(s string) string {
	return strings.ToUpper(
		strings.ReplaceAll(s, "-", "_"),
	)
}
