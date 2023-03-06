package limits

import (
	"math"
	"reflect"

	"github.com/dogmatiq/ferrite/internal/reflectx"
	"golang.org/x/exp/constraints"
)

// Of returns the minimum and maximum values of type T.
func Of[T constraints.Integer | constraints.Float]() (min, max T) {
	switch reflectx.KindOf[T]() {

	// constraints.Signed ...
	case reflect.Int:
		return limits[int, T](math.MinInt, math.MaxInt)
	case reflect.Int8:
		return limits[int8, T](math.MinInt8, math.MaxInt8)
	case reflect.Int16:
		return limits[int16, T](math.MinInt16, math.MaxInt16)
	case reflect.Int32:
		return limits[int32, T](math.MinInt32, math.MaxInt32)
	case reflect.Int64:
		return limits[int64, T](math.MinInt64, math.MaxInt64)

	// constraints.Unsigned ...
	case reflect.Uint:
		return limits[uint, T](0, math.MaxUint)
	case reflect.Uint8:
		return limits[uint8, T](0, math.MaxUint8)
	case reflect.Uint16:
		return limits[uint16, T](0, math.MaxUint16)
	case reflect.Uint32:
		return limits[uint32, T](0, math.MaxUint32)
	case reflect.Uint64:
		return limits[uint64, T](0, math.MaxUint64)

	// constraints.Float ...
	case reflect.Float32:
		return limits[float32, T](-math.MaxFloat32, +math.MaxFloat32)
	case reflect.Float64:
		return limits[float32, T](-math.MaxFloat32, +math.MaxFloat32)

	default:
		panic("not implemented")
	}
}

// limits converts constants of type C to values of type U.
//
// It is used in place of the language's built-in type-conversion syntax to
// avoid compile time errors about overflows when converting constants. This
// allows use of numeric constants in generic code.
func limits[C, T constraints.Integer | constraints.Float](min, max C) (T, T) {
	return T(min), T(max)
}
