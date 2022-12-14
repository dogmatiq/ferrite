package ferrite_test

import (
	"fmt"
	"os"
	"time"

	"github.com/dogmatiq/ferrite"
)

func ExampleDuration_required() {
	setUp()
	defer tearDown()

	v := ferrite.
		Duration("FERRITE_DURATION", "example duration variable").
		Required()

	os.Setenv("FERRITE_DURATION", "630s")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is 10m30s
}

func ExampleDuration_default() {
	setUp()
	defer tearDown()

	v := ferrite.
		Duration("FERRITE_DURATION", "example duration variable").
		WithDefault(630 * time.Second).
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is 10m30s
}

func ExampleDuration_optional() {
	setUp()
	defer tearDown()

	v := ferrite.
		Duration("FERRITE_DURATION", "example duration variable").
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

func ExampleDuration_limits() {
	setUp()
	defer tearDown()

	v := ferrite.
		Duration("FERRITE_DURATION", "example duration variable").
		WithMinimum(5 * time.Minute).
		WithMaximum(10 * time.Minute).
		Required()

	os.Setenv("FERRITE_DURATION", "5m")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is 5m0s
}
