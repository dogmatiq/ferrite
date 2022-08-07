package ferrite_test

// import (
// 	"os"

// 	"github.com/dogmatiq/ferrite"
// )

// func ExampleValidateEnvironment_failure() {
// 	setUp()
// 	defer tearDown()

// 	ferrite.
// 		String("FERRITE_UNDEFINED", "example undefined string")

// 	ferrite.
// 		String("FERRITE_DEFAULTABLE", "string with a default").
// 		WithDefault("the default")

// 	ferrite.
// 		String("FERRITE_EXPLICIT", "example explicit string")

// 	os.Setenv("FERRITE_EXPLICIT", "explicit value")
// 	ferrite.ValidateEnvironment()

// 	// Note, the report shows environment variables in alphabetical order, not
// 	// the order in which they were defined.

// 	// Output:
// 	// ENVIRONMENT VARIABLES:
// 	//    FERRITE_DEFAULTABLE  [string] = "the default"  string with a default     ✓ using default value
// 	//    FERRITE_EXPLICIT     [string]                  example explicit string   ✓ set to "explicit value"
// 	//  ❯ FERRITE_UNDEFINED    [string]                  example undefined string  ✗ must not be empty
// }
