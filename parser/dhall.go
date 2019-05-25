
package parser

import (
"bytes"
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


var g = &grammar {
	rules: []*rule{
{
	name: "DhallFile",
	pos: position{line: 36, col: 1, offset: 621},
	expr: &actionExpr{
	pos: position{line: 36, col: 13, offset: 635},
	run: (*parser).callonDhallFile1,
	expr: &seqExpr{
	pos: position{line: 36, col: 13, offset: 635},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 36, col: 13, offset: 635},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 36, col: 15, offset: 637},
	name: "CompleteExpression",
},
},
&ruleRefExpr{
	pos: position{line: 36, col: 34, offset: 656},
	name: "EOF",
},
	},
},
},
},
{
	name: "CompleteExpression",
	pos: position{line: 38, col: 1, offset: 679},
	expr: &actionExpr{
	pos: position{line: 38, col: 22, offset: 702},
	run: (*parser).callonCompleteExpression1,
	expr: &seqExpr{
	pos: position{line: 38, col: 22, offset: 702},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 38, col: 22, offset: 702},
	name: "_",
},
&labeledExpr{
	pos: position{line: 38, col: 24, offset: 704},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 38, col: 26, offset: 706},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 38, col: 37, offset: 717},
	name: "_",
},
	},
},
},
},
{
	name: "EOL",
	pos: position{line: 40, col: 1, offset: 738},
	expr: &choiceExpr{
	pos: position{line: 40, col: 7, offset: 746},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 40, col: 7, offset: 746},
	val: "\n",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 40, col: 14, offset: 753},
	run: (*parser).callonEOL3,
	expr: &litMatcher{
	pos: position{line: 40, col: 14, offset: 753},
	val: "\r\n",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "BlockComment",
	pos: position{line: 42, col: 1, offset: 790},
	expr: &seqExpr{
	pos: position{line: 42, col: 16, offset: 807},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 42, col: 16, offset: 807},
	val: "{-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 42, col: 21, offset: 812},
	name: "BlockCommentContinue",
},
	},
},
},
{
	name: "BlockCommentChunk",
	pos: position{line: 44, col: 1, offset: 834},
	expr: &choiceExpr{
	pos: position{line: 45, col: 5, offset: 860},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 45, col: 5, offset: 860},
	name: "BlockComment",
},
&charClassMatcher{
	pos: position{line: 46, col: 5, offset: 877},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 47, col: 5, offset: 903},
	name: "EOL",
},
	},
},
},
{
	name: "BlockCommentContinue",
	pos: position{line: 49, col: 1, offset: 908},
	expr: &choiceExpr{
	pos: position{line: 49, col: 24, offset: 933},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 49, col: 24, offset: 933},
	val: "-}",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 49, col: 31, offset: 940},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 49, col: 31, offset: 940},
	name: "BlockCommentChunk",
},
&ruleRefExpr{
	pos: position{line: 49, col: 49, offset: 958},
	name: "BlockCommentContinue",
},
	},
},
	},
},
},
{
	name: "NotEOL",
	pos: position{line: 51, col: 1, offset: 980},
	expr: &charClassMatcher{
	pos: position{line: 51, col: 10, offset: 991},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "LineComment",
	pos: position{line: 53, col: 1, offset: 1014},
	expr: &actionExpr{
	pos: position{line: 53, col: 15, offset: 1030},
	run: (*parser).callonLineComment1,
	expr: &seqExpr{
	pos: position{line: 53, col: 15, offset: 1030},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 53, col: 15, offset: 1030},
	val: "--",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 53, col: 20, offset: 1035},
	label: "content",
	expr: &actionExpr{
	pos: position{line: 53, col: 29, offset: 1044},
	run: (*parser).callonLineComment5,
	expr: &zeroOrMoreExpr{
	pos: position{line: 53, col: 29, offset: 1044},
	expr: &ruleRefExpr{
	pos: position{line: 53, col: 29, offset: 1044},
	name: "NotEOL",
},
},
},
},
&ruleRefExpr{
	pos: position{line: 53, col: 68, offset: 1083},
	name: "EOL",
},
	},
},
},
},
{
	name: "WhitespaceChunk",
	pos: position{line: 55, col: 1, offset: 1112},
	expr: &choiceExpr{
	pos: position{line: 55, col: 19, offset: 1132},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 55, col: 19, offset: 1132},
	val: " ",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 55, col: 25, offset: 1138},
	val: "\t",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 55, col: 32, offset: 1145},
	name: "EOL",
},
&ruleRefExpr{
	pos: position{line: 55, col: 38, offset: 1151},
	name: "LineComment",
},
&ruleRefExpr{
	pos: position{line: 55, col: 52, offset: 1165},
	name: "BlockComment",
},
	},
},
},
{
	name: "_",
	pos: position{line: 57, col: 1, offset: 1179},
	expr: &zeroOrMoreExpr{
	pos: position{line: 57, col: 5, offset: 1185},
	expr: &ruleRefExpr{
	pos: position{line: 57, col: 5, offset: 1185},
	name: "WhitespaceChunk",
},
},
},
{
	name: "_1",
	pos: position{line: 59, col: 1, offset: 1203},
	expr: &oneOrMoreExpr{
	pos: position{line: 59, col: 6, offset: 1210},
	expr: &ruleRefExpr{
	pos: position{line: 59, col: 6, offset: 1210},
	name: "WhitespaceChunk",
},
},
},
{
	name: "Digit",
	pos: position{line: 61, col: 1, offset: 1228},
	expr: &charClassMatcher{
	pos: position{line: 61, col: 9, offset: 1238},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "HexDig",
	pos: position{line: 63, col: 1, offset: 1245},
	expr: &choiceExpr{
	pos: position{line: 63, col: 10, offset: 1256},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 63, col: 10, offset: 1256},
	name: "Digit",
},
&charClassMatcher{
	pos: position{line: 63, col: 18, offset: 1264},
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
	pos: position{line: 65, col: 1, offset: 1272},
	expr: &charClassMatcher{
	pos: position{line: 65, col: 24, offset: 1297},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SimpleLabelNextChar",
	pos: position{line: 66, col: 1, offset: 1307},
	expr: &charClassMatcher{
	pos: position{line: 66, col: 23, offset: 1331},
	val: "[A-Za-z0-9_/-]",
	chars: []rune{'_','/','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SimpleLabel",
	pos: position{line: 67, col: 1, offset: 1346},
	expr: &choiceExpr{
	pos: position{line: 67, col: 15, offset: 1362},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 67, col: 15, offset: 1362},
	run: (*parser).callonSimpleLabel2,
	expr: &seqExpr{
	pos: position{line: 67, col: 15, offset: 1362},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 67, col: 15, offset: 1362},
	name: "Keyword",
},
&oneOrMoreExpr{
	pos: position{line: 67, col: 23, offset: 1370},
	expr: &ruleRefExpr{
	pos: position{line: 67, col: 23, offset: 1370},
	name: "SimpleLabelNextChar",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 68, col: 13, offset: 1434},
	run: (*parser).callonSimpleLabel7,
	expr: &seqExpr{
	pos: position{line: 68, col: 13, offset: 1434},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 68, col: 13, offset: 1434},
	expr: &ruleRefExpr{
	pos: position{line: 68, col: 14, offset: 1435},
	name: "Keyword",
},
},
&ruleRefExpr{
	pos: position{line: 68, col: 22, offset: 1443},
	name: "SimpleLabelFirstChar",
},
&zeroOrMoreExpr{
	pos: position{line: 68, col: 43, offset: 1464},
	expr: &ruleRefExpr{
	pos: position{line: 68, col: 43, offset: 1464},
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
	pos: position{line: 73, col: 1, offset: 1549},
	expr: &charClassMatcher{
	pos: position{line: 73, col: 19, offset: 1569},
	val: "[\\x20-\\x5f\\x61-\\x7e]",
	ranges: []rune{' ','_','a','~',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "QuotedLabel",
	pos: position{line: 74, col: 1, offset: 1590},
	expr: &actionExpr{
	pos: position{line: 74, col: 15, offset: 1606},
	run: (*parser).callonQuotedLabel1,
	expr: &oneOrMoreExpr{
	pos: position{line: 74, col: 15, offset: 1606},
	expr: &ruleRefExpr{
	pos: position{line: 74, col: 15, offset: 1606},
	name: "QuotedLabelChar",
},
},
},
},
{
	name: "Label",
	pos: position{line: 76, col: 1, offset: 1655},
	expr: &choiceExpr{
	pos: position{line: 76, col: 9, offset: 1665},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 76, col: 9, offset: 1665},
	run: (*parser).callonLabel2,
	expr: &seqExpr{
	pos: position{line: 76, col: 9, offset: 1665},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 76, col: 9, offset: 1665},
	val: "`",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 76, col: 13, offset: 1669},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 76, col: 19, offset: 1675},
	name: "QuotedLabel",
},
},
&litMatcher{
	pos: position{line: 76, col: 31, offset: 1687},
	val: "`",
	ignoreCase: false,
},
	},
},
},
&actionExpr{
	pos: position{line: 77, col: 9, offset: 1721},
	run: (*parser).callonLabel8,
	expr: &labeledExpr{
	pos: position{line: 77, col: 9, offset: 1721},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 77, col: 15, offset: 1727},
	name: "SimpleLabel",
},
},
},
	},
},
},
{
	name: "NonreservedLabel",
	pos: position{line: 79, col: 1, offset: 1762},
	expr: &choiceExpr{
	pos: position{line: 79, col: 20, offset: 1783},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 79, col: 20, offset: 1783},
	run: (*parser).callonNonreservedLabel2,
	expr: &seqExpr{
	pos: position{line: 79, col: 20, offset: 1783},
	exprs: []interface{}{
&andExpr{
	pos: position{line: 79, col: 20, offset: 1783},
	expr: &seqExpr{
	pos: position{line: 79, col: 22, offset: 1785},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 79, col: 22, offset: 1785},
	name: "Reserved",
},
&ruleRefExpr{
	pos: position{line: 79, col: 31, offset: 1794},
	name: "SimpleLabelNextChar",
},
	},
},
},
&labeledExpr{
	pos: position{line: 79, col: 52, offset: 1815},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 79, col: 58, offset: 1821},
	name: "Label",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 80, col: 19, offset: 1867},
	run: (*parser).callonNonreservedLabel10,
	expr: &seqExpr{
	pos: position{line: 80, col: 19, offset: 1867},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 80, col: 19, offset: 1867},
	expr: &ruleRefExpr{
	pos: position{line: 80, col: 20, offset: 1868},
	name: "Reserved",
},
},
&labeledExpr{
	pos: position{line: 80, col: 29, offset: 1877},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 80, col: 35, offset: 1883},
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
	pos: position{line: 82, col: 1, offset: 1912},
	expr: &ruleRefExpr{
	pos: position{line: 82, col: 12, offset: 1925},
	name: "Label",
},
},
{
	name: "DoubleQuoteChunk",
	pos: position{line: 85, col: 1, offset: 1933},
	expr: &choiceExpr{
	pos: position{line: 86, col: 6, offset: 1959},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 86, col: 6, offset: 1959},
	name: "Interpolation",
},
&actionExpr{
	pos: position{line: 87, col: 6, offset: 1978},
	run: (*parser).callonDoubleQuoteChunk3,
	expr: &seqExpr{
	pos: position{line: 87, col: 6, offset: 1978},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 87, col: 6, offset: 1978},
	val: "\\",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 87, col: 11, offset: 1983},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 87, col: 13, offset: 1985},
	name: "DoubleQuoteEscaped",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 88, col: 6, offset: 2027},
	name: "DoubleQuoteChar",
},
	},
},
},
{
	name: "DoubleQuoteEscaped",
	pos: position{line: 90, col: 1, offset: 2044},
	expr: &choiceExpr{
	pos: position{line: 91, col: 8, offset: 2074},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 91, col: 8, offset: 2074},
	val: "\"",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 92, col: 8, offset: 2085},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 93, col: 8, offset: 2096},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 94, col: 8, offset: 2108},
	val: "/",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 95, col: 8, offset: 2119},
	run: (*parser).callonDoubleQuoteEscaped6,
	expr: &litMatcher{
	pos: position{line: 95, col: 8, offset: 2119},
	val: "b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 96, col: 8, offset: 2159},
	run: (*parser).callonDoubleQuoteEscaped8,
	expr: &litMatcher{
	pos: position{line: 96, col: 8, offset: 2159},
	val: "f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 97, col: 8, offset: 2199},
	run: (*parser).callonDoubleQuoteEscaped10,
	expr: &litMatcher{
	pos: position{line: 97, col: 8, offset: 2199},
	val: "n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 98, col: 8, offset: 2239},
	run: (*parser).callonDoubleQuoteEscaped12,
	expr: &litMatcher{
	pos: position{line: 98, col: 8, offset: 2239},
	val: "r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 99, col: 8, offset: 2279},
	run: (*parser).callonDoubleQuoteEscaped14,
	expr: &litMatcher{
	pos: position{line: 99, col: 8, offset: 2279},
	val: "t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 100, col: 8, offset: 2319},
	run: (*parser).callonDoubleQuoteEscaped16,
	expr: &seqExpr{
	pos: position{line: 100, col: 8, offset: 2319},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 100, col: 8, offset: 2319},
	val: "u",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 100, col: 12, offset: 2323},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 100, col: 19, offset: 2330},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 100, col: 26, offset: 2337},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 100, col: 33, offset: 2344},
	name: "HexDig",
},
	},
},
},
	},
},
},
{
	name: "DoubleQuoteChar",
	pos: position{line: 105, col: 1, offset: 2476},
	expr: &choiceExpr{
	pos: position{line: 106, col: 6, offset: 2501},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 106, col: 6, offset: 2501},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 107, col: 6, offset: 2518},
	val: "[\\x23-\\x5b]",
	ranges: []rune{'#','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 108, col: 6, offset: 2535},
	val: "[\\x5d-\\U0010ffff]",
	ranges: []rune{']','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
	},
},
},
{
	name: "DoubleQuoteLiteral",
	pos: position{line: 110, col: 1, offset: 2554},
	expr: &actionExpr{
	pos: position{line: 110, col: 22, offset: 2577},
	run: (*parser).callonDoubleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 110, col: 22, offset: 2577},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 110, col: 22, offset: 2577},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 110, col: 26, offset: 2581},
	label: "chunks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 110, col: 33, offset: 2588},
	expr: &ruleRefExpr{
	pos: position{line: 110, col: 33, offset: 2588},
	name: "DoubleQuoteChunk",
},
},
},
&litMatcher{
	pos: position{line: 110, col: 51, offset: 2606},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "SingleQuoteContinue",
	pos: position{line: 127, col: 1, offset: 3074},
	expr: &choiceExpr{
	pos: position{line: 128, col: 7, offset: 3104},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 128, col: 7, offset: 3104},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 128, col: 7, offset: 3104},
	name: "Interpolation",
},
&ruleRefExpr{
	pos: position{line: 128, col: 21, offset: 3118},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 129, col: 7, offset: 3144},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 129, col: 7, offset: 3144},
	name: "EscapedQuotePair",
},
&ruleRefExpr{
	pos: position{line: 129, col: 24, offset: 3161},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 130, col: 7, offset: 3187},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 130, col: 7, offset: 3187},
	name: "EscapedInterpolation",
},
&ruleRefExpr{
	pos: position{line: 130, col: 28, offset: 3208},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 131, col: 7, offset: 3234},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 131, col: 7, offset: 3234},
	name: "SingleQuoteChar",
},
&ruleRefExpr{
	pos: position{line: 131, col: 23, offset: 3250},
	name: "SingleQuoteContinue",
},
	},
},
&litMatcher{
	pos: position{line: 132, col: 7, offset: 3276},
	val: "''",
	ignoreCase: false,
},
	},
},
},
{
	name: "EscapedQuotePair",
	pos: position{line: 134, col: 1, offset: 3282},
	expr: &actionExpr{
	pos: position{line: 134, col: 20, offset: 3303},
	run: (*parser).callonEscapedQuotePair1,
	expr: &litMatcher{
	pos: position{line: 134, col: 20, offset: 3303},
	val: "'''",
	ignoreCase: false,
},
},
},
{
	name: "EscapedInterpolation",
	pos: position{line: 138, col: 1, offset: 3438},
	expr: &actionExpr{
	pos: position{line: 138, col: 24, offset: 3463},
	run: (*parser).callonEscapedInterpolation1,
	expr: &litMatcher{
	pos: position{line: 138, col: 24, offset: 3463},
	val: "''${",
	ignoreCase: false,
},
},
},
{
	name: "SingleQuoteChar",
	pos: position{line: 140, col: 1, offset: 3505},
	expr: &choiceExpr{
	pos: position{line: 141, col: 6, offset: 3530},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 141, col: 6, offset: 3530},
	val: "[\\x20-\\U0010ffff]",
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 142, col: 6, offset: 3553},
	val: "\t",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 143, col: 6, offset: 3563},
	name: "EOL",
},
	},
},
},
{
	name: "SingleQuoteLiteral",
	pos: position{line: 145, col: 1, offset: 3568},
	expr: &actionExpr{
	pos: position{line: 145, col: 22, offset: 3591},
	run: (*parser).callonSingleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 145, col: 22, offset: 3591},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 145, col: 22, offset: 3591},
	val: "''",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 145, col: 27, offset: 3596},
	name: "EOL",
},
&labeledExpr{
	pos: position{line: 145, col: 31, offset: 3600},
	label: "content",
	expr: &ruleRefExpr{
	pos: position{line: 145, col: 39, offset: 3608},
	name: "SingleQuoteContinue",
},
},
	},
},
},
},
{
	name: "Interpolation",
	pos: position{line: 163, col: 1, offset: 4158},
	expr: &actionExpr{
	pos: position{line: 163, col: 17, offset: 4176},
	run: (*parser).callonInterpolation1,
	expr: &seqExpr{
	pos: position{line: 163, col: 17, offset: 4176},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 163, col: 17, offset: 4176},
	val: "${",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 163, col: 22, offset: 4181},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 163, col: 24, offset: 4183},
	name: "CompleteExpression",
},
},
&litMatcher{
	pos: position{line: 163, col: 43, offset: 4202},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TextLiteral",
	pos: position{line: 165, col: 1, offset: 4225},
	expr: &choiceExpr{
	pos: position{line: 165, col: 15, offset: 4241},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 165, col: 15, offset: 4241},
	name: "DoubleQuoteLiteral",
},
&ruleRefExpr{
	pos: position{line: 165, col: 36, offset: 4262},
	name: "SingleQuoteLiteral",
},
	},
},
},
{
	name: "Reserved",
	pos: position{line: 168, col: 1, offset: 4367},
	expr: &choiceExpr{
	pos: position{line: 169, col: 5, offset: 4384},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 169, col: 5, offset: 4384},
	run: (*parser).callonReserved2,
	expr: &litMatcher{
	pos: position{line: 169, col: 5, offset: 4384},
	val: "Natural/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 170, col: 5, offset: 4433},
	run: (*parser).callonReserved4,
	expr: &litMatcher{
	pos: position{line: 170, col: 5, offset: 4433},
	val: "Natural/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 171, col: 5, offset: 4480},
	run: (*parser).callonReserved6,
	expr: &litMatcher{
	pos: position{line: 171, col: 5, offset: 4480},
	val: "Natural/isZero",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 172, col: 5, offset: 4531},
	run: (*parser).callonReserved8,
	expr: &litMatcher{
	pos: position{line: 172, col: 5, offset: 4531},
	val: "Natural/even",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 173, col: 5, offset: 4578},
	run: (*parser).callonReserved10,
	expr: &litMatcher{
	pos: position{line: 173, col: 5, offset: 4578},
	val: "Natural/odd",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 174, col: 5, offset: 4623},
	run: (*parser).callonReserved12,
	expr: &litMatcher{
	pos: position{line: 174, col: 5, offset: 4623},
	val: "Natural/toInteger",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 175, col: 5, offset: 4680},
	run: (*parser).callonReserved14,
	expr: &litMatcher{
	pos: position{line: 175, col: 5, offset: 4680},
	val: "Natural/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 176, col: 5, offset: 4727},
	run: (*parser).callonReserved16,
	expr: &litMatcher{
	pos: position{line: 176, col: 5, offset: 4727},
	val: "Integer/toDouble",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 177, col: 5, offset: 4782},
	run: (*parser).callonReserved18,
	expr: &litMatcher{
	pos: position{line: 177, col: 5, offset: 4782},
	val: "Integer/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 178, col: 5, offset: 4829},
	run: (*parser).callonReserved20,
	expr: &litMatcher{
	pos: position{line: 178, col: 5, offset: 4829},
	val: "Double/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 179, col: 5, offset: 4874},
	run: (*parser).callonReserved22,
	expr: &litMatcher{
	pos: position{line: 179, col: 5, offset: 4874},
	val: "List/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 180, col: 5, offset: 4917},
	run: (*parser).callonReserved24,
	expr: &litMatcher{
	pos: position{line: 180, col: 5, offset: 4917},
	val: "List/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 181, col: 5, offset: 4958},
	run: (*parser).callonReserved26,
	expr: &litMatcher{
	pos: position{line: 181, col: 5, offset: 4958},
	val: "List/length",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 182, col: 5, offset: 5003},
	run: (*parser).callonReserved28,
	expr: &litMatcher{
	pos: position{line: 182, col: 5, offset: 5003},
	val: "List/head",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 183, col: 5, offset: 5044},
	run: (*parser).callonReserved30,
	expr: &litMatcher{
	pos: position{line: 183, col: 5, offset: 5044},
	val: "List/last",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 184, col: 5, offset: 5085},
	run: (*parser).callonReserved32,
	expr: &litMatcher{
	pos: position{line: 184, col: 5, offset: 5085},
	val: "List/indexed",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 185, col: 5, offset: 5132},
	run: (*parser).callonReserved34,
	expr: &litMatcher{
	pos: position{line: 185, col: 5, offset: 5132},
	val: "List/reverse",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 186, col: 5, offset: 5179},
	run: (*parser).callonReserved36,
	expr: &litMatcher{
	pos: position{line: 186, col: 5, offset: 5179},
	val: "Optional/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 187, col: 5, offset: 5230},
	run: (*parser).callonReserved38,
	expr: &litMatcher{
	pos: position{line: 187, col: 5, offset: 5230},
	val: "Optional/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 188, col: 5, offset: 5279},
	run: (*parser).callonReserved40,
	expr: &litMatcher{
	pos: position{line: 188, col: 5, offset: 5279},
	val: "Text/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 189, col: 5, offset: 5320},
	run: (*parser).callonReserved42,
	expr: &litMatcher{
	pos: position{line: 189, col: 5, offset: 5320},
	val: "Bool",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 190, col: 5, offset: 5352},
	run: (*parser).callonReserved44,
	expr: &litMatcher{
	pos: position{line: 190, col: 5, offset: 5352},
	val: "True",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 191, col: 5, offset: 5384},
	run: (*parser).callonReserved46,
	expr: &litMatcher{
	pos: position{line: 191, col: 5, offset: 5384},
	val: "False",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 192, col: 5, offset: 5418},
	run: (*parser).callonReserved48,
	expr: &litMatcher{
	pos: position{line: 192, col: 5, offset: 5418},
	val: "Optional",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 193, col: 5, offset: 5458},
	run: (*parser).callonReserved50,
	expr: &litMatcher{
	pos: position{line: 193, col: 5, offset: 5458},
	val: "Natural",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 194, col: 5, offset: 5496},
	run: (*parser).callonReserved52,
	expr: &litMatcher{
	pos: position{line: 194, col: 5, offset: 5496},
	val: "Integer",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 195, col: 5, offset: 5534},
	run: (*parser).callonReserved54,
	expr: &litMatcher{
	pos: position{line: 195, col: 5, offset: 5534},
	val: "Double",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 196, col: 5, offset: 5570},
	run: (*parser).callonReserved56,
	expr: &litMatcher{
	pos: position{line: 196, col: 5, offset: 5570},
	val: "Text",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 197, col: 5, offset: 5602},
	run: (*parser).callonReserved58,
	expr: &litMatcher{
	pos: position{line: 197, col: 5, offset: 5602},
	val: "List",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 198, col: 5, offset: 5634},
	run: (*parser).callonReserved60,
	expr: &litMatcher{
	pos: position{line: 198, col: 5, offset: 5634},
	val: "None",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 199, col: 5, offset: 5666},
	run: (*parser).callonReserved62,
	expr: &litMatcher{
	pos: position{line: 199, col: 5, offset: 5666},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 200, col: 5, offset: 5698},
	run: (*parser).callonReserved64,
	expr: &litMatcher{
	pos: position{line: 200, col: 5, offset: 5698},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 201, col: 5, offset: 5730},
	run: (*parser).callonReserved66,
	expr: &litMatcher{
	pos: position{line: 201, col: 5, offset: 5730},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "If",
	pos: position{line: 203, col: 1, offset: 5759},
	expr: &litMatcher{
	pos: position{line: 203, col: 6, offset: 5766},
	val: "if",
	ignoreCase: false,
},
},
{
	name: "Then",
	pos: position{line: 204, col: 1, offset: 5771},
	expr: &litMatcher{
	pos: position{line: 204, col: 8, offset: 5780},
	val: "then",
	ignoreCase: false,
},
},
{
	name: "Else",
	pos: position{line: 205, col: 1, offset: 5787},
	expr: &litMatcher{
	pos: position{line: 205, col: 8, offset: 5796},
	val: "else",
	ignoreCase: false,
},
},
{
	name: "Let",
	pos: position{line: 206, col: 1, offset: 5803},
	expr: &litMatcher{
	pos: position{line: 206, col: 7, offset: 5811},
	val: "let",
	ignoreCase: false,
},
},
{
	name: "In",
	pos: position{line: 207, col: 1, offset: 5817},
	expr: &litMatcher{
	pos: position{line: 207, col: 6, offset: 5824},
	val: "in",
	ignoreCase: false,
},
},
{
	name: "As",
	pos: position{line: 208, col: 1, offset: 5829},
	expr: &litMatcher{
	pos: position{line: 208, col: 6, offset: 5836},
	val: "as",
	ignoreCase: false,
},
},
{
	name: "Using",
	pos: position{line: 209, col: 1, offset: 5841},
	expr: &litMatcher{
	pos: position{line: 209, col: 9, offset: 5851},
	val: "using",
	ignoreCase: false,
},
},
{
	name: "Merge",
	pos: position{line: 210, col: 1, offset: 5859},
	expr: &litMatcher{
	pos: position{line: 210, col: 9, offset: 5869},
	val: "merge",
	ignoreCase: false,
},
},
{
	name: "Missing",
	pos: position{line: 211, col: 1, offset: 5877},
	expr: &actionExpr{
	pos: position{line: 211, col: 11, offset: 5889},
	run: (*parser).callonMissing1,
	expr: &litMatcher{
	pos: position{line: 211, col: 11, offset: 5889},
	val: "missing",
	ignoreCase: false,
},
},
},
{
	name: "True",
	pos: position{line: 212, col: 1, offset: 5925},
	expr: &litMatcher{
	pos: position{line: 212, col: 8, offset: 5934},
	val: "True",
	ignoreCase: false,
},
},
{
	name: "False",
	pos: position{line: 213, col: 1, offset: 5941},
	expr: &litMatcher{
	pos: position{line: 213, col: 9, offset: 5951},
	val: "False",
	ignoreCase: false,
},
},
{
	name: "Infinity",
	pos: position{line: 214, col: 1, offset: 5959},
	expr: &litMatcher{
	pos: position{line: 214, col: 12, offset: 5972},
	val: "Infinity",
	ignoreCase: false,
},
},
{
	name: "NaN",
	pos: position{line: 215, col: 1, offset: 5983},
	expr: &litMatcher{
	pos: position{line: 215, col: 7, offset: 5991},
	val: "NaN",
	ignoreCase: false,
},
},
{
	name: "Some",
	pos: position{line: 216, col: 1, offset: 5997},
	expr: &litMatcher{
	pos: position{line: 216, col: 8, offset: 6006},
	val: "Some",
	ignoreCase: false,
},
},
{
	name: "Keyword",
	pos: position{line: 218, col: 1, offset: 6014},
	expr: &choiceExpr{
	pos: position{line: 219, col: 5, offset: 6030},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 219, col: 5, offset: 6030},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 219, col: 10, offset: 6035},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 219, col: 17, offset: 6042},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 220, col: 5, offset: 6051},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 220, col: 11, offset: 6057},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 221, col: 5, offset: 6064},
	name: "Using",
},
&ruleRefExpr{
	pos: position{line: 221, col: 13, offset: 6072},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 221, col: 23, offset: 6082},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 222, col: 5, offset: 6089},
	name: "True",
},
&ruleRefExpr{
	pos: position{line: 222, col: 12, offset: 6096},
	name: "False",
},
&ruleRefExpr{
	pos: position{line: 223, col: 5, offset: 6106},
	name: "Infinity",
},
&ruleRefExpr{
	pos: position{line: 223, col: 16, offset: 6117},
	name: "NaN",
},
&ruleRefExpr{
	pos: position{line: 224, col: 5, offset: 6125},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 224, col: 13, offset: 6133},
	name: "Some",
},
	},
},
},
{
	name: "Optional",
	pos: position{line: 226, col: 1, offset: 6139},
	expr: &litMatcher{
	pos: position{line: 226, col: 12, offset: 6152},
	val: "Optional",
	ignoreCase: false,
},
},
{
	name: "Text",
	pos: position{line: 227, col: 1, offset: 6163},
	expr: &litMatcher{
	pos: position{line: 227, col: 8, offset: 6172},
	val: "Text",
	ignoreCase: false,
},
},
{
	name: "List",
	pos: position{line: 228, col: 1, offset: 6179},
	expr: &litMatcher{
	pos: position{line: 228, col: 8, offset: 6188},
	val: "List",
	ignoreCase: false,
},
},
{
	name: "Combine",
	pos: position{line: 230, col: 1, offset: 6196},
	expr: &choiceExpr{
	pos: position{line: 230, col: 11, offset: 6208},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 230, col: 11, offset: 6208},
	val: "/\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 230, col: 19, offset: 6216},
	val: "∧",
	ignoreCase: false,
},
	},
},
},
{
	name: "CombineTypes",
	pos: position{line: 231, col: 1, offset: 6222},
	expr: &choiceExpr{
	pos: position{line: 231, col: 16, offset: 6239},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 231, col: 16, offset: 6239},
	val: "//\\\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 231, col: 27, offset: 6250},
	val: "⩓",
	ignoreCase: false,
},
	},
},
},
{
	name: "Prefer",
	pos: position{line: 232, col: 1, offset: 6256},
	expr: &choiceExpr{
	pos: position{line: 232, col: 10, offset: 6267},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 232, col: 10, offset: 6267},
	val: "//",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 232, col: 17, offset: 6274},
	val: "⫽",
	ignoreCase: false,
},
	},
},
},
{
	name: "Lambda",
	pos: position{line: 233, col: 1, offset: 6280},
	expr: &choiceExpr{
	pos: position{line: 233, col: 10, offset: 6291},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 233, col: 10, offset: 6291},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 233, col: 17, offset: 6298},
	val: "λ",
	ignoreCase: false,
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 234, col: 1, offset: 6303},
	expr: &choiceExpr{
	pos: position{line: 234, col: 10, offset: 6314},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 234, col: 10, offset: 6314},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 234, col: 21, offset: 6325},
	val: "∀",
	ignoreCase: false,
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 235, col: 1, offset: 6331},
	expr: &choiceExpr{
	pos: position{line: 235, col: 9, offset: 6341},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 235, col: 9, offset: 6341},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 235, col: 16, offset: 6348},
	val: "→",
	ignoreCase: false,
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 237, col: 1, offset: 6355},
	expr: &seqExpr{
	pos: position{line: 237, col: 12, offset: 6368},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 237, col: 12, offset: 6368},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 237, col: 17, offset: 6373},
	expr: &charClassMatcher{
	pos: position{line: 237, col: 17, offset: 6373},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 237, col: 23, offset: 6379},
	expr: &ruleRefExpr{
	pos: position{line: 237, col: 23, offset: 6379},
	name: "Digit",
},
},
	},
},
},
{
	name: "NumericDoubleLiteral",
	pos: position{line: 239, col: 1, offset: 6387},
	expr: &actionExpr{
	pos: position{line: 239, col: 24, offset: 6412},
	run: (*parser).callonNumericDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 239, col: 24, offset: 6412},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 239, col: 24, offset: 6412},
	expr: &charClassMatcher{
	pos: position{line: 239, col: 24, offset: 6412},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 239, col: 30, offset: 6418},
	expr: &ruleRefExpr{
	pos: position{line: 239, col: 30, offset: 6418},
	name: "Digit",
},
},
&choiceExpr{
	pos: position{line: 239, col: 39, offset: 6427},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 239, col: 39, offset: 6427},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 239, col: 39, offset: 6427},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 239, col: 43, offset: 6431},
	expr: &ruleRefExpr{
	pos: position{line: 239, col: 43, offset: 6431},
	name: "Digit",
},
},
&zeroOrOneExpr{
	pos: position{line: 239, col: 50, offset: 6438},
	expr: &ruleRefExpr{
	pos: position{line: 239, col: 50, offset: 6438},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 239, col: 62, offset: 6450},
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
	pos: position{line: 247, col: 1, offset: 6606},
	expr: &choiceExpr{
	pos: position{line: 247, col: 17, offset: 6624},
	alternatives: []interface{}{
&labeledExpr{
	pos: position{line: 247, col: 17, offset: 6624},
	label: "d",
	expr: &ruleRefExpr{
	pos: position{line: 247, col: 19, offset: 6626},
	name: "NumericDoubleLiteral",
},
},
&actionExpr{
	pos: position{line: 248, col: 5, offset: 6651},
	run: (*parser).callonDoubleLiteral4,
	expr: &ruleRefExpr{
	pos: position{line: 248, col: 5, offset: 6651},
	name: "Infinity",
},
},
&actionExpr{
	pos: position{line: 249, col: 5, offset: 6703},
	run: (*parser).callonDoubleLiteral6,
	expr: &seqExpr{
	pos: position{line: 249, col: 5, offset: 6703},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 249, col: 5, offset: 6703},
	val: "-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 249, col: 9, offset: 6707},
	name: "Infinity",
},
	},
},
},
&actionExpr{
	pos: position{line: 250, col: 5, offset: 6760},
	run: (*parser).callonDoubleLiteral10,
	expr: &ruleRefExpr{
	pos: position{line: 250, col: 5, offset: 6760},
	name: "NaN",
},
},
	},
},
},
{
	name: "NaturalLiteral",
	pos: position{line: 252, col: 1, offset: 6803},
	expr: &actionExpr{
	pos: position{line: 252, col: 18, offset: 6822},
	run: (*parser).callonNaturalLiteral1,
	expr: &oneOrMoreExpr{
	pos: position{line: 252, col: 18, offset: 6822},
	expr: &ruleRefExpr{
	pos: position{line: 252, col: 18, offset: 6822},
	name: "Digit",
},
},
},
},
{
	name: "IntegerLiteral",
	pos: position{line: 257, col: 1, offset: 6911},
	expr: &actionExpr{
	pos: position{line: 257, col: 18, offset: 6930},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 257, col: 18, offset: 6930},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 257, col: 18, offset: 6930},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 257, col: 22, offset: 6934},
	name: "NaturalLiteral",
},
	},
},
},
},
{
	name: "DeBruijn",
	pos: position{line: 265, col: 1, offset: 7086},
	expr: &actionExpr{
	pos: position{line: 265, col: 12, offset: 7099},
	run: (*parser).callonDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 265, col: 12, offset: 7099},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 265, col: 12, offset: 7099},
	name: "_",
},
&litMatcher{
	pos: position{line: 265, col: 14, offset: 7101},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 265, col: 18, offset: 7105},
	name: "_",
},
&labeledExpr{
	pos: position{line: 265, col: 20, offset: 7107},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 265, col: 26, offset: 7113},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Variable",
	pos: position{line: 267, col: 1, offset: 7169},
	expr: &actionExpr{
	pos: position{line: 267, col: 12, offset: 7182},
	run: (*parser).callonVariable1,
	expr: &seqExpr{
	pos: position{line: 267, col: 12, offset: 7182},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 267, col: 12, offset: 7182},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 267, col: 17, offset: 7187},
	name: "NonreservedLabel",
},
},
&labeledExpr{
	pos: position{line: 267, col: 34, offset: 7204},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 267, col: 40, offset: 7210},
	expr: &ruleRefExpr{
	pos: position{line: 267, col: 40, offset: 7210},
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
	pos: position{line: 275, col: 1, offset: 7373},
	expr: &choiceExpr{
	pos: position{line: 275, col: 14, offset: 7388},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 275, col: 14, offset: 7388},
	name: "Variable",
},
&ruleRefExpr{
	pos: position{line: 275, col: 25, offset: 7399},
	name: "Reserved",
},
	},
},
},
{
	name: "PathCharacter",
	pos: position{line: 277, col: 1, offset: 7409},
	expr: &choiceExpr{
	pos: position{line: 278, col: 6, offset: 7432},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 278, col: 6, offset: 7432},
	val: "!",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 279, col: 6, offset: 7444},
	val: "[\\x24-\\x27]",
	ranges: []rune{'$','\'',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 280, col: 6, offset: 7461},
	val: "[\\x2a-\\x2b]",
	ranges: []rune{'*','+',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 281, col: 6, offset: 7478},
	val: "[\\x2d-\\x2e]",
	ranges: []rune{'-','.',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 282, col: 6, offset: 7495},
	val: "[\\x30-\\x3b]",
	ranges: []rune{'0',';',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 283, col: 6, offset: 7512},
	val: "=",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 284, col: 6, offset: 7524},
	val: "[\\x40-\\x5a]",
	ranges: []rune{'@','Z',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 285, col: 6, offset: 7541},
	val: "[\\x5e-\\x7a]",
	ranges: []rune{'^','z',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 286, col: 6, offset: 7558},
	val: "|",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 287, col: 6, offset: 7570},
	val: "~",
	ignoreCase: false,
},
	},
},
},
{
	name: "UnquotedPathComponent",
	pos: position{line: 289, col: 1, offset: 7578},
	expr: &actionExpr{
	pos: position{line: 289, col: 25, offset: 7604},
	run: (*parser).callonUnquotedPathComponent1,
	expr: &oneOrMoreExpr{
	pos: position{line: 289, col: 25, offset: 7604},
	expr: &ruleRefExpr{
	pos: position{line: 289, col: 25, offset: 7604},
	name: "PathCharacter",
},
},
},
},
{
	name: "PathComponent",
	pos: position{line: 291, col: 1, offset: 7651},
	expr: &actionExpr{
	pos: position{line: 291, col: 17, offset: 7669},
	run: (*parser).callonPathComponent1,
	expr: &seqExpr{
	pos: position{line: 291, col: 17, offset: 7669},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 291, col: 17, offset: 7669},
	val: "/",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 291, col: 21, offset: 7673},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 291, col: 23, offset: 7675},
	name: "UnquotedPathComponent",
},
},
	},
},
},
},
{
	name: "Path",
	pos: position{line: 293, col: 1, offset: 7716},
	expr: &actionExpr{
	pos: position{line: 293, col: 8, offset: 7725},
	run: (*parser).callonPath1,
	expr: &labeledExpr{
	pos: position{line: 293, col: 8, offset: 7725},
	label: "cs",
	expr: &oneOrMoreExpr{
	pos: position{line: 293, col: 11, offset: 7728},
	expr: &ruleRefExpr{
	pos: position{line: 293, col: 11, offset: 7728},
	name: "PathComponent",
},
},
},
},
},
{
	name: "Local",
	pos: position{line: 302, col: 1, offset: 8002},
	expr: &choiceExpr{
	pos: position{line: 302, col: 9, offset: 8012},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 302, col: 9, offset: 8012},
	name: "ParentPath",
},
&ruleRefExpr{
	pos: position{line: 302, col: 22, offset: 8025},
	name: "HerePath",
},
&ruleRefExpr{
	pos: position{line: 302, col: 33, offset: 8036},
	name: "HomePath",
},
&ruleRefExpr{
	pos: position{line: 302, col: 44, offset: 8047},
	name: "AbsolutePath",
},
	},
},
},
{
	name: "ParentPath",
	pos: position{line: 304, col: 1, offset: 8061},
	expr: &actionExpr{
	pos: position{line: 304, col: 14, offset: 8076},
	run: (*parser).callonParentPath1,
	expr: &seqExpr{
	pos: position{line: 304, col: 14, offset: 8076},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 304, col: 14, offset: 8076},
	val: "..",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 304, col: 19, offset: 8081},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 304, col: 21, offset: 8083},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HerePath",
	pos: position{line: 305, col: 1, offset: 8139},
	expr: &actionExpr{
	pos: position{line: 305, col: 12, offset: 8152},
	run: (*parser).callonHerePath1,
	expr: &seqExpr{
	pos: position{line: 305, col: 12, offset: 8152},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 305, col: 12, offset: 8152},
	val: ".",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 305, col: 16, offset: 8156},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 305, col: 18, offset: 8158},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HomePath",
	pos: position{line: 306, col: 1, offset: 8197},
	expr: &actionExpr{
	pos: position{line: 306, col: 12, offset: 8210},
	run: (*parser).callonHomePath1,
	expr: &seqExpr{
	pos: position{line: 306, col: 12, offset: 8210},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 306, col: 12, offset: 8210},
	val: "~",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 306, col: 16, offset: 8214},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 306, col: 18, offset: 8216},
	name: "Path",
},
},
	},
},
},
},
{
	name: "AbsolutePath",
	pos: position{line: 307, col: 1, offset: 8271},
	expr: &actionExpr{
	pos: position{line: 307, col: 16, offset: 8288},
	run: (*parser).callonAbsolutePath1,
	expr: &labeledExpr{
	pos: position{line: 307, col: 16, offset: 8288},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 307, col: 18, offset: 8290},
	name: "Path",
},
},
},
},
{
	name: "Scheme",
	pos: position{line: 309, col: 1, offset: 8346},
	expr: &seqExpr{
	pos: position{line: 309, col: 10, offset: 8357},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 309, col: 10, offset: 8357},
	val: "http",
	ignoreCase: false,
},
&zeroOrOneExpr{
	pos: position{line: 309, col: 17, offset: 8364},
	expr: &litMatcher{
	pos: position{line: 309, col: 17, offset: 8364},
	val: "s",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "HttpRaw",
	pos: position{line: 311, col: 1, offset: 8370},
	expr: &actionExpr{
	pos: position{line: 311, col: 11, offset: 8382},
	run: (*parser).callonHttpRaw1,
	expr: &seqExpr{
	pos: position{line: 311, col: 11, offset: 8382},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 311, col: 11, offset: 8382},
	name: "Scheme",
},
&litMatcher{
	pos: position{line: 311, col: 18, offset: 8389},
	val: "://",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 311, col: 24, offset: 8395},
	name: "Authority",
},
&ruleRefExpr{
	pos: position{line: 311, col: 34, offset: 8405},
	name: "Path",
},
&zeroOrOneExpr{
	pos: position{line: 311, col: 39, offset: 8410},
	expr: &seqExpr{
	pos: position{line: 311, col: 41, offset: 8412},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 311, col: 41, offset: 8412},
	val: "?",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 311, col: 45, offset: 8416},
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
	name: "Authority",
	pos: position{line: 313, col: 1, offset: 8473},
	expr: &seqExpr{
	pos: position{line: 313, col: 13, offset: 8487},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 313, col: 13, offset: 8487},
	expr: &seqExpr{
	pos: position{line: 313, col: 14, offset: 8488},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 313, col: 14, offset: 8488},
	name: "Userinfo",
},
&litMatcher{
	pos: position{line: 313, col: 23, offset: 8497},
	val: "@",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 313, col: 29, offset: 8503},
	name: "Host",
},
&zeroOrOneExpr{
	pos: position{line: 313, col: 34, offset: 8508},
	expr: &seqExpr{
	pos: position{line: 313, col: 35, offset: 8509},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 313, col: 35, offset: 8509},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 313, col: 39, offset: 8513},
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
	pos: position{line: 315, col: 1, offset: 8521},
	expr: &zeroOrMoreExpr{
	pos: position{line: 315, col: 12, offset: 8534},
	expr: &choiceExpr{
	pos: position{line: 315, col: 14, offset: 8536},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 315, col: 14, offset: 8536},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 315, col: 27, offset: 8549},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 315, col: 40, offset: 8562},
	name: "SubDelims",
},
&litMatcher{
	pos: position{line: 315, col: 52, offset: 8574},
	val: ":",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "Host",
	pos: position{line: 317, col: 1, offset: 8582},
	expr: &choiceExpr{
	pos: position{line: 317, col: 8, offset: 8591},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 317, col: 8, offset: 8591},
	name: "IPLiteral",
},
&ruleRefExpr{
	pos: position{line: 317, col: 20, offset: 8603},
	name: "RegName",
},
	},
},
},
{
	name: "Port",
	pos: position{line: 319, col: 1, offset: 8612},
	expr: &zeroOrMoreExpr{
	pos: position{line: 319, col: 8, offset: 8621},
	expr: &ruleRefExpr{
	pos: position{line: 319, col: 8, offset: 8621},
	name: "Digit",
},
},
},
{
	name: "IPLiteral",
	pos: position{line: 321, col: 1, offset: 8629},
	expr: &seqExpr{
	pos: position{line: 321, col: 13, offset: 8643},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 321, col: 13, offset: 8643},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 321, col: 17, offset: 8647},
	name: "IPv6address",
},
&litMatcher{
	pos: position{line: 321, col: 29, offset: 8659},
	val: "]",
	ignoreCase: false,
},
	},
},
},
{
	name: "IPv6address",
	pos: position{line: 323, col: 1, offset: 8664},
	expr: &actionExpr{
	pos: position{line: 323, col: 15, offset: 8680},
	run: (*parser).callonIPv6address1,
	expr: &seqExpr{
	pos: position{line: 323, col: 15, offset: 8680},
	exprs: []interface{}{
&zeroOrMoreExpr{
	pos: position{line: 323, col: 15, offset: 8680},
	expr: &ruleRefExpr{
	pos: position{line: 323, col: 16, offset: 8681},
	name: "HexDig",
},
},
&litMatcher{
	pos: position{line: 323, col: 25, offset: 8690},
	val: ":",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 323, col: 29, offset: 8694},
	expr: &choiceExpr{
	pos: position{line: 323, col: 30, offset: 8695},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 323, col: 30, offset: 8695},
	name: "HexDig",
},
&litMatcher{
	pos: position{line: 323, col: 39, offset: 8704},
	val: ":",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 323, col: 45, offset: 8710},
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
	pos: position{line: 329, col: 1, offset: 8864},
	expr: &zeroOrMoreExpr{
	pos: position{line: 329, col: 11, offset: 8876},
	expr: &choiceExpr{
	pos: position{line: 329, col: 12, offset: 8877},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 329, col: 12, offset: 8877},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 329, col: 25, offset: 8890},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 329, col: 38, offset: 8903},
	name: "SubDelims",
},
	},
},
},
},
{
	name: "PChar",
	pos: position{line: 331, col: 1, offset: 8916},
	expr: &choiceExpr{
	pos: position{line: 331, col: 9, offset: 8926},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 331, col: 9, offset: 8926},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 331, col: 22, offset: 8939},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 331, col: 35, offset: 8952},
	name: "SubDelims",
},
&charClassMatcher{
	pos: position{line: 331, col: 47, offset: 8964},
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
	pos: position{line: 333, col: 1, offset: 8970},
	expr: &zeroOrMoreExpr{
	pos: position{line: 333, col: 9, offset: 8980},
	expr: &choiceExpr{
	pos: position{line: 333, col: 10, offset: 8981},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 333, col: 10, offset: 8981},
	name: "PChar",
},
&charClassMatcher{
	pos: position{line: 333, col: 18, offset: 8989},
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
	pos: position{line: 335, col: 1, offset: 8997},
	expr: &seqExpr{
	pos: position{line: 335, col: 14, offset: 9012},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 335, col: 14, offset: 9012},
	val: "%",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 335, col: 18, offset: 9016},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 335, col: 25, offset: 9023},
	name: "HexDig",
},
	},
},
},
{
	name: "Unreserved",
	pos: position{line: 337, col: 1, offset: 9031},
	expr: &charClassMatcher{
	pos: position{line: 337, col: 14, offset: 9046},
	val: "[A-Za-z0-9._~-]",
	chars: []rune{'.','_','~','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SubDelims",
	pos: position{line: 339, col: 1, offset: 9063},
	expr: &choiceExpr{
	pos: position{line: 339, col: 13, offset: 9077},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 339, col: 13, offset: 9077},
	val: "!",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 339, col: 19, offset: 9083},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 339, col: 25, offset: 9089},
	val: "&",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 339, col: 31, offset: 9095},
	val: "'",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 339, col: 37, offset: 9101},
	val: "(",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 339, col: 43, offset: 9107},
	val: ")",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 339, col: 49, offset: 9113},
	val: "*",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 339, col: 55, offset: 9119},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 339, col: 61, offset: 9125},
	val: ",",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 339, col: 67, offset: 9131},
	val: ";",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 339, col: 73, offset: 9137},
	val: "=",
	ignoreCase: false,
},
	},
},
},
{
	name: "Http",
	pos: position{line: 341, col: 1, offset: 9142},
	expr: &actionExpr{
	pos: position{line: 341, col: 8, offset: 9151},
	run: (*parser).callonHttp1,
	expr: &labeledExpr{
	pos: position{line: 341, col: 8, offset: 9151},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 341, col: 10, offset: 9153},
	name: "HttpRaw",
},
},
},
},
{
	name: "Env",
	pos: position{line: 343, col: 1, offset: 9198},
	expr: &actionExpr{
	pos: position{line: 343, col: 7, offset: 9206},
	run: (*parser).callonEnv1,
	expr: &seqExpr{
	pos: position{line: 343, col: 7, offset: 9206},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 343, col: 7, offset: 9206},
	val: "env:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 343, col: 14, offset: 9213},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 343, col: 17, offset: 9216},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 343, col: 17, offset: 9216},
	name: "BashEnvironmentVariable",
},
&ruleRefExpr{
	pos: position{line: 343, col: 43, offset: 9242},
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
	pos: position{line: 345, col: 1, offset: 9287},
	expr: &actionExpr{
	pos: position{line: 345, col: 27, offset: 9315},
	run: (*parser).callonBashEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 345, col: 27, offset: 9315},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 345, col: 27, offset: 9315},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 345, col: 36, offset: 9324},
	expr: &charClassMatcher{
	pos: position{line: 345, col: 36, offset: 9324},
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
	pos: position{line: 349, col: 1, offset: 9380},
	expr: &actionExpr{
	pos: position{line: 349, col: 28, offset: 9409},
	run: (*parser).callonPosixEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 349, col: 28, offset: 9409},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 349, col: 28, offset: 9409},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 349, col: 32, offset: 9413},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 349, col: 34, offset: 9415},
	name: "PosixEnvironmentVariableContent",
},
},
&litMatcher{
	pos: position{line: 349, col: 66, offset: 9447},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PosixEnvironmentVariableContent",
	pos: position{line: 353, col: 1, offset: 9472},
	expr: &actionExpr{
	pos: position{line: 353, col: 35, offset: 9508},
	run: (*parser).callonPosixEnvironmentVariableContent1,
	expr: &labeledExpr{
	pos: position{line: 353, col: 35, offset: 9508},
	label: "v",
	expr: &oneOrMoreExpr{
	pos: position{line: 353, col: 37, offset: 9510},
	expr: &ruleRefExpr{
	pos: position{line: 353, col: 37, offset: 9510},
	name: "PosixEnvironmentVariableCharacter",
},
},
},
},
},
{
	name: "PosixEnvironmentVariableCharacter",
	pos: position{line: 362, col: 1, offset: 9723},
	expr: &choiceExpr{
	pos: position{line: 363, col: 7, offset: 9767},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 363, col: 7, offset: 9767},
	run: (*parser).callonPosixEnvironmentVariableCharacter2,
	expr: &litMatcher{
	pos: position{line: 363, col: 7, offset: 9767},
	val: "\\\"",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 364, col: 7, offset: 9807},
	run: (*parser).callonPosixEnvironmentVariableCharacter4,
	expr: &litMatcher{
	pos: position{line: 364, col: 7, offset: 9807},
	val: "\\\\",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 365, col: 7, offset: 9847},
	run: (*parser).callonPosixEnvironmentVariableCharacter6,
	expr: &litMatcher{
	pos: position{line: 365, col: 7, offset: 9847},
	val: "\\a",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 366, col: 7, offset: 9887},
	run: (*parser).callonPosixEnvironmentVariableCharacter8,
	expr: &litMatcher{
	pos: position{line: 366, col: 7, offset: 9887},
	val: "\\b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 367, col: 7, offset: 9927},
	run: (*parser).callonPosixEnvironmentVariableCharacter10,
	expr: &litMatcher{
	pos: position{line: 367, col: 7, offset: 9927},
	val: "\\f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 368, col: 7, offset: 9967},
	run: (*parser).callonPosixEnvironmentVariableCharacter12,
	expr: &litMatcher{
	pos: position{line: 368, col: 7, offset: 9967},
	val: "\\n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 369, col: 7, offset: 10007},
	run: (*parser).callonPosixEnvironmentVariableCharacter14,
	expr: &litMatcher{
	pos: position{line: 369, col: 7, offset: 10007},
	val: "\\r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 370, col: 7, offset: 10047},
	run: (*parser).callonPosixEnvironmentVariableCharacter16,
	expr: &litMatcher{
	pos: position{line: 370, col: 7, offset: 10047},
	val: "\\t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 371, col: 7, offset: 10087},
	run: (*parser).callonPosixEnvironmentVariableCharacter18,
	expr: &litMatcher{
	pos: position{line: 371, col: 7, offset: 10087},
	val: "\\v",
	ignoreCase: false,
},
},
&charClassMatcher{
	pos: position{line: 372, col: 7, offset: 10127},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 373, col: 7, offset: 10145},
	val: "[\\x23-\\x3c]",
	ranges: []rune{'#','<',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 374, col: 7, offset: 10163},
	val: "[\\x3e-\\x5b]",
	ranges: []rune{'>','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 375, col: 7, offset: 10181},
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
	pos: position{line: 377, col: 1, offset: 10194},
	expr: &choiceExpr{
	pos: position{line: 377, col: 14, offset: 10209},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 377, col: 14, offset: 10209},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 377, col: 24, offset: 10219},
	name: "Local",
},
&ruleRefExpr{
	pos: position{line: 377, col: 32, offset: 10227},
	name: "Http",
},
&ruleRefExpr{
	pos: position{line: 377, col: 39, offset: 10234},
	name: "Env",
},
	},
},
},
{
	name: "ImportHashed",
	pos: position{line: 379, col: 1, offset: 10239},
	expr: &actionExpr{
	pos: position{line: 379, col: 16, offset: 10256},
	run: (*parser).callonImportHashed1,
	expr: &labeledExpr{
	pos: position{line: 379, col: 16, offset: 10256},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 379, col: 18, offset: 10258},
	name: "ImportType",
},
},
},
},
{
	name: "Import",
	pos: position{line: 381, col: 1, offset: 10325},
	expr: &choiceExpr{
	pos: position{line: 381, col: 10, offset: 10336},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 381, col: 10, offset: 10336},
	run: (*parser).callonImport2,
	expr: &seqExpr{
	pos: position{line: 381, col: 10, offset: 10336},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 381, col: 10, offset: 10336},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 381, col: 12, offset: 10338},
	name: "ImportHashed",
},
},
&ruleRefExpr{
	pos: position{line: 381, col: 25, offset: 10351},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 381, col: 27, offset: 10353},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 381, col: 30, offset: 10356},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 381, col: 33, offset: 10359},
	name: "Text",
},
	},
},
},
&actionExpr{
	pos: position{line: 382, col: 10, offset: 10456},
	run: (*parser).callonImport10,
	expr: &labeledExpr{
	pos: position{line: 382, col: 10, offset: 10456},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 382, col: 12, offset: 10458},
	name: "ImportHashed",
},
},
},
	},
},
},
{
	name: "LetBinding",
	pos: position{line: 385, col: 1, offset: 10553},
	expr: &actionExpr{
	pos: position{line: 385, col: 14, offset: 10568},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 385, col: 14, offset: 10568},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 385, col: 14, offset: 10568},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 385, col: 18, offset: 10572},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 385, col: 21, offset: 10575},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 385, col: 27, offset: 10581},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 385, col: 44, offset: 10598},
	name: "_",
},
&labeledExpr{
	pos: position{line: 385, col: 46, offset: 10600},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 385, col: 48, offset: 10602},
	expr: &seqExpr{
	pos: position{line: 385, col: 49, offset: 10603},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 385, col: 49, offset: 10603},
	name: "Annotation",
},
&ruleRefExpr{
	pos: position{line: 385, col: 60, offset: 10614},
	name: "_",
},
	},
},
},
},
&litMatcher{
	pos: position{line: 386, col: 13, offset: 10630},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 386, col: 17, offset: 10634},
	name: "_",
},
&labeledExpr{
	pos: position{line: 386, col: 19, offset: 10636},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 386, col: 21, offset: 10638},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 386, col: 32, offset: 10649},
	name: "_",
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 401, col: 1, offset: 10958},
	expr: &choiceExpr{
	pos: position{line: 402, col: 7, offset: 10979},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 402, col: 7, offset: 10979},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 402, col: 7, offset: 10979},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 402, col: 7, offset: 10979},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 402, col: 14, offset: 10986},
	name: "_",
},
&litMatcher{
	pos: position{line: 402, col: 16, offset: 10988},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 402, col: 20, offset: 10992},
	name: "_",
},
&labeledExpr{
	pos: position{line: 402, col: 22, offset: 10994},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 402, col: 28, offset: 11000},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 402, col: 45, offset: 11017},
	name: "_",
},
&litMatcher{
	pos: position{line: 402, col: 47, offset: 11019},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 402, col: 51, offset: 11023},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 402, col: 54, offset: 11026},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 402, col: 56, offset: 11028},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 402, col: 67, offset: 11039},
	name: "_",
},
&litMatcher{
	pos: position{line: 402, col: 69, offset: 11041},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 402, col: 73, offset: 11045},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 402, col: 75, offset: 11047},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 402, col: 81, offset: 11053},
	name: "_",
},
&labeledExpr{
	pos: position{line: 402, col: 83, offset: 11055},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 402, col: 88, offset: 11060},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 405, col: 7, offset: 11177},
	run: (*parser).callonExpression22,
	expr: &seqExpr{
	pos: position{line: 405, col: 7, offset: 11177},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 405, col: 7, offset: 11177},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 405, col: 10, offset: 11180},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 405, col: 13, offset: 11183},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 405, col: 18, offset: 11188},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 405, col: 29, offset: 11199},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 405, col: 31, offset: 11201},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 405, col: 36, offset: 11206},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 405, col: 39, offset: 11209},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 405, col: 41, offset: 11211},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 405, col: 52, offset: 11222},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 405, col: 54, offset: 11224},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 405, col: 59, offset: 11229},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 405, col: 62, offset: 11232},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 405, col: 64, offset: 11234},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 408, col: 7, offset: 11320},
	run: (*parser).callonExpression38,
	expr: &seqExpr{
	pos: position{line: 408, col: 7, offset: 11320},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 408, col: 7, offset: 11320},
	label: "bindings",
	expr: &oneOrMoreExpr{
	pos: position{line: 408, col: 16, offset: 11329},
	expr: &ruleRefExpr{
	pos: position{line: 408, col: 16, offset: 11329},
	name: "LetBinding",
},
},
},
&ruleRefExpr{
	pos: position{line: 408, col: 28, offset: 11341},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 408, col: 31, offset: 11344},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 408, col: 34, offset: 11347},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 408, col: 36, offset: 11349},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 415, col: 7, offset: 11589},
	run: (*parser).callonExpression47,
	expr: &seqExpr{
	pos: position{line: 415, col: 7, offset: 11589},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 415, col: 7, offset: 11589},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 415, col: 14, offset: 11596},
	name: "_",
},
&litMatcher{
	pos: position{line: 415, col: 16, offset: 11598},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 415, col: 20, offset: 11602},
	name: "_",
},
&labeledExpr{
	pos: position{line: 415, col: 22, offset: 11604},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 415, col: 28, offset: 11610},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 415, col: 45, offset: 11627},
	name: "_",
},
&litMatcher{
	pos: position{line: 415, col: 47, offset: 11629},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 415, col: 51, offset: 11633},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 415, col: 54, offset: 11636},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 415, col: 56, offset: 11638},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 415, col: 67, offset: 11649},
	name: "_",
},
&litMatcher{
	pos: position{line: 415, col: 69, offset: 11651},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 415, col: 73, offset: 11655},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 415, col: 75, offset: 11657},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 415, col: 81, offset: 11663},
	name: "_",
},
&labeledExpr{
	pos: position{line: 415, col: 83, offset: 11665},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 415, col: 88, offset: 11670},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 418, col: 7, offset: 11779},
	run: (*parser).callonExpression67,
	expr: &seqExpr{
	pos: position{line: 418, col: 7, offset: 11779},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 418, col: 7, offset: 11779},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 418, col: 9, offset: 11781},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 418, col: 28, offset: 11800},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 418, col: 30, offset: 11802},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 418, col: 36, offset: 11808},
	name: "_",
},
&labeledExpr{
	pos: position{line: 418, col: 38, offset: 11810},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 418, col: 40, offset: 11812},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 419, col: 7, offset: 11871},
	run: (*parser).callonExpression76,
	expr: &seqExpr{
	pos: position{line: 419, col: 7, offset: 11871},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 419, col: 7, offset: 11871},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 419, col: 13, offset: 11877},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 419, col: 16, offset: 11880},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 419, col: 18, offset: 11882},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 419, col: 35, offset: 11899},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 419, col: 38, offset: 11902},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 419, col: 40, offset: 11904},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 419, col: 57, offset: 11921},
	name: "_",
},
&litMatcher{
	pos: position{line: 419, col: 59, offset: 11923},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 419, col: 63, offset: 11927},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 419, col: 66, offset: 11930},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 419, col: 68, offset: 11932},
	name: "ApplicationExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 422, col: 7, offset: 12053},
	name: "EmptyList",
},
&ruleRefExpr{
	pos: position{line: 423, col: 7, offset: 12069},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 425, col: 1, offset: 12090},
	expr: &actionExpr{
	pos: position{line: 425, col: 14, offset: 12105},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 425, col: 14, offset: 12105},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 425, col: 14, offset: 12105},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 425, col: 18, offset: 12109},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 425, col: 21, offset: 12112},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 425, col: 23, offset: 12114},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 427, col: 1, offset: 12144},
	expr: &actionExpr{
	pos: position{line: 428, col: 1, offset: 12168},
	run: (*parser).callonAnnotatedExpression1,
	expr: &seqExpr{
	pos: position{line: 428, col: 1, offset: 12168},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 428, col: 1, offset: 12168},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 428, col: 3, offset: 12170},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 428, col: 22, offset: 12189},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 428, col: 24, offset: 12191},
	expr: &seqExpr{
	pos: position{line: 428, col: 25, offset: 12192},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 428, col: 25, offset: 12192},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 428, col: 27, offset: 12194},
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
	pos: position{line: 433, col: 1, offset: 12319},
	expr: &actionExpr{
	pos: position{line: 433, col: 13, offset: 12333},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 433, col: 13, offset: 12333},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 433, col: 13, offset: 12333},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 433, col: 17, offset: 12337},
	name: "_",
},
&litMatcher{
	pos: position{line: 433, col: 19, offset: 12339},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 433, col: 23, offset: 12343},
	name: "_",
},
&litMatcher{
	pos: position{line: 433, col: 25, offset: 12345},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 433, col: 29, offset: 12349},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 433, col: 32, offset: 12352},
	name: "List",
},
&ruleRefExpr{
	pos: position{line: 433, col: 37, offset: 12357},
	name: "_",
},
&labeledExpr{
	pos: position{line: 433, col: 39, offset: 12359},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 433, col: 41, offset: 12361},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 437, col: 1, offset: 12424},
	expr: &ruleRefExpr{
	pos: position{line: 437, col: 22, offset: 12447},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 439, col: 1, offset: 12468},
	expr: &ruleRefExpr{
	pos: position{line: 439, col: 24, offset: 12493},
	name: "OrExpression",
},
},
{
	name: "OrExpression",
	pos: position{line: 441, col: 1, offset: 12507},
	expr: &actionExpr{
	pos: position{line: 441, col: 26, offset: 12534},
	run: (*parser).callonOrExpression1,
	expr: &seqExpr{
	pos: position{line: 441, col: 26, offset: 12534},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 441, col: 26, offset: 12534},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 441, col: 32, offset: 12540},
	name: "PlusExpression",
},
},
&labeledExpr{
	pos: position{line: 441, col: 55, offset: 12563},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 441, col: 60, offset: 12568},
	expr: &seqExpr{
	pos: position{line: 441, col: 61, offset: 12569},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 441, col: 61, offset: 12569},
	name: "_",
},
&litMatcher{
	pos: position{line: 441, col: 63, offset: 12571},
	val: "||",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 441, col: 68, offset: 12576},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 441, col: 70, offset: 12578},
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
	pos: position{line: 443, col: 1, offset: 12644},
	expr: &actionExpr{
	pos: position{line: 443, col: 26, offset: 12671},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 443, col: 26, offset: 12671},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 443, col: 26, offset: 12671},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 443, col: 32, offset: 12677},
	name: "TextAppendExpression",
},
},
&labeledExpr{
	pos: position{line: 443, col: 55, offset: 12700},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 443, col: 60, offset: 12705},
	expr: &seqExpr{
	pos: position{line: 443, col: 61, offset: 12706},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 443, col: 61, offset: 12706},
	name: "_",
},
&litMatcher{
	pos: position{line: 443, col: 63, offset: 12708},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 443, col: 67, offset: 12712},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 443, col: 70, offset: 12715},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 443, col: 72, offset: 12717},
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
	pos: position{line: 445, col: 1, offset: 12791},
	expr: &actionExpr{
	pos: position{line: 445, col: 26, offset: 12818},
	run: (*parser).callonTextAppendExpression1,
	expr: &seqExpr{
	pos: position{line: 445, col: 26, offset: 12818},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 445, col: 26, offset: 12818},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 445, col: 32, offset: 12824},
	name: "ListAppendExpression",
},
},
&labeledExpr{
	pos: position{line: 445, col: 55, offset: 12847},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 445, col: 60, offset: 12852},
	expr: &seqExpr{
	pos: position{line: 445, col: 61, offset: 12853},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 445, col: 61, offset: 12853},
	name: "_",
},
&litMatcher{
	pos: position{line: 445, col: 63, offset: 12855},
	val: "++",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 445, col: 68, offset: 12860},
	name: "_",
},
&labeledExpr{
	pos: position{line: 445, col: 70, offset: 12862},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 445, col: 72, offset: 12864},
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
	pos: position{line: 447, col: 1, offset: 12944},
	expr: &actionExpr{
	pos: position{line: 447, col: 26, offset: 12971},
	run: (*parser).callonListAppendExpression1,
	expr: &seqExpr{
	pos: position{line: 447, col: 26, offset: 12971},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 447, col: 26, offset: 12971},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 447, col: 32, offset: 12977},
	name: "AndExpression",
},
},
&labeledExpr{
	pos: position{line: 447, col: 55, offset: 13000},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 447, col: 60, offset: 13005},
	expr: &seqExpr{
	pos: position{line: 447, col: 61, offset: 13006},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 447, col: 61, offset: 13006},
	name: "_",
},
&litMatcher{
	pos: position{line: 447, col: 63, offset: 13008},
	val: "#",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 447, col: 67, offset: 13012},
	name: "_",
},
&labeledExpr{
	pos: position{line: 447, col: 69, offset: 13014},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 447, col: 71, offset: 13016},
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
	pos: position{line: 449, col: 1, offset: 13089},
	expr: &actionExpr{
	pos: position{line: 449, col: 26, offset: 13116},
	run: (*parser).callonAndExpression1,
	expr: &seqExpr{
	pos: position{line: 449, col: 26, offset: 13116},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 449, col: 26, offset: 13116},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 449, col: 32, offset: 13122},
	name: "CombineExpression",
},
},
&labeledExpr{
	pos: position{line: 449, col: 55, offset: 13145},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 449, col: 60, offset: 13150},
	expr: &seqExpr{
	pos: position{line: 449, col: 61, offset: 13151},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 449, col: 61, offset: 13151},
	name: "_",
},
&litMatcher{
	pos: position{line: 449, col: 63, offset: 13153},
	val: "&&",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 449, col: 68, offset: 13158},
	name: "_",
},
&labeledExpr{
	pos: position{line: 449, col: 70, offset: 13160},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 449, col: 72, offset: 13162},
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
	pos: position{line: 451, col: 1, offset: 13232},
	expr: &actionExpr{
	pos: position{line: 451, col: 26, offset: 13259},
	run: (*parser).callonCombineExpression1,
	expr: &seqExpr{
	pos: position{line: 451, col: 26, offset: 13259},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 451, col: 26, offset: 13259},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 451, col: 32, offset: 13265},
	name: "PreferExpression",
},
},
&labeledExpr{
	pos: position{line: 451, col: 55, offset: 13288},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 451, col: 60, offset: 13293},
	expr: &seqExpr{
	pos: position{line: 451, col: 61, offset: 13294},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 451, col: 61, offset: 13294},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 451, col: 63, offset: 13296},
	name: "Combine",
},
&ruleRefExpr{
	pos: position{line: 451, col: 71, offset: 13304},
	name: "_",
},
&labeledExpr{
	pos: position{line: 451, col: 73, offset: 13306},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 451, col: 75, offset: 13308},
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
	pos: position{line: 453, col: 1, offset: 13385},
	expr: &actionExpr{
	pos: position{line: 453, col: 26, offset: 13412},
	run: (*parser).callonPreferExpression1,
	expr: &seqExpr{
	pos: position{line: 453, col: 26, offset: 13412},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 453, col: 26, offset: 13412},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 453, col: 32, offset: 13418},
	name: "CombineTypesExpression",
},
},
&labeledExpr{
	pos: position{line: 453, col: 55, offset: 13441},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 453, col: 60, offset: 13446},
	expr: &seqExpr{
	pos: position{line: 453, col: 61, offset: 13447},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 453, col: 61, offset: 13447},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 453, col: 63, offset: 13449},
	name: "Prefer",
},
&ruleRefExpr{
	pos: position{line: 453, col: 70, offset: 13456},
	name: "_",
},
&labeledExpr{
	pos: position{line: 453, col: 72, offset: 13458},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 453, col: 74, offset: 13460},
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
	pos: position{line: 455, col: 1, offset: 13554},
	expr: &actionExpr{
	pos: position{line: 455, col: 26, offset: 13581},
	run: (*parser).callonCombineTypesExpression1,
	expr: &seqExpr{
	pos: position{line: 455, col: 26, offset: 13581},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 455, col: 26, offset: 13581},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 455, col: 32, offset: 13587},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 455, col: 55, offset: 13610},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 455, col: 60, offset: 13615},
	expr: &seqExpr{
	pos: position{line: 455, col: 61, offset: 13616},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 455, col: 61, offset: 13616},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 455, col: 63, offset: 13618},
	name: "CombineTypes",
},
&ruleRefExpr{
	pos: position{line: 455, col: 76, offset: 13631},
	name: "_",
},
&labeledExpr{
	pos: position{line: 455, col: 78, offset: 13633},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 455, col: 80, offset: 13635},
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
	pos: position{line: 457, col: 1, offset: 13715},
	expr: &actionExpr{
	pos: position{line: 457, col: 26, offset: 13742},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 457, col: 26, offset: 13742},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 457, col: 26, offset: 13742},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 457, col: 32, offset: 13748},
	name: "EqualExpression",
},
},
&labeledExpr{
	pos: position{line: 457, col: 55, offset: 13771},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 457, col: 60, offset: 13776},
	expr: &seqExpr{
	pos: position{line: 457, col: 61, offset: 13777},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 457, col: 61, offset: 13777},
	name: "_",
},
&litMatcher{
	pos: position{line: 457, col: 63, offset: 13779},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 457, col: 67, offset: 13783},
	name: "_",
},
&labeledExpr{
	pos: position{line: 457, col: 69, offset: 13785},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 457, col: 71, offset: 13787},
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
	pos: position{line: 459, col: 1, offset: 13857},
	expr: &actionExpr{
	pos: position{line: 459, col: 26, offset: 13884},
	run: (*parser).callonEqualExpression1,
	expr: &seqExpr{
	pos: position{line: 459, col: 26, offset: 13884},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 459, col: 26, offset: 13884},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 459, col: 32, offset: 13890},
	name: "NotEqualExpression",
},
},
&labeledExpr{
	pos: position{line: 459, col: 55, offset: 13913},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 459, col: 60, offset: 13918},
	expr: &seqExpr{
	pos: position{line: 459, col: 61, offset: 13919},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 459, col: 61, offset: 13919},
	name: "_",
},
&litMatcher{
	pos: position{line: 459, col: 63, offset: 13921},
	val: "==",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 459, col: 68, offset: 13926},
	name: "_",
},
&labeledExpr{
	pos: position{line: 459, col: 70, offset: 13928},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 459, col: 72, offset: 13930},
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
	pos: position{line: 461, col: 1, offset: 14000},
	expr: &actionExpr{
	pos: position{line: 461, col: 26, offset: 14027},
	run: (*parser).callonNotEqualExpression1,
	expr: &seqExpr{
	pos: position{line: 461, col: 26, offset: 14027},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 461, col: 26, offset: 14027},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 461, col: 32, offset: 14033},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 461, col: 55, offset: 14056},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 461, col: 60, offset: 14061},
	expr: &seqExpr{
	pos: position{line: 461, col: 61, offset: 14062},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 461, col: 61, offset: 14062},
	name: "_",
},
&litMatcher{
	pos: position{line: 461, col: 63, offset: 14064},
	val: "!=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 461, col: 68, offset: 14069},
	name: "_",
},
&labeledExpr{
	pos: position{line: 461, col: 70, offset: 14071},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 461, col: 72, offset: 14073},
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
	pos: position{line: 464, col: 1, offset: 14147},
	expr: &actionExpr{
	pos: position{line: 464, col: 25, offset: 14173},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 464, col: 25, offset: 14173},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 464, col: 25, offset: 14173},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 464, col: 27, offset: 14175},
	name: "FirstApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 464, col: 54, offset: 14202},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 464, col: 59, offset: 14207},
	expr: &seqExpr{
	pos: position{line: 464, col: 60, offset: 14208},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 464, col: 60, offset: 14208},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 464, col: 63, offset: 14211},
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
	pos: position{line: 473, col: 1, offset: 14454},
	expr: &choiceExpr{
	pos: position{line: 474, col: 8, offset: 14492},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 474, col: 8, offset: 14492},
	run: (*parser).callonFirstApplicationExpression2,
	expr: &seqExpr{
	pos: position{line: 474, col: 8, offset: 14492},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 474, col: 8, offset: 14492},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 474, col: 14, offset: 14498},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 474, col: 17, offset: 14501},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 474, col: 19, offset: 14503},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 474, col: 36, offset: 14520},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 474, col: 39, offset: 14523},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 474, col: 41, offset: 14525},
	name: "ImportExpression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 477, col: 8, offset: 14628},
	run: (*parser).callonFirstApplicationExpression11,
	expr: &seqExpr{
	pos: position{line: 477, col: 8, offset: 14628},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 477, col: 8, offset: 14628},
	name: "Some",
},
&ruleRefExpr{
	pos: position{line: 477, col: 13, offset: 14633},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 477, col: 16, offset: 14636},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 477, col: 18, offset: 14638},
	name: "ImportExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 478, col: 8, offset: 14693},
	name: "ImportExpression",
},
	},
},
},
{
	name: "ImportExpression",
	pos: position{line: 480, col: 1, offset: 14711},
	expr: &choiceExpr{
	pos: position{line: 480, col: 20, offset: 14732},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 480, col: 20, offset: 14732},
	name: "Import",
},
&ruleRefExpr{
	pos: position{line: 480, col: 29, offset: 14741},
	name: "SelectorExpression",
},
	},
},
},
{
	name: "SelectorExpression",
	pos: position{line: 482, col: 1, offset: 14761},
	expr: &actionExpr{
	pos: position{line: 482, col: 22, offset: 14784},
	run: (*parser).callonSelectorExpression1,
	expr: &seqExpr{
	pos: position{line: 482, col: 22, offset: 14784},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 482, col: 22, offset: 14784},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 482, col: 24, offset: 14786},
	name: "PrimitiveExpression",
},
},
&labeledExpr{
	pos: position{line: 482, col: 44, offset: 14806},
	label: "ls",
	expr: &zeroOrMoreExpr{
	pos: position{line: 482, col: 47, offset: 14809},
	expr: &seqExpr{
	pos: position{line: 482, col: 48, offset: 14810},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 482, col: 48, offset: 14810},
	name: "_",
},
&litMatcher{
	pos: position{line: 482, col: 50, offset: 14812},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 482, col: 54, offset: 14816},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 482, col: 56, offset: 14818},
	name: "AnyLabel",
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
	name: "PrimitiveExpression",
	pos: position{line: 492, col: 1, offset: 15051},
	expr: &choiceExpr{
	pos: position{line: 493, col: 7, offset: 15081},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 493, col: 7, offset: 15081},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 494, col: 7, offset: 15101},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 495, col: 7, offset: 15122},
	name: "IntegerLiteral",
},
&ruleRefExpr{
	pos: position{line: 496, col: 7, offset: 15143},
	name: "TextLiteral",
},
&actionExpr{
	pos: position{line: 497, col: 7, offset: 15161},
	run: (*parser).callonPrimitiveExpression6,
	expr: &seqExpr{
	pos: position{line: 497, col: 7, offset: 15161},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 497, col: 7, offset: 15161},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 497, col: 11, offset: 15165},
	name: "_",
},
&labeledExpr{
	pos: position{line: 497, col: 13, offset: 15167},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 497, col: 15, offset: 15169},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 497, col: 35, offset: 15189},
	name: "_",
},
&litMatcher{
	pos: position{line: 497, col: 37, offset: 15191},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&actionExpr{
	pos: position{line: 498, col: 7, offset: 15219},
	run: (*parser).callonPrimitiveExpression14,
	expr: &seqExpr{
	pos: position{line: 498, col: 7, offset: 15219},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 498, col: 7, offset: 15219},
	val: "<",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 498, col: 11, offset: 15223},
	name: "_",
},
&labeledExpr{
	pos: position{line: 498, col: 13, offset: 15225},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 498, col: 15, offset: 15227},
	name: "UnionType",
},
},
&ruleRefExpr{
	pos: position{line: 498, col: 25, offset: 15237},
	name: "_",
},
&litMatcher{
	pos: position{line: 498, col: 27, offset: 15239},
	val: ">",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 499, col: 7, offset: 15267},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 500, col: 7, offset: 15293},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 501, col: 7, offset: 15310},
	run: (*parser).callonPrimitiveExpression24,
	expr: &seqExpr{
	pos: position{line: 501, col: 7, offset: 15310},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 501, col: 7, offset: 15310},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 501, col: 11, offset: 15314},
	name: "_",
},
&labeledExpr{
	pos: position{line: 501, col: 14, offset: 15317},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 501, col: 16, offset: 15319},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 501, col: 27, offset: 15330},
	name: "_",
},
&litMatcher{
	pos: position{line: 501, col: 29, offset: 15332},
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
	pos: position{line: 503, col: 1, offset: 15355},
	expr: &choiceExpr{
	pos: position{line: 504, col: 7, offset: 15385},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 504, col: 7, offset: 15385},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &litMatcher{
	pos: position{line: 504, col: 7, offset: 15385},
	val: "=",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 505, col: 7, offset: 15440},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 506, col: 7, offset: 15465},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 507, col: 7, offset: 15493},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 507, col: 7, offset: 15493},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 509, col: 1, offset: 15539},
	expr: &actionExpr{
	pos: position{line: 509, col: 19, offset: 15559},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 509, col: 19, offset: 15559},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 509, col: 19, offset: 15559},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 509, col: 24, offset: 15564},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 509, col: 33, offset: 15573},
	name: "_",
},
&litMatcher{
	pos: position{line: 509, col: 35, offset: 15575},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 509, col: 39, offset: 15579},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 509, col: 42, offset: 15582},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 509, col: 47, offset: 15587},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 512, col: 1, offset: 15644},
	expr: &actionExpr{
	pos: position{line: 512, col: 18, offset: 15663},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 512, col: 18, offset: 15663},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 512, col: 18, offset: 15663},
	name: "_",
},
&litMatcher{
	pos: position{line: 512, col: 20, offset: 15665},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 512, col: 24, offset: 15669},
	name: "_",
},
&labeledExpr{
	pos: position{line: 512, col: 26, offset: 15671},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 512, col: 28, offset: 15673},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 513, col: 1, offset: 15705},
	expr: &actionExpr{
	pos: position{line: 514, col: 7, offset: 15734},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 514, col: 7, offset: 15734},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 514, col: 7, offset: 15734},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 514, col: 13, offset: 15740},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 514, col: 29, offset: 15756},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 514, col: 34, offset: 15761},
	expr: &ruleRefExpr{
	pos: position{line: 514, col: 34, offset: 15761},
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
	pos: position{line: 528, col: 1, offset: 16345},
	expr: &actionExpr{
	pos: position{line: 528, col: 22, offset: 16368},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 528, col: 22, offset: 16368},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 528, col: 22, offset: 16368},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 528, col: 27, offset: 16373},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 528, col: 36, offset: 16382},
	name: "_",
},
&litMatcher{
	pos: position{line: 528, col: 38, offset: 16384},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 528, col: 42, offset: 16388},
	name: "_",
},
&labeledExpr{
	pos: position{line: 528, col: 44, offset: 16390},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 528, col: 49, offset: 16395},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 531, col: 1, offset: 16452},
	expr: &actionExpr{
	pos: position{line: 531, col: 21, offset: 16474},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 531, col: 21, offset: 16474},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 531, col: 21, offset: 16474},
	name: "_",
},
&litMatcher{
	pos: position{line: 531, col: 23, offset: 16476},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 531, col: 27, offset: 16480},
	name: "_",
},
&labeledExpr{
	pos: position{line: 531, col: 29, offset: 16482},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 531, col: 31, offset: 16484},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 532, col: 1, offset: 16519},
	expr: &actionExpr{
	pos: position{line: 533, col: 7, offset: 16551},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 533, col: 7, offset: 16551},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 533, col: 7, offset: 16551},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 533, col: 13, offset: 16557},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 533, col: 32, offset: 16576},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 533, col: 37, offset: 16581},
	expr: &ruleRefExpr{
	pos: position{line: 533, col: 37, offset: 16581},
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
	pos: position{line: 547, col: 1, offset: 17171},
	expr: &choiceExpr{
	pos: position{line: 547, col: 13, offset: 17185},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 547, col: 13, offset: 17185},
	name: "NonEmptyUnionType",
},
&ruleRefExpr{
	pos: position{line: 547, col: 33, offset: 17205},
	name: "EmptyUnionType",
},
	},
},
},
{
	name: "EmptyUnionType",
	pos: position{line: 549, col: 1, offset: 17221},
	expr: &actionExpr{
	pos: position{line: 549, col: 18, offset: 17240},
	run: (*parser).callonEmptyUnionType1,
	expr: &litMatcher{
	pos: position{line: 549, col: 18, offset: 17240},
	val: "",
	ignoreCase: false,
},
},
},
{
	name: "NonEmptyUnionType",
	pos: position{line: 551, col: 1, offset: 17272},
	expr: &actionExpr{
	pos: position{line: 551, col: 21, offset: 17294},
	run: (*parser).callonNonEmptyUnionType1,
	expr: &seqExpr{
	pos: position{line: 551, col: 21, offset: 17294},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 551, col: 21, offset: 17294},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 551, col: 27, offset: 17300},
	name: "UnionVariant",
},
},
&labeledExpr{
	pos: position{line: 551, col: 40, offset: 17313},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 551, col: 45, offset: 17318},
	expr: &seqExpr{
	pos: position{line: 551, col: 46, offset: 17319},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 551, col: 46, offset: 17319},
	name: "_",
},
&litMatcher{
	pos: position{line: 551, col: 48, offset: 17321},
	val: "|",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 551, col: 52, offset: 17325},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 551, col: 54, offset: 17327},
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
	pos: position{line: 571, col: 1, offset: 18049},
	expr: &seqExpr{
	pos: position{line: 571, col: 16, offset: 18066},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 571, col: 16, offset: 18066},
	name: "AnyLabel",
},
&zeroOrOneExpr{
	pos: position{line: 571, col: 25, offset: 18075},
	expr: &seqExpr{
	pos: position{line: 571, col: 26, offset: 18076},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 571, col: 26, offset: 18076},
	name: "_",
},
&litMatcher{
	pos: position{line: 571, col: 28, offset: 18078},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 571, col: 32, offset: 18082},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 571, col: 35, offset: 18085},
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
	pos: position{line: 573, col: 1, offset: 18099},
	expr: &actionExpr{
	pos: position{line: 573, col: 12, offset: 18112},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 573, col: 12, offset: 18112},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 573, col: 12, offset: 18112},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 573, col: 16, offset: 18116},
	name: "_",
},
&labeledExpr{
	pos: position{line: 573, col: 18, offset: 18118},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 573, col: 20, offset: 18120},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 573, col: 31, offset: 18131},
	name: "_",
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 575, col: 1, offset: 18150},
	expr: &actionExpr{
	pos: position{line: 576, col: 7, offset: 18180},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 576, col: 7, offset: 18180},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 576, col: 7, offset: 18180},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 576, col: 11, offset: 18184},
	name: "_",
},
&labeledExpr{
	pos: position{line: 576, col: 13, offset: 18186},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 576, col: 19, offset: 18192},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 576, col: 30, offset: 18203},
	name: "_",
},
&labeledExpr{
	pos: position{line: 576, col: 32, offset: 18205},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 576, col: 37, offset: 18210},
	expr: &ruleRefExpr{
	pos: position{line: 576, col: 37, offset: 18210},
	name: "MoreList",
},
},
},
&litMatcher{
	pos: position{line: 576, col: 47, offset: 18220},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 586, col: 1, offset: 18496},
	expr: &notExpr{
	pos: position{line: 586, col: 7, offset: 18504},
	expr: &anyMatcher{
	line: 586, col: 8, offset: 18505,
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

func (c *current) onDoubleQuoteEscaped16() (interface{}, error) {
        i, err := strconv.ParseInt(string(c.text[1:]), 16, 32)
        return []byte(string([]rune{rune(i)})), err
     
}

func (p *parser) callonDoubleQuoteEscaped16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleQuoteEscaped16()
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

func (c *current) onPathComponent1(u interface{}) (interface{}, error) {
 return u, nil 
}

func (p *parser) callonPathComponent1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPathComponent1(stack["u"])
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
 return MakeRemote(u.(*url.URL)) 
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

func (c *current) onImportHashed1(i interface{}) (interface{}, error) {
 return ImportHashed{Fetchable: i.(Fetchable)}, nil 
}

func (p *parser) callonImportHashed1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onImportHashed1(stack["i"])
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
 return Embed(Import{ImportHashed: i.(ImportHashed), ImportMode: Code}), nil 
}

func (p *parser) callonImport10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onImport10(stack["i"])
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

func (c *current) onSelectorExpression1(e, ls interface{}) (interface{}, error) {
    expr := e.(Expr)
    labels := ls.([]interface{})
    for _, labelSelector := range labels {
        label := labelSelector.([]interface{})[3]
        expr = Field{expr, label.(string)}
    }
    return expr, nil
}

func (p *parser) callonSelectorExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSelectorExpression1(stack["e"], stack["ls"])
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
 return RecordLit(map[string]Expr{}), nil 
}

func (p *parser) callonRecordTypeOrLiteral2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRecordTypeOrLiteral2()
}

func (c *current) onRecordTypeOrLiteral6() (interface{}, error) {
 return Record(map[string]Expr{}), nil 
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
          content := make(map[string]Expr, len(fields)+1)
          content[first.([]interface{})[0].(string)] = first.([]interface{})[1].(Expr)
          for _, field := range(fields) {
              fieldName := field.([]interface{})[0].(string)
              if _, ok := content[fieldName]; ok {
                  return nil, fmt.Errorf("Duplicate field %s in record", fieldName)
              }
              content[fieldName] = field.([]interface{})[1].(Expr)
          }
          return Record(content), nil
      
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
          content := make(map[string]Expr, len(fields)+1)
          content[first.([]interface{})[0].(string)] = first.([]interface{})[1].(Expr)
          for _, field := range(fields) {
              fieldName := field.([]interface{})[0].(string)
              if _, ok := content[fieldName]; ok {
                  return nil, fmt.Errorf("Duplicate field %s in record", fieldName)
              }
              content[fieldName] = field.([]interface{})[1].(Expr)
          }
          return RecordLit(content), nil
      
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
    alternatives := make(map[string]Expr)
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
    return UnionType(alternatives), nil
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
          content := make([]Expr, len(exprs)+1)
          content[0] = first.(Expr)
          for i, expr := range(exprs) {
              content[i+1] = expr.(Expr)
          }
          return NonEmptyList(content), nil
      
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

