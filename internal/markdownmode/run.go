package markdownmode

import (
	"github.com/dogmatiq/ferrite/variable"
)

// Run generates environment variable usage instructions in markdown format.
func Run(app string, reg *variable.Registry) string {
	r := renderer{
		App:   app,
		Specs: reg.Specs(),
	}

	return r.Render()
}

// func renderMarkdownSpecLink(w *strings.Builder, s spec.Spec) {
// 	fmt.Fprintf(w, "- [`%s`](#%s) — %s\n", s.Name, s.Name, s.Description)
// }

// func renderMarkdownSpec(w *strings.Builder, s spec.Spec) {
// 	fmt.Fprintf(w, "### `%s`\n", s.Name)
// 	fmt.Fprintln(w, "")
// 	fmt.Fprintf(w, "> %s\n", s.Description)
// 	fmt.Fprintln(w, "")

// 	if s.IsOptional {

// 	} else if s.HasDefault {
// 		fmt.Fprintln(w, "This variable is **required**, although a default is provided.")
// 	}

// 	s.Schema.AcceptVisitor(&markdownSchemaRenderer{w})
// }

// type markdownSchemaRenderer struct {
// 	Output *strings.Builder
// }

// func (r *markdownSchemaRenderer) VisitOneOf(s spec.OneOf) {
// 	for i, c := range s {
// 		if i > 0 {
// 			r.Output.WriteString(" | ")
// 		}

// 		c.AcceptVisitor(r)
// 	}
// }

// func (r *markdownSchemaRenderer) VisitLiteral(s spec.Literal) {
// 	fmt.Fprintf(r.Output, "%s", s)
// }

// func (r *markdownSchemaRenderer) VisitType(s spec.Type) {
// 	fmt.Fprintf(r.Output, "<%s>", s.Type)
// }

//	func (r *markdownSchemaRenderer) VisitRange(s spec.Range) {
//		if s.Min != "" && s.Max != "" {
//			fmt.Fprintf(r.Output, "%s .. %s", s.Min, s.Max)
//		} else if s.Max != "" {
//			fmt.Fprintf(r.Output, "... %s", s.Max)
//		} else {
//			fmt.Fprintf(r.Output, "%s ...", s.Min)
//		}

// }
