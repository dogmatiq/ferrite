package markdownmode

import (
	"github.com/dogmatiq/ferrite/variable"
	"github.com/mattn/go-runewidth"
)

func (r *renderer) renderExamples(s variable.Spec) {
	r.line("```bash")

	width := 0
	for _, eg := range s.Examples() {
		w := runewidth.StringWidth(eg.Canonical.Quote())
		if w > width {
			width = w
		}
	}

	for _, eg := range s.Examples() {
		comment := ""
		if variable.IsDefault(s, eg.Canonical) {
			comment = "(default)"
		} else if !eg.IsNormative {
			comment = "(non-normative)"
		}

		if eg.Description != "" {
			if comment != "" {
				comment += " "
			}
			comment += eg.Description
		}

		if len(comment) == 0 {
			r.line(
				"export %s=%s",
				s.Name(),
				eg.Canonical.Quote(),
			)
		} else {
			r.line(
				"export %s=%-*s # %s",
				s.Name(),
				width,
				eg.Canonical.Quote(),
				comment,
			)
		}
	}

	r.line("```")
}

func (r *renderer) renderNonNormativeExampleWarning() {
	r.line("⚠️ Some of the variables have **non-normative** examples. These examples are")
	r.line("syntactically correct but may not be meaningful values for this application.")
}
