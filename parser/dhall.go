
package parser

import (
"bytes"
"errors"
"fmt"
"io"
"io/ioutil"
"math"
"os"
"path"
"strconv"
"strings"
"unicode"
"unicode/utf8"
)
import . "github.com/philandstuff/dhall-golang/ast"


var g = &grammar {
	rules: []*rule{
{
	name: "DhallFile",
	pos: position{line: 22, col: 1, offset: 189},
	expr: &actionExpr{
	pos: position{line: 22, col: 13, offset: 203},
	run: (*parser).callonDhallFile1,
	expr: &seqExpr{
	pos: position{line: 22, col: 13, offset: 203},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 22, col: 13, offset: 203},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 22, col: 15, offset: 205},
	name: "CompleteExpression",
},
},
&ruleRefExpr{
	pos: position{line: 22, col: 34, offset: 224},
	name: "EOF",
},
	},
},
},
},
{
	name: "CompleteExpression",
	pos: position{line: 24, col: 1, offset: 247},
	expr: &actionExpr{
	pos: position{line: 24, col: 22, offset: 270},
	run: (*parser).callonCompleteExpression1,
	expr: &seqExpr{
	pos: position{line: 24, col: 22, offset: 270},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 24, col: 22, offset: 270},
	name: "_",
},
&labeledExpr{
	pos: position{line: 24, col: 24, offset: 272},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 24, col: 26, offset: 274},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 24, col: 37, offset: 285},
	name: "_",
},
	},
},
},
},
{
	name: "EOL",
	pos: position{line: 26, col: 1, offset: 306},
	expr: &choiceExpr{
	pos: position{line: 26, col: 7, offset: 314},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 26, col: 7, offset: 314},
	val: "\n",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 26, col: 14, offset: 321},
	val: "\r\n",
	ignoreCase: false,
},
	},
},
},
{
	name: "BlockComment",
	pos: position{line: 28, col: 1, offset: 329},
	expr: &seqExpr{
	pos: position{line: 28, col: 16, offset: 346},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 28, col: 16, offset: 346},
	val: "{-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 28, col: 21, offset: 351},
	name: "BlockCommentContinue",
},
	},
},
},
{
	name: "BlockCommentChunk",
	pos: position{line: 30, col: 1, offset: 373},
	expr: &choiceExpr{
	pos: position{line: 31, col: 5, offset: 399},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 31, col: 5, offset: 399},
	name: "BlockComment",
},
&charClassMatcher{
	pos: position{line: 32, col: 5, offset: 416},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 33, col: 5, offset: 442},
	name: "EOL",
},
	},
},
},
{
	name: "BlockCommentContinue",
	pos: position{line: 35, col: 1, offset: 447},
	expr: &choiceExpr{
	pos: position{line: 35, col: 24, offset: 472},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 35, col: 24, offset: 472},
	val: "-}",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 35, col: 31, offset: 479},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 35, col: 31, offset: 479},
	name: "BlockCommentChunk",
},
&ruleRefExpr{
	pos: position{line: 35, col: 49, offset: 497},
	name: "BlockCommentContinue",
},
	},
},
	},
},
},
{
	name: "NotEOL",
	pos: position{line: 37, col: 1, offset: 519},
	expr: &charClassMatcher{
	pos: position{line: 37, col: 10, offset: 530},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "LineComment",
	pos: position{line: 39, col: 1, offset: 553},
	expr: &actionExpr{
	pos: position{line: 39, col: 15, offset: 569},
	run: (*parser).callonLineComment1,
	expr: &seqExpr{
	pos: position{line: 39, col: 15, offset: 569},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 39, col: 15, offset: 569},
	val: "--",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 39, col: 20, offset: 574},
	label: "content",
	expr: &actionExpr{
	pos: position{line: 39, col: 29, offset: 583},
	run: (*parser).callonLineComment5,
	expr: &zeroOrMoreExpr{
	pos: position{line: 39, col: 29, offset: 583},
	expr: &ruleRefExpr{
	pos: position{line: 39, col: 29, offset: 583},
	name: "NotEOL",
},
},
},
},
&ruleRefExpr{
	pos: position{line: 39, col: 68, offset: 622},
	name: "EOL",
},
	},
},
},
},
{
	name: "WhitespaceChunk",
	pos: position{line: 41, col: 1, offset: 651},
	expr: &choiceExpr{
	pos: position{line: 41, col: 19, offset: 671},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 41, col: 19, offset: 671},
	val: " ",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 41, col: 25, offset: 677},
	val: "\t",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 41, col: 32, offset: 684},
	name: "EOL",
},
&ruleRefExpr{
	pos: position{line: 41, col: 38, offset: 690},
	name: "LineComment",
},
&ruleRefExpr{
	pos: position{line: 41, col: 52, offset: 704},
	name: "BlockComment",
},
	},
},
},
{
	name: "_",
	pos: position{line: 43, col: 1, offset: 718},
	expr: &zeroOrMoreExpr{
	pos: position{line: 43, col: 5, offset: 724},
	expr: &ruleRefExpr{
	pos: position{line: 43, col: 5, offset: 724},
	name: "WhitespaceChunk",
},
},
},
{
	name: "_1",
	pos: position{line: 45, col: 1, offset: 742},
	expr: &oneOrMoreExpr{
	pos: position{line: 45, col: 6, offset: 749},
	expr: &ruleRefExpr{
	pos: position{line: 45, col: 6, offset: 749},
	name: "WhitespaceChunk",
},
},
},
{
	name: "HexDig",
	pos: position{line: 47, col: 1, offset: 767},
	expr: &charClassMatcher{
	pos: position{line: 47, col: 10, offset: 778},
	val: "[0-9a-f]i",
	ranges: []rune{'0','9','a','f',},
	ignoreCase: true,
	inverted: false,
},
},
{
	name: "SimpleLabelFirstChar",
	pos: position{line: 49, col: 1, offset: 789},
	expr: &charClassMatcher{
	pos: position{line: 49, col: 24, offset: 814},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SimpleLabelNextChar",
	pos: position{line: 50, col: 1, offset: 824},
	expr: &charClassMatcher{
	pos: position{line: 50, col: 23, offset: 848},
	val: "[A-Za-z0-9_/-]",
	chars: []rune{'_','/','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SimpleLabel",
	pos: position{line: 51, col: 1, offset: 863},
	expr: &choiceExpr{
	pos: position{line: 51, col: 15, offset: 879},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 51, col: 15, offset: 879},
	run: (*parser).callonSimpleLabel2,
	expr: &seqExpr{
	pos: position{line: 51, col: 15, offset: 879},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 51, col: 15, offset: 879},
	name: "Keyword",
},
&oneOrMoreExpr{
	pos: position{line: 51, col: 23, offset: 887},
	expr: &ruleRefExpr{
	pos: position{line: 51, col: 23, offset: 887},
	name: "SimpleLabelNextChar",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 52, col: 13, offset: 951},
	run: (*parser).callonSimpleLabel7,
	expr: &seqExpr{
	pos: position{line: 52, col: 13, offset: 951},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 52, col: 13, offset: 951},
	expr: &ruleRefExpr{
	pos: position{line: 52, col: 14, offset: 952},
	name: "Keyword",
},
},
&ruleRefExpr{
	pos: position{line: 52, col: 22, offset: 960},
	name: "SimpleLabelFirstChar",
},
&zeroOrMoreExpr{
	pos: position{line: 52, col: 43, offset: 981},
	expr: &ruleRefExpr{
	pos: position{line: 52, col: 43, offset: 981},
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
	pos: position{line: 59, col: 1, offset: 1082},
	expr: &actionExpr{
	pos: position{line: 59, col: 9, offset: 1092},
	run: (*parser).callonLabel1,
	expr: &labeledExpr{
	pos: position{line: 59, col: 9, offset: 1092},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 59, col: 15, offset: 1098},
	name: "SimpleLabel",
},
},
},
},
{
	name: "NonreservedLabel",
	pos: position{line: 61, col: 1, offset: 1133},
	expr: &choiceExpr{
	pos: position{line: 61, col: 20, offset: 1154},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 61, col: 20, offset: 1154},
	run: (*parser).callonNonreservedLabel2,
	expr: &seqExpr{
	pos: position{line: 61, col: 20, offset: 1154},
	exprs: []interface{}{
&andExpr{
	pos: position{line: 61, col: 20, offset: 1154},
	expr: &seqExpr{
	pos: position{line: 61, col: 22, offset: 1156},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 61, col: 22, offset: 1156},
	name: "Reserved",
},
&ruleRefExpr{
	pos: position{line: 61, col: 31, offset: 1165},
	name: "SimpleLabelNextChar",
},
	},
},
},
&labeledExpr{
	pos: position{line: 61, col: 52, offset: 1186},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 61, col: 58, offset: 1192},
	name: "Label",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 62, col: 19, offset: 1238},
	run: (*parser).callonNonreservedLabel10,
	expr: &seqExpr{
	pos: position{line: 62, col: 19, offset: 1238},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 62, col: 19, offset: 1238},
	expr: &ruleRefExpr{
	pos: position{line: 62, col: 20, offset: 1239},
	name: "Reserved",
},
},
&labeledExpr{
	pos: position{line: 62, col: 29, offset: 1248},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 62, col: 35, offset: 1254},
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
	pos: position{line: 64, col: 1, offset: 1283},
	expr: &ruleRefExpr{
	pos: position{line: 64, col: 12, offset: 1296},
	name: "Label",
},
},
{
	name: "DoubleQuoteChunk",
	pos: position{line: 67, col: 1, offset: 1304},
	expr: &choiceExpr{
	pos: position{line: 68, col: 6, offset: 1330},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 68, col: 6, offset: 1330},
	name: "Interpolation",
},
&actionExpr{
	pos: position{line: 69, col: 6, offset: 1349},
	run: (*parser).callonDoubleQuoteChunk3,
	expr: &seqExpr{
	pos: position{line: 69, col: 6, offset: 1349},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 69, col: 6, offset: 1349},
	val: "\\",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 69, col: 11, offset: 1354},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 69, col: 13, offset: 1356},
	name: "DoubleQuoteEscaped",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 70, col: 6, offset: 1398},
	name: "DoubleQuoteChar",
},
	},
},
},
{
	name: "DoubleQuoteEscaped",
	pos: position{line: 72, col: 1, offset: 1415},
	expr: &choiceExpr{
	pos: position{line: 73, col: 8, offset: 1445},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 73, col: 8, offset: 1445},
	val: "\"",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 74, col: 8, offset: 1456},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 75, col: 8, offset: 1467},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 76, col: 8, offset: 1479},
	val: "/",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 77, col: 8, offset: 1490},
	run: (*parser).callonDoubleQuoteEscaped6,
	expr: &litMatcher{
	pos: position{line: 77, col: 8, offset: 1490},
	val: "b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 78, col: 8, offset: 1530},
	run: (*parser).callonDoubleQuoteEscaped8,
	expr: &litMatcher{
	pos: position{line: 78, col: 8, offset: 1530},
	val: "f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 79, col: 8, offset: 1570},
	run: (*parser).callonDoubleQuoteEscaped10,
	expr: &litMatcher{
	pos: position{line: 79, col: 8, offset: 1570},
	val: "n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 80, col: 8, offset: 1610},
	run: (*parser).callonDoubleQuoteEscaped12,
	expr: &litMatcher{
	pos: position{line: 80, col: 8, offset: 1610},
	val: "r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 81, col: 8, offset: 1650},
	run: (*parser).callonDoubleQuoteEscaped14,
	expr: &litMatcher{
	pos: position{line: 81, col: 8, offset: 1650},
	val: "t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 82, col: 8, offset: 1690},
	run: (*parser).callonDoubleQuoteEscaped16,
	expr: &seqExpr{
	pos: position{line: 82, col: 8, offset: 1690},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 82, col: 8, offset: 1690},
	val: "u",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 82, col: 12, offset: 1694},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 82, col: 19, offset: 1701},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 82, col: 26, offset: 1708},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 82, col: 33, offset: 1715},
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
	pos: position{line: 87, col: 1, offset: 1847},
	expr: &choiceExpr{
	pos: position{line: 88, col: 6, offset: 1872},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 88, col: 6, offset: 1872},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 89, col: 6, offset: 1889},
	val: "[\\x23-\\x5b]",
	ranges: []rune{'#','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 90, col: 6, offset: 1906},
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
	pos: position{line: 92, col: 1, offset: 1925},
	expr: &actionExpr{
	pos: position{line: 92, col: 22, offset: 1948},
	run: (*parser).callonDoubleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 92, col: 22, offset: 1948},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 92, col: 22, offset: 1948},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 92, col: 26, offset: 1952},
	label: "chunks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 92, col: 33, offset: 1959},
	expr: &ruleRefExpr{
	pos: position{line: 92, col: 33, offset: 1959},
	name: "DoubleQuoteChunk",
},
},
},
&litMatcher{
	pos: position{line: 92, col: 51, offset: 1977},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "Interpolation",
	pos: position{line: 109, col: 1, offset: 2445},
	expr: &actionExpr{
	pos: position{line: 109, col: 17, offset: 2463},
	run: (*parser).callonInterpolation1,
	expr: &seqExpr{
	pos: position{line: 109, col: 17, offset: 2463},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 109, col: 17, offset: 2463},
	val: "${",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 109, col: 22, offset: 2468},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 109, col: 24, offset: 2470},
	name: "CompleteExpression",
},
},
&litMatcher{
	pos: position{line: 109, col: 43, offset: 2489},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TextLiteral",
	pos: position{line: 111, col: 1, offset: 2512},
	expr: &actionExpr{
	pos: position{line: 111, col: 15, offset: 2528},
	run: (*parser).callonTextLiteral1,
	expr: &labeledExpr{
	pos: position{line: 111, col: 15, offset: 2528},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 111, col: 17, offset: 2530},
	name: "DoubleQuoteLiteral",
},
},
},
},
{
	name: "Reserved",
	pos: position{line: 114, col: 1, offset: 2653},
	expr: &choiceExpr{
	pos: position{line: 115, col: 5, offset: 2670},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 115, col: 5, offset: 2670},
	run: (*parser).callonReserved2,
	expr: &litMatcher{
	pos: position{line: 115, col: 5, offset: 2670},
	val: "Natural/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 116, col: 5, offset: 2748},
	run: (*parser).callonReserved4,
	expr: &litMatcher{
	pos: position{line: 116, col: 5, offset: 2748},
	val: "Natural/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 117, col: 5, offset: 2824},
	run: (*parser).callonReserved6,
	expr: &litMatcher{
	pos: position{line: 117, col: 5, offset: 2824},
	val: "Natural/isZero",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 118, col: 5, offset: 2904},
	run: (*parser).callonReserved8,
	expr: &litMatcher{
	pos: position{line: 118, col: 5, offset: 2904},
	val: "Natural/even",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 119, col: 5, offset: 2980},
	run: (*parser).callonReserved10,
	expr: &litMatcher{
	pos: position{line: 119, col: 5, offset: 2980},
	val: "Natural/odd",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 120, col: 5, offset: 3054},
	run: (*parser).callonReserved12,
	expr: &litMatcher{
	pos: position{line: 120, col: 5, offset: 3054},
	val: "Natural/toInteger",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 121, col: 5, offset: 3140},
	run: (*parser).callonReserved14,
	expr: &litMatcher{
	pos: position{line: 121, col: 5, offset: 3140},
	val: "Natural/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 122, col: 5, offset: 3216},
	run: (*parser).callonReserved16,
	expr: &litMatcher{
	pos: position{line: 122, col: 5, offset: 3216},
	val: "Integer/toDouble",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 123, col: 5, offset: 3300},
	run: (*parser).callonReserved18,
	expr: &litMatcher{
	pos: position{line: 123, col: 5, offset: 3300},
	val: "Integer/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 124, col: 5, offset: 3376},
	run: (*parser).callonReserved20,
	expr: &litMatcher{
	pos: position{line: 124, col: 5, offset: 3376},
	val: "Double/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 125, col: 5, offset: 3450},
	run: (*parser).callonReserved22,
	expr: &litMatcher{
	pos: position{line: 125, col: 5, offset: 3450},
	val: "List/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 126, col: 5, offset: 3522},
	run: (*parser).callonReserved24,
	expr: &litMatcher{
	pos: position{line: 126, col: 5, offset: 3522},
	val: "List/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 127, col: 5, offset: 3592},
	run: (*parser).callonReserved26,
	expr: &litMatcher{
	pos: position{line: 127, col: 5, offset: 3592},
	val: "List/length",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 128, col: 5, offset: 3666},
	run: (*parser).callonReserved28,
	expr: &litMatcher{
	pos: position{line: 128, col: 5, offset: 3666},
	val: "List/head",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 129, col: 5, offset: 3736},
	run: (*parser).callonReserved30,
	expr: &litMatcher{
	pos: position{line: 129, col: 5, offset: 3736},
	val: "List/last",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 130, col: 5, offset: 3806},
	run: (*parser).callonReserved32,
	expr: &litMatcher{
	pos: position{line: 130, col: 5, offset: 3806},
	val: "List/indexed",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 131, col: 5, offset: 3882},
	run: (*parser).callonReserved34,
	expr: &litMatcher{
	pos: position{line: 131, col: 5, offset: 3882},
	val: "List/reverse",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 132, col: 5, offset: 3958},
	run: (*parser).callonReserved36,
	expr: &litMatcher{
	pos: position{line: 132, col: 5, offset: 3958},
	val: "Optional/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 133, col: 5, offset: 4038},
	run: (*parser).callonReserved38,
	expr: &litMatcher{
	pos: position{line: 133, col: 5, offset: 4038},
	val: "Optional/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 134, col: 5, offset: 4116},
	run: (*parser).callonReserved40,
	expr: &litMatcher{
	pos: position{line: 134, col: 5, offset: 4116},
	val: "Text/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 135, col: 5, offset: 4186},
	run: (*parser).callonReserved42,
	expr: &litMatcher{
	pos: position{line: 135, col: 5, offset: 4186},
	val: "Bool",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 136, col: 5, offset: 4218},
	run: (*parser).callonReserved44,
	expr: &litMatcher{
	pos: position{line: 136, col: 5, offset: 4218},
	val: "True",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 137, col: 5, offset: 4250},
	run: (*parser).callonReserved46,
	expr: &litMatcher{
	pos: position{line: 137, col: 5, offset: 4250},
	val: "False",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 138, col: 5, offset: 4284},
	run: (*parser).callonReserved48,
	expr: &litMatcher{
	pos: position{line: 138, col: 5, offset: 4284},
	val: "Optional",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 139, col: 5, offset: 4324},
	run: (*parser).callonReserved50,
	expr: &litMatcher{
	pos: position{line: 139, col: 5, offset: 4324},
	val: "Natural",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 140, col: 5, offset: 4362},
	run: (*parser).callonReserved52,
	expr: &litMatcher{
	pos: position{line: 140, col: 5, offset: 4362},
	val: "Integer",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 141, col: 5, offset: 4400},
	run: (*parser).callonReserved54,
	expr: &litMatcher{
	pos: position{line: 141, col: 5, offset: 4400},
	val: "Double",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 142, col: 5, offset: 4436},
	run: (*parser).callonReserved56,
	expr: &litMatcher{
	pos: position{line: 142, col: 5, offset: 4436},
	val: "Text",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 143, col: 5, offset: 4468},
	run: (*parser).callonReserved58,
	expr: &litMatcher{
	pos: position{line: 143, col: 5, offset: 4468},
	val: "List",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 144, col: 5, offset: 4500},
	run: (*parser).callonReserved60,
	expr: &litMatcher{
	pos: position{line: 144, col: 5, offset: 4500},
	val: "None",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 145, col: 5, offset: 4532},
	run: (*parser).callonReserved62,
	expr: &litMatcher{
	pos: position{line: 145, col: 5, offset: 4532},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 146, col: 5, offset: 4564},
	run: (*parser).callonReserved64,
	expr: &litMatcher{
	pos: position{line: 146, col: 5, offset: 4564},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 147, col: 5, offset: 4596},
	run: (*parser).callonReserved66,
	expr: &litMatcher{
	pos: position{line: 147, col: 5, offset: 4596},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "If",
	pos: position{line: 149, col: 1, offset: 4625},
	expr: &litMatcher{
	pos: position{line: 149, col: 6, offset: 4632},
	val: "if",
	ignoreCase: false,
},
},
{
	name: "Then",
	pos: position{line: 150, col: 1, offset: 4637},
	expr: &litMatcher{
	pos: position{line: 150, col: 8, offset: 4646},
	val: "then",
	ignoreCase: false,
},
},
{
	name: "Else",
	pos: position{line: 151, col: 1, offset: 4653},
	expr: &litMatcher{
	pos: position{line: 151, col: 8, offset: 4662},
	val: "else",
	ignoreCase: false,
},
},
{
	name: "Let",
	pos: position{line: 152, col: 1, offset: 4669},
	expr: &litMatcher{
	pos: position{line: 152, col: 7, offset: 4677},
	val: "let",
	ignoreCase: false,
},
},
{
	name: "In",
	pos: position{line: 153, col: 1, offset: 4683},
	expr: &litMatcher{
	pos: position{line: 153, col: 6, offset: 4690},
	val: "in",
	ignoreCase: false,
},
},
{
	name: "As",
	pos: position{line: 154, col: 1, offset: 4695},
	expr: &litMatcher{
	pos: position{line: 154, col: 6, offset: 4702},
	val: "as",
	ignoreCase: false,
},
},
{
	name: "Using",
	pos: position{line: 155, col: 1, offset: 4707},
	expr: &litMatcher{
	pos: position{line: 155, col: 9, offset: 4717},
	val: "using",
	ignoreCase: false,
},
},
{
	name: "Merge",
	pos: position{line: 156, col: 1, offset: 4725},
	expr: &litMatcher{
	pos: position{line: 156, col: 9, offset: 4735},
	val: "merge",
	ignoreCase: false,
},
},
{
	name: "Missing",
	pos: position{line: 157, col: 1, offset: 4743},
	expr: &litMatcher{
	pos: position{line: 157, col: 11, offset: 4755},
	val: "missing",
	ignoreCase: false,
},
},
{
	name: "True",
	pos: position{line: 158, col: 1, offset: 4765},
	expr: &litMatcher{
	pos: position{line: 158, col: 8, offset: 4774},
	val: "True",
	ignoreCase: false,
},
},
{
	name: "False",
	pos: position{line: 159, col: 1, offset: 4781},
	expr: &litMatcher{
	pos: position{line: 159, col: 9, offset: 4791},
	val: "False",
	ignoreCase: false,
},
},
{
	name: "Infinity",
	pos: position{line: 160, col: 1, offset: 4799},
	expr: &litMatcher{
	pos: position{line: 160, col: 12, offset: 4812},
	val: "Infinity",
	ignoreCase: false,
},
},
{
	name: "NaN",
	pos: position{line: 161, col: 1, offset: 4823},
	expr: &litMatcher{
	pos: position{line: 161, col: 7, offset: 4831},
	val: "NaN",
	ignoreCase: false,
},
},
{
	name: "Some",
	pos: position{line: 162, col: 1, offset: 4837},
	expr: &litMatcher{
	pos: position{line: 162, col: 8, offset: 4846},
	val: "Some",
	ignoreCase: false,
},
},
{
	name: "Keyword",
	pos: position{line: 164, col: 1, offset: 4854},
	expr: &choiceExpr{
	pos: position{line: 165, col: 5, offset: 4870},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 165, col: 5, offset: 4870},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 165, col: 10, offset: 4875},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 165, col: 17, offset: 4882},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 166, col: 5, offset: 4891},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 166, col: 11, offset: 4897},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 167, col: 5, offset: 4904},
	name: "Using",
},
&ruleRefExpr{
	pos: position{line: 167, col: 13, offset: 4912},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 167, col: 23, offset: 4922},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 168, col: 5, offset: 4929},
	name: "True",
},
&ruleRefExpr{
	pos: position{line: 168, col: 12, offset: 4936},
	name: "False",
},
&ruleRefExpr{
	pos: position{line: 169, col: 5, offset: 4946},
	name: "Infinity",
},
&ruleRefExpr{
	pos: position{line: 169, col: 16, offset: 4957},
	name: "NaN",
},
&ruleRefExpr{
	pos: position{line: 170, col: 5, offset: 4965},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 170, col: 13, offset: 4973},
	name: "Some",
},
	},
},
},
{
	name: "Optional",
	pos: position{line: 172, col: 1, offset: 4979},
	expr: &litMatcher{
	pos: position{line: 172, col: 12, offset: 4992},
	val: "Optional",
	ignoreCase: false,
},
},
{
	name: "Text",
	pos: position{line: 173, col: 1, offset: 5003},
	expr: &litMatcher{
	pos: position{line: 173, col: 8, offset: 5012},
	val: "Text",
	ignoreCase: false,
},
},
{
	name: "List",
	pos: position{line: 174, col: 1, offset: 5019},
	expr: &litMatcher{
	pos: position{line: 174, col: 8, offset: 5028},
	val: "List",
	ignoreCase: false,
},
},
{
	name: "Lambda",
	pos: position{line: 176, col: 1, offset: 5036},
	expr: &choiceExpr{
	pos: position{line: 176, col: 11, offset: 5048},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 176, col: 11, offset: 5048},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 176, col: 18, offset: 5055},
	val: "λ",
	ignoreCase: false,
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 177, col: 1, offset: 5061},
	expr: &choiceExpr{
	pos: position{line: 177, col: 11, offset: 5073},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 177, col: 11, offset: 5073},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 177, col: 22, offset: 5084},
	val: "∀",
	ignoreCase: false,
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 178, col: 1, offset: 5091},
	expr: &choiceExpr{
	pos: position{line: 178, col: 10, offset: 5102},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 178, col: 10, offset: 5102},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 178, col: 17, offset: 5109},
	val: "→",
	ignoreCase: false,
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 180, col: 1, offset: 5117},
	expr: &seqExpr{
	pos: position{line: 180, col: 12, offset: 5130},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 180, col: 12, offset: 5130},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 180, col: 17, offset: 5135},
	expr: &charClassMatcher{
	pos: position{line: 180, col: 17, offset: 5135},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 180, col: 23, offset: 5141},
	expr: &charClassMatcher{
	pos: position{line: 180, col: 23, offset: 5141},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
	},
},
},
{
	name: "NumericDoubleLiteral",
	pos: position{line: 182, col: 1, offset: 5149},
	expr: &actionExpr{
	pos: position{line: 182, col: 24, offset: 5174},
	run: (*parser).callonNumericDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 182, col: 24, offset: 5174},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 182, col: 24, offset: 5174},
	expr: &charClassMatcher{
	pos: position{line: 182, col: 24, offset: 5174},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 182, col: 30, offset: 5180},
	expr: &charClassMatcher{
	pos: position{line: 182, col: 30, offset: 5180},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&choiceExpr{
	pos: position{line: 182, col: 39, offset: 5189},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 182, col: 39, offset: 5189},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 182, col: 39, offset: 5189},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 182, col: 43, offset: 5193},
	expr: &charClassMatcher{
	pos: position{line: 182, col: 43, offset: 5193},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&zeroOrOneExpr{
	pos: position{line: 182, col: 50, offset: 5200},
	expr: &ruleRefExpr{
	pos: position{line: 182, col: 50, offset: 5200},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 182, col: 62, offset: 5212},
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
	pos: position{line: 190, col: 1, offset: 5368},
	expr: &choiceExpr{
	pos: position{line: 190, col: 17, offset: 5386},
	alternatives: []interface{}{
&labeledExpr{
	pos: position{line: 190, col: 17, offset: 5386},
	label: "d",
	expr: &ruleRefExpr{
	pos: position{line: 190, col: 19, offset: 5388},
	name: "NumericDoubleLiteral",
},
},
&actionExpr{
	pos: position{line: 191, col: 5, offset: 5413},
	run: (*parser).callonDoubleLiteral4,
	expr: &ruleRefExpr{
	pos: position{line: 191, col: 5, offset: 5413},
	name: "Infinity",
},
},
&actionExpr{
	pos: position{line: 192, col: 5, offset: 5465},
	run: (*parser).callonDoubleLiteral6,
	expr: &seqExpr{
	pos: position{line: 192, col: 5, offset: 5465},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 192, col: 5, offset: 5465},
	val: "-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 192, col: 9, offset: 5469},
	name: "Infinity",
},
	},
},
},
&actionExpr{
	pos: position{line: 193, col: 5, offset: 5522},
	run: (*parser).callonDoubleLiteral10,
	expr: &ruleRefExpr{
	pos: position{line: 193, col: 5, offset: 5522},
	name: "NaN",
},
},
	},
},
},
{
	name: "NaturalLiteral",
	pos: position{line: 195, col: 1, offset: 5565},
	expr: &actionExpr{
	pos: position{line: 195, col: 18, offset: 5584},
	run: (*parser).callonNaturalLiteral1,
	expr: &oneOrMoreExpr{
	pos: position{line: 195, col: 18, offset: 5584},
	expr: &charClassMatcher{
	pos: position{line: 195, col: 18, offset: 5584},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
},
},
{
	name: "IntegerLiteral",
	pos: position{line: 200, col: 1, offset: 5673},
	expr: &actionExpr{
	pos: position{line: 200, col: 18, offset: 5692},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 200, col: 18, offset: 5692},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 200, col: 18, offset: 5692},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&oneOrMoreExpr{
	pos: position{line: 200, col: 22, offset: 5696},
	expr: &charClassMatcher{
	pos: position{line: 200, col: 22, offset: 5696},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
	},
},
},
},
{
	name: "DeBruijn",
	pos: position{line: 208, col: 1, offset: 5840},
	expr: &actionExpr{
	pos: position{line: 208, col: 12, offset: 5853},
	run: (*parser).callonDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 208, col: 12, offset: 5853},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 208, col: 12, offset: 5853},
	name: "_",
},
&litMatcher{
	pos: position{line: 208, col: 14, offset: 5855},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 208, col: 18, offset: 5859},
	name: "_",
},
&labeledExpr{
	pos: position{line: 208, col: 20, offset: 5861},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 208, col: 26, offset: 5867},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Variable",
	pos: position{line: 210, col: 1, offset: 5923},
	expr: &actionExpr{
	pos: position{line: 210, col: 12, offset: 5936},
	run: (*parser).callonVariable1,
	expr: &seqExpr{
	pos: position{line: 210, col: 12, offset: 5936},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 210, col: 12, offset: 5936},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 210, col: 17, offset: 5941},
	name: "NonreservedLabel",
},
},
&labeledExpr{
	pos: position{line: 210, col: 34, offset: 5958},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 210, col: 40, offset: 5964},
	expr: &ruleRefExpr{
	pos: position{line: 210, col: 40, offset: 5964},
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
	pos: position{line: 218, col: 1, offset: 6127},
	expr: &choiceExpr{
	pos: position{line: 218, col: 14, offset: 6142},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 218, col: 14, offset: 6142},
	name: "Variable",
},
&ruleRefExpr{
	pos: position{line: 218, col: 25, offset: 6153},
	name: "Reserved",
},
	},
},
},
{
	name: "PathCharacter",
	pos: position{line: 220, col: 1, offset: 6163},
	expr: &choiceExpr{
	pos: position{line: 221, col: 6, offset: 6186},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 221, col: 6, offset: 6186},
	val: "!",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 222, col: 6, offset: 6198},
	val: "[\\x24-\\x27]",
	ranges: []rune{'$','\'',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 223, col: 6, offset: 6215},
	val: "[\\x2a-\\x2b]",
	ranges: []rune{'*','+',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 224, col: 6, offset: 6232},
	val: "[\\x2d-\\x2e]",
	ranges: []rune{'-','.',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 225, col: 6, offset: 6249},
	val: "[\\x30-\\x3b]",
	ranges: []rune{'0',';',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 226, col: 6, offset: 6266},
	val: "=",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 227, col: 6, offset: 6278},
	val: "[\\x40-\\x5a]",
	ranges: []rune{'@','Z',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 228, col: 6, offset: 6295},
	val: "[\\x5e-\\x7a]",
	ranges: []rune{'^','z',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 229, col: 6, offset: 6312},
	val: "|",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 230, col: 6, offset: 6324},
	val: "~",
	ignoreCase: false,
},
	},
},
},
{
	name: "UnquotedPathComponent",
	pos: position{line: 232, col: 1, offset: 6332},
	expr: &actionExpr{
	pos: position{line: 232, col: 25, offset: 6358},
	run: (*parser).callonUnquotedPathComponent1,
	expr: &oneOrMoreExpr{
	pos: position{line: 232, col: 25, offset: 6358},
	expr: &ruleRefExpr{
	pos: position{line: 232, col: 25, offset: 6358},
	name: "PathCharacter",
},
},
},
},
{
	name: "PathComponent",
	pos: position{line: 234, col: 1, offset: 6405},
	expr: &actionExpr{
	pos: position{line: 234, col: 17, offset: 6423},
	run: (*parser).callonPathComponent1,
	expr: &seqExpr{
	pos: position{line: 234, col: 17, offset: 6423},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 234, col: 17, offset: 6423},
	val: "/",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 234, col: 21, offset: 6427},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 234, col: 23, offset: 6429},
	name: "UnquotedPathComponent",
},
},
	},
},
},
},
{
	name: "Path",
	pos: position{line: 236, col: 1, offset: 6470},
	expr: &actionExpr{
	pos: position{line: 236, col: 8, offset: 6479},
	run: (*parser).callonPath1,
	expr: &labeledExpr{
	pos: position{line: 236, col: 8, offset: 6479},
	label: "cs",
	expr: &oneOrMoreExpr{
	pos: position{line: 236, col: 11, offset: 6482},
	expr: &ruleRefExpr{
	pos: position{line: 236, col: 11, offset: 6482},
	name: "PathComponent",
},
},
},
},
},
{
	name: "Local",
	pos: position{line: 245, col: 1, offset: 6756},
	expr: &choiceExpr{
	pos: position{line: 245, col: 9, offset: 6766},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 245, col: 9, offset: 6766},
	name: "ParentPath",
},
&ruleRefExpr{
	pos: position{line: 245, col: 22, offset: 6779},
	name: "HerePath",
},
	},
},
},
{
	name: "ParentPath",
	pos: position{line: 247, col: 1, offset: 6789},
	expr: &actionExpr{
	pos: position{line: 247, col: 14, offset: 6804},
	run: (*parser).callonParentPath1,
	expr: &seqExpr{
	pos: position{line: 247, col: 14, offset: 6804},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 247, col: 14, offset: 6804},
	val: "..",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 247, col: 19, offset: 6809},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 247, col: 21, offset: 6811},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HerePath",
	pos: position{line: 248, col: 1, offset: 6867},
	expr: &actionExpr{
	pos: position{line: 248, col: 12, offset: 6880},
	run: (*parser).callonHerePath1,
	expr: &seqExpr{
	pos: position{line: 248, col: 12, offset: 6880},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 248, col: 12, offset: 6880},
	val: ".",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 248, col: 16, offset: 6884},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 248, col: 18, offset: 6886},
	name: "Path",
},
},
	},
},
},
},
{
	name: "Env",
	pos: position{line: 250, col: 1, offset: 6926},
	expr: &actionExpr{
	pos: position{line: 250, col: 7, offset: 6934},
	run: (*parser).callonEnv1,
	expr: &seqExpr{
	pos: position{line: 250, col: 7, offset: 6934},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 250, col: 7, offset: 6934},
	val: "env:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 250, col: 14, offset: 6941},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 250, col: 17, offset: 6944},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 250, col: 17, offset: 6944},
	name: "BashEnvironmentVariable",
},
&ruleRefExpr{
	pos: position{line: 250, col: 43, offset: 6970},
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
	pos: position{line: 252, col: 1, offset: 7015},
	expr: &actionExpr{
	pos: position{line: 252, col: 27, offset: 7043},
	run: (*parser).callonBashEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 252, col: 27, offset: 7043},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 252, col: 27, offset: 7043},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 252, col: 36, offset: 7052},
	expr: &charClassMatcher{
	pos: position{line: 252, col: 36, offset: 7052},
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
	pos: position{line: 256, col: 1, offset: 7108},
	expr: &actionExpr{
	pos: position{line: 256, col: 28, offset: 7137},
	run: (*parser).callonPosixEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 256, col: 28, offset: 7137},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 256, col: 28, offset: 7137},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 256, col: 32, offset: 7141},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 256, col: 34, offset: 7143},
	name: "PosixEnvironmentVariableContent",
},
},
&litMatcher{
	pos: position{line: 256, col: 66, offset: 7175},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PosixEnvironmentVariableContent",
	pos: position{line: 260, col: 1, offset: 7200},
	expr: &actionExpr{
	pos: position{line: 260, col: 35, offset: 7236},
	run: (*parser).callonPosixEnvironmentVariableContent1,
	expr: &labeledExpr{
	pos: position{line: 260, col: 35, offset: 7236},
	label: "v",
	expr: &oneOrMoreExpr{
	pos: position{line: 260, col: 37, offset: 7238},
	expr: &ruleRefExpr{
	pos: position{line: 260, col: 37, offset: 7238},
	name: "PosixEnvironmentVariableCharacter",
},
},
},
},
},
{
	name: "PosixEnvironmentVariableCharacter",
	pos: position{line: 269, col: 1, offset: 7451},
	expr: &choiceExpr{
	pos: position{line: 270, col: 7, offset: 7495},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 270, col: 7, offset: 7495},
	run: (*parser).callonPosixEnvironmentVariableCharacter2,
	expr: &litMatcher{
	pos: position{line: 270, col: 7, offset: 7495},
	val: "\\\"",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 271, col: 7, offset: 7535},
	run: (*parser).callonPosixEnvironmentVariableCharacter4,
	expr: &litMatcher{
	pos: position{line: 271, col: 7, offset: 7535},
	val: "\\\\",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 272, col: 7, offset: 7575},
	run: (*parser).callonPosixEnvironmentVariableCharacter6,
	expr: &litMatcher{
	pos: position{line: 272, col: 7, offset: 7575},
	val: "\\a",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 273, col: 7, offset: 7615},
	run: (*parser).callonPosixEnvironmentVariableCharacter8,
	expr: &litMatcher{
	pos: position{line: 273, col: 7, offset: 7615},
	val: "\\b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 274, col: 7, offset: 7655},
	run: (*parser).callonPosixEnvironmentVariableCharacter10,
	expr: &litMatcher{
	pos: position{line: 274, col: 7, offset: 7655},
	val: "\\f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 275, col: 7, offset: 7695},
	run: (*parser).callonPosixEnvironmentVariableCharacter12,
	expr: &litMatcher{
	pos: position{line: 275, col: 7, offset: 7695},
	val: "\\n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 276, col: 7, offset: 7735},
	run: (*parser).callonPosixEnvironmentVariableCharacter14,
	expr: &litMatcher{
	pos: position{line: 276, col: 7, offset: 7735},
	val: "\\r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 277, col: 7, offset: 7775},
	run: (*parser).callonPosixEnvironmentVariableCharacter16,
	expr: &litMatcher{
	pos: position{line: 277, col: 7, offset: 7775},
	val: "\\t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 278, col: 7, offset: 7815},
	run: (*parser).callonPosixEnvironmentVariableCharacter18,
	expr: &litMatcher{
	pos: position{line: 278, col: 7, offset: 7815},
	val: "\\v",
	ignoreCase: false,
},
},
&charClassMatcher{
	pos: position{line: 279, col: 7, offset: 7855},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 280, col: 7, offset: 7873},
	val: "[\\x23-\\x3c]",
	ranges: []rune{'#','<',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 281, col: 7, offset: 7891},
	val: "[\\x3e-\\x5b]",
	ranges: []rune{'>','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 282, col: 7, offset: 7909},
	val: "[\\x5d-\\x7e]",
	ranges: []rune{']','~',},
	ignoreCase: false,
	inverted: false,
},
	},
},
},
{
	name: "Missing",
	pos: position{line: 284, col: 1, offset: 7922},
	expr: &actionExpr{
	pos: position{line: 284, col: 11, offset: 7934},
	run: (*parser).callonMissing1,
	expr: &seqExpr{
	pos: position{line: 284, col: 11, offset: 7934},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 284, col: 11, offset: 7934},
	val: "missing",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 284, col: 21, offset: 7944},
	name: "_",
},
	},
},
},
},
{
	name: "ImportType",
	pos: position{line: 286, col: 1, offset: 7980},
	expr: &choiceExpr{
	pos: position{line: 286, col: 14, offset: 7995},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 286, col: 14, offset: 7995},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 286, col: 24, offset: 8005},
	name: "Local",
},
&ruleRefExpr{
	pos: position{line: 286, col: 32, offset: 8013},
	name: "Env",
},
	},
},
},
{
	name: "ImportHashed",
	pos: position{line: 288, col: 1, offset: 8018},
	expr: &actionExpr{
	pos: position{line: 288, col: 16, offset: 8035},
	run: (*parser).callonImportHashed1,
	expr: &labeledExpr{
	pos: position{line: 288, col: 16, offset: 8035},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 288, col: 18, offset: 8037},
	name: "ImportType",
},
},
},
},
{
	name: "Import",
	pos: position{line: 290, col: 1, offset: 8106},
	expr: &choiceExpr{
	pos: position{line: 290, col: 10, offset: 8117},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 290, col: 10, offset: 8117},
	run: (*parser).callonImport2,
	expr: &seqExpr{
	pos: position{line: 290, col: 10, offset: 8117},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 290, col: 10, offset: 8117},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 290, col: 12, offset: 8119},
	name: "ImportHashed",
},
},
&ruleRefExpr{
	pos: position{line: 290, col: 25, offset: 8132},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 290, col: 27, offset: 8134},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 290, col: 30, offset: 8137},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 290, col: 33, offset: 8140},
	name: "Text",
},
	},
},
},
&actionExpr{
	pos: position{line: 291, col: 10, offset: 8237},
	run: (*parser).callonImport10,
	expr: &labeledExpr{
	pos: position{line: 291, col: 10, offset: 8237},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 291, col: 12, offset: 8239},
	name: "ImportHashed",
},
},
},
	},
},
},
{
	name: "LetBinding",
	pos: position{line: 294, col: 1, offset: 8334},
	expr: &actionExpr{
	pos: position{line: 294, col: 14, offset: 8349},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 294, col: 14, offset: 8349},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 294, col: 14, offset: 8349},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 294, col: 18, offset: 8353},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 294, col: 21, offset: 8356},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 294, col: 27, offset: 8362},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 294, col: 44, offset: 8379},
	name: "_",
},
&labeledExpr{
	pos: position{line: 294, col: 46, offset: 8381},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 294, col: 48, offset: 8383},
	expr: &seqExpr{
	pos: position{line: 294, col: 49, offset: 8384},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 294, col: 49, offset: 8384},
	name: "Annotation",
},
&ruleRefExpr{
	pos: position{line: 294, col: 60, offset: 8395},
	name: "_",
},
	},
},
},
},
&litMatcher{
	pos: position{line: 295, col: 13, offset: 8411},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 295, col: 17, offset: 8415},
	name: "_",
},
&labeledExpr{
	pos: position{line: 295, col: 19, offset: 8417},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 295, col: 21, offset: 8419},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 295, col: 32, offset: 8430},
	name: "_",
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 310, col: 1, offset: 8739},
	expr: &choiceExpr{
	pos: position{line: 311, col: 7, offset: 8760},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 311, col: 7, offset: 8760},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 311, col: 7, offset: 8760},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 311, col: 7, offset: 8760},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 311, col: 14, offset: 8767},
	name: "_",
},
&litMatcher{
	pos: position{line: 311, col: 16, offset: 8769},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 311, col: 20, offset: 8773},
	name: "_",
},
&labeledExpr{
	pos: position{line: 311, col: 22, offset: 8775},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 311, col: 28, offset: 8781},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 311, col: 45, offset: 8798},
	name: "_",
},
&litMatcher{
	pos: position{line: 311, col: 47, offset: 8800},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 311, col: 51, offset: 8804},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 311, col: 54, offset: 8807},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 311, col: 56, offset: 8809},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 311, col: 67, offset: 8820},
	name: "_",
},
&litMatcher{
	pos: position{line: 311, col: 69, offset: 8822},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 311, col: 73, offset: 8826},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 311, col: 75, offset: 8828},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 311, col: 81, offset: 8834},
	name: "_",
},
&labeledExpr{
	pos: position{line: 311, col: 83, offset: 8836},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 311, col: 88, offset: 8841},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 314, col: 7, offset: 8958},
	run: (*parser).callonExpression22,
	expr: &seqExpr{
	pos: position{line: 314, col: 7, offset: 8958},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 314, col: 7, offset: 8958},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 314, col: 10, offset: 8961},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 314, col: 13, offset: 8964},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 314, col: 18, offset: 8969},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 314, col: 29, offset: 8980},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 314, col: 31, offset: 8982},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 314, col: 36, offset: 8987},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 314, col: 39, offset: 8990},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 314, col: 41, offset: 8992},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 314, col: 52, offset: 9003},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 314, col: 54, offset: 9005},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 314, col: 59, offset: 9010},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 314, col: 62, offset: 9013},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 314, col: 64, offset: 9015},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 317, col: 7, offset: 9101},
	run: (*parser).callonExpression38,
	expr: &seqExpr{
	pos: position{line: 317, col: 7, offset: 9101},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 317, col: 7, offset: 9101},
	label: "bindings",
	expr: &oneOrMoreExpr{
	pos: position{line: 317, col: 16, offset: 9110},
	expr: &ruleRefExpr{
	pos: position{line: 317, col: 16, offset: 9110},
	name: "LetBinding",
},
},
},
&ruleRefExpr{
	pos: position{line: 317, col: 28, offset: 9122},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 317, col: 31, offset: 9125},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 317, col: 34, offset: 9128},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 317, col: 36, offset: 9130},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 324, col: 7, offset: 9370},
	run: (*parser).callonExpression47,
	expr: &seqExpr{
	pos: position{line: 324, col: 7, offset: 9370},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 324, col: 7, offset: 9370},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 324, col: 14, offset: 9377},
	name: "_",
},
&litMatcher{
	pos: position{line: 324, col: 16, offset: 9379},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 324, col: 20, offset: 9383},
	name: "_",
},
&labeledExpr{
	pos: position{line: 324, col: 22, offset: 9385},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 324, col: 28, offset: 9391},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 324, col: 45, offset: 9408},
	name: "_",
},
&litMatcher{
	pos: position{line: 324, col: 47, offset: 9410},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 324, col: 51, offset: 9414},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 324, col: 54, offset: 9417},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 324, col: 56, offset: 9419},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 324, col: 67, offset: 9430},
	name: "_",
},
&litMatcher{
	pos: position{line: 324, col: 69, offset: 9432},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 324, col: 73, offset: 9436},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 324, col: 75, offset: 9438},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 324, col: 81, offset: 9444},
	name: "_",
},
&labeledExpr{
	pos: position{line: 324, col: 83, offset: 9446},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 324, col: 88, offset: 9451},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 327, col: 7, offset: 9560},
	run: (*parser).callonExpression67,
	expr: &seqExpr{
	pos: position{line: 327, col: 7, offset: 9560},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 327, col: 7, offset: 9560},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 327, col: 9, offset: 9562},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 327, col: 28, offset: 9581},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 327, col: 30, offset: 9583},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 327, col: 36, offset: 9589},
	name: "_",
},
&labeledExpr{
	pos: position{line: 327, col: 38, offset: 9591},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 327, col: 40, offset: 9593},
	name: "Expression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 328, col: 7, offset: 9653},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 330, col: 1, offset: 9674},
	expr: &actionExpr{
	pos: position{line: 330, col: 14, offset: 9689},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 330, col: 14, offset: 9689},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 330, col: 14, offset: 9689},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 330, col: 18, offset: 9693},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 330, col: 21, offset: 9696},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 330, col: 23, offset: 9698},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 332, col: 1, offset: 9728},
	expr: &choiceExpr{
	pos: position{line: 333, col: 5, offset: 9756},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 333, col: 5, offset: 9756},
	name: "EmptyList",
},
&actionExpr{
	pos: position{line: 334, col: 5, offset: 9770},
	run: (*parser).callonAnnotatedExpression3,
	expr: &seqExpr{
	pos: position{line: 334, col: 5, offset: 9770},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 334, col: 5, offset: 9770},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 334, col: 7, offset: 9772},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 334, col: 26, offset: 9791},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 334, col: 28, offset: 9793},
	expr: &seqExpr{
	pos: position{line: 334, col: 29, offset: 9794},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 334, col: 29, offset: 9794},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 334, col: 31, offset: 9796},
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
	pos: position{line: 339, col: 1, offset: 9921},
	expr: &actionExpr{
	pos: position{line: 339, col: 13, offset: 9935},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 339, col: 13, offset: 9935},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 339, col: 13, offset: 9935},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 339, col: 17, offset: 9939},
	name: "_",
},
&litMatcher{
	pos: position{line: 339, col: 19, offset: 9941},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 339, col: 23, offset: 9945},
	name: "_",
},
&litMatcher{
	pos: position{line: 339, col: 25, offset: 9947},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 339, col: 29, offset: 9951},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 339, col: 32, offset: 9954},
	name: "List",
},
&ruleRefExpr{
	pos: position{line: 339, col: 37, offset: 9959},
	name: "_",
},
&labeledExpr{
	pos: position{line: 339, col: 39, offset: 9961},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 339, col: 41, offset: 9963},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 343, col: 1, offset: 10026},
	expr: &ruleRefExpr{
	pos: position{line: 343, col: 22, offset: 10049},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 345, col: 1, offset: 10070},
	expr: &ruleRefExpr{
	pos: position{line: 345, col: 23, offset: 10094},
	name: "PlusExpression",
},
},
{
	name: "MorePlus",
	pos: position{line: 347, col: 1, offset: 10110},
	expr: &actionExpr{
	pos: position{line: 347, col: 12, offset: 10123},
	run: (*parser).callonMorePlus1,
	expr: &seqExpr{
	pos: position{line: 347, col: 12, offset: 10123},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 347, col: 12, offset: 10123},
	name: "_",
},
&litMatcher{
	pos: position{line: 347, col: 14, offset: 10125},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 347, col: 18, offset: 10129},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 347, col: 21, offset: 10132},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 347, col: 23, offset: 10134},
	name: "TimesExpression",
},
},
	},
},
},
},
{
	name: "PlusExpression",
	pos: position{line: 348, col: 1, offset: 10168},
	expr: &actionExpr{
	pos: position{line: 349, col: 7, offset: 10193},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 349, col: 7, offset: 10193},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 349, col: 7, offset: 10193},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 349, col: 13, offset: 10199},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 349, col: 29, offset: 10215},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 349, col: 34, offset: 10220},
	expr: &ruleRefExpr{
	pos: position{line: 349, col: 34, offset: 10220},
	name: "MorePlus",
},
},
},
	},
},
},
},
{
	name: "MoreTimes",
	pos: position{line: 358, col: 1, offset: 10448},
	expr: &actionExpr{
	pos: position{line: 358, col: 13, offset: 10462},
	run: (*parser).callonMoreTimes1,
	expr: &seqExpr{
	pos: position{line: 358, col: 13, offset: 10462},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 358, col: 13, offset: 10462},
	name: "_",
},
&litMatcher{
	pos: position{line: 358, col: 15, offset: 10464},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 358, col: 19, offset: 10468},
	name: "_",
},
&labeledExpr{
	pos: position{line: 358, col: 21, offset: 10470},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 358, col: 23, offset: 10472},
	name: "ApplicationExpression",
},
},
	},
},
},
},
{
	name: "TimesExpression",
	pos: position{line: 359, col: 1, offset: 10512},
	expr: &actionExpr{
	pos: position{line: 360, col: 7, offset: 10538},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 360, col: 7, offset: 10538},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 360, col: 7, offset: 10538},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 360, col: 13, offset: 10544},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 360, col: 35, offset: 10566},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 360, col: 40, offset: 10571},
	expr: &ruleRefExpr{
	pos: position{line: 360, col: 40, offset: 10571},
	name: "MoreTimes",
},
},
},
	},
},
},
},
{
	name: "ApplicationExpression",
	pos: position{line: 369, col: 1, offset: 10801},
	expr: &actionExpr{
	pos: position{line: 369, col: 25, offset: 10827},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 369, col: 25, offset: 10827},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 369, col: 25, offset: 10827},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 369, col: 27, offset: 10829},
	name: "FirstApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 369, col: 54, offset: 10856},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 369, col: 59, offset: 10861},
	expr: &seqExpr{
	pos: position{line: 369, col: 60, offset: 10862},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 369, col: 60, offset: 10862},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 369, col: 63, offset: 10865},
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
	pos: position{line: 378, col: 1, offset: 11115},
	expr: &choiceExpr{
	pos: position{line: 379, col: 8, offset: 11153},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 379, col: 8, offset: 11153},
	run: (*parser).callonFirstApplicationExpression2,
	expr: &seqExpr{
	pos: position{line: 379, col: 8, offset: 11153},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 379, col: 8, offset: 11153},
	name: "Some",
},
&ruleRefExpr{
	pos: position{line: 379, col: 13, offset: 11158},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 379, col: 16, offset: 11161},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 379, col: 18, offset: 11163},
	name: "ImportExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 380, col: 8, offset: 11218},
	name: "ImportExpression",
},
	},
},
},
{
	name: "ImportExpression",
	pos: position{line: 382, col: 1, offset: 11236},
	expr: &choiceExpr{
	pos: position{line: 382, col: 20, offset: 11257},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 382, col: 20, offset: 11257},
	name: "Import",
},
&ruleRefExpr{
	pos: position{line: 382, col: 29, offset: 11266},
	name: "SelectorExpression",
},
	},
},
},
{
	name: "SelectorExpression",
	pos: position{line: 384, col: 1, offset: 11286},
	expr: &actionExpr{
	pos: position{line: 384, col: 22, offset: 11309},
	run: (*parser).callonSelectorExpression1,
	expr: &seqExpr{
	pos: position{line: 384, col: 22, offset: 11309},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 384, col: 22, offset: 11309},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 384, col: 24, offset: 11311},
	name: "PrimitiveExpression",
},
},
&labeledExpr{
	pos: position{line: 384, col: 44, offset: 11331},
	label: "ls",
	expr: &zeroOrMoreExpr{
	pos: position{line: 384, col: 47, offset: 11334},
	expr: &seqExpr{
	pos: position{line: 384, col: 48, offset: 11335},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 384, col: 48, offset: 11335},
	name: "_",
},
&litMatcher{
	pos: position{line: 384, col: 50, offset: 11337},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 384, col: 54, offset: 11341},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 384, col: 56, offset: 11343},
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
	pos: position{line: 394, col: 1, offset: 11576},
	expr: &choiceExpr{
	pos: position{line: 395, col: 7, offset: 11606},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 395, col: 7, offset: 11606},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 396, col: 7, offset: 11626},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 397, col: 7, offset: 11647},
	name: "IntegerLiteral",
},
&ruleRefExpr{
	pos: position{line: 398, col: 7, offset: 11668},
	name: "TextLiteral",
},
&actionExpr{
	pos: position{line: 399, col: 7, offset: 11686},
	run: (*parser).callonPrimitiveExpression6,
	expr: &seqExpr{
	pos: position{line: 399, col: 7, offset: 11686},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 399, col: 7, offset: 11686},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 399, col: 11, offset: 11690},
	name: "_",
},
&labeledExpr{
	pos: position{line: 399, col: 13, offset: 11692},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 399, col: 15, offset: 11694},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 399, col: 35, offset: 11714},
	name: "_",
},
&litMatcher{
	pos: position{line: 399, col: 37, offset: 11716},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 400, col: 7, offset: 11744},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 401, col: 7, offset: 11770},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 402, col: 7, offset: 11787},
	run: (*parser).callonPrimitiveExpression16,
	expr: &seqExpr{
	pos: position{line: 402, col: 7, offset: 11787},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 402, col: 7, offset: 11787},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 402, col: 11, offset: 11791},
	name: "_",
},
&labeledExpr{
	pos: position{line: 402, col: 14, offset: 11794},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 402, col: 16, offset: 11796},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 402, col: 27, offset: 11807},
	name: "_",
},
&litMatcher{
	pos: position{line: 402, col: 29, offset: 11809},
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
	pos: position{line: 404, col: 1, offset: 11832},
	expr: &choiceExpr{
	pos: position{line: 405, col: 7, offset: 11862},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 405, col: 7, offset: 11862},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &litMatcher{
	pos: position{line: 405, col: 7, offset: 11862},
	val: "=",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 406, col: 7, offset: 11917},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 407, col: 7, offset: 11942},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 408, col: 7, offset: 11970},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 408, col: 7, offset: 11970},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 410, col: 1, offset: 12016},
	expr: &actionExpr{
	pos: position{line: 410, col: 19, offset: 12036},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 410, col: 19, offset: 12036},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 410, col: 19, offset: 12036},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 410, col: 24, offset: 12041},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 410, col: 33, offset: 12050},
	name: "_",
},
&litMatcher{
	pos: position{line: 410, col: 35, offset: 12052},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 410, col: 39, offset: 12056},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 410, col: 42, offset: 12059},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 410, col: 47, offset: 12064},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 413, col: 1, offset: 12121},
	expr: &actionExpr{
	pos: position{line: 413, col: 18, offset: 12140},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 413, col: 18, offset: 12140},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 413, col: 18, offset: 12140},
	name: "_",
},
&litMatcher{
	pos: position{line: 413, col: 20, offset: 12142},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 413, col: 24, offset: 12146},
	name: "_",
},
&labeledExpr{
	pos: position{line: 413, col: 26, offset: 12148},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 413, col: 28, offset: 12150},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 414, col: 1, offset: 12182},
	expr: &actionExpr{
	pos: position{line: 415, col: 7, offset: 12211},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 415, col: 7, offset: 12211},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 415, col: 7, offset: 12211},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 415, col: 13, offset: 12217},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 415, col: 29, offset: 12233},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 415, col: 34, offset: 12238},
	expr: &ruleRefExpr{
	pos: position{line: 415, col: 34, offset: 12238},
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
	pos: position{line: 425, col: 1, offset: 12634},
	expr: &actionExpr{
	pos: position{line: 425, col: 22, offset: 12657},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 425, col: 22, offset: 12657},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 425, col: 22, offset: 12657},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 425, col: 27, offset: 12662},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 425, col: 36, offset: 12671},
	name: "_",
},
&litMatcher{
	pos: position{line: 425, col: 38, offset: 12673},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 425, col: 42, offset: 12677},
	name: "_",
},
&labeledExpr{
	pos: position{line: 425, col: 44, offset: 12679},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 425, col: 49, offset: 12684},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 428, col: 1, offset: 12741},
	expr: &actionExpr{
	pos: position{line: 428, col: 21, offset: 12763},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 428, col: 21, offset: 12763},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 428, col: 21, offset: 12763},
	name: "_",
},
&litMatcher{
	pos: position{line: 428, col: 23, offset: 12765},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 428, col: 27, offset: 12769},
	name: "_",
},
&labeledExpr{
	pos: position{line: 428, col: 29, offset: 12771},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 428, col: 31, offset: 12773},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 429, col: 1, offset: 12808},
	expr: &actionExpr{
	pos: position{line: 430, col: 7, offset: 12840},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 430, col: 7, offset: 12840},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 430, col: 7, offset: 12840},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 430, col: 13, offset: 12846},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 430, col: 32, offset: 12865},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 430, col: 37, offset: 12870},
	expr: &ruleRefExpr{
	pos: position{line: 430, col: 37, offset: 12870},
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
	pos: position{line: 440, col: 1, offset: 13272},
	expr: &actionExpr{
	pos: position{line: 440, col: 12, offset: 13285},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 440, col: 12, offset: 13285},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 440, col: 12, offset: 13285},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 440, col: 16, offset: 13289},
	name: "_",
},
&labeledExpr{
	pos: position{line: 440, col: 18, offset: 13291},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 440, col: 20, offset: 13293},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 440, col: 31, offset: 13304},
	name: "_",
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 442, col: 1, offset: 13323},
	expr: &actionExpr{
	pos: position{line: 443, col: 7, offset: 13353},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 443, col: 7, offset: 13353},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 443, col: 7, offset: 13353},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 443, col: 11, offset: 13357},
	name: "_",
},
&labeledExpr{
	pos: position{line: 443, col: 13, offset: 13359},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 443, col: 19, offset: 13365},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 443, col: 30, offset: 13376},
	name: "_",
},
&labeledExpr{
	pos: position{line: 443, col: 32, offset: 13378},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 443, col: 37, offset: 13383},
	expr: &ruleRefExpr{
	pos: position{line: 443, col: 37, offset: 13383},
	name: "MoreList",
},
},
},
&litMatcher{
	pos: position{line: 443, col: 47, offset: 13393},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 453, col: 1, offset: 13669},
	expr: &notExpr{
	pos: position{line: 453, col: 7, offset: 13677},
	expr: &anyMatcher{
	line: 453, col: 8, offset: 13678,
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

func (c *current) onInterpolation1(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonInterpolation1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInterpolation1(stack["e"])
}

func (c *current) onTextLiteral1(t interface{}) (interface{}, error) {
 return t, nil 
}

func (p *parser) callonTextLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTextLiteral1(stack["t"])
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

func (c *current) onMissing1() (interface{}, error) {
 var m Missing; return m, nil 
}

func (p *parser) callonMissing1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMissing1()
}

func (c *current) onImportHashed1(i interface{}) (interface{}, error) {
 return ImportHashed{Resolvable: i.(Resolvable)}, nil 
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

func (c *current) onMorePlus1(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonMorePlus1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMorePlus1(stack["e"])
}

func (c *current) onPlusExpression1(first, rest interface{}) (interface{}, error) {
          a := first.(Expr)
          if rest == nil { return a, nil }
          for _, b := range rest.([]interface{}) {
              a = NaturalPlus{L: a, R: b.(Expr)}
          }
          return a, nil
      
}

func (p *parser) callonPlusExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPlusExpression1(stack["first"], stack["rest"])
}

func (c *current) onMoreTimes1(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonMoreTimes1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMoreTimes1(stack["e"])
}

func (c *current) onTimesExpression1(first, rest interface{}) (interface{}, error) {
          a := first.(Expr)
          if rest == nil { return a, nil }
          for _, b := range rest.([]interface{}) {
              a = NaturalTimes{L: a, R: b.(Expr)}
          }
          return a, nil
      
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

