
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
	name: "QuotedPathCharacter",
	pos: position{line: 298, col: 1, offset: 7849},
	expr: &choiceExpr{
	pos: position{line: 299, col: 6, offset: 7878},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 299, col: 6, offset: 7878},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 300, col: 6, offset: 7895},
	val: "[\\x23-\\x2e]",
	ranges: []rune{'#','.',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 301, col: 6, offset: 7912},
	val: "[\\x30-\\U0010ffff]",
	ranges: []rune{'0','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
	},
},
},
{
	name: "UnquotedPathComponent",
	pos: position{line: 303, col: 1, offset: 7931},
	expr: &actionExpr{
	pos: position{line: 303, col: 25, offset: 7957},
	run: (*parser).callonUnquotedPathComponent1,
	expr: &oneOrMoreExpr{
	pos: position{line: 303, col: 25, offset: 7957},
	expr: &ruleRefExpr{
	pos: position{line: 303, col: 25, offset: 7957},
	name: "PathCharacter",
},
},
},
},
{
	name: "QuotedPathComponent",
	pos: position{line: 304, col: 1, offset: 8003},
	expr: &actionExpr{
	pos: position{line: 304, col: 23, offset: 8027},
	run: (*parser).callonQuotedPathComponent1,
	expr: &oneOrMoreExpr{
	pos: position{line: 304, col: 23, offset: 8027},
	expr: &ruleRefExpr{
	pos: position{line: 304, col: 23, offset: 8027},
	name: "QuotedPathCharacter",
},
},
},
},
{
	name: "PathComponent",
	pos: position{line: 306, col: 1, offset: 8080},
	expr: &choiceExpr{
	pos: position{line: 306, col: 17, offset: 8098},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 306, col: 17, offset: 8098},
	run: (*parser).callonPathComponent2,
	expr: &seqExpr{
	pos: position{line: 306, col: 17, offset: 8098},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 306, col: 17, offset: 8098},
	val: "/",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 306, col: 21, offset: 8102},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 306, col: 23, offset: 8104},
	name: "UnquotedPathComponent",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 307, col: 17, offset: 8160},
	run: (*parser).callonPathComponent7,
	expr: &seqExpr{
	pos: position{line: 307, col: 17, offset: 8160},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 307, col: 17, offset: 8160},
	val: "/",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 307, col: 21, offset: 8164},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 307, col: 25, offset: 8168},
	label: "q",
	expr: &ruleRefExpr{
	pos: position{line: 307, col: 27, offset: 8170},
	name: "QuotedPathComponent",
},
},
&litMatcher{
	pos: position{line: 307, col: 47, offset: 8190},
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
	pos: position{line: 309, col: 1, offset: 8213},
	expr: &actionExpr{
	pos: position{line: 309, col: 8, offset: 8222},
	run: (*parser).callonPath1,
	expr: &labeledExpr{
	pos: position{line: 309, col: 8, offset: 8222},
	label: "cs",
	expr: &oneOrMoreExpr{
	pos: position{line: 309, col: 11, offset: 8225},
	expr: &ruleRefExpr{
	pos: position{line: 309, col: 11, offset: 8225},
	name: "PathComponent",
},
},
},
},
},
{
	name: "Local",
	pos: position{line: 318, col: 1, offset: 8499},
	expr: &choiceExpr{
	pos: position{line: 318, col: 9, offset: 8509},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 318, col: 9, offset: 8509},
	name: "ParentPath",
},
&ruleRefExpr{
	pos: position{line: 318, col: 22, offset: 8522},
	name: "HerePath",
},
&ruleRefExpr{
	pos: position{line: 318, col: 33, offset: 8533},
	name: "HomePath",
},
&ruleRefExpr{
	pos: position{line: 318, col: 44, offset: 8544},
	name: "AbsolutePath",
},
	},
},
},
{
	name: "ParentPath",
	pos: position{line: 320, col: 1, offset: 8558},
	expr: &actionExpr{
	pos: position{line: 320, col: 14, offset: 8573},
	run: (*parser).callonParentPath1,
	expr: &seqExpr{
	pos: position{line: 320, col: 14, offset: 8573},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 320, col: 14, offset: 8573},
	val: "..",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 320, col: 19, offset: 8578},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 320, col: 21, offset: 8580},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HerePath",
	pos: position{line: 321, col: 1, offset: 8636},
	expr: &actionExpr{
	pos: position{line: 321, col: 12, offset: 8649},
	run: (*parser).callonHerePath1,
	expr: &seqExpr{
	pos: position{line: 321, col: 12, offset: 8649},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 321, col: 12, offset: 8649},
	val: ".",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 321, col: 16, offset: 8653},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 321, col: 18, offset: 8655},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HomePath",
	pos: position{line: 322, col: 1, offset: 8694},
	expr: &actionExpr{
	pos: position{line: 322, col: 12, offset: 8707},
	run: (*parser).callonHomePath1,
	expr: &seqExpr{
	pos: position{line: 322, col: 12, offset: 8707},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 322, col: 12, offset: 8707},
	val: "~",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 322, col: 16, offset: 8711},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 322, col: 18, offset: 8713},
	name: "Path",
},
},
	},
},
},
},
{
	name: "AbsolutePath",
	pos: position{line: 323, col: 1, offset: 8768},
	expr: &actionExpr{
	pos: position{line: 323, col: 16, offset: 8785},
	run: (*parser).callonAbsolutePath1,
	expr: &labeledExpr{
	pos: position{line: 323, col: 16, offset: 8785},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 323, col: 18, offset: 8787},
	name: "Path",
},
},
},
},
{
	name: "Scheme",
	pos: position{line: 325, col: 1, offset: 8843},
	expr: &seqExpr{
	pos: position{line: 325, col: 10, offset: 8854},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 325, col: 10, offset: 8854},
	val: "http",
	ignoreCase: false,
},
&zeroOrOneExpr{
	pos: position{line: 325, col: 17, offset: 8861},
	expr: &litMatcher{
	pos: position{line: 325, col: 17, offset: 8861},
	val: "s",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "HttpRaw",
	pos: position{line: 327, col: 1, offset: 8867},
	expr: &actionExpr{
	pos: position{line: 327, col: 11, offset: 8879},
	run: (*parser).callonHttpRaw1,
	expr: &seqExpr{
	pos: position{line: 327, col: 11, offset: 8879},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 327, col: 11, offset: 8879},
	name: "Scheme",
},
&litMatcher{
	pos: position{line: 327, col: 18, offset: 8886},
	val: "://",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 327, col: 24, offset: 8892},
	name: "Authority",
},
&ruleRefExpr{
	pos: position{line: 327, col: 34, offset: 8902},
	name: "Path",
},
&zeroOrOneExpr{
	pos: position{line: 327, col: 39, offset: 8907},
	expr: &seqExpr{
	pos: position{line: 327, col: 41, offset: 8909},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 327, col: 41, offset: 8909},
	val: "?",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 327, col: 45, offset: 8913},
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
	pos: position{line: 329, col: 1, offset: 8970},
	expr: &seqExpr{
	pos: position{line: 329, col: 13, offset: 8984},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 329, col: 13, offset: 8984},
	expr: &seqExpr{
	pos: position{line: 329, col: 14, offset: 8985},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 329, col: 14, offset: 8985},
	name: "Userinfo",
},
&litMatcher{
	pos: position{line: 329, col: 23, offset: 8994},
	val: "@",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 329, col: 29, offset: 9000},
	name: "Host",
},
&zeroOrOneExpr{
	pos: position{line: 329, col: 34, offset: 9005},
	expr: &seqExpr{
	pos: position{line: 329, col: 35, offset: 9006},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 329, col: 35, offset: 9006},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 329, col: 39, offset: 9010},
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
	pos: position{line: 331, col: 1, offset: 9018},
	expr: &zeroOrMoreExpr{
	pos: position{line: 331, col: 12, offset: 9031},
	expr: &choiceExpr{
	pos: position{line: 331, col: 14, offset: 9033},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 331, col: 14, offset: 9033},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 331, col: 27, offset: 9046},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 331, col: 40, offset: 9059},
	name: "SubDelims",
},
&litMatcher{
	pos: position{line: 331, col: 52, offset: 9071},
	val: ":",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "Host",
	pos: position{line: 333, col: 1, offset: 9079},
	expr: &choiceExpr{
	pos: position{line: 333, col: 8, offset: 9088},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 333, col: 8, offset: 9088},
	name: "IPLiteral",
},
&ruleRefExpr{
	pos: position{line: 333, col: 20, offset: 9100},
	name: "RegName",
},
	},
},
},
{
	name: "Port",
	pos: position{line: 335, col: 1, offset: 9109},
	expr: &zeroOrMoreExpr{
	pos: position{line: 335, col: 8, offset: 9118},
	expr: &ruleRefExpr{
	pos: position{line: 335, col: 8, offset: 9118},
	name: "Digit",
},
},
},
{
	name: "IPLiteral",
	pos: position{line: 337, col: 1, offset: 9126},
	expr: &seqExpr{
	pos: position{line: 337, col: 13, offset: 9140},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 337, col: 13, offset: 9140},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 337, col: 17, offset: 9144},
	name: "IPv6address",
},
&litMatcher{
	pos: position{line: 337, col: 29, offset: 9156},
	val: "]",
	ignoreCase: false,
},
	},
},
},
{
	name: "IPv6address",
	pos: position{line: 339, col: 1, offset: 9161},
	expr: &actionExpr{
	pos: position{line: 339, col: 15, offset: 9177},
	run: (*parser).callonIPv6address1,
	expr: &seqExpr{
	pos: position{line: 339, col: 15, offset: 9177},
	exprs: []interface{}{
&zeroOrMoreExpr{
	pos: position{line: 339, col: 15, offset: 9177},
	expr: &ruleRefExpr{
	pos: position{line: 339, col: 16, offset: 9178},
	name: "HexDig",
},
},
&litMatcher{
	pos: position{line: 339, col: 25, offset: 9187},
	val: ":",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 339, col: 29, offset: 9191},
	expr: &choiceExpr{
	pos: position{line: 339, col: 30, offset: 9192},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 339, col: 30, offset: 9192},
	name: "HexDig",
},
&litMatcher{
	pos: position{line: 339, col: 39, offset: 9201},
	val: ":",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 339, col: 45, offset: 9207},
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
	pos: position{line: 345, col: 1, offset: 9361},
	expr: &zeroOrMoreExpr{
	pos: position{line: 345, col: 11, offset: 9373},
	expr: &choiceExpr{
	pos: position{line: 345, col: 12, offset: 9374},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 345, col: 12, offset: 9374},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 345, col: 25, offset: 9387},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 345, col: 38, offset: 9400},
	name: "SubDelims",
},
	},
},
},
},
{
	name: "PChar",
	pos: position{line: 347, col: 1, offset: 9413},
	expr: &choiceExpr{
	pos: position{line: 347, col: 9, offset: 9423},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 347, col: 9, offset: 9423},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 347, col: 22, offset: 9436},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 347, col: 35, offset: 9449},
	name: "SubDelims",
},
&charClassMatcher{
	pos: position{line: 347, col: 47, offset: 9461},
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
	pos: position{line: 349, col: 1, offset: 9467},
	expr: &zeroOrMoreExpr{
	pos: position{line: 349, col: 9, offset: 9477},
	expr: &choiceExpr{
	pos: position{line: 349, col: 10, offset: 9478},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 349, col: 10, offset: 9478},
	name: "PChar",
},
&charClassMatcher{
	pos: position{line: 349, col: 18, offset: 9486},
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
	pos: position{line: 351, col: 1, offset: 9494},
	expr: &seqExpr{
	pos: position{line: 351, col: 14, offset: 9509},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 351, col: 14, offset: 9509},
	val: "%",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 351, col: 18, offset: 9513},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 351, col: 25, offset: 9520},
	name: "HexDig",
},
	},
},
},
{
	name: "Unreserved",
	pos: position{line: 353, col: 1, offset: 9528},
	expr: &charClassMatcher{
	pos: position{line: 353, col: 14, offset: 9543},
	val: "[A-Za-z0-9._~-]",
	chars: []rune{'.','_','~','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SubDelims",
	pos: position{line: 355, col: 1, offset: 9560},
	expr: &choiceExpr{
	pos: position{line: 355, col: 13, offset: 9574},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 355, col: 13, offset: 9574},
	val: "!",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 355, col: 19, offset: 9580},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 355, col: 25, offset: 9586},
	val: "&",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 355, col: 31, offset: 9592},
	val: "'",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 355, col: 37, offset: 9598},
	val: "(",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 355, col: 43, offset: 9604},
	val: ")",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 355, col: 49, offset: 9610},
	val: "*",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 355, col: 55, offset: 9616},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 355, col: 61, offset: 9622},
	val: ",",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 355, col: 67, offset: 9628},
	val: ";",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 355, col: 73, offset: 9634},
	val: "=",
	ignoreCase: false,
},
	},
},
},
{
	name: "Http",
	pos: position{line: 357, col: 1, offset: 9639},
	expr: &actionExpr{
	pos: position{line: 357, col: 8, offset: 9648},
	run: (*parser).callonHttp1,
	expr: &labeledExpr{
	pos: position{line: 357, col: 8, offset: 9648},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 357, col: 10, offset: 9650},
	name: "HttpRaw",
},
},
},
},
{
	name: "Env",
	pos: position{line: 359, col: 1, offset: 9695},
	expr: &actionExpr{
	pos: position{line: 359, col: 7, offset: 9703},
	run: (*parser).callonEnv1,
	expr: &seqExpr{
	pos: position{line: 359, col: 7, offset: 9703},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 359, col: 7, offset: 9703},
	val: "env:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 359, col: 14, offset: 9710},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 359, col: 17, offset: 9713},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 359, col: 17, offset: 9713},
	name: "BashEnvironmentVariable",
},
&ruleRefExpr{
	pos: position{line: 359, col: 43, offset: 9739},
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
	pos: position{line: 361, col: 1, offset: 9784},
	expr: &actionExpr{
	pos: position{line: 361, col: 27, offset: 9812},
	run: (*parser).callonBashEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 361, col: 27, offset: 9812},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 361, col: 27, offset: 9812},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 361, col: 36, offset: 9821},
	expr: &charClassMatcher{
	pos: position{line: 361, col: 36, offset: 9821},
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
	pos: position{line: 365, col: 1, offset: 9877},
	expr: &actionExpr{
	pos: position{line: 365, col: 28, offset: 9906},
	run: (*parser).callonPosixEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 365, col: 28, offset: 9906},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 365, col: 28, offset: 9906},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 365, col: 32, offset: 9910},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 365, col: 34, offset: 9912},
	name: "PosixEnvironmentVariableContent",
},
},
&litMatcher{
	pos: position{line: 365, col: 66, offset: 9944},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PosixEnvironmentVariableContent",
	pos: position{line: 369, col: 1, offset: 9969},
	expr: &actionExpr{
	pos: position{line: 369, col: 35, offset: 10005},
	run: (*parser).callonPosixEnvironmentVariableContent1,
	expr: &labeledExpr{
	pos: position{line: 369, col: 35, offset: 10005},
	label: "v",
	expr: &oneOrMoreExpr{
	pos: position{line: 369, col: 37, offset: 10007},
	expr: &ruleRefExpr{
	pos: position{line: 369, col: 37, offset: 10007},
	name: "PosixEnvironmentVariableCharacter",
},
},
},
},
},
{
	name: "PosixEnvironmentVariableCharacter",
	pos: position{line: 378, col: 1, offset: 10220},
	expr: &choiceExpr{
	pos: position{line: 379, col: 7, offset: 10264},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 379, col: 7, offset: 10264},
	run: (*parser).callonPosixEnvironmentVariableCharacter2,
	expr: &litMatcher{
	pos: position{line: 379, col: 7, offset: 10264},
	val: "\\\"",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 380, col: 7, offset: 10304},
	run: (*parser).callonPosixEnvironmentVariableCharacter4,
	expr: &litMatcher{
	pos: position{line: 380, col: 7, offset: 10304},
	val: "\\\\",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 381, col: 7, offset: 10344},
	run: (*parser).callonPosixEnvironmentVariableCharacter6,
	expr: &litMatcher{
	pos: position{line: 381, col: 7, offset: 10344},
	val: "\\a",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 382, col: 7, offset: 10384},
	run: (*parser).callonPosixEnvironmentVariableCharacter8,
	expr: &litMatcher{
	pos: position{line: 382, col: 7, offset: 10384},
	val: "\\b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 383, col: 7, offset: 10424},
	run: (*parser).callonPosixEnvironmentVariableCharacter10,
	expr: &litMatcher{
	pos: position{line: 383, col: 7, offset: 10424},
	val: "\\f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 384, col: 7, offset: 10464},
	run: (*parser).callonPosixEnvironmentVariableCharacter12,
	expr: &litMatcher{
	pos: position{line: 384, col: 7, offset: 10464},
	val: "\\n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 385, col: 7, offset: 10504},
	run: (*parser).callonPosixEnvironmentVariableCharacter14,
	expr: &litMatcher{
	pos: position{line: 385, col: 7, offset: 10504},
	val: "\\r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 386, col: 7, offset: 10544},
	run: (*parser).callonPosixEnvironmentVariableCharacter16,
	expr: &litMatcher{
	pos: position{line: 386, col: 7, offset: 10544},
	val: "\\t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 387, col: 7, offset: 10584},
	run: (*parser).callonPosixEnvironmentVariableCharacter18,
	expr: &litMatcher{
	pos: position{line: 387, col: 7, offset: 10584},
	val: "\\v",
	ignoreCase: false,
},
},
&charClassMatcher{
	pos: position{line: 388, col: 7, offset: 10624},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 389, col: 7, offset: 10642},
	val: "[\\x23-\\x3c]",
	ranges: []rune{'#','<',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 390, col: 7, offset: 10660},
	val: "[\\x3e-\\x5b]",
	ranges: []rune{'>','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 391, col: 7, offset: 10678},
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
	pos: position{line: 393, col: 1, offset: 10691},
	expr: &choiceExpr{
	pos: position{line: 393, col: 14, offset: 10706},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 393, col: 14, offset: 10706},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 393, col: 24, offset: 10716},
	name: "Local",
},
&ruleRefExpr{
	pos: position{line: 393, col: 32, offset: 10724},
	name: "Http",
},
&ruleRefExpr{
	pos: position{line: 393, col: 39, offset: 10731},
	name: "Env",
},
	},
},
},
{
	name: "HashValue",
	pos: position{line: 396, col: 1, offset: 10804},
	expr: &actionExpr{
	pos: position{line: 396, col: 13, offset: 10816},
	run: (*parser).callonHashValue1,
	expr: &seqExpr{
	pos: position{line: 396, col: 13, offset: 10816},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 396, col: 13, offset: 10816},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 396, col: 20, offset: 10823},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 396, col: 27, offset: 10830},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 396, col: 34, offset: 10837},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 396, col: 41, offset: 10844},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 396, col: 48, offset: 10851},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 396, col: 55, offset: 10858},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 396, col: 62, offset: 10865},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 397, col: 13, offset: 10884},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 397, col: 20, offset: 10891},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 397, col: 27, offset: 10898},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 397, col: 34, offset: 10905},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 397, col: 41, offset: 10912},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 397, col: 48, offset: 10919},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 397, col: 55, offset: 10926},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 397, col: 62, offset: 10933},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 398, col: 13, offset: 10952},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 398, col: 20, offset: 10959},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 398, col: 27, offset: 10966},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 398, col: 34, offset: 10973},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 398, col: 41, offset: 10980},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 398, col: 48, offset: 10987},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 398, col: 55, offset: 10994},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 398, col: 62, offset: 11001},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 399, col: 13, offset: 11020},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 399, col: 20, offset: 11027},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 399, col: 27, offset: 11034},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 399, col: 34, offset: 11041},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 399, col: 41, offset: 11048},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 399, col: 48, offset: 11055},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 399, col: 55, offset: 11062},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 399, col: 62, offset: 11069},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 400, col: 13, offset: 11088},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 400, col: 20, offset: 11095},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 400, col: 27, offset: 11102},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 400, col: 34, offset: 11109},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 400, col: 41, offset: 11116},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 400, col: 48, offset: 11123},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 400, col: 55, offset: 11130},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 400, col: 62, offset: 11137},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 401, col: 13, offset: 11156},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 401, col: 20, offset: 11163},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 401, col: 27, offset: 11170},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 401, col: 34, offset: 11177},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 401, col: 41, offset: 11184},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 401, col: 48, offset: 11191},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 401, col: 55, offset: 11198},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 401, col: 62, offset: 11205},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 402, col: 13, offset: 11224},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 402, col: 20, offset: 11231},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 402, col: 27, offset: 11238},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 402, col: 34, offset: 11245},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 402, col: 41, offset: 11252},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 402, col: 48, offset: 11259},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 402, col: 55, offset: 11266},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 402, col: 62, offset: 11273},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 403, col: 13, offset: 11292},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 403, col: 20, offset: 11299},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 403, col: 27, offset: 11306},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 403, col: 34, offset: 11313},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 403, col: 41, offset: 11320},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 403, col: 48, offset: 11327},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 403, col: 55, offset: 11334},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 403, col: 62, offset: 11341},
	name: "HexDig",
},
	},
},
},
},
{
	name: "Hash",
	pos: position{line: 409, col: 1, offset: 11485},
	expr: &actionExpr{
	pos: position{line: 409, col: 8, offset: 11492},
	run: (*parser).callonHash1,
	expr: &seqExpr{
	pos: position{line: 409, col: 8, offset: 11492},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 409, col: 8, offset: 11492},
	val: "sha256:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 409, col: 18, offset: 11502},
	label: "val",
	expr: &ruleRefExpr{
	pos: position{line: 409, col: 22, offset: 11506},
	name: "HashValue",
},
},
	},
},
},
},
{
	name: "ImportHashed",
	pos: position{line: 411, col: 1, offset: 11576},
	expr: &actionExpr{
	pos: position{line: 411, col: 16, offset: 11593},
	run: (*parser).callonImportHashed1,
	expr: &seqExpr{
	pos: position{line: 411, col: 16, offset: 11593},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 411, col: 16, offset: 11593},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 411, col: 18, offset: 11595},
	name: "ImportType",
},
},
&labeledExpr{
	pos: position{line: 411, col: 29, offset: 11606},
	label: "h",
	expr: &zeroOrOneExpr{
	pos: position{line: 411, col: 31, offset: 11608},
	expr: &seqExpr{
	pos: position{line: 411, col: 32, offset: 11609},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 411, col: 32, offset: 11609},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 411, col: 35, offset: 11612},
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
	pos: position{line: 419, col: 1, offset: 11767},
	expr: &choiceExpr{
	pos: position{line: 419, col: 10, offset: 11778},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 419, col: 10, offset: 11778},
	run: (*parser).callonImport2,
	expr: &seqExpr{
	pos: position{line: 419, col: 10, offset: 11778},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 419, col: 10, offset: 11778},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 419, col: 12, offset: 11780},
	name: "ImportHashed",
},
},
&ruleRefExpr{
	pos: position{line: 419, col: 25, offset: 11793},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 419, col: 27, offset: 11795},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 419, col: 30, offset: 11798},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 419, col: 33, offset: 11801},
	name: "Text",
},
	},
},
},
&actionExpr{
	pos: position{line: 420, col: 10, offset: 11898},
	run: (*parser).callonImport10,
	expr: &labeledExpr{
	pos: position{line: 420, col: 10, offset: 11898},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 420, col: 12, offset: 11900},
	name: "ImportHashed",
},
},
},
	},
},
},
{
	name: "LetBinding",
	pos: position{line: 423, col: 1, offset: 11995},
	expr: &actionExpr{
	pos: position{line: 423, col: 14, offset: 12010},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 423, col: 14, offset: 12010},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 423, col: 14, offset: 12010},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 423, col: 18, offset: 12014},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 423, col: 21, offset: 12017},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 423, col: 27, offset: 12023},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 423, col: 44, offset: 12040},
	name: "_",
},
&labeledExpr{
	pos: position{line: 423, col: 46, offset: 12042},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 423, col: 48, offset: 12044},
	expr: &seqExpr{
	pos: position{line: 423, col: 49, offset: 12045},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 423, col: 49, offset: 12045},
	name: "Annotation",
},
&ruleRefExpr{
	pos: position{line: 423, col: 60, offset: 12056},
	name: "_",
},
	},
},
},
},
&litMatcher{
	pos: position{line: 424, col: 13, offset: 12072},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 424, col: 17, offset: 12076},
	name: "_",
},
&labeledExpr{
	pos: position{line: 424, col: 19, offset: 12078},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 424, col: 21, offset: 12080},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 424, col: 32, offset: 12091},
	name: "_",
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 439, col: 1, offset: 12400},
	expr: &choiceExpr{
	pos: position{line: 440, col: 7, offset: 12421},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 440, col: 7, offset: 12421},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 440, col: 7, offset: 12421},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 440, col: 7, offset: 12421},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 440, col: 14, offset: 12428},
	name: "_",
},
&litMatcher{
	pos: position{line: 440, col: 16, offset: 12430},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 440, col: 20, offset: 12434},
	name: "_",
},
&labeledExpr{
	pos: position{line: 440, col: 22, offset: 12436},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 440, col: 28, offset: 12442},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 440, col: 45, offset: 12459},
	name: "_",
},
&litMatcher{
	pos: position{line: 440, col: 47, offset: 12461},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 440, col: 51, offset: 12465},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 440, col: 54, offset: 12468},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 440, col: 56, offset: 12470},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 440, col: 67, offset: 12481},
	name: "_",
},
&litMatcher{
	pos: position{line: 440, col: 69, offset: 12483},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 440, col: 73, offset: 12487},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 440, col: 75, offset: 12489},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 440, col: 81, offset: 12495},
	name: "_",
},
&labeledExpr{
	pos: position{line: 440, col: 83, offset: 12497},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 440, col: 88, offset: 12502},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 443, col: 7, offset: 12619},
	run: (*parser).callonExpression22,
	expr: &seqExpr{
	pos: position{line: 443, col: 7, offset: 12619},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 443, col: 7, offset: 12619},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 443, col: 10, offset: 12622},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 443, col: 13, offset: 12625},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 443, col: 18, offset: 12630},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 443, col: 29, offset: 12641},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 443, col: 31, offset: 12643},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 443, col: 36, offset: 12648},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 443, col: 39, offset: 12651},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 443, col: 41, offset: 12653},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 443, col: 52, offset: 12664},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 443, col: 54, offset: 12666},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 443, col: 59, offset: 12671},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 443, col: 62, offset: 12674},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 443, col: 64, offset: 12676},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 446, col: 7, offset: 12762},
	run: (*parser).callonExpression38,
	expr: &seqExpr{
	pos: position{line: 446, col: 7, offset: 12762},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 446, col: 7, offset: 12762},
	label: "bindings",
	expr: &oneOrMoreExpr{
	pos: position{line: 446, col: 16, offset: 12771},
	expr: &ruleRefExpr{
	pos: position{line: 446, col: 16, offset: 12771},
	name: "LetBinding",
},
},
},
&ruleRefExpr{
	pos: position{line: 446, col: 28, offset: 12783},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 446, col: 31, offset: 12786},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 446, col: 34, offset: 12789},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 446, col: 36, offset: 12791},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 453, col: 7, offset: 13031},
	run: (*parser).callonExpression47,
	expr: &seqExpr{
	pos: position{line: 453, col: 7, offset: 13031},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 453, col: 7, offset: 13031},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 453, col: 14, offset: 13038},
	name: "_",
},
&litMatcher{
	pos: position{line: 453, col: 16, offset: 13040},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 453, col: 20, offset: 13044},
	name: "_",
},
&labeledExpr{
	pos: position{line: 453, col: 22, offset: 13046},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 453, col: 28, offset: 13052},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 453, col: 45, offset: 13069},
	name: "_",
},
&litMatcher{
	pos: position{line: 453, col: 47, offset: 13071},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 453, col: 51, offset: 13075},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 453, col: 54, offset: 13078},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 453, col: 56, offset: 13080},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 453, col: 67, offset: 13091},
	name: "_",
},
&litMatcher{
	pos: position{line: 453, col: 69, offset: 13093},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 453, col: 73, offset: 13097},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 453, col: 75, offset: 13099},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 453, col: 81, offset: 13105},
	name: "_",
},
&labeledExpr{
	pos: position{line: 453, col: 83, offset: 13107},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 453, col: 88, offset: 13112},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 456, col: 7, offset: 13221},
	run: (*parser).callonExpression67,
	expr: &seqExpr{
	pos: position{line: 456, col: 7, offset: 13221},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 456, col: 7, offset: 13221},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 456, col: 9, offset: 13223},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 456, col: 28, offset: 13242},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 456, col: 30, offset: 13244},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 456, col: 36, offset: 13250},
	name: "_",
},
&labeledExpr{
	pos: position{line: 456, col: 38, offset: 13252},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 456, col: 40, offset: 13254},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 457, col: 7, offset: 13313},
	run: (*parser).callonExpression76,
	expr: &seqExpr{
	pos: position{line: 457, col: 7, offset: 13313},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 457, col: 7, offset: 13313},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 457, col: 13, offset: 13319},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 457, col: 16, offset: 13322},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 457, col: 18, offset: 13324},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 457, col: 35, offset: 13341},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 457, col: 38, offset: 13344},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 457, col: 40, offset: 13346},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 457, col: 57, offset: 13363},
	name: "_",
},
&litMatcher{
	pos: position{line: 457, col: 59, offset: 13365},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 457, col: 63, offset: 13369},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 457, col: 66, offset: 13372},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 457, col: 68, offset: 13374},
	name: "ApplicationExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 460, col: 7, offset: 13495},
	name: "EmptyList",
},
&ruleRefExpr{
	pos: position{line: 461, col: 7, offset: 13511},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 463, col: 1, offset: 13532},
	expr: &actionExpr{
	pos: position{line: 463, col: 14, offset: 13547},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 463, col: 14, offset: 13547},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 463, col: 14, offset: 13547},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 463, col: 18, offset: 13551},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 463, col: 21, offset: 13554},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 463, col: 23, offset: 13556},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 465, col: 1, offset: 13586},
	expr: &actionExpr{
	pos: position{line: 466, col: 1, offset: 13610},
	run: (*parser).callonAnnotatedExpression1,
	expr: &seqExpr{
	pos: position{line: 466, col: 1, offset: 13610},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 466, col: 1, offset: 13610},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 466, col: 3, offset: 13612},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 466, col: 22, offset: 13631},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 466, col: 24, offset: 13633},
	expr: &seqExpr{
	pos: position{line: 466, col: 25, offset: 13634},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 466, col: 25, offset: 13634},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 466, col: 27, offset: 13636},
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
	pos: position{line: 471, col: 1, offset: 13761},
	expr: &actionExpr{
	pos: position{line: 471, col: 13, offset: 13775},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 471, col: 13, offset: 13775},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 471, col: 13, offset: 13775},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 471, col: 17, offset: 13779},
	name: "_",
},
&litMatcher{
	pos: position{line: 471, col: 19, offset: 13781},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 471, col: 23, offset: 13785},
	name: "_",
},
&litMatcher{
	pos: position{line: 471, col: 25, offset: 13787},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 471, col: 29, offset: 13791},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 471, col: 32, offset: 13794},
	name: "List",
},
&ruleRefExpr{
	pos: position{line: 471, col: 37, offset: 13799},
	name: "_",
},
&labeledExpr{
	pos: position{line: 471, col: 39, offset: 13801},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 471, col: 41, offset: 13803},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 475, col: 1, offset: 13866},
	expr: &ruleRefExpr{
	pos: position{line: 475, col: 22, offset: 13889},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 477, col: 1, offset: 13910},
	expr: &actionExpr{
	pos: position{line: 477, col: 26, offset: 13937},
	run: (*parser).callonImportAltExpression1,
	expr: &seqExpr{
	pos: position{line: 477, col: 26, offset: 13937},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 477, col: 26, offset: 13937},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 477, col: 32, offset: 13943},
	name: "OrExpression",
},
},
&labeledExpr{
	pos: position{line: 477, col: 55, offset: 13966},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 477, col: 60, offset: 13971},
	expr: &seqExpr{
	pos: position{line: 477, col: 61, offset: 13972},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 477, col: 61, offset: 13972},
	name: "_",
},
&litMatcher{
	pos: position{line: 477, col: 63, offset: 13974},
	val: "?",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 477, col: 67, offset: 13978},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 477, col: 69, offset: 13980},
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
	pos: position{line: 479, col: 1, offset: 14051},
	expr: &actionExpr{
	pos: position{line: 479, col: 26, offset: 14078},
	run: (*parser).callonOrExpression1,
	expr: &seqExpr{
	pos: position{line: 479, col: 26, offset: 14078},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 479, col: 26, offset: 14078},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 479, col: 32, offset: 14084},
	name: "PlusExpression",
},
},
&labeledExpr{
	pos: position{line: 479, col: 55, offset: 14107},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 479, col: 60, offset: 14112},
	expr: &seqExpr{
	pos: position{line: 479, col: 61, offset: 14113},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 479, col: 61, offset: 14113},
	name: "_",
},
&litMatcher{
	pos: position{line: 479, col: 63, offset: 14115},
	val: "||",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 479, col: 68, offset: 14120},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 479, col: 70, offset: 14122},
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
	pos: position{line: 481, col: 1, offset: 14188},
	expr: &actionExpr{
	pos: position{line: 481, col: 26, offset: 14215},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 481, col: 26, offset: 14215},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 481, col: 26, offset: 14215},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 481, col: 32, offset: 14221},
	name: "TextAppendExpression",
},
},
&labeledExpr{
	pos: position{line: 481, col: 55, offset: 14244},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 481, col: 60, offset: 14249},
	expr: &seqExpr{
	pos: position{line: 481, col: 61, offset: 14250},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 481, col: 61, offset: 14250},
	name: "_",
},
&litMatcher{
	pos: position{line: 481, col: 63, offset: 14252},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 481, col: 67, offset: 14256},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 481, col: 70, offset: 14259},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 481, col: 72, offset: 14261},
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
	pos: position{line: 483, col: 1, offset: 14335},
	expr: &actionExpr{
	pos: position{line: 483, col: 26, offset: 14362},
	run: (*parser).callonTextAppendExpression1,
	expr: &seqExpr{
	pos: position{line: 483, col: 26, offset: 14362},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 483, col: 26, offset: 14362},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 483, col: 32, offset: 14368},
	name: "ListAppendExpression",
},
},
&labeledExpr{
	pos: position{line: 483, col: 55, offset: 14391},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 483, col: 60, offset: 14396},
	expr: &seqExpr{
	pos: position{line: 483, col: 61, offset: 14397},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 483, col: 61, offset: 14397},
	name: "_",
},
&litMatcher{
	pos: position{line: 483, col: 63, offset: 14399},
	val: "++",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 483, col: 68, offset: 14404},
	name: "_",
},
&labeledExpr{
	pos: position{line: 483, col: 70, offset: 14406},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 483, col: 72, offset: 14408},
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
	pos: position{line: 485, col: 1, offset: 14488},
	expr: &actionExpr{
	pos: position{line: 485, col: 26, offset: 14515},
	run: (*parser).callonListAppendExpression1,
	expr: &seqExpr{
	pos: position{line: 485, col: 26, offset: 14515},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 485, col: 26, offset: 14515},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 485, col: 32, offset: 14521},
	name: "AndExpression",
},
},
&labeledExpr{
	pos: position{line: 485, col: 55, offset: 14544},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 485, col: 60, offset: 14549},
	expr: &seqExpr{
	pos: position{line: 485, col: 61, offset: 14550},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 485, col: 61, offset: 14550},
	name: "_",
},
&litMatcher{
	pos: position{line: 485, col: 63, offset: 14552},
	val: "#",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 485, col: 67, offset: 14556},
	name: "_",
},
&labeledExpr{
	pos: position{line: 485, col: 69, offset: 14558},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 485, col: 71, offset: 14560},
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
	pos: position{line: 487, col: 1, offset: 14633},
	expr: &actionExpr{
	pos: position{line: 487, col: 26, offset: 14660},
	run: (*parser).callonAndExpression1,
	expr: &seqExpr{
	pos: position{line: 487, col: 26, offset: 14660},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 487, col: 26, offset: 14660},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 487, col: 32, offset: 14666},
	name: "CombineExpression",
},
},
&labeledExpr{
	pos: position{line: 487, col: 55, offset: 14689},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 487, col: 60, offset: 14694},
	expr: &seqExpr{
	pos: position{line: 487, col: 61, offset: 14695},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 487, col: 61, offset: 14695},
	name: "_",
},
&litMatcher{
	pos: position{line: 487, col: 63, offset: 14697},
	val: "&&",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 487, col: 68, offset: 14702},
	name: "_",
},
&labeledExpr{
	pos: position{line: 487, col: 70, offset: 14704},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 487, col: 72, offset: 14706},
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
	pos: position{line: 489, col: 1, offset: 14776},
	expr: &actionExpr{
	pos: position{line: 489, col: 26, offset: 14803},
	run: (*parser).callonCombineExpression1,
	expr: &seqExpr{
	pos: position{line: 489, col: 26, offset: 14803},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 489, col: 26, offset: 14803},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 489, col: 32, offset: 14809},
	name: "PreferExpression",
},
},
&labeledExpr{
	pos: position{line: 489, col: 55, offset: 14832},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 489, col: 60, offset: 14837},
	expr: &seqExpr{
	pos: position{line: 489, col: 61, offset: 14838},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 489, col: 61, offset: 14838},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 489, col: 63, offset: 14840},
	name: "Combine",
},
&ruleRefExpr{
	pos: position{line: 489, col: 71, offset: 14848},
	name: "_",
},
&labeledExpr{
	pos: position{line: 489, col: 73, offset: 14850},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 489, col: 75, offset: 14852},
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
	pos: position{line: 491, col: 1, offset: 14929},
	expr: &actionExpr{
	pos: position{line: 491, col: 26, offset: 14956},
	run: (*parser).callonPreferExpression1,
	expr: &seqExpr{
	pos: position{line: 491, col: 26, offset: 14956},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 491, col: 26, offset: 14956},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 491, col: 32, offset: 14962},
	name: "CombineTypesExpression",
},
},
&labeledExpr{
	pos: position{line: 491, col: 55, offset: 14985},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 491, col: 60, offset: 14990},
	expr: &seqExpr{
	pos: position{line: 491, col: 61, offset: 14991},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 491, col: 61, offset: 14991},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 491, col: 63, offset: 14993},
	name: "Prefer",
},
&ruleRefExpr{
	pos: position{line: 491, col: 70, offset: 15000},
	name: "_",
},
&labeledExpr{
	pos: position{line: 491, col: 72, offset: 15002},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 491, col: 74, offset: 15004},
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
	pos: position{line: 493, col: 1, offset: 15098},
	expr: &actionExpr{
	pos: position{line: 493, col: 26, offset: 15125},
	run: (*parser).callonCombineTypesExpression1,
	expr: &seqExpr{
	pos: position{line: 493, col: 26, offset: 15125},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 493, col: 26, offset: 15125},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 493, col: 32, offset: 15131},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 493, col: 55, offset: 15154},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 493, col: 60, offset: 15159},
	expr: &seqExpr{
	pos: position{line: 493, col: 61, offset: 15160},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 493, col: 61, offset: 15160},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 493, col: 63, offset: 15162},
	name: "CombineTypes",
},
&ruleRefExpr{
	pos: position{line: 493, col: 76, offset: 15175},
	name: "_",
},
&labeledExpr{
	pos: position{line: 493, col: 78, offset: 15177},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 493, col: 80, offset: 15179},
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
	pos: position{line: 495, col: 1, offset: 15259},
	expr: &actionExpr{
	pos: position{line: 495, col: 26, offset: 15286},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 495, col: 26, offset: 15286},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 495, col: 26, offset: 15286},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 495, col: 32, offset: 15292},
	name: "EqualExpression",
},
},
&labeledExpr{
	pos: position{line: 495, col: 55, offset: 15315},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 495, col: 60, offset: 15320},
	expr: &seqExpr{
	pos: position{line: 495, col: 61, offset: 15321},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 495, col: 61, offset: 15321},
	name: "_",
},
&litMatcher{
	pos: position{line: 495, col: 63, offset: 15323},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 495, col: 67, offset: 15327},
	name: "_",
},
&labeledExpr{
	pos: position{line: 495, col: 69, offset: 15329},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 495, col: 71, offset: 15331},
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
	pos: position{line: 497, col: 1, offset: 15401},
	expr: &actionExpr{
	pos: position{line: 497, col: 26, offset: 15428},
	run: (*parser).callonEqualExpression1,
	expr: &seqExpr{
	pos: position{line: 497, col: 26, offset: 15428},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 497, col: 26, offset: 15428},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 497, col: 32, offset: 15434},
	name: "NotEqualExpression",
},
},
&labeledExpr{
	pos: position{line: 497, col: 55, offset: 15457},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 497, col: 60, offset: 15462},
	expr: &seqExpr{
	pos: position{line: 497, col: 61, offset: 15463},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 497, col: 61, offset: 15463},
	name: "_",
},
&litMatcher{
	pos: position{line: 497, col: 63, offset: 15465},
	val: "==",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 497, col: 68, offset: 15470},
	name: "_",
},
&labeledExpr{
	pos: position{line: 497, col: 70, offset: 15472},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 497, col: 72, offset: 15474},
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
	pos: position{line: 499, col: 1, offset: 15544},
	expr: &actionExpr{
	pos: position{line: 499, col: 26, offset: 15571},
	run: (*parser).callonNotEqualExpression1,
	expr: &seqExpr{
	pos: position{line: 499, col: 26, offset: 15571},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 499, col: 26, offset: 15571},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 499, col: 32, offset: 15577},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 499, col: 55, offset: 15600},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 499, col: 60, offset: 15605},
	expr: &seqExpr{
	pos: position{line: 499, col: 61, offset: 15606},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 499, col: 61, offset: 15606},
	name: "_",
},
&litMatcher{
	pos: position{line: 499, col: 63, offset: 15608},
	val: "!=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 499, col: 68, offset: 15613},
	name: "_",
},
&labeledExpr{
	pos: position{line: 499, col: 70, offset: 15615},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 499, col: 72, offset: 15617},
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
	pos: position{line: 502, col: 1, offset: 15691},
	expr: &actionExpr{
	pos: position{line: 502, col: 25, offset: 15717},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 502, col: 25, offset: 15717},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 502, col: 25, offset: 15717},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 502, col: 27, offset: 15719},
	name: "FirstApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 502, col: 54, offset: 15746},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 502, col: 59, offset: 15751},
	expr: &seqExpr{
	pos: position{line: 502, col: 60, offset: 15752},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 502, col: 60, offset: 15752},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 502, col: 63, offset: 15755},
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
	pos: position{line: 511, col: 1, offset: 15998},
	expr: &choiceExpr{
	pos: position{line: 512, col: 8, offset: 16036},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 512, col: 8, offset: 16036},
	run: (*parser).callonFirstApplicationExpression2,
	expr: &seqExpr{
	pos: position{line: 512, col: 8, offset: 16036},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 512, col: 8, offset: 16036},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 512, col: 14, offset: 16042},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 512, col: 17, offset: 16045},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 512, col: 19, offset: 16047},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 512, col: 36, offset: 16064},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 512, col: 39, offset: 16067},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 512, col: 41, offset: 16069},
	name: "ImportExpression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 515, col: 8, offset: 16172},
	run: (*parser).callonFirstApplicationExpression11,
	expr: &seqExpr{
	pos: position{line: 515, col: 8, offset: 16172},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 515, col: 8, offset: 16172},
	name: "Some",
},
&ruleRefExpr{
	pos: position{line: 515, col: 13, offset: 16177},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 515, col: 16, offset: 16180},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 515, col: 18, offset: 16182},
	name: "ImportExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 516, col: 8, offset: 16237},
	name: "ImportExpression",
},
	},
},
},
{
	name: "ImportExpression",
	pos: position{line: 518, col: 1, offset: 16255},
	expr: &choiceExpr{
	pos: position{line: 518, col: 20, offset: 16276},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 518, col: 20, offset: 16276},
	name: "Import",
},
&ruleRefExpr{
	pos: position{line: 518, col: 29, offset: 16285},
	name: "SelectorExpression",
},
	},
},
},
{
	name: "SelectorExpression",
	pos: position{line: 520, col: 1, offset: 16305},
	expr: &actionExpr{
	pos: position{line: 520, col: 22, offset: 16328},
	run: (*parser).callonSelectorExpression1,
	expr: &seqExpr{
	pos: position{line: 520, col: 22, offset: 16328},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 520, col: 22, offset: 16328},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 520, col: 24, offset: 16330},
	name: "PrimitiveExpression",
},
},
&labeledExpr{
	pos: position{line: 520, col: 44, offset: 16350},
	label: "ls",
	expr: &zeroOrMoreExpr{
	pos: position{line: 520, col: 47, offset: 16353},
	expr: &seqExpr{
	pos: position{line: 520, col: 48, offset: 16354},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 520, col: 48, offset: 16354},
	name: "_",
},
&litMatcher{
	pos: position{line: 520, col: 50, offset: 16356},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 520, col: 54, offset: 16360},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 520, col: 56, offset: 16362},
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
	pos: position{line: 539, col: 1, offset: 16915},
	expr: &choiceExpr{
	pos: position{line: 539, col: 12, offset: 16928},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 539, col: 12, offset: 16928},
	name: "AnyLabel",
},
&ruleRefExpr{
	pos: position{line: 539, col: 23, offset: 16939},
	name: "Labels",
},
&ruleRefExpr{
	pos: position{line: 539, col: 32, offset: 16948},
	name: "TypeSelector",
},
	},
},
},
{
	name: "Labels",
	pos: position{line: 541, col: 1, offset: 16962},
	expr: &actionExpr{
	pos: position{line: 541, col: 10, offset: 16973},
	run: (*parser).callonLabels1,
	expr: &seqExpr{
	pos: position{line: 541, col: 10, offset: 16973},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 541, col: 10, offset: 16973},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 541, col: 14, offset: 16977},
	name: "_",
},
&labeledExpr{
	pos: position{line: 541, col: 16, offset: 16979},
	label: "optclauses",
	expr: &zeroOrOneExpr{
	pos: position{line: 541, col: 27, offset: 16990},
	expr: &seqExpr{
	pos: position{line: 541, col: 29, offset: 16992},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 541, col: 29, offset: 16992},
	name: "AnyLabel",
},
&ruleRefExpr{
	pos: position{line: 541, col: 38, offset: 17001},
	name: "_",
},
&zeroOrMoreExpr{
	pos: position{line: 541, col: 40, offset: 17003},
	expr: &seqExpr{
	pos: position{line: 541, col: 41, offset: 17004},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 541, col: 41, offset: 17004},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 541, col: 45, offset: 17008},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 541, col: 47, offset: 17010},
	name: "AnyLabel",
},
&ruleRefExpr{
	pos: position{line: 541, col: 56, offset: 17019},
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
	pos: position{line: 541, col: 64, offset: 17027},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TypeSelector",
	pos: position{line: 551, col: 1, offset: 17323},
	expr: &actionExpr{
	pos: position{line: 551, col: 16, offset: 17340},
	run: (*parser).callonTypeSelector1,
	expr: &seqExpr{
	pos: position{line: 551, col: 16, offset: 17340},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 551, col: 16, offset: 17340},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 551, col: 20, offset: 17344},
	name: "_",
},
&labeledExpr{
	pos: position{line: 551, col: 22, offset: 17346},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 551, col: 24, offset: 17348},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 551, col: 35, offset: 17359},
	name: "_",
},
&litMatcher{
	pos: position{line: 551, col: 37, offset: 17361},
	val: ")",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PrimitiveExpression",
	pos: position{line: 553, col: 1, offset: 17384},
	expr: &choiceExpr{
	pos: position{line: 554, col: 7, offset: 17414},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 554, col: 7, offset: 17414},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 555, col: 7, offset: 17434},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 556, col: 7, offset: 17455},
	name: "IntegerLiteral",
},
&ruleRefExpr{
	pos: position{line: 557, col: 7, offset: 17476},
	name: "TextLiteral",
},
&actionExpr{
	pos: position{line: 558, col: 7, offset: 17494},
	run: (*parser).callonPrimitiveExpression6,
	expr: &seqExpr{
	pos: position{line: 558, col: 7, offset: 17494},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 558, col: 7, offset: 17494},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 558, col: 11, offset: 17498},
	name: "_",
},
&labeledExpr{
	pos: position{line: 558, col: 13, offset: 17500},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 558, col: 15, offset: 17502},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 558, col: 35, offset: 17522},
	name: "_",
},
&litMatcher{
	pos: position{line: 558, col: 37, offset: 17524},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&actionExpr{
	pos: position{line: 559, col: 7, offset: 17552},
	run: (*parser).callonPrimitiveExpression14,
	expr: &seqExpr{
	pos: position{line: 559, col: 7, offset: 17552},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 559, col: 7, offset: 17552},
	val: "<",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 559, col: 11, offset: 17556},
	name: "_",
},
&labeledExpr{
	pos: position{line: 559, col: 13, offset: 17558},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 559, col: 15, offset: 17560},
	name: "UnionType",
},
},
&ruleRefExpr{
	pos: position{line: 559, col: 25, offset: 17570},
	name: "_",
},
&litMatcher{
	pos: position{line: 559, col: 27, offset: 17572},
	val: ">",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 560, col: 7, offset: 17600},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 561, col: 7, offset: 17626},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 562, col: 7, offset: 17643},
	run: (*parser).callonPrimitiveExpression24,
	expr: &seqExpr{
	pos: position{line: 562, col: 7, offset: 17643},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 562, col: 7, offset: 17643},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 562, col: 11, offset: 17647},
	name: "_",
},
&labeledExpr{
	pos: position{line: 562, col: 14, offset: 17650},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 562, col: 16, offset: 17652},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 562, col: 27, offset: 17663},
	name: "_",
},
&litMatcher{
	pos: position{line: 562, col: 29, offset: 17665},
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
	pos: position{line: 564, col: 1, offset: 17688},
	expr: &choiceExpr{
	pos: position{line: 565, col: 7, offset: 17718},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 565, col: 7, offset: 17718},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &litMatcher{
	pos: position{line: 565, col: 7, offset: 17718},
	val: "=",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 566, col: 7, offset: 17773},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 567, col: 7, offset: 17798},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 568, col: 7, offset: 17826},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 568, col: 7, offset: 17826},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 570, col: 1, offset: 17872},
	expr: &actionExpr{
	pos: position{line: 570, col: 19, offset: 17892},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 570, col: 19, offset: 17892},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 570, col: 19, offset: 17892},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 570, col: 24, offset: 17897},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 570, col: 33, offset: 17906},
	name: "_",
},
&litMatcher{
	pos: position{line: 570, col: 35, offset: 17908},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 570, col: 39, offset: 17912},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 570, col: 42, offset: 17915},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 570, col: 47, offset: 17920},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 573, col: 1, offset: 17977},
	expr: &actionExpr{
	pos: position{line: 573, col: 18, offset: 17996},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 573, col: 18, offset: 17996},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 573, col: 18, offset: 17996},
	name: "_",
},
&litMatcher{
	pos: position{line: 573, col: 20, offset: 17998},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 573, col: 24, offset: 18002},
	name: "_",
},
&labeledExpr{
	pos: position{line: 573, col: 26, offset: 18004},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 573, col: 28, offset: 18006},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 574, col: 1, offset: 18038},
	expr: &actionExpr{
	pos: position{line: 575, col: 7, offset: 18067},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 575, col: 7, offset: 18067},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 575, col: 7, offset: 18067},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 575, col: 13, offset: 18073},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 575, col: 29, offset: 18089},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 575, col: 34, offset: 18094},
	expr: &ruleRefExpr{
	pos: position{line: 575, col: 34, offset: 18094},
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
	pos: position{line: 589, col: 1, offset: 18678},
	expr: &actionExpr{
	pos: position{line: 589, col: 22, offset: 18701},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 589, col: 22, offset: 18701},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 589, col: 22, offset: 18701},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 589, col: 27, offset: 18706},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 589, col: 36, offset: 18715},
	name: "_",
},
&litMatcher{
	pos: position{line: 589, col: 38, offset: 18717},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 589, col: 42, offset: 18721},
	name: "_",
},
&labeledExpr{
	pos: position{line: 589, col: 44, offset: 18723},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 589, col: 49, offset: 18728},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 592, col: 1, offset: 18785},
	expr: &actionExpr{
	pos: position{line: 592, col: 21, offset: 18807},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 592, col: 21, offset: 18807},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 592, col: 21, offset: 18807},
	name: "_",
},
&litMatcher{
	pos: position{line: 592, col: 23, offset: 18809},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 592, col: 27, offset: 18813},
	name: "_",
},
&labeledExpr{
	pos: position{line: 592, col: 29, offset: 18815},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 592, col: 31, offset: 18817},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 593, col: 1, offset: 18852},
	expr: &actionExpr{
	pos: position{line: 594, col: 7, offset: 18884},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 594, col: 7, offset: 18884},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 594, col: 7, offset: 18884},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 594, col: 13, offset: 18890},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 594, col: 32, offset: 18909},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 594, col: 37, offset: 18914},
	expr: &ruleRefExpr{
	pos: position{line: 594, col: 37, offset: 18914},
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
	pos: position{line: 608, col: 1, offset: 19504},
	expr: &choiceExpr{
	pos: position{line: 608, col: 13, offset: 19518},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 608, col: 13, offset: 19518},
	name: "NonEmptyUnionType",
},
&ruleRefExpr{
	pos: position{line: 608, col: 33, offset: 19538},
	name: "EmptyUnionType",
},
	},
},
},
{
	name: "EmptyUnionType",
	pos: position{line: 610, col: 1, offset: 19554},
	expr: &actionExpr{
	pos: position{line: 610, col: 18, offset: 19573},
	run: (*parser).callonEmptyUnionType1,
	expr: &litMatcher{
	pos: position{line: 610, col: 18, offset: 19573},
	val: "",
	ignoreCase: false,
},
},
},
{
	name: "NonEmptyUnionType",
	pos: position{line: 612, col: 1, offset: 19605},
	expr: &actionExpr{
	pos: position{line: 612, col: 21, offset: 19627},
	run: (*parser).callonNonEmptyUnionType1,
	expr: &seqExpr{
	pos: position{line: 612, col: 21, offset: 19627},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 612, col: 21, offset: 19627},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 612, col: 27, offset: 19633},
	name: "UnionVariant",
},
},
&labeledExpr{
	pos: position{line: 612, col: 40, offset: 19646},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 612, col: 45, offset: 19651},
	expr: &seqExpr{
	pos: position{line: 612, col: 46, offset: 19652},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 612, col: 46, offset: 19652},
	name: "_",
},
&litMatcher{
	pos: position{line: 612, col: 48, offset: 19654},
	val: "|",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 612, col: 52, offset: 19658},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 612, col: 54, offset: 19660},
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
	pos: position{line: 632, col: 1, offset: 20382},
	expr: &seqExpr{
	pos: position{line: 632, col: 16, offset: 20399},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 632, col: 16, offset: 20399},
	name: "AnyLabel",
},
&zeroOrOneExpr{
	pos: position{line: 632, col: 25, offset: 20408},
	expr: &seqExpr{
	pos: position{line: 632, col: 26, offset: 20409},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 632, col: 26, offset: 20409},
	name: "_",
},
&litMatcher{
	pos: position{line: 632, col: 28, offset: 20411},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 632, col: 32, offset: 20415},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 632, col: 35, offset: 20418},
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
	pos: position{line: 634, col: 1, offset: 20432},
	expr: &actionExpr{
	pos: position{line: 634, col: 12, offset: 20445},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 634, col: 12, offset: 20445},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 634, col: 12, offset: 20445},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 634, col: 16, offset: 20449},
	name: "_",
},
&labeledExpr{
	pos: position{line: 634, col: 18, offset: 20451},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 634, col: 20, offset: 20453},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 634, col: 31, offset: 20464},
	name: "_",
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 636, col: 1, offset: 20483},
	expr: &actionExpr{
	pos: position{line: 637, col: 7, offset: 20513},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 637, col: 7, offset: 20513},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 637, col: 7, offset: 20513},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 637, col: 11, offset: 20517},
	name: "_",
},
&labeledExpr{
	pos: position{line: 637, col: 13, offset: 20519},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 637, col: 19, offset: 20525},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 637, col: 30, offset: 20536},
	name: "_",
},
&labeledExpr{
	pos: position{line: 637, col: 32, offset: 20538},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 637, col: 37, offset: 20543},
	expr: &ruleRefExpr{
	pos: position{line: 637, col: 37, offset: 20543},
	name: "MoreList",
},
},
},
&litMatcher{
	pos: position{line: 637, col: 47, offset: 20553},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 647, col: 1, offset: 20829},
	expr: &notExpr{
	pos: position{line: 647, col: 7, offset: 20837},
	expr: &anyMatcher{
	line: 647, col: 8, offset: 20838,
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

