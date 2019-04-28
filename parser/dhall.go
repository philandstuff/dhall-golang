
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
	name: "Digit",
	pos: position{line: 47, col: 1, offset: 767},
	expr: &charClassMatcher{
	pos: position{line: 47, col: 9, offset: 777},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "HexDig",
	pos: position{line: 49, col: 1, offset: 784},
	expr: &choiceExpr{
	pos: position{line: 49, col: 10, offset: 795},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 49, col: 10, offset: 795},
	name: "Digit",
},
&charClassMatcher{
	pos: position{line: 49, col: 18, offset: 803},
	val: "[a-f]i",
	ranges: []rune{'a','f',},
	ignoreCase: true,
	inverted: false,
},
	},
},
},
{
	name: "SimpleLabelFirstChar",
	pos: position{line: 51, col: 1, offset: 811},
	expr: &charClassMatcher{
	pos: position{line: 51, col: 24, offset: 836},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SimpleLabelNextChar",
	pos: position{line: 52, col: 1, offset: 846},
	expr: &charClassMatcher{
	pos: position{line: 52, col: 23, offset: 870},
	val: "[A-Za-z0-9_/-]",
	chars: []rune{'_','/','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SimpleLabel",
	pos: position{line: 53, col: 1, offset: 885},
	expr: &choiceExpr{
	pos: position{line: 53, col: 15, offset: 901},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 53, col: 15, offset: 901},
	run: (*parser).callonSimpleLabel2,
	expr: &seqExpr{
	pos: position{line: 53, col: 15, offset: 901},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 53, col: 15, offset: 901},
	name: "Keyword",
},
&oneOrMoreExpr{
	pos: position{line: 53, col: 23, offset: 909},
	expr: &ruleRefExpr{
	pos: position{line: 53, col: 23, offset: 909},
	name: "SimpleLabelNextChar",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 54, col: 13, offset: 973},
	run: (*parser).callonSimpleLabel7,
	expr: &seqExpr{
	pos: position{line: 54, col: 13, offset: 973},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 54, col: 13, offset: 973},
	expr: &ruleRefExpr{
	pos: position{line: 54, col: 14, offset: 974},
	name: "Keyword",
},
},
&ruleRefExpr{
	pos: position{line: 54, col: 22, offset: 982},
	name: "SimpleLabelFirstChar",
},
&zeroOrMoreExpr{
	pos: position{line: 54, col: 43, offset: 1003},
	expr: &ruleRefExpr{
	pos: position{line: 54, col: 43, offset: 1003},
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
	pos: position{line: 61, col: 1, offset: 1104},
	expr: &actionExpr{
	pos: position{line: 61, col: 9, offset: 1114},
	run: (*parser).callonLabel1,
	expr: &labeledExpr{
	pos: position{line: 61, col: 9, offset: 1114},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 61, col: 15, offset: 1120},
	name: "SimpleLabel",
},
},
},
},
{
	name: "NonreservedLabel",
	pos: position{line: 63, col: 1, offset: 1155},
	expr: &choiceExpr{
	pos: position{line: 63, col: 20, offset: 1176},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 63, col: 20, offset: 1176},
	run: (*parser).callonNonreservedLabel2,
	expr: &seqExpr{
	pos: position{line: 63, col: 20, offset: 1176},
	exprs: []interface{}{
&andExpr{
	pos: position{line: 63, col: 20, offset: 1176},
	expr: &seqExpr{
	pos: position{line: 63, col: 22, offset: 1178},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 63, col: 22, offset: 1178},
	name: "Reserved",
},
&ruleRefExpr{
	pos: position{line: 63, col: 31, offset: 1187},
	name: "SimpleLabelNextChar",
},
	},
},
},
&labeledExpr{
	pos: position{line: 63, col: 52, offset: 1208},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 63, col: 58, offset: 1214},
	name: "Label",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 64, col: 19, offset: 1260},
	run: (*parser).callonNonreservedLabel10,
	expr: &seqExpr{
	pos: position{line: 64, col: 19, offset: 1260},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 64, col: 19, offset: 1260},
	expr: &ruleRefExpr{
	pos: position{line: 64, col: 20, offset: 1261},
	name: "Reserved",
},
},
&labeledExpr{
	pos: position{line: 64, col: 29, offset: 1270},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 64, col: 35, offset: 1276},
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
	pos: position{line: 66, col: 1, offset: 1305},
	expr: &ruleRefExpr{
	pos: position{line: 66, col: 12, offset: 1318},
	name: "Label",
},
},
{
	name: "DoubleQuoteChunk",
	pos: position{line: 69, col: 1, offset: 1326},
	expr: &choiceExpr{
	pos: position{line: 70, col: 6, offset: 1352},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 70, col: 6, offset: 1352},
	name: "Interpolation",
},
&actionExpr{
	pos: position{line: 71, col: 6, offset: 1371},
	run: (*parser).callonDoubleQuoteChunk3,
	expr: &seqExpr{
	pos: position{line: 71, col: 6, offset: 1371},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 71, col: 6, offset: 1371},
	val: "\\",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 71, col: 11, offset: 1376},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 71, col: 13, offset: 1378},
	name: "DoubleQuoteEscaped",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 72, col: 6, offset: 1420},
	name: "DoubleQuoteChar",
},
	},
},
},
{
	name: "DoubleQuoteEscaped",
	pos: position{line: 74, col: 1, offset: 1437},
	expr: &choiceExpr{
	pos: position{line: 75, col: 8, offset: 1467},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 75, col: 8, offset: 1467},
	val: "\"",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 76, col: 8, offset: 1478},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 77, col: 8, offset: 1489},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 78, col: 8, offset: 1501},
	val: "/",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 79, col: 8, offset: 1512},
	run: (*parser).callonDoubleQuoteEscaped6,
	expr: &litMatcher{
	pos: position{line: 79, col: 8, offset: 1512},
	val: "b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 80, col: 8, offset: 1552},
	run: (*parser).callonDoubleQuoteEscaped8,
	expr: &litMatcher{
	pos: position{line: 80, col: 8, offset: 1552},
	val: "f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 81, col: 8, offset: 1592},
	run: (*parser).callonDoubleQuoteEscaped10,
	expr: &litMatcher{
	pos: position{line: 81, col: 8, offset: 1592},
	val: "n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 82, col: 8, offset: 1632},
	run: (*parser).callonDoubleQuoteEscaped12,
	expr: &litMatcher{
	pos: position{line: 82, col: 8, offset: 1632},
	val: "r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 83, col: 8, offset: 1672},
	run: (*parser).callonDoubleQuoteEscaped14,
	expr: &litMatcher{
	pos: position{line: 83, col: 8, offset: 1672},
	val: "t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 84, col: 8, offset: 1712},
	run: (*parser).callonDoubleQuoteEscaped16,
	expr: &seqExpr{
	pos: position{line: 84, col: 8, offset: 1712},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 84, col: 8, offset: 1712},
	val: "u",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 84, col: 12, offset: 1716},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 84, col: 19, offset: 1723},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 84, col: 26, offset: 1730},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 84, col: 33, offset: 1737},
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
	pos: position{line: 89, col: 1, offset: 1869},
	expr: &choiceExpr{
	pos: position{line: 90, col: 6, offset: 1894},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 90, col: 6, offset: 1894},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 91, col: 6, offset: 1911},
	val: "[\\x23-\\x5b]",
	ranges: []rune{'#','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 92, col: 6, offset: 1928},
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
	pos: position{line: 94, col: 1, offset: 1947},
	expr: &actionExpr{
	pos: position{line: 94, col: 22, offset: 1970},
	run: (*parser).callonDoubleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 94, col: 22, offset: 1970},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 94, col: 22, offset: 1970},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 94, col: 26, offset: 1974},
	label: "chunks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 94, col: 33, offset: 1981},
	expr: &ruleRefExpr{
	pos: position{line: 94, col: 33, offset: 1981},
	name: "DoubleQuoteChunk",
},
},
},
&litMatcher{
	pos: position{line: 94, col: 51, offset: 1999},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "SingleQuoteContinue",
	pos: position{line: 111, col: 1, offset: 2467},
	expr: &choiceExpr{
	pos: position{line: 112, col: 7, offset: 2497},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 112, col: 7, offset: 2497},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 112, col: 7, offset: 2497},
	name: "Interpolation",
},
&ruleRefExpr{
	pos: position{line: 112, col: 21, offset: 2511},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 113, col: 7, offset: 2537},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 113, col: 7, offset: 2537},
	name: "EscapedQuotePair",
},
&ruleRefExpr{
	pos: position{line: 113, col: 24, offset: 2554},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 114, col: 7, offset: 2580},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 114, col: 7, offset: 2580},
	name: "EscapedInterpolation",
},
&ruleRefExpr{
	pos: position{line: 114, col: 28, offset: 2601},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 115, col: 7, offset: 2627},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 115, col: 7, offset: 2627},
	name: "SingleQuoteChar",
},
&ruleRefExpr{
	pos: position{line: 115, col: 23, offset: 2643},
	name: "SingleQuoteContinue",
},
	},
},
&litMatcher{
	pos: position{line: 116, col: 7, offset: 2669},
	val: "''",
	ignoreCase: false,
},
	},
},
},
{
	name: "EscapedQuotePair",
	pos: position{line: 118, col: 1, offset: 2675},
	expr: &actionExpr{
	pos: position{line: 118, col: 20, offset: 2696},
	run: (*parser).callonEscapedQuotePair1,
	expr: &litMatcher{
	pos: position{line: 118, col: 20, offset: 2696},
	val: "'''",
	ignoreCase: false,
},
},
},
{
	name: "EscapedInterpolation",
	pos: position{line: 122, col: 1, offset: 2831},
	expr: &actionExpr{
	pos: position{line: 122, col: 24, offset: 2856},
	run: (*parser).callonEscapedInterpolation1,
	expr: &litMatcher{
	pos: position{line: 122, col: 24, offset: 2856},
	val: "''${",
	ignoreCase: false,
},
},
},
{
	name: "SingleQuoteChar",
	pos: position{line: 124, col: 1, offset: 2898},
	expr: &choiceExpr{
	pos: position{line: 125, col: 6, offset: 2923},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 125, col: 6, offset: 2923},
	val: "[\\x20-\\U0010ffff]",
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 126, col: 6, offset: 2946},
	val: "\t",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 127, col: 6, offset: 2956},
	val: "\n",
	ignoreCase: false,
},
	},
},
},
{
	name: "SingleQuoteLiteral",
	pos: position{line: 129, col: 1, offset: 2962},
	expr: &actionExpr{
	pos: position{line: 129, col: 22, offset: 2985},
	run: (*parser).callonSingleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 129, col: 22, offset: 2985},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 129, col: 22, offset: 2985},
	val: "''",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 129, col: 27, offset: 2990},
	name: "EOL",
},
&labeledExpr{
	pos: position{line: 129, col: 31, offset: 2994},
	label: "content",
	expr: &ruleRefExpr{
	pos: position{line: 129, col: 39, offset: 3002},
	name: "SingleQuoteContinue",
},
},
	},
},
},
},
{
	name: "Interpolation",
	pos: position{line: 147, col: 1, offset: 3525},
	expr: &actionExpr{
	pos: position{line: 147, col: 17, offset: 3543},
	run: (*parser).callonInterpolation1,
	expr: &seqExpr{
	pos: position{line: 147, col: 17, offset: 3543},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 147, col: 17, offset: 3543},
	val: "${",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 147, col: 22, offset: 3548},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 147, col: 24, offset: 3550},
	name: "CompleteExpression",
},
},
&litMatcher{
	pos: position{line: 147, col: 43, offset: 3569},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TextLiteral",
	pos: position{line: 149, col: 1, offset: 3592},
	expr: &choiceExpr{
	pos: position{line: 149, col: 15, offset: 3608},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 149, col: 15, offset: 3608},
	name: "DoubleQuoteLiteral",
},
&ruleRefExpr{
	pos: position{line: 149, col: 36, offset: 3629},
	name: "SingleQuoteLiteral",
},
	},
},
},
{
	name: "Reserved",
	pos: position{line: 152, col: 1, offset: 3734},
	expr: &choiceExpr{
	pos: position{line: 153, col: 5, offset: 3751},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 153, col: 5, offset: 3751},
	run: (*parser).callonReserved2,
	expr: &litMatcher{
	pos: position{line: 153, col: 5, offset: 3751},
	val: "Natural/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 154, col: 5, offset: 3829},
	run: (*parser).callonReserved4,
	expr: &litMatcher{
	pos: position{line: 154, col: 5, offset: 3829},
	val: "Natural/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 155, col: 5, offset: 3905},
	run: (*parser).callonReserved6,
	expr: &litMatcher{
	pos: position{line: 155, col: 5, offset: 3905},
	val: "Natural/isZero",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 156, col: 5, offset: 3985},
	run: (*parser).callonReserved8,
	expr: &litMatcher{
	pos: position{line: 156, col: 5, offset: 3985},
	val: "Natural/even",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 157, col: 5, offset: 4061},
	run: (*parser).callonReserved10,
	expr: &litMatcher{
	pos: position{line: 157, col: 5, offset: 4061},
	val: "Natural/odd",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 158, col: 5, offset: 4135},
	run: (*parser).callonReserved12,
	expr: &litMatcher{
	pos: position{line: 158, col: 5, offset: 4135},
	val: "Natural/toInteger",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 159, col: 5, offset: 4221},
	run: (*parser).callonReserved14,
	expr: &litMatcher{
	pos: position{line: 159, col: 5, offset: 4221},
	val: "Natural/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 160, col: 5, offset: 4297},
	run: (*parser).callonReserved16,
	expr: &litMatcher{
	pos: position{line: 160, col: 5, offset: 4297},
	val: "Integer/toDouble",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 161, col: 5, offset: 4381},
	run: (*parser).callonReserved18,
	expr: &litMatcher{
	pos: position{line: 161, col: 5, offset: 4381},
	val: "Integer/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 162, col: 5, offset: 4457},
	run: (*parser).callonReserved20,
	expr: &litMatcher{
	pos: position{line: 162, col: 5, offset: 4457},
	val: "Double/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 163, col: 5, offset: 4531},
	run: (*parser).callonReserved22,
	expr: &litMatcher{
	pos: position{line: 163, col: 5, offset: 4531},
	val: "List/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 164, col: 5, offset: 4603},
	run: (*parser).callonReserved24,
	expr: &litMatcher{
	pos: position{line: 164, col: 5, offset: 4603},
	val: "List/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 165, col: 5, offset: 4673},
	run: (*parser).callonReserved26,
	expr: &litMatcher{
	pos: position{line: 165, col: 5, offset: 4673},
	val: "List/length",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 166, col: 5, offset: 4747},
	run: (*parser).callonReserved28,
	expr: &litMatcher{
	pos: position{line: 166, col: 5, offset: 4747},
	val: "List/head",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 167, col: 5, offset: 4817},
	run: (*parser).callonReserved30,
	expr: &litMatcher{
	pos: position{line: 167, col: 5, offset: 4817},
	val: "List/last",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 168, col: 5, offset: 4887},
	run: (*parser).callonReserved32,
	expr: &litMatcher{
	pos: position{line: 168, col: 5, offset: 4887},
	val: "List/indexed",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 169, col: 5, offset: 4963},
	run: (*parser).callonReserved34,
	expr: &litMatcher{
	pos: position{line: 169, col: 5, offset: 4963},
	val: "List/reverse",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 170, col: 5, offset: 5039},
	run: (*parser).callonReserved36,
	expr: &litMatcher{
	pos: position{line: 170, col: 5, offset: 5039},
	val: "Optional/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 171, col: 5, offset: 5119},
	run: (*parser).callonReserved38,
	expr: &litMatcher{
	pos: position{line: 171, col: 5, offset: 5119},
	val: "Optional/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 172, col: 5, offset: 5197},
	run: (*parser).callonReserved40,
	expr: &litMatcher{
	pos: position{line: 172, col: 5, offset: 5197},
	val: "Text/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 173, col: 5, offset: 5267},
	run: (*parser).callonReserved42,
	expr: &litMatcher{
	pos: position{line: 173, col: 5, offset: 5267},
	val: "Bool",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 174, col: 5, offset: 5299},
	run: (*parser).callonReserved44,
	expr: &litMatcher{
	pos: position{line: 174, col: 5, offset: 5299},
	val: "True",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 175, col: 5, offset: 5331},
	run: (*parser).callonReserved46,
	expr: &litMatcher{
	pos: position{line: 175, col: 5, offset: 5331},
	val: "False",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 176, col: 5, offset: 5365},
	run: (*parser).callonReserved48,
	expr: &litMatcher{
	pos: position{line: 176, col: 5, offset: 5365},
	val: "Optional",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 177, col: 5, offset: 5405},
	run: (*parser).callonReserved50,
	expr: &litMatcher{
	pos: position{line: 177, col: 5, offset: 5405},
	val: "Natural",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 178, col: 5, offset: 5443},
	run: (*parser).callonReserved52,
	expr: &litMatcher{
	pos: position{line: 178, col: 5, offset: 5443},
	val: "Integer",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 179, col: 5, offset: 5481},
	run: (*parser).callonReserved54,
	expr: &litMatcher{
	pos: position{line: 179, col: 5, offset: 5481},
	val: "Double",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 180, col: 5, offset: 5517},
	run: (*parser).callonReserved56,
	expr: &litMatcher{
	pos: position{line: 180, col: 5, offset: 5517},
	val: "Text",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 181, col: 5, offset: 5549},
	run: (*parser).callonReserved58,
	expr: &litMatcher{
	pos: position{line: 181, col: 5, offset: 5549},
	val: "List",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 182, col: 5, offset: 5581},
	run: (*parser).callonReserved60,
	expr: &litMatcher{
	pos: position{line: 182, col: 5, offset: 5581},
	val: "None",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 183, col: 5, offset: 5613},
	run: (*parser).callonReserved62,
	expr: &litMatcher{
	pos: position{line: 183, col: 5, offset: 5613},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 184, col: 5, offset: 5645},
	run: (*parser).callonReserved64,
	expr: &litMatcher{
	pos: position{line: 184, col: 5, offset: 5645},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 185, col: 5, offset: 5677},
	run: (*parser).callonReserved66,
	expr: &litMatcher{
	pos: position{line: 185, col: 5, offset: 5677},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "If",
	pos: position{line: 187, col: 1, offset: 5706},
	expr: &litMatcher{
	pos: position{line: 187, col: 6, offset: 5713},
	val: "if",
	ignoreCase: false,
},
},
{
	name: "Then",
	pos: position{line: 188, col: 1, offset: 5718},
	expr: &litMatcher{
	pos: position{line: 188, col: 8, offset: 5727},
	val: "then",
	ignoreCase: false,
},
},
{
	name: "Else",
	pos: position{line: 189, col: 1, offset: 5734},
	expr: &litMatcher{
	pos: position{line: 189, col: 8, offset: 5743},
	val: "else",
	ignoreCase: false,
},
},
{
	name: "Let",
	pos: position{line: 190, col: 1, offset: 5750},
	expr: &litMatcher{
	pos: position{line: 190, col: 7, offset: 5758},
	val: "let",
	ignoreCase: false,
},
},
{
	name: "In",
	pos: position{line: 191, col: 1, offset: 5764},
	expr: &litMatcher{
	pos: position{line: 191, col: 6, offset: 5771},
	val: "in",
	ignoreCase: false,
},
},
{
	name: "As",
	pos: position{line: 192, col: 1, offset: 5776},
	expr: &litMatcher{
	pos: position{line: 192, col: 6, offset: 5783},
	val: "as",
	ignoreCase: false,
},
},
{
	name: "Using",
	pos: position{line: 193, col: 1, offset: 5788},
	expr: &litMatcher{
	pos: position{line: 193, col: 9, offset: 5798},
	val: "using",
	ignoreCase: false,
},
},
{
	name: "Merge",
	pos: position{line: 194, col: 1, offset: 5806},
	expr: &litMatcher{
	pos: position{line: 194, col: 9, offset: 5816},
	val: "merge",
	ignoreCase: false,
},
},
{
	name: "Missing",
	pos: position{line: 195, col: 1, offset: 5824},
	expr: &actionExpr{
	pos: position{line: 195, col: 11, offset: 5836},
	run: (*parser).callonMissing1,
	expr: &litMatcher{
	pos: position{line: 195, col: 11, offset: 5836},
	val: "missing",
	ignoreCase: false,
},
},
},
{
	name: "True",
	pos: position{line: 196, col: 1, offset: 5872},
	expr: &litMatcher{
	pos: position{line: 196, col: 8, offset: 5881},
	val: "True",
	ignoreCase: false,
},
},
{
	name: "False",
	pos: position{line: 197, col: 1, offset: 5888},
	expr: &litMatcher{
	pos: position{line: 197, col: 9, offset: 5898},
	val: "False",
	ignoreCase: false,
},
},
{
	name: "Infinity",
	pos: position{line: 198, col: 1, offset: 5906},
	expr: &litMatcher{
	pos: position{line: 198, col: 12, offset: 5919},
	val: "Infinity",
	ignoreCase: false,
},
},
{
	name: "NaN",
	pos: position{line: 199, col: 1, offset: 5930},
	expr: &litMatcher{
	pos: position{line: 199, col: 7, offset: 5938},
	val: "NaN",
	ignoreCase: false,
},
},
{
	name: "Some",
	pos: position{line: 200, col: 1, offset: 5944},
	expr: &litMatcher{
	pos: position{line: 200, col: 8, offset: 5953},
	val: "Some",
	ignoreCase: false,
},
},
{
	name: "Keyword",
	pos: position{line: 202, col: 1, offset: 5961},
	expr: &choiceExpr{
	pos: position{line: 203, col: 5, offset: 5977},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 203, col: 5, offset: 5977},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 203, col: 10, offset: 5982},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 203, col: 17, offset: 5989},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 204, col: 5, offset: 5998},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 204, col: 11, offset: 6004},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 205, col: 5, offset: 6011},
	name: "Using",
},
&ruleRefExpr{
	pos: position{line: 205, col: 13, offset: 6019},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 205, col: 23, offset: 6029},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 206, col: 5, offset: 6036},
	name: "True",
},
&ruleRefExpr{
	pos: position{line: 206, col: 12, offset: 6043},
	name: "False",
},
&ruleRefExpr{
	pos: position{line: 207, col: 5, offset: 6053},
	name: "Infinity",
},
&ruleRefExpr{
	pos: position{line: 207, col: 16, offset: 6064},
	name: "NaN",
},
&ruleRefExpr{
	pos: position{line: 208, col: 5, offset: 6072},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 208, col: 13, offset: 6080},
	name: "Some",
},
	},
},
},
{
	name: "Optional",
	pos: position{line: 210, col: 1, offset: 6086},
	expr: &litMatcher{
	pos: position{line: 210, col: 12, offset: 6099},
	val: "Optional",
	ignoreCase: false,
},
},
{
	name: "Text",
	pos: position{line: 211, col: 1, offset: 6110},
	expr: &litMatcher{
	pos: position{line: 211, col: 8, offset: 6119},
	val: "Text",
	ignoreCase: false,
},
},
{
	name: "List",
	pos: position{line: 212, col: 1, offset: 6126},
	expr: &litMatcher{
	pos: position{line: 212, col: 8, offset: 6135},
	val: "List",
	ignoreCase: false,
},
},
{
	name: "Lambda",
	pos: position{line: 214, col: 1, offset: 6143},
	expr: &choiceExpr{
	pos: position{line: 214, col: 11, offset: 6155},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 214, col: 11, offset: 6155},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 214, col: 18, offset: 6162},
	val: "λ",
	ignoreCase: false,
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 215, col: 1, offset: 6168},
	expr: &choiceExpr{
	pos: position{line: 215, col: 11, offset: 6180},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 215, col: 11, offset: 6180},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 215, col: 22, offset: 6191},
	val: "∀",
	ignoreCase: false,
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 216, col: 1, offset: 6198},
	expr: &choiceExpr{
	pos: position{line: 216, col: 10, offset: 6209},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 216, col: 10, offset: 6209},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 216, col: 17, offset: 6216},
	val: "→",
	ignoreCase: false,
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 218, col: 1, offset: 6224},
	expr: &seqExpr{
	pos: position{line: 218, col: 12, offset: 6237},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 218, col: 12, offset: 6237},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 218, col: 17, offset: 6242},
	expr: &charClassMatcher{
	pos: position{line: 218, col: 17, offset: 6242},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 218, col: 23, offset: 6248},
	expr: &ruleRefExpr{
	pos: position{line: 218, col: 23, offset: 6248},
	name: "Digit",
},
},
	},
},
},
{
	name: "NumericDoubleLiteral",
	pos: position{line: 220, col: 1, offset: 6256},
	expr: &actionExpr{
	pos: position{line: 220, col: 24, offset: 6281},
	run: (*parser).callonNumericDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 220, col: 24, offset: 6281},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 220, col: 24, offset: 6281},
	expr: &charClassMatcher{
	pos: position{line: 220, col: 24, offset: 6281},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 220, col: 30, offset: 6287},
	expr: &ruleRefExpr{
	pos: position{line: 220, col: 30, offset: 6287},
	name: "Digit",
},
},
&choiceExpr{
	pos: position{line: 220, col: 39, offset: 6296},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 220, col: 39, offset: 6296},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 220, col: 39, offset: 6296},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 220, col: 43, offset: 6300},
	expr: &ruleRefExpr{
	pos: position{line: 220, col: 43, offset: 6300},
	name: "Digit",
},
},
&zeroOrOneExpr{
	pos: position{line: 220, col: 50, offset: 6307},
	expr: &ruleRefExpr{
	pos: position{line: 220, col: 50, offset: 6307},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 220, col: 62, offset: 6319},
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
	pos: position{line: 228, col: 1, offset: 6475},
	expr: &choiceExpr{
	pos: position{line: 228, col: 17, offset: 6493},
	alternatives: []interface{}{
&labeledExpr{
	pos: position{line: 228, col: 17, offset: 6493},
	label: "d",
	expr: &ruleRefExpr{
	pos: position{line: 228, col: 19, offset: 6495},
	name: "NumericDoubleLiteral",
},
},
&actionExpr{
	pos: position{line: 229, col: 5, offset: 6520},
	run: (*parser).callonDoubleLiteral4,
	expr: &ruleRefExpr{
	pos: position{line: 229, col: 5, offset: 6520},
	name: "Infinity",
},
},
&actionExpr{
	pos: position{line: 230, col: 5, offset: 6572},
	run: (*parser).callonDoubleLiteral6,
	expr: &seqExpr{
	pos: position{line: 230, col: 5, offset: 6572},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 230, col: 5, offset: 6572},
	val: "-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 230, col: 9, offset: 6576},
	name: "Infinity",
},
	},
},
},
&actionExpr{
	pos: position{line: 231, col: 5, offset: 6629},
	run: (*parser).callonDoubleLiteral10,
	expr: &ruleRefExpr{
	pos: position{line: 231, col: 5, offset: 6629},
	name: "NaN",
},
},
	},
},
},
{
	name: "NaturalLiteral",
	pos: position{line: 233, col: 1, offset: 6672},
	expr: &actionExpr{
	pos: position{line: 233, col: 18, offset: 6691},
	run: (*parser).callonNaturalLiteral1,
	expr: &oneOrMoreExpr{
	pos: position{line: 233, col: 18, offset: 6691},
	expr: &ruleRefExpr{
	pos: position{line: 233, col: 18, offset: 6691},
	name: "Digit",
},
},
},
},
{
	name: "IntegerLiteral",
	pos: position{line: 238, col: 1, offset: 6780},
	expr: &actionExpr{
	pos: position{line: 238, col: 18, offset: 6799},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 238, col: 18, offset: 6799},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 238, col: 18, offset: 6799},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 238, col: 22, offset: 6803},
	name: "NaturalLiteral",
},
	},
},
},
},
{
	name: "DeBruijn",
	pos: position{line: 246, col: 1, offset: 6955},
	expr: &actionExpr{
	pos: position{line: 246, col: 12, offset: 6968},
	run: (*parser).callonDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 246, col: 12, offset: 6968},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 246, col: 12, offset: 6968},
	name: "_",
},
&litMatcher{
	pos: position{line: 246, col: 14, offset: 6970},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 246, col: 18, offset: 6974},
	name: "_",
},
&labeledExpr{
	pos: position{line: 246, col: 20, offset: 6976},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 246, col: 26, offset: 6982},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Variable",
	pos: position{line: 248, col: 1, offset: 7038},
	expr: &actionExpr{
	pos: position{line: 248, col: 12, offset: 7051},
	run: (*parser).callonVariable1,
	expr: &seqExpr{
	pos: position{line: 248, col: 12, offset: 7051},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 248, col: 12, offset: 7051},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 248, col: 17, offset: 7056},
	name: "NonreservedLabel",
},
},
&labeledExpr{
	pos: position{line: 248, col: 34, offset: 7073},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 248, col: 40, offset: 7079},
	expr: &ruleRefExpr{
	pos: position{line: 248, col: 40, offset: 7079},
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
	pos: position{line: 256, col: 1, offset: 7242},
	expr: &choiceExpr{
	pos: position{line: 256, col: 14, offset: 7257},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 256, col: 14, offset: 7257},
	name: "Variable",
},
&ruleRefExpr{
	pos: position{line: 256, col: 25, offset: 7268},
	name: "Reserved",
},
	},
},
},
{
	name: "PathCharacter",
	pos: position{line: 258, col: 1, offset: 7278},
	expr: &choiceExpr{
	pos: position{line: 259, col: 6, offset: 7301},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 259, col: 6, offset: 7301},
	val: "!",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 260, col: 6, offset: 7313},
	val: "[\\x24-\\x27]",
	ranges: []rune{'$','\'',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 261, col: 6, offset: 7330},
	val: "[\\x2a-\\x2b]",
	ranges: []rune{'*','+',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 262, col: 6, offset: 7347},
	val: "[\\x2d-\\x2e]",
	ranges: []rune{'-','.',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 263, col: 6, offset: 7364},
	val: "[\\x30-\\x3b]",
	ranges: []rune{'0',';',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 264, col: 6, offset: 7381},
	val: "=",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 265, col: 6, offset: 7393},
	val: "[\\x40-\\x5a]",
	ranges: []rune{'@','Z',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 266, col: 6, offset: 7410},
	val: "[\\x5e-\\x7a]",
	ranges: []rune{'^','z',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 267, col: 6, offset: 7427},
	val: "|",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 268, col: 6, offset: 7439},
	val: "~",
	ignoreCase: false,
},
	},
},
},
{
	name: "UnquotedPathComponent",
	pos: position{line: 270, col: 1, offset: 7447},
	expr: &actionExpr{
	pos: position{line: 270, col: 25, offset: 7473},
	run: (*parser).callonUnquotedPathComponent1,
	expr: &oneOrMoreExpr{
	pos: position{line: 270, col: 25, offset: 7473},
	expr: &ruleRefExpr{
	pos: position{line: 270, col: 25, offset: 7473},
	name: "PathCharacter",
},
},
},
},
{
	name: "PathComponent",
	pos: position{line: 272, col: 1, offset: 7520},
	expr: &actionExpr{
	pos: position{line: 272, col: 17, offset: 7538},
	run: (*parser).callonPathComponent1,
	expr: &seqExpr{
	pos: position{line: 272, col: 17, offset: 7538},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 272, col: 17, offset: 7538},
	val: "/",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 272, col: 21, offset: 7542},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 272, col: 23, offset: 7544},
	name: "UnquotedPathComponent",
},
},
	},
},
},
},
{
	name: "Path",
	pos: position{line: 274, col: 1, offset: 7585},
	expr: &actionExpr{
	pos: position{line: 274, col: 8, offset: 7594},
	run: (*parser).callonPath1,
	expr: &labeledExpr{
	pos: position{line: 274, col: 8, offset: 7594},
	label: "cs",
	expr: &oneOrMoreExpr{
	pos: position{line: 274, col: 11, offset: 7597},
	expr: &ruleRefExpr{
	pos: position{line: 274, col: 11, offset: 7597},
	name: "PathComponent",
},
},
},
},
},
{
	name: "Local",
	pos: position{line: 283, col: 1, offset: 7871},
	expr: &choiceExpr{
	pos: position{line: 283, col: 9, offset: 7881},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 283, col: 9, offset: 7881},
	name: "ParentPath",
},
&ruleRefExpr{
	pos: position{line: 283, col: 22, offset: 7894},
	name: "HerePath",
},
&ruleRefExpr{
	pos: position{line: 283, col: 33, offset: 7905},
	name: "HomePath",
},
&ruleRefExpr{
	pos: position{line: 283, col: 44, offset: 7916},
	name: "AbsolutePath",
},
	},
},
},
{
	name: "ParentPath",
	pos: position{line: 285, col: 1, offset: 7930},
	expr: &actionExpr{
	pos: position{line: 285, col: 14, offset: 7945},
	run: (*parser).callonParentPath1,
	expr: &seqExpr{
	pos: position{line: 285, col: 14, offset: 7945},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 285, col: 14, offset: 7945},
	val: "..",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 285, col: 19, offset: 7950},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 285, col: 21, offset: 7952},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HerePath",
	pos: position{line: 286, col: 1, offset: 8008},
	expr: &actionExpr{
	pos: position{line: 286, col: 12, offset: 8021},
	run: (*parser).callonHerePath1,
	expr: &seqExpr{
	pos: position{line: 286, col: 12, offset: 8021},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 286, col: 12, offset: 8021},
	val: ".",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 286, col: 16, offset: 8025},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 286, col: 18, offset: 8027},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HomePath",
	pos: position{line: 287, col: 1, offset: 8066},
	expr: &actionExpr{
	pos: position{line: 287, col: 12, offset: 8079},
	run: (*parser).callonHomePath1,
	expr: &seqExpr{
	pos: position{line: 287, col: 12, offset: 8079},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 287, col: 12, offset: 8079},
	val: "~",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 287, col: 16, offset: 8083},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 287, col: 18, offset: 8085},
	name: "Path",
},
},
	},
},
},
},
{
	name: "AbsolutePath",
	pos: position{line: 288, col: 1, offset: 8140},
	expr: &actionExpr{
	pos: position{line: 288, col: 16, offset: 8157},
	run: (*parser).callonAbsolutePath1,
	expr: &labeledExpr{
	pos: position{line: 288, col: 16, offset: 8157},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 288, col: 18, offset: 8159},
	name: "Path",
},
},
},
},
{
	name: "Scheme",
	pos: position{line: 290, col: 1, offset: 8215},
	expr: &seqExpr{
	pos: position{line: 290, col: 10, offset: 8226},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 290, col: 10, offset: 8226},
	val: "http",
	ignoreCase: false,
},
&zeroOrOneExpr{
	pos: position{line: 290, col: 17, offset: 8233},
	expr: &litMatcher{
	pos: position{line: 290, col: 17, offset: 8233},
	val: "s",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "HttpRaw",
	pos: position{line: 292, col: 1, offset: 8239},
	expr: &actionExpr{
	pos: position{line: 292, col: 11, offset: 8251},
	run: (*parser).callonHttpRaw1,
	expr: &seqExpr{
	pos: position{line: 292, col: 11, offset: 8251},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 292, col: 11, offset: 8251},
	name: "Scheme",
},
&litMatcher{
	pos: position{line: 292, col: 18, offset: 8258},
	val: "://",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 292, col: 24, offset: 8264},
	name: "Authority",
},
&ruleRefExpr{
	pos: position{line: 292, col: 34, offset: 8274},
	name: "Path",
},
&zeroOrOneExpr{
	pos: position{line: 292, col: 39, offset: 8279},
	expr: &seqExpr{
	pos: position{line: 292, col: 41, offset: 8281},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 292, col: 41, offset: 8281},
	val: "?",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 292, col: 45, offset: 8285},
	name: "Query",
},
	},
},
},
	},
},
},
},
{
	name: "Authority",
	pos: position{line: 294, col: 1, offset: 8326},
	expr: &seqExpr{
	pos: position{line: 294, col: 13, offset: 8340},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 294, col: 13, offset: 8340},
	name: "Host",
},
&zeroOrOneExpr{
	pos: position{line: 294, col: 18, offset: 8345},
	expr: &seqExpr{
	pos: position{line: 294, col: 19, offset: 8346},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 294, col: 19, offset: 8346},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 294, col: 23, offset: 8350},
	name: "Port",
},
	},
},
},
	},
},
},
{
	name: "Host",
	pos: position{line: 296, col: 1, offset: 8358},
	expr: &ruleRefExpr{
	pos: position{line: 296, col: 8, offset: 8367},
	name: "RegName",
},
},
{
	name: "Port",
	pos: position{line: 298, col: 1, offset: 8376},
	expr: &zeroOrMoreExpr{
	pos: position{line: 298, col: 8, offset: 8385},
	expr: &ruleRefExpr{
	pos: position{line: 298, col: 8, offset: 8385},
	name: "Digit",
},
},
},
{
	name: "RegName",
	pos: position{line: 300, col: 1, offset: 8393},
	expr: &zeroOrMoreExpr{
	pos: position{line: 300, col: 11, offset: 8405},
	expr: &choiceExpr{
	pos: position{line: 300, col: 12, offset: 8406},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 300, col: 12, offset: 8406},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 300, col: 25, offset: 8419},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 300, col: 38, offset: 8432},
	name: "SubDelims",
},
	},
},
},
},
{
	name: "PChar",
	pos: position{line: 302, col: 1, offset: 8445},
	expr: &choiceExpr{
	pos: position{line: 302, col: 9, offset: 8455},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 302, col: 9, offset: 8455},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 302, col: 22, offset: 8468},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 302, col: 35, offset: 8481},
	name: "SubDelims",
},
&charClassMatcher{
	pos: position{line: 302, col: 47, offset: 8493},
	val: "[:@]",
	chars: []rune{':','@',},
	ignoreCase: false,
	inverted: false,
},
	},
},
},
{
	name: "Query",
	pos: position{line: 304, col: 1, offset: 8499},
	expr: &zeroOrMoreExpr{
	pos: position{line: 304, col: 9, offset: 8509},
	expr: &choiceExpr{
	pos: position{line: 304, col: 10, offset: 8510},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 304, col: 10, offset: 8510},
	name: "PChar",
},
&charClassMatcher{
	pos: position{line: 304, col: 18, offset: 8518},
	val: "[/?]",
	chars: []rune{'/','?',},
	ignoreCase: false,
	inverted: false,
},
	},
},
},
},
{
	name: "PctEncoded",
	pos: position{line: 306, col: 1, offset: 8526},
	expr: &seqExpr{
	pos: position{line: 306, col: 14, offset: 8541},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 306, col: 14, offset: 8541},
	val: "%",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 306, col: 18, offset: 8545},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 306, col: 25, offset: 8552},
	name: "HexDig",
},
	},
},
},
{
	name: "Unreserved",
	pos: position{line: 308, col: 1, offset: 8560},
	expr: &charClassMatcher{
	pos: position{line: 308, col: 14, offset: 8575},
	val: "[A-Za-z0-9._~-]",
	chars: []rune{'.','_','~','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SubDelims",
	pos: position{line: 310, col: 1, offset: 8592},
	expr: &choiceExpr{
	pos: position{line: 310, col: 13, offset: 8606},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 310, col: 13, offset: 8606},
	val: "!",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 310, col: 19, offset: 8612},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 310, col: 25, offset: 8618},
	val: "&",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 310, col: 31, offset: 8624},
	val: "'",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 310, col: 37, offset: 8630},
	val: "(",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 310, col: 43, offset: 8636},
	val: ")",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 310, col: 49, offset: 8642},
	val: "*",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 310, col: 55, offset: 8648},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 310, col: 61, offset: 8654},
	val: ",",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 310, col: 67, offset: 8660},
	val: ";",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 310, col: 73, offset: 8666},
	val: "=",
	ignoreCase: false,
},
	},
},
},
{
	name: "Http",
	pos: position{line: 312, col: 1, offset: 8671},
	expr: &actionExpr{
	pos: position{line: 312, col: 8, offset: 8680},
	run: (*parser).callonHttp1,
	expr: &labeledExpr{
	pos: position{line: 312, col: 8, offset: 8680},
	label: "url",
	expr: &ruleRefExpr{
	pos: position{line: 312, col: 12, offset: 8684},
	name: "HttpRaw",
},
},
},
},
{
	name: "Env",
	pos: position{line: 314, col: 1, offset: 8730},
	expr: &actionExpr{
	pos: position{line: 314, col: 7, offset: 8738},
	run: (*parser).callonEnv1,
	expr: &seqExpr{
	pos: position{line: 314, col: 7, offset: 8738},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 314, col: 7, offset: 8738},
	val: "env:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 314, col: 14, offset: 8745},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 314, col: 17, offset: 8748},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 314, col: 17, offset: 8748},
	name: "BashEnvironmentVariable",
},
&ruleRefExpr{
	pos: position{line: 314, col: 43, offset: 8774},
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
	pos: position{line: 316, col: 1, offset: 8819},
	expr: &actionExpr{
	pos: position{line: 316, col: 27, offset: 8847},
	run: (*parser).callonBashEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 316, col: 27, offset: 8847},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 316, col: 27, offset: 8847},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 316, col: 36, offset: 8856},
	expr: &charClassMatcher{
	pos: position{line: 316, col: 36, offset: 8856},
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
	pos: position{line: 320, col: 1, offset: 8912},
	expr: &actionExpr{
	pos: position{line: 320, col: 28, offset: 8941},
	run: (*parser).callonPosixEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 320, col: 28, offset: 8941},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 320, col: 28, offset: 8941},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 320, col: 32, offset: 8945},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 320, col: 34, offset: 8947},
	name: "PosixEnvironmentVariableContent",
},
},
&litMatcher{
	pos: position{line: 320, col: 66, offset: 8979},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PosixEnvironmentVariableContent",
	pos: position{line: 324, col: 1, offset: 9004},
	expr: &actionExpr{
	pos: position{line: 324, col: 35, offset: 9040},
	run: (*parser).callonPosixEnvironmentVariableContent1,
	expr: &labeledExpr{
	pos: position{line: 324, col: 35, offset: 9040},
	label: "v",
	expr: &oneOrMoreExpr{
	pos: position{line: 324, col: 37, offset: 9042},
	expr: &ruleRefExpr{
	pos: position{line: 324, col: 37, offset: 9042},
	name: "PosixEnvironmentVariableCharacter",
},
},
},
},
},
{
	name: "PosixEnvironmentVariableCharacter",
	pos: position{line: 333, col: 1, offset: 9255},
	expr: &choiceExpr{
	pos: position{line: 334, col: 7, offset: 9299},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 334, col: 7, offset: 9299},
	run: (*parser).callonPosixEnvironmentVariableCharacter2,
	expr: &litMatcher{
	pos: position{line: 334, col: 7, offset: 9299},
	val: "\\\"",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 335, col: 7, offset: 9339},
	run: (*parser).callonPosixEnvironmentVariableCharacter4,
	expr: &litMatcher{
	pos: position{line: 335, col: 7, offset: 9339},
	val: "\\\\",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 336, col: 7, offset: 9379},
	run: (*parser).callonPosixEnvironmentVariableCharacter6,
	expr: &litMatcher{
	pos: position{line: 336, col: 7, offset: 9379},
	val: "\\a",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 337, col: 7, offset: 9419},
	run: (*parser).callonPosixEnvironmentVariableCharacter8,
	expr: &litMatcher{
	pos: position{line: 337, col: 7, offset: 9419},
	val: "\\b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 338, col: 7, offset: 9459},
	run: (*parser).callonPosixEnvironmentVariableCharacter10,
	expr: &litMatcher{
	pos: position{line: 338, col: 7, offset: 9459},
	val: "\\f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 339, col: 7, offset: 9499},
	run: (*parser).callonPosixEnvironmentVariableCharacter12,
	expr: &litMatcher{
	pos: position{line: 339, col: 7, offset: 9499},
	val: "\\n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 340, col: 7, offset: 9539},
	run: (*parser).callonPosixEnvironmentVariableCharacter14,
	expr: &litMatcher{
	pos: position{line: 340, col: 7, offset: 9539},
	val: "\\r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 341, col: 7, offset: 9579},
	run: (*parser).callonPosixEnvironmentVariableCharacter16,
	expr: &litMatcher{
	pos: position{line: 341, col: 7, offset: 9579},
	val: "\\t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 342, col: 7, offset: 9619},
	run: (*parser).callonPosixEnvironmentVariableCharacter18,
	expr: &litMatcher{
	pos: position{line: 342, col: 7, offset: 9619},
	val: "\\v",
	ignoreCase: false,
},
},
&charClassMatcher{
	pos: position{line: 343, col: 7, offset: 9659},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 344, col: 7, offset: 9677},
	val: "[\\x23-\\x3c]",
	ranges: []rune{'#','<',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 345, col: 7, offset: 9695},
	val: "[\\x3e-\\x5b]",
	ranges: []rune{'>','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 346, col: 7, offset: 9713},
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
	pos: position{line: 348, col: 1, offset: 9726},
	expr: &choiceExpr{
	pos: position{line: 348, col: 14, offset: 9741},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 348, col: 14, offset: 9741},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 348, col: 24, offset: 9751},
	name: "Local",
},
&ruleRefExpr{
	pos: position{line: 348, col: 32, offset: 9759},
	name: "Http",
},
&ruleRefExpr{
	pos: position{line: 348, col: 39, offset: 9766},
	name: "Env",
},
	},
},
},
{
	name: "ImportHashed",
	pos: position{line: 350, col: 1, offset: 9771},
	expr: &actionExpr{
	pos: position{line: 350, col: 16, offset: 9788},
	run: (*parser).callonImportHashed1,
	expr: &labeledExpr{
	pos: position{line: 350, col: 16, offset: 9788},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 350, col: 18, offset: 9790},
	name: "ImportType",
},
},
},
},
{
	name: "Import",
	pos: position{line: 352, col: 1, offset: 9859},
	expr: &choiceExpr{
	pos: position{line: 352, col: 10, offset: 9870},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 352, col: 10, offset: 9870},
	run: (*parser).callonImport2,
	expr: &seqExpr{
	pos: position{line: 352, col: 10, offset: 9870},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 352, col: 10, offset: 9870},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 352, col: 12, offset: 9872},
	name: "ImportHashed",
},
},
&ruleRefExpr{
	pos: position{line: 352, col: 25, offset: 9885},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 352, col: 27, offset: 9887},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 352, col: 30, offset: 9890},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 352, col: 33, offset: 9893},
	name: "Text",
},
	},
},
},
&actionExpr{
	pos: position{line: 353, col: 10, offset: 9990},
	run: (*parser).callonImport10,
	expr: &labeledExpr{
	pos: position{line: 353, col: 10, offset: 9990},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 353, col: 12, offset: 9992},
	name: "ImportHashed",
},
},
},
	},
},
},
{
	name: "LetBinding",
	pos: position{line: 356, col: 1, offset: 10087},
	expr: &actionExpr{
	pos: position{line: 356, col: 14, offset: 10102},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 356, col: 14, offset: 10102},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 356, col: 14, offset: 10102},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 356, col: 18, offset: 10106},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 356, col: 21, offset: 10109},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 356, col: 27, offset: 10115},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 356, col: 44, offset: 10132},
	name: "_",
},
&labeledExpr{
	pos: position{line: 356, col: 46, offset: 10134},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 356, col: 48, offset: 10136},
	expr: &seqExpr{
	pos: position{line: 356, col: 49, offset: 10137},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 356, col: 49, offset: 10137},
	name: "Annotation",
},
&ruleRefExpr{
	pos: position{line: 356, col: 60, offset: 10148},
	name: "_",
},
	},
},
},
},
&litMatcher{
	pos: position{line: 357, col: 13, offset: 10164},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 357, col: 17, offset: 10168},
	name: "_",
},
&labeledExpr{
	pos: position{line: 357, col: 19, offset: 10170},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 357, col: 21, offset: 10172},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 357, col: 32, offset: 10183},
	name: "_",
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 372, col: 1, offset: 10492},
	expr: &choiceExpr{
	pos: position{line: 373, col: 7, offset: 10513},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 373, col: 7, offset: 10513},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 373, col: 7, offset: 10513},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 373, col: 7, offset: 10513},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 373, col: 14, offset: 10520},
	name: "_",
},
&litMatcher{
	pos: position{line: 373, col: 16, offset: 10522},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 373, col: 20, offset: 10526},
	name: "_",
},
&labeledExpr{
	pos: position{line: 373, col: 22, offset: 10528},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 373, col: 28, offset: 10534},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 373, col: 45, offset: 10551},
	name: "_",
},
&litMatcher{
	pos: position{line: 373, col: 47, offset: 10553},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 373, col: 51, offset: 10557},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 373, col: 54, offset: 10560},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 373, col: 56, offset: 10562},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 373, col: 67, offset: 10573},
	name: "_",
},
&litMatcher{
	pos: position{line: 373, col: 69, offset: 10575},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 373, col: 73, offset: 10579},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 373, col: 75, offset: 10581},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 373, col: 81, offset: 10587},
	name: "_",
},
&labeledExpr{
	pos: position{line: 373, col: 83, offset: 10589},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 373, col: 88, offset: 10594},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 376, col: 7, offset: 10711},
	run: (*parser).callonExpression22,
	expr: &seqExpr{
	pos: position{line: 376, col: 7, offset: 10711},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 376, col: 7, offset: 10711},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 376, col: 10, offset: 10714},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 376, col: 13, offset: 10717},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 376, col: 18, offset: 10722},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 376, col: 29, offset: 10733},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 376, col: 31, offset: 10735},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 376, col: 36, offset: 10740},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 376, col: 39, offset: 10743},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 376, col: 41, offset: 10745},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 376, col: 52, offset: 10756},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 376, col: 54, offset: 10758},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 376, col: 59, offset: 10763},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 376, col: 62, offset: 10766},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 376, col: 64, offset: 10768},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 379, col: 7, offset: 10854},
	run: (*parser).callonExpression38,
	expr: &seqExpr{
	pos: position{line: 379, col: 7, offset: 10854},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 379, col: 7, offset: 10854},
	label: "bindings",
	expr: &oneOrMoreExpr{
	pos: position{line: 379, col: 16, offset: 10863},
	expr: &ruleRefExpr{
	pos: position{line: 379, col: 16, offset: 10863},
	name: "LetBinding",
},
},
},
&ruleRefExpr{
	pos: position{line: 379, col: 28, offset: 10875},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 379, col: 31, offset: 10878},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 379, col: 34, offset: 10881},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 379, col: 36, offset: 10883},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 386, col: 7, offset: 11123},
	run: (*parser).callonExpression47,
	expr: &seqExpr{
	pos: position{line: 386, col: 7, offset: 11123},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 386, col: 7, offset: 11123},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 386, col: 14, offset: 11130},
	name: "_",
},
&litMatcher{
	pos: position{line: 386, col: 16, offset: 11132},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 386, col: 20, offset: 11136},
	name: "_",
},
&labeledExpr{
	pos: position{line: 386, col: 22, offset: 11138},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 386, col: 28, offset: 11144},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 386, col: 45, offset: 11161},
	name: "_",
},
&litMatcher{
	pos: position{line: 386, col: 47, offset: 11163},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 386, col: 51, offset: 11167},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 386, col: 54, offset: 11170},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 386, col: 56, offset: 11172},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 386, col: 67, offset: 11183},
	name: "_",
},
&litMatcher{
	pos: position{line: 386, col: 69, offset: 11185},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 386, col: 73, offset: 11189},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 386, col: 75, offset: 11191},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 386, col: 81, offset: 11197},
	name: "_",
},
&labeledExpr{
	pos: position{line: 386, col: 83, offset: 11199},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 386, col: 88, offset: 11204},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 389, col: 7, offset: 11313},
	run: (*parser).callonExpression67,
	expr: &seqExpr{
	pos: position{line: 389, col: 7, offset: 11313},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 389, col: 7, offset: 11313},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 389, col: 9, offset: 11315},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 389, col: 28, offset: 11334},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 389, col: 30, offset: 11336},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 389, col: 36, offset: 11342},
	name: "_",
},
&labeledExpr{
	pos: position{line: 389, col: 38, offset: 11344},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 389, col: 40, offset: 11346},
	name: "Expression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 390, col: 7, offset: 11406},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 392, col: 1, offset: 11427},
	expr: &actionExpr{
	pos: position{line: 392, col: 14, offset: 11442},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 392, col: 14, offset: 11442},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 392, col: 14, offset: 11442},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 392, col: 18, offset: 11446},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 392, col: 21, offset: 11449},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 392, col: 23, offset: 11451},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 394, col: 1, offset: 11481},
	expr: &choiceExpr{
	pos: position{line: 395, col: 5, offset: 11509},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 395, col: 5, offset: 11509},
	name: "EmptyList",
},
&actionExpr{
	pos: position{line: 396, col: 5, offset: 11523},
	run: (*parser).callonAnnotatedExpression3,
	expr: &seqExpr{
	pos: position{line: 396, col: 5, offset: 11523},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 396, col: 5, offset: 11523},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 396, col: 7, offset: 11525},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 396, col: 26, offset: 11544},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 396, col: 28, offset: 11546},
	expr: &seqExpr{
	pos: position{line: 396, col: 29, offset: 11547},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 396, col: 29, offset: 11547},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 396, col: 31, offset: 11549},
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
	pos: position{line: 401, col: 1, offset: 11674},
	expr: &actionExpr{
	pos: position{line: 401, col: 13, offset: 11688},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 401, col: 13, offset: 11688},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 401, col: 13, offset: 11688},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 401, col: 17, offset: 11692},
	name: "_",
},
&litMatcher{
	pos: position{line: 401, col: 19, offset: 11694},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 401, col: 23, offset: 11698},
	name: "_",
},
&litMatcher{
	pos: position{line: 401, col: 25, offset: 11700},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 401, col: 29, offset: 11704},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 401, col: 32, offset: 11707},
	name: "List",
},
&ruleRefExpr{
	pos: position{line: 401, col: 37, offset: 11712},
	name: "_",
},
&labeledExpr{
	pos: position{line: 401, col: 39, offset: 11714},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 401, col: 41, offset: 11716},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 405, col: 1, offset: 11779},
	expr: &ruleRefExpr{
	pos: position{line: 405, col: 22, offset: 11802},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 407, col: 1, offset: 11823},
	expr: &ruleRefExpr{
	pos: position{line: 407, col: 23, offset: 11847},
	name: "PlusExpression",
},
},
{
	name: "MorePlus",
	pos: position{line: 409, col: 1, offset: 11863},
	expr: &actionExpr{
	pos: position{line: 409, col: 12, offset: 11876},
	run: (*parser).callonMorePlus1,
	expr: &seqExpr{
	pos: position{line: 409, col: 12, offset: 11876},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 409, col: 12, offset: 11876},
	name: "_",
},
&litMatcher{
	pos: position{line: 409, col: 14, offset: 11878},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 409, col: 18, offset: 11882},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 409, col: 21, offset: 11885},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 409, col: 23, offset: 11887},
	name: "TimesExpression",
},
},
	},
},
},
},
{
	name: "PlusExpression",
	pos: position{line: 410, col: 1, offset: 11921},
	expr: &actionExpr{
	pos: position{line: 411, col: 7, offset: 11946},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 411, col: 7, offset: 11946},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 411, col: 7, offset: 11946},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 411, col: 13, offset: 11952},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 411, col: 29, offset: 11968},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 411, col: 34, offset: 11973},
	expr: &ruleRefExpr{
	pos: position{line: 411, col: 34, offset: 11973},
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
	pos: position{line: 420, col: 1, offset: 12201},
	expr: &actionExpr{
	pos: position{line: 420, col: 13, offset: 12215},
	run: (*parser).callonMoreTimes1,
	expr: &seqExpr{
	pos: position{line: 420, col: 13, offset: 12215},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 420, col: 13, offset: 12215},
	name: "_",
},
&litMatcher{
	pos: position{line: 420, col: 15, offset: 12217},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 420, col: 19, offset: 12221},
	name: "_",
},
&labeledExpr{
	pos: position{line: 420, col: 21, offset: 12223},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 420, col: 23, offset: 12225},
	name: "ApplicationExpression",
},
},
	},
},
},
},
{
	name: "TimesExpression",
	pos: position{line: 421, col: 1, offset: 12265},
	expr: &actionExpr{
	pos: position{line: 422, col: 7, offset: 12291},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 422, col: 7, offset: 12291},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 422, col: 7, offset: 12291},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 422, col: 13, offset: 12297},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 422, col: 35, offset: 12319},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 422, col: 40, offset: 12324},
	expr: &ruleRefExpr{
	pos: position{line: 422, col: 40, offset: 12324},
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
	pos: position{line: 431, col: 1, offset: 12554},
	expr: &actionExpr{
	pos: position{line: 431, col: 25, offset: 12580},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 431, col: 25, offset: 12580},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 431, col: 25, offset: 12580},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 431, col: 27, offset: 12582},
	name: "FirstApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 431, col: 54, offset: 12609},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 431, col: 59, offset: 12614},
	expr: &seqExpr{
	pos: position{line: 431, col: 60, offset: 12615},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 431, col: 60, offset: 12615},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 431, col: 63, offset: 12618},
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
	pos: position{line: 440, col: 1, offset: 12868},
	expr: &choiceExpr{
	pos: position{line: 441, col: 8, offset: 12906},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 441, col: 8, offset: 12906},
	run: (*parser).callonFirstApplicationExpression2,
	expr: &seqExpr{
	pos: position{line: 441, col: 8, offset: 12906},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 441, col: 8, offset: 12906},
	name: "Some",
},
&ruleRefExpr{
	pos: position{line: 441, col: 13, offset: 12911},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 441, col: 16, offset: 12914},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 441, col: 18, offset: 12916},
	name: "ImportExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 442, col: 8, offset: 12971},
	name: "ImportExpression",
},
	},
},
},
{
	name: "ImportExpression",
	pos: position{line: 444, col: 1, offset: 12989},
	expr: &choiceExpr{
	pos: position{line: 444, col: 20, offset: 13010},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 444, col: 20, offset: 13010},
	name: "Import",
},
&ruleRefExpr{
	pos: position{line: 444, col: 29, offset: 13019},
	name: "SelectorExpression",
},
	},
},
},
{
	name: "SelectorExpression",
	pos: position{line: 446, col: 1, offset: 13039},
	expr: &actionExpr{
	pos: position{line: 446, col: 22, offset: 13062},
	run: (*parser).callonSelectorExpression1,
	expr: &seqExpr{
	pos: position{line: 446, col: 22, offset: 13062},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 446, col: 22, offset: 13062},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 446, col: 24, offset: 13064},
	name: "PrimitiveExpression",
},
},
&labeledExpr{
	pos: position{line: 446, col: 44, offset: 13084},
	label: "ls",
	expr: &zeroOrMoreExpr{
	pos: position{line: 446, col: 47, offset: 13087},
	expr: &seqExpr{
	pos: position{line: 446, col: 48, offset: 13088},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 446, col: 48, offset: 13088},
	name: "_",
},
&litMatcher{
	pos: position{line: 446, col: 50, offset: 13090},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 446, col: 54, offset: 13094},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 446, col: 56, offset: 13096},
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
	pos: position{line: 456, col: 1, offset: 13329},
	expr: &choiceExpr{
	pos: position{line: 457, col: 7, offset: 13359},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 457, col: 7, offset: 13359},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 458, col: 7, offset: 13379},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 459, col: 7, offset: 13400},
	name: "IntegerLiteral",
},
&ruleRefExpr{
	pos: position{line: 460, col: 7, offset: 13421},
	name: "TextLiteral",
},
&actionExpr{
	pos: position{line: 461, col: 7, offset: 13439},
	run: (*parser).callonPrimitiveExpression6,
	expr: &seqExpr{
	pos: position{line: 461, col: 7, offset: 13439},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 461, col: 7, offset: 13439},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 461, col: 11, offset: 13443},
	name: "_",
},
&labeledExpr{
	pos: position{line: 461, col: 13, offset: 13445},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 461, col: 15, offset: 13447},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 461, col: 35, offset: 13467},
	name: "_",
},
&litMatcher{
	pos: position{line: 461, col: 37, offset: 13469},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 462, col: 7, offset: 13497},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 463, col: 7, offset: 13523},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 464, col: 7, offset: 13540},
	run: (*parser).callonPrimitiveExpression16,
	expr: &seqExpr{
	pos: position{line: 464, col: 7, offset: 13540},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 464, col: 7, offset: 13540},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 464, col: 11, offset: 13544},
	name: "_",
},
&labeledExpr{
	pos: position{line: 464, col: 14, offset: 13547},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 464, col: 16, offset: 13549},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 464, col: 27, offset: 13560},
	name: "_",
},
&litMatcher{
	pos: position{line: 464, col: 29, offset: 13562},
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
	pos: position{line: 466, col: 1, offset: 13585},
	expr: &choiceExpr{
	pos: position{line: 467, col: 7, offset: 13615},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 467, col: 7, offset: 13615},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &litMatcher{
	pos: position{line: 467, col: 7, offset: 13615},
	val: "=",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 468, col: 7, offset: 13670},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 469, col: 7, offset: 13695},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 470, col: 7, offset: 13723},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 470, col: 7, offset: 13723},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 472, col: 1, offset: 13769},
	expr: &actionExpr{
	pos: position{line: 472, col: 19, offset: 13789},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 472, col: 19, offset: 13789},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 472, col: 19, offset: 13789},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 472, col: 24, offset: 13794},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 472, col: 33, offset: 13803},
	name: "_",
},
&litMatcher{
	pos: position{line: 472, col: 35, offset: 13805},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 472, col: 39, offset: 13809},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 472, col: 42, offset: 13812},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 472, col: 47, offset: 13817},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 475, col: 1, offset: 13874},
	expr: &actionExpr{
	pos: position{line: 475, col: 18, offset: 13893},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 475, col: 18, offset: 13893},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 475, col: 18, offset: 13893},
	name: "_",
},
&litMatcher{
	pos: position{line: 475, col: 20, offset: 13895},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 475, col: 24, offset: 13899},
	name: "_",
},
&labeledExpr{
	pos: position{line: 475, col: 26, offset: 13901},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 475, col: 28, offset: 13903},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 476, col: 1, offset: 13935},
	expr: &actionExpr{
	pos: position{line: 477, col: 7, offset: 13964},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 477, col: 7, offset: 13964},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 477, col: 7, offset: 13964},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 477, col: 13, offset: 13970},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 477, col: 29, offset: 13986},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 477, col: 34, offset: 13991},
	expr: &ruleRefExpr{
	pos: position{line: 477, col: 34, offset: 13991},
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
	pos: position{line: 487, col: 1, offset: 14387},
	expr: &actionExpr{
	pos: position{line: 487, col: 22, offset: 14410},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 487, col: 22, offset: 14410},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 487, col: 22, offset: 14410},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 487, col: 27, offset: 14415},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 487, col: 36, offset: 14424},
	name: "_",
},
&litMatcher{
	pos: position{line: 487, col: 38, offset: 14426},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 487, col: 42, offset: 14430},
	name: "_",
},
&labeledExpr{
	pos: position{line: 487, col: 44, offset: 14432},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 487, col: 49, offset: 14437},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 490, col: 1, offset: 14494},
	expr: &actionExpr{
	pos: position{line: 490, col: 21, offset: 14516},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 490, col: 21, offset: 14516},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 490, col: 21, offset: 14516},
	name: "_",
},
&litMatcher{
	pos: position{line: 490, col: 23, offset: 14518},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 490, col: 27, offset: 14522},
	name: "_",
},
&labeledExpr{
	pos: position{line: 490, col: 29, offset: 14524},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 490, col: 31, offset: 14526},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 491, col: 1, offset: 14561},
	expr: &actionExpr{
	pos: position{line: 492, col: 7, offset: 14593},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 492, col: 7, offset: 14593},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 492, col: 7, offset: 14593},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 492, col: 13, offset: 14599},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 492, col: 32, offset: 14618},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 492, col: 37, offset: 14623},
	expr: &ruleRefExpr{
	pos: position{line: 492, col: 37, offset: 14623},
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
	pos: position{line: 502, col: 1, offset: 15025},
	expr: &actionExpr{
	pos: position{line: 502, col: 12, offset: 15038},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 502, col: 12, offset: 15038},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 502, col: 12, offset: 15038},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 502, col: 16, offset: 15042},
	name: "_",
},
&labeledExpr{
	pos: position{line: 502, col: 18, offset: 15044},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 502, col: 20, offset: 15046},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 502, col: 31, offset: 15057},
	name: "_",
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 504, col: 1, offset: 15076},
	expr: &actionExpr{
	pos: position{line: 505, col: 7, offset: 15106},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 505, col: 7, offset: 15106},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 505, col: 7, offset: 15106},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 505, col: 11, offset: 15110},
	name: "_",
},
&labeledExpr{
	pos: position{line: 505, col: 13, offset: 15112},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 505, col: 19, offset: 15118},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 505, col: 30, offset: 15129},
	name: "_",
},
&labeledExpr{
	pos: position{line: 505, col: 32, offset: 15131},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 505, col: 37, offset: 15136},
	expr: &ruleRefExpr{
	pos: position{line: 505, col: 37, offset: 15136},
	name: "MoreList",
},
},
},
&litMatcher{
	pos: position{line: 505, col: 47, offset: 15146},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 515, col: 1, offset: 15422},
	expr: &notExpr{
	pos: position{line: 515, col: 7, offset: 15430},
	expr: &anyMatcher{
	line: 515, col: 8, offset: 15431,
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

func (c *current) onHomePath1(p interface{}) (interface{}, error) {
 return Local(path.Join("~", p.(string))), nil 
}

func (p *parser) callonHomePath1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHomePath1(stack["p"])
}

func (c *current) onAbsolutePath1(p interface{}) (interface{}, error) {
 return Local(path.Join("/", p.(string))), nil 
}

func (p *parser) callonAbsolutePath1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAbsolutePath1(stack["p"])
}

func (c *current) onHttpRaw1() (interface{}, error) {
 return string(c.text), nil 
}

func (p *parser) callonHttpRaw1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHttpRaw1()
}

func (c *current) onHttp1(url interface{}) (interface{}, error) {
 return Remote(url.(string)), nil 
}

func (p *parser) callonHttp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHttp1(stack["url"])
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

