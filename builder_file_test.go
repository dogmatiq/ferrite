package ferrite_test

import (
	"fmt"
	"io"
	"os"

	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type FileBuilder", func() {
	var builder *FileBuilder

	BeforeEach(func() {
		builder = File("FERRITE_FILE", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	It("panics if the name is empty", func() {
		Expect(func() {
			File("", "<desc>").Optional()
		}).To(PanicWith("invalid specification: variable name must not be empty"))
	})

	It("panics if the description is empty", func() {
		Expect(func() {
			File("FERRITE_FILE", "").Optional()
		}).To(PanicWith("specification for FERRITE_FILE is invalid: variable description must not be empty"))
	})

	When("the variable is required", func() {
		When("the value is not empty", func() {
			Describe("func Value()", func() {
				It("returns the value ", func() {
					os.Setenv("FERRITE_FILE", "/path/to/file")

					v := builder.
						Required().
						Value()

					Expect(v).To(Equal(FileName("/path/to/file")))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v := builder.
							WithDefault("/path/to/file").
							Required().
							Value()

						Expect(v).To(Equal(FileName("/path/to/file")))
					})
				})
			})

			When("there is no default value", func() {
				Describe("func Value()", func() {
					It("panics", func() {
						Expect(func() {
							builder.
								Required().
								Value()
						}).To(PanicWith(
							"FERRITE_FILE is undefined and does not have a default value",
						))
					})
				})
			})
		})
	})

	When("the variable is optional", func() {
		When("the value is not empty", func() {
			Describe("func Value()", func() {
				It("returns the value ", func() {
					os.Setenv("FERRITE_FILE", "/path/to/file")

					v, ok := builder.
						Optional().
						Value()

					Expect(ok).To(BeTrue())
					Expect(v).To(Equal(FileName("/path/to/file")))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v, ok := builder.
							WithDefault("/path/to/file").
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(FileName("/path/to/file")))
					})
				})
			})

			When("there is no default value", func() {
				Describe("func Value()", func() {
					It("returns with ok == false", func() {
						_, ok := builder.
							Optional().
							Value()

						Expect(ok).To(BeFalse())
					})
				})
			})
		})
	})

	When("the must-exist constraint is applied", func() {
		It("panics if the file does not exist", func() {
			os.Setenv("FERRITE_FILE", "/path/to/file")

			Expect(func() {
				builder.
					WithMustExist().
					Required().
					Value()
			}).To(PanicWith("value of FERRITE_FILE (/path/to/file) is invalid: the file does not exist"))
		})

		It("panics if the path refers to a directory", func() {
			os.Setenv("FERRITE_FILE", "testdata/dir")

			Expect(func() {
				builder.
					WithMustExist().
					Required().
					Value()
			}).To(PanicWith("value of FERRITE_FILE (testdata/dir) is invalid: the path refers to a directory, expected a file"))
		})
	})
})

var _ = Describe("type FileName", func() {
	var filename = FileName("testdata/hello.txt")

	Describe("func Reader()", func() {
		It("returns a reader for the file", func() {
			r, err := filename.Reader()
			Expect(err).ShouldNot(HaveOccurred())
			defer r.Close()

			data, err := io.ReadAll(r)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).To(Equal("Hello, world!\n"))
		})
	})

	Describe("func ReadBytes()", func() {
		It("returns the file contents as a byte-slice", func() {
			data, err := filename.ReadBytes()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).To(Equal("Hello, world!\n"))
		})
	})

	Describe("func ReadString()", func() {
		It("returns the file contents as a string", func() {
			data, err := filename.ReadString()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data).To(Equal("Hello, world!\n"))
		})
	})
})

func ExampleFile_required() {
	defer example()()

	v := ferrite.
		File("FERRITE_FILE", "example file variable").
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
		File("FERRITE_FILE", "example file variable").
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
		File("FERRITE_FILE", "example file variable").
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
		File("FERRITE_FILE", "example file variable").
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
		File("FERRITE_FILE", "example file variable").
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
		File("FERRITE_FILE", "example file variable").
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

func ExampleFile_deprecated() {
	defer example()()

	v := ferrite.
		File("FERRITE_FILE", "example file variable").
		Deprecated()

	os.Setenv("FERRITE_FILE", "testdata/hello.txt")
	ferrite.Init()

	if x, ok := v.DeprecatedValue(); ok {
		fmt.Println("value is", x)
	} else {
		fmt.Println("value is undefined")
	}

	// Output:
	// Environment Variables:
	//
	//  ❯ FERRITE_FILE  example file variable  [ <string> ]  ⚠ deprecated variable set to testdata/hello.txt
	//
	// value is testdata/hello.txt
}
