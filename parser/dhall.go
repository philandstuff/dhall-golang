
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
	name: "CompleteExpression",
	pos: position{line: 21, col: 1, offset: 180},
	expr: &actionExpr{
	pos: position{line: 21, col: 22, offset: 203},
	run: (*parser).callonCompleteExpression1,
	expr: &seqExpr{
	pos: position{line: 21, col: 22, offset: 203},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 21, col: 22, offset: 203},
	name: "_",
},
&labeledExpr{
	pos: position{line: 21, col: 24, offset: 205},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 21, col: 26, offset: 207},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 21, col: 37, offset: 218},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 21, col: 39, offset: 220},
	name: "EOF",
},
	},
},
},
},
{
	name: "EOL",
	pos: position{line: 23, col: 1, offset: 243},
	expr: &choiceExpr{
	pos: position{line: 23, col: 7, offset: 251},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 23, col: 7, offset: 251},
	val: "\n",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 23, col: 14, offset: 258},
	val: "\r\n",
	ignoreCase: false,
},
	},
},
},
{
	name: "BlockComment",
	pos: position{line: 25, col: 1, offset: 266},
	expr: &seqExpr{
	pos: position{line: 25, col: 16, offset: 283},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 25, col: 16, offset: 283},
	val: "{-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 25, col: 21, offset: 288},
	name: "BlockCommentContinue",
},
	},
},
},
{
	name: "BlockCommentChunk",
	pos: position{line: 27, col: 1, offset: 310},
	expr: &choiceExpr{
	pos: position{line: 28, col: 5, offset: 336},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 28, col: 5, offset: 336},
	name: "BlockComment",
},
&charClassMatcher{
	pos: position{line: 29, col: 5, offset: 353},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 30, col: 5, offset: 379},
	name: "EOL",
},
	},
},
},
{
	name: "BlockCommentContinue",
	pos: position{line: 32, col: 1, offset: 384},
	expr: &choiceExpr{
	pos: position{line: 32, col: 24, offset: 409},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 32, col: 24, offset: 409},
	val: "-}",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 32, col: 31, offset: 416},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 32, col: 31, offset: 416},
	name: "BlockCommentChunk",
},
&ruleRefExpr{
	pos: position{line: 32, col: 49, offset: 434},
	name: "BlockCommentContinue",
},
	},
},
	},
},
},
{
	name: "NotEOL",
	pos: position{line: 34, col: 1, offset: 456},
	expr: &charClassMatcher{
	pos: position{line: 34, col: 10, offset: 467},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "LineComment",
	pos: position{line: 36, col: 1, offset: 490},
	expr: &actionExpr{
	pos: position{line: 36, col: 15, offset: 506},
	run: (*parser).callonLineComment1,
	expr: &seqExpr{
	pos: position{line: 36, col: 15, offset: 506},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 36, col: 15, offset: 506},
	val: "--",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 36, col: 20, offset: 511},
	label: "content",
	expr: &actionExpr{
	pos: position{line: 36, col: 29, offset: 520},
	run: (*parser).callonLineComment5,
	expr: &zeroOrMoreExpr{
	pos: position{line: 36, col: 29, offset: 520},
	expr: &ruleRefExpr{
	pos: position{line: 36, col: 29, offset: 520},
	name: "NotEOL",
},
},
},
},
&ruleRefExpr{
	pos: position{line: 36, col: 68, offset: 559},
	name: "EOL",
},
	},
},
},
},
{
	name: "WhitespaceChunk",
	pos: position{line: 38, col: 1, offset: 588},
	expr: &choiceExpr{
	pos: position{line: 38, col: 19, offset: 608},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 38, col: 19, offset: 608},
	val: " ",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 38, col: 25, offset: 614},
	val: "\t",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 38, col: 32, offset: 621},
	name: "EOL",
},
&ruleRefExpr{
	pos: position{line: 38, col: 38, offset: 627},
	name: "LineComment",
},
&ruleRefExpr{
	pos: position{line: 38, col: 52, offset: 641},
	name: "BlockComment",
},
	},
},
},
{
	name: "_",
	pos: position{line: 40, col: 1, offset: 655},
	expr: &zeroOrMoreExpr{
	pos: position{line: 40, col: 5, offset: 661},
	expr: &ruleRefExpr{
	pos: position{line: 40, col: 5, offset: 661},
	name: "WhitespaceChunk",
},
},
},
{
	name: "NonemptyWhitespace",
	pos: position{line: 42, col: 1, offset: 679},
	expr: &oneOrMoreExpr{
	pos: position{line: 42, col: 22, offset: 702},
	expr: &ruleRefExpr{
	pos: position{line: 42, col: 22, offset: 702},
	name: "WhitespaceChunk",
},
},
},
{
	name: "SimpleLabel",
	pos: position{line: 44, col: 1, offset: 720},
	expr: &actionExpr{
	pos: position{line: 44, col: 15, offset: 736},
	run: (*parser).callonSimpleLabel1,
	expr: &seqExpr{
	pos: position{line: 44, col: 15, offset: 736},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 44, col: 15, offset: 736},
	expr: &ruleRefExpr{
	pos: position{line: 44, col: 16, offset: 737},
	name: "Keyword",
},
},
&charClassMatcher{
	pos: position{line: 45, col: 13, offset: 757},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 45, col: 23, offset: 767},
	expr: &charClassMatcher{
	pos: position{line: 45, col: 23, offset: 767},
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
	pos: position{line: 49, col: 1, offset: 831},
	expr: &actionExpr{
	pos: position{line: 49, col: 9, offset: 841},
	run: (*parser).callonLabel1,
	expr: &labeledExpr{
	pos: position{line: 49, col: 9, offset: 841},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 49, col: 15, offset: 847},
	name: "SimpleLabel",
},
},
},
},
{
	name: "Reserved",
	pos: position{line: 51, col: 1, offset: 882},
	expr: &choiceExpr{
	pos: position{line: 51, col: 12, offset: 895},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 51, col: 12, offset: 895},
	run: (*parser).callonReserved2,
	expr: &litMatcher{
	pos: position{line: 51, col: 12, offset: 895},
	val: "Bool",
	ignoreCase: false,
},
},
&litMatcher{
	pos: position{line: 52, col: 5, offset: 931},
	val: "Optional",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 53, col: 5, offset: 946},
	val: "None",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 54, col: 5, offset: 957},
	run: (*parser).callonReserved6,
	expr: &litMatcher{
	pos: position{line: 54, col: 5, offset: 957},
	val: "Natural",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 55, col: 5, offset: 999},
	run: (*parser).callonReserved8,
	expr: &litMatcher{
	pos: position{line: 55, col: 5, offset: 999},
	val: "Integer",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 56, col: 5, offset: 1041},
	run: (*parser).callonReserved10,
	expr: &litMatcher{
	pos: position{line: 56, col: 5, offset: 1041},
	val: "Double",
	ignoreCase: false,
},
},
&litMatcher{
	pos: position{line: 57, col: 5, offset: 1081},
	val: "Text",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 58, col: 5, offset: 1092},
	run: (*parser).callonReserved13,
	expr: &litMatcher{
	pos: position{line: 58, col: 5, offset: 1092},
	val: "List",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 59, col: 5, offset: 1128},
	run: (*parser).callonReserved15,
	expr: &litMatcher{
	pos: position{line: 59, col: 5, offset: 1128},
	val: "True",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 60, col: 5, offset: 1164},
	run: (*parser).callonReserved17,
	expr: &litMatcher{
	pos: position{line: 60, col: 5, offset: 1164},
	val: "False",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 61, col: 5, offset: 1202},
	run: (*parser).callonReserved19,
	expr: &litMatcher{
	pos: position{line: 61, col: 5, offset: 1202},
	val: "NaN",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 62, col: 5, offset: 1254},
	run: (*parser).callonReserved21,
	expr: &litMatcher{
	pos: position{line: 62, col: 5, offset: 1254},
	val: "Infinity",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 63, col: 5, offset: 1312},
	run: (*parser).callonReserved23,
	expr: &litMatcher{
	pos: position{line: 63, col: 5, offset: 1312},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 64, col: 5, offset: 1348},
	run: (*parser).callonReserved25,
	expr: &litMatcher{
	pos: position{line: 64, col: 5, offset: 1348},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 65, col: 5, offset: 1384},
	run: (*parser).callonReserved27,
	expr: &litMatcher{
	pos: position{line: 65, col: 5, offset: 1384},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "Keyword",
	pos: position{line: 67, col: 1, offset: 1417},
	expr: &choiceExpr{
	pos: position{line: 67, col: 11, offset: 1429},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 67, col: 11, offset: 1429},
	val: "if",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 68, col: 5, offset: 1438},
	val: "then",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 69, col: 5, offset: 1449},
	val: "else",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 70, col: 5, offset: 1460},
	val: "let",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 71, col: 5, offset: 1470},
	val: "in",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 72, col: 5, offset: 1479},
	val: "as",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 73, col: 5, offset: 1488},
	val: "using",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 74, col: 5, offset: 1500},
	val: "merge",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 75, col: 5, offset: 1512},
	val: "constructors",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 76, col: 5, offset: 1531},
	val: "Some",
	ignoreCase: false,
},
	},
},
},
{
	name: "ColonSpace",
	pos: position{line: 78, col: 1, offset: 1539},
	expr: &seqExpr{
	pos: position{line: 78, col: 14, offset: 1554},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 78, col: 14, offset: 1554},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 78, col: 18, offset: 1558},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Lambda",
	pos: position{line: 80, col: 1, offset: 1578},
	expr: &choiceExpr{
	pos: position{line: 80, col: 11, offset: 1590},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 80, col: 11, offset: 1590},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 80, col: 18, offset: 1597},
	val: "λ",
	ignoreCase: false,
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 81, col: 1, offset: 1603},
	expr: &choiceExpr{
	pos: position{line: 81, col: 11, offset: 1615},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 81, col: 11, offset: 1615},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 81, col: 22, offset: 1626},
	val: "∀",
	ignoreCase: false,
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 82, col: 1, offset: 1633},
	expr: &choiceExpr{
	pos: position{line: 82, col: 10, offset: 1644},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 82, col: 10, offset: 1644},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 82, col: 17, offset: 1651},
	val: "→",
	ignoreCase: false,
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 84, col: 1, offset: 1659},
	expr: &seqExpr{
	pos: position{line: 84, col: 12, offset: 1672},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 84, col: 12, offset: 1672},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 84, col: 17, offset: 1677},
	expr: &charClassMatcher{
	pos: position{line: 84, col: 17, offset: 1677},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 84, col: 23, offset: 1683},
	expr: &charClassMatcher{
	pos: position{line: 84, col: 23, offset: 1683},
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
	pos: position{line: 86, col: 1, offset: 1691},
	expr: &actionExpr{
	pos: position{line: 86, col: 17, offset: 1709},
	run: (*parser).callonDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 86, col: 17, offset: 1709},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 86, col: 17, offset: 1709},
	expr: &charClassMatcher{
	pos: position{line: 86, col: 17, offset: 1709},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 86, col: 23, offset: 1715},
	expr: &charClassMatcher{
	pos: position{line: 86, col: 23, offset: 1715},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&choiceExpr{
	pos: position{line: 86, col: 32, offset: 1724},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 86, col: 32, offset: 1724},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 86, col: 32, offset: 1724},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 86, col: 36, offset: 1728},
	expr: &charClassMatcher{
	pos: position{line: 86, col: 36, offset: 1728},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&zeroOrOneExpr{
	pos: position{line: 86, col: 43, offset: 1735},
	expr: &ruleRefExpr{
	pos: position{line: 86, col: 43, offset: 1735},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 86, col: 55, offset: 1747},
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
	pos: position{line: 94, col: 1, offset: 1907},
	expr: &actionExpr{
	pos: position{line: 94, col: 18, offset: 1926},
	run: (*parser).callonNaturalLiteral1,
	expr: &oneOrMoreExpr{
	pos: position{line: 94, col: 18, offset: 1926},
	expr: &charClassMatcher{
	pos: position{line: 94, col: 18, offset: 1926},
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
	pos: position{line: 102, col: 1, offset: 2074},
	expr: &actionExpr{
	pos: position{line: 102, col: 18, offset: 2093},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 102, col: 18, offset: 2093},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 102, col: 18, offset: 2093},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&oneOrMoreExpr{
	pos: position{line: 102, col: 22, offset: 2097},
	expr: &charClassMatcher{
	pos: position{line: 102, col: 22, offset: 2097},
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
	pos: position{line: 110, col: 1, offset: 2245},
	expr: &actionExpr{
	pos: position{line: 110, col: 17, offset: 2263},
	run: (*parser).callonSpaceDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 110, col: 17, offset: 2263},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 110, col: 17, offset: 2263},
	name: "_",
},
&litMatcher{
	pos: position{line: 110, col: 19, offset: 2265},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 110, col: 23, offset: 2269},
	name: "_",
},
&labeledExpr{
	pos: position{line: 110, col: 25, offset: 2271},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 110, col: 31, offset: 2277},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Identifier",
	pos: position{line: 112, col: 1, offset: 2337},
	expr: &actionExpr{
	pos: position{line: 112, col: 14, offset: 2352},
	run: (*parser).callonIdentifier1,
	expr: &seqExpr{
	pos: position{line: 112, col: 14, offset: 2352},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 112, col: 14, offset: 2352},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 112, col: 19, offset: 2357},
	name: "Label",
},
},
&labeledExpr{
	pos: position{line: 112, col: 25, offset: 2363},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 112, col: 31, offset: 2369},
	expr: &ruleRefExpr{
	pos: position{line: 112, col: 31, offset: 2369},
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
	pos: position{line: 120, col: 1, offset: 2545},
	expr: &actionExpr{
	pos: position{line: 121, col: 10, offset: 2583},
	run: (*parser).callonIdentifierReservedPrefix1,
	expr: &seqExpr{
	pos: position{line: 121, col: 10, offset: 2583},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 121, col: 10, offset: 2583},
	label: "name",
	expr: &actionExpr{
	pos: position{line: 121, col: 16, offset: 2589},
	run: (*parser).callonIdentifierReservedPrefix4,
	expr: &seqExpr{
	pos: position{line: 121, col: 16, offset: 2589},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 121, col: 16, offset: 2589},
	name: "Reserved",
},
&oneOrMoreExpr{
	pos: position{line: 121, col: 25, offset: 2598},
	expr: &charClassMatcher{
	pos: position{line: 121, col: 25, offset: 2598},
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
	pos: position{line: 122, col: 10, offset: 2654},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 122, col: 16, offset: 2660},
	expr: &ruleRefExpr{
	pos: position{line: 122, col: 16, offset: 2660},
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
	pos: position{line: 139, col: 1, offset: 3348},
	expr: &actionExpr{
	pos: position{line: 139, col: 14, offset: 3363},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 139, col: 14, offset: 3363},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 139, col: 14, offset: 3363},
	val: "let",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 139, col: 20, offset: 3369},
	name: "_",
},
&labeledExpr{
	pos: position{line: 139, col: 22, offset: 3371},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 139, col: 28, offset: 3377},
	name: "Label",
},
},
&labeledExpr{
	pos: position{line: 139, col: 34, offset: 3383},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 139, col: 36, offset: 3385},
	expr: &ruleRefExpr{
	pos: position{line: 139, col: 36, offset: 3385},
	name: "Annotation",
},
},
},
&ruleRefExpr{
	pos: position{line: 139, col: 48, offset: 3397},
	name: "_",
},
&litMatcher{
	pos: position{line: 139, col: 50, offset: 3399},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 139, col: 54, offset: 3403},
	name: "_",
},
&labeledExpr{
	pos: position{line: 139, col: 56, offset: 3405},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 139, col: 58, offset: 3407},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 139, col: 69, offset: 3418},
	name: "_",
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 154, col: 1, offset: 3728},
	expr: &choiceExpr{
	pos: position{line: 155, col: 7, offset: 3749},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 155, col: 7, offset: 3749},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 155, col: 7, offset: 3749},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 155, col: 7, offset: 3749},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 155, col: 14, offset: 3756},
	name: "_",
},
&litMatcher{
	pos: position{line: 155, col: 16, offset: 3758},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 155, col: 20, offset: 3762},
	name: "_",
},
&labeledExpr{
	pos: position{line: 155, col: 22, offset: 3764},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 155, col: 28, offset: 3770},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 155, col: 34, offset: 3776},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 155, col: 36, offset: 3778},
	name: "ColonSpace",
},
&labeledExpr{
	pos: position{line: 155, col: 47, offset: 3789},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 155, col: 49, offset: 3791},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 155, col: 60, offset: 3802},
	name: "_",
},
&litMatcher{
	pos: position{line: 155, col: 62, offset: 3804},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 155, col: 66, offset: 3808},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 155, col: 68, offset: 3810},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 155, col: 74, offset: 3816},
	name: "_",
},
&labeledExpr{
	pos: position{line: 155, col: 76, offset: 3818},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 155, col: 81, offset: 3823},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 158, col: 7, offset: 3952},
	run: (*parser).callonExpression21,
	expr: &seqExpr{
	pos: position{line: 158, col: 7, offset: 3952},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 158, col: 7, offset: 3952},
	val: "if",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 158, col: 12, offset: 3957},
	name: "_",
},
&labeledExpr{
	pos: position{line: 158, col: 14, offset: 3959},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 158, col: 19, offset: 3964},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 158, col: 30, offset: 3975},
	name: "_",
},
&litMatcher{
	pos: position{line: 158, col: 32, offset: 3977},
	val: "then",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 158, col: 39, offset: 3984},
	name: "_",
},
&labeledExpr{
	pos: position{line: 158, col: 41, offset: 3986},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 158, col: 43, offset: 3988},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 158, col: 54, offset: 3999},
	name: "_",
},
&litMatcher{
	pos: position{line: 158, col: 56, offset: 4001},
	val: "else",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 158, col: 63, offset: 4008},
	name: "_",
},
&labeledExpr{
	pos: position{line: 158, col: 65, offset: 4010},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 158, col: 67, offset: 4012},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 161, col: 7, offset: 4114},
	run: (*parser).callonExpression37,
	expr: &seqExpr{
	pos: position{line: 161, col: 7, offset: 4114},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 161, col: 7, offset: 4114},
	label: "bindings",
	expr: &zeroOrMoreExpr{
	pos: position{line: 161, col: 16, offset: 4123},
	expr: &ruleRefExpr{
	pos: position{line: 161, col: 16, offset: 4123},
	name: "LetBinding",
},
},
},
&litMatcher{
	pos: position{line: 161, col: 28, offset: 4135},
	val: "in",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 161, col: 33, offset: 4140},
	name: "_",
},
&labeledExpr{
	pos: position{line: 161, col: 35, offset: 4142},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 161, col: 37, offset: 4144},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 168, col: 7, offset: 4400},
	run: (*parser).callonExpression46,
	expr: &seqExpr{
	pos: position{line: 168, col: 7, offset: 4400},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 168, col: 7, offset: 4400},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 168, col: 14, offset: 4407},
	name: "_",
},
&litMatcher{
	pos: position{line: 168, col: 16, offset: 4409},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 168, col: 20, offset: 4413},
	name: "_",
},
&labeledExpr{
	pos: position{line: 168, col: 22, offset: 4415},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 168, col: 28, offset: 4421},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 168, col: 34, offset: 4427},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 168, col: 36, offset: 4429},
	name: "ColonSpace",
},
&labeledExpr{
	pos: position{line: 168, col: 47, offset: 4440},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 168, col: 49, offset: 4442},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 168, col: 60, offset: 4453},
	name: "_",
},
&litMatcher{
	pos: position{line: 168, col: 62, offset: 4455},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 168, col: 66, offset: 4459},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 168, col: 68, offset: 4461},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 168, col: 74, offset: 4467},
	name: "_",
},
&labeledExpr{
	pos: position{line: 168, col: 76, offset: 4469},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 168, col: 81, offset: 4474},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 171, col: 7, offset: 4595},
	run: (*parser).callonExpression65,
	expr: &seqExpr{
	pos: position{line: 171, col: 7, offset: 4595},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 171, col: 7, offset: 4595},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 171, col: 9, offset: 4597},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 171, col: 28, offset: 4616},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 171, col: 30, offset: 4618},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 171, col: 36, offset: 4624},
	name: "_",
},
&labeledExpr{
	pos: position{line: 171, col: 38, offset: 4626},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 171, col: 40, offset: 4628},
	name: "Expression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 172, col: 7, offset: 4700},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 174, col: 1, offset: 4721},
	expr: &actionExpr{
	pos: position{line: 174, col: 14, offset: 4736},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 174, col: 14, offset: 4736},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 174, col: 14, offset: 4736},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 174, col: 16, offset: 4738},
	name: "ColonSpace",
},
&labeledExpr{
	pos: position{line: 174, col: 27, offset: 4749},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 174, col: 29, offset: 4751},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 176, col: 1, offset: 4781},
	expr: &choiceExpr{
	pos: position{line: 177, col: 5, offset: 4809},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 177, col: 5, offset: 4809},
	name: "EmptyList",
},
&actionExpr{
	pos: position{line: 178, col: 5, offset: 4823},
	run: (*parser).callonAnnotatedExpression3,
	expr: &seqExpr{
	pos: position{line: 178, col: 5, offset: 4823},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 178, col: 5, offset: 4823},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 178, col: 7, offset: 4825},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 178, col: 26, offset: 4844},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 178, col: 28, offset: 4846},
	expr: &ruleRefExpr{
	pos: position{line: 178, col: 28, offset: 4846},
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
	pos: position{line: 183, col: 1, offset: 4963},
	expr: &actionExpr{
	pos: position{line: 183, col: 13, offset: 4977},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 183, col: 13, offset: 4977},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 183, col: 13, offset: 4977},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 183, col: 17, offset: 4981},
	name: "_",
},
&litMatcher{
	pos: position{line: 183, col: 19, offset: 4983},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 183, col: 23, offset: 4987},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 183, col: 25, offset: 4989},
	name: "ColonSpace",
},
&litMatcher{
	pos: position{line: 183, col: 36, offset: 5000},
	val: "List",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 183, col: 43, offset: 5007},
	name: "_",
},
&labeledExpr{
	pos: position{line: 183, col: 45, offset: 5009},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 183, col: 47, offset: 5011},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 187, col: 1, offset: 5082},
	expr: &ruleRefExpr{
	pos: position{line: 187, col: 22, offset: 5105},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 189, col: 1, offset: 5126},
	expr: &ruleRefExpr{
	pos: position{line: 189, col: 23, offset: 5150},
	name: "PlusExpression",
},
},
{
	name: "MorePlus",
	pos: position{line: 191, col: 1, offset: 5166},
	expr: &actionExpr{
	pos: position{line: 191, col: 12, offset: 5179},
	run: (*parser).callonMorePlus1,
	expr: &seqExpr{
	pos: position{line: 191, col: 12, offset: 5179},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 191, col: 12, offset: 5179},
	name: "_",
},
&litMatcher{
	pos: position{line: 191, col: 14, offset: 5181},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 191, col: 18, offset: 5185},
	name: "NonemptyWhitespace",
},
&labeledExpr{
	pos: position{line: 191, col: 37, offset: 5204},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 191, col: 39, offset: 5206},
	name: "TimesExpression",
},
},
	},
},
},
},
{
	name: "PlusExpression",
	pos: position{line: 192, col: 1, offset: 5240},
	expr: &actionExpr{
	pos: position{line: 193, col: 7, offset: 5265},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 193, col: 7, offset: 5265},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 193, col: 7, offset: 5265},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 193, col: 13, offset: 5271},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 193, col: 29, offset: 5287},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 193, col: 34, offset: 5292},
	expr: &ruleRefExpr{
	pos: position{line: 193, col: 34, offset: 5292},
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
	pos: position{line: 202, col: 1, offset: 5532},
	expr: &actionExpr{
	pos: position{line: 202, col: 13, offset: 5546},
	run: (*parser).callonMoreTimes1,
	expr: &seqExpr{
	pos: position{line: 202, col: 13, offset: 5546},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 202, col: 13, offset: 5546},
	name: "_",
},
&litMatcher{
	pos: position{line: 202, col: 15, offset: 5548},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 202, col: 19, offset: 5552},
	name: "_",
},
&labeledExpr{
	pos: position{line: 202, col: 21, offset: 5554},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 202, col: 23, offset: 5556},
	name: "ApplicationExpression",
},
},
	},
},
},
},
{
	name: "TimesExpression",
	pos: position{line: 203, col: 1, offset: 5596},
	expr: &actionExpr{
	pos: position{line: 204, col: 7, offset: 5622},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 204, col: 7, offset: 5622},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 204, col: 7, offset: 5622},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 204, col: 13, offset: 5628},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 204, col: 35, offset: 5650},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 204, col: 40, offset: 5655},
	expr: &ruleRefExpr{
	pos: position{line: 204, col: 40, offset: 5655},
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
	pos: position{line: 213, col: 1, offset: 5897},
	expr: &actionExpr{
	pos: position{line: 213, col: 9, offset: 5905},
	run: (*parser).callonAnArg1,
	expr: &seqExpr{
	pos: position{line: 213, col: 9, offset: 5905},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 213, col: 9, offset: 5905},
	name: "NonemptyWhitespace",
},
&labeledExpr{
	pos: position{line: 213, col: 28, offset: 5924},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 213, col: 30, offset: 5926},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "ApplicationExpression",
	pos: position{line: 215, col: 1, offset: 5962},
	expr: &actionExpr{
	pos: position{line: 215, col: 25, offset: 5988},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 215, col: 25, offset: 5988},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 215, col: 25, offset: 5988},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 215, col: 27, offset: 5990},
	name: "ImportExpression",
},
},
&labeledExpr{
	pos: position{line: 215, col: 44, offset: 6007},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 215, col: 49, offset: 6012},
	expr: &ruleRefExpr{
	pos: position{line: 215, col: 49, offset: 6012},
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
	pos: position{line: 224, col: 1, offset: 6243},
	expr: &ruleRefExpr{
	pos: position{line: 224, col: 20, offset: 6264},
	name: "SelectorExpression",
},
},
{
	name: "SelectorExpression",
	pos: position{line: 226, col: 1, offset: 6284},
	expr: &ruleRefExpr{
	pos: position{line: 226, col: 22, offset: 6307},
	name: "PrimitiveExpression",
},
},
{
	name: "PrimitiveExpression",
	pos: position{line: 228, col: 1, offset: 6328},
	expr: &choiceExpr{
	pos: position{line: 229, col: 7, offset: 6358},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 229, col: 7, offset: 6358},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 230, col: 7, offset: 6378},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 231, col: 7, offset: 6399},
	name: "IntegerLiteral",
},
&actionExpr{
	pos: position{line: 232, col: 7, offset: 6420},
	run: (*parser).callonPrimitiveExpression5,
	expr: &litMatcher{
	pos: position{line: 232, col: 7, offset: 6420},
	val: "-Infinity",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 233, col: 7, offset: 6482},
	run: (*parser).callonPrimitiveExpression7,
	expr: &seqExpr{
	pos: position{line: 233, col: 7, offset: 6482},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 233, col: 7, offset: 6482},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 233, col: 11, offset: 6486},
	name: "_",
},
&labeledExpr{
	pos: position{line: 233, col: 13, offset: 6488},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 233, col: 15, offset: 6490},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 233, col: 35, offset: 6510},
	name: "_",
},
&litMatcher{
	pos: position{line: 233, col: 37, offset: 6512},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 234, col: 7, offset: 6540},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 235, col: 7, offset: 6566},
	name: "IdentifierReservedPrefix",
},
&ruleRefExpr{
	pos: position{line: 236, col: 7, offset: 6597},
	name: "Reserved",
},
&ruleRefExpr{
	pos: position{line: 237, col: 7, offset: 6612},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 238, col: 7, offset: 6629},
	run: (*parser).callonPrimitiveExpression19,
	expr: &seqExpr{
	pos: position{line: 238, col: 7, offset: 6629},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 238, col: 7, offset: 6629},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 238, col: 11, offset: 6633},
	name: "_",
},
&labeledExpr{
	pos: position{line: 238, col: 13, offset: 6635},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 238, col: 15, offset: 6637},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 238, col: 26, offset: 6648},
	name: "_",
},
&litMatcher{
	pos: position{line: 238, col: 28, offset: 6650},
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
	pos: position{line: 240, col: 1, offset: 6673},
	expr: &choiceExpr{
	pos: position{line: 241, col: 7, offset: 6703},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 241, col: 7, offset: 6703},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &litMatcher{
	pos: position{line: 241, col: 7, offset: 6703},
	val: "=",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 242, col: 7, offset: 6766},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 243, col: 7, offset: 6791},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 244, col: 7, offset: 6819},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 244, col: 7, offset: 6819},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 246, col: 1, offset: 6873},
	expr: &actionExpr{
	pos: position{line: 246, col: 19, offset: 6893},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 246, col: 19, offset: 6893},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 246, col: 19, offset: 6893},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 246, col: 24, offset: 6898},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 246, col: 30, offset: 6904},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 246, col: 32, offset: 6906},
	name: "ColonSpace",
},
&labeledExpr{
	pos: position{line: 246, col: 43, offset: 6917},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 246, col: 48, offset: 6922},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 249, col: 1, offset: 6979},
	expr: &actionExpr{
	pos: position{line: 249, col: 18, offset: 6998},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 249, col: 18, offset: 6998},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 249, col: 18, offset: 6998},
	name: "_",
},
&litMatcher{
	pos: position{line: 249, col: 20, offset: 7000},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 249, col: 24, offset: 7004},
	name: "_",
},
&labeledExpr{
	pos: position{line: 249, col: 26, offset: 7006},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 249, col: 28, offset: 7008},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 250, col: 1, offset: 7040},
	expr: &actionExpr{
	pos: position{line: 251, col: 7, offset: 7069},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 251, col: 7, offset: 7069},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 251, col: 7, offset: 7069},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 251, col: 13, offset: 7075},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 251, col: 29, offset: 7091},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 251, col: 34, offset: 7096},
	expr: &ruleRefExpr{
	pos: position{line: 251, col: 34, offset: 7096},
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
	pos: position{line: 261, col: 1, offset: 7508},
	expr: &actionExpr{
	pos: position{line: 261, col: 22, offset: 7531},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 261, col: 22, offset: 7531},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 261, col: 22, offset: 7531},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 261, col: 27, offset: 7536},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 261, col: 33, offset: 7542},
	name: "_",
},
&litMatcher{
	pos: position{line: 261, col: 35, offset: 7544},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 261, col: 39, offset: 7548},
	name: "_",
},
&labeledExpr{
	pos: position{line: 261, col: 41, offset: 7550},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 261, col: 46, offset: 7555},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 264, col: 1, offset: 7612},
	expr: &actionExpr{
	pos: position{line: 264, col: 21, offset: 7634},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 264, col: 21, offset: 7634},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 264, col: 21, offset: 7634},
	name: "_",
},
&litMatcher{
	pos: position{line: 264, col: 23, offset: 7636},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 264, col: 27, offset: 7640},
	name: "_",
},
&labeledExpr{
	pos: position{line: 264, col: 29, offset: 7642},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 264, col: 31, offset: 7644},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 265, col: 1, offset: 7679},
	expr: &actionExpr{
	pos: position{line: 266, col: 7, offset: 7711},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 266, col: 7, offset: 7711},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 266, col: 7, offset: 7711},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 266, col: 13, offset: 7717},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 266, col: 32, offset: 7736},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 266, col: 37, offset: 7741},
	expr: &ruleRefExpr{
	pos: position{line: 266, col: 37, offset: 7741},
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
	pos: position{line: 276, col: 1, offset: 8159},
	expr: &actionExpr{
	pos: position{line: 276, col: 12, offset: 8172},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 276, col: 12, offset: 8172},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 276, col: 12, offset: 8172},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 276, col: 16, offset: 8176},
	name: "_",
},
&labeledExpr{
	pos: position{line: 276, col: 18, offset: 8178},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 276, col: 20, offset: 8180},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 276, col: 31, offset: 8191},
	name: "_",
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 278, col: 1, offset: 8210},
	expr: &actionExpr{
	pos: position{line: 279, col: 7, offset: 8240},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 279, col: 7, offset: 8240},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 279, col: 7, offset: 8240},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 279, col: 11, offset: 8244},
	name: "_",
},
&labeledExpr{
	pos: position{line: 279, col: 13, offset: 8246},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 279, col: 19, offset: 8252},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 279, col: 30, offset: 8263},
	name: "_",
},
&labeledExpr{
	pos: position{line: 279, col: 32, offset: 8265},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 279, col: 37, offset: 8270},
	expr: &ruleRefExpr{
	pos: position{line: 279, col: 37, offset: 8270},
	name: "MoreList",
},
},
},
&litMatcher{
	pos: position{line: 279, col: 47, offset: 8280},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 289, col: 1, offset: 8572},
	expr: &notExpr{
	pos: position{line: 289, col: 7, offset: 8580},
	expr: &anyMatcher{
	line: 289, col: 8, offset: 8581,
},
},
},
	},
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

func (c *current) onReserved13() (interface{}, error) {
 return ast.List, nil 
}

func (p *parser) callonReserved13() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved13()
}

func (c *current) onReserved15() (interface{}, error) {
 return ast.True, nil 
}

func (p *parser) callonReserved15() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved15()
}

func (c *current) onReserved17() (interface{}, error) {
 return ast.False, nil 
}

func (p *parser) callonReserved17() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved17()
}

func (c *current) onReserved19() (interface{}, error) {
 return ast.DoubleLit(math.NaN()), nil 
}

func (p *parser) callonReserved19() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved19()
}

func (c *current) onReserved21() (interface{}, error) {
 return ast.DoubleLit(math.Inf(1)), nil 
}

func (p *parser) callonReserved21() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved21()
}

func (c *current) onReserved23() (interface{}, error) {
 return ast.Type, nil 
}

func (p *parser) callonReserved23() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved23()
}

func (c *current) onReserved25() (interface{}, error) {
 return ast.Kind, nil 
}

func (p *parser) callonReserved25() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved25()
}

func (c *current) onReserved27() (interface{}, error) {
 return ast.Sort, nil 
}

func (p *parser) callonReserved27() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved27()
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

func (c *current) onPrimitiveExpression7(r interface{}) (interface{}, error) {
 return r, nil 
}

func (p *parser) callonPrimitiveExpression7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimitiveExpression7(stack["r"])
}

func (c *current) onPrimitiveExpression19(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonPrimitiveExpression19() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimitiveExpression19(stack["e"])
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

