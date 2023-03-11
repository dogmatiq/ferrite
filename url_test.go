package ferrite_test

import (
	"os"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type URLBuilder", func() {
	var builder *URLBuilder

	BeforeEach(func() {
		builder = URL("FERRITE_URL", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	It("panics if the name is empty", func() {
		Expect(func() {
			URL("", "<desc>").Optional()
		}).To(PanicWith("invalid specification: variable name must not be empty"))
	})

	It("panics if the description is empty", func() {
		Expect(func() {
			URL("FERRITE_URL", "").Optional()
		}).To(PanicWith("specification for FERRITE_URL is invalid: variable description must not be empty"))
	})

	When("the variable is required", func() {
		When("the value is not empty", func() {
			Describe("func Value()", func() {
				It("returns the value ", func() {
					os.Setenv("FERRITE_URL", "https://example.org/path")

					v := builder.
						Required().
						Value()

					Expect(v.String()).To(Equal("https://example.org/path"))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v := builder.
							WithDefault("https://example.org/path").
							Required().
							Value()

						Expect(v.String()).To(Equal("https://example.org/path"))
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
							"FERRITE_URL is undefined and does not have a default value",
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
					os.Setenv("FERRITE_URL", "https://example.org/path")

					v, ok := builder.
						Optional().
						Value()

					Expect(ok).To(BeTrue())
					Expect(v.String()).To(Equal("https://example.org/path"))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v, ok := builder.
							WithDefault("https://example.org/path").
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v.String()).To(Equal("https://example.org/path"))
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
