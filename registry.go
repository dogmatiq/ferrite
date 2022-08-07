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
	// valid is the icon displayed next to valid environment variables.
	valid = "✓"

	// invalid is the icon displayed next to invalid environment variables.
	invalid = "✗"

	// chevron is the icon used to draw attention to invalid environment
	// variables.
	chevron = "❯"
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
			name += chevron
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
			status += invalid + " " + v.Result.Error.Error()
		} else if v.Result.UsedDefault {
			status += fmt.Sprintf("%s using default value", valid)
		} else {
			status += fmt.Sprintf("%s set to %q", valid, v.Result.Value)
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
