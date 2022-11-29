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
						`value of FERRITE_DURATION (630) is invalid: missing unit`,
					),
					Entry(
						"unknown units",
						"630q",
						`value of FERRITE_DURATION (630q) is invalid: unknown unit "q"`,
					),
					Entry(
						"less than the minimum",
						"0s",
						`value of FERRITE_DURATION (0s) is invalid: too low, expected 1ns or greater`,
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
						`value of FERRITE_DURATION (630) is invalid: missing unit`,
					),
					Entry(
						"unknown units",
						"630q",
						`value of FERRITE_DURATION (630q) is invalid: unknown unit "q"`,
					),
					Entry(
						"less than the minimum",
						"0s",
						`value of FERRITE_DURATION (0s) is invalid: too low, expected 1ns or greater`,
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

	When("the value is lower than the minimum limit", func() {
		It("panics", func() {
			Expect(func() {
				os.Setenv("FERRITE_DURATION", "1s")

				builder.
					WithMinimum(5 * time.Second).
					Required().
					Value()
			}).To(PanicWith(
				`value of FERRITE_DURATION (1s) is invalid: too low, expected 5s or greater`,
			))
		})
	})

	When("the value is greater than the maximum limit", func() {
		It("panics", func() {
			Expect(func() {
				os.Setenv("FERRITE_DURATION", "10s")

				builder.
					WithMaximum(5 * time.Second).
					Required().
					Value()
			}).To(PanicWith(
				`value of FERRITE_DURATION (10s) is invalid: too high, expected between 1ns and 5s`,
			))
		})
	})
})
