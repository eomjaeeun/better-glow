package ui

import (
	"strings"

	runewidth "github.com/mattn/go-runewidth"
)

type searchMatch struct {
	Line     int
	StartCol int // display width start
	EndCol   int // display width end
}

func (m *pagerModel) findMatches() {
	m.searchMatches = nil
	m.currentMatchIdx = 0
	if m.searchQuery == "" {
		return
	}

	queryLower := strings.ToLower(m.searchQuery)

	for lineIdx, line := range m.contentLines {
		stripped := stripAnsi(line)
		strippedLower := strings.ToLower(stripped)

		searchFrom := 0
		for {
			idx := strings.Index(strippedLower[searchFrom:], queryLower)
			if idx < 0 {
				break
			}
			byteStart := searchFrom + idx
			byteEnd := byteStart + len(queryLower)

			startCol := displayWidthAt(stripped, byteStart)
			endCol := displayWidthAt(stripped, byteEnd)

			m.searchMatches = append(m.searchMatches, searchMatch{
				Line:     lineIdx,
				StartCol: startCol,
				EndCol:   endCol,
			})

			searchFrom = byteEnd
		}
	}
}

func (m *pagerModel) jumpToMatch(idx int) {
	if idx < 0 || idx >= len(m.searchMatches) {
		return
	}
	m.currentMatchIdx = idx
	match := m.searchMatches[idx]

	// Scroll so the match line is visible, centered if possible
	targetOffset := match.Line - m.viewport.Height()/2
	if targetOffset < 0 {
		targetOffset = 0
	}
	m.viewport.SetYOffset(targetOffset)
}

func (m *pagerModel) clearSearch() {
	m.searchQuery = ""
	m.searchInput = ""
	m.searchMatches = nil
	m.currentMatchIdx = 0
	m.searchMode = false
}

func (m pagerModel) viewWithSearch() string {
	viewContent := m.viewport.View()
	if len(m.searchMatches) == 0 {
		return viewContent
	}

	lines := strings.Split(viewContent, "\n")
	offset := m.viewport.YOffset()

	// Build a map of line -> matches for visible lines
	for i := range lines {
		contentLine := offset + i
		var lineMatches []hlRange
		for matchIdx, match := range m.searchMatches {
			if match.Line != contentLine {
				continue
			}
			style := "\033[48;5;58m"  // dark olive background
			off := "\033[49m"
			if matchIdx == m.currentMatchIdx {
				style = "\033[48;5;30m" // teal background (current match)
			}
			lineMatches = append(lineMatches, hlRange{
				Start: match.StartCol,
				End:   match.EndCol,
				On:    style,
				Off:   off,
			})
		}
		if len(lineMatches) > 0 {
			lines[i] = highlightRanges(lines[i], lineMatches)
		}
	}

	return strings.Join(lines, "\n")
}

// displayWidthAt returns the display width of s[:bytePos].
func displayWidthAt(s string, bytePos int) int {
	width := 0
	for i, r := range s {
		if i >= bytePos {
			break
		}
		width += runewidth.RuneWidth(r)
	}
	return width
}
