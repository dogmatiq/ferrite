package ferrite_test

import (
	"io"
	"os"
	"strings"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func Setup() {
	SetExitBehavior(os.Stdout, func(code int) {})
}

func Teardown() {
	SetExitBehavior(os.Stderr, os.Exit)
	DefaultRegistry.Reset()
	unsetTestVariables()
}

func unsetTestVariables() {
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "FERRITE_") {
			i := strings.Index(env, "=")
			os.Unsetenv(env[:i])
		}
	}
}

var _ = Describe("func ValidateEnvironment()", func() {
	AfterEach(func() {
		Teardown()
	})

	It("validates the default registry", func() {
		v := Bool("FERRITE_REG", "<desc>")

		os.Setenv("FERRITE_REG", "true")
		ValidateEnvironment()

		Expect(v.Value()).To(BeTrue())
	})

	It("writes a report and exits if the registry can not be validated", func() {
		called := false
		// w := gbytes.NewBuffer()

		SetExitBehavior(
			io.Discard,
			func(code int) {
				Expect(code).To(Equal(1))
				called = true
			},
		)

		Bool("FERRITE_REG", "<desc>")

		ValidateEnvironment()
		// w.Close()

		// Expect(w).To(gbytes.Say(`x`))
		Expect(called).To(BeTrue())
	})
})
