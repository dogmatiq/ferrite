package markdownmode

import (
	"github.com/dogmatiq/ferrite/variable"
)

func (r *renderer) renderUsageExamples(s variable.Spec, v variable.Literal) {
	r.line("<details>")
	r.line("<summary>Usage Examples</summary>")

	r.line("")
	r.renderK8sContainerExample(s, v)

	r.line("")
	r.renderK8sConfigMapExample(s, v)

	r.line("")
	r.renderDockerExample(s, v)

	r.line("")
	r.renderGHAExample(s, v)

	r.line("")
	r.line("</details>")
}

func (r *renderer) renderK8sContainerExample(s variable.Spec, v variable.Literal) {
	r.line("#### Kubernetes Container")
	r.line("")
	r.yaml(
		map[any]any{
			"env": []map[any]any{
				{
					"name":  s.Name(),
					"value": v.String,
				},
			},
		},
	)
}

func (r *renderer) renderK8sConfigMapExample(s variable.Spec, v variable.Literal) {
	r.line("#### Kubernetes Config Map")
	r.line("")
	r.yaml(
		map[string]any{
			"data": map[any]any{
				s.Name(): v.String,
			},
		},
	)
}

func (r *renderer) renderDockerExample(s variable.Spec, v variable.Literal) {
	r.line("#### Docker Service")
	r.line("")
}

func (r *renderer) renderGHAExample(s variable.Spec, v variable.Literal) {
	r.line("#### GitHub Actions Workflow")
	r.line("")
}

// #### Kubernetes Container

// #### Kubernetes Config Map

// ```yaml
// data:
//   DEBUG: "true"
// ```

// #### Docker Service

// ```yaml
// environment:
//   DEBUG: "true"
// ```

// #### GitHub Actions Workflow

// ```yaml
// env:
//   DEBUG: "true"
// ```
