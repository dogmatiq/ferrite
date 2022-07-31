package ferrite_test

import (
	"os"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type BoolSpec", func() {
	var (
		reg  *Registry
		spec *BoolSpec
	)

	BeforeEach(func() {
		reg = &Registry{}

		spec = Bool(
			"FERRITE_TEST",
			"<desc>",
			WithRegistry(reg),
		)
	})

	Describe("func Resolve()", func() {
		When("the value is not a valid literal", func() {
			It("causes an error", func() {
				os.Setenv("FERRITE_TEST", "<invalid>")
				defer os.Unsetenv("FERRITE_TEST")

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
					os.Setenv("FERRITE_TEST", "true")
					defer os.Unsetenv("FERRITE_TEST")

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

	Describe("func Value()", func() {
		When("the spec has not been resolved", func() {
			It("panics", func() {
				Expect(func() {
					spec.Value()
				}).To(PanicWith("environment has not been resolved"))
			})
		})
	})
})
