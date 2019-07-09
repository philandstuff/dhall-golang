
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
	name: "Location",
	pos: position{line: 238, col: 1, offset: 6466},
	expr: &litMatcher{
	pos: position{line: 238, col: 12, offset: 6479},
	val: "Location",
	ignoreCase: false,
},
},
{
	name: "Combine",
	pos: position{line: 240, col: 1, offset: 6491},
	expr: &choiceExpr{
	pos: position{line: 240, col: 11, offset: 6503},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 240, col: 11, offset: 6503},
	val: "/\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 240, col: 19, offset: 6511},
	val: "∧",
	ignoreCase: false,
},
	},
},
},
{
	name: "CombineTypes",
	pos: position{line: 241, col: 1, offset: 6517},
	expr: &choiceExpr{
	pos: position{line: 241, col: 16, offset: 6534},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 241, col: 16, offset: 6534},
	val: "//\\\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 241, col: 27, offset: 6545},
	val: "⩓",
	ignoreCase: false,
},
	},
},
},
{
	name: "Prefer",
	pos: position{line: 242, col: 1, offset: 6551},
	expr: &choiceExpr{
	pos: position{line: 242, col: 10, offset: 6562},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 242, col: 10, offset: 6562},
	val: "//",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 242, col: 17, offset: 6569},
	val: "⫽",
	ignoreCase: false,
},
	},
},
},
{
	name: "Lambda",
	pos: position{line: 243, col: 1, offset: 6575},
	expr: &choiceExpr{
	pos: position{line: 243, col: 10, offset: 6586},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 243, col: 10, offset: 6586},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 243, col: 17, offset: 6593},
	val: "λ",
	ignoreCase: false,
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 244, col: 1, offset: 6598},
	expr: &choiceExpr{
	pos: position{line: 244, col: 10, offset: 6609},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 244, col: 10, offset: 6609},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 244, col: 21, offset: 6620},
	val: "∀",
	ignoreCase: false,
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 245, col: 1, offset: 6626},
	expr: &choiceExpr{
	pos: position{line: 245, col: 9, offset: 6636},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 245, col: 9, offset: 6636},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 245, col: 16, offset: 6643},
	val: "→",
	ignoreCase: false,
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 247, col: 1, offset: 6650},
	expr: &seqExpr{
	pos: position{line: 247, col: 12, offset: 6663},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 247, col: 12, offset: 6663},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 247, col: 17, offset: 6668},
	expr: &charClassMatcher{
	pos: position{line: 247, col: 17, offset: 6668},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 247, col: 23, offset: 6674},
	expr: &ruleRefExpr{
	pos: position{line: 247, col: 23, offset: 6674},
	name: "Digit",
},
},
	},
},
},
{
	name: "NumericDoubleLiteral",
	pos: position{line: 249, col: 1, offset: 6682},
	expr: &actionExpr{
	pos: position{line: 249, col: 24, offset: 6707},
	run: (*parser).callonNumericDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 249, col: 24, offset: 6707},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 249, col: 24, offset: 6707},
	expr: &charClassMatcher{
	pos: position{line: 249, col: 24, offset: 6707},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 249, col: 30, offset: 6713},
	expr: &ruleRefExpr{
	pos: position{line: 249, col: 30, offset: 6713},
	name: "Digit",
},
},
&choiceExpr{
	pos: position{line: 249, col: 39, offset: 6722},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 249, col: 39, offset: 6722},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 249, col: 39, offset: 6722},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 249, col: 43, offset: 6726},
	expr: &ruleRefExpr{
	pos: position{line: 249, col: 43, offset: 6726},
	name: "Digit",
},
},
&zeroOrOneExpr{
	pos: position{line: 249, col: 50, offset: 6733},
	expr: &ruleRefExpr{
	pos: position{line: 249, col: 50, offset: 6733},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 249, col: 62, offset: 6745},
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
	pos: position{line: 257, col: 1, offset: 6901},
	expr: &choiceExpr{
	pos: position{line: 257, col: 17, offset: 6919},
	alternatives: []interface{}{
&labeledExpr{
	pos: position{line: 257, col: 17, offset: 6919},
	label: "d",
	expr: &ruleRefExpr{
	pos: position{line: 257, col: 19, offset: 6921},
	name: "NumericDoubleLiteral",
},
},
&actionExpr{
	pos: position{line: 258, col: 5, offset: 6946},
	run: (*parser).callonDoubleLiteral4,
	expr: &ruleRefExpr{
	pos: position{line: 258, col: 5, offset: 6946},
	name: "Infinity",
},
},
&actionExpr{
	pos: position{line: 259, col: 5, offset: 6998},
	run: (*parser).callonDoubleLiteral6,
	expr: &seqExpr{
	pos: position{line: 259, col: 5, offset: 6998},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 259, col: 5, offset: 6998},
	val: "-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 259, col: 9, offset: 7002},
	name: "Infinity",
},
	},
},
},
&actionExpr{
	pos: position{line: 260, col: 5, offset: 7055},
	run: (*parser).callonDoubleLiteral10,
	expr: &ruleRefExpr{
	pos: position{line: 260, col: 5, offset: 7055},
	name: "NaN",
},
},
	},
},
},
{
	name: "NaturalLiteral",
	pos: position{line: 262, col: 1, offset: 7098},
	expr: &actionExpr{
	pos: position{line: 262, col: 18, offset: 7117},
	run: (*parser).callonNaturalLiteral1,
	expr: &oneOrMoreExpr{
	pos: position{line: 262, col: 18, offset: 7117},
	expr: &ruleRefExpr{
	pos: position{line: 262, col: 18, offset: 7117},
	name: "Digit",
},
},
},
},
{
	name: "IntegerLiteral",
	pos: position{line: 267, col: 1, offset: 7206},
	expr: &actionExpr{
	pos: position{line: 267, col: 18, offset: 7225},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 267, col: 18, offset: 7225},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 267, col: 18, offset: 7225},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 267, col: 22, offset: 7229},
	name: "NaturalLiteral",
},
	},
},
},
},
{
	name: "DeBruijn",
	pos: position{line: 275, col: 1, offset: 7381},
	expr: &actionExpr{
	pos: position{line: 275, col: 12, offset: 7394},
	run: (*parser).callonDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 275, col: 12, offset: 7394},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 275, col: 12, offset: 7394},
	name: "_",
},
&litMatcher{
	pos: position{line: 275, col: 14, offset: 7396},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 275, col: 18, offset: 7400},
	name: "_",
},
&labeledExpr{
	pos: position{line: 275, col: 20, offset: 7402},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 275, col: 26, offset: 7408},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Variable",
	pos: position{line: 277, col: 1, offset: 7464},
	expr: &actionExpr{
	pos: position{line: 277, col: 12, offset: 7477},
	run: (*parser).callonVariable1,
	expr: &seqExpr{
	pos: position{line: 277, col: 12, offset: 7477},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 277, col: 12, offset: 7477},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 277, col: 17, offset: 7482},
	name: "NonreservedLabel",
},
},
&labeledExpr{
	pos: position{line: 277, col: 34, offset: 7499},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 277, col: 40, offset: 7505},
	expr: &ruleRefExpr{
	pos: position{line: 277, col: 40, offset: 7505},
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
	pos: position{line: 285, col: 1, offset: 7668},
	expr: &choiceExpr{
	pos: position{line: 285, col: 14, offset: 7683},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 285, col: 14, offset: 7683},
	name: "Variable",
},
&ruleRefExpr{
	pos: position{line: 285, col: 25, offset: 7694},
	name: "Reserved",
},
	},
},
},
{
	name: "PathCharacter",
	pos: position{line: 287, col: 1, offset: 7704},
	expr: &choiceExpr{
	pos: position{line: 288, col: 6, offset: 7727},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 288, col: 6, offset: 7727},
	val: "!",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 289, col: 6, offset: 7739},
	val: "[\\x24-\\x27]",
	ranges: []rune{'$','\'',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 290, col: 6, offset: 7756},
	val: "[\\x2a-\\x2b]",
	ranges: []rune{'*','+',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 291, col: 6, offset: 7773},
	val: "[\\x2d-\\x2e]",
	ranges: []rune{'-','.',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 292, col: 6, offset: 7790},
	val: "[\\x30-\\x3b]",
	ranges: []rune{'0',';',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 293, col: 6, offset: 7807},
	val: "=",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 294, col: 6, offset: 7819},
	val: "[\\x40-\\x5a]",
	ranges: []rune{'@','Z',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 295, col: 6, offset: 7836},
	val: "[\\x5e-\\x7a]",
	ranges: []rune{'^','z',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 296, col: 6, offset: 7853},
	val: "|",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 297, col: 6, offset: 7865},
	val: "~",
	ignoreCase: false,
},
	},
},
},
{
	name: "QuotedPathCharacter",
	pos: position{line: 299, col: 1, offset: 7873},
	expr: &choiceExpr{
	pos: position{line: 300, col: 6, offset: 7902},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 300, col: 6, offset: 7902},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 301, col: 6, offset: 7919},
	val: "[\\x23-\\x2e]",
	ranges: []rune{'#','.',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 302, col: 6, offset: 7936},
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
	pos: position{line: 304, col: 1, offset: 7955},
	expr: &actionExpr{
	pos: position{line: 304, col: 25, offset: 7981},
	run: (*parser).callonUnquotedPathComponent1,
	expr: &oneOrMoreExpr{
	pos: position{line: 304, col: 25, offset: 7981},
	expr: &ruleRefExpr{
	pos: position{line: 304, col: 25, offset: 7981},
	name: "PathCharacter",
},
},
},
},
{
	name: "QuotedPathComponent",
	pos: position{line: 305, col: 1, offset: 8027},
	expr: &actionExpr{
	pos: position{line: 305, col: 23, offset: 8051},
	run: (*parser).callonQuotedPathComponent1,
	expr: &oneOrMoreExpr{
	pos: position{line: 305, col: 23, offset: 8051},
	expr: &ruleRefExpr{
	pos: position{line: 305, col: 23, offset: 8051},
	name: "QuotedPathCharacter",
},
},
},
},
{
	name: "PathComponent",
	pos: position{line: 307, col: 1, offset: 8104},
	expr: &choiceExpr{
	pos: position{line: 307, col: 17, offset: 8122},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 307, col: 17, offset: 8122},
	run: (*parser).callonPathComponent2,
	expr: &seqExpr{
	pos: position{line: 307, col: 17, offset: 8122},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 307, col: 17, offset: 8122},
	val: "/",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 307, col: 21, offset: 8126},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 307, col: 23, offset: 8128},
	name: "UnquotedPathComponent",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 308, col: 17, offset: 8184},
	run: (*parser).callonPathComponent7,
	expr: &seqExpr{
	pos: position{line: 308, col: 17, offset: 8184},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 308, col: 17, offset: 8184},
	val: "/",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 308, col: 21, offset: 8188},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 308, col: 25, offset: 8192},
	label: "q",
	expr: &ruleRefExpr{
	pos: position{line: 308, col: 27, offset: 8194},
	name: "QuotedPathComponent",
},
},
&litMatcher{
	pos: position{line: 308, col: 47, offset: 8214},
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
	pos: position{line: 310, col: 1, offset: 8237},
	expr: &actionExpr{
	pos: position{line: 310, col: 8, offset: 8246},
	run: (*parser).callonPath1,
	expr: &labeledExpr{
	pos: position{line: 310, col: 8, offset: 8246},
	label: "cs",
	expr: &oneOrMoreExpr{
	pos: position{line: 310, col: 11, offset: 8249},
	expr: &ruleRefExpr{
	pos: position{line: 310, col: 11, offset: 8249},
	name: "PathComponent",
},
},
},
},
},
{
	name: "Local",
	pos: position{line: 319, col: 1, offset: 8523},
	expr: &choiceExpr{
	pos: position{line: 319, col: 9, offset: 8533},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 319, col: 9, offset: 8533},
	name: "ParentPath",
},
&ruleRefExpr{
	pos: position{line: 319, col: 22, offset: 8546},
	name: "HerePath",
},
&ruleRefExpr{
	pos: position{line: 319, col: 33, offset: 8557},
	name: "HomePath",
},
&ruleRefExpr{
	pos: position{line: 319, col: 44, offset: 8568},
	name: "AbsolutePath",
},
	},
},
},
{
	name: "ParentPath",
	pos: position{line: 321, col: 1, offset: 8582},
	expr: &actionExpr{
	pos: position{line: 321, col: 14, offset: 8597},
	run: (*parser).callonParentPath1,
	expr: &seqExpr{
	pos: position{line: 321, col: 14, offset: 8597},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 321, col: 14, offset: 8597},
	val: "..",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 321, col: 19, offset: 8602},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 321, col: 21, offset: 8604},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HerePath",
	pos: position{line: 322, col: 1, offset: 8660},
	expr: &actionExpr{
	pos: position{line: 322, col: 12, offset: 8673},
	run: (*parser).callonHerePath1,
	expr: &seqExpr{
	pos: position{line: 322, col: 12, offset: 8673},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 322, col: 12, offset: 8673},
	val: ".",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 322, col: 16, offset: 8677},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 322, col: 18, offset: 8679},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HomePath",
	pos: position{line: 323, col: 1, offset: 8718},
	expr: &actionExpr{
	pos: position{line: 323, col: 12, offset: 8731},
	run: (*parser).callonHomePath1,
	expr: &seqExpr{
	pos: position{line: 323, col: 12, offset: 8731},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 323, col: 12, offset: 8731},
	val: "~",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 323, col: 16, offset: 8735},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 323, col: 18, offset: 8737},
	name: "Path",
},
},
	},
},
},
},
{
	name: "AbsolutePath",
	pos: position{line: 324, col: 1, offset: 8792},
	expr: &actionExpr{
	pos: position{line: 324, col: 16, offset: 8809},
	run: (*parser).callonAbsolutePath1,
	expr: &labeledExpr{
	pos: position{line: 324, col: 16, offset: 8809},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 324, col: 18, offset: 8811},
	name: "Path",
},
},
},
},
{
	name: "Scheme",
	pos: position{line: 326, col: 1, offset: 8867},
	expr: &seqExpr{
	pos: position{line: 326, col: 10, offset: 8878},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 326, col: 10, offset: 8878},
	val: "http",
	ignoreCase: false,
},
&zeroOrOneExpr{
	pos: position{line: 326, col: 17, offset: 8885},
	expr: &litMatcher{
	pos: position{line: 326, col: 17, offset: 8885},
	val: "s",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "HttpRaw",
	pos: position{line: 328, col: 1, offset: 8891},
	expr: &actionExpr{
	pos: position{line: 328, col: 11, offset: 8903},
	run: (*parser).callonHttpRaw1,
	expr: &seqExpr{
	pos: position{line: 328, col: 11, offset: 8903},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 328, col: 11, offset: 8903},
	name: "Scheme",
},
&litMatcher{
	pos: position{line: 328, col: 18, offset: 8910},
	val: "://",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 328, col: 24, offset: 8916},
	name: "Authority",
},
&ruleRefExpr{
	pos: position{line: 328, col: 34, offset: 8926},
	name: "UrlPath",
},
&zeroOrOneExpr{
	pos: position{line: 328, col: 42, offset: 8934},
	expr: &seqExpr{
	pos: position{line: 328, col: 44, offset: 8936},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 328, col: 44, offset: 8936},
	val: "?",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 328, col: 48, offset: 8940},
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
	pos: position{line: 330, col: 1, offset: 8997},
	expr: &zeroOrMoreExpr{
	pos: position{line: 330, col: 11, offset: 9009},
	expr: &choiceExpr{
	pos: position{line: 330, col: 12, offset: 9010},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 330, col: 12, offset: 9010},
	name: "PathComponent",
},
&seqExpr{
	pos: position{line: 330, col: 28, offset: 9026},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 330, col: 28, offset: 9026},
	val: "/",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 330, col: 32, offset: 9030},
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
	pos: position{line: 332, col: 1, offset: 9041},
	expr: &seqExpr{
	pos: position{line: 332, col: 13, offset: 9055},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 332, col: 13, offset: 9055},
	expr: &seqExpr{
	pos: position{line: 332, col: 14, offset: 9056},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 332, col: 14, offset: 9056},
	name: "Userinfo",
},
&litMatcher{
	pos: position{line: 332, col: 23, offset: 9065},
	val: "@",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 332, col: 29, offset: 9071},
	name: "Host",
},
&zeroOrOneExpr{
	pos: position{line: 332, col: 34, offset: 9076},
	expr: &seqExpr{
	pos: position{line: 332, col: 35, offset: 9077},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 332, col: 35, offset: 9077},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 332, col: 39, offset: 9081},
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
	pos: position{line: 334, col: 1, offset: 9089},
	expr: &zeroOrMoreExpr{
	pos: position{line: 334, col: 12, offset: 9102},
	expr: &choiceExpr{
	pos: position{line: 334, col: 14, offset: 9104},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 334, col: 14, offset: 9104},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 334, col: 27, offset: 9117},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 334, col: 40, offset: 9130},
	name: "SubDelims",
},
&litMatcher{
	pos: position{line: 334, col: 52, offset: 9142},
	val: ":",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "Host",
	pos: position{line: 336, col: 1, offset: 9150},
	expr: &choiceExpr{
	pos: position{line: 336, col: 8, offset: 9159},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 336, col: 8, offset: 9159},
	name: "IPLiteral",
},
&ruleRefExpr{
	pos: position{line: 336, col: 20, offset: 9171},
	name: "RegName",
},
	},
},
},
{
	name: "Port",
	pos: position{line: 338, col: 1, offset: 9180},
	expr: &zeroOrMoreExpr{
	pos: position{line: 338, col: 8, offset: 9189},
	expr: &ruleRefExpr{
	pos: position{line: 338, col: 8, offset: 9189},
	name: "Digit",
},
},
},
{
	name: "IPLiteral",
	pos: position{line: 340, col: 1, offset: 9197},
	expr: &seqExpr{
	pos: position{line: 340, col: 13, offset: 9211},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 340, col: 13, offset: 9211},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 340, col: 17, offset: 9215},
	name: "IPv6address",
},
&litMatcher{
	pos: position{line: 340, col: 29, offset: 9227},
	val: "]",
	ignoreCase: false,
},
	},
},
},
{
	name: "IPv6address",
	pos: position{line: 342, col: 1, offset: 9232},
	expr: &actionExpr{
	pos: position{line: 342, col: 15, offset: 9248},
	run: (*parser).callonIPv6address1,
	expr: &seqExpr{
	pos: position{line: 342, col: 15, offset: 9248},
	exprs: []interface{}{
&zeroOrMoreExpr{
	pos: position{line: 342, col: 15, offset: 9248},
	expr: &ruleRefExpr{
	pos: position{line: 342, col: 16, offset: 9249},
	name: "HexDig",
},
},
&litMatcher{
	pos: position{line: 342, col: 25, offset: 9258},
	val: ":",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 342, col: 29, offset: 9262},
	expr: &choiceExpr{
	pos: position{line: 342, col: 30, offset: 9263},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 342, col: 30, offset: 9263},
	name: "HexDig",
},
&litMatcher{
	pos: position{line: 342, col: 39, offset: 9272},
	val: ":",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 342, col: 45, offset: 9278},
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
	pos: position{line: 348, col: 1, offset: 9432},
	expr: &zeroOrMoreExpr{
	pos: position{line: 348, col: 11, offset: 9444},
	expr: &choiceExpr{
	pos: position{line: 348, col: 12, offset: 9445},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 348, col: 12, offset: 9445},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 348, col: 25, offset: 9458},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 348, col: 38, offset: 9471},
	name: "SubDelims",
},
	},
},
},
},
{
	name: "Segment",
	pos: position{line: 350, col: 1, offset: 9484},
	expr: &zeroOrMoreExpr{
	pos: position{line: 350, col: 11, offset: 9496},
	expr: &ruleRefExpr{
	pos: position{line: 350, col: 11, offset: 9496},
	name: "PChar",
},
},
},
{
	name: "PChar",
	pos: position{line: 352, col: 1, offset: 9504},
	expr: &choiceExpr{
	pos: position{line: 352, col: 9, offset: 9514},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 352, col: 9, offset: 9514},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 352, col: 22, offset: 9527},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 352, col: 35, offset: 9540},
	name: "SubDelims",
},
&charClassMatcher{
	pos: position{line: 352, col: 47, offset: 9552},
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
	pos: position{line: 354, col: 1, offset: 9558},
	expr: &zeroOrMoreExpr{
	pos: position{line: 354, col: 9, offset: 9568},
	expr: &choiceExpr{
	pos: position{line: 354, col: 10, offset: 9569},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 354, col: 10, offset: 9569},
	name: "PChar",
},
&charClassMatcher{
	pos: position{line: 354, col: 18, offset: 9577},
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
	pos: position{line: 356, col: 1, offset: 9585},
	expr: &seqExpr{
	pos: position{line: 356, col: 14, offset: 9600},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 356, col: 14, offset: 9600},
	val: "%",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 356, col: 18, offset: 9604},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 356, col: 25, offset: 9611},
	name: "HexDig",
},
	},
},
},
{
	name: "Unreserved",
	pos: position{line: 358, col: 1, offset: 9619},
	expr: &charClassMatcher{
	pos: position{line: 358, col: 14, offset: 9634},
	val: "[A-Za-z0-9._~-]",
	chars: []rune{'.','_','~','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SubDelims",
	pos: position{line: 360, col: 1, offset: 9651},
	expr: &choiceExpr{
	pos: position{line: 360, col: 13, offset: 9665},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 360, col: 13, offset: 9665},
	val: "!",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 360, col: 19, offset: 9671},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 360, col: 25, offset: 9677},
	val: "&",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 360, col: 31, offset: 9683},
	val: "'",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 360, col: 37, offset: 9689},
	val: "*",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 360, col: 43, offset: 9695},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 360, col: 49, offset: 9701},
	val: ";",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 360, col: 55, offset: 9707},
	val: "=",
	ignoreCase: false,
},
	},
},
},
{
	name: "Http",
	pos: position{line: 362, col: 1, offset: 9712},
	expr: &actionExpr{
	pos: position{line: 362, col: 8, offset: 9721},
	run: (*parser).callonHttp1,
	expr: &labeledExpr{
	pos: position{line: 362, col: 8, offset: 9721},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 362, col: 10, offset: 9723},
	name: "HttpRaw",
},
},
},
},
{
	name: "Env",
	pos: position{line: 364, col: 1, offset: 9768},
	expr: &actionExpr{
	pos: position{line: 364, col: 7, offset: 9776},
	run: (*parser).callonEnv1,
	expr: &seqExpr{
	pos: position{line: 364, col: 7, offset: 9776},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 364, col: 7, offset: 9776},
	val: "env:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 364, col: 14, offset: 9783},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 364, col: 17, offset: 9786},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 364, col: 17, offset: 9786},
	name: "BashEnvironmentVariable",
},
&ruleRefExpr{
	pos: position{line: 364, col: 43, offset: 9812},
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
	pos: position{line: 366, col: 1, offset: 9857},
	expr: &actionExpr{
	pos: position{line: 366, col: 27, offset: 9885},
	run: (*parser).callonBashEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 366, col: 27, offset: 9885},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 366, col: 27, offset: 9885},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 366, col: 36, offset: 9894},
	expr: &charClassMatcher{
	pos: position{line: 366, col: 36, offset: 9894},
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
	pos: position{line: 370, col: 1, offset: 9950},
	expr: &actionExpr{
	pos: position{line: 370, col: 28, offset: 9979},
	run: (*parser).callonPosixEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 370, col: 28, offset: 9979},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 370, col: 28, offset: 9979},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 370, col: 32, offset: 9983},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 370, col: 34, offset: 9985},
	name: "PosixEnvironmentVariableContent",
},
},
&litMatcher{
	pos: position{line: 370, col: 66, offset: 10017},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PosixEnvironmentVariableContent",
	pos: position{line: 374, col: 1, offset: 10042},
	expr: &actionExpr{
	pos: position{line: 374, col: 35, offset: 10078},
	run: (*parser).callonPosixEnvironmentVariableContent1,
	expr: &labeledExpr{
	pos: position{line: 374, col: 35, offset: 10078},
	label: "v",
	expr: &oneOrMoreExpr{
	pos: position{line: 374, col: 37, offset: 10080},
	expr: &ruleRefExpr{
	pos: position{line: 374, col: 37, offset: 10080},
	name: "PosixEnvironmentVariableCharacter",
},
},
},
},
},
{
	name: "PosixEnvironmentVariableCharacter",
	pos: position{line: 383, col: 1, offset: 10293},
	expr: &choiceExpr{
	pos: position{line: 384, col: 7, offset: 10337},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 384, col: 7, offset: 10337},
	run: (*parser).callonPosixEnvironmentVariableCharacter2,
	expr: &litMatcher{
	pos: position{line: 384, col: 7, offset: 10337},
	val: "\\\"",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 385, col: 7, offset: 10377},
	run: (*parser).callonPosixEnvironmentVariableCharacter4,
	expr: &litMatcher{
	pos: position{line: 385, col: 7, offset: 10377},
	val: "\\\\",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 386, col: 7, offset: 10417},
	run: (*parser).callonPosixEnvironmentVariableCharacter6,
	expr: &litMatcher{
	pos: position{line: 386, col: 7, offset: 10417},
	val: "\\a",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 387, col: 7, offset: 10457},
	run: (*parser).callonPosixEnvironmentVariableCharacter8,
	expr: &litMatcher{
	pos: position{line: 387, col: 7, offset: 10457},
	val: "\\b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 388, col: 7, offset: 10497},
	run: (*parser).callonPosixEnvironmentVariableCharacter10,
	expr: &litMatcher{
	pos: position{line: 388, col: 7, offset: 10497},
	val: "\\f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 389, col: 7, offset: 10537},
	run: (*parser).callonPosixEnvironmentVariableCharacter12,
	expr: &litMatcher{
	pos: position{line: 389, col: 7, offset: 10537},
	val: "\\n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 390, col: 7, offset: 10577},
	run: (*parser).callonPosixEnvironmentVariableCharacter14,
	expr: &litMatcher{
	pos: position{line: 390, col: 7, offset: 10577},
	val: "\\r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 391, col: 7, offset: 10617},
	run: (*parser).callonPosixEnvironmentVariableCharacter16,
	expr: &litMatcher{
	pos: position{line: 391, col: 7, offset: 10617},
	val: "\\t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 392, col: 7, offset: 10657},
	run: (*parser).callonPosixEnvironmentVariableCharacter18,
	expr: &litMatcher{
	pos: position{line: 392, col: 7, offset: 10657},
	val: "\\v",
	ignoreCase: false,
},
},
&charClassMatcher{
	pos: position{line: 393, col: 7, offset: 10697},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 394, col: 7, offset: 10715},
	val: "[\\x23-\\x3c]",
	ranges: []rune{'#','<',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 395, col: 7, offset: 10733},
	val: "[\\x3e-\\x5b]",
	ranges: []rune{'>','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 396, col: 7, offset: 10751},
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
	pos: position{line: 398, col: 1, offset: 10764},
	expr: &choiceExpr{
	pos: position{line: 398, col: 14, offset: 10779},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 398, col: 14, offset: 10779},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 398, col: 24, offset: 10789},
	name: "Local",
},
&ruleRefExpr{
	pos: position{line: 398, col: 32, offset: 10797},
	name: "Http",
},
&ruleRefExpr{
	pos: position{line: 398, col: 39, offset: 10804},
	name: "Env",
},
	},
},
},
{
	name: "HashValue",
	pos: position{line: 401, col: 1, offset: 10877},
	expr: &actionExpr{
	pos: position{line: 401, col: 13, offset: 10889},
	run: (*parser).callonHashValue1,
	expr: &seqExpr{
	pos: position{line: 401, col: 13, offset: 10889},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 401, col: 13, offset: 10889},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 401, col: 20, offset: 10896},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 401, col: 27, offset: 10903},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 401, col: 34, offset: 10910},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 401, col: 41, offset: 10917},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 401, col: 48, offset: 10924},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 401, col: 55, offset: 10931},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 401, col: 62, offset: 10938},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 402, col: 13, offset: 10957},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 402, col: 20, offset: 10964},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 402, col: 27, offset: 10971},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 402, col: 34, offset: 10978},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 402, col: 41, offset: 10985},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 402, col: 48, offset: 10992},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 402, col: 55, offset: 10999},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 402, col: 62, offset: 11006},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 403, col: 13, offset: 11025},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 403, col: 20, offset: 11032},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 403, col: 27, offset: 11039},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 403, col: 34, offset: 11046},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 403, col: 41, offset: 11053},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 403, col: 48, offset: 11060},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 403, col: 55, offset: 11067},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 403, col: 62, offset: 11074},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 404, col: 13, offset: 11093},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 404, col: 20, offset: 11100},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 404, col: 27, offset: 11107},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 404, col: 34, offset: 11114},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 404, col: 41, offset: 11121},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 404, col: 48, offset: 11128},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 404, col: 55, offset: 11135},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 404, col: 62, offset: 11142},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 405, col: 13, offset: 11161},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 405, col: 20, offset: 11168},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 405, col: 27, offset: 11175},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 405, col: 34, offset: 11182},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 405, col: 41, offset: 11189},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 405, col: 48, offset: 11196},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 405, col: 55, offset: 11203},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 405, col: 62, offset: 11210},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 406, col: 13, offset: 11229},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 406, col: 20, offset: 11236},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 406, col: 27, offset: 11243},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 406, col: 34, offset: 11250},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 406, col: 41, offset: 11257},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 406, col: 48, offset: 11264},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 406, col: 55, offset: 11271},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 406, col: 62, offset: 11278},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 407, col: 13, offset: 11297},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 407, col: 20, offset: 11304},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 407, col: 27, offset: 11311},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 407, col: 34, offset: 11318},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 407, col: 41, offset: 11325},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 407, col: 48, offset: 11332},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 407, col: 55, offset: 11339},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 407, col: 62, offset: 11346},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 408, col: 13, offset: 11365},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 408, col: 20, offset: 11372},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 408, col: 27, offset: 11379},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 408, col: 34, offset: 11386},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 408, col: 41, offset: 11393},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 408, col: 48, offset: 11400},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 408, col: 55, offset: 11407},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 408, col: 62, offset: 11414},
	name: "HexDig",
},
	},
},
},
},
{
	name: "Hash",
	pos: position{line: 414, col: 1, offset: 11558},
	expr: &actionExpr{
	pos: position{line: 414, col: 8, offset: 11565},
	run: (*parser).callonHash1,
	expr: &seqExpr{
	pos: position{line: 414, col: 8, offset: 11565},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 414, col: 8, offset: 11565},
	val: "sha256:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 414, col: 18, offset: 11575},
	label: "val",
	expr: &ruleRefExpr{
	pos: position{line: 414, col: 22, offset: 11579},
	name: "HashValue",
},
},
	},
},
},
},
{
	name: "ImportHashed",
	pos: position{line: 416, col: 1, offset: 11649},
	expr: &actionExpr{
	pos: position{line: 416, col: 16, offset: 11666},
	run: (*parser).callonImportHashed1,
	expr: &seqExpr{
	pos: position{line: 416, col: 16, offset: 11666},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 416, col: 16, offset: 11666},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 416, col: 18, offset: 11668},
	name: "ImportType",
},
},
&labeledExpr{
	pos: position{line: 416, col: 29, offset: 11679},
	label: "h",
	expr: &zeroOrOneExpr{
	pos: position{line: 416, col: 31, offset: 11681},
	expr: &seqExpr{
	pos: position{line: 416, col: 32, offset: 11682},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 416, col: 32, offset: 11682},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 416, col: 35, offset: 11685},
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
	pos: position{line: 424, col: 1, offset: 11840},
	expr: &choiceExpr{
	pos: position{line: 424, col: 10, offset: 11851},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 424, col: 10, offset: 11851},
	run: (*parser).callonImport2,
	expr: &seqExpr{
	pos: position{line: 424, col: 10, offset: 11851},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 424, col: 10, offset: 11851},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 424, col: 12, offset: 11853},
	name: "ImportHashed",
},
},
&ruleRefExpr{
	pos: position{line: 424, col: 25, offset: 11866},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 424, col: 27, offset: 11868},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 424, col: 30, offset: 11871},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 424, col: 33, offset: 11874},
	name: "Text",
},
	},
},
},
&actionExpr{
	pos: position{line: 425, col: 10, offset: 11971},
	run: (*parser).callonImport10,
	expr: &seqExpr{
	pos: position{line: 425, col: 10, offset: 11971},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 425, col: 10, offset: 11971},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 425, col: 12, offset: 11973},
	name: "ImportHashed",
},
},
&ruleRefExpr{
	pos: position{line: 425, col: 25, offset: 11986},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 425, col: 27, offset: 11988},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 425, col: 30, offset: 11991},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 425, col: 33, offset: 11994},
	name: "Location",
},
	},
},
},
&actionExpr{
	pos: position{line: 426, col: 10, offset: 12096},
	run: (*parser).callonImport18,
	expr: &labeledExpr{
	pos: position{line: 426, col: 10, offset: 12096},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 426, col: 12, offset: 12098},
	name: "ImportHashed",
},
},
},
	},
},
},
{
	name: "LetBinding",
	pos: position{line: 429, col: 1, offset: 12193},
	expr: &actionExpr{
	pos: position{line: 429, col: 14, offset: 12208},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 429, col: 14, offset: 12208},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 429, col: 14, offset: 12208},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 429, col: 18, offset: 12212},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 429, col: 21, offset: 12215},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 429, col: 27, offset: 12221},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 429, col: 44, offset: 12238},
	name: "_",
},
&labeledExpr{
	pos: position{line: 429, col: 46, offset: 12240},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 429, col: 48, offset: 12242},
	expr: &seqExpr{
	pos: position{line: 429, col: 49, offset: 12243},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 429, col: 49, offset: 12243},
	name: "Annotation",
},
&ruleRefExpr{
	pos: position{line: 429, col: 60, offset: 12254},
	name: "_",
},
	},
},
},
},
&litMatcher{
	pos: position{line: 430, col: 13, offset: 12270},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 430, col: 17, offset: 12274},
	name: "_",
},
&labeledExpr{
	pos: position{line: 430, col: 19, offset: 12276},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 430, col: 21, offset: 12278},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 430, col: 32, offset: 12289},
	name: "_",
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 445, col: 1, offset: 12598},
	expr: &choiceExpr{
	pos: position{line: 446, col: 7, offset: 12619},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 446, col: 7, offset: 12619},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 446, col: 7, offset: 12619},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 446, col: 7, offset: 12619},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 446, col: 14, offset: 12626},
	name: "_",
},
&litMatcher{
	pos: position{line: 446, col: 16, offset: 12628},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 446, col: 20, offset: 12632},
	name: "_",
},
&labeledExpr{
	pos: position{line: 446, col: 22, offset: 12634},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 446, col: 28, offset: 12640},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 446, col: 45, offset: 12657},
	name: "_",
},
&litMatcher{
	pos: position{line: 446, col: 47, offset: 12659},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 446, col: 51, offset: 12663},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 446, col: 54, offset: 12666},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 446, col: 56, offset: 12668},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 446, col: 67, offset: 12679},
	name: "_",
},
&litMatcher{
	pos: position{line: 446, col: 69, offset: 12681},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 446, col: 73, offset: 12685},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 446, col: 75, offset: 12687},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 446, col: 81, offset: 12693},
	name: "_",
},
&labeledExpr{
	pos: position{line: 446, col: 83, offset: 12695},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 446, col: 88, offset: 12700},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 449, col: 7, offset: 12817},
	run: (*parser).callonExpression22,
	expr: &seqExpr{
	pos: position{line: 449, col: 7, offset: 12817},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 449, col: 7, offset: 12817},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 449, col: 10, offset: 12820},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 449, col: 13, offset: 12823},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 449, col: 18, offset: 12828},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 449, col: 29, offset: 12839},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 449, col: 31, offset: 12841},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 449, col: 36, offset: 12846},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 449, col: 39, offset: 12849},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 449, col: 41, offset: 12851},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 449, col: 52, offset: 12862},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 449, col: 54, offset: 12864},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 449, col: 59, offset: 12869},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 449, col: 62, offset: 12872},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 449, col: 64, offset: 12874},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 452, col: 7, offset: 12960},
	run: (*parser).callonExpression38,
	expr: &seqExpr{
	pos: position{line: 452, col: 7, offset: 12960},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 452, col: 7, offset: 12960},
	label: "bindings",
	expr: &oneOrMoreExpr{
	pos: position{line: 452, col: 16, offset: 12969},
	expr: &ruleRefExpr{
	pos: position{line: 452, col: 16, offset: 12969},
	name: "LetBinding",
},
},
},
&ruleRefExpr{
	pos: position{line: 452, col: 28, offset: 12981},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 452, col: 31, offset: 12984},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 452, col: 34, offset: 12987},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 452, col: 36, offset: 12989},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 459, col: 7, offset: 13229},
	run: (*parser).callonExpression47,
	expr: &seqExpr{
	pos: position{line: 459, col: 7, offset: 13229},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 459, col: 7, offset: 13229},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 459, col: 14, offset: 13236},
	name: "_",
},
&litMatcher{
	pos: position{line: 459, col: 16, offset: 13238},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 459, col: 20, offset: 13242},
	name: "_",
},
&labeledExpr{
	pos: position{line: 459, col: 22, offset: 13244},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 459, col: 28, offset: 13250},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 459, col: 45, offset: 13267},
	name: "_",
},
&litMatcher{
	pos: position{line: 459, col: 47, offset: 13269},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 459, col: 51, offset: 13273},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 459, col: 54, offset: 13276},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 459, col: 56, offset: 13278},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 459, col: 67, offset: 13289},
	name: "_",
},
&litMatcher{
	pos: position{line: 459, col: 69, offset: 13291},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 459, col: 73, offset: 13295},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 459, col: 75, offset: 13297},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 459, col: 81, offset: 13303},
	name: "_",
},
&labeledExpr{
	pos: position{line: 459, col: 83, offset: 13305},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 459, col: 88, offset: 13310},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 462, col: 7, offset: 13419},
	run: (*parser).callonExpression67,
	expr: &seqExpr{
	pos: position{line: 462, col: 7, offset: 13419},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 462, col: 7, offset: 13419},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 462, col: 9, offset: 13421},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 462, col: 28, offset: 13440},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 462, col: 30, offset: 13442},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 462, col: 36, offset: 13448},
	name: "_",
},
&labeledExpr{
	pos: position{line: 462, col: 38, offset: 13450},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 462, col: 40, offset: 13452},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 463, col: 7, offset: 13511},
	run: (*parser).callonExpression76,
	expr: &seqExpr{
	pos: position{line: 463, col: 7, offset: 13511},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 463, col: 7, offset: 13511},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 463, col: 13, offset: 13517},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 463, col: 16, offset: 13520},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 463, col: 18, offset: 13522},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 463, col: 35, offset: 13539},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 463, col: 38, offset: 13542},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 463, col: 40, offset: 13544},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 463, col: 57, offset: 13561},
	name: "_",
},
&litMatcher{
	pos: position{line: 463, col: 59, offset: 13563},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 463, col: 63, offset: 13567},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 463, col: 66, offset: 13570},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 463, col: 68, offset: 13572},
	name: "ApplicationExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 466, col: 7, offset: 13693},
	name: "EmptyList",
},
&ruleRefExpr{
	pos: position{line: 467, col: 7, offset: 13709},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 469, col: 1, offset: 13730},
	expr: &actionExpr{
	pos: position{line: 469, col: 14, offset: 13745},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 469, col: 14, offset: 13745},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 469, col: 14, offset: 13745},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 469, col: 18, offset: 13749},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 469, col: 21, offset: 13752},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 469, col: 23, offset: 13754},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 471, col: 1, offset: 13784},
	expr: &actionExpr{
	pos: position{line: 472, col: 1, offset: 13808},
	run: (*parser).callonAnnotatedExpression1,
	expr: &seqExpr{
	pos: position{line: 472, col: 1, offset: 13808},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 472, col: 1, offset: 13808},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 472, col: 3, offset: 13810},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 472, col: 22, offset: 13829},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 472, col: 24, offset: 13831},
	expr: &seqExpr{
	pos: position{line: 472, col: 25, offset: 13832},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 472, col: 25, offset: 13832},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 472, col: 27, offset: 13834},
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
	pos: position{line: 477, col: 1, offset: 13959},
	expr: &actionExpr{
	pos: position{line: 477, col: 13, offset: 13973},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 477, col: 13, offset: 13973},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 477, col: 13, offset: 13973},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 477, col: 17, offset: 13977},
	name: "_",
},
&litMatcher{
	pos: position{line: 477, col: 19, offset: 13979},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 477, col: 23, offset: 13983},
	name: "_",
},
&litMatcher{
	pos: position{line: 477, col: 25, offset: 13985},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 477, col: 29, offset: 13989},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 477, col: 32, offset: 13992},
	name: "List",
},
&ruleRefExpr{
	pos: position{line: 477, col: 37, offset: 13997},
	name: "_",
},
&labeledExpr{
	pos: position{line: 477, col: 39, offset: 13999},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 477, col: 41, offset: 14001},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 481, col: 1, offset: 14064},
	expr: &ruleRefExpr{
	pos: position{line: 481, col: 22, offset: 14087},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 483, col: 1, offset: 14108},
	expr: &actionExpr{
	pos: position{line: 483, col: 26, offset: 14135},
	run: (*parser).callonImportAltExpression1,
	expr: &seqExpr{
	pos: position{line: 483, col: 26, offset: 14135},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 483, col: 26, offset: 14135},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 483, col: 32, offset: 14141},
	name: "OrExpression",
},
},
&labeledExpr{
	pos: position{line: 483, col: 55, offset: 14164},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 483, col: 60, offset: 14169},
	expr: &seqExpr{
	pos: position{line: 483, col: 61, offset: 14170},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 483, col: 61, offset: 14170},
	name: "_",
},
&litMatcher{
	pos: position{line: 483, col: 63, offset: 14172},
	val: "?",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 483, col: 67, offset: 14176},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 483, col: 69, offset: 14178},
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
	pos: position{line: 485, col: 1, offset: 14249},
	expr: &actionExpr{
	pos: position{line: 485, col: 26, offset: 14276},
	run: (*parser).callonOrExpression1,
	expr: &seqExpr{
	pos: position{line: 485, col: 26, offset: 14276},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 485, col: 26, offset: 14276},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 485, col: 32, offset: 14282},
	name: "PlusExpression",
},
},
&labeledExpr{
	pos: position{line: 485, col: 55, offset: 14305},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 485, col: 60, offset: 14310},
	expr: &seqExpr{
	pos: position{line: 485, col: 61, offset: 14311},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 485, col: 61, offset: 14311},
	name: "_",
},
&litMatcher{
	pos: position{line: 485, col: 63, offset: 14313},
	val: "||",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 485, col: 68, offset: 14318},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 485, col: 70, offset: 14320},
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
	pos: position{line: 487, col: 1, offset: 14386},
	expr: &actionExpr{
	pos: position{line: 487, col: 26, offset: 14413},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 487, col: 26, offset: 14413},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 487, col: 26, offset: 14413},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 487, col: 32, offset: 14419},
	name: "TextAppendExpression",
},
},
&labeledExpr{
	pos: position{line: 487, col: 55, offset: 14442},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 487, col: 60, offset: 14447},
	expr: &seqExpr{
	pos: position{line: 487, col: 61, offset: 14448},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 487, col: 61, offset: 14448},
	name: "_",
},
&litMatcher{
	pos: position{line: 487, col: 63, offset: 14450},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 487, col: 67, offset: 14454},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 487, col: 70, offset: 14457},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 487, col: 72, offset: 14459},
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
	pos: position{line: 489, col: 1, offset: 14533},
	expr: &actionExpr{
	pos: position{line: 489, col: 26, offset: 14560},
	run: (*parser).callonTextAppendExpression1,
	expr: &seqExpr{
	pos: position{line: 489, col: 26, offset: 14560},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 489, col: 26, offset: 14560},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 489, col: 32, offset: 14566},
	name: "ListAppendExpression",
},
},
&labeledExpr{
	pos: position{line: 489, col: 55, offset: 14589},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 489, col: 60, offset: 14594},
	expr: &seqExpr{
	pos: position{line: 489, col: 61, offset: 14595},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 489, col: 61, offset: 14595},
	name: "_",
},
&litMatcher{
	pos: position{line: 489, col: 63, offset: 14597},
	val: "++",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 489, col: 68, offset: 14602},
	name: "_",
},
&labeledExpr{
	pos: position{line: 489, col: 70, offset: 14604},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 489, col: 72, offset: 14606},
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
	pos: position{line: 491, col: 1, offset: 14686},
	expr: &actionExpr{
	pos: position{line: 491, col: 26, offset: 14713},
	run: (*parser).callonListAppendExpression1,
	expr: &seqExpr{
	pos: position{line: 491, col: 26, offset: 14713},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 491, col: 26, offset: 14713},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 491, col: 32, offset: 14719},
	name: "AndExpression",
},
},
&labeledExpr{
	pos: position{line: 491, col: 55, offset: 14742},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 491, col: 60, offset: 14747},
	expr: &seqExpr{
	pos: position{line: 491, col: 61, offset: 14748},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 491, col: 61, offset: 14748},
	name: "_",
},
&litMatcher{
	pos: position{line: 491, col: 63, offset: 14750},
	val: "#",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 491, col: 67, offset: 14754},
	name: "_",
},
&labeledExpr{
	pos: position{line: 491, col: 69, offset: 14756},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 491, col: 71, offset: 14758},
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
	pos: position{line: 493, col: 1, offset: 14831},
	expr: &actionExpr{
	pos: position{line: 493, col: 26, offset: 14858},
	run: (*parser).callonAndExpression1,
	expr: &seqExpr{
	pos: position{line: 493, col: 26, offset: 14858},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 493, col: 26, offset: 14858},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 493, col: 32, offset: 14864},
	name: "CombineExpression",
},
},
&labeledExpr{
	pos: position{line: 493, col: 55, offset: 14887},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 493, col: 60, offset: 14892},
	expr: &seqExpr{
	pos: position{line: 493, col: 61, offset: 14893},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 493, col: 61, offset: 14893},
	name: "_",
},
&litMatcher{
	pos: position{line: 493, col: 63, offset: 14895},
	val: "&&",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 493, col: 68, offset: 14900},
	name: "_",
},
&labeledExpr{
	pos: position{line: 493, col: 70, offset: 14902},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 493, col: 72, offset: 14904},
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
	pos: position{line: 495, col: 1, offset: 14974},
	expr: &actionExpr{
	pos: position{line: 495, col: 26, offset: 15001},
	run: (*parser).callonCombineExpression1,
	expr: &seqExpr{
	pos: position{line: 495, col: 26, offset: 15001},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 495, col: 26, offset: 15001},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 495, col: 32, offset: 15007},
	name: "PreferExpression",
},
},
&labeledExpr{
	pos: position{line: 495, col: 55, offset: 15030},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 495, col: 60, offset: 15035},
	expr: &seqExpr{
	pos: position{line: 495, col: 61, offset: 15036},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 495, col: 61, offset: 15036},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 495, col: 63, offset: 15038},
	name: "Combine",
},
&ruleRefExpr{
	pos: position{line: 495, col: 71, offset: 15046},
	name: "_",
},
&labeledExpr{
	pos: position{line: 495, col: 73, offset: 15048},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 495, col: 75, offset: 15050},
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
	pos: position{line: 497, col: 1, offset: 15127},
	expr: &actionExpr{
	pos: position{line: 497, col: 26, offset: 15154},
	run: (*parser).callonPreferExpression1,
	expr: &seqExpr{
	pos: position{line: 497, col: 26, offset: 15154},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 497, col: 26, offset: 15154},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 497, col: 32, offset: 15160},
	name: "CombineTypesExpression",
},
},
&labeledExpr{
	pos: position{line: 497, col: 55, offset: 15183},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 497, col: 60, offset: 15188},
	expr: &seqExpr{
	pos: position{line: 497, col: 61, offset: 15189},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 497, col: 61, offset: 15189},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 497, col: 63, offset: 15191},
	name: "Prefer",
},
&ruleRefExpr{
	pos: position{line: 497, col: 70, offset: 15198},
	name: "_",
},
&labeledExpr{
	pos: position{line: 497, col: 72, offset: 15200},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 497, col: 74, offset: 15202},
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
	pos: position{line: 499, col: 1, offset: 15296},
	expr: &actionExpr{
	pos: position{line: 499, col: 26, offset: 15323},
	run: (*parser).callonCombineTypesExpression1,
	expr: &seqExpr{
	pos: position{line: 499, col: 26, offset: 15323},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 499, col: 26, offset: 15323},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 499, col: 32, offset: 15329},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 499, col: 55, offset: 15352},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 499, col: 60, offset: 15357},
	expr: &seqExpr{
	pos: position{line: 499, col: 61, offset: 15358},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 499, col: 61, offset: 15358},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 499, col: 63, offset: 15360},
	name: "CombineTypes",
},
&ruleRefExpr{
	pos: position{line: 499, col: 76, offset: 15373},
	name: "_",
},
&labeledExpr{
	pos: position{line: 499, col: 78, offset: 15375},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 499, col: 80, offset: 15377},
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
	pos: position{line: 501, col: 1, offset: 15457},
	expr: &actionExpr{
	pos: position{line: 501, col: 26, offset: 15484},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 501, col: 26, offset: 15484},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 501, col: 26, offset: 15484},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 501, col: 32, offset: 15490},
	name: "EqualExpression",
},
},
&labeledExpr{
	pos: position{line: 501, col: 55, offset: 15513},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 501, col: 60, offset: 15518},
	expr: &seqExpr{
	pos: position{line: 501, col: 61, offset: 15519},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 501, col: 61, offset: 15519},
	name: "_",
},
&litMatcher{
	pos: position{line: 501, col: 63, offset: 15521},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 501, col: 67, offset: 15525},
	name: "_",
},
&labeledExpr{
	pos: position{line: 501, col: 69, offset: 15527},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 501, col: 71, offset: 15529},
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
	pos: position{line: 503, col: 1, offset: 15599},
	expr: &actionExpr{
	pos: position{line: 503, col: 26, offset: 15626},
	run: (*parser).callonEqualExpression1,
	expr: &seqExpr{
	pos: position{line: 503, col: 26, offset: 15626},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 503, col: 26, offset: 15626},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 503, col: 32, offset: 15632},
	name: "NotEqualExpression",
},
},
&labeledExpr{
	pos: position{line: 503, col: 55, offset: 15655},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 503, col: 60, offset: 15660},
	expr: &seqExpr{
	pos: position{line: 503, col: 61, offset: 15661},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 503, col: 61, offset: 15661},
	name: "_",
},
&litMatcher{
	pos: position{line: 503, col: 63, offset: 15663},
	val: "==",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 503, col: 68, offset: 15668},
	name: "_",
},
&labeledExpr{
	pos: position{line: 503, col: 70, offset: 15670},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 503, col: 72, offset: 15672},
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
	pos: position{line: 505, col: 1, offset: 15742},
	expr: &actionExpr{
	pos: position{line: 505, col: 26, offset: 15769},
	run: (*parser).callonNotEqualExpression1,
	expr: &seqExpr{
	pos: position{line: 505, col: 26, offset: 15769},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 505, col: 26, offset: 15769},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 505, col: 32, offset: 15775},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 505, col: 55, offset: 15798},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 505, col: 60, offset: 15803},
	expr: &seqExpr{
	pos: position{line: 505, col: 61, offset: 15804},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 505, col: 61, offset: 15804},
	name: "_",
},
&litMatcher{
	pos: position{line: 505, col: 63, offset: 15806},
	val: "!=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 505, col: 68, offset: 15811},
	name: "_",
},
&labeledExpr{
	pos: position{line: 505, col: 70, offset: 15813},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 505, col: 72, offset: 15815},
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
	pos: position{line: 508, col: 1, offset: 15889},
	expr: &actionExpr{
	pos: position{line: 508, col: 25, offset: 15915},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 508, col: 25, offset: 15915},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 508, col: 25, offset: 15915},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 508, col: 27, offset: 15917},
	name: "FirstApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 508, col: 54, offset: 15944},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 508, col: 59, offset: 15949},
	expr: &seqExpr{
	pos: position{line: 508, col: 60, offset: 15950},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 508, col: 60, offset: 15950},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 508, col: 63, offset: 15953},
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
	pos: position{line: 517, col: 1, offset: 16196},
	expr: &choiceExpr{
	pos: position{line: 518, col: 8, offset: 16234},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 518, col: 8, offset: 16234},
	run: (*parser).callonFirstApplicationExpression2,
	expr: &seqExpr{
	pos: position{line: 518, col: 8, offset: 16234},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 518, col: 8, offset: 16234},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 518, col: 14, offset: 16240},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 518, col: 17, offset: 16243},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 518, col: 19, offset: 16245},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 518, col: 36, offset: 16262},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 518, col: 39, offset: 16265},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 518, col: 41, offset: 16267},
	name: "ImportExpression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 521, col: 8, offset: 16370},
	run: (*parser).callonFirstApplicationExpression11,
	expr: &seqExpr{
	pos: position{line: 521, col: 8, offset: 16370},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 521, col: 8, offset: 16370},
	name: "Some",
},
&ruleRefExpr{
	pos: position{line: 521, col: 13, offset: 16375},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 521, col: 16, offset: 16378},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 521, col: 18, offset: 16380},
	name: "ImportExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 522, col: 8, offset: 16435},
	name: "ImportExpression",
},
	},
},
},
{
	name: "ImportExpression",
	pos: position{line: 524, col: 1, offset: 16453},
	expr: &choiceExpr{
	pos: position{line: 524, col: 20, offset: 16474},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 524, col: 20, offset: 16474},
	name: "Import",
},
&ruleRefExpr{
	pos: position{line: 524, col: 29, offset: 16483},
	name: "SelectorExpression",
},
	},
},
},
{
	name: "SelectorExpression",
	pos: position{line: 526, col: 1, offset: 16503},
	expr: &actionExpr{
	pos: position{line: 526, col: 22, offset: 16526},
	run: (*parser).callonSelectorExpression1,
	expr: &seqExpr{
	pos: position{line: 526, col: 22, offset: 16526},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 526, col: 22, offset: 16526},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 526, col: 24, offset: 16528},
	name: "PrimitiveExpression",
},
},
&labeledExpr{
	pos: position{line: 526, col: 44, offset: 16548},
	label: "ls",
	expr: &zeroOrMoreExpr{
	pos: position{line: 526, col: 47, offset: 16551},
	expr: &seqExpr{
	pos: position{line: 526, col: 48, offset: 16552},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 526, col: 48, offset: 16552},
	name: "_",
},
&litMatcher{
	pos: position{line: 526, col: 50, offset: 16554},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 526, col: 54, offset: 16558},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 526, col: 56, offset: 16560},
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
	pos: position{line: 545, col: 1, offset: 17113},
	expr: &choiceExpr{
	pos: position{line: 545, col: 12, offset: 17126},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 545, col: 12, offset: 17126},
	name: "AnyLabel",
},
&ruleRefExpr{
	pos: position{line: 545, col: 23, offset: 17137},
	name: "Labels",
},
&ruleRefExpr{
	pos: position{line: 545, col: 32, offset: 17146},
	name: "TypeSelector",
},
	},
},
},
{
	name: "Labels",
	pos: position{line: 547, col: 1, offset: 17160},
	expr: &actionExpr{
	pos: position{line: 547, col: 10, offset: 17171},
	run: (*parser).callonLabels1,
	expr: &seqExpr{
	pos: position{line: 547, col: 10, offset: 17171},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 547, col: 10, offset: 17171},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 547, col: 14, offset: 17175},
	name: "_",
},
&labeledExpr{
	pos: position{line: 547, col: 16, offset: 17177},
	label: "optclauses",
	expr: &zeroOrOneExpr{
	pos: position{line: 547, col: 27, offset: 17188},
	expr: &seqExpr{
	pos: position{line: 547, col: 29, offset: 17190},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 547, col: 29, offset: 17190},
	name: "AnyLabel",
},
&ruleRefExpr{
	pos: position{line: 547, col: 38, offset: 17199},
	name: "_",
},
&zeroOrMoreExpr{
	pos: position{line: 547, col: 40, offset: 17201},
	expr: &seqExpr{
	pos: position{line: 547, col: 41, offset: 17202},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 547, col: 41, offset: 17202},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 547, col: 45, offset: 17206},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 547, col: 47, offset: 17208},
	name: "AnyLabel",
},
&ruleRefExpr{
	pos: position{line: 547, col: 56, offset: 17217},
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
	pos: position{line: 547, col: 64, offset: 17225},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TypeSelector",
	pos: position{line: 557, col: 1, offset: 17521},
	expr: &actionExpr{
	pos: position{line: 557, col: 16, offset: 17538},
	run: (*parser).callonTypeSelector1,
	expr: &seqExpr{
	pos: position{line: 557, col: 16, offset: 17538},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 557, col: 16, offset: 17538},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 557, col: 20, offset: 17542},
	name: "_",
},
&labeledExpr{
	pos: position{line: 557, col: 22, offset: 17544},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 557, col: 24, offset: 17546},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 557, col: 35, offset: 17557},
	name: "_",
},
&litMatcher{
	pos: position{line: 557, col: 37, offset: 17559},
	val: ")",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PrimitiveExpression",
	pos: position{line: 559, col: 1, offset: 17582},
	expr: &choiceExpr{
	pos: position{line: 560, col: 7, offset: 17612},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 560, col: 7, offset: 17612},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 561, col: 7, offset: 17632},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 562, col: 7, offset: 17653},
	name: "IntegerLiteral",
},
&ruleRefExpr{
	pos: position{line: 563, col: 7, offset: 17674},
	name: "TextLiteral",
},
&actionExpr{
	pos: position{line: 564, col: 7, offset: 17692},
	run: (*parser).callonPrimitiveExpression6,
	expr: &seqExpr{
	pos: position{line: 564, col: 7, offset: 17692},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 564, col: 7, offset: 17692},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 564, col: 11, offset: 17696},
	name: "_",
},
&labeledExpr{
	pos: position{line: 564, col: 13, offset: 17698},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 564, col: 15, offset: 17700},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 564, col: 35, offset: 17720},
	name: "_",
},
&litMatcher{
	pos: position{line: 564, col: 37, offset: 17722},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&actionExpr{
	pos: position{line: 565, col: 7, offset: 17750},
	run: (*parser).callonPrimitiveExpression14,
	expr: &seqExpr{
	pos: position{line: 565, col: 7, offset: 17750},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 565, col: 7, offset: 17750},
	val: "<",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 565, col: 11, offset: 17754},
	name: "_",
},
&labeledExpr{
	pos: position{line: 565, col: 13, offset: 17756},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 565, col: 15, offset: 17758},
	name: "UnionType",
},
},
&ruleRefExpr{
	pos: position{line: 565, col: 25, offset: 17768},
	name: "_",
},
&litMatcher{
	pos: position{line: 565, col: 27, offset: 17770},
	val: ">",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 566, col: 7, offset: 17798},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 567, col: 7, offset: 17824},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 568, col: 7, offset: 17841},
	run: (*parser).callonPrimitiveExpression24,
	expr: &seqExpr{
	pos: position{line: 568, col: 7, offset: 17841},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 568, col: 7, offset: 17841},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 568, col: 11, offset: 17845},
	name: "_",
},
&labeledExpr{
	pos: position{line: 568, col: 14, offset: 17848},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 568, col: 16, offset: 17850},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 568, col: 27, offset: 17861},
	name: "_",
},
&litMatcher{
	pos: position{line: 568, col: 29, offset: 17863},
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
	pos: position{line: 570, col: 1, offset: 17886},
	expr: &choiceExpr{
	pos: position{line: 571, col: 7, offset: 17916},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 571, col: 7, offset: 17916},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &litMatcher{
	pos: position{line: 571, col: 7, offset: 17916},
	val: "=",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 572, col: 7, offset: 17971},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 573, col: 7, offset: 17996},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 574, col: 7, offset: 18024},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 574, col: 7, offset: 18024},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 576, col: 1, offset: 18070},
	expr: &actionExpr{
	pos: position{line: 576, col: 19, offset: 18090},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 576, col: 19, offset: 18090},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 576, col: 19, offset: 18090},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 576, col: 24, offset: 18095},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 576, col: 33, offset: 18104},
	name: "_",
},
&litMatcher{
	pos: position{line: 576, col: 35, offset: 18106},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 576, col: 39, offset: 18110},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 576, col: 42, offset: 18113},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 576, col: 47, offset: 18118},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 579, col: 1, offset: 18175},
	expr: &actionExpr{
	pos: position{line: 579, col: 18, offset: 18194},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 579, col: 18, offset: 18194},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 579, col: 18, offset: 18194},
	name: "_",
},
&litMatcher{
	pos: position{line: 579, col: 20, offset: 18196},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 579, col: 24, offset: 18200},
	name: "_",
},
&labeledExpr{
	pos: position{line: 579, col: 26, offset: 18202},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 579, col: 28, offset: 18204},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 580, col: 1, offset: 18236},
	expr: &actionExpr{
	pos: position{line: 581, col: 7, offset: 18265},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 581, col: 7, offset: 18265},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 581, col: 7, offset: 18265},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 581, col: 13, offset: 18271},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 581, col: 29, offset: 18287},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 581, col: 34, offset: 18292},
	expr: &ruleRefExpr{
	pos: position{line: 581, col: 34, offset: 18292},
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
	pos: position{line: 595, col: 1, offset: 18876},
	expr: &actionExpr{
	pos: position{line: 595, col: 22, offset: 18899},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 595, col: 22, offset: 18899},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 595, col: 22, offset: 18899},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 595, col: 27, offset: 18904},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 595, col: 36, offset: 18913},
	name: "_",
},
&litMatcher{
	pos: position{line: 595, col: 38, offset: 18915},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 595, col: 42, offset: 18919},
	name: "_",
},
&labeledExpr{
	pos: position{line: 595, col: 44, offset: 18921},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 595, col: 49, offset: 18926},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 598, col: 1, offset: 18983},
	expr: &actionExpr{
	pos: position{line: 598, col: 21, offset: 19005},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 598, col: 21, offset: 19005},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 598, col: 21, offset: 19005},
	name: "_",
},
&litMatcher{
	pos: position{line: 598, col: 23, offset: 19007},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 598, col: 27, offset: 19011},
	name: "_",
},
&labeledExpr{
	pos: position{line: 598, col: 29, offset: 19013},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 598, col: 31, offset: 19015},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 599, col: 1, offset: 19050},
	expr: &actionExpr{
	pos: position{line: 600, col: 7, offset: 19082},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 600, col: 7, offset: 19082},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 600, col: 7, offset: 19082},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 600, col: 13, offset: 19088},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 600, col: 32, offset: 19107},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 600, col: 37, offset: 19112},
	expr: &ruleRefExpr{
	pos: position{line: 600, col: 37, offset: 19112},
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
	pos: position{line: 614, col: 1, offset: 19702},
	expr: &choiceExpr{
	pos: position{line: 614, col: 13, offset: 19716},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 614, col: 13, offset: 19716},
	name: "NonEmptyUnionType",
},
&ruleRefExpr{
	pos: position{line: 614, col: 33, offset: 19736},
	name: "EmptyUnionType",
},
	},
},
},
{
	name: "EmptyUnionType",
	pos: position{line: 616, col: 1, offset: 19752},
	expr: &actionExpr{
	pos: position{line: 616, col: 18, offset: 19771},
	run: (*parser).callonEmptyUnionType1,
	expr: &litMatcher{
	pos: position{line: 616, col: 18, offset: 19771},
	val: "",
	ignoreCase: false,
},
},
},
{
	name: "NonEmptyUnionType",
	pos: position{line: 618, col: 1, offset: 19803},
	expr: &actionExpr{
	pos: position{line: 618, col: 21, offset: 19825},
	run: (*parser).callonNonEmptyUnionType1,
	expr: &seqExpr{
	pos: position{line: 618, col: 21, offset: 19825},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 618, col: 21, offset: 19825},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 618, col: 27, offset: 19831},
	name: "UnionVariant",
},
},
&labeledExpr{
	pos: position{line: 618, col: 40, offset: 19844},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 618, col: 45, offset: 19849},
	expr: &seqExpr{
	pos: position{line: 618, col: 46, offset: 19850},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 618, col: 46, offset: 19850},
	name: "_",
},
&litMatcher{
	pos: position{line: 618, col: 48, offset: 19852},
	val: "|",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 618, col: 52, offset: 19856},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 618, col: 54, offset: 19858},
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
	pos: position{line: 638, col: 1, offset: 20580},
	expr: &seqExpr{
	pos: position{line: 638, col: 16, offset: 20597},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 638, col: 16, offset: 20597},
	name: "AnyLabel",
},
&zeroOrOneExpr{
	pos: position{line: 638, col: 25, offset: 20606},
	expr: &seqExpr{
	pos: position{line: 638, col: 26, offset: 20607},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 638, col: 26, offset: 20607},
	name: "_",
},
&litMatcher{
	pos: position{line: 638, col: 28, offset: 20609},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 638, col: 32, offset: 20613},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 638, col: 35, offset: 20616},
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
	pos: position{line: 640, col: 1, offset: 20630},
	expr: &actionExpr{
	pos: position{line: 640, col: 12, offset: 20643},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 640, col: 12, offset: 20643},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 640, col: 12, offset: 20643},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 640, col: 16, offset: 20647},
	name: "_",
},
&labeledExpr{
	pos: position{line: 640, col: 18, offset: 20649},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 640, col: 20, offset: 20651},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 640, col: 31, offset: 20662},
	name: "_",
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 642, col: 1, offset: 20681},
	expr: &actionExpr{
	pos: position{line: 643, col: 7, offset: 20711},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 643, col: 7, offset: 20711},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 643, col: 7, offset: 20711},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 643, col: 11, offset: 20715},
	name: "_",
},
&labeledExpr{
	pos: position{line: 643, col: 13, offset: 20717},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 643, col: 19, offset: 20723},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 643, col: 30, offset: 20734},
	name: "_",
},
&labeledExpr{
	pos: position{line: 643, col: 32, offset: 20736},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 643, col: 37, offset: 20741},
	expr: &ruleRefExpr{
	pos: position{line: 643, col: 37, offset: 20741},
	name: "MoreList",
},
},
},
&litMatcher{
	pos: position{line: 643, col: 47, offset: 20751},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 653, col: 1, offset: 21027},
	expr: &notExpr{
	pos: position{line: 653, col: 7, offset: 21035},
	expr: &anyMatcher{
	line: 653, col: 8, offset: 21036,
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

