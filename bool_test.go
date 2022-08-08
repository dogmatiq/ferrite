package ferrite_test

import (
	"os"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type userDefinedBool bool

func (v userDefinedBool) String() string {
	if v {
		return "yes"
	}
	return "no"
}

var _ = Describe("type BoolBuilder", func() {
	var builder BoolBuilder[userDefinedBool]

	BeforeEach(func() {
		builder = BoolAs[userDefinedBool]("FERRITE_BOOL", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	When("the variable is required", func() {
		When("the value is one of the accepted literals", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"it returns the value associated with the literal",
					func(value string, expect userDefinedBool) {
						os.Setenv("FERRITE_BOOL", value)

						v := builder.
							Required().
							Value()

						Expect(v).To(Equal(expect))
					},
					Entry("true", "yes", userDefinedBool(true)),
					Entry("false", "no", userDefinedBool(false)),
				)
			})
		})

		When("the value is invalid", func() {
			BeforeEach(func() {
				// we don't accept true/false for the userDefinedBool type
				os.Setenv("FERRITE_BOOL", "true")
			})

			Describe("func Value()", func() {
				It("panics", func() {
					Expect(func() {
						builder.
							Required().
							Value()
					}).To(PanicWith(
						`FERRITE_BOOL (true) is invalid: must be either yes or no`,
					))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					DescribeTable(
						"it returns the default",
						func(expect userDefinedBool) {
							v := builder.
								WithDefault(expect).
								Required().
								Value()

							Expect(v).To(Equal(expect))
						},
						Entry("true", userDefinedBool(true)),
						Entry("false", userDefinedBool(false)),
					)
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
							"FERRITE_BOOL is undefined and does not have a default value",
						))
					})
				})
			})
		})
	})

	When("the variable is optional", func() {
		When("the value is invalid", func() {
			BeforeEach(func() {
				// we don't accept true/false for the userDefinedBool type
				os.Setenv("FERRITE_BOOL", "true")
			})

			Describe("func Value()", func() {
				It("panics", func() {
					Expect(func() {
						builder.
							Optional().
							Value()
					}).To(PanicWith(
						`FERRITE_BOOL (true) is invalid: must be either yes or no`,
					))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					DescribeTable(
						"it returns the default",
						func(expect userDefinedBool) {
							v, ok := builder.
								WithDefault(expect).
								Optional().
								Value()

							Expect(ok).To(BeTrue())
							Expect(v).To(Equal(expect))
						},
						Entry("true", userDefinedBool(true)),
						Entry("false", userDefinedBool(false)),
					)
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

	Describe("func WithLiterals()", func() {
		DescribeTable(
			"it overrides the default literals",
			func(value string, expect userDefinedBool) {
				os.Setenv("FERRITE_BOOL", value)

				v := builder.
					WithLiterals("true", "false").
					Required().
					Value()

				Expect(v).To(Equal(expect))
			},
			Entry("true", "true", userDefinedBool(true)),
			Entry("false", "false", userDefinedBool(false)),
		)

		When("the true literal is empty", func() {
			It("panics", func() {
				Expect(func() {
					builder.WithLiterals("", "no")
				}).To(PanicWith(
					"specification for FERRITE_BOOL is invalid: boolean literals must not be zero-length",
				))
			})
		})

		When("the true literal is empty", func() {
			It("panics", func() {
				Expect(func() {
					builder.WithLiterals("yes", "")
				}).To(PanicWith(
					"specification for FERRITE_BOOL is invalid: boolean literals must not be zero-length",
				))
			})
		})
	})
})
