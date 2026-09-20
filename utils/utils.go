// Package utils provides utility functions.
package utils

import (
	_ "embed"
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
		return true
	}

	for _, v := range markdownExtensions {
		if strings.EqualFold(ext, v) {
			return true
		}
	}

	return false
}

//go:embed blow.json
var blowStyleJSON []byte

// GlamourStyle returns a glamour.TermRendererOption based on the given style.
func GlamourStyle(style string, isCode bool) glamour.TermRendererOption {
	if !isCode {
		if style == "auto" {
			return glamour.WithStylesFromJSONBytes(blowStyleJSON)
		}
		return glamour.WithStylePath(style)
	}

	var styleConfig ansi.StyleConfig

	switch style {
	case "auto":
		return glamour.WithStylesFromJSONBytes(blowStyleJSON)
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
