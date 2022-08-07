package spec

import "reflect"

// Schema describes the valid values of an environment variable.
type Schema interface {
	// AcceptVisitor calls the method on v that is related to this schema type.
	AcceptVisitor(v Visitor)
}

// Visitor dispatches logic based on a schema element.
type Visitor interface {
	VisitOneOf(OneOf)
	VisitLiteral(Literal)
	VisitType(Type)
	VisitRange(Range)
}

// OneOf is a schema that allows any of a number of other schemas.
type OneOf []Schema

// AcceptVisitor calls the method on v that is related to this schema type.
func (s OneOf) AcceptVisitor(v Visitor) {
	v.VisitOneOf(s)
}

// Literal is a Schema that allows a specific string value.
type Literal string

// AcceptVisitor calls the method on v that is related to this schema type.
func (s Literal) AcceptVisitor(v Visitor) {
	v.VisitLiteral(s)
}

// Type is a schema that can be any value of a specific type.
type Type struct {
	Type reflect.Type
}

// OfType returns a schema that allows any value of type T.
func OfType[T any]() Type {
	return Type{
		Type: reflect.TypeOf([...]T{}).Elem(),
	}
}

// AcceptVisitor calls the method on v that is related to this schema type.
func (s Type) AcceptVisitor(v Visitor) {
	v.VisitType(s)
}

// Range is a schema that allows values within a specific range.
type Range struct {
	Min string
	Max string
}

// AcceptVisitor calls the method on v that is related to this schema type.
func (s Range) AcceptVisitor(v Visitor) {
	v.VisitRange(s)
}
