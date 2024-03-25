package markdown

import (
	"fmt"
	"io"
	"strings"

	"github.com/dogmatiq/ferrite/internal/variable"
	"github.com/dogmatiq/ferrite/internal/wordwrap"
	"gopkg.in/yaml.v3"
)

type renderer struct {
	App       string
	Variables []variable.RegisteredVariable
	Output    io.Writer

	withoutExplanatoryText bool
	withoutIndex           bool
	withoutUsageExamples   bool
	links                  map[string]string
}

func (r *renderer) Render() {
	r.line("# Environment Variables")

	if !r.withoutExplanatoryText {
		r.gap()
		r.line("This document describes the environment variables used by `%s`.", r.App)
	}

	if !r.withoutIndex {
		r.gap()
		r.renderIndex()
	}

	if !r.withoutExplanatoryText {
		r.alertf(
			"WARNING",
			"This document only shows environment variables declared using %s.",
			"`%s` may consume other undocumented environment variables.",
		)(
			r.linkToURL(
				"Ferrite",
				"https://github.com/dogmatiq/ferrite",
			),
			r.App,
		)
	}

	if len(r.Variables) != 0 {
		r.gap()
		r.line("## Specification")

		if !r.withoutExplanatoryText {
			r.paragraphf(
				"All environment variables described below must meet the stated requirements.",
				"Otherwise, `%s` prints usage information to `STDERR` then exits.",
				"**Undefined** variables and **empty** values are equivalent.",
			)(
				r.App,
			)

			if r.hasNonNormativeExamples() {
				r.paragraphf(
					"⚠️ This section includes **non-normative** example values.",
					"These examples are syntactically valid, but may not be meaningful to `%s`.",
				)(
					r.App,
				)
			}

			r.paragraphf(
				"The key words **MUST**, **MUST NOT**, **REQUIRED**, **SHALL**, **SHALL NOT**,",
				"**SHOULD**, **SHOULD NOT**, **RECOMMENDED**, **MAY**, and **OPTIONAL** in this",
				"document are to be interpreted as described in %s.",
			)(
				r.linkToRFC(2119),
			)
		}

		for _, v := range r.Variables {
			sr := specRenderer{
				r,
				v.Spec(),
				v.Registry,
			}
			sr.Render()
		}
	}

	r.renderLinkRefs()
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

func (r *renderer) alertf(t string, text ...string) func(...any) {
	return func(v ...any) {
		text := fmt.Sprintf(strings.Join(text, " "), v...)

		r.gap()
		r.line("> [!%s]", t)
		for _, line := range wordwrap.Wrap(text, 78) {
			r.line("> %s", line)
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
