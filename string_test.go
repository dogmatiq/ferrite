package ferrite_test

import (
	"errors"
	"os"

	. "github.com/dogmatiq/ferrite"
	"github.com/dogmatiq/ferrite/schema"
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
		tearDown()
	})

	Describe("func Describe()", func() {
		It("describes the variable", func() {
			Expect(spec.Describe()).To(ConsistOf(
				VariableXXX{
					Name:        "FERRITE_STRING",
					Description: "<desc>",
					Schema:      schema.Type[customString](),
				},
			))
		})
	})

	When("the environment variable is not empty", func() {
		BeforeEach(func() {
			os.Setenv("FERRITE_STRING", "<value>")
		})

		Describe("func Value()", func() {
			It("returns the raw string value", func() {
				Expect(spec.Value()).To(Equal(customString("<value>")))
			})
		})

		Describe("func Validate()", func() {
			It("returns a successful result", func() {
				Expect(spec.Validate()).To(ConsistOf(
					ValidationResult{
						Name:  "FERRITE_STRING",
						Value: `<value>`,
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
					Expect(spec.Value()).To(Equal(customString("<value>")))
				})
			})

			Describe("func Describe()", func() {
				It("describes the variable", func() {
					Expect(spec.Describe()).To(ConsistOf(
						VariableXXX{
							Name:        "FERRITE_STRING",
							Description: "<desc>",
							Schema:      schema.Type[customString](),
							Default:     "<value>",
						},
					))
				})
			})

			Describe("func Validate()", func() {
				It("returns a success result", func() {
					Expect(spec.Validate()).To(ConsistOf(
						ValidationResult{
							Name:        "FERRITE_STRING",
							Value:       "<value>",
							UsedDefault: true,
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
					}).To(PanicWith("FERRITE_STRING is invalid: must not be empty"))
				})
			})

			Describe("func Validate()", func() {
				It("returns a failure result", func() {
					Expect(spec.Validate()).To(ConsistOf(
						ValidationResult{
							Name:  "FERRITE_STRING",
							Error: errors.New(`must not be empty`),
						},
					))
				})
			})
		})
	})
})
