
package parser

import (
"bytes"
"errors"
"fmt"
"io"
"io/ioutil"
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
	pos: position{line: 20, col: 1, offset: 173},
	expr: &actionExpr{
	pos: position{line: 20, col: 22, offset: 196},
	run: (*parser).callonCompleteExpression1,
	expr: &seqExpr{
	pos: position{line: 20, col: 22, offset: 196},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 20, col: 22, offset: 196},
	name: "_",
},
&labeledExpr{
	pos: position{line: 20, col: 24, offset: 198},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 20, col: 26, offset: 200},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 20, col: 37, offset: 211},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 20, col: 39, offset: 213},
	name: "EOF",
},
	},
},
},
},
{
	name: "EOL",
	pos: position{line: 22, col: 1, offset: 236},
	expr: &choiceExpr{
	pos: position{line: 22, col: 7, offset: 244},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 22, col: 7, offset: 244},
	val: "\n",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 22, col: 14, offset: 251},
	val: "\r\n",
	ignoreCase: false,
},
	},
},
},
{
	name: "NotEOL",
	pos: position{line: 24, col: 1, offset: 259},
	expr: &charClassMatcher{
	pos: position{line: 24, col: 10, offset: 270},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "LineComment",
	pos: position{line: 26, col: 1, offset: 293},
	expr: &actionExpr{
	pos: position{line: 26, col: 15, offset: 309},
	run: (*parser).callonLineComment1,
	expr: &seqExpr{
	pos: position{line: 26, col: 15, offset: 309},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 26, col: 15, offset: 309},
	val: "--",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 26, col: 20, offset: 314},
	label: "content",
	expr: &actionExpr{
	pos: position{line: 26, col: 29, offset: 323},
	run: (*parser).callonLineComment5,
	expr: &zeroOrMoreExpr{
	pos: position{line: 26, col: 29, offset: 323},
	expr: &ruleRefExpr{
	pos: position{line: 26, col: 29, offset: 323},
	name: "NotEOL",
},
},
},
},
&ruleRefExpr{
	pos: position{line: 26, col: 68, offset: 362},
	name: "EOL",
},
	},
},
},
},
{
	name: "WhitespaceChunk",
	pos: position{line: 28, col: 1, offset: 391},
	expr: &choiceExpr{
	pos: position{line: 28, col: 19, offset: 411},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 28, col: 19, offset: 411},
	val: " ",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 28, col: 25, offset: 417},
	val: "\t",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 28, col: 32, offset: 424},
	name: "EOL",
},
&ruleRefExpr{
	pos: position{line: 28, col: 38, offset: 430},
	name: "LineComment",
},
	},
},
},
{
	name: "_",
	pos: position{line: 30, col: 1, offset: 461},
	expr: &zeroOrMoreExpr{
	pos: position{line: 30, col: 5, offset: 467},
	expr: &ruleRefExpr{
	pos: position{line: 30, col: 5, offset: 467},
	name: "WhitespaceChunk",
},
},
},
{
	name: "NonemptyWhitespace",
	pos: position{line: 32, col: 1, offset: 485},
	expr: &oneOrMoreExpr{
	pos: position{line: 32, col: 22, offset: 508},
	expr: &ruleRefExpr{
	pos: position{line: 32, col: 22, offset: 508},
	name: "WhitespaceChunk",
},
},
},
{
	name: "SimpleLabel",
	pos: position{line: 34, col: 1, offset: 526},
	expr: &actionExpr{
	pos: position{line: 34, col: 15, offset: 542},
	run: (*parser).callonSimpleLabel1,
	expr: &seqExpr{
	pos: position{line: 34, col: 15, offset: 542},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 34, col: 15, offset: 542},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 34, col: 25, offset: 552},
	expr: &charClassMatcher{
	pos: position{line: 34, col: 25, offset: 552},
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
	pos: position{line: 38, col: 1, offset: 616},
	expr: &actionExpr{
	pos: position{line: 38, col: 9, offset: 626},
	run: (*parser).callonLabel1,
	expr: &labeledExpr{
	pos: position{line: 38, col: 9, offset: 626},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 38, col: 15, offset: 632},
	name: "SimpleLabel",
},
},
},
},
{
	name: "Reserved",
	pos: position{line: 40, col: 1, offset: 667},
	expr: &choiceExpr{
	pos: position{line: 40, col: 12, offset: 680},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 40, col: 12, offset: 680},
	val: "Bool",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 41, col: 5, offset: 691},
	val: "Optional",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 42, col: 5, offset: 706},
	val: "None",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 43, col: 5, offset: 717},
	run: (*parser).callonReserved5,
	expr: &litMatcher{
	pos: position{line: 43, col: 5, offset: 717},
	val: "Natural",
	ignoreCase: false,
},
},
&litMatcher{
	pos: position{line: 44, col: 5, offset: 759},
	val: "Integer",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 45, col: 5, offset: 773},
	val: "Double",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 46, col: 5, offset: 786},
	val: "Text",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 47, col: 5, offset: 797},
	val: "List",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 48, col: 5, offset: 808},
	val: "True",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 49, col: 5, offset: 819},
	val: "False",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 50, col: 5, offset: 831},
	val: "NaN",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 51, col: 5, offset: 841},
	val: "Infinity",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 52, col: 5, offset: 856},
	run: (*parser).callonReserved15,
	expr: &litMatcher{
	pos: position{line: 52, col: 5, offset: 856},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 53, col: 5, offset: 892},
	run: (*parser).callonReserved17,
	expr: &litMatcher{
	pos: position{line: 53, col: 5, offset: 892},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 54, col: 5, offset: 928},
	run: (*parser).callonReserved19,
	expr: &litMatcher{
	pos: position{line: 54, col: 5, offset: 928},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "ColonSpace",
	pos: position{line: 56, col: 1, offset: 961},
	expr: &seqExpr{
	pos: position{line: 56, col: 14, offset: 976},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 56, col: 14, offset: 976},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 56, col: 18, offset: 980},
	name: "NonemptyWhitespace",
},
	},
},
},
{
	name: "Lambda",
	pos: position{line: 58, col: 1, offset: 1000},
	expr: &choiceExpr{
	pos: position{line: 58, col: 11, offset: 1012},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 58, col: 11, offset: 1012},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 58, col: 18, offset: 1019},
	val: "λ",
	ignoreCase: false,
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 59, col: 1, offset: 1025},
	expr: &choiceExpr{
	pos: position{line: 59, col: 11, offset: 1037},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 59, col: 11, offset: 1037},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 59, col: 22, offset: 1048},
	val: "∀",
	ignoreCase: false,
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 60, col: 1, offset: 1055},
	expr: &choiceExpr{
	pos: position{line: 60, col: 10, offset: 1066},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 60, col: 10, offset: 1066},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 60, col: 17, offset: 1073},
	val: "→",
	ignoreCase: false,
},
	},
},
},
{
	name: "NaturalLiteral",
	pos: position{line: 64, col: 1, offset: 1113},
	expr: &actionExpr{
	pos: position{line: 64, col: 18, offset: 1132},
	run: (*parser).callonNaturalLiteral1,
	expr: &oneOrMoreExpr{
	pos: position{line: 64, col: 18, offset: 1132},
	expr: &charClassMatcher{
	pos: position{line: 64, col: 18, offset: 1132},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
},
},
{
	name: "SpaceDeBruijn",
	pos: position{line: 72, col: 1, offset: 1280},
	expr: &actionExpr{
	pos: position{line: 72, col: 17, offset: 1298},
	run: (*parser).callonSpaceDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 72, col: 17, offset: 1298},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 72, col: 17, offset: 1298},
	name: "_",
},
&litMatcher{
	pos: position{line: 72, col: 19, offset: 1300},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 72, col: 23, offset: 1304},
	name: "_",
},
&labeledExpr{
	pos: position{line: 72, col: 25, offset: 1306},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 72, col: 31, offset: 1312},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Identifier",
	pos: position{line: 74, col: 1, offset: 1356},
	expr: &actionExpr{
	pos: position{line: 74, col: 14, offset: 1371},
	run: (*parser).callonIdentifier1,
	expr: &seqExpr{
	pos: position{line: 74, col: 14, offset: 1371},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 74, col: 14, offset: 1371},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 74, col: 19, offset: 1376},
	name: "Label",
},
},
&labeledExpr{
	pos: position{line: 74, col: 25, offset: 1382},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 74, col: 31, offset: 1388},
	expr: &ruleRefExpr{
	pos: position{line: 74, col: 31, offset: 1388},
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
	pos: position{line: 82, col: 1, offset: 1564},
	expr: &actionExpr{
	pos: position{line: 82, col: 28, offset: 1593},
	run: (*parser).callonIdentifierReservedPrefix1,
	expr: &seqExpr{
	pos: position{line: 82, col: 28, offset: 1593},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 82, col: 28, offset: 1593},
	label: "reserved",
	expr: &ruleRefExpr{
	pos: position{line: 82, col: 37, offset: 1602},
	name: "Reserved",
},
},
&labeledExpr{
	pos: position{line: 82, col: 46, offset: 1611},
	label: "suffix",
	expr: &actionExpr{
	pos: position{line: 82, col: 54, offset: 1619},
	run: (*parser).callonIdentifierReservedPrefix6,
	expr: &oneOrMoreExpr{
	pos: position{line: 82, col: 54, offset: 1619},
	expr: &charClassMatcher{
	pos: position{line: 82, col: 54, offset: 1619},
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
	pos: position{line: 82, col: 102, offset: 1667},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 82, col: 108, offset: 1673},
	expr: &ruleRefExpr{
	pos: position{line: 82, col: 108, offset: 1673},
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
	pos: position{line: 100, col: 1, offset: 2391},
	expr: &choiceExpr{
	pos: position{line: 101, col: 7, offset: 2412},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 101, col: 7, offset: 2412},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 101, col: 7, offset: 2412},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 101, col: 7, offset: 2412},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 101, col: 14, offset: 2419},
	name: "_",
},
&litMatcher{
	pos: position{line: 101, col: 16, offset: 2421},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 101, col: 20, offset: 2425},
	name: "_",
},
&labeledExpr{
	pos: position{line: 101, col: 22, offset: 2427},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 101, col: 28, offset: 2433},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 101, col: 34, offset: 2439},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 101, col: 36, offset: 2441},
	name: "ColonSpace",
},
&labeledExpr{
	pos: position{line: 101, col: 47, offset: 2452},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 101, col: 49, offset: 2454},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 101, col: 60, offset: 2465},
	name: "_",
},
&litMatcher{
	pos: position{line: 101, col: 62, offset: 2467},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 101, col: 66, offset: 2471},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 101, col: 68, offset: 2473},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 101, col: 74, offset: 2479},
	name: "_",
},
&labeledExpr{
	pos: position{line: 101, col: 76, offset: 2481},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 101, col: 81, offset: 2486},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 104, col: 7, offset: 2615},
	run: (*parser).callonExpression21,
	expr: &seqExpr{
	pos: position{line: 104, col: 7, offset: 2615},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 104, col: 7, offset: 2615},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 104, col: 14, offset: 2622},
	name: "_",
},
&litMatcher{
	pos: position{line: 104, col: 16, offset: 2624},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 104, col: 20, offset: 2628},
	name: "_",
},
&labeledExpr{
	pos: position{line: 104, col: 22, offset: 2630},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 104, col: 28, offset: 2636},
	name: "Label",
},
},
&ruleRefExpr{
	pos: position{line: 104, col: 34, offset: 2642},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 104, col: 36, offset: 2644},
	name: "ColonSpace",
},
&labeledExpr{
	pos: position{line: 104, col: 47, offset: 2655},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 104, col: 49, offset: 2657},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 104, col: 60, offset: 2668},
	name: "_",
},
&litMatcher{
	pos: position{line: 104, col: 62, offset: 2670},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 104, col: 66, offset: 2674},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 104, col: 68, offset: 2676},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 104, col: 74, offset: 2682},
	name: "_",
},
&labeledExpr{
	pos: position{line: 104, col: 76, offset: 2684},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 104, col: 81, offset: 2689},
	name: "Expression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 107, col: 7, offset: 2810},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 109, col: 1, offset: 2831},
	expr: &ruleRefExpr{
	pos: position{line: 109, col: 23, offset: 2855},
	name: "OperatorExpression",
},
},
{
	name: "OperatorExpression",
	pos: position{line: 111, col: 1, offset: 2875},
	expr: &ruleRefExpr{
	pos: position{line: 111, col: 22, offset: 2898},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 113, col: 1, offset: 2919},
	expr: &ruleRefExpr{
	pos: position{line: 114, col: 7, offset: 2985},
	name: "ApplicationExpression",
},
},
{
	name: "AnArg",
	pos: position{line: 116, col: 1, offset: 3008},
	expr: &actionExpr{
	pos: position{line: 116, col: 9, offset: 3016},
	run: (*parser).callonAnArg1,
	expr: &seqExpr{
	pos: position{line: 116, col: 9, offset: 3016},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 116, col: 9, offset: 3016},
	name: "NonemptyWhitespace",
},
&labeledExpr{
	pos: position{line: 116, col: 28, offset: 3035},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 116, col: 30, offset: 3037},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "ApplicationExpression",
	pos: position{line: 118, col: 1, offset: 3073},
	expr: &actionExpr{
	pos: position{line: 118, col: 25, offset: 3099},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 118, col: 25, offset: 3099},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 118, col: 25, offset: 3099},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 118, col: 27, offset: 3101},
	name: "ImportExpression",
},
},
&labeledExpr{
	pos: position{line: 118, col: 44, offset: 3118},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 118, col: 49, offset: 3123},
	expr: &ruleRefExpr{
	pos: position{line: 118, col: 49, offset: 3123},
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
	pos: position{line: 127, col: 1, offset: 3354},
	expr: &ruleRefExpr{
	pos: position{line: 127, col: 20, offset: 3375},
	name: "SelectorExpression",
},
},
{
	name: "SelectorExpression",
	pos: position{line: 129, col: 1, offset: 3395},
	expr: &ruleRefExpr{
	pos: position{line: 129, col: 22, offset: 3418},
	name: "PrimitiveExpression",
},
},
{
	name: "PrimitiveExpression",
	pos: position{line: 131, col: 1, offset: 3439},
	expr: &choiceExpr{
	pos: position{line: 132, col: 7, offset: 3469},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 132, col: 7, offset: 3469},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 133, col: 7, offset: 3490},
	name: "IdentifierReservedPrefix",
},
&ruleRefExpr{
	pos: position{line: 134, col: 7, offset: 3521},
	name: "Reserved",
},
&ruleRefExpr{
	pos: position{line: 135, col: 7, offset: 3536},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 136, col: 7, offset: 3553},
	run: (*parser).callonPrimitiveExpression6,
	expr: &seqExpr{
	pos: position{line: 136, col: 7, offset: 3553},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 136, col: 7, offset: 3553},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 136, col: 11, offset: 3557},
	name: "_",
},
&labeledExpr{
	pos: position{line: 136, col: 13, offset: 3559},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 136, col: 15, offset: 3561},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 136, col: 26, offset: 3572},
	name: "_",
},
&litMatcher{
	pos: position{line: 136, col: 28, offset: 3574},
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
	pos: position{line: 138, col: 1, offset: 3597},
	expr: &notExpr{
	pos: position{line: 138, col: 7, offset: 3605},
	expr: &anyMatcher{
	line: 138, col: 8, offset: 3606,
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

func (c *current) onReserved15() (interface{}, error) {
 return ast.Type, nil 
}

func (p *parser) callonReserved15() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved15()
}

func (c *current) onReserved17() (interface{}, error) {
 return ast.Kind, nil 
}

func (p *parser) callonReserved17() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved17()
}

func (c *current) onReserved19() (interface{}, error) {
 return ast.Sort, nil 
}

func (p *parser) callonReserved19() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved19()
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

func (c *current) onPrimitiveExpression6(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonPrimitiveExpression6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimitiveExpression6(stack["e"])
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

