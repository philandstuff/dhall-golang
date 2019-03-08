
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
	pos: position{line: 59, col: 5, offset: 1019},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 59, col: 5, offset: 1019},
	val: "\"",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 60, col: 10, offset: 1032},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 61, col: 10, offset: 1045},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 62, col: 10, offset: 1059},
	val: "/",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 63, col: 10, offset: 1072},
	val: "b",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 64, col: 10, offset: 1085},
	val: "f",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 65, col: 10, offset: 1098},
	val: "n",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 66, col: 10, offset: 1111},
	val: "r",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 67, col: 10, offset: 1124},
	val: "t",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 68, col: 10, offset: 1137},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 68, col: 10, offset: 1137},
	val: "u",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 68, col: 14, offset: 1141},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 68, col: 21, offset: 1148},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 68, col: 28, offset: 1155},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 68, col: 35, offset: 1162},
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
	pos: position{line: 89, col: 1, offset: 1605},
	expr: &choiceExpr{
	pos: position{line: 90, col: 6, offset: 1631},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 90, col: 6, offset: 1631},
	run: (*parser).callonDoubleQuoteChunk2,
	expr: &seqExpr{
	pos: position{line: 90, col: 6, offset: 1631},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 90, col: 6, offset: 1631},
	val: "${",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 90, col: 11, offset: 1636},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 90, col: 13, offset: 1638},
	name: "CompleteExpression",
},
},
&litMatcher{
	pos: position{line: 90, col: 32, offset: 1657},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 91, col: 6, offset: 1684},
	name: "EscapedChar",
},
&charClassMatcher{
	pos: position{line: 92, col: 6, offset: 1701},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 93, col: 6, offset: 1718},
	val: "[\\x23-\\x5b]",
	ranges: []rune{'#','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 94, col: 6, offset: 1735},
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
	pos: position{line: 96, col: 1, offset: 1754},
	expr: &actionExpr{
	pos: position{line: 96, col: 22, offset: 1777},
	run: (*parser).callonDoubleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 96, col: 22, offset: 1777},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 96, col: 22, offset: 1777},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 96, col: 26, offset: 1781},
	label: "chunks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 96, col: 33, offset: 1788},
	expr: &ruleRefExpr{
	pos: position{line: 96, col: 33, offset: 1788},
	name: "DoubleQuoteChunk",
},
},
},
&litMatcher{
	pos: position{line: 96, col: 51, offset: 1806},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TextLiteral",
	pos: position{line: 113, col: 1, offset: 2290},
	expr: &ruleRefExpr{
	pos: position{line: 113, col: 15, offset: 2306},
	name: "DoubleQuoteLiteral",
},
},
{
	name: "Reserved",
	pos: position{line: 115, col: 1, offset: 2326},
	expr: &choiceExpr{
	pos: position{line: 115, col: 12, offset: 2339},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 115, col: 12, offset: 2339},
	run: (*parser).callonReserved2,
	expr: &litMatcher{
	pos: position{line: 115, col: 12, offset: 2339},
	val: "Bool",
	ignoreCase: false,
},
},
&litMatcher{
	pos: position{line: 116, col: 5, offset: 2375},
	val: "Optional",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 117, col: 5, offset: 2390},
	val: "None",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 118, col: 5, offset: 2401},
	run: (*parser).callonReserved6,
	expr: &litMatcher{
	pos: position{line: 118, col: 5, offset: 2401},
	val: "Natural",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 119, col: 5, offset: 2443},
	run: (*parser).callonReserved8,
	expr: &litMatcher{
	pos: position{line: 119, col: 5, offset: 2443},
	val: "Integer",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 120, col: 5, offset: 2485},
	run: (*parser).callonReserved10,
	expr: &litMatcher{
	pos: position{line: 120, col: 5, offset: 2485},
	val: "Double",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 121, col: 5, offset: 2525},
	run: (*parser).callonReserved12,
	expr: &litMatcher{
	pos: position{line: 121, col: 5, offset: 2525},
	val: "Text",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 122, col: 5, offset: 2561},
	run: (*parser).callonReserved14,
	expr: &litMatcher{
	pos: position{line: 122, col: 5, offset: 2561},
	val: "List",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 123, col: 5, offset: 2597},
	run: (*parser).callonReserved16,
	expr: &litMatcher{
	pos: position{line: 123, col: 5, offset: 2597},
	val: "True",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 124, col: 5, offset: 2633},
	run: (*parser).callonReserved18,
	expr: &litMatcher{
	pos: position{line: 124, col: 5, offset: 2633},
	val: "False",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 125, col: 5, offset: 2671},
	run: (*parser).callonReserved20,
	expr: &litMatcher{
	pos: position{line: 125, col: 5, offset: 2671},
	val: "NaN",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 126, col: 5, offset: 2723},
	run: (*parser).callonReserved22,
	expr: &litMatcher{
	pos: position{line: 126, col: 5, offset: 2723},
	val: "Infinity",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 127, col: 5, offset: 2781},
	run: (*parser).callonReserved24,
	expr: &litMatcher{
	pos: position{line: 127, col: 5, offset: 2781},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 128, col: 5, offset: 2817},
	run: (*parser).callonReserved26,
	expr: &litMatcher{
	pos: position{line: 128, col: 5, offset: 2817},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 129, col: 5, offset: 2853},
	run: (*parser).callonReserved28,
	expr: &litMatcher{
	pos: position{line: 129, col: 5, offset: 2853},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "Keyword",
	pos: position{line: 131, col: 1, offset: 2886},
	expr: &choiceExpr{
	pos: position{line: 131, col: 11, offset: 2898},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 131, col: 11, offset: 2898},
	val: "if",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 132, col: 5, offset: 2907},
	val: "then",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 133, col: 5, offset: 2918},
	val: "else",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 134, col: 5, offset: 2929},
	val: "let",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 135, col: 5, offset: 2939},
	val: "in",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 136, col: 5, offset: 2948},
	val: "as",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 137, col: 5, offset: 2957},
	val: "using",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 138, col: 5, offset: 2969},
	val: "merge",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 139, col: 5, offset: 2981},
	val: "constructors",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 140, col: 5, offset: 3000},
	val: "Some",
	ignoreCase: false,
},
	},
},
},
{
	name: "ColonSpace",
	pos: position{line: 142, col: 1, offset: 3008},
	expr: &seqExpr{
	pos: position{line: 142, col: 14, offset: 3023},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 142, col: 14, offset: 3023},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 142, col: 18, offset: 3027},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Lambda",
	pos: position{line: 144, col: 1, offset: 3047},
	expr: &choiceExpr{
	pos: position{line: 144, col: 11, offset: 3059},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 144, col: 11, offset: 3059},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 144, col: 18, offset: 3066},
	val: "λ",
	ignoreCase: false,
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 145, col: 1, offset: 3072},
	expr: &choiceExpr{
	pos: position{line: 145, col: 11, offset: 3084},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 145, col: 11, offset: 3084},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 145, col: 22, offset: 3095},
	val: "∀",
	ignoreCase: false,
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 146, col: 1, offset: 3102},
	expr: &choiceExpr{
	pos: position{line: 146, col: 10, offset: 3113},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 146, col: 10, offset: 3113},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 146, col: 17, offset: 3120},
	val: "→",
	ignoreCase: false,
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 148, col: 1, offset: 3128},
	expr: &seqExpr{
	pos: position{line: 148, col: 12, offset: 3141},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 148, col: 12, offset: 3141},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 148, col: 17, offset: 3146},
	expr: &charClassMatcher{
	pos: position{line: 148, col: 17, offset: 3146},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 148, col: 23, offset: 3152},
	expr: &charClassMatcher{
	pos: position{line: 148, col: 23, offset: 3152},
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
	pos: position{line: 150, col: 1, offset: 3160},
	expr: &actionExpr{
	pos: position{line: 150, col: 17, offset: 3178},
	run: (*parser).callonDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 150, col: 17, offset: 3178},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 150, col: 17, offset: 3178},
	expr: &charClassMatcher{
	pos: position{line: 150, col: 17, offset: 3178},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 150, col: 23, offset: 3184},
	expr: &charClassMatcher{
	pos: position{line: 150, col: 23, offset: 3184},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&choiceExpr{
	pos: position{line: 150, col: 32, offset: 3193},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 150, col: 32, offset: 3193},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 150, col: 32, offset: 3193},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 150, col: 36, offset: 3197},
	expr: &charClassMatcher{
	pos: position{line: 150, col: 36, offset: 3197},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&zeroOrOneExpr{
	pos: position{line: 150, col: 43, offset: 3204},
	expr: &ruleRefExpr{
	pos: position{line: 150, col: 43, offset: 3204},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 150, col: 55, offset: 3216},
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
	pos: position{line: 158, col: 1, offset: 3376},
	expr: &actionExpr{
	pos: position{line: 158, col: 18, offset: 3395},
	run: (*parser).callonNaturalLiteral1,
	expr: &oneOrMoreExpr{
	pos: position{line: 158, col: 18, offset: 3395},
	expr: &charClassMatcher{
	pos: position{line: 158, col: 18, offset: 3395},
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
	pos: position{line: 166, col: 1, offset: 3543},
	expr: &actionExpr{
	pos: position{line: 166, col: 18, offset: 3562},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 166, col: 18, offset: 3562},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 166, col: 18, offset: 3562},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&oneOrMoreExpr{
	pos: position{line: 166, col: 22, offset: 3566},
	expr: &charClassMatcher{
	pos: position{line: 166, col: 22, offset: 3566},
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
	pos: position{line: 174, col: 1, offset: 3714},
	expr: &actionExpr{
	pos: position{line: 174, col: 17, offset: 3732},
	run: (*parser).callonSpaceDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 174, col: 17, offset: 3732},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 174, col: 17, offset: 3732},
	name: "_",
},
&litMatcher{
	pos: position{line: 174, col: 19, offset: 3734},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 174, col: 23, offset: 3738},
	name: "_",
},
&labeledExpr{
	pos: position{line: 174, col: 25, offset: 3740},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 174, col: 31, offset: 3746},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Identifier",
	pos: position{line: 176, col: 1, offset: 3806},
	expr: &actionExpr{
	pos: position{line: 176, col: 14, offset: 3821},
	run: (*parser).callonIdentifier1,
	expr: &seqExpr{
	pos: position{line: 176, col: 14, offset: 3821},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 176, col: 14, offset: 3821},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 176, col: 19, offset: 3826},
	name: "Label",
},
},
&labeledExpr{
	pos: position{line: 176, col: 25, offset: 3832},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 176, col: 31, offset: 3838},
	expr: &ruleRefExpr{
	pos: position{line: 176, col: 31, offset: 3838},
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
	pos: position{line: 184, col: 1, offset: 4014},
	expr: &actionExpr{
	pos: position{line: 185, col: 10, offset: 4052},
	run: (*parser).callonIdentifierReservedPrefix1,
	expr: &seqExpr{
	pos: position{line: 185, col: 10, offset: 4052},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 185, col: 10, offset: 4052},
	label: "name",
	expr: &actionExpr{
	pos: position{line: 185, col: 16, offset: 4058},
	run: (*parser).callonIdentifierReservedPrefix4,
	expr: &seqExpr{
	pos: position{line: 185, col: 16, offset: 4058},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 185, col: 16, offset: 4058},
	name: "Reserved",
},
&oneOrMoreExpr{
	pos: position{line: 185, col: 25, offset: 4067},
	expr: &charClassMatcher{
	pos: position{line: 185, col: 25, offset: 4067},
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
	pos: position{line: 186, col: 10, offset: 4123},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 186, col: 16, offset: 4129},
	expr: &ruleRefExpr{
	pos: position{line: 186, col: 16, offset: 4129},
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
	pos: position{line: 203, col: 1, offset: 4817},
	expr: &actionExpr{
	pos: position{line: 203, col: 14, offset: 4832},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 203, col: 14, offset: 4832},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 203, col: 14, offset: 4832},
	val: "let",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 203, col: 20, offset: 4838},
	name: "_",
},
&labeledExpr{
	pos: position{line: 203, col: 22, offset: 4840},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 203, col: 28, offset: 4846},
	name: "Label",
},
},
&labeledExpr{
	pos: position{line: 203, col: 34, offset: 4852},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 203, col: 36, offset: 4854},
	expr: &ruleRefExpr{
	pos: position{line: 203, col: 36, offset: 4854},
	name: "Annotation",
},
},
},
&ruleRefExpr{
	pos: position{line: 203, col: 48, offset: 4866},
	name: "_",
},
&litMatcher{
	pos: position{line: 203, col: 50, offset: 4868},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 203, col: 54, offset: 4872},
	name: "_",
},
&labeledExpr{
	pos: position{line: 203, col: 56, offset: 4874},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 203, col: 58, offset: 4876},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 203, col: 69, offset: 4887},
	name: "_",
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 218, col: 1, offset: 5197},
	expr: &choiceExpr{
	pos: position{line: 219, col: 7, offset: 5218},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 219, col: 7, offset: 5218},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 219, col: 7, offset: 5218},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 219, col: 7, offset: 5218},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 219, col: 14, offset: 5225},
	name: "_",
},
&litMatcher{
	pos: position{line: 219, col: 16, offset: 5227},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 219, col: 20, offset: 5231},
	name: "_",
},
&labeledExpr{
	pos: position{line: 219, col: 22, offset: 5233},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 219, col: 28, offset: 5239},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 219, col: 34, offset: 5245},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 219, col: 36, offset: 5247},
	name: "ColonSpace",
},
&labeledExpr{
	pos: position{line: 219, col: 47, offset: 5258},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 219, col: 49, offset: 5260},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 219, col: 60, offset: 5271},
	name: "_",
},
&litMatcher{
	pos: position{line: 219, col: 62, offset: 5273},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 219, col: 66, offset: 5277},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 219, col: 68, offset: 5279},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 219, col: 74, offset: 5285},
	name: "_",
},
&labeledExpr{
	pos: position{line: 219, col: 76, offset: 5287},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 219, col: 81, offset: 5292},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 222, col: 7, offset: 5421},
	run: (*parser).callonExpression21,
	expr: &seqExpr{
	pos: position{line: 222, col: 7, offset: 5421},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 222, col: 7, offset: 5421},
	val: "if",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 222, col: 12, offset: 5426},
	name: "_",
},
&labeledExpr{
	pos: position{line: 222, col: 14, offset: 5428},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 222, col: 19, offset: 5433},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 222, col: 30, offset: 5444},
	name: "_",
},
&litMatcher{
	pos: position{line: 222, col: 32, offset: 5446},
	val: "then",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 222, col: 39, offset: 5453},
	name: "_",
},
&labeledExpr{
	pos: position{line: 222, col: 41, offset: 5455},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 222, col: 43, offset: 5457},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 222, col: 54, offset: 5468},
	name: "_",
},
&litMatcher{
	pos: position{line: 222, col: 56, offset: 5470},
	val: "else",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 222, col: 63, offset: 5477},
	name: "_",
},
&labeledExpr{
	pos: position{line: 222, col: 65, offset: 5479},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 222, col: 67, offset: 5481},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 225, col: 7, offset: 5583},
	run: (*parser).callonExpression37,
	expr: &seqExpr{
	pos: position{line: 225, col: 7, offset: 5583},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 225, col: 7, offset: 5583},
	label: "bindings",
	expr: &zeroOrMoreExpr{
	pos: position{line: 225, col: 16, offset: 5592},
	expr: &ruleRefExpr{
	pos: position{line: 225, col: 16, offset: 5592},
	name: "LetBinding",
},
},
},
&litMatcher{
	pos: position{line: 225, col: 28, offset: 5604},
	val: "in",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 225, col: 33, offset: 5609},
	name: "_",
},
&labeledExpr{
	pos: position{line: 225, col: 35, offset: 5611},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 225, col: 37, offset: 5613},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 232, col: 7, offset: 5869},
	run: (*parser).callonExpression46,
	expr: &seqExpr{
	pos: position{line: 232, col: 7, offset: 5869},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 232, col: 7, offset: 5869},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 232, col: 14, offset: 5876},
	name: "_",
},
&litMatcher{
	pos: position{line: 232, col: 16, offset: 5878},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 232, col: 20, offset: 5882},
	name: "_",
},
&labeledExpr{
	pos: position{line: 232, col: 22, offset: 5884},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 232, col: 28, offset: 5890},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 232, col: 34, offset: 5896},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 232, col: 36, offset: 5898},
	name: "ColonSpace",
},
&labeledExpr{
	pos: position{line: 232, col: 47, offset: 5909},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 232, col: 49, offset: 5911},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 232, col: 60, offset: 5922},
	name: "_",
},
&litMatcher{
	pos: position{line: 232, col: 62, offset: 5924},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 232, col: 66, offset: 5928},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 232, col: 68, offset: 5930},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 232, col: 74, offset: 5936},
	name: "_",
},
&labeledExpr{
	pos: position{line: 232, col: 76, offset: 5938},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 232, col: 81, offset: 5943},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 235, col: 7, offset: 6064},
	run: (*parser).callonExpression65,
	expr: &seqExpr{
	pos: position{line: 235, col: 7, offset: 6064},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 235, col: 7, offset: 6064},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 235, col: 9, offset: 6066},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 235, col: 28, offset: 6085},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 235, col: 30, offset: 6087},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 235, col: 36, offset: 6093},
	name: "_",
},
&labeledExpr{
	pos: position{line: 235, col: 38, offset: 6095},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 235, col: 40, offset: 6097},
	name: "Expression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 236, col: 7, offset: 6169},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 238, col: 1, offset: 6190},
	expr: &actionExpr{
	pos: position{line: 238, col: 14, offset: 6205},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 238, col: 14, offset: 6205},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 238, col: 14, offset: 6205},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 238, col: 16, offset: 6207},
	name: "ColonSpace",
},
&labeledExpr{
	pos: position{line: 238, col: 27, offset: 6218},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 238, col: 29, offset: 6220},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 240, col: 1, offset: 6250},
	expr: &choiceExpr{
	pos: position{line: 241, col: 5, offset: 6278},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 241, col: 5, offset: 6278},
	name: "EmptyList",
},
&actionExpr{
	pos: position{line: 242, col: 5, offset: 6292},
	run: (*parser).callonAnnotatedExpression3,
	expr: &seqExpr{
	pos: position{line: 242, col: 5, offset: 6292},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 242, col: 5, offset: 6292},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 242, col: 7, offset: 6294},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 242, col: 26, offset: 6313},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 242, col: 28, offset: 6315},
	expr: &ruleRefExpr{
	pos: position{line: 242, col: 28, offset: 6315},
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
	pos: position{line: 247, col: 1, offset: 6432},
	expr: &actionExpr{
	pos: position{line: 247, col: 13, offset: 6446},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 247, col: 13, offset: 6446},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 247, col: 13, offset: 6446},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 247, col: 17, offset: 6450},
	name: "_",
},
&litMatcher{
	pos: position{line: 247, col: 19, offset: 6452},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 247, col: 23, offset: 6456},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 247, col: 25, offset: 6458},
	name: "ColonSpace",
},
&litMatcher{
	pos: position{line: 247, col: 36, offset: 6469},
	val: "List",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 247, col: 43, offset: 6476},
	name: "_",
},
&labeledExpr{
	pos: position{line: 247, col: 45, offset: 6478},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 247, col: 47, offset: 6480},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 251, col: 1, offset: 6551},
	expr: &ruleRefExpr{
	pos: position{line: 251, col: 22, offset: 6574},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 253, col: 1, offset: 6595},
	expr: &ruleRefExpr{
	pos: position{line: 253, col: 23, offset: 6619},
	name: "PlusExpression",
},
},
{
	name: "MorePlus",
	pos: position{line: 255, col: 1, offset: 6635},
	expr: &actionExpr{
	pos: position{line: 255, col: 12, offset: 6648},
	run: (*parser).callonMorePlus1,
	expr: &seqExpr{
	pos: position{line: 255, col: 12, offset: 6648},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 255, col: 12, offset: 6648},
	name: "_",
},
&litMatcher{
	pos: position{line: 255, col: 14, offset: 6650},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 255, col: 18, offset: 6654},
	name: "NonemptyWhitespace",
},
&labeledExpr{
	pos: position{line: 255, col: 37, offset: 6673},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 255, col: 39, offset: 6675},
	name: "TimesExpression",
},
},
	},
},
},
},
{
	name: "PlusExpression",
	pos: position{line: 256, col: 1, offset: 6709},
	expr: &actionExpr{
	pos: position{line: 257, col: 7, offset: 6734},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 257, col: 7, offset: 6734},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 257, col: 7, offset: 6734},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 257, col: 13, offset: 6740},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 257, col: 29, offset: 6756},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 257, col: 34, offset: 6761},
	expr: &ruleRefExpr{
	pos: position{line: 257, col: 34, offset: 6761},
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
	pos: position{line: 266, col: 1, offset: 7001},
	expr: &actionExpr{
	pos: position{line: 266, col: 13, offset: 7015},
	run: (*parser).callonMoreTimes1,
	expr: &seqExpr{
	pos: position{line: 266, col: 13, offset: 7015},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 266, col: 13, offset: 7015},
	name: "_",
},
&litMatcher{
	pos: position{line: 266, col: 15, offset: 7017},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 266, col: 19, offset: 7021},
	name: "_",
},
&labeledExpr{
	pos: position{line: 266, col: 21, offset: 7023},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 266, col: 23, offset: 7025},
	name: "ApplicationExpression",
},
},
	},
},
},
},
{
	name: "TimesExpression",
	pos: position{line: 267, col: 1, offset: 7065},
	expr: &actionExpr{
	pos: position{line: 268, col: 7, offset: 7091},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 268, col: 7, offset: 7091},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 268, col: 7, offset: 7091},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 268, col: 13, offset: 7097},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 268, col: 35, offset: 7119},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 268, col: 40, offset: 7124},
	expr: &ruleRefExpr{
	pos: position{line: 268, col: 40, offset: 7124},
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
	pos: position{line: 277, col: 1, offset: 7366},
	expr: &actionExpr{
	pos: position{line: 277, col: 9, offset: 7374},
	run: (*parser).callonAnArg1,
	expr: &seqExpr{
	pos: position{line: 277, col: 9, offset: 7374},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 277, col: 9, offset: 7374},
	name: "NonemptyWhitespace",
},
&labeledExpr{
	pos: position{line: 277, col: 28, offset: 7393},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 277, col: 30, offset: 7395},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "ApplicationExpression",
	pos: position{line: 279, col: 1, offset: 7431},
	expr: &actionExpr{
	pos: position{line: 279, col: 25, offset: 7457},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 279, col: 25, offset: 7457},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 279, col: 25, offset: 7457},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 279, col: 27, offset: 7459},
	name: "ImportExpression",
},
},
&labeledExpr{
	pos: position{line: 279, col: 44, offset: 7476},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 279, col: 49, offset: 7481},
	expr: &ruleRefExpr{
	pos: position{line: 279, col: 49, offset: 7481},
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
	pos: position{line: 288, col: 1, offset: 7712},
	expr: &ruleRefExpr{
	pos: position{line: 288, col: 20, offset: 7733},
	name: "SelectorExpression",
},
},
{
	name: "SelectorExpression",
	pos: position{line: 290, col: 1, offset: 7753},
	expr: &ruleRefExpr{
	pos: position{line: 290, col: 22, offset: 7776},
	name: "PrimitiveExpression",
},
},
{
	name: "PrimitiveExpression",
	pos: position{line: 292, col: 1, offset: 7797},
	expr: &choiceExpr{
	pos: position{line: 293, col: 7, offset: 7827},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 293, col: 7, offset: 7827},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 294, col: 7, offset: 7847},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 295, col: 7, offset: 7868},
	name: "IntegerLiteral",
},
&actionExpr{
	pos: position{line: 296, col: 7, offset: 7889},
	run: (*parser).callonPrimitiveExpression5,
	expr: &litMatcher{
	pos: position{line: 296, col: 7, offset: 7889},
	val: "-Infinity",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 297, col: 7, offset: 7951},
	name: "TextLiteral",
},
&actionExpr{
	pos: position{line: 298, col: 7, offset: 7969},
	run: (*parser).callonPrimitiveExpression8,
	expr: &seqExpr{
	pos: position{line: 298, col: 7, offset: 7969},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 298, col: 7, offset: 7969},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 298, col: 11, offset: 7973},
	name: "_",
},
&labeledExpr{
	pos: position{line: 298, col: 13, offset: 7975},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 298, col: 15, offset: 7977},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 298, col: 35, offset: 7997},
	name: "_",
},
&litMatcher{
	pos: position{line: 298, col: 37, offset: 7999},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 299, col: 7, offset: 8027},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 300, col: 7, offset: 8053},
	name: "IdentifierReservedPrefix",
},
&ruleRefExpr{
	pos: position{line: 301, col: 7, offset: 8084},
	name: "Reserved",
},
&ruleRefExpr{
	pos: position{line: 302, col: 7, offset: 8099},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 303, col: 7, offset: 8116},
	run: (*parser).callonPrimitiveExpression20,
	expr: &seqExpr{
	pos: position{line: 303, col: 7, offset: 8116},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 303, col: 7, offset: 8116},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 303, col: 11, offset: 8120},
	name: "_",
},
&labeledExpr{
	pos: position{line: 303, col: 13, offset: 8122},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 303, col: 15, offset: 8124},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 303, col: 26, offset: 8135},
	name: "_",
},
&litMatcher{
	pos: position{line: 303, col: 28, offset: 8137},
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
	pos: position{line: 305, col: 1, offset: 8160},
	expr: &choiceExpr{
	pos: position{line: 306, col: 7, offset: 8190},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 306, col: 7, offset: 8190},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &litMatcher{
	pos: position{line: 306, col: 7, offset: 8190},
	val: "=",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 307, col: 7, offset: 8253},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 308, col: 7, offset: 8278},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 309, col: 7, offset: 8306},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 309, col: 7, offset: 8306},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 311, col: 1, offset: 8360},
	expr: &actionExpr{
	pos: position{line: 311, col: 19, offset: 8380},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 311, col: 19, offset: 8380},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 311, col: 19, offset: 8380},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 311, col: 24, offset: 8385},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 311, col: 30, offset: 8391},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 311, col: 32, offset: 8393},
	name: "ColonSpace",
},
&labeledExpr{
	pos: position{line: 311, col: 43, offset: 8404},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 311, col: 48, offset: 8409},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 314, col: 1, offset: 8466},
	expr: &actionExpr{
	pos: position{line: 314, col: 18, offset: 8485},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 314, col: 18, offset: 8485},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 314, col: 18, offset: 8485},
	name: "_",
},
&litMatcher{
	pos: position{line: 314, col: 20, offset: 8487},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 314, col: 24, offset: 8491},
	name: "_",
},
&labeledExpr{
	pos: position{line: 314, col: 26, offset: 8493},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 314, col: 28, offset: 8495},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 315, col: 1, offset: 8527},
	expr: &actionExpr{
	pos: position{line: 316, col: 7, offset: 8556},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 316, col: 7, offset: 8556},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 316, col: 7, offset: 8556},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 316, col: 13, offset: 8562},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 316, col: 29, offset: 8578},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 316, col: 34, offset: 8583},
	expr: &ruleRefExpr{
	pos: position{line: 316, col: 34, offset: 8583},
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
	pos: position{line: 326, col: 1, offset: 8995},
	expr: &actionExpr{
	pos: position{line: 326, col: 22, offset: 9018},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 326, col: 22, offset: 9018},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 326, col: 22, offset: 9018},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 326, col: 27, offset: 9023},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 326, col: 33, offset: 9029},
	name: "_",
},
&litMatcher{
	pos: position{line: 326, col: 35, offset: 9031},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 326, col: 39, offset: 9035},
	name: "_",
},
&labeledExpr{
	pos: position{line: 326, col: 41, offset: 9037},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 326, col: 46, offset: 9042},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 329, col: 1, offset: 9099},
	expr: &actionExpr{
	pos: position{line: 329, col: 21, offset: 9121},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 329, col: 21, offset: 9121},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 329, col: 21, offset: 9121},
	name: "_",
},
&litMatcher{
	pos: position{line: 329, col: 23, offset: 9123},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 329, col: 27, offset: 9127},
	name: "_",
},
&labeledExpr{
	pos: position{line: 329, col: 29, offset: 9129},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 329, col: 31, offset: 9131},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 330, col: 1, offset: 9166},
	expr: &actionExpr{
	pos: position{line: 331, col: 7, offset: 9198},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 331, col: 7, offset: 9198},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 331, col: 7, offset: 9198},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 331, col: 13, offset: 9204},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 331, col: 32, offset: 9223},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 331, col: 37, offset: 9228},
	expr: &ruleRefExpr{
	pos: position{line: 331, col: 37, offset: 9228},
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
	pos: position{line: 341, col: 1, offset: 9646},
	expr: &actionExpr{
	pos: position{line: 341, col: 12, offset: 9659},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 341, col: 12, offset: 9659},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 341, col: 12, offset: 9659},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 341, col: 16, offset: 9663},
	name: "_",
},
&labeledExpr{
	pos: position{line: 341, col: 18, offset: 9665},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 341, col: 20, offset: 9667},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 341, col: 31, offset: 9678},
	name: "_",
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 343, col: 1, offset: 9697},
	expr: &actionExpr{
	pos: position{line: 344, col: 7, offset: 9727},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 344, col: 7, offset: 9727},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 344, col: 7, offset: 9727},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 344, col: 11, offset: 9731},
	name: "_",
},
&labeledExpr{
	pos: position{line: 344, col: 13, offset: 9733},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 344, col: 19, offset: 9739},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 344, col: 30, offset: 9750},
	name: "_",
},
&labeledExpr{
	pos: position{line: 344, col: 32, offset: 9752},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 344, col: 37, offset: 9757},
	expr: &ruleRefExpr{
	pos: position{line: 344, col: 37, offset: 9757},
	name: "MoreList",
},
},
},
&litMatcher{
	pos: position{line: 344, col: 47, offset: 9767},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 354, col: 1, offset: 10059},
	expr: &notExpr{
	pos: position{line: 354, col: 7, offset: 10067},
	expr: &anyMatcher{
	line: 354, col: 8, offset: 10068,
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

