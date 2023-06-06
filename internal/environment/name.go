package environment

// EqualNames returns true if the given names refer to the same environment
// variable.
func EqualNames(a, b string) bool {
	return normalizeName(a) == normalizeName(b)
}

// NormalizeName normalizes an environment variable name for use as a key within
// a map.
func NormalizeName(n string) string {
	return normalizeName(n)
}
