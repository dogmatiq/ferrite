package ferrite_test

import (
	"errors"
	"os"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type SignedSpec", func() {
	type customNumeric int

	var spec *SignedSpec[customNumeric]

	BeforeEach(func() {
		spec = Signed[customNumeric]("FERRITE_SIGNED", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	When("the environment variable is not empty", func() {
		BeforeEach(func() {
			os.Setenv("FERRITE_SIGNED", "-123")
		})

		Describe("func Value()", func() {
			It("returns the numeric value", func() {
				Expect(spec.Value()).To(Equal(customNumeric(-123)))
			})
		})

		Describe("func Validate()", func() {
			It("returns a successful result", func() {
				Expect(spec.Validate()).To(ConsistOf(
					ValidationResult{
						Name:          "FERRITE_SIGNED",
						Description:   "<desc>",
						ValidInput:    "[ferrite_test.customNumeric]",
						ExplicitValue: "-123",
					},
				))
			})
		})
	})

	When("the environment variable is empty", func() {
		When("there is a default value", func() {
			BeforeEach(func() {
				spec.WithDefault(-123)
			})

			Describe("func Value()", func() {
				It("returns the default", func() {
					Expect(spec.Value()).To(Equal(customNumeric(-123)))
				})
			})

			Describe("func Validate()", func() {
				It("returns a success result", func() {
					Expect(spec.Validate()).To(ConsistOf(
						ValidationResult{
							Name:         "FERRITE_SIGNED",
							Description:  "<desc>",
							ValidInput:   "[ferrite_test.customNumeric]",
							DefaultValue: "-123",
							UsingDefault: true,
						},
					))
				})
			})
		})

		When("there is no default value", func() {
			Describe("func Value()", func() {
				It("panics", func() {
					Expect(func() {
						spec.Value()
					}).To(PanicWith("FERRITE_SIGNED is invalid: must not be empty"))
				})
			})

			Describe("func Validate()", func() {
				It("returns a failure result", func() {
					Expect(spec.Validate()).To(ConsistOf(
						ValidationResult{
							Name:        "FERRITE_SIGNED",
							Description: "<desc>",
							ValidInput:  "[ferrite_test.customNumeric]",
							Error:       errors.New(`must not be empty`),
						},
					))
				})
			})
		})
	})
})
