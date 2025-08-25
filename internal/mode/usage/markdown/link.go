package markdown

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dogmatiq/ferrite/internal/variable"
)

// linkToURL returns the markdown that links to the given URL.
func (r *renderer) linkToURL(text, url string, optionalRef ...string) string {
	var ref string
	switch len(optionalRef) {
	case 0:
		ref = strings.ToLower(text)
	case 1:
		ref = strings.ToLower(optionalRef[0])
	default:
		panic("too many arguments")
	}

	if r.links == nil {
		r.links = map[string]string{ref: url}
	} else if existing, ok := r.links[ref]; !ok {
		r.links[ref] = url
	} else if existing != url {
		panic(fmt.Sprintf(
			"link ref %q refers to both %q and %q",
			ref, existing,
			url,
		))
	}

	if strings.EqualFold(text, ref) {
		return fmt.Sprintf("[%s]", text)
	}

	return fmt.Sprintf("[%s](%s)", text, ref)
}

// linkToSpec returns markdown that links to the given variable specification.
func (r *renderer) linkToSpec(s variable.Spec) string {
	return r.linkToURL(
		fmt.Sprintf("`%s`", s.Name()),
		fmt.Sprintf("#%s", s.Name()),
	)
}

// linkToRegistry returns markdown that links to the given registry.
//
// If the registry does not have a documentation URL, it returns the registry
// name without any link.
func (r *renderer) linkToRegistry(reg *variable.Registry) string {
	if reg.URL == nil {
		return reg.Name
	}

	return r.linkToURL(
		reg.Name,
		reg.URL.String(),
		"registry:"+reg.Key,
	)
}

// renderLinkRefs renders the map of reference names to URLs.
func (r *renderer) renderLinkRefs() {
	if len(r.links) == 0 {
		return
	}

	r.gap()
	r.line("<!-- references -->")
	r.gap()

	refs := make([]string, 0, len(r.links))
	for ref := range r.links {
		refs = append(refs, ref)
	}

	sort.Slice(
		refs,
		func(i, j int) bool {
			a := strings.Trim(refs[i], "`")
			b := strings.Trim(refs[j], "`")
			return a < b
		},
	)

	for _, ref := range refs {
		r.line("[%s]: %s", ref, r.links[ref])
	}
}
