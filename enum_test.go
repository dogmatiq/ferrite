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
			"FERRITE_TEST",
			"<desc>",
			WithRegistry(reg),
		).Members(
			member0,
			member1,
			member2,
		)
	})

	Describe("func Value()", func() {
		DescribeTable(
			"it returns the value associated with the member key",
			func(value string, expect enumMember) {
				os.Setenv("FERRITE_TEST", value)
				defer os.Unsetenv("FERRITE_TEST")

				err := reg.Validate()
				Expect(err).ShouldNot(HaveOccurred())
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
						"FERRITE_TEST",
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
				}).To(PanicWith("default value must be one of the enum members"))
			})
		})
	})

	Describe("func Validate()", func() {
		When("the value is not one of the member keys", func() {
			It("returns an error", func() {
				os.Setenv("FERRITE_TEST", "<invalid>")
				defer os.Unsetenv("FERRITE_TEST")

				expectErr(
					reg.Validate(),
					`ENVIRONMENT VARIABLES`,
					` ✗ FERRITE_TEST [enumMember enum] (<desc>)`,
					`   ✓ must be set explicitly`,
					`   ✗ must be one of the enum members, got "<invalid>"`,
					`     • <member-0>`,
					`     • <member-1>`,
					`     • <member-2>`,
				)
			})
		})

		When("the variable is not defined", func() {
			It("returns an error", func() {
				expectErr(
					reg.Validate(),
					`ENVIRONMENT VARIABLES`,
					` ✗ FERRITE_TEST [enumMember enum] (<desc>)`,
					`   ✗ must be set explicitly`,
					`   • must be one of the enum members`,
					`     • <member-0>`,
					`     • <member-1>`,
					`     • <member-2>`,
				)
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
					err := reg.Validate()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(spec.Value()).To(Equal(member1))
				})
			})
		})
	})
})
