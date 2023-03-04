package markdownmode

func (r *renderer) renderPreamble() {
	r.line("This document describes the environment variables used by `%s`.", r.App)
	r.line("")

	if len(r.Specs) == 0 {
		r.line("**There do not appear to be any environment variables.**")
	} else {
		r.line("If any of the environment variable values do not meet the requirements herein,")
		r.line("the application will print usage information to `STDERR` then exit with a")
		r.line("non-zero exit code. Please note that **undefined** variables and **empty**")
		r.line("values are considered equivalent.")

		if r.hasNonNormativeExamples() {
			r.line("")
			r.line("⚠️ This document includes **non-normative** example values. While these values")
			r.line("are syntactically correct, they may not be meaningful to this application.")
		}
	}

	r.line("")
	r.line("⚠️ The application may consume other undocumented environment variables; this")
	r.line("document only shows those variables declared using %s.", r.link("Ferrite"))

	if len(r.Specs) != 0 {
		r.line("")
		r.line("The key words **MUST**, **MUST NOT**, **REQUIRED**, **SHALL**, **SHALL NOT**,")
		r.line("**SHOULD**, **SHOULD NOT**, **RECOMMENDED**, **MAY**, and **OPTIONAL** in this")
		r.line("document are to be interpreted as described in %s.", r.link("RFC 2119"))
	}
}
