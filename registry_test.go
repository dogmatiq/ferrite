package ferrite_test

import "github.com/dogmatiq/ferrite"

func ExampleRegistry() {
	defer example()()

	// Create a custom registry.
	reg := ferrite.NewRegistry("custom-registry")

	// Define an environment variable specification within that registry.
	ferrite.
		String("FERRITE_STRING", "example string variable").
		Required(
			ferrite.WithRegistry(reg),
		)

	// Import the custom registry when initializing Ferrite. Multiple custom
	// registries can be imported. The default global registry is always used.
	ferrite.Init(
		ferrite.WithRegistry(reg),
	)

	// Output:
	// Environment Variables:
	//
	//  ❯ FERRITE_STRING  example string variable    <string>    ✗ undefined
	//
	// <process exited with error code 1>
}
