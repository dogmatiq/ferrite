package numeric

import (
	"unsafe"

	"golang.org/x/exp/constraints"
)

// BitSize returns the number of bits used to represent T.
func BitSize[T constraints.Integer | constraints.Float]() int {
	var zero T
	return int(unsafe.Sizeof(zero)) * 8
}
