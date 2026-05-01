package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type TextAsBuilder", func() {
	var builder *TextAsBuilder[textValue]

	BeforeEach(func() {
		builder = TextAs[textValue]("FERRITE_TEXT", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	It("panics if the name is empty", func() {
		Expect(func() {
			TextAs[textValue]("", "<desc>").Optional()
		}).To(PanicWith("invalid specification: variable name must not be empty"))
	})

	It("panics if the description is empty", func() {
		Expect(func() {
			TextAs[textValue]("FERRITE_TEXT", "").Optional()
		}).To(PanicWith("specification for FERRITE_TEXT is invalid: variable description must not be empty"))
	})

	When("the variable is required", func() {
		When("the value is not empty", func() {
			Describe("func Value()", func() {
				It("returns the parsed value", func() {
					os.Setenv("FERRITE_TEXT", "text:hello")

					v := builder.
						Required().
						Value()

					Expect(v).To(Equal(textValue{Data: "hello"}))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v := builder.
							WithDefault(textValue{Data: "fallback"}).
							Required().
							Value()

						Expect(v).To(Equal(textValue{Data: "fallback"}))
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
							"FERRITE_TEXT is undefined and does not have a default value",
						))
					})
				})
			})
		})
	})

	When("the variable is optional", func() {
		When("the value is not empty", func() {
			Describe("func Value()", func() {
				It("returns the parsed value", func() {
					os.Setenv("FERRITE_TEXT", "text:hello")

					v, ok := builder.
						Optional().
						Value()

					Expect(ok).To(BeTrue())
					Expect(v).To(Equal(textValue{Data: "hello"}))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v, ok := builder.
							WithDefault(textValue{Data: "fallback"}).
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(textValue{Data: "fallback"}))
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
})

var _ = Describe("func TextAsP", func() {
	var builder *TextAsBuilder[*textValue]

	BeforeEach(func() {
		builder = TextAsP[*textValue]("FERRITE_TEXT_P", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	When("the variable is required", func() {
		When("the value is not empty", func() {
			Describe("func Value()", func() {
				It("returns the parsed value as a pointer", func() {
					os.Setenv("FERRITE_TEXT_P", "text:hello")

					v := builder.
						Required().
						Value()

					Expect(v).To(Equal(&textValue{Data: "hello"}))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v := builder.
							WithDefault(&textValue{Data: "fallback"}).
							Required().
							Value()

						Expect(v).To(Equal(&textValue{Data: "fallback"}))
					})
				})
			})
		})
	})

	When("the variable is optional", func() {
		When("the value is not empty", func() {
			Describe("func Value()", func() {
				It("returns the parsed value as a pointer", func() {
					os.Setenv("FERRITE_TEXT_P", "text:hello")

					v, ok := builder.
						Optional().
						Value()

					Expect(ok).To(BeTrue())
					Expect(v).To(Equal(&textValue{Data: "hello"}))
				})
			})
		})

		When("the value is empty", func() {
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
})

func ExampleTextAs_required() {
	defer example()()

	v := ferrite.
		TextAs[textValue]("FERRITE_TEXT", "example text-marshaled variable").
		Required()

	os.Setenv("FERRITE_TEXT", "text:hello")
	ferrite.Init()

	fmt.Println("value is", v.Value().Data)

	// Output:
	// value is hello
}

func ExampleTextAs_default() {
	defer example()()

	v := ferrite.
		TextAs[textValue]("FERRITE_TEXT", "example text-marshaled variable").
		WithDefault(textValue{Data: "fallback"}).
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value().Data)

	// Output:
	// value is fallback
}

func ExampleTextAs_optional() {
	defer example()()

	v := ferrite.
		TextAs[textValue]("FERRITE_TEXT", "example text-marshaled variable").
		Optional()

	ferrite.Init()

	if x, ok := v.Value(); ok {
		fmt.Println("value is", x.Data)
	} else {
		fmt.Println("value is undefined")
	}

	// Output:
	// value is undefined
}

func ExampleTextAsP_required() {
	defer example()()

	v := ferrite.
		TextAsP[*textValue]("FERRITE_TEXT_P", "example text-marshaled pointer variable").
		Required()

	os.Setenv("FERRITE_TEXT_P", "text:hello")
	ferrite.Init()

	fmt.Println("value is", v.Value().Data)

	// Output:
	// value is hello
}

// textValue is a test type that implements encoding.TextMarshaler via a value
// receiver and encoding.TextUnmarshaler via a pointer receiver.
type textValue struct {
	Data string
}

func (v textValue) MarshalText() ([]byte, error) {
	return []byte("text:" + v.Data), nil
}

func (v *textValue) UnmarshalText(data []byte) error {
	if len(data) < 5 || string(data[:5]) != "text:" {
		return fmt.Errorf("expected text: prefix")
	}
	v.Data = string(data[5:])
	return nil
}
