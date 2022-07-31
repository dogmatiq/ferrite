package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleString() {
	ferrite.DefaultRegistry.Reset()
	os.Setenv("FERRITE_STRING", "<value>")
	defer os.Unsetenv("FERRITE_STRING")

	value := ferrite.
		String(
			"FERRITE_STRING",
			"example string variable",
		)

	ferrite.ValidateEnvironment()

	fmt.Println("value is", value.Value())

	// Output:
	// value is <value>
}
