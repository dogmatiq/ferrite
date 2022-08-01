package ferrite_test

import (
	"errors"
	"os"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type StringSpec", func() {
	type customString string

	var spec *StringSpec[customString]

	BeforeEach(func() {
		spec = StringAs[customString]("FERRITE_STRING", "<desc>")
	})

	AfterEach(func() {
		Teardown()
	})

	When("the environment variable is not empty", func() {
		BeforeEach(func() {
			os.Setenv("FERRITE_STRING", "<value>")
		})

		Describe("func Value()", func() {
			It("returns the raw string value", func() {
				res := spec.Validate()
				Expect(res.Error).ShouldNot(HaveOccurred())
				Expect(spec.Value()).To(Equal(customString("<value>")))
			})
		})

		Describe("func Validate()", func() {
			It("returns a successful result", func() {
				Expect(spec.Validate()).To(Equal(
					VariableValidationResult{
						Name:          "FERRITE_STRING",
						Description:   "<desc>",
						ValidInput:    "[ferrite_test.customString]",
						DefaultValue:  "",
						ExplicitValue: `"<value>"`,
						Error:         nil,
					},
				))
			})
		})
	})

	When("the environment variable is empty", func() {
		When("there is a default value", func() {
			BeforeEach(func() {
				spec.WithDefault("<value>")
			})

			Describe("func Value()", func() {
				It("returns the default", func() {
					res := spec.Validate()
					Expect(res.Error).ShouldNot(HaveOccurred())
					Expect(spec.Value()).To(Equal(customString("<value>")))
				})
			})

			Describe("func Validate()", func() {
				It("returns a success result", func() {
					Expect(spec.Validate()).To(Equal(
						VariableValidationResult{
							Name:          "FERRITE_STRING",
							Description:   "<desc>",
							ValidInput:    "[ferrite_test.customString]",
							DefaultValue:  `"<value>"`,
							ExplicitValue: `""`,
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
					Expect(spec.Validate()).To(Equal(
						VariableValidationResult{
							Name:          "FERRITE_STRING",
							Description:   "<desc>",
							ValidInput:    "[ferrite_test.customString]",
							DefaultValue:  "",
							ExplicitValue: `""`,
							UsingDefault:  false,
							Error:         errors.New(`must not be empty`),
						},
					))
				})
			})
		})
	})
})
