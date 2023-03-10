package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleKubernetesService_required() {
	defer example()()

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
	defer example()()

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
	defer example()()

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
	defer example()()

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

func ExampleKubernetesService_deprecated() {
	defer example()()

	os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org")
	os.Setenv("FERRITE_SVC_SERVICE_PORT", "12345")
	ferrite.
		KubernetesService("ferrite-svc").
		Deprecated()

	ferrite.Init()

	fmt.Println("<execution continues>")

	// Output:
	// Environment Variables:
	//
	//  ❯ FERRITE_SVC_SERVICE_HOST  kubernetes "ferrite-svc" service host  [ <string> ]  ⚠ deprecated variable set to host.example.org
	//  ❯ FERRITE_SVC_SERVICE_PORT  kubernetes "ferrite-svc" service port  [ <string> ]  ⚠ deprecated variable set to 12345
	//
	// <execution continues>
}
