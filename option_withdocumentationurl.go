package ferrite

import (
	"net/url"

	"github.com/dogmatiq/ferrite/internal/variable"
)

// WithDocumentationURL is a [RegistryOption] that adds a link to some
// supplementary documentation.
func WithDocumentationURL(u string) RegistryOption {
	x, err := url.Parse(u)
	if err != nil {
		panic("invalid URL: " + err.Error())
	}

	return option{
		ApplyToRegistry: func(r *variable.Registry) {
			r.URL = x
		},
	}
}
