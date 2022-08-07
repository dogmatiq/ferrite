package ferrite_test

import (
	"os"
	"strings"

	. "github.com/dogmatiq/ferrite"
	. "github.com/jmalloc/gomegax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("func ValidateEnvironment()", func() {
	AfterEach(func() {
		tearDown()
	})

	It("invokes the registered validators", func() {
		v := String("FERRITE_REG", "<desc>")

		os.Setenv("FERRITE_REG", "<value>")
		ValidateEnvironment()

		Expect(v.Value()).To(Equal("<value>"))
	})

	It("writes a report and exits if a validator fails", func() {
		stderr := &strings.Builder{}
		called := false

		SetExitBehavior(
			stderr,
			func(code int) {
				Expect(code).To(Equal(1))
				called = true
			},
		)

		String("FERRITE_REG", "<desc>")

		ValidateEnvironment()

		Expect(called).To(BeTrue())
		expectLines(
			stderr.String(),
			`ENVIRONMENT VARIABLES:`,
			` ❯ FERRITE_REG    [string]  <desc>  ✗ must not be empty`,
		)
	})
})

// expectLines verifies that text consists of the given lines.
func expectLines(actual string, lines ...string) {
	expect := strings.Join(lines, "\n") + "\n"
	ExpectWithOffset(1, actual).To(EqualX(expect))
}
