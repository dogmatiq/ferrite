package markdown

func inflect(word string, n int) string {
	if n == 1 {
		return word
	}
	return word + "s"
}
