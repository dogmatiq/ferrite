package markdownmode

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dogmatiq/ferrite/variable"
	"gopkg.in/yaml.v3"
)

type renderer struct {
	App   string
	Specs []variable.Spec

	w                    strings.Builder
	withoutPreamble      bool
	withoutIndex         bool
	withoutUsageExamples bool
	refs                 map[string]struct{}
}

func (r *renderer) Render() string {
	r.line("# Environment Variables")

	if !r.withoutPreamble {
		r.line("")
		r.renderPreamble()
	}

	if len(r.Specs) != 0 {
		if !r.withoutIndex {
			r.line("")
			r.renderIndex()
		}

		r.line("")
		r.renderSpecs()

		if !r.withoutUsageExamples {
			r.line("")
			r.renderUsage()
		}
	}

	if len(r.refs) != 0 {
		r.line("")
		r.renderRefs()
	}

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

func (r *renderer) paragraph(text ...string) {
	const width = 80

	t := strings.Join(text, " ")

	for len(t) > width {
		for i := width; i >= 0; i-- {
			if t[i] == ' ' {
				r.w.WriteString(t[:i])
				r.w.WriteByte('\n')
				t = t[i+1:]
				break
			}
		}
	}

	r.w.WriteString(t)
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

// rfc2119Keywords is a pattern that matches the keywords defined in RFC 2199.
//
// See https://www.rfc-editor.org/rfc/rfc2119.
var rfc2119Keywords = regexp.MustCompile(`\b((?:MUST|REQUIRED|SHALL|SHOULD|RECOMMENDED|MAY)(?:\s+NOT)?)\b`)

// asMarkdown converts plain text to Markdown.
func asMarkdown(text string) string {
	return rfc2119Keywords.ReplaceAllString(text, "**$1**")
}
