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

func (r *renderer) paragraphf(text ...string) func(...any) {
	return func(v ...any) {
		text := fmt.Sprintf(strings.Join(text, " "), v...)

		r.gap()
		for _, line := range wordwrap.Wrap(text, 80) {
			r.line("%s", line)
		}
	}
}

func (r *renderer) paragraph(
	fn func(
		func(string, ...any),
	),
) {
	var w strings.Builder

	fn(func(f string, v ...any) {
		fmt.Fprintf(&w, f, v...)
	})

	r.paragraphf("%s")(w.String())
}

func (r *renderer) yaml(v string) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(data))
}

func andList[T any](
	parts []T,
	format func(T) string,
) string {
	return inlineList(parts, format, ", ", " and ")
}

func orList[T any](
	parts []T,
	format func(T) string,
) string {
	return inlineList(parts, format, ", ", " or ")
}

func inlineList[T any](
	parts []T,
	format func(T) string,
	sep, lastSep string,
) string {
	w := &strings.Builder{}

	for i, s := range parts {
		if i > 0 {
			if i == len(parts)-1 {
				w.WriteString(lastSep)
			} else {
				w.WriteString(sep)
			}
		}

		w.WriteString(format(s))
	}

	return w.String()
}
