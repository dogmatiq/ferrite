package ferrite_test

import (
	"os"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("boolean values", func() {
	BeforeEach(func() {
		DefaultRegistry.Reset()
	})

	When("the value is not a valid literal", func() {
		It("panics", func() {
			os.Setenv("FERRITE_TEST", "<invalid>")
			defer os.Unsetenv("FERRITE_TEST")

			Bool("FERRITE_TEST", "<desc>")

			Expect(func() {
				ResolveEnvironment()
			}).To(PanicWith(
				`ENVIRONMENT VARIABLES
	✗ FERRITE_TEST [bool] (<desc>)
		✓ must be set explicitly
		✗ must be either "true" or "false", got "<invalid>"`,
			))
		})
	})

	When("the value is not a valid custom literal", func() {
		It("panics", func() {
			os.Setenv("FERRITE_TEST", "true")
			defer os.Unsetenv("FERRITE_TEST")

			Bool("FERRITE_TEST", "<desc>").
				Literals("yes", "no")

			Expect(func() {
				ResolveEnvironment()
			}).To(PanicWith(
				`ENVIRONMENT VARIABLES
	✗ FERRITE_TEST [bool] (<desc>)
		✓ must be set explicitly
		✗ must be either "yes" or "no", got "true"`,
			))
		})
	})
})
