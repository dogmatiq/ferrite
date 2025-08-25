package variable

type Default[T any] interface {
	finalize() error
}

func ConstDefault[T any](v T) Default[T] {
	return constDefault[T]{Native: v}
}

type constDefault[T any] struct {
	Native T
}

func (d constDefault[T]) x() {}

type defaultFromBuilder[T any] struct {
}

func DefaultFrom[T any]() {
}

// type defaultFromVariableSet[T any] struct {
// }

// func (d DefaultFromVariableSet[T]) x() {}
