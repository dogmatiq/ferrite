package markdown

func (r *specRenderer) renderRegistry() {
	if r.reg.IsDefault {
		return
	}

	r.ren.gap()

	link := r.reg.Name
	if r.reg.URL != nil {
		link = r.ren.linkToURL(
			r.reg.Name,
			r.reg.URL.String(),
			"registry:"+r.reg.Key,
		)
	}

	r.ren.line("This variable is imported from %s.", link)
}
