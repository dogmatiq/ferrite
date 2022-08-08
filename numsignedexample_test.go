package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleSigned_required() {
	setUp()
	defer tearDown()

	v := ferrite.
		Signed[int]("FERRITE_SIGNED", "example signed integer variable").
		Required()

	os.Setenv("FERRITE_SIGNED", "-123")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is -123
}

func ExampleSigned_default() {
	setUp()
	defer tearDown()

	v := ferrite.
		Signed[int]("FERRITE_SIGNED", "example signed integer variable").
		WithDefault(-123).
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is -123
}

func ExampleSigned_optional() {
	setUp()
	defer tearDown()

	v := ferrite.
		Signed[int]("FERRITE_SIGNED", "example signed integer variable").
		Optional()

	ferrite.Init()

	if x, ok := v.Value(); ok {
		fmt.Println("value is", x)
	} else {
		fmt.Println("value is undefined")
	}

	// Output:
	// value is undefined
}
