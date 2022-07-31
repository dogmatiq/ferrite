package ferrite_test

import (
	"strings"

	. "github.com/jmalloc/gomegax"
	. "github.com/onsi/gomega"
)

func expectErr(err error, expect ...string) {
	ExpectWithOffset(1, err).Should(HaveOccurred())

	actual := strings.Split(err.Error(), "\n")
	ExpectWithOffset(1, actual).To(EqualX(expect))
}
