package parser

import (
	"strings"

	"github.com/philandstuff/dhall-golang/ast"
)

//RemoveLeadingCommonIndent
// removes the common leading indent from a TextLit, as defined in standard/multiline.md
func RemoveLeadingCommonIndent(text ast.TextLit) ast.TextLit {
	prefix := longestCommonIndentPrefix(allTheLines(text))
	trimmedText := ast.TextLit{Suffix: strings.ReplaceAll(text.Suffix, "\n"+prefix, "\n")}
	for _, chunk := range text.Chunks {
		// trim lines other than the first (everywhere after a '\n')
		trimmedText.Chunks = append(trimmedText.Chunks, ast.Chunk{
			Prefix: strings.ReplaceAll(chunk.Prefix, "\n"+prefix, "\n"),
			Expr:   chunk.Expr,
		})
	}
	// trim the first line
	if trimmedText.Chunks != nil {
		trimmedText.Chunks[0].Prefix = strings.Replace(trimmedText.Chunks[0].Prefix, prefix, "", 1)
	} else {
		trimmedText.Suffix = strings.Replace(trimmedText.Suffix, prefix, "", 1)
	}
	return trimmedText
}

func allTheLines(text ast.TextLit) []string {
	chunkLines := make([]string, 0)
	for _, chunk := range text.Chunks {
		chunkLines = append(chunkLines, strings.Split(chunk.Prefix, "\n")...)
	}
	return append(chunkLines, strings.Split(text.Suffix, "\n")...)
}

// lines must be a nonempty slice
func longestCommonIndentPrefix(lines []string) string {
	var prefix strings.Builder
	firstLine := lines[0]
	otherLines := lines[1:]
	// we range over bytes rather than runes because a) indexing byte
	// arrays is easier than indexing runes in strings, and b) we know
	// we are only looking for '\t' or ' ', which are both one-byte
	// characters in utf8, and we will stop at the first sign of a
	// multibyte character
Loop:
	for i, ch := range []byte(firstLine) {
		if ch != '\t' && ch != ' ' {
			break Loop
		}
		for j, line := range otherLines {
			if line == "" && j+1 != len(otherLines) {
				// ignore blank lines unless they are the last line
				continue
			}
			if i >= len(line) || line[i] != ch {
				break Loop
			}
		}
		prefix.WriteByte(ch)
	}
	return prefix.String()
}
