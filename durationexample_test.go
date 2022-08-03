package ferrite_test

import (
	"fmt"
	"os"
	"time"

	"github.com/dogmatiq/ferrite"
)

func ExampleDuration() {
	setUp()
	defer tearDown()

	value := ferrite.
		Duration("FERRITE_DURATION", "example duration variable")

	os.Setenv("FERRITE_DURATION", "630s")
	ferrite.ValidateEnvironment()

	fmt.Println("value is", value.Value())

	// Output:
	// value is 10m30s
}

func ExampleDuration_default() {
	setUp()
	defer tearDown()

	value := ferrite.
		Duration("FERRITE_DURATION", "example duration variable").
		WithDefault(630 * time.Second)

	ferrite.ValidateEnvironment()

	fmt.Println("value is", value.Value())

	// Output:
	// value is 10m30s
}
