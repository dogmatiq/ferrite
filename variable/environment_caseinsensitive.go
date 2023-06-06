//go:build windows

package variable

import "strings"

// normalizeVariableName normalizes an environment variable name, such that it
// compares as equals to other variables with the same name.
func normalizeVariableName(n string) string {
	return strings.ToUpper(n)
}
