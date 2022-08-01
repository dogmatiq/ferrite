package ferrite

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

// ValidateEnvironment validates all environment variables.
func ValidateEnvironment() {
	if result, ok := validate(); !ok {
		io.WriteString(output, result)
		exit(1)
	}
}

// Register adds a variable specification to the register.
//
// It can be used to register custom specifications with Ferrite's validation
// system.
func Register(name string, v Validator) {
	if registry == nil {
		registry = map[string]Validator{}
	}

	registry[name] = v
}

// Validator is an interface used to validate environment variables
type Validator interface {
	// Validate validates the environment variable.
	Validate(name, value string) ValidationResult
}

// ValidationResult is the result of validating an environment variable.
type ValidationResult struct {
	// Name is the name of the environment variable.
	Name string

	// Description is a human-readable description of the environment variable.
	Description string

	// ValidInput is a human-readable description of what constitutes valid
	// input for this environment variable.
	//
	// It is free-form text, but spec implementations SHOULD adhere to the
	// following conventions:
	//
	//  - separate enum style inputs with pipe, e.g. "true|false"
	//  - render type names inside square barckets, e.g. "[string]"
	ValidInput string

	// DefaultValue is the environment variable's default value, rendered as it
	// should be displayed in the console.
	//
	// This is not necessarily equal to a raw environment variable value. For
	// example, StringSpec renders strings with surrounding quotes.
	//
	// It must be non-empty if the environment variable has a default value;
	// otherwise it must be empty.
	DefaultValue string

	// ExplicitValue is the environment variable's value as captured from the
	// environment, rendered as it should be displayed in the console.
	//
	// This is not necessarily equal to the raw environment variable value. For
	// example, StringSpec renders strings with surrounding quotes.
	ExplicitValue string

	// UsingDefault is true if the environment variable's default value
	// would be retuned by the specs Value() method.
	UsingDefault bool

	// Error is an error describing why the validation failed.
	//
	// If it is nil, the validation is considered successful.
	Error error
}

var (
	// registry is a global registry of environment variable specs.
	registry map[string]Validator

	// output is the writer to which the validation result is written.
	output io.Writer = os.Stderr

	// exit is called to exit the process when validation fails.
	exit = os.Exit
)

// validate parses and validates all environment variables.
func validate() (string, bool) {
	var results []ValidationResult

	ok := true
	for n, s := range registry {
		res := s.Validate(n, os.Getenv(n))
		if res.Error != nil {
			ok = false
		}

		results = append(results, res)
	}

	return renderResults(results), ok
}

const (
	// valid is the icon displayed next to valid environment variables.
	valid = "✓"

	// invalid is the icon displayed next to invalid environment variables.
	invalid = "✗"

	// chevron is the icon used to draw attention to invalid environment
	// variables.
	chevron = "❯"
)

// renderResults renders a set of validation results as a human-readable string.
func renderResults(results []ValidationResult) string {
	slices.SortFunc(
		results,
		func(a, b ValidationResult) bool {
			return a.Name < b.Name
		},
	)

	var t table

	for _, v := range results {
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
