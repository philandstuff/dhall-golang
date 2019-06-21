
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


var g = &grammar {
	rules: []*rule{
{
	name: "DhallFile",
	pos: position{line: 38, col: 1, offset: 652},
	expr: &actionExpr{
	pos: position{line: 38, col: 13, offset: 666},
	run: (*parser).callonDhallFile1,
	expr: &seqExpr{
	pos: position{line: 38, col: 13, offset: 666},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 38, col: 13, offset: 666},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 38, col: 15, offset: 668},
	name: "CompleteExpression",
},
},
&ruleRefExpr{
	pos: position{line: 38, col: 34, offset: 687},
	name: "EOF",
},
	},
},
},
},
{
	name: "CompleteExpression",
	pos: position{line: 40, col: 1, offset: 710},
	expr: &actionExpr{
	pos: position{line: 40, col: 22, offset: 733},
	run: (*parser).callonCompleteExpression1,
	expr: &seqExpr{
	pos: position{line: 40, col: 22, offset: 733},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 40, col: 22, offset: 733},
	name: "_",
},
&labeledExpr{
	pos: position{line: 40, col: 24, offset: 735},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 40, col: 26, offset: 737},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 40, col: 37, offset: 748},
	name: "_",
},
	},
},
},
},
{
	name: "EOL",
	pos: position{line: 42, col: 1, offset: 769},
	expr: &choiceExpr{
	pos: position{line: 42, col: 7, offset: 777},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 42, col: 7, offset: 777},
	val: "\n",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 42, col: 14, offset: 784},
	run: (*parser).callonEOL3,
	expr: &litMatcher{
	pos: position{line: 42, col: 14, offset: 784},
	val: "\r\n",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "BlockComment",
	pos: position{line: 44, col: 1, offset: 821},
	expr: &seqExpr{
	pos: position{line: 44, col: 16, offset: 838},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 44, col: 16, offset: 838},
	val: "{-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 44, col: 21, offset: 843},
	name: "BlockCommentContinue",
},
	},
},
},
{
	name: "BlockCommentChunk",
	pos: position{line: 46, col: 1, offset: 865},
	expr: &choiceExpr{
	pos: position{line: 47, col: 5, offset: 891},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 47, col: 5, offset: 891},
	name: "BlockComment",
},
&charClassMatcher{
	pos: position{line: 48, col: 5, offset: 908},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 49, col: 5, offset: 934},
	name: "EOL",
},
	},
},
},
{
	name: "BlockCommentContinue",
	pos: position{line: 51, col: 1, offset: 939},
	expr: &choiceExpr{
	pos: position{line: 51, col: 24, offset: 964},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 51, col: 24, offset: 964},
	val: "-}",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 51, col: 31, offset: 971},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 51, col: 31, offset: 971},
	name: "BlockCommentChunk",
},
&ruleRefExpr{
	pos: position{line: 51, col: 49, offset: 989},
	name: "BlockCommentContinue",
},
	},
},
	},
},
},
{
	name: "NotEOL",
	pos: position{line: 53, col: 1, offset: 1011},
	expr: &charClassMatcher{
	pos: position{line: 53, col: 10, offset: 1022},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "LineComment",
	pos: position{line: 55, col: 1, offset: 1045},
	expr: &actionExpr{
	pos: position{line: 55, col: 15, offset: 1061},
	run: (*parser).callonLineComment1,
	expr: &seqExpr{
	pos: position{line: 55, col: 15, offset: 1061},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 55, col: 15, offset: 1061},
	val: "--",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 55, col: 20, offset: 1066},
	label: "content",
	expr: &actionExpr{
	pos: position{line: 55, col: 29, offset: 1075},
	run: (*parser).callonLineComment5,
	expr: &zeroOrMoreExpr{
	pos: position{line: 55, col: 29, offset: 1075},
	expr: &ruleRefExpr{
	pos: position{line: 55, col: 29, offset: 1075},
	name: "NotEOL",
},
},
},
},
&ruleRefExpr{
	pos: position{line: 55, col: 68, offset: 1114},
	name: "EOL",
},
	},
},
},
},
{
	name: "WhitespaceChunk",
	pos: position{line: 57, col: 1, offset: 1143},
	expr: &choiceExpr{
	pos: position{line: 57, col: 19, offset: 1163},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 57, col: 19, offset: 1163},
	val: " ",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 57, col: 25, offset: 1169},
	val: "\t",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 57, col: 32, offset: 1176},
	name: "EOL",
},
&ruleRefExpr{
	pos: position{line: 57, col: 38, offset: 1182},
	name: "LineComment",
},
&ruleRefExpr{
	pos: position{line: 57, col: 52, offset: 1196},
	name: "BlockComment",
},
	},
},
},
{
	name: "_",
	pos: position{line: 59, col: 1, offset: 1210},
	expr: &zeroOrMoreExpr{
	pos: position{line: 59, col: 5, offset: 1216},
	expr: &ruleRefExpr{
	pos: position{line: 59, col: 5, offset: 1216},
	name: "WhitespaceChunk",
},
},
},
{
	name: "_1",
	pos: position{line: 61, col: 1, offset: 1234},
	expr: &oneOrMoreExpr{
	pos: position{line: 61, col: 6, offset: 1241},
	expr: &ruleRefExpr{
	pos: position{line: 61, col: 6, offset: 1241},
	name: "WhitespaceChunk",
},
},
},
{
	name: "Digit",
	pos: position{line: 63, col: 1, offset: 1259},
	expr: &charClassMatcher{
	pos: position{line: 63, col: 9, offset: 1269},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "HexDig",
	pos: position{line: 65, col: 1, offset: 1276},
	expr: &choiceExpr{
	pos: position{line: 65, col: 10, offset: 1287},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 65, col: 10, offset: 1287},
	name: "Digit",
},
&charClassMatcher{
	pos: position{line: 65, col: 18, offset: 1295},
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
	pos: position{line: 67, col: 1, offset: 1303},
	expr: &charClassMatcher{
	pos: position{line: 67, col: 24, offset: 1328},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SimpleLabelNextChar",
	pos: position{line: 68, col: 1, offset: 1338},
	expr: &charClassMatcher{
	pos: position{line: 68, col: 23, offset: 1362},
	val: "[A-Za-z0-9_/-]",
	chars: []rune{'_','/','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SimpleLabel",
	pos: position{line: 69, col: 1, offset: 1377},
	expr: &choiceExpr{
	pos: position{line: 69, col: 15, offset: 1393},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 69, col: 15, offset: 1393},
	run: (*parser).callonSimpleLabel2,
	expr: &seqExpr{
	pos: position{line: 69, col: 15, offset: 1393},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 69, col: 15, offset: 1393},
	name: "Keyword",
},
&oneOrMoreExpr{
	pos: position{line: 69, col: 23, offset: 1401},
	expr: &ruleRefExpr{
	pos: position{line: 69, col: 23, offset: 1401},
	name: "SimpleLabelNextChar",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 70, col: 13, offset: 1465},
	run: (*parser).callonSimpleLabel7,
	expr: &seqExpr{
	pos: position{line: 70, col: 13, offset: 1465},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 70, col: 13, offset: 1465},
	expr: &ruleRefExpr{
	pos: position{line: 70, col: 14, offset: 1466},
	name: "Keyword",
},
},
&ruleRefExpr{
	pos: position{line: 70, col: 22, offset: 1474},
	name: "SimpleLabelFirstChar",
},
&zeroOrMoreExpr{
	pos: position{line: 70, col: 43, offset: 1495},
	expr: &ruleRefExpr{
	pos: position{line: 70, col: 43, offset: 1495},
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
	pos: position{line: 75, col: 1, offset: 1580},
	expr: &charClassMatcher{
	pos: position{line: 75, col: 19, offset: 1600},
	val: "[\\x20-\\x5f\\x61-\\x7e]",
	ranges: []rune{' ','_','a','~',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "QuotedLabel",
	pos: position{line: 76, col: 1, offset: 1621},
	expr: &actionExpr{
	pos: position{line: 76, col: 15, offset: 1637},
	run: (*parser).callonQuotedLabel1,
	expr: &oneOrMoreExpr{
	pos: position{line: 76, col: 15, offset: 1637},
	expr: &ruleRefExpr{
	pos: position{line: 76, col: 15, offset: 1637},
	name: "QuotedLabelChar",
},
},
},
},
{
	name: "Label",
	pos: position{line: 78, col: 1, offset: 1686},
	expr: &choiceExpr{
	pos: position{line: 78, col: 9, offset: 1696},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 78, col: 9, offset: 1696},
	run: (*parser).callonLabel2,
	expr: &seqExpr{
	pos: position{line: 78, col: 9, offset: 1696},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 78, col: 9, offset: 1696},
	val: "`",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 78, col: 13, offset: 1700},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 78, col: 19, offset: 1706},
	name: "QuotedLabel",
},
},
&litMatcher{
	pos: position{line: 78, col: 31, offset: 1718},
	val: "`",
	ignoreCase: false,
},
	},
},
},
&actionExpr{
	pos: position{line: 79, col: 9, offset: 1752},
	run: (*parser).callonLabel8,
	expr: &labeledExpr{
	pos: position{line: 79, col: 9, offset: 1752},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 79, col: 15, offset: 1758},
	name: "SimpleLabel",
},
},
},
	},
},
},
{
	name: "NonreservedLabel",
	pos: position{line: 81, col: 1, offset: 1793},
	expr: &choiceExpr{
	pos: position{line: 81, col: 20, offset: 1814},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 81, col: 20, offset: 1814},
	run: (*parser).callonNonreservedLabel2,
	expr: &seqExpr{
	pos: position{line: 81, col: 20, offset: 1814},
	exprs: []interface{}{
&andExpr{
	pos: position{line: 81, col: 20, offset: 1814},
	expr: &seqExpr{
	pos: position{line: 81, col: 22, offset: 1816},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 81, col: 22, offset: 1816},
	name: "Reserved",
},
&ruleRefExpr{
	pos: position{line: 81, col: 31, offset: 1825},
	name: "SimpleLabelNextChar",
},
	},
},
},
&labeledExpr{
	pos: position{line: 81, col: 52, offset: 1846},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 81, col: 58, offset: 1852},
	name: "Label",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 82, col: 19, offset: 1898},
	run: (*parser).callonNonreservedLabel10,
	expr: &seqExpr{
	pos: position{line: 82, col: 19, offset: 1898},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 82, col: 19, offset: 1898},
	expr: &ruleRefExpr{
	pos: position{line: 82, col: 20, offset: 1899},
	name: "Reserved",
},
},
&labeledExpr{
	pos: position{line: 82, col: 29, offset: 1908},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 82, col: 35, offset: 1914},
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
	pos: position{line: 84, col: 1, offset: 1943},
	expr: &ruleRefExpr{
	pos: position{line: 84, col: 12, offset: 1956},
	name: "Label",
},
},
{
	name: "DoubleQuoteChunk",
	pos: position{line: 87, col: 1, offset: 1964},
	expr: &choiceExpr{
	pos: position{line: 88, col: 6, offset: 1990},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 88, col: 6, offset: 1990},
	name: "Interpolation",
},
&actionExpr{
	pos: position{line: 89, col: 6, offset: 2009},
	run: (*parser).callonDoubleQuoteChunk3,
	expr: &seqExpr{
	pos: position{line: 89, col: 6, offset: 2009},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 89, col: 6, offset: 2009},
	val: "\\",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 89, col: 11, offset: 2014},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 89, col: 13, offset: 2016},
	name: "DoubleQuoteEscaped",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 90, col: 6, offset: 2058},
	name: "DoubleQuoteChar",
},
	},
},
},
{
	name: "DoubleQuoteEscaped",
	pos: position{line: 92, col: 1, offset: 2075},
	expr: &choiceExpr{
	pos: position{line: 93, col: 8, offset: 2105},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 93, col: 8, offset: 2105},
	val: "\"",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 94, col: 8, offset: 2116},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 95, col: 8, offset: 2127},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 96, col: 8, offset: 2139},
	val: "/",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 97, col: 8, offset: 2150},
	run: (*parser).callonDoubleQuoteEscaped6,
	expr: &litMatcher{
	pos: position{line: 97, col: 8, offset: 2150},
	val: "b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 98, col: 8, offset: 2190},
	run: (*parser).callonDoubleQuoteEscaped8,
	expr: &litMatcher{
	pos: position{line: 98, col: 8, offset: 2190},
	val: "f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 99, col: 8, offset: 2230},
	run: (*parser).callonDoubleQuoteEscaped10,
	expr: &litMatcher{
	pos: position{line: 99, col: 8, offset: 2230},
	val: "n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 100, col: 8, offset: 2270},
	run: (*parser).callonDoubleQuoteEscaped12,
	expr: &litMatcher{
	pos: position{line: 100, col: 8, offset: 2270},
	val: "r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 101, col: 8, offset: 2310},
	run: (*parser).callonDoubleQuoteEscaped14,
	expr: &litMatcher{
	pos: position{line: 101, col: 8, offset: 2310},
	val: "t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 102, col: 8, offset: 2350},
	run: (*parser).callonDoubleQuoteEscaped16,
	expr: &seqExpr{
	pos: position{line: 102, col: 8, offset: 2350},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 102, col: 8, offset: 2350},
	val: "u",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 102, col: 12, offset: 2354},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 102, col: 14, offset: 2356},
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
	pos: position{line: 104, col: 1, offset: 2389},
	expr: &choiceExpr{
	pos: position{line: 105, col: 9, offset: 2415},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 105, col: 9, offset: 2415},
	run: (*parser).callonUnicodeEscape2,
	expr: &seqExpr{
	pos: position{line: 105, col: 9, offset: 2415},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 105, col: 9, offset: 2415},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 105, col: 16, offset: 2422},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 105, col: 23, offset: 2429},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 105, col: 30, offset: 2436},
	name: "HexDig",
},
	},
},
},
&actionExpr{
	pos: position{line: 109, col: 9, offset: 2582},
	run: (*parser).callonUnicodeEscape8,
	expr: &seqExpr{
	pos: position{line: 109, col: 9, offset: 2582},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 109, col: 9, offset: 2582},
	val: "{",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 109, col: 13, offset: 2586},
	expr: &ruleRefExpr{
	pos: position{line: 109, col: 13, offset: 2586},
	name: "HexDig",
},
},
&litMatcher{
	pos: position{line: 109, col: 21, offset: 2594},
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
	pos: position{line: 114, col: 1, offset: 2747},
	expr: &choiceExpr{
	pos: position{line: 115, col: 6, offset: 2772},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 115, col: 6, offset: 2772},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 116, col: 6, offset: 2789},
	val: "[\\x23-\\x5b]",
	ranges: []rune{'#','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 117, col: 6, offset: 2806},
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
	pos: position{line: 119, col: 1, offset: 2825},
	expr: &actionExpr{
	pos: position{line: 119, col: 22, offset: 2848},
	run: (*parser).callonDoubleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 119, col: 22, offset: 2848},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 119, col: 22, offset: 2848},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 119, col: 26, offset: 2852},
	label: "chunks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 119, col: 33, offset: 2859},
	expr: &ruleRefExpr{
	pos: position{line: 119, col: 33, offset: 2859},
	name: "DoubleQuoteChunk",
},
},
},
&litMatcher{
	pos: position{line: 119, col: 51, offset: 2877},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "SingleQuoteContinue",
	pos: position{line: 136, col: 1, offset: 3345},
	expr: &choiceExpr{
	pos: position{line: 137, col: 7, offset: 3375},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 137, col: 7, offset: 3375},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 137, col: 7, offset: 3375},
	name: "Interpolation",
},
&ruleRefExpr{
	pos: position{line: 137, col: 21, offset: 3389},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 138, col: 7, offset: 3415},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 138, col: 7, offset: 3415},
	name: "EscapedQuotePair",
},
&ruleRefExpr{
	pos: position{line: 138, col: 24, offset: 3432},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 139, col: 7, offset: 3458},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 139, col: 7, offset: 3458},
	name: "EscapedInterpolation",
},
&ruleRefExpr{
	pos: position{line: 139, col: 28, offset: 3479},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 140, col: 7, offset: 3505},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 140, col: 7, offset: 3505},
	name: "SingleQuoteChar",
},
&ruleRefExpr{
	pos: position{line: 140, col: 23, offset: 3521},
	name: "SingleQuoteContinue",
},
	},
},
&litMatcher{
	pos: position{line: 141, col: 7, offset: 3547},
	val: "''",
	ignoreCase: false,
},
	},
},
},
{
	name: "EscapedQuotePair",
	pos: position{line: 143, col: 1, offset: 3553},
	expr: &actionExpr{
	pos: position{line: 143, col: 20, offset: 3574},
	run: (*parser).callonEscapedQuotePair1,
	expr: &litMatcher{
	pos: position{line: 143, col: 20, offset: 3574},
	val: "'''",
	ignoreCase: false,
},
},
},
{
	name: "EscapedInterpolation",
	pos: position{line: 147, col: 1, offset: 3709},
	expr: &actionExpr{
	pos: position{line: 147, col: 24, offset: 3734},
	run: (*parser).callonEscapedInterpolation1,
	expr: &litMatcher{
	pos: position{line: 147, col: 24, offset: 3734},
	val: "''${",
	ignoreCase: false,
},
},
},
{
	name: "SingleQuoteChar",
	pos: position{line: 149, col: 1, offset: 3776},
	expr: &choiceExpr{
	pos: position{line: 150, col: 6, offset: 3801},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 150, col: 6, offset: 3801},
	val: "[\\x20-\\U0010ffff]",
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 151, col: 6, offset: 3824},
	val: "\t",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 152, col: 6, offset: 3834},
	name: "EOL",
},
	},
},
},
{
	name: "SingleQuoteLiteral",
	pos: position{line: 154, col: 1, offset: 3839},
	expr: &actionExpr{
	pos: position{line: 154, col: 22, offset: 3862},
	run: (*parser).callonSingleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 154, col: 22, offset: 3862},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 154, col: 22, offset: 3862},
	val: "''",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 154, col: 27, offset: 3867},
	name: "EOL",
},
&labeledExpr{
	pos: position{line: 154, col: 31, offset: 3871},
	label: "content",
	expr: &ruleRefExpr{
	pos: position{line: 154, col: 39, offset: 3879},
	name: "SingleQuoteContinue",
},
},
	},
},
},
},
{
	name: "Interpolation",
	pos: position{line: 172, col: 1, offset: 4429},
	expr: &actionExpr{
	pos: position{line: 172, col: 17, offset: 4447},
	run: (*parser).callonInterpolation1,
	expr: &seqExpr{
	pos: position{line: 172, col: 17, offset: 4447},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 172, col: 17, offset: 4447},
	val: "${",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 172, col: 22, offset: 4452},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 172, col: 24, offset: 4454},
	name: "CompleteExpression",
},
},
&litMatcher{
	pos: position{line: 172, col: 43, offset: 4473},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TextLiteral",
	pos: position{line: 174, col: 1, offset: 4496},
	expr: &choiceExpr{
	pos: position{line: 174, col: 15, offset: 4512},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 174, col: 15, offset: 4512},
	name: "DoubleQuoteLiteral",
},
&ruleRefExpr{
	pos: position{line: 174, col: 36, offset: 4533},
	name: "SingleQuoteLiteral",
},
	},
},
},
{
	name: "Reserved",
	pos: position{line: 177, col: 1, offset: 4638},
	expr: &choiceExpr{
	pos: position{line: 178, col: 5, offset: 4655},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 178, col: 5, offset: 4655},
	run: (*parser).callonReserved2,
	expr: &litMatcher{
	pos: position{line: 178, col: 5, offset: 4655},
	val: "Natural/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 179, col: 5, offset: 4704},
	run: (*parser).callonReserved4,
	expr: &litMatcher{
	pos: position{line: 179, col: 5, offset: 4704},
	val: "Natural/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 180, col: 5, offset: 4751},
	run: (*parser).callonReserved6,
	expr: &litMatcher{
	pos: position{line: 180, col: 5, offset: 4751},
	val: "Natural/isZero",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 181, col: 5, offset: 4802},
	run: (*parser).callonReserved8,
	expr: &litMatcher{
	pos: position{line: 181, col: 5, offset: 4802},
	val: "Natural/even",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 182, col: 5, offset: 4849},
	run: (*parser).callonReserved10,
	expr: &litMatcher{
	pos: position{line: 182, col: 5, offset: 4849},
	val: "Natural/odd",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 183, col: 5, offset: 4894},
	run: (*parser).callonReserved12,
	expr: &litMatcher{
	pos: position{line: 183, col: 5, offset: 4894},
	val: "Natural/toInteger",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 184, col: 5, offset: 4951},
	run: (*parser).callonReserved14,
	expr: &litMatcher{
	pos: position{line: 184, col: 5, offset: 4951},
	val: "Natural/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 185, col: 5, offset: 4998},
	run: (*parser).callonReserved16,
	expr: &litMatcher{
	pos: position{line: 185, col: 5, offset: 4998},
	val: "Integer/toDouble",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 186, col: 5, offset: 5053},
	run: (*parser).callonReserved18,
	expr: &litMatcher{
	pos: position{line: 186, col: 5, offset: 5053},
	val: "Integer/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 187, col: 5, offset: 5100},
	run: (*parser).callonReserved20,
	expr: &litMatcher{
	pos: position{line: 187, col: 5, offset: 5100},
	val: "Double/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 188, col: 5, offset: 5145},
	run: (*parser).callonReserved22,
	expr: &litMatcher{
	pos: position{line: 188, col: 5, offset: 5145},
	val: "List/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 189, col: 5, offset: 5188},
	run: (*parser).callonReserved24,
	expr: &litMatcher{
	pos: position{line: 189, col: 5, offset: 5188},
	val: "List/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 190, col: 5, offset: 5229},
	run: (*parser).callonReserved26,
	expr: &litMatcher{
	pos: position{line: 190, col: 5, offset: 5229},
	val: "List/length",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 191, col: 5, offset: 5274},
	run: (*parser).callonReserved28,
	expr: &litMatcher{
	pos: position{line: 191, col: 5, offset: 5274},
	val: "List/head",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 192, col: 5, offset: 5315},
	run: (*parser).callonReserved30,
	expr: &litMatcher{
	pos: position{line: 192, col: 5, offset: 5315},
	val: "List/last",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 193, col: 5, offset: 5356},
	run: (*parser).callonReserved32,
	expr: &litMatcher{
	pos: position{line: 193, col: 5, offset: 5356},
	val: "List/indexed",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 194, col: 5, offset: 5403},
	run: (*parser).callonReserved34,
	expr: &litMatcher{
	pos: position{line: 194, col: 5, offset: 5403},
	val: "List/reverse",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 195, col: 5, offset: 5450},
	run: (*parser).callonReserved36,
	expr: &litMatcher{
	pos: position{line: 195, col: 5, offset: 5450},
	val: "Optional/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 196, col: 5, offset: 5501},
	run: (*parser).callonReserved38,
	expr: &litMatcher{
	pos: position{line: 196, col: 5, offset: 5501},
	val: "Optional/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 197, col: 5, offset: 5550},
	run: (*parser).callonReserved40,
	expr: &litMatcher{
	pos: position{line: 197, col: 5, offset: 5550},
	val: "Text/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 198, col: 5, offset: 5591},
	run: (*parser).callonReserved42,
	expr: &litMatcher{
	pos: position{line: 198, col: 5, offset: 5591},
	val: "Bool",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 199, col: 5, offset: 5623},
	run: (*parser).callonReserved44,
	expr: &litMatcher{
	pos: position{line: 199, col: 5, offset: 5623},
	val: "True",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 200, col: 5, offset: 5655},
	run: (*parser).callonReserved46,
	expr: &litMatcher{
	pos: position{line: 200, col: 5, offset: 5655},
	val: "False",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 201, col: 5, offset: 5689},
	run: (*parser).callonReserved48,
	expr: &litMatcher{
	pos: position{line: 201, col: 5, offset: 5689},
	val: "Optional",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 202, col: 5, offset: 5729},
	run: (*parser).callonReserved50,
	expr: &litMatcher{
	pos: position{line: 202, col: 5, offset: 5729},
	val: "Natural",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 203, col: 5, offset: 5767},
	run: (*parser).callonReserved52,
	expr: &litMatcher{
	pos: position{line: 203, col: 5, offset: 5767},
	val: "Integer",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 204, col: 5, offset: 5805},
	run: (*parser).callonReserved54,
	expr: &litMatcher{
	pos: position{line: 204, col: 5, offset: 5805},
	val: "Double",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 205, col: 5, offset: 5841},
	run: (*parser).callonReserved56,
	expr: &litMatcher{
	pos: position{line: 205, col: 5, offset: 5841},
	val: "Text",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 206, col: 5, offset: 5873},
	run: (*parser).callonReserved58,
	expr: &litMatcher{
	pos: position{line: 206, col: 5, offset: 5873},
	val: "List",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 207, col: 5, offset: 5905},
	run: (*parser).callonReserved60,
	expr: &litMatcher{
	pos: position{line: 207, col: 5, offset: 5905},
	val: "None",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 208, col: 5, offset: 5937},
	run: (*parser).callonReserved62,
	expr: &litMatcher{
	pos: position{line: 208, col: 5, offset: 5937},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 209, col: 5, offset: 5969},
	run: (*parser).callonReserved64,
	expr: &litMatcher{
	pos: position{line: 209, col: 5, offset: 5969},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 210, col: 5, offset: 6001},
	run: (*parser).callonReserved66,
	expr: &litMatcher{
	pos: position{line: 210, col: 5, offset: 6001},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "If",
	pos: position{line: 212, col: 1, offset: 6030},
	expr: &litMatcher{
	pos: position{line: 212, col: 6, offset: 6037},
	val: "if",
	ignoreCase: false,
},
},
{
	name: "Then",
	pos: position{line: 213, col: 1, offset: 6042},
	expr: &litMatcher{
	pos: position{line: 213, col: 8, offset: 6051},
	val: "then",
	ignoreCase: false,
},
},
{
	name: "Else",
	pos: position{line: 214, col: 1, offset: 6058},
	expr: &litMatcher{
	pos: position{line: 214, col: 8, offset: 6067},
	val: "else",
	ignoreCase: false,
},
},
{
	name: "Let",
	pos: position{line: 215, col: 1, offset: 6074},
	expr: &litMatcher{
	pos: position{line: 215, col: 7, offset: 6082},
	val: "let",
	ignoreCase: false,
},
},
{
	name: "In",
	pos: position{line: 216, col: 1, offset: 6088},
	expr: &litMatcher{
	pos: position{line: 216, col: 6, offset: 6095},
	val: "in",
	ignoreCase: false,
},
},
{
	name: "As",
	pos: position{line: 217, col: 1, offset: 6100},
	expr: &litMatcher{
	pos: position{line: 217, col: 6, offset: 6107},
	val: "as",
	ignoreCase: false,
},
},
{
	name: "Using",
	pos: position{line: 218, col: 1, offset: 6112},
	expr: &litMatcher{
	pos: position{line: 218, col: 9, offset: 6122},
	val: "using",
	ignoreCase: false,
},
},
{
	name: "Merge",
	pos: position{line: 219, col: 1, offset: 6130},
	expr: &litMatcher{
	pos: position{line: 219, col: 9, offset: 6140},
	val: "merge",
	ignoreCase: false,
},
},
{
	name: "Missing",
	pos: position{line: 220, col: 1, offset: 6148},
	expr: &actionExpr{
	pos: position{line: 220, col: 11, offset: 6160},
	run: (*parser).callonMissing1,
	expr: &litMatcher{
	pos: position{line: 220, col: 11, offset: 6160},
	val: "missing",
	ignoreCase: false,
},
},
},
{
	name: "True",
	pos: position{line: 221, col: 1, offset: 6196},
	expr: &litMatcher{
	pos: position{line: 221, col: 8, offset: 6205},
	val: "True",
	ignoreCase: false,
},
},
{
	name: "False",
	pos: position{line: 222, col: 1, offset: 6212},
	expr: &litMatcher{
	pos: position{line: 222, col: 9, offset: 6222},
	val: "False",
	ignoreCase: false,
},
},
{
	name: "Infinity",
	pos: position{line: 223, col: 1, offset: 6230},
	expr: &litMatcher{
	pos: position{line: 223, col: 12, offset: 6243},
	val: "Infinity",
	ignoreCase: false,
},
},
{
	name: "NaN",
	pos: position{line: 224, col: 1, offset: 6254},
	expr: &litMatcher{
	pos: position{line: 224, col: 7, offset: 6262},
	val: "NaN",
	ignoreCase: false,
},
},
{
	name: "Some",
	pos: position{line: 225, col: 1, offset: 6268},
	expr: &litMatcher{
	pos: position{line: 225, col: 8, offset: 6277},
	val: "Some",
	ignoreCase: false,
},
},
{
	name: "Keyword",
	pos: position{line: 227, col: 1, offset: 6285},
	expr: &choiceExpr{
	pos: position{line: 228, col: 5, offset: 6301},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 228, col: 5, offset: 6301},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 228, col: 10, offset: 6306},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 228, col: 17, offset: 6313},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 229, col: 5, offset: 6322},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 229, col: 11, offset: 6328},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 230, col: 5, offset: 6335},
	name: "Using",
},
&ruleRefExpr{
	pos: position{line: 230, col: 13, offset: 6343},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 230, col: 23, offset: 6353},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 231, col: 5, offset: 6360},
	name: "True",
},
&ruleRefExpr{
	pos: position{line: 231, col: 12, offset: 6367},
	name: "False",
},
&ruleRefExpr{
	pos: position{line: 232, col: 5, offset: 6377},
	name: "Infinity",
},
&ruleRefExpr{
	pos: position{line: 232, col: 16, offset: 6388},
	name: "NaN",
},
&ruleRefExpr{
	pos: position{line: 233, col: 5, offset: 6396},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 233, col: 13, offset: 6404},
	name: "Some",
},
	},
},
},
{
	name: "Optional",
	pos: position{line: 235, col: 1, offset: 6410},
	expr: &litMatcher{
	pos: position{line: 235, col: 12, offset: 6423},
	val: "Optional",
	ignoreCase: false,
},
},
{
	name: "Text",
	pos: position{line: 236, col: 1, offset: 6434},
	expr: &litMatcher{
	pos: position{line: 236, col: 8, offset: 6443},
	val: "Text",
	ignoreCase: false,
},
},
{
	name: "List",
	pos: position{line: 237, col: 1, offset: 6450},
	expr: &litMatcher{
	pos: position{line: 237, col: 8, offset: 6459},
	val: "List",
	ignoreCase: false,
},
},
{
	name: "Combine",
	pos: position{line: 239, col: 1, offset: 6467},
	expr: &choiceExpr{
	pos: position{line: 239, col: 11, offset: 6479},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 239, col: 11, offset: 6479},
	val: "/\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 239, col: 19, offset: 6487},
	val: "∧",
	ignoreCase: false,
},
	},
},
},
{
	name: "CombineTypes",
	pos: position{line: 240, col: 1, offset: 6493},
	expr: &choiceExpr{
	pos: position{line: 240, col: 16, offset: 6510},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 240, col: 16, offset: 6510},
	val: "//\\\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 240, col: 27, offset: 6521},
	val: "⩓",
	ignoreCase: false,
},
	},
},
},
{
	name: "Prefer",
	pos: position{line: 241, col: 1, offset: 6527},
	expr: &choiceExpr{
	pos: position{line: 241, col: 10, offset: 6538},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 241, col: 10, offset: 6538},
	val: "//",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 241, col: 17, offset: 6545},
	val: "⫽",
	ignoreCase: false,
},
	},
},
},
{
	name: "Lambda",
	pos: position{line: 242, col: 1, offset: 6551},
	expr: &choiceExpr{
	pos: position{line: 242, col: 10, offset: 6562},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 242, col: 10, offset: 6562},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 242, col: 17, offset: 6569},
	val: "λ",
	ignoreCase: false,
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 243, col: 1, offset: 6574},
	expr: &choiceExpr{
	pos: position{line: 243, col: 10, offset: 6585},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 243, col: 10, offset: 6585},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 243, col: 21, offset: 6596},
	val: "∀",
	ignoreCase: false,
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 244, col: 1, offset: 6602},
	expr: &choiceExpr{
	pos: position{line: 244, col: 9, offset: 6612},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 244, col: 9, offset: 6612},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 244, col: 16, offset: 6619},
	val: "→",
	ignoreCase: false,
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 246, col: 1, offset: 6626},
	expr: &seqExpr{
	pos: position{line: 246, col: 12, offset: 6639},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 246, col: 12, offset: 6639},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 246, col: 17, offset: 6644},
	expr: &charClassMatcher{
	pos: position{line: 246, col: 17, offset: 6644},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 246, col: 23, offset: 6650},
	expr: &ruleRefExpr{
	pos: position{line: 246, col: 23, offset: 6650},
	name: "Digit",
},
},
	},
},
},
{
	name: "NumericDoubleLiteral",
	pos: position{line: 248, col: 1, offset: 6658},
	expr: &actionExpr{
	pos: position{line: 248, col: 24, offset: 6683},
	run: (*parser).callonNumericDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 248, col: 24, offset: 6683},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 248, col: 24, offset: 6683},
	expr: &charClassMatcher{
	pos: position{line: 248, col: 24, offset: 6683},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 248, col: 30, offset: 6689},
	expr: &ruleRefExpr{
	pos: position{line: 248, col: 30, offset: 6689},
	name: "Digit",
},
},
&choiceExpr{
	pos: position{line: 248, col: 39, offset: 6698},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 248, col: 39, offset: 6698},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 248, col: 39, offset: 6698},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 248, col: 43, offset: 6702},
	expr: &ruleRefExpr{
	pos: position{line: 248, col: 43, offset: 6702},
	name: "Digit",
},
},
&zeroOrOneExpr{
	pos: position{line: 248, col: 50, offset: 6709},
	expr: &ruleRefExpr{
	pos: position{line: 248, col: 50, offset: 6709},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 248, col: 62, offset: 6721},
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
	pos: position{line: 256, col: 1, offset: 6877},
	expr: &choiceExpr{
	pos: position{line: 256, col: 17, offset: 6895},
	alternatives: []interface{}{
&labeledExpr{
	pos: position{line: 256, col: 17, offset: 6895},
	label: "d",
	expr: &ruleRefExpr{
	pos: position{line: 256, col: 19, offset: 6897},
	name: "NumericDoubleLiteral",
},
},
&actionExpr{
	pos: position{line: 257, col: 5, offset: 6922},
	run: (*parser).callonDoubleLiteral4,
	expr: &ruleRefExpr{
	pos: position{line: 257, col: 5, offset: 6922},
	name: "Infinity",
},
},
&actionExpr{
	pos: position{line: 258, col: 5, offset: 6974},
	run: (*parser).callonDoubleLiteral6,
	expr: &seqExpr{
	pos: position{line: 258, col: 5, offset: 6974},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 258, col: 5, offset: 6974},
	val: "-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 258, col: 9, offset: 6978},
	name: "Infinity",
},
	},
},
},
&actionExpr{
	pos: position{line: 259, col: 5, offset: 7031},
	run: (*parser).callonDoubleLiteral10,
	expr: &ruleRefExpr{
	pos: position{line: 259, col: 5, offset: 7031},
	name: "NaN",
},
},
	},
},
},
{
	name: "NaturalLiteral",
	pos: position{line: 261, col: 1, offset: 7074},
	expr: &actionExpr{
	pos: position{line: 261, col: 18, offset: 7093},
	run: (*parser).callonNaturalLiteral1,
	expr: &oneOrMoreExpr{
	pos: position{line: 261, col: 18, offset: 7093},
	expr: &ruleRefExpr{
	pos: position{line: 261, col: 18, offset: 7093},
	name: "Digit",
},
},
},
},
{
	name: "IntegerLiteral",
	pos: position{line: 266, col: 1, offset: 7182},
	expr: &actionExpr{
	pos: position{line: 266, col: 18, offset: 7201},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 266, col: 18, offset: 7201},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 266, col: 18, offset: 7201},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 266, col: 22, offset: 7205},
	name: "NaturalLiteral",
},
	},
},
},
},
{
	name: "DeBruijn",
	pos: position{line: 274, col: 1, offset: 7357},
	expr: &actionExpr{
	pos: position{line: 274, col: 12, offset: 7370},
	run: (*parser).callonDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 274, col: 12, offset: 7370},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 274, col: 12, offset: 7370},
	name: "_",
},
&litMatcher{
	pos: position{line: 274, col: 14, offset: 7372},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 274, col: 18, offset: 7376},
	name: "_",
},
&labeledExpr{
	pos: position{line: 274, col: 20, offset: 7378},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 274, col: 26, offset: 7384},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Variable",
	pos: position{line: 276, col: 1, offset: 7440},
	expr: &actionExpr{
	pos: position{line: 276, col: 12, offset: 7453},
	run: (*parser).callonVariable1,
	expr: &seqExpr{
	pos: position{line: 276, col: 12, offset: 7453},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 276, col: 12, offset: 7453},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 276, col: 17, offset: 7458},
	name: "NonreservedLabel",
},
},
&labeledExpr{
	pos: position{line: 276, col: 34, offset: 7475},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 276, col: 40, offset: 7481},
	expr: &ruleRefExpr{
	pos: position{line: 276, col: 40, offset: 7481},
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
	pos: position{line: 284, col: 1, offset: 7644},
	expr: &choiceExpr{
	pos: position{line: 284, col: 14, offset: 7659},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 284, col: 14, offset: 7659},
	name: "Variable",
},
&ruleRefExpr{
	pos: position{line: 284, col: 25, offset: 7670},
	name: "Reserved",
},
	},
},
},
{
	name: "PathCharacter",
	pos: position{line: 286, col: 1, offset: 7680},
	expr: &choiceExpr{
	pos: position{line: 287, col: 6, offset: 7703},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 287, col: 6, offset: 7703},
	val: "!",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 288, col: 6, offset: 7715},
	val: "[\\x24-\\x27]",
	ranges: []rune{'$','\'',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 289, col: 6, offset: 7732},
	val: "[\\x2a-\\x2b]",
	ranges: []rune{'*','+',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 290, col: 6, offset: 7749},
	val: "[\\x2d-\\x2e]",
	ranges: []rune{'-','.',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 291, col: 6, offset: 7766},
	val: "[\\x30-\\x3b]",
	ranges: []rune{'0',';',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 292, col: 6, offset: 7783},
	val: "=",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 293, col: 6, offset: 7795},
	val: "[\\x40-\\x5a]",
	ranges: []rune{'@','Z',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 294, col: 6, offset: 7812},
	val: "[\\x5e-\\x7a]",
	ranges: []rune{'^','z',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 295, col: 6, offset: 7829},
	val: "|",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 296, col: 6, offset: 7841},
	val: "~",
	ignoreCase: false,
},
	},
},
},
{
	name: "UnquotedPathComponent",
	pos: position{line: 298, col: 1, offset: 7849},
	expr: &actionExpr{
	pos: position{line: 298, col: 25, offset: 7875},
	run: (*parser).callonUnquotedPathComponent1,
	expr: &oneOrMoreExpr{
	pos: position{line: 298, col: 25, offset: 7875},
	expr: &ruleRefExpr{
	pos: position{line: 298, col: 25, offset: 7875},
	name: "PathCharacter",
},
},
},
},
{
	name: "PathComponent",
	pos: position{line: 300, col: 1, offset: 7922},
	expr: &actionExpr{
	pos: position{line: 300, col: 17, offset: 7940},
	run: (*parser).callonPathComponent1,
	expr: &seqExpr{
	pos: position{line: 300, col: 17, offset: 7940},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 300, col: 17, offset: 7940},
	val: "/",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 300, col: 21, offset: 7944},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 300, col: 23, offset: 7946},
	name: "UnquotedPathComponent",
},
},
	},
},
},
},
{
	name: "Path",
	pos: position{line: 302, col: 1, offset: 7987},
	expr: &actionExpr{
	pos: position{line: 302, col: 8, offset: 7996},
	run: (*parser).callonPath1,
	expr: &labeledExpr{
	pos: position{line: 302, col: 8, offset: 7996},
	label: "cs",
	expr: &oneOrMoreExpr{
	pos: position{line: 302, col: 11, offset: 7999},
	expr: &ruleRefExpr{
	pos: position{line: 302, col: 11, offset: 7999},
	name: "PathComponent",
},
},
},
},
},
{
	name: "Local",
	pos: position{line: 311, col: 1, offset: 8273},
	expr: &choiceExpr{
	pos: position{line: 311, col: 9, offset: 8283},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 311, col: 9, offset: 8283},
	name: "ParentPath",
},
&ruleRefExpr{
	pos: position{line: 311, col: 22, offset: 8296},
	name: "HerePath",
},
&ruleRefExpr{
	pos: position{line: 311, col: 33, offset: 8307},
	name: "HomePath",
},
&ruleRefExpr{
	pos: position{line: 311, col: 44, offset: 8318},
	name: "AbsolutePath",
},
	},
},
},
{
	name: "ParentPath",
	pos: position{line: 313, col: 1, offset: 8332},
	expr: &actionExpr{
	pos: position{line: 313, col: 14, offset: 8347},
	run: (*parser).callonParentPath1,
	expr: &seqExpr{
	pos: position{line: 313, col: 14, offset: 8347},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 313, col: 14, offset: 8347},
	val: "..",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 313, col: 19, offset: 8352},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 313, col: 21, offset: 8354},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HerePath",
	pos: position{line: 314, col: 1, offset: 8410},
	expr: &actionExpr{
	pos: position{line: 314, col: 12, offset: 8423},
	run: (*parser).callonHerePath1,
	expr: &seqExpr{
	pos: position{line: 314, col: 12, offset: 8423},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 314, col: 12, offset: 8423},
	val: ".",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 314, col: 16, offset: 8427},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 314, col: 18, offset: 8429},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HomePath",
	pos: position{line: 315, col: 1, offset: 8468},
	expr: &actionExpr{
	pos: position{line: 315, col: 12, offset: 8481},
	run: (*parser).callonHomePath1,
	expr: &seqExpr{
	pos: position{line: 315, col: 12, offset: 8481},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 315, col: 12, offset: 8481},
	val: "~",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 315, col: 16, offset: 8485},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 315, col: 18, offset: 8487},
	name: "Path",
},
},
	},
},
},
},
{
	name: "AbsolutePath",
	pos: position{line: 316, col: 1, offset: 8542},
	expr: &actionExpr{
	pos: position{line: 316, col: 16, offset: 8559},
	run: (*parser).callonAbsolutePath1,
	expr: &labeledExpr{
	pos: position{line: 316, col: 16, offset: 8559},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 316, col: 18, offset: 8561},
	name: "Path",
},
},
},
},
{
	name: "Scheme",
	pos: position{line: 318, col: 1, offset: 8617},
	expr: &seqExpr{
	pos: position{line: 318, col: 10, offset: 8628},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 318, col: 10, offset: 8628},
	val: "http",
	ignoreCase: false,
},
&zeroOrOneExpr{
	pos: position{line: 318, col: 17, offset: 8635},
	expr: &litMatcher{
	pos: position{line: 318, col: 17, offset: 8635},
	val: "s",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "HttpRaw",
	pos: position{line: 320, col: 1, offset: 8641},
	expr: &actionExpr{
	pos: position{line: 320, col: 11, offset: 8653},
	run: (*parser).callonHttpRaw1,
	expr: &seqExpr{
	pos: position{line: 320, col: 11, offset: 8653},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 320, col: 11, offset: 8653},
	name: "Scheme",
},
&litMatcher{
	pos: position{line: 320, col: 18, offset: 8660},
	val: "://",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 320, col: 24, offset: 8666},
	name: "Authority",
},
&ruleRefExpr{
	pos: position{line: 320, col: 34, offset: 8676},
	name: "Path",
},
&zeroOrOneExpr{
	pos: position{line: 320, col: 39, offset: 8681},
	expr: &seqExpr{
	pos: position{line: 320, col: 41, offset: 8683},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 320, col: 41, offset: 8683},
	val: "?",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 320, col: 45, offset: 8687},
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
	pos: position{line: 322, col: 1, offset: 8744},
	expr: &seqExpr{
	pos: position{line: 322, col: 13, offset: 8758},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 322, col: 13, offset: 8758},
	expr: &seqExpr{
	pos: position{line: 322, col: 14, offset: 8759},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 322, col: 14, offset: 8759},
	name: "Userinfo",
},
&litMatcher{
	pos: position{line: 322, col: 23, offset: 8768},
	val: "@",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 322, col: 29, offset: 8774},
	name: "Host",
},
&zeroOrOneExpr{
	pos: position{line: 322, col: 34, offset: 8779},
	expr: &seqExpr{
	pos: position{line: 322, col: 35, offset: 8780},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 322, col: 35, offset: 8780},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 322, col: 39, offset: 8784},
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
	pos: position{line: 324, col: 1, offset: 8792},
	expr: &zeroOrMoreExpr{
	pos: position{line: 324, col: 12, offset: 8805},
	expr: &choiceExpr{
	pos: position{line: 324, col: 14, offset: 8807},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 324, col: 14, offset: 8807},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 324, col: 27, offset: 8820},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 324, col: 40, offset: 8833},
	name: "SubDelims",
},
&litMatcher{
	pos: position{line: 324, col: 52, offset: 8845},
	val: ":",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "Host",
	pos: position{line: 326, col: 1, offset: 8853},
	expr: &choiceExpr{
	pos: position{line: 326, col: 8, offset: 8862},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 326, col: 8, offset: 8862},
	name: "IPLiteral",
},
&ruleRefExpr{
	pos: position{line: 326, col: 20, offset: 8874},
	name: "RegName",
},
	},
},
},
{
	name: "Port",
	pos: position{line: 328, col: 1, offset: 8883},
	expr: &zeroOrMoreExpr{
	pos: position{line: 328, col: 8, offset: 8892},
	expr: &ruleRefExpr{
	pos: position{line: 328, col: 8, offset: 8892},
	name: "Digit",
},
},
},
{
	name: "IPLiteral",
	pos: position{line: 330, col: 1, offset: 8900},
	expr: &seqExpr{
	pos: position{line: 330, col: 13, offset: 8914},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 330, col: 13, offset: 8914},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 330, col: 17, offset: 8918},
	name: "IPv6address",
},
&litMatcher{
	pos: position{line: 330, col: 29, offset: 8930},
	val: "]",
	ignoreCase: false,
},
	},
},
},
{
	name: "IPv6address",
	pos: position{line: 332, col: 1, offset: 8935},
	expr: &actionExpr{
	pos: position{line: 332, col: 15, offset: 8951},
	run: (*parser).callonIPv6address1,
	expr: &seqExpr{
	pos: position{line: 332, col: 15, offset: 8951},
	exprs: []interface{}{
&zeroOrMoreExpr{
	pos: position{line: 332, col: 15, offset: 8951},
	expr: &ruleRefExpr{
	pos: position{line: 332, col: 16, offset: 8952},
	name: "HexDig",
},
},
&litMatcher{
	pos: position{line: 332, col: 25, offset: 8961},
	val: ":",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 332, col: 29, offset: 8965},
	expr: &choiceExpr{
	pos: position{line: 332, col: 30, offset: 8966},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 332, col: 30, offset: 8966},
	name: "HexDig",
},
&litMatcher{
	pos: position{line: 332, col: 39, offset: 8975},
	val: ":",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 332, col: 45, offset: 8981},
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
	pos: position{line: 338, col: 1, offset: 9135},
	expr: &zeroOrMoreExpr{
	pos: position{line: 338, col: 11, offset: 9147},
	expr: &choiceExpr{
	pos: position{line: 338, col: 12, offset: 9148},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 338, col: 12, offset: 9148},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 338, col: 25, offset: 9161},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 338, col: 38, offset: 9174},
	name: "SubDelims",
},
	},
},
},
},
{
	name: "PChar",
	pos: position{line: 340, col: 1, offset: 9187},
	expr: &choiceExpr{
	pos: position{line: 340, col: 9, offset: 9197},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 340, col: 9, offset: 9197},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 340, col: 22, offset: 9210},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 340, col: 35, offset: 9223},
	name: "SubDelims",
},
&charClassMatcher{
	pos: position{line: 340, col: 47, offset: 9235},
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
	pos: position{line: 342, col: 1, offset: 9241},
	expr: &zeroOrMoreExpr{
	pos: position{line: 342, col: 9, offset: 9251},
	expr: &choiceExpr{
	pos: position{line: 342, col: 10, offset: 9252},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 342, col: 10, offset: 9252},
	name: "PChar",
},
&charClassMatcher{
	pos: position{line: 342, col: 18, offset: 9260},
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
	pos: position{line: 344, col: 1, offset: 9268},
	expr: &seqExpr{
	pos: position{line: 344, col: 14, offset: 9283},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 344, col: 14, offset: 9283},
	val: "%",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 344, col: 18, offset: 9287},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 344, col: 25, offset: 9294},
	name: "HexDig",
},
	},
},
},
{
	name: "Unreserved",
	pos: position{line: 346, col: 1, offset: 9302},
	expr: &charClassMatcher{
	pos: position{line: 346, col: 14, offset: 9317},
	val: "[A-Za-z0-9._~-]",
	chars: []rune{'.','_','~','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SubDelims",
	pos: position{line: 348, col: 1, offset: 9334},
	expr: &choiceExpr{
	pos: position{line: 348, col: 13, offset: 9348},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 348, col: 13, offset: 9348},
	val: "!",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 348, col: 19, offset: 9354},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 348, col: 25, offset: 9360},
	val: "&",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 348, col: 31, offset: 9366},
	val: "'",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 348, col: 37, offset: 9372},
	val: "(",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 348, col: 43, offset: 9378},
	val: ")",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 348, col: 49, offset: 9384},
	val: "*",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 348, col: 55, offset: 9390},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 348, col: 61, offset: 9396},
	val: ",",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 348, col: 67, offset: 9402},
	val: ";",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 348, col: 73, offset: 9408},
	val: "=",
	ignoreCase: false,
},
	},
},
},
{
	name: "Http",
	pos: position{line: 350, col: 1, offset: 9413},
	expr: &actionExpr{
	pos: position{line: 350, col: 8, offset: 9422},
	run: (*parser).callonHttp1,
	expr: &labeledExpr{
	pos: position{line: 350, col: 8, offset: 9422},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 350, col: 10, offset: 9424},
	name: "HttpRaw",
},
},
},
},
{
	name: "Env",
	pos: position{line: 352, col: 1, offset: 9469},
	expr: &actionExpr{
	pos: position{line: 352, col: 7, offset: 9477},
	run: (*parser).callonEnv1,
	expr: &seqExpr{
	pos: position{line: 352, col: 7, offset: 9477},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 352, col: 7, offset: 9477},
	val: "env:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 352, col: 14, offset: 9484},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 352, col: 17, offset: 9487},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 352, col: 17, offset: 9487},
	name: "BashEnvironmentVariable",
},
&ruleRefExpr{
	pos: position{line: 352, col: 43, offset: 9513},
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
	pos: position{line: 354, col: 1, offset: 9558},
	expr: &actionExpr{
	pos: position{line: 354, col: 27, offset: 9586},
	run: (*parser).callonBashEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 354, col: 27, offset: 9586},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 354, col: 27, offset: 9586},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 354, col: 36, offset: 9595},
	expr: &charClassMatcher{
	pos: position{line: 354, col: 36, offset: 9595},
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
	pos: position{line: 358, col: 1, offset: 9651},
	expr: &actionExpr{
	pos: position{line: 358, col: 28, offset: 9680},
	run: (*parser).callonPosixEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 358, col: 28, offset: 9680},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 358, col: 28, offset: 9680},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 358, col: 32, offset: 9684},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 358, col: 34, offset: 9686},
	name: "PosixEnvironmentVariableContent",
},
},
&litMatcher{
	pos: position{line: 358, col: 66, offset: 9718},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PosixEnvironmentVariableContent",
	pos: position{line: 362, col: 1, offset: 9743},
	expr: &actionExpr{
	pos: position{line: 362, col: 35, offset: 9779},
	run: (*parser).callonPosixEnvironmentVariableContent1,
	expr: &labeledExpr{
	pos: position{line: 362, col: 35, offset: 9779},
	label: "v",
	expr: &oneOrMoreExpr{
	pos: position{line: 362, col: 37, offset: 9781},
	expr: &ruleRefExpr{
	pos: position{line: 362, col: 37, offset: 9781},
	name: "PosixEnvironmentVariableCharacter",
},
},
},
},
},
{
	name: "PosixEnvironmentVariableCharacter",
	pos: position{line: 371, col: 1, offset: 9994},
	expr: &choiceExpr{
	pos: position{line: 372, col: 7, offset: 10038},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 372, col: 7, offset: 10038},
	run: (*parser).callonPosixEnvironmentVariableCharacter2,
	expr: &litMatcher{
	pos: position{line: 372, col: 7, offset: 10038},
	val: "\\\"",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 373, col: 7, offset: 10078},
	run: (*parser).callonPosixEnvironmentVariableCharacter4,
	expr: &litMatcher{
	pos: position{line: 373, col: 7, offset: 10078},
	val: "\\\\",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 374, col: 7, offset: 10118},
	run: (*parser).callonPosixEnvironmentVariableCharacter6,
	expr: &litMatcher{
	pos: position{line: 374, col: 7, offset: 10118},
	val: "\\a",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 375, col: 7, offset: 10158},
	run: (*parser).callonPosixEnvironmentVariableCharacter8,
	expr: &litMatcher{
	pos: position{line: 375, col: 7, offset: 10158},
	val: "\\b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 376, col: 7, offset: 10198},
	run: (*parser).callonPosixEnvironmentVariableCharacter10,
	expr: &litMatcher{
	pos: position{line: 376, col: 7, offset: 10198},
	val: "\\f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 377, col: 7, offset: 10238},
	run: (*parser).callonPosixEnvironmentVariableCharacter12,
	expr: &litMatcher{
	pos: position{line: 377, col: 7, offset: 10238},
	val: "\\n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 378, col: 7, offset: 10278},
	run: (*parser).callonPosixEnvironmentVariableCharacter14,
	expr: &litMatcher{
	pos: position{line: 378, col: 7, offset: 10278},
	val: "\\r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 379, col: 7, offset: 10318},
	run: (*parser).callonPosixEnvironmentVariableCharacter16,
	expr: &litMatcher{
	pos: position{line: 379, col: 7, offset: 10318},
	val: "\\t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 380, col: 7, offset: 10358},
	run: (*parser).callonPosixEnvironmentVariableCharacter18,
	expr: &litMatcher{
	pos: position{line: 380, col: 7, offset: 10358},
	val: "\\v",
	ignoreCase: false,
},
},
&charClassMatcher{
	pos: position{line: 381, col: 7, offset: 10398},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 382, col: 7, offset: 10416},
	val: "[\\x23-\\x3c]",
	ranges: []rune{'#','<',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 383, col: 7, offset: 10434},
	val: "[\\x3e-\\x5b]",
	ranges: []rune{'>','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 384, col: 7, offset: 10452},
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
	pos: position{line: 386, col: 1, offset: 10465},
	expr: &choiceExpr{
	pos: position{line: 386, col: 14, offset: 10480},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 386, col: 14, offset: 10480},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 386, col: 24, offset: 10490},
	name: "Local",
},
&ruleRefExpr{
	pos: position{line: 386, col: 32, offset: 10498},
	name: "Http",
},
&ruleRefExpr{
	pos: position{line: 386, col: 39, offset: 10505},
	name: "Env",
},
	},
},
},
{
	name: "HashValue",
	pos: position{line: 389, col: 1, offset: 10578},
	expr: &actionExpr{
	pos: position{line: 389, col: 13, offset: 10590},
	run: (*parser).callonHashValue1,
	expr: &seqExpr{
	pos: position{line: 389, col: 13, offset: 10590},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 389, col: 13, offset: 10590},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 389, col: 20, offset: 10597},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 389, col: 27, offset: 10604},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 389, col: 34, offset: 10611},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 389, col: 41, offset: 10618},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 389, col: 48, offset: 10625},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 389, col: 55, offset: 10632},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 389, col: 62, offset: 10639},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 390, col: 13, offset: 10658},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 390, col: 20, offset: 10665},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 390, col: 27, offset: 10672},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 390, col: 34, offset: 10679},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 390, col: 41, offset: 10686},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 390, col: 48, offset: 10693},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 390, col: 55, offset: 10700},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 390, col: 62, offset: 10707},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 391, col: 13, offset: 10726},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 391, col: 20, offset: 10733},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 391, col: 27, offset: 10740},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 391, col: 34, offset: 10747},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 391, col: 41, offset: 10754},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 391, col: 48, offset: 10761},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 391, col: 55, offset: 10768},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 391, col: 62, offset: 10775},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 392, col: 13, offset: 10794},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 392, col: 20, offset: 10801},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 392, col: 27, offset: 10808},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 392, col: 34, offset: 10815},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 392, col: 41, offset: 10822},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 392, col: 48, offset: 10829},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 392, col: 55, offset: 10836},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 392, col: 62, offset: 10843},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 393, col: 13, offset: 10862},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 393, col: 20, offset: 10869},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 393, col: 27, offset: 10876},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 393, col: 34, offset: 10883},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 393, col: 41, offset: 10890},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 393, col: 48, offset: 10897},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 393, col: 55, offset: 10904},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 393, col: 62, offset: 10911},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 394, col: 13, offset: 10930},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 394, col: 20, offset: 10937},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 394, col: 27, offset: 10944},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 394, col: 34, offset: 10951},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 394, col: 41, offset: 10958},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 394, col: 48, offset: 10965},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 394, col: 55, offset: 10972},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 394, col: 62, offset: 10979},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 395, col: 13, offset: 10998},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 395, col: 20, offset: 11005},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 395, col: 27, offset: 11012},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 395, col: 34, offset: 11019},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 395, col: 41, offset: 11026},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 395, col: 48, offset: 11033},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 395, col: 55, offset: 11040},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 395, col: 62, offset: 11047},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 396, col: 13, offset: 11066},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 396, col: 20, offset: 11073},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 396, col: 27, offset: 11080},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 396, col: 34, offset: 11087},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 396, col: 41, offset: 11094},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 396, col: 48, offset: 11101},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 396, col: 55, offset: 11108},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 396, col: 62, offset: 11115},
	name: "HexDig",
},
	},
},
},
},
{
	name: "Hash",
	pos: position{line: 402, col: 1, offset: 11259},
	expr: &actionExpr{
	pos: position{line: 402, col: 8, offset: 11266},
	run: (*parser).callonHash1,
	expr: &seqExpr{
	pos: position{line: 402, col: 8, offset: 11266},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 402, col: 8, offset: 11266},
	val: "sha256:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 402, col: 18, offset: 11276},
	label: "val",
	expr: &ruleRefExpr{
	pos: position{line: 402, col: 22, offset: 11280},
	name: "HashValue",
},
},
	},
},
},
},
{
	name: "ImportHashed",
	pos: position{line: 404, col: 1, offset: 11350},
	expr: &actionExpr{
	pos: position{line: 404, col: 16, offset: 11367},
	run: (*parser).callonImportHashed1,
	expr: &seqExpr{
	pos: position{line: 404, col: 16, offset: 11367},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 404, col: 16, offset: 11367},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 404, col: 18, offset: 11369},
	name: "ImportType",
},
},
&labeledExpr{
	pos: position{line: 404, col: 29, offset: 11380},
	label: "h",
	expr: &zeroOrOneExpr{
	pos: position{line: 404, col: 31, offset: 11382},
	expr: &seqExpr{
	pos: position{line: 404, col: 32, offset: 11383},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 404, col: 32, offset: 11383},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 404, col: 35, offset: 11386},
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
	pos: position{line: 412, col: 1, offset: 11541},
	expr: &choiceExpr{
	pos: position{line: 412, col: 10, offset: 11552},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 412, col: 10, offset: 11552},
	run: (*parser).callonImport2,
	expr: &seqExpr{
	pos: position{line: 412, col: 10, offset: 11552},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 412, col: 10, offset: 11552},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 412, col: 12, offset: 11554},
	name: "ImportHashed",
},
},
&ruleRefExpr{
	pos: position{line: 412, col: 25, offset: 11567},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 412, col: 27, offset: 11569},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 412, col: 30, offset: 11572},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 412, col: 33, offset: 11575},
	name: "Text",
},
	},
},
},
&actionExpr{
	pos: position{line: 413, col: 10, offset: 11672},
	run: (*parser).callonImport10,
	expr: &labeledExpr{
	pos: position{line: 413, col: 10, offset: 11672},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 413, col: 12, offset: 11674},
	name: "ImportHashed",
},
},
},
	},
},
},
{
	name: "LetBinding",
	pos: position{line: 416, col: 1, offset: 11769},
	expr: &actionExpr{
	pos: position{line: 416, col: 14, offset: 11784},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 416, col: 14, offset: 11784},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 416, col: 14, offset: 11784},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 416, col: 18, offset: 11788},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 416, col: 21, offset: 11791},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 416, col: 27, offset: 11797},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 416, col: 44, offset: 11814},
	name: "_",
},
&labeledExpr{
	pos: position{line: 416, col: 46, offset: 11816},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 416, col: 48, offset: 11818},
	expr: &seqExpr{
	pos: position{line: 416, col: 49, offset: 11819},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 416, col: 49, offset: 11819},
	name: "Annotation",
},
&ruleRefExpr{
	pos: position{line: 416, col: 60, offset: 11830},
	name: "_",
},
	},
},
},
},
&litMatcher{
	pos: position{line: 417, col: 13, offset: 11846},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 417, col: 17, offset: 11850},
	name: "_",
},
&labeledExpr{
	pos: position{line: 417, col: 19, offset: 11852},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 417, col: 21, offset: 11854},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 417, col: 32, offset: 11865},
	name: "_",
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 432, col: 1, offset: 12174},
	expr: &choiceExpr{
	pos: position{line: 433, col: 7, offset: 12195},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 433, col: 7, offset: 12195},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 433, col: 7, offset: 12195},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 433, col: 7, offset: 12195},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 433, col: 14, offset: 12202},
	name: "_",
},
&litMatcher{
	pos: position{line: 433, col: 16, offset: 12204},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 433, col: 20, offset: 12208},
	name: "_",
},
&labeledExpr{
	pos: position{line: 433, col: 22, offset: 12210},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 433, col: 28, offset: 12216},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 433, col: 45, offset: 12233},
	name: "_",
},
&litMatcher{
	pos: position{line: 433, col: 47, offset: 12235},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 433, col: 51, offset: 12239},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 433, col: 54, offset: 12242},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 433, col: 56, offset: 12244},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 433, col: 67, offset: 12255},
	name: "_",
},
&litMatcher{
	pos: position{line: 433, col: 69, offset: 12257},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 433, col: 73, offset: 12261},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 433, col: 75, offset: 12263},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 433, col: 81, offset: 12269},
	name: "_",
},
&labeledExpr{
	pos: position{line: 433, col: 83, offset: 12271},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 433, col: 88, offset: 12276},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 436, col: 7, offset: 12393},
	run: (*parser).callonExpression22,
	expr: &seqExpr{
	pos: position{line: 436, col: 7, offset: 12393},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 436, col: 7, offset: 12393},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 436, col: 10, offset: 12396},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 436, col: 13, offset: 12399},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 436, col: 18, offset: 12404},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 436, col: 29, offset: 12415},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 436, col: 31, offset: 12417},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 436, col: 36, offset: 12422},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 436, col: 39, offset: 12425},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 436, col: 41, offset: 12427},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 436, col: 52, offset: 12438},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 436, col: 54, offset: 12440},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 436, col: 59, offset: 12445},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 436, col: 62, offset: 12448},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 436, col: 64, offset: 12450},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 439, col: 7, offset: 12536},
	run: (*parser).callonExpression38,
	expr: &seqExpr{
	pos: position{line: 439, col: 7, offset: 12536},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 439, col: 7, offset: 12536},
	label: "bindings",
	expr: &oneOrMoreExpr{
	pos: position{line: 439, col: 16, offset: 12545},
	expr: &ruleRefExpr{
	pos: position{line: 439, col: 16, offset: 12545},
	name: "LetBinding",
},
},
},
&ruleRefExpr{
	pos: position{line: 439, col: 28, offset: 12557},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 439, col: 31, offset: 12560},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 439, col: 34, offset: 12563},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 439, col: 36, offset: 12565},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 446, col: 7, offset: 12805},
	run: (*parser).callonExpression47,
	expr: &seqExpr{
	pos: position{line: 446, col: 7, offset: 12805},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 446, col: 7, offset: 12805},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 446, col: 14, offset: 12812},
	name: "_",
},
&litMatcher{
	pos: position{line: 446, col: 16, offset: 12814},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 446, col: 20, offset: 12818},
	name: "_",
},
&labeledExpr{
	pos: position{line: 446, col: 22, offset: 12820},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 446, col: 28, offset: 12826},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 446, col: 45, offset: 12843},
	name: "_",
},
&litMatcher{
	pos: position{line: 446, col: 47, offset: 12845},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 446, col: 51, offset: 12849},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 446, col: 54, offset: 12852},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 446, col: 56, offset: 12854},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 446, col: 67, offset: 12865},
	name: "_",
},
&litMatcher{
	pos: position{line: 446, col: 69, offset: 12867},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 446, col: 73, offset: 12871},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 446, col: 75, offset: 12873},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 446, col: 81, offset: 12879},
	name: "_",
},
&labeledExpr{
	pos: position{line: 446, col: 83, offset: 12881},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 446, col: 88, offset: 12886},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 449, col: 7, offset: 12995},
	run: (*parser).callonExpression67,
	expr: &seqExpr{
	pos: position{line: 449, col: 7, offset: 12995},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 449, col: 7, offset: 12995},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 449, col: 9, offset: 12997},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 449, col: 28, offset: 13016},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 449, col: 30, offset: 13018},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 449, col: 36, offset: 13024},
	name: "_",
},
&labeledExpr{
	pos: position{line: 449, col: 38, offset: 13026},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 449, col: 40, offset: 13028},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 450, col: 7, offset: 13087},
	run: (*parser).callonExpression76,
	expr: &seqExpr{
	pos: position{line: 450, col: 7, offset: 13087},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 450, col: 7, offset: 13087},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 450, col: 13, offset: 13093},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 450, col: 16, offset: 13096},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 450, col: 18, offset: 13098},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 450, col: 35, offset: 13115},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 450, col: 38, offset: 13118},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 450, col: 40, offset: 13120},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 450, col: 57, offset: 13137},
	name: "_",
},
&litMatcher{
	pos: position{line: 450, col: 59, offset: 13139},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 450, col: 63, offset: 13143},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 450, col: 66, offset: 13146},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 450, col: 68, offset: 13148},
	name: "ApplicationExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 453, col: 7, offset: 13269},
	name: "EmptyList",
},
&ruleRefExpr{
	pos: position{line: 454, col: 7, offset: 13285},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 456, col: 1, offset: 13306},
	expr: &actionExpr{
	pos: position{line: 456, col: 14, offset: 13321},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 456, col: 14, offset: 13321},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 456, col: 14, offset: 13321},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 456, col: 18, offset: 13325},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 456, col: 21, offset: 13328},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 456, col: 23, offset: 13330},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 458, col: 1, offset: 13360},
	expr: &actionExpr{
	pos: position{line: 459, col: 1, offset: 13384},
	run: (*parser).callonAnnotatedExpression1,
	expr: &seqExpr{
	pos: position{line: 459, col: 1, offset: 13384},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 459, col: 1, offset: 13384},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 459, col: 3, offset: 13386},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 459, col: 22, offset: 13405},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 459, col: 24, offset: 13407},
	expr: &seqExpr{
	pos: position{line: 459, col: 25, offset: 13408},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 459, col: 25, offset: 13408},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 459, col: 27, offset: 13410},
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
	pos: position{line: 464, col: 1, offset: 13535},
	expr: &actionExpr{
	pos: position{line: 464, col: 13, offset: 13549},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 464, col: 13, offset: 13549},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 464, col: 13, offset: 13549},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 464, col: 17, offset: 13553},
	name: "_",
},
&litMatcher{
	pos: position{line: 464, col: 19, offset: 13555},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 464, col: 23, offset: 13559},
	name: "_",
},
&litMatcher{
	pos: position{line: 464, col: 25, offset: 13561},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 464, col: 29, offset: 13565},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 464, col: 32, offset: 13568},
	name: "List",
},
&ruleRefExpr{
	pos: position{line: 464, col: 37, offset: 13573},
	name: "_",
},
&labeledExpr{
	pos: position{line: 464, col: 39, offset: 13575},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 464, col: 41, offset: 13577},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 468, col: 1, offset: 13640},
	expr: &ruleRefExpr{
	pos: position{line: 468, col: 22, offset: 13663},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 470, col: 1, offset: 13684},
	expr: &actionExpr{
	pos: position{line: 470, col: 26, offset: 13711},
	run: (*parser).callonImportAltExpression1,
	expr: &seqExpr{
	pos: position{line: 470, col: 26, offset: 13711},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 470, col: 26, offset: 13711},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 470, col: 32, offset: 13717},
	name: "OrExpression",
},
},
&labeledExpr{
	pos: position{line: 470, col: 55, offset: 13740},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 470, col: 60, offset: 13745},
	expr: &seqExpr{
	pos: position{line: 470, col: 61, offset: 13746},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 470, col: 61, offset: 13746},
	name: "_",
},
&litMatcher{
	pos: position{line: 470, col: 63, offset: 13748},
	val: "?",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 470, col: 67, offset: 13752},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 470, col: 69, offset: 13754},
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
	pos: position{line: 472, col: 1, offset: 13825},
	expr: &actionExpr{
	pos: position{line: 472, col: 26, offset: 13852},
	run: (*parser).callonOrExpression1,
	expr: &seqExpr{
	pos: position{line: 472, col: 26, offset: 13852},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 472, col: 26, offset: 13852},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 472, col: 32, offset: 13858},
	name: "PlusExpression",
},
},
&labeledExpr{
	pos: position{line: 472, col: 55, offset: 13881},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 472, col: 60, offset: 13886},
	expr: &seqExpr{
	pos: position{line: 472, col: 61, offset: 13887},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 472, col: 61, offset: 13887},
	name: "_",
},
&litMatcher{
	pos: position{line: 472, col: 63, offset: 13889},
	val: "||",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 472, col: 68, offset: 13894},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 472, col: 70, offset: 13896},
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
	pos: position{line: 474, col: 1, offset: 13962},
	expr: &actionExpr{
	pos: position{line: 474, col: 26, offset: 13989},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 474, col: 26, offset: 13989},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 474, col: 26, offset: 13989},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 474, col: 32, offset: 13995},
	name: "TextAppendExpression",
},
},
&labeledExpr{
	pos: position{line: 474, col: 55, offset: 14018},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 474, col: 60, offset: 14023},
	expr: &seqExpr{
	pos: position{line: 474, col: 61, offset: 14024},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 474, col: 61, offset: 14024},
	name: "_",
},
&litMatcher{
	pos: position{line: 474, col: 63, offset: 14026},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 474, col: 67, offset: 14030},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 474, col: 70, offset: 14033},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 474, col: 72, offset: 14035},
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
	pos: position{line: 476, col: 1, offset: 14109},
	expr: &actionExpr{
	pos: position{line: 476, col: 26, offset: 14136},
	run: (*parser).callonTextAppendExpression1,
	expr: &seqExpr{
	pos: position{line: 476, col: 26, offset: 14136},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 476, col: 26, offset: 14136},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 476, col: 32, offset: 14142},
	name: "ListAppendExpression",
},
},
&labeledExpr{
	pos: position{line: 476, col: 55, offset: 14165},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 476, col: 60, offset: 14170},
	expr: &seqExpr{
	pos: position{line: 476, col: 61, offset: 14171},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 476, col: 61, offset: 14171},
	name: "_",
},
&litMatcher{
	pos: position{line: 476, col: 63, offset: 14173},
	val: "++",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 476, col: 68, offset: 14178},
	name: "_",
},
&labeledExpr{
	pos: position{line: 476, col: 70, offset: 14180},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 476, col: 72, offset: 14182},
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
	pos: position{line: 478, col: 1, offset: 14262},
	expr: &actionExpr{
	pos: position{line: 478, col: 26, offset: 14289},
	run: (*parser).callonListAppendExpression1,
	expr: &seqExpr{
	pos: position{line: 478, col: 26, offset: 14289},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 478, col: 26, offset: 14289},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 478, col: 32, offset: 14295},
	name: "AndExpression",
},
},
&labeledExpr{
	pos: position{line: 478, col: 55, offset: 14318},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 478, col: 60, offset: 14323},
	expr: &seqExpr{
	pos: position{line: 478, col: 61, offset: 14324},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 478, col: 61, offset: 14324},
	name: "_",
},
&litMatcher{
	pos: position{line: 478, col: 63, offset: 14326},
	val: "#",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 478, col: 67, offset: 14330},
	name: "_",
},
&labeledExpr{
	pos: position{line: 478, col: 69, offset: 14332},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 478, col: 71, offset: 14334},
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
	pos: position{line: 480, col: 1, offset: 14407},
	expr: &actionExpr{
	pos: position{line: 480, col: 26, offset: 14434},
	run: (*parser).callonAndExpression1,
	expr: &seqExpr{
	pos: position{line: 480, col: 26, offset: 14434},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 480, col: 26, offset: 14434},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 480, col: 32, offset: 14440},
	name: "CombineExpression",
},
},
&labeledExpr{
	pos: position{line: 480, col: 55, offset: 14463},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 480, col: 60, offset: 14468},
	expr: &seqExpr{
	pos: position{line: 480, col: 61, offset: 14469},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 480, col: 61, offset: 14469},
	name: "_",
},
&litMatcher{
	pos: position{line: 480, col: 63, offset: 14471},
	val: "&&",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 480, col: 68, offset: 14476},
	name: "_",
},
&labeledExpr{
	pos: position{line: 480, col: 70, offset: 14478},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 480, col: 72, offset: 14480},
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
	pos: position{line: 482, col: 1, offset: 14550},
	expr: &actionExpr{
	pos: position{line: 482, col: 26, offset: 14577},
	run: (*parser).callonCombineExpression1,
	expr: &seqExpr{
	pos: position{line: 482, col: 26, offset: 14577},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 482, col: 26, offset: 14577},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 482, col: 32, offset: 14583},
	name: "PreferExpression",
},
},
&labeledExpr{
	pos: position{line: 482, col: 55, offset: 14606},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 482, col: 60, offset: 14611},
	expr: &seqExpr{
	pos: position{line: 482, col: 61, offset: 14612},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 482, col: 61, offset: 14612},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 482, col: 63, offset: 14614},
	name: "Combine",
},
&ruleRefExpr{
	pos: position{line: 482, col: 71, offset: 14622},
	name: "_",
},
&labeledExpr{
	pos: position{line: 482, col: 73, offset: 14624},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 482, col: 75, offset: 14626},
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
	pos: position{line: 484, col: 1, offset: 14703},
	expr: &actionExpr{
	pos: position{line: 484, col: 26, offset: 14730},
	run: (*parser).callonPreferExpression1,
	expr: &seqExpr{
	pos: position{line: 484, col: 26, offset: 14730},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 484, col: 26, offset: 14730},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 484, col: 32, offset: 14736},
	name: "CombineTypesExpression",
},
},
&labeledExpr{
	pos: position{line: 484, col: 55, offset: 14759},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 484, col: 60, offset: 14764},
	expr: &seqExpr{
	pos: position{line: 484, col: 61, offset: 14765},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 484, col: 61, offset: 14765},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 484, col: 63, offset: 14767},
	name: "Prefer",
},
&ruleRefExpr{
	pos: position{line: 484, col: 70, offset: 14774},
	name: "_",
},
&labeledExpr{
	pos: position{line: 484, col: 72, offset: 14776},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 484, col: 74, offset: 14778},
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
	pos: position{line: 486, col: 1, offset: 14872},
	expr: &actionExpr{
	pos: position{line: 486, col: 26, offset: 14899},
	run: (*parser).callonCombineTypesExpression1,
	expr: &seqExpr{
	pos: position{line: 486, col: 26, offset: 14899},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 486, col: 26, offset: 14899},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 486, col: 32, offset: 14905},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 486, col: 55, offset: 14928},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 486, col: 60, offset: 14933},
	expr: &seqExpr{
	pos: position{line: 486, col: 61, offset: 14934},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 486, col: 61, offset: 14934},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 486, col: 63, offset: 14936},
	name: "CombineTypes",
},
&ruleRefExpr{
	pos: position{line: 486, col: 76, offset: 14949},
	name: "_",
},
&labeledExpr{
	pos: position{line: 486, col: 78, offset: 14951},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 486, col: 80, offset: 14953},
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
	pos: position{line: 488, col: 1, offset: 15033},
	expr: &actionExpr{
	pos: position{line: 488, col: 26, offset: 15060},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 488, col: 26, offset: 15060},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 488, col: 26, offset: 15060},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 488, col: 32, offset: 15066},
	name: "EqualExpression",
},
},
&labeledExpr{
	pos: position{line: 488, col: 55, offset: 15089},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 488, col: 60, offset: 15094},
	expr: &seqExpr{
	pos: position{line: 488, col: 61, offset: 15095},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 488, col: 61, offset: 15095},
	name: "_",
},
&litMatcher{
	pos: position{line: 488, col: 63, offset: 15097},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 488, col: 67, offset: 15101},
	name: "_",
},
&labeledExpr{
	pos: position{line: 488, col: 69, offset: 15103},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 488, col: 71, offset: 15105},
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
	pos: position{line: 490, col: 1, offset: 15175},
	expr: &actionExpr{
	pos: position{line: 490, col: 26, offset: 15202},
	run: (*parser).callonEqualExpression1,
	expr: &seqExpr{
	pos: position{line: 490, col: 26, offset: 15202},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 490, col: 26, offset: 15202},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 490, col: 32, offset: 15208},
	name: "NotEqualExpression",
},
},
&labeledExpr{
	pos: position{line: 490, col: 55, offset: 15231},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 490, col: 60, offset: 15236},
	expr: &seqExpr{
	pos: position{line: 490, col: 61, offset: 15237},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 490, col: 61, offset: 15237},
	name: "_",
},
&litMatcher{
	pos: position{line: 490, col: 63, offset: 15239},
	val: "==",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 490, col: 68, offset: 15244},
	name: "_",
},
&labeledExpr{
	pos: position{line: 490, col: 70, offset: 15246},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 490, col: 72, offset: 15248},
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
	pos: position{line: 492, col: 1, offset: 15318},
	expr: &actionExpr{
	pos: position{line: 492, col: 26, offset: 15345},
	run: (*parser).callonNotEqualExpression1,
	expr: &seqExpr{
	pos: position{line: 492, col: 26, offset: 15345},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 492, col: 26, offset: 15345},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 492, col: 32, offset: 15351},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 492, col: 55, offset: 15374},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 492, col: 60, offset: 15379},
	expr: &seqExpr{
	pos: position{line: 492, col: 61, offset: 15380},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 492, col: 61, offset: 15380},
	name: "_",
},
&litMatcher{
	pos: position{line: 492, col: 63, offset: 15382},
	val: "!=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 492, col: 68, offset: 15387},
	name: "_",
},
&labeledExpr{
	pos: position{line: 492, col: 70, offset: 15389},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 492, col: 72, offset: 15391},
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
	pos: position{line: 495, col: 1, offset: 15465},
	expr: &actionExpr{
	pos: position{line: 495, col: 25, offset: 15491},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 495, col: 25, offset: 15491},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 495, col: 25, offset: 15491},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 495, col: 27, offset: 15493},
	name: "FirstApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 495, col: 54, offset: 15520},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 495, col: 59, offset: 15525},
	expr: &seqExpr{
	pos: position{line: 495, col: 60, offset: 15526},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 495, col: 60, offset: 15526},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 495, col: 63, offset: 15529},
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
	pos: position{line: 504, col: 1, offset: 15772},
	expr: &choiceExpr{
	pos: position{line: 505, col: 8, offset: 15810},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 505, col: 8, offset: 15810},
	run: (*parser).callonFirstApplicationExpression2,
	expr: &seqExpr{
	pos: position{line: 505, col: 8, offset: 15810},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 505, col: 8, offset: 15810},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 505, col: 14, offset: 15816},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 505, col: 17, offset: 15819},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 505, col: 19, offset: 15821},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 505, col: 36, offset: 15838},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 505, col: 39, offset: 15841},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 505, col: 41, offset: 15843},
	name: "ImportExpression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 508, col: 8, offset: 15946},
	run: (*parser).callonFirstApplicationExpression11,
	expr: &seqExpr{
	pos: position{line: 508, col: 8, offset: 15946},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 508, col: 8, offset: 15946},
	name: "Some",
},
&ruleRefExpr{
	pos: position{line: 508, col: 13, offset: 15951},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 508, col: 16, offset: 15954},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 508, col: 18, offset: 15956},
	name: "ImportExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 509, col: 8, offset: 16011},
	name: "ImportExpression",
},
	},
},
},
{
	name: "ImportExpression",
	pos: position{line: 511, col: 1, offset: 16029},
	expr: &choiceExpr{
	pos: position{line: 511, col: 20, offset: 16050},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 511, col: 20, offset: 16050},
	name: "Import",
},
&ruleRefExpr{
	pos: position{line: 511, col: 29, offset: 16059},
	name: "SelectorExpression",
},
	},
},
},
{
	name: "SelectorExpression",
	pos: position{line: 513, col: 1, offset: 16079},
	expr: &actionExpr{
	pos: position{line: 513, col: 22, offset: 16102},
	run: (*parser).callonSelectorExpression1,
	expr: &seqExpr{
	pos: position{line: 513, col: 22, offset: 16102},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 513, col: 22, offset: 16102},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 513, col: 24, offset: 16104},
	name: "PrimitiveExpression",
},
},
&labeledExpr{
	pos: position{line: 513, col: 44, offset: 16124},
	label: "ls",
	expr: &zeroOrMoreExpr{
	pos: position{line: 513, col: 47, offset: 16127},
	expr: &seqExpr{
	pos: position{line: 513, col: 48, offset: 16128},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 513, col: 48, offset: 16128},
	name: "_",
},
&litMatcher{
	pos: position{line: 513, col: 50, offset: 16130},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 513, col: 54, offset: 16134},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 513, col: 56, offset: 16136},
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
	pos: position{line: 530, col: 1, offset: 16615},
	expr: &choiceExpr{
	pos: position{line: 530, col: 12, offset: 16628},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 530, col: 12, offset: 16628},
	name: "AnyLabel",
},
&ruleRefExpr{
	pos: position{line: 530, col: 23, offset: 16639},
	name: "Labels",
},
	},
},
},
{
	name: "Labels",
	pos: position{line: 532, col: 1, offset: 16647},
	expr: &actionExpr{
	pos: position{line: 532, col: 10, offset: 16658},
	run: (*parser).callonLabels1,
	expr: &seqExpr{
	pos: position{line: 532, col: 10, offset: 16658},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 532, col: 10, offset: 16658},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 532, col: 14, offset: 16662},
	name: "_",
},
&labeledExpr{
	pos: position{line: 532, col: 16, offset: 16664},
	label: "optclauses",
	expr: &zeroOrOneExpr{
	pos: position{line: 532, col: 27, offset: 16675},
	expr: &seqExpr{
	pos: position{line: 532, col: 29, offset: 16677},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 532, col: 29, offset: 16677},
	name: "AnyLabel",
},
&ruleRefExpr{
	pos: position{line: 532, col: 38, offset: 16686},
	name: "_",
},
&zeroOrMoreExpr{
	pos: position{line: 532, col: 40, offset: 16688},
	expr: &seqExpr{
	pos: position{line: 532, col: 41, offset: 16689},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 532, col: 41, offset: 16689},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 532, col: 45, offset: 16693},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 532, col: 47, offset: 16695},
	name: "AnyLabel",
},
&ruleRefExpr{
	pos: position{line: 532, col: 56, offset: 16704},
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
	pos: position{line: 532, col: 64, offset: 16712},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PrimitiveExpression",
	pos: position{line: 542, col: 1, offset: 17008},
	expr: &choiceExpr{
	pos: position{line: 543, col: 7, offset: 17038},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 543, col: 7, offset: 17038},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 544, col: 7, offset: 17058},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 545, col: 7, offset: 17079},
	name: "IntegerLiteral",
},
&ruleRefExpr{
	pos: position{line: 546, col: 7, offset: 17100},
	name: "TextLiteral",
},
&actionExpr{
	pos: position{line: 547, col: 7, offset: 17118},
	run: (*parser).callonPrimitiveExpression6,
	expr: &seqExpr{
	pos: position{line: 547, col: 7, offset: 17118},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 547, col: 7, offset: 17118},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 547, col: 11, offset: 17122},
	name: "_",
},
&labeledExpr{
	pos: position{line: 547, col: 13, offset: 17124},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 547, col: 15, offset: 17126},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 547, col: 35, offset: 17146},
	name: "_",
},
&litMatcher{
	pos: position{line: 547, col: 37, offset: 17148},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&actionExpr{
	pos: position{line: 548, col: 7, offset: 17176},
	run: (*parser).callonPrimitiveExpression14,
	expr: &seqExpr{
	pos: position{line: 548, col: 7, offset: 17176},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 548, col: 7, offset: 17176},
	val: "<",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 548, col: 11, offset: 17180},
	name: "_",
},
&labeledExpr{
	pos: position{line: 548, col: 13, offset: 17182},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 548, col: 15, offset: 17184},
	name: "UnionType",
},
},
&ruleRefExpr{
	pos: position{line: 548, col: 25, offset: 17194},
	name: "_",
},
&litMatcher{
	pos: position{line: 548, col: 27, offset: 17196},
	val: ">",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 549, col: 7, offset: 17224},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 550, col: 7, offset: 17250},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 551, col: 7, offset: 17267},
	run: (*parser).callonPrimitiveExpression24,
	expr: &seqExpr{
	pos: position{line: 551, col: 7, offset: 17267},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 551, col: 7, offset: 17267},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 551, col: 11, offset: 17271},
	name: "_",
},
&labeledExpr{
	pos: position{line: 551, col: 14, offset: 17274},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 551, col: 16, offset: 17276},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 551, col: 27, offset: 17287},
	name: "_",
},
&litMatcher{
	pos: position{line: 551, col: 29, offset: 17289},
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
	pos: position{line: 553, col: 1, offset: 17312},
	expr: &choiceExpr{
	pos: position{line: 554, col: 7, offset: 17342},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 554, col: 7, offset: 17342},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &litMatcher{
	pos: position{line: 554, col: 7, offset: 17342},
	val: "=",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 555, col: 7, offset: 17397},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 556, col: 7, offset: 17422},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 557, col: 7, offset: 17450},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 557, col: 7, offset: 17450},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 559, col: 1, offset: 17496},
	expr: &actionExpr{
	pos: position{line: 559, col: 19, offset: 17516},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 559, col: 19, offset: 17516},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 559, col: 19, offset: 17516},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 559, col: 24, offset: 17521},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 559, col: 33, offset: 17530},
	name: "_",
},
&litMatcher{
	pos: position{line: 559, col: 35, offset: 17532},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 559, col: 39, offset: 17536},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 559, col: 42, offset: 17539},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 559, col: 47, offset: 17544},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 562, col: 1, offset: 17601},
	expr: &actionExpr{
	pos: position{line: 562, col: 18, offset: 17620},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 562, col: 18, offset: 17620},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 562, col: 18, offset: 17620},
	name: "_",
},
&litMatcher{
	pos: position{line: 562, col: 20, offset: 17622},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 562, col: 24, offset: 17626},
	name: "_",
},
&labeledExpr{
	pos: position{line: 562, col: 26, offset: 17628},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 562, col: 28, offset: 17630},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 563, col: 1, offset: 17662},
	expr: &actionExpr{
	pos: position{line: 564, col: 7, offset: 17691},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 564, col: 7, offset: 17691},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 564, col: 7, offset: 17691},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 564, col: 13, offset: 17697},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 564, col: 29, offset: 17713},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 564, col: 34, offset: 17718},
	expr: &ruleRefExpr{
	pos: position{line: 564, col: 34, offset: 17718},
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
	pos: position{line: 578, col: 1, offset: 18302},
	expr: &actionExpr{
	pos: position{line: 578, col: 22, offset: 18325},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 578, col: 22, offset: 18325},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 578, col: 22, offset: 18325},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 578, col: 27, offset: 18330},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 578, col: 36, offset: 18339},
	name: "_",
},
&litMatcher{
	pos: position{line: 578, col: 38, offset: 18341},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 578, col: 42, offset: 18345},
	name: "_",
},
&labeledExpr{
	pos: position{line: 578, col: 44, offset: 18347},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 578, col: 49, offset: 18352},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 581, col: 1, offset: 18409},
	expr: &actionExpr{
	pos: position{line: 581, col: 21, offset: 18431},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 581, col: 21, offset: 18431},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 581, col: 21, offset: 18431},
	name: "_",
},
&litMatcher{
	pos: position{line: 581, col: 23, offset: 18433},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 581, col: 27, offset: 18437},
	name: "_",
},
&labeledExpr{
	pos: position{line: 581, col: 29, offset: 18439},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 581, col: 31, offset: 18441},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 582, col: 1, offset: 18476},
	expr: &actionExpr{
	pos: position{line: 583, col: 7, offset: 18508},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 583, col: 7, offset: 18508},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 583, col: 7, offset: 18508},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 583, col: 13, offset: 18514},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 583, col: 32, offset: 18533},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 583, col: 37, offset: 18538},
	expr: &ruleRefExpr{
	pos: position{line: 583, col: 37, offset: 18538},
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
	pos: position{line: 597, col: 1, offset: 19128},
	expr: &choiceExpr{
	pos: position{line: 597, col: 13, offset: 19142},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 597, col: 13, offset: 19142},
	name: "NonEmptyUnionType",
},
&ruleRefExpr{
	pos: position{line: 597, col: 33, offset: 19162},
	name: "EmptyUnionType",
},
	},
},
},
{
	name: "EmptyUnionType",
	pos: position{line: 599, col: 1, offset: 19178},
	expr: &actionExpr{
	pos: position{line: 599, col: 18, offset: 19197},
	run: (*parser).callonEmptyUnionType1,
	expr: &litMatcher{
	pos: position{line: 599, col: 18, offset: 19197},
	val: "",
	ignoreCase: false,
},
},
},
{
	name: "NonEmptyUnionType",
	pos: position{line: 601, col: 1, offset: 19229},
	expr: &actionExpr{
	pos: position{line: 601, col: 21, offset: 19251},
	run: (*parser).callonNonEmptyUnionType1,
	expr: &seqExpr{
	pos: position{line: 601, col: 21, offset: 19251},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 601, col: 21, offset: 19251},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 601, col: 27, offset: 19257},
	name: "UnionVariant",
},
},
&labeledExpr{
	pos: position{line: 601, col: 40, offset: 19270},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 601, col: 45, offset: 19275},
	expr: &seqExpr{
	pos: position{line: 601, col: 46, offset: 19276},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 601, col: 46, offset: 19276},
	name: "_",
},
&litMatcher{
	pos: position{line: 601, col: 48, offset: 19278},
	val: "|",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 601, col: 52, offset: 19282},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 601, col: 54, offset: 19284},
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
	pos: position{line: 621, col: 1, offset: 20006},
	expr: &seqExpr{
	pos: position{line: 621, col: 16, offset: 20023},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 621, col: 16, offset: 20023},
	name: "AnyLabel",
},
&zeroOrOneExpr{
	pos: position{line: 621, col: 25, offset: 20032},
	expr: &seqExpr{
	pos: position{line: 621, col: 26, offset: 20033},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 621, col: 26, offset: 20033},
	name: "_",
},
&litMatcher{
	pos: position{line: 621, col: 28, offset: 20035},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 621, col: 32, offset: 20039},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 621, col: 35, offset: 20042},
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
	pos: position{line: 623, col: 1, offset: 20056},
	expr: &actionExpr{
	pos: position{line: 623, col: 12, offset: 20069},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 623, col: 12, offset: 20069},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 623, col: 12, offset: 20069},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 623, col: 16, offset: 20073},
	name: "_",
},
&labeledExpr{
	pos: position{line: 623, col: 18, offset: 20075},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 623, col: 20, offset: 20077},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 623, col: 31, offset: 20088},
	name: "_",
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 625, col: 1, offset: 20107},
	expr: &actionExpr{
	pos: position{line: 626, col: 7, offset: 20137},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 626, col: 7, offset: 20137},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 626, col: 7, offset: 20137},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 626, col: 11, offset: 20141},
	name: "_",
},
&labeledExpr{
	pos: position{line: 626, col: 13, offset: 20143},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 626, col: 19, offset: 20149},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 626, col: 30, offset: 20160},
	name: "_",
},
&labeledExpr{
	pos: position{line: 626, col: 32, offset: 20162},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 626, col: 37, offset: 20167},
	expr: &ruleRefExpr{
	pos: position{line: 626, col: 37, offset: 20167},
	name: "MoreList",
},
},
},
&litMatcher{
	pos: position{line: 626, col: 47, offset: 20177},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 636, col: 1, offset: 20453},
	expr: &notExpr{
	pos: position{line: 636, col: 7, offset: 20461},
	expr: &anyMatcher{
	line: 636, col: 8, offset: 20462,
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
            i, err := strconv.ParseInt(string(c.text), 16, 32)
            return []byte(string([]rune{rune(i)})), err
        
}

func (p *parser) callonUnicodeEscape2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnicodeEscape2()
}

func (c *current) onUnicodeEscape8() (interface{}, error) {
            i, err := strconv.ParseInt(string(c.text[1:len(c.text)-1]), 16, 32)
            return []byte(string([]rune{rune(i)})), err
        
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

