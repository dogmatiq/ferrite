package ferrite

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/dogmatiq/ferrite/internal/table"
	"github.com/dogmatiq/ferrite/schema"
	"golang.org/x/exp/slices"
)

// ValidateEnvironment validates all environment variables.
func ValidateEnvironment() {
	if result, ok := validate(); !ok {
		io.WriteString(output, result)
		exit(1)
	}
}

// Register adds a validator to the global validation system.
func Register(s Spec) {
	specsM.Lock()
	specs = append(specs, s)
	specsM.Unlock()
}

var (
	// specs is a global set of variable specifications that are invoked by
	// ValidateEnvironment().
	specsM sync.Mutex
	specs  []Spec

	// output is the writer to which the validation result is written.
	output io.Writer = os.Stderr

	// exit is called to exit the process when validation fails.
	exit = os.Exit
)

type pair struct {
	Variable VariableXXX
	Result   ValidationResult
}

// validate parses and validates all environment variables.
func validate() (string, bool) {
	specsM.Lock()
	defer specsM.Unlock()

	var results []pair
	ok := true

	for _, s := range specs {
		vars := s.Describe()

		for _, res := range s.Validate() {
			var p pair
			p.Result = res

			for _, v := range vars {
				if v.Name == res.Name {
					p.Variable = v
				}
			}

			if res.Error != nil {
				ok = false
			}

			results = append(results, p)
		}
	}

	return renderResults(results), ok
}

const (
	// validIcon is the icon displayed next to validIcon environment variables.
	validIcon = "✓"

	// invalidIcon is the icon displayed next to invalidIcon environment variables.
	invalidIcon = "✗"

	// chevronIcon is the icon used to draw attention to invalid environment
	// variables.
	chevronIcon = "❯"
)

// renderResults renders a set of validation results as a human-readable string.
func renderResults(results []pair) string {
	slices.SortFunc(
		results,
		func(a, b pair) bool {
			return a.Variable.Name < b.Variable.Name
		},
	)

	var t table.Table

	for _, v := range results {
		name := " "
		if v.Result.Error != nil {
			name += chevronIcon
		} else {
			name += " "
		}
		name += " " + v.Variable.Name

		renderer := &validateSchemaRenderer{}
		v.Variable.Schema.AcceptVisitor(renderer)

		input := renderer.Output.String()
		if v.Variable.Default != "" {
			input += fmt.Sprintf(" = %q", v.Variable.Default)
		}

		status := ""
		if v.Result.Error != nil {
			status += invalidIcon + " " + v.Result.Error.Error()
		} else if v.Result.UsedDefault {
			status += fmt.Sprintf("%s using default value", validIcon)
		} else {
			status += fmt.Sprintf("%s set to %q", validIcon, v.Result.Value)
		}

		t.AddRow(name, input, v.Variable.Description, status)
	}

	return "ENVIRONMENT VARIABLES:\n" + t.String()
}

type validateSchemaRenderer struct {
	Output strings.Builder
}

func (r *validateSchemaRenderer) VisitOneOf(s schema.OneOf) {
	for i, c := range s {
		if i > 0 {
			r.Output.WriteString("|")
		}

		c.AcceptVisitor(r)
	}
}

func (r *validateSchemaRenderer) VisitLiteral(s schema.Literal) {
	r.Output.WriteString(string(s))
}

func (r *validateSchemaRenderer) VisitType(s schema.TypeSchema) {
	fmt.Fprintf(&r.Output, "[%s]", s.Type)
}

func (r *validateSchemaRenderer) VisitRange(s schema.Range) {
	if s.Min != "" && s.Max != "" {
		fmt.Fprintf(&r.Output, "(%s..%s)", s.Min, s.Max)
	} else if s.Max != "" {
		fmt.Fprintf(&r.Output, "(...%s)", s.Max)
	} else {
		fmt.Fprintf(&r.Output, "(%s...)", s.Min)
	}
}

// Spec is a specification for an environment variable.
type Spec interface {
	// Describe returns a description of the environment variable(s) described
	// by this spec.
	Describe() []VariableXXX

	// Validate validates the environment variable(s) described by this spec.
	Validate() []ValidationResult
}

// VariableXXX describes an environment variable.
type VariableXXX struct {
	// Name is the name of the environment variable.
	Name string

	// Description is a human-readable description of the environment variable.
	Description string

	// Schema describes the valid values for this environment variable.
	Schema schema.Schema

	// Default is the environment variable's default value.
	//
	// It must be non-empty if the environment variable has a default value;
	// otherwise it must be empty.
	Default string
}

// ValidationResult is the result of validating an environment variable.
type ValidationResult struct {
	// Name is the name of the environment variable.
	Name string

	// Value is the environment variable's value.
	Value string

	// UsedDefault is true if the default value of the environment variable was
	// used to populate the value.
	UsedDefault bool

	// Error is an error describing why the validation failed.
	//
	// If it is nil, the validation is considered successful.
	Error error
}
