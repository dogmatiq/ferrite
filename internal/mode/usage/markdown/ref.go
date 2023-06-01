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
		"ferrite": "https://github.com/dogmatiq/ferrite",
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
			n, ok := cutPrefix(k, "rfc ")
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

func cutPrefix(s, prefix string) (after string, found bool) {
	if !strings.HasPrefix(s, prefix) {
		return s, false
	}
	return s[len(prefix):], true
}
