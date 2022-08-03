package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleKubeService() {
	setUp()
	defer tearDown()

	value := ferrite.
		KubeService("ferrite-svc")

	os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org")
	os.Setenv("FERRITE_SVC_SERVICE_PORT", "12345")
	ferrite.ValidateEnvironment()

	fmt.Println("address is", value.Address())

	// Output:
	// address is host.example.org:12345
}

func ExampleKubeService_namedPort() {
	setUp()
	defer tearDown()

	value := ferrite.
		KubeService("ferrite-svc").
		WithNamedPort("api")

	os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org")
	os.Setenv("FERRITE_SVC_SERVICE_PORT_API", "12345") // note _API suffix
	ferrite.ValidateEnvironment()

	fmt.Println("address is", value.Address())

	// Output:
	// address is host.example.org:12345
}
