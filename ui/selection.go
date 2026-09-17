package ui

import (
	"strings"

	runewidth "github.com/mattn/go-runewidth"
)

const (
	// Selection uses a dark blue background, preserving Glamour's foreground colors.
	selBgOn  = "\033[48;5;24m"
	selBgOff = "\033[49m"
)

type selectionPos struct {
	Line int // absolute content line (viewport offset + visible line index)
	Col  int // display column width
}

func normalizeSelection(a, b selectionPos) (start, end selectionPos) {
	if a.Line < b.Line || (a.Line == b.Line && a.Col <= b.Col) {
		return a, b
	}
	return b, a
}

func stripAnsi(s string) string {
	var b strings.Builder
	inEsc := false
	for _, r := range s {
		if r == '\033' {
			inEsc = true
			continue
		}
		if inEsc {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEsc = false
			}
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// highlightLine applies a background color highlight between startCol and endCol
// (display width units). endCol < 0 means highlight to end of line.
// Re-applies the background after any SGR reset from Glamour.
func highlightLine(line string, startCol, endCol int) string {
	var result strings.Builder
	displayPos := 0
	inEsc := false
	highlighted := false

	for _, r := range line {
		if r == '\033' {
			inEsc = true
			result.WriteRune(r)
			continue
		}
		if inEsc {
			result.WriteRune(r)
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEsc = false
				// Re-apply selection background after any SGR sequence (e.g. \033[0m)
				if highlighted && r == 'm' {
					result.WriteString(selBgOn)
				}
			}
			continue
		}

		w := runewidth.RuneWidth(r)

		if !highlighted && displayPos >= startCol {
			result.WriteString(selBgOn)
			highlighted = true
		}
		if highlighted && endCol >= 0 && displayPos >= endCol {
			result.WriteString(selBgOff)
			highlighted = false
		}

		result.WriteRune(r)
		displayPos += w
	}

	if highlighted {
		result.WriteString(selBgOff)
	}

	return result.String()
}

// hlRange defines a highlight range with custom ANSI style.
type hlRange struct {
	Start int
	End   int
	On    string
	Off   string
}

// highlightRanges applies multiple highlight ranges to a line.
// Ranges must be sorted by Start and non-overlapping.
// Re-applies the style after any SGR reset within the highlighted region.
func highlightRanges(line string, ranges []hlRange) string {
	if len(ranges) == 0 {
		return line
	}

	var result strings.Builder
	displayPos := 0
	inEsc := false
	rangeIdx := 0
	active := false

	for _, r := range line {
		if r == '\033' {
			inEsc = true
			result.WriteRune(r)
			continue
		}
		if inEsc {
			result.WriteRune(r)
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEsc = false
				if active && r == 'm' && rangeIdx < len(ranges) {
					result.WriteString(ranges[rangeIdx].On)
				}
			}
			continue
		}

		w := runewidth.RuneWidth(r)

		if active && rangeIdx < len(ranges) && displayPos >= ranges[rangeIdx].End {
			result.WriteString(ranges[rangeIdx].Off)
			active = false
			rangeIdx++
		}

		if !active && rangeIdx < len(ranges) && displayPos >= ranges[rangeIdx].Start {
			result.WriteString(ranges[rangeIdx].On)
			active = true
		}

		result.WriteRune(r)
		displayPos += w
	}

	if active && rangeIdx < len(ranges) {
		result.WriteString(ranges[rangeIdx].Off)
	}

	return result.String()
}

// substringByWidth extracts a substring by display width positions.
// endWidth < 0 means to end of string.
func substringByWidth(s string, startWidth, endWidth int) string {
	var result strings.Builder
	pos := 0
	for _, r := range s {
		w := runewidth.RuneWidth(r)
		if endWidth >= 0 && pos >= endWidth {
			break
		}
		if pos+w > startWidth {
			result.WriteRune(r)
		}
		pos += w
	}
	return result.String()
}
