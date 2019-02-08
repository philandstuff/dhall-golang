
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
&litMatcher{
	pos: position{line: 51, col: 12, offset: 895},
	val: "Bool",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 52, col: 5, offset: 906},
	val: "Optional",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 53, col: 5, offset: 921},
	val: "None",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 54, col: 5, offset: 932},
	run: (*parser).callonReserved5,
	expr: &litMatcher{
	pos: position{line: 54, col: 5, offset: 932},
	val: "Natural",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 55, col: 5, offset: 974},
	run: (*parser).callonReserved7,
	expr: &litMatcher{
	pos: position{line: 55, col: 5, offset: 974},
	val: "Integer",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 56, col: 5, offset: 1016},
	run: (*parser).callonReserved9,
	expr: &litMatcher{
	pos: position{line: 56, col: 5, offset: 1016},
	val: "Double",
	ignoreCase: false,
},
},
&litMatcher{
	pos: position{line: 57, col: 5, offset: 1056},
	val: "Text",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 58, col: 5, offset: 1067},
	val: "List",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 59, col: 5, offset: 1078},
	val: "True",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 60, col: 5, offset: 1089},
	val: "False",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 61, col: 5, offset: 1101},
	run: (*parser).callonReserved15,
	expr: &litMatcher{
	pos: position{line: 61, col: 5, offset: 1101},
	val: "NaN",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 62, col: 5, offset: 1153},
	run: (*parser).callonReserved17,
	expr: &litMatcher{
	pos: position{line: 62, col: 5, offset: 1153},
	val: "Infinity",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 63, col: 5, offset: 1211},
	run: (*parser).callonReserved19,
	expr: &litMatcher{
	pos: position{line: 63, col: 5, offset: 1211},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 64, col: 5, offset: 1247},
	run: (*parser).callonReserved21,
	expr: &litMatcher{
	pos: position{line: 64, col: 5, offset: 1247},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 65, col: 5, offset: 1283},
	run: (*parser).callonReserved23,
	expr: &litMatcher{
	pos: position{line: 65, col: 5, offset: 1283},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "Keyword",
	pos: position{line: 67, col: 1, offset: 1316},
	expr: &choiceExpr{
	pos: position{line: 67, col: 11, offset: 1328},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 67, col: 11, offset: 1328},
	val: "if",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 68, col: 5, offset: 1337},
	val: "then",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 69, col: 5, offset: 1348},
	val: "else",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 70, col: 5, offset: 1359},
	val: "let",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 71, col: 5, offset: 1369},
	val: "in",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 72, col: 5, offset: 1378},
	val: "as",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 73, col: 5, offset: 1387},
	val: "using",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 74, col: 5, offset: 1399},
	val: "merge",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 75, col: 5, offset: 1411},
	val: "constructors",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 76, col: 5, offset: 1430},
	val: "Some",
	ignoreCase: false,
},
	},
},
},
{
	name: "ColonSpace",
	pos: position{line: 78, col: 1, offset: 1438},
	expr: &seqExpr{
	pos: position{line: 78, col: 14, offset: 1453},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 78, col: 14, offset: 1453},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 78, col: 18, offset: 1457},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Lambda",
	pos: position{line: 80, col: 1, offset: 1477},
	expr: &choiceExpr{
	pos: position{line: 80, col: 11, offset: 1489},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 80, col: 11, offset: 1489},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 80, col: 18, offset: 1496},
	val: "λ",
	ignoreCase: false,
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 81, col: 1, offset: 1502},
	expr: &choiceExpr{
	pos: position{line: 81, col: 11, offset: 1514},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 81, col: 11, offset: 1514},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 81, col: 22, offset: 1525},
	val: "∀",
	ignoreCase: false,
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 82, col: 1, offset: 1532},
	expr: &choiceExpr{
	pos: position{line: 82, col: 10, offset: 1543},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 82, col: 10, offset: 1543},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 82, col: 17, offset: 1550},
	val: "→",
	ignoreCase: false,
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 84, col: 1, offset: 1558},
	expr: &seqExpr{
	pos: position{line: 84, col: 12, offset: 1571},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 84, col: 12, offset: 1571},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 84, col: 17, offset: 1576},
	expr: &charClassMatcher{
	pos: position{line: 84, col: 17, offset: 1576},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 84, col: 23, offset: 1582},
	expr: &charClassMatcher{
	pos: position{line: 84, col: 23, offset: 1582},
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
	pos: position{line: 86, col: 1, offset: 1590},
	expr: &actionExpr{
	pos: position{line: 86, col: 17, offset: 1608},
	run: (*parser).callonDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 86, col: 17, offset: 1608},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 86, col: 17, offset: 1608},
	expr: &charClassMatcher{
	pos: position{line: 86, col: 17, offset: 1608},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 86, col: 23, offset: 1614},
	expr: &charClassMatcher{
	pos: position{line: 86, col: 23, offset: 1614},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&choiceExpr{
	pos: position{line: 86, col: 32, offset: 1623},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 86, col: 32, offset: 1623},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 86, col: 32, offset: 1623},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 86, col: 36, offset: 1627},
	expr: &charClassMatcher{
	pos: position{line: 86, col: 36, offset: 1627},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
&zeroOrOneExpr{
	pos: position{line: 86, col: 43, offset: 1634},
	expr: &ruleRefExpr{
	pos: position{line: 86, col: 43, offset: 1634},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 86, col: 55, offset: 1646},
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
	pos: position{line: 94, col: 1, offset: 1806},
	expr: &actionExpr{
	pos: position{line: 94, col: 18, offset: 1825},
	run: (*parser).callonNaturalLiteral1,
	expr: &oneOrMoreExpr{
	pos: position{line: 94, col: 18, offset: 1825},
	expr: &charClassMatcher{
	pos: position{line: 94, col: 18, offset: 1825},
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
	pos: position{line: 102, col: 1, offset: 1973},
	expr: &actionExpr{
	pos: position{line: 102, col: 18, offset: 1992},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 102, col: 18, offset: 1992},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 102, col: 18, offset: 1992},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&oneOrMoreExpr{
	pos: position{line: 102, col: 22, offset: 1996},
	expr: &charClassMatcher{
	pos: position{line: 102, col: 22, offset: 1996},
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
	pos: position{line: 110, col: 1, offset: 2144},
	expr: &actionExpr{
	pos: position{line: 110, col: 17, offset: 2162},
	run: (*parser).callonSpaceDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 110, col: 17, offset: 2162},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 110, col: 17, offset: 2162},
	name: "_",
},
&litMatcher{
	pos: position{line: 110, col: 19, offset: 2164},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 110, col: 23, offset: 2168},
	name: "_",
},
&labeledExpr{
	pos: position{line: 110, col: 25, offset: 2170},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 110, col: 31, offset: 2176},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Identifier",
	pos: position{line: 112, col: 1, offset: 2220},
	expr: &actionExpr{
	pos: position{line: 112, col: 14, offset: 2235},
	run: (*parser).callonIdentifier1,
	expr: &seqExpr{
	pos: position{line: 112, col: 14, offset: 2235},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 112, col: 14, offset: 2235},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 112, col: 19, offset: 2240},
	name: "Label",
},
},
&labeledExpr{
	pos: position{line: 112, col: 25, offset: 2246},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 112, col: 31, offset: 2252},
	expr: &ruleRefExpr{
	pos: position{line: 112, col: 31, offset: 2252},
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
	pos: position{line: 120, col: 1, offset: 2428},
	expr: &actionExpr{
	pos: position{line: 120, col: 28, offset: 2457},
	run: (*parser).callonIdentifierReservedPrefix1,
	expr: &seqExpr{
	pos: position{line: 120, col: 28, offset: 2457},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 120, col: 28, offset: 2457},
	label: "reserved",
	expr: &ruleRefExpr{
	pos: position{line: 120, col: 37, offset: 2466},
	name: "Reserved",
},
},
&labeledExpr{
	pos: position{line: 120, col: 46, offset: 2475},
	label: "suffix",
	expr: &actionExpr{
	pos: position{line: 120, col: 54, offset: 2483},
	run: (*parser).callonIdentifierReservedPrefix6,
	expr: &oneOrMoreExpr{
	pos: position{line: 120, col: 54, offset: 2483},
	expr: &charClassMatcher{
	pos: position{line: 120, col: 54, offset: 2483},
	val: "[A-Za-z0-9/_-]",
	chars: []rune{'/','_','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
},
},
&labeledExpr{
	pos: position{line: 120, col: 102, offset: 2531},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 120, col: 108, offset: 2537},
	expr: &ruleRefExpr{
	pos: position{line: 120, col: 108, offset: 2537},
	name: "SpaceDeBruijn",
},
},
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 138, col: 1, offset: 3255},
	expr: &choiceExpr{
	pos: position{line: 139, col: 7, offset: 3276},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 139, col: 7, offset: 3276},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 139, col: 7, offset: 3276},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 139, col: 7, offset: 3276},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 139, col: 14, offset: 3283},
	name: "_",
},
&litMatcher{
	pos: position{line: 139, col: 16, offset: 3285},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 139, col: 20, offset: 3289},
	name: "_",
},
&labeledExpr{
	pos: position{line: 139, col: 22, offset: 3291},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 139, col: 28, offset: 3297},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 139, col: 34, offset: 3303},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 139, col: 36, offset: 3305},
	name: "ColonSpace",
},
&labeledExpr{
	pos: position{line: 139, col: 47, offset: 3316},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 139, col: 49, offset: 3318},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 139, col: 60, offset: 3329},
	name: "_",
},
&litMatcher{
	pos: position{line: 139, col: 62, offset: 3331},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 139, col: 66, offset: 3335},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 139, col: 68, offset: 3337},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 139, col: 74, offset: 3343},
	name: "_",
},
&labeledExpr{
	pos: position{line: 139, col: 76, offset: 3345},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 139, col: 81, offset: 3350},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 142, col: 7, offset: 3479},
	run: (*parser).callonExpression21,
	expr: &seqExpr{
	pos: position{line: 142, col: 7, offset: 3479},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 142, col: 7, offset: 3479},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 142, col: 14, offset: 3486},
	name: "_",
},
&litMatcher{
	pos: position{line: 142, col: 16, offset: 3488},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 142, col: 20, offset: 3492},
	name: "_",
},
&labeledExpr{
	pos: position{line: 142, col: 22, offset: 3494},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 142, col: 28, offset: 3500},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 142, col: 34, offset: 3506},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 142, col: 36, offset: 3508},
	name: "ColonSpace",
},
&labeledExpr{
	pos: position{line: 142, col: 47, offset: 3519},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 142, col: 49, offset: 3521},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 142, col: 60, offset: 3532},
	name: "_",
},
&litMatcher{
	pos: position{line: 142, col: 62, offset: 3534},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 142, col: 66, offset: 3538},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 142, col: 68, offset: 3540},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 142, col: 74, offset: 3546},
	name: "_",
},
&labeledExpr{
	pos: position{line: 142, col: 76, offset: 3548},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 142, col: 81, offset: 3553},
	name: "Expression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 145, col: 7, offset: 3674},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 147, col: 1, offset: 3695},
	expr: &ruleRefExpr{
	pos: position{line: 147, col: 23, offset: 3719},
	name: "OperatorExpression",
},
},
{
	name: "OperatorExpression",
	pos: position{line: 149, col: 1, offset: 3739},
	expr: &ruleRefExpr{
	pos: position{line: 149, col: 22, offset: 3762},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 151, col: 1, offset: 3783},
	expr: &ruleRefExpr{
	pos: position{line: 152, col: 7, offset: 3849},
	name: "ApplicationExpression",
},
},
{
	name: "AnArg",
	pos: position{line: 154, col: 1, offset: 3872},
	expr: &actionExpr{
	pos: position{line: 154, col: 9, offset: 3880},
	run: (*parser).callonAnArg1,
	expr: &seqExpr{
	pos: position{line: 154, col: 9, offset: 3880},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 154, col: 9, offset: 3880},
	name: "NonemptyWhitespace",
},
&labeledExpr{
	pos: position{line: 154, col: 28, offset: 3899},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 154, col: 30, offset: 3901},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "ApplicationExpression",
	pos: position{line: 156, col: 1, offset: 3937},
	expr: &actionExpr{
	pos: position{line: 156, col: 25, offset: 3963},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 156, col: 25, offset: 3963},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 156, col: 25, offset: 3963},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 156, col: 27, offset: 3965},
	name: "ImportExpression",
},
},
&labeledExpr{
	pos: position{line: 156, col: 44, offset: 3982},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 156, col: 49, offset: 3987},
	expr: &ruleRefExpr{
	pos: position{line: 156, col: 49, offset: 3987},
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
	pos: position{line: 165, col: 1, offset: 4218},
	expr: &ruleRefExpr{
	pos: position{line: 165, col: 20, offset: 4239},
	name: "SelectorExpression",
},
},
{
	name: "SelectorExpression",
	pos: position{line: 167, col: 1, offset: 4259},
	expr: &ruleRefExpr{
	pos: position{line: 167, col: 22, offset: 4282},
	name: "PrimitiveExpression",
},
},
{
	name: "PrimitiveExpression",
	pos: position{line: 169, col: 1, offset: 4303},
	expr: &choiceExpr{
	pos: position{line: 170, col: 7, offset: 4333},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 170, col: 7, offset: 4333},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 171, col: 7, offset: 4353},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 172, col: 7, offset: 4374},
	name: "IntegerLiteral",
},
&actionExpr{
	pos: position{line: 173, col: 7, offset: 4395},
	run: (*parser).callonPrimitiveExpression5,
	expr: &litMatcher{
	pos: position{line: 173, col: 7, offset: 4395},
	val: "-Infinity",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 174, col: 7, offset: 4457},
	name: "IdentifierReservedPrefix",
},
&ruleRefExpr{
	pos: position{line: 175, col: 7, offset: 4488},
	name: "Reserved",
},
&ruleRefExpr{
	pos: position{line: 176, col: 7, offset: 4503},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 177, col: 7, offset: 4520},
	run: (*parser).callonPrimitiveExpression10,
	expr: &seqExpr{
	pos: position{line: 177, col: 7, offset: 4520},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 177, col: 7, offset: 4520},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 177, col: 11, offset: 4524},
	name: "_",
},
&labeledExpr{
	pos: position{line: 177, col: 13, offset: 4526},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 177, col: 15, offset: 4528},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 177, col: 26, offset: 4539},
	name: "_",
},
&litMatcher{
	pos: position{line: 177, col: 28, offset: 4541},
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
	name: "EOF",
	pos: position{line: 179, col: 1, offset: 4564},
	expr: &notExpr{
	pos: position{line: 179, col: 7, offset: 4572},
	expr: &anyMatcher{
	line: 179, col: 8, offset: 4573,
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

func (c *current) onReserved5() (interface{}, error) {
 return ast.Natural, nil 
}

func (p *parser) callonReserved5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved5()
}

func (c *current) onReserved7() (interface{}, error) {
 return ast.Integer, nil 
}

func (p *parser) callonReserved7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved7()
}

func (c *current) onReserved9() (interface{}, error) {
 return ast.Double, nil 
}

func (p *parser) callonReserved9() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved9()
}

func (c *current) onReserved15() (interface{}, error) {
 return ast.DoubleLit(math.NaN()), nil 
}

func (p *parser) callonReserved15() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved15()
}

func (c *current) onReserved17() (interface{}, error) {
 return ast.DoubleLit(math.Inf(1)), nil 
}

func (p *parser) callonReserved17() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved17()
}

func (c *current) onReserved19() (interface{}, error) {
 return ast.Type, nil 
}

func (p *parser) callonReserved19() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved19()
}

func (c *current) onReserved21() (interface{}, error) {
 return ast.Kind, nil 
}

func (p *parser) callonReserved21() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved21()
}

func (c *current) onReserved23() (interface{}, error) {
 return ast.Sort, nil 
}

func (p *parser) callonReserved23() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved23()
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
 return index.(int), nil 
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

func (c *current) onIdentifierReservedPrefix6() (interface{}, error) {
 return string(c.text), nil 
}

func (p *parser) callonIdentifierReservedPrefix6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifierReservedPrefix6()
}

func (c *current) onIdentifierReservedPrefix1(reserved, suffix, index interface{}) (interface{}, error) {
    name := reserved.(string) + suffix.(string)
    if index != nil {
        return ast.Var{Name:name, Index:index.(int)}, nil
    } else {
        return ast.Var{Name:name}, nil
    }
}

func (p *parser) callonIdentifierReservedPrefix1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifierReservedPrefix1(stack["reserved"], stack["suffix"], stack["index"])
}

func (c *current) onExpression2(label, t, body interface{}) (interface{}, error) {
          return &ast.LambdaExpr{Label:label.(string), Type:t.(ast.Expr), Body: body.(ast.Expr)}, nil
      
}

func (p *parser) callonExpression2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression2(stack["label"], stack["t"], stack["body"])
}

func (c *current) onExpression21(label, t, body interface{}) (interface{}, error) {
          return &ast.Pi{Label:label.(string), Type:t.(ast.Expr), Body: body.(ast.Expr)}, nil
      
}

func (p *parser) callonExpression21() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression21(stack["label"], stack["t"], stack["body"])
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

func (c *current) onPrimitiveExpression10(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonPrimitiveExpression10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimitiveExpression10(stack["e"])
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

