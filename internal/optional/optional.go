package optional

// Optional holds an Optional value of type T.
type Optional[T any] struct {
	v  T
	ok bool
}

// With returns an optional the value v.
func With[T any](v T) Optional[T] {
	return Optional[T]{v, true}
}

// Without returns an empty optional value of type T.
func Without[T any]() Optional[T] {
	return Optional[T]{}
}

// Get returns the optional value if it is present.
func (o Optional[T]) Get() (T, bool) {
	return o.v, o.ok
}

// Coalesce returns o if it has a value, otherwise it returns an optional
// populated with v.
func (o Optional[T]) Coalesce(v T) Optional[T] {
	if o.ok {
		return o
	}

	return With(v)
}
