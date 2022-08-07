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

// without returns an empty optional value of type T.
func without[T any]() Optional[T] {
	return Optional[T]{}
}

// TryGet returns the optional value if it is present.
func (o Optional[T]) TryGet() (T, bool) {
	return o.v, o.ok
}

// Get returns the optional value, or panics if the value is not present.
func (o Optional[T]) Get() T {
	if !o.ok {
		panic("value is not present")
	}

	return o.v
}
