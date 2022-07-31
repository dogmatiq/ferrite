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
		spec *BoolSpec[bool]
	)

	BeforeEach(func() {
		reg = &Registry{}

		spec = Bool(
			"FERRITE_TEST",
			"<desc>",
			WithRegistry(reg),
		)
	})

	Describe("func Value()", func() {
		DescribeTable(
			"it returns the value associated with the literal",
			func(value string, expect bool) {
				os.Setenv("FERRITE_TEST", value)
				defer os.Unsetenv("FERRITE_TEST")

				err := reg.Validate()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(spec.Value()).To(Equal(expect))
			},
			Entry("true", "true", true),
			Entry("false", "false", false),
		)
	})

	Describe("func Validate()", func() {
		When("the value is not a valid literal", func() {
			It("returns an error", func() {
				os.Setenv("FERRITE_TEST", "<invalid>")
				defer os.Unsetenv("FERRITE_TEST")

				expectErr(
					reg.Validate(),
					`ENVIRONMENT VARIABLES`,
					` ✗ FERRITE_TEST [bool] (<desc>)`,
					`   ✓ must be set explicitly`,
					`   ✗ must be either "true" or "false", got "<invalid>"`,
				)
			})
		})

		When("the variable is not defined", func() {
			It("returns an error", func() {
				expectErr(
					reg.Validate(),
					`ENVIRONMENT VARIABLES`,
					` ✗ FERRITE_TEST [bool] (<desc>)`,
					`   ✗ must be set explicitly`,
					`   - must be either "true" or "false"`,
				)
			})
		})
	})

	When("there is a default value", func() {
		Describe("func Value()", func() {
			When("the variable is not defined", func() {
				DescribeTable(
					"it returns the default",
					func(expect bool) {
						spec.Default(expect)

						err := reg.Validate()
						Expect(err).ShouldNot(HaveOccurred())
						Expect(spec.Value()).To(Equal(expect))
					},
					Entry("true", true),
					Entry("false", false),
				)
			})
		})
	})

	When("there are custom literals", func() {
		BeforeEach(func() {
			spec.Literals("yes", "no")
		})

		Describe("func Value()", func() {
			DescribeTable(
				"it returns the value associated with the literal",
				func(value string, expect bool) {
					os.Setenv("FERRITE_TEST", value)
					defer os.Unsetenv("FERRITE_TEST")

					err := reg.Validate()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(spec.Value()).To(Equal(expect))
				},
				Entry("true", "yes", true),
				Entry("false", "no", false),
			)
		})

		Describe("func Validate()", func() {
			When("the value is not a valid literal", func() {
				It("returns an error", func() {
					os.Setenv("FERRITE_TEST", "true")
					defer os.Unsetenv("FERRITE_TEST")

					expectErr(
						reg.Validate(),
						`ENVIRONMENT VARIABLES`,
						` ✗ FERRITE_TEST [bool] (<desc>)`,
						`   ✓ must be set explicitly`,
						`   ✗ must be either "yes" or "no", got "true"`,
					)
				})
			})
		})
	})
})
