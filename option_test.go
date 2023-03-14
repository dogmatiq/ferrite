package ferrite_test

import (
	"github.com/dogmatiq/ferrite"
)

func Example_seeAlso() {
	defer example()()

	verbose := ferrite.
		Bool("FERRITE_VERBOSE", "enable verbose logging").
		Optional()

	ferrite.
		Bool("FERRITE_DEBUG", "enable or disable debugging features").
		Optional(ferrite.SeeAlso(verbose))

	// Output:
}

func Example_supersededBy() {
	defer example()()

	verbose := ferrite.
		Bool("FERRITE_VERBOSE", "enable verbose logging").
		Optional()

	ferrite.
		Bool("FERRITE_DEBUG", "enable debug logging").
		Deprecated(ferrite.SupersededBy(verbose))

	// Output:
}
