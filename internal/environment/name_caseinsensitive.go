//go:build windows

package environment

import "strings"

// normalizeName normalizes an environment variable name, such that it compares
// as equals to other variables with the same name.
func normalizeName(n string) string {
	return strings.ToUpper(n)
}
