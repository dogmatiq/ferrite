package ferrite

import (
	"unsafe"
)

// bitSize returns the number of bits used to represent T.
func bitSize[T any]() int {
	var zero T
	return int(unsafe.Sizeof(zero)) * 8
}
