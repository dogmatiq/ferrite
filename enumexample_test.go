package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleEnum() {
	ferrite.DefaultRegistry.Reset()
	os.Setenv("FERRITE_ENUM", "red")
	defer os.Unsetenv("FERRITE_ENUM")

	value := ferrite.Enum[string](
		"FERRITE_ENUM",
		"example enum variable",
	).
		Members("red", "green", "blue")

	ferrite.ValidateEnvironment()

	fmt.Println("value is", value.Value())

	// Output:
	// value is red
}

func ExampleEnum_default() {
	ferrite.DefaultRegistry.Reset()

	value := ferrite.Enum[string](
		"FERRITE_ENUM",
		"example enum variable",
	).
		Members("red", "green", "blue").
		Default("green")

	ferrite.ValidateEnvironment()

	fmt.Println("value is", value.Value())

	// Output:
	// value is green
}
