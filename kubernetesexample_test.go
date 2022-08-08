package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleKubernetesService_required() {
	setUp()
	defer tearDown()

	v := ferrite.
		KubernetesService("ferrite-svc").
		Required()

	os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org")
	os.Setenv("FERRITE_SVC_SERVICE_PORT", "12345")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is host.example.org:12345
}

func ExampleKubernetesService_default() {
	setUp()
	defer tearDown()

	v := ferrite.
		KubernetesService("ferrite-svc").
		WithDefault("host.example.org", "12345").
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is host.example.org:12345
}

func ExampleKubernetesService_optional() {
	setUp()
	defer tearDown()

	v := ferrite.
		KubernetesService("ferrite-svc").
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

func ExampleKubernetesService_namedPort() {
	setUp()
	defer tearDown()

	v := ferrite.
		KubernetesService("ferrite-svc").
		WithNamedPort("api").
		Required()

	os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org")
	os.Setenv("FERRITE_SVC_SERVICE_PORT_API", "12345") // note _API suffix
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is host.example.org:12345
}
