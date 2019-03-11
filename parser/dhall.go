
package parser

import (
"bytes"
"errors"
"fmt"
"io"
"io/ioutil"
"math"
"os"
"strconv"
"strings"
"unicode"
"unicode/utf8"
)
import "github.com/philandstuff/dhall-golang/ast"


var g = &grammar {
	rules: []*rule{
{
	name: "DhallFile",
	pos: position{line: 21, col: 1, offset: 180},
	expr: &actionExpr{
	pos: position{line: 21, col: 13, offset: 194},
	run: (*parser).callonDhallFile1,
	expr: &seqExpr{
	pos: position{line: 21, col: 13, offset: 194},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 21, col: 13, offset: 194},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 21, col: 15, offset: 196},
	name: "CompleteExpression",
},
},
&ruleRefExpr{
	pos: position{line: 21, col: 34, offset: 215},
	name: "EOF",
},
	},
},
},
},
{
	name: "CompleteExpression",
	pos: position{line: 23, col: 1, offset: 238},
	expr: &actionExpr{
	pos: position{line: 23, col: 22, offset: 261},
	run: (*parser).callonCompleteExpression1,
	expr: &seqExpr{
	pos: position{line: 23, col: 22, offset: 261},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 23, col: 22, offset: 261},
	name: "_",
},
&labeledExpr{
	pos: position{line: 23, col: 24, offset: 263},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 23, col: 26, offset: 265},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "EOL",
	pos: position{line: 25, col: 1, offset: 295},
	expr: &choiceExpr{
	pos: position{line: 25, col: 7, offset: 303},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 25, col: 7, offset: 303},
	val: "\n",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 25, col: 14, offset: 310},
	val: "\r\n",
	ignoreCase: false,
},
	},
},
},
{
	name: "BlockComment",
	pos: position{line: 27, col: 1, offset: 318},
	expr: &seqExpr{
	pos: position{line: 27, col: 16, offset: 335},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 27, col: 16, offset: 335},
	val: "{-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 27, col: 21, offset: 340},
	name: "BlockCommentContinue",
},
	},
},
},
{
	name: "BlockCommentChunk",
	pos: position{line: 29, col: 1, offset: 362},
	expr: &choiceExpr{
	pos: position{line: 30, col: 5, offset: 388},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 30, col: 5, offset: 388},
	name: "BlockComment",
},
&charClassMatcher{
	pos: position{line: 31, col: 5, offset: 405},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 32, col: 5, offset: 431},
	name: "EOL",
},
	},
},
},
{
	name: "BlockCommentContinue",
	pos: position{line: 34, col: 1, offset: 436},
	expr: &choiceExpr{
	pos: position{line: 34, col: 24, offset: 461},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 34, col: 24, offset: 461},
	val: "-}",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 34, col: 31, offset: 468},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 34, col: 31, offset: 468},
	name: "BlockCommentChunk",
},
&ruleRefExpr{
	pos: position{line: 34, col: 49, offset: 486},
	name: "BlockCommentContinue",
},
	},
},
	},
},
},
{
	name: "NotEOL",
	pos: position{line: 36, col: 1, offset: 508},
	expr: &charClassMatcher{
	pos: position{line: 36, col: 10, offset: 519},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "LineComment",
	pos: position{line: 38, col: 1, offset: 542},
	expr: &actionExpr{
	pos: position{line: 38, col: 15, offset: 558},
	run: (*parser).callonLineComment1,
	expr: &seqExpr{
	pos: position{line: 38, col: 15, offset: 558},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 38, col: 15, offset: 558},
	val: "--",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 38, col: 20, offset: 563},
	label: "content",
	expr: &actionExpr{
	pos: position{line: 38, col: 29, offset: 572},
	run: (*parser).callonLineComment5,
	expr: &zeroOrMoreExpr{
	pos: position{line: 38, col: 29, offset: 572},
	expr: &ruleRefExpr{
	pos: position{line: 38, col: 29, offset: 572},
	name: "NotEOL",
},
},
},
},
&ruleRefExpr{
	pos: position{line: 38, col: 68, offset: 611},
	name: "EOL",
},
	},
},
},
},
{
	name: "WhitespaceChunk",
	pos: position{line: 40, col: 1, offset: 640},
	expr: &choiceExpr{
	pos: position{line: 40, col: 19, offset: 660},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 40, col: 19, offset: 660},
	val: " ",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 40, col: 25, offset: 666},
	val: "\t",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 40, col: 32, offset: 673},
	name: "EOL",
},
&ruleRefExpr{
	pos: position{line: 40, col: 38, offset: 679},
	name: "LineComment",
},
&ruleRefExpr{
	pos: position{line: 40, col: 52, offset: 693},
	name: "BlockComment",
},
	},
},
},
{
	name: "_",
	pos: position{line: 42, col: 1, offset: 707},
	expr: &zeroOrMoreExpr{
	pos: position{line: 42, col: 5, offset: 713},
	expr: &ruleRefExpr{
	pos: position{line: 42, col: 5, offset: 713},
	name: "WhitespaceChunk",
},
},
},
{
	name: "NonemptyWhitespace",
	pos: position{line: 44, col: 1, offset: 731},
	expr: &oneOrMoreExpr{
	pos: position{line: 44, col: 22, offset: 754},
	expr: &ruleRefExpr{
	pos: position{line: 44, col: 22, offset: 754},
	name: "WhitespaceChunk",
},
},
},
{
	name: "HexDig",
	pos: position{line: 46, col: 1, offset: 772},
	expr: &charClassMatcher{
	pos: position{line: 46, col: 10, offset: 783},
	val: "[0-9a-f]i",
	ranges: []rune{'0','9','a','f',},
	ignoreCase: true,
	inverted: false,
},
},
{
	name: "SimpleLabel",
	pos: position{line: 48, col: 1, offset: 794},
	expr: &actionExpr{
	pos: position{line: 48, col: 15, offset: 810},
	run: (*parser).callonSimpleLabel1,
	expr: &seqExpr{
	pos: position{line: 48, col: 15, offset: 810},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 48, col: 15, offset: 810},
	expr: &ruleRefExpr{
	pos: position{line: 48, col: 16, offset: 811},
	name: "KeywordRaw",
},
},
&charClassMatcher{
	pos: position{line: 49, col: 13, offset: 834},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 49, col: 23, offset: 844},
	expr: &charClassMatcher{
	pos: position{line: 49, col: 23, offset: 844},
	val: "[A-Za-z0-9_/-]",
	chars: []rune{'_','/','-',},
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
	name: "Label",
	pos: position{line: 53, col: 1, offset: 908},
	expr: &actionExpr{
	pos: position{line: 53, col: 9, offset: 918},
	run: (*parser).callonLabel1,
	expr: &seqExpr{
	pos: position{line: 53, col: 9, offset: 918},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 53, col: 9, offset: 918},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 53, col: 15, offset: 924},
	name: "SimpleLabel",
},
},
&ruleRefExpr{
	pos: position{line: 53, col: 27, offset: 936},
	name: "_",
},
	},
},
},
},
{
	name: "EscapedChar",
	pos: position{line: 57, col: 1, offset: 994},
	expr: &actionExpr{
	pos: position{line: 58, col: 3, offset: 1012},
	run: (*parser).callonEscapedChar1,
	expr: &seqExpr{
	pos: position{line: 58, col: 3, offset: 1012},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 58, col: 3, offset: 1012},
	val: "\\",
	ignoreCase: false,
},
&choiceExpr{
	pos: position{line: 59, col: 5, offset: 1021},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 59, col: 5, offset: 1021},
	val: "\"",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 60, col: 10, offset: 1034},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 61, col: 10, offset: 1047},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 62, col: 10, offset: 1061},
	val: "/",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 63, col: 10, offset: 1074},
	val: "b",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 64, col: 10, offset: 1087},
	val: "f",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 65, col: 10, offset: 1100},
	val: "n",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 66, col: 10, offset: 1113},
	val: "r",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 67, col: 10, offset: 1126},
	val: "t",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 68, col: 10, offset: 1139},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 68, col: 10, offset: 1139},
	val: "u",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 68, col: 14, offset: 1143},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 68, col: 21, offset: 1150},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 68, col: 28, offset: 1157},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 68, col: 35, offset: 1164},
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
	pos: position{line: 89, col: 1, offset: 1607},
	expr: &choiceExpr{
	pos: position{line: 90, col: 6, offset: 1633},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 90, col: 6, offset: 1633},
	run: (*parser).callonDoubleQuoteChunk2,
	expr: &seqExpr{
	pos: position{line: 90, col: 6, offset: 1633},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 90, col: 6, offset: 1633},
	val: "${",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 90, col: 11, offset: 1638},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 90, col: 13, offset: 1640},
	name: "CompleteExpression",
},
},
&litMatcher{
	pos: position{line: 90, col: 32, offset: 1659},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 91, col: 6, offset: 1686},
	name: "EscapedChar",
},
&charClassMatcher{
	pos: position{line: 92, col: 6, offset: 1703},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 93, col: 6, offset: 1720},
	val: "[\\x23-\\x5b]",
	ranges: []rune{'#','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 94, col: 6, offset: 1737},
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
	pos: position{line: 96, col: 1, offset: 1756},
	expr: &actionExpr{
	pos: position{line: 96, col: 22, offset: 1779},
	run: (*parser).callonDoubleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 96, col: 22, offset: 1779},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 96, col: 22, offset: 1779},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 96, col: 26, offset: 1783},
	label: "chunks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 96, col: 33, offset: 1790},
	expr: &ruleRefExpr{
	pos: position{line: 96, col: 33, offset: 1790},
	name: "DoubleQuoteChunk",
},
},
},
&litMatcher{
	pos: position{line: 96, col: 51, offset: 1808},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TextLiteral",
	pos: position{line: 113, col: 1, offset: 2292},
	expr: &actionExpr{
	pos: position{line: 113, col: 15, offset: 2308},
	run: (*parser).callonTextLiteral1,
	expr: &seqExpr{
	pos: position{line: 113, col: 15, offset: 2308},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 113, col: 15, offset: 2308},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 113, col: 17, offset: 2310},
	name: "DoubleQuoteLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 113, col: 36, offset: 2329},
	name: "_",
},
	},
},
},
},
{
	name: "ReservedRaw",
	pos: position{line: 115, col: 1, offset: 2350},
	expr: &choiceExpr{
	pos: position{line: 115, col: 15, offset: 2366},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 115, col: 15, offset: 2366},
	run: (*parser).callonReservedRaw2,
	expr: &litMatcher{
	pos: position{line: 115, col: 15, offset: 2366},
	val: "Bool",
	ignoreCase: false,
},
},
&litMatcher{
	pos: position{line: 116, col: 5, offset: 2402},
	val: "Optional",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 117, col: 5, offset: 2417},
	val: "None",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 118, col: 5, offset: 2428},
	run: (*parser).callonReservedRaw6,
	expr: &litMatcher{
	pos: position{line: 118, col: 5, offset: 2428},
	val: "Natural",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 119, col: 5, offset: 2470},
	run: (*parser).callonReservedRaw8,
	expr: &litMatcher{
	pos: position{line: 119, col: 5, offset: 2470},
	val: "Integer",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 120, col: 5, offset: 2512},
	run: (*parser).callonReservedRaw10,
	expr: &litMatcher{
	pos: position{line: 120, col: 5, offset: 2512},
	val: "Double",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 121, col: 5, offset: 2552},
	run: (*parser).callonReservedRaw12,
	expr: &litMatcher{
	pos: position{line: 121, col: 5, offset: 2552},
	val: "Text",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 122, col: 5, offset: 2588},
	run: (*parser).callonReservedRaw14,
	expr: &litMatcher{
	pos: position{line: 122, col: 5, offset: 2588},
	val: "List",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 123, col: 5, offset: 2624},
	run: (*parser).callonReservedRaw16,
	expr: &litMatcher{
	pos: position{line: 123, col: 5, offset: 2624},
	val: "True",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 124, col: 5, offset: 2660},
	run: (*parser).callonReservedRaw18,
	expr: &litMatcher{
	pos: position{line: 124, col: 5, offset: 2660},
	val: "False",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 125, col: 5, offset: 2698},
	run: (*parser).callonReservedRaw20,
	expr: &litMatcher{
	pos: position{line: 125, col: 5, offset: 2698},
	val: "NaN",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 126, col: 5, offset: 2750},
	run: (*parser).callonReservedRaw22,
	expr: &litMatcher{
	pos: position{line: 126, col: 5, offset: 2750},
	val: "Infinity",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 127, col: 5, offset: 2808},
	run: (*parser).callonReservedRaw24,
	expr: &litMatcher{
	pos: position{line: 127, col: 5, offset: 2808},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 128, col: 5, offset: 2844},
	run: (*parser).callonReservedRaw26,
	expr: &litMatcher{
	pos: position{line: 128, col: 5, offset: 2844},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 129, col: 5, offset: 2880},
	run: (*parser).callonReservedRaw28,
	expr: &litMatcher{
	pos: position{line: 129, col: 5, offset: 2880},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "Reserved",
	pos: position{line: 131, col: 1, offset: 2913},
	expr: &actionExpr{
	pos: position{line: 131, col: 12, offset: 2926},
	run: (*parser).callonReserved1,
	expr: &seqExpr{
	pos: position{line: 131, col: 12, offset: 2926},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 131, col: 12, offset: 2926},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 131, col: 14, offset: 2928},
	name: "ReservedRaw",
},
},
&ruleRefExpr{
	pos: position{line: 131, col: 26, offset: 2940},
	name: "_",
},
	},
},
},
},
{
	name: "KeywordRaw",
	pos: position{line: 133, col: 1, offset: 2961},
	expr: &choiceExpr{
	pos: position{line: 133, col: 14, offset: 2976},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 133, col: 14, offset: 2976},
	val: "if",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 134, col: 5, offset: 2985},
	val: "then",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 135, col: 5, offset: 2996},
	val: "else",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 136, col: 5, offset: 3007},
	val: "let",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 137, col: 5, offset: 3017},
	val: "in",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 138, col: 5, offset: 3026},
	val: "as",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 139, col: 5, offset: 3035},
	val: "using",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 140, col: 5, offset: 3047},
	val: "merge",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 141, col: 5, offset: 3059},
	val: "constructors",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 142, col: 5, offset: 3078},
	val: "Some",
	ignoreCase: false,
},
	},
},
},
{
	name: "If",
	pos: position{line: 144, col: 1, offset: 3086},
	expr: &seqExpr{
	pos: position{line: 144, col: 6, offset: 3093},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 144, col: 6, offset: 3093},
	val: "if",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 144, col: 11, offset: 3098},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Then",
	pos: position{line: 145, col: 1, offset: 3117},
	expr: &seqExpr{
	pos: position{line: 145, col: 8, offset: 3126},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 145, col: 8, offset: 3126},
	val: "then",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 145, col: 15, offset: 3133},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Else",
	pos: position{line: 146, col: 1, offset: 3152},
	expr: &seqExpr{
	pos: position{line: 146, col: 8, offset: 3161},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 146, col: 8, offset: 3161},
	val: "else",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 146, col: 15, offset: 3168},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Let",
	pos: position{line: 147, col: 1, offset: 3187},
	expr: &seqExpr{
	pos: position{line: 147, col: 7, offset: 3195},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 147, col: 7, offset: 3195},
	val: "let",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 147, col: 13, offset: 3201},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "In",
	pos: position{line: 148, col: 1, offset: 3220},
	expr: &seqExpr{
	pos: position{line: 148, col: 6, offset: 3227},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 148, col: 6, offset: 3227},
	val: "in",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 148, col: 11, offset: 3232},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "As",
	pos: position{line: 149, col: 1, offset: 3251},
	expr: &seqExpr{
	pos: position{line: 149, col: 6, offset: 3258},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 149, col: 6, offset: 3258},
	val: "as",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 149, col: 11, offset: 3263},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Using",
	pos: position{line: 150, col: 1, offset: 3282},
	expr: &seqExpr{
	pos: position{line: 150, col: 9, offset: 3292},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 150, col: 9, offset: 3292},
	val: "using",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 150, col: 17, offset: 3300},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Merge",
	pos: position{line: 151, col: 1, offset: 3319},
	expr: &seqExpr{
	pos: position{line: 151, col: 9, offset: 3329},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 151, col: 9, offset: 3329},
	val: "merge",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 151, col: 17, offset: 3337},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Some",
	pos: position{line: 152, col: 1, offset: 3356},
	expr: &seqExpr{
	pos: position{line: 152, col: 8, offset: 3365},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 152, col: 8, offset: 3365},
	val: "Some",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 152, col: 15, offset: 3372},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Optional",
	pos: position{line: 153, col: 1, offset: 3391},
	expr: &seqExpr{
	pos: position{line: 153, col: 12, offset: 3404},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 153, col: 12, offset: 3404},
	val: "Optional",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 153, col: 23, offset: 3415},
	name: "_",
},
	},
},
},
{
	name: "Text",
	pos: position{line: 154, col: 1, offset: 3417},
	expr: &seqExpr{
	pos: position{line: 154, col: 8, offset: 3426},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 154, col: 8, offset: 3426},
	val: "Text",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 154, col: 15, offset: 3433},
	name: "_",
},
	},
},
},
{
	name: "List",
	pos: position{line: 155, col: 1, offset: 3435},
	expr: &seqExpr{
	pos: position{line: 155, col: 8, offset: 3444},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 155, col: 8, offset: 3444},
	val: "List",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 155, col: 15, offset: 3451},
	name: "_",
},
	},
},
},
{
	name: "Equal",
	pos: position{line: 157, col: 1, offset: 3454},
	expr: &seqExpr{
	pos: position{line: 157, col: 9, offset: 3464},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 157, col: 9, offset: 3464},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 157, col: 13, offset: 3468},
	name: "_",
},
	},
},
},
{
	name: "Plus",
	pos: position{line: 158, col: 1, offset: 3470},
	expr: &seqExpr{
	pos: position{line: 158, col: 8, offset: 3479},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 158, col: 8, offset: 3479},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 158, col: 12, offset: 3483},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Times",
	pos: position{line: 159, col: 1, offset: 3502},
	expr: &seqExpr{
	pos: position{line: 159, col: 9, offset: 3512},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 159, col: 9, offset: 3512},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 159, col: 13, offset: 3516},
	name: "_",
},
	},
},
},
{
	name: "Dot",
	pos: position{line: 160, col: 1, offset: 3518},
	expr: &seqExpr{
	pos: position{line: 160, col: 7, offset: 3526},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 160, col: 7, offset: 3526},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 160, col: 11, offset: 3530},
	name: "_",
},
	},
},
},
{
	name: "OpenBrace",
	pos: position{line: 161, col: 1, offset: 3532},
	expr: &seqExpr{
	pos: position{line: 161, col: 13, offset: 3546},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 161, col: 13, offset: 3546},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 161, col: 17, offset: 3550},
	name: "_",
},
	},
},
},
{
	name: "CloseBrace",
	pos: position{line: 162, col: 1, offset: 3552},
	expr: &seqExpr{
	pos: position{line: 162, col: 14, offset: 3567},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 162, col: 14, offset: 3567},
	val: "}",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 162, col: 18, offset: 3571},
	name: "_",
},
	},
},
},
{
	name: "OpenBracket",
	pos: position{line: 163, col: 1, offset: 3573},
	expr: &seqExpr{
	pos: position{line: 163, col: 15, offset: 3589},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 163, col: 15, offset: 3589},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 163, col: 19, offset: 3593},
	name: "_",
},
	},
},
},
{
	name: "CloseBracket",
	pos: position{line: 164, col: 1, offset: 3595},
	expr: &seqExpr{
	pos: position{line: 164, col: 16, offset: 3612},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 164, col: 16, offset: 3612},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 164, col: 20, offset: 3616},
	name: "_",
},
	},
},
},
{
	name: "Comma",
	pos: position{line: 165, col: 1, offset: 3618},
	expr: &seqExpr{
	pos: position{line: 165, col: 9, offset: 3628},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 165, col: 9, offset: 3628},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 165, col: 13, offset: 3632},
	name: "_",
},
	},
},
},
{
	name: "OpenParens",
	pos: position{line: 166, col: 1, offset: 3634},
	expr: &seqExpr{
	pos: position{line: 166, col: 14, offset: 3649},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 166, col: 14, offset: 3649},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 166, col: 18, offset: 3653},
	name: "_",
},
	},
},
},
{
	name: "CloseParens",
	pos: position{line: 167, col: 1, offset: 3655},
	expr: &seqExpr{
	pos: position{line: 167, col: 15, offset: 3671},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 167, col: 15, offset: 3671},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 167, col: 19, offset: 3675},
	name: "_",
},
	},
},
},
{
	name: "At",
	pos: position{line: 168, col: 1, offset: 3677},
	expr: &seqExpr{
	pos: position{line: 168, col: 6, offset: 3684},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 168, col: 6, offset: 3684},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 168, col: 10, offset: 3688},
	name: "_",
},
	},
},
},
{
	name: "Colon",
	pos: position{line: 169, col: 1, offset: 3690},
	expr: &seqExpr{
	pos: position{line: 169, col: 9, offset: 3700},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 169, col: 9, offset: 3700},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 169, col: 13, offset: 3704},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Lambda",
	pos: position{line: 171, col: 1, offset: 3724},
	expr: &seqExpr{
	pos: position{line: 171, col: 10, offset: 3735},
	exprs: []interface{}{
&choiceExpr{
	pos: position{line: 171, col: 11, offset: 3736},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 171, col: 11, offset: 3736},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 171, col: 18, offset: 3743},
	val: "λ",
	ignoreCase: false,
},
	},
},
&ruleRefExpr{
	pos: position{line: 171, col: 23, offset: 3749},
	name: "_",
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 172, col: 1, offset: 3751},
	expr: &seqExpr{
	pos: position{line: 172, col: 10, offset: 3762},
	exprs: []interface{}{
&choiceExpr{
	pos: position{line: 172, col: 11, offset: 3763},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 172, col: 11, offset: 3763},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 172, col: 22, offset: 3774},
	val: "∀",
	ignoreCase: false,
},
	},
},
&ruleRefExpr{
	pos: position{line: 172, col: 27, offset: 3781},
	name: "_",
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 173, col: 1, offset: 3783},
	expr: &seqExpr{
	pos: position{line: 173, col: 9, offset: 3793},
	exprs: []interface{}{
&choiceExpr{
	pos: position{line: 173, col: 10, offset: 3794},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 173, col: 10, offset: 3794},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 173, col: 17, offset: 3801},
	val: "→",
	ignoreCase: false,
},
	},
},
&ruleRefExpr{
	pos: position{line: 173, col: 22, offset: 3808},
	name: "_",
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 175, col: 1, offset: 3811},
	expr: &seqExpr{
	pos: position{line: 175, col: 12, offset: 3824},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 175, col: 12, offset: 3824},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 175, col: 17, offset: 3829},
	expr: &charClassMatcher{
	pos: position{line: 175, col: 17, offset: 3829},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 175, col: 23, offset: 3835},
	expr: &charClassMatcher{
	pos: position{line: 175, col: 23, offset: 3835},
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
	name: "DoubleLiteralRaw",
	pos: position{line: 177, col: 1, offset: 3843},
	expr: &actionExpr{
	pos: position{line: 177, col: 20, offset: 3864},
	run: (*parser).callonDoubleLiteralRaw1,
	expr: &seqExpr{
	pos: position{line: 177, col: 20, offset: 3864},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 177, col: 20, offset: 3864},
	expr: &charClassMatcher{
	pos: position{line: 177, col: 20, offset: 3864},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 177, col: 26, offset: 3870},
	expr: &charClassMatcher{
	pos: position{line: 177, col: 26, offset: 3870},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&choiceExpr{
	pos: position{line: 177, col: 35, offset: 3879},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 177, col: 35, offset: 3879},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 177, col: 35, offset: 3879},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 177, col: 39, offset: 3883},
	expr: &charClassMatcher{
	pos: position{line: 177, col: 39, offset: 3883},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&zeroOrOneExpr{
	pos: position{line: 177, col: 46, offset: 3890},
	expr: &ruleRefExpr{
	pos: position{line: 177, col: 46, offset: 3890},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 177, col: 58, offset: 3902},
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
	pos: position{line: 185, col: 1, offset: 4062},
	expr: &actionExpr{
	pos: position{line: 185, col: 17, offset: 4080},
	run: (*parser).callonDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 185, col: 17, offset: 4080},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 185, col: 17, offset: 4080},
	label: "d",
	expr: &ruleRefExpr{
	pos: position{line: 185, col: 19, offset: 4082},
	name: "DoubleLiteralRaw",
},
},
&ruleRefExpr{
	pos: position{line: 185, col: 36, offset: 4099},
	name: "_",
},
	},
},
},
},
{
	name: "NaturalLiteralRaw",
	pos: position{line: 187, col: 1, offset: 4120},
	expr: &actionExpr{
	pos: position{line: 187, col: 21, offset: 4142},
	run: (*parser).callonNaturalLiteralRaw1,
	expr: &oneOrMoreExpr{
	pos: position{line: 187, col: 21, offset: 4142},
	expr: &charClassMatcher{
	pos: position{line: 187, col: 21, offset: 4142},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
},
},
{
	name: "NaturalLiteral",
	pos: position{line: 192, col: 1, offset: 4235},
	expr: &actionExpr{
	pos: position{line: 192, col: 18, offset: 4254},
	run: (*parser).callonNaturalLiteral1,
	expr: &seqExpr{
	pos: position{line: 192, col: 18, offset: 4254},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 192, col: 18, offset: 4254},
	label: "n",
	expr: &ruleRefExpr{
	pos: position{line: 192, col: 20, offset: 4256},
	name: "NaturalLiteralRaw",
},
},
&ruleRefExpr{
	pos: position{line: 192, col: 38, offset: 4274},
	name: "_",
},
	},
},
},
},
{
	name: "IntegerLiteralRaw",
	pos: position{line: 194, col: 1, offset: 4295},
	expr: &actionExpr{
	pos: position{line: 194, col: 21, offset: 4317},
	run: (*parser).callonIntegerLiteralRaw1,
	expr: &seqExpr{
	pos: position{line: 194, col: 21, offset: 4317},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 194, col: 21, offset: 4317},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&oneOrMoreExpr{
	pos: position{line: 194, col: 25, offset: 4321},
	expr: &charClassMatcher{
	pos: position{line: 194, col: 25, offset: 4321},
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
	name: "IntegerLiteral",
	pos: position{line: 202, col: 1, offset: 4469},
	expr: &actionExpr{
	pos: position{line: 202, col: 18, offset: 4488},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 202, col: 18, offset: 4488},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 202, col: 18, offset: 4488},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 202, col: 21, offset: 4491},
	name: "IntegerLiteralRaw",
},
},
&ruleRefExpr{
	pos: position{line: 202, col: 39, offset: 4509},
	name: "_",
},
	},
},
},
},
{
	name: "DeBruijn",
	pos: position{line: 204, col: 1, offset: 4530},
	expr: &actionExpr{
	pos: position{line: 204, col: 12, offset: 4543},
	run: (*parser).callonDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 204, col: 12, offset: 4543},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 204, col: 12, offset: 4543},
	name: "At",
},
&labeledExpr{
	pos: position{line: 204, col: 15, offset: 4546},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 204, col: 21, offset: 4552},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Identifier",
	pos: position{line: 206, col: 1, offset: 4612},
	expr: &actionExpr{
	pos: position{line: 206, col: 14, offset: 4627},
	run: (*parser).callonIdentifier1,
	expr: &seqExpr{
	pos: position{line: 206, col: 14, offset: 4627},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 206, col: 14, offset: 4627},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 206, col: 19, offset: 4632},
	name: "Label",
},
},
&labeledExpr{
	pos: position{line: 206, col: 25, offset: 4638},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 206, col: 31, offset: 4644},
	expr: &ruleRefExpr{
	pos: position{line: 206, col: 31, offset: 4644},
	name: "DeBruijn",
},
},
},
	},
},
},
},
{
	name: "IdentifierReservedPrefix",
	pos: position{line: 214, col: 1, offset: 4815},
	expr: &actionExpr{
	pos: position{line: 215, col: 10, offset: 4853},
	run: (*parser).callonIdentifierReservedPrefix1,
	expr: &seqExpr{
	pos: position{line: 215, col: 10, offset: 4853},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 215, col: 10, offset: 4853},
	label: "name",
	expr: &actionExpr{
	pos: position{line: 215, col: 16, offset: 4859},
	run: (*parser).callonIdentifierReservedPrefix4,
	expr: &seqExpr{
	pos: position{line: 215, col: 16, offset: 4859},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 215, col: 16, offset: 4859},
	name: "ReservedRaw",
},
&oneOrMoreExpr{
	pos: position{line: 215, col: 28, offset: 4871},
	expr: &charClassMatcher{
	pos: position{line: 215, col: 28, offset: 4871},
	val: "[A-Za-z0-9/_-]",
	chars: []rune{'/','_','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
	},
},
},
},
&ruleRefExpr{
	pos: position{line: 215, col: 75, offset: 4918},
	name: "_",
},
&labeledExpr{
	pos: position{line: 216, col: 10, offset: 4929},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 216, col: 16, offset: 4935},
	expr: &ruleRefExpr{
	pos: position{line: 216, col: 16, offset: 4935},
	name: "DeBruijn",
},
},
},
	},
},
},
},
{
	name: "LetBinding",
	pos: position{line: 228, col: 1, offset: 5292},
	expr: &actionExpr{
	pos: position{line: 228, col: 14, offset: 5307},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 228, col: 14, offset: 5307},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 228, col: 14, offset: 5307},
	name: "Let",
},
&labeledExpr{
	pos: position{line: 228, col: 18, offset: 5311},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 228, col: 24, offset: 5317},
	name: "Label",
},
},
&labeledExpr{
	pos: position{line: 228, col: 30, offset: 5323},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 228, col: 32, offset: 5325},
	expr: &ruleRefExpr{
	pos: position{line: 228, col: 32, offset: 5325},
	name: "Annotation",
},
},
},
&ruleRefExpr{
	pos: position{line: 228, col: 44, offset: 5337},
	name: "Equal",
},
&labeledExpr{
	pos: position{line: 228, col: 50, offset: 5343},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 228, col: 52, offset: 5345},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 243, col: 1, offset: 5664},
	expr: &choiceExpr{
	pos: position{line: 244, col: 7, offset: 5685},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 244, col: 7, offset: 5685},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 244, col: 7, offset: 5685},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 244, col: 7, offset: 5685},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 244, col: 14, offset: 5692},
	name: "OpenParens",
},
&labeledExpr{
	pos: position{line: 244, col: 25, offset: 5703},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 244, col: 31, offset: 5709},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 244, col: 37, offset: 5715},
	name: "Colon",
},
&labeledExpr{
	pos: position{line: 244, col: 43, offset: 5721},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 244, col: 45, offset: 5723},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 244, col: 56, offset: 5734},
	name: "CloseParens",
},
&ruleRefExpr{
	pos: position{line: 244, col: 68, offset: 5746},
	name: "Arrow",
},
&labeledExpr{
	pos: position{line: 244, col: 74, offset: 5752},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 244, col: 79, offset: 5757},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 247, col: 7, offset: 5886},
	run: (*parser).callonExpression15,
	expr: &seqExpr{
	pos: position{line: 247, col: 7, offset: 5886},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 247, col: 7, offset: 5886},
	name: "If",
},
&labeledExpr{
	pos: position{line: 247, col: 10, offset: 5889},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 247, col: 15, offset: 5894},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 247, col: 26, offset: 5905},
	name: "Then",
},
&labeledExpr{
	pos: position{line: 247, col: 31, offset: 5910},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 247, col: 33, offset: 5912},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 247, col: 44, offset: 5923},
	name: "Else",
},
&labeledExpr{
	pos: position{line: 247, col: 49, offset: 5928},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 247, col: 51, offset: 5930},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 250, col: 7, offset: 6032},
	run: (*parser).callonExpression26,
	expr: &seqExpr{
	pos: position{line: 250, col: 7, offset: 6032},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 250, col: 7, offset: 6032},
	label: "bindings",
	expr: &zeroOrMoreExpr{
	pos: position{line: 250, col: 16, offset: 6041},
	expr: &ruleRefExpr{
	pos: position{line: 250, col: 16, offset: 6041},
	name: "LetBinding",
},
},
},
&ruleRefExpr{
	pos: position{line: 250, col: 28, offset: 6053},
	name: "In",
},
&labeledExpr{
	pos: position{line: 250, col: 31, offset: 6056},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 250, col: 33, offset: 6058},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 257, col: 7, offset: 6314},
	run: (*parser).callonExpression34,
	expr: &seqExpr{
	pos: position{line: 257, col: 7, offset: 6314},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 257, col: 7, offset: 6314},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 257, col: 14, offset: 6321},
	name: "OpenParens",
},
&labeledExpr{
	pos: position{line: 257, col: 25, offset: 6332},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 257, col: 31, offset: 6338},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 257, col: 37, offset: 6344},
	name: "Colon",
},
&labeledExpr{
	pos: position{line: 257, col: 43, offset: 6350},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 257, col: 45, offset: 6352},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 257, col: 56, offset: 6363},
	name: "CloseParens",
},
&ruleRefExpr{
	pos: position{line: 257, col: 68, offset: 6375},
	name: "Arrow",
},
&labeledExpr{
	pos: position{line: 257, col: 74, offset: 6381},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 257, col: 79, offset: 6386},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 260, col: 7, offset: 6507},
	run: (*parser).callonExpression47,
	expr: &seqExpr{
	pos: position{line: 260, col: 7, offset: 6507},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 260, col: 7, offset: 6507},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 260, col: 9, offset: 6509},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 260, col: 28, offset: 6528},
	name: "Arrow",
},
&labeledExpr{
	pos: position{line: 260, col: 34, offset: 6534},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 260, col: 36, offset: 6536},
	name: "Expression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 261, col: 7, offset: 6608},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 263, col: 1, offset: 6629},
	expr: &actionExpr{
	pos: position{line: 263, col: 14, offset: 6644},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 263, col: 14, offset: 6644},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 263, col: 14, offset: 6644},
	name: "Colon",
},
&labeledExpr{
	pos: position{line: 263, col: 20, offset: 6650},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 263, col: 22, offset: 6652},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 265, col: 1, offset: 6682},
	expr: &choiceExpr{
	pos: position{line: 266, col: 5, offset: 6710},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 266, col: 5, offset: 6710},
	name: "EmptyList",
},
&actionExpr{
	pos: position{line: 267, col: 5, offset: 6724},
	run: (*parser).callonAnnotatedExpression3,
	expr: &seqExpr{
	pos: position{line: 267, col: 5, offset: 6724},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 267, col: 5, offset: 6724},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 267, col: 7, offset: 6726},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 267, col: 26, offset: 6745},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 267, col: 28, offset: 6747},
	expr: &ruleRefExpr{
	pos: position{line: 267, col: 28, offset: 6747},
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
	pos: position{line: 272, col: 1, offset: 6864},
	expr: &actionExpr{
	pos: position{line: 272, col: 13, offset: 6878},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 272, col: 13, offset: 6878},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 272, col: 13, offset: 6878},
	name: "OpenBracket",
},
&ruleRefExpr{
	pos: position{line: 272, col: 25, offset: 6890},
	name: "CloseBracket",
},
&ruleRefExpr{
	pos: position{line: 272, col: 38, offset: 6903},
	name: "Colon",
},
&ruleRefExpr{
	pos: position{line: 272, col: 44, offset: 6909},
	name: "List",
},
&labeledExpr{
	pos: position{line: 272, col: 49, offset: 6914},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 272, col: 51, offset: 6916},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 276, col: 1, offset: 6987},
	expr: &ruleRefExpr{
	pos: position{line: 276, col: 22, offset: 7010},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 278, col: 1, offset: 7031},
	expr: &ruleRefExpr{
	pos: position{line: 278, col: 23, offset: 7055},
	name: "PlusExpression",
},
},
{
	name: "MorePlus",
	pos: position{line: 280, col: 1, offset: 7071},
	expr: &actionExpr{
	pos: position{line: 280, col: 12, offset: 7084},
	run: (*parser).callonMorePlus1,
	expr: &seqExpr{
	pos: position{line: 280, col: 12, offset: 7084},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 280, col: 12, offset: 7084},
	name: "Plus",
},
&labeledExpr{
	pos: position{line: 280, col: 17, offset: 7089},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 280, col: 19, offset: 7091},
	name: "TimesExpression",
},
},
	},
},
},
},
{
	name: "PlusExpression",
	pos: position{line: 281, col: 1, offset: 7125},
	expr: &actionExpr{
	pos: position{line: 282, col: 7, offset: 7150},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 282, col: 7, offset: 7150},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 282, col: 7, offset: 7150},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 282, col: 13, offset: 7156},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 282, col: 29, offset: 7172},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 282, col: 34, offset: 7177},
	expr: &ruleRefExpr{
	pos: position{line: 282, col: 34, offset: 7177},
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
	pos: position{line: 291, col: 1, offset: 7417},
	expr: &actionExpr{
	pos: position{line: 291, col: 13, offset: 7431},
	run: (*parser).callonMoreTimes1,
	expr: &seqExpr{
	pos: position{line: 291, col: 13, offset: 7431},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 291, col: 13, offset: 7431},
	name: "Times",
},
&labeledExpr{
	pos: position{line: 291, col: 19, offset: 7437},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 291, col: 21, offset: 7439},
	name: "ApplicationExpression",
},
},
	},
},
},
},
{
	name: "TimesExpression",
	pos: position{line: 292, col: 1, offset: 7479},
	expr: &actionExpr{
	pos: position{line: 293, col: 7, offset: 7505},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 293, col: 7, offset: 7505},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 293, col: 7, offset: 7505},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 293, col: 13, offset: 7511},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 293, col: 35, offset: 7533},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 293, col: 40, offset: 7538},
	expr: &ruleRefExpr{
	pos: position{line: 293, col: 40, offset: 7538},
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
	pos: position{line: 302, col: 1, offset: 7780},
	expr: &actionExpr{
	pos: position{line: 302, col: 25, offset: 7806},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 302, col: 25, offset: 7806},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 302, col: 25, offset: 7806},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 302, col: 27, offset: 7808},
	name: "ImportExpression",
},
},
&labeledExpr{
	pos: position{line: 302, col: 44, offset: 7825},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 302, col: 49, offset: 7830},
	expr: &ruleRefExpr{
	pos: position{line: 302, col: 49, offset: 7830},
	name: "ImportExpression",
},
},
},
	},
},
},
},
{
	name: "ImportExpression",
	pos: position{line: 311, col: 1, offset: 8072},
	expr: &ruleRefExpr{
	pos: position{line: 311, col: 20, offset: 8093},
	name: "SelectorExpression",
},
},
{
	name: "SelectorExpression",
	pos: position{line: 313, col: 1, offset: 8113},
	expr: &actionExpr{
	pos: position{line: 313, col: 22, offset: 8136},
	run: (*parser).callonSelectorExpression1,
	expr: &seqExpr{
	pos: position{line: 313, col: 22, offset: 8136},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 313, col: 22, offset: 8136},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 313, col: 24, offset: 8138},
	name: "PrimitiveExpression",
},
},
&labeledExpr{
	pos: position{line: 313, col: 44, offset: 8158},
	label: "ls",
	expr: &zeroOrMoreExpr{
	pos: position{line: 313, col: 47, offset: 8161},
	expr: &seqExpr{
	pos: position{line: 313, col: 48, offset: 8162},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 313, col: 48, offset: 8162},
	name: "Dot",
},
&ruleRefExpr{
	pos: position{line: 313, col: 52, offset: 8166},
	name: "Label",
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
	pos: position{line: 323, col: 1, offset: 8404},
	expr: &choiceExpr{
	pos: position{line: 324, col: 7, offset: 8434},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 324, col: 7, offset: 8434},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 325, col: 7, offset: 8454},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 326, col: 7, offset: 8475},
	name: "IntegerLiteral",
},
&actionExpr{
	pos: position{line: 327, col: 7, offset: 8496},
	run: (*parser).callonPrimitiveExpression5,
	expr: &litMatcher{
	pos: position{line: 327, col: 7, offset: 8496},
	val: "-Infinity",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 328, col: 7, offset: 8558},
	name: "TextLiteral",
},
&actionExpr{
	pos: position{line: 329, col: 7, offset: 8576},
	run: (*parser).callonPrimitiveExpression8,
	expr: &seqExpr{
	pos: position{line: 329, col: 7, offset: 8576},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 329, col: 7, offset: 8576},
	name: "OpenBrace",
},
&labeledExpr{
	pos: position{line: 329, col: 17, offset: 8586},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 329, col: 19, offset: 8588},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 329, col: 39, offset: 8608},
	name: "CloseBrace",
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 330, col: 7, offset: 8643},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 331, col: 7, offset: 8669},
	name: "IdentifierReservedPrefix",
},
&ruleRefExpr{
	pos: position{line: 332, col: 7, offset: 8700},
	name: "Reserved",
},
&ruleRefExpr{
	pos: position{line: 333, col: 7, offset: 8715},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 334, col: 7, offset: 8732},
	run: (*parser).callonPrimitiveExpression18,
	expr: &seqExpr{
	pos: position{line: 334, col: 7, offset: 8732},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 334, col: 7, offset: 8732},
	name: "OpenParens",
},
&labeledExpr{
	pos: position{line: 334, col: 18, offset: 8743},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 334, col: 20, offset: 8745},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 334, col: 31, offset: 8756},
	name: "CloseParens",
},
	},
},
},
	},
},
},
{
	name: "RecordTypeOrLiteral",
	pos: position{line: 336, col: 1, offset: 8787},
	expr: &choiceExpr{
	pos: position{line: 337, col: 7, offset: 8817},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 337, col: 7, offset: 8817},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &ruleRefExpr{
	pos: position{line: 337, col: 7, offset: 8817},
	name: "Equal",
},
},
&ruleRefExpr{
	pos: position{line: 338, col: 7, offset: 8882},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 339, col: 7, offset: 8907},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 340, col: 7, offset: 8935},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 340, col: 7, offset: 8935},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 342, col: 1, offset: 8989},
	expr: &actionExpr{
	pos: position{line: 342, col: 19, offset: 9009},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 342, col: 19, offset: 9009},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 342, col: 19, offset: 9009},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 342, col: 24, offset: 9014},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 342, col: 30, offset: 9020},
	name: "Colon",
},
&labeledExpr{
	pos: position{line: 342, col: 36, offset: 9026},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 342, col: 41, offset: 9031},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 345, col: 1, offset: 9088},
	expr: &actionExpr{
	pos: position{line: 345, col: 18, offset: 9107},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 345, col: 18, offset: 9107},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 345, col: 18, offset: 9107},
	name: "Comma",
},
&labeledExpr{
	pos: position{line: 345, col: 24, offset: 9113},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 345, col: 26, offset: 9115},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 346, col: 1, offset: 9147},
	expr: &actionExpr{
	pos: position{line: 347, col: 7, offset: 9176},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 347, col: 7, offset: 9176},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 347, col: 7, offset: 9176},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 347, col: 13, offset: 9182},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 347, col: 29, offset: 9198},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 347, col: 34, offset: 9203},
	expr: &ruleRefExpr{
	pos: position{line: 347, col: 34, offset: 9203},
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
	pos: position{line: 357, col: 1, offset: 9615},
	expr: &actionExpr{
	pos: position{line: 357, col: 22, offset: 9638},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 357, col: 22, offset: 9638},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 357, col: 22, offset: 9638},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 357, col: 27, offset: 9643},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 357, col: 33, offset: 9649},
	name: "Equal",
},
&labeledExpr{
	pos: position{line: 357, col: 39, offset: 9655},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 357, col: 44, offset: 9660},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 360, col: 1, offset: 9717},
	expr: &actionExpr{
	pos: position{line: 360, col: 21, offset: 9739},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 360, col: 21, offset: 9739},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 360, col: 21, offset: 9739},
	name: "Comma",
},
&labeledExpr{
	pos: position{line: 360, col: 27, offset: 9745},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 360, col: 29, offset: 9747},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 361, col: 1, offset: 9782},
	expr: &actionExpr{
	pos: position{line: 362, col: 7, offset: 9814},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 362, col: 7, offset: 9814},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 362, col: 7, offset: 9814},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 362, col: 13, offset: 9820},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 362, col: 32, offset: 9839},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 362, col: 37, offset: 9844},
	expr: &ruleRefExpr{
	pos: position{line: 362, col: 37, offset: 9844},
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
	pos: position{line: 372, col: 1, offset: 10262},
	expr: &actionExpr{
	pos: position{line: 372, col: 12, offset: 10275},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 372, col: 12, offset: 10275},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 372, col: 12, offset: 10275},
	name: "Comma",
},
&labeledExpr{
	pos: position{line: 372, col: 18, offset: 10281},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 372, col: 20, offset: 10283},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 374, col: 1, offset: 10311},
	expr: &actionExpr{
	pos: position{line: 375, col: 7, offset: 10341},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 375, col: 7, offset: 10341},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 375, col: 7, offset: 10341},
	name: "OpenBracket",
},
&labeledExpr{
	pos: position{line: 375, col: 19, offset: 10353},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 375, col: 25, offset: 10359},
	name: "Expression",
},
},
&labeledExpr{
	pos: position{line: 375, col: 36, offset: 10370},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 375, col: 41, offset: 10375},
	expr: &ruleRefExpr{
	pos: position{line: 375, col: 41, offset: 10375},
	name: "MoreList",
},
},
},
&ruleRefExpr{
	pos: position{line: 375, col: 51, offset: 10385},
	name: "CloseBracket",
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 385, col: 1, offset: 10686},
	expr: &notExpr{
	pos: position{line: 385, col: 7, offset: 10694},
	expr: &anyMatcher{
	line: 385, col: 8, offset: 10695,
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

func (c *current) onSimpleLabel1() (interface{}, error) {
 return string(c.text), nil 
}

func (p *parser) callonSimpleLabel1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSimpleLabel1()
}

func (c *current) onLabel1(label interface{}) (interface{}, error) {
 return label, nil 
}

func (p *parser) callonLabel1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLabel1(stack["label"])
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
    var outChunks ast.Chunks
    for _, chunk := range chunks.([]interface{}) {
        switch e := chunk.(type) {
        case []byte:
                str.Write(e)
        case ast.Expr:
                outChunks = append(outChunks, ast.Chunk{str.String(), e})
                str.Reset()
        default:
                return nil, errors.New("can't happen")
        }
    }
    return ast.TextLit{Chunks: outChunks, Suffix: str.String()}, nil
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

func (c *current) onReservedRaw2() (interface{}, error) {
 return ast.Bool, nil 
}

func (p *parser) callonReservedRaw2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw2()
}

func (c *current) onReservedRaw6() (interface{}, error) {
 return ast.Natural, nil 
}

func (p *parser) callonReservedRaw6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw6()
}

func (c *current) onReservedRaw8() (interface{}, error) {
 return ast.Integer, nil 
}

func (p *parser) callonReservedRaw8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw8()
}

func (c *current) onReservedRaw10() (interface{}, error) {
 return ast.Double, nil 
}

func (p *parser) callonReservedRaw10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw10()
}

func (c *current) onReservedRaw12() (interface{}, error) {
 return ast.Text, nil 
}

func (p *parser) callonReservedRaw12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw12()
}

func (c *current) onReservedRaw14() (interface{}, error) {
 return ast.List, nil 
}

func (p *parser) callonReservedRaw14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw14()
}

func (c *current) onReservedRaw16() (interface{}, error) {
 return ast.True, nil 
}

func (p *parser) callonReservedRaw16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw16()
}

func (c *current) onReservedRaw18() (interface{}, error) {
 return ast.False, nil 
}

func (p *parser) callonReservedRaw18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw18()
}

func (c *current) onReservedRaw20() (interface{}, error) {
 return ast.DoubleLit(math.NaN()), nil 
}

func (p *parser) callonReservedRaw20() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw20()
}

func (c *current) onReservedRaw22() (interface{}, error) {
 return ast.DoubleLit(math.Inf(1)), nil 
}

func (p *parser) callonReservedRaw22() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw22()
}

func (c *current) onReservedRaw24() (interface{}, error) {
 return ast.Type, nil 
}

func (p *parser) callonReservedRaw24() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw24()
}

func (c *current) onReservedRaw26() (interface{}, error) {
 return ast.Kind, nil 
}

func (p *parser) callonReservedRaw26() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw26()
}

func (c *current) onReservedRaw28() (interface{}, error) {
 return ast.Sort, nil 
}

func (p *parser) callonReservedRaw28() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw28()
}

func (c *current) onReserved1(r interface{}) (interface{}, error) {
 return r, nil 
}

func (p *parser) callonReserved1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved1(stack["r"])
}

func (c *current) onDoubleLiteralRaw1() (interface{}, error) {
      d, err := strconv.ParseFloat(string(c.text), 64)
      if err != nil {
         return nil, err
      }
      return ast.DoubleLit(d), nil
}

func (p *parser) callonDoubleLiteralRaw1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleLiteralRaw1()
}

func (c *current) onDoubleLiteral1(d interface{}) (interface{}, error) {
 return d, nil 
}

func (p *parser) callonDoubleLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleLiteral1(stack["d"])
}

func (c *current) onNaturalLiteralRaw1() (interface{}, error) {
      i, err := strconv.Atoi(string(c.text))
      return ast.NaturalLit(i), err
}

func (p *parser) callonNaturalLiteralRaw1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNaturalLiteralRaw1()
}

func (c *current) onNaturalLiteral1(n interface{}) (interface{}, error) {
 return n, nil 
}

func (p *parser) callonNaturalLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNaturalLiteral1(stack["n"])
}

func (c *current) onIntegerLiteralRaw1() (interface{}, error) {
      i, err := strconv.Atoi(string(c.text))
      if err != nil {
         return nil, err
      }
      return ast.IntegerLit(i), nil
}

func (p *parser) callonIntegerLiteralRaw1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIntegerLiteralRaw1()
}

func (c *current) onIntegerLiteral1(i interface{}) (interface{}, error) {
 return i, nil 
}

func (p *parser) callonIntegerLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIntegerLiteral1(stack["i"])
}

func (c *current) onDeBruijn1(index interface{}) (interface{}, error) {
 return int(index.(ast.NaturalLit)), nil 
}

func (p *parser) callonDeBruijn1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDeBruijn1(stack["index"])
}

func (c *current) onIdentifier1(name, index interface{}) (interface{}, error) {
    if index != nil {
        return ast.Var{Name:name.(string), Index:index.(int)}, nil
    } else {
        return ast.Var{Name:name.(string)}, nil
    }
}

func (p *parser) callonIdentifier1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifier1(stack["name"], stack["index"])
}

func (c *current) onIdentifierReservedPrefix4() (interface{}, error) {
 return string(c.text),nil 
}

func (p *parser) callonIdentifierReservedPrefix4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifierReservedPrefix4()
}

func (c *current) onIdentifierReservedPrefix1(name, index interface{}) (interface{}, error) {
    if index != nil {
        return ast.Var{Name:name.(string), Index:index.(int)}, nil
    } else {
        return ast.Var{Name:name.(string)}, nil
    }
}

func (p *parser) callonIdentifierReservedPrefix1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifierReservedPrefix1(stack["name"], stack["index"])
}

func (c *current) onLetBinding1(label, a, v interface{}) (interface{}, error) {
    if a != nil {
        return ast.Binding{
            Variable: label.(string),
            Annotation: a.(ast.Expr),
            Value: v.(ast.Expr),
        }, nil
    } else {
        return ast.Binding{
            Variable: label.(string),
            Value: v.(ast.Expr),
        }, nil
    }
}

func (p *parser) callonLetBinding1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLetBinding1(stack["label"], stack["a"], stack["v"])
}

func (c *current) onExpression2(label, t, body interface{}) (interface{}, error) {
          return &ast.LambdaExpr{Label:label.(string), Type:t.(ast.Expr), Body: body.(ast.Expr)}, nil
      
}

func (p *parser) callonExpression2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression2(stack["label"], stack["t"], stack["body"])
}

func (c *current) onExpression15(cond, t, f interface{}) (interface{}, error) {
          return ast.BoolIf{cond.(ast.Expr),t.(ast.Expr),f.(ast.Expr)},nil
      
}

func (p *parser) callonExpression15() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression15(stack["cond"], stack["t"], stack["f"])
}

func (c *current) onExpression26(bindings, b interface{}) (interface{}, error) {
        bs := make([]ast.Binding, len(bindings.([]interface{})))
        for i, binding := range bindings.([]interface{}) {
            bs[i] = binding.(ast.Binding)
        }
        return ast.MakeLet(b.(ast.Expr), bs...), nil
      
}

func (p *parser) callonExpression26() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression26(stack["bindings"], stack["b"])
}

func (c *current) onExpression34(label, t, body interface{}) (interface{}, error) {
          return &ast.Pi{Label:label.(string), Type:t.(ast.Expr), Body: body.(ast.Expr)}, nil
      
}

func (p *parser) callonExpression34() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression34(stack["label"], stack["t"], stack["body"])
}

func (c *current) onExpression47(o, e interface{}) (interface{}, error) {
 return &ast.Pi{"_",o.(ast.Expr),e.(ast.Expr)}, nil 
}

func (p *parser) callonExpression47() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression47(stack["o"], stack["e"])
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
        return ast.Annot{e.(ast.Expr), a.(ast.Expr)}, nil
    
}

func (p *parser) callonAnnotatedExpression3() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAnnotatedExpression3(stack["e"], stack["a"])
}

func (c *current) onEmptyList1(a interface{}) (interface{}, error) {
          return ast.EmptyList{a.(ast.Expr)},nil
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
          a := first.(ast.Expr)
          if rest == nil { return a, nil }
          for _, b := range rest.([]interface{}) {
              a = ast.NaturalPlus{L: a, R: b.(ast.Expr)}
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
          a := first.(ast.Expr)
          if rest == nil { return a, nil }
          for _, b := range rest.([]interface{}) {
              a = ast.NaturalTimes{L: a, R: b.(ast.Expr)}
          }
          return a, nil
      
}

func (p *parser) callonTimesExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTimesExpression1(stack["first"], stack["rest"])
}

func (c *current) onApplicationExpression1(f, rest interface{}) (interface{}, error) {
          e := f.(ast.Expr)
          if rest == nil { return e, nil }
          for _, arg := range rest.([]interface{}) {
              e = &ast.App{Fn:e, Arg: arg.(ast.Expr)}
          }
          return e,nil
      
}

func (p *parser) callonApplicationExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onApplicationExpression1(stack["f"], stack["rest"])
}

func (c *current) onSelectorExpression1(e, ls interface{}) (interface{}, error) {
    expr := e.(ast.Expr)
    labels := ls.([]interface{})
    for _, labelSelector := range labels {
        label := labelSelector.([]interface{})[1]
        expr = ast.Field{expr, label.(string)}
    }
    return expr, nil
}

func (p *parser) callonSelectorExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSelectorExpression1(stack["e"], stack["ls"])
}

func (c *current) onPrimitiveExpression5() (interface{}, error) {
 return ast.DoubleLit(math.Inf(-1)), nil 
}

func (p *parser) callonPrimitiveExpression5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimitiveExpression5()
}

func (c *current) onPrimitiveExpression8(r interface{}) (interface{}, error) {
 return r, nil 
}

func (p *parser) callonPrimitiveExpression8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimitiveExpression8(stack["r"])
}

func (c *current) onPrimitiveExpression18(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonPrimitiveExpression18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimitiveExpression18(stack["e"])
}

func (c *current) onRecordTypeOrLiteral2() (interface{}, error) {
 return ast.RecordLit(map[string]ast.Expr{}), nil 
}

func (p *parser) callonRecordTypeOrLiteral2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRecordTypeOrLiteral2()
}

func (c *current) onRecordTypeOrLiteral6() (interface{}, error) {
 return ast.Record(map[string]ast.Expr{}), nil 
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
          content := make(map[string]ast.Expr, len(fields)+1)
          content[first.([]interface{})[0].(string)] = first.([]interface{})[1].(ast.Expr)
          for _, field := range(fields) {
              content[field.([]interface{})[0].(string)] = field.([]interface{})[1].(ast.Expr)
          }
          return ast.Record(content), nil
      
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
          content := make(map[string]ast.Expr, len(fields)+1)
          content[first.([]interface{})[0].(string)] = first.([]interface{})[1].(ast.Expr)
          for _, field := range(fields) {
              content[field.([]interface{})[0].(string)] = field.([]interface{})[1].(ast.Expr)
          }
          return ast.RecordLit(content), nil
      
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
          content := make([]ast.Expr, len(exprs)+1)
          content[0] = first.(ast.Expr)
          for i, expr := range(exprs) {
              content[i+1] = expr.(ast.Expr)
          }
          return ast.NonEmptyList(content), nil
      
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

