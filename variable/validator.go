package variable

// Validator validates values of type T.
type Validator[T any] interface {
	Validate(T) error
}

// validatorFunc is a function that implements the Validator interface.
type validatorFunc[T any] func(T) error

// ValidatorFunc returns a validator implemented by fn.
func ValidatorFunc[T any](fn func(T) error) Validator[T] {
	return validatorFunc[T](fn)
}

// Validate returns an error if v is invalid.
func (f validatorFunc[T]) Validate(v T) error {
	return f(v)
}
