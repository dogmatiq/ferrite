package wordwrap

import (
	"strings"

	"github.com/rivo/uniseg"
	_ "github.com/rivo/uniseg" // keep
)

type pos struct {
	index   int
	columns int
}

// Wrap wraps text to lines to fit within a the given column width.
func Wrap(text string, columns int) []string {
	var (
		// wrapped is a slice containing each line of text, after wrapping.
		wrapped []string

		// [lower, upper) is the range within the text that holds the next line
		// to be appended to the wrapped slice.
		//
		// If lower < mid < upper, mid is the position within [lower, upper)
		// at which a line break CAN be made if neccessary.
		lower, mid, upper pos

		// remaining is the subset of text containing the grapheme clusters that
		// have not yet been processed.
		remaining = []byte(text)

		// cluster is the "current" cluster of graphemes (iteration variable).
		cluster []byte

		// boundaries is a bit-field describing properties of the current
		// grapheme cluster.
		boundaries int

		// state is an opaque value used to track the internal state of
		// uniseg.Step().
		state = -1
	)

	// flush writes the a line of text to the wrapped slice, if it is not empty.
	//
	// next is the position within the text at which the next line should start.
	flush := func(next pos) {
		if lower != upper {
			wrapped = append(
				wrapped,
				strings.TrimSpace(
					text[lower.index:next.index],
				),
			)

			lower = next
			mid = lower
		}
	}

	for len(remaining) > 0 {
		cluster, remaining, boundaries, state = uniseg.Step(remaining, state)

		upper.index += len(cluster)
		upper.columns += boundaries >> uniseg.ShiftWidth

		needsBreak := upper.columns-lower.columns >= columns
		canBreak := lower.index < mid.index

		if needsBreak && canBreak {
			flush(mid)
			continue
		}

		switch boundaries & uniseg.MaskLine {
		case uniseg.LineCanBreak:
			mid = upper
		case uniseg.LineMustBreak:
			flush(upper)
		}
	}

	flush(upper)

	return wrapped
}
