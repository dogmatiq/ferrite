package reflectx

import "reflect"

// TypeOf returns the type of T.
func TypeOf[T any]() reflect.Type {
	return reflect.TypeOf([...]T{}).Elem()
}

// KindOf returns the kind of T.
func KindOf[T any]() reflect.Kind {
	return TypeOf[T]().Kind()
}
