package ferrite

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

const (
	// valid is the icon displayed next to valid environment variables.
	valid = "✓"

	// invalid is the icon displayed next to invalid environment variables.
	invalid = "✗"

	// chevron is the icon used to draw attention to invalid environment
	// variables.
	chevron = "❯"
)

// RegistryValidationResult is the result of validating environment variables.
type RegistryValidationResult struct {
	Variables []VariableValidationResult
}

// IsValid returns true if all environment variables in the registry are valid.
func (r RegistryValidationResult) IsValid() bool {
	for _, v := range r.Variables {
		if v.Error != nil {
			return false
		}
	}

	return true
}

// String returns a human-readable representation of the validation result.
func (r RegistryValidationResult) String() string {
	slices.SortFunc(
		r.Variables,
		func(a, b VariableValidationResult) bool {
			return a.Name < b.Name
		},
	)

	var t table
	for _, v := range r.Variables {
		name := " "
		if v.Error != nil {
			name += chevron
		} else {
			name += " "
		}
		name += " " + v.Name

		input := v.ValidInput
		if v.DefaultValue != "" {
			input += " = " + v.DefaultValue
		}

		status := ""
		if v.Error != nil {
			status += invalid + " " + v.Error.Error()
		} else if v.UsingDefault {
			status += valid + " using default value"
		} else {
			status += valid + " set to " + v.ExplicitValue
		}

		t.AddRow(name, input, v.Description, status)
	}

	return "ENVIRONMENT VARIABLES:\n" + t.String()
}

// VariableValidationResult contains information about a single environment
// variable on a RegistryValidationResult.
type VariableValidationResult struct {
	Name          string
	Description   string
	ValidInput    string
	DefaultValue  string
	ExplicitValue string
	UsingDefault  bool
	Error         error
}

// table renders a column-aligned table.
type table struct {
	widths []int
	rows   [][]string
}

// AddRow adds a row to the table.
func (t *table) AddRow(columns ...string) {
	for len(t.widths) < len(columns) {
		t.widths = append(t.widths, 0)
	}

	for i, col := range columns {
		if len(col) > t.widths[i] {
			t.widths[i] = len(col)
		}
	}

	t.rows = append(t.rows, columns)
}

// String returns the rendered table.
func (t *table) String() string {
	var w strings.Builder

	for _, columns := range t.rows {
		n := len(columns) - 1

		for i, col := range columns[:n] {
			fmt.Fprintf(
				&w,
				"%-*s  ",
				t.widths[i],
				col,
			)
		}

		fmt.Fprintln(&w, columns[n])
	}

	return w.String()
}
