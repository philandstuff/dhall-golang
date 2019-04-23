
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
	name: "SingleQuoteContinue",
	pos: position{line: 109, col: 1, offset: 2445},
	expr: &choiceExpr{
	pos: position{line: 110, col: 7, offset: 2475},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 110, col: 7, offset: 2475},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 110, col: 7, offset: 2475},
	name: "Interpolation",
},
&ruleRefExpr{
	pos: position{line: 110, col: 21, offset: 2489},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 111, col: 7, offset: 2515},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 111, col: 7, offset: 2515},
	name: "EscapedQuotePair",
},
&ruleRefExpr{
	pos: position{line: 111, col: 24, offset: 2532},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 112, col: 7, offset: 2558},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 112, col: 7, offset: 2558},
	name: "EscapedInterpolation",
},
&ruleRefExpr{
	pos: position{line: 112, col: 28, offset: 2579},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 113, col: 7, offset: 2605},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 113, col: 7, offset: 2605},
	name: "SingleQuoteChar",
},
&ruleRefExpr{
	pos: position{line: 113, col: 23, offset: 2621},
	name: "SingleQuoteContinue",
},
	},
},
&litMatcher{
	pos: position{line: 114, col: 7, offset: 2647},
	val: "''",
	ignoreCase: false,
},
	},
},
},
{
	name: "EscapedQuotePair",
	pos: position{line: 116, col: 1, offset: 2653},
	expr: &actionExpr{
	pos: position{line: 116, col: 20, offset: 2674},
	run: (*parser).callonEscapedQuotePair1,
	expr: &litMatcher{
	pos: position{line: 116, col: 20, offset: 2674},
	val: "'''",
	ignoreCase: false,
},
},
},
{
	name: "EscapedInterpolation",
	pos: position{line: 120, col: 1, offset: 2809},
	expr: &actionExpr{
	pos: position{line: 120, col: 24, offset: 2834},
	run: (*parser).callonEscapedInterpolation1,
	expr: &litMatcher{
	pos: position{line: 120, col: 24, offset: 2834},
	val: "''${",
	ignoreCase: false,
},
},
},
{
	name: "SingleQuoteChar",
	pos: position{line: 122, col: 1, offset: 2876},
	expr: &choiceExpr{
	pos: position{line: 123, col: 6, offset: 2901},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 123, col: 6, offset: 2901},
	val: "[\\x20-\\U0010ffff]",
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 124, col: 6, offset: 2924},
	val: "\t",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 125, col: 6, offset: 2934},
	val: "\n",
	ignoreCase: false,
},
	},
},
},
{
	name: "SingleQuoteLiteral",
	pos: position{line: 127, col: 1, offset: 2940},
	expr: &actionExpr{
	pos: position{line: 127, col: 22, offset: 2963},
	run: (*parser).callonSingleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 127, col: 22, offset: 2963},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 127, col: 22, offset: 2963},
	val: "''",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 127, col: 27, offset: 2968},
	name: "EOL",
},
&labeledExpr{
	pos: position{line: 127, col: 31, offset: 2972},
	label: "content",
	expr: &ruleRefExpr{
	pos: position{line: 127, col: 39, offset: 2980},
	name: "SingleQuoteContinue",
},
},
	},
},
},
},
{
	name: "Interpolation",
	pos: position{line: 145, col: 1, offset: 3503},
	expr: &actionExpr{
	pos: position{line: 145, col: 17, offset: 3521},
	run: (*parser).callonInterpolation1,
	expr: &seqExpr{
	pos: position{line: 145, col: 17, offset: 3521},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 145, col: 17, offset: 3521},
	val: "${",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 145, col: 22, offset: 3526},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 145, col: 24, offset: 3528},
	name: "CompleteExpression",
},
},
&litMatcher{
	pos: position{line: 145, col: 43, offset: 3547},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TextLiteral",
	pos: position{line: 147, col: 1, offset: 3570},
	expr: &choiceExpr{
	pos: position{line: 147, col: 15, offset: 3586},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 147, col: 15, offset: 3586},
	name: "DoubleQuoteLiteral",
},
&ruleRefExpr{
	pos: position{line: 147, col: 36, offset: 3607},
	name: "SingleQuoteLiteral",
},
	},
},
},
{
	name: "Reserved",
	pos: position{line: 150, col: 1, offset: 3712},
	expr: &choiceExpr{
	pos: position{line: 151, col: 5, offset: 3729},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 151, col: 5, offset: 3729},
	run: (*parser).callonReserved2,
	expr: &litMatcher{
	pos: position{line: 151, col: 5, offset: 3729},
	val: "Natural/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 152, col: 5, offset: 3807},
	run: (*parser).callonReserved4,
	expr: &litMatcher{
	pos: position{line: 152, col: 5, offset: 3807},
	val: "Natural/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 153, col: 5, offset: 3883},
	run: (*parser).callonReserved6,
	expr: &litMatcher{
	pos: position{line: 153, col: 5, offset: 3883},
	val: "Natural/isZero",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 154, col: 5, offset: 3963},
	run: (*parser).callonReserved8,
	expr: &litMatcher{
	pos: position{line: 154, col: 5, offset: 3963},
	val: "Natural/even",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 155, col: 5, offset: 4039},
	run: (*parser).callonReserved10,
	expr: &litMatcher{
	pos: position{line: 155, col: 5, offset: 4039},
	val: "Natural/odd",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 156, col: 5, offset: 4113},
	run: (*parser).callonReserved12,
	expr: &litMatcher{
	pos: position{line: 156, col: 5, offset: 4113},
	val: "Natural/toInteger",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 157, col: 5, offset: 4199},
	run: (*parser).callonReserved14,
	expr: &litMatcher{
	pos: position{line: 157, col: 5, offset: 4199},
	val: "Natural/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 158, col: 5, offset: 4275},
	run: (*parser).callonReserved16,
	expr: &litMatcher{
	pos: position{line: 158, col: 5, offset: 4275},
	val: "Integer/toDouble",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 159, col: 5, offset: 4359},
	run: (*parser).callonReserved18,
	expr: &litMatcher{
	pos: position{line: 159, col: 5, offset: 4359},
	val: "Integer/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 160, col: 5, offset: 4435},
	run: (*parser).callonReserved20,
	expr: &litMatcher{
	pos: position{line: 160, col: 5, offset: 4435},
	val: "Double/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 161, col: 5, offset: 4509},
	run: (*parser).callonReserved22,
	expr: &litMatcher{
	pos: position{line: 161, col: 5, offset: 4509},
	val: "List/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 162, col: 5, offset: 4581},
	run: (*parser).callonReserved24,
	expr: &litMatcher{
	pos: position{line: 162, col: 5, offset: 4581},
	val: "List/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 163, col: 5, offset: 4651},
	run: (*parser).callonReserved26,
	expr: &litMatcher{
	pos: position{line: 163, col: 5, offset: 4651},
	val: "List/length",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 164, col: 5, offset: 4725},
	run: (*parser).callonReserved28,
	expr: &litMatcher{
	pos: position{line: 164, col: 5, offset: 4725},
	val: "List/head",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 165, col: 5, offset: 4795},
	run: (*parser).callonReserved30,
	expr: &litMatcher{
	pos: position{line: 165, col: 5, offset: 4795},
	val: "List/last",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 166, col: 5, offset: 4865},
	run: (*parser).callonReserved32,
	expr: &litMatcher{
	pos: position{line: 166, col: 5, offset: 4865},
	val: "List/indexed",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 167, col: 5, offset: 4941},
	run: (*parser).callonReserved34,
	expr: &litMatcher{
	pos: position{line: 167, col: 5, offset: 4941},
	val: "List/reverse",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 168, col: 5, offset: 5017},
	run: (*parser).callonReserved36,
	expr: &litMatcher{
	pos: position{line: 168, col: 5, offset: 5017},
	val: "Optional/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 169, col: 5, offset: 5097},
	run: (*parser).callonReserved38,
	expr: &litMatcher{
	pos: position{line: 169, col: 5, offset: 5097},
	val: "Optional/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 170, col: 5, offset: 5175},
	run: (*parser).callonReserved40,
	expr: &litMatcher{
	pos: position{line: 170, col: 5, offset: 5175},
	val: "Text/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 171, col: 5, offset: 5245},
	run: (*parser).callonReserved42,
	expr: &litMatcher{
	pos: position{line: 171, col: 5, offset: 5245},
	val: "Bool",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 172, col: 5, offset: 5277},
	run: (*parser).callonReserved44,
	expr: &litMatcher{
	pos: position{line: 172, col: 5, offset: 5277},
	val: "True",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 173, col: 5, offset: 5309},
	run: (*parser).callonReserved46,
	expr: &litMatcher{
	pos: position{line: 173, col: 5, offset: 5309},
	val: "False",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 174, col: 5, offset: 5343},
	run: (*parser).callonReserved48,
	expr: &litMatcher{
	pos: position{line: 174, col: 5, offset: 5343},
	val: "Optional",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 175, col: 5, offset: 5383},
	run: (*parser).callonReserved50,
	expr: &litMatcher{
	pos: position{line: 175, col: 5, offset: 5383},
	val: "Natural",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 176, col: 5, offset: 5421},
	run: (*parser).callonReserved52,
	expr: &litMatcher{
	pos: position{line: 176, col: 5, offset: 5421},
	val: "Integer",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 177, col: 5, offset: 5459},
	run: (*parser).callonReserved54,
	expr: &litMatcher{
	pos: position{line: 177, col: 5, offset: 5459},
	val: "Double",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 178, col: 5, offset: 5495},
	run: (*parser).callonReserved56,
	expr: &litMatcher{
	pos: position{line: 178, col: 5, offset: 5495},
	val: "Text",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 179, col: 5, offset: 5527},
	run: (*parser).callonReserved58,
	expr: &litMatcher{
	pos: position{line: 179, col: 5, offset: 5527},
	val: "List",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 180, col: 5, offset: 5559},
	run: (*parser).callonReserved60,
	expr: &litMatcher{
	pos: position{line: 180, col: 5, offset: 5559},
	val: "None",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 181, col: 5, offset: 5591},
	run: (*parser).callonReserved62,
	expr: &litMatcher{
	pos: position{line: 181, col: 5, offset: 5591},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 182, col: 5, offset: 5623},
	run: (*parser).callonReserved64,
	expr: &litMatcher{
	pos: position{line: 182, col: 5, offset: 5623},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 183, col: 5, offset: 5655},
	run: (*parser).callonReserved66,
	expr: &litMatcher{
	pos: position{line: 183, col: 5, offset: 5655},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "If",
	pos: position{line: 185, col: 1, offset: 5684},
	expr: &litMatcher{
	pos: position{line: 185, col: 6, offset: 5691},
	val: "if",
	ignoreCase: false,
},
},
{
	name: "Then",
	pos: position{line: 186, col: 1, offset: 5696},
	expr: &litMatcher{
	pos: position{line: 186, col: 8, offset: 5705},
	val: "then",
	ignoreCase: false,
},
},
{
	name: "Else",
	pos: position{line: 187, col: 1, offset: 5712},
	expr: &litMatcher{
	pos: position{line: 187, col: 8, offset: 5721},
	val: "else",
	ignoreCase: false,
},
},
{
	name: "Let",
	pos: position{line: 188, col: 1, offset: 5728},
	expr: &litMatcher{
	pos: position{line: 188, col: 7, offset: 5736},
	val: "let",
	ignoreCase: false,
},
},
{
	name: "In",
	pos: position{line: 189, col: 1, offset: 5742},
	expr: &litMatcher{
	pos: position{line: 189, col: 6, offset: 5749},
	val: "in",
	ignoreCase: false,
},
},
{
	name: "As",
	pos: position{line: 190, col: 1, offset: 5754},
	expr: &litMatcher{
	pos: position{line: 190, col: 6, offset: 5761},
	val: "as",
	ignoreCase: false,
},
},
{
	name: "Using",
	pos: position{line: 191, col: 1, offset: 5766},
	expr: &litMatcher{
	pos: position{line: 191, col: 9, offset: 5776},
	val: "using",
	ignoreCase: false,
},
},
{
	name: "Merge",
	pos: position{line: 192, col: 1, offset: 5784},
	expr: &litMatcher{
	pos: position{line: 192, col: 9, offset: 5794},
	val: "merge",
	ignoreCase: false,
},
},
{
	name: "Missing",
	pos: position{line: 193, col: 1, offset: 5802},
	expr: &actionExpr{
	pos: position{line: 193, col: 11, offset: 5814},
	run: (*parser).callonMissing1,
	expr: &litMatcher{
	pos: position{line: 193, col: 11, offset: 5814},
	val: "missing",
	ignoreCase: false,
},
},
},
{
	name: "True",
	pos: position{line: 194, col: 1, offset: 5850},
	expr: &litMatcher{
	pos: position{line: 194, col: 8, offset: 5859},
	val: "True",
	ignoreCase: false,
},
},
{
	name: "False",
	pos: position{line: 195, col: 1, offset: 5866},
	expr: &litMatcher{
	pos: position{line: 195, col: 9, offset: 5876},
	val: "False",
	ignoreCase: false,
},
},
{
	name: "Infinity",
	pos: position{line: 196, col: 1, offset: 5884},
	expr: &litMatcher{
	pos: position{line: 196, col: 12, offset: 5897},
	val: "Infinity",
	ignoreCase: false,
},
},
{
	name: "NaN",
	pos: position{line: 197, col: 1, offset: 5908},
	expr: &litMatcher{
	pos: position{line: 197, col: 7, offset: 5916},
	val: "NaN",
	ignoreCase: false,
},
},
{
	name: "Some",
	pos: position{line: 198, col: 1, offset: 5922},
	expr: &litMatcher{
	pos: position{line: 198, col: 8, offset: 5931},
	val: "Some",
	ignoreCase: false,
},
},
{
	name: "Keyword",
	pos: position{line: 200, col: 1, offset: 5939},
	expr: &choiceExpr{
	pos: position{line: 201, col: 5, offset: 5955},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 201, col: 5, offset: 5955},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 201, col: 10, offset: 5960},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 201, col: 17, offset: 5967},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 202, col: 5, offset: 5976},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 202, col: 11, offset: 5982},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 203, col: 5, offset: 5989},
	name: "Using",
},
&ruleRefExpr{
	pos: position{line: 203, col: 13, offset: 5997},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 203, col: 23, offset: 6007},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 204, col: 5, offset: 6014},
	name: "True",
},
&ruleRefExpr{
	pos: position{line: 204, col: 12, offset: 6021},
	name: "False",
},
&ruleRefExpr{
	pos: position{line: 205, col: 5, offset: 6031},
	name: "Infinity",
},
&ruleRefExpr{
	pos: position{line: 205, col: 16, offset: 6042},
	name: "NaN",
},
&ruleRefExpr{
	pos: position{line: 206, col: 5, offset: 6050},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 206, col: 13, offset: 6058},
	name: "Some",
},
	},
},
},
{
	name: "Optional",
	pos: position{line: 208, col: 1, offset: 6064},
	expr: &litMatcher{
	pos: position{line: 208, col: 12, offset: 6077},
	val: "Optional",
	ignoreCase: false,
},
},
{
	name: "Text",
	pos: position{line: 209, col: 1, offset: 6088},
	expr: &litMatcher{
	pos: position{line: 209, col: 8, offset: 6097},
	val: "Text",
	ignoreCase: false,
},
},
{
	name: "List",
	pos: position{line: 210, col: 1, offset: 6104},
	expr: &litMatcher{
	pos: position{line: 210, col: 8, offset: 6113},
	val: "List",
	ignoreCase: false,
},
},
{
	name: "Lambda",
	pos: position{line: 212, col: 1, offset: 6121},
	expr: &choiceExpr{
	pos: position{line: 212, col: 11, offset: 6133},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 212, col: 11, offset: 6133},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 212, col: 18, offset: 6140},
	val: "λ",
	ignoreCase: false,
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 213, col: 1, offset: 6146},
	expr: &choiceExpr{
	pos: position{line: 213, col: 11, offset: 6158},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 213, col: 11, offset: 6158},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 213, col: 22, offset: 6169},
	val: "∀",
	ignoreCase: false,
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 214, col: 1, offset: 6176},
	expr: &choiceExpr{
	pos: position{line: 214, col: 10, offset: 6187},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 214, col: 10, offset: 6187},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 214, col: 17, offset: 6194},
	val: "→",
	ignoreCase: false,
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 216, col: 1, offset: 6202},
	expr: &seqExpr{
	pos: position{line: 216, col: 12, offset: 6215},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 216, col: 12, offset: 6215},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 216, col: 17, offset: 6220},
	expr: &charClassMatcher{
	pos: position{line: 216, col: 17, offset: 6220},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 216, col: 23, offset: 6226},
	expr: &charClassMatcher{
	pos: position{line: 216, col: 23, offset: 6226},
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
	pos: position{line: 218, col: 1, offset: 6234},
	expr: &actionExpr{
	pos: position{line: 218, col: 24, offset: 6259},
	run: (*parser).callonNumericDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 218, col: 24, offset: 6259},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 218, col: 24, offset: 6259},
	expr: &charClassMatcher{
	pos: position{line: 218, col: 24, offset: 6259},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 218, col: 30, offset: 6265},
	expr: &charClassMatcher{
	pos: position{line: 218, col: 30, offset: 6265},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&choiceExpr{
	pos: position{line: 218, col: 39, offset: 6274},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 218, col: 39, offset: 6274},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 218, col: 39, offset: 6274},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 218, col: 43, offset: 6278},
	expr: &charClassMatcher{
	pos: position{line: 218, col: 43, offset: 6278},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&zeroOrOneExpr{
	pos: position{line: 218, col: 50, offset: 6285},
	expr: &ruleRefExpr{
	pos: position{line: 218, col: 50, offset: 6285},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 218, col: 62, offset: 6297},
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
	pos: position{line: 226, col: 1, offset: 6453},
	expr: &choiceExpr{
	pos: position{line: 226, col: 17, offset: 6471},
	alternatives: []interface{}{
&labeledExpr{
	pos: position{line: 226, col: 17, offset: 6471},
	label: "d",
	expr: &ruleRefExpr{
	pos: position{line: 226, col: 19, offset: 6473},
	name: "NumericDoubleLiteral",
},
},
&actionExpr{
	pos: position{line: 227, col: 5, offset: 6498},
	run: (*parser).callonDoubleLiteral4,
	expr: &ruleRefExpr{
	pos: position{line: 227, col: 5, offset: 6498},
	name: "Infinity",
},
},
&actionExpr{
	pos: position{line: 228, col: 5, offset: 6550},
	run: (*parser).callonDoubleLiteral6,
	expr: &seqExpr{
	pos: position{line: 228, col: 5, offset: 6550},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 228, col: 5, offset: 6550},
	val: "-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 228, col: 9, offset: 6554},
	name: "Infinity",
},
	},
},
},
&actionExpr{
	pos: position{line: 229, col: 5, offset: 6607},
	run: (*parser).callonDoubleLiteral10,
	expr: &ruleRefExpr{
	pos: position{line: 229, col: 5, offset: 6607},
	name: "NaN",
},
},
	},
},
},
{
	name: "NaturalLiteral",
	pos: position{line: 231, col: 1, offset: 6650},
	expr: &actionExpr{
	pos: position{line: 231, col: 18, offset: 6669},
	run: (*parser).callonNaturalLiteral1,
	expr: &oneOrMoreExpr{
	pos: position{line: 231, col: 18, offset: 6669},
	expr: &charClassMatcher{
	pos: position{line: 231, col: 18, offset: 6669},
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
	pos: position{line: 236, col: 1, offset: 6758},
	expr: &actionExpr{
	pos: position{line: 236, col: 18, offset: 6777},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 236, col: 18, offset: 6777},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 236, col: 18, offset: 6777},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&oneOrMoreExpr{
	pos: position{line: 236, col: 22, offset: 6781},
	expr: &charClassMatcher{
	pos: position{line: 236, col: 22, offset: 6781},
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
	pos: position{line: 244, col: 1, offset: 6925},
	expr: &actionExpr{
	pos: position{line: 244, col: 12, offset: 6938},
	run: (*parser).callonDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 244, col: 12, offset: 6938},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 244, col: 12, offset: 6938},
	name: "_",
},
&litMatcher{
	pos: position{line: 244, col: 14, offset: 6940},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 244, col: 18, offset: 6944},
	name: "_",
},
&labeledExpr{
	pos: position{line: 244, col: 20, offset: 6946},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 244, col: 26, offset: 6952},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Variable",
	pos: position{line: 246, col: 1, offset: 7008},
	expr: &actionExpr{
	pos: position{line: 246, col: 12, offset: 7021},
	run: (*parser).callonVariable1,
	expr: &seqExpr{
	pos: position{line: 246, col: 12, offset: 7021},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 246, col: 12, offset: 7021},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 246, col: 17, offset: 7026},
	name: "NonreservedLabel",
},
},
&labeledExpr{
	pos: position{line: 246, col: 34, offset: 7043},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 246, col: 40, offset: 7049},
	expr: &ruleRefExpr{
	pos: position{line: 246, col: 40, offset: 7049},
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
	pos: position{line: 254, col: 1, offset: 7212},
	expr: &choiceExpr{
	pos: position{line: 254, col: 14, offset: 7227},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 254, col: 14, offset: 7227},
	name: "Variable",
},
&ruleRefExpr{
	pos: position{line: 254, col: 25, offset: 7238},
	name: "Reserved",
},
	},
},
},
{
	name: "PathCharacter",
	pos: position{line: 256, col: 1, offset: 7248},
	expr: &choiceExpr{
	pos: position{line: 257, col: 6, offset: 7271},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 257, col: 6, offset: 7271},
	val: "!",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 258, col: 6, offset: 7283},
	val: "[\\x24-\\x27]",
	ranges: []rune{'$','\'',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 259, col: 6, offset: 7300},
	val: "[\\x2a-\\x2b]",
	ranges: []rune{'*','+',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 260, col: 6, offset: 7317},
	val: "[\\x2d-\\x2e]",
	ranges: []rune{'-','.',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 261, col: 6, offset: 7334},
	val: "[\\x30-\\x3b]",
	ranges: []rune{'0',';',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 262, col: 6, offset: 7351},
	val: "=",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 263, col: 6, offset: 7363},
	val: "[\\x40-\\x5a]",
	ranges: []rune{'@','Z',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 264, col: 6, offset: 7380},
	val: "[\\x5e-\\x7a]",
	ranges: []rune{'^','z',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 265, col: 6, offset: 7397},
	val: "|",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 266, col: 6, offset: 7409},
	val: "~",
	ignoreCase: false,
},
	},
},
},
{
	name: "UnquotedPathComponent",
	pos: position{line: 268, col: 1, offset: 7417},
	expr: &actionExpr{
	pos: position{line: 268, col: 25, offset: 7443},
	run: (*parser).callonUnquotedPathComponent1,
	expr: &oneOrMoreExpr{
	pos: position{line: 268, col: 25, offset: 7443},
	expr: &ruleRefExpr{
	pos: position{line: 268, col: 25, offset: 7443},
	name: "PathCharacter",
},
},
},
},
{
	name: "PathComponent",
	pos: position{line: 270, col: 1, offset: 7490},
	expr: &actionExpr{
	pos: position{line: 270, col: 17, offset: 7508},
	run: (*parser).callonPathComponent1,
	expr: &seqExpr{
	pos: position{line: 270, col: 17, offset: 7508},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 270, col: 17, offset: 7508},
	val: "/",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 270, col: 21, offset: 7512},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 270, col: 23, offset: 7514},
	name: "UnquotedPathComponent",
},
},
	},
},
},
},
{
	name: "Path",
	pos: position{line: 272, col: 1, offset: 7555},
	expr: &actionExpr{
	pos: position{line: 272, col: 8, offset: 7564},
	run: (*parser).callonPath1,
	expr: &labeledExpr{
	pos: position{line: 272, col: 8, offset: 7564},
	label: "cs",
	expr: &oneOrMoreExpr{
	pos: position{line: 272, col: 11, offset: 7567},
	expr: &ruleRefExpr{
	pos: position{line: 272, col: 11, offset: 7567},
	name: "PathComponent",
},
},
},
},
},
{
	name: "Local",
	pos: position{line: 281, col: 1, offset: 7841},
	expr: &choiceExpr{
	pos: position{line: 281, col: 9, offset: 7851},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 281, col: 9, offset: 7851},
	name: "ParentPath",
},
&ruleRefExpr{
	pos: position{line: 281, col: 22, offset: 7864},
	name: "HerePath",
},
&ruleRefExpr{
	pos: position{line: 281, col: 33, offset: 7875},
	name: "AbsolutePath",
},
	},
},
},
{
	name: "ParentPath",
	pos: position{line: 283, col: 1, offset: 7889},
	expr: &actionExpr{
	pos: position{line: 283, col: 14, offset: 7904},
	run: (*parser).callonParentPath1,
	expr: &seqExpr{
	pos: position{line: 283, col: 14, offset: 7904},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 283, col: 14, offset: 7904},
	val: "..",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 283, col: 19, offset: 7909},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 283, col: 21, offset: 7911},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HerePath",
	pos: position{line: 284, col: 1, offset: 7967},
	expr: &actionExpr{
	pos: position{line: 284, col: 12, offset: 7980},
	run: (*parser).callonHerePath1,
	expr: &seqExpr{
	pos: position{line: 284, col: 12, offset: 7980},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 284, col: 12, offset: 7980},
	val: ".",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 284, col: 16, offset: 7984},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 284, col: 18, offset: 7986},
	name: "Path",
},
},
	},
},
},
},
{
	name: "AbsolutePath",
	pos: position{line: 285, col: 1, offset: 8025},
	expr: &actionExpr{
	pos: position{line: 285, col: 16, offset: 8042},
	run: (*parser).callonAbsolutePath1,
	expr: &labeledExpr{
	pos: position{line: 285, col: 16, offset: 8042},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 285, col: 18, offset: 8044},
	name: "Path",
},
},
},
},
{
	name: "Env",
	pos: position{line: 287, col: 1, offset: 8100},
	expr: &actionExpr{
	pos: position{line: 287, col: 7, offset: 8108},
	run: (*parser).callonEnv1,
	expr: &seqExpr{
	pos: position{line: 287, col: 7, offset: 8108},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 287, col: 7, offset: 8108},
	val: "env:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 287, col: 14, offset: 8115},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 287, col: 17, offset: 8118},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 287, col: 17, offset: 8118},
	name: "BashEnvironmentVariable",
},
&ruleRefExpr{
	pos: position{line: 287, col: 43, offset: 8144},
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
	pos: position{line: 289, col: 1, offset: 8189},
	expr: &actionExpr{
	pos: position{line: 289, col: 27, offset: 8217},
	run: (*parser).callonBashEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 289, col: 27, offset: 8217},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 289, col: 27, offset: 8217},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 289, col: 36, offset: 8226},
	expr: &charClassMatcher{
	pos: position{line: 289, col: 36, offset: 8226},
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
	pos: position{line: 293, col: 1, offset: 8282},
	expr: &actionExpr{
	pos: position{line: 293, col: 28, offset: 8311},
	run: (*parser).callonPosixEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 293, col: 28, offset: 8311},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 293, col: 28, offset: 8311},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 293, col: 32, offset: 8315},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 293, col: 34, offset: 8317},
	name: "PosixEnvironmentVariableContent",
},
},
&litMatcher{
	pos: position{line: 293, col: 66, offset: 8349},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PosixEnvironmentVariableContent",
	pos: position{line: 297, col: 1, offset: 8374},
	expr: &actionExpr{
	pos: position{line: 297, col: 35, offset: 8410},
	run: (*parser).callonPosixEnvironmentVariableContent1,
	expr: &labeledExpr{
	pos: position{line: 297, col: 35, offset: 8410},
	label: "v",
	expr: &oneOrMoreExpr{
	pos: position{line: 297, col: 37, offset: 8412},
	expr: &ruleRefExpr{
	pos: position{line: 297, col: 37, offset: 8412},
	name: "PosixEnvironmentVariableCharacter",
},
},
},
},
},
{
	name: "PosixEnvironmentVariableCharacter",
	pos: position{line: 306, col: 1, offset: 8625},
	expr: &choiceExpr{
	pos: position{line: 307, col: 7, offset: 8669},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 307, col: 7, offset: 8669},
	run: (*parser).callonPosixEnvironmentVariableCharacter2,
	expr: &litMatcher{
	pos: position{line: 307, col: 7, offset: 8669},
	val: "\\\"",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 308, col: 7, offset: 8709},
	run: (*parser).callonPosixEnvironmentVariableCharacter4,
	expr: &litMatcher{
	pos: position{line: 308, col: 7, offset: 8709},
	val: "\\\\",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 309, col: 7, offset: 8749},
	run: (*parser).callonPosixEnvironmentVariableCharacter6,
	expr: &litMatcher{
	pos: position{line: 309, col: 7, offset: 8749},
	val: "\\a",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 310, col: 7, offset: 8789},
	run: (*parser).callonPosixEnvironmentVariableCharacter8,
	expr: &litMatcher{
	pos: position{line: 310, col: 7, offset: 8789},
	val: "\\b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 311, col: 7, offset: 8829},
	run: (*parser).callonPosixEnvironmentVariableCharacter10,
	expr: &litMatcher{
	pos: position{line: 311, col: 7, offset: 8829},
	val: "\\f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 312, col: 7, offset: 8869},
	run: (*parser).callonPosixEnvironmentVariableCharacter12,
	expr: &litMatcher{
	pos: position{line: 312, col: 7, offset: 8869},
	val: "\\n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 313, col: 7, offset: 8909},
	run: (*parser).callonPosixEnvironmentVariableCharacter14,
	expr: &litMatcher{
	pos: position{line: 313, col: 7, offset: 8909},
	val: "\\r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 314, col: 7, offset: 8949},
	run: (*parser).callonPosixEnvironmentVariableCharacter16,
	expr: &litMatcher{
	pos: position{line: 314, col: 7, offset: 8949},
	val: "\\t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 315, col: 7, offset: 8989},
	run: (*parser).callonPosixEnvironmentVariableCharacter18,
	expr: &litMatcher{
	pos: position{line: 315, col: 7, offset: 8989},
	val: "\\v",
	ignoreCase: false,
},
},
&charClassMatcher{
	pos: position{line: 316, col: 7, offset: 9029},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 317, col: 7, offset: 9047},
	val: "[\\x23-\\x3c]",
	ranges: []rune{'#','<',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 318, col: 7, offset: 9065},
	val: "[\\x3e-\\x5b]",
	ranges: []rune{'>','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 319, col: 7, offset: 9083},
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
	pos: position{line: 321, col: 1, offset: 9096},
	expr: &choiceExpr{
	pos: position{line: 321, col: 14, offset: 9111},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 321, col: 14, offset: 9111},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 321, col: 24, offset: 9121},
	name: "Local",
},
&ruleRefExpr{
	pos: position{line: 321, col: 32, offset: 9129},
	name: "Env",
},
	},
},
},
{
	name: "ImportHashed",
	pos: position{line: 323, col: 1, offset: 9134},
	expr: &actionExpr{
	pos: position{line: 323, col: 16, offset: 9151},
	run: (*parser).callonImportHashed1,
	expr: &labeledExpr{
	pos: position{line: 323, col: 16, offset: 9151},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 323, col: 18, offset: 9153},
	name: "ImportType",
},
},
},
},
{
	name: "Import",
	pos: position{line: 325, col: 1, offset: 9222},
	expr: &choiceExpr{
	pos: position{line: 325, col: 10, offset: 9233},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 325, col: 10, offset: 9233},
	run: (*parser).callonImport2,
	expr: &seqExpr{
	pos: position{line: 325, col: 10, offset: 9233},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 325, col: 10, offset: 9233},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 325, col: 12, offset: 9235},
	name: "ImportHashed",
},
},
&ruleRefExpr{
	pos: position{line: 325, col: 25, offset: 9248},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 325, col: 27, offset: 9250},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 325, col: 30, offset: 9253},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 325, col: 33, offset: 9256},
	name: "Text",
},
	},
},
},
&actionExpr{
	pos: position{line: 326, col: 10, offset: 9353},
	run: (*parser).callonImport10,
	expr: &labeledExpr{
	pos: position{line: 326, col: 10, offset: 9353},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 326, col: 12, offset: 9355},
	name: "ImportHashed",
},
},
},
	},
},
},
{
	name: "LetBinding",
	pos: position{line: 329, col: 1, offset: 9450},
	expr: &actionExpr{
	pos: position{line: 329, col: 14, offset: 9465},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 329, col: 14, offset: 9465},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 329, col: 14, offset: 9465},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 329, col: 18, offset: 9469},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 329, col: 21, offset: 9472},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 329, col: 27, offset: 9478},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 329, col: 44, offset: 9495},
	name: "_",
},
&labeledExpr{
	pos: position{line: 329, col: 46, offset: 9497},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 329, col: 48, offset: 9499},
	expr: &seqExpr{
	pos: position{line: 329, col: 49, offset: 9500},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 329, col: 49, offset: 9500},
	name: "Annotation",
},
&ruleRefExpr{
	pos: position{line: 329, col: 60, offset: 9511},
	name: "_",
},
	},
},
},
},
&litMatcher{
	pos: position{line: 330, col: 13, offset: 9527},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 330, col: 17, offset: 9531},
	name: "_",
},
&labeledExpr{
	pos: position{line: 330, col: 19, offset: 9533},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 330, col: 21, offset: 9535},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 330, col: 32, offset: 9546},
	name: "_",
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 345, col: 1, offset: 9855},
	expr: &choiceExpr{
	pos: position{line: 346, col: 7, offset: 9876},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 346, col: 7, offset: 9876},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 346, col: 7, offset: 9876},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 346, col: 7, offset: 9876},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 346, col: 14, offset: 9883},
	name: "_",
},
&litMatcher{
	pos: position{line: 346, col: 16, offset: 9885},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 346, col: 20, offset: 9889},
	name: "_",
},
&labeledExpr{
	pos: position{line: 346, col: 22, offset: 9891},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 346, col: 28, offset: 9897},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 346, col: 45, offset: 9914},
	name: "_",
},
&litMatcher{
	pos: position{line: 346, col: 47, offset: 9916},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 346, col: 51, offset: 9920},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 346, col: 54, offset: 9923},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 346, col: 56, offset: 9925},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 346, col: 67, offset: 9936},
	name: "_",
},
&litMatcher{
	pos: position{line: 346, col: 69, offset: 9938},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 346, col: 73, offset: 9942},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 346, col: 75, offset: 9944},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 346, col: 81, offset: 9950},
	name: "_",
},
&labeledExpr{
	pos: position{line: 346, col: 83, offset: 9952},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 346, col: 88, offset: 9957},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 349, col: 7, offset: 10074},
	run: (*parser).callonExpression22,
	expr: &seqExpr{
	pos: position{line: 349, col: 7, offset: 10074},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 349, col: 7, offset: 10074},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 349, col: 10, offset: 10077},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 349, col: 13, offset: 10080},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 349, col: 18, offset: 10085},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 349, col: 29, offset: 10096},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 349, col: 31, offset: 10098},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 349, col: 36, offset: 10103},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 349, col: 39, offset: 10106},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 349, col: 41, offset: 10108},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 349, col: 52, offset: 10119},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 349, col: 54, offset: 10121},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 349, col: 59, offset: 10126},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 349, col: 62, offset: 10129},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 349, col: 64, offset: 10131},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 352, col: 7, offset: 10217},
	run: (*parser).callonExpression38,
	expr: &seqExpr{
	pos: position{line: 352, col: 7, offset: 10217},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 352, col: 7, offset: 10217},
	label: "bindings",
	expr: &oneOrMoreExpr{
	pos: position{line: 352, col: 16, offset: 10226},
	expr: &ruleRefExpr{
	pos: position{line: 352, col: 16, offset: 10226},
	name: "LetBinding",
},
},
},
&ruleRefExpr{
	pos: position{line: 352, col: 28, offset: 10238},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 352, col: 31, offset: 10241},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 352, col: 34, offset: 10244},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 352, col: 36, offset: 10246},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 359, col: 7, offset: 10486},
	run: (*parser).callonExpression47,
	expr: &seqExpr{
	pos: position{line: 359, col: 7, offset: 10486},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 359, col: 7, offset: 10486},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 359, col: 14, offset: 10493},
	name: "_",
},
&litMatcher{
	pos: position{line: 359, col: 16, offset: 10495},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 359, col: 20, offset: 10499},
	name: "_",
},
&labeledExpr{
	pos: position{line: 359, col: 22, offset: 10501},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 359, col: 28, offset: 10507},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 359, col: 45, offset: 10524},
	name: "_",
},
&litMatcher{
	pos: position{line: 359, col: 47, offset: 10526},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 359, col: 51, offset: 10530},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 359, col: 54, offset: 10533},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 359, col: 56, offset: 10535},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 359, col: 67, offset: 10546},
	name: "_",
},
&litMatcher{
	pos: position{line: 359, col: 69, offset: 10548},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 359, col: 73, offset: 10552},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 359, col: 75, offset: 10554},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 359, col: 81, offset: 10560},
	name: "_",
},
&labeledExpr{
	pos: position{line: 359, col: 83, offset: 10562},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 359, col: 88, offset: 10567},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 362, col: 7, offset: 10676},
	run: (*parser).callonExpression67,
	expr: &seqExpr{
	pos: position{line: 362, col: 7, offset: 10676},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 362, col: 7, offset: 10676},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 362, col: 9, offset: 10678},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 362, col: 28, offset: 10697},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 362, col: 30, offset: 10699},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 362, col: 36, offset: 10705},
	name: "_",
},
&labeledExpr{
	pos: position{line: 362, col: 38, offset: 10707},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 362, col: 40, offset: 10709},
	name: "Expression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 363, col: 7, offset: 10769},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 365, col: 1, offset: 10790},
	expr: &actionExpr{
	pos: position{line: 365, col: 14, offset: 10805},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 365, col: 14, offset: 10805},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 365, col: 14, offset: 10805},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 365, col: 18, offset: 10809},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 365, col: 21, offset: 10812},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 365, col: 23, offset: 10814},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 367, col: 1, offset: 10844},
	expr: &choiceExpr{
	pos: position{line: 368, col: 5, offset: 10872},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 368, col: 5, offset: 10872},
	name: "EmptyList",
},
&actionExpr{
	pos: position{line: 369, col: 5, offset: 10886},
	run: (*parser).callonAnnotatedExpression3,
	expr: &seqExpr{
	pos: position{line: 369, col: 5, offset: 10886},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 369, col: 5, offset: 10886},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 369, col: 7, offset: 10888},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 369, col: 26, offset: 10907},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 369, col: 28, offset: 10909},
	expr: &seqExpr{
	pos: position{line: 369, col: 29, offset: 10910},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 369, col: 29, offset: 10910},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 369, col: 31, offset: 10912},
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
	pos: position{line: 374, col: 1, offset: 11037},
	expr: &actionExpr{
	pos: position{line: 374, col: 13, offset: 11051},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 374, col: 13, offset: 11051},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 374, col: 13, offset: 11051},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 374, col: 17, offset: 11055},
	name: "_",
},
&litMatcher{
	pos: position{line: 374, col: 19, offset: 11057},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 374, col: 23, offset: 11061},
	name: "_",
},
&litMatcher{
	pos: position{line: 374, col: 25, offset: 11063},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 374, col: 29, offset: 11067},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 374, col: 32, offset: 11070},
	name: "List",
},
&ruleRefExpr{
	pos: position{line: 374, col: 37, offset: 11075},
	name: "_",
},
&labeledExpr{
	pos: position{line: 374, col: 39, offset: 11077},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 374, col: 41, offset: 11079},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 378, col: 1, offset: 11142},
	expr: &ruleRefExpr{
	pos: position{line: 378, col: 22, offset: 11165},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 380, col: 1, offset: 11186},
	expr: &ruleRefExpr{
	pos: position{line: 380, col: 23, offset: 11210},
	name: "PlusExpression",
},
},
{
	name: "MorePlus",
	pos: position{line: 382, col: 1, offset: 11226},
	expr: &actionExpr{
	pos: position{line: 382, col: 12, offset: 11239},
	run: (*parser).callonMorePlus1,
	expr: &seqExpr{
	pos: position{line: 382, col: 12, offset: 11239},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 382, col: 12, offset: 11239},
	name: "_",
},
&litMatcher{
	pos: position{line: 382, col: 14, offset: 11241},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 382, col: 18, offset: 11245},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 382, col: 21, offset: 11248},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 382, col: 23, offset: 11250},
	name: "TimesExpression",
},
},
	},
},
},
},
{
	name: "PlusExpression",
	pos: position{line: 383, col: 1, offset: 11284},
	expr: &actionExpr{
	pos: position{line: 384, col: 7, offset: 11309},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 384, col: 7, offset: 11309},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 384, col: 7, offset: 11309},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 384, col: 13, offset: 11315},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 384, col: 29, offset: 11331},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 384, col: 34, offset: 11336},
	expr: &ruleRefExpr{
	pos: position{line: 384, col: 34, offset: 11336},
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
	pos: position{line: 393, col: 1, offset: 11564},
	expr: &actionExpr{
	pos: position{line: 393, col: 13, offset: 11578},
	run: (*parser).callonMoreTimes1,
	expr: &seqExpr{
	pos: position{line: 393, col: 13, offset: 11578},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 393, col: 13, offset: 11578},
	name: "_",
},
&litMatcher{
	pos: position{line: 393, col: 15, offset: 11580},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 393, col: 19, offset: 11584},
	name: "_",
},
&labeledExpr{
	pos: position{line: 393, col: 21, offset: 11586},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 393, col: 23, offset: 11588},
	name: "ApplicationExpression",
},
},
	},
},
},
},
{
	name: "TimesExpression",
	pos: position{line: 394, col: 1, offset: 11628},
	expr: &actionExpr{
	pos: position{line: 395, col: 7, offset: 11654},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 395, col: 7, offset: 11654},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 395, col: 7, offset: 11654},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 395, col: 13, offset: 11660},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 395, col: 35, offset: 11682},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 395, col: 40, offset: 11687},
	expr: &ruleRefExpr{
	pos: position{line: 395, col: 40, offset: 11687},
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
	pos: position{line: 404, col: 1, offset: 11917},
	expr: &actionExpr{
	pos: position{line: 404, col: 25, offset: 11943},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 404, col: 25, offset: 11943},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 404, col: 25, offset: 11943},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 404, col: 27, offset: 11945},
	name: "FirstApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 404, col: 54, offset: 11972},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 404, col: 59, offset: 11977},
	expr: &seqExpr{
	pos: position{line: 404, col: 60, offset: 11978},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 404, col: 60, offset: 11978},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 404, col: 63, offset: 11981},
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
	pos: position{line: 413, col: 1, offset: 12231},
	expr: &choiceExpr{
	pos: position{line: 414, col: 8, offset: 12269},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 414, col: 8, offset: 12269},
	run: (*parser).callonFirstApplicationExpression2,
	expr: &seqExpr{
	pos: position{line: 414, col: 8, offset: 12269},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 414, col: 8, offset: 12269},
	name: "Some",
},
&ruleRefExpr{
	pos: position{line: 414, col: 13, offset: 12274},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 414, col: 16, offset: 12277},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 414, col: 18, offset: 12279},
	name: "ImportExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 415, col: 8, offset: 12334},
	name: "ImportExpression",
},
	},
},
},
{
	name: "ImportExpression",
	pos: position{line: 417, col: 1, offset: 12352},
	expr: &choiceExpr{
	pos: position{line: 417, col: 20, offset: 12373},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 417, col: 20, offset: 12373},
	name: "Import",
},
&ruleRefExpr{
	pos: position{line: 417, col: 29, offset: 12382},
	name: "SelectorExpression",
},
	},
},
},
{
	name: "SelectorExpression",
	pos: position{line: 419, col: 1, offset: 12402},
	expr: &actionExpr{
	pos: position{line: 419, col: 22, offset: 12425},
	run: (*parser).callonSelectorExpression1,
	expr: &seqExpr{
	pos: position{line: 419, col: 22, offset: 12425},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 419, col: 22, offset: 12425},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 419, col: 24, offset: 12427},
	name: "PrimitiveExpression",
},
},
&labeledExpr{
	pos: position{line: 419, col: 44, offset: 12447},
	label: "ls",
	expr: &zeroOrMoreExpr{
	pos: position{line: 419, col: 47, offset: 12450},
	expr: &seqExpr{
	pos: position{line: 419, col: 48, offset: 12451},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 419, col: 48, offset: 12451},
	name: "_",
},
&litMatcher{
	pos: position{line: 419, col: 50, offset: 12453},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 419, col: 54, offset: 12457},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 419, col: 56, offset: 12459},
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
	pos: position{line: 429, col: 1, offset: 12692},
	expr: &choiceExpr{
	pos: position{line: 430, col: 7, offset: 12722},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 430, col: 7, offset: 12722},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 431, col: 7, offset: 12742},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 432, col: 7, offset: 12763},
	name: "IntegerLiteral",
},
&ruleRefExpr{
	pos: position{line: 433, col: 7, offset: 12784},
	name: "TextLiteral",
},
&actionExpr{
	pos: position{line: 434, col: 7, offset: 12802},
	run: (*parser).callonPrimitiveExpression6,
	expr: &seqExpr{
	pos: position{line: 434, col: 7, offset: 12802},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 434, col: 7, offset: 12802},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 434, col: 11, offset: 12806},
	name: "_",
},
&labeledExpr{
	pos: position{line: 434, col: 13, offset: 12808},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 434, col: 15, offset: 12810},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 434, col: 35, offset: 12830},
	name: "_",
},
&litMatcher{
	pos: position{line: 434, col: 37, offset: 12832},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 435, col: 7, offset: 12860},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 436, col: 7, offset: 12886},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 437, col: 7, offset: 12903},
	run: (*parser).callonPrimitiveExpression16,
	expr: &seqExpr{
	pos: position{line: 437, col: 7, offset: 12903},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 437, col: 7, offset: 12903},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 437, col: 11, offset: 12907},
	name: "_",
},
&labeledExpr{
	pos: position{line: 437, col: 14, offset: 12910},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 437, col: 16, offset: 12912},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 437, col: 27, offset: 12923},
	name: "_",
},
&litMatcher{
	pos: position{line: 437, col: 29, offset: 12925},
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
	pos: position{line: 439, col: 1, offset: 12948},
	expr: &choiceExpr{
	pos: position{line: 440, col: 7, offset: 12978},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 440, col: 7, offset: 12978},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &litMatcher{
	pos: position{line: 440, col: 7, offset: 12978},
	val: "=",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 441, col: 7, offset: 13033},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 442, col: 7, offset: 13058},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 443, col: 7, offset: 13086},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 443, col: 7, offset: 13086},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 445, col: 1, offset: 13132},
	expr: &actionExpr{
	pos: position{line: 445, col: 19, offset: 13152},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 445, col: 19, offset: 13152},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 445, col: 19, offset: 13152},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 445, col: 24, offset: 13157},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 445, col: 33, offset: 13166},
	name: "_",
},
&litMatcher{
	pos: position{line: 445, col: 35, offset: 13168},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 445, col: 39, offset: 13172},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 445, col: 42, offset: 13175},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 445, col: 47, offset: 13180},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 448, col: 1, offset: 13237},
	expr: &actionExpr{
	pos: position{line: 448, col: 18, offset: 13256},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 448, col: 18, offset: 13256},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 448, col: 18, offset: 13256},
	name: "_",
},
&litMatcher{
	pos: position{line: 448, col: 20, offset: 13258},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 448, col: 24, offset: 13262},
	name: "_",
},
&labeledExpr{
	pos: position{line: 448, col: 26, offset: 13264},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 448, col: 28, offset: 13266},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 449, col: 1, offset: 13298},
	expr: &actionExpr{
	pos: position{line: 450, col: 7, offset: 13327},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 450, col: 7, offset: 13327},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 450, col: 7, offset: 13327},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 450, col: 13, offset: 13333},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 450, col: 29, offset: 13349},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 450, col: 34, offset: 13354},
	expr: &ruleRefExpr{
	pos: position{line: 450, col: 34, offset: 13354},
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
	pos: position{line: 460, col: 1, offset: 13750},
	expr: &actionExpr{
	pos: position{line: 460, col: 22, offset: 13773},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 460, col: 22, offset: 13773},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 460, col: 22, offset: 13773},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 460, col: 27, offset: 13778},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 460, col: 36, offset: 13787},
	name: "_",
},
&litMatcher{
	pos: position{line: 460, col: 38, offset: 13789},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 460, col: 42, offset: 13793},
	name: "_",
},
&labeledExpr{
	pos: position{line: 460, col: 44, offset: 13795},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 460, col: 49, offset: 13800},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 463, col: 1, offset: 13857},
	expr: &actionExpr{
	pos: position{line: 463, col: 21, offset: 13879},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 463, col: 21, offset: 13879},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 463, col: 21, offset: 13879},
	name: "_",
},
&litMatcher{
	pos: position{line: 463, col: 23, offset: 13881},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 463, col: 27, offset: 13885},
	name: "_",
},
&labeledExpr{
	pos: position{line: 463, col: 29, offset: 13887},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 463, col: 31, offset: 13889},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 464, col: 1, offset: 13924},
	expr: &actionExpr{
	pos: position{line: 465, col: 7, offset: 13956},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 465, col: 7, offset: 13956},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 465, col: 7, offset: 13956},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 465, col: 13, offset: 13962},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 465, col: 32, offset: 13981},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 465, col: 37, offset: 13986},
	expr: &ruleRefExpr{
	pos: position{line: 465, col: 37, offset: 13986},
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
	pos: position{line: 475, col: 1, offset: 14388},
	expr: &actionExpr{
	pos: position{line: 475, col: 12, offset: 14401},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 475, col: 12, offset: 14401},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 475, col: 12, offset: 14401},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 475, col: 16, offset: 14405},
	name: "_",
},
&labeledExpr{
	pos: position{line: 475, col: 18, offset: 14407},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 475, col: 20, offset: 14409},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 475, col: 31, offset: 14420},
	name: "_",
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 477, col: 1, offset: 14439},
	expr: &actionExpr{
	pos: position{line: 478, col: 7, offset: 14469},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 478, col: 7, offset: 14469},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 478, col: 7, offset: 14469},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 478, col: 11, offset: 14473},
	name: "_",
},
&labeledExpr{
	pos: position{line: 478, col: 13, offset: 14475},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 478, col: 19, offset: 14481},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 478, col: 30, offset: 14492},
	name: "_",
},
&labeledExpr{
	pos: position{line: 478, col: 32, offset: 14494},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 478, col: 37, offset: 14499},
	expr: &ruleRefExpr{
	pos: position{line: 478, col: 37, offset: 14499},
	name: "MoreList",
},
},
},
&litMatcher{
	pos: position{line: 478, col: 47, offset: 14509},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 488, col: 1, offset: 14785},
	expr: &notExpr{
	pos: position{line: 488, col: 7, offset: 14793},
	expr: &anyMatcher{
	line: 488, col: 8, offset: 14794,
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
    return TextLit{Chunks: outChunks, Suffix: str.String()}, nil
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

func (c *current) onAbsolutePath1(p interface{}) (interface{}, error) {
 return Local(path.Join("/", p.(string))), nil 
}

func (p *parser) callonAbsolutePath1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAbsolutePath1(stack["p"])
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

