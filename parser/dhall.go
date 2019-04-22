
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
	name: "EscapedChar",
	pos: position{line: 68, col: 1, offset: 1336},
	expr: &actionExpr{
	pos: position{line: 69, col: 3, offset: 1354},
	run: (*parser).callonEscapedChar1,
	expr: &seqExpr{
	pos: position{line: 69, col: 3, offset: 1354},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 69, col: 3, offset: 1354},
	val: "\\",
	ignoreCase: false,
},
&choiceExpr{
	pos: position{line: 70, col: 5, offset: 1363},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 70, col: 5, offset: 1363},
	val: "\"",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 71, col: 10, offset: 1376},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 72, col: 10, offset: 1389},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 73, col: 10, offset: 1403},
	val: "/",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 74, col: 10, offset: 1416},
	val: "b",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 75, col: 10, offset: 1429},
	val: "f",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 76, col: 10, offset: 1442},
	val: "n",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 77, col: 10, offset: 1455},
	val: "r",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 78, col: 10, offset: 1468},
	val: "t",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 79, col: 10, offset: 1481},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 79, col: 10, offset: 1481},
	val: "u",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 79, col: 14, offset: 1485},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 79, col: 21, offset: 1492},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 79, col: 28, offset: 1499},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 79, col: 35, offset: 1506},
	name: "HexDig",
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
	name: "DoubleQuoteChunk",
	pos: position{line: 100, col: 1, offset: 1949},
	expr: &choiceExpr{
	pos: position{line: 101, col: 6, offset: 1975},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 101, col: 6, offset: 1975},
	run: (*parser).callonDoubleQuoteChunk2,
	expr: &seqExpr{
	pos: position{line: 101, col: 6, offset: 1975},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 101, col: 6, offset: 1975},
	val: "${",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 101, col: 11, offset: 1980},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 101, col: 13, offset: 1982},
	name: "CompleteExpression",
},
},
&litMatcher{
	pos: position{line: 101, col: 32, offset: 2001},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 102, col: 6, offset: 2028},
	name: "EscapedChar",
},
&charClassMatcher{
	pos: position{line: 103, col: 6, offset: 2045},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 104, col: 6, offset: 2062},
	val: "[\\x23-\\x5b]",
	ranges: []rune{'#','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 105, col: 6, offset: 2079},
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
	pos: position{line: 107, col: 1, offset: 2098},
	expr: &actionExpr{
	pos: position{line: 107, col: 22, offset: 2121},
	run: (*parser).callonDoubleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 107, col: 22, offset: 2121},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 107, col: 22, offset: 2121},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 107, col: 26, offset: 2125},
	label: "chunks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 107, col: 33, offset: 2132},
	expr: &ruleRefExpr{
	pos: position{line: 107, col: 33, offset: 2132},
	name: "DoubleQuoteChunk",
},
},
},
&litMatcher{
	pos: position{line: 107, col: 51, offset: 2150},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TextLiteral",
	pos: position{line: 124, col: 1, offset: 2618},
	expr: &actionExpr{
	pos: position{line: 124, col: 15, offset: 2634},
	run: (*parser).callonTextLiteral1,
	expr: &labeledExpr{
	pos: position{line: 124, col: 15, offset: 2634},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 124, col: 17, offset: 2636},
	name: "DoubleQuoteLiteral",
},
},
},
},
{
	name: "Reserved",
	pos: position{line: 127, col: 1, offset: 2759},
	expr: &choiceExpr{
	pos: position{line: 128, col: 5, offset: 2776},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 128, col: 5, offset: 2776},
	run: (*parser).callonReserved2,
	expr: &litMatcher{
	pos: position{line: 128, col: 5, offset: 2776},
	val: "Natural/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 129, col: 5, offset: 2854},
	run: (*parser).callonReserved4,
	expr: &litMatcher{
	pos: position{line: 129, col: 5, offset: 2854},
	val: "Natural/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 130, col: 5, offset: 2930},
	run: (*parser).callonReserved6,
	expr: &litMatcher{
	pos: position{line: 130, col: 5, offset: 2930},
	val: "Natural/isZero",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 131, col: 5, offset: 3010},
	run: (*parser).callonReserved8,
	expr: &litMatcher{
	pos: position{line: 131, col: 5, offset: 3010},
	val: "Natural/even",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 132, col: 5, offset: 3086},
	run: (*parser).callonReserved10,
	expr: &litMatcher{
	pos: position{line: 132, col: 5, offset: 3086},
	val: "Natural/odd",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 133, col: 5, offset: 3160},
	run: (*parser).callonReserved12,
	expr: &litMatcher{
	pos: position{line: 133, col: 5, offset: 3160},
	val: "Natural/toInteger",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 134, col: 5, offset: 3246},
	run: (*parser).callonReserved14,
	expr: &litMatcher{
	pos: position{line: 134, col: 5, offset: 3246},
	val: "Natural/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 135, col: 5, offset: 3322},
	run: (*parser).callonReserved16,
	expr: &litMatcher{
	pos: position{line: 135, col: 5, offset: 3322},
	val: "Integer/toDouble",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 136, col: 5, offset: 3406},
	run: (*parser).callonReserved18,
	expr: &litMatcher{
	pos: position{line: 136, col: 5, offset: 3406},
	val: "Integer/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 137, col: 5, offset: 3482},
	run: (*parser).callonReserved20,
	expr: &litMatcher{
	pos: position{line: 137, col: 5, offset: 3482},
	val: "Double/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 138, col: 5, offset: 3556},
	run: (*parser).callonReserved22,
	expr: &litMatcher{
	pos: position{line: 138, col: 5, offset: 3556},
	val: "List/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 139, col: 5, offset: 3628},
	run: (*parser).callonReserved24,
	expr: &litMatcher{
	pos: position{line: 139, col: 5, offset: 3628},
	val: "List/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 140, col: 5, offset: 3698},
	run: (*parser).callonReserved26,
	expr: &litMatcher{
	pos: position{line: 140, col: 5, offset: 3698},
	val: "List/length",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 141, col: 5, offset: 3772},
	run: (*parser).callonReserved28,
	expr: &litMatcher{
	pos: position{line: 141, col: 5, offset: 3772},
	val: "List/head",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 142, col: 5, offset: 3842},
	run: (*parser).callonReserved30,
	expr: &litMatcher{
	pos: position{line: 142, col: 5, offset: 3842},
	val: "List/last",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 143, col: 5, offset: 3912},
	run: (*parser).callonReserved32,
	expr: &litMatcher{
	pos: position{line: 143, col: 5, offset: 3912},
	val: "List/indexed",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 144, col: 5, offset: 3988},
	run: (*parser).callonReserved34,
	expr: &litMatcher{
	pos: position{line: 144, col: 5, offset: 3988},
	val: "List/reverse",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 145, col: 5, offset: 4064},
	run: (*parser).callonReserved36,
	expr: &litMatcher{
	pos: position{line: 145, col: 5, offset: 4064},
	val: "Optional/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 146, col: 5, offset: 4144},
	run: (*parser).callonReserved38,
	expr: &litMatcher{
	pos: position{line: 146, col: 5, offset: 4144},
	val: "Optional/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 147, col: 5, offset: 4222},
	run: (*parser).callonReserved40,
	expr: &litMatcher{
	pos: position{line: 147, col: 5, offset: 4222},
	val: "Text/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 148, col: 5, offset: 4292},
	run: (*parser).callonReserved42,
	expr: &litMatcher{
	pos: position{line: 148, col: 5, offset: 4292},
	val: "Bool",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 149, col: 5, offset: 4324},
	run: (*parser).callonReserved44,
	expr: &litMatcher{
	pos: position{line: 149, col: 5, offset: 4324},
	val: "True",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 150, col: 5, offset: 4356},
	run: (*parser).callonReserved46,
	expr: &litMatcher{
	pos: position{line: 150, col: 5, offset: 4356},
	val: "False",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 151, col: 5, offset: 4390},
	run: (*parser).callonReserved48,
	expr: &litMatcher{
	pos: position{line: 151, col: 5, offset: 4390},
	val: "Optional",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 152, col: 5, offset: 4430},
	run: (*parser).callonReserved50,
	expr: &litMatcher{
	pos: position{line: 152, col: 5, offset: 4430},
	val: "Natural",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 153, col: 5, offset: 4468},
	run: (*parser).callonReserved52,
	expr: &litMatcher{
	pos: position{line: 153, col: 5, offset: 4468},
	val: "Integer",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 154, col: 5, offset: 4506},
	run: (*parser).callonReserved54,
	expr: &litMatcher{
	pos: position{line: 154, col: 5, offset: 4506},
	val: "Double",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 155, col: 5, offset: 4542},
	run: (*parser).callonReserved56,
	expr: &litMatcher{
	pos: position{line: 155, col: 5, offset: 4542},
	val: "Text",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 156, col: 5, offset: 4574},
	run: (*parser).callonReserved58,
	expr: &litMatcher{
	pos: position{line: 156, col: 5, offset: 4574},
	val: "List",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 157, col: 5, offset: 4606},
	run: (*parser).callonReserved60,
	expr: &litMatcher{
	pos: position{line: 157, col: 5, offset: 4606},
	val: "None",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 158, col: 5, offset: 4638},
	run: (*parser).callonReserved62,
	expr: &litMatcher{
	pos: position{line: 158, col: 5, offset: 4638},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 159, col: 5, offset: 4670},
	run: (*parser).callonReserved64,
	expr: &litMatcher{
	pos: position{line: 159, col: 5, offset: 4670},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 160, col: 5, offset: 4702},
	run: (*parser).callonReserved66,
	expr: &litMatcher{
	pos: position{line: 160, col: 5, offset: 4702},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "If",
	pos: position{line: 162, col: 1, offset: 4731},
	expr: &litMatcher{
	pos: position{line: 162, col: 6, offset: 4738},
	val: "if",
	ignoreCase: false,
},
},
{
	name: "Then",
	pos: position{line: 163, col: 1, offset: 4743},
	expr: &litMatcher{
	pos: position{line: 163, col: 8, offset: 4752},
	val: "then",
	ignoreCase: false,
},
},
{
	name: "Else",
	pos: position{line: 164, col: 1, offset: 4759},
	expr: &litMatcher{
	pos: position{line: 164, col: 8, offset: 4768},
	val: "else",
	ignoreCase: false,
},
},
{
	name: "Let",
	pos: position{line: 165, col: 1, offset: 4775},
	expr: &litMatcher{
	pos: position{line: 165, col: 7, offset: 4783},
	val: "let",
	ignoreCase: false,
},
},
{
	name: "In",
	pos: position{line: 166, col: 1, offset: 4789},
	expr: &litMatcher{
	pos: position{line: 166, col: 6, offset: 4796},
	val: "in",
	ignoreCase: false,
},
},
{
	name: "As",
	pos: position{line: 167, col: 1, offset: 4801},
	expr: &litMatcher{
	pos: position{line: 167, col: 6, offset: 4808},
	val: "as",
	ignoreCase: false,
},
},
{
	name: "Using",
	pos: position{line: 168, col: 1, offset: 4813},
	expr: &litMatcher{
	pos: position{line: 168, col: 9, offset: 4823},
	val: "using",
	ignoreCase: false,
},
},
{
	name: "Merge",
	pos: position{line: 169, col: 1, offset: 4831},
	expr: &litMatcher{
	pos: position{line: 169, col: 9, offset: 4841},
	val: "merge",
	ignoreCase: false,
},
},
{
	name: "Missing",
	pos: position{line: 170, col: 1, offset: 4849},
	expr: &litMatcher{
	pos: position{line: 170, col: 11, offset: 4861},
	val: "missing",
	ignoreCase: false,
},
},
{
	name: "True",
	pos: position{line: 171, col: 1, offset: 4871},
	expr: &litMatcher{
	pos: position{line: 171, col: 8, offset: 4880},
	val: "True",
	ignoreCase: false,
},
},
{
	name: "False",
	pos: position{line: 172, col: 1, offset: 4887},
	expr: &litMatcher{
	pos: position{line: 172, col: 9, offset: 4897},
	val: "False",
	ignoreCase: false,
},
},
{
	name: "Infinity",
	pos: position{line: 173, col: 1, offset: 4905},
	expr: &litMatcher{
	pos: position{line: 173, col: 12, offset: 4918},
	val: "Infinity",
	ignoreCase: false,
},
},
{
	name: "NaN",
	pos: position{line: 174, col: 1, offset: 4929},
	expr: &litMatcher{
	pos: position{line: 174, col: 7, offset: 4937},
	val: "NaN",
	ignoreCase: false,
},
},
{
	name: "Some",
	pos: position{line: 175, col: 1, offset: 4943},
	expr: &litMatcher{
	pos: position{line: 175, col: 8, offset: 4952},
	val: "Some",
	ignoreCase: false,
},
},
{
	name: "Keyword",
	pos: position{line: 177, col: 1, offset: 4960},
	expr: &choiceExpr{
	pos: position{line: 178, col: 5, offset: 4976},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 178, col: 5, offset: 4976},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 178, col: 10, offset: 4981},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 178, col: 17, offset: 4988},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 179, col: 5, offset: 4997},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 179, col: 11, offset: 5003},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 180, col: 5, offset: 5010},
	name: "Using",
},
&ruleRefExpr{
	pos: position{line: 180, col: 13, offset: 5018},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 180, col: 23, offset: 5028},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 181, col: 5, offset: 5035},
	name: "True",
},
&ruleRefExpr{
	pos: position{line: 181, col: 12, offset: 5042},
	name: "False",
},
&ruleRefExpr{
	pos: position{line: 182, col: 5, offset: 5052},
	name: "Infinity",
},
&ruleRefExpr{
	pos: position{line: 182, col: 16, offset: 5063},
	name: "NaN",
},
&ruleRefExpr{
	pos: position{line: 183, col: 5, offset: 5071},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 183, col: 13, offset: 5079},
	name: "Some",
},
	},
},
},
{
	name: "Optional",
	pos: position{line: 185, col: 1, offset: 5085},
	expr: &litMatcher{
	pos: position{line: 185, col: 12, offset: 5098},
	val: "Optional",
	ignoreCase: false,
},
},
{
	name: "Text",
	pos: position{line: 186, col: 1, offset: 5109},
	expr: &litMatcher{
	pos: position{line: 186, col: 8, offset: 5118},
	val: "Text",
	ignoreCase: false,
},
},
{
	name: "List",
	pos: position{line: 187, col: 1, offset: 5125},
	expr: &litMatcher{
	pos: position{line: 187, col: 8, offset: 5134},
	val: "List",
	ignoreCase: false,
},
},
{
	name: "Lambda",
	pos: position{line: 189, col: 1, offset: 5142},
	expr: &choiceExpr{
	pos: position{line: 189, col: 11, offset: 5154},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 189, col: 11, offset: 5154},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 189, col: 18, offset: 5161},
	val: "λ",
	ignoreCase: false,
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 190, col: 1, offset: 5167},
	expr: &choiceExpr{
	pos: position{line: 190, col: 11, offset: 5179},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 190, col: 11, offset: 5179},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 190, col: 22, offset: 5190},
	val: "∀",
	ignoreCase: false,
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 191, col: 1, offset: 5197},
	expr: &choiceExpr{
	pos: position{line: 191, col: 10, offset: 5208},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 191, col: 10, offset: 5208},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 191, col: 17, offset: 5215},
	val: "→",
	ignoreCase: false,
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 193, col: 1, offset: 5223},
	expr: &seqExpr{
	pos: position{line: 193, col: 12, offset: 5236},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 193, col: 12, offset: 5236},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 193, col: 17, offset: 5241},
	expr: &charClassMatcher{
	pos: position{line: 193, col: 17, offset: 5241},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 193, col: 23, offset: 5247},
	expr: &charClassMatcher{
	pos: position{line: 193, col: 23, offset: 5247},
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
	pos: position{line: 195, col: 1, offset: 5255},
	expr: &actionExpr{
	pos: position{line: 195, col: 24, offset: 5280},
	run: (*parser).callonNumericDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 195, col: 24, offset: 5280},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 195, col: 24, offset: 5280},
	expr: &charClassMatcher{
	pos: position{line: 195, col: 24, offset: 5280},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 195, col: 30, offset: 5286},
	expr: &charClassMatcher{
	pos: position{line: 195, col: 30, offset: 5286},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&choiceExpr{
	pos: position{line: 195, col: 39, offset: 5295},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 195, col: 39, offset: 5295},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 195, col: 39, offset: 5295},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 195, col: 43, offset: 5299},
	expr: &charClassMatcher{
	pos: position{line: 195, col: 43, offset: 5299},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&zeroOrOneExpr{
	pos: position{line: 195, col: 50, offset: 5306},
	expr: &ruleRefExpr{
	pos: position{line: 195, col: 50, offset: 5306},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 195, col: 62, offset: 5318},
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
	pos: position{line: 203, col: 1, offset: 5474},
	expr: &choiceExpr{
	pos: position{line: 203, col: 17, offset: 5492},
	alternatives: []interface{}{
&labeledExpr{
	pos: position{line: 203, col: 17, offset: 5492},
	label: "d",
	expr: &ruleRefExpr{
	pos: position{line: 203, col: 19, offset: 5494},
	name: "NumericDoubleLiteral",
},
},
&actionExpr{
	pos: position{line: 204, col: 5, offset: 5519},
	run: (*parser).callonDoubleLiteral4,
	expr: &ruleRefExpr{
	pos: position{line: 204, col: 5, offset: 5519},
	name: "Infinity",
},
},
&actionExpr{
	pos: position{line: 205, col: 5, offset: 5571},
	run: (*parser).callonDoubleLiteral6,
	expr: &seqExpr{
	pos: position{line: 205, col: 5, offset: 5571},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 205, col: 5, offset: 5571},
	val: "-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 205, col: 9, offset: 5575},
	name: "Infinity",
},
	},
},
},
&actionExpr{
	pos: position{line: 206, col: 5, offset: 5628},
	run: (*parser).callonDoubleLiteral10,
	expr: &ruleRefExpr{
	pos: position{line: 206, col: 5, offset: 5628},
	name: "NaN",
},
},
	},
},
},
{
	name: "NaturalLiteral",
	pos: position{line: 208, col: 1, offset: 5671},
	expr: &actionExpr{
	pos: position{line: 208, col: 18, offset: 5690},
	run: (*parser).callonNaturalLiteral1,
	expr: &oneOrMoreExpr{
	pos: position{line: 208, col: 18, offset: 5690},
	expr: &charClassMatcher{
	pos: position{line: 208, col: 18, offset: 5690},
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
	pos: position{line: 213, col: 1, offset: 5779},
	expr: &actionExpr{
	pos: position{line: 213, col: 18, offset: 5798},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 213, col: 18, offset: 5798},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 213, col: 18, offset: 5798},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&oneOrMoreExpr{
	pos: position{line: 213, col: 22, offset: 5802},
	expr: &charClassMatcher{
	pos: position{line: 213, col: 22, offset: 5802},
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
	pos: position{line: 221, col: 1, offset: 5946},
	expr: &actionExpr{
	pos: position{line: 221, col: 12, offset: 5959},
	run: (*parser).callonDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 221, col: 12, offset: 5959},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 221, col: 12, offset: 5959},
	name: "_",
},
&litMatcher{
	pos: position{line: 221, col: 14, offset: 5961},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 221, col: 18, offset: 5965},
	name: "_",
},
&labeledExpr{
	pos: position{line: 221, col: 20, offset: 5967},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 221, col: 26, offset: 5973},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Variable",
	pos: position{line: 223, col: 1, offset: 6029},
	expr: &actionExpr{
	pos: position{line: 223, col: 12, offset: 6042},
	run: (*parser).callonVariable1,
	expr: &seqExpr{
	pos: position{line: 223, col: 12, offset: 6042},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 223, col: 12, offset: 6042},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 223, col: 17, offset: 6047},
	name: "NonreservedLabel",
},
},
&labeledExpr{
	pos: position{line: 223, col: 34, offset: 6064},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 223, col: 40, offset: 6070},
	expr: &ruleRefExpr{
	pos: position{line: 223, col: 40, offset: 6070},
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
	pos: position{line: 231, col: 1, offset: 6233},
	expr: &choiceExpr{
	pos: position{line: 231, col: 14, offset: 6248},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 231, col: 14, offset: 6248},
	name: "Variable",
},
&ruleRefExpr{
	pos: position{line: 231, col: 25, offset: 6259},
	name: "Reserved",
},
	},
},
},
{
	name: "PathCharacter",
	pos: position{line: 233, col: 1, offset: 6269},
	expr: &choiceExpr{
	pos: position{line: 234, col: 6, offset: 6292},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 234, col: 6, offset: 6292},
	val: "!",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 235, col: 6, offset: 6304},
	val: "[\\x24-\\x27]",
	ranges: []rune{'$','\'',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 236, col: 6, offset: 6321},
	val: "[\\x2a-\\x2b]",
	ranges: []rune{'*','+',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 237, col: 6, offset: 6338},
	val: "[\\x2d-\\x2e]",
	ranges: []rune{'-','.',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 238, col: 6, offset: 6355},
	val: "[\\x30-\\x3b]",
	ranges: []rune{'0',';',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 239, col: 6, offset: 6372},
	val: "=",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 240, col: 6, offset: 6384},
	val: "[\\x40-\\x5a]",
	ranges: []rune{'@','Z',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 241, col: 6, offset: 6401},
	val: "[\\x5e-\\x7a]",
	ranges: []rune{'^','z',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 242, col: 6, offset: 6418},
	val: "|",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 243, col: 6, offset: 6430},
	val: "~",
	ignoreCase: false,
},
	},
},
},
{
	name: "UnquotedPathComponent",
	pos: position{line: 245, col: 1, offset: 6438},
	expr: &actionExpr{
	pos: position{line: 245, col: 25, offset: 6464},
	run: (*parser).callonUnquotedPathComponent1,
	expr: &oneOrMoreExpr{
	pos: position{line: 245, col: 25, offset: 6464},
	expr: &ruleRefExpr{
	pos: position{line: 245, col: 25, offset: 6464},
	name: "PathCharacter",
},
},
},
},
{
	name: "PathComponent",
	pos: position{line: 247, col: 1, offset: 6511},
	expr: &actionExpr{
	pos: position{line: 247, col: 17, offset: 6529},
	run: (*parser).callonPathComponent1,
	expr: &seqExpr{
	pos: position{line: 247, col: 17, offset: 6529},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 247, col: 17, offset: 6529},
	val: "/",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 247, col: 21, offset: 6533},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 247, col: 23, offset: 6535},
	name: "UnquotedPathComponent",
},
},
	},
},
},
},
{
	name: "Path",
	pos: position{line: 249, col: 1, offset: 6576},
	expr: &actionExpr{
	pos: position{line: 249, col: 8, offset: 6585},
	run: (*parser).callonPath1,
	expr: &labeledExpr{
	pos: position{line: 249, col: 8, offset: 6585},
	label: "cs",
	expr: &oneOrMoreExpr{
	pos: position{line: 249, col: 11, offset: 6588},
	expr: &ruleRefExpr{
	pos: position{line: 249, col: 11, offset: 6588},
	name: "PathComponent",
},
},
},
},
},
{
	name: "Local",
	pos: position{line: 258, col: 1, offset: 6862},
	expr: &choiceExpr{
	pos: position{line: 258, col: 9, offset: 6872},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 258, col: 9, offset: 6872},
	name: "ParentPath",
},
&ruleRefExpr{
	pos: position{line: 258, col: 22, offset: 6885},
	name: "HerePath",
},
	},
},
},
{
	name: "ParentPath",
	pos: position{line: 260, col: 1, offset: 6895},
	expr: &actionExpr{
	pos: position{line: 260, col: 14, offset: 6910},
	run: (*parser).callonParentPath1,
	expr: &seqExpr{
	pos: position{line: 260, col: 14, offset: 6910},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 260, col: 14, offset: 6910},
	val: "..",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 260, col: 19, offset: 6915},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 260, col: 21, offset: 6917},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HerePath",
	pos: position{line: 261, col: 1, offset: 6973},
	expr: &actionExpr{
	pos: position{line: 261, col: 12, offset: 6986},
	run: (*parser).callonHerePath1,
	expr: &seqExpr{
	pos: position{line: 261, col: 12, offset: 6986},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 261, col: 12, offset: 6986},
	val: ".",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 261, col: 16, offset: 6990},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 261, col: 18, offset: 6992},
	name: "Path",
},
},
	},
},
},
},
{
	name: "Env",
	pos: position{line: 263, col: 1, offset: 7032},
	expr: &actionExpr{
	pos: position{line: 263, col: 7, offset: 7040},
	run: (*parser).callonEnv1,
	expr: &seqExpr{
	pos: position{line: 263, col: 7, offset: 7040},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 263, col: 7, offset: 7040},
	val: "env:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 263, col: 14, offset: 7047},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 263, col: 17, offset: 7050},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 263, col: 17, offset: 7050},
	name: "BashEnvironmentVariable",
},
&ruleRefExpr{
	pos: position{line: 263, col: 43, offset: 7076},
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
	pos: position{line: 265, col: 1, offset: 7121},
	expr: &actionExpr{
	pos: position{line: 265, col: 27, offset: 7149},
	run: (*parser).callonBashEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 265, col: 27, offset: 7149},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 265, col: 27, offset: 7149},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 265, col: 36, offset: 7158},
	expr: &charClassMatcher{
	pos: position{line: 265, col: 36, offset: 7158},
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
	pos: position{line: 269, col: 1, offset: 7214},
	expr: &actionExpr{
	pos: position{line: 269, col: 28, offset: 7243},
	run: (*parser).callonPosixEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 269, col: 28, offset: 7243},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 269, col: 28, offset: 7243},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 269, col: 32, offset: 7247},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 269, col: 34, offset: 7249},
	name: "PosixEnvironmentVariableContent",
},
},
&litMatcher{
	pos: position{line: 269, col: 66, offset: 7281},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PosixEnvironmentVariableContent",
	pos: position{line: 273, col: 1, offset: 7306},
	expr: &actionExpr{
	pos: position{line: 273, col: 35, offset: 7342},
	run: (*parser).callonPosixEnvironmentVariableContent1,
	expr: &labeledExpr{
	pos: position{line: 273, col: 35, offset: 7342},
	label: "v",
	expr: &oneOrMoreExpr{
	pos: position{line: 273, col: 37, offset: 7344},
	expr: &ruleRefExpr{
	pos: position{line: 273, col: 37, offset: 7344},
	name: "PosixEnvironmentVariableCharacter",
},
},
},
},
},
{
	name: "PosixEnvironmentVariableCharacter",
	pos: position{line: 282, col: 1, offset: 7557},
	expr: &choiceExpr{
	pos: position{line: 283, col: 7, offset: 7601},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 283, col: 7, offset: 7601},
	run: (*parser).callonPosixEnvironmentVariableCharacter2,
	expr: &litMatcher{
	pos: position{line: 283, col: 7, offset: 7601},
	val: "\\\"",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 284, col: 7, offset: 7641},
	run: (*parser).callonPosixEnvironmentVariableCharacter4,
	expr: &litMatcher{
	pos: position{line: 284, col: 7, offset: 7641},
	val: "\\\\",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 285, col: 7, offset: 7681},
	run: (*parser).callonPosixEnvironmentVariableCharacter6,
	expr: &litMatcher{
	pos: position{line: 285, col: 7, offset: 7681},
	val: "\\a",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 286, col: 7, offset: 7721},
	run: (*parser).callonPosixEnvironmentVariableCharacter8,
	expr: &litMatcher{
	pos: position{line: 286, col: 7, offset: 7721},
	val: "\\b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 287, col: 7, offset: 7761},
	run: (*parser).callonPosixEnvironmentVariableCharacter10,
	expr: &litMatcher{
	pos: position{line: 287, col: 7, offset: 7761},
	val: "\\f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 288, col: 7, offset: 7801},
	run: (*parser).callonPosixEnvironmentVariableCharacter12,
	expr: &litMatcher{
	pos: position{line: 288, col: 7, offset: 7801},
	val: "\\n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 289, col: 7, offset: 7841},
	run: (*parser).callonPosixEnvironmentVariableCharacter14,
	expr: &litMatcher{
	pos: position{line: 289, col: 7, offset: 7841},
	val: "\\r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 290, col: 7, offset: 7881},
	run: (*parser).callonPosixEnvironmentVariableCharacter16,
	expr: &litMatcher{
	pos: position{line: 290, col: 7, offset: 7881},
	val: "\\t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 291, col: 7, offset: 7921},
	run: (*parser).callonPosixEnvironmentVariableCharacter18,
	expr: &litMatcher{
	pos: position{line: 291, col: 7, offset: 7921},
	val: "\\v",
	ignoreCase: false,
},
},
&charClassMatcher{
	pos: position{line: 292, col: 7, offset: 7961},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 293, col: 7, offset: 7979},
	val: "[\\x23-\\x3c]",
	ranges: []rune{'#','<',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 294, col: 7, offset: 7997},
	val: "[\\x3e-\\x5b]",
	ranges: []rune{'>','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 295, col: 7, offset: 8015},
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
	pos: position{line: 297, col: 1, offset: 8028},
	expr: &actionExpr{
	pos: position{line: 297, col: 11, offset: 8040},
	run: (*parser).callonMissing1,
	expr: &seqExpr{
	pos: position{line: 297, col: 11, offset: 8040},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 297, col: 11, offset: 8040},
	val: "missing",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 297, col: 21, offset: 8050},
	name: "_",
},
	},
},
},
},
{
	name: "ImportType",
	pos: position{line: 299, col: 1, offset: 8086},
	expr: &choiceExpr{
	pos: position{line: 299, col: 14, offset: 8101},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 299, col: 14, offset: 8101},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 299, col: 24, offset: 8111},
	name: "Local",
},
&ruleRefExpr{
	pos: position{line: 299, col: 32, offset: 8119},
	name: "Env",
},
	},
},
},
{
	name: "ImportHashed",
	pos: position{line: 301, col: 1, offset: 8124},
	expr: &actionExpr{
	pos: position{line: 301, col: 16, offset: 8141},
	run: (*parser).callonImportHashed1,
	expr: &labeledExpr{
	pos: position{line: 301, col: 16, offset: 8141},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 301, col: 18, offset: 8143},
	name: "ImportType",
},
},
},
},
{
	name: "Import",
	pos: position{line: 303, col: 1, offset: 8212},
	expr: &choiceExpr{
	pos: position{line: 303, col: 10, offset: 8223},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 303, col: 10, offset: 8223},
	run: (*parser).callonImport2,
	expr: &seqExpr{
	pos: position{line: 303, col: 10, offset: 8223},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 303, col: 10, offset: 8223},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 303, col: 12, offset: 8225},
	name: "ImportHashed",
},
},
&ruleRefExpr{
	pos: position{line: 303, col: 25, offset: 8238},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 303, col: 27, offset: 8240},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 303, col: 30, offset: 8243},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 303, col: 33, offset: 8246},
	name: "Text",
},
	},
},
},
&actionExpr{
	pos: position{line: 304, col: 10, offset: 8343},
	run: (*parser).callonImport10,
	expr: &labeledExpr{
	pos: position{line: 304, col: 10, offset: 8343},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 304, col: 12, offset: 8345},
	name: "ImportHashed",
},
},
},
	},
},
},
{
	name: "LetBinding",
	pos: position{line: 307, col: 1, offset: 8440},
	expr: &actionExpr{
	pos: position{line: 307, col: 14, offset: 8455},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 307, col: 14, offset: 8455},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 307, col: 14, offset: 8455},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 307, col: 18, offset: 8459},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 307, col: 21, offset: 8462},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 307, col: 27, offset: 8468},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 307, col: 44, offset: 8485},
	name: "_",
},
&labeledExpr{
	pos: position{line: 307, col: 46, offset: 8487},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 307, col: 48, offset: 8489},
	expr: &seqExpr{
	pos: position{line: 307, col: 49, offset: 8490},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 307, col: 49, offset: 8490},
	name: "Annotation",
},
&ruleRefExpr{
	pos: position{line: 307, col: 60, offset: 8501},
	name: "_",
},
	},
},
},
},
&litMatcher{
	pos: position{line: 308, col: 13, offset: 8517},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 308, col: 17, offset: 8521},
	name: "_",
},
&labeledExpr{
	pos: position{line: 308, col: 19, offset: 8523},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 308, col: 21, offset: 8525},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 308, col: 32, offset: 8536},
	name: "_",
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 323, col: 1, offset: 8845},
	expr: &choiceExpr{
	pos: position{line: 324, col: 7, offset: 8866},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 324, col: 7, offset: 8866},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 324, col: 7, offset: 8866},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 324, col: 7, offset: 8866},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 324, col: 14, offset: 8873},
	name: "_",
},
&litMatcher{
	pos: position{line: 324, col: 16, offset: 8875},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 324, col: 20, offset: 8879},
	name: "_",
},
&labeledExpr{
	pos: position{line: 324, col: 22, offset: 8881},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 324, col: 28, offset: 8887},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 324, col: 45, offset: 8904},
	name: "_",
},
&litMatcher{
	pos: position{line: 324, col: 47, offset: 8906},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 324, col: 51, offset: 8910},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 324, col: 54, offset: 8913},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 324, col: 56, offset: 8915},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 324, col: 67, offset: 8926},
	name: "_",
},
&litMatcher{
	pos: position{line: 324, col: 69, offset: 8928},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 324, col: 73, offset: 8932},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 324, col: 75, offset: 8934},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 324, col: 81, offset: 8940},
	name: "_",
},
&labeledExpr{
	pos: position{line: 324, col: 83, offset: 8942},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 324, col: 88, offset: 8947},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 327, col: 7, offset: 9064},
	run: (*parser).callonExpression22,
	expr: &seqExpr{
	pos: position{line: 327, col: 7, offset: 9064},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 327, col: 7, offset: 9064},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 327, col: 10, offset: 9067},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 327, col: 13, offset: 9070},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 327, col: 18, offset: 9075},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 327, col: 29, offset: 9086},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 327, col: 31, offset: 9088},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 327, col: 36, offset: 9093},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 327, col: 39, offset: 9096},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 327, col: 41, offset: 9098},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 327, col: 52, offset: 9109},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 327, col: 54, offset: 9111},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 327, col: 59, offset: 9116},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 327, col: 62, offset: 9119},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 327, col: 64, offset: 9121},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 330, col: 7, offset: 9207},
	run: (*parser).callonExpression38,
	expr: &seqExpr{
	pos: position{line: 330, col: 7, offset: 9207},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 330, col: 7, offset: 9207},
	label: "bindings",
	expr: &oneOrMoreExpr{
	pos: position{line: 330, col: 16, offset: 9216},
	expr: &ruleRefExpr{
	pos: position{line: 330, col: 16, offset: 9216},
	name: "LetBinding",
},
},
},
&ruleRefExpr{
	pos: position{line: 330, col: 28, offset: 9228},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 330, col: 31, offset: 9231},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 330, col: 34, offset: 9234},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 330, col: 36, offset: 9236},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 337, col: 7, offset: 9476},
	run: (*parser).callonExpression47,
	expr: &seqExpr{
	pos: position{line: 337, col: 7, offset: 9476},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 337, col: 7, offset: 9476},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 337, col: 14, offset: 9483},
	name: "_",
},
&litMatcher{
	pos: position{line: 337, col: 16, offset: 9485},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 337, col: 20, offset: 9489},
	name: "_",
},
&labeledExpr{
	pos: position{line: 337, col: 22, offset: 9491},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 337, col: 28, offset: 9497},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 337, col: 45, offset: 9514},
	name: "_",
},
&litMatcher{
	pos: position{line: 337, col: 47, offset: 9516},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 337, col: 51, offset: 9520},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 337, col: 54, offset: 9523},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 337, col: 56, offset: 9525},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 337, col: 67, offset: 9536},
	name: "_",
},
&litMatcher{
	pos: position{line: 337, col: 69, offset: 9538},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 337, col: 73, offset: 9542},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 337, col: 75, offset: 9544},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 337, col: 81, offset: 9550},
	name: "_",
},
&labeledExpr{
	pos: position{line: 337, col: 83, offset: 9552},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 337, col: 88, offset: 9557},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 340, col: 7, offset: 9666},
	run: (*parser).callonExpression67,
	expr: &seqExpr{
	pos: position{line: 340, col: 7, offset: 9666},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 340, col: 7, offset: 9666},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 340, col: 9, offset: 9668},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 340, col: 28, offset: 9687},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 340, col: 30, offset: 9689},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 340, col: 36, offset: 9695},
	name: "_",
},
&labeledExpr{
	pos: position{line: 340, col: 38, offset: 9697},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 340, col: 40, offset: 9699},
	name: "Expression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 341, col: 7, offset: 9759},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 343, col: 1, offset: 9780},
	expr: &actionExpr{
	pos: position{line: 343, col: 14, offset: 9795},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 343, col: 14, offset: 9795},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 343, col: 14, offset: 9795},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 343, col: 18, offset: 9799},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 343, col: 21, offset: 9802},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 343, col: 23, offset: 9804},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 345, col: 1, offset: 9834},
	expr: &choiceExpr{
	pos: position{line: 346, col: 5, offset: 9862},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 346, col: 5, offset: 9862},
	name: "EmptyList",
},
&actionExpr{
	pos: position{line: 347, col: 5, offset: 9876},
	run: (*parser).callonAnnotatedExpression3,
	expr: &seqExpr{
	pos: position{line: 347, col: 5, offset: 9876},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 347, col: 5, offset: 9876},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 347, col: 7, offset: 9878},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 347, col: 26, offset: 9897},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 347, col: 28, offset: 9899},
	expr: &seqExpr{
	pos: position{line: 347, col: 29, offset: 9900},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 347, col: 29, offset: 9900},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 347, col: 31, offset: 9902},
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
	pos: position{line: 352, col: 1, offset: 10027},
	expr: &actionExpr{
	pos: position{line: 352, col: 13, offset: 10041},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 352, col: 13, offset: 10041},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 352, col: 13, offset: 10041},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 352, col: 17, offset: 10045},
	name: "_",
},
&litMatcher{
	pos: position{line: 352, col: 19, offset: 10047},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 352, col: 23, offset: 10051},
	name: "_",
},
&litMatcher{
	pos: position{line: 352, col: 25, offset: 10053},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 352, col: 29, offset: 10057},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 352, col: 32, offset: 10060},
	name: "List",
},
&ruleRefExpr{
	pos: position{line: 352, col: 37, offset: 10065},
	name: "_",
},
&labeledExpr{
	pos: position{line: 352, col: 39, offset: 10067},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 352, col: 41, offset: 10069},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 356, col: 1, offset: 10132},
	expr: &ruleRefExpr{
	pos: position{line: 356, col: 22, offset: 10155},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 358, col: 1, offset: 10176},
	expr: &ruleRefExpr{
	pos: position{line: 358, col: 23, offset: 10200},
	name: "PlusExpression",
},
},
{
	name: "MorePlus",
	pos: position{line: 360, col: 1, offset: 10216},
	expr: &actionExpr{
	pos: position{line: 360, col: 12, offset: 10229},
	run: (*parser).callonMorePlus1,
	expr: &seqExpr{
	pos: position{line: 360, col: 12, offset: 10229},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 360, col: 12, offset: 10229},
	name: "_",
},
&litMatcher{
	pos: position{line: 360, col: 14, offset: 10231},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 360, col: 18, offset: 10235},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 360, col: 21, offset: 10238},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 360, col: 23, offset: 10240},
	name: "TimesExpression",
},
},
	},
},
},
},
{
	name: "PlusExpression",
	pos: position{line: 361, col: 1, offset: 10274},
	expr: &actionExpr{
	pos: position{line: 362, col: 7, offset: 10299},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 362, col: 7, offset: 10299},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 362, col: 7, offset: 10299},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 362, col: 13, offset: 10305},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 362, col: 29, offset: 10321},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 362, col: 34, offset: 10326},
	expr: &ruleRefExpr{
	pos: position{line: 362, col: 34, offset: 10326},
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
	pos: position{line: 371, col: 1, offset: 10554},
	expr: &actionExpr{
	pos: position{line: 371, col: 13, offset: 10568},
	run: (*parser).callonMoreTimes1,
	expr: &seqExpr{
	pos: position{line: 371, col: 13, offset: 10568},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 371, col: 13, offset: 10568},
	name: "_",
},
&litMatcher{
	pos: position{line: 371, col: 15, offset: 10570},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 371, col: 19, offset: 10574},
	name: "_",
},
&labeledExpr{
	pos: position{line: 371, col: 21, offset: 10576},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 371, col: 23, offset: 10578},
	name: "ApplicationExpression",
},
},
	},
},
},
},
{
	name: "TimesExpression",
	pos: position{line: 372, col: 1, offset: 10618},
	expr: &actionExpr{
	pos: position{line: 373, col: 7, offset: 10644},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 373, col: 7, offset: 10644},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 373, col: 7, offset: 10644},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 373, col: 13, offset: 10650},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 373, col: 35, offset: 10672},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 373, col: 40, offset: 10677},
	expr: &ruleRefExpr{
	pos: position{line: 373, col: 40, offset: 10677},
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
	pos: position{line: 382, col: 1, offset: 10907},
	expr: &actionExpr{
	pos: position{line: 382, col: 25, offset: 10933},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 382, col: 25, offset: 10933},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 382, col: 25, offset: 10933},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 382, col: 27, offset: 10935},
	name: "FirstApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 382, col: 54, offset: 10962},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 382, col: 59, offset: 10967},
	expr: &seqExpr{
	pos: position{line: 382, col: 60, offset: 10968},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 382, col: 60, offset: 10968},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 382, col: 63, offset: 10971},
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
	pos: position{line: 391, col: 1, offset: 11221},
	expr: &choiceExpr{
	pos: position{line: 392, col: 8, offset: 11259},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 392, col: 8, offset: 11259},
	run: (*parser).callonFirstApplicationExpression2,
	expr: &seqExpr{
	pos: position{line: 392, col: 8, offset: 11259},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 392, col: 8, offset: 11259},
	name: "Some",
},
&ruleRefExpr{
	pos: position{line: 392, col: 13, offset: 11264},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 392, col: 16, offset: 11267},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 392, col: 18, offset: 11269},
	name: "ImportExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 393, col: 8, offset: 11324},
	name: "ImportExpression",
},
	},
},
},
{
	name: "ImportExpression",
	pos: position{line: 395, col: 1, offset: 11342},
	expr: &choiceExpr{
	pos: position{line: 395, col: 20, offset: 11363},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 395, col: 20, offset: 11363},
	name: "Import",
},
&ruleRefExpr{
	pos: position{line: 395, col: 29, offset: 11372},
	name: "SelectorExpression",
},
	},
},
},
{
	name: "SelectorExpression",
	pos: position{line: 397, col: 1, offset: 11392},
	expr: &actionExpr{
	pos: position{line: 397, col: 22, offset: 11415},
	run: (*parser).callonSelectorExpression1,
	expr: &seqExpr{
	pos: position{line: 397, col: 22, offset: 11415},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 397, col: 22, offset: 11415},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 397, col: 24, offset: 11417},
	name: "PrimitiveExpression",
},
},
&labeledExpr{
	pos: position{line: 397, col: 44, offset: 11437},
	label: "ls",
	expr: &zeroOrMoreExpr{
	pos: position{line: 397, col: 47, offset: 11440},
	expr: &seqExpr{
	pos: position{line: 397, col: 48, offset: 11441},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 397, col: 48, offset: 11441},
	name: "_",
},
&litMatcher{
	pos: position{line: 397, col: 50, offset: 11443},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 397, col: 54, offset: 11447},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 397, col: 56, offset: 11449},
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
	pos: position{line: 407, col: 1, offset: 11682},
	expr: &choiceExpr{
	pos: position{line: 408, col: 7, offset: 11712},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 408, col: 7, offset: 11712},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 409, col: 7, offset: 11732},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 410, col: 7, offset: 11753},
	name: "IntegerLiteral",
},
&ruleRefExpr{
	pos: position{line: 411, col: 7, offset: 11774},
	name: "TextLiteral",
},
&actionExpr{
	pos: position{line: 412, col: 7, offset: 11792},
	run: (*parser).callonPrimitiveExpression6,
	expr: &seqExpr{
	pos: position{line: 412, col: 7, offset: 11792},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 412, col: 7, offset: 11792},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 412, col: 11, offset: 11796},
	name: "_",
},
&labeledExpr{
	pos: position{line: 412, col: 13, offset: 11798},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 412, col: 15, offset: 11800},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 412, col: 35, offset: 11820},
	name: "_",
},
&litMatcher{
	pos: position{line: 412, col: 37, offset: 11822},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 413, col: 7, offset: 11850},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 414, col: 7, offset: 11876},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 415, col: 7, offset: 11893},
	run: (*parser).callonPrimitiveExpression16,
	expr: &seqExpr{
	pos: position{line: 415, col: 7, offset: 11893},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 415, col: 7, offset: 11893},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 415, col: 11, offset: 11897},
	name: "_",
},
&labeledExpr{
	pos: position{line: 415, col: 14, offset: 11900},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 415, col: 16, offset: 11902},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 415, col: 27, offset: 11913},
	name: "_",
},
&litMatcher{
	pos: position{line: 415, col: 29, offset: 11915},
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
	pos: position{line: 417, col: 1, offset: 11938},
	expr: &choiceExpr{
	pos: position{line: 418, col: 7, offset: 11968},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 418, col: 7, offset: 11968},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &litMatcher{
	pos: position{line: 418, col: 7, offset: 11968},
	val: "=",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 419, col: 7, offset: 12023},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 420, col: 7, offset: 12048},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 421, col: 7, offset: 12076},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 421, col: 7, offset: 12076},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 423, col: 1, offset: 12122},
	expr: &actionExpr{
	pos: position{line: 423, col: 19, offset: 12142},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 423, col: 19, offset: 12142},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 423, col: 19, offset: 12142},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 423, col: 24, offset: 12147},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 423, col: 33, offset: 12156},
	name: "_",
},
&litMatcher{
	pos: position{line: 423, col: 35, offset: 12158},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 423, col: 39, offset: 12162},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 423, col: 42, offset: 12165},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 423, col: 47, offset: 12170},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 426, col: 1, offset: 12227},
	expr: &actionExpr{
	pos: position{line: 426, col: 18, offset: 12246},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 426, col: 18, offset: 12246},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 426, col: 18, offset: 12246},
	name: "_",
},
&litMatcher{
	pos: position{line: 426, col: 20, offset: 12248},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 426, col: 24, offset: 12252},
	name: "_",
},
&labeledExpr{
	pos: position{line: 426, col: 26, offset: 12254},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 426, col: 28, offset: 12256},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 427, col: 1, offset: 12288},
	expr: &actionExpr{
	pos: position{line: 428, col: 7, offset: 12317},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 428, col: 7, offset: 12317},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 428, col: 7, offset: 12317},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 428, col: 13, offset: 12323},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 428, col: 29, offset: 12339},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 428, col: 34, offset: 12344},
	expr: &ruleRefExpr{
	pos: position{line: 428, col: 34, offset: 12344},
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
	pos: position{line: 438, col: 1, offset: 12740},
	expr: &actionExpr{
	pos: position{line: 438, col: 22, offset: 12763},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 438, col: 22, offset: 12763},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 438, col: 22, offset: 12763},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 438, col: 27, offset: 12768},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 438, col: 36, offset: 12777},
	name: "_",
},
&litMatcher{
	pos: position{line: 438, col: 38, offset: 12779},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 438, col: 42, offset: 12783},
	name: "_",
},
&labeledExpr{
	pos: position{line: 438, col: 44, offset: 12785},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 438, col: 49, offset: 12790},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 441, col: 1, offset: 12847},
	expr: &actionExpr{
	pos: position{line: 441, col: 21, offset: 12869},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 441, col: 21, offset: 12869},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 441, col: 21, offset: 12869},
	name: "_",
},
&litMatcher{
	pos: position{line: 441, col: 23, offset: 12871},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 441, col: 27, offset: 12875},
	name: "_",
},
&labeledExpr{
	pos: position{line: 441, col: 29, offset: 12877},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 441, col: 31, offset: 12879},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 442, col: 1, offset: 12914},
	expr: &actionExpr{
	pos: position{line: 443, col: 7, offset: 12946},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 443, col: 7, offset: 12946},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 443, col: 7, offset: 12946},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 443, col: 13, offset: 12952},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 443, col: 32, offset: 12971},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 443, col: 37, offset: 12976},
	expr: &ruleRefExpr{
	pos: position{line: 443, col: 37, offset: 12976},
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
	pos: position{line: 453, col: 1, offset: 13378},
	expr: &actionExpr{
	pos: position{line: 453, col: 12, offset: 13391},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 453, col: 12, offset: 13391},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 453, col: 12, offset: 13391},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 453, col: 16, offset: 13395},
	name: "_",
},
&labeledExpr{
	pos: position{line: 453, col: 18, offset: 13397},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 453, col: 20, offset: 13399},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 453, col: 31, offset: 13410},
	name: "_",
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 455, col: 1, offset: 13429},
	expr: &actionExpr{
	pos: position{line: 456, col: 7, offset: 13459},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 456, col: 7, offset: 13459},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 456, col: 7, offset: 13459},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 456, col: 11, offset: 13463},
	name: "_",
},
&labeledExpr{
	pos: position{line: 456, col: 13, offset: 13465},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 456, col: 19, offset: 13471},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 456, col: 30, offset: 13482},
	name: "_",
},
&labeledExpr{
	pos: position{line: 456, col: 32, offset: 13484},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 456, col: 37, offset: 13489},
	expr: &ruleRefExpr{
	pos: position{line: 456, col: 37, offset: 13489},
	name: "MoreList",
},
},
},
&litMatcher{
	pos: position{line: 456, col: 47, offset: 13499},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 466, col: 1, offset: 13775},
	expr: &notExpr{
	pos: position{line: 466, col: 7, offset: 13783},
	expr: &anyMatcher{
	line: 466, col: 8, offset: 13784,
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

func (c *current) onEscapedChar1() (interface{}, error) {
    switch c.text[1] {
    case 'b':
        return []byte("\b"), nil
    case 'f':
        return []byte("\f"), nil
    case 'n':
        return []byte("\n"), nil
    case 'r':
        return []byte("\r"), nil
    case 't':
        return []byte("\t"), nil
    case 'u':
        i, err := strconv.ParseInt(string(c.text[2:]), 16, 32)
        return []byte(string([]rune{rune(i)})), err
    }
    return c.text[1:2], nil
}

func (p *parser) callonEscapedChar1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEscapedChar1()
}

func (c *current) onDoubleQuoteChunk2(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonDoubleQuoteChunk2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleQuoteChunk2(stack["e"])
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

