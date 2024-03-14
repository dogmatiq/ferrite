package inflect

// Pluralize returns the plural form of word if n is not 1.
func Pluralize(word string, n int) string {
	if n == 1 {
		return word
	}
	return word + "s"
}
