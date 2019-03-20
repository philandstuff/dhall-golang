
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
import . "github.com/philandstuff/dhall-golang/ast"


var g = &grammar {
	rules: []*rule{
{
	name: "DhallFile",
	pos: position{line: 21, col: 1, offset: 182},
	expr: &actionExpr{
	pos: position{line: 21, col: 13, offset: 196},
	run: (*parser).callonDhallFile1,
	expr: &seqExpr{
	pos: position{line: 21, col: 13, offset: 196},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 21, col: 13, offset: 196},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 21, col: 15, offset: 198},
	name: "CompleteExpression",
},
},
&ruleRefExpr{
	pos: position{line: 21, col: 34, offset: 217},
	name: "EOF",
},
	},
},
},
},
{
	name: "CompleteExpression",
	pos: position{line: 23, col: 1, offset: 240},
	expr: &actionExpr{
	pos: position{line: 23, col: 22, offset: 263},
	run: (*parser).callonCompleteExpression1,
	expr: &seqExpr{
	pos: position{line: 23, col: 22, offset: 263},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 23, col: 22, offset: 263},
	name: "_",
},
&labeledExpr{
	pos: position{line: 23, col: 24, offset: 265},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 23, col: 26, offset: 267},
	name: "Expression",
},
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
	name: "KeywordRaw",
},
},
&charClassMatcher{
	pos: position{line: 49, col: 13, offset: 836},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 49, col: 23, offset: 846},
	expr: &charClassMatcher{
	pos: position{line: 49, col: 23, offset: 846},
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
	pos: position{line: 53, col: 1, offset: 910},
	expr: &actionExpr{
	pos: position{line: 53, col: 9, offset: 920},
	run: (*parser).callonLabel1,
	expr: &seqExpr{
	pos: position{line: 53, col: 9, offset: 920},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 53, col: 9, offset: 920},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 53, col: 15, offset: 926},
	name: "SimpleLabel",
},
},
&ruleRefExpr{
	pos: position{line: 53, col: 27, offset: 938},
	name: "_",
},
	},
},
},
},
{
	name: "EscapedChar",
	pos: position{line: 57, col: 1, offset: 996},
	expr: &actionExpr{
	pos: position{line: 58, col: 3, offset: 1014},
	run: (*parser).callonEscapedChar1,
	expr: &seqExpr{
	pos: position{line: 58, col: 3, offset: 1014},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 58, col: 3, offset: 1014},
	val: "\\",
	ignoreCase: false,
},
&choiceExpr{
	pos: position{line: 59, col: 5, offset: 1023},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 59, col: 5, offset: 1023},
	val: "\"",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 60, col: 10, offset: 1036},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 61, col: 10, offset: 1049},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 62, col: 10, offset: 1063},
	val: "/",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 63, col: 10, offset: 1076},
	val: "b",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 64, col: 10, offset: 1089},
	val: "f",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 65, col: 10, offset: 1102},
	val: "n",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 66, col: 10, offset: 1115},
	val: "r",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 67, col: 10, offset: 1128},
	val: "t",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 68, col: 10, offset: 1141},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 68, col: 10, offset: 1141},
	val: "u",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 68, col: 14, offset: 1145},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 68, col: 21, offset: 1152},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 68, col: 28, offset: 1159},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 68, col: 35, offset: 1166},
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
	pos: position{line: 89, col: 1, offset: 1609},
	expr: &choiceExpr{
	pos: position{line: 90, col: 6, offset: 1635},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 90, col: 6, offset: 1635},
	run: (*parser).callonDoubleQuoteChunk2,
	expr: &seqExpr{
	pos: position{line: 90, col: 6, offset: 1635},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 90, col: 6, offset: 1635},
	val: "${",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 90, col: 11, offset: 1640},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 90, col: 13, offset: 1642},
	name: "CompleteExpression",
},
},
&litMatcher{
	pos: position{line: 90, col: 32, offset: 1661},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 91, col: 6, offset: 1688},
	name: "EscapedChar",
},
&charClassMatcher{
	pos: position{line: 92, col: 6, offset: 1705},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 93, col: 6, offset: 1722},
	val: "[\\x23-\\x5b]",
	ranges: []rune{'#','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 94, col: 6, offset: 1739},
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
	pos: position{line: 96, col: 1, offset: 1758},
	expr: &actionExpr{
	pos: position{line: 96, col: 22, offset: 1781},
	run: (*parser).callonDoubleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 96, col: 22, offset: 1781},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 96, col: 22, offset: 1781},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 96, col: 26, offset: 1785},
	label: "chunks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 96, col: 33, offset: 1792},
	expr: &ruleRefExpr{
	pos: position{line: 96, col: 33, offset: 1792},
	name: "DoubleQuoteChunk",
},
},
},
&litMatcher{
	pos: position{line: 96, col: 51, offset: 1810},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TextLiteral",
	pos: position{line: 113, col: 1, offset: 2278},
	expr: &actionExpr{
	pos: position{line: 113, col: 15, offset: 2294},
	run: (*parser).callonTextLiteral1,
	expr: &seqExpr{
	pos: position{line: 113, col: 15, offset: 2294},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 113, col: 15, offset: 2294},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 113, col: 17, offset: 2296},
	name: "DoubleQuoteLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 113, col: 36, offset: 2315},
	name: "_",
},
	},
},
},
},
{
	name: "ReservedRaw",
	pos: position{line: 115, col: 1, offset: 2336},
	expr: &choiceExpr{
	pos: position{line: 115, col: 15, offset: 2352},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 115, col: 15, offset: 2352},
	run: (*parser).callonReservedRaw2,
	expr: &litMatcher{
	pos: position{line: 115, col: 15, offset: 2352},
	val: "Bool",
	ignoreCase: false,
},
},
&litMatcher{
	pos: position{line: 116, col: 5, offset: 2384},
	val: "Optional",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 117, col: 5, offset: 2399},
	val: "None",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 118, col: 5, offset: 2410},
	run: (*parser).callonReservedRaw6,
	expr: &litMatcher{
	pos: position{line: 118, col: 5, offset: 2410},
	val: "Natural",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 119, col: 5, offset: 2448},
	run: (*parser).callonReservedRaw8,
	expr: &litMatcher{
	pos: position{line: 119, col: 5, offset: 2448},
	val: "Integer",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 120, col: 5, offset: 2486},
	run: (*parser).callonReservedRaw10,
	expr: &litMatcher{
	pos: position{line: 120, col: 5, offset: 2486},
	val: "Double",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 121, col: 5, offset: 2522},
	run: (*parser).callonReservedRaw12,
	expr: &litMatcher{
	pos: position{line: 121, col: 5, offset: 2522},
	val: "Text",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 122, col: 5, offset: 2554},
	run: (*parser).callonReservedRaw14,
	expr: &litMatcher{
	pos: position{line: 122, col: 5, offset: 2554},
	val: "List",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 123, col: 5, offset: 2586},
	run: (*parser).callonReservedRaw16,
	expr: &litMatcher{
	pos: position{line: 123, col: 5, offset: 2586},
	val: "True",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 124, col: 5, offset: 2618},
	run: (*parser).callonReservedRaw18,
	expr: &litMatcher{
	pos: position{line: 124, col: 5, offset: 2618},
	val: "False",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 125, col: 5, offset: 2652},
	run: (*parser).callonReservedRaw20,
	expr: &litMatcher{
	pos: position{line: 125, col: 5, offset: 2652},
	val: "NaN",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 126, col: 5, offset: 2700},
	run: (*parser).callonReservedRaw22,
	expr: &litMatcher{
	pos: position{line: 126, col: 5, offset: 2700},
	val: "Infinity",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 127, col: 5, offset: 2754},
	run: (*parser).callonReservedRaw24,
	expr: &litMatcher{
	pos: position{line: 127, col: 5, offset: 2754},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 128, col: 5, offset: 2786},
	run: (*parser).callonReservedRaw26,
	expr: &litMatcher{
	pos: position{line: 128, col: 5, offset: 2786},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 129, col: 5, offset: 2818},
	run: (*parser).callonReservedRaw28,
	expr: &litMatcher{
	pos: position{line: 129, col: 5, offset: 2818},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "Reserved",
	pos: position{line: 131, col: 1, offset: 2847},
	expr: &actionExpr{
	pos: position{line: 131, col: 12, offset: 2860},
	run: (*parser).callonReserved1,
	expr: &seqExpr{
	pos: position{line: 131, col: 12, offset: 2860},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 131, col: 12, offset: 2860},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 131, col: 14, offset: 2862},
	name: "ReservedRaw",
},
},
&ruleRefExpr{
	pos: position{line: 131, col: 26, offset: 2874},
	name: "_",
},
	},
},
},
},
{
	name: "KeywordRaw",
	pos: position{line: 133, col: 1, offset: 2895},
	expr: &choiceExpr{
	pos: position{line: 133, col: 14, offset: 2910},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 133, col: 14, offset: 2910},
	val: "if",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 134, col: 5, offset: 2919},
	val: "then",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 135, col: 5, offset: 2930},
	val: "else",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 136, col: 5, offset: 2941},
	val: "let",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 137, col: 5, offset: 2951},
	val: "in",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 138, col: 5, offset: 2960},
	val: "as",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 139, col: 5, offset: 2969},
	val: "using",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 140, col: 5, offset: 2981},
	val: "merge",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 141, col: 5, offset: 2993},
	val: "constructors",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 142, col: 5, offset: 3012},
	val: "Some",
	ignoreCase: false,
},
	},
},
},
{
	name: "If",
	pos: position{line: 144, col: 1, offset: 3020},
	expr: &seqExpr{
	pos: position{line: 144, col: 6, offset: 3027},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 144, col: 6, offset: 3027},
	val: "if",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 144, col: 11, offset: 3032},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Then",
	pos: position{line: 145, col: 1, offset: 3051},
	expr: &seqExpr{
	pos: position{line: 145, col: 8, offset: 3060},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 145, col: 8, offset: 3060},
	val: "then",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 145, col: 15, offset: 3067},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Else",
	pos: position{line: 146, col: 1, offset: 3086},
	expr: &seqExpr{
	pos: position{line: 146, col: 8, offset: 3095},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 146, col: 8, offset: 3095},
	val: "else",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 146, col: 15, offset: 3102},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Let",
	pos: position{line: 147, col: 1, offset: 3121},
	expr: &seqExpr{
	pos: position{line: 147, col: 7, offset: 3129},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 147, col: 7, offset: 3129},
	val: "let",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 147, col: 13, offset: 3135},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "In",
	pos: position{line: 148, col: 1, offset: 3154},
	expr: &seqExpr{
	pos: position{line: 148, col: 6, offset: 3161},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 148, col: 6, offset: 3161},
	val: "in",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 148, col: 11, offset: 3166},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "As",
	pos: position{line: 149, col: 1, offset: 3185},
	expr: &seqExpr{
	pos: position{line: 149, col: 6, offset: 3192},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 149, col: 6, offset: 3192},
	val: "as",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 149, col: 11, offset: 3197},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Using",
	pos: position{line: 150, col: 1, offset: 3216},
	expr: &seqExpr{
	pos: position{line: 150, col: 9, offset: 3226},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 150, col: 9, offset: 3226},
	val: "using",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 150, col: 17, offset: 3234},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Merge",
	pos: position{line: 151, col: 1, offset: 3253},
	expr: &seqExpr{
	pos: position{line: 151, col: 9, offset: 3263},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 151, col: 9, offset: 3263},
	val: "merge",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 151, col: 17, offset: 3271},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Some",
	pos: position{line: 152, col: 1, offset: 3290},
	expr: &seqExpr{
	pos: position{line: 152, col: 8, offset: 3299},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 152, col: 8, offset: 3299},
	val: "Some",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 152, col: 15, offset: 3306},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Optional",
	pos: position{line: 153, col: 1, offset: 3325},
	expr: &seqExpr{
	pos: position{line: 153, col: 12, offset: 3338},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 153, col: 12, offset: 3338},
	val: "Optional",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 153, col: 23, offset: 3349},
	name: "_",
},
	},
},
},
{
	name: "Text",
	pos: position{line: 154, col: 1, offset: 3351},
	expr: &seqExpr{
	pos: position{line: 154, col: 8, offset: 3360},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 154, col: 8, offset: 3360},
	val: "Text",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 154, col: 15, offset: 3367},
	name: "_",
},
	},
},
},
{
	name: "List",
	pos: position{line: 155, col: 1, offset: 3369},
	expr: &seqExpr{
	pos: position{line: 155, col: 8, offset: 3378},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 155, col: 8, offset: 3378},
	val: "List",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 155, col: 15, offset: 3385},
	name: "_",
},
	},
},
},
{
	name: "Equal",
	pos: position{line: 157, col: 1, offset: 3388},
	expr: &seqExpr{
	pos: position{line: 157, col: 9, offset: 3398},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 157, col: 9, offset: 3398},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 157, col: 13, offset: 3402},
	name: "_",
},
	},
},
},
{
	name: "Plus",
	pos: position{line: 158, col: 1, offset: 3404},
	expr: &seqExpr{
	pos: position{line: 158, col: 8, offset: 3413},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 158, col: 8, offset: 3413},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 158, col: 12, offset: 3417},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Times",
	pos: position{line: 159, col: 1, offset: 3436},
	expr: &seqExpr{
	pos: position{line: 159, col: 9, offset: 3446},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 159, col: 9, offset: 3446},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 159, col: 13, offset: 3450},
	name: "_",
},
	},
},
},
{
	name: "Dot",
	pos: position{line: 160, col: 1, offset: 3452},
	expr: &seqExpr{
	pos: position{line: 160, col: 7, offset: 3460},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 160, col: 7, offset: 3460},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 160, col: 11, offset: 3464},
	name: "_",
},
	},
},
},
{
	name: "OpenBrace",
	pos: position{line: 161, col: 1, offset: 3466},
	expr: &seqExpr{
	pos: position{line: 161, col: 13, offset: 3480},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 161, col: 13, offset: 3480},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 161, col: 17, offset: 3484},
	name: "_",
},
	},
},
},
{
	name: "CloseBrace",
	pos: position{line: 162, col: 1, offset: 3486},
	expr: &seqExpr{
	pos: position{line: 162, col: 14, offset: 3501},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 162, col: 14, offset: 3501},
	val: "}",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 162, col: 18, offset: 3505},
	name: "_",
},
	},
},
},
{
	name: "OpenBracket",
	pos: position{line: 163, col: 1, offset: 3507},
	expr: &seqExpr{
	pos: position{line: 163, col: 15, offset: 3523},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 163, col: 15, offset: 3523},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 163, col: 19, offset: 3527},
	name: "_",
},
	},
},
},
{
	name: "CloseBracket",
	pos: position{line: 164, col: 1, offset: 3529},
	expr: &seqExpr{
	pos: position{line: 164, col: 16, offset: 3546},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 164, col: 16, offset: 3546},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 164, col: 20, offset: 3550},
	name: "_",
},
	},
},
},
{
	name: "Comma",
	pos: position{line: 165, col: 1, offset: 3552},
	expr: &seqExpr{
	pos: position{line: 165, col: 9, offset: 3562},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 165, col: 9, offset: 3562},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 165, col: 13, offset: 3566},
	name: "_",
},
	},
},
},
{
	name: "OpenParens",
	pos: position{line: 166, col: 1, offset: 3568},
	expr: &seqExpr{
	pos: position{line: 166, col: 14, offset: 3583},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 166, col: 14, offset: 3583},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 166, col: 18, offset: 3587},
	name: "_",
},
	},
},
},
{
	name: "CloseParens",
	pos: position{line: 167, col: 1, offset: 3589},
	expr: &seqExpr{
	pos: position{line: 167, col: 15, offset: 3605},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 167, col: 15, offset: 3605},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 167, col: 19, offset: 3609},
	name: "_",
},
	},
},
},
{
	name: "At",
	pos: position{line: 168, col: 1, offset: 3611},
	expr: &seqExpr{
	pos: position{line: 168, col: 6, offset: 3618},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 168, col: 6, offset: 3618},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 168, col: 10, offset: 3622},
	name: "_",
},
	},
},
},
{
	name: "Colon",
	pos: position{line: 169, col: 1, offset: 3624},
	expr: &seqExpr{
	pos: position{line: 169, col: 9, offset: 3634},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 169, col: 9, offset: 3634},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 169, col: 13, offset: 3638},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Lambda",
	pos: position{line: 171, col: 1, offset: 3658},
	expr: &seqExpr{
	pos: position{line: 171, col: 10, offset: 3669},
	exprs: []interface{}{
&choiceExpr{
	pos: position{line: 171, col: 11, offset: 3670},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 171, col: 11, offset: 3670},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 171, col: 18, offset: 3677},
	val: "λ",
	ignoreCase: false,
},
	},
},
&ruleRefExpr{
	pos: position{line: 171, col: 23, offset: 3683},
	name: "_",
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 172, col: 1, offset: 3685},
	expr: &seqExpr{
	pos: position{line: 172, col: 10, offset: 3696},
	exprs: []interface{}{
&choiceExpr{
	pos: position{line: 172, col: 11, offset: 3697},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 172, col: 11, offset: 3697},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 172, col: 22, offset: 3708},
	val: "∀",
	ignoreCase: false,
},
	},
},
&ruleRefExpr{
	pos: position{line: 172, col: 27, offset: 3715},
	name: "_",
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 173, col: 1, offset: 3717},
	expr: &seqExpr{
	pos: position{line: 173, col: 9, offset: 3727},
	exprs: []interface{}{
&choiceExpr{
	pos: position{line: 173, col: 10, offset: 3728},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 173, col: 10, offset: 3728},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 173, col: 17, offset: 3735},
	val: "→",
	ignoreCase: false,
},
	},
},
&ruleRefExpr{
	pos: position{line: 173, col: 22, offset: 3742},
	name: "_",
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 175, col: 1, offset: 3745},
	expr: &seqExpr{
	pos: position{line: 175, col: 12, offset: 3758},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 175, col: 12, offset: 3758},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 175, col: 17, offset: 3763},
	expr: &charClassMatcher{
	pos: position{line: 175, col: 17, offset: 3763},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 175, col: 23, offset: 3769},
	expr: &charClassMatcher{
	pos: position{line: 175, col: 23, offset: 3769},
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
	pos: position{line: 177, col: 1, offset: 3777},
	expr: &actionExpr{
	pos: position{line: 177, col: 20, offset: 3798},
	run: (*parser).callonDoubleLiteralRaw1,
	expr: &seqExpr{
	pos: position{line: 177, col: 20, offset: 3798},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 177, col: 20, offset: 3798},
	expr: &charClassMatcher{
	pos: position{line: 177, col: 20, offset: 3798},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 177, col: 26, offset: 3804},
	expr: &charClassMatcher{
	pos: position{line: 177, col: 26, offset: 3804},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&choiceExpr{
	pos: position{line: 177, col: 35, offset: 3813},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 177, col: 35, offset: 3813},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 177, col: 35, offset: 3813},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 177, col: 39, offset: 3817},
	expr: &charClassMatcher{
	pos: position{line: 177, col: 39, offset: 3817},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&zeroOrOneExpr{
	pos: position{line: 177, col: 46, offset: 3824},
	expr: &ruleRefExpr{
	pos: position{line: 177, col: 46, offset: 3824},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 177, col: 58, offset: 3836},
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
	pos: position{line: 185, col: 1, offset: 3992},
	expr: &actionExpr{
	pos: position{line: 185, col: 17, offset: 4010},
	run: (*parser).callonDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 185, col: 17, offset: 4010},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 185, col: 17, offset: 4010},
	label: "d",
	expr: &ruleRefExpr{
	pos: position{line: 185, col: 19, offset: 4012},
	name: "DoubleLiteralRaw",
},
},
&ruleRefExpr{
	pos: position{line: 185, col: 36, offset: 4029},
	name: "_",
},
	},
},
},
},
{
	name: "NaturalLiteralRaw",
	pos: position{line: 187, col: 1, offset: 4050},
	expr: &actionExpr{
	pos: position{line: 187, col: 21, offset: 4072},
	run: (*parser).callonNaturalLiteralRaw1,
	expr: &oneOrMoreExpr{
	pos: position{line: 187, col: 21, offset: 4072},
	expr: &charClassMatcher{
	pos: position{line: 187, col: 21, offset: 4072},
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
	pos: position{line: 192, col: 1, offset: 4161},
	expr: &actionExpr{
	pos: position{line: 192, col: 18, offset: 4180},
	run: (*parser).callonNaturalLiteral1,
	expr: &seqExpr{
	pos: position{line: 192, col: 18, offset: 4180},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 192, col: 18, offset: 4180},
	label: "n",
	expr: &ruleRefExpr{
	pos: position{line: 192, col: 20, offset: 4182},
	name: "NaturalLiteralRaw",
},
},
&ruleRefExpr{
	pos: position{line: 192, col: 38, offset: 4200},
	name: "_",
},
	},
},
},
},
{
	name: "IntegerLiteralRaw",
	pos: position{line: 194, col: 1, offset: 4221},
	expr: &actionExpr{
	pos: position{line: 194, col: 21, offset: 4243},
	run: (*parser).callonIntegerLiteralRaw1,
	expr: &seqExpr{
	pos: position{line: 194, col: 21, offset: 4243},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 194, col: 21, offset: 4243},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&oneOrMoreExpr{
	pos: position{line: 194, col: 25, offset: 4247},
	expr: &charClassMatcher{
	pos: position{line: 194, col: 25, offset: 4247},
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
	pos: position{line: 202, col: 1, offset: 4391},
	expr: &actionExpr{
	pos: position{line: 202, col: 18, offset: 4410},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 202, col: 18, offset: 4410},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 202, col: 18, offset: 4410},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 202, col: 21, offset: 4413},
	name: "IntegerLiteralRaw",
},
},
&ruleRefExpr{
	pos: position{line: 202, col: 39, offset: 4431},
	name: "_",
},
	},
},
},
},
{
	name: "DeBruijn",
	pos: position{line: 204, col: 1, offset: 4452},
	expr: &actionExpr{
	pos: position{line: 204, col: 12, offset: 4465},
	run: (*parser).callonDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 204, col: 12, offset: 4465},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 204, col: 12, offset: 4465},
	name: "At",
},
&labeledExpr{
	pos: position{line: 204, col: 15, offset: 4468},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 204, col: 21, offset: 4474},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Identifier",
	pos: position{line: 206, col: 1, offset: 4530},
	expr: &actionExpr{
	pos: position{line: 206, col: 14, offset: 4545},
	run: (*parser).callonIdentifier1,
	expr: &seqExpr{
	pos: position{line: 206, col: 14, offset: 4545},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 206, col: 14, offset: 4545},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 206, col: 19, offset: 4550},
	name: "Label",
},
},
&labeledExpr{
	pos: position{line: 206, col: 25, offset: 4556},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 206, col: 31, offset: 4562},
	expr: &ruleRefExpr{
	pos: position{line: 206, col: 31, offset: 4562},
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
	pos: position{line: 214, col: 1, offset: 4725},
	expr: &actionExpr{
	pos: position{line: 215, col: 10, offset: 4763},
	run: (*parser).callonIdentifierReservedPrefix1,
	expr: &seqExpr{
	pos: position{line: 215, col: 10, offset: 4763},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 215, col: 10, offset: 4763},
	label: "name",
	expr: &actionExpr{
	pos: position{line: 215, col: 16, offset: 4769},
	run: (*parser).callonIdentifierReservedPrefix4,
	expr: &seqExpr{
	pos: position{line: 215, col: 16, offset: 4769},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 215, col: 16, offset: 4769},
	name: "ReservedRaw",
},
&oneOrMoreExpr{
	pos: position{line: 215, col: 28, offset: 4781},
	expr: &charClassMatcher{
	pos: position{line: 215, col: 28, offset: 4781},
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
	pos: position{line: 215, col: 75, offset: 4828},
	name: "_",
},
&labeledExpr{
	pos: position{line: 216, col: 10, offset: 4839},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 216, col: 16, offset: 4845},
	expr: &ruleRefExpr{
	pos: position{line: 216, col: 16, offset: 4845},
	name: "DeBruijn",
},
},
},
	},
},
},
},
{
	name: "Env",
	pos: position{line: 224, col: 1, offset: 5008},
	expr: &actionExpr{
	pos: position{line: 224, col: 7, offset: 5016},
	run: (*parser).callonEnv1,
	expr: &seqExpr{
	pos: position{line: 224, col: 7, offset: 5016},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 224, col: 7, offset: 5016},
	val: "env:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 224, col: 14, offset: 5023},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 224, col: 17, offset: 5026},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 224, col: 17, offset: 5026},
	name: "BashEnvironmentVariable",
},
&ruleRefExpr{
	pos: position{line: 224, col: 43, offset: 5052},
	name: "PosixEnvironmentVariable",
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 224, col: 69, offset: 5078},
	name: "_",
},
	},
},
},
},
{
	name: "BashEnvironmentVariable",
	pos: position{line: 226, col: 1, offset: 5099},
	expr: &actionExpr{
	pos: position{line: 226, col: 27, offset: 5127},
	run: (*parser).callonBashEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 226, col: 27, offset: 5127},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 226, col: 27, offset: 5127},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 226, col: 36, offset: 5136},
	expr: &charClassMatcher{
	pos: position{line: 226, col: 36, offset: 5136},
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
	pos: position{line: 230, col: 1, offset: 5207},
	expr: &actionExpr{
	pos: position{line: 230, col: 28, offset: 5236},
	run: (*parser).callonPosixEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 230, col: 28, offset: 5236},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 230, col: 28, offset: 5236},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 230, col: 32, offset: 5240},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 230, col: 34, offset: 5242},
	name: "PosixEnvironmentVariableContent",
},
},
&litMatcher{
	pos: position{line: 230, col: 66, offset: 5274},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PosixEnvironmentVariableContent",
	pos: position{line: 234, col: 1, offset: 5299},
	expr: &actionExpr{
	pos: position{line: 234, col: 35, offset: 5335},
	run: (*parser).callonPosixEnvironmentVariableContent1,
	expr: &labeledExpr{
	pos: position{line: 234, col: 35, offset: 5335},
	label: "v",
	expr: &oneOrMoreExpr{
	pos: position{line: 234, col: 37, offset: 5337},
	expr: &ruleRefExpr{
	pos: position{line: 234, col: 37, offset: 5337},
	name: "PosixEnvironmentVariableCharacter",
},
},
},
},
},
{
	name: "PosixEnvironmentVariableCharacter",
	pos: position{line: 243, col: 1, offset: 5565},
	expr: &choiceExpr{
	pos: position{line: 244, col: 7, offset: 5609},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 244, col: 7, offset: 5609},
	run: (*parser).callonPosixEnvironmentVariableCharacter2,
	expr: &litMatcher{
	pos: position{line: 244, col: 7, offset: 5609},
	val: "\\\"",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 245, col: 7, offset: 5649},
	run: (*parser).callonPosixEnvironmentVariableCharacter4,
	expr: &litMatcher{
	pos: position{line: 245, col: 7, offset: 5649},
	val: "\\\\",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 246, col: 7, offset: 5689},
	run: (*parser).callonPosixEnvironmentVariableCharacter6,
	expr: &litMatcher{
	pos: position{line: 246, col: 7, offset: 5689},
	val: "\\a",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 247, col: 7, offset: 5729},
	run: (*parser).callonPosixEnvironmentVariableCharacter8,
	expr: &litMatcher{
	pos: position{line: 247, col: 7, offset: 5729},
	val: "\\b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 248, col: 7, offset: 5769},
	run: (*parser).callonPosixEnvironmentVariableCharacter10,
	expr: &litMatcher{
	pos: position{line: 248, col: 7, offset: 5769},
	val: "\\f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 249, col: 7, offset: 5809},
	run: (*parser).callonPosixEnvironmentVariableCharacter12,
	expr: &litMatcher{
	pos: position{line: 249, col: 7, offset: 5809},
	val: "\\n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 250, col: 7, offset: 5849},
	run: (*parser).callonPosixEnvironmentVariableCharacter14,
	expr: &litMatcher{
	pos: position{line: 250, col: 7, offset: 5849},
	val: "\\r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 251, col: 7, offset: 5889},
	run: (*parser).callonPosixEnvironmentVariableCharacter16,
	expr: &litMatcher{
	pos: position{line: 251, col: 7, offset: 5889},
	val: "\\t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 252, col: 7, offset: 5929},
	run: (*parser).callonPosixEnvironmentVariableCharacter18,
	expr: &litMatcher{
	pos: position{line: 252, col: 7, offset: 5929},
	val: "\\v",
	ignoreCase: false,
},
},
&charClassMatcher{
	pos: position{line: 253, col: 7, offset: 5969},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 254, col: 7, offset: 5987},
	val: "[\\x23-\\x3c]",
	ranges: []rune{'#','<',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 255, col: 7, offset: 6005},
	val: "[\\x3e-\\x5b]",
	ranges: []rune{'>','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 256, col: 7, offset: 6023},
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
	pos: position{line: 258, col: 1, offset: 6036},
	expr: &ruleRefExpr{
	pos: position{line: 258, col: 14, offset: 6051},
	name: "Env",
},
},
{
	name: "ImportHashed",
	pos: position{line: 260, col: 1, offset: 6056},
	expr: &ruleRefExpr{
	pos: position{line: 260, col: 16, offset: 6073},
	name: "ImportType",
},
},
{
	name: "Import",
	pos: position{line: 262, col: 1, offset: 6085},
	expr: &actionExpr{
	pos: position{line: 262, col: 10, offset: 6096},
	run: (*parser).callonImport1,
	expr: &seqExpr{
	pos: position{line: 262, col: 10, offset: 6096},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 262, col: 10, offset: 6096},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 262, col: 12, offset: 6098},
	name: "ImportHashed",
},
},
&ruleRefExpr{
	pos: position{line: 262, col: 25, offset: 6111},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 262, col: 28, offset: 6114},
	name: "Text",
},
	},
},
},
},
{
	name: "LetBinding",
	pos: position{line: 268, col: 1, offset: 6324},
	expr: &actionExpr{
	pos: position{line: 268, col: 14, offset: 6339},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 268, col: 14, offset: 6339},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 268, col: 14, offset: 6339},
	name: "Let",
},
&labeledExpr{
	pos: position{line: 268, col: 18, offset: 6343},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 268, col: 24, offset: 6349},
	name: "Label",
},
},
&labeledExpr{
	pos: position{line: 268, col: 30, offset: 6355},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 268, col: 32, offset: 6357},
	expr: &ruleRefExpr{
	pos: position{line: 268, col: 32, offset: 6357},
	name: "Annotation",
},
},
},
&ruleRefExpr{
	pos: position{line: 268, col: 44, offset: 6369},
	name: "Equal",
},
&labeledExpr{
	pos: position{line: 268, col: 50, offset: 6375},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 268, col: 52, offset: 6377},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 283, col: 1, offset: 6676},
	expr: &choiceExpr{
	pos: position{line: 284, col: 7, offset: 6697},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 284, col: 7, offset: 6697},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 284, col: 7, offset: 6697},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 284, col: 7, offset: 6697},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 284, col: 14, offset: 6704},
	name: "OpenParens",
},
&labeledExpr{
	pos: position{line: 284, col: 25, offset: 6715},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 284, col: 31, offset: 6721},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 284, col: 37, offset: 6727},
	name: "Colon",
},
&labeledExpr{
	pos: position{line: 284, col: 43, offset: 6733},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 284, col: 45, offset: 6735},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 284, col: 56, offset: 6746},
	name: "CloseParens",
},
&ruleRefExpr{
	pos: position{line: 284, col: 68, offset: 6758},
	name: "Arrow",
},
&labeledExpr{
	pos: position{line: 284, col: 74, offset: 6764},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 284, col: 79, offset: 6769},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 287, col: 7, offset: 6886},
	run: (*parser).callonExpression15,
	expr: &seqExpr{
	pos: position{line: 287, col: 7, offset: 6886},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 287, col: 7, offset: 6886},
	name: "If",
},
&labeledExpr{
	pos: position{line: 287, col: 10, offset: 6889},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 287, col: 15, offset: 6894},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 287, col: 26, offset: 6905},
	name: "Then",
},
&labeledExpr{
	pos: position{line: 287, col: 31, offset: 6910},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 287, col: 33, offset: 6912},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 287, col: 44, offset: 6923},
	name: "Else",
},
&labeledExpr{
	pos: position{line: 287, col: 49, offset: 6928},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 287, col: 51, offset: 6930},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 290, col: 7, offset: 7016},
	run: (*parser).callonExpression26,
	expr: &seqExpr{
	pos: position{line: 290, col: 7, offset: 7016},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 290, col: 7, offset: 7016},
	label: "bindings",
	expr: &zeroOrMoreExpr{
	pos: position{line: 290, col: 16, offset: 7025},
	expr: &ruleRefExpr{
	pos: position{line: 290, col: 16, offset: 7025},
	name: "LetBinding",
},
},
},
&ruleRefExpr{
	pos: position{line: 290, col: 28, offset: 7037},
	name: "In",
},
&labeledExpr{
	pos: position{line: 290, col: 31, offset: 7040},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 290, col: 33, offset: 7042},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 297, col: 7, offset: 7282},
	run: (*parser).callonExpression34,
	expr: &seqExpr{
	pos: position{line: 297, col: 7, offset: 7282},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 297, col: 7, offset: 7282},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 297, col: 14, offset: 7289},
	name: "OpenParens",
},
&labeledExpr{
	pos: position{line: 297, col: 25, offset: 7300},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 297, col: 31, offset: 7306},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 297, col: 37, offset: 7312},
	name: "Colon",
},
&labeledExpr{
	pos: position{line: 297, col: 43, offset: 7318},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 297, col: 45, offset: 7320},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 297, col: 56, offset: 7331},
	name: "CloseParens",
},
&ruleRefExpr{
	pos: position{line: 297, col: 68, offset: 7343},
	name: "Arrow",
},
&labeledExpr{
	pos: position{line: 297, col: 74, offset: 7349},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 297, col: 79, offset: 7354},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 300, col: 7, offset: 7463},
	run: (*parser).callonExpression47,
	expr: &seqExpr{
	pos: position{line: 300, col: 7, offset: 7463},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 300, col: 7, offset: 7463},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 300, col: 9, offset: 7465},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 300, col: 28, offset: 7484},
	name: "Arrow",
},
&labeledExpr{
	pos: position{line: 300, col: 34, offset: 7490},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 300, col: 36, offset: 7492},
	name: "Expression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 301, col: 7, offset: 7552},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 303, col: 1, offset: 7573},
	expr: &actionExpr{
	pos: position{line: 303, col: 14, offset: 7588},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 303, col: 14, offset: 7588},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 303, col: 14, offset: 7588},
	name: "Colon",
},
&labeledExpr{
	pos: position{line: 303, col: 20, offset: 7594},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 303, col: 22, offset: 7596},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 305, col: 1, offset: 7626},
	expr: &choiceExpr{
	pos: position{line: 306, col: 5, offset: 7654},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 306, col: 5, offset: 7654},
	name: "EmptyList",
},
&actionExpr{
	pos: position{line: 307, col: 5, offset: 7668},
	run: (*parser).callonAnnotatedExpression3,
	expr: &seqExpr{
	pos: position{line: 307, col: 5, offset: 7668},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 307, col: 5, offset: 7668},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 307, col: 7, offset: 7670},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 307, col: 26, offset: 7689},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 307, col: 28, offset: 7691},
	expr: &ruleRefExpr{
	pos: position{line: 307, col: 28, offset: 7691},
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
	pos: position{line: 312, col: 1, offset: 7796},
	expr: &actionExpr{
	pos: position{line: 312, col: 13, offset: 7810},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 312, col: 13, offset: 7810},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 312, col: 13, offset: 7810},
	name: "OpenBracket",
},
&ruleRefExpr{
	pos: position{line: 312, col: 25, offset: 7822},
	name: "CloseBracket",
},
&ruleRefExpr{
	pos: position{line: 312, col: 38, offset: 7835},
	name: "Colon",
},
&ruleRefExpr{
	pos: position{line: 312, col: 44, offset: 7841},
	name: "List",
},
&labeledExpr{
	pos: position{line: 312, col: 49, offset: 7846},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 312, col: 51, offset: 7848},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 316, col: 1, offset: 7911},
	expr: &ruleRefExpr{
	pos: position{line: 316, col: 22, offset: 7934},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 318, col: 1, offset: 7955},
	expr: &ruleRefExpr{
	pos: position{line: 318, col: 23, offset: 7979},
	name: "PlusExpression",
},
},
{
	name: "MorePlus",
	pos: position{line: 320, col: 1, offset: 7995},
	expr: &actionExpr{
	pos: position{line: 320, col: 12, offset: 8008},
	run: (*parser).callonMorePlus1,
	expr: &seqExpr{
	pos: position{line: 320, col: 12, offset: 8008},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 320, col: 12, offset: 8008},
	name: "Plus",
},
&labeledExpr{
	pos: position{line: 320, col: 17, offset: 8013},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 320, col: 19, offset: 8015},
	name: "TimesExpression",
},
},
	},
},
},
},
{
	name: "PlusExpression",
	pos: position{line: 321, col: 1, offset: 8049},
	expr: &actionExpr{
	pos: position{line: 322, col: 7, offset: 8074},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 322, col: 7, offset: 8074},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 322, col: 7, offset: 8074},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 322, col: 13, offset: 8080},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 322, col: 29, offset: 8096},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 322, col: 34, offset: 8101},
	expr: &ruleRefExpr{
	pos: position{line: 322, col: 34, offset: 8101},
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
	pos: position{line: 331, col: 1, offset: 8329},
	expr: &actionExpr{
	pos: position{line: 331, col: 13, offset: 8343},
	run: (*parser).callonMoreTimes1,
	expr: &seqExpr{
	pos: position{line: 331, col: 13, offset: 8343},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 331, col: 13, offset: 8343},
	name: "Times",
},
&labeledExpr{
	pos: position{line: 331, col: 19, offset: 8349},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 331, col: 21, offset: 8351},
	name: "ApplicationExpression",
},
},
	},
},
},
},
{
	name: "TimesExpression",
	pos: position{line: 332, col: 1, offset: 8391},
	expr: &actionExpr{
	pos: position{line: 333, col: 7, offset: 8417},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 333, col: 7, offset: 8417},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 333, col: 7, offset: 8417},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 333, col: 13, offset: 8423},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 333, col: 35, offset: 8445},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 333, col: 40, offset: 8450},
	expr: &ruleRefExpr{
	pos: position{line: 333, col: 40, offset: 8450},
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
	pos: position{line: 342, col: 1, offset: 8680},
	expr: &actionExpr{
	pos: position{line: 342, col: 25, offset: 8706},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 342, col: 25, offset: 8706},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 342, col: 25, offset: 8706},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 342, col: 27, offset: 8708},
	name: "ImportExpression",
},
},
&labeledExpr{
	pos: position{line: 342, col: 44, offset: 8725},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 342, col: 49, offset: 8730},
	expr: &ruleRefExpr{
	pos: position{line: 342, col: 49, offset: 8730},
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
	pos: position{line: 351, col: 1, offset: 8960},
	expr: &choiceExpr{
	pos: position{line: 351, col: 20, offset: 8981},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 351, col: 20, offset: 8981},
	name: "Import",
},
&ruleRefExpr{
	pos: position{line: 351, col: 29, offset: 8990},
	name: "SelectorExpression",
},
	},
},
},
{
	name: "SelectorExpression",
	pos: position{line: 353, col: 1, offset: 9010},
	expr: &actionExpr{
	pos: position{line: 353, col: 22, offset: 9033},
	run: (*parser).callonSelectorExpression1,
	expr: &seqExpr{
	pos: position{line: 353, col: 22, offset: 9033},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 353, col: 22, offset: 9033},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 353, col: 24, offset: 9035},
	name: "PrimitiveExpression",
},
},
&labeledExpr{
	pos: position{line: 353, col: 44, offset: 9055},
	label: "ls",
	expr: &zeroOrMoreExpr{
	pos: position{line: 353, col: 47, offset: 9058},
	expr: &seqExpr{
	pos: position{line: 353, col: 48, offset: 9059},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 353, col: 48, offset: 9059},
	name: "Dot",
},
&ruleRefExpr{
	pos: position{line: 353, col: 52, offset: 9063},
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
	pos: position{line: 363, col: 1, offset: 9293},
	expr: &choiceExpr{
	pos: position{line: 364, col: 7, offset: 9323},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 364, col: 7, offset: 9323},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 365, col: 7, offset: 9343},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 366, col: 7, offset: 9364},
	name: "IntegerLiteral",
},
&actionExpr{
	pos: position{line: 367, col: 7, offset: 9385},
	run: (*parser).callonPrimitiveExpression5,
	expr: &litMatcher{
	pos: position{line: 367, col: 7, offset: 9385},
	val: "-Infinity",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 368, col: 7, offset: 9443},
	name: "TextLiteral",
},
&actionExpr{
	pos: position{line: 369, col: 7, offset: 9461},
	run: (*parser).callonPrimitiveExpression8,
	expr: &seqExpr{
	pos: position{line: 369, col: 7, offset: 9461},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 369, col: 7, offset: 9461},
	name: "OpenBrace",
},
&labeledExpr{
	pos: position{line: 369, col: 17, offset: 9471},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 369, col: 19, offset: 9473},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 369, col: 39, offset: 9493},
	name: "CloseBrace",
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 370, col: 7, offset: 9528},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 371, col: 7, offset: 9554},
	name: "IdentifierReservedPrefix",
},
&ruleRefExpr{
	pos: position{line: 372, col: 7, offset: 9585},
	name: "Reserved",
},
&ruleRefExpr{
	pos: position{line: 373, col: 7, offset: 9600},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 374, col: 7, offset: 9617},
	run: (*parser).callonPrimitiveExpression18,
	expr: &seqExpr{
	pos: position{line: 374, col: 7, offset: 9617},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 374, col: 7, offset: 9617},
	name: "OpenParens",
},
&labeledExpr{
	pos: position{line: 374, col: 18, offset: 9628},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 374, col: 20, offset: 9630},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 374, col: 31, offset: 9641},
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
	pos: position{line: 376, col: 1, offset: 9672},
	expr: &choiceExpr{
	pos: position{line: 377, col: 7, offset: 9702},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 377, col: 7, offset: 9702},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &ruleRefExpr{
	pos: position{line: 377, col: 7, offset: 9702},
	name: "Equal",
},
},
&ruleRefExpr{
	pos: position{line: 378, col: 7, offset: 9759},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 379, col: 7, offset: 9784},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 380, col: 7, offset: 9812},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 380, col: 7, offset: 9812},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 382, col: 1, offset: 9858},
	expr: &actionExpr{
	pos: position{line: 382, col: 19, offset: 9878},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 382, col: 19, offset: 9878},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 382, col: 19, offset: 9878},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 382, col: 24, offset: 9883},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 382, col: 30, offset: 9889},
	name: "Colon",
},
&labeledExpr{
	pos: position{line: 382, col: 36, offset: 9895},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 382, col: 41, offset: 9900},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 385, col: 1, offset: 9957},
	expr: &actionExpr{
	pos: position{line: 385, col: 18, offset: 9976},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 385, col: 18, offset: 9976},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 385, col: 18, offset: 9976},
	name: "Comma",
},
&labeledExpr{
	pos: position{line: 385, col: 24, offset: 9982},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 385, col: 26, offset: 9984},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 386, col: 1, offset: 10016},
	expr: &actionExpr{
	pos: position{line: 387, col: 7, offset: 10045},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 387, col: 7, offset: 10045},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 387, col: 7, offset: 10045},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 387, col: 13, offset: 10051},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 387, col: 29, offset: 10067},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 387, col: 34, offset: 10072},
	expr: &ruleRefExpr{
	pos: position{line: 387, col: 34, offset: 10072},
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
	pos: position{line: 397, col: 1, offset: 10468},
	expr: &actionExpr{
	pos: position{line: 397, col: 22, offset: 10491},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 397, col: 22, offset: 10491},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 397, col: 22, offset: 10491},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 397, col: 27, offset: 10496},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 397, col: 33, offset: 10502},
	name: "Equal",
},
&labeledExpr{
	pos: position{line: 397, col: 39, offset: 10508},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 397, col: 44, offset: 10513},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 400, col: 1, offset: 10570},
	expr: &actionExpr{
	pos: position{line: 400, col: 21, offset: 10592},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 400, col: 21, offset: 10592},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 400, col: 21, offset: 10592},
	name: "Comma",
},
&labeledExpr{
	pos: position{line: 400, col: 27, offset: 10598},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 400, col: 29, offset: 10600},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 401, col: 1, offset: 10635},
	expr: &actionExpr{
	pos: position{line: 402, col: 7, offset: 10667},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 402, col: 7, offset: 10667},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 402, col: 7, offset: 10667},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 402, col: 13, offset: 10673},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 402, col: 32, offset: 10692},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 402, col: 37, offset: 10697},
	expr: &ruleRefExpr{
	pos: position{line: 402, col: 37, offset: 10697},
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
	pos: position{line: 412, col: 1, offset: 11099},
	expr: &actionExpr{
	pos: position{line: 412, col: 12, offset: 11112},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 412, col: 12, offset: 11112},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 412, col: 12, offset: 11112},
	name: "Comma",
},
&labeledExpr{
	pos: position{line: 412, col: 18, offset: 11118},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 412, col: 20, offset: 11120},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 414, col: 1, offset: 11148},
	expr: &actionExpr{
	pos: position{line: 415, col: 7, offset: 11178},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 415, col: 7, offset: 11178},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 415, col: 7, offset: 11178},
	name: "OpenBracket",
},
&labeledExpr{
	pos: position{line: 415, col: 19, offset: 11190},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 415, col: 25, offset: 11196},
	name: "Expression",
},
},
&labeledExpr{
	pos: position{line: 415, col: 36, offset: 11207},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 415, col: 41, offset: 11212},
	expr: &ruleRefExpr{
	pos: position{line: 415, col: 41, offset: 11212},
	name: "MoreList",
},
},
},
&ruleRefExpr{
	pos: position{line: 415, col: 51, offset: 11222},
	name: "CloseBracket",
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 425, col: 1, offset: 11507},
	expr: &notExpr{
	pos: position{line: 425, col: 7, offset: 11515},
	expr: &anyMatcher{
	line: 425, col: 8, offset: 11516,
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

func (c *current) onReservedRaw2() (interface{}, error) {
 return Bool, nil 
}

func (p *parser) callonReservedRaw2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw2()
}

func (c *current) onReservedRaw6() (interface{}, error) {
 return Natural, nil 
}

func (p *parser) callonReservedRaw6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw6()
}

func (c *current) onReservedRaw8() (interface{}, error) {
 return Integer, nil 
}

func (p *parser) callonReservedRaw8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw8()
}

func (c *current) onReservedRaw10() (interface{}, error) {
 return Double, nil 
}

func (p *parser) callonReservedRaw10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw10()
}

func (c *current) onReservedRaw12() (interface{}, error) {
 return Text, nil 
}

func (p *parser) callonReservedRaw12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw12()
}

func (c *current) onReservedRaw14() (interface{}, error) {
 return List, nil 
}

func (p *parser) callonReservedRaw14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw14()
}

func (c *current) onReservedRaw16() (interface{}, error) {
 return True, nil 
}

func (p *parser) callonReservedRaw16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw16()
}

func (c *current) onReservedRaw18() (interface{}, error) {
 return False, nil 
}

func (p *parser) callonReservedRaw18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw18()
}

func (c *current) onReservedRaw20() (interface{}, error) {
 return DoubleLit(math.NaN()), nil 
}

func (p *parser) callonReservedRaw20() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw20()
}

func (c *current) onReservedRaw22() (interface{}, error) {
 return DoubleLit(math.Inf(1)), nil 
}

func (p *parser) callonReservedRaw22() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw22()
}

func (c *current) onReservedRaw24() (interface{}, error) {
 return Type, nil 
}

func (p *parser) callonReservedRaw24() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw24()
}

func (c *current) onReservedRaw26() (interface{}, error) {
 return Kind, nil 
}

func (p *parser) callonReservedRaw26() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw26()
}

func (c *current) onReservedRaw28() (interface{}, error) {
 return Sort, nil 
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
      return DoubleLit(d), nil
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
      return NaturalLit(i), err
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
      return IntegerLit(i), nil
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
 return int(index.(NaturalLit)), nil 
}

func (p *parser) callonDeBruijn1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDeBruijn1(stack["index"])
}

func (c *current) onIdentifier1(name, index interface{}) (interface{}, error) {
    if index != nil {
        return Var{Name:name.(string), Index:index.(int)}, nil
    } else {
        return Var{Name:name.(string)}, nil
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
        return Var{Name:name.(string), Index:index.(int)}, nil
    } else {
        return Var{Name:name.(string)}, nil
    }
}

func (p *parser) callonIdentifierReservedPrefix1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifierReservedPrefix1(stack["name"], stack["index"])
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
  return Embed(Import{EnvVar: string(c.text)}), nil
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
  return Embed(Import{EnvVar: b.String()}), nil
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

func (c *current) onImport1(i interface{}) (interface{}, error) {
 return i, nil 
}

func (p *parser) callonImport1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onImport1(stack["i"])
}

func (c *current) onLetBinding1(label, a, v interface{}) (interface{}, error) {
    if a != nil {
        return Binding{
            Variable: label.(string),
            Annotation: a.(Expr),
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

func (c *current) onExpression15(cond, t, f interface{}) (interface{}, error) {
          return BoolIf{cond.(Expr),t.(Expr),f.(Expr)},nil
      
}

func (p *parser) callonExpression15() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression15(stack["cond"], stack["t"], stack["f"])
}

func (c *current) onExpression26(bindings, b interface{}) (interface{}, error) {
        bs := make([]Binding, len(bindings.([]interface{})))
        for i, binding := range bindings.([]interface{}) {
            bs[i] = binding.(Binding)
        }
        return MakeLet(b.(Expr), bs...), nil
      
}

func (p *parser) callonExpression26() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression26(stack["bindings"], stack["b"])
}

func (c *current) onExpression34(label, t, body interface{}) (interface{}, error) {
          return &Pi{Label:label.(string), Type:t.(Expr), Body: body.(Expr)}, nil
      
}

func (p *parser) callonExpression34() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression34(stack["label"], stack["t"], stack["body"])
}

func (c *current) onExpression47(o, e interface{}) (interface{}, error) {
 return &Pi{"_",o.(Expr),e.(Expr)}, nil 
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
        return Annot{e.(Expr), a.(Expr)}, nil
    
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
              e = &App{Fn:e, Arg: arg.(Expr)}
          }
          return e,nil
      
}

func (p *parser) callonApplicationExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onApplicationExpression1(stack["f"], stack["rest"])
}

func (c *current) onSelectorExpression1(e, ls interface{}) (interface{}, error) {
    expr := e.(Expr)
    labels := ls.([]interface{})
    for _, labelSelector := range labels {
        label := labelSelector.([]interface{})[1]
        expr = Field{expr, label.(string)}
    }
    return expr, nil
}

func (p *parser) callonSelectorExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSelectorExpression1(stack["e"], stack["ls"])
}

func (c *current) onPrimitiveExpression5() (interface{}, error) {
 return DoubleLit(math.Inf(-1)), nil 
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

