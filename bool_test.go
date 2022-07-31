package ferrite_test

import (
	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Context("boolean values", func() {
	var (
		env  MemoryEnvironment
		reg  *Registry
		spec *BoolSpec
	)

	BeforeEach(func() {
		env = MemoryEnvironment{}

		reg = &Registry{
			Environment: env,
		}

		spec = Bool(
			"FERRITE_TEST",
			"<desc>",
			WithRegistry(reg),
		)
	})

	When("the value is not a valid literal", func() {
		It("causes an error", func() {
			env["FERRITE_TEST"] = "<invalid>"

			expectErr(
				reg.Resolve(),
				`ENVIRONMENT VARIABLES`,
				` ✗ FERRITE_TEST [bool] (<desc>)`,
				`   ✓ must be set explicitly`,
				`   ✗ must be either "true" or "false", got "<invalid>"`,
			)
		})
	})

	When("the value has custom literals", func() {
		BeforeEach(func() {
			spec.Literals("yes", "no")
		})

		When("the value is not a valid literal", func() {
			It("causes an error", func() {
				env["FERRITE_TEST"] = "true"

				expectErr(
					reg.Resolve(),
					`ENVIRONMENT VARIABLES`,
					` ✗ FERRITE_TEST [bool] (<desc>)`,
					`   ✓ must be set explicitly`,
					`   ✗ must be either "yes" or "no", got "true"`,
				)
			})
		})
	})
})
