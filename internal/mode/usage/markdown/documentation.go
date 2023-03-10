package markdown

func (r *specRenderer) renderImportantDocumentation() {
	for _, d := range r.spec.Documentation() {
		if d.IsImportant {
			for _, p := range d.Paragraphs {
				r.ren.paragraph(p)()
			}
		}
	}
}

func (r *specRenderer) renderUnimportantDocumentation() {
	for _, d := range r.spec.Documentation() {
		if d.IsImportant {
			continue
		}

		r.ren.gap()
		r.ren.line("<details>")

		if d.Summary != "" {
			r.ren.line("<summary>%s</summary>", d.Summary)
		}

		for _, p := range d.Paragraphs {
			r.ren.paragraph(p)()
		}

		r.ren.gap()
		r.ren.line("</details>")
	}
}
