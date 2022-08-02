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
// name is the name of the Kubernetes service, NOT the environment variable.
func KubeService(name, desc string) *KubeServiceSpec {
	s := &KubeServiceSpec{
		service: name,
		hostResult: ValidationResult{
			Name: fmt.Sprintf("%s_SERVICE_HOST", kubeToEnv(name)),
		},
		portResult: ValidationResult{
			Name: fmt.Sprintf("%s_SERVICE_PORT", kubeToEnv(name)),
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

func (s *KubeServiceSpec) WithNamedPort(name string) {
	s.portResult.Name = fmt.Sprintf(
		"%s_SERVICE_PORT_%s",
		kubeToEnv(s.service),
		kubeToEnv(name),
	)
}

// Address returns the address of the Kubernetes service.
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

// Validate validates the environment variables.
func (s *KubeServiceSpec) Validate() []ValidationResult {
	return []ValidationResult{
		s.hostResult,
		s.portResult,
	}
}

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
