package ferrite_test

import (
	"io"
	"os"

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
