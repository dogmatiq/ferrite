package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleString_required() {
	setUp()
	defer tearDown()

	v := ferrite.
		String("FERRITE_STRING", "example string variable").
		Required()

	os.Setenv("FERRITE_STRING", "<value>")
	ferrite.ValidateEnvironment()

	fmt.Println("value is", v.Value())

	// Output:
	// value is <value>
}

func ExampleString_default() {
	setUp()
	defer tearDown()

	v := ferrite.
		String("FERRITE_STRING", "example string variable").
		WithDefault("<default>").
		Required()

	ferrite.ValidateEnvironment()

	fmt.Println("value is", v.Value())

	// Output:
	// value is <default>
}

func ExampleString_optional() {
	setUp()
	defer tearDown()

	v := ferrite.
		String("FERRITE_STRING", "example string variable").
		Optional()

	ferrite.ValidateEnvironment()

	if x, ok := v.Value(); ok {
		fmt.Println("value is", x)
	} else {
		fmt.Println("value is undefined")
	}

	// Output:
	// value is undefined
}
