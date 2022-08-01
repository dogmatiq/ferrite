package ferrite_test

import (
	"os"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type BoolSpec", func() {
	type customBool bool

	var (
		reg  *Registry
		spec *BoolSpec[customBool]
	)

	BeforeEach(func() {
		reg = &Registry{}

		spec = BoolAs[customBool](
			"FERRITE_BOOL",
			"<desc>",
			WithRegistry(reg),
		)
	})

	AfterEach(func() {
		Teardown()
	})

	When("the environment variable is set to one of the standard literals", func() {
		Describe("func Value()", func() {
			DescribeTable(
				"it returns the value associated with the literal",
				func(value string, expect customBool) {
					os.Setenv("FERRITE_BOOL", value)

					res := spec.Validate()
					Expect(res.Error).ShouldNot(HaveOccurred())
					Expect(spec.Value()).To(Equal(expect))
				},
				Entry("true", "true", customBool(true)),
				Entry("false", "false", customBool(false)),
			)
		})

		Describe("func Validate()", func() {
			It("returns a successful result", func() {
				os.Setenv("FERRITE_BOOL", "true")

				Expect(spec.Validate()).To(Equal(
					VariableValidationResult{
						Name:          "FERRITE_BOOL",
						Description:   "<desc>",
						ValidInput:    "true|false",
						DefaultValue:  "",
						ExplicitValue: "true",
						Error:         nil,
					},
				))
			})
		})
	})

	When("the environment variable is empty", func() {
		When("there is a default value", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"it returns the default",
					func(expect customBool) {
						spec.Default(expect)

						res := spec.Validate()
						Expect(res.Error).ShouldNot(HaveOccurred())
						Expect(spec.Value()).To(Equal(expect))
					},
					Entry("true", customBool(true)),
					Entry("false", customBool(false)),
				)
			})
		})

		When("there is no default value", func() {
			Describe("func Validate()", func() {
				It("it returns an error", func() {
					res := spec.Validate()
					Expect(res.Error).Should(MatchError(`must not be empty`))
				})
			})
		})
	})

	When("the environment variable is set to some other value", func() {
		Describe("func Validate()", func() {
			It("returns an error", func() {
				os.Setenv("FERRITE_BOOL", "<invalid>")

				res := spec.Validate()
				Expect(res.Error).Should(MatchError(`must be either "true" or "false"`))
			})
		})
	})

	When("there are custom literals", func() {
		BeforeEach(func() {
			spec.Literals("yes", "no")
		})

		When("the environment variable is set to one of the custom literals", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"it returns the value associated with the literal",
					func(value string, expect customBool) {
						os.Setenv("FERRITE_BOOL", value)

						res := spec.Validate()
						Expect(res.Error).ShouldNot(HaveOccurred())
						Expect(spec.Value()).To(Equal(expect))
					},
					Entry("true", "yes", customBool(true)),
					Entry("false", "no", customBool(false)),
				)
			})
		})

		When("the environment variable is set to some other value", func() {
			Describe("func Validate()", func() {
				DescribeTable(
					"it returns an error",
					func(value string) {
						os.Setenv("FERRITE_BOOL", value)

						res := spec.Validate()
						Expect(res.Error).Should(MatchError(`must be either "yes" or "no"`))
					},
					Entry("true", "true"),
					Entry("false", "false"),
				)
			})
		})
	})
})
