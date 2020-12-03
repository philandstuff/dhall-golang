package internal

import (
	"strings"

	"github.com/philandstuff/dhall-golang/v6/term"
)

// removeLeadingCommonIndent removes the common leading indent from a
// TextLit, as defined in standard/multiline.md
func removeLeadingCommonIndent(text term.TextLit) term.TextLit {
	prefix := longestCommonIndentPrefix(text)
	trimmedText := term.TextLit{Suffix: strings.ReplaceAll(text.Suffix, "\n"+prefix, "\n")}
	for _, chunk := range text.Chunks {
		// trim lines other than the first (everywhere after a '\n')
		trimmedText.Chunks = append(trimmedText.Chunks, term.Chunk{
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

func nonBlankLinesPlusPossibleBlankLastLine(s string) []string {
	lines := strings.FieldsFunc(s, func(r rune) bool { return r == '\n' })
	if len(s) > 0 && s[len(s)-1] == '\n' {
		lines = append(lines, "")
	}
	return lines
}

func allNonBlankLines(text term.TextLit) []string {
	if len(text.Chunks) == 0 {
		return nonBlankLinesPlusPossibleBlankLastLine(text.Suffix)
	}

	lines := make([]string, 0)
	for i, chunk := range text.Chunks {
		thisChunkLines := nonBlankLinesPlusPossibleBlankLastLine(chunk.Prefix)
		if i == 0 {
			lines = append(lines, thisChunkLines...)
		} else {
			// skip the first line, because it's a trailing portion of the previous line
			if len(thisChunkLines) >= 1 {
				lines = append(lines, thisChunkLines[1:]...)
			}
		}
	}
	// skip the first line, because it's a trailing portion of the previous line
	suffixLines := nonBlankLinesPlusPossibleBlankLastLine(text.Suffix)
	if len(suffixLines) >= 1 {
		lines = append(lines, suffixLines[1:]...)
	}
	return lines
}

func longestCommonIndentPrefix(text term.TextLit) string {
	lines := allNonBlankLines(text)
	if len(lines) == 0 {
		return ""
	}
	var prefix strings.Builder
	firstLine := lines[0]
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
		for _, line := range otherLines {
			if i >= len(line) || line[i] != ch {
				break Loop
			}
		}
		prefix.WriteByte(ch)
	}
	return prefix.String()
}
