package parser

import (
	"strings"

	"github.com/philandstuff/dhall-golang/ast"
)

//RemoveLeadingCommonIndent
// removes the common leading indent from a TextLit, as defined in standard/multiline.md
func RemoveLeadingCommonIndent(text ast.TextLit) ast.TextLit {
	prefix := longestCommonIndentPrefix(text)
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
	if len(text.Chunks) == 0 {
		return strings.Split(text.Suffix, "\n")
	}

	chunkLines := make([]string, 0)
	for i, chunk := range text.Chunks {
		if i == 0 {
			chunkLines = append(chunkLines, strings.Split(chunk.Prefix, "\n")...)
		} else {
			// skip the first split, because it's a trailing portion of the previous line
			chunkLines = append(chunkLines, strings.Split(chunk.Prefix, "\n")[1:]...)
		}
	}
	// skip the first split, because it's a trailing portion of the previous line
	return append(chunkLines, strings.Split(text.Suffix, "\n")[1:]...)
}

// lines must be a nonempty slice
func longestCommonIndentPrefix(text ast.TextLit) string {
	lines := allTheLines(text)
	var prefix strings.Builder
	firstLine := strings.Split(text.Suffix, "\n")[0]
	if text.Chunks != nil {
		firstLine = strings.Split(text.Chunks[0].Prefix, "\n")[0]
	}
	otherLines := lines[1:]
	// we range over bytes rather than runes because a) indexing byte
	// arrays is easier than indexing runes in strings, and b) we know
	// we are only looking for '\n', '\t' or ' ', which are all
	// one-byte characters in utf8
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
