package ferrite_test

import (
	"errors"
	"os"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type BoolSpec", func() {
	type customBool bool

	var spec *BoolSpec[customBool]

	BeforeEach(func() {
		spec = BoolAs[customBool]("FERRITE_BOOL", "<desc>")
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

					res := spec.Validate("FERRITE_BOOL")
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

				Expect(spec.Validate("FERRITE_BOOL")).To(Equal(
					ValidationResult{
						Name:          "FERRITE_BOOL",
						Description:   "<desc>",
						ValidInput:    "true|false",
						DefaultValue:  "",
						ExplicitValue: "true",
						UsingDefault:  false,
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
						spec.WithDefault(expect)

						res := spec.Validate("FERRITE_BOOL")
						Expect(res.Error).ShouldNot(HaveOccurred())
						Expect(spec.Value()).To(Equal(expect))
					},
					Entry("true", customBool(true)),
					Entry("false", customBool(false)),
				)
			})

			Describe("func Validate()", func() {
				It("returns a success result", func() {
					spec.WithDefault(true)

					Expect(spec.Validate("FERRITE_BOOL")).To(Equal(
						ValidationResult{
							Name:          "FERRITE_BOOL",
							Description:   "<desc>",
							ValidInput:    "true|false",
							DefaultValue:  "true",
							ExplicitValue: "",
							UsingDefault:  true,
							Error:         nil,
						},
					))
				})
			})
		})

		When("there is no default value", func() {
			Describe("func Validate()", func() {
				It("returns a failure result", func() {
					Expect(spec.Validate("FERRITE_BOOL")).To(Equal(
						ValidationResult{
							Name:          "FERRITE_BOOL",
							Description:   "<desc>",
							ValidInput:    "true|false",
							DefaultValue:  "",
							ExplicitValue: "",
							UsingDefault:  false,
							Error:         errors.New(`must not be empty`),
						},
					))
				})
			})
		})
	})

	When("the environment variable is set to some other value", func() {
		Describe("func Validate()", func() {
			It("returns a failure result", func() {
				os.Setenv("FERRITE_BOOL", "<invalid>")

				Expect(spec.Validate("FERRITE_BOOL")).To(Equal(
					ValidationResult{
						Name:          "FERRITE_BOOL",
						Description:   "<desc>",
						ValidInput:    "true|false",
						DefaultValue:  "",
						ExplicitValue: "<invalid>",
						UsingDefault:  false,
						Error:         errors.New(`must be either "true" or "false"`),
					},
				))
			})
		})
	})

	When("there are custom literals", func() {
		BeforeEach(func() {
			spec.WithLiterals("yes", "no")
		})

		When("the environment variable is set to one of the custom literals", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"it returns the value associated with the literal",
					func(value string, expect customBool) {
						os.Setenv("FERRITE_BOOL", value)

						res := spec.Validate("FERRITE_BOOL")
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
					"it returns a failure result",
					func(value string) {
						os.Setenv("FERRITE_BOOL", value)

						Expect(spec.Validate("FERRITE_BOOL")).To(Equal(
							ValidationResult{
								Name:          "FERRITE_BOOL",
								Description:   "<desc>",
								ValidInput:    "yes|no",
								DefaultValue:  "",
								ExplicitValue: value,
								UsingDefault:  false,
								Error:         errors.New(`must be either "yes" or "no"`),
							},
						))
					},
					Entry("true", "true"),
					Entry("false", "false"),
				)
			})
		})
	})

	Describe("func WithLiterals()", func() {
		When("the true literal is empty", func() {
			It("panics", func() {
				Expect(func() {
					spec.WithLiterals("", "no")
				}).To(PanicWith("boolean literals must not be zero-length"))
			})
		})

		When("the true literal is empty", func() {
			It("panics", func() {
				Expect(func() {
					spec.WithLiterals("yes", "")
				}).To(PanicWith("boolean literals must not be zero-length"))
			})
		})
	})
})
