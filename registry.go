package ferrite

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// ValidateEnvironment validates all environment variables.
//
// It panics if any of the defined variables are invalid.
func ValidateEnvironment() {
	var w bytes.Buffer
	if err := DefaultRegistry.Validate(&w); err != nil {
		io.Copy(stderr, &w)
		exit(1)
	}
}

var (
	stderr io.Writer = os.Stderr
	exit             = os.Exit
)

// DefaultRegistry is the default environment variable registry.
var DefaultRegistry = Registry{}

// Registry is a container of environment variable specifications.
type Registry struct {
	specs map[string]Spec
}

// Register adds a variable specification to the register.
func (r *Registry) Register(s Spec) {
	if r.specs == nil {
		r.specs = map[string]Spec{}
	}

	r.specs[s.Name()] = s
}

// Reset removes all variable specifications from the registry.
func (r *Registry) Reset() {
	r.specs = nil
}

// Validate parses and validates all environment variables in the registry,
// allowing their associated values to be obtained.
func (r *Registry) Validate(w io.Writer) error {
	var rows [][]string
	var cause error

	specs := maps.Values(r.specs)
	slices.SortFunc(specs, func(a, b Spec) bool {
		return a.Name() < b.Name()
	})

	for _, s := range specs {
		marker := " "
		status := ""

		value, isDefault, err := s.Validate()
		if err == nil {
			if isDefault {
				status = fmt.Sprintf("%s using default value", pass)
			} else {
				status = fmt.Sprintf("%s set to %s", pass, value)
			}
		} else {
			marker = chevron
			status = fmt.Sprintf("%s %s", fail, err)

			if cause == nil {
				cause = err
			}
		}

		desc := s.Describe()

		input := desc.Input
		if desc.Default != "" {
			input += " = " + desc.Default
		}

		rows = append(rows, []string{
			" " + marker + " " + s.Name(),
			input,
			desc.Variable,
			status,
		})
	}

	if w == nil {
		w = io.Discard
	}

	if _, err := fmt.Fprintln(w, "ENVIRONMENT VARIABLES:"); err != nil {
		return err
	}

	if err := renderTable(w, rows); err != nil {
		return err
	}

	return cause
}

// register adds s to the registry configured by the given options.
func register(s Spec, options []SpecOption) {
	opts := specOptions{
		Registry: &DefaultRegistry,
	}

	for _, opt := range options {
		opt(&opts)
	}

	opts.Registry.Register(s)
}
