package ferrite_test

import (
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleValidateEnvironment_failure() {
	Setup()
	defer Teardown()

	ferrite.
		String("FERRITE_STRING", "example string")

	ferrite.
		String("FERRITE_DEFAULTABLE", "string with a default").
		Default("the default")

	ferrite.
		Bool("FERRITE_BOOL", "example boolean")

	os.Setenv("FERRITE_BOOL", "false")
	ferrite.ValidateEnvironment()

	// Note, the report shows environment variables in alphabetical order, not
	// the order in which they were defined.

	// Output:
	// ENVIRONMENT VARIABLES:
	//    FERRITE_BOOL         true|false                example boolean        ✓ set to false
	//    FERRITE_DEFAULTABLE  [string] = "the default"  string with a default  ✓ using default value
	//  ❯ FERRITE_STRING       [string]                  example string         ✗ must not be empty
}
