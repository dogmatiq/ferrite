package ferrite_test

import (
	"fmt"
	"io"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleFile_required() {
	defer example()()

	v := ferrite.
		File("FERRITE_FILE", "example file").
		Required()

	os.Setenv("FERRITE_FILE", "testdata/hello.txt")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is testdata/hello.txt
}

func ExampleFile_default() {
	defer example()()

	v := ferrite.
		File("FERRITE_FILE", "example file").
		WithDefault("testdata/hello.txt").
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is testdata/hello.txt
}

func ExampleFile_optional() {
	defer example()()

	v := ferrite.
		File("FERRITE_FILE", "example file").
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

func ExampleFile_contentAsReader() {
	defer example()()

	v := ferrite.
		File("FERRITE_FILE", "example file").
		Required()

	os.Setenv("FERRITE_FILE", "testdata/hello.txt")
	ferrite.Init()

	r, err := v.Value().Reader()
	if err != nil {
		panic(err)
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	fmt.Printf("file content is %#v\n", string(data))

	// Output:
	// file content is "Hello, world!\n"
}

func ExampleFile_contentAsBytes() {
	defer example()()

	v := ferrite.
		File("FERRITE_FILE", "example file").
		Required()

	os.Setenv("FERRITE_FILE", "testdata/hello.txt")
	ferrite.Init()

	data, err := v.Value().ReadBytes()
	if err != nil {
		panic(err)
	}

	fmt.Printf("file content is %#v\n", string(data))

	// Output:
	// file content is "Hello, world!\n"
}

func ExampleFile_contentAsString() {
	defer example()()

	v := ferrite.
		File("FERRITE_FILE", "example file").
		Required()

	os.Setenv("FERRITE_FILE", "testdata/hello.txt")
	ferrite.Init()

	data, err := v.Value().ReadString()
	if err != nil {
		panic(err)
	}

	fmt.Printf("file content is %#v\n", data)

	// Output:
	// file content is "Hello, world!\n"
}
