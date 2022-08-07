package ferrite_test

import (
	"os"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type userDefinedString string

var _ = Describe("type StringSpec", func() {
	var builder StringBuilder[userDefinedString]

	BeforeEach(func() {
		builder = StringAs[userDefinedString]("FERRITE_STRING", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	When("the variable is required", func() {
		When("the value is not empty", func() {
			Describe("func Value()", func() {
				It("returns the value ", func() {
					os.Setenv("FERRITE_STRING", "<value>")

					v := builder.
						Required().
						Value()

					Expect(v).To(Equal(userDefinedString("<value>")))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v := builder.
							WithDefault("<value>").
							Required().
							Value()

						Expect(v).To(Equal(userDefinedString("<value>")))
					})
				})
			})

			When("there is no default value", func() {
				Describe("func Value()", func() {
					It("panics", func() {
						Expect(func() {
							builder.
								Required().
								Value()
						}).To(PanicWith(
							"FERRITE_STRING is undefined and does not have a default value",
						))
					})
				})
			})
		})
	})

	When("the variable is optional", func() {
		When("the value is not empty", func() {
			Describe("func Value()", func() {
				It("returns the value ", func() {
					os.Setenv("FERRITE_STRING", "<value>")

					v, ok := builder.
						Optional().
						Value()

					Expect(ok).To(BeTrue())
					Expect(v).To(Equal(userDefinedString("<value>")))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v, ok := builder.
							WithDefault("<value>").
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(userDefinedString("<value>")))
					})
				})
			})

			When("there is no default value", func() {
				Describe("func Value()", func() {
					It("returns with ok == false", func() {
						_, ok := builder.
							Optional().
							Value()

						Expect(ok).To(BeFalse())
					})
				})
			})
		})
	})
})
