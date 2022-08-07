package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type enumMember int

const (
	member0 enumMember = iota
	member1
	member2
)

// String returns the value to use as the enum's key.
func (m enumMember) String() string {
	return fmt.Sprintf("<member-%d>", m)
}

var _ = Describe("type EnumSpec", func() {
	var builder EnumBuilder[enumMember]

	BeforeEach(func() {
		builder = EnumAs[enumMember]("FERRITE_ENUM", "<desc>").
			WithMembers(
				member0,
				member1,
				member2,
			)
	})

	AfterEach(func() {
		tearDown()
	})

	When("the variable is required", func() {
		When("the value is one of the accepted literals", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"it returns the value associated with the literal",
					func(value string, expect enumMember) {
						os.Setenv("FERRITE_ENUM", value)

						v := builder.
							Required().
							Value()

						Expect(v).To(Equal(expect))
					},
					Entry("member 0", "<member-0>", member0),
					Entry("member 1", "<member-1>", member1),
					Entry("member 2", "<member-2>", member2),
				)
			})
		})

		When("the value is invalid", func() {
			BeforeEach(func() {
				os.Setenv("FERRITE_ENUM", "<non-member>")
			})

			Describe("func Value()", func() {
				It("panics", func() {
					Expect(func() {
						builder.
							Required().
							Value()
					}).To(PanicWith(
						`FERRITE_ENUM must be one of one of the enum members, got "<non-member>"`,
					))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v := builder.
							WithDefault(member1).
							Required().
							Value()

						Expect(v).To(Equal(member1))
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
							"FERRITE_ENUM is undefined and does not have a default value",
						))
					})
				})
			})
		})
	})

	When("the variable is optional", func() {
		When("the value is one of the accepted literals", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"it returns the value associated with the literal",
					func(value string, expect enumMember) {
						os.Setenv("FERRITE_ENUM", value)

						v, ok := builder.
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(expect))
					},
					Entry("member 0", "<member-0>", member0),
					Entry("member 1", "<member-1>", member1),
					Entry("member 2", "<member-2>", member2),
				)
			})
		})

		When("the value is invalid", func() {
			BeforeEach(func() {
				os.Setenv("FERRITE_ENUM", "<non-member>")
			})

			Describe("func Value()", func() {
				It("panics", func() {
					Expect(func() {
						builder.
							Optional().
							Value()
					}).To(PanicWith(
						`FERRITE_ENUM must be one of one of the enum members, got "<non-member>"`,
					))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v, ok := builder.
							WithDefault(member1).
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(member1))
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

	When("one of the members has an empty literal representation", func() {
		It("uses that member as the default", func() {
			v := ferrite.
				EnumAs[string]("FERRITE_ENUM", "<desc>").
				WithMembers("foo", "", "bar").
				Required().
				Value()

			Expect(v).To(Equal(""))
		})

		It("does not take precedence over an explicit default", func() {
			v := ferrite.
				EnumAs[string]("FERRITE_ENUM", "<desc>").
				WithDefault("foo").
				WithMembers("foo", "", "bar").
				Required().
				Value()

			Expect(v).To(Equal("foo"))
		})
	})

	When("multiple members have the same literal representation", func() {
		It("panics", func() {
			Expect(func() {
				builder.
					WithMembers(member1, member1).
					Required()
			}).To(PanicWith(
				`specification for FERRITE_ENUM is invalid: multiple members use "<member-1>" as their literal representation`,
			))
		})
	})

	When("there are no members", func() {
		It("panics", func() {
			Expect(func() {
				ferrite.
					EnumAs[string]("FERRITE_ENUM", "<desc>").
					Required()
			}).To(PanicWith(
				`specification for FERRITE_ENUM is invalid: no enum members are defined`,
			))
		})
	})

	When("the default value is not a member of the enum", func() {
		It("panics", func() {
			Expect(func() {
				builder.
					WithDefault(enumMember(100)).
					Required()
			}).To(PanicWith(
				`specification for FERRITE_ENUM is invalid: the default value must be one of the enum members, got "<member-100>"`,
			))
		})
	})
})
