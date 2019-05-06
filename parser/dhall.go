
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
&litMatcher{
	pos: position{line: 40, col: 14, offset: 753},
	val: "\r\n",
	ignoreCase: false,
},
	},
},
},
{
	name: "BlockComment",
	pos: position{line: 42, col: 1, offset: 761},
	expr: &seqExpr{
	pos: position{line: 42, col: 16, offset: 778},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 42, col: 16, offset: 778},
	val: "{-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 42, col: 21, offset: 783},
	name: "BlockCommentContinue",
},
	},
},
},
{
	name: "BlockCommentChunk",
	pos: position{line: 44, col: 1, offset: 805},
	expr: &choiceExpr{
	pos: position{line: 45, col: 5, offset: 831},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 45, col: 5, offset: 831},
	name: "BlockComment",
},
&charClassMatcher{
	pos: position{line: 46, col: 5, offset: 848},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 47, col: 5, offset: 874},
	name: "EOL",
},
	},
},
},
{
	name: "BlockCommentContinue",
	pos: position{line: 49, col: 1, offset: 879},
	expr: &choiceExpr{
	pos: position{line: 49, col: 24, offset: 904},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 49, col: 24, offset: 904},
	val: "-}",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 49, col: 31, offset: 911},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 49, col: 31, offset: 911},
	name: "BlockCommentChunk",
},
&ruleRefExpr{
	pos: position{line: 49, col: 49, offset: 929},
	name: "BlockCommentContinue",
},
	},
},
	},
},
},
{
	name: "NotEOL",
	pos: position{line: 51, col: 1, offset: 951},
	expr: &charClassMatcher{
	pos: position{line: 51, col: 10, offset: 962},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "LineComment",
	pos: position{line: 53, col: 1, offset: 985},
	expr: &actionExpr{
	pos: position{line: 53, col: 15, offset: 1001},
	run: (*parser).callonLineComment1,
	expr: &seqExpr{
	pos: position{line: 53, col: 15, offset: 1001},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 53, col: 15, offset: 1001},
	val: "--",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 53, col: 20, offset: 1006},
	label: "content",
	expr: &actionExpr{
	pos: position{line: 53, col: 29, offset: 1015},
	run: (*parser).callonLineComment5,
	expr: &zeroOrMoreExpr{
	pos: position{line: 53, col: 29, offset: 1015},
	expr: &ruleRefExpr{
	pos: position{line: 53, col: 29, offset: 1015},
	name: "NotEOL",
},
},
},
},
&ruleRefExpr{
	pos: position{line: 53, col: 68, offset: 1054},
	name: "EOL",
},
	},
},
},
},
{
	name: "WhitespaceChunk",
	pos: position{line: 55, col: 1, offset: 1083},
	expr: &choiceExpr{
	pos: position{line: 55, col: 19, offset: 1103},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 55, col: 19, offset: 1103},
	val: " ",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 55, col: 25, offset: 1109},
	val: "\t",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 55, col: 32, offset: 1116},
	name: "EOL",
},
&ruleRefExpr{
	pos: position{line: 55, col: 38, offset: 1122},
	name: "LineComment",
},
&ruleRefExpr{
	pos: position{line: 55, col: 52, offset: 1136},
	name: "BlockComment",
},
	},
},
},
{
	name: "_",
	pos: position{line: 57, col: 1, offset: 1150},
	expr: &zeroOrMoreExpr{
	pos: position{line: 57, col: 5, offset: 1156},
	expr: &ruleRefExpr{
	pos: position{line: 57, col: 5, offset: 1156},
	name: "WhitespaceChunk",
},
},
},
{
	name: "_1",
	pos: position{line: 59, col: 1, offset: 1174},
	expr: &oneOrMoreExpr{
	pos: position{line: 59, col: 6, offset: 1181},
	expr: &ruleRefExpr{
	pos: position{line: 59, col: 6, offset: 1181},
	name: "WhitespaceChunk",
},
},
},
{
	name: "Digit",
	pos: position{line: 61, col: 1, offset: 1199},
	expr: &charClassMatcher{
	pos: position{line: 61, col: 9, offset: 1209},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "HexDig",
	pos: position{line: 63, col: 1, offset: 1216},
	expr: &choiceExpr{
	pos: position{line: 63, col: 10, offset: 1227},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 63, col: 10, offset: 1227},
	name: "Digit",
},
&charClassMatcher{
	pos: position{line: 63, col: 18, offset: 1235},
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
	pos: position{line: 65, col: 1, offset: 1243},
	expr: &charClassMatcher{
	pos: position{line: 65, col: 24, offset: 1268},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SimpleLabelNextChar",
	pos: position{line: 66, col: 1, offset: 1278},
	expr: &charClassMatcher{
	pos: position{line: 66, col: 23, offset: 1302},
	val: "[A-Za-z0-9_/-]",
	chars: []rune{'_','/','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SimpleLabel",
	pos: position{line: 67, col: 1, offset: 1317},
	expr: &choiceExpr{
	pos: position{line: 67, col: 15, offset: 1333},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 67, col: 15, offset: 1333},
	run: (*parser).callonSimpleLabel2,
	expr: &seqExpr{
	pos: position{line: 67, col: 15, offset: 1333},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 67, col: 15, offset: 1333},
	name: "Keyword",
},
&oneOrMoreExpr{
	pos: position{line: 67, col: 23, offset: 1341},
	expr: &ruleRefExpr{
	pos: position{line: 67, col: 23, offset: 1341},
	name: "SimpleLabelNextChar",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 68, col: 13, offset: 1405},
	run: (*parser).callonSimpleLabel7,
	expr: &seqExpr{
	pos: position{line: 68, col: 13, offset: 1405},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 68, col: 13, offset: 1405},
	expr: &ruleRefExpr{
	pos: position{line: 68, col: 14, offset: 1406},
	name: "Keyword",
},
},
&ruleRefExpr{
	pos: position{line: 68, col: 22, offset: 1414},
	name: "SimpleLabelFirstChar",
},
&zeroOrMoreExpr{
	pos: position{line: 68, col: 43, offset: 1435},
	expr: &ruleRefExpr{
	pos: position{line: 68, col: 43, offset: 1435},
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
	name: "Label",
	pos: position{line: 75, col: 1, offset: 1536},
	expr: &actionExpr{
	pos: position{line: 75, col: 9, offset: 1546},
	run: (*parser).callonLabel1,
	expr: &labeledExpr{
	pos: position{line: 75, col: 9, offset: 1546},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 75, col: 15, offset: 1552},
	name: "SimpleLabel",
},
},
},
},
{
	name: "NonreservedLabel",
	pos: position{line: 77, col: 1, offset: 1587},
	expr: &choiceExpr{
	pos: position{line: 77, col: 20, offset: 1608},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 77, col: 20, offset: 1608},
	run: (*parser).callonNonreservedLabel2,
	expr: &seqExpr{
	pos: position{line: 77, col: 20, offset: 1608},
	exprs: []interface{}{
&andExpr{
	pos: position{line: 77, col: 20, offset: 1608},
	expr: &seqExpr{
	pos: position{line: 77, col: 22, offset: 1610},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 77, col: 22, offset: 1610},
	name: "Reserved",
},
&ruleRefExpr{
	pos: position{line: 77, col: 31, offset: 1619},
	name: "SimpleLabelNextChar",
},
	},
},
},
&labeledExpr{
	pos: position{line: 77, col: 52, offset: 1640},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 77, col: 58, offset: 1646},
	name: "Label",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 78, col: 19, offset: 1692},
	run: (*parser).callonNonreservedLabel10,
	expr: &seqExpr{
	pos: position{line: 78, col: 19, offset: 1692},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 78, col: 19, offset: 1692},
	expr: &ruleRefExpr{
	pos: position{line: 78, col: 20, offset: 1693},
	name: "Reserved",
},
},
&labeledExpr{
	pos: position{line: 78, col: 29, offset: 1702},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 78, col: 35, offset: 1708},
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
	pos: position{line: 80, col: 1, offset: 1737},
	expr: &ruleRefExpr{
	pos: position{line: 80, col: 12, offset: 1750},
	name: "Label",
},
},
{
	name: "DoubleQuoteChunk",
	pos: position{line: 83, col: 1, offset: 1758},
	expr: &choiceExpr{
	pos: position{line: 84, col: 6, offset: 1784},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 84, col: 6, offset: 1784},
	name: "Interpolation",
},
&actionExpr{
	pos: position{line: 85, col: 6, offset: 1803},
	run: (*parser).callonDoubleQuoteChunk3,
	expr: &seqExpr{
	pos: position{line: 85, col: 6, offset: 1803},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 85, col: 6, offset: 1803},
	val: "\\",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 85, col: 11, offset: 1808},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 85, col: 13, offset: 1810},
	name: "DoubleQuoteEscaped",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 86, col: 6, offset: 1852},
	name: "DoubleQuoteChar",
},
	},
},
},
{
	name: "DoubleQuoteEscaped",
	pos: position{line: 88, col: 1, offset: 1869},
	expr: &choiceExpr{
	pos: position{line: 89, col: 8, offset: 1899},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 89, col: 8, offset: 1899},
	val: "\"",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 90, col: 8, offset: 1910},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 91, col: 8, offset: 1921},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 92, col: 8, offset: 1933},
	val: "/",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 93, col: 8, offset: 1944},
	run: (*parser).callonDoubleQuoteEscaped6,
	expr: &litMatcher{
	pos: position{line: 93, col: 8, offset: 1944},
	val: "b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 94, col: 8, offset: 1984},
	run: (*parser).callonDoubleQuoteEscaped8,
	expr: &litMatcher{
	pos: position{line: 94, col: 8, offset: 1984},
	val: "f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 95, col: 8, offset: 2024},
	run: (*parser).callonDoubleQuoteEscaped10,
	expr: &litMatcher{
	pos: position{line: 95, col: 8, offset: 2024},
	val: "n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 96, col: 8, offset: 2064},
	run: (*parser).callonDoubleQuoteEscaped12,
	expr: &litMatcher{
	pos: position{line: 96, col: 8, offset: 2064},
	val: "r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 97, col: 8, offset: 2104},
	run: (*parser).callonDoubleQuoteEscaped14,
	expr: &litMatcher{
	pos: position{line: 97, col: 8, offset: 2104},
	val: "t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 98, col: 8, offset: 2144},
	run: (*parser).callonDoubleQuoteEscaped16,
	expr: &seqExpr{
	pos: position{line: 98, col: 8, offset: 2144},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 98, col: 8, offset: 2144},
	val: "u",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 98, col: 12, offset: 2148},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 98, col: 19, offset: 2155},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 98, col: 26, offset: 2162},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 98, col: 33, offset: 2169},
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
	pos: position{line: 103, col: 1, offset: 2301},
	expr: &choiceExpr{
	pos: position{line: 104, col: 6, offset: 2326},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 104, col: 6, offset: 2326},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 105, col: 6, offset: 2343},
	val: "[\\x23-\\x5b]",
	ranges: []rune{'#','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 106, col: 6, offset: 2360},
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
	pos: position{line: 108, col: 1, offset: 2379},
	expr: &actionExpr{
	pos: position{line: 108, col: 22, offset: 2402},
	run: (*parser).callonDoubleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 108, col: 22, offset: 2402},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 108, col: 22, offset: 2402},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 108, col: 26, offset: 2406},
	label: "chunks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 108, col: 33, offset: 2413},
	expr: &ruleRefExpr{
	pos: position{line: 108, col: 33, offset: 2413},
	name: "DoubleQuoteChunk",
},
},
},
&litMatcher{
	pos: position{line: 108, col: 51, offset: 2431},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "SingleQuoteContinue",
	pos: position{line: 125, col: 1, offset: 2899},
	expr: &choiceExpr{
	pos: position{line: 126, col: 7, offset: 2929},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 126, col: 7, offset: 2929},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 126, col: 7, offset: 2929},
	name: "Interpolation",
},
&ruleRefExpr{
	pos: position{line: 126, col: 21, offset: 2943},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 127, col: 7, offset: 2969},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 127, col: 7, offset: 2969},
	name: "EscapedQuotePair",
},
&ruleRefExpr{
	pos: position{line: 127, col: 24, offset: 2986},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 128, col: 7, offset: 3012},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 128, col: 7, offset: 3012},
	name: "EscapedInterpolation",
},
&ruleRefExpr{
	pos: position{line: 128, col: 28, offset: 3033},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 129, col: 7, offset: 3059},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 129, col: 7, offset: 3059},
	name: "SingleQuoteChar",
},
&ruleRefExpr{
	pos: position{line: 129, col: 23, offset: 3075},
	name: "SingleQuoteContinue",
},
	},
},
&litMatcher{
	pos: position{line: 130, col: 7, offset: 3101},
	val: "''",
	ignoreCase: false,
},
	},
},
},
{
	name: "EscapedQuotePair",
	pos: position{line: 132, col: 1, offset: 3107},
	expr: &actionExpr{
	pos: position{line: 132, col: 20, offset: 3128},
	run: (*parser).callonEscapedQuotePair1,
	expr: &litMatcher{
	pos: position{line: 132, col: 20, offset: 3128},
	val: "'''",
	ignoreCase: false,
},
},
},
{
	name: "EscapedInterpolation",
	pos: position{line: 136, col: 1, offset: 3263},
	expr: &actionExpr{
	pos: position{line: 136, col: 24, offset: 3288},
	run: (*parser).callonEscapedInterpolation1,
	expr: &litMatcher{
	pos: position{line: 136, col: 24, offset: 3288},
	val: "''${",
	ignoreCase: false,
},
},
},
{
	name: "SingleQuoteChar",
	pos: position{line: 138, col: 1, offset: 3330},
	expr: &choiceExpr{
	pos: position{line: 139, col: 6, offset: 3355},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 139, col: 6, offset: 3355},
	val: "[\\x20-\\U0010ffff]",
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 140, col: 6, offset: 3378},
	val: "\t",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 141, col: 6, offset: 3388},
	val: "\n",
	ignoreCase: false,
},
	},
},
},
{
	name: "SingleQuoteLiteral",
	pos: position{line: 143, col: 1, offset: 3394},
	expr: &actionExpr{
	pos: position{line: 143, col: 22, offset: 3417},
	run: (*parser).callonSingleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 143, col: 22, offset: 3417},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 143, col: 22, offset: 3417},
	val: "''",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 143, col: 27, offset: 3422},
	name: "EOL",
},
&labeledExpr{
	pos: position{line: 143, col: 31, offset: 3426},
	label: "content",
	expr: &ruleRefExpr{
	pos: position{line: 143, col: 39, offset: 3434},
	name: "SingleQuoteContinue",
},
},
	},
},
},
},
{
	name: "Interpolation",
	pos: position{line: 161, col: 1, offset: 3984},
	expr: &actionExpr{
	pos: position{line: 161, col: 17, offset: 4002},
	run: (*parser).callonInterpolation1,
	expr: &seqExpr{
	pos: position{line: 161, col: 17, offset: 4002},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 161, col: 17, offset: 4002},
	val: "${",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 161, col: 22, offset: 4007},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 161, col: 24, offset: 4009},
	name: "CompleteExpression",
},
},
&litMatcher{
	pos: position{line: 161, col: 43, offset: 4028},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TextLiteral",
	pos: position{line: 163, col: 1, offset: 4051},
	expr: &choiceExpr{
	pos: position{line: 163, col: 15, offset: 4067},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 163, col: 15, offset: 4067},
	name: "DoubleQuoteLiteral",
},
&ruleRefExpr{
	pos: position{line: 163, col: 36, offset: 4088},
	name: "SingleQuoteLiteral",
},
	},
},
},
{
	name: "Reserved",
	pos: position{line: 166, col: 1, offset: 4193},
	expr: &choiceExpr{
	pos: position{line: 167, col: 5, offset: 4210},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 167, col: 5, offset: 4210},
	run: (*parser).callonReserved2,
	expr: &litMatcher{
	pos: position{line: 167, col: 5, offset: 4210},
	val: "Natural/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 168, col: 5, offset: 4288},
	run: (*parser).callonReserved4,
	expr: &litMatcher{
	pos: position{line: 168, col: 5, offset: 4288},
	val: "Natural/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 169, col: 5, offset: 4364},
	run: (*parser).callonReserved6,
	expr: &litMatcher{
	pos: position{line: 169, col: 5, offset: 4364},
	val: "Natural/isZero",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 170, col: 5, offset: 4444},
	run: (*parser).callonReserved8,
	expr: &litMatcher{
	pos: position{line: 170, col: 5, offset: 4444},
	val: "Natural/even",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 171, col: 5, offset: 4520},
	run: (*parser).callonReserved10,
	expr: &litMatcher{
	pos: position{line: 171, col: 5, offset: 4520},
	val: "Natural/odd",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 172, col: 5, offset: 4594},
	run: (*parser).callonReserved12,
	expr: &litMatcher{
	pos: position{line: 172, col: 5, offset: 4594},
	val: "Natural/toInteger",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 173, col: 5, offset: 4680},
	run: (*parser).callonReserved14,
	expr: &litMatcher{
	pos: position{line: 173, col: 5, offset: 4680},
	val: "Natural/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 174, col: 5, offset: 4756},
	run: (*parser).callonReserved16,
	expr: &litMatcher{
	pos: position{line: 174, col: 5, offset: 4756},
	val: "Integer/toDouble",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 175, col: 5, offset: 4840},
	run: (*parser).callonReserved18,
	expr: &litMatcher{
	pos: position{line: 175, col: 5, offset: 4840},
	val: "Integer/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 176, col: 5, offset: 4916},
	run: (*parser).callonReserved20,
	expr: &litMatcher{
	pos: position{line: 176, col: 5, offset: 4916},
	val: "Double/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 177, col: 5, offset: 4990},
	run: (*parser).callonReserved22,
	expr: &litMatcher{
	pos: position{line: 177, col: 5, offset: 4990},
	val: "List/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 178, col: 5, offset: 5062},
	run: (*parser).callonReserved24,
	expr: &litMatcher{
	pos: position{line: 178, col: 5, offset: 5062},
	val: "List/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 179, col: 5, offset: 5132},
	run: (*parser).callonReserved26,
	expr: &litMatcher{
	pos: position{line: 179, col: 5, offset: 5132},
	val: "List/length",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 180, col: 5, offset: 5206},
	run: (*parser).callonReserved28,
	expr: &litMatcher{
	pos: position{line: 180, col: 5, offset: 5206},
	val: "List/head",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 181, col: 5, offset: 5276},
	run: (*parser).callonReserved30,
	expr: &litMatcher{
	pos: position{line: 181, col: 5, offset: 5276},
	val: "List/last",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 182, col: 5, offset: 5346},
	run: (*parser).callonReserved32,
	expr: &litMatcher{
	pos: position{line: 182, col: 5, offset: 5346},
	val: "List/indexed",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 183, col: 5, offset: 5422},
	run: (*parser).callonReserved34,
	expr: &litMatcher{
	pos: position{line: 183, col: 5, offset: 5422},
	val: "List/reverse",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 184, col: 5, offset: 5498},
	run: (*parser).callonReserved36,
	expr: &litMatcher{
	pos: position{line: 184, col: 5, offset: 5498},
	val: "Optional/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 185, col: 5, offset: 5578},
	run: (*parser).callonReserved38,
	expr: &litMatcher{
	pos: position{line: 185, col: 5, offset: 5578},
	val: "Optional/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 186, col: 5, offset: 5656},
	run: (*parser).callonReserved40,
	expr: &litMatcher{
	pos: position{line: 186, col: 5, offset: 5656},
	val: "Text/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 187, col: 5, offset: 5726},
	run: (*parser).callonReserved42,
	expr: &litMatcher{
	pos: position{line: 187, col: 5, offset: 5726},
	val: "Bool",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 188, col: 5, offset: 5758},
	run: (*parser).callonReserved44,
	expr: &litMatcher{
	pos: position{line: 188, col: 5, offset: 5758},
	val: "True",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 189, col: 5, offset: 5790},
	run: (*parser).callonReserved46,
	expr: &litMatcher{
	pos: position{line: 189, col: 5, offset: 5790},
	val: "False",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 190, col: 5, offset: 5824},
	run: (*parser).callonReserved48,
	expr: &litMatcher{
	pos: position{line: 190, col: 5, offset: 5824},
	val: "Optional",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 191, col: 5, offset: 5864},
	run: (*parser).callonReserved50,
	expr: &litMatcher{
	pos: position{line: 191, col: 5, offset: 5864},
	val: "Natural",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 192, col: 5, offset: 5902},
	run: (*parser).callonReserved52,
	expr: &litMatcher{
	pos: position{line: 192, col: 5, offset: 5902},
	val: "Integer",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 193, col: 5, offset: 5940},
	run: (*parser).callonReserved54,
	expr: &litMatcher{
	pos: position{line: 193, col: 5, offset: 5940},
	val: "Double",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 194, col: 5, offset: 5976},
	run: (*parser).callonReserved56,
	expr: &litMatcher{
	pos: position{line: 194, col: 5, offset: 5976},
	val: "Text",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 195, col: 5, offset: 6008},
	run: (*parser).callonReserved58,
	expr: &litMatcher{
	pos: position{line: 195, col: 5, offset: 6008},
	val: "List",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 196, col: 5, offset: 6040},
	run: (*parser).callonReserved60,
	expr: &litMatcher{
	pos: position{line: 196, col: 5, offset: 6040},
	val: "None",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 197, col: 5, offset: 6072},
	run: (*parser).callonReserved62,
	expr: &litMatcher{
	pos: position{line: 197, col: 5, offset: 6072},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 198, col: 5, offset: 6104},
	run: (*parser).callonReserved64,
	expr: &litMatcher{
	pos: position{line: 198, col: 5, offset: 6104},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 199, col: 5, offset: 6136},
	run: (*parser).callonReserved66,
	expr: &litMatcher{
	pos: position{line: 199, col: 5, offset: 6136},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "If",
	pos: position{line: 201, col: 1, offset: 6165},
	expr: &litMatcher{
	pos: position{line: 201, col: 6, offset: 6172},
	val: "if",
	ignoreCase: false,
},
},
{
	name: "Then",
	pos: position{line: 202, col: 1, offset: 6177},
	expr: &litMatcher{
	pos: position{line: 202, col: 8, offset: 6186},
	val: "then",
	ignoreCase: false,
},
},
{
	name: "Else",
	pos: position{line: 203, col: 1, offset: 6193},
	expr: &litMatcher{
	pos: position{line: 203, col: 8, offset: 6202},
	val: "else",
	ignoreCase: false,
},
},
{
	name: "Let",
	pos: position{line: 204, col: 1, offset: 6209},
	expr: &litMatcher{
	pos: position{line: 204, col: 7, offset: 6217},
	val: "let",
	ignoreCase: false,
},
},
{
	name: "In",
	pos: position{line: 205, col: 1, offset: 6223},
	expr: &litMatcher{
	pos: position{line: 205, col: 6, offset: 6230},
	val: "in",
	ignoreCase: false,
},
},
{
	name: "As",
	pos: position{line: 206, col: 1, offset: 6235},
	expr: &litMatcher{
	pos: position{line: 206, col: 6, offset: 6242},
	val: "as",
	ignoreCase: false,
},
},
{
	name: "Using",
	pos: position{line: 207, col: 1, offset: 6247},
	expr: &litMatcher{
	pos: position{line: 207, col: 9, offset: 6257},
	val: "using",
	ignoreCase: false,
},
},
{
	name: "Merge",
	pos: position{line: 208, col: 1, offset: 6265},
	expr: &litMatcher{
	pos: position{line: 208, col: 9, offset: 6275},
	val: "merge",
	ignoreCase: false,
},
},
{
	name: "Missing",
	pos: position{line: 209, col: 1, offset: 6283},
	expr: &actionExpr{
	pos: position{line: 209, col: 11, offset: 6295},
	run: (*parser).callonMissing1,
	expr: &litMatcher{
	pos: position{line: 209, col: 11, offset: 6295},
	val: "missing",
	ignoreCase: false,
},
},
},
{
	name: "True",
	pos: position{line: 210, col: 1, offset: 6331},
	expr: &litMatcher{
	pos: position{line: 210, col: 8, offset: 6340},
	val: "True",
	ignoreCase: false,
},
},
{
	name: "False",
	pos: position{line: 211, col: 1, offset: 6347},
	expr: &litMatcher{
	pos: position{line: 211, col: 9, offset: 6357},
	val: "False",
	ignoreCase: false,
},
},
{
	name: "Infinity",
	pos: position{line: 212, col: 1, offset: 6365},
	expr: &litMatcher{
	pos: position{line: 212, col: 12, offset: 6378},
	val: "Infinity",
	ignoreCase: false,
},
},
{
	name: "NaN",
	pos: position{line: 213, col: 1, offset: 6389},
	expr: &litMatcher{
	pos: position{line: 213, col: 7, offset: 6397},
	val: "NaN",
	ignoreCase: false,
},
},
{
	name: "Some",
	pos: position{line: 214, col: 1, offset: 6403},
	expr: &litMatcher{
	pos: position{line: 214, col: 8, offset: 6412},
	val: "Some",
	ignoreCase: false,
},
},
{
	name: "Keyword",
	pos: position{line: 216, col: 1, offset: 6420},
	expr: &choiceExpr{
	pos: position{line: 217, col: 5, offset: 6436},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 217, col: 5, offset: 6436},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 217, col: 10, offset: 6441},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 217, col: 17, offset: 6448},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 218, col: 5, offset: 6457},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 218, col: 11, offset: 6463},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 219, col: 5, offset: 6470},
	name: "Using",
},
&ruleRefExpr{
	pos: position{line: 219, col: 13, offset: 6478},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 219, col: 23, offset: 6488},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 220, col: 5, offset: 6495},
	name: "True",
},
&ruleRefExpr{
	pos: position{line: 220, col: 12, offset: 6502},
	name: "False",
},
&ruleRefExpr{
	pos: position{line: 221, col: 5, offset: 6512},
	name: "Infinity",
},
&ruleRefExpr{
	pos: position{line: 221, col: 16, offset: 6523},
	name: "NaN",
},
&ruleRefExpr{
	pos: position{line: 222, col: 5, offset: 6531},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 222, col: 13, offset: 6539},
	name: "Some",
},
	},
},
},
{
	name: "Optional",
	pos: position{line: 224, col: 1, offset: 6545},
	expr: &litMatcher{
	pos: position{line: 224, col: 12, offset: 6558},
	val: "Optional",
	ignoreCase: false,
},
},
{
	name: "Text",
	pos: position{line: 225, col: 1, offset: 6569},
	expr: &litMatcher{
	pos: position{line: 225, col: 8, offset: 6578},
	val: "Text",
	ignoreCase: false,
},
},
{
	name: "List",
	pos: position{line: 226, col: 1, offset: 6585},
	expr: &litMatcher{
	pos: position{line: 226, col: 8, offset: 6594},
	val: "List",
	ignoreCase: false,
},
},
{
	name: "Lambda",
	pos: position{line: 228, col: 1, offset: 6602},
	expr: &choiceExpr{
	pos: position{line: 228, col: 11, offset: 6614},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 228, col: 11, offset: 6614},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 228, col: 18, offset: 6621},
	val: "λ",
	ignoreCase: false,
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 229, col: 1, offset: 6627},
	expr: &choiceExpr{
	pos: position{line: 229, col: 11, offset: 6639},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 229, col: 11, offset: 6639},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 229, col: 22, offset: 6650},
	val: "∀",
	ignoreCase: false,
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 230, col: 1, offset: 6657},
	expr: &choiceExpr{
	pos: position{line: 230, col: 10, offset: 6668},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 230, col: 10, offset: 6668},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 230, col: 17, offset: 6675},
	val: "→",
	ignoreCase: false,
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 232, col: 1, offset: 6683},
	expr: &seqExpr{
	pos: position{line: 232, col: 12, offset: 6696},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 232, col: 12, offset: 6696},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 232, col: 17, offset: 6701},
	expr: &charClassMatcher{
	pos: position{line: 232, col: 17, offset: 6701},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 232, col: 23, offset: 6707},
	expr: &ruleRefExpr{
	pos: position{line: 232, col: 23, offset: 6707},
	name: "Digit",
},
},
	},
},
},
{
	name: "NumericDoubleLiteral",
	pos: position{line: 234, col: 1, offset: 6715},
	expr: &actionExpr{
	pos: position{line: 234, col: 24, offset: 6740},
	run: (*parser).callonNumericDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 234, col: 24, offset: 6740},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 234, col: 24, offset: 6740},
	expr: &charClassMatcher{
	pos: position{line: 234, col: 24, offset: 6740},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 234, col: 30, offset: 6746},
	expr: &ruleRefExpr{
	pos: position{line: 234, col: 30, offset: 6746},
	name: "Digit",
},
},
&choiceExpr{
	pos: position{line: 234, col: 39, offset: 6755},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 234, col: 39, offset: 6755},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 234, col: 39, offset: 6755},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 234, col: 43, offset: 6759},
	expr: &ruleRefExpr{
	pos: position{line: 234, col: 43, offset: 6759},
	name: "Digit",
},
},
&zeroOrOneExpr{
	pos: position{line: 234, col: 50, offset: 6766},
	expr: &ruleRefExpr{
	pos: position{line: 234, col: 50, offset: 6766},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 234, col: 62, offset: 6778},
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
	pos: position{line: 242, col: 1, offset: 6934},
	expr: &choiceExpr{
	pos: position{line: 242, col: 17, offset: 6952},
	alternatives: []interface{}{
&labeledExpr{
	pos: position{line: 242, col: 17, offset: 6952},
	label: "d",
	expr: &ruleRefExpr{
	pos: position{line: 242, col: 19, offset: 6954},
	name: "NumericDoubleLiteral",
},
},
&actionExpr{
	pos: position{line: 243, col: 5, offset: 6979},
	run: (*parser).callonDoubleLiteral4,
	expr: &ruleRefExpr{
	pos: position{line: 243, col: 5, offset: 6979},
	name: "Infinity",
},
},
&actionExpr{
	pos: position{line: 244, col: 5, offset: 7031},
	run: (*parser).callonDoubleLiteral6,
	expr: &seqExpr{
	pos: position{line: 244, col: 5, offset: 7031},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 244, col: 5, offset: 7031},
	val: "-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 244, col: 9, offset: 7035},
	name: "Infinity",
},
	},
},
},
&actionExpr{
	pos: position{line: 245, col: 5, offset: 7088},
	run: (*parser).callonDoubleLiteral10,
	expr: &ruleRefExpr{
	pos: position{line: 245, col: 5, offset: 7088},
	name: "NaN",
},
},
	},
},
},
{
	name: "NaturalLiteral",
	pos: position{line: 247, col: 1, offset: 7131},
	expr: &actionExpr{
	pos: position{line: 247, col: 18, offset: 7150},
	run: (*parser).callonNaturalLiteral1,
	expr: &oneOrMoreExpr{
	pos: position{line: 247, col: 18, offset: 7150},
	expr: &ruleRefExpr{
	pos: position{line: 247, col: 18, offset: 7150},
	name: "Digit",
},
},
},
},
{
	name: "IntegerLiteral",
	pos: position{line: 252, col: 1, offset: 7239},
	expr: &actionExpr{
	pos: position{line: 252, col: 18, offset: 7258},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 252, col: 18, offset: 7258},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 252, col: 18, offset: 7258},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 252, col: 22, offset: 7262},
	name: "NaturalLiteral",
},
	},
},
},
},
{
	name: "DeBruijn",
	pos: position{line: 260, col: 1, offset: 7414},
	expr: &actionExpr{
	pos: position{line: 260, col: 12, offset: 7427},
	run: (*parser).callonDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 260, col: 12, offset: 7427},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 260, col: 12, offset: 7427},
	name: "_",
},
&litMatcher{
	pos: position{line: 260, col: 14, offset: 7429},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 260, col: 18, offset: 7433},
	name: "_",
},
&labeledExpr{
	pos: position{line: 260, col: 20, offset: 7435},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 260, col: 26, offset: 7441},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Variable",
	pos: position{line: 262, col: 1, offset: 7497},
	expr: &actionExpr{
	pos: position{line: 262, col: 12, offset: 7510},
	run: (*parser).callonVariable1,
	expr: &seqExpr{
	pos: position{line: 262, col: 12, offset: 7510},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 262, col: 12, offset: 7510},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 262, col: 17, offset: 7515},
	name: "NonreservedLabel",
},
},
&labeledExpr{
	pos: position{line: 262, col: 34, offset: 7532},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 262, col: 40, offset: 7538},
	expr: &ruleRefExpr{
	pos: position{line: 262, col: 40, offset: 7538},
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
	pos: position{line: 270, col: 1, offset: 7701},
	expr: &choiceExpr{
	pos: position{line: 270, col: 14, offset: 7716},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 270, col: 14, offset: 7716},
	name: "Variable",
},
&ruleRefExpr{
	pos: position{line: 270, col: 25, offset: 7727},
	name: "Reserved",
},
	},
},
},
{
	name: "PathCharacter",
	pos: position{line: 272, col: 1, offset: 7737},
	expr: &choiceExpr{
	pos: position{line: 273, col: 6, offset: 7760},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 273, col: 6, offset: 7760},
	val: "!",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 274, col: 6, offset: 7772},
	val: "[\\x24-\\x27]",
	ranges: []rune{'$','\'',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 275, col: 6, offset: 7789},
	val: "[\\x2a-\\x2b]",
	ranges: []rune{'*','+',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 276, col: 6, offset: 7806},
	val: "[\\x2d-\\x2e]",
	ranges: []rune{'-','.',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 277, col: 6, offset: 7823},
	val: "[\\x30-\\x3b]",
	ranges: []rune{'0',';',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 278, col: 6, offset: 7840},
	val: "=",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 279, col: 6, offset: 7852},
	val: "[\\x40-\\x5a]",
	ranges: []rune{'@','Z',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 280, col: 6, offset: 7869},
	val: "[\\x5e-\\x7a]",
	ranges: []rune{'^','z',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 281, col: 6, offset: 7886},
	val: "|",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 282, col: 6, offset: 7898},
	val: "~",
	ignoreCase: false,
},
	},
},
},
{
	name: "UnquotedPathComponent",
	pos: position{line: 284, col: 1, offset: 7906},
	expr: &actionExpr{
	pos: position{line: 284, col: 25, offset: 7932},
	run: (*parser).callonUnquotedPathComponent1,
	expr: &oneOrMoreExpr{
	pos: position{line: 284, col: 25, offset: 7932},
	expr: &ruleRefExpr{
	pos: position{line: 284, col: 25, offset: 7932},
	name: "PathCharacter",
},
},
},
},
{
	name: "PathComponent",
	pos: position{line: 286, col: 1, offset: 7979},
	expr: &actionExpr{
	pos: position{line: 286, col: 17, offset: 7997},
	run: (*parser).callonPathComponent1,
	expr: &seqExpr{
	pos: position{line: 286, col: 17, offset: 7997},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 286, col: 17, offset: 7997},
	val: "/",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 286, col: 21, offset: 8001},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 286, col: 23, offset: 8003},
	name: "UnquotedPathComponent",
},
},
	},
},
},
},
{
	name: "Path",
	pos: position{line: 288, col: 1, offset: 8044},
	expr: &actionExpr{
	pos: position{line: 288, col: 8, offset: 8053},
	run: (*parser).callonPath1,
	expr: &labeledExpr{
	pos: position{line: 288, col: 8, offset: 8053},
	label: "cs",
	expr: &oneOrMoreExpr{
	pos: position{line: 288, col: 11, offset: 8056},
	expr: &ruleRefExpr{
	pos: position{line: 288, col: 11, offset: 8056},
	name: "PathComponent",
},
},
},
},
},
{
	name: "Local",
	pos: position{line: 297, col: 1, offset: 8330},
	expr: &choiceExpr{
	pos: position{line: 297, col: 9, offset: 8340},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 297, col: 9, offset: 8340},
	name: "ParentPath",
},
&ruleRefExpr{
	pos: position{line: 297, col: 22, offset: 8353},
	name: "HerePath",
},
&ruleRefExpr{
	pos: position{line: 297, col: 33, offset: 8364},
	name: "HomePath",
},
&ruleRefExpr{
	pos: position{line: 297, col: 44, offset: 8375},
	name: "AbsolutePath",
},
	},
},
},
{
	name: "ParentPath",
	pos: position{line: 299, col: 1, offset: 8389},
	expr: &actionExpr{
	pos: position{line: 299, col: 14, offset: 8404},
	run: (*parser).callonParentPath1,
	expr: &seqExpr{
	pos: position{line: 299, col: 14, offset: 8404},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 299, col: 14, offset: 8404},
	val: "..",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 299, col: 19, offset: 8409},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 299, col: 21, offset: 8411},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HerePath",
	pos: position{line: 300, col: 1, offset: 8467},
	expr: &actionExpr{
	pos: position{line: 300, col: 12, offset: 8480},
	run: (*parser).callonHerePath1,
	expr: &seqExpr{
	pos: position{line: 300, col: 12, offset: 8480},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 300, col: 12, offset: 8480},
	val: ".",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 300, col: 16, offset: 8484},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 300, col: 18, offset: 8486},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HomePath",
	pos: position{line: 301, col: 1, offset: 8525},
	expr: &actionExpr{
	pos: position{line: 301, col: 12, offset: 8538},
	run: (*parser).callonHomePath1,
	expr: &seqExpr{
	pos: position{line: 301, col: 12, offset: 8538},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 301, col: 12, offset: 8538},
	val: "~",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 301, col: 16, offset: 8542},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 301, col: 18, offset: 8544},
	name: "Path",
},
},
	},
},
},
},
{
	name: "AbsolutePath",
	pos: position{line: 302, col: 1, offset: 8599},
	expr: &actionExpr{
	pos: position{line: 302, col: 16, offset: 8616},
	run: (*parser).callonAbsolutePath1,
	expr: &labeledExpr{
	pos: position{line: 302, col: 16, offset: 8616},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 302, col: 18, offset: 8618},
	name: "Path",
},
},
},
},
{
	name: "Scheme",
	pos: position{line: 304, col: 1, offset: 8674},
	expr: &seqExpr{
	pos: position{line: 304, col: 10, offset: 8685},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 304, col: 10, offset: 8685},
	val: "http",
	ignoreCase: false,
},
&zeroOrOneExpr{
	pos: position{line: 304, col: 17, offset: 8692},
	expr: &litMatcher{
	pos: position{line: 304, col: 17, offset: 8692},
	val: "s",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "HttpRaw",
	pos: position{line: 306, col: 1, offset: 8698},
	expr: &actionExpr{
	pos: position{line: 306, col: 11, offset: 8710},
	run: (*parser).callonHttpRaw1,
	expr: &seqExpr{
	pos: position{line: 306, col: 11, offset: 8710},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 306, col: 11, offset: 8710},
	name: "Scheme",
},
&litMatcher{
	pos: position{line: 306, col: 18, offset: 8717},
	val: "://",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 306, col: 24, offset: 8723},
	name: "Authority",
},
&ruleRefExpr{
	pos: position{line: 306, col: 34, offset: 8733},
	name: "Path",
},
&zeroOrOneExpr{
	pos: position{line: 306, col: 39, offset: 8738},
	expr: &seqExpr{
	pos: position{line: 306, col: 41, offset: 8740},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 306, col: 41, offset: 8740},
	val: "?",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 306, col: 45, offset: 8744},
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
	pos: position{line: 308, col: 1, offset: 8801},
	expr: &seqExpr{
	pos: position{line: 308, col: 13, offset: 8815},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 308, col: 13, offset: 8815},
	expr: &seqExpr{
	pos: position{line: 308, col: 14, offset: 8816},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 308, col: 14, offset: 8816},
	name: "Userinfo",
},
&litMatcher{
	pos: position{line: 308, col: 23, offset: 8825},
	val: "@",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 308, col: 29, offset: 8831},
	name: "Host",
},
&zeroOrOneExpr{
	pos: position{line: 308, col: 34, offset: 8836},
	expr: &seqExpr{
	pos: position{line: 308, col: 35, offset: 8837},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 308, col: 35, offset: 8837},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 308, col: 39, offset: 8841},
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
	pos: position{line: 310, col: 1, offset: 8849},
	expr: &zeroOrMoreExpr{
	pos: position{line: 310, col: 12, offset: 8862},
	expr: &choiceExpr{
	pos: position{line: 310, col: 14, offset: 8864},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 310, col: 14, offset: 8864},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 310, col: 27, offset: 8877},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 310, col: 40, offset: 8890},
	name: "SubDelims",
},
&litMatcher{
	pos: position{line: 310, col: 52, offset: 8902},
	val: ":",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "Host",
	pos: position{line: 312, col: 1, offset: 8910},
	expr: &choiceExpr{
	pos: position{line: 312, col: 8, offset: 8919},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 312, col: 8, offset: 8919},
	name: "IPLiteral",
},
&ruleRefExpr{
	pos: position{line: 312, col: 20, offset: 8931},
	name: "RegName",
},
	},
},
},
{
	name: "Port",
	pos: position{line: 314, col: 1, offset: 8940},
	expr: &zeroOrMoreExpr{
	pos: position{line: 314, col: 8, offset: 8949},
	expr: &ruleRefExpr{
	pos: position{line: 314, col: 8, offset: 8949},
	name: "Digit",
},
},
},
{
	name: "IPLiteral",
	pos: position{line: 316, col: 1, offset: 8957},
	expr: &seqExpr{
	pos: position{line: 316, col: 13, offset: 8971},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 316, col: 13, offset: 8971},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 316, col: 17, offset: 8975},
	name: "IPv6address",
},
&litMatcher{
	pos: position{line: 316, col: 29, offset: 8987},
	val: "]",
	ignoreCase: false,
},
	},
},
},
{
	name: "IPv6address",
	pos: position{line: 318, col: 1, offset: 8992},
	expr: &actionExpr{
	pos: position{line: 318, col: 15, offset: 9008},
	run: (*parser).callonIPv6address1,
	expr: &seqExpr{
	pos: position{line: 318, col: 15, offset: 9008},
	exprs: []interface{}{
&zeroOrMoreExpr{
	pos: position{line: 318, col: 15, offset: 9008},
	expr: &ruleRefExpr{
	pos: position{line: 318, col: 16, offset: 9009},
	name: "HexDig",
},
},
&litMatcher{
	pos: position{line: 318, col: 25, offset: 9018},
	val: ":",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 318, col: 29, offset: 9022},
	expr: &choiceExpr{
	pos: position{line: 318, col: 30, offset: 9023},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 318, col: 30, offset: 9023},
	name: "HexDig",
},
&litMatcher{
	pos: position{line: 318, col: 39, offset: 9032},
	val: ":",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 318, col: 45, offset: 9038},
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
	pos: position{line: 324, col: 1, offset: 9192},
	expr: &zeroOrMoreExpr{
	pos: position{line: 324, col: 11, offset: 9204},
	expr: &choiceExpr{
	pos: position{line: 324, col: 12, offset: 9205},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 324, col: 12, offset: 9205},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 324, col: 25, offset: 9218},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 324, col: 38, offset: 9231},
	name: "SubDelims",
},
	},
},
},
},
{
	name: "PChar",
	pos: position{line: 326, col: 1, offset: 9244},
	expr: &choiceExpr{
	pos: position{line: 326, col: 9, offset: 9254},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 326, col: 9, offset: 9254},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 326, col: 22, offset: 9267},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 326, col: 35, offset: 9280},
	name: "SubDelims",
},
&charClassMatcher{
	pos: position{line: 326, col: 47, offset: 9292},
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
	pos: position{line: 328, col: 1, offset: 9298},
	expr: &zeroOrMoreExpr{
	pos: position{line: 328, col: 9, offset: 9308},
	expr: &choiceExpr{
	pos: position{line: 328, col: 10, offset: 9309},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 328, col: 10, offset: 9309},
	name: "PChar",
},
&charClassMatcher{
	pos: position{line: 328, col: 18, offset: 9317},
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
	pos: position{line: 330, col: 1, offset: 9325},
	expr: &seqExpr{
	pos: position{line: 330, col: 14, offset: 9340},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 330, col: 14, offset: 9340},
	val: "%",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 330, col: 18, offset: 9344},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 330, col: 25, offset: 9351},
	name: "HexDig",
},
	},
},
},
{
	name: "Unreserved",
	pos: position{line: 332, col: 1, offset: 9359},
	expr: &charClassMatcher{
	pos: position{line: 332, col: 14, offset: 9374},
	val: "[A-Za-z0-9._~-]",
	chars: []rune{'.','_','~','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SubDelims",
	pos: position{line: 334, col: 1, offset: 9391},
	expr: &choiceExpr{
	pos: position{line: 334, col: 13, offset: 9405},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 334, col: 13, offset: 9405},
	val: "!",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 334, col: 19, offset: 9411},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 334, col: 25, offset: 9417},
	val: "&",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 334, col: 31, offset: 9423},
	val: "'",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 334, col: 37, offset: 9429},
	val: "(",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 334, col: 43, offset: 9435},
	val: ")",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 334, col: 49, offset: 9441},
	val: "*",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 334, col: 55, offset: 9447},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 334, col: 61, offset: 9453},
	val: ",",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 334, col: 67, offset: 9459},
	val: ";",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 334, col: 73, offset: 9465},
	val: "=",
	ignoreCase: false,
},
	},
},
},
{
	name: "Http",
	pos: position{line: 336, col: 1, offset: 9470},
	expr: &actionExpr{
	pos: position{line: 336, col: 8, offset: 9479},
	run: (*parser).callonHttp1,
	expr: &labeledExpr{
	pos: position{line: 336, col: 8, offset: 9479},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 336, col: 10, offset: 9481},
	name: "HttpRaw",
},
},
},
},
{
	name: "Env",
	pos: position{line: 338, col: 1, offset: 9526},
	expr: &actionExpr{
	pos: position{line: 338, col: 7, offset: 9534},
	run: (*parser).callonEnv1,
	expr: &seqExpr{
	pos: position{line: 338, col: 7, offset: 9534},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 338, col: 7, offset: 9534},
	val: "env:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 338, col: 14, offset: 9541},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 338, col: 17, offset: 9544},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 338, col: 17, offset: 9544},
	name: "BashEnvironmentVariable",
},
&ruleRefExpr{
	pos: position{line: 338, col: 43, offset: 9570},
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
	pos: position{line: 340, col: 1, offset: 9615},
	expr: &actionExpr{
	pos: position{line: 340, col: 27, offset: 9643},
	run: (*parser).callonBashEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 340, col: 27, offset: 9643},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 340, col: 27, offset: 9643},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 340, col: 36, offset: 9652},
	expr: &charClassMatcher{
	pos: position{line: 340, col: 36, offset: 9652},
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
	pos: position{line: 344, col: 1, offset: 9708},
	expr: &actionExpr{
	pos: position{line: 344, col: 28, offset: 9737},
	run: (*parser).callonPosixEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 344, col: 28, offset: 9737},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 344, col: 28, offset: 9737},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 344, col: 32, offset: 9741},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 344, col: 34, offset: 9743},
	name: "PosixEnvironmentVariableContent",
},
},
&litMatcher{
	pos: position{line: 344, col: 66, offset: 9775},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PosixEnvironmentVariableContent",
	pos: position{line: 348, col: 1, offset: 9800},
	expr: &actionExpr{
	pos: position{line: 348, col: 35, offset: 9836},
	run: (*parser).callonPosixEnvironmentVariableContent1,
	expr: &labeledExpr{
	pos: position{line: 348, col: 35, offset: 9836},
	label: "v",
	expr: &oneOrMoreExpr{
	pos: position{line: 348, col: 37, offset: 9838},
	expr: &ruleRefExpr{
	pos: position{line: 348, col: 37, offset: 9838},
	name: "PosixEnvironmentVariableCharacter",
},
},
},
},
},
{
	name: "PosixEnvironmentVariableCharacter",
	pos: position{line: 357, col: 1, offset: 10051},
	expr: &choiceExpr{
	pos: position{line: 358, col: 7, offset: 10095},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 358, col: 7, offset: 10095},
	run: (*parser).callonPosixEnvironmentVariableCharacter2,
	expr: &litMatcher{
	pos: position{line: 358, col: 7, offset: 10095},
	val: "\\\"",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 359, col: 7, offset: 10135},
	run: (*parser).callonPosixEnvironmentVariableCharacter4,
	expr: &litMatcher{
	pos: position{line: 359, col: 7, offset: 10135},
	val: "\\\\",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 360, col: 7, offset: 10175},
	run: (*parser).callonPosixEnvironmentVariableCharacter6,
	expr: &litMatcher{
	pos: position{line: 360, col: 7, offset: 10175},
	val: "\\a",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 361, col: 7, offset: 10215},
	run: (*parser).callonPosixEnvironmentVariableCharacter8,
	expr: &litMatcher{
	pos: position{line: 361, col: 7, offset: 10215},
	val: "\\b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 362, col: 7, offset: 10255},
	run: (*parser).callonPosixEnvironmentVariableCharacter10,
	expr: &litMatcher{
	pos: position{line: 362, col: 7, offset: 10255},
	val: "\\f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 363, col: 7, offset: 10295},
	run: (*parser).callonPosixEnvironmentVariableCharacter12,
	expr: &litMatcher{
	pos: position{line: 363, col: 7, offset: 10295},
	val: "\\n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 364, col: 7, offset: 10335},
	run: (*parser).callonPosixEnvironmentVariableCharacter14,
	expr: &litMatcher{
	pos: position{line: 364, col: 7, offset: 10335},
	val: "\\r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 365, col: 7, offset: 10375},
	run: (*parser).callonPosixEnvironmentVariableCharacter16,
	expr: &litMatcher{
	pos: position{line: 365, col: 7, offset: 10375},
	val: "\\t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 366, col: 7, offset: 10415},
	run: (*parser).callonPosixEnvironmentVariableCharacter18,
	expr: &litMatcher{
	pos: position{line: 366, col: 7, offset: 10415},
	val: "\\v",
	ignoreCase: false,
},
},
&charClassMatcher{
	pos: position{line: 367, col: 7, offset: 10455},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 368, col: 7, offset: 10473},
	val: "[\\x23-\\x3c]",
	ranges: []rune{'#','<',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 369, col: 7, offset: 10491},
	val: "[\\x3e-\\x5b]",
	ranges: []rune{'>','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 370, col: 7, offset: 10509},
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
	pos: position{line: 372, col: 1, offset: 10522},
	expr: &choiceExpr{
	pos: position{line: 372, col: 14, offset: 10537},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 372, col: 14, offset: 10537},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 372, col: 24, offset: 10547},
	name: "Local",
},
&ruleRefExpr{
	pos: position{line: 372, col: 32, offset: 10555},
	name: "Http",
},
&ruleRefExpr{
	pos: position{line: 372, col: 39, offset: 10562},
	name: "Env",
},
	},
},
},
{
	name: "ImportHashed",
	pos: position{line: 374, col: 1, offset: 10567},
	expr: &actionExpr{
	pos: position{line: 374, col: 16, offset: 10584},
	run: (*parser).callonImportHashed1,
	expr: &labeledExpr{
	pos: position{line: 374, col: 16, offset: 10584},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 374, col: 18, offset: 10586},
	name: "ImportType",
},
},
},
},
{
	name: "Import",
	pos: position{line: 376, col: 1, offset: 10653},
	expr: &choiceExpr{
	pos: position{line: 376, col: 10, offset: 10664},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 376, col: 10, offset: 10664},
	run: (*parser).callonImport2,
	expr: &seqExpr{
	pos: position{line: 376, col: 10, offset: 10664},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 376, col: 10, offset: 10664},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 376, col: 12, offset: 10666},
	name: "ImportHashed",
},
},
&ruleRefExpr{
	pos: position{line: 376, col: 25, offset: 10679},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 376, col: 27, offset: 10681},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 376, col: 30, offset: 10684},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 376, col: 33, offset: 10687},
	name: "Text",
},
	},
},
},
&actionExpr{
	pos: position{line: 377, col: 10, offset: 10784},
	run: (*parser).callonImport10,
	expr: &labeledExpr{
	pos: position{line: 377, col: 10, offset: 10784},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 377, col: 12, offset: 10786},
	name: "ImportHashed",
},
},
},
	},
},
},
{
	name: "LetBinding",
	pos: position{line: 380, col: 1, offset: 10881},
	expr: &actionExpr{
	pos: position{line: 380, col: 14, offset: 10896},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 380, col: 14, offset: 10896},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 380, col: 14, offset: 10896},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 380, col: 18, offset: 10900},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 380, col: 21, offset: 10903},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 380, col: 27, offset: 10909},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 380, col: 44, offset: 10926},
	name: "_",
},
&labeledExpr{
	pos: position{line: 380, col: 46, offset: 10928},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 380, col: 48, offset: 10930},
	expr: &seqExpr{
	pos: position{line: 380, col: 49, offset: 10931},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 380, col: 49, offset: 10931},
	name: "Annotation",
},
&ruleRefExpr{
	pos: position{line: 380, col: 60, offset: 10942},
	name: "_",
},
	},
},
},
},
&litMatcher{
	pos: position{line: 381, col: 13, offset: 10958},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 381, col: 17, offset: 10962},
	name: "_",
},
&labeledExpr{
	pos: position{line: 381, col: 19, offset: 10964},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 381, col: 21, offset: 10966},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 381, col: 32, offset: 10977},
	name: "_",
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 396, col: 1, offset: 11286},
	expr: &choiceExpr{
	pos: position{line: 397, col: 7, offset: 11307},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 397, col: 7, offset: 11307},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 397, col: 7, offset: 11307},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 397, col: 7, offset: 11307},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 397, col: 14, offset: 11314},
	name: "_",
},
&litMatcher{
	pos: position{line: 397, col: 16, offset: 11316},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 397, col: 20, offset: 11320},
	name: "_",
},
&labeledExpr{
	pos: position{line: 397, col: 22, offset: 11322},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 397, col: 28, offset: 11328},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 397, col: 45, offset: 11345},
	name: "_",
},
&litMatcher{
	pos: position{line: 397, col: 47, offset: 11347},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 397, col: 51, offset: 11351},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 397, col: 54, offset: 11354},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 397, col: 56, offset: 11356},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 397, col: 67, offset: 11367},
	name: "_",
},
&litMatcher{
	pos: position{line: 397, col: 69, offset: 11369},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 397, col: 73, offset: 11373},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 397, col: 75, offset: 11375},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 397, col: 81, offset: 11381},
	name: "_",
},
&labeledExpr{
	pos: position{line: 397, col: 83, offset: 11383},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 397, col: 88, offset: 11388},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 400, col: 7, offset: 11505},
	run: (*parser).callonExpression22,
	expr: &seqExpr{
	pos: position{line: 400, col: 7, offset: 11505},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 400, col: 7, offset: 11505},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 400, col: 10, offset: 11508},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 400, col: 13, offset: 11511},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 400, col: 18, offset: 11516},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 400, col: 29, offset: 11527},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 400, col: 31, offset: 11529},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 400, col: 36, offset: 11534},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 400, col: 39, offset: 11537},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 400, col: 41, offset: 11539},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 400, col: 52, offset: 11550},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 400, col: 54, offset: 11552},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 400, col: 59, offset: 11557},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 400, col: 62, offset: 11560},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 400, col: 64, offset: 11562},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 403, col: 7, offset: 11648},
	run: (*parser).callonExpression38,
	expr: &seqExpr{
	pos: position{line: 403, col: 7, offset: 11648},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 403, col: 7, offset: 11648},
	label: "bindings",
	expr: &oneOrMoreExpr{
	pos: position{line: 403, col: 16, offset: 11657},
	expr: &ruleRefExpr{
	pos: position{line: 403, col: 16, offset: 11657},
	name: "LetBinding",
},
},
},
&ruleRefExpr{
	pos: position{line: 403, col: 28, offset: 11669},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 403, col: 31, offset: 11672},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 403, col: 34, offset: 11675},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 403, col: 36, offset: 11677},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 410, col: 7, offset: 11917},
	run: (*parser).callonExpression47,
	expr: &seqExpr{
	pos: position{line: 410, col: 7, offset: 11917},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 410, col: 7, offset: 11917},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 410, col: 14, offset: 11924},
	name: "_",
},
&litMatcher{
	pos: position{line: 410, col: 16, offset: 11926},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 410, col: 20, offset: 11930},
	name: "_",
},
&labeledExpr{
	pos: position{line: 410, col: 22, offset: 11932},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 410, col: 28, offset: 11938},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 410, col: 45, offset: 11955},
	name: "_",
},
&litMatcher{
	pos: position{line: 410, col: 47, offset: 11957},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 410, col: 51, offset: 11961},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 410, col: 54, offset: 11964},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 410, col: 56, offset: 11966},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 410, col: 67, offset: 11977},
	name: "_",
},
&litMatcher{
	pos: position{line: 410, col: 69, offset: 11979},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 410, col: 73, offset: 11983},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 410, col: 75, offset: 11985},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 410, col: 81, offset: 11991},
	name: "_",
},
&labeledExpr{
	pos: position{line: 410, col: 83, offset: 11993},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 410, col: 88, offset: 11998},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 413, col: 7, offset: 12107},
	run: (*parser).callonExpression67,
	expr: &seqExpr{
	pos: position{line: 413, col: 7, offset: 12107},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 413, col: 7, offset: 12107},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 413, col: 9, offset: 12109},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 413, col: 28, offset: 12128},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 413, col: 30, offset: 12130},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 413, col: 36, offset: 12136},
	name: "_",
},
&labeledExpr{
	pos: position{line: 413, col: 38, offset: 12138},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 413, col: 40, offset: 12140},
	name: "Expression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 414, col: 7, offset: 12200},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 416, col: 1, offset: 12221},
	expr: &actionExpr{
	pos: position{line: 416, col: 14, offset: 12236},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 416, col: 14, offset: 12236},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 416, col: 14, offset: 12236},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 416, col: 18, offset: 12240},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 416, col: 21, offset: 12243},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 416, col: 23, offset: 12245},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 418, col: 1, offset: 12275},
	expr: &choiceExpr{
	pos: position{line: 419, col: 5, offset: 12303},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 419, col: 5, offset: 12303},
	name: "EmptyList",
},
&actionExpr{
	pos: position{line: 420, col: 5, offset: 12317},
	run: (*parser).callonAnnotatedExpression3,
	expr: &seqExpr{
	pos: position{line: 420, col: 5, offset: 12317},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 420, col: 5, offset: 12317},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 420, col: 7, offset: 12319},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 420, col: 26, offset: 12338},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 420, col: 28, offset: 12340},
	expr: &seqExpr{
	pos: position{line: 420, col: 29, offset: 12341},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 420, col: 29, offset: 12341},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 420, col: 31, offset: 12343},
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
},
},
{
	name: "EmptyList",
	pos: position{line: 425, col: 1, offset: 12468},
	expr: &actionExpr{
	pos: position{line: 425, col: 13, offset: 12482},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 425, col: 13, offset: 12482},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 425, col: 13, offset: 12482},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 425, col: 17, offset: 12486},
	name: "_",
},
&litMatcher{
	pos: position{line: 425, col: 19, offset: 12488},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 425, col: 23, offset: 12492},
	name: "_",
},
&litMatcher{
	pos: position{line: 425, col: 25, offset: 12494},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 425, col: 29, offset: 12498},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 425, col: 32, offset: 12501},
	name: "List",
},
&ruleRefExpr{
	pos: position{line: 425, col: 37, offset: 12506},
	name: "_",
},
&labeledExpr{
	pos: position{line: 425, col: 39, offset: 12508},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 425, col: 41, offset: 12510},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 429, col: 1, offset: 12573},
	expr: &ruleRefExpr{
	pos: position{line: 429, col: 22, offset: 12596},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 431, col: 1, offset: 12617},
	expr: &ruleRefExpr{
	pos: position{line: 431, col: 24, offset: 12642},
	name: "OrExpression",
},
},
{
	name: "OrExpression",
	pos: position{line: 433, col: 1, offset: 12656},
	expr: &actionExpr{
	pos: position{line: 433, col: 24, offset: 12681},
	run: (*parser).callonOrExpression1,
	expr: &seqExpr{
	pos: position{line: 433, col: 24, offset: 12681},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 433, col: 24, offset: 12681},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 433, col: 30, offset: 12687},
	name: "PlusExpression",
},
},
&labeledExpr{
	pos: position{line: 433, col: 52, offset: 12709},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 433, col: 57, offset: 12714},
	expr: &seqExpr{
	pos: position{line: 433, col: 58, offset: 12715},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 433, col: 58, offset: 12715},
	name: "_",
},
&litMatcher{
	pos: position{line: 433, col: 60, offset: 12717},
	val: "||",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 433, col: 65, offset: 12722},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 433, col: 67, offset: 12724},
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
	pos: position{line: 435, col: 1, offset: 12790},
	expr: &actionExpr{
	pos: position{line: 435, col: 24, offset: 12815},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 435, col: 24, offset: 12815},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 435, col: 24, offset: 12815},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 435, col: 30, offset: 12821},
	name: "ListAppendExpression",
},
},
&labeledExpr{
	pos: position{line: 435, col: 52, offset: 12843},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 435, col: 57, offset: 12848},
	expr: &seqExpr{
	pos: position{line: 435, col: 58, offset: 12849},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 435, col: 58, offset: 12849},
	name: "_",
},
&litMatcher{
	pos: position{line: 435, col: 60, offset: 12851},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 435, col: 64, offset: 12855},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 435, col: 67, offset: 12858},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 435, col: 69, offset: 12860},
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
	pos: position{line: 437, col: 1, offset: 12934},
	expr: &actionExpr{
	pos: position{line: 437, col: 24, offset: 12959},
	run: (*parser).callonListAppendExpression1,
	expr: &seqExpr{
	pos: position{line: 437, col: 24, offset: 12959},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 437, col: 24, offset: 12959},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 437, col: 30, offset: 12965},
	name: "AndExpression",
},
},
&labeledExpr{
	pos: position{line: 437, col: 52, offset: 12987},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 437, col: 57, offset: 12992},
	expr: &seqExpr{
	pos: position{line: 437, col: 58, offset: 12993},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 437, col: 58, offset: 12993},
	name: "_",
},
&litMatcher{
	pos: position{line: 437, col: 60, offset: 12995},
	val: "#",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 437, col: 64, offset: 12999},
	name: "_",
},
&labeledExpr{
	pos: position{line: 437, col: 66, offset: 13001},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 437, col: 68, offset: 13003},
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
	pos: position{line: 439, col: 1, offset: 13076},
	expr: &actionExpr{
	pos: position{line: 439, col: 24, offset: 13101},
	run: (*parser).callonAndExpression1,
	expr: &seqExpr{
	pos: position{line: 439, col: 24, offset: 13101},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 439, col: 24, offset: 13101},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 439, col: 30, offset: 13107},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 439, col: 52, offset: 13129},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 439, col: 57, offset: 13134},
	expr: &seqExpr{
	pos: position{line: 439, col: 58, offset: 13135},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 439, col: 58, offset: 13135},
	name: "_",
},
&litMatcher{
	pos: position{line: 439, col: 60, offset: 13137},
	val: "&&",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 439, col: 65, offset: 13142},
	name: "_",
},
&labeledExpr{
	pos: position{line: 439, col: 67, offset: 13144},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 439, col: 69, offset: 13146},
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
	pos: position{line: 441, col: 1, offset: 13214},
	expr: &actionExpr{
	pos: position{line: 441, col: 24, offset: 13239},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 441, col: 24, offset: 13239},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 441, col: 24, offset: 13239},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 441, col: 30, offset: 13245},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 441, col: 52, offset: 13267},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 441, col: 57, offset: 13272},
	expr: &seqExpr{
	pos: position{line: 441, col: 58, offset: 13273},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 441, col: 58, offset: 13273},
	name: "_",
},
&litMatcher{
	pos: position{line: 441, col: 60, offset: 13275},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 441, col: 64, offset: 13279},
	name: "_",
},
&labeledExpr{
	pos: position{line: 441, col: 66, offset: 13281},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 441, col: 68, offset: 13283},
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
	pos: position{line: 444, col: 1, offset: 13360},
	expr: &actionExpr{
	pos: position{line: 444, col: 25, offset: 13386},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 444, col: 25, offset: 13386},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 444, col: 25, offset: 13386},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 444, col: 27, offset: 13388},
	name: "FirstApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 444, col: 54, offset: 13415},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 444, col: 59, offset: 13420},
	expr: &seqExpr{
	pos: position{line: 444, col: 60, offset: 13421},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 444, col: 60, offset: 13421},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 444, col: 63, offset: 13424},
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
	pos: position{line: 453, col: 1, offset: 13674},
	expr: &choiceExpr{
	pos: position{line: 454, col: 8, offset: 13712},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 454, col: 8, offset: 13712},
	run: (*parser).callonFirstApplicationExpression2,
	expr: &seqExpr{
	pos: position{line: 454, col: 8, offset: 13712},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 454, col: 8, offset: 13712},
	name: "Some",
},
&ruleRefExpr{
	pos: position{line: 454, col: 13, offset: 13717},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 454, col: 16, offset: 13720},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 454, col: 18, offset: 13722},
	name: "ImportExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 455, col: 8, offset: 13777},
	name: "ImportExpression",
},
	},
},
},
{
	name: "ImportExpression",
	pos: position{line: 457, col: 1, offset: 13795},
	expr: &choiceExpr{
	pos: position{line: 457, col: 20, offset: 13816},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 457, col: 20, offset: 13816},
	name: "Import",
},
&ruleRefExpr{
	pos: position{line: 457, col: 29, offset: 13825},
	name: "SelectorExpression",
},
	},
},
},
{
	name: "SelectorExpression",
	pos: position{line: 459, col: 1, offset: 13845},
	expr: &actionExpr{
	pos: position{line: 459, col: 22, offset: 13868},
	run: (*parser).callonSelectorExpression1,
	expr: &seqExpr{
	pos: position{line: 459, col: 22, offset: 13868},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 459, col: 22, offset: 13868},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 459, col: 24, offset: 13870},
	name: "PrimitiveExpression",
},
},
&labeledExpr{
	pos: position{line: 459, col: 44, offset: 13890},
	label: "ls",
	expr: &zeroOrMoreExpr{
	pos: position{line: 459, col: 47, offset: 13893},
	expr: &seqExpr{
	pos: position{line: 459, col: 48, offset: 13894},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 459, col: 48, offset: 13894},
	name: "_",
},
&litMatcher{
	pos: position{line: 459, col: 50, offset: 13896},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 459, col: 54, offset: 13900},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 459, col: 56, offset: 13902},
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
	pos: position{line: 469, col: 1, offset: 14135},
	expr: &choiceExpr{
	pos: position{line: 470, col: 7, offset: 14165},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 470, col: 7, offset: 14165},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 471, col: 7, offset: 14185},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 472, col: 7, offset: 14206},
	name: "IntegerLiteral",
},
&ruleRefExpr{
	pos: position{line: 473, col: 7, offset: 14227},
	name: "TextLiteral",
},
&actionExpr{
	pos: position{line: 474, col: 7, offset: 14245},
	run: (*parser).callonPrimitiveExpression6,
	expr: &seqExpr{
	pos: position{line: 474, col: 7, offset: 14245},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 474, col: 7, offset: 14245},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 474, col: 11, offset: 14249},
	name: "_",
},
&labeledExpr{
	pos: position{line: 474, col: 13, offset: 14251},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 474, col: 15, offset: 14253},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 474, col: 35, offset: 14273},
	name: "_",
},
&litMatcher{
	pos: position{line: 474, col: 37, offset: 14275},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 475, col: 7, offset: 14303},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 476, col: 7, offset: 14329},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 477, col: 7, offset: 14346},
	run: (*parser).callonPrimitiveExpression16,
	expr: &seqExpr{
	pos: position{line: 477, col: 7, offset: 14346},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 477, col: 7, offset: 14346},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 477, col: 11, offset: 14350},
	name: "_",
},
&labeledExpr{
	pos: position{line: 477, col: 14, offset: 14353},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 477, col: 16, offset: 14355},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 477, col: 27, offset: 14366},
	name: "_",
},
&litMatcher{
	pos: position{line: 477, col: 29, offset: 14368},
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
	pos: position{line: 479, col: 1, offset: 14391},
	expr: &choiceExpr{
	pos: position{line: 480, col: 7, offset: 14421},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 480, col: 7, offset: 14421},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &litMatcher{
	pos: position{line: 480, col: 7, offset: 14421},
	val: "=",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 481, col: 7, offset: 14476},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 482, col: 7, offset: 14501},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 483, col: 7, offset: 14529},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 483, col: 7, offset: 14529},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 485, col: 1, offset: 14575},
	expr: &actionExpr{
	pos: position{line: 485, col: 19, offset: 14595},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 485, col: 19, offset: 14595},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 485, col: 19, offset: 14595},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 485, col: 24, offset: 14600},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 485, col: 33, offset: 14609},
	name: "_",
},
&litMatcher{
	pos: position{line: 485, col: 35, offset: 14611},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 485, col: 39, offset: 14615},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 485, col: 42, offset: 14618},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 485, col: 47, offset: 14623},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 488, col: 1, offset: 14680},
	expr: &actionExpr{
	pos: position{line: 488, col: 18, offset: 14699},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 488, col: 18, offset: 14699},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 488, col: 18, offset: 14699},
	name: "_",
},
&litMatcher{
	pos: position{line: 488, col: 20, offset: 14701},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 488, col: 24, offset: 14705},
	name: "_",
},
&labeledExpr{
	pos: position{line: 488, col: 26, offset: 14707},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 488, col: 28, offset: 14709},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 489, col: 1, offset: 14741},
	expr: &actionExpr{
	pos: position{line: 490, col: 7, offset: 14770},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 490, col: 7, offset: 14770},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 490, col: 7, offset: 14770},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 490, col: 13, offset: 14776},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 490, col: 29, offset: 14792},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 490, col: 34, offset: 14797},
	expr: &ruleRefExpr{
	pos: position{line: 490, col: 34, offset: 14797},
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
	pos: position{line: 500, col: 1, offset: 15193},
	expr: &actionExpr{
	pos: position{line: 500, col: 22, offset: 15216},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 500, col: 22, offset: 15216},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 500, col: 22, offset: 15216},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 500, col: 27, offset: 15221},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 500, col: 36, offset: 15230},
	name: "_",
},
&litMatcher{
	pos: position{line: 500, col: 38, offset: 15232},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 500, col: 42, offset: 15236},
	name: "_",
},
&labeledExpr{
	pos: position{line: 500, col: 44, offset: 15238},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 500, col: 49, offset: 15243},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 503, col: 1, offset: 15300},
	expr: &actionExpr{
	pos: position{line: 503, col: 21, offset: 15322},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 503, col: 21, offset: 15322},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 503, col: 21, offset: 15322},
	name: "_",
},
&litMatcher{
	pos: position{line: 503, col: 23, offset: 15324},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 503, col: 27, offset: 15328},
	name: "_",
},
&labeledExpr{
	pos: position{line: 503, col: 29, offset: 15330},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 503, col: 31, offset: 15332},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 504, col: 1, offset: 15367},
	expr: &actionExpr{
	pos: position{line: 505, col: 7, offset: 15399},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 505, col: 7, offset: 15399},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 505, col: 7, offset: 15399},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 505, col: 13, offset: 15405},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 505, col: 32, offset: 15424},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 505, col: 37, offset: 15429},
	expr: &ruleRefExpr{
	pos: position{line: 505, col: 37, offset: 15429},
	name: "MoreRecordLiteral",
},
},
},
	},
},
},
},
{
	name: "MoreList",
	pos: position{line: 515, col: 1, offset: 15831},
	expr: &actionExpr{
	pos: position{line: 515, col: 12, offset: 15844},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 515, col: 12, offset: 15844},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 515, col: 12, offset: 15844},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 515, col: 16, offset: 15848},
	name: "_",
},
&labeledExpr{
	pos: position{line: 515, col: 18, offset: 15850},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 515, col: 20, offset: 15852},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 515, col: 31, offset: 15863},
	name: "_",
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 517, col: 1, offset: 15882},
	expr: &actionExpr{
	pos: position{line: 518, col: 7, offset: 15912},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 518, col: 7, offset: 15912},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 518, col: 7, offset: 15912},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 518, col: 11, offset: 15916},
	name: "_",
},
&labeledExpr{
	pos: position{line: 518, col: 13, offset: 15918},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 518, col: 19, offset: 15924},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 518, col: 30, offset: 15935},
	name: "_",
},
&labeledExpr{
	pos: position{line: 518, col: 32, offset: 15937},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 518, col: 37, offset: 15942},
	expr: &ruleRefExpr{
	pos: position{line: 518, col: 37, offset: 15942},
	name: "MoreList",
},
},
},
&litMatcher{
	pos: position{line: 518, col: 47, offset: 15952},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 528, col: 1, offset: 16228},
	expr: &notExpr{
	pos: position{line: 528, col: 7, offset: 16236},
	expr: &anyMatcher{
	line: 528, col: 8, offset: 16237,
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

func (c *current) onLabel1(label interface{}) (interface{}, error) {
 return label, nil 
}

func (p *parser) callonLabel1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLabel1(stack["label"])
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
 return nil, errors.New("Natural/build unimplemented") 
}

func (p *parser) callonReserved2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved2()
}

func (c *current) onReserved4() (interface{}, error) {
 return nil, errors.New("Natural/fold unimplemented") 
}

func (p *parser) callonReserved4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved4()
}

func (c *current) onReserved6() (interface{}, error) {
 return nil, errors.New("Natural/isZero unimplemented") 
}

func (p *parser) callonReserved6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved6()
}

func (c *current) onReserved8() (interface{}, error) {
 return nil, errors.New("Natural/even unimplemented") 
}

func (p *parser) callonReserved8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved8()
}

func (c *current) onReserved10() (interface{}, error) {
 return nil, errors.New("Natural/odd unimplemented") 
}

func (p *parser) callonReserved10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved10()
}

func (c *current) onReserved12() (interface{}, error) {
 return nil, errors.New("Natural/toInteger unimplemented") 
}

func (p *parser) callonReserved12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved12()
}

func (c *current) onReserved14() (interface{}, error) {
 return nil, errors.New("Natural/show unimplemented") 
}

func (p *parser) callonReserved14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved14()
}

func (c *current) onReserved16() (interface{}, error) {
 return nil, errors.New("Integer/toDouble unimplemented") 
}

func (p *parser) callonReserved16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved16()
}

func (c *current) onReserved18() (interface{}, error) {
 return nil, errors.New("Integer/show unimplemented") 
}

func (p *parser) callonReserved18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved18()
}

func (c *current) onReserved20() (interface{}, error) {
 return nil, errors.New("Double/show unimplemented") 
}

func (p *parser) callonReserved20() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved20()
}

func (c *current) onReserved22() (interface{}, error) {
 return nil, errors.New("List/build unimplemented") 
}

func (p *parser) callonReserved22() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved22()
}

func (c *current) onReserved24() (interface{}, error) {
 return nil, errors.New("List/fold unimplemented") 
}

func (p *parser) callonReserved24() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved24()
}

func (c *current) onReserved26() (interface{}, error) {
 return nil, errors.New("List/length unimplemented") 
}

func (p *parser) callonReserved26() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved26()
}

func (c *current) onReserved28() (interface{}, error) {
 return nil, errors.New("List/head unimplemented") 
}

func (p *parser) callonReserved28() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved28()
}

func (c *current) onReserved30() (interface{}, error) {
 return nil, errors.New("List/last unimplemented") 
}

func (p *parser) callonReserved30() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved30()
}

func (c *current) onReserved32() (interface{}, error) {
 return nil, errors.New("List/indexed unimplemented") 
}

func (p *parser) callonReserved32() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved32()
}

func (c *current) onReserved34() (interface{}, error) {
 return nil, errors.New("List/reverse unimplemented") 
}

func (p *parser) callonReserved34() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved34()
}

func (c *current) onReserved36() (interface{}, error) {
 return nil, errors.New("Optional/build unimplemented") 
}

func (p *parser) callonReserved36() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved36()
}

func (c *current) onReserved38() (interface{}, error) {
 return nil, errors.New("Optional/fold unimplemented") 
}

func (p *parser) callonReserved38() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved38()
}

func (c *current) onReserved40() (interface{}, error) {
 return nil, errors.New("Text/show unimplemented") 
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
 return &Pi{"_",o.(Expr),e.(Expr)}, nil 
}

func (p *parser) callonExpression67() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression67(stack["o"], stack["e"])
}

func (c *current) onAnnotation1(a interface{}) (interface{}, error) {
 return a, nil 
}

func (p *parser) callonAnnotation1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAnnotation1(stack["a"])
}

func (c *current) onAnnotatedExpression3(e, a interface{}) (interface{}, error) {
        if a == nil { return e, nil }
        return Annot{e.(Expr), a.([]interface{})[1].(Expr)}, nil
    
}

func (p *parser) callonAnnotatedExpression3() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAnnotatedExpression3(stack["e"], stack["a"])
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

func (c *current) onTimesExpression1(first, rest interface{}) (interface{}, error) {
return ParseOperator(TimesOp, first, rest), nil
}

func (p *parser) callonTimesExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTimesExpression1(stack["first"], stack["rest"])
}

func (c *current) onApplicationExpression1(f, rest interface{}) (interface{}, error) {
          e := f.(Expr)
          if rest == nil { return e, nil }
          for _, arg := range rest.([]interface{}) {
              e = &App{Fn:e, Arg: arg.([]interface{})[1].(Expr)}
          }
          return e,nil
      
}

func (p *parser) callonApplicationExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onApplicationExpression1(stack["f"], stack["rest"])
}

func (c *current) onFirstApplicationExpression2(e interface{}) (interface{}, error) {
 return Some{e.(Expr)}, nil 
}

func (p *parser) callonFirstApplicationExpression2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFirstApplicationExpression2(stack["e"])
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

func (c *current) onPrimitiveExpression16(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonPrimitiveExpression16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimitiveExpression16(stack["e"])
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
              content[field.([]interface{})[0].(string)] = field.([]interface{})[1].(Expr)
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
              content[field.([]interface{})[0].(string)] = field.([]interface{})[1].(Expr)
          }
          return RecordLit(content), nil
      
}

func (p *parser) callonNonEmptyRecordLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNonEmptyRecordLiteral1(stack["first"], stack["rest"])
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

