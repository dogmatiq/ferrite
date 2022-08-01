package ferrite

import (
	"fmt"
	"io"
	"reflect"
)

// typeName returns T's name.
func typeName[T any]() string {
	return reflect.
		TypeOf((*T)(nil)).
		Elem().
		Name()
}

const (
	pass    = "✓"
	fail    = "✗"
	neutral = "•"
	chevron = "❯"
)

func renderTable(w io.Writer, rows [][]string) error {
	var widths []int

	for _, row := range rows {
		for i, col := range row {
			for len(widths) < i+1 {
				widths = append(widths, 0)
			}

			if len(col) > widths[i] {
				widths[i] = len(col)
			}
		}
	}

	for _, row := range rows {
		n := len(row) - 1

		for i, col := range row[:n] {
			if _, err := fmt.Fprintf(
				w,
				"%-*s  ",
				widths[i],
				col,
			); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprintln(w, row[n]); err != nil {
			return err
		}
	}

	return nil
}
