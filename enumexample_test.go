package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleEnum() {
	ferrite.Setup()
	defer ferrite.Teardown()

	value := ferrite.
		Enum[string]("FERRITE_ENUM", "example enum variable").
		WithMembers("red", "green", "blue")

	os.Setenv("FERRITE_ENUM", "red")
	ferrite.ValidateEnvironment()

	fmt.Println("value is", value.Value())

	// Output:
	// value is red
}

func ExampleEnum_default() {
	ferrite.Setup()
	defer ferrite.Teardown()

	value := ferrite.
		Enum[string]("FERRITE_ENUM", "example enum variable").
		WithMembers("red", "green", "blue").
		WithDefault("green")

	ferrite.ValidateEnvironment()

	fmt.Println("value is", value.Value())

	// Output:
	// value is green
}
