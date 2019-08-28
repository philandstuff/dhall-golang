
package parser

import (
"bytes"
"crypto/sha256"
"encoding/hex"
"errors"
"fmt"
"io"
"io/ioutil"
"math"
"net"
"net/url"
"os"
"path"
"strconv"
"strings"
"unicode"
"unicode/utf8"
)
import . "github.com/philandstuff/dhall-golang/ast"

// Helper function for parsing all the operator parsing blocks
// see OrExpression for an example of how this is used
func ParseOperator(opcode int, first, rest interface{}) Expr {
    out := first.(Expr)
    if rest == nil { return out }
    for _, b := range rest.([]interface{}) {
        nextExpr := b.([]interface{})[3].(Expr)
        out = Operator{OpCode: opcode, L: out, R: nextExpr}
    }
    return out
}

func IsNonCharacter(r rune) bool {
     return r & 0xfffe == 0xfffe
}

func ValidCodepoint(r rune) bool {
     return utf8.ValidRune(r) && !IsNonCharacter(r)
}

// Helper for parsing unicode code points
func ParseCodepoint(codepointText string) ([]byte, error) {
    i, err := strconv.ParseInt(codepointText, 16, 32)
    if err != nil { return nil, err }
    r := rune(i)
    if !ValidCodepoint(r) {
        return nil, fmt.Errorf("%s is not a valid unicode code point", codepointText)
    }
    return []byte(string([]rune{r})), nil
}


var g = &grammar {
	rules: []*rule{
{
	name: "DhallFile",
	pos: position{line: 57, col: 1, offset: 1189},
	expr: &actionExpr{
	pos: position{line: 57, col: 13, offset: 1203},
	run: (*parser).callonDhallFile1,
	expr: &seqExpr{
	pos: position{line: 57, col: 13, offset: 1203},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 57, col: 13, offset: 1203},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 57, col: 15, offset: 1205},
	name: "CompleteExpression",
},
},
&ruleRefExpr{
	pos: position{line: 57, col: 34, offset: 1224},
	name: "EOF",
},
	},
},
},
},
{
	name: "CompleteExpression",
	pos: position{line: 59, col: 1, offset: 1247},
	expr: &actionExpr{
	pos: position{line: 59, col: 22, offset: 1270},
	run: (*parser).callonCompleteExpression1,
	expr: &seqExpr{
	pos: position{line: 59, col: 22, offset: 1270},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 59, col: 22, offset: 1270},
	name: "_",
},
&labeledExpr{
	pos: position{line: 59, col: 24, offset: 1272},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 59, col: 26, offset: 1274},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 59, col: 37, offset: 1285},
	name: "_",
},
	},
},
},
},
{
	name: "EOL",
	pos: position{line: 61, col: 1, offset: 1306},
	expr: &choiceExpr{
	pos: position{line: 61, col: 7, offset: 1314},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 61, col: 7, offset: 1314},
	val: "\n",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 61, col: 14, offset: 1321},
	run: (*parser).callonEOL3,
	expr: &litMatcher{
	pos: position{line: 61, col: 14, offset: 1321},
	val: "\r\n",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "ValidNonAscii",
	pos: position{line: 63, col: 1, offset: 1358},
	expr: &choiceExpr{
	pos: position{line: 64, col: 5, offset: 1380},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 64, col: 5, offset: 1380},
	val: "[\\u0080-\\uD7FF]",
	ranges: []rune{'\u0080','\ud7ff',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 65, col: 5, offset: 1400},
	val: "[\\uE000-\\uFFFD]",
	ranges: []rune{'\ue000','�',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 66, col: 5, offset: 1420},
	val: "[\\U00010000-\\U0001FFFD]",
	ranges: []rune{'𐀀','\U0001fffd',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 67, col: 5, offset: 1448},
	val: "[\\U00020000-\\U0002FFFD]",
	ranges: []rune{'𠀀','\U0002fffd',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 68, col: 5, offset: 1476},
	val: "[\\U00030000-\\U0003FFFD]",
	ranges: []rune{'\U00030000','\U0003fffd',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 69, col: 5, offset: 1504},
	val: "[\\U00040000-\\U0004FFFD]",
	ranges: []rune{'\U00040000','\U0004fffd',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 70, col: 5, offset: 1532},
	val: "[\\U00050000-\\U0005FFFD]",
	ranges: []rune{'\U00050000','\U0005fffd',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 71, col: 5, offset: 1560},
	val: "[\\U00060000-\\U0006FFFD]",
	ranges: []rune{'\U00060000','\U0006fffd',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 72, col: 5, offset: 1588},
	val: "[\\U00070000-\\U0007FFFD]",
	ranges: []rune{'\U00070000','\U0007fffd',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 73, col: 5, offset: 1616},
	val: "[\\U00080000-\\U0008FFFD]",
	ranges: []rune{'\U00080000','\U0008fffd',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 74, col: 5, offset: 1644},
	val: "[\\U00090000-\\U0009FFFD]",
	ranges: []rune{'\U00090000','\U0009fffd',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 75, col: 5, offset: 1672},
	val: "[\\U000A0000-\\U000AFFFD]",
	ranges: []rune{'\U000a0000','\U000afffd',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 76, col: 5, offset: 1700},
	val: "[\\U000B0000-\\U000BFFFD]",
	ranges: []rune{'\U000b0000','\U000bfffd',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 77, col: 5, offset: 1728},
	val: "[\\U000C0000-\\U000CFFFD]",
	ranges: []rune{'\U000c0000','\U000cfffd',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 78, col: 5, offset: 1756},
	val: "[\\U000D0000-\\U000DFFFD]",
	ranges: []rune{'\U000d0000','\U000dfffd',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 79, col: 5, offset: 1784},
	val: "[\\U000E0000-\\U000EFFFD]",
	ranges: []rune{'\U000e0000','\U000efffd',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 80, col: 5, offset: 1812},
	val: "[\\U000F0000-\\U000FFFFD]",
	ranges: []rune{'\U000f0000','\U000ffffd',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 81, col: 5, offset: 1840},
	val: "[\\U000100000-\\U00010FFFD]",
	chars: []rune{'𐀀','D',},
	ranges: []rune{'0','\U00010fff',},
	ignoreCase: false,
	inverted: false,
},
	},
},
},
{
	name: "BlockComment",
	pos: position{line: 83, col: 1, offset: 1867},
	expr: &seqExpr{
	pos: position{line: 83, col: 16, offset: 1884},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 83, col: 16, offset: 1884},
	val: "{-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 83, col: 21, offset: 1889},
	name: "BlockCommentContinue",
},
	},
},
},
{
	name: "BlockCommentChar",
	pos: position{line: 85, col: 1, offset: 1911},
	expr: &choiceExpr{
	pos: position{line: 86, col: 5, offset: 1936},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 86, col: 5, offset: 1936},
	val: "[\\x20-\\x7f]",
	ranges: []rune{' ','\u007f',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 87, col: 5, offset: 1952},
	name: "ValidNonAscii",
},
&litMatcher{
	pos: position{line: 88, col: 5, offset: 1970},
	val: "\t",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 89, col: 5, offset: 1979},
	name: "EOL",
},
	},
},
},
{
	name: "BlockCommentContinue",
	pos: position{line: 91, col: 1, offset: 1984},
	expr: &choiceExpr{
	pos: position{line: 92, col: 7, offset: 2015},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 92, col: 7, offset: 2015},
	val: "-}",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 93, col: 7, offset: 2026},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 93, col: 7, offset: 2026},
	name: "BlockComment",
},
&ruleRefExpr{
	pos: position{line: 93, col: 20, offset: 2039},
	name: "BlockCommentContinue",
},
	},
},
&seqExpr{
	pos: position{line: 94, col: 7, offset: 2066},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 94, col: 7, offset: 2066},
	name: "BlockCommentChar",
},
&ruleRefExpr{
	pos: position{line: 94, col: 24, offset: 2083},
	name: "BlockCommentContinue",
},
	},
},
	},
},
},
{
	name: "NotEOL",
	pos: position{line: 96, col: 1, offset: 2105},
	expr: &choiceExpr{
	pos: position{line: 96, col: 10, offset: 2116},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 96, col: 10, offset: 2116},
	val: "[\\x20-\\x7f]",
	ranges: []rune{' ','\u007f',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 96, col: 24, offset: 2130},
	name: "ValidNonAscii",
},
&litMatcher{
	pos: position{line: 96, col: 40, offset: 2146},
	val: "\t",
	ignoreCase: false,
},
	},
},
},
{
	name: "LineComment",
	pos: position{line: 98, col: 1, offset: 2152},
	expr: &actionExpr{
	pos: position{line: 98, col: 15, offset: 2168},
	run: (*parser).callonLineComment1,
	expr: &seqExpr{
	pos: position{line: 98, col: 15, offset: 2168},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 98, col: 15, offset: 2168},
	val: "--",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 98, col: 20, offset: 2173},
	label: "content",
	expr: &actionExpr{
	pos: position{line: 98, col: 29, offset: 2182},
	run: (*parser).callonLineComment5,
	expr: &zeroOrMoreExpr{
	pos: position{line: 98, col: 29, offset: 2182},
	expr: &ruleRefExpr{
	pos: position{line: 98, col: 29, offset: 2182},
	name: "NotEOL",
},
},
},
},
&ruleRefExpr{
	pos: position{line: 98, col: 68, offset: 2221},
	name: "EOL",
},
	},
},
},
},
{
	name: "WhitespaceChunk",
	pos: position{line: 100, col: 1, offset: 2250},
	expr: &choiceExpr{
	pos: position{line: 100, col: 19, offset: 2270},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 100, col: 19, offset: 2270},
	val: " ",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 100, col: 25, offset: 2276},
	val: "\t",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 100, col: 32, offset: 2283},
	name: "EOL",
},
&ruleRefExpr{
	pos: position{line: 100, col: 38, offset: 2289},
	name: "LineComment",
},
&ruleRefExpr{
	pos: position{line: 100, col: 52, offset: 2303},
	name: "BlockComment",
},
	},
},
},
{
	name: "_",
	pos: position{line: 102, col: 1, offset: 2317},
	expr: &zeroOrMoreExpr{
	pos: position{line: 102, col: 5, offset: 2323},
	expr: &ruleRefExpr{
	pos: position{line: 102, col: 5, offset: 2323},
	name: "WhitespaceChunk",
},
},
},
{
	name: "_1",
	pos: position{line: 104, col: 1, offset: 2341},
	expr: &oneOrMoreExpr{
	pos: position{line: 104, col: 6, offset: 2348},
	expr: &ruleRefExpr{
	pos: position{line: 104, col: 6, offset: 2348},
	name: "WhitespaceChunk",
},
},
},
{
	name: "Digit",
	pos: position{line: 106, col: 1, offset: 2366},
	expr: &charClassMatcher{
	pos: position{line: 106, col: 9, offset: 2376},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "HexDig",
	pos: position{line: 108, col: 1, offset: 2383},
	expr: &choiceExpr{
	pos: position{line: 108, col: 10, offset: 2394},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 108, col: 10, offset: 2394},
	name: "Digit",
},
&charClassMatcher{
	pos: position{line: 108, col: 18, offset: 2402},
	val: "[a-f]i",
	ranges: []rune{'a','f',},
	ignoreCase: true,
	inverted: false,
},
	},
},
},
{
	name: "SimpleLabelFirstChar",
	pos: position{line: 110, col: 1, offset: 2410},
	expr: &charClassMatcher{
	pos: position{line: 110, col: 24, offset: 2435},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SimpleLabelNextChar",
	pos: position{line: 111, col: 1, offset: 2445},
	expr: &charClassMatcher{
	pos: position{line: 111, col: 23, offset: 2469},
	val: "[A-Za-z0-9_/-]",
	chars: []rune{'_','/','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SimpleLabel",
	pos: position{line: 112, col: 1, offset: 2484},
	expr: &choiceExpr{
	pos: position{line: 112, col: 15, offset: 2500},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 112, col: 15, offset: 2500},
	run: (*parser).callonSimpleLabel2,
	expr: &seqExpr{
	pos: position{line: 112, col: 15, offset: 2500},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 112, col: 15, offset: 2500},
	name: "Keyword",
},
&oneOrMoreExpr{
	pos: position{line: 112, col: 23, offset: 2508},
	expr: &ruleRefExpr{
	pos: position{line: 112, col: 23, offset: 2508},
	name: "SimpleLabelNextChar",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 113, col: 13, offset: 2572},
	run: (*parser).callonSimpleLabel7,
	expr: &seqExpr{
	pos: position{line: 113, col: 13, offset: 2572},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 113, col: 13, offset: 2572},
	expr: &ruleRefExpr{
	pos: position{line: 113, col: 14, offset: 2573},
	name: "Keyword",
},
},
&ruleRefExpr{
	pos: position{line: 113, col: 22, offset: 2581},
	name: "SimpleLabelFirstChar",
},
&zeroOrMoreExpr{
	pos: position{line: 113, col: 43, offset: 2602},
	expr: &ruleRefExpr{
	pos: position{line: 113, col: 43, offset: 2602},
	name: "SimpleLabelNextChar",
},
},
	},
},
},
	},
},
},
{
	name: "QuotedLabelChar",
	pos: position{line: 118, col: 1, offset: 2687},
	expr: &charClassMatcher{
	pos: position{line: 118, col: 19, offset: 2707},
	val: "[\\x20-\\x5f\\x61-\\x7e]",
	ranges: []rune{' ','_','a','~',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "QuotedLabel",
	pos: position{line: 119, col: 1, offset: 2728},
	expr: &actionExpr{
	pos: position{line: 119, col: 15, offset: 2744},
	run: (*parser).callonQuotedLabel1,
	expr: &oneOrMoreExpr{
	pos: position{line: 119, col: 15, offset: 2744},
	expr: &ruleRefExpr{
	pos: position{line: 119, col: 15, offset: 2744},
	name: "QuotedLabelChar",
},
},
},
},
{
	name: "Label",
	pos: position{line: 121, col: 1, offset: 2793},
	expr: &choiceExpr{
	pos: position{line: 121, col: 9, offset: 2803},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 121, col: 9, offset: 2803},
	run: (*parser).callonLabel2,
	expr: &seqExpr{
	pos: position{line: 121, col: 9, offset: 2803},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 121, col: 9, offset: 2803},
	val: "`",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 121, col: 13, offset: 2807},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 121, col: 19, offset: 2813},
	name: "QuotedLabel",
},
},
&litMatcher{
	pos: position{line: 121, col: 31, offset: 2825},
	val: "`",
	ignoreCase: false,
},
	},
},
},
&actionExpr{
	pos: position{line: 122, col: 9, offset: 2859},
	run: (*parser).callonLabel8,
	expr: &labeledExpr{
	pos: position{line: 122, col: 9, offset: 2859},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 122, col: 15, offset: 2865},
	name: "SimpleLabel",
},
},
},
	},
},
},
{
	name: "NonreservedLabel",
	pos: position{line: 124, col: 1, offset: 2900},
	expr: &choiceExpr{
	pos: position{line: 124, col: 20, offset: 2921},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 124, col: 20, offset: 2921},
	run: (*parser).callonNonreservedLabel2,
	expr: &seqExpr{
	pos: position{line: 124, col: 20, offset: 2921},
	exprs: []interface{}{
&andExpr{
	pos: position{line: 124, col: 20, offset: 2921},
	expr: &seqExpr{
	pos: position{line: 124, col: 22, offset: 2923},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 124, col: 22, offset: 2923},
	name: "Reserved",
},
&ruleRefExpr{
	pos: position{line: 124, col: 31, offset: 2932},
	name: "SimpleLabelNextChar",
},
	},
},
},
&labeledExpr{
	pos: position{line: 124, col: 52, offset: 2953},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 124, col: 58, offset: 2959},
	name: "Label",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 125, col: 19, offset: 3005},
	run: (*parser).callonNonreservedLabel10,
	expr: &seqExpr{
	pos: position{line: 125, col: 19, offset: 3005},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 125, col: 19, offset: 3005},
	expr: &ruleRefExpr{
	pos: position{line: 125, col: 20, offset: 3006},
	name: "Reserved",
},
},
&labeledExpr{
	pos: position{line: 125, col: 29, offset: 3015},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 125, col: 35, offset: 3021},
	name: "Label",
},
},
	},
},
},
	},
},
},
{
	name: "AnyLabel",
	pos: position{line: 127, col: 1, offset: 3050},
	expr: &ruleRefExpr{
	pos: position{line: 127, col: 12, offset: 3063},
	name: "Label",
},
},
{
	name: "DoubleQuoteChunk",
	pos: position{line: 130, col: 1, offset: 3071},
	expr: &choiceExpr{
	pos: position{line: 131, col: 6, offset: 3097},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 131, col: 6, offset: 3097},
	name: "Interpolation",
},
&actionExpr{
	pos: position{line: 132, col: 6, offset: 3116},
	run: (*parser).callonDoubleQuoteChunk3,
	expr: &seqExpr{
	pos: position{line: 132, col: 6, offset: 3116},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 132, col: 6, offset: 3116},
	val: "\\",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 132, col: 11, offset: 3121},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 132, col: 13, offset: 3123},
	name: "DoubleQuoteEscaped",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 133, col: 6, offset: 3165},
	name: "DoubleQuoteChar",
},
	},
},
},
{
	name: "DoubleQuoteEscaped",
	pos: position{line: 135, col: 1, offset: 3182},
	expr: &choiceExpr{
	pos: position{line: 136, col: 8, offset: 3212},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 136, col: 8, offset: 3212},
	val: "\"",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 137, col: 8, offset: 3223},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 138, col: 8, offset: 3234},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 139, col: 8, offset: 3246},
	val: "/",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 140, col: 8, offset: 3257},
	run: (*parser).callonDoubleQuoteEscaped6,
	expr: &litMatcher{
	pos: position{line: 140, col: 8, offset: 3257},
	val: "b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 141, col: 8, offset: 3297},
	run: (*parser).callonDoubleQuoteEscaped8,
	expr: &litMatcher{
	pos: position{line: 141, col: 8, offset: 3297},
	val: "f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 142, col: 8, offset: 3337},
	run: (*parser).callonDoubleQuoteEscaped10,
	expr: &litMatcher{
	pos: position{line: 142, col: 8, offset: 3337},
	val: "n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 143, col: 8, offset: 3377},
	run: (*parser).callonDoubleQuoteEscaped12,
	expr: &litMatcher{
	pos: position{line: 143, col: 8, offset: 3377},
	val: "r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 144, col: 8, offset: 3417},
	run: (*parser).callonDoubleQuoteEscaped14,
	expr: &litMatcher{
	pos: position{line: 144, col: 8, offset: 3417},
	val: "t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 145, col: 8, offset: 3457},
	run: (*parser).callonDoubleQuoteEscaped16,
	expr: &seqExpr{
	pos: position{line: 145, col: 8, offset: 3457},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 145, col: 8, offset: 3457},
	val: "u",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 145, col: 12, offset: 3461},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 145, col: 14, offset: 3463},
	name: "UnicodeEscape",
},
},
	},
},
},
	},
},
},
{
	name: "UnicodeEscape",
	pos: position{line: 147, col: 1, offset: 3496},
	expr: &choiceExpr{
	pos: position{line: 148, col: 9, offset: 3522},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 148, col: 9, offset: 3522},
	run: (*parser).callonUnicodeEscape2,
	expr: &seqExpr{
	pos: position{line: 148, col: 9, offset: 3522},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 148, col: 9, offset: 3522},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 148, col: 16, offset: 3529},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 148, col: 23, offset: 3536},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 148, col: 30, offset: 3543},
	name: "HexDig",
},
	},
},
},
&actionExpr{
	pos: position{line: 151, col: 9, offset: 3620},
	run: (*parser).callonUnicodeEscape8,
	expr: &seqExpr{
	pos: position{line: 151, col: 9, offset: 3620},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 151, col: 9, offset: 3620},
	val: "{",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 151, col: 13, offset: 3624},
	expr: &ruleRefExpr{
	pos: position{line: 151, col: 13, offset: 3624},
	name: "HexDig",
},
},
&litMatcher{
	pos: position{line: 151, col: 21, offset: 3632},
	val: "}",
	ignoreCase: false,
},
	},
},
},
	},
},
},
{
	name: "DoubleQuoteChar",
	pos: position{line: 155, col: 1, offset: 3716},
	expr: &choiceExpr{
	pos: position{line: 156, col: 6, offset: 3741},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 156, col: 6, offset: 3741},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 157, col: 6, offset: 3758},
	val: "[\\x23-\\x5b]",
	ranges: []rune{'#','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 158, col: 6, offset: 3775},
	val: "[\\x5d-\\x7f]",
	ranges: []rune{']','\u007f',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 159, col: 6, offset: 3792},
	name: "ValidNonAscii",
},
	},
},
},
{
	name: "DoubleQuoteLiteral",
	pos: position{line: 161, col: 1, offset: 3807},
	expr: &actionExpr{
	pos: position{line: 161, col: 22, offset: 3830},
	run: (*parser).callonDoubleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 161, col: 22, offset: 3830},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 161, col: 22, offset: 3830},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 161, col: 26, offset: 3834},
	label: "chunks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 161, col: 33, offset: 3841},
	expr: &ruleRefExpr{
	pos: position{line: 161, col: 33, offset: 3841},
	name: "DoubleQuoteChunk",
},
},
},
&litMatcher{
	pos: position{line: 161, col: 51, offset: 3859},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "SingleQuoteContinue",
	pos: position{line: 178, col: 1, offset: 4327},
	expr: &choiceExpr{
	pos: position{line: 179, col: 7, offset: 4357},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 179, col: 7, offset: 4357},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 179, col: 7, offset: 4357},
	name: "Interpolation",
},
&ruleRefExpr{
	pos: position{line: 179, col: 21, offset: 4371},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 180, col: 7, offset: 4397},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 180, col: 7, offset: 4397},
	name: "EscapedQuotePair",
},
&ruleRefExpr{
	pos: position{line: 180, col: 24, offset: 4414},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 181, col: 7, offset: 4440},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 181, col: 7, offset: 4440},
	name: "EscapedInterpolation",
},
&ruleRefExpr{
	pos: position{line: 181, col: 28, offset: 4461},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 182, col: 7, offset: 4487},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 182, col: 7, offset: 4487},
	name: "SingleQuoteChar",
},
&ruleRefExpr{
	pos: position{line: 182, col: 23, offset: 4503},
	name: "SingleQuoteContinue",
},
	},
},
&litMatcher{
	pos: position{line: 183, col: 7, offset: 4529},
	val: "''",
	ignoreCase: false,
},
	},
},
},
{
	name: "EscapedQuotePair",
	pos: position{line: 185, col: 1, offset: 4535},
	expr: &actionExpr{
	pos: position{line: 185, col: 20, offset: 4556},
	run: (*parser).callonEscapedQuotePair1,
	expr: &litMatcher{
	pos: position{line: 185, col: 20, offset: 4556},
	val: "'''",
	ignoreCase: false,
},
},
},
{
	name: "EscapedInterpolation",
	pos: position{line: 189, col: 1, offset: 4691},
	expr: &actionExpr{
	pos: position{line: 189, col: 24, offset: 4716},
	run: (*parser).callonEscapedInterpolation1,
	expr: &litMatcher{
	pos: position{line: 189, col: 24, offset: 4716},
	val: "''${",
	ignoreCase: false,
},
},
},
{
	name: "SingleQuoteChar",
	pos: position{line: 191, col: 1, offset: 4758},
	expr: &choiceExpr{
	pos: position{line: 192, col: 6, offset: 4783},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 192, col: 6, offset: 4783},
	val: "[\\x20-\\x7f]",
	ranges: []rune{' ','\u007f',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 193, col: 6, offset: 4800},
	name: "ValidNonAscii",
},
&litMatcher{
	pos: position{line: 194, col: 6, offset: 4819},
	val: "\t",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 195, col: 6, offset: 4829},
	name: "EOL",
},
	},
},
},
{
	name: "SingleQuoteLiteral",
	pos: position{line: 197, col: 1, offset: 4834},
	expr: &actionExpr{
	pos: position{line: 197, col: 22, offset: 4857},
	run: (*parser).callonSingleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 197, col: 22, offset: 4857},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 197, col: 22, offset: 4857},
	val: "''",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 197, col: 27, offset: 4862},
	name: "EOL",
},
&labeledExpr{
	pos: position{line: 197, col: 31, offset: 4866},
	label: "content",
	expr: &ruleRefExpr{
	pos: position{line: 197, col: 39, offset: 4874},
	name: "SingleQuoteContinue",
},
},
	},
},
},
},
{
	name: "Interpolation",
	pos: position{line: 215, col: 1, offset: 5424},
	expr: &actionExpr{
	pos: position{line: 215, col: 17, offset: 5442},
	run: (*parser).callonInterpolation1,
	expr: &seqExpr{
	pos: position{line: 215, col: 17, offset: 5442},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 215, col: 17, offset: 5442},
	val: "${",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 215, col: 22, offset: 5447},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 215, col: 24, offset: 5449},
	name: "CompleteExpression",
},
},
&litMatcher{
	pos: position{line: 215, col: 43, offset: 5468},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TextLiteral",
	pos: position{line: 217, col: 1, offset: 5491},
	expr: &choiceExpr{
	pos: position{line: 217, col: 15, offset: 5507},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 217, col: 15, offset: 5507},
	name: "DoubleQuoteLiteral",
},
&ruleRefExpr{
	pos: position{line: 217, col: 36, offset: 5528},
	name: "SingleQuoteLiteral",
},
	},
},
},
{
	name: "Reserved",
	pos: position{line: 220, col: 1, offset: 5633},
	expr: &choiceExpr{
	pos: position{line: 221, col: 5, offset: 5650},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 221, col: 5, offset: 5650},
	run: (*parser).callonReserved2,
	expr: &litMatcher{
	pos: position{line: 221, col: 5, offset: 5650},
	val: "Natural/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 222, col: 5, offset: 5699},
	run: (*parser).callonReserved4,
	expr: &litMatcher{
	pos: position{line: 222, col: 5, offset: 5699},
	val: "Natural/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 223, col: 5, offset: 5746},
	run: (*parser).callonReserved6,
	expr: &litMatcher{
	pos: position{line: 223, col: 5, offset: 5746},
	val: "Natural/isZero",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 224, col: 5, offset: 5797},
	run: (*parser).callonReserved8,
	expr: &litMatcher{
	pos: position{line: 224, col: 5, offset: 5797},
	val: "Natural/even",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 225, col: 5, offset: 5844},
	run: (*parser).callonReserved10,
	expr: &litMatcher{
	pos: position{line: 225, col: 5, offset: 5844},
	val: "Natural/odd",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 226, col: 5, offset: 5889},
	run: (*parser).callonReserved12,
	expr: &litMatcher{
	pos: position{line: 226, col: 5, offset: 5889},
	val: "Natural/toInteger",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 227, col: 5, offset: 5946},
	run: (*parser).callonReserved14,
	expr: &litMatcher{
	pos: position{line: 227, col: 5, offset: 5946},
	val: "Natural/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 228, col: 5, offset: 5993},
	run: (*parser).callonReserved16,
	expr: &litMatcher{
	pos: position{line: 228, col: 5, offset: 5993},
	val: "Integer/toDouble",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 229, col: 5, offset: 6048},
	run: (*parser).callonReserved18,
	expr: &litMatcher{
	pos: position{line: 229, col: 5, offset: 6048},
	val: "Integer/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 230, col: 5, offset: 6095},
	run: (*parser).callonReserved20,
	expr: &litMatcher{
	pos: position{line: 230, col: 5, offset: 6095},
	val: "Double/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 231, col: 5, offset: 6140},
	run: (*parser).callonReserved22,
	expr: &litMatcher{
	pos: position{line: 231, col: 5, offset: 6140},
	val: "List/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 232, col: 5, offset: 6183},
	run: (*parser).callonReserved24,
	expr: &litMatcher{
	pos: position{line: 232, col: 5, offset: 6183},
	val: "List/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 233, col: 5, offset: 6224},
	run: (*parser).callonReserved26,
	expr: &litMatcher{
	pos: position{line: 233, col: 5, offset: 6224},
	val: "List/length",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 234, col: 5, offset: 6269},
	run: (*parser).callonReserved28,
	expr: &litMatcher{
	pos: position{line: 234, col: 5, offset: 6269},
	val: "List/head",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 235, col: 5, offset: 6310},
	run: (*parser).callonReserved30,
	expr: &litMatcher{
	pos: position{line: 235, col: 5, offset: 6310},
	val: "List/last",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 236, col: 5, offset: 6351},
	run: (*parser).callonReserved32,
	expr: &litMatcher{
	pos: position{line: 236, col: 5, offset: 6351},
	val: "List/indexed",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 237, col: 5, offset: 6398},
	run: (*parser).callonReserved34,
	expr: &litMatcher{
	pos: position{line: 237, col: 5, offset: 6398},
	val: "List/reverse",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 238, col: 5, offset: 6445},
	run: (*parser).callonReserved36,
	expr: &litMatcher{
	pos: position{line: 238, col: 5, offset: 6445},
	val: "Optional/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 239, col: 5, offset: 6496},
	run: (*parser).callonReserved38,
	expr: &litMatcher{
	pos: position{line: 239, col: 5, offset: 6496},
	val: "Optional/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 240, col: 5, offset: 6545},
	run: (*parser).callonReserved40,
	expr: &litMatcher{
	pos: position{line: 240, col: 5, offset: 6545},
	val: "Text/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 241, col: 5, offset: 6586},
	run: (*parser).callonReserved42,
	expr: &litMatcher{
	pos: position{line: 241, col: 5, offset: 6586},
	val: "Bool",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 242, col: 5, offset: 6618},
	run: (*parser).callonReserved44,
	expr: &litMatcher{
	pos: position{line: 242, col: 5, offset: 6618},
	val: "True",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 243, col: 5, offset: 6650},
	run: (*parser).callonReserved46,
	expr: &litMatcher{
	pos: position{line: 243, col: 5, offset: 6650},
	val: "False",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 244, col: 5, offset: 6684},
	run: (*parser).callonReserved48,
	expr: &litMatcher{
	pos: position{line: 244, col: 5, offset: 6684},
	val: "Optional",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 245, col: 5, offset: 6724},
	run: (*parser).callonReserved50,
	expr: &litMatcher{
	pos: position{line: 245, col: 5, offset: 6724},
	val: "Natural",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 246, col: 5, offset: 6762},
	run: (*parser).callonReserved52,
	expr: &litMatcher{
	pos: position{line: 246, col: 5, offset: 6762},
	val: "Integer",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 247, col: 5, offset: 6800},
	run: (*parser).callonReserved54,
	expr: &litMatcher{
	pos: position{line: 247, col: 5, offset: 6800},
	val: "Double",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 248, col: 5, offset: 6836},
	run: (*parser).callonReserved56,
	expr: &litMatcher{
	pos: position{line: 248, col: 5, offset: 6836},
	val: "Text",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 249, col: 5, offset: 6868},
	run: (*parser).callonReserved58,
	expr: &litMatcher{
	pos: position{line: 249, col: 5, offset: 6868},
	val: "List",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 250, col: 5, offset: 6900},
	run: (*parser).callonReserved60,
	expr: &litMatcher{
	pos: position{line: 250, col: 5, offset: 6900},
	val: "None",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 251, col: 5, offset: 6932},
	run: (*parser).callonReserved62,
	expr: &litMatcher{
	pos: position{line: 251, col: 5, offset: 6932},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 252, col: 5, offset: 6964},
	run: (*parser).callonReserved64,
	expr: &litMatcher{
	pos: position{line: 252, col: 5, offset: 6964},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 253, col: 5, offset: 6996},
	run: (*parser).callonReserved66,
	expr: &litMatcher{
	pos: position{line: 253, col: 5, offset: 6996},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "If",
	pos: position{line: 255, col: 1, offset: 7025},
	expr: &litMatcher{
	pos: position{line: 255, col: 6, offset: 7032},
	val: "if",
	ignoreCase: false,
},
},
{
	name: "Then",
	pos: position{line: 256, col: 1, offset: 7037},
	expr: &litMatcher{
	pos: position{line: 256, col: 8, offset: 7046},
	val: "then",
	ignoreCase: false,
},
},
{
	name: "Else",
	pos: position{line: 257, col: 1, offset: 7053},
	expr: &litMatcher{
	pos: position{line: 257, col: 8, offset: 7062},
	val: "else",
	ignoreCase: false,
},
},
{
	name: "Let",
	pos: position{line: 258, col: 1, offset: 7069},
	expr: &litMatcher{
	pos: position{line: 258, col: 7, offset: 7077},
	val: "let",
	ignoreCase: false,
},
},
{
	name: "In",
	pos: position{line: 259, col: 1, offset: 7083},
	expr: &litMatcher{
	pos: position{line: 259, col: 6, offset: 7090},
	val: "in",
	ignoreCase: false,
},
},
{
	name: "As",
	pos: position{line: 260, col: 1, offset: 7095},
	expr: &litMatcher{
	pos: position{line: 260, col: 6, offset: 7102},
	val: "as",
	ignoreCase: false,
},
},
{
	name: "Using",
	pos: position{line: 261, col: 1, offset: 7107},
	expr: &litMatcher{
	pos: position{line: 261, col: 9, offset: 7117},
	val: "using",
	ignoreCase: false,
},
},
{
	name: "Merge",
	pos: position{line: 262, col: 1, offset: 7125},
	expr: &litMatcher{
	pos: position{line: 262, col: 9, offset: 7135},
	val: "merge",
	ignoreCase: false,
},
},
{
	name: "Missing",
	pos: position{line: 263, col: 1, offset: 7143},
	expr: &actionExpr{
	pos: position{line: 263, col: 11, offset: 7155},
	run: (*parser).callonMissing1,
	expr: &litMatcher{
	pos: position{line: 263, col: 11, offset: 7155},
	val: "missing",
	ignoreCase: false,
},
},
},
{
	name: "True",
	pos: position{line: 264, col: 1, offset: 7191},
	expr: &litMatcher{
	pos: position{line: 264, col: 8, offset: 7200},
	val: "True",
	ignoreCase: false,
},
},
{
	name: "False",
	pos: position{line: 265, col: 1, offset: 7207},
	expr: &litMatcher{
	pos: position{line: 265, col: 9, offset: 7217},
	val: "False",
	ignoreCase: false,
},
},
{
	name: "Infinity",
	pos: position{line: 266, col: 1, offset: 7225},
	expr: &litMatcher{
	pos: position{line: 266, col: 12, offset: 7238},
	val: "Infinity",
	ignoreCase: false,
},
},
{
	name: "NaN",
	pos: position{line: 267, col: 1, offset: 7249},
	expr: &litMatcher{
	pos: position{line: 267, col: 7, offset: 7257},
	val: "NaN",
	ignoreCase: false,
},
},
{
	name: "Some",
	pos: position{line: 268, col: 1, offset: 7263},
	expr: &litMatcher{
	pos: position{line: 268, col: 8, offset: 7272},
	val: "Some",
	ignoreCase: false,
},
},
{
	name: "toMap",
	pos: position{line: 269, col: 1, offset: 7279},
	expr: &litMatcher{
	pos: position{line: 269, col: 9, offset: 7289},
	val: "toMap",
	ignoreCase: false,
},
},
{
	name: "Keyword",
	pos: position{line: 271, col: 1, offset: 7298},
	expr: &choiceExpr{
	pos: position{line: 272, col: 5, offset: 7314},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 272, col: 5, offset: 7314},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 272, col: 10, offset: 7319},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 272, col: 17, offset: 7326},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 273, col: 5, offset: 7335},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 273, col: 11, offset: 7341},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 274, col: 5, offset: 7348},
	name: "Using",
},
&ruleRefExpr{
	pos: position{line: 274, col: 13, offset: 7356},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 274, col: 23, offset: 7366},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 275, col: 5, offset: 7373},
	name: "True",
},
&ruleRefExpr{
	pos: position{line: 275, col: 12, offset: 7380},
	name: "False",
},
&ruleRefExpr{
	pos: position{line: 276, col: 5, offset: 7390},
	name: "Infinity",
},
&ruleRefExpr{
	pos: position{line: 276, col: 16, offset: 7401},
	name: "NaN",
},
&ruleRefExpr{
	pos: position{line: 277, col: 5, offset: 7409},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 277, col: 13, offset: 7417},
	name: "Some",
},
&ruleRefExpr{
	pos: position{line: 277, col: 20, offset: 7424},
	name: "toMap",
},
	},
},
},
{
	name: "Optional",
	pos: position{line: 279, col: 1, offset: 7431},
	expr: &litMatcher{
	pos: position{line: 279, col: 12, offset: 7444},
	val: "Optional",
	ignoreCase: false,
},
},
{
	name: "Text",
	pos: position{line: 280, col: 1, offset: 7455},
	expr: &litMatcher{
	pos: position{line: 280, col: 8, offset: 7464},
	val: "Text",
	ignoreCase: false,
},
},
{
	name: "List",
	pos: position{line: 281, col: 1, offset: 7471},
	expr: &litMatcher{
	pos: position{line: 281, col: 8, offset: 7480},
	val: "List",
	ignoreCase: false,
},
},
{
	name: "Location",
	pos: position{line: 282, col: 1, offset: 7487},
	expr: &litMatcher{
	pos: position{line: 282, col: 12, offset: 7500},
	val: "Location",
	ignoreCase: false,
},
},
{
	name: "Combine",
	pos: position{line: 284, col: 1, offset: 7512},
	expr: &choiceExpr{
	pos: position{line: 284, col: 11, offset: 7524},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 284, col: 11, offset: 7524},
	val: "/\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 284, col: 19, offset: 7532},
	val: "∧",
	ignoreCase: false,
},
	},
},
},
{
	name: "CombineTypes",
	pos: position{line: 285, col: 1, offset: 7538},
	expr: &choiceExpr{
	pos: position{line: 285, col: 16, offset: 7555},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 285, col: 16, offset: 7555},
	val: "//\\\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 285, col: 27, offset: 7566},
	val: "⩓",
	ignoreCase: false,
},
	},
},
},
{
	name: "Prefer",
	pos: position{line: 286, col: 1, offset: 7572},
	expr: &choiceExpr{
	pos: position{line: 286, col: 10, offset: 7583},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 286, col: 10, offset: 7583},
	val: "//",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 286, col: 17, offset: 7590},
	val: "⫽",
	ignoreCase: false,
},
	},
},
},
{
	name: "Lambda",
	pos: position{line: 287, col: 1, offset: 7596},
	expr: &choiceExpr{
	pos: position{line: 287, col: 10, offset: 7607},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 287, col: 10, offset: 7607},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 287, col: 17, offset: 7614},
	val: "λ",
	ignoreCase: false,
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 288, col: 1, offset: 7619},
	expr: &choiceExpr{
	pos: position{line: 288, col: 10, offset: 7630},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 288, col: 10, offset: 7630},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 288, col: 21, offset: 7641},
	val: "∀",
	ignoreCase: false,
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 289, col: 1, offset: 7647},
	expr: &choiceExpr{
	pos: position{line: 289, col: 9, offset: 7657},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 289, col: 9, offset: 7657},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 289, col: 16, offset: 7664},
	val: "→",
	ignoreCase: false,
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 291, col: 1, offset: 7671},
	expr: &seqExpr{
	pos: position{line: 291, col: 12, offset: 7684},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 291, col: 12, offset: 7684},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 291, col: 17, offset: 7689},
	expr: &charClassMatcher{
	pos: position{line: 291, col: 17, offset: 7689},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 291, col: 23, offset: 7695},
	expr: &ruleRefExpr{
	pos: position{line: 291, col: 23, offset: 7695},
	name: "Digit",
},
},
	},
},
},
{
	name: "NumericDoubleLiteral",
	pos: position{line: 293, col: 1, offset: 7703},
	expr: &actionExpr{
	pos: position{line: 293, col: 24, offset: 7728},
	run: (*parser).callonNumericDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 293, col: 24, offset: 7728},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 293, col: 24, offset: 7728},
	expr: &charClassMatcher{
	pos: position{line: 293, col: 24, offset: 7728},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 293, col: 30, offset: 7734},
	expr: &ruleRefExpr{
	pos: position{line: 293, col: 30, offset: 7734},
	name: "Digit",
},
},
&choiceExpr{
	pos: position{line: 293, col: 39, offset: 7743},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 293, col: 39, offset: 7743},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 293, col: 39, offset: 7743},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 293, col: 43, offset: 7747},
	expr: &ruleRefExpr{
	pos: position{line: 293, col: 43, offset: 7747},
	name: "Digit",
},
},
&zeroOrOneExpr{
	pos: position{line: 293, col: 50, offset: 7754},
	expr: &ruleRefExpr{
	pos: position{line: 293, col: 50, offset: 7754},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 293, col: 62, offset: 7766},
	name: "Exponent",
},
	},
},
	},
},
},
},
{
	name: "DoubleLiteral",
	pos: position{line: 301, col: 1, offset: 7922},
	expr: &choiceExpr{
	pos: position{line: 301, col: 17, offset: 7940},
	alternatives: []interface{}{
&labeledExpr{
	pos: position{line: 301, col: 17, offset: 7940},
	label: "d",
	expr: &ruleRefExpr{
	pos: position{line: 301, col: 19, offset: 7942},
	name: "NumericDoubleLiteral",
},
},
&actionExpr{
	pos: position{line: 302, col: 5, offset: 7967},
	run: (*parser).callonDoubleLiteral4,
	expr: &ruleRefExpr{
	pos: position{line: 302, col: 5, offset: 7967},
	name: "Infinity",
},
},
&actionExpr{
	pos: position{line: 303, col: 5, offset: 8019},
	run: (*parser).callonDoubleLiteral6,
	expr: &seqExpr{
	pos: position{line: 303, col: 5, offset: 8019},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 303, col: 5, offset: 8019},
	val: "-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 303, col: 9, offset: 8023},
	name: "Infinity",
},
	},
},
},
&actionExpr{
	pos: position{line: 304, col: 5, offset: 8076},
	run: (*parser).callonDoubleLiteral10,
	expr: &ruleRefExpr{
	pos: position{line: 304, col: 5, offset: 8076},
	name: "NaN",
},
},
	},
},
},
{
	name: "NaturalLiteral",
	pos: position{line: 306, col: 1, offset: 8119},
	expr: &actionExpr{
	pos: position{line: 306, col: 18, offset: 8138},
	run: (*parser).callonNaturalLiteral1,
	expr: &oneOrMoreExpr{
	pos: position{line: 306, col: 18, offset: 8138},
	expr: &ruleRefExpr{
	pos: position{line: 306, col: 18, offset: 8138},
	name: "Digit",
},
},
},
},
{
	name: "IntegerLiteral",
	pos: position{line: 311, col: 1, offset: 8227},
	expr: &actionExpr{
	pos: position{line: 311, col: 18, offset: 8246},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 311, col: 18, offset: 8246},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 311, col: 18, offset: 8246},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 311, col: 22, offset: 8250},
	name: "NaturalLiteral",
},
	},
},
},
},
{
	name: "DeBruijn",
	pos: position{line: 319, col: 1, offset: 8402},
	expr: &actionExpr{
	pos: position{line: 319, col: 12, offset: 8415},
	run: (*parser).callonDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 319, col: 12, offset: 8415},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 319, col: 12, offset: 8415},
	name: "_",
},
&litMatcher{
	pos: position{line: 319, col: 14, offset: 8417},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 319, col: 18, offset: 8421},
	name: "_",
},
&labeledExpr{
	pos: position{line: 319, col: 20, offset: 8423},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 319, col: 26, offset: 8429},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Variable",
	pos: position{line: 321, col: 1, offset: 8485},
	expr: &actionExpr{
	pos: position{line: 321, col: 12, offset: 8498},
	run: (*parser).callonVariable1,
	expr: &seqExpr{
	pos: position{line: 321, col: 12, offset: 8498},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 321, col: 12, offset: 8498},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 321, col: 17, offset: 8503},
	name: "NonreservedLabel",
},
},
&labeledExpr{
	pos: position{line: 321, col: 34, offset: 8520},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 321, col: 40, offset: 8526},
	expr: &ruleRefExpr{
	pos: position{line: 321, col: 40, offset: 8526},
	name: "DeBruijn",
},
},
},
	},
},
},
},
{
	name: "Identifier",
	pos: position{line: 329, col: 1, offset: 8689},
	expr: &choiceExpr{
	pos: position{line: 329, col: 14, offset: 8704},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 329, col: 14, offset: 8704},
	name: "Variable",
},
&ruleRefExpr{
	pos: position{line: 329, col: 25, offset: 8715},
	name: "Reserved",
},
	},
},
},
{
	name: "PathCharacter",
	pos: position{line: 331, col: 1, offset: 8725},
	expr: &choiceExpr{
	pos: position{line: 332, col: 6, offset: 8748},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 332, col: 6, offset: 8748},
	val: "!",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 333, col: 6, offset: 8760},
	val: "[\\x24-\\x27]",
	ranges: []rune{'$','\'',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 334, col: 6, offset: 8777},
	val: "[\\x2a-\\x2b]",
	ranges: []rune{'*','+',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 335, col: 6, offset: 8794},
	val: "[\\x2d-\\x2e]",
	ranges: []rune{'-','.',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 336, col: 6, offset: 8811},
	val: "[\\x30-\\x3b]",
	ranges: []rune{'0',';',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 337, col: 6, offset: 8828},
	val: "=",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 338, col: 6, offset: 8840},
	val: "[\\x40-\\x5a]",
	ranges: []rune{'@','Z',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 339, col: 6, offset: 8857},
	val: "[\\x5e-\\x7a]",
	ranges: []rune{'^','z',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 340, col: 6, offset: 8874},
	val: "|",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 341, col: 6, offset: 8886},
	val: "~",
	ignoreCase: false,
},
	},
},
},
{
	name: "QuotedPathCharacter",
	pos: position{line: 343, col: 1, offset: 8894},
	expr: &choiceExpr{
	pos: position{line: 344, col: 6, offset: 8923},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 344, col: 6, offset: 8923},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 345, col: 6, offset: 8940},
	val: "[\\x23-\\x2e]",
	ranges: []rune{'#','.',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 346, col: 6, offset: 8957},
	val: "[\\x30-\\x7f]",
	ranges: []rune{'0','\u007f',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 347, col: 6, offset: 8974},
	name: "ValidNonAscii",
},
	},
},
},
{
	name: "UnquotedPathComponent",
	pos: position{line: 349, col: 1, offset: 8989},
	expr: &actionExpr{
	pos: position{line: 349, col: 25, offset: 9015},
	run: (*parser).callonUnquotedPathComponent1,
	expr: &oneOrMoreExpr{
	pos: position{line: 349, col: 25, offset: 9015},
	expr: &ruleRefExpr{
	pos: position{line: 349, col: 25, offset: 9015},
	name: "PathCharacter",
},
},
},
},
{
	name: "QuotedPathComponent",
	pos: position{line: 350, col: 1, offset: 9061},
	expr: &actionExpr{
	pos: position{line: 350, col: 23, offset: 9085},
	run: (*parser).callonQuotedPathComponent1,
	expr: &oneOrMoreExpr{
	pos: position{line: 350, col: 23, offset: 9085},
	expr: &ruleRefExpr{
	pos: position{line: 350, col: 23, offset: 9085},
	name: "QuotedPathCharacter",
},
},
},
},
{
	name: "PathComponent",
	pos: position{line: 352, col: 1, offset: 9138},
	expr: &choiceExpr{
	pos: position{line: 352, col: 17, offset: 9156},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 352, col: 17, offset: 9156},
	run: (*parser).callonPathComponent2,
	expr: &seqExpr{
	pos: position{line: 352, col: 17, offset: 9156},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 352, col: 17, offset: 9156},
	val: "/",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 352, col: 21, offset: 9160},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 352, col: 23, offset: 9162},
	name: "UnquotedPathComponent",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 353, col: 17, offset: 9218},
	run: (*parser).callonPathComponent7,
	expr: &seqExpr{
	pos: position{line: 353, col: 17, offset: 9218},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 353, col: 17, offset: 9218},
	val: "/",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 353, col: 21, offset: 9222},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 353, col: 25, offset: 9226},
	label: "q",
	expr: &ruleRefExpr{
	pos: position{line: 353, col: 27, offset: 9228},
	name: "QuotedPathComponent",
},
},
&litMatcher{
	pos: position{line: 353, col: 47, offset: 9248},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
	},
},
},
{
	name: "Path",
	pos: position{line: 355, col: 1, offset: 9271},
	expr: &actionExpr{
	pos: position{line: 355, col: 8, offset: 9280},
	run: (*parser).callonPath1,
	expr: &labeledExpr{
	pos: position{line: 355, col: 8, offset: 9280},
	label: "cs",
	expr: &oneOrMoreExpr{
	pos: position{line: 355, col: 11, offset: 9283},
	expr: &ruleRefExpr{
	pos: position{line: 355, col: 11, offset: 9283},
	name: "PathComponent",
},
},
},
},
},
{
	name: "Local",
	pos: position{line: 364, col: 1, offset: 9557},
	expr: &choiceExpr{
	pos: position{line: 364, col: 9, offset: 9567},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 364, col: 9, offset: 9567},
	name: "ParentPath",
},
&ruleRefExpr{
	pos: position{line: 364, col: 22, offset: 9580},
	name: "HerePath",
},
&ruleRefExpr{
	pos: position{line: 364, col: 33, offset: 9591},
	name: "HomePath",
},
&ruleRefExpr{
	pos: position{line: 364, col: 44, offset: 9602},
	name: "AbsolutePath",
},
	},
},
},
{
	name: "ParentPath",
	pos: position{line: 366, col: 1, offset: 9616},
	expr: &actionExpr{
	pos: position{line: 366, col: 14, offset: 9631},
	run: (*parser).callonParentPath1,
	expr: &seqExpr{
	pos: position{line: 366, col: 14, offset: 9631},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 366, col: 14, offset: 9631},
	val: "..",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 366, col: 19, offset: 9636},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 366, col: 21, offset: 9638},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HerePath",
	pos: position{line: 367, col: 1, offset: 9694},
	expr: &actionExpr{
	pos: position{line: 367, col: 12, offset: 9707},
	run: (*parser).callonHerePath1,
	expr: &seqExpr{
	pos: position{line: 367, col: 12, offset: 9707},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 367, col: 12, offset: 9707},
	val: ".",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 367, col: 16, offset: 9711},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 367, col: 18, offset: 9713},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HomePath",
	pos: position{line: 368, col: 1, offset: 9752},
	expr: &actionExpr{
	pos: position{line: 368, col: 12, offset: 9765},
	run: (*parser).callonHomePath1,
	expr: &seqExpr{
	pos: position{line: 368, col: 12, offset: 9765},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 368, col: 12, offset: 9765},
	val: "~",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 368, col: 16, offset: 9769},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 368, col: 18, offset: 9771},
	name: "Path",
},
},
	},
},
},
},
{
	name: "AbsolutePath",
	pos: position{line: 369, col: 1, offset: 9826},
	expr: &actionExpr{
	pos: position{line: 369, col: 16, offset: 9843},
	run: (*parser).callonAbsolutePath1,
	expr: &labeledExpr{
	pos: position{line: 369, col: 16, offset: 9843},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 369, col: 18, offset: 9845},
	name: "Path",
},
},
},
},
{
	name: "Scheme",
	pos: position{line: 371, col: 1, offset: 9901},
	expr: &seqExpr{
	pos: position{line: 371, col: 10, offset: 9912},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 371, col: 10, offset: 9912},
	val: "http",
	ignoreCase: false,
},
&zeroOrOneExpr{
	pos: position{line: 371, col: 17, offset: 9919},
	expr: &litMatcher{
	pos: position{line: 371, col: 17, offset: 9919},
	val: "s",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "HttpRaw",
	pos: position{line: 373, col: 1, offset: 9925},
	expr: &actionExpr{
	pos: position{line: 373, col: 11, offset: 9937},
	run: (*parser).callonHttpRaw1,
	expr: &seqExpr{
	pos: position{line: 373, col: 11, offset: 9937},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 373, col: 11, offset: 9937},
	name: "Scheme",
},
&litMatcher{
	pos: position{line: 373, col: 18, offset: 9944},
	val: "://",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 373, col: 24, offset: 9950},
	name: "Authority",
},
&ruleRefExpr{
	pos: position{line: 373, col: 34, offset: 9960},
	name: "UrlPath",
},
&zeroOrOneExpr{
	pos: position{line: 373, col: 42, offset: 9968},
	expr: &seqExpr{
	pos: position{line: 373, col: 44, offset: 9970},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 373, col: 44, offset: 9970},
	val: "?",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 373, col: 48, offset: 9974},
	name: "Query",
},
	},
},
},
	},
},
},
},
{
	name: "UrlPath",
	pos: position{line: 375, col: 1, offset: 10031},
	expr: &zeroOrMoreExpr{
	pos: position{line: 375, col: 11, offset: 10043},
	expr: &choiceExpr{
	pos: position{line: 375, col: 12, offset: 10044},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 375, col: 12, offset: 10044},
	name: "PathComponent",
},
&seqExpr{
	pos: position{line: 375, col: 28, offset: 10060},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 375, col: 28, offset: 10060},
	val: "/",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 375, col: 32, offset: 10064},
	name: "Segment",
},
	},
},
	},
},
},
},
{
	name: "Authority",
	pos: position{line: 377, col: 1, offset: 10075},
	expr: &seqExpr{
	pos: position{line: 377, col: 13, offset: 10089},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 377, col: 13, offset: 10089},
	expr: &seqExpr{
	pos: position{line: 377, col: 14, offset: 10090},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 377, col: 14, offset: 10090},
	name: "Userinfo",
},
&litMatcher{
	pos: position{line: 377, col: 23, offset: 10099},
	val: "@",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 377, col: 29, offset: 10105},
	name: "Host",
},
&zeroOrOneExpr{
	pos: position{line: 377, col: 34, offset: 10110},
	expr: &seqExpr{
	pos: position{line: 377, col: 35, offset: 10111},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 377, col: 35, offset: 10111},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 377, col: 39, offset: 10115},
	name: "Port",
},
	},
},
},
	},
},
},
{
	name: "Userinfo",
	pos: position{line: 379, col: 1, offset: 10123},
	expr: &zeroOrMoreExpr{
	pos: position{line: 379, col: 12, offset: 10136},
	expr: &choiceExpr{
	pos: position{line: 379, col: 14, offset: 10138},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 379, col: 14, offset: 10138},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 379, col: 27, offset: 10151},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 379, col: 40, offset: 10164},
	name: "SubDelims",
},
&litMatcher{
	pos: position{line: 379, col: 52, offset: 10176},
	val: ":",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "Host",
	pos: position{line: 381, col: 1, offset: 10184},
	expr: &choiceExpr{
	pos: position{line: 381, col: 8, offset: 10193},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 381, col: 8, offset: 10193},
	name: "IPLiteral",
},
&ruleRefExpr{
	pos: position{line: 381, col: 20, offset: 10205},
	name: "RegName",
},
	},
},
},
{
	name: "Port",
	pos: position{line: 383, col: 1, offset: 10214},
	expr: &zeroOrMoreExpr{
	pos: position{line: 383, col: 8, offset: 10223},
	expr: &ruleRefExpr{
	pos: position{line: 383, col: 8, offset: 10223},
	name: "Digit",
},
},
},
{
	name: "IPLiteral",
	pos: position{line: 385, col: 1, offset: 10231},
	expr: &seqExpr{
	pos: position{line: 385, col: 13, offset: 10245},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 385, col: 13, offset: 10245},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 385, col: 17, offset: 10249},
	name: "IPv6address",
},
&litMatcher{
	pos: position{line: 385, col: 29, offset: 10261},
	val: "]",
	ignoreCase: false,
},
	},
},
},
{
	name: "IPv6address",
	pos: position{line: 387, col: 1, offset: 10266},
	expr: &actionExpr{
	pos: position{line: 387, col: 15, offset: 10282},
	run: (*parser).callonIPv6address1,
	expr: &seqExpr{
	pos: position{line: 387, col: 15, offset: 10282},
	exprs: []interface{}{
&zeroOrMoreExpr{
	pos: position{line: 387, col: 15, offset: 10282},
	expr: &ruleRefExpr{
	pos: position{line: 387, col: 16, offset: 10283},
	name: "HexDig",
},
},
&litMatcher{
	pos: position{line: 387, col: 25, offset: 10292},
	val: ":",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 387, col: 29, offset: 10296},
	expr: &choiceExpr{
	pos: position{line: 387, col: 30, offset: 10297},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 387, col: 30, offset: 10297},
	name: "HexDig",
},
&litMatcher{
	pos: position{line: 387, col: 39, offset: 10306},
	val: ":",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 387, col: 45, offset: 10312},
	val: ".",
	ignoreCase: false,
},
	},
},
},
	},
},
},
},
{
	name: "RegName",
	pos: position{line: 393, col: 1, offset: 10466},
	expr: &zeroOrMoreExpr{
	pos: position{line: 393, col: 11, offset: 10478},
	expr: &choiceExpr{
	pos: position{line: 393, col: 12, offset: 10479},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 393, col: 12, offset: 10479},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 393, col: 25, offset: 10492},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 393, col: 38, offset: 10505},
	name: "SubDelims",
},
	},
},
},
},
{
	name: "Segment",
	pos: position{line: 395, col: 1, offset: 10518},
	expr: &zeroOrMoreExpr{
	pos: position{line: 395, col: 11, offset: 10530},
	expr: &ruleRefExpr{
	pos: position{line: 395, col: 11, offset: 10530},
	name: "PChar",
},
},
},
{
	name: "PChar",
	pos: position{line: 397, col: 1, offset: 10538},
	expr: &choiceExpr{
	pos: position{line: 397, col: 9, offset: 10548},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 397, col: 9, offset: 10548},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 397, col: 22, offset: 10561},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 397, col: 35, offset: 10574},
	name: "SubDelims",
},
&charClassMatcher{
	pos: position{line: 397, col: 47, offset: 10586},
	val: "[:@]",
	chars: []rune{':','@',},
	ignoreCase: false,
	inverted: false,
},
	},
},
},
{
	name: "Query",
	pos: position{line: 399, col: 1, offset: 10592},
	expr: &zeroOrMoreExpr{
	pos: position{line: 399, col: 9, offset: 10602},
	expr: &choiceExpr{
	pos: position{line: 399, col: 10, offset: 10603},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 399, col: 10, offset: 10603},
	name: "PChar",
},
&charClassMatcher{
	pos: position{line: 399, col: 18, offset: 10611},
	val: "[/?]",
	chars: []rune{'/','?',},
	ignoreCase: false,
	inverted: false,
},
	},
},
},
},
{
	name: "PctEncoded",
	pos: position{line: 401, col: 1, offset: 10619},
	expr: &seqExpr{
	pos: position{line: 401, col: 14, offset: 10634},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 401, col: 14, offset: 10634},
	val: "%",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 401, col: 18, offset: 10638},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 401, col: 25, offset: 10645},
	name: "HexDig",
},
	},
},
},
{
	name: "Unreserved",
	pos: position{line: 403, col: 1, offset: 10653},
	expr: &charClassMatcher{
	pos: position{line: 403, col: 14, offset: 10668},
	val: "[A-Za-z0-9._~-]",
	chars: []rune{'.','_','~','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SubDelims",
	pos: position{line: 405, col: 1, offset: 10685},
	expr: &choiceExpr{
	pos: position{line: 405, col: 13, offset: 10699},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 405, col: 13, offset: 10699},
	val: "!",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 405, col: 19, offset: 10705},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 405, col: 25, offset: 10711},
	val: "&",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 405, col: 31, offset: 10717},
	val: "'",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 405, col: 37, offset: 10723},
	val: "*",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 405, col: 43, offset: 10729},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 405, col: 49, offset: 10735},
	val: ";",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 405, col: 55, offset: 10741},
	val: "=",
	ignoreCase: false,
},
	},
},
},
{
	name: "Http",
	pos: position{line: 407, col: 1, offset: 10746},
	expr: &actionExpr{
	pos: position{line: 407, col: 8, offset: 10755},
	run: (*parser).callonHttp1,
	expr: &labeledExpr{
	pos: position{line: 407, col: 8, offset: 10755},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 407, col: 10, offset: 10757},
	name: "HttpRaw",
},
},
},
},
{
	name: "Env",
	pos: position{line: 409, col: 1, offset: 10807},
	expr: &actionExpr{
	pos: position{line: 409, col: 7, offset: 10815},
	run: (*parser).callonEnv1,
	expr: &seqExpr{
	pos: position{line: 409, col: 7, offset: 10815},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 409, col: 7, offset: 10815},
	val: "env:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 409, col: 14, offset: 10822},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 409, col: 17, offset: 10825},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 409, col: 17, offset: 10825},
	name: "BashEnvironmentVariable",
},
&ruleRefExpr{
	pos: position{line: 409, col: 43, offset: 10851},
	name: "PosixEnvironmentVariable",
},
	},
},
},
	},
},
},
},
{
	name: "BashEnvironmentVariable",
	pos: position{line: 411, col: 1, offset: 10896},
	expr: &actionExpr{
	pos: position{line: 411, col: 27, offset: 10924},
	run: (*parser).callonBashEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 411, col: 27, offset: 10924},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 411, col: 27, offset: 10924},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 411, col: 36, offset: 10933},
	expr: &charClassMatcher{
	pos: position{line: 411, col: 36, offset: 10933},
	val: "[A-Za-z0-9_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
	},
},
},
},
{
	name: "PosixEnvironmentVariable",
	pos: position{line: 415, col: 1, offset: 10989},
	expr: &actionExpr{
	pos: position{line: 415, col: 28, offset: 11018},
	run: (*parser).callonPosixEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 415, col: 28, offset: 11018},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 415, col: 28, offset: 11018},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 415, col: 32, offset: 11022},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 415, col: 34, offset: 11024},
	name: "PosixEnvironmentVariableContent",
},
},
&litMatcher{
	pos: position{line: 415, col: 66, offset: 11056},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PosixEnvironmentVariableContent",
	pos: position{line: 419, col: 1, offset: 11081},
	expr: &actionExpr{
	pos: position{line: 419, col: 35, offset: 11117},
	run: (*parser).callonPosixEnvironmentVariableContent1,
	expr: &labeledExpr{
	pos: position{line: 419, col: 35, offset: 11117},
	label: "v",
	expr: &oneOrMoreExpr{
	pos: position{line: 419, col: 37, offset: 11119},
	expr: &ruleRefExpr{
	pos: position{line: 419, col: 37, offset: 11119},
	name: "PosixEnvironmentVariableCharacter",
},
},
},
},
},
{
	name: "PosixEnvironmentVariableCharacter",
	pos: position{line: 428, col: 1, offset: 11332},
	expr: &choiceExpr{
	pos: position{line: 429, col: 7, offset: 11376},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 429, col: 7, offset: 11376},
	run: (*parser).callonPosixEnvironmentVariableCharacter2,
	expr: &litMatcher{
	pos: position{line: 429, col: 7, offset: 11376},
	val: "\\\"",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 430, col: 7, offset: 11416},
	run: (*parser).callonPosixEnvironmentVariableCharacter4,
	expr: &litMatcher{
	pos: position{line: 430, col: 7, offset: 11416},
	val: "\\\\",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 431, col: 7, offset: 11456},
	run: (*parser).callonPosixEnvironmentVariableCharacter6,
	expr: &litMatcher{
	pos: position{line: 431, col: 7, offset: 11456},
	val: "\\a",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 432, col: 7, offset: 11496},
	run: (*parser).callonPosixEnvironmentVariableCharacter8,
	expr: &litMatcher{
	pos: position{line: 432, col: 7, offset: 11496},
	val: "\\b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 433, col: 7, offset: 11536},
	run: (*parser).callonPosixEnvironmentVariableCharacter10,
	expr: &litMatcher{
	pos: position{line: 433, col: 7, offset: 11536},
	val: "\\f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 434, col: 7, offset: 11576},
	run: (*parser).callonPosixEnvironmentVariableCharacter12,
	expr: &litMatcher{
	pos: position{line: 434, col: 7, offset: 11576},
	val: "\\n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 435, col: 7, offset: 11616},
	run: (*parser).callonPosixEnvironmentVariableCharacter14,
	expr: &litMatcher{
	pos: position{line: 435, col: 7, offset: 11616},
	val: "\\r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 436, col: 7, offset: 11656},
	run: (*parser).callonPosixEnvironmentVariableCharacter16,
	expr: &litMatcher{
	pos: position{line: 436, col: 7, offset: 11656},
	val: "\\t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 437, col: 7, offset: 11696},
	run: (*parser).callonPosixEnvironmentVariableCharacter18,
	expr: &litMatcher{
	pos: position{line: 437, col: 7, offset: 11696},
	val: "\\v",
	ignoreCase: false,
},
},
&charClassMatcher{
	pos: position{line: 438, col: 7, offset: 11736},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 439, col: 7, offset: 11754},
	val: "[\\x23-\\x3c]",
	ranges: []rune{'#','<',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 440, col: 7, offset: 11772},
	val: "[\\x3e-\\x5b]",
	ranges: []rune{'>','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 441, col: 7, offset: 11790},
	val: "[\\x5d-\\x7e]",
	ranges: []rune{']','~',},
	ignoreCase: false,
	inverted: false,
},
	},
},
},
{
	name: "ImportType",
	pos: position{line: 443, col: 1, offset: 11803},
	expr: &choiceExpr{
	pos: position{line: 443, col: 14, offset: 11818},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 443, col: 14, offset: 11818},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 443, col: 24, offset: 11828},
	name: "Local",
},
&ruleRefExpr{
	pos: position{line: 443, col: 32, offset: 11836},
	name: "Http",
},
&ruleRefExpr{
	pos: position{line: 443, col: 39, offset: 11843},
	name: "Env",
},
	},
},
},
{
	name: "HashValue",
	pos: position{line: 446, col: 1, offset: 11916},
	expr: &actionExpr{
	pos: position{line: 446, col: 13, offset: 11928},
	run: (*parser).callonHashValue1,
	expr: &seqExpr{
	pos: position{line: 446, col: 13, offset: 11928},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 446, col: 13, offset: 11928},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 446, col: 20, offset: 11935},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 446, col: 27, offset: 11942},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 446, col: 34, offset: 11949},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 446, col: 41, offset: 11956},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 446, col: 48, offset: 11963},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 446, col: 55, offset: 11970},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 446, col: 62, offset: 11977},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 447, col: 13, offset: 11996},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 447, col: 20, offset: 12003},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 447, col: 27, offset: 12010},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 447, col: 34, offset: 12017},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 447, col: 41, offset: 12024},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 447, col: 48, offset: 12031},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 447, col: 55, offset: 12038},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 447, col: 62, offset: 12045},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 448, col: 13, offset: 12064},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 448, col: 20, offset: 12071},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 448, col: 27, offset: 12078},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 448, col: 34, offset: 12085},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 448, col: 41, offset: 12092},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 448, col: 48, offset: 12099},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 448, col: 55, offset: 12106},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 448, col: 62, offset: 12113},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 449, col: 13, offset: 12132},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 449, col: 20, offset: 12139},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 449, col: 27, offset: 12146},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 449, col: 34, offset: 12153},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 449, col: 41, offset: 12160},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 449, col: 48, offset: 12167},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 449, col: 55, offset: 12174},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 449, col: 62, offset: 12181},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 450, col: 13, offset: 12200},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 450, col: 20, offset: 12207},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 450, col: 27, offset: 12214},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 450, col: 34, offset: 12221},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 450, col: 41, offset: 12228},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 450, col: 48, offset: 12235},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 450, col: 55, offset: 12242},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 450, col: 62, offset: 12249},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 451, col: 13, offset: 12268},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 451, col: 20, offset: 12275},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 451, col: 27, offset: 12282},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 451, col: 34, offset: 12289},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 451, col: 41, offset: 12296},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 451, col: 48, offset: 12303},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 451, col: 55, offset: 12310},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 451, col: 62, offset: 12317},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 452, col: 13, offset: 12336},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 452, col: 20, offset: 12343},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 452, col: 27, offset: 12350},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 452, col: 34, offset: 12357},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 452, col: 41, offset: 12364},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 452, col: 48, offset: 12371},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 452, col: 55, offset: 12378},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 452, col: 62, offset: 12385},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 453, col: 13, offset: 12404},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 453, col: 20, offset: 12411},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 453, col: 27, offset: 12418},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 453, col: 34, offset: 12425},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 453, col: 41, offset: 12432},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 453, col: 48, offset: 12439},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 453, col: 55, offset: 12446},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 453, col: 62, offset: 12453},
	name: "HexDig",
},
	},
},
},
},
{
	name: "Hash",
	pos: position{line: 459, col: 1, offset: 12597},
	expr: &actionExpr{
	pos: position{line: 459, col: 8, offset: 12604},
	run: (*parser).callonHash1,
	expr: &seqExpr{
	pos: position{line: 459, col: 8, offset: 12604},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 459, col: 8, offset: 12604},
	val: "sha256:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 459, col: 18, offset: 12614},
	label: "val",
	expr: &ruleRefExpr{
	pos: position{line: 459, col: 22, offset: 12618},
	name: "HashValue",
},
},
	},
},
},
},
{
	name: "ImportHashed",
	pos: position{line: 461, col: 1, offset: 12688},
	expr: &actionExpr{
	pos: position{line: 461, col: 16, offset: 12705},
	run: (*parser).callonImportHashed1,
	expr: &seqExpr{
	pos: position{line: 461, col: 16, offset: 12705},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 461, col: 16, offset: 12705},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 461, col: 18, offset: 12707},
	name: "ImportType",
},
},
&labeledExpr{
	pos: position{line: 461, col: 29, offset: 12718},
	label: "h",
	expr: &zeroOrOneExpr{
	pos: position{line: 461, col: 31, offset: 12720},
	expr: &seqExpr{
	pos: position{line: 461, col: 32, offset: 12721},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 461, col: 32, offset: 12721},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 461, col: 35, offset: 12724},
	name: "Hash",
},
	},
},
},
},
	},
},
},
},
{
	name: "Import",
	pos: position{line: 469, col: 1, offset: 12879},
	expr: &choiceExpr{
	pos: position{line: 469, col: 10, offset: 12890},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 469, col: 10, offset: 12890},
	run: (*parser).callonImport2,
	expr: &seqExpr{
	pos: position{line: 469, col: 10, offset: 12890},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 469, col: 10, offset: 12890},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 469, col: 12, offset: 12892},
	name: "ImportHashed",
},
},
&ruleRefExpr{
	pos: position{line: 469, col: 25, offset: 12905},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 469, col: 27, offset: 12907},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 469, col: 30, offset: 12910},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 469, col: 33, offset: 12913},
	name: "Text",
},
	},
},
},
&actionExpr{
	pos: position{line: 470, col: 10, offset: 13003},
	run: (*parser).callonImport10,
	expr: &seqExpr{
	pos: position{line: 470, col: 10, offset: 13003},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 470, col: 10, offset: 13003},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 470, col: 12, offset: 13005},
	name: "ImportHashed",
},
},
&ruleRefExpr{
	pos: position{line: 470, col: 25, offset: 13018},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 470, col: 27, offset: 13020},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 470, col: 30, offset: 13023},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 470, col: 33, offset: 13026},
	name: "Location",
},
	},
},
},
&actionExpr{
	pos: position{line: 471, col: 10, offset: 13121},
	run: (*parser).callonImport18,
	expr: &labeledExpr{
	pos: position{line: 471, col: 10, offset: 13121},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 471, col: 12, offset: 13123},
	name: "ImportHashed",
},
},
},
	},
},
},
{
	name: "LetBinding",
	pos: position{line: 474, col: 1, offset: 13211},
	expr: &actionExpr{
	pos: position{line: 474, col: 14, offset: 13226},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 474, col: 14, offset: 13226},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 474, col: 14, offset: 13226},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 474, col: 18, offset: 13230},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 474, col: 21, offset: 13233},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 474, col: 27, offset: 13239},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 474, col: 44, offset: 13256},
	name: "_",
},
&labeledExpr{
	pos: position{line: 474, col: 46, offset: 13258},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 474, col: 48, offset: 13260},
	expr: &seqExpr{
	pos: position{line: 474, col: 49, offset: 13261},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 474, col: 49, offset: 13261},
	name: "Annotation",
},
&ruleRefExpr{
	pos: position{line: 474, col: 60, offset: 13272},
	name: "_",
},
	},
},
},
},
&litMatcher{
	pos: position{line: 475, col: 13, offset: 13288},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 475, col: 17, offset: 13292},
	name: "_",
},
&labeledExpr{
	pos: position{line: 475, col: 19, offset: 13294},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 475, col: 21, offset: 13296},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 475, col: 32, offset: 13307},
	name: "_",
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 490, col: 1, offset: 13616},
	expr: &choiceExpr{
	pos: position{line: 491, col: 7, offset: 13637},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 491, col: 7, offset: 13637},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 491, col: 7, offset: 13637},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 491, col: 7, offset: 13637},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 491, col: 14, offset: 13644},
	name: "_",
},
&litMatcher{
	pos: position{line: 491, col: 16, offset: 13646},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 491, col: 20, offset: 13650},
	name: "_",
},
&labeledExpr{
	pos: position{line: 491, col: 22, offset: 13652},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 491, col: 28, offset: 13658},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 491, col: 45, offset: 13675},
	name: "_",
},
&litMatcher{
	pos: position{line: 491, col: 47, offset: 13677},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 491, col: 51, offset: 13681},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 491, col: 54, offset: 13684},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 491, col: 56, offset: 13686},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 491, col: 67, offset: 13697},
	name: "_",
},
&litMatcher{
	pos: position{line: 491, col: 69, offset: 13699},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 491, col: 73, offset: 13703},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 491, col: 75, offset: 13705},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 491, col: 81, offset: 13711},
	name: "_",
},
&labeledExpr{
	pos: position{line: 491, col: 83, offset: 13713},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 491, col: 88, offset: 13718},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 494, col: 7, offset: 13835},
	run: (*parser).callonExpression22,
	expr: &seqExpr{
	pos: position{line: 494, col: 7, offset: 13835},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 494, col: 7, offset: 13835},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 494, col: 10, offset: 13838},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 494, col: 13, offset: 13841},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 494, col: 18, offset: 13846},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 494, col: 29, offset: 13857},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 494, col: 31, offset: 13859},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 494, col: 36, offset: 13864},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 494, col: 39, offset: 13867},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 494, col: 41, offset: 13869},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 494, col: 52, offset: 13880},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 494, col: 54, offset: 13882},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 494, col: 59, offset: 13887},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 494, col: 62, offset: 13890},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 494, col: 64, offset: 13892},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 497, col: 7, offset: 13978},
	run: (*parser).callonExpression38,
	expr: &seqExpr{
	pos: position{line: 497, col: 7, offset: 13978},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 497, col: 7, offset: 13978},
	label: "bindings",
	expr: &oneOrMoreExpr{
	pos: position{line: 497, col: 16, offset: 13987},
	expr: &ruleRefExpr{
	pos: position{line: 497, col: 16, offset: 13987},
	name: "LetBinding",
},
},
},
&ruleRefExpr{
	pos: position{line: 497, col: 28, offset: 13999},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 497, col: 31, offset: 14002},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 497, col: 34, offset: 14005},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 497, col: 36, offset: 14007},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 504, col: 7, offset: 14247},
	run: (*parser).callonExpression47,
	expr: &seqExpr{
	pos: position{line: 504, col: 7, offset: 14247},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 504, col: 7, offset: 14247},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 504, col: 14, offset: 14254},
	name: "_",
},
&litMatcher{
	pos: position{line: 504, col: 16, offset: 14256},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 504, col: 20, offset: 14260},
	name: "_",
},
&labeledExpr{
	pos: position{line: 504, col: 22, offset: 14262},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 504, col: 28, offset: 14268},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 504, col: 45, offset: 14285},
	name: "_",
},
&litMatcher{
	pos: position{line: 504, col: 47, offset: 14287},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 504, col: 51, offset: 14291},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 504, col: 54, offset: 14294},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 504, col: 56, offset: 14296},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 504, col: 67, offset: 14307},
	name: "_",
},
&litMatcher{
	pos: position{line: 504, col: 69, offset: 14309},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 504, col: 73, offset: 14313},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 504, col: 75, offset: 14315},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 504, col: 81, offset: 14321},
	name: "_",
},
&labeledExpr{
	pos: position{line: 504, col: 83, offset: 14323},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 504, col: 88, offset: 14328},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 507, col: 7, offset: 14437},
	run: (*parser).callonExpression67,
	expr: &seqExpr{
	pos: position{line: 507, col: 7, offset: 14437},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 507, col: 7, offset: 14437},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 507, col: 9, offset: 14439},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 507, col: 28, offset: 14458},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 507, col: 30, offset: 14460},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 507, col: 36, offset: 14466},
	name: "_",
},
&labeledExpr{
	pos: position{line: 507, col: 38, offset: 14468},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 507, col: 40, offset: 14470},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 508, col: 7, offset: 14529},
	run: (*parser).callonExpression76,
	expr: &seqExpr{
	pos: position{line: 508, col: 7, offset: 14529},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 508, col: 7, offset: 14529},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 508, col: 13, offset: 14535},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 508, col: 16, offset: 14538},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 508, col: 18, offset: 14540},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 508, col: 35, offset: 14557},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 508, col: 38, offset: 14560},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 508, col: 40, offset: 14562},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 508, col: 57, offset: 14579},
	name: "_",
},
&litMatcher{
	pos: position{line: 508, col: 59, offset: 14581},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 508, col: 63, offset: 14585},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 508, col: 66, offset: 14588},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 508, col: 68, offset: 14590},
	name: "ApplicationExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 511, col: 7, offset: 14711},
	name: "EmptyList",
},
&actionExpr{
	pos: position{line: 512, col: 7, offset: 14727},
	run: (*parser).callonExpression91,
	expr: &seqExpr{
	pos: position{line: 512, col: 7, offset: 14727},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 512, col: 7, offset: 14727},
	name: "toMap",
},
&ruleRefExpr{
	pos: position{line: 512, col: 13, offset: 14733},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 512, col: 16, offset: 14736},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 512, col: 18, offset: 14738},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 512, col: 35, offset: 14755},
	name: "_",
},
&litMatcher{
	pos: position{line: 512, col: 37, offset: 14757},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 512, col: 41, offset: 14761},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 512, col: 44, offset: 14764},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 512, col: 46, offset: 14766},
	name: "ApplicationExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 513, col: 7, offset: 14836},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 515, col: 1, offset: 14857},
	expr: &actionExpr{
	pos: position{line: 515, col: 14, offset: 14872},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 515, col: 14, offset: 14872},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 515, col: 14, offset: 14872},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 515, col: 18, offset: 14876},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 515, col: 21, offset: 14879},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 515, col: 23, offset: 14881},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 517, col: 1, offset: 14911},
	expr: &actionExpr{
	pos: position{line: 518, col: 1, offset: 14935},
	run: (*parser).callonAnnotatedExpression1,
	expr: &seqExpr{
	pos: position{line: 518, col: 1, offset: 14935},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 518, col: 1, offset: 14935},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 518, col: 3, offset: 14937},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 518, col: 22, offset: 14956},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 518, col: 24, offset: 14958},
	expr: &seqExpr{
	pos: position{line: 518, col: 25, offset: 14959},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 518, col: 25, offset: 14959},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 518, col: 27, offset: 14961},
	name: "Annotation",
},
	},
},
},
},
	},
},
},
},
{
	name: "EmptyList",
	pos: position{line: 523, col: 1, offset: 15086},
	expr: &actionExpr{
	pos: position{line: 523, col: 13, offset: 15100},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 523, col: 13, offset: 15100},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 523, col: 13, offset: 15100},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 523, col: 17, offset: 15104},
	name: "_",
},
&litMatcher{
	pos: position{line: 523, col: 19, offset: 15106},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 523, col: 23, offset: 15110},
	name: "_",
},
&litMatcher{
	pos: position{line: 523, col: 25, offset: 15112},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 523, col: 29, offset: 15116},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 523, col: 32, offset: 15119},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 523, col: 34, offset: 15121},
	name: "ApplicationExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 527, col: 1, offset: 15189},
	expr: &ruleRefExpr{
	pos: position{line: 527, col: 22, offset: 15212},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 529, col: 1, offset: 15233},
	expr: &actionExpr{
	pos: position{line: 529, col: 26, offset: 15260},
	run: (*parser).callonImportAltExpression1,
	expr: &seqExpr{
	pos: position{line: 529, col: 26, offset: 15260},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 529, col: 26, offset: 15260},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 529, col: 32, offset: 15266},
	name: "OrExpression",
},
},
&labeledExpr{
	pos: position{line: 529, col: 55, offset: 15289},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 529, col: 60, offset: 15294},
	expr: &seqExpr{
	pos: position{line: 529, col: 61, offset: 15295},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 529, col: 61, offset: 15295},
	name: "_",
},
&litMatcher{
	pos: position{line: 529, col: 63, offset: 15297},
	val: "?",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 529, col: 67, offset: 15301},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 529, col: 69, offset: 15303},
	name: "OrExpression",
},
	},
},
},
},
	},
},
},
},
{
	name: "OrExpression",
	pos: position{line: 531, col: 1, offset: 15374},
	expr: &actionExpr{
	pos: position{line: 531, col: 26, offset: 15401},
	run: (*parser).callonOrExpression1,
	expr: &seqExpr{
	pos: position{line: 531, col: 26, offset: 15401},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 531, col: 26, offset: 15401},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 531, col: 32, offset: 15407},
	name: "PlusExpression",
},
},
&labeledExpr{
	pos: position{line: 531, col: 55, offset: 15430},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 531, col: 60, offset: 15435},
	expr: &seqExpr{
	pos: position{line: 531, col: 61, offset: 15436},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 531, col: 61, offset: 15436},
	name: "_",
},
&litMatcher{
	pos: position{line: 531, col: 63, offset: 15438},
	val: "||",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 531, col: 68, offset: 15443},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 531, col: 70, offset: 15445},
	name: "PlusExpression",
},
	},
},
},
},
	},
},
},
},
{
	name: "PlusExpression",
	pos: position{line: 533, col: 1, offset: 15511},
	expr: &actionExpr{
	pos: position{line: 533, col: 26, offset: 15538},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 533, col: 26, offset: 15538},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 533, col: 26, offset: 15538},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 533, col: 32, offset: 15544},
	name: "TextAppendExpression",
},
},
&labeledExpr{
	pos: position{line: 533, col: 55, offset: 15567},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 533, col: 60, offset: 15572},
	expr: &seqExpr{
	pos: position{line: 533, col: 61, offset: 15573},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 533, col: 61, offset: 15573},
	name: "_",
},
&litMatcher{
	pos: position{line: 533, col: 63, offset: 15575},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 533, col: 67, offset: 15579},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 533, col: 70, offset: 15582},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 533, col: 72, offset: 15584},
	name: "TextAppendExpression",
},
},
	},
},
},
},
	},
},
},
},
{
	name: "TextAppendExpression",
	pos: position{line: 535, col: 1, offset: 15658},
	expr: &actionExpr{
	pos: position{line: 535, col: 26, offset: 15685},
	run: (*parser).callonTextAppendExpression1,
	expr: &seqExpr{
	pos: position{line: 535, col: 26, offset: 15685},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 535, col: 26, offset: 15685},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 535, col: 32, offset: 15691},
	name: "ListAppendExpression",
},
},
&labeledExpr{
	pos: position{line: 535, col: 55, offset: 15714},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 535, col: 60, offset: 15719},
	expr: &seqExpr{
	pos: position{line: 535, col: 61, offset: 15720},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 535, col: 61, offset: 15720},
	name: "_",
},
&litMatcher{
	pos: position{line: 535, col: 63, offset: 15722},
	val: "++",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 535, col: 68, offset: 15727},
	name: "_",
},
&labeledExpr{
	pos: position{line: 535, col: 70, offset: 15729},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 535, col: 72, offset: 15731},
	name: "ListAppendExpression",
},
},
	},
},
},
},
	},
},
},
},
{
	name: "ListAppendExpression",
	pos: position{line: 537, col: 1, offset: 15811},
	expr: &actionExpr{
	pos: position{line: 537, col: 26, offset: 15838},
	run: (*parser).callonListAppendExpression1,
	expr: &seqExpr{
	pos: position{line: 537, col: 26, offset: 15838},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 537, col: 26, offset: 15838},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 537, col: 32, offset: 15844},
	name: "AndExpression",
},
},
&labeledExpr{
	pos: position{line: 537, col: 55, offset: 15867},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 537, col: 60, offset: 15872},
	expr: &seqExpr{
	pos: position{line: 537, col: 61, offset: 15873},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 537, col: 61, offset: 15873},
	name: "_",
},
&litMatcher{
	pos: position{line: 537, col: 63, offset: 15875},
	val: "#",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 537, col: 67, offset: 15879},
	name: "_",
},
&labeledExpr{
	pos: position{line: 537, col: 69, offset: 15881},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 537, col: 71, offset: 15883},
	name: "AndExpression",
},
},
	},
},
},
},
	},
},
},
},
{
	name: "AndExpression",
	pos: position{line: 539, col: 1, offset: 15956},
	expr: &actionExpr{
	pos: position{line: 539, col: 26, offset: 15983},
	run: (*parser).callonAndExpression1,
	expr: &seqExpr{
	pos: position{line: 539, col: 26, offset: 15983},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 539, col: 26, offset: 15983},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 539, col: 32, offset: 15989},
	name: "CombineExpression",
},
},
&labeledExpr{
	pos: position{line: 539, col: 55, offset: 16012},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 539, col: 60, offset: 16017},
	expr: &seqExpr{
	pos: position{line: 539, col: 61, offset: 16018},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 539, col: 61, offset: 16018},
	name: "_",
},
&litMatcher{
	pos: position{line: 539, col: 63, offset: 16020},
	val: "&&",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 539, col: 68, offset: 16025},
	name: "_",
},
&labeledExpr{
	pos: position{line: 539, col: 70, offset: 16027},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 539, col: 72, offset: 16029},
	name: "CombineExpression",
},
},
	},
},
},
},
	},
},
},
},
{
	name: "CombineExpression",
	pos: position{line: 541, col: 1, offset: 16099},
	expr: &actionExpr{
	pos: position{line: 541, col: 26, offset: 16126},
	run: (*parser).callonCombineExpression1,
	expr: &seqExpr{
	pos: position{line: 541, col: 26, offset: 16126},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 541, col: 26, offset: 16126},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 541, col: 32, offset: 16132},
	name: "PreferExpression",
},
},
&labeledExpr{
	pos: position{line: 541, col: 55, offset: 16155},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 541, col: 60, offset: 16160},
	expr: &seqExpr{
	pos: position{line: 541, col: 61, offset: 16161},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 541, col: 61, offset: 16161},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 541, col: 63, offset: 16163},
	name: "Combine",
},
&ruleRefExpr{
	pos: position{line: 541, col: 71, offset: 16171},
	name: "_",
},
&labeledExpr{
	pos: position{line: 541, col: 73, offset: 16173},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 541, col: 75, offset: 16175},
	name: "PreferExpression",
},
},
	},
},
},
},
	},
},
},
},
{
	name: "PreferExpression",
	pos: position{line: 543, col: 1, offset: 16252},
	expr: &actionExpr{
	pos: position{line: 543, col: 26, offset: 16279},
	run: (*parser).callonPreferExpression1,
	expr: &seqExpr{
	pos: position{line: 543, col: 26, offset: 16279},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 543, col: 26, offset: 16279},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 543, col: 32, offset: 16285},
	name: "CombineTypesExpression",
},
},
&labeledExpr{
	pos: position{line: 543, col: 55, offset: 16308},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 543, col: 60, offset: 16313},
	expr: &seqExpr{
	pos: position{line: 543, col: 61, offset: 16314},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 543, col: 61, offset: 16314},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 543, col: 63, offset: 16316},
	name: "Prefer",
},
&ruleRefExpr{
	pos: position{line: 543, col: 70, offset: 16323},
	name: "_",
},
&labeledExpr{
	pos: position{line: 543, col: 72, offset: 16325},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 543, col: 74, offset: 16327},
	name: "CombineTypesExpression",
},
},
	},
},
},
},
	},
},
},
},
{
	name: "CombineTypesExpression",
	pos: position{line: 545, col: 1, offset: 16421},
	expr: &actionExpr{
	pos: position{line: 545, col: 26, offset: 16448},
	run: (*parser).callonCombineTypesExpression1,
	expr: &seqExpr{
	pos: position{line: 545, col: 26, offset: 16448},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 545, col: 26, offset: 16448},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 545, col: 32, offset: 16454},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 545, col: 55, offset: 16477},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 545, col: 60, offset: 16482},
	expr: &seqExpr{
	pos: position{line: 545, col: 61, offset: 16483},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 545, col: 61, offset: 16483},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 545, col: 63, offset: 16485},
	name: "CombineTypes",
},
&ruleRefExpr{
	pos: position{line: 545, col: 76, offset: 16498},
	name: "_",
},
&labeledExpr{
	pos: position{line: 545, col: 78, offset: 16500},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 545, col: 80, offset: 16502},
	name: "TimesExpression",
},
},
	},
},
},
},
	},
},
},
},
{
	name: "TimesExpression",
	pos: position{line: 547, col: 1, offset: 16582},
	expr: &actionExpr{
	pos: position{line: 547, col: 26, offset: 16609},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 547, col: 26, offset: 16609},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 547, col: 26, offset: 16609},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 547, col: 32, offset: 16615},
	name: "EqualExpression",
},
},
&labeledExpr{
	pos: position{line: 547, col: 55, offset: 16638},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 547, col: 60, offset: 16643},
	expr: &seqExpr{
	pos: position{line: 547, col: 61, offset: 16644},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 547, col: 61, offset: 16644},
	name: "_",
},
&litMatcher{
	pos: position{line: 547, col: 63, offset: 16646},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 547, col: 67, offset: 16650},
	name: "_",
},
&labeledExpr{
	pos: position{line: 547, col: 69, offset: 16652},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 547, col: 71, offset: 16654},
	name: "EqualExpression",
},
},
	},
},
},
},
	},
},
},
},
{
	name: "EqualExpression",
	pos: position{line: 549, col: 1, offset: 16724},
	expr: &actionExpr{
	pos: position{line: 549, col: 26, offset: 16751},
	run: (*parser).callonEqualExpression1,
	expr: &seqExpr{
	pos: position{line: 549, col: 26, offset: 16751},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 549, col: 26, offset: 16751},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 549, col: 32, offset: 16757},
	name: "NotEqualExpression",
},
},
&labeledExpr{
	pos: position{line: 549, col: 55, offset: 16780},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 549, col: 60, offset: 16785},
	expr: &seqExpr{
	pos: position{line: 549, col: 61, offset: 16786},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 549, col: 61, offset: 16786},
	name: "_",
},
&litMatcher{
	pos: position{line: 549, col: 63, offset: 16788},
	val: "==",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 549, col: 68, offset: 16793},
	name: "_",
},
&labeledExpr{
	pos: position{line: 549, col: 70, offset: 16795},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 549, col: 72, offset: 16797},
	name: "NotEqualExpression",
},
},
	},
},
},
},
	},
},
},
},
{
	name: "NotEqualExpression",
	pos: position{line: 551, col: 1, offset: 16867},
	expr: &actionExpr{
	pos: position{line: 551, col: 26, offset: 16894},
	run: (*parser).callonNotEqualExpression1,
	expr: &seqExpr{
	pos: position{line: 551, col: 26, offset: 16894},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 551, col: 26, offset: 16894},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 551, col: 32, offset: 16900},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 551, col: 55, offset: 16923},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 551, col: 60, offset: 16928},
	expr: &seqExpr{
	pos: position{line: 551, col: 61, offset: 16929},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 551, col: 61, offset: 16929},
	name: "_",
},
&litMatcher{
	pos: position{line: 551, col: 63, offset: 16931},
	val: "!=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 551, col: 68, offset: 16936},
	name: "_",
},
&labeledExpr{
	pos: position{line: 551, col: 70, offset: 16938},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 551, col: 72, offset: 16940},
	name: "ApplicationExpression",
},
},
	},
},
},
},
	},
},
},
},
{
	name: "ApplicationExpression",
	pos: position{line: 554, col: 1, offset: 17014},
	expr: &actionExpr{
	pos: position{line: 554, col: 25, offset: 17040},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 554, col: 25, offset: 17040},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 554, col: 25, offset: 17040},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 554, col: 27, offset: 17042},
	name: "FirstApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 554, col: 54, offset: 17069},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 554, col: 59, offset: 17074},
	expr: &seqExpr{
	pos: position{line: 554, col: 60, offset: 17075},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 554, col: 60, offset: 17075},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 554, col: 63, offset: 17078},
	name: "ImportExpression",
},
	},
},
},
},
	},
},
},
},
{
	name: "FirstApplicationExpression",
	pos: position{line: 563, col: 1, offset: 17321},
	expr: &choiceExpr{
	pos: position{line: 564, col: 8, offset: 17359},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 564, col: 8, offset: 17359},
	run: (*parser).callonFirstApplicationExpression2,
	expr: &seqExpr{
	pos: position{line: 564, col: 8, offset: 17359},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 564, col: 8, offset: 17359},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 564, col: 14, offset: 17365},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 564, col: 17, offset: 17368},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 564, col: 19, offset: 17370},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 564, col: 36, offset: 17387},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 564, col: 39, offset: 17390},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 564, col: 41, offset: 17392},
	name: "ImportExpression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 567, col: 8, offset: 17495},
	run: (*parser).callonFirstApplicationExpression11,
	expr: &seqExpr{
	pos: position{line: 567, col: 8, offset: 17495},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 567, col: 8, offset: 17495},
	name: "Some",
},
&ruleRefExpr{
	pos: position{line: 567, col: 13, offset: 17500},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 567, col: 16, offset: 17503},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 567, col: 18, offset: 17505},
	name: "ImportExpression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 568, col: 8, offset: 17560},
	run: (*parser).callonFirstApplicationExpression17,
	expr: &seqExpr{
	pos: position{line: 568, col: 8, offset: 17560},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 568, col: 8, offset: 17560},
	name: "toMap",
},
&ruleRefExpr{
	pos: position{line: 568, col: 14, offset: 17566},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 568, col: 17, offset: 17569},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 568, col: 19, offset: 17571},
	name: "ImportExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 569, col: 8, offset: 17635},
	name: "ImportExpression",
},
	},
},
},
{
	name: "ImportExpression",
	pos: position{line: 571, col: 1, offset: 17653},
	expr: &choiceExpr{
	pos: position{line: 571, col: 20, offset: 17674},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 571, col: 20, offset: 17674},
	name: "Import",
},
&ruleRefExpr{
	pos: position{line: 571, col: 29, offset: 17683},
	name: "SelectorExpression",
},
	},
},
},
{
	name: "SelectorExpression",
	pos: position{line: 573, col: 1, offset: 17703},
	expr: &actionExpr{
	pos: position{line: 573, col: 22, offset: 17726},
	run: (*parser).callonSelectorExpression1,
	expr: &seqExpr{
	pos: position{line: 573, col: 22, offset: 17726},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 573, col: 22, offset: 17726},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 573, col: 24, offset: 17728},
	name: "PrimitiveExpression",
},
},
&labeledExpr{
	pos: position{line: 573, col: 44, offset: 17748},
	label: "ls",
	expr: &zeroOrMoreExpr{
	pos: position{line: 573, col: 47, offset: 17751},
	expr: &seqExpr{
	pos: position{line: 573, col: 48, offset: 17752},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 573, col: 48, offset: 17752},
	name: "_",
},
&litMatcher{
	pos: position{line: 573, col: 50, offset: 17754},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 573, col: 54, offset: 17758},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 573, col: 56, offset: 17760},
	name: "Selector",
},
	},
},
},
},
	},
},
},
},
{
	name: "Selector",
	pos: position{line: 592, col: 1, offset: 18313},
	expr: &choiceExpr{
	pos: position{line: 592, col: 12, offset: 18326},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 592, col: 12, offset: 18326},
	name: "AnyLabel",
},
&ruleRefExpr{
	pos: position{line: 592, col: 23, offset: 18337},
	name: "Labels",
},
&ruleRefExpr{
	pos: position{line: 592, col: 32, offset: 18346},
	name: "TypeSelector",
},
	},
},
},
{
	name: "Labels",
	pos: position{line: 594, col: 1, offset: 18360},
	expr: &actionExpr{
	pos: position{line: 594, col: 10, offset: 18371},
	run: (*parser).callonLabels1,
	expr: &seqExpr{
	pos: position{line: 594, col: 10, offset: 18371},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 594, col: 10, offset: 18371},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 594, col: 14, offset: 18375},
	name: "_",
},
&labeledExpr{
	pos: position{line: 594, col: 16, offset: 18377},
	label: "optclauses",
	expr: &zeroOrOneExpr{
	pos: position{line: 594, col: 27, offset: 18388},
	expr: &seqExpr{
	pos: position{line: 594, col: 29, offset: 18390},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 594, col: 29, offset: 18390},
	name: "AnyLabel",
},
&ruleRefExpr{
	pos: position{line: 594, col: 38, offset: 18399},
	name: "_",
},
&zeroOrMoreExpr{
	pos: position{line: 594, col: 40, offset: 18401},
	expr: &seqExpr{
	pos: position{line: 594, col: 41, offset: 18402},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 594, col: 41, offset: 18402},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 594, col: 45, offset: 18406},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 594, col: 47, offset: 18408},
	name: "AnyLabel",
},
&ruleRefExpr{
	pos: position{line: 594, col: 56, offset: 18417},
	name: "_",
},
	},
},
},
	},
},
},
},
&litMatcher{
	pos: position{line: 594, col: 64, offset: 18425},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TypeSelector",
	pos: position{line: 604, col: 1, offset: 18721},
	expr: &actionExpr{
	pos: position{line: 604, col: 16, offset: 18738},
	run: (*parser).callonTypeSelector1,
	expr: &seqExpr{
	pos: position{line: 604, col: 16, offset: 18738},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 604, col: 16, offset: 18738},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 604, col: 20, offset: 18742},
	name: "_",
},
&labeledExpr{
	pos: position{line: 604, col: 22, offset: 18744},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 604, col: 24, offset: 18746},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 604, col: 35, offset: 18757},
	name: "_",
},
&litMatcher{
	pos: position{line: 604, col: 37, offset: 18759},
	val: ")",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PrimitiveExpression",
	pos: position{line: 606, col: 1, offset: 18782},
	expr: &choiceExpr{
	pos: position{line: 607, col: 7, offset: 18812},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 607, col: 7, offset: 18812},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 608, col: 7, offset: 18832},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 609, col: 7, offset: 18853},
	name: "IntegerLiteral",
},
&ruleRefExpr{
	pos: position{line: 610, col: 7, offset: 18874},
	name: "TextLiteral",
},
&actionExpr{
	pos: position{line: 611, col: 7, offset: 18892},
	run: (*parser).callonPrimitiveExpression6,
	expr: &seqExpr{
	pos: position{line: 611, col: 7, offset: 18892},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 611, col: 7, offset: 18892},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 611, col: 11, offset: 18896},
	name: "_",
},
&labeledExpr{
	pos: position{line: 611, col: 13, offset: 18898},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 611, col: 15, offset: 18900},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 611, col: 35, offset: 18920},
	name: "_",
},
&litMatcher{
	pos: position{line: 611, col: 37, offset: 18922},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&actionExpr{
	pos: position{line: 612, col: 7, offset: 18950},
	run: (*parser).callonPrimitiveExpression14,
	expr: &seqExpr{
	pos: position{line: 612, col: 7, offset: 18950},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 612, col: 7, offset: 18950},
	val: "<",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 612, col: 11, offset: 18954},
	name: "_",
},
&labeledExpr{
	pos: position{line: 612, col: 13, offset: 18956},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 612, col: 15, offset: 18958},
	name: "UnionType",
},
},
&ruleRefExpr{
	pos: position{line: 612, col: 25, offset: 18968},
	name: "_",
},
&litMatcher{
	pos: position{line: 612, col: 27, offset: 18970},
	val: ">",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 613, col: 7, offset: 18998},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 614, col: 7, offset: 19024},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 615, col: 7, offset: 19041},
	run: (*parser).callonPrimitiveExpression24,
	expr: &seqExpr{
	pos: position{line: 615, col: 7, offset: 19041},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 615, col: 7, offset: 19041},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 615, col: 11, offset: 19045},
	name: "_",
},
&labeledExpr{
	pos: position{line: 615, col: 14, offset: 19048},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 615, col: 16, offset: 19050},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 615, col: 27, offset: 19061},
	name: "_",
},
&litMatcher{
	pos: position{line: 615, col: 29, offset: 19063},
	val: ")",
	ignoreCase: false,
},
	},
},
},
	},
},
},
{
	name: "RecordTypeOrLiteral",
	pos: position{line: 617, col: 1, offset: 19086},
	expr: &choiceExpr{
	pos: position{line: 618, col: 7, offset: 19116},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 618, col: 7, offset: 19116},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &litMatcher{
	pos: position{line: 618, col: 7, offset: 19116},
	val: "=",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 619, col: 7, offset: 19154},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 620, col: 7, offset: 19179},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 621, col: 7, offset: 19207},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 621, col: 7, offset: 19207},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 623, col: 1, offset: 19236},
	expr: &actionExpr{
	pos: position{line: 623, col: 19, offset: 19256},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 623, col: 19, offset: 19256},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 623, col: 19, offset: 19256},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 623, col: 24, offset: 19261},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 623, col: 33, offset: 19270},
	name: "_",
},
&litMatcher{
	pos: position{line: 623, col: 35, offset: 19272},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 623, col: 39, offset: 19276},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 623, col: 42, offset: 19279},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 623, col: 47, offset: 19284},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 626, col: 1, offset: 19341},
	expr: &actionExpr{
	pos: position{line: 626, col: 18, offset: 19360},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 626, col: 18, offset: 19360},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 626, col: 18, offset: 19360},
	name: "_",
},
&litMatcher{
	pos: position{line: 626, col: 20, offset: 19362},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 626, col: 24, offset: 19366},
	name: "_",
},
&labeledExpr{
	pos: position{line: 626, col: 26, offset: 19368},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 626, col: 28, offset: 19370},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 627, col: 1, offset: 19402},
	expr: &actionExpr{
	pos: position{line: 628, col: 7, offset: 19431},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 628, col: 7, offset: 19431},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 628, col: 7, offset: 19431},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 628, col: 13, offset: 19437},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 628, col: 29, offset: 19453},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 628, col: 34, offset: 19458},
	expr: &ruleRefExpr{
	pos: position{line: 628, col: 34, offset: 19458},
	name: "MoreRecordType",
},
},
},
	},
},
},
},
{
	name: "RecordLiteralField",
	pos: position{line: 642, col: 1, offset: 20025},
	expr: &actionExpr{
	pos: position{line: 642, col: 22, offset: 20048},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 642, col: 22, offset: 20048},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 642, col: 22, offset: 20048},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 642, col: 27, offset: 20053},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 642, col: 36, offset: 20062},
	name: "_",
},
&litMatcher{
	pos: position{line: 642, col: 38, offset: 20064},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 642, col: 42, offset: 20068},
	name: "_",
},
&labeledExpr{
	pos: position{line: 642, col: 44, offset: 20070},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 642, col: 49, offset: 20075},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 645, col: 1, offset: 20132},
	expr: &actionExpr{
	pos: position{line: 645, col: 21, offset: 20154},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 645, col: 21, offset: 20154},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 645, col: 21, offset: 20154},
	name: "_",
},
&litMatcher{
	pos: position{line: 645, col: 23, offset: 20156},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 645, col: 27, offset: 20160},
	name: "_",
},
&labeledExpr{
	pos: position{line: 645, col: 29, offset: 20162},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 645, col: 31, offset: 20164},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 646, col: 1, offset: 20199},
	expr: &actionExpr{
	pos: position{line: 647, col: 7, offset: 20231},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 647, col: 7, offset: 20231},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 647, col: 7, offset: 20231},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 647, col: 13, offset: 20237},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 647, col: 32, offset: 20256},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 647, col: 37, offset: 20261},
	expr: &ruleRefExpr{
	pos: position{line: 647, col: 37, offset: 20261},
	name: "MoreRecordLiteral",
},
},
},
	},
},
},
},
{
	name: "UnionType",
	pos: position{line: 661, col: 1, offset: 20834},
	expr: &choiceExpr{
	pos: position{line: 661, col: 13, offset: 20848},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 661, col: 13, offset: 20848},
	name: "NonEmptyUnionType",
},
&ruleRefExpr{
	pos: position{line: 661, col: 33, offset: 20868},
	name: "EmptyUnionType",
},
	},
},
},
{
	name: "EmptyUnionType",
	pos: position{line: 663, col: 1, offset: 20884},
	expr: &actionExpr{
	pos: position{line: 663, col: 18, offset: 20903},
	run: (*parser).callonEmptyUnionType1,
	expr: &litMatcher{
	pos: position{line: 663, col: 18, offset: 20903},
	val: "",
	ignoreCase: false,
},
},
},
{
	name: "NonEmptyUnionType",
	pos: position{line: 665, col: 1, offset: 20935},
	expr: &actionExpr{
	pos: position{line: 665, col: 21, offset: 20957},
	run: (*parser).callonNonEmptyUnionType1,
	expr: &seqExpr{
	pos: position{line: 665, col: 21, offset: 20957},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 665, col: 21, offset: 20957},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 665, col: 27, offset: 20963},
	name: "UnionVariant",
},
},
&labeledExpr{
	pos: position{line: 665, col: 40, offset: 20976},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 665, col: 45, offset: 20981},
	expr: &seqExpr{
	pos: position{line: 665, col: 46, offset: 20982},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 665, col: 46, offset: 20982},
	name: "_",
},
&litMatcher{
	pos: position{line: 665, col: 48, offset: 20984},
	val: "|",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 665, col: 52, offset: 20988},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 665, col: 54, offset: 20990},
	name: "UnionVariant",
},
	},
},
},
},
	},
},
},
},
{
	name: "UnionVariant",
	pos: position{line: 685, col: 1, offset: 21695},
	expr: &seqExpr{
	pos: position{line: 685, col: 16, offset: 21712},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 685, col: 16, offset: 21712},
	name: "AnyLabel",
},
&zeroOrOneExpr{
	pos: position{line: 685, col: 25, offset: 21721},
	expr: &seqExpr{
	pos: position{line: 685, col: 26, offset: 21722},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 685, col: 26, offset: 21722},
	name: "_",
},
&litMatcher{
	pos: position{line: 685, col: 28, offset: 21724},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 685, col: 32, offset: 21728},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 685, col: 35, offset: 21731},
	name: "Expression",
},
	},
},
},
	},
},
},
{
	name: "MoreList",
	pos: position{line: 687, col: 1, offset: 21745},
	expr: &actionExpr{
	pos: position{line: 687, col: 12, offset: 21758},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 687, col: 12, offset: 21758},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 687, col: 12, offset: 21758},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 687, col: 16, offset: 21762},
	name: "_",
},
&labeledExpr{
	pos: position{line: 687, col: 18, offset: 21764},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 687, col: 20, offset: 21766},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 687, col: 31, offset: 21777},
	name: "_",
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 689, col: 1, offset: 21796},
	expr: &actionExpr{
	pos: position{line: 690, col: 7, offset: 21826},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 690, col: 7, offset: 21826},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 690, col: 7, offset: 21826},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 690, col: 11, offset: 21830},
	name: "_",
},
&labeledExpr{
	pos: position{line: 690, col: 13, offset: 21832},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 690, col: 19, offset: 21838},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 690, col: 30, offset: 21849},
	name: "_",
},
&labeledExpr{
	pos: position{line: 690, col: 32, offset: 21851},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 690, col: 37, offset: 21856},
	expr: &ruleRefExpr{
	pos: position{line: 690, col: 37, offset: 21856},
	name: "MoreList",
},
},
},
&litMatcher{
	pos: position{line: 690, col: 47, offset: 21866},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 700, col: 1, offset: 22134},
	expr: &notExpr{
	pos: position{line: 700, col: 7, offset: 22142},
	expr: &anyMatcher{
	line: 700, col: 8, offset: 22143,
},
},
},
	},
}
func (c *current) onDhallFile1(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonDhallFile1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDhallFile1(stack["e"])
}

func (c *current) onCompleteExpression1(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonCompleteExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCompleteExpression1(stack["e"])
}

func (c *current) onEOL3() (interface{}, error) {
 return []byte{'\n'}, nil 
}

func (p *parser) callonEOL3() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEOL3()
}

func (c *current) onLineComment5() (interface{}, error) {
 return string(c.text), nil
}

func (p *parser) callonLineComment5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLineComment5()
}

func (c *current) onLineComment1(content interface{}) (interface{}, error) {
 return content, nil 
}

func (p *parser) callonLineComment1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLineComment1(stack["content"])
}

func (c *current) onSimpleLabel2() (interface{}, error) {
 return string(c.text), nil 
}

func (p *parser) callonSimpleLabel2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSimpleLabel2()
}

func (c *current) onSimpleLabel7() (interface{}, error) {
            return string(c.text), nil
          
}

func (p *parser) callonSimpleLabel7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSimpleLabel7()
}

func (c *current) onQuotedLabel1() (interface{}, error) {
 return string(c.text), nil 
}

func (p *parser) callonQuotedLabel1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onQuotedLabel1()
}

func (c *current) onLabel2(label interface{}) (interface{}, error) {
 return label, nil 
}

func (p *parser) callonLabel2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLabel2(stack["label"])
}

func (c *current) onLabel8(label interface{}) (interface{}, error) {
 return label, nil 
}

func (p *parser) callonLabel8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLabel8(stack["label"])
}

func (c *current) onNonreservedLabel2(label interface{}) (interface{}, error) {
 return label, nil 
}

func (p *parser) callonNonreservedLabel2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNonreservedLabel2(stack["label"])
}

func (c *current) onNonreservedLabel10(label interface{}) (interface{}, error) {
 return label, nil 
}

func (p *parser) callonNonreservedLabel10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNonreservedLabel10(stack["label"])
}

func (c *current) onDoubleQuoteChunk3(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonDoubleQuoteChunk3() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleQuoteChunk3(stack["e"])
}

func (c *current) onDoubleQuoteEscaped6() (interface{}, error) {
 return []byte("\b"), nil 
}

func (p *parser) callonDoubleQuoteEscaped6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleQuoteEscaped6()
}

func (c *current) onDoubleQuoteEscaped8() (interface{}, error) {
 return []byte("\f"), nil 
}

func (p *parser) callonDoubleQuoteEscaped8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleQuoteEscaped8()
}

func (c *current) onDoubleQuoteEscaped10() (interface{}, error) {
 return []byte("\n"), nil 
}

func (p *parser) callonDoubleQuoteEscaped10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleQuoteEscaped10()
}

func (c *current) onDoubleQuoteEscaped12() (interface{}, error) {
 return []byte("\r"), nil 
}

func (p *parser) callonDoubleQuoteEscaped12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleQuoteEscaped12()
}

func (c *current) onDoubleQuoteEscaped14() (interface{}, error) {
 return []byte("\t"), nil 
}

func (p *parser) callonDoubleQuoteEscaped14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleQuoteEscaped14()
}

func (c *current) onDoubleQuoteEscaped16(u interface{}) (interface{}, error) {
 return u, nil 
}

func (p *parser) callonDoubleQuoteEscaped16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleQuoteEscaped16(stack["u"])
}

func (c *current) onUnicodeEscape2() (interface{}, error) {
            return ParseCodepoint(string(c.text))
        
}

func (p *parser) callonUnicodeEscape2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnicodeEscape2()
}

func (c *current) onUnicodeEscape8() (interface{}, error) {
            return ParseCodepoint(string(c.text[1:len(c.text)-1]))
        
}

func (p *parser) callonUnicodeEscape8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnicodeEscape8()
}

func (c *current) onDoubleQuoteLiteral1(chunks interface{}) (interface{}, error) {
    var str strings.Builder
    var outChunks Chunks
    for _, chunk := range chunks.([]interface{}) {
        switch e := chunk.(type) {
        case []byte:
                str.Write(e)
        case Expr:
                outChunks = append(outChunks, Chunk{str.String(), e})
                str.Reset()
        default:
                return nil, errors.New("can't happen")
        }
    }
    return TextLit{Chunks: outChunks, Suffix: str.String()}, nil
}

func (p *parser) callonDoubleQuoteLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleQuoteLiteral1(stack["chunks"])
}

func (c *current) onEscapedQuotePair1() (interface{}, error) {
 return []byte("''"), nil 
}

func (p *parser) callonEscapedQuotePair1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEscapedQuotePair1()
}

func (c *current) onEscapedInterpolation1() (interface{}, error) {
 return []byte("$\u007b"), nil 
}

func (p *parser) callonEscapedInterpolation1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEscapedInterpolation1()
}

func (c *current) onSingleQuoteLiteral1(content interface{}) (interface{}, error) {
    var str strings.Builder
    var outChunks Chunks
    chunk, ok := content.([]interface{})
    for ; ok; chunk, ok = chunk[1].([]interface{}) {
        switch e := chunk[0].(type) {
        case []byte:
            str.Write(e)
        case Expr:
                outChunks = append(outChunks, Chunk{str.String(), e})
                str.Reset()
        default:
            return nil, errors.New("unimplemented")
        }
    }
    return RemoveLeadingCommonIndent(TextLit{Chunks: outChunks, Suffix: str.String()}), nil
}

func (p *parser) callonSingleQuoteLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSingleQuoteLiteral1(stack["content"])
}

func (c *current) onInterpolation1(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonInterpolation1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInterpolation1(stack["e"])
}

func (c *current) onReserved2() (interface{}, error) {
 return NaturalBuild, nil 
}

func (p *parser) callonReserved2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved2()
}

func (c *current) onReserved4() (interface{}, error) {
 return NaturalFold, nil 
}

func (p *parser) callonReserved4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved4()
}

func (c *current) onReserved6() (interface{}, error) {
 return NaturalIsZero, nil 
}

func (p *parser) callonReserved6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved6()
}

func (c *current) onReserved8() (interface{}, error) {
 return NaturalEven, nil 
}

func (p *parser) callonReserved8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved8()
}

func (c *current) onReserved10() (interface{}, error) {
 return NaturalOdd, nil 
}

func (p *parser) callonReserved10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved10()
}

func (c *current) onReserved12() (interface{}, error) {
 return NaturalToInteger, nil 
}

func (p *parser) callonReserved12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved12()
}

func (c *current) onReserved14() (interface{}, error) {
 return NaturalShow, nil 
}

func (p *parser) callonReserved14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved14()
}

func (c *current) onReserved16() (interface{}, error) {
 return IntegerToDouble, nil 
}

func (p *parser) callonReserved16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved16()
}

func (c *current) onReserved18() (interface{}, error) {
 return IntegerShow, nil 
}

func (p *parser) callonReserved18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved18()
}

func (c *current) onReserved20() (interface{}, error) {
 return DoubleShow, nil 
}

func (p *parser) callonReserved20() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved20()
}

func (c *current) onReserved22() (interface{}, error) {
 return ListBuild, nil 
}

func (p *parser) callonReserved22() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved22()
}

func (c *current) onReserved24() (interface{}, error) {
 return ListFold, nil 
}

func (p *parser) callonReserved24() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved24()
}

func (c *current) onReserved26() (interface{}, error) {
 return ListLength, nil 
}

func (p *parser) callonReserved26() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved26()
}

func (c *current) onReserved28() (interface{}, error) {
 return ListHead, nil 
}

func (p *parser) callonReserved28() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved28()
}

func (c *current) onReserved30() (interface{}, error) {
 return ListLast, nil 
}

func (p *parser) callonReserved30() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved30()
}

func (c *current) onReserved32() (interface{}, error) {
 return ListIndexed, nil 
}

func (p *parser) callonReserved32() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved32()
}

func (c *current) onReserved34() (interface{}, error) {
 return ListReverse, nil 
}

func (p *parser) callonReserved34() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved34()
}

func (c *current) onReserved36() (interface{}, error) {
 return OptionalBuild, nil 
}

func (p *parser) callonReserved36() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved36()
}

func (c *current) onReserved38() (interface{}, error) {
 return OptionalFold, nil 
}

func (p *parser) callonReserved38() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved38()
}

func (c *current) onReserved40() (interface{}, error) {
 return TextShow, nil 
}

func (p *parser) callonReserved40() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved40()
}

func (c *current) onReserved42() (interface{}, error) {
 return Bool, nil 
}

func (p *parser) callonReserved42() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved42()
}

func (c *current) onReserved44() (interface{}, error) {
 return True, nil 
}

func (p *parser) callonReserved44() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved44()
}

func (c *current) onReserved46() (interface{}, error) {
 return False, nil 
}

func (p *parser) callonReserved46() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved46()
}

func (c *current) onReserved48() (interface{}, error) {
 return Optional, nil 
}

func (p *parser) callonReserved48() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved48()
}

func (c *current) onReserved50() (interface{}, error) {
 return Natural, nil 
}

func (p *parser) callonReserved50() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved50()
}

func (c *current) onReserved52() (interface{}, error) {
 return Integer, nil 
}

func (p *parser) callonReserved52() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved52()
}

func (c *current) onReserved54() (interface{}, error) {
 return Double, nil 
}

func (p *parser) callonReserved54() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved54()
}

func (c *current) onReserved56() (interface{}, error) {
 return Text, nil 
}

func (p *parser) callonReserved56() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved56()
}

func (c *current) onReserved58() (interface{}, error) {
 return List, nil 
}

func (p *parser) callonReserved58() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved58()
}

func (c *current) onReserved60() (interface{}, error) {
 return None, nil 
}

func (p *parser) callonReserved60() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved60()
}

func (c *current) onReserved62() (interface{}, error) {
 return Type, nil 
}

func (p *parser) callonReserved62() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved62()
}

func (c *current) onReserved64() (interface{}, error) {
 return Kind, nil 
}

func (p *parser) callonReserved64() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved64()
}

func (c *current) onReserved66() (interface{}, error) {
 return Sort, nil 
}

func (p *parser) callonReserved66() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved66()
}

func (c *current) onMissing1() (interface{}, error) {
 return Missing{}, nil 
}

func (p *parser) callonMissing1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMissing1()
}

func (c *current) onNumericDoubleLiteral1() (interface{}, error) {
      d, err := strconv.ParseFloat(string(c.text), 64)
      if err != nil {
         return nil, err
      }
      return DoubleLit(d), nil
}

func (p *parser) callonNumericDoubleLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNumericDoubleLiteral1()
}

func (c *current) onDoubleLiteral4() (interface{}, error) {
 return DoubleLit(math.Inf(1)), nil 
}

func (p *parser) callonDoubleLiteral4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleLiteral4()
}

func (c *current) onDoubleLiteral6() (interface{}, error) {
 return DoubleLit(math.Inf(-1)), nil 
}

func (p *parser) callonDoubleLiteral6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleLiteral6()
}

func (c *current) onDoubleLiteral10() (interface{}, error) {
 return DoubleLit(math.NaN()), nil 
}

func (p *parser) callonDoubleLiteral10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleLiteral10()
}

func (c *current) onNaturalLiteral1() (interface{}, error) {
      i, err := strconv.Atoi(string(c.text))
      return NaturalLit(i), err
}

func (p *parser) callonNaturalLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNaturalLiteral1()
}

func (c *current) onIntegerLiteral1() (interface{}, error) {
      i, err := strconv.Atoi(string(c.text))
      if err != nil {
         return nil, err
      }
      return IntegerLit(i), nil
}

func (p *parser) callonIntegerLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIntegerLiteral1()
}

func (c *current) onDeBruijn1(index interface{}) (interface{}, error) {
 return int(index.(NaturalLit)), nil 
}

func (p *parser) callonDeBruijn1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDeBruijn1(stack["index"])
}

func (c *current) onVariable1(name, index interface{}) (interface{}, error) {
    if index != nil {
        return Var{Name:name.(string), Index:index.(int)}, nil
    } else {
        return Var{Name:name.(string)}, nil
    }
}

func (p *parser) callonVariable1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVariable1(stack["name"], stack["index"])
}

func (c *current) onUnquotedPathComponent1() (interface{}, error) {
 return string(c.text), nil 
}

func (p *parser) callonUnquotedPathComponent1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnquotedPathComponent1()
}

func (c *current) onQuotedPathComponent1() (interface{}, error) {
 return string(c.text), nil 
}

func (p *parser) callonQuotedPathComponent1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onQuotedPathComponent1()
}

func (c *current) onPathComponent2(u interface{}) (interface{}, error) {
 return u, nil 
}

func (p *parser) callonPathComponent2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPathComponent2(stack["u"])
}

func (c *current) onPathComponent7(q interface{}) (interface{}, error) {
 return q, nil 
}

func (p *parser) callonPathComponent7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPathComponent7(stack["q"])
}

func (c *current) onPath1(cs interface{}) (interface{}, error) {
    // urgh, have to convert []interface{} to []string
    components := make([]string, len(cs.([]interface{})))
    for i, component := range cs.([]interface{}) {
        components[i] = component.(string)
    }
    return path.Join(components...), nil
}

func (p *parser) callonPath1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPath1(stack["cs"])
}

func (c *current) onParentPath1(p interface{}) (interface{}, error) {
 return Local(path.Join("..", p.(string))), nil 
}

func (p *parser) callonParentPath1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onParentPath1(stack["p"])
}

func (c *current) onHerePath1(p interface{}) (interface{}, error) {
 return Local(p.(string)), nil 
}

func (p *parser) callonHerePath1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHerePath1(stack["p"])
}

func (c *current) onHomePath1(p interface{}) (interface{}, error) {
 return Local(path.Join("~", p.(string))), nil 
}

func (p *parser) callonHomePath1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHomePath1(stack["p"])
}

func (c *current) onAbsolutePath1(p interface{}) (interface{}, error) {
 return Local(path.Join("/", p.(string))), nil 
}

func (p *parser) callonAbsolutePath1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAbsolutePath1(stack["p"])
}

func (c *current) onHttpRaw1() (interface{}, error) {
 return url.ParseRequestURI(string(c.text)) 
}

func (p *parser) callonHttpRaw1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHttpRaw1()
}

func (c *current) onIPv6address1() (interface{}, error) {
    addr := net.ParseIP(string(c.text))
    if addr == nil { return nil, errors.New("Malformed IPv6 address") }
    return string(c.text), nil
}

func (p *parser) callonIPv6address1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIPv6address1()
}

func (c *current) onHttp1(u interface{}) (interface{}, error) {
 return MakeRemote(u.(*url.URL)), nil 
}

func (p *parser) callonHttp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHttp1(stack["u"])
}

func (c *current) onEnv1(v interface{}) (interface{}, error) {
 return v, nil 
}

func (p *parser) callonEnv1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEnv1(stack["v"])
}

func (c *current) onBashEnvironmentVariable1() (interface{}, error) {
  return EnvVar(string(c.text)), nil
}

func (p *parser) callonBashEnvironmentVariable1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBashEnvironmentVariable1()
}

func (c *current) onPosixEnvironmentVariable1(v interface{}) (interface{}, error) {
  return v, nil
}

func (p *parser) callonPosixEnvironmentVariable1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariable1(stack["v"])
}

func (c *current) onPosixEnvironmentVariableContent1(v interface{}) (interface{}, error) {
  var b strings.Builder
  for _, c := range v.([]interface{}) {
    _, err := b.Write(c.([]byte))
    if err != nil { return nil, err }
  }
  return EnvVar(b.String()), nil
}

func (p *parser) callonPosixEnvironmentVariableContent1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariableContent1(stack["v"])
}

func (c *current) onPosixEnvironmentVariableCharacter2() (interface{}, error) {
 return []byte{0x22}, nil 
}

func (p *parser) callonPosixEnvironmentVariableCharacter2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariableCharacter2()
}

func (c *current) onPosixEnvironmentVariableCharacter4() (interface{}, error) {
 return []byte{0x5c}, nil 
}

func (p *parser) callonPosixEnvironmentVariableCharacter4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariableCharacter4()
}

func (c *current) onPosixEnvironmentVariableCharacter6() (interface{}, error) {
 return []byte{0x07}, nil 
}

func (p *parser) callonPosixEnvironmentVariableCharacter6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariableCharacter6()
}

func (c *current) onPosixEnvironmentVariableCharacter8() (interface{}, error) {
 return []byte{0x08}, nil 
}

func (p *parser) callonPosixEnvironmentVariableCharacter8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariableCharacter8()
}

func (c *current) onPosixEnvironmentVariableCharacter10() (interface{}, error) {
 return []byte{0x0c}, nil 
}

func (p *parser) callonPosixEnvironmentVariableCharacter10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariableCharacter10()
}

func (c *current) onPosixEnvironmentVariableCharacter12() (interface{}, error) {
 return []byte{0x0a}, nil 
}

func (p *parser) callonPosixEnvironmentVariableCharacter12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariableCharacter12()
}

func (c *current) onPosixEnvironmentVariableCharacter14() (interface{}, error) {
 return []byte{0x0d}, nil 
}

func (p *parser) callonPosixEnvironmentVariableCharacter14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariableCharacter14()
}

func (c *current) onPosixEnvironmentVariableCharacter16() (interface{}, error) {
 return []byte{0x09}, nil 
}

func (p *parser) callonPosixEnvironmentVariableCharacter16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariableCharacter16()
}

func (c *current) onPosixEnvironmentVariableCharacter18() (interface{}, error) {
 return []byte{0x0b}, nil 
}

func (p *parser) callonPosixEnvironmentVariableCharacter18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPosixEnvironmentVariableCharacter18()
}

func (c *current) onHashValue1() (interface{}, error) {
    out := make([]byte, sha256.Size)
    _, err := hex.Decode(out, c.text)
    if err != nil { return nil, err }
    return out, nil
}

func (p *parser) callonHashValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHashValue1()
}

func (c *current) onHash1(val interface{}) (interface{}, error) {
 return append([]byte{0x12,0x20}, val.([]byte)...), nil 
}

func (p *parser) callonHash1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHash1(stack["val"])
}

func (c *current) onImportHashed1(i, h interface{}) (interface{}, error) {
    out := ImportHashed{Fetchable: i.(Fetchable)}
    if h != nil {
        out.Hash = h.([]interface{})[1].([]byte)
    }
    return out, nil
}

func (p *parser) callonImportHashed1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onImportHashed1(stack["i"], stack["h"])
}

func (c *current) onImport2(i interface{}) (interface{}, error) {
 return Import{ImportHashed: i.(ImportHashed), ImportMode: RawText}, nil 
}

func (p *parser) callonImport2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onImport2(stack["i"])
}

func (c *current) onImport10(i interface{}) (interface{}, error) {
 return Import{ImportHashed: i.(ImportHashed), ImportMode: Location}, nil 
}

func (p *parser) callonImport10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onImport10(stack["i"])
}

func (c *current) onImport18(i interface{}) (interface{}, error) {
 return Import{ImportHashed: i.(ImportHashed), ImportMode: Code}, nil 
}

func (p *parser) callonImport18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onImport18(stack["i"])
}

func (c *current) onLetBinding1(label, a, v interface{}) (interface{}, error) {
    if a != nil {
        return Binding{
            Variable: label.(string),
            Annotation: a.([]interface{})[0].(Expr),
            Value: v.(Expr),
        }, nil
    } else {
        return Binding{
            Variable: label.(string),
            Value: v.(Expr),
        }, nil
    }
}

func (p *parser) callonLetBinding1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLetBinding1(stack["label"], stack["a"], stack["v"])
}

func (c *current) onExpression2(label, t, body interface{}) (interface{}, error) {
          return &LambdaExpr{Label:label.(string), Type:t.(Expr), Body: body.(Expr)}, nil
      
}

func (p *parser) callonExpression2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression2(stack["label"], stack["t"], stack["body"])
}

func (c *current) onExpression22(cond, t, f interface{}) (interface{}, error) {
          return BoolIf{cond.(Expr),t.(Expr),f.(Expr)},nil
      
}

func (p *parser) callonExpression22() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression22(stack["cond"], stack["t"], stack["f"])
}

func (c *current) onExpression38(bindings, b interface{}) (interface{}, error) {
        bs := make([]Binding, len(bindings.([]interface{})))
        for i, binding := range bindings.([]interface{}) {
            bs[i] = binding.(Binding)
        }
        return MakeLet(b.(Expr), bs...), nil
      
}

func (p *parser) callonExpression38() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression38(stack["bindings"], stack["b"])
}

func (c *current) onExpression47(label, t, body interface{}) (interface{}, error) {
          return &Pi{Label:label.(string), Type:t.(Expr), Body: body.(Expr)}, nil
      
}

func (p *parser) callonExpression47() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression47(stack["label"], stack["t"], stack["body"])
}

func (c *current) onExpression67(o, e interface{}) (interface{}, error) {
 return FnType(o.(Expr),e.(Expr)), nil 
}

func (p *parser) callonExpression67() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression67(stack["o"], stack["e"])
}

func (c *current) onExpression76(h, u, a interface{}) (interface{}, error) {
          return Merge{Handler:h.(Expr), Union:u.(Expr), Annotation:a.(Expr)}, nil
      
}

func (p *parser) callonExpression76() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression76(stack["h"], stack["u"], stack["a"])
}

func (c *current) onExpression91(e, t interface{}) (interface{}, error) {
 return ToMap{e.(Expr), t.(Expr)}, nil 
}

func (p *parser) callonExpression91() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression91(stack["e"], stack["t"])
}

func (c *current) onAnnotation1(a interface{}) (interface{}, error) {
 return a, nil 
}

func (p *parser) callonAnnotation1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAnnotation1(stack["a"])
}

func (c *current) onAnnotatedExpression1(e, a interface{}) (interface{}, error) {
        if a == nil { return e, nil }
        return Annot{e.(Expr), a.([]interface{})[1].(Expr)}, nil
    
}

func (p *parser) callonAnnotatedExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAnnotatedExpression1(stack["e"], stack["a"])
}

func (c *current) onEmptyList1(a interface{}) (interface{}, error) {
          return EmptyList{a.(Expr)},nil
}

func (p *parser) callonEmptyList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEmptyList1(stack["a"])
}

func (c *current) onImportAltExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(ImportAltOp, first, rest), nil
}

func (p *parser) callonImportAltExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onImportAltExpression1(stack["first"], stack["rest"])
}

func (c *current) onOrExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(OrOp, first, rest), nil
}

func (p *parser) callonOrExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOrExpression1(stack["first"], stack["rest"])
}

func (c *current) onPlusExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(PlusOp, first, rest), nil
}

func (p *parser) callonPlusExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPlusExpression1(stack["first"], stack["rest"])
}

func (c *current) onTextAppendExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(TextAppendOp, first, rest), nil
}

func (p *parser) callonTextAppendExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTextAppendExpression1(stack["first"], stack["rest"])
}

func (c *current) onListAppendExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(ListAppendOp, first, rest), nil
}

func (p *parser) callonListAppendExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onListAppendExpression1(stack["first"], stack["rest"])
}

func (c *current) onAndExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(AndOp, first, rest), nil
}

func (p *parser) callonAndExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAndExpression1(stack["first"], stack["rest"])
}

func (c *current) onCombineExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(RecordMergeOp, first, rest), nil
}

func (p *parser) callonCombineExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCombineExpression1(stack["first"], stack["rest"])
}

func (c *current) onPreferExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(RightBiasedRecordMergeOp, first, rest), nil
}

func (p *parser) callonPreferExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPreferExpression1(stack["first"], stack["rest"])
}

func (c *current) onCombineTypesExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(RecordTypeMergeOp, first, rest), nil
}

func (p *parser) callonCombineTypesExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCombineTypesExpression1(stack["first"], stack["rest"])
}

func (c *current) onTimesExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(TimesOp, first, rest), nil
}

func (p *parser) callonTimesExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTimesExpression1(stack["first"], stack["rest"])
}

func (c *current) onEqualExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(EqOp, first, rest), nil
}

func (p *parser) callonEqualExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEqualExpression1(stack["first"], stack["rest"])
}

func (c *current) onNotEqualExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(NeOp, first, rest), nil
}

func (p *parser) callonNotEqualExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNotEqualExpression1(stack["first"], stack["rest"])
}

func (c *current) onApplicationExpression1(f, rest interface{}) (interface{}, error) {
          e := f.(Expr)
          if rest == nil { return e, nil }
          for _, arg := range rest.([]interface{}) {
              e = Apply(e, arg.([]interface{})[1].(Expr))
          }
          return e,nil
      
}

func (p *parser) callonApplicationExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onApplicationExpression1(stack["f"], stack["rest"])
}

func (c *current) onFirstApplicationExpression2(h, u interface{}) (interface{}, error) {
             return Merge{Handler:h.(Expr), Union:u.(Expr)}, nil
          
}

func (p *parser) callonFirstApplicationExpression2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFirstApplicationExpression2(stack["h"], stack["u"])
}

func (c *current) onFirstApplicationExpression11(e interface{}) (interface{}, error) {
 return Some{e.(Expr)}, nil 
}

func (p *parser) callonFirstApplicationExpression11() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFirstApplicationExpression11(stack["e"])
}

func (c *current) onFirstApplicationExpression17(e interface{}) (interface{}, error) {
 return ToMap{Record: e.(Expr)}, nil 
}

func (p *parser) callonFirstApplicationExpression17() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFirstApplicationExpression17(stack["e"])
}

func (c *current) onSelectorExpression1(e, ls interface{}) (interface{}, error) {
    expr := e.(Expr)
    labels := ls.([]interface{})
    for _, labelSelector := range labels {
        selectorIface := labelSelector.([]interface{})[3]
        switch selector := selectorIface.(type) {
            case string:
                expr = Field{expr, selector}
            case []string:
                expr = Project{expr, selector}
            case Expr:
                expr = ProjectType{expr, selector}
            default:
                return nil, errors.New("unimplemented")
        }
    }
    return expr, nil
}

func (p *parser) callonSelectorExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSelectorExpression1(stack["e"], stack["ls"])
}

func (c *current) onLabels1(optclauses interface{}) (interface{}, error) {
    if optclauses == nil { return []string{}, nil }
    clauses := optclauses.([]interface{})
    labels := []string{clauses[0].(string)}
    for _, next := range clauses[2].([]interface{}) {
        labels = append(labels, next.([]interface{})[2].(string))
    }
    return labels, nil
}

func (p *parser) callonLabels1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLabels1(stack["optclauses"])
}

func (c *current) onTypeSelector1(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonTypeSelector1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTypeSelector1(stack["e"])
}

func (c *current) onPrimitiveExpression6(r interface{}) (interface{}, error) {
 return r, nil 
}

func (p *parser) callonPrimitiveExpression6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimitiveExpression6(stack["r"])
}

func (c *current) onPrimitiveExpression14(u interface{}) (interface{}, error) {
 return u, nil 
}

func (p *parser) callonPrimitiveExpression14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimitiveExpression14(stack["u"])
}

func (c *current) onPrimitiveExpression24(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonPrimitiveExpression24() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimitiveExpression24(stack["e"])
}

func (c *current) onRecordTypeOrLiteral2() (interface{}, error) {
 return RecordLit{}, nil 
}

func (p *parser) callonRecordTypeOrLiteral2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRecordTypeOrLiteral2()
}

func (c *current) onRecordTypeOrLiteral6() (interface{}, error) {
 return Record{}, nil 
}

func (p *parser) callonRecordTypeOrLiteral6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRecordTypeOrLiteral6()
}

func (c *current) onRecordTypeField1(name, expr interface{}) (interface{}, error) {
    return []interface{}{name, expr}, nil
}

func (p *parser) callonRecordTypeField1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRecordTypeField1(stack["name"], stack["expr"])
}

func (c *current) onMoreRecordType1(f interface{}) (interface{}, error) {
return f, nil
}

func (p *parser) callonMoreRecordType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMoreRecordType1(stack["f"])
}

func (c *current) onNonEmptyRecordType1(first, rest interface{}) (interface{}, error) {
          fields := rest.([]interface{})
          content := make(Record, len(fields)+1)
          content[first.([]interface{})[0].(string)] = first.([]interface{})[1].(Expr)
          for _, field := range(fields) {
              fieldName := field.([]interface{})[0].(string)
              if _, ok := content[fieldName]; ok {
                  return nil, fmt.Errorf("Duplicate field %s in record", fieldName)
              }
              content[fieldName] = field.([]interface{})[1].(Expr)
          }
          return content, nil
      
}

func (p *parser) callonNonEmptyRecordType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNonEmptyRecordType1(stack["first"], stack["rest"])
}

func (c *current) onRecordLiteralField1(name, expr interface{}) (interface{}, error) {
    return []interface{}{name, expr}, nil
}

func (p *parser) callonRecordLiteralField1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRecordLiteralField1(stack["name"], stack["expr"])
}

func (c *current) onMoreRecordLiteral1(f interface{}) (interface{}, error) {
return f, nil
}

func (p *parser) callonMoreRecordLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMoreRecordLiteral1(stack["f"])
}

func (c *current) onNonEmptyRecordLiteral1(first, rest interface{}) (interface{}, error) {
          fields := rest.([]interface{})
          content := make(RecordLit, len(fields)+1)
          content[first.([]interface{})[0].(string)] = first.([]interface{})[1].(Expr)
          for _, field := range(fields) {
              fieldName := field.([]interface{})[0].(string)
              if _, ok := content[fieldName]; ok {
                  return nil, fmt.Errorf("Duplicate field %s in record", fieldName)
              }
              content[fieldName] = field.([]interface{})[1].(Expr)
          }
          return content, nil
      
}

func (p *parser) callonNonEmptyRecordLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNonEmptyRecordLiteral1(stack["first"], stack["rest"])
}

func (c *current) onEmptyUnionType1() (interface{}, error) {
 return UnionType{}, nil 
}

func (p *parser) callonEmptyUnionType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEmptyUnionType1()
}

func (c *current) onNonEmptyUnionType1(first, rest interface{}) (interface{}, error) {
    alternatives := make(UnionType)
    first2 := first.([]interface{})
    if first2[1] == nil {
        alternatives[first2[0].(string)] = nil
    } else {
        alternatives[first2[0].(string)] = first2[1].([]interface{})[3].(Expr)
    }
    if rest == nil { return UnionType(alternatives), nil }
    for _, alternativeSyntax := range rest.([]interface{}) {
        alternative := alternativeSyntax.([]interface{})[3].([]interface{})
        if alternative[1] == nil {
            alternatives[alternative[0].(string)] = nil
        } else {
            alternatives[alternative[0].(string)] = alternative[1].([]interface{})[3].(Expr)
        }
    }
    return alternatives, nil
}

func (p *parser) callonNonEmptyUnionType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNonEmptyUnionType1(stack["first"], stack["rest"])
}

func (c *current) onMoreList1(e interface{}) (interface{}, error) {
return e, nil
}

func (p *parser) callonMoreList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMoreList1(stack["e"])
}

func (c *current) onNonEmptyListLiteral1(first, rest interface{}) (interface{}, error) {
          exprs := rest.([]interface{})
          content := make(NonEmptyList, len(exprs)+1)
          content[0] = first.(Expr)
          for i, expr := range(exprs) {
              content[i+1] = expr.(Expr)
          }
          return content, nil
      
}

func (p *parser) callonNonEmptyListLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNonEmptyListLiteral1(stack["first"], stack["rest"])
}


var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule          = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch         = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos    position
	expr   interface{}
	run    func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs: new(errList),
		data: b,
		pt: savepoint{position: position{line: 1}},
		recover: true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v interface{}
	b bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug bool
	depth  int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules  map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth) + ">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth) + "<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n == 1 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth) + "MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth) + "MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	// can't match EOF
	if cur == utf8.RuneError {
		return nil, false
	}
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}

