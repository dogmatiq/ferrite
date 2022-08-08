package ferrite_test

import (
	"os"
	"time"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type DurationBuilder", func() {
	var builder DurationBuilder

	BeforeEach(func() {
		builder = Duration("FERRITE_DURATION", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	When("the variable is required", func() {
		When("the value is a valid duration", func() {
			Describe("func Value()", func() {
				It("returns the value", func() {
					os.Setenv("FERRITE_DURATION", "630s")

					v := builder.
						Required().
						Value()

					Expect(v).To(Equal(630 * time.Second))
				})
			})
		})

		When("the value is invalid", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"it panics",
					func(value, expect string) {
						os.Setenv("FERRITE_DURATION", value)

						Expect(func() {
							builder.
								Required().
								Value()
						}).To(PanicWith(expect))
					},
					Entry(
						"missing units",
						"630",
						`FERRITE_DURATION ("630") is invalid: time: missing unit in duration "630"`,
					),
					Entry(
						"less than the minimum",
						"0s",
						`FERRITE_DURATION ("0s") is invalid: must be 1ns or greater`,
					),
				)
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						expect := 10*time.Minute + 30*time.Second

						v := builder.
							WithDefault(expect).
							Required().
							Value()

						Expect(v).To(Equal(expect))
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
							"FERRITE_DURATION is undefined and does not have a default value",
						))
					})
				})
			})
		})
	})

	When("the variable is optional", func() {
		When("the value is a valid duration", func() {
			Describe("func Value()", func() {
				It("returns the value", func() {
					os.Setenv("FERRITE_DURATION", "630s")

					v, ok := builder.
						Optional().
						Value()

					Expect(ok).To(BeTrue())
					Expect(v).To(Equal(630 * time.Second))
				})
			})
		})

		When("the value is invalid", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"it panics",
					func(value, expect string) {
						os.Setenv("FERRITE_DURATION", value)

						Expect(func() {
							builder.
								Optional().
								Value()
						}).To(PanicWith(expect))
					},
					Entry(
						"missing units",
						"630",
						`FERRITE_DURATION ("630") is invalid: time: missing unit in duration "630"`,
					),
					Entry(
						"less than the minimum",
						"0s",
						`FERRITE_DURATION ("0s") is invalid: must be 1ns or greater`,
					),
				)
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						expect := 10*time.Minute + 30*time.Second

						v, ok := builder.
							WithDefault(expect).
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(expect))
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
