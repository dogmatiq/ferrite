package markdown

import (
	"fmt"
	"io"
	"strings"

	"github.com/rivo/uniseg"
)

// table renders a column-aligned markdown table.
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

	for i, columns := range t.rows {
		n, err := t.writeRow(w, columns)
		count += int64(n)
		if err != nil {
			return count, err
		}

		if i == 0 {
			var values []string
			for _, width := range t.widths {
				values = append(values, strings.Repeat("-", width))
			}
			n, err := t.writeRow(w, values)
			count += int64(n)
			if err != nil {
				return count, err
			}
		}

	}

	return count, nil
}

func (t *table) writeRow(w io.Writer, columns []string) (int64, error) {
	var count int64

	for index, text := range columns {
		size, err := fmt.Fprintf(
			w,
			"| %-*s ",
			t.widths[index],
			text,
		)
		count += int64(size)
		if err != nil {
			return count, err
		}
	}

	n, err := fmt.Fprintln(w, "|")
	count += int64(n)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (t *table) String() string {
	var buf strings.Builder
	t.WriteTo(&buf)
	return buf.String()
}
