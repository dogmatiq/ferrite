package markdown

func (r *specRenderer) renderRegistry() {
	if !r.reg.IsDefault {
		r.ren.gap()
		r.ren.line(
			"This variable is imported from %s.",
			r.ren.linkToRegistry(r.reg),
		)
	}
}
