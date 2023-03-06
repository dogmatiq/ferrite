package variable

import (
	"github.com/dogmatiq/ferrite/maybe"
)

// Marshaler converts values to/from their native/literal representations.
type Marshaler[T any] interface {
	// Marshal converts a value to its literal representation.
	Marshal(T) (Literal, error)

	// Unmarshal converts a literal value to it's native representation.
	Unmarshal(Literal) (T, error)
}

// marshal marshals a native "maybe" value to a literal.
func marshal[T any, M Marshaler[T]](
	m M,
	v maybe.Value[T],
) (maybe.Value[Literal], error) {
	return maybe.TryMap(
		v,
		m.Marshal,
	)
}

// mustMarshal marshals a native "maybe" value to a literal or panics if unable
// to do so.
func mustMarshal[T any, M Marshaler[T]](
	m M,
	v maybe.Value[T],
) maybe.Value[Literal] {
	return maybe.Map(
		v,
		func(n T) Literal {
			lit, err := m.Marshal(n)
			if err != nil {
				panic(err)
			}

			return lit
		},
	)
}
