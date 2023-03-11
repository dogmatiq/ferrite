package markdown

import (
	"fmt"
	"io"
	"strings"

	"github.com/dogmatiq/ferrite/internal/wordwrap"
	"github.com/dogmatiq/ferrite/variable"
	"gopkg.in/yaml.v3"
)

type renderer struct {
	App    string
	Specs  []variable.Spec
	Output io.Writer

	withoutPreamble      bool
	withoutIndex         bool
	withoutUsageExamples bool
	refs                 map[string]struct{}
}

func (r *renderer) Render() {
	r.line("# Environment Variables")

	if !r.withoutPreamble {
		r.gap()
		r.renderPreamble()
	}

	if len(r.Specs) != 0 {
		if !r.withoutIndex {
			r.gap()
			r.renderIndex()
		}

		r.gap()
		r.line("## Specification")

		for _, s := range r.Specs {
			sr := specRenderer{r, s}
			sr.Render()
		}

		if !r.withoutUsageExamples {
			r.gap()
			r.renderPlatformExamples()
		}
	}

	if len(r.refs) != 0 {
		r.gap()
		r.renderRefs()
	}
}

func (r *renderer) gap() {
	r.line("")
}

func (r *renderer) line(f string, v ...any) {
	if _, err := fmt.Fprintf(r.Output, f+"\n", v...); err != nil {
		panic(err)
	}
}

func (r *renderer) paragraph(text ...string) func(...any) {
	return func(v ...any) {
		text := fmt.Sprintf(strings.Join(text, " "), v...)

		r.gap()
		for _, line := range wordwrap.Wrap(text, 80) {
			r.line("%s", line)
		}
	}
}

func (r *renderer) yaml(v string) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(data))
}
