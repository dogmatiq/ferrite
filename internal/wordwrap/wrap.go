package wordwrap

import (
	"strings"
	"unicode"
	"unicode/utf8"

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
		// If lower < mid < upper, mid is the position within [lower, upper) at
		// which a line break CAN be made if neccessary.
		//
		// [lower, nonspace) is the range within [lower, upper) without trailing
		// whitespace.
		lower, mid, nonspace, upper pos

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

	wrap := func(force bool) {
		if lower.index >= mid.index {
			// wrap not possible
			return
		}

		if !force && nonspace.width-lower.width <= width {
			// wrap not needed
			return
		}

		line := text[lower.index:mid.index]
		lower = mid
		wrapped = append(wrapped, strings.TrimSpace(line))
	}

	for len(remaining) > 0 {
		cluster, remaining, boundaries, state = uniseg.Step(remaining, state)

		clusterWidth := boundaries >> uniseg.ShiftWidth
		upper.index += len(cluster)
		upper.width += clusterWidth
		if !isSpace(cluster) {
			nonspace = upper
		}

		wrap(false)

		switch boundaries & uniseg.MaskLine {
		case uniseg.LineCanBreak:
			mid = upper
		case uniseg.LineMustBreak:
			mid = upper
			wrap(true)
		}
	}

	return wrapped
}

func isSpace(cluster []byte) bool {
	for len(cluster) > 0 {
		r, n := utf8.DecodeRune(cluster)
		if !unicode.IsSpace(r) {
			return false
		}

		cluster = cluster[n:]
	}

	return true
}
