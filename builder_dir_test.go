package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type DirBuilder", func() {
	var builder *DirBuilder

	BeforeEach(func() {
		builder = Dir("FERRITE_DIR", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	It("panics if the name is empty", func() {
		Expect(func() {
			Dir("", "<desc>").Optional()
		}).To(PanicWith("invalid specification: variable name must not be empty"))
	})

	It("panics if the description is empty", func() {
		Expect(func() {
			Dir("FERRITE_DIR", "").Optional()
		}).To(PanicWith("specification for FERRITE_DIR is invalid: variable description must not be empty"))
	})

	When("the variable is required", func() {
		When("the value is not empty", func() {
			Describe("func Value()", func() {
				It("returns the value ", func() {
					os.Setenv("FERRITE_DIR", "/path/to/dir")

					v := builder.
						Required().
						Value()

					Expect(v).To(Equal(DirName("/path/to/dir")))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v := builder.
							WithDefault("/path/to/dir").
							Required().
							Value()

						Expect(v).To(Equal(DirName("/path/to/dir")))
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
							"FERRITE_DIR is undefined and does not have a default value",
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
					os.Setenv("FERRITE_DIR", "/path/to/dir")

					v, ok := builder.
						Optional().
						Value()

					Expect(ok).To(BeTrue())
					Expect(v).To(Equal(DirName("/path/to/dir")))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v, ok := builder.
							WithDefault("/path/to/dir").
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(DirName("/path/to/dir")))
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
		It("panics if the directory does not exist", func() {
			os.Setenv("FERRITE_DIR", "/path/to/dir")

			Expect(func() {
				builder.
					WithMustExist().
					Required().
					Value()
			}).To(PanicWith("value of FERRITE_DIR (/path/to/dir) is invalid: expected the directory to exist"))
		})

		It("panics if the path refers to a non-directory", func() {
			os.Setenv("FERRITE_DIR", "testdata/hello.txt")

			Expect(func() {
				builder.
					WithMustExist().
					Required().
					Value()
			}).To(PanicWith("value of FERRITE_DIR (testdata/hello.txt) is invalid: the path refers to a file, expected a directory"))
		})
	})
})

func ExampleDir_required() {
	defer example()()

	v := ferrite.
		File("FERRITE_DIR", "example directory variable").
		Required()

	os.Setenv("FERRITE_DIR", "testdata/dir")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is testdata/dir
}

func ExampleDir_default() {
	defer example()()

	v := ferrite.
		File("FERRITE_DIR", "example directory variable").
		WithDefault("testdata/dir").
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is testdata/dir
}

func ExampleDir_optional() {
	defer example()()

	v := ferrite.
		File("FERRITE_DIR", "example directory variable").
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

func ExampleDir_deprecated() {
	defer example()()

	v := ferrite.
		File("FERRITE_DIR", "example directory variable").
		Deprecated()

	os.Setenv("FERRITE_DIR", "testdata/dir")
	ferrite.Init()

	if x, ok := v.DeprecatedValue(); ok {
		fmt.Println("value is", x)
	} else {
		fmt.Println("value is undefined")
	}

	// Output:
	// Environment Variables:
	//
	//  ❯ FERRITE_DIR  example directory variable  [ <string> ]  ⚠ deprecated variable set to testdata/dir
	//
	// value is testdata/dir
}
