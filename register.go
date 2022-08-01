package ferrite

var (
	// registry is a global registry of environment variable specs.
	registry map[string]Spec
)

// Register adds a variable specification to the register.
func Register[T any, S SpecFor[T]](s S) {
	if registry == nil {
		registry = map[string]Spec{}
	}

	registry[s.Name()] = s
}
