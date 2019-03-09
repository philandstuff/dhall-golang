
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
&ruleRefExpr{
	pos: position{line: 23, col: 37, offset: 276},
	name: "_",
},
	},
},
},
},
{
	name: "EOL",
	pos: position{line: 25, col: 1, offset: 297},
	expr: &choiceExpr{
	pos: position{line: 25, col: 7, offset: 305},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 25, col: 7, offset: 305},
	val: "\n",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 25, col: 14, offset: 312},
	val: "\r\n",
	ignoreCase: false,
},
	},
},
},
{
	name: "BlockComment",
	pos: position{line: 27, col: 1, offset: 320},
	expr: &seqExpr{
	pos: position{line: 27, col: 16, offset: 337},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 27, col: 16, offset: 337},
	val: "{-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 27, col: 21, offset: 342},
	name: "BlockCommentContinue",
},
	},
},
},
{
	name: "BlockCommentChunk",
	pos: position{line: 29, col: 1, offset: 364},
	expr: &choiceExpr{
	pos: position{line: 30, col: 5, offset: 390},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 30, col: 5, offset: 390},
	name: "BlockComment",
},
&charClassMatcher{
	pos: position{line: 31, col: 5, offset: 407},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 32, col: 5, offset: 433},
	name: "EOL",
},
	},
},
},
{
	name: "BlockCommentContinue",
	pos: position{line: 34, col: 1, offset: 438},
	expr: &choiceExpr{
	pos: position{line: 34, col: 24, offset: 463},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 34, col: 24, offset: 463},
	val: "-}",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 34, col: 31, offset: 470},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 34, col: 31, offset: 470},
	name: "BlockCommentChunk",
},
&ruleRefExpr{
	pos: position{line: 34, col: 49, offset: 488},
	name: "BlockCommentContinue",
},
	},
},
	},
},
},
{
	name: "NotEOL",
	pos: position{line: 36, col: 1, offset: 510},
	expr: &charClassMatcher{
	pos: position{line: 36, col: 10, offset: 521},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "LineComment",
	pos: position{line: 38, col: 1, offset: 544},
	expr: &actionExpr{
	pos: position{line: 38, col: 15, offset: 560},
	run: (*parser).callonLineComment1,
	expr: &seqExpr{
	pos: position{line: 38, col: 15, offset: 560},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 38, col: 15, offset: 560},
	val: "--",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 38, col: 20, offset: 565},
	label: "content",
	expr: &actionExpr{
	pos: position{line: 38, col: 29, offset: 574},
	run: (*parser).callonLineComment5,
	expr: &zeroOrMoreExpr{
	pos: position{line: 38, col: 29, offset: 574},
	expr: &ruleRefExpr{
	pos: position{line: 38, col: 29, offset: 574},
	name: "NotEOL",
},
},
},
},
&ruleRefExpr{
	pos: position{line: 38, col: 68, offset: 613},
	name: "EOL",
},
	},
},
},
},
{
	name: "WhitespaceChunk",
	pos: position{line: 40, col: 1, offset: 642},
	expr: &choiceExpr{
	pos: position{line: 40, col: 19, offset: 662},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 40, col: 19, offset: 662},
	val: " ",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 40, col: 25, offset: 668},
	val: "\t",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 40, col: 32, offset: 675},
	name: "EOL",
},
&ruleRefExpr{
	pos: position{line: 40, col: 38, offset: 681},
	name: "LineComment",
},
&ruleRefExpr{
	pos: position{line: 40, col: 52, offset: 695},
	name: "BlockComment",
},
	},
},
},
{
	name: "_",
	pos: position{line: 42, col: 1, offset: 709},
	expr: &zeroOrMoreExpr{
	pos: position{line: 42, col: 5, offset: 715},
	expr: &ruleRefExpr{
	pos: position{line: 42, col: 5, offset: 715},
	name: "WhitespaceChunk",
},
},
},
{
	name: "NonemptyWhitespace",
	pos: position{line: 44, col: 1, offset: 733},
	expr: &oneOrMoreExpr{
	pos: position{line: 44, col: 22, offset: 756},
	expr: &ruleRefExpr{
	pos: position{line: 44, col: 22, offset: 756},
	name: "WhitespaceChunk",
},
},
},
{
	name: "HexDig",
	pos: position{line: 46, col: 1, offset: 774},
	expr: &charClassMatcher{
	pos: position{line: 46, col: 10, offset: 785},
	val: "[0-9a-f]i",
	ranges: []rune{'0','9','a','f',},
	ignoreCase: true,
	inverted: false,
},
},
{
	name: "SimpleLabel",
	pos: position{line: 48, col: 1, offset: 796},
	expr: &actionExpr{
	pos: position{line: 48, col: 15, offset: 812},
	run: (*parser).callonSimpleLabel1,
	expr: &seqExpr{
	pos: position{line: 48, col: 15, offset: 812},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 48, col: 15, offset: 812},
	expr: &ruleRefExpr{
	pos: position{line: 48, col: 16, offset: 813},
	name: "Keyword",
},
},
&charClassMatcher{
	pos: position{line: 49, col: 13, offset: 833},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 49, col: 23, offset: 843},
	expr: &charClassMatcher{
	pos: position{line: 49, col: 23, offset: 843},
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
	pos: position{line: 53, col: 1, offset: 907},
	expr: &actionExpr{
	pos: position{line: 53, col: 9, offset: 917},
	run: (*parser).callonLabel1,
	expr: &labeledExpr{
	pos: position{line: 53, col: 9, offset: 917},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 53, col: 15, offset: 923},
	name: "SimpleLabel",
},
},
},
},
{
	name: "EscapedChar",
	pos: position{line: 57, col: 1, offset: 991},
	expr: &actionExpr{
	pos: position{line: 58, col: 3, offset: 1009},
	run: (*parser).callonEscapedChar1,
	expr: &seqExpr{
	pos: position{line: 58, col: 3, offset: 1009},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 58, col: 3, offset: 1009},
	val: "\\",
	ignoreCase: false,
},
&choiceExpr{
	pos: position{line: 59, col: 5, offset: 1018},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 59, col: 5, offset: 1018},
	val: "\"",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 60, col: 10, offset: 1031},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 61, col: 10, offset: 1044},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 62, col: 10, offset: 1058},
	val: "/",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 63, col: 10, offset: 1071},
	val: "b",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 64, col: 10, offset: 1084},
	val: "f",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 65, col: 10, offset: 1097},
	val: "n",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 66, col: 10, offset: 1110},
	val: "r",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 67, col: 10, offset: 1123},
	val: "t",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 68, col: 10, offset: 1136},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 68, col: 10, offset: 1136},
	val: "u",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 68, col: 14, offset: 1140},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 68, col: 21, offset: 1147},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 68, col: 28, offset: 1154},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 68, col: 35, offset: 1161},
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
	pos: position{line: 89, col: 1, offset: 1604},
	expr: &choiceExpr{
	pos: position{line: 90, col: 6, offset: 1630},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 90, col: 6, offset: 1630},
	run: (*parser).callonDoubleQuoteChunk2,
	expr: &seqExpr{
	pos: position{line: 90, col: 6, offset: 1630},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 90, col: 6, offset: 1630},
	val: "${",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 90, col: 11, offset: 1635},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 90, col: 13, offset: 1637},
	name: "CompleteExpression",
},
},
&litMatcher{
	pos: position{line: 90, col: 32, offset: 1656},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 91, col: 6, offset: 1683},
	name: "EscapedChar",
},
&charClassMatcher{
	pos: position{line: 92, col: 6, offset: 1700},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 93, col: 6, offset: 1717},
	val: "[\\x23-\\x5b]",
	ranges: []rune{'#','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 94, col: 6, offset: 1734},
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
	pos: position{line: 96, col: 1, offset: 1753},
	expr: &actionExpr{
	pos: position{line: 96, col: 22, offset: 1776},
	run: (*parser).callonDoubleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 96, col: 22, offset: 1776},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 96, col: 22, offset: 1776},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 96, col: 26, offset: 1780},
	label: "chunks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 96, col: 33, offset: 1787},
	expr: &ruleRefExpr{
	pos: position{line: 96, col: 33, offset: 1787},
	name: "DoubleQuoteChunk",
},
},
},
&litMatcher{
	pos: position{line: 96, col: 51, offset: 1805},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TextLiteral",
	pos: position{line: 113, col: 1, offset: 2289},
	expr: &ruleRefExpr{
	pos: position{line: 113, col: 15, offset: 2305},
	name: "DoubleQuoteLiteral",
},
},
{
	name: "Reserved",
	pos: position{line: 115, col: 1, offset: 2325},
	expr: &choiceExpr{
	pos: position{line: 115, col: 12, offset: 2338},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 115, col: 12, offset: 2338},
	run: (*parser).callonReserved2,
	expr: &litMatcher{
	pos: position{line: 115, col: 12, offset: 2338},
	val: "Bool",
	ignoreCase: false,
},
},
&litMatcher{
	pos: position{line: 116, col: 5, offset: 2374},
	val: "Optional",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 117, col: 5, offset: 2389},
	val: "None",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 118, col: 5, offset: 2400},
	run: (*parser).callonReserved6,
	expr: &litMatcher{
	pos: position{line: 118, col: 5, offset: 2400},
	val: "Natural",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 119, col: 5, offset: 2442},
	run: (*parser).callonReserved8,
	expr: &litMatcher{
	pos: position{line: 119, col: 5, offset: 2442},
	val: "Integer",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 120, col: 5, offset: 2484},
	run: (*parser).callonReserved10,
	expr: &litMatcher{
	pos: position{line: 120, col: 5, offset: 2484},
	val: "Double",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 121, col: 5, offset: 2524},
	run: (*parser).callonReserved12,
	expr: &litMatcher{
	pos: position{line: 121, col: 5, offset: 2524},
	val: "Text",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 122, col: 5, offset: 2560},
	run: (*parser).callonReserved14,
	expr: &litMatcher{
	pos: position{line: 122, col: 5, offset: 2560},
	val: "List",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 123, col: 5, offset: 2596},
	run: (*parser).callonReserved16,
	expr: &litMatcher{
	pos: position{line: 123, col: 5, offset: 2596},
	val: "True",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 124, col: 5, offset: 2632},
	run: (*parser).callonReserved18,
	expr: &litMatcher{
	pos: position{line: 124, col: 5, offset: 2632},
	val: "False",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 125, col: 5, offset: 2670},
	run: (*parser).callonReserved20,
	expr: &litMatcher{
	pos: position{line: 125, col: 5, offset: 2670},
	val: "NaN",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 126, col: 5, offset: 2722},
	run: (*parser).callonReserved22,
	expr: &litMatcher{
	pos: position{line: 126, col: 5, offset: 2722},
	val: "Infinity",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 127, col: 5, offset: 2780},
	run: (*parser).callonReserved24,
	expr: &litMatcher{
	pos: position{line: 127, col: 5, offset: 2780},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 128, col: 5, offset: 2816},
	run: (*parser).callonReserved26,
	expr: &litMatcher{
	pos: position{line: 128, col: 5, offset: 2816},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 129, col: 5, offset: 2852},
	run: (*parser).callonReserved28,
	expr: &litMatcher{
	pos: position{line: 129, col: 5, offset: 2852},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "Keyword",
	pos: position{line: 131, col: 1, offset: 2885},
	expr: &choiceExpr{
	pos: position{line: 131, col: 11, offset: 2897},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 131, col: 11, offset: 2897},
	val: "if",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 132, col: 5, offset: 2906},
	val: "then",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 133, col: 5, offset: 2917},
	val: "else",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 134, col: 5, offset: 2928},
	val: "let",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 135, col: 5, offset: 2938},
	val: "in",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 136, col: 5, offset: 2947},
	val: "as",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 137, col: 5, offset: 2956},
	val: "using",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 138, col: 5, offset: 2968},
	val: "merge",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 139, col: 5, offset: 2980},
	val: "constructors",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 140, col: 5, offset: 2999},
	val: "Some",
	ignoreCase: false,
},
	},
},
},
{
	name: "ColonSpace",
	pos: position{line: 142, col: 1, offset: 3007},
	expr: &seqExpr{
	pos: position{line: 142, col: 14, offset: 3022},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 142, col: 14, offset: 3022},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 142, col: 18, offset: 3026},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Lambda",
	pos: position{line: 144, col: 1, offset: 3046},
	expr: &choiceExpr{
	pos: position{line: 144, col: 11, offset: 3058},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 144, col: 11, offset: 3058},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 144, col: 18, offset: 3065},
	val: "λ",
	ignoreCase: false,
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 145, col: 1, offset: 3071},
	expr: &choiceExpr{
	pos: position{line: 145, col: 11, offset: 3083},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 145, col: 11, offset: 3083},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 145, col: 22, offset: 3094},
	val: "∀",
	ignoreCase: false,
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 146, col: 1, offset: 3101},
	expr: &choiceExpr{
	pos: position{line: 146, col: 10, offset: 3112},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 146, col: 10, offset: 3112},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 146, col: 17, offset: 3119},
	val: "→",
	ignoreCase: false,
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 148, col: 1, offset: 3127},
	expr: &seqExpr{
	pos: position{line: 148, col: 12, offset: 3140},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 148, col: 12, offset: 3140},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 148, col: 17, offset: 3145},
	expr: &charClassMatcher{
	pos: position{line: 148, col: 17, offset: 3145},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 148, col: 23, offset: 3151},
	expr: &charClassMatcher{
	pos: position{line: 148, col: 23, offset: 3151},
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
	name: "DoubleLiteral",
	pos: position{line: 150, col: 1, offset: 3159},
	expr: &actionExpr{
	pos: position{line: 150, col: 17, offset: 3177},
	run: (*parser).callonDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 150, col: 17, offset: 3177},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 150, col: 17, offset: 3177},
	expr: &charClassMatcher{
	pos: position{line: 150, col: 17, offset: 3177},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 150, col: 23, offset: 3183},
	expr: &charClassMatcher{
	pos: position{line: 150, col: 23, offset: 3183},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&choiceExpr{
	pos: position{line: 150, col: 32, offset: 3192},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 150, col: 32, offset: 3192},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 150, col: 32, offset: 3192},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 150, col: 36, offset: 3196},
	expr: &charClassMatcher{
	pos: position{line: 150, col: 36, offset: 3196},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&zeroOrOneExpr{
	pos: position{line: 150, col: 43, offset: 3203},
	expr: &ruleRefExpr{
	pos: position{line: 150, col: 43, offset: 3203},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 150, col: 55, offset: 3215},
	name: "Exponent",
},
	},
},
	},
},
},
},
{
	name: "NaturalLiteral",
	pos: position{line: 158, col: 1, offset: 3375},
	expr: &actionExpr{
	pos: position{line: 158, col: 18, offset: 3394},
	run: (*parser).callonNaturalLiteral1,
	expr: &oneOrMoreExpr{
	pos: position{line: 158, col: 18, offset: 3394},
	expr: &charClassMatcher{
	pos: position{line: 158, col: 18, offset: 3394},
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
	pos: position{line: 166, col: 1, offset: 3542},
	expr: &actionExpr{
	pos: position{line: 166, col: 18, offset: 3561},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 166, col: 18, offset: 3561},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 166, col: 18, offset: 3561},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&oneOrMoreExpr{
	pos: position{line: 166, col: 22, offset: 3565},
	expr: &charClassMatcher{
	pos: position{line: 166, col: 22, offset: 3565},
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
	name: "SpaceDeBruijn",
	pos: position{line: 174, col: 1, offset: 3713},
	expr: &actionExpr{
	pos: position{line: 174, col: 17, offset: 3731},
	run: (*parser).callonSpaceDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 174, col: 17, offset: 3731},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 174, col: 17, offset: 3731},
	name: "_",
},
&litMatcher{
	pos: position{line: 174, col: 19, offset: 3733},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 174, col: 23, offset: 3737},
	name: "_",
},
&labeledExpr{
	pos: position{line: 174, col: 25, offset: 3739},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 174, col: 31, offset: 3745},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Identifier",
	pos: position{line: 176, col: 1, offset: 3805},
	expr: &actionExpr{
	pos: position{line: 176, col: 14, offset: 3820},
	run: (*parser).callonIdentifier1,
	expr: &seqExpr{
	pos: position{line: 176, col: 14, offset: 3820},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 176, col: 14, offset: 3820},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 176, col: 19, offset: 3825},
	name: "Label",
},
},
&labeledExpr{
	pos: position{line: 176, col: 25, offset: 3831},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 176, col: 31, offset: 3837},
	expr: &ruleRefExpr{
	pos: position{line: 176, col: 31, offset: 3837},
	name: "SpaceDeBruijn",
},
},
},
	},
},
},
},
{
	name: "IdentifierReservedPrefix",
	pos: position{line: 184, col: 1, offset: 4013},
	expr: &actionExpr{
	pos: position{line: 185, col: 10, offset: 4051},
	run: (*parser).callonIdentifierReservedPrefix1,
	expr: &seqExpr{
	pos: position{line: 185, col: 10, offset: 4051},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 185, col: 10, offset: 4051},
	label: "name",
	expr: &actionExpr{
	pos: position{line: 185, col: 16, offset: 4057},
	run: (*parser).callonIdentifierReservedPrefix4,
	expr: &seqExpr{
	pos: position{line: 185, col: 16, offset: 4057},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 185, col: 16, offset: 4057},
	name: "Reserved",
},
&oneOrMoreExpr{
	pos: position{line: 185, col: 25, offset: 4066},
	expr: &charClassMatcher{
	pos: position{line: 185, col: 25, offset: 4066},
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
&labeledExpr{
	pos: position{line: 186, col: 10, offset: 4122},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 186, col: 16, offset: 4128},
	expr: &ruleRefExpr{
	pos: position{line: 186, col: 16, offset: 4128},
	name: "SpaceDeBruijn",
},
},
},
	},
},
},
},
{
	name: "LetBinding",
	pos: position{line: 203, col: 1, offset: 4816},
	expr: &actionExpr{
	pos: position{line: 203, col: 14, offset: 4831},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 203, col: 14, offset: 4831},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 203, col: 14, offset: 4831},
	val: "let",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 203, col: 20, offset: 4837},
	name: "_",
},
&labeledExpr{
	pos: position{line: 203, col: 22, offset: 4839},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 203, col: 28, offset: 4845},
	name: "Label",
},
},
&labeledExpr{
	pos: position{line: 203, col: 34, offset: 4851},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 203, col: 36, offset: 4853},
	expr: &ruleRefExpr{
	pos: position{line: 203, col: 36, offset: 4853},
	name: "Annotation",
},
},
},
&ruleRefExpr{
	pos: position{line: 203, col: 48, offset: 4865},
	name: "_",
},
&litMatcher{
	pos: position{line: 203, col: 50, offset: 4867},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 203, col: 54, offset: 4871},
	name: "_",
},
&labeledExpr{
	pos: position{line: 203, col: 56, offset: 4873},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 203, col: 58, offset: 4875},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 203, col: 69, offset: 4886},
	name: "_",
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 218, col: 1, offset: 5196},
	expr: &choiceExpr{
	pos: position{line: 219, col: 7, offset: 5217},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 219, col: 7, offset: 5217},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 219, col: 7, offset: 5217},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 219, col: 7, offset: 5217},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 219, col: 14, offset: 5224},
	name: "_",
},
&litMatcher{
	pos: position{line: 219, col: 16, offset: 5226},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 219, col: 20, offset: 5230},
	name: "_",
},
&labeledExpr{
	pos: position{line: 219, col: 22, offset: 5232},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 219, col: 28, offset: 5238},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 219, col: 34, offset: 5244},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 219, col: 36, offset: 5246},
	name: "ColonSpace",
},
&labeledExpr{
	pos: position{line: 219, col: 47, offset: 5257},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 219, col: 49, offset: 5259},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 219, col: 60, offset: 5270},
	name: "_",
},
&litMatcher{
	pos: position{line: 219, col: 62, offset: 5272},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 219, col: 66, offset: 5276},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 219, col: 68, offset: 5278},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 219, col: 74, offset: 5284},
	name: "_",
},
&labeledExpr{
	pos: position{line: 219, col: 76, offset: 5286},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 219, col: 81, offset: 5291},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 222, col: 7, offset: 5420},
	run: (*parser).callonExpression21,
	expr: &seqExpr{
	pos: position{line: 222, col: 7, offset: 5420},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 222, col: 7, offset: 5420},
	val: "if",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 222, col: 12, offset: 5425},
	name: "_",
},
&labeledExpr{
	pos: position{line: 222, col: 14, offset: 5427},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 222, col: 19, offset: 5432},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 222, col: 30, offset: 5443},
	name: "_",
},
&litMatcher{
	pos: position{line: 222, col: 32, offset: 5445},
	val: "then",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 222, col: 39, offset: 5452},
	name: "_",
},
&labeledExpr{
	pos: position{line: 222, col: 41, offset: 5454},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 222, col: 43, offset: 5456},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 222, col: 54, offset: 5467},
	name: "_",
},
&litMatcher{
	pos: position{line: 222, col: 56, offset: 5469},
	val: "else",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 222, col: 63, offset: 5476},
	name: "_",
},
&labeledExpr{
	pos: position{line: 222, col: 65, offset: 5478},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 222, col: 67, offset: 5480},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 225, col: 7, offset: 5582},
	run: (*parser).callonExpression37,
	expr: &seqExpr{
	pos: position{line: 225, col: 7, offset: 5582},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 225, col: 7, offset: 5582},
	label: "bindings",
	expr: &zeroOrMoreExpr{
	pos: position{line: 225, col: 16, offset: 5591},
	expr: &ruleRefExpr{
	pos: position{line: 225, col: 16, offset: 5591},
	name: "LetBinding",
},
},
},
&litMatcher{
	pos: position{line: 225, col: 28, offset: 5603},
	val: "in",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 225, col: 33, offset: 5608},
	name: "_",
},
&labeledExpr{
	pos: position{line: 225, col: 35, offset: 5610},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 225, col: 37, offset: 5612},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 232, col: 7, offset: 5868},
	run: (*parser).callonExpression46,
	expr: &seqExpr{
	pos: position{line: 232, col: 7, offset: 5868},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 232, col: 7, offset: 5868},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 232, col: 14, offset: 5875},
	name: "_",
},
&litMatcher{
	pos: position{line: 232, col: 16, offset: 5877},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 232, col: 20, offset: 5881},
	name: "_",
},
&labeledExpr{
	pos: position{line: 232, col: 22, offset: 5883},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 232, col: 28, offset: 5889},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 232, col: 34, offset: 5895},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 232, col: 36, offset: 5897},
	name: "ColonSpace",
},
&labeledExpr{
	pos: position{line: 232, col: 47, offset: 5908},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 232, col: 49, offset: 5910},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 232, col: 60, offset: 5921},
	name: "_",
},
&litMatcher{
	pos: position{line: 232, col: 62, offset: 5923},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 232, col: 66, offset: 5927},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 232, col: 68, offset: 5929},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 232, col: 74, offset: 5935},
	name: "_",
},
&labeledExpr{
	pos: position{line: 232, col: 76, offset: 5937},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 232, col: 81, offset: 5942},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 235, col: 7, offset: 6063},
	run: (*parser).callonExpression65,
	expr: &seqExpr{
	pos: position{line: 235, col: 7, offset: 6063},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 235, col: 7, offset: 6063},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 235, col: 9, offset: 6065},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 235, col: 28, offset: 6084},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 235, col: 30, offset: 6086},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 235, col: 36, offset: 6092},
	name: "_",
},
&labeledExpr{
	pos: position{line: 235, col: 38, offset: 6094},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 235, col: 40, offset: 6096},
	name: "Expression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 236, col: 7, offset: 6168},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 238, col: 1, offset: 6189},
	expr: &actionExpr{
	pos: position{line: 238, col: 14, offset: 6204},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 238, col: 14, offset: 6204},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 238, col: 14, offset: 6204},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 238, col: 16, offset: 6206},
	name: "ColonSpace",
},
&labeledExpr{
	pos: position{line: 238, col: 27, offset: 6217},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 238, col: 29, offset: 6219},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 240, col: 1, offset: 6249},
	expr: &choiceExpr{
	pos: position{line: 241, col: 5, offset: 6277},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 241, col: 5, offset: 6277},
	name: "EmptyList",
},
&actionExpr{
	pos: position{line: 242, col: 5, offset: 6291},
	run: (*parser).callonAnnotatedExpression3,
	expr: &seqExpr{
	pos: position{line: 242, col: 5, offset: 6291},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 242, col: 5, offset: 6291},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 242, col: 7, offset: 6293},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 242, col: 26, offset: 6312},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 242, col: 28, offset: 6314},
	expr: &ruleRefExpr{
	pos: position{line: 242, col: 28, offset: 6314},
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
	pos: position{line: 247, col: 1, offset: 6431},
	expr: &actionExpr{
	pos: position{line: 247, col: 13, offset: 6445},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 247, col: 13, offset: 6445},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 247, col: 13, offset: 6445},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 247, col: 17, offset: 6449},
	name: "_",
},
&litMatcher{
	pos: position{line: 247, col: 19, offset: 6451},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 247, col: 23, offset: 6455},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 247, col: 25, offset: 6457},
	name: "ColonSpace",
},
&litMatcher{
	pos: position{line: 247, col: 36, offset: 6468},
	val: "List",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 247, col: 43, offset: 6475},
	name: "_",
},
&labeledExpr{
	pos: position{line: 247, col: 45, offset: 6477},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 247, col: 47, offset: 6479},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 251, col: 1, offset: 6550},
	expr: &ruleRefExpr{
	pos: position{line: 251, col: 22, offset: 6573},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 253, col: 1, offset: 6594},
	expr: &ruleRefExpr{
	pos: position{line: 253, col: 23, offset: 6618},
	name: "PlusExpression",
},
},
{
	name: "MorePlus",
	pos: position{line: 255, col: 1, offset: 6634},
	expr: &actionExpr{
	pos: position{line: 255, col: 12, offset: 6647},
	run: (*parser).callonMorePlus1,
	expr: &seqExpr{
	pos: position{line: 255, col: 12, offset: 6647},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 255, col: 12, offset: 6647},
	name: "_",
},
&litMatcher{
	pos: position{line: 255, col: 14, offset: 6649},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 255, col: 18, offset: 6653},
	name: "NonemptyWhitespace",
},
&labeledExpr{
	pos: position{line: 255, col: 37, offset: 6672},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 255, col: 39, offset: 6674},
	name: "TimesExpression",
},
},
	},
},
},
},
{
	name: "PlusExpression",
	pos: position{line: 256, col: 1, offset: 6708},
	expr: &actionExpr{
	pos: position{line: 257, col: 7, offset: 6733},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 257, col: 7, offset: 6733},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 257, col: 7, offset: 6733},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 257, col: 13, offset: 6739},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 257, col: 29, offset: 6755},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 257, col: 34, offset: 6760},
	expr: &ruleRefExpr{
	pos: position{line: 257, col: 34, offset: 6760},
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
	pos: position{line: 266, col: 1, offset: 7000},
	expr: &actionExpr{
	pos: position{line: 266, col: 13, offset: 7014},
	run: (*parser).callonMoreTimes1,
	expr: &seqExpr{
	pos: position{line: 266, col: 13, offset: 7014},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 266, col: 13, offset: 7014},
	name: "_",
},
&litMatcher{
	pos: position{line: 266, col: 15, offset: 7016},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 266, col: 19, offset: 7020},
	name: "_",
},
&labeledExpr{
	pos: position{line: 266, col: 21, offset: 7022},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 266, col: 23, offset: 7024},
	name: "ApplicationExpression",
},
},
	},
},
},
},
{
	name: "TimesExpression",
	pos: position{line: 267, col: 1, offset: 7064},
	expr: &actionExpr{
	pos: position{line: 268, col: 7, offset: 7090},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 268, col: 7, offset: 7090},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 268, col: 7, offset: 7090},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 268, col: 13, offset: 7096},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 268, col: 35, offset: 7118},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 268, col: 40, offset: 7123},
	expr: &ruleRefExpr{
	pos: position{line: 268, col: 40, offset: 7123},
	name: "MoreTimes",
},
},
},
	},
},
},
},
{
	name: "AnArg",
	pos: position{line: 277, col: 1, offset: 7365},
	expr: &actionExpr{
	pos: position{line: 277, col: 9, offset: 7373},
	run: (*parser).callonAnArg1,
	expr: &seqExpr{
	pos: position{line: 277, col: 9, offset: 7373},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 277, col: 9, offset: 7373},
	name: "NonemptyWhitespace",
},
&labeledExpr{
	pos: position{line: 277, col: 28, offset: 7392},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 277, col: 30, offset: 7394},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "ApplicationExpression",
	pos: position{line: 279, col: 1, offset: 7430},
	expr: &actionExpr{
	pos: position{line: 279, col: 25, offset: 7456},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 279, col: 25, offset: 7456},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 279, col: 25, offset: 7456},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 279, col: 27, offset: 7458},
	name: "ImportExpression",
},
},
&labeledExpr{
	pos: position{line: 279, col: 44, offset: 7475},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 279, col: 49, offset: 7480},
	expr: &ruleRefExpr{
	pos: position{line: 279, col: 49, offset: 7480},
	name: "AnArg",
},
},
},
	},
},
},
},
{
	name: "ImportExpression",
	pos: position{line: 288, col: 1, offset: 7711},
	expr: &ruleRefExpr{
	pos: position{line: 288, col: 20, offset: 7732},
	name: "SelectorExpression",
},
},
{
	name: "SelectorExpression",
	pos: position{line: 290, col: 1, offset: 7752},
	expr: &actionExpr{
	pos: position{line: 290, col: 22, offset: 7775},
	run: (*parser).callonSelectorExpression1,
	expr: &seqExpr{
	pos: position{line: 290, col: 22, offset: 7775},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 290, col: 22, offset: 7775},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 290, col: 24, offset: 7777},
	name: "PrimitiveExpression",
},
},
&labeledExpr{
	pos: position{line: 290, col: 44, offset: 7797},
	label: "ls",
	expr: &zeroOrMoreExpr{
	pos: position{line: 290, col: 47, offset: 7800},
	expr: &seqExpr{
	pos: position{line: 290, col: 48, offset: 7801},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 290, col: 48, offset: 7801},
	name: "_",
},
&litMatcher{
	pos: position{line: 290, col: 50, offset: 7803},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 290, col: 54, offset: 7807},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 290, col: 56, offset: 7809},
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
	pos: position{line: 300, col: 1, offset: 8047},
	expr: &choiceExpr{
	pos: position{line: 301, col: 7, offset: 8077},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 301, col: 7, offset: 8077},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 302, col: 7, offset: 8097},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 303, col: 7, offset: 8118},
	name: "IntegerLiteral",
},
&actionExpr{
	pos: position{line: 304, col: 7, offset: 8139},
	run: (*parser).callonPrimitiveExpression5,
	expr: &litMatcher{
	pos: position{line: 304, col: 7, offset: 8139},
	val: "-Infinity",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 305, col: 7, offset: 8201},
	name: "TextLiteral",
},
&actionExpr{
	pos: position{line: 306, col: 7, offset: 8219},
	run: (*parser).callonPrimitiveExpression8,
	expr: &seqExpr{
	pos: position{line: 306, col: 7, offset: 8219},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 306, col: 7, offset: 8219},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 306, col: 11, offset: 8223},
	name: "_",
},
&labeledExpr{
	pos: position{line: 306, col: 13, offset: 8225},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 306, col: 15, offset: 8227},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 306, col: 35, offset: 8247},
	name: "_",
},
&litMatcher{
	pos: position{line: 306, col: 37, offset: 8249},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 307, col: 7, offset: 8277},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 308, col: 7, offset: 8303},
	name: "IdentifierReservedPrefix",
},
&ruleRefExpr{
	pos: position{line: 309, col: 7, offset: 8334},
	name: "Reserved",
},
&ruleRefExpr{
	pos: position{line: 310, col: 7, offset: 8349},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 311, col: 7, offset: 8366},
	run: (*parser).callonPrimitiveExpression20,
	expr: &seqExpr{
	pos: position{line: 311, col: 7, offset: 8366},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 311, col: 7, offset: 8366},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 311, col: 11, offset: 8370},
	name: "_",
},
&labeledExpr{
	pos: position{line: 311, col: 13, offset: 8372},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 311, col: 15, offset: 8374},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 311, col: 26, offset: 8385},
	name: "_",
},
&litMatcher{
	pos: position{line: 311, col: 28, offset: 8387},
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
	pos: position{line: 313, col: 1, offset: 8410},
	expr: &choiceExpr{
	pos: position{line: 314, col: 7, offset: 8440},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 314, col: 7, offset: 8440},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &litMatcher{
	pos: position{line: 314, col: 7, offset: 8440},
	val: "=",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 315, col: 7, offset: 8503},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 316, col: 7, offset: 8528},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 317, col: 7, offset: 8556},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 317, col: 7, offset: 8556},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 319, col: 1, offset: 8610},
	expr: &actionExpr{
	pos: position{line: 319, col: 19, offset: 8630},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 319, col: 19, offset: 8630},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 319, col: 19, offset: 8630},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 319, col: 24, offset: 8635},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 319, col: 30, offset: 8641},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 319, col: 32, offset: 8643},
	name: "ColonSpace",
},
&labeledExpr{
	pos: position{line: 319, col: 43, offset: 8654},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 319, col: 48, offset: 8659},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 322, col: 1, offset: 8716},
	expr: &actionExpr{
	pos: position{line: 322, col: 18, offset: 8735},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 322, col: 18, offset: 8735},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 322, col: 18, offset: 8735},
	name: "_",
},
&litMatcher{
	pos: position{line: 322, col: 20, offset: 8737},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 322, col: 24, offset: 8741},
	name: "_",
},
&labeledExpr{
	pos: position{line: 322, col: 26, offset: 8743},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 322, col: 28, offset: 8745},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 323, col: 1, offset: 8777},
	expr: &actionExpr{
	pos: position{line: 324, col: 7, offset: 8806},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 324, col: 7, offset: 8806},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 324, col: 7, offset: 8806},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 324, col: 13, offset: 8812},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 324, col: 29, offset: 8828},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 324, col: 34, offset: 8833},
	expr: &ruleRefExpr{
	pos: position{line: 324, col: 34, offset: 8833},
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
	pos: position{line: 334, col: 1, offset: 9245},
	expr: &actionExpr{
	pos: position{line: 334, col: 22, offset: 9268},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 334, col: 22, offset: 9268},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 334, col: 22, offset: 9268},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 334, col: 27, offset: 9273},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 334, col: 33, offset: 9279},
	name: "_",
},
&litMatcher{
	pos: position{line: 334, col: 35, offset: 9281},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 334, col: 39, offset: 9285},
	name: "_",
},
&labeledExpr{
	pos: position{line: 334, col: 41, offset: 9287},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 334, col: 46, offset: 9292},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 337, col: 1, offset: 9349},
	expr: &actionExpr{
	pos: position{line: 337, col: 21, offset: 9371},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 337, col: 21, offset: 9371},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 337, col: 21, offset: 9371},
	name: "_",
},
&litMatcher{
	pos: position{line: 337, col: 23, offset: 9373},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 337, col: 27, offset: 9377},
	name: "_",
},
&labeledExpr{
	pos: position{line: 337, col: 29, offset: 9379},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 337, col: 31, offset: 9381},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 338, col: 1, offset: 9416},
	expr: &actionExpr{
	pos: position{line: 339, col: 7, offset: 9448},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 339, col: 7, offset: 9448},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 339, col: 7, offset: 9448},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 339, col: 13, offset: 9454},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 339, col: 32, offset: 9473},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 339, col: 37, offset: 9478},
	expr: &ruleRefExpr{
	pos: position{line: 339, col: 37, offset: 9478},
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
	pos: position{line: 349, col: 1, offset: 9896},
	expr: &actionExpr{
	pos: position{line: 349, col: 12, offset: 9909},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 349, col: 12, offset: 9909},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 349, col: 12, offset: 9909},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 349, col: 16, offset: 9913},
	name: "_",
},
&labeledExpr{
	pos: position{line: 349, col: 18, offset: 9915},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 349, col: 20, offset: 9917},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 349, col: 31, offset: 9928},
	name: "_",
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 351, col: 1, offset: 9947},
	expr: &actionExpr{
	pos: position{line: 352, col: 7, offset: 9977},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 352, col: 7, offset: 9977},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 352, col: 7, offset: 9977},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 352, col: 11, offset: 9981},
	name: "_",
},
&labeledExpr{
	pos: position{line: 352, col: 13, offset: 9983},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 352, col: 19, offset: 9989},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 352, col: 30, offset: 10000},
	name: "_",
},
&labeledExpr{
	pos: position{line: 352, col: 32, offset: 10002},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 352, col: 37, offset: 10007},
	expr: &ruleRefExpr{
	pos: position{line: 352, col: 37, offset: 10007},
	name: "MoreList",
},
},
},
&litMatcher{
	pos: position{line: 352, col: 47, offset: 10017},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 362, col: 1, offset: 10309},
	expr: &notExpr{
	pos: position{line: 362, col: 7, offset: 10317},
	expr: &anyMatcher{
	line: 362, col: 8, offset: 10318,
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

func (c *current) onReserved2() (interface{}, error) {
 return ast.Bool, nil 
}

func (p *parser) callonReserved2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved2()
}

func (c *current) onReserved6() (interface{}, error) {
 return ast.Natural, nil 
}

func (p *parser) callonReserved6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved6()
}

func (c *current) onReserved8() (interface{}, error) {
 return ast.Integer, nil 
}

func (p *parser) callonReserved8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved8()
}

func (c *current) onReserved10() (interface{}, error) {
 return ast.Double, nil 
}

func (p *parser) callonReserved10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved10()
}

func (c *current) onReserved12() (interface{}, error) {
 return ast.Text, nil 
}

func (p *parser) callonReserved12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved12()
}

func (c *current) onReserved14() (interface{}, error) {
 return ast.List, nil 
}

func (p *parser) callonReserved14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved14()
}

func (c *current) onReserved16() (interface{}, error) {
 return ast.True, nil 
}

func (p *parser) callonReserved16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved16()
}

func (c *current) onReserved18() (interface{}, error) {
 return ast.False, nil 
}

func (p *parser) callonReserved18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved18()
}

func (c *current) onReserved20() (interface{}, error) {
 return ast.DoubleLit(math.NaN()), nil 
}

func (p *parser) callonReserved20() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved20()
}

func (c *current) onReserved22() (interface{}, error) {
 return ast.DoubleLit(math.Inf(1)), nil 
}

func (p *parser) callonReserved22() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved22()
}

func (c *current) onReserved24() (interface{}, error) {
 return ast.Type, nil 
}

func (p *parser) callonReserved24() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved24()
}

func (c *current) onReserved26() (interface{}, error) {
 return ast.Kind, nil 
}

func (p *parser) callonReserved26() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved26()
}

func (c *current) onReserved28() (interface{}, error) {
 return ast.Sort, nil 
}

func (p *parser) callonReserved28() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved28()
}

func (c *current) onDoubleLiteral1() (interface{}, error) {
      d, err := strconv.ParseFloat(string(c.text), 64)
      if err != nil {
         return nil, err
      }
      return ast.DoubleLit(d), nil
}

func (p *parser) callonDoubleLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleLiteral1()
}

func (c *current) onNaturalLiteral1() (interface{}, error) {
      i, err := strconv.Atoi(string(c.text))
      if err != nil {
         return nil, err
      }
      return ast.NaturalLit(i), nil
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
      return ast.IntegerLit(i), nil
}

func (p *parser) callonIntegerLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIntegerLiteral1()
}

func (c *current) onSpaceDeBruijn1(index interface{}) (interface{}, error) {
 return int(index.(ast.NaturalLit)), nil 
}

func (p *parser) callonSpaceDeBruijn1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSpaceDeBruijn1(stack["index"])
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

func (c *current) onExpression21(cond, t, f interface{}) (interface{}, error) {
          return ast.BoolIf{cond.(ast.Expr),t.(ast.Expr),f.(ast.Expr)},nil
      
}

func (p *parser) callonExpression21() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression21(stack["cond"], stack["t"], stack["f"])
}

func (c *current) onExpression37(bindings, b interface{}) (interface{}, error) {
        bs := make([]ast.Binding, len(bindings.([]interface{})))
        for i, binding := range bindings.([]interface{}) {
            bs[i] = binding.(ast.Binding)
        }
        return ast.MakeLet(b.(ast.Expr), bs...), nil
      
}

func (p *parser) callonExpression37() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression37(stack["bindings"], stack["b"])
}

func (c *current) onExpression46(label, t, body interface{}) (interface{}, error) {
          return &ast.Pi{Label:label.(string), Type:t.(ast.Expr), Body: body.(ast.Expr)}, nil
      
}

func (p *parser) callonExpression46() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression46(stack["label"], stack["t"], stack["body"])
}

func (c *current) onExpression65(o, e interface{}) (interface{}, error) {
 return &ast.Pi{"_",o.(ast.Expr),e.(ast.Expr)}, nil 
}

func (p *parser) callonExpression65() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression65(stack["o"], stack["e"])
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

func (c *current) onAnArg1(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonAnArg1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAnArg1(stack["e"])
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
        label := labelSelector.([]interface{})[3]
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

func (c *current) onPrimitiveExpression20(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonPrimitiveExpression20() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimitiveExpression20(stack["e"])
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

