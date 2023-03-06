package reflectx

import "unsafe"

// BitSize returns the number of bits used to represent T.
func BitSize[T any]() int {
	var zero T
	return int(unsafe.Sizeof(zero)) * 8
}
