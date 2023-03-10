package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleNetworkPort_required() {
	defer example()()

	v := ferrite.
		NetworkPort("FERRITE_NETWORK_PORT", "example network port variable").
		Required()

	os.Setenv("FERRITE_NETWORK_PORT", "https")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is https
}

func ExampleNetworkPort_default() {
	defer example()()

	v := ferrite.
		NetworkPort("FERRITE_NETWORK_PORT", "example network port variable").
		WithDefault("12345").
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is 12345
}

func ExampleNetworkPort_optional() {
	defer example()()

	v := ferrite.
		NetworkPort("FERRITE_NETWORK_PORT", "example network port variable").
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

func ExampleNetworkPort_deprecated() {
	defer example()()

	os.Setenv("FERRITE_NETWORK_PORT", "https")
	ferrite.
		NetworkPort("FERRITE_NETWORK_PORT", "example network port variable").
		Deprecated()

	ferrite.Init()

	fmt.Println("<execution continues>")

	// Output:
	// Environment Variables:
	//
	//  ❯ FERRITE_NETWORK_PORT  example network port variable  [ <string> ]  ⚠ deprecated variable set to https
	//
	// <execution continues>
}
