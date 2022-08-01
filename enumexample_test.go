package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleEnum() {
	Setup()
	defer Teardown()

	value := ferrite.
		Enum[string]("FERRITE_ENUM", "example enum variable").
		Members("red", "green", "blue")

	os.Setenv("FERRITE_ENUM", "red")
	ferrite.ValidateEnvironment()

	fmt.Println("value is", value.Value())

	// Output:
	// value is red
}

func ExampleEnum_default() {
	Setup()
	defer Teardown()

	value := ferrite.
		Enum[string]("FERRITE_ENUM", "example enum variable").
		Members("red", "green", "blue").
		Default("green")

	ferrite.ValidateEnvironment()

	fmt.Println("value is", value.Value())

	// Output:
	// value is green
}
