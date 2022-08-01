package ferrite_test

import (
	"fmt"
	"os"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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

var _ = Describe("type EnumSpec", func() {
	var (
		reg  *Registry
		spec *EnumSpec[enumMember]
	)

	BeforeEach(func() {
		reg = &Registry{}

		spec = Enum[enumMember](
			"FERRITE_ENUM",
			"<desc>",
			WithRegistry(reg),
		).
			Members(
				member0,
				member1,
				member2,
			)
	})

	AfterEach(func() {
		Teardown()
	})

	Describe("func Value()", func() {
		DescribeTable(
			"it returns the value associated with the member key",
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

	Describe("func Members()", func() {
		When("one of the members has an empty string representation", func() {
			It("panics", func() {
				Expect(func() {
					Enum[string](
						"FERRITE_ENUM",
						"<desc>",
					).Members("")
				}).To(PanicWith("enum member must not have an empty string representation"))
			})
		})
	})

	Describe("func Default()", func() {
		When("the default value is not one of the enum members", func() {
			It("panics", func() {
				Expect(func() {
					spec.Default(enumMember(100))
				}).To(PanicWith("the default value must be one of the enum members"))
			})
		})
	})

	Describe("func Validate()", func() {
		When("the value is not one of the member keys", func() {
			It("returns an error", func() {
				os.Setenv("FERRITE_ENUM", "<invalid>")

				res := spec.Validate()
				Expect(res.Error).To(MatchError(`must be one of "<member-0>", "<member-1>" or "<member-2>"`))
			})
		})

		When("the variable is not defined", func() {
			It("returns an error", func() {
				res := spec.Validate()
				Expect(res.Error).To(MatchError("must not be empty"))
			})
		})
	})

	When("there is a default value", func() {
		BeforeEach(func() {
			spec = spec.Default(member1)
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
	})
})
