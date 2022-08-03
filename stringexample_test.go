package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleString() {
	setUp()
	defer tearDown()

	value := ferrite.
		String("FERRITE_STRING", "example string variable")

	os.Setenv("FERRITE_STRING", "<value>")
	ferrite.ValidateEnvironment()

	fmt.Println("value is", value.Value())

	// Output:
	// value is <value>
}

func ExampleString_default() {
	setUp()
	defer tearDown()

	value := ferrite.
		String("FERRITE_STRING", "example string variable").
		WithDefault("<default>")

	ferrite.ValidateEnvironment()

	fmt.Println("value is", value.Value())

	// Output:
	// value is <default>
}
