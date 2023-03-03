package markdownmode

import (
	"fmt"
	"strings"

	"github.com/dogmatiq/ferrite/variable"
	"gopkg.in/yaml.v3"
)

type renderer struct {
	App   string
	Specs []variable.Spec

	w                 strings.Builder
	hideUsageExamples bool
	refs              map[string]struct{}
}

func (r *renderer) Render() string {
	r.line("# Environment Variables")
	r.line("")
	r.line("This document describes the environment variables used by `%s`.", r.App)

	if len(r.Specs) == 0 {
		r.line("")
		r.line("**There do not appear to be any environment variables.**")
	}

	if r.hasNonNormativeExamples() {
		r.line("")
		r.renderNonNormativeExampleWarning()
	}

	r.line("")
	r.line("⚠️ The application may consume other undocumented environment variables; this")
	r.line("document only shows those variables declared using %s.", r.link("Ferrite"))

	if len(r.Specs) != 0 {
		r.line("")
		r.line("Please note that **undefined** variables and **empty strings** are considered")
		r.line("equivalent.")
	}

	if len(r.Specs) != 0 {
		r.line("")
		r.renderIndex()
		r.line("")
		r.renderSpecs()

		if !r.hideUsageExamples {
			r.line("")
			r.renderUsage()
		}

	}

	r.line("")
	r.renderRefs()

	return r.w.String()
}

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

func (r *renderer) line(f string, v ...any) {
	fmt.Fprintf(&r.w, f+"\n", v...)
}

func (r *renderer) paragraph(text string) {
	const width = 80

	for len(text) > width {
		for i := width; i >= 0; i-- {
			if text[i] == ' ' {
				r.w.WriteString(text[:i])
				r.w.WriteByte('\n')
				text = text[i+1:]
				break
			}
		}
	}

	r.w.WriteString(text)
	r.w.WriteByte('\n')

	return

}

func (r *renderer) yaml(v string) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(data))
}
