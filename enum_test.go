package ferrite_test

import (
	"os"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type EnumSpec", func() {
	var (
		reg  *Registry
		spec *EnumSpec[string]
	)

	BeforeEach(func() {
		reg = &Registry{}

		spec = Enum[string](
			"FERRITE_TEST",
			"<desc>",
			WithRegistry(reg),
		).Members(
			"red",
			"green",
			"blue",
		)
	})

	Describe("func Value()", func() {
		DescribeTable(
			"it returns the value associated with the member key",
			func(value string) {
				os.Setenv("FERRITE_TEST", value)
				defer os.Unsetenv("FERRITE_TEST")

				err := reg.Validate()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(spec.Value()).To(Equal(value))
			},
			Entry("red", "red"),
			Entry("green", "green"),
			Entry("blue", "blue"),
		)
	})

	Describe("func Members()", func() {
		When("one of the members has an empty string representation", func() {
			It("panics", func() {
				Expect(func() {
					spec.Members("red", "", "blue")
				}).To(PanicWith("enum member must not have an empty string representation"))
			})
		})
	})

	Describe("func Default()", func() {
		When("the default value is not one of the enum members", func() {
			It("panics", func() {
				Expect(func() {
					spec.Default("<invalid>")
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
					` ✗ FERRITE_TEST [string enum] (<desc>)`,
					`   ✓ must be set explicitly`,
					`   ✗ must be "red", "green" or "blue", got "<invalid>"`,
				)
			})
		})

		When("the variable is not defined", func() {
			It("returns an error", func() {
				expectErr(
					reg.Validate(),
					`ENVIRONMENT VARIABLES`,
					` ✗ FERRITE_TEST [string enum] (<desc>)`,
					`   ✗ must be set explicitly`,
					`   - must be "red", "green" or "blue"`,
				)
			})
		})
	})

	When("there is a default value", func() {
		BeforeEach(func() {
			spec = spec.Default("green")
		})

		Describe("func Value()", func() {
			When("the variable is not defined", func() {
				It("returns the default", func() {
					err := reg.Validate()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(spec.Value()).To(Equal("green"))
				})
			})
		})
	})
})
