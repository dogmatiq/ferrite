package ferrite_test

import (
	"os"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type StringSpec", func() {
	type customString string

	var (
		reg  *Registry
		spec *StringSpec[customString]
	)

	BeforeEach(func() {
		reg = &Registry{}

		spec = StringAs[customString](
			"FERRITE_TEST",
			"<desc>",
			WithRegistry(reg),
		)
	})

	Describe("func Value()", func() {
		It("returns the value", func() {
			os.Setenv("FERRITE_TEST", "<value>")
			defer os.Unsetenv("FERRITE_TEST")

			err := reg.Validate()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(spec.Value()).To(Equal(customString("<value>")))
		})
	})

	Describe("func Validate()", func() {
		When("the variable is not defined", func() {
			It("returns an error", func() {
				expectErr(
					reg.Validate(),
					`ENVIRONMENT VARIABLES`,
					` ✗ FERRITE_TEST [customString] (<desc>)`,
					`   ✗ must be set explicitly`,
				)
			})
		})
	})

	When("there is a default value", func() {
		Describe("func Value()", func() {
			When("the variable is not defined", func() {
				It("returns the default", func() {
					expect := customString("<value>")
					spec.Default(expect)

					err := reg.Validate()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(spec.Value()).To(Equal(expect))
				})
			})
		})
	})
})
