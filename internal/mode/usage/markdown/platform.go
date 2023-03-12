package markdown

import (
	"fmt"
	"strings"

	"github.com/dogmatiq/ferrite/variable"
)

func (r *renderer) renderPlatformExamples() {
	r.line("## Usage Examples")
	r.gap()
	r.renderKubernetesExample()
	r.gap()
	r.renderDockerExample()
}

func (r *renderer) renderKubernetesExample() {
	r.line("<details>")
	r.line("<summary>Kubernetes</summary>")

	r.gap()
	r.line("This example shows how to define the environment variables needed by `%s`", r.App)
	r.line("on a %s within a Kubenetes deployment manifest.", r.link("Kubernetes container"))
	r.gap()

	r.line("```yaml")
	r.line("apiVersion: apps/v1")
	r.line("kind: Deployment")
	r.line("metadata:")
	r.line("  name: example-deployment")
	r.line("spec:")
	r.line("  template:")
	r.line("    spec:")
	r.line("      containers:")
	r.line("        - name: example-container")
	r.line("          env:")

	for _, s := range r.Specs {
		eg := variable.BestExample(s)

		r.line("            - name: %s %s", r.yaml(s.Name()), renderDescriptionAsYAMLComment(s))
		r.line("              value: %s", r.yaml(eg.Canonical.String))
	}

	r.line("```")

	r.gap()
	r.line("Alternatively, the environment variables can be defined within a %s", r.link("config map", "kubernetes config map"))
	r.line("then referenced a deployment manifest using `configMapRef`.")
	r.gap()

	r.line("```yaml")
	r.line("apiVersion: v1")
	r.line("kind: ConfigMap")
	r.line("metadata:")
	r.line("  name: example-config-map")
	r.line("data:")

	for _, s := range r.Specs {
		eg := variable.BestExample(s)

		r.line(
			"  %s: %s %s",
			r.yaml(s.Name()),
			r.yaml(eg.Canonical.String),
			renderDescriptionAsYAMLComment(s),
		)
	}

	r.line("---")
	r.line("apiVersion: apps/v1")
	r.line("kind: Deployment")
	r.line("metadata:")
	r.line("  name: example-deployment")
	r.line("spec:")
	r.line("  template:")
	r.line("    spec:")
	r.line("      containers:")
	r.line("        - name: example-container")
	r.line("          envFrom:")
	r.line("            - configMapRef:")
	r.line("                name: example-config-map")
	r.line("```")

	r.gap()
	r.line("</details>")
}

func (r *renderer) renderDockerExample() {
	r.line("<details>")
	r.line("<summary>Docker</summary>")
	r.gap()
	r.line("This example shows how to define the environment variables needed by `%s`", r.App)
	r.line("when running as a %s defined in a Docker compose file.", r.link("Docker service"))
	r.gap()

	r.line("```yaml")
	r.line("service:")
	r.line("  example-service:")
	r.line("    environment:")

	for _, s := range r.Specs {
		eg := variable.BestExample(s)

		r.line(
			"      %s: %s %s",
			r.yaml(s.Name()),
			r.yaml(eg.Canonical.String),
			renderDescriptionAsYAMLComment(s),
		)
	}

	r.line("```")

	r.gap()
	r.line("</details>")
}

func renderDescriptionAsYAMLComment(s variable.Spec) string {
	var w strings.Builder

	w.WriteString("# ")
	w.WriteString(s.Description())

	if s.IsDeprecated() {
		w.WriteString(" (deprecated)")
	} else if def, ok := s.Default(); ok {
		fmt.Fprintf(&w, " (defaults to %s)", def.Quote())
	} else if !s.IsRequired() {
		w.WriteString(" (optional)")
	}

	return w.String()
}
