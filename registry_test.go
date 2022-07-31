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
	BeforeEach(func() {
		DefaultRegistry.Reset()
	})

	It("validates the default registry", func() {
		os.Setenv("FERRITE_TEST", "true")
		defer os.Unsetenv("FERRITE_TEST")

		v := Bool("FERRITE_TEST", "<desc>")
		ValidateEnvironment()
		Expect(v.Value()).To(BeTrue())
	})

	It("panics if the registry can not be validated", func() {
		Bool("FERRITE_TEST", "<desc>")

		Expect(func() {
			ValidateEnvironment()
		}).To(Panic())
	})
})

func expectErr(err error, expect ...string) {
	ExpectWithOffset(1, err).Should(HaveOccurred())

	actual := strings.Split(err.Error(), "\n")
	ExpectWithOffset(1, actual).To(EqualX(expect))
}
