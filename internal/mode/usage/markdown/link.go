package markdown

import (
	"fmt"
	"strings"

	"github.com/dogmatiq/ferrite/internal/variable"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
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

// linkToRFC returns the markdown that links to the given RFC number.
func (r *renderer) linkToRFC(number int) string {
	n := fmt.Sprintf("%04d", number)

	return r.linkToURL(
		"RFC "+n,
		fmt.Sprintf("https://www.rfc-editor.org/rfc/rfc%s.html", n),
	)
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

	refs := maps.Keys(r.links)

	slices.SortFunc(
		refs,
		func(a, b string) int {
			a = strings.Trim(a, "`")
			b = strings.Trim(b, "`")
			return strings.Compare(a, b)
		},
	)

	for _, ref := range refs {
		r.line("[%s]: %s", ref, r.links[ref])
	}
}
