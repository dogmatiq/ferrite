package maybe

// Value is an optional value of type T.
type Value[T any] struct {
	value T
	ok    bool
}

// Some returns a non-empty value.
func Some[T any](v T) Value[T] {
	return Value[T]{v, true}
}

// None returns an empty value.
func None[T any]() Value[T] {
	return Value[T]{}
}

// MustGet returns the optional value or panics if m is empty.
func (m Value[T]) MustGet() T {
	if m.ok {
		return m.value
	}

	panic("maybe-value is empty")
}

// Get returns the optional value.
//
// If m is empty ok is false; otherwise, ok is true and v is the value.
func (m Value[T]) Get() (v T, ok bool) {
	return m.value, m.ok
}

// IsEmpty returns true if m is empty, meaning that it does not have a value.
func (m Value[T]) IsEmpty() bool {
	return !m.ok
}

// Map converts m to a maybe-value of type U by applying a mapping function.
func Map[T, U any](m Value[T], fn func(T) U) Value[U] {
	if v, ok := m.Get(); ok {
		return Some(fn(v))
	}

	return None[U]()
}

// TryMap converts m to a maybe-value of type U by applying a mapping function.
func TryMap[T, U any](m Value[T], fn func(T) (U, error)) (Value[U], error) {
	v, ok := m.Get()
	if !ok {
		return None[U](), nil
	}

	x, err := fn(v)
	if err != nil {
		return None[U](), err
	}

	return Some(x), nil
}
