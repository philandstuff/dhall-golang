
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
&ruleRefExpr{
	pos: position{line: 23, col: 37, offset: 278},
	name: "_",
},
	},
},
},
},
{
	name: "EOL",
	pos: position{line: 25, col: 1, offset: 299},
	expr: &choiceExpr{
	pos: position{line: 25, col: 7, offset: 307},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 25, col: 7, offset: 307},
	val: "\n",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 25, col: 14, offset: 314},
	val: "\r\n",
	ignoreCase: false,
},
	},
},
},
{
	name: "BlockComment",
	pos: position{line: 27, col: 1, offset: 322},
	expr: &seqExpr{
	pos: position{line: 27, col: 16, offset: 339},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 27, col: 16, offset: 339},
	val: "{-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 27, col: 21, offset: 344},
	name: "BlockCommentContinue",
},
	},
},
},
{
	name: "BlockCommentChunk",
	pos: position{line: 29, col: 1, offset: 366},
	expr: &choiceExpr{
	pos: position{line: 30, col: 5, offset: 392},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 30, col: 5, offset: 392},
	name: "BlockComment",
},
&charClassMatcher{
	pos: position{line: 31, col: 5, offset: 409},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 32, col: 5, offset: 435},
	name: "EOL",
},
	},
},
},
{
	name: "BlockCommentContinue",
	pos: position{line: 34, col: 1, offset: 440},
	expr: &choiceExpr{
	pos: position{line: 34, col: 24, offset: 465},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 34, col: 24, offset: 465},
	val: "-}",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 34, col: 31, offset: 472},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 34, col: 31, offset: 472},
	name: "BlockCommentChunk",
},
&ruleRefExpr{
	pos: position{line: 34, col: 49, offset: 490},
	name: "BlockCommentContinue",
},
	},
},
	},
},
},
{
	name: "NotEOL",
	pos: position{line: 36, col: 1, offset: 512},
	expr: &charClassMatcher{
	pos: position{line: 36, col: 10, offset: 523},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "LineComment",
	pos: position{line: 38, col: 1, offset: 546},
	expr: &actionExpr{
	pos: position{line: 38, col: 15, offset: 562},
	run: (*parser).callonLineComment1,
	expr: &seqExpr{
	pos: position{line: 38, col: 15, offset: 562},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 38, col: 15, offset: 562},
	val: "--",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 38, col: 20, offset: 567},
	label: "content",
	expr: &actionExpr{
	pos: position{line: 38, col: 29, offset: 576},
	run: (*parser).callonLineComment5,
	expr: &zeroOrMoreExpr{
	pos: position{line: 38, col: 29, offset: 576},
	expr: &ruleRefExpr{
	pos: position{line: 38, col: 29, offset: 576},
	name: "NotEOL",
},
},
},
},
&ruleRefExpr{
	pos: position{line: 38, col: 68, offset: 615},
	name: "EOL",
},
	},
},
},
},
{
	name: "WhitespaceChunk",
	pos: position{line: 40, col: 1, offset: 644},
	expr: &choiceExpr{
	pos: position{line: 40, col: 19, offset: 664},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 40, col: 19, offset: 664},
	val: " ",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 40, col: 25, offset: 670},
	val: "\t",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 40, col: 32, offset: 677},
	name: "EOL",
},
&ruleRefExpr{
	pos: position{line: 40, col: 38, offset: 683},
	name: "LineComment",
},
&ruleRefExpr{
	pos: position{line: 40, col: 52, offset: 697},
	name: "BlockComment",
},
	},
},
},
{
	name: "_",
	pos: position{line: 42, col: 1, offset: 711},
	expr: &zeroOrMoreExpr{
	pos: position{line: 42, col: 5, offset: 717},
	expr: &ruleRefExpr{
	pos: position{line: 42, col: 5, offset: 717},
	name: "WhitespaceChunk",
},
},
},
{
	name: "_1",
	pos: position{line: 44, col: 1, offset: 735},
	expr: &oneOrMoreExpr{
	pos: position{line: 44, col: 6, offset: 742},
	expr: &ruleRefExpr{
	pos: position{line: 44, col: 6, offset: 742},
	name: "WhitespaceChunk",
},
},
},
{
	name: "HexDig",
	pos: position{line: 46, col: 1, offset: 760},
	expr: &charClassMatcher{
	pos: position{line: 46, col: 10, offset: 771},
	val: "[0-9a-f]i",
	ranges: []rune{'0','9','a','f',},
	ignoreCase: true,
	inverted: false,
},
},
{
	name: "SimpleLabelFirstChar",
	pos: position{line: 48, col: 1, offset: 782},
	expr: &charClassMatcher{
	pos: position{line: 48, col: 24, offset: 807},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SimpleLabelNextChar",
	pos: position{line: 49, col: 1, offset: 817},
	expr: &charClassMatcher{
	pos: position{line: 49, col: 23, offset: 841},
	val: "[A-Za-z0-9_/-]",
	chars: []rune{'_','/','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SimpleLabel",
	pos: position{line: 50, col: 1, offset: 856},
	expr: &choiceExpr{
	pos: position{line: 50, col: 15, offset: 872},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 50, col: 15, offset: 872},
	run: (*parser).callonSimpleLabel2,
	expr: &seqExpr{
	pos: position{line: 50, col: 15, offset: 872},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 50, col: 15, offset: 872},
	name: "Keyword",
},
&oneOrMoreExpr{
	pos: position{line: 50, col: 23, offset: 880},
	expr: &ruleRefExpr{
	pos: position{line: 50, col: 23, offset: 880},
	name: "SimpleLabelNextChar",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 51, col: 13, offset: 944},
	run: (*parser).callonSimpleLabel7,
	expr: &seqExpr{
	pos: position{line: 51, col: 13, offset: 944},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 51, col: 13, offset: 944},
	expr: &ruleRefExpr{
	pos: position{line: 51, col: 14, offset: 945},
	name: "Keyword",
},
},
&ruleRefExpr{
	pos: position{line: 51, col: 22, offset: 953},
	name: "SimpleLabelFirstChar",
},
&zeroOrMoreExpr{
	pos: position{line: 51, col: 43, offset: 974},
	expr: &ruleRefExpr{
	pos: position{line: 51, col: 43, offset: 974},
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
	pos: position{line: 58, col: 1, offset: 1075},
	expr: &actionExpr{
	pos: position{line: 58, col: 9, offset: 1085},
	run: (*parser).callonLabel1,
	expr: &labeledExpr{
	pos: position{line: 58, col: 9, offset: 1085},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 58, col: 15, offset: 1091},
	name: "SimpleLabel",
},
},
},
},
{
	name: "NonreservedLabel",
	pos: position{line: 60, col: 1, offset: 1126},
	expr: &choiceExpr{
	pos: position{line: 60, col: 20, offset: 1147},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 60, col: 20, offset: 1147},
	run: (*parser).callonNonreservedLabel2,
	expr: &seqExpr{
	pos: position{line: 60, col: 20, offset: 1147},
	exprs: []interface{}{
&andExpr{
	pos: position{line: 60, col: 20, offset: 1147},
	expr: &seqExpr{
	pos: position{line: 60, col: 22, offset: 1149},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 60, col: 22, offset: 1149},
	name: "Reserved",
},
&ruleRefExpr{
	pos: position{line: 60, col: 31, offset: 1158},
	name: "SimpleLabelNextChar",
},
	},
},
},
&labeledExpr{
	pos: position{line: 60, col: 52, offset: 1179},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 60, col: 58, offset: 1185},
	name: "Label",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 61, col: 19, offset: 1231},
	run: (*parser).callonNonreservedLabel10,
	expr: &seqExpr{
	pos: position{line: 61, col: 19, offset: 1231},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 61, col: 19, offset: 1231},
	expr: &ruleRefExpr{
	pos: position{line: 61, col: 20, offset: 1232},
	name: "Reserved",
},
},
&labeledExpr{
	pos: position{line: 61, col: 29, offset: 1241},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 61, col: 35, offset: 1247},
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
	pos: position{line: 63, col: 1, offset: 1276},
	expr: &ruleRefExpr{
	pos: position{line: 63, col: 12, offset: 1289},
	name: "Label",
},
},
{
	name: "EscapedChar",
	pos: position{line: 67, col: 1, offset: 1329},
	expr: &actionExpr{
	pos: position{line: 68, col: 3, offset: 1347},
	run: (*parser).callonEscapedChar1,
	expr: &seqExpr{
	pos: position{line: 68, col: 3, offset: 1347},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 68, col: 3, offset: 1347},
	val: "\\",
	ignoreCase: false,
},
&choiceExpr{
	pos: position{line: 69, col: 5, offset: 1356},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 69, col: 5, offset: 1356},
	val: "\"",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 70, col: 10, offset: 1369},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 71, col: 10, offset: 1382},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 72, col: 10, offset: 1396},
	val: "/",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 73, col: 10, offset: 1409},
	val: "b",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 74, col: 10, offset: 1422},
	val: "f",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 75, col: 10, offset: 1435},
	val: "n",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 76, col: 10, offset: 1448},
	val: "r",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 77, col: 10, offset: 1461},
	val: "t",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 78, col: 10, offset: 1474},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 78, col: 10, offset: 1474},
	val: "u",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 78, col: 14, offset: 1478},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 78, col: 21, offset: 1485},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 78, col: 28, offset: 1492},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 78, col: 35, offset: 1499},
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
	pos: position{line: 99, col: 1, offset: 1942},
	expr: &choiceExpr{
	pos: position{line: 100, col: 6, offset: 1968},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 100, col: 6, offset: 1968},
	run: (*parser).callonDoubleQuoteChunk2,
	expr: &seqExpr{
	pos: position{line: 100, col: 6, offset: 1968},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 100, col: 6, offset: 1968},
	val: "${",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 100, col: 11, offset: 1973},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 100, col: 13, offset: 1975},
	name: "CompleteExpression",
},
},
&litMatcher{
	pos: position{line: 100, col: 32, offset: 1994},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 101, col: 6, offset: 2021},
	name: "EscapedChar",
},
&charClassMatcher{
	pos: position{line: 102, col: 6, offset: 2038},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 103, col: 6, offset: 2055},
	val: "[\\x23-\\x5b]",
	ranges: []rune{'#','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 104, col: 6, offset: 2072},
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
	pos: position{line: 106, col: 1, offset: 2091},
	expr: &actionExpr{
	pos: position{line: 106, col: 22, offset: 2114},
	run: (*parser).callonDoubleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 106, col: 22, offset: 2114},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 106, col: 22, offset: 2114},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 106, col: 26, offset: 2118},
	label: "chunks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 106, col: 33, offset: 2125},
	expr: &ruleRefExpr{
	pos: position{line: 106, col: 33, offset: 2125},
	name: "DoubleQuoteChunk",
},
},
},
&litMatcher{
	pos: position{line: 106, col: 51, offset: 2143},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TextLiteral",
	pos: position{line: 123, col: 1, offset: 2611},
	expr: &actionExpr{
	pos: position{line: 123, col: 15, offset: 2627},
	run: (*parser).callonTextLiteral1,
	expr: &labeledExpr{
	pos: position{line: 123, col: 15, offset: 2627},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 123, col: 17, offset: 2629},
	name: "DoubleQuoteLiteral",
},
},
},
},
{
	name: "Reserved",
	pos: position{line: 126, col: 1, offset: 2752},
	expr: &choiceExpr{
	pos: position{line: 127, col: 5, offset: 2769},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 127, col: 5, offset: 2769},
	run: (*parser).callonReserved2,
	expr: &litMatcher{
	pos: position{line: 127, col: 5, offset: 2769},
	val: "Natural/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 128, col: 5, offset: 2847},
	run: (*parser).callonReserved4,
	expr: &litMatcher{
	pos: position{line: 128, col: 5, offset: 2847},
	val: "Natural/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 129, col: 5, offset: 2923},
	run: (*parser).callonReserved6,
	expr: &litMatcher{
	pos: position{line: 129, col: 5, offset: 2923},
	val: "Natural/isZero",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 130, col: 5, offset: 3003},
	run: (*parser).callonReserved8,
	expr: &litMatcher{
	pos: position{line: 130, col: 5, offset: 3003},
	val: "Natural/even",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 131, col: 5, offset: 3079},
	run: (*parser).callonReserved10,
	expr: &litMatcher{
	pos: position{line: 131, col: 5, offset: 3079},
	val: "Natural/odd",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 132, col: 5, offset: 3153},
	run: (*parser).callonReserved12,
	expr: &litMatcher{
	pos: position{line: 132, col: 5, offset: 3153},
	val: "Natural/toInteger",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 133, col: 5, offset: 3239},
	run: (*parser).callonReserved14,
	expr: &litMatcher{
	pos: position{line: 133, col: 5, offset: 3239},
	val: "Natural/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 134, col: 5, offset: 3315},
	run: (*parser).callonReserved16,
	expr: &litMatcher{
	pos: position{line: 134, col: 5, offset: 3315},
	val: "Integer/toDouble",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 135, col: 5, offset: 3399},
	run: (*parser).callonReserved18,
	expr: &litMatcher{
	pos: position{line: 135, col: 5, offset: 3399},
	val: "Integer/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 136, col: 5, offset: 3475},
	run: (*parser).callonReserved20,
	expr: &litMatcher{
	pos: position{line: 136, col: 5, offset: 3475},
	val: "Double/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 137, col: 5, offset: 3549},
	run: (*parser).callonReserved22,
	expr: &litMatcher{
	pos: position{line: 137, col: 5, offset: 3549},
	val: "List/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 138, col: 5, offset: 3621},
	run: (*parser).callonReserved24,
	expr: &litMatcher{
	pos: position{line: 138, col: 5, offset: 3621},
	val: "List/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 139, col: 5, offset: 3691},
	run: (*parser).callonReserved26,
	expr: &litMatcher{
	pos: position{line: 139, col: 5, offset: 3691},
	val: "List/length",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 140, col: 5, offset: 3765},
	run: (*parser).callonReserved28,
	expr: &litMatcher{
	pos: position{line: 140, col: 5, offset: 3765},
	val: "List/head",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 141, col: 5, offset: 3835},
	run: (*parser).callonReserved30,
	expr: &litMatcher{
	pos: position{line: 141, col: 5, offset: 3835},
	val: "List/last",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 142, col: 5, offset: 3905},
	run: (*parser).callonReserved32,
	expr: &litMatcher{
	pos: position{line: 142, col: 5, offset: 3905},
	val: "List/indexed",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 143, col: 5, offset: 3981},
	run: (*parser).callonReserved34,
	expr: &litMatcher{
	pos: position{line: 143, col: 5, offset: 3981},
	val: "List/reverse",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 144, col: 5, offset: 4057},
	run: (*parser).callonReserved36,
	expr: &litMatcher{
	pos: position{line: 144, col: 5, offset: 4057},
	val: "Optional/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 145, col: 5, offset: 4137},
	run: (*parser).callonReserved38,
	expr: &litMatcher{
	pos: position{line: 145, col: 5, offset: 4137},
	val: "Optional/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 146, col: 5, offset: 4215},
	run: (*parser).callonReserved40,
	expr: &litMatcher{
	pos: position{line: 146, col: 5, offset: 4215},
	val: "Text/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 147, col: 5, offset: 4285},
	run: (*parser).callonReserved42,
	expr: &litMatcher{
	pos: position{line: 147, col: 5, offset: 4285},
	val: "Bool",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 148, col: 5, offset: 4317},
	run: (*parser).callonReserved44,
	expr: &litMatcher{
	pos: position{line: 148, col: 5, offset: 4317},
	val: "True",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 149, col: 5, offset: 4349},
	run: (*parser).callonReserved46,
	expr: &litMatcher{
	pos: position{line: 149, col: 5, offset: 4349},
	val: "False",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 150, col: 5, offset: 4383},
	run: (*parser).callonReserved48,
	expr: &litMatcher{
	pos: position{line: 150, col: 5, offset: 4383},
	val: "Optional",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 151, col: 5, offset: 4423},
	run: (*parser).callonReserved50,
	expr: &litMatcher{
	pos: position{line: 151, col: 5, offset: 4423},
	val: "Natural",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 152, col: 5, offset: 4461},
	run: (*parser).callonReserved52,
	expr: &litMatcher{
	pos: position{line: 152, col: 5, offset: 4461},
	val: "Integer",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 153, col: 5, offset: 4499},
	run: (*parser).callonReserved54,
	expr: &litMatcher{
	pos: position{line: 153, col: 5, offset: 4499},
	val: "Double",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 154, col: 5, offset: 4535},
	run: (*parser).callonReserved56,
	expr: &litMatcher{
	pos: position{line: 154, col: 5, offset: 4535},
	val: "Text",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 155, col: 5, offset: 4567},
	run: (*parser).callonReserved58,
	expr: &litMatcher{
	pos: position{line: 155, col: 5, offset: 4567},
	val: "List",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 156, col: 5, offset: 4599},
	run: (*parser).callonReserved60,
	expr: &litMatcher{
	pos: position{line: 156, col: 5, offset: 4599},
	val: "None",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 157, col: 5, offset: 4631},
	run: (*parser).callonReserved62,
	expr: &litMatcher{
	pos: position{line: 157, col: 5, offset: 4631},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 158, col: 5, offset: 4663},
	run: (*parser).callonReserved64,
	expr: &litMatcher{
	pos: position{line: 158, col: 5, offset: 4663},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 159, col: 5, offset: 4695},
	run: (*parser).callonReserved66,
	expr: &litMatcher{
	pos: position{line: 159, col: 5, offset: 4695},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "If",
	pos: position{line: 161, col: 1, offset: 4724},
	expr: &litMatcher{
	pos: position{line: 161, col: 6, offset: 4731},
	val: "if",
	ignoreCase: false,
},
},
{
	name: "Then",
	pos: position{line: 162, col: 1, offset: 4736},
	expr: &litMatcher{
	pos: position{line: 162, col: 8, offset: 4745},
	val: "then",
	ignoreCase: false,
},
},
{
	name: "Else",
	pos: position{line: 163, col: 1, offset: 4752},
	expr: &litMatcher{
	pos: position{line: 163, col: 8, offset: 4761},
	val: "else",
	ignoreCase: false,
},
},
{
	name: "Let",
	pos: position{line: 164, col: 1, offset: 4768},
	expr: &litMatcher{
	pos: position{line: 164, col: 7, offset: 4776},
	val: "let",
	ignoreCase: false,
},
},
{
	name: "In",
	pos: position{line: 165, col: 1, offset: 4782},
	expr: &litMatcher{
	pos: position{line: 165, col: 6, offset: 4789},
	val: "in",
	ignoreCase: false,
},
},
{
	name: "As",
	pos: position{line: 166, col: 1, offset: 4794},
	expr: &litMatcher{
	pos: position{line: 166, col: 6, offset: 4801},
	val: "as",
	ignoreCase: false,
},
},
{
	name: "Using",
	pos: position{line: 167, col: 1, offset: 4806},
	expr: &litMatcher{
	pos: position{line: 167, col: 9, offset: 4816},
	val: "using",
	ignoreCase: false,
},
},
{
	name: "Merge",
	pos: position{line: 168, col: 1, offset: 4824},
	expr: &litMatcher{
	pos: position{line: 168, col: 9, offset: 4834},
	val: "merge",
	ignoreCase: false,
},
},
{
	name: "Missing",
	pos: position{line: 169, col: 1, offset: 4842},
	expr: &litMatcher{
	pos: position{line: 169, col: 11, offset: 4854},
	val: "missing",
	ignoreCase: false,
},
},
{
	name: "True",
	pos: position{line: 170, col: 1, offset: 4864},
	expr: &litMatcher{
	pos: position{line: 170, col: 8, offset: 4873},
	val: "True",
	ignoreCase: false,
},
},
{
	name: "False",
	pos: position{line: 171, col: 1, offset: 4880},
	expr: &litMatcher{
	pos: position{line: 171, col: 9, offset: 4890},
	val: "False",
	ignoreCase: false,
},
},
{
	name: "Infinity",
	pos: position{line: 172, col: 1, offset: 4898},
	expr: &litMatcher{
	pos: position{line: 172, col: 12, offset: 4911},
	val: "Infinity",
	ignoreCase: false,
},
},
{
	name: "NaN",
	pos: position{line: 173, col: 1, offset: 4922},
	expr: &litMatcher{
	pos: position{line: 173, col: 7, offset: 4930},
	val: "NaN",
	ignoreCase: false,
},
},
{
	name: "Some",
	pos: position{line: 174, col: 1, offset: 4936},
	expr: &litMatcher{
	pos: position{line: 174, col: 8, offset: 4945},
	val: "Some",
	ignoreCase: false,
},
},
{
	name: "Keyword",
	pos: position{line: 176, col: 1, offset: 4953},
	expr: &choiceExpr{
	pos: position{line: 177, col: 5, offset: 4969},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 177, col: 5, offset: 4969},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 177, col: 10, offset: 4974},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 177, col: 17, offset: 4981},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 178, col: 5, offset: 4990},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 178, col: 11, offset: 4996},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 179, col: 5, offset: 5003},
	name: "Using",
},
&ruleRefExpr{
	pos: position{line: 179, col: 13, offset: 5011},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 179, col: 23, offset: 5021},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 180, col: 5, offset: 5028},
	name: "True",
},
&ruleRefExpr{
	pos: position{line: 180, col: 12, offset: 5035},
	name: "False",
},
&ruleRefExpr{
	pos: position{line: 181, col: 5, offset: 5045},
	name: "Infinity",
},
&ruleRefExpr{
	pos: position{line: 181, col: 16, offset: 5056},
	name: "NaN",
},
&ruleRefExpr{
	pos: position{line: 182, col: 5, offset: 5064},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 182, col: 13, offset: 5072},
	name: "Some",
},
	},
},
},
{
	name: "Optional",
	pos: position{line: 184, col: 1, offset: 5078},
	expr: &litMatcher{
	pos: position{line: 184, col: 12, offset: 5091},
	val: "Optional",
	ignoreCase: false,
},
},
{
	name: "Text",
	pos: position{line: 185, col: 1, offset: 5102},
	expr: &litMatcher{
	pos: position{line: 185, col: 8, offset: 5111},
	val: "Text",
	ignoreCase: false,
},
},
{
	name: "List",
	pos: position{line: 186, col: 1, offset: 5118},
	expr: &litMatcher{
	pos: position{line: 186, col: 8, offset: 5127},
	val: "List",
	ignoreCase: false,
},
},
{
	name: "Lambda",
	pos: position{line: 188, col: 1, offset: 5135},
	expr: &choiceExpr{
	pos: position{line: 188, col: 11, offset: 5147},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 188, col: 11, offset: 5147},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 188, col: 18, offset: 5154},
	val: "λ",
	ignoreCase: false,
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 189, col: 1, offset: 5160},
	expr: &choiceExpr{
	pos: position{line: 189, col: 11, offset: 5172},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 189, col: 11, offset: 5172},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 189, col: 22, offset: 5183},
	val: "∀",
	ignoreCase: false,
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 190, col: 1, offset: 5190},
	expr: &choiceExpr{
	pos: position{line: 190, col: 10, offset: 5201},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 190, col: 10, offset: 5201},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 190, col: 17, offset: 5208},
	val: "→",
	ignoreCase: false,
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 192, col: 1, offset: 5216},
	expr: &seqExpr{
	pos: position{line: 192, col: 12, offset: 5229},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 192, col: 12, offset: 5229},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 192, col: 17, offset: 5234},
	expr: &charClassMatcher{
	pos: position{line: 192, col: 17, offset: 5234},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 192, col: 23, offset: 5240},
	expr: &charClassMatcher{
	pos: position{line: 192, col: 23, offset: 5240},
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
	pos: position{line: 194, col: 1, offset: 5248},
	expr: &actionExpr{
	pos: position{line: 194, col: 24, offset: 5273},
	run: (*parser).callonNumericDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 194, col: 24, offset: 5273},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 194, col: 24, offset: 5273},
	expr: &charClassMatcher{
	pos: position{line: 194, col: 24, offset: 5273},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 194, col: 30, offset: 5279},
	expr: &charClassMatcher{
	pos: position{line: 194, col: 30, offset: 5279},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&choiceExpr{
	pos: position{line: 194, col: 39, offset: 5288},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 194, col: 39, offset: 5288},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 194, col: 39, offset: 5288},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 194, col: 43, offset: 5292},
	expr: &charClassMatcher{
	pos: position{line: 194, col: 43, offset: 5292},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&zeroOrOneExpr{
	pos: position{line: 194, col: 50, offset: 5299},
	expr: &ruleRefExpr{
	pos: position{line: 194, col: 50, offset: 5299},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 194, col: 62, offset: 5311},
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
	pos: position{line: 202, col: 1, offset: 5467},
	expr: &choiceExpr{
	pos: position{line: 202, col: 17, offset: 5485},
	alternatives: []interface{}{
&labeledExpr{
	pos: position{line: 202, col: 17, offset: 5485},
	label: "d",
	expr: &ruleRefExpr{
	pos: position{line: 202, col: 19, offset: 5487},
	name: "NumericDoubleLiteral",
},
},
&actionExpr{
	pos: position{line: 203, col: 5, offset: 5512},
	run: (*parser).callonDoubleLiteral4,
	expr: &ruleRefExpr{
	pos: position{line: 203, col: 5, offset: 5512},
	name: "Infinity",
},
},
&actionExpr{
	pos: position{line: 204, col: 5, offset: 5564},
	run: (*parser).callonDoubleLiteral6,
	expr: &seqExpr{
	pos: position{line: 204, col: 5, offset: 5564},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 204, col: 5, offset: 5564},
	val: "-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 204, col: 9, offset: 5568},
	name: "Infinity",
},
	},
},
},
&actionExpr{
	pos: position{line: 205, col: 5, offset: 5621},
	run: (*parser).callonDoubleLiteral10,
	expr: &ruleRefExpr{
	pos: position{line: 205, col: 5, offset: 5621},
	name: "NaN",
},
},
	},
},
},
{
	name: "NaturalLiteral",
	pos: position{line: 207, col: 1, offset: 5664},
	expr: &actionExpr{
	pos: position{line: 207, col: 18, offset: 5683},
	run: (*parser).callonNaturalLiteral1,
	expr: &oneOrMoreExpr{
	pos: position{line: 207, col: 18, offset: 5683},
	expr: &charClassMatcher{
	pos: position{line: 207, col: 18, offset: 5683},
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
	pos: position{line: 212, col: 1, offset: 5772},
	expr: &actionExpr{
	pos: position{line: 212, col: 18, offset: 5791},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 212, col: 18, offset: 5791},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 212, col: 18, offset: 5791},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&oneOrMoreExpr{
	pos: position{line: 212, col: 22, offset: 5795},
	expr: &charClassMatcher{
	pos: position{line: 212, col: 22, offset: 5795},
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
	pos: position{line: 220, col: 1, offset: 5939},
	expr: &actionExpr{
	pos: position{line: 220, col: 12, offset: 5952},
	run: (*parser).callonDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 220, col: 12, offset: 5952},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 220, col: 12, offset: 5952},
	name: "_",
},
&litMatcher{
	pos: position{line: 220, col: 14, offset: 5954},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 220, col: 18, offset: 5958},
	name: "_",
},
&labeledExpr{
	pos: position{line: 220, col: 20, offset: 5960},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 220, col: 26, offset: 5966},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Variable",
	pos: position{line: 222, col: 1, offset: 6022},
	expr: &actionExpr{
	pos: position{line: 222, col: 12, offset: 6035},
	run: (*parser).callonVariable1,
	expr: &seqExpr{
	pos: position{line: 222, col: 12, offset: 6035},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 222, col: 12, offset: 6035},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 222, col: 17, offset: 6040},
	name: "NonreservedLabel",
},
},
&labeledExpr{
	pos: position{line: 222, col: 34, offset: 6057},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 222, col: 40, offset: 6063},
	expr: &ruleRefExpr{
	pos: position{line: 222, col: 40, offset: 6063},
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
	pos: position{line: 230, col: 1, offset: 6226},
	expr: &choiceExpr{
	pos: position{line: 230, col: 14, offset: 6241},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 230, col: 14, offset: 6241},
	name: "Variable",
},
&ruleRefExpr{
	pos: position{line: 230, col: 25, offset: 6252},
	name: "Reserved",
},
	},
},
},
{
	name: "Env",
	pos: position{line: 232, col: 1, offset: 6262},
	expr: &actionExpr{
	pos: position{line: 232, col: 7, offset: 6270},
	run: (*parser).callonEnv1,
	expr: &seqExpr{
	pos: position{line: 232, col: 7, offset: 6270},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 232, col: 7, offset: 6270},
	val: "env:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 232, col: 14, offset: 6277},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 232, col: 17, offset: 6280},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 232, col: 17, offset: 6280},
	name: "BashEnvironmentVariable",
},
&ruleRefExpr{
	pos: position{line: 232, col: 43, offset: 6306},
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
	pos: position{line: 234, col: 1, offset: 6351},
	expr: &actionExpr{
	pos: position{line: 234, col: 27, offset: 6379},
	run: (*parser).callonBashEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 234, col: 27, offset: 6379},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 234, col: 27, offset: 6379},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 234, col: 36, offset: 6388},
	expr: &charClassMatcher{
	pos: position{line: 234, col: 36, offset: 6388},
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
	pos: position{line: 238, col: 1, offset: 6456},
	expr: &actionExpr{
	pos: position{line: 238, col: 28, offset: 6485},
	run: (*parser).callonPosixEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 238, col: 28, offset: 6485},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 238, col: 28, offset: 6485},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 238, col: 32, offset: 6489},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 238, col: 34, offset: 6491},
	name: "PosixEnvironmentVariableContent",
},
},
&litMatcher{
	pos: position{line: 238, col: 66, offset: 6523},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PosixEnvironmentVariableContent",
	pos: position{line: 242, col: 1, offset: 6548},
	expr: &actionExpr{
	pos: position{line: 242, col: 35, offset: 6584},
	run: (*parser).callonPosixEnvironmentVariableContent1,
	expr: &labeledExpr{
	pos: position{line: 242, col: 35, offset: 6584},
	label: "v",
	expr: &oneOrMoreExpr{
	pos: position{line: 242, col: 37, offset: 6586},
	expr: &ruleRefExpr{
	pos: position{line: 242, col: 37, offset: 6586},
	name: "PosixEnvironmentVariableCharacter",
},
},
},
},
},
{
	name: "PosixEnvironmentVariableCharacter",
	pos: position{line: 251, col: 1, offset: 6811},
	expr: &choiceExpr{
	pos: position{line: 252, col: 7, offset: 6855},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 252, col: 7, offset: 6855},
	run: (*parser).callonPosixEnvironmentVariableCharacter2,
	expr: &litMatcher{
	pos: position{line: 252, col: 7, offset: 6855},
	val: "\\\"",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 253, col: 7, offset: 6895},
	run: (*parser).callonPosixEnvironmentVariableCharacter4,
	expr: &litMatcher{
	pos: position{line: 253, col: 7, offset: 6895},
	val: "\\\\",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 254, col: 7, offset: 6935},
	run: (*parser).callonPosixEnvironmentVariableCharacter6,
	expr: &litMatcher{
	pos: position{line: 254, col: 7, offset: 6935},
	val: "\\a",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 255, col: 7, offset: 6975},
	run: (*parser).callonPosixEnvironmentVariableCharacter8,
	expr: &litMatcher{
	pos: position{line: 255, col: 7, offset: 6975},
	val: "\\b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 256, col: 7, offset: 7015},
	run: (*parser).callonPosixEnvironmentVariableCharacter10,
	expr: &litMatcher{
	pos: position{line: 256, col: 7, offset: 7015},
	val: "\\f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 257, col: 7, offset: 7055},
	run: (*parser).callonPosixEnvironmentVariableCharacter12,
	expr: &litMatcher{
	pos: position{line: 257, col: 7, offset: 7055},
	val: "\\n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 258, col: 7, offset: 7095},
	run: (*parser).callonPosixEnvironmentVariableCharacter14,
	expr: &litMatcher{
	pos: position{line: 258, col: 7, offset: 7095},
	val: "\\r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 259, col: 7, offset: 7135},
	run: (*parser).callonPosixEnvironmentVariableCharacter16,
	expr: &litMatcher{
	pos: position{line: 259, col: 7, offset: 7135},
	val: "\\t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 260, col: 7, offset: 7175},
	run: (*parser).callonPosixEnvironmentVariableCharacter18,
	expr: &litMatcher{
	pos: position{line: 260, col: 7, offset: 7175},
	val: "\\v",
	ignoreCase: false,
},
},
&charClassMatcher{
	pos: position{line: 261, col: 7, offset: 7215},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 262, col: 7, offset: 7233},
	val: "[\\x23-\\x3c]",
	ranges: []rune{'#','<',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 263, col: 7, offset: 7251},
	val: "[\\x3e-\\x5b]",
	ranges: []rune{'>','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 264, col: 7, offset: 7269},
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
	pos: position{line: 266, col: 1, offset: 7282},
	expr: &ruleRefExpr{
	pos: position{line: 266, col: 14, offset: 7297},
	name: "Env",
},
},
{
	name: "ImportHashed",
	pos: position{line: 268, col: 1, offset: 7302},
	expr: &actionExpr{
	pos: position{line: 268, col: 16, offset: 7319},
	run: (*parser).callonImportHashed1,
	expr: &labeledExpr{
	pos: position{line: 268, col: 16, offset: 7319},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 268, col: 18, offset: 7321},
	name: "ImportType",
},
},
},
},
{
	name: "Import",
	pos: position{line: 270, col: 1, offset: 7390},
	expr: &choiceExpr{
	pos: position{line: 270, col: 10, offset: 7401},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 270, col: 10, offset: 7401},
	run: (*parser).callonImport2,
	expr: &seqExpr{
	pos: position{line: 270, col: 10, offset: 7401},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 270, col: 10, offset: 7401},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 270, col: 12, offset: 7403},
	name: "ImportHashed",
},
},
&ruleRefExpr{
	pos: position{line: 270, col: 25, offset: 7416},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 270, col: 27, offset: 7418},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 270, col: 30, offset: 7421},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 270, col: 33, offset: 7424},
	name: "Text",
},
	},
},
},
&actionExpr{
	pos: position{line: 271, col: 10, offset: 7521},
	run: (*parser).callonImport10,
	expr: &labeledExpr{
	pos: position{line: 271, col: 10, offset: 7521},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 271, col: 12, offset: 7523},
	name: "ImportHashed",
},
},
},
	},
},
},
{
	name: "LetBinding",
	pos: position{line: 274, col: 1, offset: 7618},
	expr: &actionExpr{
	pos: position{line: 274, col: 14, offset: 7633},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 274, col: 14, offset: 7633},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 274, col: 14, offset: 7633},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 274, col: 18, offset: 7637},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 274, col: 21, offset: 7640},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 274, col: 27, offset: 7646},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 274, col: 44, offset: 7663},
	name: "_",
},
&labeledExpr{
	pos: position{line: 274, col: 46, offset: 7665},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 274, col: 48, offset: 7667},
	expr: &seqExpr{
	pos: position{line: 274, col: 49, offset: 7668},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 274, col: 49, offset: 7668},
	name: "Annotation",
},
&ruleRefExpr{
	pos: position{line: 274, col: 60, offset: 7679},
	name: "_",
},
	},
},
},
},
&litMatcher{
	pos: position{line: 275, col: 13, offset: 7695},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 275, col: 17, offset: 7699},
	name: "_",
},
&labeledExpr{
	pos: position{line: 275, col: 19, offset: 7701},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 275, col: 21, offset: 7703},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 275, col: 32, offset: 7714},
	name: "_",
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 290, col: 1, offset: 8023},
	expr: &choiceExpr{
	pos: position{line: 291, col: 7, offset: 8044},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 291, col: 7, offset: 8044},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 291, col: 7, offset: 8044},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 291, col: 7, offset: 8044},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 291, col: 14, offset: 8051},
	name: "_",
},
&litMatcher{
	pos: position{line: 291, col: 16, offset: 8053},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 291, col: 20, offset: 8057},
	name: "_",
},
&labeledExpr{
	pos: position{line: 291, col: 22, offset: 8059},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 291, col: 28, offset: 8065},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 291, col: 45, offset: 8082},
	name: "_",
},
&litMatcher{
	pos: position{line: 291, col: 47, offset: 8084},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 291, col: 51, offset: 8088},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 291, col: 54, offset: 8091},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 291, col: 56, offset: 8093},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 291, col: 67, offset: 8104},
	name: "_",
},
&litMatcher{
	pos: position{line: 291, col: 69, offset: 8106},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 291, col: 73, offset: 8110},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 291, col: 75, offset: 8112},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 291, col: 81, offset: 8118},
	name: "_",
},
&labeledExpr{
	pos: position{line: 291, col: 83, offset: 8120},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 291, col: 88, offset: 8125},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 294, col: 7, offset: 8242},
	run: (*parser).callonExpression22,
	expr: &seqExpr{
	pos: position{line: 294, col: 7, offset: 8242},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 294, col: 7, offset: 8242},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 294, col: 10, offset: 8245},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 294, col: 13, offset: 8248},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 294, col: 18, offset: 8253},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 294, col: 29, offset: 8264},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 294, col: 31, offset: 8266},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 294, col: 36, offset: 8271},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 294, col: 39, offset: 8274},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 294, col: 41, offset: 8276},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 294, col: 52, offset: 8287},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 294, col: 54, offset: 8289},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 294, col: 59, offset: 8294},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 294, col: 62, offset: 8297},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 294, col: 64, offset: 8299},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 297, col: 7, offset: 8385},
	run: (*parser).callonExpression38,
	expr: &seqExpr{
	pos: position{line: 297, col: 7, offset: 8385},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 297, col: 7, offset: 8385},
	label: "bindings",
	expr: &oneOrMoreExpr{
	pos: position{line: 297, col: 16, offset: 8394},
	expr: &ruleRefExpr{
	pos: position{line: 297, col: 16, offset: 8394},
	name: "LetBinding",
},
},
},
&ruleRefExpr{
	pos: position{line: 297, col: 28, offset: 8406},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 297, col: 31, offset: 8409},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 297, col: 34, offset: 8412},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 297, col: 36, offset: 8414},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 304, col: 7, offset: 8654},
	run: (*parser).callonExpression47,
	expr: &seqExpr{
	pos: position{line: 304, col: 7, offset: 8654},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 304, col: 7, offset: 8654},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 304, col: 14, offset: 8661},
	name: "_",
},
&litMatcher{
	pos: position{line: 304, col: 16, offset: 8663},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 304, col: 20, offset: 8667},
	name: "_",
},
&labeledExpr{
	pos: position{line: 304, col: 22, offset: 8669},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 304, col: 28, offset: 8675},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 304, col: 45, offset: 8692},
	name: "_",
},
&litMatcher{
	pos: position{line: 304, col: 47, offset: 8694},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 304, col: 51, offset: 8698},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 304, col: 54, offset: 8701},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 304, col: 56, offset: 8703},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 304, col: 67, offset: 8714},
	name: "_",
},
&litMatcher{
	pos: position{line: 304, col: 69, offset: 8716},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 304, col: 73, offset: 8720},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 304, col: 75, offset: 8722},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 304, col: 81, offset: 8728},
	name: "_",
},
&labeledExpr{
	pos: position{line: 304, col: 83, offset: 8730},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 304, col: 88, offset: 8735},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 307, col: 7, offset: 8844},
	run: (*parser).callonExpression67,
	expr: &seqExpr{
	pos: position{line: 307, col: 7, offset: 8844},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 307, col: 7, offset: 8844},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 307, col: 9, offset: 8846},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 307, col: 28, offset: 8865},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 307, col: 30, offset: 8867},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 307, col: 36, offset: 8873},
	name: "_",
},
&labeledExpr{
	pos: position{line: 307, col: 38, offset: 8875},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 307, col: 40, offset: 8877},
	name: "Expression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 308, col: 7, offset: 8937},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 310, col: 1, offset: 8958},
	expr: &actionExpr{
	pos: position{line: 310, col: 14, offset: 8973},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 310, col: 14, offset: 8973},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 310, col: 14, offset: 8973},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 310, col: 18, offset: 8977},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 310, col: 21, offset: 8980},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 310, col: 23, offset: 8982},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 312, col: 1, offset: 9012},
	expr: &choiceExpr{
	pos: position{line: 313, col: 5, offset: 9040},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 313, col: 5, offset: 9040},
	name: "EmptyList",
},
&actionExpr{
	pos: position{line: 314, col: 5, offset: 9054},
	run: (*parser).callonAnnotatedExpression3,
	expr: &seqExpr{
	pos: position{line: 314, col: 5, offset: 9054},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 314, col: 5, offset: 9054},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 314, col: 7, offset: 9056},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 314, col: 26, offset: 9075},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 314, col: 28, offset: 9077},
	expr: &seqExpr{
	pos: position{line: 314, col: 29, offset: 9078},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 314, col: 29, offset: 9078},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 314, col: 31, offset: 9080},
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
	pos: position{line: 319, col: 1, offset: 9205},
	expr: &actionExpr{
	pos: position{line: 319, col: 13, offset: 9219},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 319, col: 13, offset: 9219},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 319, col: 13, offset: 9219},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 319, col: 17, offset: 9223},
	name: "_",
},
&litMatcher{
	pos: position{line: 319, col: 19, offset: 9225},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 319, col: 23, offset: 9229},
	name: "_",
},
&litMatcher{
	pos: position{line: 319, col: 25, offset: 9231},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 319, col: 29, offset: 9235},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 319, col: 32, offset: 9238},
	name: "List",
},
&ruleRefExpr{
	pos: position{line: 319, col: 37, offset: 9243},
	name: "_",
},
&labeledExpr{
	pos: position{line: 319, col: 39, offset: 9245},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 319, col: 41, offset: 9247},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 323, col: 1, offset: 9310},
	expr: &ruleRefExpr{
	pos: position{line: 323, col: 22, offset: 9333},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 325, col: 1, offset: 9354},
	expr: &ruleRefExpr{
	pos: position{line: 325, col: 23, offset: 9378},
	name: "PlusExpression",
},
},
{
	name: "MorePlus",
	pos: position{line: 327, col: 1, offset: 9394},
	expr: &actionExpr{
	pos: position{line: 327, col: 12, offset: 9407},
	run: (*parser).callonMorePlus1,
	expr: &seqExpr{
	pos: position{line: 327, col: 12, offset: 9407},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 327, col: 12, offset: 9407},
	name: "_",
},
&litMatcher{
	pos: position{line: 327, col: 14, offset: 9409},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 327, col: 18, offset: 9413},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 327, col: 21, offset: 9416},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 327, col: 23, offset: 9418},
	name: "TimesExpression",
},
},
	},
},
},
},
{
	name: "PlusExpression",
	pos: position{line: 328, col: 1, offset: 9452},
	expr: &actionExpr{
	pos: position{line: 329, col: 7, offset: 9477},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 329, col: 7, offset: 9477},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 329, col: 7, offset: 9477},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 329, col: 13, offset: 9483},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 329, col: 29, offset: 9499},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 329, col: 34, offset: 9504},
	expr: &ruleRefExpr{
	pos: position{line: 329, col: 34, offset: 9504},
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
	pos: position{line: 338, col: 1, offset: 9732},
	expr: &actionExpr{
	pos: position{line: 338, col: 13, offset: 9746},
	run: (*parser).callonMoreTimes1,
	expr: &seqExpr{
	pos: position{line: 338, col: 13, offset: 9746},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 338, col: 13, offset: 9746},
	name: "_",
},
&litMatcher{
	pos: position{line: 338, col: 15, offset: 9748},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 338, col: 19, offset: 9752},
	name: "_",
},
&labeledExpr{
	pos: position{line: 338, col: 21, offset: 9754},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 338, col: 23, offset: 9756},
	name: "ApplicationExpression",
},
},
	},
},
},
},
{
	name: "TimesExpression",
	pos: position{line: 339, col: 1, offset: 9796},
	expr: &actionExpr{
	pos: position{line: 340, col: 7, offset: 9822},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 340, col: 7, offset: 9822},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 340, col: 7, offset: 9822},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 340, col: 13, offset: 9828},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 340, col: 35, offset: 9850},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 340, col: 40, offset: 9855},
	expr: &ruleRefExpr{
	pos: position{line: 340, col: 40, offset: 9855},
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
	pos: position{line: 349, col: 1, offset: 10085},
	expr: &actionExpr{
	pos: position{line: 349, col: 25, offset: 10111},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 349, col: 25, offset: 10111},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 349, col: 25, offset: 10111},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 349, col: 27, offset: 10113},
	name: "FirstApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 349, col: 54, offset: 10140},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 349, col: 59, offset: 10145},
	expr: &seqExpr{
	pos: position{line: 349, col: 60, offset: 10146},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 349, col: 60, offset: 10146},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 349, col: 63, offset: 10149},
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
	pos: position{line: 358, col: 1, offset: 10399},
	expr: &choiceExpr{
	pos: position{line: 359, col: 8, offset: 10437},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 359, col: 8, offset: 10437},
	run: (*parser).callonFirstApplicationExpression2,
	expr: &seqExpr{
	pos: position{line: 359, col: 8, offset: 10437},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 359, col: 8, offset: 10437},
	name: "Some",
},
&ruleRefExpr{
	pos: position{line: 359, col: 13, offset: 10442},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 359, col: 16, offset: 10445},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 359, col: 18, offset: 10447},
	name: "ImportExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 360, col: 8, offset: 10502},
	name: "ImportExpression",
},
	},
},
},
{
	name: "ImportExpression",
	pos: position{line: 362, col: 1, offset: 10520},
	expr: &choiceExpr{
	pos: position{line: 362, col: 20, offset: 10541},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 362, col: 20, offset: 10541},
	name: "Import",
},
&ruleRefExpr{
	pos: position{line: 362, col: 29, offset: 10550},
	name: "SelectorExpression",
},
	},
},
},
{
	name: "SelectorExpression",
	pos: position{line: 364, col: 1, offset: 10570},
	expr: &actionExpr{
	pos: position{line: 364, col: 22, offset: 10593},
	run: (*parser).callonSelectorExpression1,
	expr: &seqExpr{
	pos: position{line: 364, col: 22, offset: 10593},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 364, col: 22, offset: 10593},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 364, col: 24, offset: 10595},
	name: "PrimitiveExpression",
},
},
&labeledExpr{
	pos: position{line: 364, col: 44, offset: 10615},
	label: "ls",
	expr: &zeroOrMoreExpr{
	pos: position{line: 364, col: 47, offset: 10618},
	expr: &seqExpr{
	pos: position{line: 364, col: 48, offset: 10619},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 364, col: 48, offset: 10619},
	name: "_",
},
&litMatcher{
	pos: position{line: 364, col: 50, offset: 10621},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 364, col: 54, offset: 10625},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 364, col: 56, offset: 10627},
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
	pos: position{line: 374, col: 1, offset: 10860},
	expr: &choiceExpr{
	pos: position{line: 375, col: 7, offset: 10890},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 375, col: 7, offset: 10890},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 376, col: 7, offset: 10910},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 377, col: 7, offset: 10931},
	name: "IntegerLiteral",
},
&ruleRefExpr{
	pos: position{line: 378, col: 7, offset: 10952},
	name: "TextLiteral",
},
&actionExpr{
	pos: position{line: 379, col: 7, offset: 10970},
	run: (*parser).callonPrimitiveExpression6,
	expr: &seqExpr{
	pos: position{line: 379, col: 7, offset: 10970},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 379, col: 7, offset: 10970},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 379, col: 11, offset: 10974},
	name: "_",
},
&labeledExpr{
	pos: position{line: 379, col: 13, offset: 10976},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 379, col: 15, offset: 10978},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 379, col: 35, offset: 10998},
	name: "_",
},
&litMatcher{
	pos: position{line: 379, col: 37, offset: 11000},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 380, col: 7, offset: 11028},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 381, col: 7, offset: 11054},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 382, col: 7, offset: 11071},
	run: (*parser).callonPrimitiveExpression16,
	expr: &seqExpr{
	pos: position{line: 382, col: 7, offset: 11071},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 382, col: 7, offset: 11071},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 382, col: 11, offset: 11075},
	name: "_",
},
&labeledExpr{
	pos: position{line: 382, col: 14, offset: 11078},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 382, col: 16, offset: 11080},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 382, col: 27, offset: 11091},
	name: "_",
},
&litMatcher{
	pos: position{line: 382, col: 29, offset: 11093},
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
	pos: position{line: 384, col: 1, offset: 11116},
	expr: &choiceExpr{
	pos: position{line: 385, col: 7, offset: 11146},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 385, col: 7, offset: 11146},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &litMatcher{
	pos: position{line: 385, col: 7, offset: 11146},
	val: "=",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 386, col: 7, offset: 11201},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 387, col: 7, offset: 11226},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 388, col: 7, offset: 11254},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 388, col: 7, offset: 11254},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 390, col: 1, offset: 11300},
	expr: &actionExpr{
	pos: position{line: 390, col: 19, offset: 11320},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 390, col: 19, offset: 11320},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 390, col: 19, offset: 11320},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 390, col: 24, offset: 11325},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 390, col: 33, offset: 11334},
	name: "_",
},
&litMatcher{
	pos: position{line: 390, col: 35, offset: 11336},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 390, col: 39, offset: 11340},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 390, col: 42, offset: 11343},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 390, col: 47, offset: 11348},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 393, col: 1, offset: 11405},
	expr: &actionExpr{
	pos: position{line: 393, col: 18, offset: 11424},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 393, col: 18, offset: 11424},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 393, col: 18, offset: 11424},
	name: "_",
},
&litMatcher{
	pos: position{line: 393, col: 20, offset: 11426},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 393, col: 24, offset: 11430},
	name: "_",
},
&labeledExpr{
	pos: position{line: 393, col: 26, offset: 11432},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 393, col: 28, offset: 11434},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 394, col: 1, offset: 11466},
	expr: &actionExpr{
	pos: position{line: 395, col: 7, offset: 11495},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 395, col: 7, offset: 11495},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 395, col: 7, offset: 11495},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 395, col: 13, offset: 11501},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 395, col: 29, offset: 11517},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 395, col: 34, offset: 11522},
	expr: &ruleRefExpr{
	pos: position{line: 395, col: 34, offset: 11522},
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
	pos: position{line: 405, col: 1, offset: 11918},
	expr: &actionExpr{
	pos: position{line: 405, col: 22, offset: 11941},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 405, col: 22, offset: 11941},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 405, col: 22, offset: 11941},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 405, col: 27, offset: 11946},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 405, col: 36, offset: 11955},
	name: "_",
},
&litMatcher{
	pos: position{line: 405, col: 38, offset: 11957},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 405, col: 42, offset: 11961},
	name: "_",
},
&labeledExpr{
	pos: position{line: 405, col: 44, offset: 11963},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 405, col: 49, offset: 11968},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 408, col: 1, offset: 12025},
	expr: &actionExpr{
	pos: position{line: 408, col: 21, offset: 12047},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 408, col: 21, offset: 12047},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 408, col: 21, offset: 12047},
	name: "_",
},
&litMatcher{
	pos: position{line: 408, col: 23, offset: 12049},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 408, col: 27, offset: 12053},
	name: "_",
},
&labeledExpr{
	pos: position{line: 408, col: 29, offset: 12055},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 408, col: 31, offset: 12057},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 409, col: 1, offset: 12092},
	expr: &actionExpr{
	pos: position{line: 410, col: 7, offset: 12124},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 410, col: 7, offset: 12124},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 410, col: 7, offset: 12124},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 410, col: 13, offset: 12130},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 410, col: 32, offset: 12149},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 410, col: 37, offset: 12154},
	expr: &ruleRefExpr{
	pos: position{line: 410, col: 37, offset: 12154},
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
	pos: position{line: 420, col: 1, offset: 12556},
	expr: &actionExpr{
	pos: position{line: 420, col: 12, offset: 12569},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 420, col: 12, offset: 12569},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 420, col: 12, offset: 12569},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 420, col: 16, offset: 12573},
	name: "_",
},
&labeledExpr{
	pos: position{line: 420, col: 18, offset: 12575},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 420, col: 20, offset: 12577},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 420, col: 31, offset: 12588},
	name: "_",
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 422, col: 1, offset: 12607},
	expr: &actionExpr{
	pos: position{line: 423, col: 7, offset: 12637},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 423, col: 7, offset: 12637},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 423, col: 7, offset: 12637},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 423, col: 11, offset: 12641},
	name: "_",
},
&labeledExpr{
	pos: position{line: 423, col: 13, offset: 12643},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 423, col: 19, offset: 12649},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 423, col: 30, offset: 12660},
	name: "_",
},
&labeledExpr{
	pos: position{line: 423, col: 32, offset: 12662},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 423, col: 37, offset: 12667},
	expr: &ruleRefExpr{
	pos: position{line: 423, col: 37, offset: 12667},
	name: "MoreList",
},
},
},
&litMatcher{
	pos: position{line: 423, col: 47, offset: 12677},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 433, col: 1, offset: 12953},
	expr: &notExpr{
	pos: position{line: 433, col: 7, offset: 12961},
	expr: &anyMatcher{
	line: 433, col: 8, offset: 12962,
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

