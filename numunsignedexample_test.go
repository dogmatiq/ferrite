package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleUnsigned_required() {
	setUp()
	defer tearDown()

	v := ferrite.
		Unsigned[uint]("FERRITE_UNSIGNED", "example unsigned integer variable").
		Required()

	os.Setenv("FERRITE_UNSIGNED", "123")
	ferrite.ValidateEnvironment()

	fmt.Println("value is", v.Value())

	// Output:
	// value is 123
}

func ExampleUnsigned_default() {
	setUp()
	defer tearDown()

	v := ferrite.
		Unsigned[uint]("FERRITE_UNSIGNED", "example unsigned integer variable").
		WithDefault(123).
		Required()

	ferrite.ValidateEnvironment()

	fmt.Println("value is", v.Value())

	// Output:
	// value is 123
}

func ExampleUnsigned_optional() {
	setUp()
	defer tearDown()

	v := ferrite.
		Unsigned[uint]("FERRITE_UNSIGNED", "example unsigned integer variable").
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
