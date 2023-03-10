package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleEnum_required() {
	defer example()()

	v := ferrite.
		Enum("FERRITE_ENUM", "example enum variable").
		WithMembers("red", "green", "blue").
		Required()

	os.Setenv("FERRITE_ENUM", "red")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is red
}

func ExampleEnum_default() {
	defer example()()

	v := ferrite.
		Enum("FERRITE_ENUM", "example enum variable").
		WithMembers("red", "green", "blue").
		WithDefault("green").
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is green
}

func ExampleEnum_optional() {
	defer example()()

	v := ferrite.
		Enum("FERRITE_ENUM", "example enum variable").
		WithMembers("red", "green", "blue").
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

func ExampleEnum_descriptions() {
	defer example()()

	v := ferrite.
		Enum("FERRITE_ENUM", "example enum variable").
		WithMember("red", "the color red").
		WithMember("green", "the color green").
		WithMember("blue", "the color blue").
		Required()

	os.Setenv("FERRITE_ENUM", "red")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is red
}

func ExampleEnum_deprecated() {
	defer example()()

	os.Setenv("FERRITE_ENUM", "red")
	ferrite.
		Enum("FERRITE_ENUM", "example enum variable").
		WithMember("red", "the color red").
		WithMember("green", "the color green").
		WithMember("blue", "the color blue").
		Deprecated()

	ferrite.Init()

	fmt.Println("<execution continues>")

	// Output:
	// Environment Variables:
	//
	//  ❯ FERRITE_ENUM  example enum variable  [ red | green | blue ]  ⚠ deprecated variable set to red
	//
	// <execution continues>
}
