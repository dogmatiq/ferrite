package markdownmode

import (
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// func renderLink(w *strings.Builder, text, ref string) {
// 	if text == ref || ref == "" {
// 		fmt.Fprintf(w, "[%s]", text)
// 	} else {
// 		fmt.Fprintf(w, "[%s](%s)", text, ref)
// 	}
// }

func (r *renderer) renderRefs() {
	r.line("<!-- references -->")
	r.line("")

	keys := maps.Keys(r.refs)
	slices.SortFunc(
		keys,
		func(a, b string) bool {
			return strings.Trim(a, "`") < strings.Trim(b, "`")
		},
	)

	for _, k := range keys {
		r.line("[%s]: %s", k, r.refs[k])
	}
}
