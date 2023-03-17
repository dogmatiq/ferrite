package markdown

import (
	"fmt"
	"strings"

	"github.com/dogmatiq/ferrite/variable"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func (r *renderer) renderRefs() {
	r.line("<!-- references -->")
	r.gap()

	keys := maps.Keys(r.refs)
	urls := map[string]string{
		"docker service":        "https://docs.docker.com/compose/environment-variables/#set-environment-variables-in-containers",
		"ferrite":               "https://github.com/dogmatiq/ferrite",
		"kubernetes config map": "https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#configure-all-key-value-pairs-in-a-configmap-as-container-environment-variables",
		"kubernetes container":  "https://kubernetes.io/docs/tasks/inject-data-application/define-environment-variable-container/#define-an-environment-variable-for-a-container",
	}

	for _, s := range r.Specs {
		key := "`" + strings.ToLower(s.Name()) + "`"
		if _, ok := r.refs[key]; ok {
			urls[key] = "#" + s.Name()
		}
	}

	slices.SortFunc(
		keys,
		func(a, b string) bool {
			return strings.Trim(a, "`") < strings.Trim(b, "`")
		},
	)

	for _, k := range keys {
		u, ok := urls[k]
		if !ok {
			n, ok := strings.CutPrefix(k, "rfc ")
			if !ok {
				panic("unknown reference: " + k)
			}

			u = fmt.Sprintf("https://www.rfc-editor.org/rfc/rfc%s.html", n)
		}

		r.line("[%s]: %s", k, u)
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

func (r *renderer) linkToSpec(s variable.Spec) string {
	if r.refs == nil {
		r.refs = map[string]struct{}{}
	}

	key := fmt.Sprintf(
		"`%s`",
		strings.ToLower(s.Name()),
	)
	r.refs[key] = struct{}{}

	return fmt.Sprintf("[`%s`]", s.Name())
}
