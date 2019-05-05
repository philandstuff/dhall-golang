package parser

import (
	"regexp"
	"strings"

	"github.com/philandstuff/dhall-golang/ast"
)

//RemoveLeadingCommonIndent
// removes the common leading indent from a TextLit, as defined in standard/multiline.md
func RemoveLeadingCommonIndent(text ast.TextLit) ast.TextLit {
	lines := allTheLines(text)
	prefix := longestCommonIndentPrefix(lines)
	re := regexp.MustCompile("(?m)^" + prefix)
	var newChunks ast.Chunks
	for _, chunk := range text.Chunks {
		newChunks = append(newChunks, ast.Chunk{
			Prefix: re.ReplaceAllLiteralString(chunk.Prefix, ""),
			Expr:   chunk.Expr,
		})
	}
	return ast.TextLit{Chunks: newChunks, Suffix: re.ReplaceAllLiteralString(text.Suffix, "")}
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
