package ferrite_test

import (
	"github.com/dogmatiq/ferrite"
)

func Example_seeAlso() {
	defer example()()

	verbose := ferrite.
		Bool("VERBOSE", "enable verbose logging").
		Optional()

	ferrite.
		Bool("DEBUG", "enable or disable debugging features").
		Optional(ferrite.SeeAlso(verbose))
}
