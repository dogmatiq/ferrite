package ferrite

// String configures an environment variable as a string.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func String(name, desc string) *StringSpec[string] {
	return StringAs[string](name, desc)
}

// StringAs configures an environment variable as a string using a user-defined
// type.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func StringAs[T ~string](name, desc string) *StringSpec[T] {
	s := &StringSpec[T]{}
	s.init(s, name, desc)
	return s
}

// StringSpec is the specification for a string.
type StringSpec[T ~string] struct {
	standard[T, *StringSpec[T]]
}

// parses parses and validates the value of the environment variable.
func (s *StringSpec[T]) parse(value string, def *T) (T, ValidationResult) {
	res := ValidationResult{
		Name:          s.name,
		Description:   s.desc,
		ValidInput:    inputOfType[T](),
		DefaultValue:  renderString(def),
		ExplicitValue: renderString(&value),
	}

	if value != "" {
		return T(value), res
	}

	if def != nil {
		res.UsingDefault = true
		return *def, res
	}

	res.Error = errUndefined

	return "", res
}

// validate validates a parsed or default value.
func (s *StringSpec[T]) validate(value T) error {
	return nil
}
