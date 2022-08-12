package variable

import "github.com/dogmatiq/ferrite/maybe"

// Marshaler converts values to/from their native/literal representations.
type Marshaler[T any] interface {
	// Marshal converts a value to its literal representation.
	Marshal(T) (Literal, error)

	// Unmarshal converts a literal value to it's native representation.
	Unmarshal(Literal) (T, error)
}

func marshal[T any, M Marshaler[T]](
	m M,
	v maybe.Value[T],
) (maybe.Value[valueOf[T]], error) {
	return maybe.TryMap(
		v,
		func(n T) (valueOf[T], error) {
			lit, err := m.Marshal(n)
			if err != nil {
				return valueOf[T]{}, err
			}

			return valueOf[T]{
				canonical: lit,
				native:    n,
			}, nil
		},
	)
}
