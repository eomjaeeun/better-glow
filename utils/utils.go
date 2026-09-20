// Package utils provides utility functions.
package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"charm.land/glamour/v2"
	"charm.land/glamour/v2/ansi"
	"charm.land/glamour/v2/styles"
	"github.com/mitchellh/go-homedir"
)

// RemoveFrontmatter removes the front matter header of a markdown file.
func RemoveFrontmatter(content []byte) []byte {
	if frontmatterBoundaries := detectFrontmatter(content); frontmatterBoundaries[0] == 0 {
		return content[frontmatterBoundaries[1]:]
	}
	return content
}

var yamlPattern = regexp.MustCompile(`(?m)^---\r?\n(\s*\r?\n)?`)

func detectFrontmatter(c []byte) []int {
	if matches := yamlPattern.FindAllIndex(c, 2); len(matches) > 1 {
		return []int{matches[0][0], matches[1][1]}
	}
	return []int{-1, -1}
}

// ExpandPath expands tilde and all environment variables from the given path.
func ExpandPath(path string) string {
	s, err := homedir.Expand(path)
	if err == nil {
		return os.ExpandEnv(s)
	}
	return os.ExpandEnv(path)
}

// WrapCodeBlock wraps a string in a code block with the given language.
func WrapCodeBlock(s, language string) string {
	return "```" + language + "\n" + s + "```"
}

var markdownExtensions = []string{
	".md", ".mdown", ".mkdn", ".mkd", ".markdown",
}

// IsMarkdownFile returns whether the filename has a markdown extension.
func IsMarkdownFile(filename string) bool {
	ext := filepath.Ext(filename)

	if ext == "" {
		// By default, assume it's a markdown file.
		return true
	}

	for _, v := range markdownExtensions {
		if strings.EqualFold(ext, v) {
			return true
		}
	}

	// Has an extension but not markdown
	// so assume this is a code file.
	return false
}

// blowStyleConfig returns a Glamour style that uses only terminal-native
// attributes (bold, italic, strikethrough, underline, color) with no
// background colors on inline elements.
func blowStyleConfig() ansi.StyleConfig {
	t := true
	codeColor := stringPtr("203")
	return ansi.StyleConfig{
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "\n",
				BlockSuffix: "\n",
				Color:       stringPtr("252"),
			},
			Margin: uintPtr(2),
		},
		BlockQuote: ansi.StyleBlock{
			Indent:      uintPtr(1),
			IndentToken: stringPtr("│ "),
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "\n",
				Color:       stringPtr("39"),
				Bold:        &t,
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "# ",
				Color:  stringPtr("228"),
				Bold:   &t,
			},
		},
		H2: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{Prefix: "## "},
		},
		H3: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{Prefix: "### "},
		},
		H4: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{Prefix: "#### "},
		},
		H5: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{Prefix: "##### "},
		},
		H6: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "###### ",
				Color:  stringPtr("35"),
			},
		},
		Strikethrough: ansi.StylePrimitive{CrossedOut: &t},
		Emph:         ansi.StylePrimitive{Italic: &t},
		Strong:       ansi.StylePrimitive{Bold: &t},
		HorizontalRule: ansi.StylePrimitive{
			Color:  stringPtr("240"),
			Format: "\n--------\n",
		},
		Item:        ansi.StylePrimitive{BlockPrefix: "• "},
		Enumeration: ansi.StylePrimitive{BlockPrefix: ". "},
		Task: ansi.StyleTask{
			Ticked:   "[✓] ",
			Unticked: "[ ] ",
		},
		Link:     ansi.StylePrimitive{Color: stringPtr("30"), Underline: &t},
		LinkText: ansi.StylePrimitive{Color: stringPtr("35"), Bold: &t},
		Image:     ansi.StylePrimitive{Color: stringPtr("212"), Underline: &t},
		ImageText: ansi.StylePrimitive{Color: stringPtr("243"), Format: "Image: {{.text}} →"},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: " ",
				Suffix: " ",
				Color:  codeColor,
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{Color: stringPtr("244")},
				Margin:         uintPtr(2),
			},
			Chroma: &ansi.Chroma{
				Text:              ansi.StylePrimitive{Color: stringPtr("#C4C4C4")},
				Error:             ansi.StylePrimitive{Color: stringPtr("#F1F1F1")},
				Comment:           ansi.StylePrimitive{Color: stringPtr("#676767")},
				CommentPreproc:    ansi.StylePrimitive{Color: stringPtr("#FF875F")},
				Keyword:           ansi.StylePrimitive{Color: stringPtr("#00AAFF")},
				KeywordReserved:   ansi.StylePrimitive{Color: stringPtr("#FF5FD2")},
				KeywordNamespace:  ansi.StylePrimitive{Color: stringPtr("#FF5F87")},
				KeywordType:       ansi.StylePrimitive{Color: stringPtr("#6E6ED8")},
				Operator:          ansi.StylePrimitive{Color: stringPtr("#EF8080")},
				Punctuation:       ansi.StylePrimitive{Color: stringPtr("#E8E8A8")},
				Name:              ansi.StylePrimitive{Color: stringPtr("#C4C4C4")},
				NameBuiltin:       ansi.StylePrimitive{Color: stringPtr("#FF8EC7")},
				NameTag:           ansi.StylePrimitive{Color: stringPtr("#B083EA")},
				NameAttribute:     ansi.StylePrimitive{Color: stringPtr("#7A7AE6")},
				NameClass:         ansi.StylePrimitive{Color: stringPtr("#F1F1F1"), Underline: &t, Bold: &t},
				NameDecorator:     ansi.StylePrimitive{Color: stringPtr("#FFFF87")},
				NameFunction:      ansi.StylePrimitive{Color: stringPtr("#00D787")},
				LiteralNumber:     ansi.StylePrimitive{Color: stringPtr("#6EEFC0")},
				LiteralString:     ansi.StylePrimitive{Color: stringPtr("#C69669")},
				LiteralStringEscape: ansi.StylePrimitive{Color: stringPtr("#AFFFD7")},
				GenericDeleted:    ansi.StylePrimitive{Color: stringPtr("#FD5B5B")},
				GenericEmph:       ansi.StylePrimitive{Italic: &t},
				GenericInserted:   ansi.StylePrimitive{Color: stringPtr("#00D787")},
				GenericStrong:     ansi.StylePrimitive{Bold: &t},
				GenericSubheading: ansi.StylePrimitive{Color: stringPtr("#777777")},
				Background:        ansi.StylePrimitive{BackgroundColor: stringPtr("#373737")},
			},
		},
		DefinitionDescription: ansi.StylePrimitive{BlockPrefix: "\n🠶 "},
	}
}

func stringPtr(s string) *string { return &s }
func uintPtr(u uint) *uint       { return &u }

// GlamourStyle returns a glamour.TermRendererOption based on the given style.
func GlamourStyle(style string, isCode bool) glamour.TermRendererOption {
	if !isCode {
		if style == "auto" {
			return glamour.WithStyles(blowStyleConfig())
		}
		return glamour.WithStylePath(style)
	}

	var styleConfig ansi.StyleConfig

	switch style {
	case "auto":
		styleConfig = blowStyleConfig()
	case styles.DarkStyle:
		styleConfig = styles.DarkStyleConfig
	case styles.LightStyle:
		styleConfig = styles.LightStyleConfig
	case styles.PinkStyle:
		styleConfig = styles.PinkStyleConfig
	case styles.NoTTYStyle:
		styleConfig = styles.NoTTYStyleConfig
	case styles.DraculaStyle:
		styleConfig = styles.DraculaStyleConfig
	case styles.TokyoNightStyle:
		styleConfig = styles.DraculaStyleConfig
	default:
		return glamour.WithStylesFromJSONFile(style)
	}

	var margin uint
	styleConfig.CodeBlock.Margin = &margin

	return glamour.WithStyles(styleConfig)
}
