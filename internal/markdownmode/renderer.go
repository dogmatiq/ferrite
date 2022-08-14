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

	w    strings.Builder
	refs map[string]string
}

func (r *renderer) line(f string, v ...any) {
	fmt.Fprintf(&r.w, f+"\n", v...)
}

func (r *renderer) Render() string {
	r.refs = map[string]string{
		"dogmatiq/ferrite": "https://github.com/dogmatiq/ferrite",
	}

	r.line("# Environment Variables")
	r.line("")
	r.line("This document describes the environment variables used by `%s`.", r.App)

	if len(r.Specs) == 0 {
		r.line("")
		r.line("**There do not appear to be any environment variables.**")
	}

	r.line("")
	r.line("The application may consume other undocumented environment variables; this")
	r.line("document only shows those variables defined using [dogmatiq/ferrite].")

	if len(r.Specs) != 0 {
		r.line("")
		r.renderIndex()
		r.line("")
		r.renderSpecs()
	}

	r.line("")
	r.renderRefs()

	return r.w.String()
}

func (r *renderer) yaml(v any) {
	r.line("```yaml")

	enc := yaml.NewEncoder(&r.w)
	enc.SetIndent(2)
	defer enc.Close()

	if err := enc.Encode(v); err != nil {
		panic(err)
	}

	r.line("```")
}
