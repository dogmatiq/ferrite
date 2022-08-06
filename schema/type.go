package schema

import "reflect"

// TypeSchema is a schema that can be any value of a specific type.
type TypeSchema struct {
	Type reflect.Type
}

// Type returns a schema that allows any value of type T.
func Type[T any]() TypeSchema {
	return TypeSchema{
		Type: reflect.TypeOf([...]T{}).Elem(),
	}
}

// AcceptVisitor calls the method on v that is related to this schema type.
func (s TypeSchema) AcceptVisitor(v Visitor) {
	v.VisitType(s)
}
