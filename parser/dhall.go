
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
	name: "NonreservedLabel",
	pos: position{line: 55, col: 1, offset: 963},
	expr: &actionExpr{
	pos: position{line: 55, col: 20, offset: 984},
	run: (*parser).callonNonreservedLabel1,
	expr: &seqExpr{
	pos: position{line: 55, col: 20, offset: 984},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 55, col: 20, offset: 984},
	expr: &ruleRefExpr{
	pos: position{line: 55, col: 21, offset: 985},
	name: "ReservedRaw",
},
},
&labeledExpr{
	pos: position{line: 55, col: 33, offset: 997},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 55, col: 39, offset: 1003},
	name: "Label",
},
},
	},
},
},
},
{
	name: "EscapedChar",
	pos: position{line: 59, col: 1, offset: 1065},
	expr: &actionExpr{
	pos: position{line: 60, col: 3, offset: 1083},
	run: (*parser).callonEscapedChar1,
	expr: &seqExpr{
	pos: position{line: 60, col: 3, offset: 1083},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 60, col: 3, offset: 1083},
	val: "\\",
	ignoreCase: false,
},
&choiceExpr{
	pos: position{line: 61, col: 5, offset: 1092},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 61, col: 5, offset: 1092},
	val: "\"",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 62, col: 10, offset: 1105},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 63, col: 10, offset: 1118},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 64, col: 10, offset: 1132},
	val: "/",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 65, col: 10, offset: 1145},
	val: "b",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 66, col: 10, offset: 1158},
	val: "f",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 67, col: 10, offset: 1171},
	val: "n",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 68, col: 10, offset: 1184},
	val: "r",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 69, col: 10, offset: 1197},
	val: "t",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 70, col: 10, offset: 1210},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 70, col: 10, offset: 1210},
	val: "u",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 70, col: 14, offset: 1214},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 70, col: 21, offset: 1221},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 70, col: 28, offset: 1228},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 70, col: 35, offset: 1235},
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
	pos: position{line: 91, col: 1, offset: 1678},
	expr: &choiceExpr{
	pos: position{line: 92, col: 6, offset: 1704},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 92, col: 6, offset: 1704},
	run: (*parser).callonDoubleQuoteChunk2,
	expr: &seqExpr{
	pos: position{line: 92, col: 6, offset: 1704},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 92, col: 6, offset: 1704},
	val: "${",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 92, col: 11, offset: 1709},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 92, col: 13, offset: 1711},
	name: "CompleteExpression",
},
},
&litMatcher{
	pos: position{line: 92, col: 32, offset: 1730},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 93, col: 6, offset: 1757},
	name: "EscapedChar",
},
&charClassMatcher{
	pos: position{line: 94, col: 6, offset: 1774},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 95, col: 6, offset: 1791},
	val: "[\\x23-\\x5b]",
	ranges: []rune{'#','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 96, col: 6, offset: 1808},
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
	pos: position{line: 98, col: 1, offset: 1827},
	expr: &actionExpr{
	pos: position{line: 98, col: 22, offset: 1850},
	run: (*parser).callonDoubleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 98, col: 22, offset: 1850},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 98, col: 22, offset: 1850},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 98, col: 26, offset: 1854},
	label: "chunks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 98, col: 33, offset: 1861},
	expr: &ruleRefExpr{
	pos: position{line: 98, col: 33, offset: 1861},
	name: "DoubleQuoteChunk",
},
},
},
&litMatcher{
	pos: position{line: 98, col: 51, offset: 1879},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TextLiteral",
	pos: position{line: 115, col: 1, offset: 2347},
	expr: &actionExpr{
	pos: position{line: 115, col: 15, offset: 2363},
	run: (*parser).callonTextLiteral1,
	expr: &seqExpr{
	pos: position{line: 115, col: 15, offset: 2363},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 115, col: 15, offset: 2363},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 115, col: 17, offset: 2365},
	name: "DoubleQuoteLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 115, col: 36, offset: 2384},
	name: "_",
},
	},
},
},
},
{
	name: "ReservedRaw",
	pos: position{line: 117, col: 1, offset: 2405},
	expr: &choiceExpr{
	pos: position{line: 117, col: 15, offset: 2421},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 117, col: 15, offset: 2421},
	run: (*parser).callonReservedRaw2,
	expr: &litMatcher{
	pos: position{line: 117, col: 15, offset: 2421},
	val: "Bool",
	ignoreCase: false,
},
},
&litMatcher{
	pos: position{line: 118, col: 5, offset: 2453},
	val: "Optional",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 119, col: 5, offset: 2468},
	val: "None",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 120, col: 5, offset: 2479},
	run: (*parser).callonReservedRaw6,
	expr: &litMatcher{
	pos: position{line: 120, col: 5, offset: 2479},
	val: "Natural",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 121, col: 5, offset: 2517},
	run: (*parser).callonReservedRaw8,
	expr: &litMatcher{
	pos: position{line: 121, col: 5, offset: 2517},
	val: "Integer",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 122, col: 5, offset: 2555},
	run: (*parser).callonReservedRaw10,
	expr: &litMatcher{
	pos: position{line: 122, col: 5, offset: 2555},
	val: "Double",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 123, col: 5, offset: 2591},
	run: (*parser).callonReservedRaw12,
	expr: &litMatcher{
	pos: position{line: 123, col: 5, offset: 2591},
	val: "Text",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 124, col: 5, offset: 2623},
	run: (*parser).callonReservedRaw14,
	expr: &litMatcher{
	pos: position{line: 124, col: 5, offset: 2623},
	val: "List",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 125, col: 5, offset: 2655},
	run: (*parser).callonReservedRaw16,
	expr: &litMatcher{
	pos: position{line: 125, col: 5, offset: 2655},
	val: "True",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 126, col: 5, offset: 2687},
	run: (*parser).callonReservedRaw18,
	expr: &litMatcher{
	pos: position{line: 126, col: 5, offset: 2687},
	val: "False",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 127, col: 5, offset: 2721},
	run: (*parser).callonReservedRaw20,
	expr: &litMatcher{
	pos: position{line: 127, col: 5, offset: 2721},
	val: "NaN",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 128, col: 5, offset: 2769},
	run: (*parser).callonReservedRaw22,
	expr: &litMatcher{
	pos: position{line: 128, col: 5, offset: 2769},
	val: "Infinity",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 129, col: 5, offset: 2823},
	run: (*parser).callonReservedRaw24,
	expr: &litMatcher{
	pos: position{line: 129, col: 5, offset: 2823},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 130, col: 5, offset: 2855},
	run: (*parser).callonReservedRaw26,
	expr: &litMatcher{
	pos: position{line: 130, col: 5, offset: 2855},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 131, col: 5, offset: 2887},
	run: (*parser).callonReservedRaw28,
	expr: &litMatcher{
	pos: position{line: 131, col: 5, offset: 2887},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "Reserved",
	pos: position{line: 133, col: 1, offset: 2916},
	expr: &actionExpr{
	pos: position{line: 133, col: 12, offset: 2929},
	run: (*parser).callonReserved1,
	expr: &seqExpr{
	pos: position{line: 133, col: 12, offset: 2929},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 133, col: 12, offset: 2929},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 133, col: 14, offset: 2931},
	name: "ReservedRaw",
},
},
&ruleRefExpr{
	pos: position{line: 133, col: 26, offset: 2943},
	name: "_",
},
	},
},
},
},
{
	name: "KeywordRaw",
	pos: position{line: 135, col: 1, offset: 2964},
	expr: &choiceExpr{
	pos: position{line: 135, col: 14, offset: 2979},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 135, col: 14, offset: 2979},
	val: "if",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 136, col: 5, offset: 2988},
	val: "then",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 137, col: 5, offset: 2999},
	val: "else",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 138, col: 5, offset: 3010},
	val: "let",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 139, col: 5, offset: 3020},
	val: "in",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 140, col: 5, offset: 3029},
	val: "as",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 141, col: 5, offset: 3038},
	val: "using",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 142, col: 5, offset: 3050},
	val: "merge",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 143, col: 5, offset: 3062},
	val: "constructors",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 144, col: 5, offset: 3081},
	val: "Some",
	ignoreCase: false,
},
	},
},
},
{
	name: "If",
	pos: position{line: 146, col: 1, offset: 3089},
	expr: &seqExpr{
	pos: position{line: 146, col: 6, offset: 3096},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 146, col: 6, offset: 3096},
	val: "if",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 146, col: 11, offset: 3101},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Then",
	pos: position{line: 147, col: 1, offset: 3120},
	expr: &seqExpr{
	pos: position{line: 147, col: 8, offset: 3129},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 147, col: 8, offset: 3129},
	val: "then",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 147, col: 15, offset: 3136},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Else",
	pos: position{line: 148, col: 1, offset: 3155},
	expr: &seqExpr{
	pos: position{line: 148, col: 8, offset: 3164},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 148, col: 8, offset: 3164},
	val: "else",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 148, col: 15, offset: 3171},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Let",
	pos: position{line: 149, col: 1, offset: 3190},
	expr: &seqExpr{
	pos: position{line: 149, col: 7, offset: 3198},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 149, col: 7, offset: 3198},
	val: "let",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 149, col: 13, offset: 3204},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "In",
	pos: position{line: 150, col: 1, offset: 3223},
	expr: &seqExpr{
	pos: position{line: 150, col: 6, offset: 3230},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 150, col: 6, offset: 3230},
	val: "in",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 150, col: 11, offset: 3235},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "As",
	pos: position{line: 151, col: 1, offset: 3254},
	expr: &seqExpr{
	pos: position{line: 151, col: 6, offset: 3261},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 151, col: 6, offset: 3261},
	val: "as",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 151, col: 11, offset: 3266},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Using",
	pos: position{line: 152, col: 1, offset: 3285},
	expr: &seqExpr{
	pos: position{line: 152, col: 9, offset: 3295},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 152, col: 9, offset: 3295},
	val: "using",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 152, col: 17, offset: 3303},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Merge",
	pos: position{line: 153, col: 1, offset: 3322},
	expr: &seqExpr{
	pos: position{line: 153, col: 9, offset: 3332},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 153, col: 9, offset: 3332},
	val: "merge",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 153, col: 17, offset: 3340},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Some",
	pos: position{line: 154, col: 1, offset: 3359},
	expr: &seqExpr{
	pos: position{line: 154, col: 8, offset: 3368},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 154, col: 8, offset: 3368},
	val: "Some",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 154, col: 15, offset: 3375},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Optional",
	pos: position{line: 155, col: 1, offset: 3394},
	expr: &seqExpr{
	pos: position{line: 155, col: 12, offset: 3407},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 155, col: 12, offset: 3407},
	val: "Optional",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 155, col: 23, offset: 3418},
	name: "_",
},
	},
},
},
{
	name: "Text",
	pos: position{line: 156, col: 1, offset: 3420},
	expr: &seqExpr{
	pos: position{line: 156, col: 8, offset: 3429},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 156, col: 8, offset: 3429},
	val: "Text",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 156, col: 15, offset: 3436},
	name: "_",
},
	},
},
},
{
	name: "List",
	pos: position{line: 157, col: 1, offset: 3438},
	expr: &seqExpr{
	pos: position{line: 157, col: 8, offset: 3447},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 157, col: 8, offset: 3447},
	val: "List",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 157, col: 15, offset: 3454},
	name: "_",
},
	},
},
},
{
	name: "Equal",
	pos: position{line: 159, col: 1, offset: 3457},
	expr: &seqExpr{
	pos: position{line: 159, col: 9, offset: 3467},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 159, col: 9, offset: 3467},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 159, col: 13, offset: 3471},
	name: "_",
},
	},
},
},
{
	name: "Plus",
	pos: position{line: 160, col: 1, offset: 3473},
	expr: &seqExpr{
	pos: position{line: 160, col: 8, offset: 3482},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 160, col: 8, offset: 3482},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 160, col: 12, offset: 3486},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Times",
	pos: position{line: 161, col: 1, offset: 3505},
	expr: &seqExpr{
	pos: position{line: 161, col: 9, offset: 3515},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 161, col: 9, offset: 3515},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 161, col: 13, offset: 3519},
	name: "_",
},
	},
},
},
{
	name: "Dot",
	pos: position{line: 162, col: 1, offset: 3521},
	expr: &seqExpr{
	pos: position{line: 162, col: 7, offset: 3529},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 162, col: 7, offset: 3529},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 162, col: 11, offset: 3533},
	name: "_",
},
	},
},
},
{
	name: "OpenBrace",
	pos: position{line: 163, col: 1, offset: 3535},
	expr: &seqExpr{
	pos: position{line: 163, col: 13, offset: 3549},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 163, col: 13, offset: 3549},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 163, col: 17, offset: 3553},
	name: "_",
},
	},
},
},
{
	name: "CloseBrace",
	pos: position{line: 164, col: 1, offset: 3555},
	expr: &seqExpr{
	pos: position{line: 164, col: 14, offset: 3570},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 164, col: 14, offset: 3570},
	val: "}",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 164, col: 18, offset: 3574},
	name: "_",
},
	},
},
},
{
	name: "OpenBracket",
	pos: position{line: 165, col: 1, offset: 3576},
	expr: &seqExpr{
	pos: position{line: 165, col: 15, offset: 3592},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 165, col: 15, offset: 3592},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 165, col: 19, offset: 3596},
	name: "_",
},
	},
},
},
{
	name: "CloseBracket",
	pos: position{line: 166, col: 1, offset: 3598},
	expr: &seqExpr{
	pos: position{line: 166, col: 16, offset: 3615},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 166, col: 16, offset: 3615},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 166, col: 20, offset: 3619},
	name: "_",
},
	},
},
},
{
	name: "Comma",
	pos: position{line: 167, col: 1, offset: 3621},
	expr: &seqExpr{
	pos: position{line: 167, col: 9, offset: 3631},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 167, col: 9, offset: 3631},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 167, col: 13, offset: 3635},
	name: "_",
},
	},
},
},
{
	name: "OpenParens",
	pos: position{line: 168, col: 1, offset: 3637},
	expr: &seqExpr{
	pos: position{line: 168, col: 14, offset: 3652},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 168, col: 14, offset: 3652},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 168, col: 18, offset: 3656},
	name: "_",
},
	},
},
},
{
	name: "CloseParens",
	pos: position{line: 169, col: 1, offset: 3658},
	expr: &seqExpr{
	pos: position{line: 169, col: 15, offset: 3674},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 169, col: 15, offset: 3674},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 169, col: 19, offset: 3678},
	name: "_",
},
	},
},
},
{
	name: "At",
	pos: position{line: 170, col: 1, offset: 3680},
	expr: &seqExpr{
	pos: position{line: 170, col: 6, offset: 3687},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 170, col: 6, offset: 3687},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 170, col: 10, offset: 3691},
	name: "_",
},
	},
},
},
{
	name: "Colon",
	pos: position{line: 171, col: 1, offset: 3693},
	expr: &seqExpr{
	pos: position{line: 171, col: 9, offset: 3703},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 171, col: 9, offset: 3703},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 171, col: 13, offset: 3707},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Lambda",
	pos: position{line: 173, col: 1, offset: 3727},
	expr: &seqExpr{
	pos: position{line: 173, col: 10, offset: 3738},
	exprs: []interface{}{
&choiceExpr{
	pos: position{line: 173, col: 11, offset: 3739},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 173, col: 11, offset: 3739},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 173, col: 18, offset: 3746},
	val: "λ",
	ignoreCase: false,
},
	},
},
&ruleRefExpr{
	pos: position{line: 173, col: 23, offset: 3752},
	name: "_",
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 174, col: 1, offset: 3754},
	expr: &seqExpr{
	pos: position{line: 174, col: 10, offset: 3765},
	exprs: []interface{}{
&choiceExpr{
	pos: position{line: 174, col: 11, offset: 3766},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 174, col: 11, offset: 3766},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 174, col: 22, offset: 3777},
	val: "∀",
	ignoreCase: false,
},
	},
},
&ruleRefExpr{
	pos: position{line: 174, col: 27, offset: 3784},
	name: "_",
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 175, col: 1, offset: 3786},
	expr: &seqExpr{
	pos: position{line: 175, col: 9, offset: 3796},
	exprs: []interface{}{
&choiceExpr{
	pos: position{line: 175, col: 10, offset: 3797},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 175, col: 10, offset: 3797},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 175, col: 17, offset: 3804},
	val: "→",
	ignoreCase: false,
},
	},
},
&ruleRefExpr{
	pos: position{line: 175, col: 22, offset: 3811},
	name: "_",
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 177, col: 1, offset: 3814},
	expr: &seqExpr{
	pos: position{line: 177, col: 12, offset: 3827},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 177, col: 12, offset: 3827},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 177, col: 17, offset: 3832},
	expr: &charClassMatcher{
	pos: position{line: 177, col: 17, offset: 3832},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 177, col: 23, offset: 3838},
	expr: &charClassMatcher{
	pos: position{line: 177, col: 23, offset: 3838},
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
	pos: position{line: 179, col: 1, offset: 3846},
	expr: &actionExpr{
	pos: position{line: 179, col: 20, offset: 3867},
	run: (*parser).callonDoubleLiteralRaw1,
	expr: &seqExpr{
	pos: position{line: 179, col: 20, offset: 3867},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 179, col: 20, offset: 3867},
	expr: &charClassMatcher{
	pos: position{line: 179, col: 20, offset: 3867},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 179, col: 26, offset: 3873},
	expr: &charClassMatcher{
	pos: position{line: 179, col: 26, offset: 3873},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&choiceExpr{
	pos: position{line: 179, col: 35, offset: 3882},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 179, col: 35, offset: 3882},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 179, col: 35, offset: 3882},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 179, col: 39, offset: 3886},
	expr: &charClassMatcher{
	pos: position{line: 179, col: 39, offset: 3886},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&zeroOrOneExpr{
	pos: position{line: 179, col: 46, offset: 3893},
	expr: &ruleRefExpr{
	pos: position{line: 179, col: 46, offset: 3893},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 179, col: 58, offset: 3905},
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
	pos: position{line: 187, col: 1, offset: 4061},
	expr: &actionExpr{
	pos: position{line: 187, col: 17, offset: 4079},
	run: (*parser).callonDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 187, col: 17, offset: 4079},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 187, col: 17, offset: 4079},
	label: "d",
	expr: &ruleRefExpr{
	pos: position{line: 187, col: 19, offset: 4081},
	name: "DoubleLiteralRaw",
},
},
&ruleRefExpr{
	pos: position{line: 187, col: 36, offset: 4098},
	name: "_",
},
	},
},
},
},
{
	name: "NaturalLiteralRaw",
	pos: position{line: 189, col: 1, offset: 4119},
	expr: &actionExpr{
	pos: position{line: 189, col: 21, offset: 4141},
	run: (*parser).callonNaturalLiteralRaw1,
	expr: &oneOrMoreExpr{
	pos: position{line: 189, col: 21, offset: 4141},
	expr: &charClassMatcher{
	pos: position{line: 189, col: 21, offset: 4141},
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
	pos: position{line: 194, col: 1, offset: 4230},
	expr: &actionExpr{
	pos: position{line: 194, col: 18, offset: 4249},
	run: (*parser).callonNaturalLiteral1,
	expr: &seqExpr{
	pos: position{line: 194, col: 18, offset: 4249},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 194, col: 18, offset: 4249},
	label: "n",
	expr: &ruleRefExpr{
	pos: position{line: 194, col: 20, offset: 4251},
	name: "NaturalLiteralRaw",
},
},
&ruleRefExpr{
	pos: position{line: 194, col: 38, offset: 4269},
	name: "_",
},
	},
},
},
},
{
	name: "IntegerLiteralRaw",
	pos: position{line: 196, col: 1, offset: 4290},
	expr: &actionExpr{
	pos: position{line: 196, col: 21, offset: 4312},
	run: (*parser).callonIntegerLiteralRaw1,
	expr: &seqExpr{
	pos: position{line: 196, col: 21, offset: 4312},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 196, col: 21, offset: 4312},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&oneOrMoreExpr{
	pos: position{line: 196, col: 25, offset: 4316},
	expr: &charClassMatcher{
	pos: position{line: 196, col: 25, offset: 4316},
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
	pos: position{line: 204, col: 1, offset: 4460},
	expr: &actionExpr{
	pos: position{line: 204, col: 18, offset: 4479},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 204, col: 18, offset: 4479},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 204, col: 18, offset: 4479},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 204, col: 21, offset: 4482},
	name: "IntegerLiteralRaw",
},
},
&ruleRefExpr{
	pos: position{line: 204, col: 39, offset: 4500},
	name: "_",
},
	},
},
},
},
{
	name: "DeBruijn",
	pos: position{line: 206, col: 1, offset: 4521},
	expr: &actionExpr{
	pos: position{line: 206, col: 12, offset: 4534},
	run: (*parser).callonDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 206, col: 12, offset: 4534},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 206, col: 12, offset: 4534},
	name: "At",
},
&labeledExpr{
	pos: position{line: 206, col: 15, offset: 4537},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 206, col: 21, offset: 4543},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Identifier",
	pos: position{line: 208, col: 1, offset: 4599},
	expr: &actionExpr{
	pos: position{line: 208, col: 14, offset: 4614},
	run: (*parser).callonIdentifier1,
	expr: &seqExpr{
	pos: position{line: 208, col: 14, offset: 4614},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 208, col: 14, offset: 4614},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 208, col: 19, offset: 4619},
	name: "Label",
},
},
&labeledExpr{
	pos: position{line: 208, col: 25, offset: 4625},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 208, col: 31, offset: 4631},
	expr: &ruleRefExpr{
	pos: position{line: 208, col: 31, offset: 4631},
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
	pos: position{line: 216, col: 1, offset: 4794},
	expr: &actionExpr{
	pos: position{line: 217, col: 10, offset: 4832},
	run: (*parser).callonIdentifierReservedPrefix1,
	expr: &seqExpr{
	pos: position{line: 217, col: 10, offset: 4832},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 217, col: 10, offset: 4832},
	label: "name",
	expr: &actionExpr{
	pos: position{line: 217, col: 16, offset: 4838},
	run: (*parser).callonIdentifierReservedPrefix4,
	expr: &seqExpr{
	pos: position{line: 217, col: 16, offset: 4838},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 217, col: 16, offset: 4838},
	name: "ReservedRaw",
},
&oneOrMoreExpr{
	pos: position{line: 217, col: 28, offset: 4850},
	expr: &charClassMatcher{
	pos: position{line: 217, col: 28, offset: 4850},
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
	pos: position{line: 217, col: 75, offset: 4897},
	name: "_",
},
&labeledExpr{
	pos: position{line: 218, col: 10, offset: 4908},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 218, col: 16, offset: 4914},
	expr: &ruleRefExpr{
	pos: position{line: 218, col: 16, offset: 4914},
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
	pos: position{line: 226, col: 1, offset: 5077},
	expr: &actionExpr{
	pos: position{line: 226, col: 7, offset: 5085},
	run: (*parser).callonEnv1,
	expr: &seqExpr{
	pos: position{line: 226, col: 7, offset: 5085},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 226, col: 7, offset: 5085},
	val: "env:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 226, col: 14, offset: 5092},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 226, col: 17, offset: 5095},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 226, col: 17, offset: 5095},
	name: "BashEnvironmentVariable",
},
&ruleRefExpr{
	pos: position{line: 226, col: 43, offset: 5121},
	name: "PosixEnvironmentVariable",
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 226, col: 69, offset: 5147},
	name: "_",
},
	},
},
},
},
{
	name: "BashEnvironmentVariable",
	pos: position{line: 228, col: 1, offset: 5168},
	expr: &actionExpr{
	pos: position{line: 228, col: 27, offset: 5196},
	run: (*parser).callonBashEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 228, col: 27, offset: 5196},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 228, col: 27, offset: 5196},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 228, col: 36, offset: 5205},
	expr: &charClassMatcher{
	pos: position{line: 228, col: 36, offset: 5205},
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
	pos: position{line: 232, col: 1, offset: 5273},
	expr: &actionExpr{
	pos: position{line: 232, col: 28, offset: 5302},
	run: (*parser).callonPosixEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 232, col: 28, offset: 5302},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 232, col: 28, offset: 5302},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 232, col: 32, offset: 5306},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 232, col: 34, offset: 5308},
	name: "PosixEnvironmentVariableContent",
},
},
&litMatcher{
	pos: position{line: 232, col: 66, offset: 5340},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PosixEnvironmentVariableContent",
	pos: position{line: 236, col: 1, offset: 5365},
	expr: &actionExpr{
	pos: position{line: 236, col: 35, offset: 5401},
	run: (*parser).callonPosixEnvironmentVariableContent1,
	expr: &labeledExpr{
	pos: position{line: 236, col: 35, offset: 5401},
	label: "v",
	expr: &oneOrMoreExpr{
	pos: position{line: 236, col: 37, offset: 5403},
	expr: &ruleRefExpr{
	pos: position{line: 236, col: 37, offset: 5403},
	name: "PosixEnvironmentVariableCharacter",
},
},
},
},
},
{
	name: "PosixEnvironmentVariableCharacter",
	pos: position{line: 245, col: 1, offset: 5628},
	expr: &choiceExpr{
	pos: position{line: 246, col: 7, offset: 5672},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 246, col: 7, offset: 5672},
	run: (*parser).callonPosixEnvironmentVariableCharacter2,
	expr: &litMatcher{
	pos: position{line: 246, col: 7, offset: 5672},
	val: "\\\"",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 247, col: 7, offset: 5712},
	run: (*parser).callonPosixEnvironmentVariableCharacter4,
	expr: &litMatcher{
	pos: position{line: 247, col: 7, offset: 5712},
	val: "\\\\",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 248, col: 7, offset: 5752},
	run: (*parser).callonPosixEnvironmentVariableCharacter6,
	expr: &litMatcher{
	pos: position{line: 248, col: 7, offset: 5752},
	val: "\\a",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 249, col: 7, offset: 5792},
	run: (*parser).callonPosixEnvironmentVariableCharacter8,
	expr: &litMatcher{
	pos: position{line: 249, col: 7, offset: 5792},
	val: "\\b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 250, col: 7, offset: 5832},
	run: (*parser).callonPosixEnvironmentVariableCharacter10,
	expr: &litMatcher{
	pos: position{line: 250, col: 7, offset: 5832},
	val: "\\f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 251, col: 7, offset: 5872},
	run: (*parser).callonPosixEnvironmentVariableCharacter12,
	expr: &litMatcher{
	pos: position{line: 251, col: 7, offset: 5872},
	val: "\\n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 252, col: 7, offset: 5912},
	run: (*parser).callonPosixEnvironmentVariableCharacter14,
	expr: &litMatcher{
	pos: position{line: 252, col: 7, offset: 5912},
	val: "\\r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 253, col: 7, offset: 5952},
	run: (*parser).callonPosixEnvironmentVariableCharacter16,
	expr: &litMatcher{
	pos: position{line: 253, col: 7, offset: 5952},
	val: "\\t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 254, col: 7, offset: 5992},
	run: (*parser).callonPosixEnvironmentVariableCharacter18,
	expr: &litMatcher{
	pos: position{line: 254, col: 7, offset: 5992},
	val: "\\v",
	ignoreCase: false,
},
},
&charClassMatcher{
	pos: position{line: 255, col: 7, offset: 6032},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 256, col: 7, offset: 6050},
	val: "[\\x23-\\x3c]",
	ranges: []rune{'#','<',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 257, col: 7, offset: 6068},
	val: "[\\x3e-\\x5b]",
	ranges: []rune{'>','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 258, col: 7, offset: 6086},
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
	pos: position{line: 260, col: 1, offset: 6099},
	expr: &ruleRefExpr{
	pos: position{line: 260, col: 14, offset: 6114},
	name: "Env",
},
},
{
	name: "ImportHashed",
	pos: position{line: 262, col: 1, offset: 6119},
	expr: &actionExpr{
	pos: position{line: 262, col: 16, offset: 6136},
	run: (*parser).callonImportHashed1,
	expr: &labeledExpr{
	pos: position{line: 262, col: 16, offset: 6136},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 262, col: 18, offset: 6138},
	name: "ImportType",
},
},
},
},
{
	name: "Import",
	pos: position{line: 264, col: 1, offset: 6207},
	expr: &choiceExpr{
	pos: position{line: 264, col: 10, offset: 6218},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 264, col: 10, offset: 6218},
	run: (*parser).callonImport2,
	expr: &seqExpr{
	pos: position{line: 264, col: 10, offset: 6218},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 264, col: 10, offset: 6218},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 264, col: 12, offset: 6220},
	name: "ImportHashed",
},
},
&ruleRefExpr{
	pos: position{line: 264, col: 25, offset: 6233},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 264, col: 28, offset: 6236},
	name: "Text",
},
	},
},
},
&actionExpr{
	pos: position{line: 265, col: 10, offset: 6333},
	run: (*parser).callonImport8,
	expr: &labeledExpr{
	pos: position{line: 265, col: 10, offset: 6333},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 265, col: 12, offset: 6335},
	name: "ImportHashed",
},
},
},
	},
},
},
{
	name: "LetBinding",
	pos: position{line: 271, col: 1, offset: 6615},
	expr: &actionExpr{
	pos: position{line: 271, col: 14, offset: 6630},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 271, col: 14, offset: 6630},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 271, col: 14, offset: 6630},
	name: "Let",
},
&labeledExpr{
	pos: position{line: 271, col: 18, offset: 6634},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 271, col: 24, offset: 6640},
	name: "NonreservedLabel",
},
},
&labeledExpr{
	pos: position{line: 271, col: 41, offset: 6657},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 271, col: 43, offset: 6659},
	expr: &ruleRefExpr{
	pos: position{line: 271, col: 43, offset: 6659},
	name: "Annotation",
},
},
},
&ruleRefExpr{
	pos: position{line: 271, col: 55, offset: 6671},
	name: "Equal",
},
&labeledExpr{
	pos: position{line: 271, col: 61, offset: 6677},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 271, col: 63, offset: 6679},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 286, col: 1, offset: 6978},
	expr: &choiceExpr{
	pos: position{line: 287, col: 7, offset: 6999},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 287, col: 7, offset: 6999},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 287, col: 7, offset: 6999},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 287, col: 7, offset: 6999},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 287, col: 14, offset: 7006},
	name: "OpenParens",
},
&labeledExpr{
	pos: position{line: 287, col: 25, offset: 7017},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 287, col: 31, offset: 7023},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 287, col: 48, offset: 7040},
	name: "Colon",
},
&labeledExpr{
	pos: position{line: 287, col: 54, offset: 7046},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 287, col: 56, offset: 7048},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 287, col: 67, offset: 7059},
	name: "CloseParens",
},
&ruleRefExpr{
	pos: position{line: 287, col: 79, offset: 7071},
	name: "Arrow",
},
&labeledExpr{
	pos: position{line: 287, col: 85, offset: 7077},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 287, col: 90, offset: 7082},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 290, col: 7, offset: 7199},
	run: (*parser).callonExpression15,
	expr: &seqExpr{
	pos: position{line: 290, col: 7, offset: 7199},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 290, col: 7, offset: 7199},
	name: "If",
},
&labeledExpr{
	pos: position{line: 290, col: 10, offset: 7202},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 290, col: 15, offset: 7207},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 290, col: 26, offset: 7218},
	name: "Then",
},
&labeledExpr{
	pos: position{line: 290, col: 31, offset: 7223},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 290, col: 33, offset: 7225},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 290, col: 44, offset: 7236},
	name: "Else",
},
&labeledExpr{
	pos: position{line: 290, col: 49, offset: 7241},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 290, col: 51, offset: 7243},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 293, col: 7, offset: 7329},
	run: (*parser).callonExpression26,
	expr: &seqExpr{
	pos: position{line: 293, col: 7, offset: 7329},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 293, col: 7, offset: 7329},
	label: "bindings",
	expr: &zeroOrMoreExpr{
	pos: position{line: 293, col: 16, offset: 7338},
	expr: &ruleRefExpr{
	pos: position{line: 293, col: 16, offset: 7338},
	name: "LetBinding",
},
},
},
&ruleRefExpr{
	pos: position{line: 293, col: 28, offset: 7350},
	name: "In",
},
&labeledExpr{
	pos: position{line: 293, col: 31, offset: 7353},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 293, col: 33, offset: 7355},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 300, col: 7, offset: 7595},
	run: (*parser).callonExpression34,
	expr: &seqExpr{
	pos: position{line: 300, col: 7, offset: 7595},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 300, col: 7, offset: 7595},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 300, col: 14, offset: 7602},
	name: "OpenParens",
},
&labeledExpr{
	pos: position{line: 300, col: 25, offset: 7613},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 300, col: 31, offset: 7619},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 300, col: 48, offset: 7636},
	name: "Colon",
},
&labeledExpr{
	pos: position{line: 300, col: 54, offset: 7642},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 300, col: 56, offset: 7644},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 300, col: 67, offset: 7655},
	name: "CloseParens",
},
&ruleRefExpr{
	pos: position{line: 300, col: 79, offset: 7667},
	name: "Arrow",
},
&labeledExpr{
	pos: position{line: 300, col: 85, offset: 7673},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 300, col: 90, offset: 7678},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 303, col: 7, offset: 7787},
	run: (*parser).callonExpression47,
	expr: &seqExpr{
	pos: position{line: 303, col: 7, offset: 7787},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 303, col: 7, offset: 7787},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 303, col: 9, offset: 7789},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 303, col: 28, offset: 7808},
	name: "Arrow",
},
&labeledExpr{
	pos: position{line: 303, col: 34, offset: 7814},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 303, col: 36, offset: 7816},
	name: "Expression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 304, col: 7, offset: 7876},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 306, col: 1, offset: 7897},
	expr: &actionExpr{
	pos: position{line: 306, col: 14, offset: 7912},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 306, col: 14, offset: 7912},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 306, col: 14, offset: 7912},
	name: "Colon",
},
&labeledExpr{
	pos: position{line: 306, col: 20, offset: 7918},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 306, col: 22, offset: 7920},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 308, col: 1, offset: 7950},
	expr: &choiceExpr{
	pos: position{line: 309, col: 5, offset: 7978},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 309, col: 5, offset: 7978},
	name: "EmptyList",
},
&actionExpr{
	pos: position{line: 310, col: 5, offset: 7992},
	run: (*parser).callonAnnotatedExpression3,
	expr: &seqExpr{
	pos: position{line: 310, col: 5, offset: 7992},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 310, col: 5, offset: 7992},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 310, col: 7, offset: 7994},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 310, col: 26, offset: 8013},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 310, col: 28, offset: 8015},
	expr: &ruleRefExpr{
	pos: position{line: 310, col: 28, offset: 8015},
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
	pos: position{line: 315, col: 1, offset: 8120},
	expr: &actionExpr{
	pos: position{line: 315, col: 13, offset: 8134},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 315, col: 13, offset: 8134},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 315, col: 13, offset: 8134},
	name: "OpenBracket",
},
&ruleRefExpr{
	pos: position{line: 315, col: 25, offset: 8146},
	name: "CloseBracket",
},
&ruleRefExpr{
	pos: position{line: 315, col: 38, offset: 8159},
	name: "Colon",
},
&ruleRefExpr{
	pos: position{line: 315, col: 44, offset: 8165},
	name: "List",
},
&labeledExpr{
	pos: position{line: 315, col: 49, offset: 8170},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 315, col: 51, offset: 8172},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 319, col: 1, offset: 8235},
	expr: &ruleRefExpr{
	pos: position{line: 319, col: 22, offset: 8258},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 321, col: 1, offset: 8279},
	expr: &ruleRefExpr{
	pos: position{line: 321, col: 23, offset: 8303},
	name: "PlusExpression",
},
},
{
	name: "MorePlus",
	pos: position{line: 323, col: 1, offset: 8319},
	expr: &actionExpr{
	pos: position{line: 323, col: 12, offset: 8332},
	run: (*parser).callonMorePlus1,
	expr: &seqExpr{
	pos: position{line: 323, col: 12, offset: 8332},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 323, col: 12, offset: 8332},
	name: "Plus",
},
&labeledExpr{
	pos: position{line: 323, col: 17, offset: 8337},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 323, col: 19, offset: 8339},
	name: "TimesExpression",
},
},
	},
},
},
},
{
	name: "PlusExpression",
	pos: position{line: 324, col: 1, offset: 8373},
	expr: &actionExpr{
	pos: position{line: 325, col: 7, offset: 8398},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 325, col: 7, offset: 8398},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 325, col: 7, offset: 8398},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 325, col: 13, offset: 8404},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 325, col: 29, offset: 8420},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 325, col: 34, offset: 8425},
	expr: &ruleRefExpr{
	pos: position{line: 325, col: 34, offset: 8425},
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
	pos: position{line: 334, col: 1, offset: 8653},
	expr: &actionExpr{
	pos: position{line: 334, col: 13, offset: 8667},
	run: (*parser).callonMoreTimes1,
	expr: &seqExpr{
	pos: position{line: 334, col: 13, offset: 8667},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 334, col: 13, offset: 8667},
	name: "Times",
},
&labeledExpr{
	pos: position{line: 334, col: 19, offset: 8673},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 334, col: 21, offset: 8675},
	name: "ApplicationExpression",
},
},
	},
},
},
},
{
	name: "TimesExpression",
	pos: position{line: 335, col: 1, offset: 8715},
	expr: &actionExpr{
	pos: position{line: 336, col: 7, offset: 8741},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 336, col: 7, offset: 8741},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 336, col: 7, offset: 8741},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 336, col: 13, offset: 8747},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 336, col: 35, offset: 8769},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 336, col: 40, offset: 8774},
	expr: &ruleRefExpr{
	pos: position{line: 336, col: 40, offset: 8774},
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
	pos: position{line: 345, col: 1, offset: 9004},
	expr: &actionExpr{
	pos: position{line: 345, col: 25, offset: 9030},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 345, col: 25, offset: 9030},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 345, col: 25, offset: 9030},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 345, col: 27, offset: 9032},
	name: "ImportExpression",
},
},
&labeledExpr{
	pos: position{line: 345, col: 44, offset: 9049},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 345, col: 49, offset: 9054},
	expr: &ruleRefExpr{
	pos: position{line: 345, col: 49, offset: 9054},
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
	pos: position{line: 354, col: 1, offset: 9284},
	expr: &choiceExpr{
	pos: position{line: 354, col: 20, offset: 9305},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 354, col: 20, offset: 9305},
	name: "Import",
},
&ruleRefExpr{
	pos: position{line: 354, col: 29, offset: 9314},
	name: "SelectorExpression",
},
	},
},
},
{
	name: "SelectorExpression",
	pos: position{line: 356, col: 1, offset: 9334},
	expr: &actionExpr{
	pos: position{line: 356, col: 22, offset: 9357},
	run: (*parser).callonSelectorExpression1,
	expr: &seqExpr{
	pos: position{line: 356, col: 22, offset: 9357},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 356, col: 22, offset: 9357},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 356, col: 24, offset: 9359},
	name: "PrimitiveExpression",
},
},
&labeledExpr{
	pos: position{line: 356, col: 44, offset: 9379},
	label: "ls",
	expr: &zeroOrMoreExpr{
	pos: position{line: 356, col: 47, offset: 9382},
	expr: &seqExpr{
	pos: position{line: 356, col: 48, offset: 9383},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 356, col: 48, offset: 9383},
	name: "Dot",
},
&ruleRefExpr{
	pos: position{line: 356, col: 52, offset: 9387},
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
	pos: position{line: 366, col: 1, offset: 9617},
	expr: &choiceExpr{
	pos: position{line: 367, col: 7, offset: 9647},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 367, col: 7, offset: 9647},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 368, col: 7, offset: 9667},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 369, col: 7, offset: 9688},
	name: "IntegerLiteral",
},
&actionExpr{
	pos: position{line: 370, col: 7, offset: 9709},
	run: (*parser).callonPrimitiveExpression5,
	expr: &litMatcher{
	pos: position{line: 370, col: 7, offset: 9709},
	val: "-Infinity",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 371, col: 7, offset: 9767},
	name: "TextLiteral",
},
&actionExpr{
	pos: position{line: 372, col: 7, offset: 9785},
	run: (*parser).callonPrimitiveExpression8,
	expr: &seqExpr{
	pos: position{line: 372, col: 7, offset: 9785},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 372, col: 7, offset: 9785},
	name: "OpenBrace",
},
&labeledExpr{
	pos: position{line: 372, col: 17, offset: 9795},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 372, col: 19, offset: 9797},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 372, col: 39, offset: 9817},
	name: "CloseBrace",
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 373, col: 7, offset: 9852},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 374, col: 7, offset: 9878},
	name: "IdentifierReservedPrefix",
},
&ruleRefExpr{
	pos: position{line: 375, col: 7, offset: 9909},
	name: "Reserved",
},
&ruleRefExpr{
	pos: position{line: 376, col: 7, offset: 9924},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 377, col: 7, offset: 9941},
	run: (*parser).callonPrimitiveExpression18,
	expr: &seqExpr{
	pos: position{line: 377, col: 7, offset: 9941},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 377, col: 7, offset: 9941},
	name: "OpenParens",
},
&labeledExpr{
	pos: position{line: 377, col: 18, offset: 9952},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 377, col: 20, offset: 9954},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 377, col: 31, offset: 9965},
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
	pos: position{line: 379, col: 1, offset: 9996},
	expr: &choiceExpr{
	pos: position{line: 380, col: 7, offset: 10026},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 380, col: 7, offset: 10026},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &ruleRefExpr{
	pos: position{line: 380, col: 7, offset: 10026},
	name: "Equal",
},
},
&ruleRefExpr{
	pos: position{line: 381, col: 7, offset: 10083},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 382, col: 7, offset: 10108},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 383, col: 7, offset: 10136},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 383, col: 7, offset: 10136},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 385, col: 1, offset: 10182},
	expr: &actionExpr{
	pos: position{line: 385, col: 19, offset: 10202},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 385, col: 19, offset: 10202},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 385, col: 19, offset: 10202},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 385, col: 24, offset: 10207},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 385, col: 30, offset: 10213},
	name: "Colon",
},
&labeledExpr{
	pos: position{line: 385, col: 36, offset: 10219},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 385, col: 41, offset: 10224},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 388, col: 1, offset: 10281},
	expr: &actionExpr{
	pos: position{line: 388, col: 18, offset: 10300},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 388, col: 18, offset: 10300},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 388, col: 18, offset: 10300},
	name: "Comma",
},
&labeledExpr{
	pos: position{line: 388, col: 24, offset: 10306},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 388, col: 26, offset: 10308},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 389, col: 1, offset: 10340},
	expr: &actionExpr{
	pos: position{line: 390, col: 7, offset: 10369},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 390, col: 7, offset: 10369},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 390, col: 7, offset: 10369},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 390, col: 13, offset: 10375},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 390, col: 29, offset: 10391},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 390, col: 34, offset: 10396},
	expr: &ruleRefExpr{
	pos: position{line: 390, col: 34, offset: 10396},
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
	pos: position{line: 400, col: 1, offset: 10792},
	expr: &actionExpr{
	pos: position{line: 400, col: 22, offset: 10815},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 400, col: 22, offset: 10815},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 400, col: 22, offset: 10815},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 400, col: 27, offset: 10820},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 400, col: 33, offset: 10826},
	name: "Equal",
},
&labeledExpr{
	pos: position{line: 400, col: 39, offset: 10832},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 400, col: 44, offset: 10837},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 403, col: 1, offset: 10894},
	expr: &actionExpr{
	pos: position{line: 403, col: 21, offset: 10916},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 403, col: 21, offset: 10916},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 403, col: 21, offset: 10916},
	name: "Comma",
},
&labeledExpr{
	pos: position{line: 403, col: 27, offset: 10922},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 403, col: 29, offset: 10924},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 404, col: 1, offset: 10959},
	expr: &actionExpr{
	pos: position{line: 405, col: 7, offset: 10991},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 405, col: 7, offset: 10991},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 405, col: 7, offset: 10991},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 405, col: 13, offset: 10997},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 405, col: 32, offset: 11016},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 405, col: 37, offset: 11021},
	expr: &ruleRefExpr{
	pos: position{line: 405, col: 37, offset: 11021},
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
	pos: position{line: 415, col: 1, offset: 11423},
	expr: &actionExpr{
	pos: position{line: 415, col: 12, offset: 11436},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 415, col: 12, offset: 11436},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 415, col: 12, offset: 11436},
	name: "Comma",
},
&labeledExpr{
	pos: position{line: 415, col: 18, offset: 11442},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 415, col: 20, offset: 11444},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 417, col: 1, offset: 11472},
	expr: &actionExpr{
	pos: position{line: 418, col: 7, offset: 11502},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 418, col: 7, offset: 11502},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 418, col: 7, offset: 11502},
	name: "OpenBracket",
},
&labeledExpr{
	pos: position{line: 418, col: 19, offset: 11514},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 418, col: 25, offset: 11520},
	name: "Expression",
},
},
&labeledExpr{
	pos: position{line: 418, col: 36, offset: 11531},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 418, col: 41, offset: 11536},
	expr: &ruleRefExpr{
	pos: position{line: 418, col: 41, offset: 11536},
	name: "MoreList",
},
},
},
&ruleRefExpr{
	pos: position{line: 418, col: 51, offset: 11546},
	name: "CloseBracket",
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 428, col: 1, offset: 11831},
	expr: &notExpr{
	pos: position{line: 428, col: 7, offset: 11839},
	expr: &anyMatcher{
	line: 428, col: 8, offset: 11840,
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

func (c *current) onNonreservedLabel1(label interface{}) (interface{}, error) {
 return label, nil 
}

func (p *parser) callonNonreservedLabel1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNonreservedLabel1(stack["label"])
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
  return ImportType{EnvVar: string(c.text)}, nil
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
  return ImportType{EnvVar: b.String()}, nil
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
 return ImportHashed{ImportType: i.(ImportType)}, nil 
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

func (c *current) onImport8(i interface{}) (interface{}, error) {
 return Embed(Import{ImportHashed: i.(ImportHashed), ImportMode: Code}), nil 
}

func (p *parser) callonImport8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onImport8(stack["i"])
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

