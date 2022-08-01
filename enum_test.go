package ferrite_test

import (
	"errors"
	"fmt"
	"os"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type EnumSpec", func() {
	var spec *EnumSpec[enumMember]

	BeforeEach(func() {
		spec = Enum[enumMember]("FERRITE_ENUM", "<desc>").
			WithMembers(
				member0,
				member1,
				member2,
			)
	})

	AfterEach(func() {
		Teardown()
	})

	When("the environment variable is set to one of the enum members", func() {
		Describe("func Value()", func() {
			DescribeTable(
				"it returns the value associated with the member",
				func(value string, expect enumMember) {
					os.Setenv("FERRITE_ENUM", value)

					res := spec.Validate()
					Expect(res.Error).ShouldNot(HaveOccurred())
					Expect(spec.Value()).To(Equal(expect))
				},
				Entry("member 0", "<member-0>", member0),
				Entry("member 1", "<member-1>", member1),
				Entry("member 2", "<member-2>", member2),
			)
		})

		Describe("func Validate()", func() {
			It("returns a successful result", func() {
				os.Setenv("FERRITE_ENUM", "<member-1>")

				Expect(spec.Validate()).To(Equal(
					VariableValidationResult{
						Name:          "FERRITE_ENUM",
						Description:   "<desc>",
						ValidInput:    "<member-0>|<member-1>|<member-2>",
						DefaultValue:  "",
						ExplicitValue: "<member-1>",
						Error:         nil,
					},
				))
			})
		})
	})

	When("the environment variable is empty", func() {
		When("there is a default value", func() {
			BeforeEach(func() {
				spec = spec.WithDefault(member1)
			})

			Describe("func Value()", func() {
				When("the variable is not defined", func() {
					It("returns the default", func() {
						res := spec.Validate()
						Expect(res.Error).ShouldNot(HaveOccurred())
						Expect(spec.Value()).To(Equal(member1))
					})
				})
			})

			Describe("func Validate()", func() {
				It("returns a success result", func() {
					Expect(spec.Validate()).To(Equal(
						VariableValidationResult{
							Name:          "FERRITE_ENUM",
							Description:   "<desc>",
							ValidInput:    "<member-0>|<member-1>|<member-2>",
							DefaultValue:  "<member-1>",
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
					Expect(spec.Validate()).To(Equal(
						VariableValidationResult{
							Name:          "FERRITE_ENUM",
							Description:   "<desc>",
							ValidInput:    "<member-0>|<member-1>|<member-2>",
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
			It("returns an failure result", func() {
				os.Setenv("FERRITE_ENUM", "<invalid>")

				Expect(spec.Validate()).To(Equal(
					VariableValidationResult{
						Name:          "FERRITE_ENUM",
						Description:   "<desc>",
						ValidInput:    "<member-0>|<member-1>|<member-2>",
						DefaultValue:  "",
						ExplicitValue: "<invalid>",
						UsingDefault:  false,
						Error:         errors.New(`<invalid> is not a member of the enum`),
					},
				))
			})
		})
	})

	Describe("func WithMembers()", func() {
		When("one of the members has an empty string representation", func() {
			It("panics", func() {
				Expect(func() {
					Enum[string](
						"FERRITE_ENUM",
						"<desc>",
					).WithMembers("")
				}).To(PanicWith("enum member must not have an empty string representation"))
			})
		})
	})

	Describe("func WithDefault()", func() {
		When("the default value is not a member of the enum", func() {
			It("panics", func() {
				Expect(func() {
					spec.WithDefault(enumMember(100))
				}).To(PanicWith("the default value must be one of the enum members"))
			})
		})
	})
})

// enumMember is used to test enumerations.
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
