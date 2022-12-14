package validatemode

import (
	"fmt"
	"io"
	"strings"

	"github.com/mattn/go-runewidth"
)

// table renders a column-aligned table.
type table struct {
	widths []int
	rows   [][]string
}

// AddRow adds a row to the table.
func (t *table) AddRow(columns ...string) {
	for len(t.widths) < len(columns) {
		t.widths = append(t.widths, 0)
	}

	for i, col := range columns {
		n := runewidth.StringWidth(col)
		if n > t.widths[i] {
			t.widths[i] = n
		}
	}

	t.rows = append(t.rows, columns)
}

// WriteTo writes the table to w.
func (t *table) WriteTo(w io.Writer) (int64, error) {
	var count int64

	for _, columns := range t.rows {
		n := len(columns) - 1

		for i, col := range columns[:n] {
			n, err := fmt.Fprintf(
				w,
				"%-*s  ",
				t.widths[i],
				col,
			)
			count += int64(n)
			if err != nil {
				return count, err
			}
		}

		n, err := fmt.Fprintln(w, columns[n])
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
