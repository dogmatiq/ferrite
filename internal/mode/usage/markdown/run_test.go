package markdown_test

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/dogmatiq/ferrite"
	"github.com/dogmatiq/ferrite/internal/diff"
	"github.com/dogmatiq/ferrite/internal/environment"
	"github.com/dogmatiq/ferrite/internal/mode"
	. "github.com/dogmatiq/ferrite/internal/mode/usage/markdown"
	"github.com/dogmatiq/ferrite/internal/variable"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func tableTest(
	path string,
	options ...Option,
) func(
	file string,
	setup func(ferrite.Registry),
) {
	return func(
		file string,
		setup func(ferrite.Registry),
	) {
		reg := &variable.Registry{
			IsDefault: true,
		}

		snapshot := environment.TakeSnapshot()
		defer environment.RestoreSnapshot(snapshot)

		setup(reg)

		expect, err := os.ReadFile(
			filepath.Join(
				"testdata",
				path,
				file,
			),
		)
		Expect(err).ShouldNot(HaveOccurred())

		actual := &bytes.Buffer{}
		exited := false

		cfg := mode.Config{
			Args: []string{"<app>"},
			Out:  actual,
			Exit: func(code int) {
				exited = true
				Expect(code).To(Equal(0))
			},
		}
		cfg.Registries.Add(reg)

		Run(cfg, options...)

		if d := diff.Diff(
			"actual",
			actual.Bytes(),
			"expect",
			expect,
		); d != nil {
			ginkgo.Fail(
				"Unexpected markdown content:\n"+string(d),
				1,
			)
		}

		Expect(exited).To(BeTrue())
	}
}
