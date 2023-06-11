package ferrite_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite"
	"github.com/dogmatiq/ferrite/internal/mode"
	"github.com/dogmatiq/ferrite/internal/variable"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// example is a helper function that sets up the global state for a testable
// example. It returns a function that resets the global state after the test.
func example() func() {
	mode.DefaultConfig.Out = os.Stdout
	mode.DefaultConfig.Err = os.Stdout
	mode.DefaultConfig.Exit = func(code int) {
		if code == 0 {
			fmt.Fprintln(os.Stdout, "<process exited successfully>")
		} else {
			fmt.Fprintf(os.Stdout, "<process exited with error code %d>\n", code)
		}
	}

	return tearDown
}

// tearDown resets the environemnt and Ferrite global state after a test.
func tearDown() {
	mode.ResetDefaultConfig()
	variable.ResetDefaultRegistry()

	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "FERRITE_") {
			i := strings.Index(env, "=")
			os.Unsetenv(env[:i])
		}
	}
}

var _ = Describe("func Init()", func() {
	AfterEach(func() {
		tearDown()
	})

	It("exits with a non-zero status code if there is an invalid environment variable", func() {
		String("FERRITE_STRING", "required").
			Required()

		exited := false
		mode.DefaultConfig.Exit = func(code int) {
			Expect(code).NotTo(Equal(0))
			exited = true
		}

		ferrite.Init()
		Expect(exited).To(BeTrue())
	})
})
