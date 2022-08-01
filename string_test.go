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
			"FERRITE_STRING",
			"<desc>",
			WithRegistry(reg),
		)
	})

	AfterEach(func() {
		Teardown()
	})

	Describe("func Value()", func() {
		It("returns the value", func() {
			os.Setenv("FERRITE_STRING", "<value>")

			err := reg.Validate(nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(spec.Value()).To(Equal(customString("<value>")))
		})
	})

	Describe("func Validate()", func() {
		When("the variable is not defined", func() {
			It("returns an error", func() {
				err := reg.Validate(nil)
				Expect(err).To(MatchError("must be defined and not empty"))
			})
		})
	})

	When("there is a default value", func() {
		Describe("func Value()", func() {
			When("the variable is not defined", func() {
				It("returns the default", func() {
					expect := customString("<value>")
					spec.Default(expect)

					err := reg.Validate(nil)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(spec.Value()).To(Equal(expect))
				})
			})
		})
	})
})
