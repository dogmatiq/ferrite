package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleNetworkPort_required() {
	setUp()
	defer tearDown()

	v := ferrite.
		NetworkPort("FERRITE_NETWORK_PORT", "example network port").
		Required()

	os.Setenv("FERRITE_NETWORK_PORT", "https")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is https
}

func ExampleNetworkPort_default() {
	setUp()
	defer tearDown()

	v := ferrite.
		NetworkPort("FERRITE_NETWORK_PORT", "example network port").
		WithDefault("12345").
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is 12345
}

func ExampleNetworkPort_optional() {
	setUp()
	defer tearDown()

	v := ferrite.
		NetworkPort("FERRITE_NETWORK_PORT", "example network port").
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
