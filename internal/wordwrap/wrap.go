package wordwrap

import (
	"strings"

	"github.com/rivo/uniseg"
)

type pos struct {
	index int
	width int
}

// Wrap wraps text at unicode grapheme boundaries such that it fits within the
// specified maximum column width.
func Wrap(text string, width int) []string {
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

	// wrap writes the a line of text to the wrapped slice, if it is not empty.
	wrap := func() {
		if lower != upper {
			line := text[lower.index:mid.index]
			lower = mid
			wrapped = append(wrapped, strings.TrimSpace(line))
		}
	}

	for len(remaining) > 0 {
		cluster, remaining, boundaries, state = uniseg.Step(remaining, state)

		clusterWidth := boundaries >> uniseg.ShiftWidth

		upper.index += len(cluster)
		upper.width += clusterWidth

		lineWidth := upper.width - lower.width
		wrapNeeded := lineWidth > width
		wrapPossible := lower.index < mid.index

		switch boundaries & uniseg.MaskLine {
		case uniseg.LineCanBreak:
			mid = upper
			wrapPossible = true
		case uniseg.LineMustBreak:
			mid = upper
			wrapNeeded = true
			wrapPossible = true
		}

		if wrapNeeded && wrapPossible {
			wrap()
		}

	}

	// Add any remaining text as the last line.
	wrap()

	return wrapped
}
