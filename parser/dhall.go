
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
&ruleRefExpr{
	pos: position{line: 102, col: 12, offset: 2354},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 102, col: 19, offset: 2361},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 102, col: 26, offset: 2368},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 102, col: 33, offset: 2375},
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
	pos: position{line: 107, col: 1, offset: 2507},
	expr: &choiceExpr{
	pos: position{line: 108, col: 6, offset: 2532},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 108, col: 6, offset: 2532},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 109, col: 6, offset: 2549},
	val: "[\\x23-\\x5b]",
	ranges: []rune{'#','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 110, col: 6, offset: 2566},
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
	pos: position{line: 112, col: 1, offset: 2585},
	expr: &actionExpr{
	pos: position{line: 112, col: 22, offset: 2608},
	run: (*parser).callonDoubleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 112, col: 22, offset: 2608},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 112, col: 22, offset: 2608},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 112, col: 26, offset: 2612},
	label: "chunks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 112, col: 33, offset: 2619},
	expr: &ruleRefExpr{
	pos: position{line: 112, col: 33, offset: 2619},
	name: "DoubleQuoteChunk",
},
},
},
&litMatcher{
	pos: position{line: 112, col: 51, offset: 2637},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "SingleQuoteContinue",
	pos: position{line: 129, col: 1, offset: 3105},
	expr: &choiceExpr{
	pos: position{line: 130, col: 7, offset: 3135},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 130, col: 7, offset: 3135},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 130, col: 7, offset: 3135},
	name: "Interpolation",
},
&ruleRefExpr{
	pos: position{line: 130, col: 21, offset: 3149},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 131, col: 7, offset: 3175},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 131, col: 7, offset: 3175},
	name: "EscapedQuotePair",
},
&ruleRefExpr{
	pos: position{line: 131, col: 24, offset: 3192},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 132, col: 7, offset: 3218},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 132, col: 7, offset: 3218},
	name: "EscapedInterpolation",
},
&ruleRefExpr{
	pos: position{line: 132, col: 28, offset: 3239},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 133, col: 7, offset: 3265},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 133, col: 7, offset: 3265},
	name: "SingleQuoteChar",
},
&ruleRefExpr{
	pos: position{line: 133, col: 23, offset: 3281},
	name: "SingleQuoteContinue",
},
	},
},
&litMatcher{
	pos: position{line: 134, col: 7, offset: 3307},
	val: "''",
	ignoreCase: false,
},
	},
},
},
{
	name: "EscapedQuotePair",
	pos: position{line: 136, col: 1, offset: 3313},
	expr: &actionExpr{
	pos: position{line: 136, col: 20, offset: 3334},
	run: (*parser).callonEscapedQuotePair1,
	expr: &litMatcher{
	pos: position{line: 136, col: 20, offset: 3334},
	val: "'''",
	ignoreCase: false,
},
},
},
{
	name: "EscapedInterpolation",
	pos: position{line: 140, col: 1, offset: 3469},
	expr: &actionExpr{
	pos: position{line: 140, col: 24, offset: 3494},
	run: (*parser).callonEscapedInterpolation1,
	expr: &litMatcher{
	pos: position{line: 140, col: 24, offset: 3494},
	val: "''${",
	ignoreCase: false,
},
},
},
{
	name: "SingleQuoteChar",
	pos: position{line: 142, col: 1, offset: 3536},
	expr: &choiceExpr{
	pos: position{line: 143, col: 6, offset: 3561},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 143, col: 6, offset: 3561},
	val: "[\\x20-\\U0010ffff]",
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 144, col: 6, offset: 3584},
	val: "\t",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 145, col: 6, offset: 3594},
	name: "EOL",
},
	},
},
},
{
	name: "SingleQuoteLiteral",
	pos: position{line: 147, col: 1, offset: 3599},
	expr: &actionExpr{
	pos: position{line: 147, col: 22, offset: 3622},
	run: (*parser).callonSingleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 147, col: 22, offset: 3622},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 147, col: 22, offset: 3622},
	val: "''",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 147, col: 27, offset: 3627},
	name: "EOL",
},
&labeledExpr{
	pos: position{line: 147, col: 31, offset: 3631},
	label: "content",
	expr: &ruleRefExpr{
	pos: position{line: 147, col: 39, offset: 3639},
	name: "SingleQuoteContinue",
},
},
	},
},
},
},
{
	name: "Interpolation",
	pos: position{line: 165, col: 1, offset: 4189},
	expr: &actionExpr{
	pos: position{line: 165, col: 17, offset: 4207},
	run: (*parser).callonInterpolation1,
	expr: &seqExpr{
	pos: position{line: 165, col: 17, offset: 4207},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 165, col: 17, offset: 4207},
	val: "${",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 165, col: 22, offset: 4212},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 165, col: 24, offset: 4214},
	name: "CompleteExpression",
},
},
&litMatcher{
	pos: position{line: 165, col: 43, offset: 4233},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TextLiteral",
	pos: position{line: 167, col: 1, offset: 4256},
	expr: &choiceExpr{
	pos: position{line: 167, col: 15, offset: 4272},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 167, col: 15, offset: 4272},
	name: "DoubleQuoteLiteral",
},
&ruleRefExpr{
	pos: position{line: 167, col: 36, offset: 4293},
	name: "SingleQuoteLiteral",
},
	},
},
},
{
	name: "Reserved",
	pos: position{line: 170, col: 1, offset: 4398},
	expr: &choiceExpr{
	pos: position{line: 171, col: 5, offset: 4415},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 171, col: 5, offset: 4415},
	run: (*parser).callonReserved2,
	expr: &litMatcher{
	pos: position{line: 171, col: 5, offset: 4415},
	val: "Natural/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 172, col: 5, offset: 4464},
	run: (*parser).callonReserved4,
	expr: &litMatcher{
	pos: position{line: 172, col: 5, offset: 4464},
	val: "Natural/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 173, col: 5, offset: 4511},
	run: (*parser).callonReserved6,
	expr: &litMatcher{
	pos: position{line: 173, col: 5, offset: 4511},
	val: "Natural/isZero",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 174, col: 5, offset: 4562},
	run: (*parser).callonReserved8,
	expr: &litMatcher{
	pos: position{line: 174, col: 5, offset: 4562},
	val: "Natural/even",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 175, col: 5, offset: 4609},
	run: (*parser).callonReserved10,
	expr: &litMatcher{
	pos: position{line: 175, col: 5, offset: 4609},
	val: "Natural/odd",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 176, col: 5, offset: 4654},
	run: (*parser).callonReserved12,
	expr: &litMatcher{
	pos: position{line: 176, col: 5, offset: 4654},
	val: "Natural/toInteger",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 177, col: 5, offset: 4711},
	run: (*parser).callonReserved14,
	expr: &litMatcher{
	pos: position{line: 177, col: 5, offset: 4711},
	val: "Natural/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 178, col: 5, offset: 4758},
	run: (*parser).callonReserved16,
	expr: &litMatcher{
	pos: position{line: 178, col: 5, offset: 4758},
	val: "Integer/toDouble",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 179, col: 5, offset: 4813},
	run: (*parser).callonReserved18,
	expr: &litMatcher{
	pos: position{line: 179, col: 5, offset: 4813},
	val: "Integer/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 180, col: 5, offset: 4860},
	run: (*parser).callonReserved20,
	expr: &litMatcher{
	pos: position{line: 180, col: 5, offset: 4860},
	val: "Double/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 181, col: 5, offset: 4905},
	run: (*parser).callonReserved22,
	expr: &litMatcher{
	pos: position{line: 181, col: 5, offset: 4905},
	val: "List/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 182, col: 5, offset: 4948},
	run: (*parser).callonReserved24,
	expr: &litMatcher{
	pos: position{line: 182, col: 5, offset: 4948},
	val: "List/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 183, col: 5, offset: 4989},
	run: (*parser).callonReserved26,
	expr: &litMatcher{
	pos: position{line: 183, col: 5, offset: 4989},
	val: "List/length",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 184, col: 5, offset: 5034},
	run: (*parser).callonReserved28,
	expr: &litMatcher{
	pos: position{line: 184, col: 5, offset: 5034},
	val: "List/head",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 185, col: 5, offset: 5075},
	run: (*parser).callonReserved30,
	expr: &litMatcher{
	pos: position{line: 185, col: 5, offset: 5075},
	val: "List/last",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 186, col: 5, offset: 5116},
	run: (*parser).callonReserved32,
	expr: &litMatcher{
	pos: position{line: 186, col: 5, offset: 5116},
	val: "List/indexed",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 187, col: 5, offset: 5163},
	run: (*parser).callonReserved34,
	expr: &litMatcher{
	pos: position{line: 187, col: 5, offset: 5163},
	val: "List/reverse",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 188, col: 5, offset: 5210},
	run: (*parser).callonReserved36,
	expr: &litMatcher{
	pos: position{line: 188, col: 5, offset: 5210},
	val: "Optional/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 189, col: 5, offset: 5261},
	run: (*parser).callonReserved38,
	expr: &litMatcher{
	pos: position{line: 189, col: 5, offset: 5261},
	val: "Optional/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 190, col: 5, offset: 5310},
	run: (*parser).callonReserved40,
	expr: &litMatcher{
	pos: position{line: 190, col: 5, offset: 5310},
	val: "Text/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 191, col: 5, offset: 5351},
	run: (*parser).callonReserved42,
	expr: &litMatcher{
	pos: position{line: 191, col: 5, offset: 5351},
	val: "Bool",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 192, col: 5, offset: 5383},
	run: (*parser).callonReserved44,
	expr: &litMatcher{
	pos: position{line: 192, col: 5, offset: 5383},
	val: "True",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 193, col: 5, offset: 5415},
	run: (*parser).callonReserved46,
	expr: &litMatcher{
	pos: position{line: 193, col: 5, offset: 5415},
	val: "False",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 194, col: 5, offset: 5449},
	run: (*parser).callonReserved48,
	expr: &litMatcher{
	pos: position{line: 194, col: 5, offset: 5449},
	val: "Optional",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 195, col: 5, offset: 5489},
	run: (*parser).callonReserved50,
	expr: &litMatcher{
	pos: position{line: 195, col: 5, offset: 5489},
	val: "Natural",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 196, col: 5, offset: 5527},
	run: (*parser).callonReserved52,
	expr: &litMatcher{
	pos: position{line: 196, col: 5, offset: 5527},
	val: "Integer",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 197, col: 5, offset: 5565},
	run: (*parser).callonReserved54,
	expr: &litMatcher{
	pos: position{line: 197, col: 5, offset: 5565},
	val: "Double",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 198, col: 5, offset: 5601},
	run: (*parser).callonReserved56,
	expr: &litMatcher{
	pos: position{line: 198, col: 5, offset: 5601},
	val: "Text",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 199, col: 5, offset: 5633},
	run: (*parser).callonReserved58,
	expr: &litMatcher{
	pos: position{line: 199, col: 5, offset: 5633},
	val: "List",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 200, col: 5, offset: 5665},
	run: (*parser).callonReserved60,
	expr: &litMatcher{
	pos: position{line: 200, col: 5, offset: 5665},
	val: "None",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 201, col: 5, offset: 5697},
	run: (*parser).callonReserved62,
	expr: &litMatcher{
	pos: position{line: 201, col: 5, offset: 5697},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 202, col: 5, offset: 5729},
	run: (*parser).callonReserved64,
	expr: &litMatcher{
	pos: position{line: 202, col: 5, offset: 5729},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 203, col: 5, offset: 5761},
	run: (*parser).callonReserved66,
	expr: &litMatcher{
	pos: position{line: 203, col: 5, offset: 5761},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "If",
	pos: position{line: 205, col: 1, offset: 5790},
	expr: &litMatcher{
	pos: position{line: 205, col: 6, offset: 5797},
	val: "if",
	ignoreCase: false,
},
},
{
	name: "Then",
	pos: position{line: 206, col: 1, offset: 5802},
	expr: &litMatcher{
	pos: position{line: 206, col: 8, offset: 5811},
	val: "then",
	ignoreCase: false,
},
},
{
	name: "Else",
	pos: position{line: 207, col: 1, offset: 5818},
	expr: &litMatcher{
	pos: position{line: 207, col: 8, offset: 5827},
	val: "else",
	ignoreCase: false,
},
},
{
	name: "Let",
	pos: position{line: 208, col: 1, offset: 5834},
	expr: &litMatcher{
	pos: position{line: 208, col: 7, offset: 5842},
	val: "let",
	ignoreCase: false,
},
},
{
	name: "In",
	pos: position{line: 209, col: 1, offset: 5848},
	expr: &litMatcher{
	pos: position{line: 209, col: 6, offset: 5855},
	val: "in",
	ignoreCase: false,
},
},
{
	name: "As",
	pos: position{line: 210, col: 1, offset: 5860},
	expr: &litMatcher{
	pos: position{line: 210, col: 6, offset: 5867},
	val: "as",
	ignoreCase: false,
},
},
{
	name: "Using",
	pos: position{line: 211, col: 1, offset: 5872},
	expr: &litMatcher{
	pos: position{line: 211, col: 9, offset: 5882},
	val: "using",
	ignoreCase: false,
},
},
{
	name: "Merge",
	pos: position{line: 212, col: 1, offset: 5890},
	expr: &litMatcher{
	pos: position{line: 212, col: 9, offset: 5900},
	val: "merge",
	ignoreCase: false,
},
},
{
	name: "Missing",
	pos: position{line: 213, col: 1, offset: 5908},
	expr: &actionExpr{
	pos: position{line: 213, col: 11, offset: 5920},
	run: (*parser).callonMissing1,
	expr: &litMatcher{
	pos: position{line: 213, col: 11, offset: 5920},
	val: "missing",
	ignoreCase: false,
},
},
},
{
	name: "True",
	pos: position{line: 214, col: 1, offset: 5956},
	expr: &litMatcher{
	pos: position{line: 214, col: 8, offset: 5965},
	val: "True",
	ignoreCase: false,
},
},
{
	name: "False",
	pos: position{line: 215, col: 1, offset: 5972},
	expr: &litMatcher{
	pos: position{line: 215, col: 9, offset: 5982},
	val: "False",
	ignoreCase: false,
},
},
{
	name: "Infinity",
	pos: position{line: 216, col: 1, offset: 5990},
	expr: &litMatcher{
	pos: position{line: 216, col: 12, offset: 6003},
	val: "Infinity",
	ignoreCase: false,
},
},
{
	name: "NaN",
	pos: position{line: 217, col: 1, offset: 6014},
	expr: &litMatcher{
	pos: position{line: 217, col: 7, offset: 6022},
	val: "NaN",
	ignoreCase: false,
},
},
{
	name: "Some",
	pos: position{line: 218, col: 1, offset: 6028},
	expr: &litMatcher{
	pos: position{line: 218, col: 8, offset: 6037},
	val: "Some",
	ignoreCase: false,
},
},
{
	name: "Keyword",
	pos: position{line: 220, col: 1, offset: 6045},
	expr: &choiceExpr{
	pos: position{line: 221, col: 5, offset: 6061},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 221, col: 5, offset: 6061},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 221, col: 10, offset: 6066},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 221, col: 17, offset: 6073},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 222, col: 5, offset: 6082},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 222, col: 11, offset: 6088},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 223, col: 5, offset: 6095},
	name: "Using",
},
&ruleRefExpr{
	pos: position{line: 223, col: 13, offset: 6103},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 223, col: 23, offset: 6113},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 224, col: 5, offset: 6120},
	name: "True",
},
&ruleRefExpr{
	pos: position{line: 224, col: 12, offset: 6127},
	name: "False",
},
&ruleRefExpr{
	pos: position{line: 225, col: 5, offset: 6137},
	name: "Infinity",
},
&ruleRefExpr{
	pos: position{line: 225, col: 16, offset: 6148},
	name: "NaN",
},
&ruleRefExpr{
	pos: position{line: 226, col: 5, offset: 6156},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 226, col: 13, offset: 6164},
	name: "Some",
},
	},
},
},
{
	name: "Optional",
	pos: position{line: 228, col: 1, offset: 6170},
	expr: &litMatcher{
	pos: position{line: 228, col: 12, offset: 6183},
	val: "Optional",
	ignoreCase: false,
},
},
{
	name: "Text",
	pos: position{line: 229, col: 1, offset: 6194},
	expr: &litMatcher{
	pos: position{line: 229, col: 8, offset: 6203},
	val: "Text",
	ignoreCase: false,
},
},
{
	name: "List",
	pos: position{line: 230, col: 1, offset: 6210},
	expr: &litMatcher{
	pos: position{line: 230, col: 8, offset: 6219},
	val: "List",
	ignoreCase: false,
},
},
{
	name: "Combine",
	pos: position{line: 232, col: 1, offset: 6227},
	expr: &choiceExpr{
	pos: position{line: 232, col: 11, offset: 6239},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 232, col: 11, offset: 6239},
	val: "/\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 232, col: 19, offset: 6247},
	val: "∧",
	ignoreCase: false,
},
	},
},
},
{
	name: "CombineTypes",
	pos: position{line: 233, col: 1, offset: 6253},
	expr: &choiceExpr{
	pos: position{line: 233, col: 16, offset: 6270},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 233, col: 16, offset: 6270},
	val: "//\\\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 233, col: 27, offset: 6281},
	val: "⩓",
	ignoreCase: false,
},
	},
},
},
{
	name: "Prefer",
	pos: position{line: 234, col: 1, offset: 6287},
	expr: &choiceExpr{
	pos: position{line: 234, col: 10, offset: 6298},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 234, col: 10, offset: 6298},
	val: "//",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 234, col: 17, offset: 6305},
	val: "⫽",
	ignoreCase: false,
},
	},
},
},
{
	name: "Lambda",
	pos: position{line: 235, col: 1, offset: 6311},
	expr: &choiceExpr{
	pos: position{line: 235, col: 10, offset: 6322},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 235, col: 10, offset: 6322},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 235, col: 17, offset: 6329},
	val: "λ",
	ignoreCase: false,
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 236, col: 1, offset: 6334},
	expr: &choiceExpr{
	pos: position{line: 236, col: 10, offset: 6345},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 236, col: 10, offset: 6345},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 236, col: 21, offset: 6356},
	val: "∀",
	ignoreCase: false,
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 237, col: 1, offset: 6362},
	expr: &choiceExpr{
	pos: position{line: 237, col: 9, offset: 6372},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 237, col: 9, offset: 6372},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 237, col: 16, offset: 6379},
	val: "→",
	ignoreCase: false,
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 239, col: 1, offset: 6386},
	expr: &seqExpr{
	pos: position{line: 239, col: 12, offset: 6399},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 239, col: 12, offset: 6399},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 239, col: 17, offset: 6404},
	expr: &charClassMatcher{
	pos: position{line: 239, col: 17, offset: 6404},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 239, col: 23, offset: 6410},
	expr: &ruleRefExpr{
	pos: position{line: 239, col: 23, offset: 6410},
	name: "Digit",
},
},
	},
},
},
{
	name: "NumericDoubleLiteral",
	pos: position{line: 241, col: 1, offset: 6418},
	expr: &actionExpr{
	pos: position{line: 241, col: 24, offset: 6443},
	run: (*parser).callonNumericDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 241, col: 24, offset: 6443},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 241, col: 24, offset: 6443},
	expr: &charClassMatcher{
	pos: position{line: 241, col: 24, offset: 6443},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 241, col: 30, offset: 6449},
	expr: &ruleRefExpr{
	pos: position{line: 241, col: 30, offset: 6449},
	name: "Digit",
},
},
&choiceExpr{
	pos: position{line: 241, col: 39, offset: 6458},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 241, col: 39, offset: 6458},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 241, col: 39, offset: 6458},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 241, col: 43, offset: 6462},
	expr: &ruleRefExpr{
	pos: position{line: 241, col: 43, offset: 6462},
	name: "Digit",
},
},
&zeroOrOneExpr{
	pos: position{line: 241, col: 50, offset: 6469},
	expr: &ruleRefExpr{
	pos: position{line: 241, col: 50, offset: 6469},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 241, col: 62, offset: 6481},
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
	pos: position{line: 249, col: 1, offset: 6637},
	expr: &choiceExpr{
	pos: position{line: 249, col: 17, offset: 6655},
	alternatives: []interface{}{
&labeledExpr{
	pos: position{line: 249, col: 17, offset: 6655},
	label: "d",
	expr: &ruleRefExpr{
	pos: position{line: 249, col: 19, offset: 6657},
	name: "NumericDoubleLiteral",
},
},
&actionExpr{
	pos: position{line: 250, col: 5, offset: 6682},
	run: (*parser).callonDoubleLiteral4,
	expr: &ruleRefExpr{
	pos: position{line: 250, col: 5, offset: 6682},
	name: "Infinity",
},
},
&actionExpr{
	pos: position{line: 251, col: 5, offset: 6734},
	run: (*parser).callonDoubleLiteral6,
	expr: &seqExpr{
	pos: position{line: 251, col: 5, offset: 6734},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 251, col: 5, offset: 6734},
	val: "-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 251, col: 9, offset: 6738},
	name: "Infinity",
},
	},
},
},
&actionExpr{
	pos: position{line: 252, col: 5, offset: 6791},
	run: (*parser).callonDoubleLiteral10,
	expr: &ruleRefExpr{
	pos: position{line: 252, col: 5, offset: 6791},
	name: "NaN",
},
},
	},
},
},
{
	name: "NaturalLiteral",
	pos: position{line: 254, col: 1, offset: 6834},
	expr: &actionExpr{
	pos: position{line: 254, col: 18, offset: 6853},
	run: (*parser).callonNaturalLiteral1,
	expr: &oneOrMoreExpr{
	pos: position{line: 254, col: 18, offset: 6853},
	expr: &ruleRefExpr{
	pos: position{line: 254, col: 18, offset: 6853},
	name: "Digit",
},
},
},
},
{
	name: "IntegerLiteral",
	pos: position{line: 259, col: 1, offset: 6942},
	expr: &actionExpr{
	pos: position{line: 259, col: 18, offset: 6961},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 259, col: 18, offset: 6961},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 259, col: 18, offset: 6961},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 259, col: 22, offset: 6965},
	name: "NaturalLiteral",
},
	},
},
},
},
{
	name: "DeBruijn",
	pos: position{line: 267, col: 1, offset: 7117},
	expr: &actionExpr{
	pos: position{line: 267, col: 12, offset: 7130},
	run: (*parser).callonDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 267, col: 12, offset: 7130},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 267, col: 12, offset: 7130},
	name: "_",
},
&litMatcher{
	pos: position{line: 267, col: 14, offset: 7132},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 267, col: 18, offset: 7136},
	name: "_",
},
&labeledExpr{
	pos: position{line: 267, col: 20, offset: 7138},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 267, col: 26, offset: 7144},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Variable",
	pos: position{line: 269, col: 1, offset: 7200},
	expr: &actionExpr{
	pos: position{line: 269, col: 12, offset: 7213},
	run: (*parser).callonVariable1,
	expr: &seqExpr{
	pos: position{line: 269, col: 12, offset: 7213},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 269, col: 12, offset: 7213},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 269, col: 17, offset: 7218},
	name: "NonreservedLabel",
},
},
&labeledExpr{
	pos: position{line: 269, col: 34, offset: 7235},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 269, col: 40, offset: 7241},
	expr: &ruleRefExpr{
	pos: position{line: 269, col: 40, offset: 7241},
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
	pos: position{line: 277, col: 1, offset: 7404},
	expr: &choiceExpr{
	pos: position{line: 277, col: 14, offset: 7419},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 277, col: 14, offset: 7419},
	name: "Variable",
},
&ruleRefExpr{
	pos: position{line: 277, col: 25, offset: 7430},
	name: "Reserved",
},
	},
},
},
{
	name: "PathCharacter",
	pos: position{line: 279, col: 1, offset: 7440},
	expr: &choiceExpr{
	pos: position{line: 280, col: 6, offset: 7463},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 280, col: 6, offset: 7463},
	val: "!",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 281, col: 6, offset: 7475},
	val: "[\\x24-\\x27]",
	ranges: []rune{'$','\'',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 282, col: 6, offset: 7492},
	val: "[\\x2a-\\x2b]",
	ranges: []rune{'*','+',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 283, col: 6, offset: 7509},
	val: "[\\x2d-\\x2e]",
	ranges: []rune{'-','.',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 284, col: 6, offset: 7526},
	val: "[\\x30-\\x3b]",
	ranges: []rune{'0',';',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 285, col: 6, offset: 7543},
	val: "=",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 286, col: 6, offset: 7555},
	val: "[\\x40-\\x5a]",
	ranges: []rune{'@','Z',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 287, col: 6, offset: 7572},
	val: "[\\x5e-\\x7a]",
	ranges: []rune{'^','z',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 288, col: 6, offset: 7589},
	val: "|",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 289, col: 6, offset: 7601},
	val: "~",
	ignoreCase: false,
},
	},
},
},
{
	name: "UnquotedPathComponent",
	pos: position{line: 291, col: 1, offset: 7609},
	expr: &actionExpr{
	pos: position{line: 291, col: 25, offset: 7635},
	run: (*parser).callonUnquotedPathComponent1,
	expr: &oneOrMoreExpr{
	pos: position{line: 291, col: 25, offset: 7635},
	expr: &ruleRefExpr{
	pos: position{line: 291, col: 25, offset: 7635},
	name: "PathCharacter",
},
},
},
},
{
	name: "PathComponent",
	pos: position{line: 293, col: 1, offset: 7682},
	expr: &actionExpr{
	pos: position{line: 293, col: 17, offset: 7700},
	run: (*parser).callonPathComponent1,
	expr: &seqExpr{
	pos: position{line: 293, col: 17, offset: 7700},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 293, col: 17, offset: 7700},
	val: "/",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 293, col: 21, offset: 7704},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 293, col: 23, offset: 7706},
	name: "UnquotedPathComponent",
},
},
	},
},
},
},
{
	name: "Path",
	pos: position{line: 295, col: 1, offset: 7747},
	expr: &actionExpr{
	pos: position{line: 295, col: 8, offset: 7756},
	run: (*parser).callonPath1,
	expr: &labeledExpr{
	pos: position{line: 295, col: 8, offset: 7756},
	label: "cs",
	expr: &oneOrMoreExpr{
	pos: position{line: 295, col: 11, offset: 7759},
	expr: &ruleRefExpr{
	pos: position{line: 295, col: 11, offset: 7759},
	name: "PathComponent",
},
},
},
},
},
{
	name: "Local",
	pos: position{line: 304, col: 1, offset: 8033},
	expr: &choiceExpr{
	pos: position{line: 304, col: 9, offset: 8043},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 304, col: 9, offset: 8043},
	name: "ParentPath",
},
&ruleRefExpr{
	pos: position{line: 304, col: 22, offset: 8056},
	name: "HerePath",
},
&ruleRefExpr{
	pos: position{line: 304, col: 33, offset: 8067},
	name: "HomePath",
},
&ruleRefExpr{
	pos: position{line: 304, col: 44, offset: 8078},
	name: "AbsolutePath",
},
	},
},
},
{
	name: "ParentPath",
	pos: position{line: 306, col: 1, offset: 8092},
	expr: &actionExpr{
	pos: position{line: 306, col: 14, offset: 8107},
	run: (*parser).callonParentPath1,
	expr: &seqExpr{
	pos: position{line: 306, col: 14, offset: 8107},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 306, col: 14, offset: 8107},
	val: "..",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 306, col: 19, offset: 8112},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 306, col: 21, offset: 8114},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HerePath",
	pos: position{line: 307, col: 1, offset: 8170},
	expr: &actionExpr{
	pos: position{line: 307, col: 12, offset: 8183},
	run: (*parser).callonHerePath1,
	expr: &seqExpr{
	pos: position{line: 307, col: 12, offset: 8183},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 307, col: 12, offset: 8183},
	val: ".",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 307, col: 16, offset: 8187},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 307, col: 18, offset: 8189},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HomePath",
	pos: position{line: 308, col: 1, offset: 8228},
	expr: &actionExpr{
	pos: position{line: 308, col: 12, offset: 8241},
	run: (*parser).callonHomePath1,
	expr: &seqExpr{
	pos: position{line: 308, col: 12, offset: 8241},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 308, col: 12, offset: 8241},
	val: "~",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 308, col: 16, offset: 8245},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 308, col: 18, offset: 8247},
	name: "Path",
},
},
	},
},
},
},
{
	name: "AbsolutePath",
	pos: position{line: 309, col: 1, offset: 8302},
	expr: &actionExpr{
	pos: position{line: 309, col: 16, offset: 8319},
	run: (*parser).callonAbsolutePath1,
	expr: &labeledExpr{
	pos: position{line: 309, col: 16, offset: 8319},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 309, col: 18, offset: 8321},
	name: "Path",
},
},
},
},
{
	name: "Scheme",
	pos: position{line: 311, col: 1, offset: 8377},
	expr: &seqExpr{
	pos: position{line: 311, col: 10, offset: 8388},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 311, col: 10, offset: 8388},
	val: "http",
	ignoreCase: false,
},
&zeroOrOneExpr{
	pos: position{line: 311, col: 17, offset: 8395},
	expr: &litMatcher{
	pos: position{line: 311, col: 17, offset: 8395},
	val: "s",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "HttpRaw",
	pos: position{line: 313, col: 1, offset: 8401},
	expr: &actionExpr{
	pos: position{line: 313, col: 11, offset: 8413},
	run: (*parser).callonHttpRaw1,
	expr: &seqExpr{
	pos: position{line: 313, col: 11, offset: 8413},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 313, col: 11, offset: 8413},
	name: "Scheme",
},
&litMatcher{
	pos: position{line: 313, col: 18, offset: 8420},
	val: "://",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 313, col: 24, offset: 8426},
	name: "Authority",
},
&ruleRefExpr{
	pos: position{line: 313, col: 34, offset: 8436},
	name: "Path",
},
&zeroOrOneExpr{
	pos: position{line: 313, col: 39, offset: 8441},
	expr: &seqExpr{
	pos: position{line: 313, col: 41, offset: 8443},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 313, col: 41, offset: 8443},
	val: "?",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 313, col: 45, offset: 8447},
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
	pos: position{line: 315, col: 1, offset: 8504},
	expr: &seqExpr{
	pos: position{line: 315, col: 13, offset: 8518},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 315, col: 13, offset: 8518},
	expr: &seqExpr{
	pos: position{line: 315, col: 14, offset: 8519},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 315, col: 14, offset: 8519},
	name: "Userinfo",
},
&litMatcher{
	pos: position{line: 315, col: 23, offset: 8528},
	val: "@",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 315, col: 29, offset: 8534},
	name: "Host",
},
&zeroOrOneExpr{
	pos: position{line: 315, col: 34, offset: 8539},
	expr: &seqExpr{
	pos: position{line: 315, col: 35, offset: 8540},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 315, col: 35, offset: 8540},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 315, col: 39, offset: 8544},
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
	pos: position{line: 317, col: 1, offset: 8552},
	expr: &zeroOrMoreExpr{
	pos: position{line: 317, col: 12, offset: 8565},
	expr: &choiceExpr{
	pos: position{line: 317, col: 14, offset: 8567},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 317, col: 14, offset: 8567},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 317, col: 27, offset: 8580},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 317, col: 40, offset: 8593},
	name: "SubDelims",
},
&litMatcher{
	pos: position{line: 317, col: 52, offset: 8605},
	val: ":",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "Host",
	pos: position{line: 319, col: 1, offset: 8613},
	expr: &choiceExpr{
	pos: position{line: 319, col: 8, offset: 8622},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 319, col: 8, offset: 8622},
	name: "IPLiteral",
},
&ruleRefExpr{
	pos: position{line: 319, col: 20, offset: 8634},
	name: "RegName",
},
	},
},
},
{
	name: "Port",
	pos: position{line: 321, col: 1, offset: 8643},
	expr: &zeroOrMoreExpr{
	pos: position{line: 321, col: 8, offset: 8652},
	expr: &ruleRefExpr{
	pos: position{line: 321, col: 8, offset: 8652},
	name: "Digit",
},
},
},
{
	name: "IPLiteral",
	pos: position{line: 323, col: 1, offset: 8660},
	expr: &seqExpr{
	pos: position{line: 323, col: 13, offset: 8674},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 323, col: 13, offset: 8674},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 323, col: 17, offset: 8678},
	name: "IPv6address",
},
&litMatcher{
	pos: position{line: 323, col: 29, offset: 8690},
	val: "]",
	ignoreCase: false,
},
	},
},
},
{
	name: "IPv6address",
	pos: position{line: 325, col: 1, offset: 8695},
	expr: &actionExpr{
	pos: position{line: 325, col: 15, offset: 8711},
	run: (*parser).callonIPv6address1,
	expr: &seqExpr{
	pos: position{line: 325, col: 15, offset: 8711},
	exprs: []interface{}{
&zeroOrMoreExpr{
	pos: position{line: 325, col: 15, offset: 8711},
	expr: &ruleRefExpr{
	pos: position{line: 325, col: 16, offset: 8712},
	name: "HexDig",
},
},
&litMatcher{
	pos: position{line: 325, col: 25, offset: 8721},
	val: ":",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 325, col: 29, offset: 8725},
	expr: &choiceExpr{
	pos: position{line: 325, col: 30, offset: 8726},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 325, col: 30, offset: 8726},
	name: "HexDig",
},
&litMatcher{
	pos: position{line: 325, col: 39, offset: 8735},
	val: ":",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 325, col: 45, offset: 8741},
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
	pos: position{line: 331, col: 1, offset: 8895},
	expr: &zeroOrMoreExpr{
	pos: position{line: 331, col: 11, offset: 8907},
	expr: &choiceExpr{
	pos: position{line: 331, col: 12, offset: 8908},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 331, col: 12, offset: 8908},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 331, col: 25, offset: 8921},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 331, col: 38, offset: 8934},
	name: "SubDelims",
},
	},
},
},
},
{
	name: "PChar",
	pos: position{line: 333, col: 1, offset: 8947},
	expr: &choiceExpr{
	pos: position{line: 333, col: 9, offset: 8957},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 333, col: 9, offset: 8957},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 333, col: 22, offset: 8970},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 333, col: 35, offset: 8983},
	name: "SubDelims",
},
&charClassMatcher{
	pos: position{line: 333, col: 47, offset: 8995},
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
	pos: position{line: 335, col: 1, offset: 9001},
	expr: &zeroOrMoreExpr{
	pos: position{line: 335, col: 9, offset: 9011},
	expr: &choiceExpr{
	pos: position{line: 335, col: 10, offset: 9012},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 335, col: 10, offset: 9012},
	name: "PChar",
},
&charClassMatcher{
	pos: position{line: 335, col: 18, offset: 9020},
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
	pos: position{line: 337, col: 1, offset: 9028},
	expr: &seqExpr{
	pos: position{line: 337, col: 14, offset: 9043},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 337, col: 14, offset: 9043},
	val: "%",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 337, col: 18, offset: 9047},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 337, col: 25, offset: 9054},
	name: "HexDig",
},
	},
},
},
{
	name: "Unreserved",
	pos: position{line: 339, col: 1, offset: 9062},
	expr: &charClassMatcher{
	pos: position{line: 339, col: 14, offset: 9077},
	val: "[A-Za-z0-9._~-]",
	chars: []rune{'.','_','~','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SubDelims",
	pos: position{line: 341, col: 1, offset: 9094},
	expr: &choiceExpr{
	pos: position{line: 341, col: 13, offset: 9108},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 341, col: 13, offset: 9108},
	val: "!",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 341, col: 19, offset: 9114},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 341, col: 25, offset: 9120},
	val: "&",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 341, col: 31, offset: 9126},
	val: "'",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 341, col: 37, offset: 9132},
	val: "(",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 341, col: 43, offset: 9138},
	val: ")",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 341, col: 49, offset: 9144},
	val: "*",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 341, col: 55, offset: 9150},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 341, col: 61, offset: 9156},
	val: ",",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 341, col: 67, offset: 9162},
	val: ";",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 341, col: 73, offset: 9168},
	val: "=",
	ignoreCase: false,
},
	},
},
},
{
	name: "Http",
	pos: position{line: 343, col: 1, offset: 9173},
	expr: &actionExpr{
	pos: position{line: 343, col: 8, offset: 9182},
	run: (*parser).callonHttp1,
	expr: &labeledExpr{
	pos: position{line: 343, col: 8, offset: 9182},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 343, col: 10, offset: 9184},
	name: "HttpRaw",
},
},
},
},
{
	name: "Env",
	pos: position{line: 345, col: 1, offset: 9229},
	expr: &actionExpr{
	pos: position{line: 345, col: 7, offset: 9237},
	run: (*parser).callonEnv1,
	expr: &seqExpr{
	pos: position{line: 345, col: 7, offset: 9237},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 345, col: 7, offset: 9237},
	val: "env:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 345, col: 14, offset: 9244},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 345, col: 17, offset: 9247},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 345, col: 17, offset: 9247},
	name: "BashEnvironmentVariable",
},
&ruleRefExpr{
	pos: position{line: 345, col: 43, offset: 9273},
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
	pos: position{line: 347, col: 1, offset: 9318},
	expr: &actionExpr{
	pos: position{line: 347, col: 27, offset: 9346},
	run: (*parser).callonBashEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 347, col: 27, offset: 9346},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 347, col: 27, offset: 9346},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 347, col: 36, offset: 9355},
	expr: &charClassMatcher{
	pos: position{line: 347, col: 36, offset: 9355},
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
	pos: position{line: 351, col: 1, offset: 9411},
	expr: &actionExpr{
	pos: position{line: 351, col: 28, offset: 9440},
	run: (*parser).callonPosixEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 351, col: 28, offset: 9440},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 351, col: 28, offset: 9440},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 351, col: 32, offset: 9444},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 351, col: 34, offset: 9446},
	name: "PosixEnvironmentVariableContent",
},
},
&litMatcher{
	pos: position{line: 351, col: 66, offset: 9478},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PosixEnvironmentVariableContent",
	pos: position{line: 355, col: 1, offset: 9503},
	expr: &actionExpr{
	pos: position{line: 355, col: 35, offset: 9539},
	run: (*parser).callonPosixEnvironmentVariableContent1,
	expr: &labeledExpr{
	pos: position{line: 355, col: 35, offset: 9539},
	label: "v",
	expr: &oneOrMoreExpr{
	pos: position{line: 355, col: 37, offset: 9541},
	expr: &ruleRefExpr{
	pos: position{line: 355, col: 37, offset: 9541},
	name: "PosixEnvironmentVariableCharacter",
},
},
},
},
},
{
	name: "PosixEnvironmentVariableCharacter",
	pos: position{line: 364, col: 1, offset: 9754},
	expr: &choiceExpr{
	pos: position{line: 365, col: 7, offset: 9798},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 365, col: 7, offset: 9798},
	run: (*parser).callonPosixEnvironmentVariableCharacter2,
	expr: &litMatcher{
	pos: position{line: 365, col: 7, offset: 9798},
	val: "\\\"",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 366, col: 7, offset: 9838},
	run: (*parser).callonPosixEnvironmentVariableCharacter4,
	expr: &litMatcher{
	pos: position{line: 366, col: 7, offset: 9838},
	val: "\\\\",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 367, col: 7, offset: 9878},
	run: (*parser).callonPosixEnvironmentVariableCharacter6,
	expr: &litMatcher{
	pos: position{line: 367, col: 7, offset: 9878},
	val: "\\a",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 368, col: 7, offset: 9918},
	run: (*parser).callonPosixEnvironmentVariableCharacter8,
	expr: &litMatcher{
	pos: position{line: 368, col: 7, offset: 9918},
	val: "\\b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 369, col: 7, offset: 9958},
	run: (*parser).callonPosixEnvironmentVariableCharacter10,
	expr: &litMatcher{
	pos: position{line: 369, col: 7, offset: 9958},
	val: "\\f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 370, col: 7, offset: 9998},
	run: (*parser).callonPosixEnvironmentVariableCharacter12,
	expr: &litMatcher{
	pos: position{line: 370, col: 7, offset: 9998},
	val: "\\n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 371, col: 7, offset: 10038},
	run: (*parser).callonPosixEnvironmentVariableCharacter14,
	expr: &litMatcher{
	pos: position{line: 371, col: 7, offset: 10038},
	val: "\\r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 372, col: 7, offset: 10078},
	run: (*parser).callonPosixEnvironmentVariableCharacter16,
	expr: &litMatcher{
	pos: position{line: 372, col: 7, offset: 10078},
	val: "\\t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 373, col: 7, offset: 10118},
	run: (*parser).callonPosixEnvironmentVariableCharacter18,
	expr: &litMatcher{
	pos: position{line: 373, col: 7, offset: 10118},
	val: "\\v",
	ignoreCase: false,
},
},
&charClassMatcher{
	pos: position{line: 374, col: 7, offset: 10158},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 375, col: 7, offset: 10176},
	val: "[\\x23-\\x3c]",
	ranges: []rune{'#','<',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 376, col: 7, offset: 10194},
	val: "[\\x3e-\\x5b]",
	ranges: []rune{'>','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 377, col: 7, offset: 10212},
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
	pos: position{line: 379, col: 1, offset: 10225},
	expr: &choiceExpr{
	pos: position{line: 379, col: 14, offset: 10240},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 379, col: 14, offset: 10240},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 379, col: 24, offset: 10250},
	name: "Local",
},
&ruleRefExpr{
	pos: position{line: 379, col: 32, offset: 10258},
	name: "Http",
},
&ruleRefExpr{
	pos: position{line: 379, col: 39, offset: 10265},
	name: "Env",
},
	},
},
},
{
	name: "HashValue",
	pos: position{line: 382, col: 1, offset: 10338},
	expr: &actionExpr{
	pos: position{line: 382, col: 13, offset: 10350},
	run: (*parser).callonHashValue1,
	expr: &seqExpr{
	pos: position{line: 382, col: 13, offset: 10350},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 382, col: 13, offset: 10350},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 382, col: 20, offset: 10357},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 382, col: 27, offset: 10364},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 382, col: 34, offset: 10371},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 382, col: 41, offset: 10378},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 382, col: 48, offset: 10385},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 382, col: 55, offset: 10392},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 382, col: 62, offset: 10399},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 383, col: 13, offset: 10418},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 383, col: 20, offset: 10425},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 383, col: 27, offset: 10432},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 383, col: 34, offset: 10439},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 383, col: 41, offset: 10446},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 383, col: 48, offset: 10453},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 383, col: 55, offset: 10460},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 383, col: 62, offset: 10467},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 384, col: 13, offset: 10486},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 384, col: 20, offset: 10493},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 384, col: 27, offset: 10500},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 384, col: 34, offset: 10507},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 384, col: 41, offset: 10514},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 384, col: 48, offset: 10521},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 384, col: 55, offset: 10528},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 384, col: 62, offset: 10535},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 385, col: 13, offset: 10554},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 385, col: 20, offset: 10561},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 385, col: 27, offset: 10568},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 385, col: 34, offset: 10575},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 385, col: 41, offset: 10582},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 385, col: 48, offset: 10589},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 385, col: 55, offset: 10596},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 385, col: 62, offset: 10603},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 386, col: 13, offset: 10622},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 386, col: 20, offset: 10629},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 386, col: 27, offset: 10636},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 386, col: 34, offset: 10643},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 386, col: 41, offset: 10650},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 386, col: 48, offset: 10657},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 386, col: 55, offset: 10664},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 386, col: 62, offset: 10671},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 387, col: 13, offset: 10690},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 387, col: 20, offset: 10697},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 387, col: 27, offset: 10704},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 387, col: 34, offset: 10711},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 387, col: 41, offset: 10718},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 387, col: 48, offset: 10725},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 387, col: 55, offset: 10732},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 387, col: 62, offset: 10739},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 388, col: 13, offset: 10758},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 388, col: 20, offset: 10765},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 388, col: 27, offset: 10772},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 388, col: 34, offset: 10779},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 388, col: 41, offset: 10786},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 388, col: 48, offset: 10793},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 388, col: 55, offset: 10800},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 388, col: 62, offset: 10807},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 389, col: 13, offset: 10826},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 389, col: 20, offset: 10833},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 389, col: 27, offset: 10840},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 389, col: 34, offset: 10847},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 389, col: 41, offset: 10854},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 389, col: 48, offset: 10861},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 389, col: 55, offset: 10868},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 389, col: 62, offset: 10875},
	name: "HexDig",
},
	},
},
},
},
{
	name: "Hash",
	pos: position{line: 395, col: 1, offset: 11019},
	expr: &actionExpr{
	pos: position{line: 395, col: 8, offset: 11026},
	run: (*parser).callonHash1,
	expr: &seqExpr{
	pos: position{line: 395, col: 8, offset: 11026},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 395, col: 8, offset: 11026},
	val: "sha256:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 395, col: 18, offset: 11036},
	label: "val",
	expr: &ruleRefExpr{
	pos: position{line: 395, col: 22, offset: 11040},
	name: "HashValue",
},
},
	},
},
},
},
{
	name: "ImportHashed",
	pos: position{line: 397, col: 1, offset: 11110},
	expr: &actionExpr{
	pos: position{line: 397, col: 16, offset: 11127},
	run: (*parser).callonImportHashed1,
	expr: &seqExpr{
	pos: position{line: 397, col: 16, offset: 11127},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 397, col: 16, offset: 11127},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 397, col: 18, offset: 11129},
	name: "ImportType",
},
},
&labeledExpr{
	pos: position{line: 397, col: 29, offset: 11140},
	label: "h",
	expr: &zeroOrOneExpr{
	pos: position{line: 397, col: 31, offset: 11142},
	expr: &seqExpr{
	pos: position{line: 397, col: 32, offset: 11143},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 397, col: 32, offset: 11143},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 397, col: 35, offset: 11146},
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
	pos: position{line: 405, col: 1, offset: 11301},
	expr: &choiceExpr{
	pos: position{line: 405, col: 10, offset: 11312},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 405, col: 10, offset: 11312},
	run: (*parser).callonImport2,
	expr: &seqExpr{
	pos: position{line: 405, col: 10, offset: 11312},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 405, col: 10, offset: 11312},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 405, col: 12, offset: 11314},
	name: "ImportHashed",
},
},
&ruleRefExpr{
	pos: position{line: 405, col: 25, offset: 11327},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 405, col: 27, offset: 11329},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 405, col: 30, offset: 11332},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 405, col: 33, offset: 11335},
	name: "Text",
},
	},
},
},
&actionExpr{
	pos: position{line: 406, col: 10, offset: 11432},
	run: (*parser).callonImport10,
	expr: &labeledExpr{
	pos: position{line: 406, col: 10, offset: 11432},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 406, col: 12, offset: 11434},
	name: "ImportHashed",
},
},
},
	},
},
},
{
	name: "LetBinding",
	pos: position{line: 409, col: 1, offset: 11529},
	expr: &actionExpr{
	pos: position{line: 409, col: 14, offset: 11544},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 409, col: 14, offset: 11544},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 409, col: 14, offset: 11544},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 409, col: 18, offset: 11548},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 409, col: 21, offset: 11551},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 409, col: 27, offset: 11557},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 409, col: 44, offset: 11574},
	name: "_",
},
&labeledExpr{
	pos: position{line: 409, col: 46, offset: 11576},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 409, col: 48, offset: 11578},
	expr: &seqExpr{
	pos: position{line: 409, col: 49, offset: 11579},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 409, col: 49, offset: 11579},
	name: "Annotation",
},
&ruleRefExpr{
	pos: position{line: 409, col: 60, offset: 11590},
	name: "_",
},
	},
},
},
},
&litMatcher{
	pos: position{line: 410, col: 13, offset: 11606},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 410, col: 17, offset: 11610},
	name: "_",
},
&labeledExpr{
	pos: position{line: 410, col: 19, offset: 11612},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 410, col: 21, offset: 11614},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 410, col: 32, offset: 11625},
	name: "_",
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 425, col: 1, offset: 11934},
	expr: &choiceExpr{
	pos: position{line: 426, col: 7, offset: 11955},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 426, col: 7, offset: 11955},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 426, col: 7, offset: 11955},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 426, col: 7, offset: 11955},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 426, col: 14, offset: 11962},
	name: "_",
},
&litMatcher{
	pos: position{line: 426, col: 16, offset: 11964},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 426, col: 20, offset: 11968},
	name: "_",
},
&labeledExpr{
	pos: position{line: 426, col: 22, offset: 11970},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 426, col: 28, offset: 11976},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 426, col: 45, offset: 11993},
	name: "_",
},
&litMatcher{
	pos: position{line: 426, col: 47, offset: 11995},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 426, col: 51, offset: 11999},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 426, col: 54, offset: 12002},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 426, col: 56, offset: 12004},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 426, col: 67, offset: 12015},
	name: "_",
},
&litMatcher{
	pos: position{line: 426, col: 69, offset: 12017},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 426, col: 73, offset: 12021},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 426, col: 75, offset: 12023},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 426, col: 81, offset: 12029},
	name: "_",
},
&labeledExpr{
	pos: position{line: 426, col: 83, offset: 12031},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 426, col: 88, offset: 12036},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 429, col: 7, offset: 12153},
	run: (*parser).callonExpression22,
	expr: &seqExpr{
	pos: position{line: 429, col: 7, offset: 12153},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 429, col: 7, offset: 12153},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 429, col: 10, offset: 12156},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 429, col: 13, offset: 12159},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 429, col: 18, offset: 12164},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 429, col: 29, offset: 12175},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 429, col: 31, offset: 12177},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 429, col: 36, offset: 12182},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 429, col: 39, offset: 12185},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 429, col: 41, offset: 12187},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 429, col: 52, offset: 12198},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 429, col: 54, offset: 12200},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 429, col: 59, offset: 12205},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 429, col: 62, offset: 12208},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 429, col: 64, offset: 12210},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 432, col: 7, offset: 12296},
	run: (*parser).callonExpression38,
	expr: &seqExpr{
	pos: position{line: 432, col: 7, offset: 12296},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 432, col: 7, offset: 12296},
	label: "bindings",
	expr: &oneOrMoreExpr{
	pos: position{line: 432, col: 16, offset: 12305},
	expr: &ruleRefExpr{
	pos: position{line: 432, col: 16, offset: 12305},
	name: "LetBinding",
},
},
},
&ruleRefExpr{
	pos: position{line: 432, col: 28, offset: 12317},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 432, col: 31, offset: 12320},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 432, col: 34, offset: 12323},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 432, col: 36, offset: 12325},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 439, col: 7, offset: 12565},
	run: (*parser).callonExpression47,
	expr: &seqExpr{
	pos: position{line: 439, col: 7, offset: 12565},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 439, col: 7, offset: 12565},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 439, col: 14, offset: 12572},
	name: "_",
},
&litMatcher{
	pos: position{line: 439, col: 16, offset: 12574},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 439, col: 20, offset: 12578},
	name: "_",
},
&labeledExpr{
	pos: position{line: 439, col: 22, offset: 12580},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 439, col: 28, offset: 12586},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 439, col: 45, offset: 12603},
	name: "_",
},
&litMatcher{
	pos: position{line: 439, col: 47, offset: 12605},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 439, col: 51, offset: 12609},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 439, col: 54, offset: 12612},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 439, col: 56, offset: 12614},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 439, col: 67, offset: 12625},
	name: "_",
},
&litMatcher{
	pos: position{line: 439, col: 69, offset: 12627},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 439, col: 73, offset: 12631},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 439, col: 75, offset: 12633},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 439, col: 81, offset: 12639},
	name: "_",
},
&labeledExpr{
	pos: position{line: 439, col: 83, offset: 12641},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 439, col: 88, offset: 12646},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 442, col: 7, offset: 12755},
	run: (*parser).callonExpression67,
	expr: &seqExpr{
	pos: position{line: 442, col: 7, offset: 12755},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 442, col: 7, offset: 12755},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 442, col: 9, offset: 12757},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 442, col: 28, offset: 12776},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 442, col: 30, offset: 12778},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 442, col: 36, offset: 12784},
	name: "_",
},
&labeledExpr{
	pos: position{line: 442, col: 38, offset: 12786},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 442, col: 40, offset: 12788},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 443, col: 7, offset: 12847},
	run: (*parser).callonExpression76,
	expr: &seqExpr{
	pos: position{line: 443, col: 7, offset: 12847},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 443, col: 7, offset: 12847},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 443, col: 13, offset: 12853},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 443, col: 16, offset: 12856},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 443, col: 18, offset: 12858},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 443, col: 35, offset: 12875},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 443, col: 38, offset: 12878},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 443, col: 40, offset: 12880},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 443, col: 57, offset: 12897},
	name: "_",
},
&litMatcher{
	pos: position{line: 443, col: 59, offset: 12899},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 443, col: 63, offset: 12903},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 443, col: 66, offset: 12906},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 443, col: 68, offset: 12908},
	name: "ApplicationExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 446, col: 7, offset: 13029},
	name: "EmptyList",
},
&ruleRefExpr{
	pos: position{line: 447, col: 7, offset: 13045},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 449, col: 1, offset: 13066},
	expr: &actionExpr{
	pos: position{line: 449, col: 14, offset: 13081},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 449, col: 14, offset: 13081},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 449, col: 14, offset: 13081},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 449, col: 18, offset: 13085},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 449, col: 21, offset: 13088},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 449, col: 23, offset: 13090},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 451, col: 1, offset: 13120},
	expr: &actionExpr{
	pos: position{line: 452, col: 1, offset: 13144},
	run: (*parser).callonAnnotatedExpression1,
	expr: &seqExpr{
	pos: position{line: 452, col: 1, offset: 13144},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 452, col: 1, offset: 13144},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 452, col: 3, offset: 13146},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 452, col: 22, offset: 13165},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 452, col: 24, offset: 13167},
	expr: &seqExpr{
	pos: position{line: 452, col: 25, offset: 13168},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 452, col: 25, offset: 13168},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 452, col: 27, offset: 13170},
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
	pos: position{line: 457, col: 1, offset: 13295},
	expr: &actionExpr{
	pos: position{line: 457, col: 13, offset: 13309},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 457, col: 13, offset: 13309},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 457, col: 13, offset: 13309},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 457, col: 17, offset: 13313},
	name: "_",
},
&litMatcher{
	pos: position{line: 457, col: 19, offset: 13315},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 457, col: 23, offset: 13319},
	name: "_",
},
&litMatcher{
	pos: position{line: 457, col: 25, offset: 13321},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 457, col: 29, offset: 13325},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 457, col: 32, offset: 13328},
	name: "List",
},
&ruleRefExpr{
	pos: position{line: 457, col: 37, offset: 13333},
	name: "_",
},
&labeledExpr{
	pos: position{line: 457, col: 39, offset: 13335},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 457, col: 41, offset: 13337},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 461, col: 1, offset: 13400},
	expr: &ruleRefExpr{
	pos: position{line: 461, col: 22, offset: 13423},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 463, col: 1, offset: 13444},
	expr: &actionExpr{
	pos: position{line: 463, col: 26, offset: 13471},
	run: (*parser).callonImportAltExpression1,
	expr: &seqExpr{
	pos: position{line: 463, col: 26, offset: 13471},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 463, col: 26, offset: 13471},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 463, col: 32, offset: 13477},
	name: "OrExpression",
},
},
&labeledExpr{
	pos: position{line: 463, col: 55, offset: 13500},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 463, col: 60, offset: 13505},
	expr: &seqExpr{
	pos: position{line: 463, col: 61, offset: 13506},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 463, col: 61, offset: 13506},
	name: "_",
},
&litMatcher{
	pos: position{line: 463, col: 63, offset: 13508},
	val: "?",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 463, col: 67, offset: 13512},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 463, col: 69, offset: 13514},
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
	pos: position{line: 465, col: 1, offset: 13585},
	expr: &actionExpr{
	pos: position{line: 465, col: 26, offset: 13612},
	run: (*parser).callonOrExpression1,
	expr: &seqExpr{
	pos: position{line: 465, col: 26, offset: 13612},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 465, col: 26, offset: 13612},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 465, col: 32, offset: 13618},
	name: "PlusExpression",
},
},
&labeledExpr{
	pos: position{line: 465, col: 55, offset: 13641},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 465, col: 60, offset: 13646},
	expr: &seqExpr{
	pos: position{line: 465, col: 61, offset: 13647},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 465, col: 61, offset: 13647},
	name: "_",
},
&litMatcher{
	pos: position{line: 465, col: 63, offset: 13649},
	val: "||",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 465, col: 68, offset: 13654},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 465, col: 70, offset: 13656},
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
	pos: position{line: 467, col: 1, offset: 13722},
	expr: &actionExpr{
	pos: position{line: 467, col: 26, offset: 13749},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 467, col: 26, offset: 13749},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 467, col: 26, offset: 13749},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 467, col: 32, offset: 13755},
	name: "TextAppendExpression",
},
},
&labeledExpr{
	pos: position{line: 467, col: 55, offset: 13778},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 467, col: 60, offset: 13783},
	expr: &seqExpr{
	pos: position{line: 467, col: 61, offset: 13784},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 467, col: 61, offset: 13784},
	name: "_",
},
&litMatcher{
	pos: position{line: 467, col: 63, offset: 13786},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 467, col: 67, offset: 13790},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 467, col: 70, offset: 13793},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 467, col: 72, offset: 13795},
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
	pos: position{line: 469, col: 1, offset: 13869},
	expr: &actionExpr{
	pos: position{line: 469, col: 26, offset: 13896},
	run: (*parser).callonTextAppendExpression1,
	expr: &seqExpr{
	pos: position{line: 469, col: 26, offset: 13896},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 469, col: 26, offset: 13896},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 469, col: 32, offset: 13902},
	name: "ListAppendExpression",
},
},
&labeledExpr{
	pos: position{line: 469, col: 55, offset: 13925},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 469, col: 60, offset: 13930},
	expr: &seqExpr{
	pos: position{line: 469, col: 61, offset: 13931},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 469, col: 61, offset: 13931},
	name: "_",
},
&litMatcher{
	pos: position{line: 469, col: 63, offset: 13933},
	val: "++",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 469, col: 68, offset: 13938},
	name: "_",
},
&labeledExpr{
	pos: position{line: 469, col: 70, offset: 13940},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 469, col: 72, offset: 13942},
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
	pos: position{line: 471, col: 1, offset: 14022},
	expr: &actionExpr{
	pos: position{line: 471, col: 26, offset: 14049},
	run: (*parser).callonListAppendExpression1,
	expr: &seqExpr{
	pos: position{line: 471, col: 26, offset: 14049},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 471, col: 26, offset: 14049},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 471, col: 32, offset: 14055},
	name: "AndExpression",
},
},
&labeledExpr{
	pos: position{line: 471, col: 55, offset: 14078},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 471, col: 60, offset: 14083},
	expr: &seqExpr{
	pos: position{line: 471, col: 61, offset: 14084},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 471, col: 61, offset: 14084},
	name: "_",
},
&litMatcher{
	pos: position{line: 471, col: 63, offset: 14086},
	val: "#",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 471, col: 67, offset: 14090},
	name: "_",
},
&labeledExpr{
	pos: position{line: 471, col: 69, offset: 14092},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 471, col: 71, offset: 14094},
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
	pos: position{line: 473, col: 1, offset: 14167},
	expr: &actionExpr{
	pos: position{line: 473, col: 26, offset: 14194},
	run: (*parser).callonAndExpression1,
	expr: &seqExpr{
	pos: position{line: 473, col: 26, offset: 14194},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 473, col: 26, offset: 14194},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 473, col: 32, offset: 14200},
	name: "CombineExpression",
},
},
&labeledExpr{
	pos: position{line: 473, col: 55, offset: 14223},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 473, col: 60, offset: 14228},
	expr: &seqExpr{
	pos: position{line: 473, col: 61, offset: 14229},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 473, col: 61, offset: 14229},
	name: "_",
},
&litMatcher{
	pos: position{line: 473, col: 63, offset: 14231},
	val: "&&",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 473, col: 68, offset: 14236},
	name: "_",
},
&labeledExpr{
	pos: position{line: 473, col: 70, offset: 14238},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 473, col: 72, offset: 14240},
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
	pos: position{line: 475, col: 1, offset: 14310},
	expr: &actionExpr{
	pos: position{line: 475, col: 26, offset: 14337},
	run: (*parser).callonCombineExpression1,
	expr: &seqExpr{
	pos: position{line: 475, col: 26, offset: 14337},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 475, col: 26, offset: 14337},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 475, col: 32, offset: 14343},
	name: "PreferExpression",
},
},
&labeledExpr{
	pos: position{line: 475, col: 55, offset: 14366},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 475, col: 60, offset: 14371},
	expr: &seqExpr{
	pos: position{line: 475, col: 61, offset: 14372},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 475, col: 61, offset: 14372},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 475, col: 63, offset: 14374},
	name: "Combine",
},
&ruleRefExpr{
	pos: position{line: 475, col: 71, offset: 14382},
	name: "_",
},
&labeledExpr{
	pos: position{line: 475, col: 73, offset: 14384},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 475, col: 75, offset: 14386},
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
	pos: position{line: 477, col: 1, offset: 14463},
	expr: &actionExpr{
	pos: position{line: 477, col: 26, offset: 14490},
	run: (*parser).callonPreferExpression1,
	expr: &seqExpr{
	pos: position{line: 477, col: 26, offset: 14490},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 477, col: 26, offset: 14490},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 477, col: 32, offset: 14496},
	name: "CombineTypesExpression",
},
},
&labeledExpr{
	pos: position{line: 477, col: 55, offset: 14519},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 477, col: 60, offset: 14524},
	expr: &seqExpr{
	pos: position{line: 477, col: 61, offset: 14525},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 477, col: 61, offset: 14525},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 477, col: 63, offset: 14527},
	name: "Prefer",
},
&ruleRefExpr{
	pos: position{line: 477, col: 70, offset: 14534},
	name: "_",
},
&labeledExpr{
	pos: position{line: 477, col: 72, offset: 14536},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 477, col: 74, offset: 14538},
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
	pos: position{line: 479, col: 1, offset: 14632},
	expr: &actionExpr{
	pos: position{line: 479, col: 26, offset: 14659},
	run: (*parser).callonCombineTypesExpression1,
	expr: &seqExpr{
	pos: position{line: 479, col: 26, offset: 14659},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 479, col: 26, offset: 14659},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 479, col: 32, offset: 14665},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 479, col: 55, offset: 14688},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 479, col: 60, offset: 14693},
	expr: &seqExpr{
	pos: position{line: 479, col: 61, offset: 14694},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 479, col: 61, offset: 14694},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 479, col: 63, offset: 14696},
	name: "CombineTypes",
},
&ruleRefExpr{
	pos: position{line: 479, col: 76, offset: 14709},
	name: "_",
},
&labeledExpr{
	pos: position{line: 479, col: 78, offset: 14711},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 479, col: 80, offset: 14713},
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
	pos: position{line: 481, col: 1, offset: 14793},
	expr: &actionExpr{
	pos: position{line: 481, col: 26, offset: 14820},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 481, col: 26, offset: 14820},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 481, col: 26, offset: 14820},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 481, col: 32, offset: 14826},
	name: "EqualExpression",
},
},
&labeledExpr{
	pos: position{line: 481, col: 55, offset: 14849},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 481, col: 60, offset: 14854},
	expr: &seqExpr{
	pos: position{line: 481, col: 61, offset: 14855},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 481, col: 61, offset: 14855},
	name: "_",
},
&litMatcher{
	pos: position{line: 481, col: 63, offset: 14857},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 481, col: 67, offset: 14861},
	name: "_",
},
&labeledExpr{
	pos: position{line: 481, col: 69, offset: 14863},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 481, col: 71, offset: 14865},
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
	pos: position{line: 483, col: 1, offset: 14935},
	expr: &actionExpr{
	pos: position{line: 483, col: 26, offset: 14962},
	run: (*parser).callonEqualExpression1,
	expr: &seqExpr{
	pos: position{line: 483, col: 26, offset: 14962},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 483, col: 26, offset: 14962},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 483, col: 32, offset: 14968},
	name: "NotEqualExpression",
},
},
&labeledExpr{
	pos: position{line: 483, col: 55, offset: 14991},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 483, col: 60, offset: 14996},
	expr: &seqExpr{
	pos: position{line: 483, col: 61, offset: 14997},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 483, col: 61, offset: 14997},
	name: "_",
},
&litMatcher{
	pos: position{line: 483, col: 63, offset: 14999},
	val: "==",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 483, col: 68, offset: 15004},
	name: "_",
},
&labeledExpr{
	pos: position{line: 483, col: 70, offset: 15006},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 483, col: 72, offset: 15008},
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
	pos: position{line: 485, col: 1, offset: 15078},
	expr: &actionExpr{
	pos: position{line: 485, col: 26, offset: 15105},
	run: (*parser).callonNotEqualExpression1,
	expr: &seqExpr{
	pos: position{line: 485, col: 26, offset: 15105},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 485, col: 26, offset: 15105},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 485, col: 32, offset: 15111},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 485, col: 55, offset: 15134},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 485, col: 60, offset: 15139},
	expr: &seqExpr{
	pos: position{line: 485, col: 61, offset: 15140},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 485, col: 61, offset: 15140},
	name: "_",
},
&litMatcher{
	pos: position{line: 485, col: 63, offset: 15142},
	val: "!=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 485, col: 68, offset: 15147},
	name: "_",
},
&labeledExpr{
	pos: position{line: 485, col: 70, offset: 15149},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 485, col: 72, offset: 15151},
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
	pos: position{line: 488, col: 1, offset: 15225},
	expr: &actionExpr{
	pos: position{line: 488, col: 25, offset: 15251},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 488, col: 25, offset: 15251},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 488, col: 25, offset: 15251},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 488, col: 27, offset: 15253},
	name: "FirstApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 488, col: 54, offset: 15280},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 488, col: 59, offset: 15285},
	expr: &seqExpr{
	pos: position{line: 488, col: 60, offset: 15286},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 488, col: 60, offset: 15286},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 488, col: 63, offset: 15289},
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
	pos: position{line: 497, col: 1, offset: 15532},
	expr: &choiceExpr{
	pos: position{line: 498, col: 8, offset: 15570},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 498, col: 8, offset: 15570},
	run: (*parser).callonFirstApplicationExpression2,
	expr: &seqExpr{
	pos: position{line: 498, col: 8, offset: 15570},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 498, col: 8, offset: 15570},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 498, col: 14, offset: 15576},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 498, col: 17, offset: 15579},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 498, col: 19, offset: 15581},
	name: "ImportExpression",
},
},
&ruleRefExpr{
	pos: position{line: 498, col: 36, offset: 15598},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 498, col: 39, offset: 15601},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 498, col: 41, offset: 15603},
	name: "ImportExpression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 501, col: 8, offset: 15706},
	run: (*parser).callonFirstApplicationExpression11,
	expr: &seqExpr{
	pos: position{line: 501, col: 8, offset: 15706},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 501, col: 8, offset: 15706},
	name: "Some",
},
&ruleRefExpr{
	pos: position{line: 501, col: 13, offset: 15711},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 501, col: 16, offset: 15714},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 501, col: 18, offset: 15716},
	name: "ImportExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 502, col: 8, offset: 15771},
	name: "ImportExpression",
},
	},
},
},
{
	name: "ImportExpression",
	pos: position{line: 504, col: 1, offset: 15789},
	expr: &choiceExpr{
	pos: position{line: 504, col: 20, offset: 15810},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 504, col: 20, offset: 15810},
	name: "Import",
},
&ruleRefExpr{
	pos: position{line: 504, col: 29, offset: 15819},
	name: "SelectorExpression",
},
	},
},
},
{
	name: "SelectorExpression",
	pos: position{line: 506, col: 1, offset: 15839},
	expr: &actionExpr{
	pos: position{line: 506, col: 22, offset: 15862},
	run: (*parser).callonSelectorExpression1,
	expr: &seqExpr{
	pos: position{line: 506, col: 22, offset: 15862},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 506, col: 22, offset: 15862},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 506, col: 24, offset: 15864},
	name: "PrimitiveExpression",
},
},
&labeledExpr{
	pos: position{line: 506, col: 44, offset: 15884},
	label: "ls",
	expr: &zeroOrMoreExpr{
	pos: position{line: 506, col: 47, offset: 15887},
	expr: &seqExpr{
	pos: position{line: 506, col: 48, offset: 15888},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 506, col: 48, offset: 15888},
	name: "_",
},
&litMatcher{
	pos: position{line: 506, col: 50, offset: 15890},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 506, col: 54, offset: 15894},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 506, col: 56, offset: 15896},
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
	pos: position{line: 516, col: 1, offset: 16129},
	expr: &choiceExpr{
	pos: position{line: 517, col: 7, offset: 16159},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 517, col: 7, offset: 16159},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 518, col: 7, offset: 16179},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 519, col: 7, offset: 16200},
	name: "IntegerLiteral",
},
&ruleRefExpr{
	pos: position{line: 520, col: 7, offset: 16221},
	name: "TextLiteral",
},
&actionExpr{
	pos: position{line: 521, col: 7, offset: 16239},
	run: (*parser).callonPrimitiveExpression6,
	expr: &seqExpr{
	pos: position{line: 521, col: 7, offset: 16239},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 521, col: 7, offset: 16239},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 521, col: 11, offset: 16243},
	name: "_",
},
&labeledExpr{
	pos: position{line: 521, col: 13, offset: 16245},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 521, col: 15, offset: 16247},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 521, col: 35, offset: 16267},
	name: "_",
},
&litMatcher{
	pos: position{line: 521, col: 37, offset: 16269},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&actionExpr{
	pos: position{line: 522, col: 7, offset: 16297},
	run: (*parser).callonPrimitiveExpression14,
	expr: &seqExpr{
	pos: position{line: 522, col: 7, offset: 16297},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 522, col: 7, offset: 16297},
	val: "<",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 522, col: 11, offset: 16301},
	name: "_",
},
&labeledExpr{
	pos: position{line: 522, col: 13, offset: 16303},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 522, col: 15, offset: 16305},
	name: "UnionType",
},
},
&ruleRefExpr{
	pos: position{line: 522, col: 25, offset: 16315},
	name: "_",
},
&litMatcher{
	pos: position{line: 522, col: 27, offset: 16317},
	val: ">",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 523, col: 7, offset: 16345},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 524, col: 7, offset: 16371},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 525, col: 7, offset: 16388},
	run: (*parser).callonPrimitiveExpression24,
	expr: &seqExpr{
	pos: position{line: 525, col: 7, offset: 16388},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 525, col: 7, offset: 16388},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 525, col: 11, offset: 16392},
	name: "_",
},
&labeledExpr{
	pos: position{line: 525, col: 14, offset: 16395},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 525, col: 16, offset: 16397},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 525, col: 27, offset: 16408},
	name: "_",
},
&litMatcher{
	pos: position{line: 525, col: 29, offset: 16410},
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
	pos: position{line: 527, col: 1, offset: 16433},
	expr: &choiceExpr{
	pos: position{line: 528, col: 7, offset: 16463},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 528, col: 7, offset: 16463},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &litMatcher{
	pos: position{line: 528, col: 7, offset: 16463},
	val: "=",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 529, col: 7, offset: 16518},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 530, col: 7, offset: 16543},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 531, col: 7, offset: 16571},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 531, col: 7, offset: 16571},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 533, col: 1, offset: 16617},
	expr: &actionExpr{
	pos: position{line: 533, col: 19, offset: 16637},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 533, col: 19, offset: 16637},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 533, col: 19, offset: 16637},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 533, col: 24, offset: 16642},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 533, col: 33, offset: 16651},
	name: "_",
},
&litMatcher{
	pos: position{line: 533, col: 35, offset: 16653},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 533, col: 39, offset: 16657},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 533, col: 42, offset: 16660},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 533, col: 47, offset: 16665},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 536, col: 1, offset: 16722},
	expr: &actionExpr{
	pos: position{line: 536, col: 18, offset: 16741},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 536, col: 18, offset: 16741},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 536, col: 18, offset: 16741},
	name: "_",
},
&litMatcher{
	pos: position{line: 536, col: 20, offset: 16743},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 536, col: 24, offset: 16747},
	name: "_",
},
&labeledExpr{
	pos: position{line: 536, col: 26, offset: 16749},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 536, col: 28, offset: 16751},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 537, col: 1, offset: 16783},
	expr: &actionExpr{
	pos: position{line: 538, col: 7, offset: 16812},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 538, col: 7, offset: 16812},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 538, col: 7, offset: 16812},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 538, col: 13, offset: 16818},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 538, col: 29, offset: 16834},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 538, col: 34, offset: 16839},
	expr: &ruleRefExpr{
	pos: position{line: 538, col: 34, offset: 16839},
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
	pos: position{line: 552, col: 1, offset: 17423},
	expr: &actionExpr{
	pos: position{line: 552, col: 22, offset: 17446},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 552, col: 22, offset: 17446},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 552, col: 22, offset: 17446},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 552, col: 27, offset: 17451},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 552, col: 36, offset: 17460},
	name: "_",
},
&litMatcher{
	pos: position{line: 552, col: 38, offset: 17462},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 552, col: 42, offset: 17466},
	name: "_",
},
&labeledExpr{
	pos: position{line: 552, col: 44, offset: 17468},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 552, col: 49, offset: 17473},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 555, col: 1, offset: 17530},
	expr: &actionExpr{
	pos: position{line: 555, col: 21, offset: 17552},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 555, col: 21, offset: 17552},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 555, col: 21, offset: 17552},
	name: "_",
},
&litMatcher{
	pos: position{line: 555, col: 23, offset: 17554},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 555, col: 27, offset: 17558},
	name: "_",
},
&labeledExpr{
	pos: position{line: 555, col: 29, offset: 17560},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 555, col: 31, offset: 17562},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 556, col: 1, offset: 17597},
	expr: &actionExpr{
	pos: position{line: 557, col: 7, offset: 17629},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 557, col: 7, offset: 17629},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 557, col: 7, offset: 17629},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 557, col: 13, offset: 17635},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 557, col: 32, offset: 17654},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 557, col: 37, offset: 17659},
	expr: &ruleRefExpr{
	pos: position{line: 557, col: 37, offset: 17659},
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
	pos: position{line: 571, col: 1, offset: 18249},
	expr: &choiceExpr{
	pos: position{line: 571, col: 13, offset: 18263},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 571, col: 13, offset: 18263},
	name: "NonEmptyUnionType",
},
&ruleRefExpr{
	pos: position{line: 571, col: 33, offset: 18283},
	name: "EmptyUnionType",
},
	},
},
},
{
	name: "EmptyUnionType",
	pos: position{line: 573, col: 1, offset: 18299},
	expr: &actionExpr{
	pos: position{line: 573, col: 18, offset: 18318},
	run: (*parser).callonEmptyUnionType1,
	expr: &litMatcher{
	pos: position{line: 573, col: 18, offset: 18318},
	val: "",
	ignoreCase: false,
},
},
},
{
	name: "NonEmptyUnionType",
	pos: position{line: 575, col: 1, offset: 18350},
	expr: &actionExpr{
	pos: position{line: 575, col: 21, offset: 18372},
	run: (*parser).callonNonEmptyUnionType1,
	expr: &seqExpr{
	pos: position{line: 575, col: 21, offset: 18372},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 575, col: 21, offset: 18372},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 575, col: 27, offset: 18378},
	name: "UnionVariant",
},
},
&labeledExpr{
	pos: position{line: 575, col: 40, offset: 18391},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 575, col: 45, offset: 18396},
	expr: &seqExpr{
	pos: position{line: 575, col: 46, offset: 18397},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 575, col: 46, offset: 18397},
	name: "_",
},
&litMatcher{
	pos: position{line: 575, col: 48, offset: 18399},
	val: "|",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 575, col: 52, offset: 18403},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 575, col: 54, offset: 18405},
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
	pos: position{line: 595, col: 1, offset: 19127},
	expr: &seqExpr{
	pos: position{line: 595, col: 16, offset: 19144},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 595, col: 16, offset: 19144},
	name: "AnyLabel",
},
&zeroOrOneExpr{
	pos: position{line: 595, col: 25, offset: 19153},
	expr: &seqExpr{
	pos: position{line: 595, col: 26, offset: 19154},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 595, col: 26, offset: 19154},
	name: "_",
},
&litMatcher{
	pos: position{line: 595, col: 28, offset: 19156},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 595, col: 32, offset: 19160},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 595, col: 35, offset: 19163},
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
	pos: position{line: 597, col: 1, offset: 19177},
	expr: &actionExpr{
	pos: position{line: 597, col: 12, offset: 19190},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 597, col: 12, offset: 19190},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 597, col: 12, offset: 19190},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 597, col: 16, offset: 19194},
	name: "_",
},
&labeledExpr{
	pos: position{line: 597, col: 18, offset: 19196},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 597, col: 20, offset: 19198},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 597, col: 31, offset: 19209},
	name: "_",
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 599, col: 1, offset: 19228},
	expr: &actionExpr{
	pos: position{line: 600, col: 7, offset: 19258},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 600, col: 7, offset: 19258},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 600, col: 7, offset: 19258},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 600, col: 11, offset: 19262},
	name: "_",
},
&labeledExpr{
	pos: position{line: 600, col: 13, offset: 19264},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 600, col: 19, offset: 19270},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 600, col: 30, offset: 19281},
	name: "_",
},
&labeledExpr{
	pos: position{line: 600, col: 32, offset: 19283},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 600, col: 37, offset: 19288},
	expr: &ruleRefExpr{
	pos: position{line: 600, col: 37, offset: 19288},
	name: "MoreList",
},
},
},
&litMatcher{
	pos: position{line: 600, col: 47, offset: 19298},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 610, col: 1, offset: 19574},
	expr: &notExpr{
	pos: position{line: 610, col: 7, offset: 19582},
	expr: &anyMatcher{
	line: 610, col: 8, offset: 19583,
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

