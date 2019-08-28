
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
	val: "Natural/subtract",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 229, col: 5, offset: 6048},
	run: (*parser).callonReserved18,
	expr: &litMatcher{
	pos: position{line: 229, col: 5, offset: 6048},
	val: "Integer/toDouble",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 230, col: 5, offset: 6103},
	run: (*parser).callonReserved20,
	expr: &litMatcher{
	pos: position{line: 230, col: 5, offset: 6103},
	val: "Integer/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 231, col: 5, offset: 6150},
	run: (*parser).callonReserved22,
	expr: &litMatcher{
	pos: position{line: 231, col: 5, offset: 6150},
	val: "Double/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 232, col: 5, offset: 6195},
	run: (*parser).callonReserved24,
	expr: &litMatcher{
	pos: position{line: 232, col: 5, offset: 6195},
	val: "List/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 233, col: 5, offset: 6238},
	run: (*parser).callonReserved26,
	expr: &litMatcher{
	pos: position{line: 233, col: 5, offset: 6238},
	val: "List/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 234, col: 5, offset: 6279},
	run: (*parser).callonReserved28,
	expr: &litMatcher{
	pos: position{line: 234, col: 5, offset: 6279},
	val: "List/length",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 235, col: 5, offset: 6324},
	run: (*parser).callonReserved30,
	expr: &litMatcher{
	pos: position{line: 235, col: 5, offset: 6324},
	val: "List/head",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 236, col: 5, offset: 6365},
	run: (*parser).callonReserved32,
	expr: &litMatcher{
	pos: position{line: 236, col: 5, offset: 6365},
	val: "List/last",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 237, col: 5, offset: 6406},
	run: (*parser).callonReserved34,
	expr: &litMatcher{
	pos: position{line: 237, col: 5, offset: 6406},
	val: "List/indexed",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 238, col: 5, offset: 6453},
	run: (*parser).callonReserved36,
	expr: &litMatcher{
	pos: position{line: 238, col: 5, offset: 6453},
	val: "List/reverse",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 239, col: 5, offset: 6500},
	run: (*parser).callonReserved38,
	expr: &litMatcher{
	pos: position{line: 239, col: 5, offset: 6500},
	val: "Optional/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 240, col: 5, offset: 6551},
	run: (*parser).callonReserved40,
	expr: &litMatcher{
	pos: position{line: 240, col: 5, offset: 6551},
	val: "Optional/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 241, col: 5, offset: 6600},
	run: (*parser).callonReserved42,
	expr: &litMatcher{
	pos: position{line: 241, col: 5, offset: 6600},
	val: "Text/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 242, col: 5, offset: 6641},
	run: (*parser).callonReserved44,
	expr: &litMatcher{
	pos: position{line: 242, col: 5, offset: 6641},
	val: "Bool",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 243, col: 5, offset: 6673},
	run: (*parser).callonReserved46,
	expr: &litMatcher{
	pos: position{line: 243, col: 5, offset: 6673},
	val: "True",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 244, col: 5, offset: 6705},
	run: (*parser).callonReserved48,
	expr: &litMatcher{
	pos: position{line: 244, col: 5, offset: 6705},
	val: "False",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 245, col: 5, offset: 6739},
	run: (*parser).callonReserved50,
	expr: &litMatcher{
	pos: position{line: 245, col: 5, offset: 6739},
	val: "Optional",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 246, col: 5, offset: 6779},
	run: (*parser).callonReserved52,
	expr: &litMatcher{
	pos: position{line: 246, col: 5, offset: 6779},
	val: "Natural",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 247, col: 5, offset: 6817},
	run: (*parser).callonReserved54,
	expr: &litMatcher{
	pos: position{line: 247, col: 5, offset: 6817},
	val: "Integer",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 248, col: 5, offset: 6855},
	run: (*parser).callonReserved56,
	expr: &litMatcher{
	pos: position{line: 248, col: 5, offset: 6855},
	val: "Double",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 249, col: 5, offset: 6891},
	run: (*parser).callonReserved58,
	expr: &litMatcher{
	pos: position{line: 249, col: 5, offset: 6891},
	val: "Text",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 250, col: 5, offset: 6923},
	run: (*parser).callonReserved60,
	expr: &litMatcher{
	pos: position{line: 250, col: 5, offset: 6923},
	val: "List",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 251, col: 5, offset: 6955},
	run: (*parser).callonReserved62,
	expr: &litMatcher{
	pos: position{line: 251, col: 5, offset: 6955},
	val: "None",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 252, col: 5, offset: 6987},
	run: (*parser).callonReserved64,
	expr: &litMatcher{
	pos: position{line: 252, col: 5, offset: 6987},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 253, col: 5, offset: 7019},
	run: (*parser).callonReserved66,
	expr: &litMatcher{
	pos: position{line: 253, col: 5, offset: 7019},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 254, col: 5, offset: 7051},
	run: (*parser).callonReserved68,
	expr: &litMatcher{
	pos: position{line: 254, col: 5, offset: 7051},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "If",
	pos: position{line: 256, col: 1, offset: 7080},
	expr: &litMatcher{
	pos: position{line: 256, col: 6, offset: 7087},
	val: "if",
	ignoreCase: false,
},
},
{
	name: "Then",
	pos: position{line: 257, col: 1, offset: 7092},
	expr: &litMatcher{
	pos: position{line: 257, col: 8, offset: 7101},
	val: "then",
	ignoreCase: false,
},
},
{
	name: "Else",
	pos: position{line: 258, col: 1, offset: 7108},
	expr: &litMatcher{
	pos: position{line: 258, col: 8, offset: 7117},
	val: "else",
	ignoreCase: false,
},
},
{
	name: "Let",
	pos: position{line: 259, col: 1, offset: 7124},
	expr: &litMatcher{
	pos: position{line: 259, col: 7, offset: 7132},
	val: "let",
	ignoreCase: false,
},
},
{
	name: "In",
	pos: position{line: 260, col: 1, offset: 7138},
	expr: &litMatcher{
	pos: position{line: 260, col: 6, offset: 7145},
	val: "in",
	ignoreCase: false,
},
},
{
	name: "As",
	pos: position{line: 261, col: 1, offset: 7150},
	expr: &litMatcher{
	pos: position{line: 261, col: 6, offset: 7157},
	val: "as",
	ignoreCase: false,
},
},
{
	name: "Using",
	pos: position{line: 262, col: 1, offset: 7162},
	expr: &litMatcher{
	pos: position{line: 262, col: 9, offset: 7172},
	val: "using",
	ignoreCase: false,
},
},
{
	name: "Merge",
	pos: position{line: 263, col: 1, offset: 7180},
	expr: &litMatcher{
	pos: position{line: 263, col: 9, offset: 7190},
	val: "merge",
	ignoreCase: false,
},
},
{
	name: "Missing",
	pos: position{line: 264, col: 1, offset: 7198},
	expr: &actionExpr{
	pos: position{line: 264, col: 11, offset: 7210},
	run: (*parser).callonMissing1,
	expr: &litMatcher{
	pos: position{line: 264, col: 11, offset: 7210},
	val: "missing",
	ignoreCase: false,
},
},
},
{
	name: "True",
	pos: position{line: 265, col: 1, offset: 7246},
	expr: &litMatcher{
	pos: position{line: 265, col: 8, offset: 7255},
	val: "True",
	ignoreCase: false,
},
},
{
	name: "False",
	pos: position{line: 266, col: 1, offset: 7262},
	expr: &litMatcher{
	pos: position{line: 266, col: 9, offset: 7272},
	val: "False",
	ignoreCase: false,
},
},
{
	name: "Infinity",
	pos: position{line: 267, col: 1, offset: 7280},
	expr: &litMatcher{
	pos: position{line: 267, col: 12, offset: 7293},
	val: "Infinity",
	ignoreCase: false,
},
},
{
	name: "NaN",
	pos: position{line: 268, col: 1, offset: 7304},
	expr: &litMatcher{
	pos: position{line: 268, col: 7, offset: 7312},
	val: "NaN",
	ignoreCase: false,
},
},
{
	name: "Some",
	pos: position{line: 269, col: 1, offset: 7318},
	expr: &litMatcher{
	pos: position{line: 269, col: 8, offset: 7327},
	val: "Some",
	ignoreCase: false,
},
},
{
	name: "toMap",
	pos: position{line: 270, col: 1, offset: 7334},
	expr: &litMatcher{
	pos: position{line: 270, col: 9, offset: 7344},
	val: "toMap",
	ignoreCase: false,
},
},
{
	name: "assert",
	pos: position{line: 271, col: 1, offset: 7352},
	expr: &litMatcher{
	pos: position{line: 271, col: 10, offset: 7363},
	val: "assert",
	ignoreCase: false,
},
},
{
	name: "Keyword",
	pos: position{line: 273, col: 1, offset: 7373},
	expr: &choiceExpr{
	pos: position{line: 274, col: 5, offset: 7389},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 274, col: 5, offset: 7389},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 274, col: 10, offset: 7394},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 274, col: 17, offset: 7401},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 275, col: 5, offset: 7410},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 275, col: 11, offset: 7416},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 276, col: 5, offset: 7423},
	name: "Using",
},
&ruleRefExpr{
	pos: position{line: 276, col: 13, offset: 7431},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 276, col: 23, offset: 7441},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 277, col: 5, offset: 7448},
	name: "True",
},
&ruleRefExpr{
	pos: position{line: 277, col: 12, offset: 7455},
	name: "False",
},
&ruleRefExpr{
	pos: position{line: 278, col: 5, offset: 7465},
	name: "Infinity",
},
&ruleRefExpr{
	pos: position{line: 278, col: 16, offset: 7476},
	name: "NaN",
},
&ruleRefExpr{
	pos: position{line: 279, col: 5, offset: 7484},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 279, col: 13, offset: 7492},
	name: "Some",
},
&ruleRefExpr{
	pos: position{line: 279, col: 20, offset: 7499},
	name: "toMap",
},
	},
},
},
{
	name: "Optional",
	pos: position{line: 281, col: 1, offset: 7506},
	expr: &litMatcher{
	pos: position{line: 281, col: 12, offset: 7519},
	val: "Optional",
	ignoreCase: false,
},
},
{
	name: "Text",
	pos: position{line: 282, col: 1, offset: 7530},
	expr: &litMatcher{
	pos: position{line: 282, col: 8, offset: 7539},
	val: "Text",
	ignoreCase: false,
},
},
{
	name: "List",
	pos: position{line: 283, col: 1, offset: 7546},
	expr: &litMatcher{
	pos: position{line: 283, col: 8, offset: 7555},
	val: "List",
	ignoreCase: false,
},
},
{
	name: "Location",
	pos: position{line: 284, col: 1, offset: 7562},
	expr: &litMatcher{
	pos: position{line: 284, col: 12, offset: 7575},
	val: "Location",
	ignoreCase: false,
},
},
{
	name: "Combine",
	pos: position{line: 286, col: 1, offset: 7587},
	expr: &choiceExpr{
	pos: position{line: 286, col: 11, offset: 7599},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 286, col: 11, offset: 7599},
	val: "/\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 286, col: 19, offset: 7607},
	val: "∧",
	ignoreCase: false,
},
	},
},
},
{
	name: "CombineTypes",
	pos: position{line: 287, col: 1, offset: 7613},
	expr: &choiceExpr{
	pos: position{line: 287, col: 16, offset: 7630},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 287, col: 16, offset: 7630},
	val: "//\\\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 287, col: 27, offset: 7641},
	val: "⩓",
	ignoreCase: false,
},
	},
},
},
{
	name: "Equivalent",
	pos: position{line: 288, col: 1, offset: 7647},
	expr: &choiceExpr{
	pos: position{line: 288, col: 14, offset: 7662},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 288, col: 14, offset: 7662},
	val: "===",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 288, col: 22, offset: 7670},
	val: "≡",
	ignoreCase: false,
},
	},
},
},
{
	name: "Prefer",
	pos: position{line: 289, col: 1, offset: 7676},
	expr: &choiceExpr{
	pos: position{line: 289, col: 10, offset: 7687},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 289, col: 10, offset: 7687},
	val: "//",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 289, col: 17, offset: 7694},
	val: "⫽",
	ignoreCase: false,
},
	},
},
},
{
	name: "Lambda",
	pos: position{line: 290, col: 1, offset: 7700},
	expr: &choiceExpr{
	pos: position{line: 290, col: 10, offset: 7711},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 290, col: 10, offset: 7711},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 290, col: 17, offset: 7718},
	val: "λ",
	ignoreCase: false,
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 291, col: 1, offset: 7723},
	expr: &choiceExpr{
	pos: position{line: 291, col: 10, offset: 7734},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 291, col: 10, offset: 7734},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 291, col: 21, offset: 7745},
	val: "∀",
	ignoreCase: false,
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 292, col: 1, offset: 7751},
	expr: &choiceExpr{
	pos: position{line: 292, col: 9, offset: 7761},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 292, col: 9, offset: 7761},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 292, col: 16, offset: 7768},
	val: "→",
	ignoreCase: false,
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 294, col: 1, offset: 7775},
	expr: &seqExpr{
	pos: position{line: 294, col: 12, offset: 7788},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 294, col: 12, offset: 7788},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 294, col: 17, offset: 7793},
	expr: &charClassMatcher{
	pos: position{line: 294, col: 17, offset: 7793},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 294, col: 23, offset: 7799},
	expr: &ruleRefExpr{
	pos: position{line: 294, col: 23, offset: 7799},
	name: "Digit",
},
},
	},
},
},
{
	name: "NumericDoubleLiteral",
	pos: position{line: 296, col: 1, offset: 7807},
	expr: &actionExpr{
	pos: position{line: 296, col: 24, offset: 7832},
	run: (*parser).callonNumericDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 296, col: 24, offset: 7832},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 296, col: 24, offset: 7832},
	expr: &charClassMatcher{
	pos: position{line: 296, col: 24, offset: 7832},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 296, col: 30, offset: 7838},
	expr: &ruleRefExpr{
	pos: position{line: 296, col: 30, offset: 7838},
	name: "Digit",
},
},
&choiceExpr{
	pos: position{line: 296, col: 39, offset: 7847},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 296, col: 39, offset: 7847},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 296, col: 39, offset: 7847},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 296, col: 43, offset: 7851},
	expr: &ruleRefExpr{
	pos: position{line: 296, col: 43, offset: 7851},
	name: "Digit",
},
},
&zeroOrOneExpr{
	pos: position{line: 296, col: 50, offset: 7858},
	expr: &ruleRefExpr{
	pos: position{line: 296, col: 50, offset: 7858},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 296, col: 62, offset: 7870},
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
	pos: position{line: 304, col: 1, offset: 8026},
	expr: &choiceExpr{
	pos: position{line: 304, col: 17, offset: 8044},
	alternatives: []interface{}{
&labeledExpr{
	pos: position{line: 304, col: 17, offset: 8044},
	label: "d",
	expr: &ruleRefExpr{
	pos: position{line: 304, col: 19, offset: 8046},
	name: "NumericDoubleLiteral",
},
},
&actionExpr{
	pos: position{line: 305, col: 5, offset: 8071},
	run: (*parser).callonDoubleLiteral4,
	expr: &ruleRefExpr{
	pos: position{line: 305, col: 5, offset: 8071},
	name: "Infinity",
},
},
&actionExpr{
	pos: position{line: 306, col: 5, offset: 8123},
	run: (*parser).callonDoubleLiteral6,
	expr: &seqExpr{
	pos: position{line: 306, col: 5, offset: 8123},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 306, col: 5, offset: 8123},
	val: "-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 306, col: 9, offset: 8127},
	name: "Infinity",
},
	},
},
},
&actionExpr{
	pos: position{line: 307, col: 5, offset: 8180},
	run: (*parser).callonDoubleLiteral10,
	expr: &ruleRefExpr{
	pos: position{line: 307, col: 5, offset: 8180},
	name: "NaN",
},
},
	},
},
},
{
	name: "NaturalLiteral",
	pos: position{line: 309, col: 1, offset: 8223},
	expr: &actionExpr{
	pos: position{line: 309, col: 18, offset: 8242},
	run: (*parser).callonNaturalLiteral1,
	expr: &oneOrMoreExpr{
	pos: position{line: 309, col: 18, offset: 8242},
	expr: &ruleRefExpr{
	pos: position{line: 309, col: 18, offset: 8242},
	name: "Digit",
},
},
},
},
{
	name: "IntegerLiteral",
	pos: position{line: 314, col: 1, offset: 8331},
	expr: &actionExpr{
	pos: position{line: 314, col: 18, offset: 8350},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 314, col: 18, offset: 8350},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 314, col: 18, offset: 8350},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 314, col: 22, offset: 8354},
	name: "NaturalLiteral",
},
	},
},
},
},
{
	name: "DeBruijn",
	pos: position{line: 322, col: 1, offset: 8506},
	expr: &actionExpr{
	pos: position{line: 322, col: 12, offset: 8519},
	run: (*parser).callonDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 322, col: 12, offset: 8519},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 322, col: 12, offset: 8519},
	name: "_",
},
&litMatcher{
	pos: position{line: 322, col: 14, offset: 8521},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 322, col: 18, offset: 8525},
	name: "_",
},
&labeledExpr{
	pos: position{line: 322, col: 20, offset: 8527},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 322, col: 26, offset: 8533},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Variable",
	pos: position{line: 324, col: 1, offset: 8589},
	expr: &actionExpr{
	pos: position{line: 324, col: 12, offset: 8602},
	run: (*parser).callonVariable1,
	expr: &seqExpr{
	pos: position{line: 324, col: 12, offset: 8602},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 324, col: 12, offset: 8602},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 324, col: 17, offset: 8607},
	name: "NonreservedLabel",
},
},
&labeledExpr{
	pos: position{line: 324, col: 34, offset: 8624},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 324, col: 40, offset: 8630},
	expr: &ruleRefExpr{
	pos: position{line: 324, col: 40, offset: 8630},
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
	pos: position{line: 332, col: 1, offset: 8793},
	expr: &choiceExpr{
	pos: position{line: 332, col: 14, offset: 8808},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 332, col: 14, offset: 8808},
	name: "Variable",
},
&ruleRefExpr{
	pos: position{line: 332, col: 25, offset: 8819},
	name: "Reserved",
},
	},
},
},
{
	name: "PathCharacter",
	pos: position{line: 334, col: 1, offset: 8829},
	expr: &choiceExpr{
	pos: position{line: 335, col: 6, offset: 8852},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 335, col: 6, offset: 8852},
	val: "!",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 336, col: 6, offset: 8864},
	val: "[\\x24-\\x27]",
	ranges: []rune{'$','\'',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 337, col: 6, offset: 8881},
	val: "[\\x2a-\\x2b]",
	ranges: []rune{'*','+',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 338, col: 6, offset: 8898},
	val: "[\\x2d-\\x2e]",
	ranges: []rune{'-','.',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 339, col: 6, offset: 8915},
	val: "[\\x30-\\x3b]",
	ranges: []rune{'0',';',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 340, col: 6, offset: 8932},
	val: "=",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 341, col: 6, offset: 8944},
	val: "[\\x40-\\x5a]",
	ranges: []rune{'@','Z',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 342, col: 6, offset: 8961},
	val: "[\\x5e-\\x7a]",
	ranges: []rune{'^','z',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 343, col: 6, offset: 8978},
	val: "|",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 344, col: 6, offset: 8990},
	val: "~",
	ignoreCase: false,
},
	},
},
},
{
	name: "QuotedPathCharacter",
	pos: position{line: 346, col: 1, offset: 8998},
	expr: &choiceExpr{
	pos: position{line: 347, col: 6, offset: 9027},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 347, col: 6, offset: 9027},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 348, col: 6, offset: 9044},
	val: "[\\x23-\\x2e]",
	ranges: []rune{'#','.',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 349, col: 6, offset: 9061},
	val: "[\\x30-\\x7f]",
	ranges: []rune{'0','\u007f',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 350, col: 6, offset: 9078},
	name: "ValidNonAscii",
},
	},
},
},
{
	name: "UnquotedPathComponent",
	pos: position{line: 352, col: 1, offset: 9093},
	expr: &actionExpr{
	pos: position{line: 352, col: 25, offset: 9119},
	run: (*parser).callonUnquotedPathComponent1,
	expr: &oneOrMoreExpr{
	pos: position{line: 352, col: 25, offset: 9119},
	expr: &ruleRefExpr{
	pos: position{line: 352, col: 25, offset: 9119},
	name: "PathCharacter",
},
},
},
},
{
	name: "QuotedPathComponent",
	pos: position{line: 353, col: 1, offset: 9165},
	expr: &actionExpr{
	pos: position{line: 353, col: 23, offset: 9189},
	run: (*parser).callonQuotedPathComponent1,
	expr: &oneOrMoreExpr{
	pos: position{line: 353, col: 23, offset: 9189},
	expr: &ruleRefExpr{
	pos: position{line: 353, col: 23, offset: 9189},
	name: "QuotedPathCharacter",
},
},
},
},
{
	name: "PathComponent",
	pos: position{line: 355, col: 1, offset: 9242},
	expr: &choiceExpr{
	pos: position{line: 355, col: 17, offset: 9260},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 355, col: 17, offset: 9260},
	run: (*parser).callonPathComponent2,
	expr: &seqExpr{
	pos: position{line: 355, col: 17, offset: 9260},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 355, col: 17, offset: 9260},
	val: "/",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 355, col: 21, offset: 9264},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 355, col: 23, offset: 9266},
	name: "UnquotedPathComponent",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 356, col: 17, offset: 9322},
	run: (*parser).callonPathComponent7,
	expr: &seqExpr{
	pos: position{line: 356, col: 17, offset: 9322},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 356, col: 17, offset: 9322},
	val: "/",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 356, col: 21, offset: 9326},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 356, col: 25, offset: 9330},
	label: "q",
	expr: &ruleRefExpr{
	pos: position{line: 356, col: 27, offset: 9332},
	name: "QuotedPathComponent",
},
},
&litMatcher{
	pos: position{line: 356, col: 47, offset: 9352},
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
	pos: position{line: 358, col: 1, offset: 9375},
	expr: &actionExpr{
	pos: position{line: 358, col: 8, offset: 9384},
	run: (*parser).callonPath1,
	expr: &labeledExpr{
	pos: position{line: 358, col: 8, offset: 9384},
	label: "cs",
	expr: &oneOrMoreExpr{
	pos: position{line: 358, col: 11, offset: 9387},
	expr: &ruleRefExpr{
	pos: position{line: 358, col: 11, offset: 9387},
	name: "PathComponent",
},
},
},
},
},
{
	name: "Local",
	pos: position{line: 367, col: 1, offset: 9661},
	expr: &choiceExpr{
	pos: position{line: 367, col: 9, offset: 9671},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 367, col: 9, offset: 9671},
	name: "ParentPath",
},
&ruleRefExpr{
	pos: position{line: 367, col: 22, offset: 9684},
	name: "HerePath",
},
&ruleRefExpr{
	pos: position{line: 367, col: 33, offset: 9695},
	name: "HomePath",
},
&ruleRefExpr{
	pos: position{line: 367, col: 44, offset: 9706},
	name: "AbsolutePath",
},
	},
},
},
{
	name: "ParentPath",
	pos: position{line: 369, col: 1, offset: 9720},
	expr: &actionExpr{
	pos: position{line: 369, col: 14, offset: 9735},
	run: (*parser).callonParentPath1,
	expr: &seqExpr{
	pos: position{line: 369, col: 14, offset: 9735},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 369, col: 14, offset: 9735},
	val: "..",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 369, col: 19, offset: 9740},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 369, col: 21, offset: 9742},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HerePath",
	pos: position{line: 370, col: 1, offset: 9798},
	expr: &actionExpr{
	pos: position{line: 370, col: 12, offset: 9811},
	run: (*parser).callonHerePath1,
	expr: &seqExpr{
	pos: position{line: 370, col: 12, offset: 9811},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 370, col: 12, offset: 9811},
	val: ".",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 370, col: 16, offset: 9815},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 370, col: 18, offset: 9817},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HomePath",
	pos: position{line: 371, col: 1, offset: 9856},
	expr: &actionExpr{
	pos: position{line: 371, col: 12, offset: 9869},
	run: (*parser).callonHomePath1,
	expr: &seqExpr{
	pos: position{line: 371, col: 12, offset: 9869},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 371, col: 12, offset: 9869},
	val: "~",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 371, col: 16, offset: 9873},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 371, col: 18, offset: 9875},
	name: "Path",
},
},
	},
},
},
},
{
	name: "AbsolutePath",
	pos: position{line: 372, col: 1, offset: 9930},
	expr: &actionExpr{
	pos: position{line: 372, col: 16, offset: 9947},
	run: (*parser).callonAbsolutePath1,
	expr: &labeledExpr{
	pos: position{line: 372, col: 16, offset: 9947},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 372, col: 18, offset: 9949},
	name: "Path",
},
},
},
},
{
	name: "Scheme",
	pos: position{line: 374, col: 1, offset: 10005},
	expr: &seqExpr{
	pos: position{line: 374, col: 10, offset: 10016},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 374, col: 10, offset: 10016},
	val: "http",
	ignoreCase: false,
},
&zeroOrOneExpr{
	pos: position{line: 374, col: 17, offset: 10023},
	expr: &litMatcher{
	pos: position{line: 374, col: 17, offset: 10023},
	val: "s",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "HttpRaw",
	pos: position{line: 376, col: 1, offset: 10029},
	expr: &actionExpr{
	pos: position{line: 376, col: 11, offset: 10041},
	run: (*parser).callonHttpRaw1,
	expr: &seqExpr{
	pos: position{line: 376, col: 11, offset: 10041},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 376, col: 11, offset: 10041},
	name: "Scheme",
},
&litMatcher{
	pos: position{line: 376, col: 18, offset: 10048},
	val: "://",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 376, col: 24, offset: 10054},
	name: "Authority",
},
&ruleRefExpr{
	pos: position{line: 376, col: 34, offset: 10064},
	name: "UrlPath",
},
&zeroOrOneExpr{
	pos: position{line: 376, col: 42, offset: 10072},
	expr: &seqExpr{
	pos: position{line: 376, col: 44, offset: 10074},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 376, col: 44, offset: 10074},
	val: "?",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 376, col: 48, offset: 10078},
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
	pos: position{line: 378, col: 1, offset: 10135},
	expr: &zeroOrMoreExpr{
	pos: position{line: 378, col: 11, offset: 10147},
	expr: &choiceExpr{
	pos: position{line: 378, col: 12, offset: 10148},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 378, col: 12, offset: 10148},
	name: "PathComponent",
},
&seqExpr{
	pos: position{line: 378, col: 28, offset: 10164},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 378, col: 28, offset: 10164},
	val: "/",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 378, col: 32, offset: 10168},
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
	pos: position{line: 380, col: 1, offset: 10179},
	expr: &seqExpr{
	pos: position{line: 380, col: 13, offset: 10193},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 380, col: 13, offset: 10193},
	expr: &seqExpr{
	pos: position{line: 380, col: 14, offset: 10194},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 380, col: 14, offset: 10194},
	name: "Userinfo",
},
&litMatcher{
	pos: position{line: 380, col: 23, offset: 10203},
	val: "@",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 380, col: 29, offset: 10209},
	name: "Host",
},
&zeroOrOneExpr{
	pos: position{line: 380, col: 34, offset: 10214},
	expr: &seqExpr{
	pos: position{line: 380, col: 35, offset: 10215},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 380, col: 35, offset: 10215},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 380, col: 39, offset: 10219},
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
	pos: position{line: 382, col: 1, offset: 10227},
	expr: &zeroOrMoreExpr{
	pos: position{line: 382, col: 12, offset: 10240},
	expr: &choiceExpr{
	pos: position{line: 382, col: 14, offset: 10242},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 382, col: 14, offset: 10242},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 382, col: 27, offset: 10255},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 382, col: 40, offset: 10268},
	name: "SubDelims",
},
&litMatcher{
	pos: position{line: 382, col: 52, offset: 10280},
	val: ":",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "Host",
	pos: position{line: 384, col: 1, offset: 10288},
	expr: &choiceExpr{
	pos: position{line: 384, col: 8, offset: 10297},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 384, col: 8, offset: 10297},
	name: "IPLiteral",
},
&ruleRefExpr{
	pos: position{line: 384, col: 20, offset: 10309},
	name: "RegName",
},
	},
},
},
{
	name: "Port",
	pos: position{line: 386, col: 1, offset: 10318},
	expr: &zeroOrMoreExpr{
	pos: position{line: 386, col: 8, offset: 10327},
	expr: &ruleRefExpr{
	pos: position{line: 386, col: 8, offset: 10327},
	name: "Digit",
},
},
},
{
	name: "IPLiteral",
	pos: position{line: 388, col: 1, offset: 10335},
	expr: &seqExpr{
	pos: position{line: 388, col: 13, offset: 10349},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 388, col: 13, offset: 10349},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 388, col: 17, offset: 10353},
	name: "IPv6address",
},
&litMatcher{
	pos: position{line: 388, col: 29, offset: 10365},
	val: "]",
	ignoreCase: false,
},
	},
},
},
{
	name: "IPv6address",
	pos: position{line: 390, col: 1, offset: 10370},
	expr: &actionExpr{
	pos: position{line: 390, col: 15, offset: 10386},
	run: (*parser).callonIPv6address1,
	expr: &seqExpr{
	pos: position{line: 390, col: 15, offset: 10386},
	exprs: []interface{}{
&zeroOrMoreExpr{
	pos: position{line: 390, col: 15, offset: 10386},
	expr: &ruleRefExpr{
	pos: position{line: 390, col: 16, offset: 10387},
	name: "HexDig",
},
},
&litMatcher{
	pos: position{line: 390, col: 25, offset: 10396},
	val: ":",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 390, col: 29, offset: 10400},
	expr: &choiceExpr{
	pos: position{line: 390, col: 30, offset: 10401},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 390, col: 30, offset: 10401},
	name: "HexDig",
},
&litMatcher{
	pos: position{line: 390, col: 39, offset: 10410},
	val: ":",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 390, col: 45, offset: 10416},
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
	pos: position{line: 396, col: 1, offset: 10570},
	expr: &zeroOrMoreExpr{
	pos: position{line: 396, col: 11, offset: 10582},
	expr: &choiceExpr{
	pos: position{line: 396, col: 12, offset: 10583},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 396, col: 12, offset: 10583},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 396, col: 25, offset: 10596},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 396, col: 38, offset: 10609},
	name: "SubDelims",
},
	},
},
},
},
{
	name: "Segment",
	pos: position{line: 398, col: 1, offset: 10622},
	expr: &zeroOrMoreExpr{
	pos: position{line: 398, col: 11, offset: 10634},
	expr: &ruleRefExpr{
	pos: position{line: 398, col: 11, offset: 10634},
	name: "PChar",
},
},
},
{
	name: "PChar",
	pos: position{line: 400, col: 1, offset: 10642},
	expr: &choiceExpr{
	pos: position{line: 400, col: 9, offset: 10652},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 400, col: 9, offset: 10652},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 400, col: 22, offset: 10665},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 400, col: 35, offset: 10678},
	name: "SubDelims",
},
&charClassMatcher{
	pos: position{line: 400, col: 47, offset: 10690},
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
	pos: position{line: 402, col: 1, offset: 10696},
	expr: &zeroOrMoreExpr{
	pos: position{line: 402, col: 9, offset: 10706},
	expr: &choiceExpr{
	pos: position{line: 402, col: 10, offset: 10707},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 402, col: 10, offset: 10707},
	name: "PChar",
},
&charClassMatcher{
	pos: position{line: 402, col: 18, offset: 10715},
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
	pos: position{line: 404, col: 1, offset: 10723},
	expr: &seqExpr{
	pos: position{line: 404, col: 14, offset: 10738},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 404, col: 14, offset: 10738},
	val: "%",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 404, col: 18, offset: 10742},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 404, col: 25, offset: 10749},
	name: "HexDig",
},
	},
},
},
{
	name: "Unreserved",
	pos: position{line: 406, col: 1, offset: 10757},
	expr: &charClassMatcher{
	pos: position{line: 406, col: 14, offset: 10772},
	val: "[A-Za-z0-9._~-]",
	chars: []rune{'.','_','~','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SubDelims",
	pos: position{line: 408, col: 1, offset: 10789},
	expr: &choiceExpr{
	pos: position{line: 408, col: 13, offset: 10803},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 408, col: 13, offset: 10803},
	val: "!",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 408, col: 19, offset: 10809},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 408, col: 25, offset: 10815},
	val: "&",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 408, col: 31, offset: 10821},
	val: "'",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 408, col: 37, offset: 10827},
	val: "*",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 408, col: 43, offset: 10833},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 408, col: 49, offset: 10839},
	val: ";",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 408, col: 55, offset: 10845},
	val: "=",
	ignoreCase: false,
},
	},
},
},
{
	name: "Http",
	pos: position{line: 410, col: 1, offset: 10850},
	expr: &actionExpr{
	pos: position{line: 410, col: 8, offset: 10859},
	run: (*parser).callonHttp1,
	expr: &labeledExpr{
	pos: position{line: 410, col: 8, offset: 10859},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 410, col: 10, offset: 10861},
	name: "HttpRaw",
},
},
},
},
{
	name: "Env",
	pos: position{line: 412, col: 1, offset: 10911},
	expr: &actionExpr{
	pos: position{line: 412, col: 7, offset: 10919},
	run: (*parser).callonEnv1,
	expr: &seqExpr{
	pos: position{line: 412, col: 7, offset: 10919},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 412, col: 7, offset: 10919},
	val: "env:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 412, col: 14, offset: 10926},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 412, col: 17, offset: 10929},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 412, col: 17, offset: 10929},
	name: "BashEnvironmentVariable",
},
&ruleRefExpr{
	pos: position{line: 412, col: 43, offset: 10955},
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
	pos: position{line: 414, col: 1, offset: 11000},
	expr: &actionExpr{
	pos: position{line: 414, col: 27, offset: 11028},
	run: (*parser).callonBashEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 414, col: 27, offset: 11028},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 414, col: 27, offset: 11028},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 414, col: 36, offset: 11037},
	expr: &charClassMatcher{
	pos: position{line: 414, col: 36, offset: 11037},
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
	pos: position{line: 418, col: 1, offset: 11093},
	expr: &actionExpr{
	pos: position{line: 418, col: 28, offset: 11122},
	run: (*parser).callonPosixEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 418, col: 28, offset: 11122},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 418, col: 28, offset: 11122},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 418, col: 32, offset: 11126},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 418, col: 34, offset: 11128},
	name: "PosixEnvironmentVariableContent",
},
},
&litMatcher{
	pos: position{line: 418, col: 66, offset: 11160},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PosixEnvironmentVariableContent",
	pos: position{line: 422, col: 1, offset: 11185},
	expr: &actionExpr{
	pos: position{line: 422, col: 35, offset: 11221},
	run: (*parser).callonPosixEnvironmentVariableContent1,
	expr: &labeledExpr{
	pos: position{line: 422, col: 35, offset: 11221},
	label: "v",
	expr: &oneOrMoreExpr{
	pos: position{line: 422, col: 37, offset: 11223},
	expr: &ruleRefExpr{
	pos: position{line: 422, col: 37, offset: 11223},
	name: "PosixEnvironmentVariableCharacter",
},
},
},
},
},
{
	name: "PosixEnvironmentVariableCharacter",
	pos: position{line: 431, col: 1, offset: 11436},
	expr: &choiceExpr{
	pos: position{line: 432, col: 7, offset: 11480},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 432, col: 7, offset: 11480},
	run: (*parser).callonPosixEnvironmentVariableCharacter2,
	expr: &litMatcher{
	pos: position{line: 432, col: 7, offset: 11480},
	val: "\\\"",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 433, col: 7, offset: 11520},
	run: (*parser).callonPosixEnvironmentVariableCharacter4,
	expr: &litMatcher{
	pos: position{line: 433, col: 7, offset: 11520},
	val: "\\\\",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 434, col: 7, offset: 11560},
	run: (*parser).callonPosixEnvironmentVariableCharacter6,
	expr: &litMatcher{
	pos: position{line: 434, col: 7, offset: 11560},
	val: "\\a",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 435, col: 7, offset: 11600},
	run: (*parser).callonPosixEnvironmentVariableCharacter8,
	expr: &litMatcher{
	pos: position{line: 435, col: 7, offset: 11600},
	val: "\\b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 436, col: 7, offset: 11640},
	run: (*parser).callonPosixEnvironmentVariableCharacter10,
	expr: &litMatcher{
	pos: position{line: 436, col: 7, offset: 11640},
	val: "\\f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 437, col: 7, offset: 11680},
	run: (*parser).callonPosixEnvironmentVariableCharacter12,
	expr: &litMatcher{
	pos: position{line: 437, col: 7, offset: 11680},
	val: "\\n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 438, col: 7, offset: 11720},
	run: (*parser).callonPosixEnvironmentVariableCharacter14,
	expr: &litMatcher{
	pos: position{line: 438, col: 7, offset: 11720},
	val: "\\r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 439, col: 7, offset: 11760},
	run: (*parser).callonPosixEnvironmentVariableCharacter16,
	expr: &litMatcher{
	pos: position{line: 439, col: 7, offset: 11760},
	val: "\\t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 440, col: 7, offset: 11800},
	run: (*parser).callonPosixEnvironmentVariableCharacter18,
	expr: &litMatcher{
	pos: position{line: 440, col: 7, offset: 11800},
	val: "\\v",
	ignoreCase: false,
},
},
&charClassMatcher{
	pos: position{line: 441, col: 7, offset: 11840},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 442, col: 7, offset: 11858},
	val: "[\\x23-\\x3c]",
	ranges: []rune{'#','<',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 443, col: 7, offset: 11876},
	val: "[\\x3e-\\x5b]",
	ranges: []rune{'>','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 444, col: 7, offset: 11894},
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
	pos: position{line: 446, col: 1, offset: 11907},
	expr: &choiceExpr{
	pos: position{line: 446, col: 14, offset: 11922},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 446, col: 14, offset: 11922},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 446, col: 24, offset: 11932},
	name: "Local",
},
&ruleRefExpr{
	pos: position{line: 446, col: 32, offset: 11940},
	name: "Http",
},
&ruleRefExpr{
	pos: position{line: 446, col: 39, offset: 11947},
	name: "Env",
},
	},
},
},
{
	name: "HashValue",
	pos: position{line: 449, col: 1, offset: 12020},
	expr: &actionExpr{
	pos: position{line: 449, col: 13, offset: 12032},
	run: (*parser).callonHashValue1,
	expr: &seqExpr{
	pos: position{line: 449, col: 13, offset: 12032},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 449, col: 13, offset: 12032},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 449, col: 20, offset: 12039},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 449, col: 27, offset: 12046},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 449, col: 34, offset: 12053},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 449, col: 41, offset: 12060},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 449, col: 48, offset: 12067},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 449, col: 55, offset: 12074},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 449, col: 62, offset: 12081},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 450, col: 13, offset: 12100},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 450, col: 20, offset: 12107},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 450, col: 27, offset: 12114},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 450, col: 34, offset: 12121},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 450, col: 41, offset: 12128},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 450, col: 48, offset: 12135},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 450, col: 55, offset: 12142},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 450, col: 62, offset: 12149},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 451, col: 13, offset: 12168},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 451, col: 20, offset: 12175},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 451, col: 27, offset: 12182},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 451, col: 34, offset: 12189},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 451, col: 41, offset: 12196},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 451, col: 48, offset: 12203},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 451, col: 55, offset: 12210},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 451, col: 62, offset: 12217},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 452, col: 13, offset: 12236},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 452, col: 20, offset: 12243},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 452, col: 27, offset: 12250},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 452, col: 34, offset: 12257},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 452, col: 41, offset: 12264},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 452, col: 48, offset: 12271},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 452, col: 55, offset: 12278},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 452, col: 62, offset: 12285},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 453, col: 13, offset: 12304},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 453, col: 20, offset: 12311},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 453, col: 27, offset: 12318},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 453, col: 34, offset: 12325},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 453, col: 41, offset: 12332},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 453, col: 48, offset: 12339},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 453, col: 55, offset: 12346},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 453, col: 62, offset: 12353},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 454, col: 13, offset: 12372},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 454, col: 20, offset: 12379},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 454, col: 27, offset: 12386},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 454, col: 34, offset: 12393},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 454, col: 41, offset: 12400},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 454, col: 48, offset: 12407},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 454, col: 55, offset: 12414},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 454, col: 62, offset: 12421},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 455, col: 13, offset: 12440},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 455, col: 20, offset: 12447},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 455, col: 27, offset: 12454},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 455, col: 34, offset: 12461},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 455, col: 41, offset: 12468},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 455, col: 48, offset: 12475},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 455, col: 55, offset: 12482},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 455, col: 62, offset: 12489},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 456, col: 13, offset: 12508},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 456, col: 20, offset: 12515},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 456, col: 27, offset: 12522},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 456, col: 34, offset: 12529},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 456, col: 41, offset: 12536},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 456, col: 48, offset: 12543},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 456, col: 55, offset: 12550},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 456, col: 62, offset: 12557},
	name: "HexDig",
},
	},
},
},
},
{
	name: "Hash",
	pos: position{line: 462, col: 1, offset: 12701},
	expr: &actionExpr{
	pos: position{line: 462, col: 8, offset: 12708},
	run: (*parser).callonHash1,
	expr: &seqExpr{
	pos: position{line: 462, col: 8, offset: 12708},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 462, col: 8, offset: 12708},
	val: "sha256:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 462, col: 18, offset: 12718},
	label: "val",
	expr: &ruleRefExpr{
	pos: position{line: 462, col: 22, offset: 12722},
	name: "HashValue",
},
},
	},
},
},
},
{
	name: "ImportHashed",
	pos: position{line: 464, col: 1, offset: 12792},
	expr: &actionExpr{
	pos: position{line: 464, col: 16, offset: 12809},
	run: (*parser).callonImportHashed1,
	expr: &seqExpr{
	pos: position{line: 464, col: 16, offset: 12809},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 464, col: 16, offset: 12809},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 464, col: 18, offset: 12811},
	name: "ImportType",
},
},
&labeledExpr{
	pos: position{line: 464, col: 29, offset: 12822},
	label: "h",
	expr: &zeroOrOneExpr{
	pos: position{line: 464, col: 31, offset: 12824},
	expr: &seqExpr{
	pos: position{line: 464, col: 32, offset: 12825},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 464, col: 32, offset: 12825},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 464, col: 35, offset: 12828},
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
	pos: position{line: 472, col: 1, offset: 12983},
	expr: &choiceExpr{
	pos: position{line: 472, col: 10, offset: 12994},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 472, col: 10, offset: 12994},
	run: (*parser).callonImport2,
	expr: &seqExpr{
	pos: position{line: 472, col: 10, offset: 12994},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 472, col: 10, offset: 12994},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 472, col: 12, offset: 12996},
	name: "ImportHashed",
},
},
&ruleRefExpr{
	pos: position{line: 472, col: 25, offset: 13009},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 472, col: 27, offset: 13011},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 472, col: 30, offset: 13014},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 472, col: 33, offset: 13017},
	name: "Text",
},
	},
},
},
&actionExpr{
	pos: position{line: 473, col: 10, offset: 13114},
	run: (*parser).callonImport10,
	expr: &seqExpr{
	pos: position{line: 473, col: 10, offset: 13114},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 473, col: 10, offset: 13114},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 473, col: 12, offset: 13116},
	name: "ImportHashed",
},
},
&ruleRefExpr{
	pos: position{line: 473, col: 25, offset: 13129},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 473, col: 27, offset: 13131},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 473, col: 30, offset: 13134},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 473, col: 33, offset: 13137},
	name: "Location",
},
	},
},
},
&actionExpr{
	pos: position{line: 474, col: 10, offset: 13239},
	run: (*parser).callonImport18,
	expr: &labeledExpr{
	pos: position{line: 474, col: 10, offset: 13239},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 474, col: 12, offset: 13241},
	name: "ImportHashed",
},
},
},
	},
},
},
{
	name: "LetBinding",
	pos: position{line: 477, col: 1, offset: 13336},
	expr: &actionExpr{
	pos: position{line: 477, col: 14, offset: 13351},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 477, col: 14, offset: 13351},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 477, col: 14, offset: 13351},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 477, col: 18, offset: 13355},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 477, col: 21, offset: 13358},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 477, col: 27, offset: 13364},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 477, col: 44, offset: 13381},
	name: "_",
},
&labeledExpr{
	pos: position{line: 477, col: 46, offset: 13383},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 477, col: 48, offset: 13385},
	expr: &seqExpr{
	pos: position{line: 477, col: 49, offset: 13386},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 477, col: 49, offset: 13386},
	name: "Annotation",
},
&ruleRefExpr{
	pos: position{line: 477, col: 60, offset: 13397},
	name: "_",
},
	},
},
},
},
&litMatcher{
	pos: position{line: 478, col: 13, offset: 13413},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 478, col: 17, offset: 13417},
	name: "_",
},
&labeledExpr{
	pos: position{line: 478, col: 19, offset: 13419},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 478, col: 21, offset: 13421},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 478, col: 32, offset: 13432},
	name: "_",
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 493, col: 1, offset: 13741},
	expr: &choiceExpr{
	pos: position{line: 494, col: 7, offset: 13762},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 494, col: 7, offset: 13762},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 494, col: 7, offset: 13762},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 494, col: 7, offset: 13762},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 494, col: 14, offset: 13769},
	name: "_",
},
&litMatcher{
	pos: position{line: 494, col: 16, offset: 13771},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 494, col: 20, offset: 13775},
	name: "_",
},
&labeledExpr{
	pos: position{line: 494, col: 22, offset: 13777},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 494, col: 28, offset: 13783},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 494, col: 45, offset: 13800},
	name: "_",
},
&litMatcher{
	pos: position{line: 494, col: 47, offset: 13802},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 494, col: 51, offset: 13806},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 494, col: 54, offset: 13809},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 494, col: 56, offset: 13811},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 494, col: 67, offset: 13822},
	name: "_",
},
&litMatcher{
	pos: position{line: 494, col: 69, offset: 13824},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 494, col: 73, offset: 13828},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 494, col: 75, offset: 13830},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 494, col: 81, offset: 13836},
	name: "_",
},
&labeledExpr{
	pos: position{line: 494, col: 83, offset: 13838},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 494, col: 88, offset: 13843},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 497, col: 7, offset: 13960},
	run: (*parser).callonExpression22,
	expr: &seqExpr{
	pos: position{line: 497, col: 7, offset: 13960},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 497, col: 7, offset: 13960},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 497, col: 10, offset: 13963},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 497, col: 13, offset: 13966},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 497, col: 18, offset: 13971},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 497, col: 29, offset: 13982},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 497, col: 31, offset: 13984},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 497, col: 36, offset: 13989},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 497, col: 39, offset: 13992},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 497, col: 41, offset: 13994},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 497, col: 52, offset: 14005},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 497, col: 54, offset: 14007},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 497, col: 59, offset: 14012},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 497, col: 62, offset: 14015},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 497, col: 64, offset: 14017},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 500, col: 7, offset: 14103},
	run: (*parser).callonExpression38,
	expr: &seqExpr{
	pos: position{line: 500, col: 7, offset: 14103},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 500, col: 7, offset: 14103},
	label: "bindings",
	expr: &oneOrMoreExpr{
	pos: position{line: 500, col: 16, offset: 14112},
	expr: &ruleRefExpr{
	pos: position{line: 500, col: 16, offset: 14112},
	name: "LetBinding",
},
},
},
&ruleRefExpr{
	pos: position{line: 500, col: 28, offset: 14124},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 500, col: 31, offset: 14127},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 500, col: 34, offset: 14130},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 500, col: 36, offset: 14132},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 507, col: 7, offset: 14372},
	run: (*parser).callonExpression47,
	expr: &seqExpr{
	pos: position{line: 507, col: 7, offset: 14372},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 507, col: 7, offset: 14372},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 507, col: 14, offset: 14379},
	name: "_",
},
&litMatcher{
	pos: position{line: 507, col: 16, offset: 14381},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 507, col: 20, offset: 14385},
	name: "_",
},
&labeledExpr{
	pos: position{line: 507, col: 22, offset: 14387},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 507, col: 28, offset: 14393},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 507, col: 45, offset: 14410},
	name: "_",
},
&litMatcher{
	pos: position{line: 507, col: 47, offset: 14412},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 507, col: 51, offset: 14416},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 507, col: 54, offset: 14419},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 507, col: 56, offset: 14421},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 507, col: 67, offset: 14432},
	name: "_",
},
&litMatcher{
	pos: position{line: 507, col: 69, offset: 14434},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 507, col: 73, offset: 14438},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 507, col: 75, offset: 14440},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 507, col: 81, offset: 14446},
	name: "_",
},
&labeledExpr{
	pos: position{line: 507, col: 83, offset: 14448},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 507, col: 88, offset: 14453},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 510, col: 7, offset: 14562},
	run: (*parser).callonExpression67,
	expr: &seqExpr{
	pos: position{line: 510, col: 7, offset: 14562},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 510, col: 7, offset: 14562},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 510, col: 9, offset: 14564},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 510, col: 28, offset: 14583},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 510, col: 30, offset: 14585},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 510, col: 36, offset: 14591},
	name: "_",
},
&labeledExpr{
	pos: position{line: 510, col: 38, offset: 14593},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 510, col: 40, offset: 14595},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 511, col: 7, offset: 14654},
	run: (*parser).callonExpression76,
	expr: &seqExpr{
	pos: position{line: 511, col: 7, offset: 14654},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 511, col: 7, offset: 14654},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 511, col: 13, offset: 14660},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 511, col: 16, offset: 14663},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 511, col: 18, offset: 14665},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 511, col: 35, offset: 14682},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 511, col: 38, offset: 14685},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 511, col: 40, offset: 14687},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 511, col: 57, offset: 14704},
	name: "_",
},
&litMatcher{
	pos: position{line: 511, col: 59, offset: 14706},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 511, col: 63, offset: 14710},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 511, col: 66, offset: 14713},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 511, col: 68, offset: 14715},
	name: "ApplicationExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 514, col: 7, offset: 14836},
	name: "EmptyList",
},
&actionExpr{
	pos: position{line: 515, col: 7, offset: 14852},
	run: (*parser).callonExpression91,
	expr: &seqExpr{
	pos: position{line: 515, col: 7, offset: 14852},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 515, col: 7, offset: 14852},
	name: "toMap",
},
&ruleRefExpr{
	pos: position{line: 515, col: 13, offset: 14858},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 515, col: 16, offset: 14861},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 515, col: 18, offset: 14863},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 515, col: 35, offset: 14880},
	name: "_",
},
&litMatcher{
	pos: position{line: 515, col: 37, offset: 14882},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 515, col: 41, offset: 14886},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 515, col: 44, offset: 14889},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 515, col: 46, offset: 14891},
	name: "ApplicationExpression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 516, col: 7, offset: 14961},
	run: (*parser).callonExpression102,
	expr: &seqExpr{
	pos: position{line: 516, col: 7, offset: 14961},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 516, col: 7, offset: 14961},
	name: "assert",
},
&ruleRefExpr{
	pos: position{line: 516, col: 14, offset: 14968},
	name: "_",
},
&litMatcher{
	pos: position{line: 516, col: 16, offset: 14970},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 516, col: 20, offset: 14974},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 516, col: 23, offset: 14977},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 516, col: 25, offset: 14979},
	name: "Expression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 517, col: 7, offset: 15041},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 519, col: 1, offset: 15062},
	expr: &actionExpr{
	pos: position{line: 519, col: 14, offset: 15077},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 519, col: 14, offset: 15077},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 519, col: 14, offset: 15077},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 519, col: 18, offset: 15081},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 519, col: 21, offset: 15084},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 519, col: 23, offset: 15086},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 521, col: 1, offset: 15116},
	expr: &actionExpr{
	pos: position{line: 522, col: 1, offset: 15140},
	run: (*parser).callonAnnotatedExpression1,
	expr: &seqExpr{
	pos: position{line: 522, col: 1, offset: 15140},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 522, col: 1, offset: 15140},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 522, col: 3, offset: 15142},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 522, col: 22, offset: 15161},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 522, col: 24, offset: 15163},
	expr: &seqExpr{
	pos: position{line: 522, col: 25, offset: 15164},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 522, col: 25, offset: 15164},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 522, col: 27, offset: 15166},
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
	pos: position{line: 527, col: 1, offset: 15291},
	expr: &actionExpr{
	pos: position{line: 527, col: 13, offset: 15305},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 527, col: 13, offset: 15305},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 527, col: 13, offset: 15305},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 527, col: 17, offset: 15309},
	name: "_",
},
&litMatcher{
	pos: position{line: 527, col: 19, offset: 15311},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 527, col: 23, offset: 15315},
	name: "_",
},
&litMatcher{
	pos: position{line: 527, col: 25, offset: 15317},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 527, col: 29, offset: 15321},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 527, col: 32, offset: 15324},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 527, col: 34, offset: 15326},
	name: "ApplicationExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 531, col: 1, offset: 15394},
	expr: &ruleRefExpr{
	pos: position{line: 531, col: 22, offset: 15417},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 533, col: 1, offset: 15438},
	expr: &actionExpr{
	pos: position{line: 533, col: 26, offset: 15465},
	run: (*parser).callonImportAltExpression1,
	expr: &seqExpr{
	pos: position{line: 533, col: 26, offset: 15465},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 533, col: 26, offset: 15465},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 533, col: 32, offset: 15471},
	name: "OrExpression",
},
},
&labeledExpr{
	pos: position{line: 533, col: 55, offset: 15494},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 533, col: 60, offset: 15499},
	expr: &seqExpr{
	pos: position{line: 533, col: 61, offset: 15500},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 533, col: 61, offset: 15500},
	name: "_",
},
&litMatcher{
	pos: position{line: 533, col: 63, offset: 15502},
	val: "?",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 533, col: 67, offset: 15506},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 533, col: 70, offset: 15509},
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
	pos: position{line: 535, col: 1, offset: 15580},
	expr: &actionExpr{
	pos: position{line: 535, col: 26, offset: 15607},
	run: (*parser).callonOrExpression1,
	expr: &seqExpr{
	pos: position{line: 535, col: 26, offset: 15607},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 535, col: 26, offset: 15607},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 535, col: 32, offset: 15613},
	name: "PlusExpression",
},
},
&labeledExpr{
	pos: position{line: 535, col: 55, offset: 15636},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 535, col: 60, offset: 15641},
	expr: &seqExpr{
	pos: position{line: 535, col: 61, offset: 15642},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 535, col: 61, offset: 15642},
	name: "_",
},
&litMatcher{
	pos: position{line: 535, col: 63, offset: 15644},
	val: "||",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 535, col: 68, offset: 15649},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 535, col: 70, offset: 15651},
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
	pos: position{line: 537, col: 1, offset: 15717},
	expr: &actionExpr{
	pos: position{line: 537, col: 26, offset: 15744},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 537, col: 26, offset: 15744},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 537, col: 26, offset: 15744},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 537, col: 32, offset: 15750},
	name: "TextAppendExpression",
},
},
&labeledExpr{
	pos: position{line: 537, col: 55, offset: 15773},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 537, col: 60, offset: 15778},
	expr: &seqExpr{
	pos: position{line: 537, col: 61, offset: 15779},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 537, col: 61, offset: 15779},
	name: "_",
},
&litMatcher{
	pos: position{line: 537, col: 63, offset: 15781},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 537, col: 67, offset: 15785},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 537, col: 70, offset: 15788},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 537, col: 72, offset: 15790},
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
	pos: position{line: 539, col: 1, offset: 15864},
	expr: &actionExpr{
	pos: position{line: 539, col: 26, offset: 15891},
	run: (*parser).callonTextAppendExpression1,
	expr: &seqExpr{
	pos: position{line: 539, col: 26, offset: 15891},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 539, col: 26, offset: 15891},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 539, col: 32, offset: 15897},
	name: "ListAppendExpression",
},
},
&labeledExpr{
	pos: position{line: 539, col: 55, offset: 15920},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 539, col: 60, offset: 15925},
	expr: &seqExpr{
	pos: position{line: 539, col: 61, offset: 15926},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 539, col: 61, offset: 15926},
	name: "_",
},
&litMatcher{
	pos: position{line: 539, col: 63, offset: 15928},
	val: "++",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 539, col: 68, offset: 15933},
	name: "_",
},
&labeledExpr{
	pos: position{line: 539, col: 70, offset: 15935},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 539, col: 72, offset: 15937},
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
	pos: position{line: 541, col: 1, offset: 16017},
	expr: &actionExpr{
	pos: position{line: 541, col: 26, offset: 16044},
	run: (*parser).callonListAppendExpression1,
	expr: &seqExpr{
	pos: position{line: 541, col: 26, offset: 16044},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 541, col: 26, offset: 16044},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 541, col: 32, offset: 16050},
	name: "AndExpression",
},
},
&labeledExpr{
	pos: position{line: 541, col: 55, offset: 16073},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 541, col: 60, offset: 16078},
	expr: &seqExpr{
	pos: position{line: 541, col: 61, offset: 16079},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 541, col: 61, offset: 16079},
	name: "_",
},
&litMatcher{
	pos: position{line: 541, col: 63, offset: 16081},
	val: "#",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 541, col: 67, offset: 16085},
	name: "_",
},
&labeledExpr{
	pos: position{line: 541, col: 69, offset: 16087},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 541, col: 71, offset: 16089},
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
	pos: position{line: 543, col: 1, offset: 16162},
	expr: &actionExpr{
	pos: position{line: 543, col: 26, offset: 16189},
	run: (*parser).callonAndExpression1,
	expr: &seqExpr{
	pos: position{line: 543, col: 26, offset: 16189},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 543, col: 26, offset: 16189},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 543, col: 32, offset: 16195},
	name: "CombineExpression",
},
},
&labeledExpr{
	pos: position{line: 543, col: 55, offset: 16218},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 543, col: 60, offset: 16223},
	expr: &seqExpr{
	pos: position{line: 543, col: 61, offset: 16224},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 543, col: 61, offset: 16224},
	name: "_",
},
&litMatcher{
	pos: position{line: 543, col: 63, offset: 16226},
	val: "&&",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 543, col: 68, offset: 16231},
	name: "_",
},
&labeledExpr{
	pos: position{line: 543, col: 70, offset: 16233},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 543, col: 72, offset: 16235},
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
	pos: position{line: 545, col: 1, offset: 16305},
	expr: &actionExpr{
	pos: position{line: 545, col: 26, offset: 16332},
	run: (*parser).callonCombineExpression1,
	expr: &seqExpr{
	pos: position{line: 545, col: 26, offset: 16332},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 545, col: 26, offset: 16332},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 545, col: 32, offset: 16338},
	name: "PreferExpression",
},
},
&labeledExpr{
	pos: position{line: 545, col: 55, offset: 16361},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 545, col: 60, offset: 16366},
	expr: &seqExpr{
	pos: position{line: 545, col: 61, offset: 16367},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 545, col: 61, offset: 16367},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 545, col: 63, offset: 16369},
	name: "Combine",
},
&ruleRefExpr{
	pos: position{line: 545, col: 71, offset: 16377},
	name: "_",
},
&labeledExpr{
	pos: position{line: 545, col: 73, offset: 16379},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 545, col: 75, offset: 16381},
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
	pos: position{line: 547, col: 1, offset: 16458},
	expr: &actionExpr{
	pos: position{line: 547, col: 26, offset: 16485},
	run: (*parser).callonPreferExpression1,
	expr: &seqExpr{
	pos: position{line: 547, col: 26, offset: 16485},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 547, col: 26, offset: 16485},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 547, col: 32, offset: 16491},
	name: "CombineTypesExpression",
},
},
&labeledExpr{
	pos: position{line: 547, col: 55, offset: 16514},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 547, col: 60, offset: 16519},
	expr: &seqExpr{
	pos: position{line: 547, col: 61, offset: 16520},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 547, col: 61, offset: 16520},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 547, col: 63, offset: 16522},
	name: "Prefer",
},
&ruleRefExpr{
	pos: position{line: 547, col: 70, offset: 16529},
	name: "_",
},
&labeledExpr{
	pos: position{line: 547, col: 72, offset: 16531},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 547, col: 74, offset: 16533},
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
	pos: position{line: 549, col: 1, offset: 16627},
	expr: &actionExpr{
	pos: position{line: 549, col: 26, offset: 16654},
	run: (*parser).callonCombineTypesExpression1,
	expr: &seqExpr{
	pos: position{line: 549, col: 26, offset: 16654},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 549, col: 26, offset: 16654},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 549, col: 32, offset: 16660},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 549, col: 55, offset: 16683},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 549, col: 60, offset: 16688},
	expr: &seqExpr{
	pos: position{line: 549, col: 61, offset: 16689},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 549, col: 61, offset: 16689},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 549, col: 63, offset: 16691},
	name: "CombineTypes",
},
&ruleRefExpr{
	pos: position{line: 549, col: 76, offset: 16704},
	name: "_",
},
&labeledExpr{
	pos: position{line: 549, col: 78, offset: 16706},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 549, col: 80, offset: 16708},
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
	pos: position{line: 551, col: 1, offset: 16788},
	expr: &actionExpr{
	pos: position{line: 551, col: 26, offset: 16815},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 551, col: 26, offset: 16815},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 551, col: 26, offset: 16815},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 551, col: 32, offset: 16821},
	name: "EqualExpression",
},
},
&labeledExpr{
	pos: position{line: 551, col: 55, offset: 16844},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 551, col: 60, offset: 16849},
	expr: &seqExpr{
	pos: position{line: 551, col: 61, offset: 16850},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 551, col: 61, offset: 16850},
	name: "_",
},
&litMatcher{
	pos: position{line: 551, col: 63, offset: 16852},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 551, col: 67, offset: 16856},
	name: "_",
},
&labeledExpr{
	pos: position{line: 551, col: 69, offset: 16858},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 551, col: 71, offset: 16860},
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
	pos: position{line: 553, col: 1, offset: 16930},
	expr: &actionExpr{
	pos: position{line: 553, col: 26, offset: 16957},
	run: (*parser).callonEqualExpression1,
	expr: &seqExpr{
	pos: position{line: 553, col: 26, offset: 16957},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 553, col: 26, offset: 16957},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 553, col: 32, offset: 16963},
	name: "NotEqualExpression",
},
},
&labeledExpr{
	pos: position{line: 553, col: 55, offset: 16986},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 553, col: 60, offset: 16991},
	expr: &seqExpr{
	pos: position{line: 553, col: 61, offset: 16992},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 553, col: 61, offset: 16992},
	name: "_",
},
&litMatcher{
	pos: position{line: 553, col: 63, offset: 16994},
	val: "==",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 553, col: 68, offset: 16999},
	name: "_",
},
&labeledExpr{
	pos: position{line: 553, col: 70, offset: 17001},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 553, col: 72, offset: 17003},
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
	pos: position{line: 555, col: 1, offset: 17073},
	expr: &actionExpr{
	pos: position{line: 555, col: 26, offset: 17100},
	run: (*parser).callonNotEqualExpression1,
	expr: &seqExpr{
	pos: position{line: 555, col: 26, offset: 17100},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 555, col: 26, offset: 17100},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 555, col: 32, offset: 17106},
	name: "EquivalentExpression",
},
},
&labeledExpr{
	pos: position{line: 555, col: 54, offset: 17128},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 555, col: 59, offset: 17133},
	expr: &seqExpr{
	pos: position{line: 555, col: 60, offset: 17134},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 555, col: 60, offset: 17134},
	name: "_",
},
&litMatcher{
	pos: position{line: 555, col: 62, offset: 17136},
	val: "!=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 555, col: 67, offset: 17141},
	name: "_",
},
&labeledExpr{
	pos: position{line: 555, col: 69, offset: 17143},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 555, col: 71, offset: 17145},
	name: "EquivalentExpression",
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
	name: "EquivalentExpression",
	pos: position{line: 557, col: 1, offset: 17217},
	expr: &actionExpr{
	pos: position{line: 557, col: 28, offset: 17246},
	run: (*parser).callonEquivalentExpression1,
	expr: &seqExpr{
	pos: position{line: 557, col: 28, offset: 17246},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 557, col: 28, offset: 17246},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 557, col: 34, offset: 17252},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 557, col: 57, offset: 17275},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 557, col: 62, offset: 17280},
	expr: &seqExpr{
	pos: position{line: 557, col: 63, offset: 17281},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 557, col: 63, offset: 17281},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 557, col: 65, offset: 17283},
	name: "Equivalent",
},
&ruleRefExpr{
	pos: position{line: 557, col: 76, offset: 17294},
	name: "_",
},
&labeledExpr{
	pos: position{line: 557, col: 78, offset: 17296},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 557, col: 80, offset: 17298},
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
	pos: position{line: 560, col: 1, offset: 17375},
	expr: &actionExpr{
	pos: position{line: 560, col: 25, offset: 17401},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 560, col: 25, offset: 17401},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 560, col: 25, offset: 17401},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 560, col: 27, offset: 17403},
	name: "FirstApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 560, col: 54, offset: 17430},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 560, col: 59, offset: 17435},
	expr: &seqExpr{
	pos: position{line: 560, col: 60, offset: 17436},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 560, col: 60, offset: 17436},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 560, col: 63, offset: 17439},
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
	pos: position{line: 569, col: 1, offset: 17682},
	expr: &choiceExpr{
	pos: position{line: 570, col: 8, offset: 17720},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 570, col: 8, offset: 17720},
	run: (*parser).callonFirstApplicationExpression2,
	expr: &seqExpr{
	pos: position{line: 570, col: 8, offset: 17720},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 570, col: 8, offset: 17720},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 570, col: 14, offset: 17726},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 570, col: 17, offset: 17729},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 570, col: 19, offset: 17731},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 570, col: 36, offset: 17748},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 570, col: 39, offset: 17751},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 570, col: 41, offset: 17753},
	name: "ImportExpression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 573, col: 8, offset: 17856},
	run: (*parser).callonFirstApplicationExpression11,
	expr: &seqExpr{
	pos: position{line: 573, col: 8, offset: 17856},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 573, col: 8, offset: 17856},
	name: "Some",
},
&ruleRefExpr{
	pos: position{line: 573, col: 13, offset: 17861},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 573, col: 16, offset: 17864},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 573, col: 18, offset: 17866},
	name: "ImportExpression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 574, col: 8, offset: 17921},
	run: (*parser).callonFirstApplicationExpression17,
	expr: &seqExpr{
	pos: position{line: 574, col: 8, offset: 17921},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 574, col: 8, offset: 17921},
	name: "toMap",
},
&ruleRefExpr{
	pos: position{line: 574, col: 14, offset: 17927},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 574, col: 17, offset: 17930},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 574, col: 19, offset: 17932},
	name: "ImportExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 575, col: 8, offset: 17996},
	name: "ImportExpression",
},
	},
},
},
{
	name: "ImportExpression",
	pos: position{line: 577, col: 1, offset: 18014},
	expr: &choiceExpr{
	pos: position{line: 577, col: 20, offset: 18035},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 577, col: 20, offset: 18035},
	name: "Import",
},
&ruleRefExpr{
	pos: position{line: 577, col: 29, offset: 18044},
	name: "SelectorExpression",
},
	},
},
},
{
	name: "SelectorExpression",
	pos: position{line: 579, col: 1, offset: 18064},
	expr: &actionExpr{
	pos: position{line: 579, col: 22, offset: 18087},
	run: (*parser).callonSelectorExpression1,
	expr: &seqExpr{
	pos: position{line: 579, col: 22, offset: 18087},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 579, col: 22, offset: 18087},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 579, col: 24, offset: 18089},
	name: "PrimitiveExpression",
},
},
&labeledExpr{
	pos: position{line: 579, col: 44, offset: 18109},
	label: "ls",
	expr: &zeroOrMoreExpr{
	pos: position{line: 579, col: 47, offset: 18112},
	expr: &seqExpr{
	pos: position{line: 579, col: 48, offset: 18113},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 579, col: 48, offset: 18113},
	name: "_",
},
&litMatcher{
	pos: position{line: 579, col: 50, offset: 18115},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 579, col: 54, offset: 18119},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 579, col: 56, offset: 18121},
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
	pos: position{line: 598, col: 1, offset: 18674},
	expr: &choiceExpr{
	pos: position{line: 598, col: 12, offset: 18687},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 598, col: 12, offset: 18687},
	name: "AnyLabel",
},
&ruleRefExpr{
	pos: position{line: 598, col: 23, offset: 18698},
	name: "Labels",
},
&ruleRefExpr{
	pos: position{line: 598, col: 32, offset: 18707},
	name: "TypeSelector",
},
	},
},
},
{
	name: "Labels",
	pos: position{line: 600, col: 1, offset: 18721},
	expr: &actionExpr{
	pos: position{line: 600, col: 10, offset: 18732},
	run: (*parser).callonLabels1,
	expr: &seqExpr{
	pos: position{line: 600, col: 10, offset: 18732},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 600, col: 10, offset: 18732},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 600, col: 14, offset: 18736},
	name: "_",
},
&labeledExpr{
	pos: position{line: 600, col: 16, offset: 18738},
	label: "optclauses",
	expr: &zeroOrOneExpr{
	pos: position{line: 600, col: 27, offset: 18749},
	expr: &seqExpr{
	pos: position{line: 600, col: 29, offset: 18751},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 600, col: 29, offset: 18751},
	name: "AnyLabel",
},
&ruleRefExpr{
	pos: position{line: 600, col: 38, offset: 18760},
	name: "_",
},
&zeroOrMoreExpr{
	pos: position{line: 600, col: 40, offset: 18762},
	expr: &seqExpr{
	pos: position{line: 600, col: 41, offset: 18763},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 600, col: 41, offset: 18763},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 600, col: 45, offset: 18767},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 600, col: 47, offset: 18769},
	name: "AnyLabel",
},
&ruleRefExpr{
	pos: position{line: 600, col: 56, offset: 18778},
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
	pos: position{line: 600, col: 64, offset: 18786},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TypeSelector",
	pos: position{line: 610, col: 1, offset: 19082},
	expr: &actionExpr{
	pos: position{line: 610, col: 16, offset: 19099},
	run: (*parser).callonTypeSelector1,
	expr: &seqExpr{
	pos: position{line: 610, col: 16, offset: 19099},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 610, col: 16, offset: 19099},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 610, col: 20, offset: 19103},
	name: "_",
},
&labeledExpr{
	pos: position{line: 610, col: 22, offset: 19105},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 610, col: 24, offset: 19107},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 610, col: 35, offset: 19118},
	name: "_",
},
&litMatcher{
	pos: position{line: 610, col: 37, offset: 19120},
	val: ")",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PrimitiveExpression",
	pos: position{line: 612, col: 1, offset: 19143},
	expr: &choiceExpr{
	pos: position{line: 613, col: 7, offset: 19173},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 613, col: 7, offset: 19173},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 614, col: 7, offset: 19193},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 615, col: 7, offset: 19214},
	name: "IntegerLiteral",
},
&ruleRefExpr{
	pos: position{line: 616, col: 7, offset: 19235},
	name: "TextLiteral",
},
&actionExpr{
	pos: position{line: 617, col: 7, offset: 19253},
	run: (*parser).callonPrimitiveExpression6,
	expr: &seqExpr{
	pos: position{line: 617, col: 7, offset: 19253},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 617, col: 7, offset: 19253},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 617, col: 11, offset: 19257},
	name: "_",
},
&labeledExpr{
	pos: position{line: 617, col: 13, offset: 19259},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 617, col: 15, offset: 19261},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 617, col: 35, offset: 19281},
	name: "_",
},
&litMatcher{
	pos: position{line: 617, col: 37, offset: 19283},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&actionExpr{
	pos: position{line: 618, col: 7, offset: 19311},
	run: (*parser).callonPrimitiveExpression14,
	expr: &seqExpr{
	pos: position{line: 618, col: 7, offset: 19311},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 618, col: 7, offset: 19311},
	val: "<",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 618, col: 11, offset: 19315},
	name: "_",
},
&labeledExpr{
	pos: position{line: 618, col: 13, offset: 19317},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 618, col: 15, offset: 19319},
	name: "UnionType",
},
},
&ruleRefExpr{
	pos: position{line: 618, col: 25, offset: 19329},
	name: "_",
},
&litMatcher{
	pos: position{line: 618, col: 27, offset: 19331},
	val: ">",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 619, col: 7, offset: 19359},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 620, col: 7, offset: 19385},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 621, col: 7, offset: 19402},
	run: (*parser).callonPrimitiveExpression24,
	expr: &seqExpr{
	pos: position{line: 621, col: 7, offset: 19402},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 621, col: 7, offset: 19402},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 621, col: 11, offset: 19406},
	name: "_",
},
&labeledExpr{
	pos: position{line: 621, col: 14, offset: 19409},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 621, col: 16, offset: 19411},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 621, col: 27, offset: 19422},
	name: "_",
},
&litMatcher{
	pos: position{line: 621, col: 29, offset: 19424},
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
	pos: position{line: 623, col: 1, offset: 19447},
	expr: &choiceExpr{
	pos: position{line: 624, col: 7, offset: 19477},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 624, col: 7, offset: 19477},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &litMatcher{
	pos: position{line: 624, col: 7, offset: 19477},
	val: "=",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 625, col: 7, offset: 19515},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 626, col: 7, offset: 19540},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 627, col: 7, offset: 19568},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 627, col: 7, offset: 19568},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 629, col: 1, offset: 19597},
	expr: &actionExpr{
	pos: position{line: 629, col: 19, offset: 19617},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 629, col: 19, offset: 19617},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 629, col: 19, offset: 19617},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 629, col: 24, offset: 19622},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 629, col: 33, offset: 19631},
	name: "_",
},
&litMatcher{
	pos: position{line: 629, col: 35, offset: 19633},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 629, col: 39, offset: 19637},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 629, col: 42, offset: 19640},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 629, col: 47, offset: 19645},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 632, col: 1, offset: 19702},
	expr: &actionExpr{
	pos: position{line: 632, col: 18, offset: 19721},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 632, col: 18, offset: 19721},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 632, col: 18, offset: 19721},
	name: "_",
},
&litMatcher{
	pos: position{line: 632, col: 20, offset: 19723},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 632, col: 24, offset: 19727},
	name: "_",
},
&labeledExpr{
	pos: position{line: 632, col: 26, offset: 19729},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 632, col: 28, offset: 19731},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 633, col: 1, offset: 19763},
	expr: &actionExpr{
	pos: position{line: 634, col: 7, offset: 19792},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 634, col: 7, offset: 19792},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 634, col: 7, offset: 19792},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 634, col: 13, offset: 19798},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 634, col: 29, offset: 19814},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 634, col: 34, offset: 19819},
	expr: &ruleRefExpr{
	pos: position{line: 634, col: 34, offset: 19819},
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
	pos: position{line: 648, col: 1, offset: 20386},
	expr: &actionExpr{
	pos: position{line: 648, col: 22, offset: 20409},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 648, col: 22, offset: 20409},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 648, col: 22, offset: 20409},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 648, col: 27, offset: 20414},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 648, col: 36, offset: 20423},
	name: "_",
},
&litMatcher{
	pos: position{line: 648, col: 38, offset: 20425},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 648, col: 42, offset: 20429},
	name: "_",
},
&labeledExpr{
	pos: position{line: 648, col: 44, offset: 20431},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 648, col: 49, offset: 20436},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 651, col: 1, offset: 20493},
	expr: &actionExpr{
	pos: position{line: 651, col: 21, offset: 20515},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 651, col: 21, offset: 20515},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 651, col: 21, offset: 20515},
	name: "_",
},
&litMatcher{
	pos: position{line: 651, col: 23, offset: 20517},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 651, col: 27, offset: 20521},
	name: "_",
},
&labeledExpr{
	pos: position{line: 651, col: 29, offset: 20523},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 651, col: 31, offset: 20525},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 652, col: 1, offset: 20560},
	expr: &actionExpr{
	pos: position{line: 653, col: 7, offset: 20592},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 653, col: 7, offset: 20592},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 653, col: 7, offset: 20592},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 653, col: 13, offset: 20598},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 653, col: 32, offset: 20617},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 653, col: 37, offset: 20622},
	expr: &ruleRefExpr{
	pos: position{line: 653, col: 37, offset: 20622},
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
	pos: position{line: 667, col: 1, offset: 21195},
	expr: &choiceExpr{
	pos: position{line: 667, col: 13, offset: 21209},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 667, col: 13, offset: 21209},
	name: "NonEmptyUnionType",
},
&ruleRefExpr{
	pos: position{line: 667, col: 33, offset: 21229},
	name: "EmptyUnionType",
},
	},
},
},
{
	name: "EmptyUnionType",
	pos: position{line: 669, col: 1, offset: 21245},
	expr: &actionExpr{
	pos: position{line: 669, col: 18, offset: 21264},
	run: (*parser).callonEmptyUnionType1,
	expr: &litMatcher{
	pos: position{line: 669, col: 18, offset: 21264},
	val: "",
	ignoreCase: false,
},
},
},
{
	name: "NonEmptyUnionType",
	pos: position{line: 671, col: 1, offset: 21296},
	expr: &actionExpr{
	pos: position{line: 671, col: 21, offset: 21318},
	run: (*parser).callonNonEmptyUnionType1,
	expr: &seqExpr{
	pos: position{line: 671, col: 21, offset: 21318},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 671, col: 21, offset: 21318},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 671, col: 27, offset: 21324},
	name: "UnionVariant",
},
},
&labeledExpr{
	pos: position{line: 671, col: 40, offset: 21337},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 671, col: 45, offset: 21342},
	expr: &seqExpr{
	pos: position{line: 671, col: 46, offset: 21343},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 671, col: 46, offset: 21343},
	name: "_",
},
&litMatcher{
	pos: position{line: 671, col: 48, offset: 21345},
	val: "|",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 671, col: 52, offset: 21349},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 671, col: 54, offset: 21351},
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
	pos: position{line: 696, col: 1, offset: 22192},
	expr: &seqExpr{
	pos: position{line: 696, col: 16, offset: 22209},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 696, col: 16, offset: 22209},
	name: "AnyLabel",
},
&zeroOrOneExpr{
	pos: position{line: 696, col: 25, offset: 22218},
	expr: &seqExpr{
	pos: position{line: 696, col: 26, offset: 22219},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 696, col: 26, offset: 22219},
	name: "_",
},
&litMatcher{
	pos: position{line: 696, col: 28, offset: 22221},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 696, col: 32, offset: 22225},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 696, col: 35, offset: 22228},
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
	pos: position{line: 698, col: 1, offset: 22242},
	expr: &actionExpr{
	pos: position{line: 698, col: 12, offset: 22255},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 698, col: 12, offset: 22255},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 698, col: 12, offset: 22255},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 698, col: 16, offset: 22259},
	name: "_",
},
&labeledExpr{
	pos: position{line: 698, col: 18, offset: 22261},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 698, col: 20, offset: 22263},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 698, col: 31, offset: 22274},
	name: "_",
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 700, col: 1, offset: 22293},
	expr: &actionExpr{
	pos: position{line: 701, col: 7, offset: 22323},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 701, col: 7, offset: 22323},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 701, col: 7, offset: 22323},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 701, col: 11, offset: 22327},
	name: "_",
},
&labeledExpr{
	pos: position{line: 701, col: 13, offset: 22329},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 701, col: 19, offset: 22335},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 701, col: 30, offset: 22346},
	name: "_",
},
&labeledExpr{
	pos: position{line: 701, col: 32, offset: 22348},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 701, col: 37, offset: 22353},
	expr: &ruleRefExpr{
	pos: position{line: 701, col: 37, offset: 22353},
	name: "MoreList",
},
},
},
&litMatcher{
	pos: position{line: 701, col: 47, offset: 22363},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 711, col: 1, offset: 22631},
	expr: &notExpr{
	pos: position{line: 711, col: 7, offset: 22639},
	expr: &anyMatcher{
	line: 711, col: 8, offset: 22640,
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
 return NaturalSubtract, nil 
}

func (p *parser) callonReserved16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved16()
}

func (c *current) onReserved18() (interface{}, error) {
 return IntegerToDouble, nil 
}

func (p *parser) callonReserved18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved18()
}

func (c *current) onReserved20() (interface{}, error) {
 return IntegerShow, nil 
}

func (p *parser) callonReserved20() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved20()
}

func (c *current) onReserved22() (interface{}, error) {
 return DoubleShow, nil 
}

func (p *parser) callonReserved22() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved22()
}

func (c *current) onReserved24() (interface{}, error) {
 return ListBuild, nil 
}

func (p *parser) callonReserved24() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved24()
}

func (c *current) onReserved26() (interface{}, error) {
 return ListFold, nil 
}

func (p *parser) callonReserved26() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved26()
}

func (c *current) onReserved28() (interface{}, error) {
 return ListLength, nil 
}

func (p *parser) callonReserved28() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved28()
}

func (c *current) onReserved30() (interface{}, error) {
 return ListHead, nil 
}

func (p *parser) callonReserved30() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved30()
}

func (c *current) onReserved32() (interface{}, error) {
 return ListLast, nil 
}

func (p *parser) callonReserved32() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved32()
}

func (c *current) onReserved34() (interface{}, error) {
 return ListIndexed, nil 
}

func (p *parser) callonReserved34() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved34()
}

func (c *current) onReserved36() (interface{}, error) {
 return ListReverse, nil 
}

func (p *parser) callonReserved36() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved36()
}

func (c *current) onReserved38() (interface{}, error) {
 return OptionalBuild, nil 
}

func (p *parser) callonReserved38() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved38()
}

func (c *current) onReserved40() (interface{}, error) {
 return OptionalFold, nil 
}

func (p *parser) callonReserved40() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved40()
}

func (c *current) onReserved42() (interface{}, error) {
 return TextShow, nil 
}

func (p *parser) callonReserved42() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved42()
}

func (c *current) onReserved44() (interface{}, error) {
 return Bool, nil 
}

func (p *parser) callonReserved44() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved44()
}

func (c *current) onReserved46() (interface{}, error) {
 return True, nil 
}

func (p *parser) callonReserved46() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved46()
}

func (c *current) onReserved48() (interface{}, error) {
 return False, nil 
}

func (p *parser) callonReserved48() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved48()
}

func (c *current) onReserved50() (interface{}, error) {
 return Optional, nil 
}

func (p *parser) callonReserved50() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved50()
}

func (c *current) onReserved52() (interface{}, error) {
 return Natural, nil 
}

func (p *parser) callonReserved52() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved52()
}

func (c *current) onReserved54() (interface{}, error) {
 return Integer, nil 
}

func (p *parser) callonReserved54() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved54()
}

func (c *current) onReserved56() (interface{}, error) {
 return Double, nil 
}

func (p *parser) callonReserved56() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved56()
}

func (c *current) onReserved58() (interface{}, error) {
 return Text, nil 
}

func (p *parser) callonReserved58() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved58()
}

func (c *current) onReserved60() (interface{}, error) {
 return List, nil 
}

func (p *parser) callonReserved60() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved60()
}

func (c *current) onReserved62() (interface{}, error) {
 return None, nil 
}

func (p *parser) callonReserved62() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved62()
}

func (c *current) onReserved64() (interface{}, error) {
 return Type, nil 
}

func (p *parser) callonReserved64() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved64()
}

func (c *current) onReserved66() (interface{}, error) {
 return Kind, nil 
}

func (p *parser) callonReserved66() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved66()
}

func (c *current) onReserved68() (interface{}, error) {
 return Sort, nil 
}

func (p *parser) callonReserved68() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved68()
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
 return Embed(Import{ImportHashed: i.(ImportHashed), ImportMode: RawText}), nil 
}

func (p *parser) callonImport2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onImport2(stack["i"])
}

func (c *current) onImport10(i interface{}) (interface{}, error) {
 return Embed(Import{ImportHashed: i.(ImportHashed), ImportMode: Location}), nil 
}

func (p *parser) callonImport10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onImport10(stack["i"])
}

func (c *current) onImport18(i interface{}) (interface{}, error) {
 return Embed(Import{ImportHashed: i.(ImportHashed), ImportMode: Code}), nil 
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

func (c *current) onExpression102(a interface{}) (interface{}, error) {
 return Assert{Annotation: a.(Expr)}, nil 
}

func (p *parser) callonExpression102() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression102(stack["a"])
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

func (c *current) onEquivalentExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(EquivOp, first, rest), nil
}

func (p *parser) callonEquivalentExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEquivalentExpression1(stack["first"], stack["rest"])
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
        name := alternative[0].(string)
        if _, ok := alternatives[name]; ok {
            return nil, fmt.Errorf("Duplicate alternative %s in union", name)
        }

        if alternative[1] == nil {
            alternatives[name] = nil
        } else {
            alternatives[name] = alternative[1].([]interface{})[3].(Expr)
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

