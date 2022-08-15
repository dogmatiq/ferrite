package markdownmode

import (
	"github.com/dogmatiq/ferrite/variable"
)

// Run generates environment variable usage instructions in markdown format.
func Run(app string, reg *variable.Registry, usage bool) string {
	r := renderer{
		App:         app,
		RenderUsage: usage,
		Specs:       reg.Specs(),
	}

	return r.Render()
}
