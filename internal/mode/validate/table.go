package validate

import (
	"fmt"
	"io"
	"strings"

	"github.com/rivo/uniseg"
)

// table renders a column-aligned table.
type table struct {

	// rows is the rows of the table, each containing a slice of text with an
	// element for each column.
	rows [][]string

	// columns the number of columns in the table, excluding the right-most
	// columns that are empty in every row.
	columns int

	// widths is the visual width of each column (assuming a monospace font).
	widths []int
}

// AddRow adds a row to the table.
func (t *table) AddRow(columns ...string) {
	for index, text := range columns {
		width := uniseg.GraphemeClusterCount(text)

		if width == 0 {
			continue
		}

		count := index + 1
		for len(t.widths) <= index {
			t.widths = append(t.widths, 0)
		}
		if count > t.columns {
			t.columns = count
		}

		if width > t.widths[index] {
			t.widths[index] = width
		}
	}

	t.rows = append(t.rows, columns)
}

// WriteTo writes the table to w.
func (t *table) WriteTo(w io.Writer) (int64, error) {
	var count int64

	for _, columns := range t.rows {
		for index, text := range columns[:t.columns-1] {
			// width := t.widths[index]

			size, err := fmt.Fprintf(
				w,
				"%-*s  ",
				t.widths[index],
				text,
			)
			count += int64(size)
			if err != nil {
				return count, err
			}
		}

		n, err := fmt.Fprintln(w, columns[t.columns-1])
		count += int64(n)
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (t *table) String() string {
	var buf strings.Builder
	t.WriteTo(&buf)
	return buf.String()
}
