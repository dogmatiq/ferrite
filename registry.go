package ferrite

// ResolveEnvironment validates all environment variables.
//
// It panics if any of the defined variables are invalid.
func ResolveEnvironment() {
	if err := DefaultRegistry.Resolve(); err != nil {
		panic(err.Error())
	}
}

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

// Resolve parses and validates all environment variables in the registry,
// allowing their associated values to be obtained.
func (r *Registry) Resolve() error {
	for _, s := range r.specs {
		if err := s.Resolve(); err != nil {
			return err
		}
	}

	return nil
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
