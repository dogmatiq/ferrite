package markdown

import (
	"github.com/dogmatiq/ferrite/variable"
	"github.com/mattn/go-runewidth"
)

// hasNonNormativeExamples returns true if any of the specs have non-normative
// examples.
func (r *renderer) hasNonNormativeExamples() bool {
	for _, s := range r.Specs {
		for _, eg := range s.Examples() {
			if !eg.IsNormative {
				return true
			}
		}
	}

	return false
}

func (r *specRenderer) renderExamples() {
	r.ren.gap()
	r.ren.line("```bash")

	width := 0
	for _, eg := range r.spec.Examples() {
		w := runewidth.StringWidth(eg.Canonical.Quote())
		if w > width {
			width = w
		}
	}

	for _, eg := range r.spec.Examples() {
		comment := ""
		if variable.IsDefault(r.spec, eg.Canonical) {
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
			r.ren.line(
				"export %s=%s",
				r.spec.Name(),
				eg.Canonical.Quote(),
			)
		} else {
			r.ren.line(
				"export %s=%-*s # %s",
				r.spec.Name(),
				width,
				eg.Canonical.Quote(),
				comment,
			)
		}
	}

	r.ren.line("```")
}
