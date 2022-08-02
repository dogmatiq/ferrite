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
	service    string
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

	s.portResult.Description = fmt.Sprintf(
		`Network port of the "%s" service's "%s" port.`,
		s.service,
		port,
	)
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
// It may be a port number of an IANA registered port name (e.g. "https").
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

func (s *KubeServiceSpec) resolve() {
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
}

// WithDefault sets a default value to use when the environment variables are
// undefined.
//
// port may be a port number or an IANA registered port name (such as "https").
// The IANA name is not to be confused with the Kubernetes port name.
func (s *KubeServiceSpec) WithDefault(host, port string) *KubeServiceSpec {
	// TODO: validate host/port
	s.host = host
	s.port = port

	s.hostResult.DefaultValue = host
	s.portResult.DefaultValue = port

	return s
}

// kubeToEnv converts a kubernetes resource name to an environment variable
// name, as per Kubernetes own behavior.
func kubeToEnv(s string) string {
	return strings.ToUpper(
		strings.ReplaceAll(s, "-", "_"),
	)
}
