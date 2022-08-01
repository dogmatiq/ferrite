package ferrite

import (
	"reflect"
)

// typeName returns T's name.
func typeName[T any]() string {
	return reflect.
		TypeOf((*T)(nil)).
		Elem().
		Name()
}
