package markdownmode

import (
	"fmt"
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

	urls := map[string]string{
		"docker service":        "https://docs.docker.com/compose/environment-variables/#set-environment-variables-in-containers",
		"ferrite":               "https://github.com/dogmatiq/ferrite",
		"kubernetes config map": "https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#configure-all-key-value-pairs-in-a-configmap-as-container-environment-variables",
		"kubernetes container":  "https://kubernetes.io/docs/tasks/inject-data-application/define-environment-variable-container/#define-an-environment-variable-for-a-container",
	}

	keys := maps.Keys(r.refs)
	slices.SortFunc(
		keys,
		func(a, b string) bool {
			return strings.Trim(a, "`") < strings.Trim(b, "`")
		},
	)

	for _, k := range keys {
		r.line("[%s]: %s", k, urls[k])
	}
}

func (r *renderer) link(text string, ref ...string) string {
	if r.refs == nil {
		r.refs = map[string]struct{}{}
	}

	switch len(ref) {
	case 0:
		ref := strings.ToLower(text)
		r.refs[ref] = struct{}{}
		return fmt.Sprintf("[%s]", text)
	case 1:
		ref := strings.ToLower(ref[0])
		r.refs[ref] = struct{}{}
		return fmt.Sprintf("[%s][%s]", text, ref)
	default:
		panic("too many arguments")
	}
}
