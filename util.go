package ferrite

// ptr returns a pointer to a new variable with the given value.
func ptr[T any](v T) *T {
	return &v
}
