
package parser

import (
"bytes"
"errors"
"fmt"
"io"
"io/ioutil"
"math"
"net"
"net/url"
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
	pos: position{line: 24, col: 1, offset: 205},
	expr: &actionExpr{
	pos: position{line: 24, col: 13, offset: 219},
	run: (*parser).callonDhallFile1,
	expr: &seqExpr{
	pos: position{line: 24, col: 13, offset: 219},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 24, col: 13, offset: 219},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 24, col: 15, offset: 221},
	name: "CompleteExpression",
},
},
&ruleRefExpr{
	pos: position{line: 24, col: 34, offset: 240},
	name: "EOF",
},
	},
},
},
},
{
	name: "CompleteExpression",
	pos: position{line: 26, col: 1, offset: 263},
	expr: &actionExpr{
	pos: position{line: 26, col: 22, offset: 286},
	run: (*parser).callonCompleteExpression1,
	expr: &seqExpr{
	pos: position{line: 26, col: 22, offset: 286},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 26, col: 22, offset: 286},
	name: "_",
},
&labeledExpr{
	pos: position{line: 26, col: 24, offset: 288},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 26, col: 26, offset: 290},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 26, col: 37, offset: 301},
	name: "_",
},
	},
},
},
},
{
	name: "EOL",
	pos: position{line: 28, col: 1, offset: 322},
	expr: &choiceExpr{
	pos: position{line: 28, col: 7, offset: 330},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 28, col: 7, offset: 330},
	val: "\n",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 28, col: 14, offset: 337},
	val: "\r\n",
	ignoreCase: false,
},
	},
},
},
{
	name: "BlockComment",
	pos: position{line: 30, col: 1, offset: 345},
	expr: &seqExpr{
	pos: position{line: 30, col: 16, offset: 362},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 30, col: 16, offset: 362},
	val: "{-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 30, col: 21, offset: 367},
	name: "BlockCommentContinue",
},
	},
},
},
{
	name: "BlockCommentChunk",
	pos: position{line: 32, col: 1, offset: 389},
	expr: &choiceExpr{
	pos: position{line: 33, col: 5, offset: 415},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 33, col: 5, offset: 415},
	name: "BlockComment",
},
&charClassMatcher{
	pos: position{line: 34, col: 5, offset: 432},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 35, col: 5, offset: 458},
	name: "EOL",
},
	},
},
},
{
	name: "BlockCommentContinue",
	pos: position{line: 37, col: 1, offset: 463},
	expr: &choiceExpr{
	pos: position{line: 37, col: 24, offset: 488},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 37, col: 24, offset: 488},
	val: "-}",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 37, col: 31, offset: 495},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 37, col: 31, offset: 495},
	name: "BlockCommentChunk",
},
&ruleRefExpr{
	pos: position{line: 37, col: 49, offset: 513},
	name: "BlockCommentContinue",
},
	},
},
	},
},
},
{
	name: "NotEOL",
	pos: position{line: 39, col: 1, offset: 535},
	expr: &charClassMatcher{
	pos: position{line: 39, col: 10, offset: 546},
	val: "[\\t\\u0020-\\U0010ffff]",
	chars: []rune{'\t',},
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "LineComment",
	pos: position{line: 41, col: 1, offset: 569},
	expr: &actionExpr{
	pos: position{line: 41, col: 15, offset: 585},
	run: (*parser).callonLineComment1,
	expr: &seqExpr{
	pos: position{line: 41, col: 15, offset: 585},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 41, col: 15, offset: 585},
	val: "--",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 41, col: 20, offset: 590},
	label: "content",
	expr: &actionExpr{
	pos: position{line: 41, col: 29, offset: 599},
	run: (*parser).callonLineComment5,
	expr: &zeroOrMoreExpr{
	pos: position{line: 41, col: 29, offset: 599},
	expr: &ruleRefExpr{
	pos: position{line: 41, col: 29, offset: 599},
	name: "NotEOL",
},
},
},
},
&ruleRefExpr{
	pos: position{line: 41, col: 68, offset: 638},
	name: "EOL",
},
	},
},
},
},
{
	name: "WhitespaceChunk",
	pos: position{line: 43, col: 1, offset: 667},
	expr: &choiceExpr{
	pos: position{line: 43, col: 19, offset: 687},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 43, col: 19, offset: 687},
	val: " ",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 43, col: 25, offset: 693},
	val: "\t",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 43, col: 32, offset: 700},
	name: "EOL",
},
&ruleRefExpr{
	pos: position{line: 43, col: 38, offset: 706},
	name: "LineComment",
},
&ruleRefExpr{
	pos: position{line: 43, col: 52, offset: 720},
	name: "BlockComment",
},
	},
},
},
{
	name: "_",
	pos: position{line: 45, col: 1, offset: 734},
	expr: &zeroOrMoreExpr{
	pos: position{line: 45, col: 5, offset: 740},
	expr: &ruleRefExpr{
	pos: position{line: 45, col: 5, offset: 740},
	name: "WhitespaceChunk",
},
},
},
{
	name: "_1",
	pos: position{line: 47, col: 1, offset: 758},
	expr: &oneOrMoreExpr{
	pos: position{line: 47, col: 6, offset: 765},
	expr: &ruleRefExpr{
	pos: position{line: 47, col: 6, offset: 765},
	name: "WhitespaceChunk",
},
},
},
{
	name: "Digit",
	pos: position{line: 49, col: 1, offset: 783},
	expr: &charClassMatcher{
	pos: position{line: 49, col: 9, offset: 793},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "HexDig",
	pos: position{line: 51, col: 1, offset: 800},
	expr: &choiceExpr{
	pos: position{line: 51, col: 10, offset: 811},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 51, col: 10, offset: 811},
	name: "Digit",
},
&charClassMatcher{
	pos: position{line: 51, col: 18, offset: 819},
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
	pos: position{line: 53, col: 1, offset: 827},
	expr: &charClassMatcher{
	pos: position{line: 53, col: 24, offset: 852},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SimpleLabelNextChar",
	pos: position{line: 54, col: 1, offset: 862},
	expr: &charClassMatcher{
	pos: position{line: 54, col: 23, offset: 886},
	val: "[A-Za-z0-9_/-]",
	chars: []rune{'_','/','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SimpleLabel",
	pos: position{line: 55, col: 1, offset: 901},
	expr: &choiceExpr{
	pos: position{line: 55, col: 15, offset: 917},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 55, col: 15, offset: 917},
	run: (*parser).callonSimpleLabel2,
	expr: &seqExpr{
	pos: position{line: 55, col: 15, offset: 917},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 55, col: 15, offset: 917},
	name: "Keyword",
},
&oneOrMoreExpr{
	pos: position{line: 55, col: 23, offset: 925},
	expr: &ruleRefExpr{
	pos: position{line: 55, col: 23, offset: 925},
	name: "SimpleLabelNextChar",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 56, col: 13, offset: 989},
	run: (*parser).callonSimpleLabel7,
	expr: &seqExpr{
	pos: position{line: 56, col: 13, offset: 989},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 56, col: 13, offset: 989},
	expr: &ruleRefExpr{
	pos: position{line: 56, col: 14, offset: 990},
	name: "Keyword",
},
},
&ruleRefExpr{
	pos: position{line: 56, col: 22, offset: 998},
	name: "SimpleLabelFirstChar",
},
&zeroOrMoreExpr{
	pos: position{line: 56, col: 43, offset: 1019},
	expr: &ruleRefExpr{
	pos: position{line: 56, col: 43, offset: 1019},
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
	pos: position{line: 63, col: 1, offset: 1120},
	expr: &actionExpr{
	pos: position{line: 63, col: 9, offset: 1130},
	run: (*parser).callonLabel1,
	expr: &labeledExpr{
	pos: position{line: 63, col: 9, offset: 1130},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 63, col: 15, offset: 1136},
	name: "SimpleLabel",
},
},
},
},
{
	name: "NonreservedLabel",
	pos: position{line: 65, col: 1, offset: 1171},
	expr: &choiceExpr{
	pos: position{line: 65, col: 20, offset: 1192},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 65, col: 20, offset: 1192},
	run: (*parser).callonNonreservedLabel2,
	expr: &seqExpr{
	pos: position{line: 65, col: 20, offset: 1192},
	exprs: []interface{}{
&andExpr{
	pos: position{line: 65, col: 20, offset: 1192},
	expr: &seqExpr{
	pos: position{line: 65, col: 22, offset: 1194},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 65, col: 22, offset: 1194},
	name: "Reserved",
},
&ruleRefExpr{
	pos: position{line: 65, col: 31, offset: 1203},
	name: "SimpleLabelNextChar",
},
	},
},
},
&labeledExpr{
	pos: position{line: 65, col: 52, offset: 1224},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 65, col: 58, offset: 1230},
	name: "Label",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 66, col: 19, offset: 1276},
	run: (*parser).callonNonreservedLabel10,
	expr: &seqExpr{
	pos: position{line: 66, col: 19, offset: 1276},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 66, col: 19, offset: 1276},
	expr: &ruleRefExpr{
	pos: position{line: 66, col: 20, offset: 1277},
	name: "Reserved",
},
},
&labeledExpr{
	pos: position{line: 66, col: 29, offset: 1286},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 66, col: 35, offset: 1292},
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
	pos: position{line: 68, col: 1, offset: 1321},
	expr: &ruleRefExpr{
	pos: position{line: 68, col: 12, offset: 1334},
	name: "Label",
},
},
{
	name: "DoubleQuoteChunk",
	pos: position{line: 71, col: 1, offset: 1342},
	expr: &choiceExpr{
	pos: position{line: 72, col: 6, offset: 1368},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 72, col: 6, offset: 1368},
	name: "Interpolation",
},
&actionExpr{
	pos: position{line: 73, col: 6, offset: 1387},
	run: (*parser).callonDoubleQuoteChunk3,
	expr: &seqExpr{
	pos: position{line: 73, col: 6, offset: 1387},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 73, col: 6, offset: 1387},
	val: "\\",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 73, col: 11, offset: 1392},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 73, col: 13, offset: 1394},
	name: "DoubleQuoteEscaped",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 74, col: 6, offset: 1436},
	name: "DoubleQuoteChar",
},
	},
},
},
{
	name: "DoubleQuoteEscaped",
	pos: position{line: 76, col: 1, offset: 1453},
	expr: &choiceExpr{
	pos: position{line: 77, col: 8, offset: 1483},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 77, col: 8, offset: 1483},
	val: "\"",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 78, col: 8, offset: 1494},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 79, col: 8, offset: 1505},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 80, col: 8, offset: 1517},
	val: "/",
	ignoreCase: false,
},
&actionExpr{
	pos: position{line: 81, col: 8, offset: 1528},
	run: (*parser).callonDoubleQuoteEscaped6,
	expr: &litMatcher{
	pos: position{line: 81, col: 8, offset: 1528},
	val: "b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 82, col: 8, offset: 1568},
	run: (*parser).callonDoubleQuoteEscaped8,
	expr: &litMatcher{
	pos: position{line: 82, col: 8, offset: 1568},
	val: "f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 83, col: 8, offset: 1608},
	run: (*parser).callonDoubleQuoteEscaped10,
	expr: &litMatcher{
	pos: position{line: 83, col: 8, offset: 1608},
	val: "n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 84, col: 8, offset: 1648},
	run: (*parser).callonDoubleQuoteEscaped12,
	expr: &litMatcher{
	pos: position{line: 84, col: 8, offset: 1648},
	val: "r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 85, col: 8, offset: 1688},
	run: (*parser).callonDoubleQuoteEscaped14,
	expr: &litMatcher{
	pos: position{line: 85, col: 8, offset: 1688},
	val: "t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 86, col: 8, offset: 1728},
	run: (*parser).callonDoubleQuoteEscaped16,
	expr: &seqExpr{
	pos: position{line: 86, col: 8, offset: 1728},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 86, col: 8, offset: 1728},
	val: "u",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 86, col: 12, offset: 1732},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 86, col: 19, offset: 1739},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 86, col: 26, offset: 1746},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 86, col: 33, offset: 1753},
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
	pos: position{line: 91, col: 1, offset: 1885},
	expr: &choiceExpr{
	pos: position{line: 92, col: 6, offset: 1910},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 92, col: 6, offset: 1910},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 93, col: 6, offset: 1927},
	val: "[\\x23-\\x5b]",
	ranges: []rune{'#','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 94, col: 6, offset: 1944},
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
	pos: position{line: 96, col: 1, offset: 1963},
	expr: &actionExpr{
	pos: position{line: 96, col: 22, offset: 1986},
	run: (*parser).callonDoubleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 96, col: 22, offset: 1986},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 96, col: 22, offset: 1986},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 96, col: 26, offset: 1990},
	label: "chunks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 96, col: 33, offset: 1997},
	expr: &ruleRefExpr{
	pos: position{line: 96, col: 33, offset: 1997},
	name: "DoubleQuoteChunk",
},
},
},
&litMatcher{
	pos: position{line: 96, col: 51, offset: 2015},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "SingleQuoteContinue",
	pos: position{line: 113, col: 1, offset: 2483},
	expr: &choiceExpr{
	pos: position{line: 114, col: 7, offset: 2513},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 114, col: 7, offset: 2513},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 114, col: 7, offset: 2513},
	name: "Interpolation",
},
&ruleRefExpr{
	pos: position{line: 114, col: 21, offset: 2527},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 115, col: 7, offset: 2553},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 115, col: 7, offset: 2553},
	name: "EscapedQuotePair",
},
&ruleRefExpr{
	pos: position{line: 115, col: 24, offset: 2570},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 116, col: 7, offset: 2596},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 116, col: 7, offset: 2596},
	name: "EscapedInterpolation",
},
&ruleRefExpr{
	pos: position{line: 116, col: 28, offset: 2617},
	name: "SingleQuoteContinue",
},
	},
},
&seqExpr{
	pos: position{line: 117, col: 7, offset: 2643},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 117, col: 7, offset: 2643},
	name: "SingleQuoteChar",
},
&ruleRefExpr{
	pos: position{line: 117, col: 23, offset: 2659},
	name: "SingleQuoteContinue",
},
	},
},
&litMatcher{
	pos: position{line: 118, col: 7, offset: 2685},
	val: "''",
	ignoreCase: false,
},
	},
},
},
{
	name: "EscapedQuotePair",
	pos: position{line: 120, col: 1, offset: 2691},
	expr: &actionExpr{
	pos: position{line: 120, col: 20, offset: 2712},
	run: (*parser).callonEscapedQuotePair1,
	expr: &litMatcher{
	pos: position{line: 120, col: 20, offset: 2712},
	val: "'''",
	ignoreCase: false,
},
},
},
{
	name: "EscapedInterpolation",
	pos: position{line: 124, col: 1, offset: 2847},
	expr: &actionExpr{
	pos: position{line: 124, col: 24, offset: 2872},
	run: (*parser).callonEscapedInterpolation1,
	expr: &litMatcher{
	pos: position{line: 124, col: 24, offset: 2872},
	val: "''${",
	ignoreCase: false,
},
},
},
{
	name: "SingleQuoteChar",
	pos: position{line: 126, col: 1, offset: 2914},
	expr: &choiceExpr{
	pos: position{line: 127, col: 6, offset: 2939},
	alternatives: []interface{}{
&charClassMatcher{
	pos: position{line: 127, col: 6, offset: 2939},
	val: "[\\x20-\\U0010ffff]",
	ranges: []rune{' ','\U0010ffff',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 128, col: 6, offset: 2962},
	val: "\t",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 129, col: 6, offset: 2972},
	val: "\n",
	ignoreCase: false,
},
	},
},
},
{
	name: "SingleQuoteLiteral",
	pos: position{line: 131, col: 1, offset: 2978},
	expr: &actionExpr{
	pos: position{line: 131, col: 22, offset: 3001},
	run: (*parser).callonSingleQuoteLiteral1,
	expr: &seqExpr{
	pos: position{line: 131, col: 22, offset: 3001},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 131, col: 22, offset: 3001},
	val: "''",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 131, col: 27, offset: 3006},
	name: "EOL",
},
&labeledExpr{
	pos: position{line: 131, col: 31, offset: 3010},
	label: "content",
	expr: &ruleRefExpr{
	pos: position{line: 131, col: 39, offset: 3018},
	name: "SingleQuoteContinue",
},
},
	},
},
},
},
{
	name: "Interpolation",
	pos: position{line: 149, col: 1, offset: 3541},
	expr: &actionExpr{
	pos: position{line: 149, col: 17, offset: 3559},
	run: (*parser).callonInterpolation1,
	expr: &seqExpr{
	pos: position{line: 149, col: 17, offset: 3559},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 149, col: 17, offset: 3559},
	val: "${",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 149, col: 22, offset: 3564},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 149, col: 24, offset: 3566},
	name: "CompleteExpression",
},
},
&litMatcher{
	pos: position{line: 149, col: 43, offset: 3585},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TextLiteral",
	pos: position{line: 151, col: 1, offset: 3608},
	expr: &choiceExpr{
	pos: position{line: 151, col: 15, offset: 3624},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 151, col: 15, offset: 3624},
	name: "DoubleQuoteLiteral",
},
&ruleRefExpr{
	pos: position{line: 151, col: 36, offset: 3645},
	name: "SingleQuoteLiteral",
},
	},
},
},
{
	name: "Reserved",
	pos: position{line: 154, col: 1, offset: 3750},
	expr: &choiceExpr{
	pos: position{line: 155, col: 5, offset: 3767},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 155, col: 5, offset: 3767},
	run: (*parser).callonReserved2,
	expr: &litMatcher{
	pos: position{line: 155, col: 5, offset: 3767},
	val: "Natural/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 156, col: 5, offset: 3845},
	run: (*parser).callonReserved4,
	expr: &litMatcher{
	pos: position{line: 156, col: 5, offset: 3845},
	val: "Natural/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 157, col: 5, offset: 3921},
	run: (*parser).callonReserved6,
	expr: &litMatcher{
	pos: position{line: 157, col: 5, offset: 3921},
	val: "Natural/isZero",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 158, col: 5, offset: 4001},
	run: (*parser).callonReserved8,
	expr: &litMatcher{
	pos: position{line: 158, col: 5, offset: 4001},
	val: "Natural/even",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 159, col: 5, offset: 4077},
	run: (*parser).callonReserved10,
	expr: &litMatcher{
	pos: position{line: 159, col: 5, offset: 4077},
	val: "Natural/odd",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 160, col: 5, offset: 4151},
	run: (*parser).callonReserved12,
	expr: &litMatcher{
	pos: position{line: 160, col: 5, offset: 4151},
	val: "Natural/toInteger",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 161, col: 5, offset: 4237},
	run: (*parser).callonReserved14,
	expr: &litMatcher{
	pos: position{line: 161, col: 5, offset: 4237},
	val: "Natural/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 162, col: 5, offset: 4313},
	run: (*parser).callonReserved16,
	expr: &litMatcher{
	pos: position{line: 162, col: 5, offset: 4313},
	val: "Integer/toDouble",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 163, col: 5, offset: 4397},
	run: (*parser).callonReserved18,
	expr: &litMatcher{
	pos: position{line: 163, col: 5, offset: 4397},
	val: "Integer/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 164, col: 5, offset: 4473},
	run: (*parser).callonReserved20,
	expr: &litMatcher{
	pos: position{line: 164, col: 5, offset: 4473},
	val: "Double/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 165, col: 5, offset: 4547},
	run: (*parser).callonReserved22,
	expr: &litMatcher{
	pos: position{line: 165, col: 5, offset: 4547},
	val: "List/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 166, col: 5, offset: 4619},
	run: (*parser).callonReserved24,
	expr: &litMatcher{
	pos: position{line: 166, col: 5, offset: 4619},
	val: "List/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 167, col: 5, offset: 4689},
	run: (*parser).callonReserved26,
	expr: &litMatcher{
	pos: position{line: 167, col: 5, offset: 4689},
	val: "List/length",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 168, col: 5, offset: 4763},
	run: (*parser).callonReserved28,
	expr: &litMatcher{
	pos: position{line: 168, col: 5, offset: 4763},
	val: "List/head",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 169, col: 5, offset: 4833},
	run: (*parser).callonReserved30,
	expr: &litMatcher{
	pos: position{line: 169, col: 5, offset: 4833},
	val: "List/last",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 170, col: 5, offset: 4903},
	run: (*parser).callonReserved32,
	expr: &litMatcher{
	pos: position{line: 170, col: 5, offset: 4903},
	val: "List/indexed",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 171, col: 5, offset: 4979},
	run: (*parser).callonReserved34,
	expr: &litMatcher{
	pos: position{line: 171, col: 5, offset: 4979},
	val: "List/reverse",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 172, col: 5, offset: 5055},
	run: (*parser).callonReserved36,
	expr: &litMatcher{
	pos: position{line: 172, col: 5, offset: 5055},
	val: "Optional/build",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 173, col: 5, offset: 5135},
	run: (*parser).callonReserved38,
	expr: &litMatcher{
	pos: position{line: 173, col: 5, offset: 5135},
	val: "Optional/fold",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 174, col: 5, offset: 5213},
	run: (*parser).callonReserved40,
	expr: &litMatcher{
	pos: position{line: 174, col: 5, offset: 5213},
	val: "Text/show",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 175, col: 5, offset: 5283},
	run: (*parser).callonReserved42,
	expr: &litMatcher{
	pos: position{line: 175, col: 5, offset: 5283},
	val: "Bool",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 176, col: 5, offset: 5315},
	run: (*parser).callonReserved44,
	expr: &litMatcher{
	pos: position{line: 176, col: 5, offset: 5315},
	val: "True",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 177, col: 5, offset: 5347},
	run: (*parser).callonReserved46,
	expr: &litMatcher{
	pos: position{line: 177, col: 5, offset: 5347},
	val: "False",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 178, col: 5, offset: 5381},
	run: (*parser).callonReserved48,
	expr: &litMatcher{
	pos: position{line: 178, col: 5, offset: 5381},
	val: "Optional",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 179, col: 5, offset: 5421},
	run: (*parser).callonReserved50,
	expr: &litMatcher{
	pos: position{line: 179, col: 5, offset: 5421},
	val: "Natural",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 180, col: 5, offset: 5459},
	run: (*parser).callonReserved52,
	expr: &litMatcher{
	pos: position{line: 180, col: 5, offset: 5459},
	val: "Integer",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 181, col: 5, offset: 5497},
	run: (*parser).callonReserved54,
	expr: &litMatcher{
	pos: position{line: 181, col: 5, offset: 5497},
	val: "Double",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 182, col: 5, offset: 5533},
	run: (*parser).callonReserved56,
	expr: &litMatcher{
	pos: position{line: 182, col: 5, offset: 5533},
	val: "Text",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 183, col: 5, offset: 5565},
	run: (*parser).callonReserved58,
	expr: &litMatcher{
	pos: position{line: 183, col: 5, offset: 5565},
	val: "List",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 184, col: 5, offset: 5597},
	run: (*parser).callonReserved60,
	expr: &litMatcher{
	pos: position{line: 184, col: 5, offset: 5597},
	val: "None",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 185, col: 5, offset: 5629},
	run: (*parser).callonReserved62,
	expr: &litMatcher{
	pos: position{line: 185, col: 5, offset: 5629},
	val: "Type",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 186, col: 5, offset: 5661},
	run: (*parser).callonReserved64,
	expr: &litMatcher{
	pos: position{line: 186, col: 5, offset: 5661},
	val: "Kind",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 187, col: 5, offset: 5693},
	run: (*parser).callonReserved66,
	expr: &litMatcher{
	pos: position{line: 187, col: 5, offset: 5693},
	val: "Sort",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "If",
	pos: position{line: 189, col: 1, offset: 5722},
	expr: &litMatcher{
	pos: position{line: 189, col: 6, offset: 5729},
	val: "if",
	ignoreCase: false,
},
},
{
	name: "Then",
	pos: position{line: 190, col: 1, offset: 5734},
	expr: &litMatcher{
	pos: position{line: 190, col: 8, offset: 5743},
	val: "then",
	ignoreCase: false,
},
},
{
	name: "Else",
	pos: position{line: 191, col: 1, offset: 5750},
	expr: &litMatcher{
	pos: position{line: 191, col: 8, offset: 5759},
	val: "else",
	ignoreCase: false,
},
},
{
	name: "Let",
	pos: position{line: 192, col: 1, offset: 5766},
	expr: &litMatcher{
	pos: position{line: 192, col: 7, offset: 5774},
	val: "let",
	ignoreCase: false,
},
},
{
	name: "In",
	pos: position{line: 193, col: 1, offset: 5780},
	expr: &litMatcher{
	pos: position{line: 193, col: 6, offset: 5787},
	val: "in",
	ignoreCase: false,
},
},
{
	name: "As",
	pos: position{line: 194, col: 1, offset: 5792},
	expr: &litMatcher{
	pos: position{line: 194, col: 6, offset: 5799},
	val: "as",
	ignoreCase: false,
},
},
{
	name: "Using",
	pos: position{line: 195, col: 1, offset: 5804},
	expr: &litMatcher{
	pos: position{line: 195, col: 9, offset: 5814},
	val: "using",
	ignoreCase: false,
},
},
{
	name: "Merge",
	pos: position{line: 196, col: 1, offset: 5822},
	expr: &litMatcher{
	pos: position{line: 196, col: 9, offset: 5832},
	val: "merge",
	ignoreCase: false,
},
},
{
	name: "Missing",
	pos: position{line: 197, col: 1, offset: 5840},
	expr: &actionExpr{
	pos: position{line: 197, col: 11, offset: 5852},
	run: (*parser).callonMissing1,
	expr: &litMatcher{
	pos: position{line: 197, col: 11, offset: 5852},
	val: "missing",
	ignoreCase: false,
},
},
},
{
	name: "True",
	pos: position{line: 198, col: 1, offset: 5888},
	expr: &litMatcher{
	pos: position{line: 198, col: 8, offset: 5897},
	val: "True",
	ignoreCase: false,
},
},
{
	name: "False",
	pos: position{line: 199, col: 1, offset: 5904},
	expr: &litMatcher{
	pos: position{line: 199, col: 9, offset: 5914},
	val: "False",
	ignoreCase: false,
},
},
{
	name: "Infinity",
	pos: position{line: 200, col: 1, offset: 5922},
	expr: &litMatcher{
	pos: position{line: 200, col: 12, offset: 5935},
	val: "Infinity",
	ignoreCase: false,
},
},
{
	name: "NaN",
	pos: position{line: 201, col: 1, offset: 5946},
	expr: &litMatcher{
	pos: position{line: 201, col: 7, offset: 5954},
	val: "NaN",
	ignoreCase: false,
},
},
{
	name: "Some",
	pos: position{line: 202, col: 1, offset: 5960},
	expr: &litMatcher{
	pos: position{line: 202, col: 8, offset: 5969},
	val: "Some",
	ignoreCase: false,
},
},
{
	name: "Keyword",
	pos: position{line: 204, col: 1, offset: 5977},
	expr: &choiceExpr{
	pos: position{line: 205, col: 5, offset: 5993},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 205, col: 5, offset: 5993},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 205, col: 10, offset: 5998},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 205, col: 17, offset: 6005},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 206, col: 5, offset: 6014},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 206, col: 11, offset: 6020},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 207, col: 5, offset: 6027},
	name: "Using",
},
&ruleRefExpr{
	pos: position{line: 207, col: 13, offset: 6035},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 207, col: 23, offset: 6045},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 208, col: 5, offset: 6052},
	name: "True",
},
&ruleRefExpr{
	pos: position{line: 208, col: 12, offset: 6059},
	name: "False",
},
&ruleRefExpr{
	pos: position{line: 209, col: 5, offset: 6069},
	name: "Infinity",
},
&ruleRefExpr{
	pos: position{line: 209, col: 16, offset: 6080},
	name: "NaN",
},
&ruleRefExpr{
	pos: position{line: 210, col: 5, offset: 6088},
	name: "Merge",
},
&ruleRefExpr{
	pos: position{line: 210, col: 13, offset: 6096},
	name: "Some",
},
	},
},
},
{
	name: "Optional",
	pos: position{line: 212, col: 1, offset: 6102},
	expr: &litMatcher{
	pos: position{line: 212, col: 12, offset: 6115},
	val: "Optional",
	ignoreCase: false,
},
},
{
	name: "Text",
	pos: position{line: 213, col: 1, offset: 6126},
	expr: &litMatcher{
	pos: position{line: 213, col: 8, offset: 6135},
	val: "Text",
	ignoreCase: false,
},
},
{
	name: "List",
	pos: position{line: 214, col: 1, offset: 6142},
	expr: &litMatcher{
	pos: position{line: 214, col: 8, offset: 6151},
	val: "List",
	ignoreCase: false,
},
},
{
	name: "Lambda",
	pos: position{line: 216, col: 1, offset: 6159},
	expr: &choiceExpr{
	pos: position{line: 216, col: 11, offset: 6171},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 216, col: 11, offset: 6171},
	val: "\\",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 216, col: 18, offset: 6178},
	val: "λ",
	ignoreCase: false,
},
	},
},
},
{
	name: "Forall",
	pos: position{line: 217, col: 1, offset: 6184},
	expr: &choiceExpr{
	pos: position{line: 217, col: 11, offset: 6196},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 217, col: 11, offset: 6196},
	val: "forall",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 217, col: 22, offset: 6207},
	val: "∀",
	ignoreCase: false,
},
	},
},
},
{
	name: "Arrow",
	pos: position{line: 218, col: 1, offset: 6214},
	expr: &choiceExpr{
	pos: position{line: 218, col: 10, offset: 6225},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 218, col: 10, offset: 6225},
	val: "->",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 218, col: 17, offset: 6232},
	val: "→",
	ignoreCase: false,
},
	},
},
},
{
	name: "Exponent",
	pos: position{line: 220, col: 1, offset: 6240},
	expr: &seqExpr{
	pos: position{line: 220, col: 12, offset: 6253},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 220, col: 12, offset: 6253},
	val: "e",
	ignoreCase: true,
},
&zeroOrOneExpr{
	pos: position{line: 220, col: 17, offset: 6258},
	expr: &charClassMatcher{
	pos: position{line: 220, col: 17, offset: 6258},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 220, col: 23, offset: 6264},
	expr: &ruleRefExpr{
	pos: position{line: 220, col: 23, offset: 6264},
	name: "Digit",
},
},
	},
},
},
{
	name: "NumericDoubleLiteral",
	pos: position{line: 222, col: 1, offset: 6272},
	expr: &actionExpr{
	pos: position{line: 222, col: 24, offset: 6297},
	run: (*parser).callonNumericDoubleLiteral1,
	expr: &seqExpr{
	pos: position{line: 222, col: 24, offset: 6297},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 222, col: 24, offset: 6297},
	expr: &charClassMatcher{
	pos: position{line: 222, col: 24, offset: 6297},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
},
&oneOrMoreExpr{
	pos: position{line: 222, col: 30, offset: 6303},
	expr: &ruleRefExpr{
	pos: position{line: 222, col: 30, offset: 6303},
	name: "Digit",
},
},
&choiceExpr{
	pos: position{line: 222, col: 39, offset: 6312},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 222, col: 39, offset: 6312},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 222, col: 39, offset: 6312},
	val: ".",
	ignoreCase: false,
},
&oneOrMoreExpr{
	pos: position{line: 222, col: 43, offset: 6316},
	expr: &ruleRefExpr{
	pos: position{line: 222, col: 43, offset: 6316},
	name: "Digit",
},
},
&zeroOrOneExpr{
	pos: position{line: 222, col: 50, offset: 6323},
	expr: &ruleRefExpr{
	pos: position{line: 222, col: 50, offset: 6323},
	name: "Exponent",
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 222, col: 62, offset: 6335},
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
	pos: position{line: 230, col: 1, offset: 6491},
	expr: &choiceExpr{
	pos: position{line: 230, col: 17, offset: 6509},
	alternatives: []interface{}{
&labeledExpr{
	pos: position{line: 230, col: 17, offset: 6509},
	label: "d",
	expr: &ruleRefExpr{
	pos: position{line: 230, col: 19, offset: 6511},
	name: "NumericDoubleLiteral",
},
},
&actionExpr{
	pos: position{line: 231, col: 5, offset: 6536},
	run: (*parser).callonDoubleLiteral4,
	expr: &ruleRefExpr{
	pos: position{line: 231, col: 5, offset: 6536},
	name: "Infinity",
},
},
&actionExpr{
	pos: position{line: 232, col: 5, offset: 6588},
	run: (*parser).callonDoubleLiteral6,
	expr: &seqExpr{
	pos: position{line: 232, col: 5, offset: 6588},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 232, col: 5, offset: 6588},
	val: "-",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 232, col: 9, offset: 6592},
	name: "Infinity",
},
	},
},
},
&actionExpr{
	pos: position{line: 233, col: 5, offset: 6645},
	run: (*parser).callonDoubleLiteral10,
	expr: &ruleRefExpr{
	pos: position{line: 233, col: 5, offset: 6645},
	name: "NaN",
},
},
	},
},
},
{
	name: "NaturalLiteral",
	pos: position{line: 235, col: 1, offset: 6688},
	expr: &actionExpr{
	pos: position{line: 235, col: 18, offset: 6707},
	run: (*parser).callonNaturalLiteral1,
	expr: &oneOrMoreExpr{
	pos: position{line: 235, col: 18, offset: 6707},
	expr: &ruleRefExpr{
	pos: position{line: 235, col: 18, offset: 6707},
	name: "Digit",
},
},
},
},
{
	name: "IntegerLiteral",
	pos: position{line: 240, col: 1, offset: 6796},
	expr: &actionExpr{
	pos: position{line: 240, col: 18, offset: 6815},
	run: (*parser).callonIntegerLiteral1,
	expr: &seqExpr{
	pos: position{line: 240, col: 18, offset: 6815},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 240, col: 18, offset: 6815},
	val: "[+-]",
	chars: []rune{'+','-',},
	ignoreCase: false,
	inverted: false,
},
&ruleRefExpr{
	pos: position{line: 240, col: 22, offset: 6819},
	name: "NaturalLiteral",
},
	},
},
},
},
{
	name: "DeBruijn",
	pos: position{line: 248, col: 1, offset: 6971},
	expr: &actionExpr{
	pos: position{line: 248, col: 12, offset: 6984},
	run: (*parser).callonDeBruijn1,
	expr: &seqExpr{
	pos: position{line: 248, col: 12, offset: 6984},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 248, col: 12, offset: 6984},
	name: "_",
},
&litMatcher{
	pos: position{line: 248, col: 14, offset: 6986},
	val: "@",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 248, col: 18, offset: 6990},
	name: "_",
},
&labeledExpr{
	pos: position{line: 248, col: 20, offset: 6992},
	label: "index",
	expr: &ruleRefExpr{
	pos: position{line: 248, col: 26, offset: 6998},
	name: "NaturalLiteral",
},
},
	},
},
},
},
{
	name: "Variable",
	pos: position{line: 250, col: 1, offset: 7054},
	expr: &actionExpr{
	pos: position{line: 250, col: 12, offset: 7067},
	run: (*parser).callonVariable1,
	expr: &seqExpr{
	pos: position{line: 250, col: 12, offset: 7067},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 250, col: 12, offset: 7067},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 250, col: 17, offset: 7072},
	name: "NonreservedLabel",
},
},
&labeledExpr{
	pos: position{line: 250, col: 34, offset: 7089},
	label: "index",
	expr: &zeroOrOneExpr{
	pos: position{line: 250, col: 40, offset: 7095},
	expr: &ruleRefExpr{
	pos: position{line: 250, col: 40, offset: 7095},
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
	pos: position{line: 258, col: 1, offset: 7258},
	expr: &choiceExpr{
	pos: position{line: 258, col: 14, offset: 7273},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 258, col: 14, offset: 7273},
	name: "Variable",
},
&ruleRefExpr{
	pos: position{line: 258, col: 25, offset: 7284},
	name: "Reserved",
},
	},
},
},
{
	name: "PathCharacter",
	pos: position{line: 260, col: 1, offset: 7294},
	expr: &choiceExpr{
	pos: position{line: 261, col: 6, offset: 7317},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 261, col: 6, offset: 7317},
	val: "!",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 262, col: 6, offset: 7329},
	val: "[\\x24-\\x27]",
	ranges: []rune{'$','\'',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 263, col: 6, offset: 7346},
	val: "[\\x2a-\\x2b]",
	ranges: []rune{'*','+',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 264, col: 6, offset: 7363},
	val: "[\\x2d-\\x2e]",
	ranges: []rune{'-','.',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 265, col: 6, offset: 7380},
	val: "[\\x30-\\x3b]",
	ranges: []rune{'0',';',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 266, col: 6, offset: 7397},
	val: "=",
	ignoreCase: false,
},
&charClassMatcher{
	pos: position{line: 267, col: 6, offset: 7409},
	val: "[\\x40-\\x5a]",
	ranges: []rune{'@','Z',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 268, col: 6, offset: 7426},
	val: "[\\x5e-\\x7a]",
	ranges: []rune{'^','z',},
	ignoreCase: false,
	inverted: false,
},
&litMatcher{
	pos: position{line: 269, col: 6, offset: 7443},
	val: "|",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 270, col: 6, offset: 7455},
	val: "~",
	ignoreCase: false,
},
	},
},
},
{
	name: "UnquotedPathComponent",
	pos: position{line: 272, col: 1, offset: 7463},
	expr: &actionExpr{
	pos: position{line: 272, col: 25, offset: 7489},
	run: (*parser).callonUnquotedPathComponent1,
	expr: &oneOrMoreExpr{
	pos: position{line: 272, col: 25, offset: 7489},
	expr: &ruleRefExpr{
	pos: position{line: 272, col: 25, offset: 7489},
	name: "PathCharacter",
},
},
},
},
{
	name: "PathComponent",
	pos: position{line: 274, col: 1, offset: 7536},
	expr: &actionExpr{
	pos: position{line: 274, col: 17, offset: 7554},
	run: (*parser).callonPathComponent1,
	expr: &seqExpr{
	pos: position{line: 274, col: 17, offset: 7554},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 274, col: 17, offset: 7554},
	val: "/",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 274, col: 21, offset: 7558},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 274, col: 23, offset: 7560},
	name: "UnquotedPathComponent",
},
},
	},
},
},
},
{
	name: "Path",
	pos: position{line: 276, col: 1, offset: 7601},
	expr: &actionExpr{
	pos: position{line: 276, col: 8, offset: 7610},
	run: (*parser).callonPath1,
	expr: &labeledExpr{
	pos: position{line: 276, col: 8, offset: 7610},
	label: "cs",
	expr: &oneOrMoreExpr{
	pos: position{line: 276, col: 11, offset: 7613},
	expr: &ruleRefExpr{
	pos: position{line: 276, col: 11, offset: 7613},
	name: "PathComponent",
},
},
},
},
},
{
	name: "Local",
	pos: position{line: 285, col: 1, offset: 7887},
	expr: &choiceExpr{
	pos: position{line: 285, col: 9, offset: 7897},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 285, col: 9, offset: 7897},
	name: "ParentPath",
},
&ruleRefExpr{
	pos: position{line: 285, col: 22, offset: 7910},
	name: "HerePath",
},
&ruleRefExpr{
	pos: position{line: 285, col: 33, offset: 7921},
	name: "HomePath",
},
&ruleRefExpr{
	pos: position{line: 285, col: 44, offset: 7932},
	name: "AbsolutePath",
},
	},
},
},
{
	name: "ParentPath",
	pos: position{line: 287, col: 1, offset: 7946},
	expr: &actionExpr{
	pos: position{line: 287, col: 14, offset: 7961},
	run: (*parser).callonParentPath1,
	expr: &seqExpr{
	pos: position{line: 287, col: 14, offset: 7961},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 287, col: 14, offset: 7961},
	val: "..",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 287, col: 19, offset: 7966},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 287, col: 21, offset: 7968},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HerePath",
	pos: position{line: 288, col: 1, offset: 8024},
	expr: &actionExpr{
	pos: position{line: 288, col: 12, offset: 8037},
	run: (*parser).callonHerePath1,
	expr: &seqExpr{
	pos: position{line: 288, col: 12, offset: 8037},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 288, col: 12, offset: 8037},
	val: ".",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 288, col: 16, offset: 8041},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 288, col: 18, offset: 8043},
	name: "Path",
},
},
	},
},
},
},
{
	name: "HomePath",
	pos: position{line: 289, col: 1, offset: 8082},
	expr: &actionExpr{
	pos: position{line: 289, col: 12, offset: 8095},
	run: (*parser).callonHomePath1,
	expr: &seqExpr{
	pos: position{line: 289, col: 12, offset: 8095},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 289, col: 12, offset: 8095},
	val: "~",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 289, col: 16, offset: 8099},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 289, col: 18, offset: 8101},
	name: "Path",
},
},
	},
},
},
},
{
	name: "AbsolutePath",
	pos: position{line: 290, col: 1, offset: 8156},
	expr: &actionExpr{
	pos: position{line: 290, col: 16, offset: 8173},
	run: (*parser).callonAbsolutePath1,
	expr: &labeledExpr{
	pos: position{line: 290, col: 16, offset: 8173},
	label: "p",
	expr: &ruleRefExpr{
	pos: position{line: 290, col: 18, offset: 8175},
	name: "Path",
},
},
},
},
{
	name: "Scheme",
	pos: position{line: 292, col: 1, offset: 8231},
	expr: &seqExpr{
	pos: position{line: 292, col: 10, offset: 8242},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 292, col: 10, offset: 8242},
	val: "http",
	ignoreCase: false,
},
&zeroOrOneExpr{
	pos: position{line: 292, col: 17, offset: 8249},
	expr: &litMatcher{
	pos: position{line: 292, col: 17, offset: 8249},
	val: "s",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "HttpRaw",
	pos: position{line: 294, col: 1, offset: 8255},
	expr: &actionExpr{
	pos: position{line: 294, col: 11, offset: 8267},
	run: (*parser).callonHttpRaw1,
	expr: &seqExpr{
	pos: position{line: 294, col: 11, offset: 8267},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 294, col: 11, offset: 8267},
	name: "Scheme",
},
&litMatcher{
	pos: position{line: 294, col: 18, offset: 8274},
	val: "://",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 294, col: 24, offset: 8280},
	name: "Authority",
},
&ruleRefExpr{
	pos: position{line: 294, col: 34, offset: 8290},
	name: "Path",
},
&zeroOrOneExpr{
	pos: position{line: 294, col: 39, offset: 8295},
	expr: &seqExpr{
	pos: position{line: 294, col: 41, offset: 8297},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 294, col: 41, offset: 8297},
	val: "?",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 294, col: 45, offset: 8301},
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
	pos: position{line: 296, col: 1, offset: 8358},
	expr: &seqExpr{
	pos: position{line: 296, col: 13, offset: 8372},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 296, col: 13, offset: 8372},
	expr: &seqExpr{
	pos: position{line: 296, col: 14, offset: 8373},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 296, col: 14, offset: 8373},
	name: "Userinfo",
},
&litMatcher{
	pos: position{line: 296, col: 23, offset: 8382},
	val: "@",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 296, col: 29, offset: 8388},
	name: "Host",
},
&zeroOrOneExpr{
	pos: position{line: 296, col: 34, offset: 8393},
	expr: &seqExpr{
	pos: position{line: 296, col: 35, offset: 8394},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 296, col: 35, offset: 8394},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 296, col: 39, offset: 8398},
	name: "Port",
},
	},
},
},
	},
},
},
{
	name: "Userinfo",
	pos: position{line: 298, col: 1, offset: 8406},
	expr: &zeroOrMoreExpr{
	pos: position{line: 298, col: 12, offset: 8419},
	expr: &choiceExpr{
	pos: position{line: 298, col: 14, offset: 8421},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 298, col: 14, offset: 8421},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 298, col: 27, offset: 8434},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 298, col: 40, offset: 8447},
	name: "SubDelims",
},
&litMatcher{
	pos: position{line: 298, col: 52, offset: 8459},
	val: ":",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "Host",
	pos: position{line: 300, col: 1, offset: 8467},
	expr: &choiceExpr{
	pos: position{line: 300, col: 8, offset: 8476},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 300, col: 8, offset: 8476},
	name: "IPLiteral",
},
&ruleRefExpr{
	pos: position{line: 300, col: 20, offset: 8488},
	name: "RegName",
},
	},
},
},
{
	name: "Port",
	pos: position{line: 302, col: 1, offset: 8497},
	expr: &zeroOrMoreExpr{
	pos: position{line: 302, col: 8, offset: 8506},
	expr: &ruleRefExpr{
	pos: position{line: 302, col: 8, offset: 8506},
	name: "Digit",
},
},
},
{
	name: "IPLiteral",
	pos: position{line: 304, col: 1, offset: 8514},
	expr: &seqExpr{
	pos: position{line: 304, col: 13, offset: 8528},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 304, col: 13, offset: 8528},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 304, col: 17, offset: 8532},
	name: "IPv6address",
},
&litMatcher{
	pos: position{line: 304, col: 29, offset: 8544},
	val: "]",
	ignoreCase: false,
},
	},
},
},
{
	name: "IPv6address",
	pos: position{line: 306, col: 1, offset: 8549},
	expr: &actionExpr{
	pos: position{line: 306, col: 15, offset: 8565},
	run: (*parser).callonIPv6address1,
	expr: &seqExpr{
	pos: position{line: 306, col: 15, offset: 8565},
	exprs: []interface{}{
&zeroOrMoreExpr{
	pos: position{line: 306, col: 15, offset: 8565},
	expr: &ruleRefExpr{
	pos: position{line: 306, col: 16, offset: 8566},
	name: "HexDig",
},
},
&litMatcher{
	pos: position{line: 306, col: 25, offset: 8575},
	val: ":",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 306, col: 29, offset: 8579},
	expr: &choiceExpr{
	pos: position{line: 306, col: 30, offset: 8580},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 306, col: 30, offset: 8580},
	name: "HexDig",
},
&litMatcher{
	pos: position{line: 306, col: 39, offset: 8589},
	val: ":",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 306, col: 45, offset: 8595},
	val: ".",
	ignoreCase: false,
},
	},
},
},
	},
},
},
},
{
	name: "RegName",
	pos: position{line: 312, col: 1, offset: 8749},
	expr: &zeroOrMoreExpr{
	pos: position{line: 312, col: 11, offset: 8761},
	expr: &choiceExpr{
	pos: position{line: 312, col: 12, offset: 8762},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 312, col: 12, offset: 8762},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 312, col: 25, offset: 8775},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 312, col: 38, offset: 8788},
	name: "SubDelims",
},
	},
},
},
},
{
	name: "PChar",
	pos: position{line: 314, col: 1, offset: 8801},
	expr: &choiceExpr{
	pos: position{line: 314, col: 9, offset: 8811},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 314, col: 9, offset: 8811},
	name: "Unreserved",
},
&ruleRefExpr{
	pos: position{line: 314, col: 22, offset: 8824},
	name: "PctEncoded",
},
&ruleRefExpr{
	pos: position{line: 314, col: 35, offset: 8837},
	name: "SubDelims",
},
&charClassMatcher{
	pos: position{line: 314, col: 47, offset: 8849},
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
	pos: position{line: 316, col: 1, offset: 8855},
	expr: &zeroOrMoreExpr{
	pos: position{line: 316, col: 9, offset: 8865},
	expr: &choiceExpr{
	pos: position{line: 316, col: 10, offset: 8866},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 316, col: 10, offset: 8866},
	name: "PChar",
},
&charClassMatcher{
	pos: position{line: 316, col: 18, offset: 8874},
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
	pos: position{line: 318, col: 1, offset: 8882},
	expr: &seqExpr{
	pos: position{line: 318, col: 14, offset: 8897},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 318, col: 14, offset: 8897},
	val: "%",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 318, col: 18, offset: 8901},
	name: "HexDig",
},
&ruleRefExpr{
	pos: position{line: 318, col: 25, offset: 8908},
	name: "HexDig",
},
	},
},
},
{
	name: "Unreserved",
	pos: position{line: 320, col: 1, offset: 8916},
	expr: &charClassMatcher{
	pos: position{line: 320, col: 14, offset: 8931},
	val: "[A-Za-z0-9._~-]",
	chars: []rune{'.','_','~','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SubDelims",
	pos: position{line: 322, col: 1, offset: 8948},
	expr: &choiceExpr{
	pos: position{line: 322, col: 13, offset: 8962},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 322, col: 13, offset: 8962},
	val: "!",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 322, col: 19, offset: 8968},
	val: "$",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 322, col: 25, offset: 8974},
	val: "&",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 322, col: 31, offset: 8980},
	val: "'",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 322, col: 37, offset: 8986},
	val: "(",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 322, col: 43, offset: 8992},
	val: ")",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 322, col: 49, offset: 8998},
	val: "*",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 322, col: 55, offset: 9004},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 322, col: 61, offset: 9010},
	val: ",",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 322, col: 67, offset: 9016},
	val: ";",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 322, col: 73, offset: 9022},
	val: "=",
	ignoreCase: false,
},
	},
},
},
{
	name: "Http",
	pos: position{line: 324, col: 1, offset: 9027},
	expr: &actionExpr{
	pos: position{line: 324, col: 8, offset: 9036},
	run: (*parser).callonHttp1,
	expr: &labeledExpr{
	pos: position{line: 324, col: 8, offset: 9036},
	label: "u",
	expr: &ruleRefExpr{
	pos: position{line: 324, col: 10, offset: 9038},
	name: "HttpRaw",
},
},
},
},
{
	name: "Env",
	pos: position{line: 326, col: 1, offset: 9083},
	expr: &actionExpr{
	pos: position{line: 326, col: 7, offset: 9091},
	run: (*parser).callonEnv1,
	expr: &seqExpr{
	pos: position{line: 326, col: 7, offset: 9091},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 326, col: 7, offset: 9091},
	val: "env:",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 326, col: 14, offset: 9098},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 326, col: 17, offset: 9101},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 326, col: 17, offset: 9101},
	name: "BashEnvironmentVariable",
},
&ruleRefExpr{
	pos: position{line: 326, col: 43, offset: 9127},
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
	pos: position{line: 328, col: 1, offset: 9172},
	expr: &actionExpr{
	pos: position{line: 328, col: 27, offset: 9200},
	run: (*parser).callonBashEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 328, col: 27, offset: 9200},
	exprs: []interface{}{
&charClassMatcher{
	pos: position{line: 328, col: 27, offset: 9200},
	val: "[A-Za-z_]",
	chars: []rune{'_',},
	ranges: []rune{'A','Z','a','z',},
	ignoreCase: false,
	inverted: false,
},
&zeroOrMoreExpr{
	pos: position{line: 328, col: 36, offset: 9209},
	expr: &charClassMatcher{
	pos: position{line: 328, col: 36, offset: 9209},
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
	pos: position{line: 332, col: 1, offset: 9265},
	expr: &actionExpr{
	pos: position{line: 332, col: 28, offset: 9294},
	run: (*parser).callonPosixEnvironmentVariable1,
	expr: &seqExpr{
	pos: position{line: 332, col: 28, offset: 9294},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 332, col: 28, offset: 9294},
	val: "\"",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 332, col: 32, offset: 9298},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 332, col: 34, offset: 9300},
	name: "PosixEnvironmentVariableContent",
},
},
&litMatcher{
	pos: position{line: 332, col: 66, offset: 9332},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "PosixEnvironmentVariableContent",
	pos: position{line: 336, col: 1, offset: 9357},
	expr: &actionExpr{
	pos: position{line: 336, col: 35, offset: 9393},
	run: (*parser).callonPosixEnvironmentVariableContent1,
	expr: &labeledExpr{
	pos: position{line: 336, col: 35, offset: 9393},
	label: "v",
	expr: &oneOrMoreExpr{
	pos: position{line: 336, col: 37, offset: 9395},
	expr: &ruleRefExpr{
	pos: position{line: 336, col: 37, offset: 9395},
	name: "PosixEnvironmentVariableCharacter",
},
},
},
},
},
{
	name: "PosixEnvironmentVariableCharacter",
	pos: position{line: 345, col: 1, offset: 9608},
	expr: &choiceExpr{
	pos: position{line: 346, col: 7, offset: 9652},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 346, col: 7, offset: 9652},
	run: (*parser).callonPosixEnvironmentVariableCharacter2,
	expr: &litMatcher{
	pos: position{line: 346, col: 7, offset: 9652},
	val: "\\\"",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 347, col: 7, offset: 9692},
	run: (*parser).callonPosixEnvironmentVariableCharacter4,
	expr: &litMatcher{
	pos: position{line: 347, col: 7, offset: 9692},
	val: "\\\\",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 348, col: 7, offset: 9732},
	run: (*parser).callonPosixEnvironmentVariableCharacter6,
	expr: &litMatcher{
	pos: position{line: 348, col: 7, offset: 9732},
	val: "\\a",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 349, col: 7, offset: 9772},
	run: (*parser).callonPosixEnvironmentVariableCharacter8,
	expr: &litMatcher{
	pos: position{line: 349, col: 7, offset: 9772},
	val: "\\b",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 350, col: 7, offset: 9812},
	run: (*parser).callonPosixEnvironmentVariableCharacter10,
	expr: &litMatcher{
	pos: position{line: 350, col: 7, offset: 9812},
	val: "\\f",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 351, col: 7, offset: 9852},
	run: (*parser).callonPosixEnvironmentVariableCharacter12,
	expr: &litMatcher{
	pos: position{line: 351, col: 7, offset: 9852},
	val: "\\n",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 352, col: 7, offset: 9892},
	run: (*parser).callonPosixEnvironmentVariableCharacter14,
	expr: &litMatcher{
	pos: position{line: 352, col: 7, offset: 9892},
	val: "\\r",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 353, col: 7, offset: 9932},
	run: (*parser).callonPosixEnvironmentVariableCharacter16,
	expr: &litMatcher{
	pos: position{line: 353, col: 7, offset: 9932},
	val: "\\t",
	ignoreCase: false,
},
},
&actionExpr{
	pos: position{line: 354, col: 7, offset: 9972},
	run: (*parser).callonPosixEnvironmentVariableCharacter18,
	expr: &litMatcher{
	pos: position{line: 354, col: 7, offset: 9972},
	val: "\\v",
	ignoreCase: false,
},
},
&charClassMatcher{
	pos: position{line: 355, col: 7, offset: 10012},
	val: "[\\x20-\\x21]",
	ranges: []rune{' ','!',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 356, col: 7, offset: 10030},
	val: "[\\x23-\\x3c]",
	ranges: []rune{'#','<',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 357, col: 7, offset: 10048},
	val: "[\\x3e-\\x5b]",
	ranges: []rune{'>','[',},
	ignoreCase: false,
	inverted: false,
},
&charClassMatcher{
	pos: position{line: 358, col: 7, offset: 10066},
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
	pos: position{line: 360, col: 1, offset: 10079},
	expr: &choiceExpr{
	pos: position{line: 360, col: 14, offset: 10094},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 360, col: 14, offset: 10094},
	name: "Missing",
},
&ruleRefExpr{
	pos: position{line: 360, col: 24, offset: 10104},
	name: "Local",
},
&ruleRefExpr{
	pos: position{line: 360, col: 32, offset: 10112},
	name: "Http",
},
&ruleRefExpr{
	pos: position{line: 360, col: 39, offset: 10119},
	name: "Env",
},
	},
},
},
{
	name: "ImportHashed",
	pos: position{line: 362, col: 1, offset: 10124},
	expr: &actionExpr{
	pos: position{line: 362, col: 16, offset: 10141},
	run: (*parser).callonImportHashed1,
	expr: &labeledExpr{
	pos: position{line: 362, col: 16, offset: 10141},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 362, col: 18, offset: 10143},
	name: "ImportType",
},
},
},
},
{
	name: "Import",
	pos: position{line: 364, col: 1, offset: 10212},
	expr: &choiceExpr{
	pos: position{line: 364, col: 10, offset: 10223},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 364, col: 10, offset: 10223},
	run: (*parser).callonImport2,
	expr: &seqExpr{
	pos: position{line: 364, col: 10, offset: 10223},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 364, col: 10, offset: 10223},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 364, col: 12, offset: 10225},
	name: "ImportHashed",
},
},
&ruleRefExpr{
	pos: position{line: 364, col: 25, offset: 10238},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 364, col: 27, offset: 10240},
	name: "As",
},
&ruleRefExpr{
	pos: position{line: 364, col: 30, offset: 10243},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 364, col: 33, offset: 10246},
	name: "Text",
},
	},
},
},
&actionExpr{
	pos: position{line: 365, col: 10, offset: 10343},
	run: (*parser).callonImport10,
	expr: &labeledExpr{
	pos: position{line: 365, col: 10, offset: 10343},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 365, col: 12, offset: 10345},
	name: "ImportHashed",
},
},
},
	},
},
},
{
	name: "LetBinding",
	pos: position{line: 368, col: 1, offset: 10440},
	expr: &actionExpr{
	pos: position{line: 368, col: 14, offset: 10455},
	run: (*parser).callonLetBinding1,
	expr: &seqExpr{
	pos: position{line: 368, col: 14, offset: 10455},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 368, col: 14, offset: 10455},
	name: "Let",
},
&ruleRefExpr{
	pos: position{line: 368, col: 18, offset: 10459},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 368, col: 21, offset: 10462},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 368, col: 27, offset: 10468},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 368, col: 44, offset: 10485},
	name: "_",
},
&labeledExpr{
	pos: position{line: 368, col: 46, offset: 10487},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 368, col: 48, offset: 10489},
	expr: &seqExpr{
	pos: position{line: 368, col: 49, offset: 10490},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 368, col: 49, offset: 10490},
	name: "Annotation",
},
&ruleRefExpr{
	pos: position{line: 368, col: 60, offset: 10501},
	name: "_",
},
	},
},
},
},
&litMatcher{
	pos: position{line: 369, col: 13, offset: 10517},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 369, col: 17, offset: 10521},
	name: "_",
},
&labeledExpr{
	pos: position{line: 369, col: 19, offset: 10523},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 369, col: 21, offset: 10525},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 369, col: 32, offset: 10536},
	name: "_",
},
	},
},
},
},
{
	name: "Expression",
	pos: position{line: 384, col: 1, offset: 10845},
	expr: &choiceExpr{
	pos: position{line: 385, col: 7, offset: 10866},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 385, col: 7, offset: 10866},
	run: (*parser).callonExpression2,
	expr: &seqExpr{
	pos: position{line: 385, col: 7, offset: 10866},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 385, col: 7, offset: 10866},
	name: "Lambda",
},
&ruleRefExpr{
	pos: position{line: 385, col: 14, offset: 10873},
	name: "_",
},
&litMatcher{
	pos: position{line: 385, col: 16, offset: 10875},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 385, col: 20, offset: 10879},
	name: "_",
},
&labeledExpr{
	pos: position{line: 385, col: 22, offset: 10881},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 385, col: 28, offset: 10887},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 385, col: 45, offset: 10904},
	name: "_",
},
&litMatcher{
	pos: position{line: 385, col: 47, offset: 10906},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 385, col: 51, offset: 10910},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 385, col: 54, offset: 10913},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 385, col: 56, offset: 10915},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 385, col: 67, offset: 10926},
	name: "_",
},
&litMatcher{
	pos: position{line: 385, col: 69, offset: 10928},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 385, col: 73, offset: 10932},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 385, col: 75, offset: 10934},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 385, col: 81, offset: 10940},
	name: "_",
},
&labeledExpr{
	pos: position{line: 385, col: 83, offset: 10942},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 385, col: 88, offset: 10947},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 388, col: 7, offset: 11064},
	run: (*parser).callonExpression22,
	expr: &seqExpr{
	pos: position{line: 388, col: 7, offset: 11064},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 388, col: 7, offset: 11064},
	name: "If",
},
&ruleRefExpr{
	pos: position{line: 388, col: 10, offset: 11067},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 388, col: 13, offset: 11070},
	label: "cond",
	expr: &ruleRefExpr{
	pos: position{line: 388, col: 18, offset: 11075},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 388, col: 29, offset: 11086},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 388, col: 31, offset: 11088},
	name: "Then",
},
&ruleRefExpr{
	pos: position{line: 388, col: 36, offset: 11093},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 388, col: 39, offset: 11096},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 388, col: 41, offset: 11098},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 388, col: 52, offset: 11109},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 388, col: 54, offset: 11111},
	name: "Else",
},
&ruleRefExpr{
	pos: position{line: 388, col: 59, offset: 11116},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 388, col: 62, offset: 11119},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 388, col: 64, offset: 11121},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 391, col: 7, offset: 11207},
	run: (*parser).callonExpression38,
	expr: &seqExpr{
	pos: position{line: 391, col: 7, offset: 11207},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 391, col: 7, offset: 11207},
	label: "bindings",
	expr: &oneOrMoreExpr{
	pos: position{line: 391, col: 16, offset: 11216},
	expr: &ruleRefExpr{
	pos: position{line: 391, col: 16, offset: 11216},
	name: "LetBinding",
},
},
},
&ruleRefExpr{
	pos: position{line: 391, col: 28, offset: 11228},
	name: "In",
},
&ruleRefExpr{
	pos: position{line: 391, col: 31, offset: 11231},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 391, col: 34, offset: 11234},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 391, col: 36, offset: 11236},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 398, col: 7, offset: 11476},
	run: (*parser).callonExpression47,
	expr: &seqExpr{
	pos: position{line: 398, col: 7, offset: 11476},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 398, col: 7, offset: 11476},
	name: "Forall",
},
&ruleRefExpr{
	pos: position{line: 398, col: 14, offset: 11483},
	name: "_",
},
&litMatcher{
	pos: position{line: 398, col: 16, offset: 11485},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 398, col: 20, offset: 11489},
	name: "_",
},
&labeledExpr{
	pos: position{line: 398, col: 22, offset: 11491},
	label: "label",
	expr: &ruleRefExpr{
	pos: position{line: 398, col: 28, offset: 11497},
	name: "NonreservedLabel",
},
},
&ruleRefExpr{
	pos: position{line: 398, col: 45, offset: 11514},
	name: "_",
},
&litMatcher{
	pos: position{line: 398, col: 47, offset: 11516},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 398, col: 51, offset: 11520},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 398, col: 54, offset: 11523},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 398, col: 56, offset: 11525},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 398, col: 67, offset: 11536},
	name: "_",
},
&litMatcher{
	pos: position{line: 398, col: 69, offset: 11538},
	val: ")",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 398, col: 73, offset: 11542},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 398, col: 75, offset: 11544},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 398, col: 81, offset: 11550},
	name: "_",
},
&labeledExpr{
	pos: position{line: 398, col: 83, offset: 11552},
	label: "body",
	expr: &ruleRefExpr{
	pos: position{line: 398, col: 88, offset: 11557},
	name: "Expression",
},
},
	},
},
},
&actionExpr{
	pos: position{line: 401, col: 7, offset: 11666},
	run: (*parser).callonExpression67,
	expr: &seqExpr{
	pos: position{line: 401, col: 7, offset: 11666},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 401, col: 7, offset: 11666},
	label: "o",
	expr: &ruleRefExpr{
	pos: position{line: 401, col: 9, offset: 11668},
	name: "OperatorExpression",
},
},
&ruleRefExpr{
	pos: position{line: 401, col: 28, offset: 11687},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 401, col: 30, offset: 11689},
	name: "Arrow",
},
&ruleRefExpr{
	pos: position{line: 401, col: 36, offset: 11695},
	name: "_",
},
&labeledExpr{
	pos: position{line: 401, col: 38, offset: 11697},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 401, col: 40, offset: 11699},
	name: "Expression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 402, col: 7, offset: 11759},
	name: "AnnotatedExpression",
},
	},
},
},
{
	name: "Annotation",
	pos: position{line: 404, col: 1, offset: 11780},
	expr: &actionExpr{
	pos: position{line: 404, col: 14, offset: 11795},
	run: (*parser).callonAnnotation1,
	expr: &seqExpr{
	pos: position{line: 404, col: 14, offset: 11795},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 404, col: 14, offset: 11795},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 404, col: 18, offset: 11799},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 404, col: 21, offset: 11802},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 404, col: 23, offset: 11804},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "AnnotatedExpression",
	pos: position{line: 406, col: 1, offset: 11834},
	expr: &choiceExpr{
	pos: position{line: 407, col: 5, offset: 11862},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 407, col: 5, offset: 11862},
	name: "EmptyList",
},
&actionExpr{
	pos: position{line: 408, col: 5, offset: 11876},
	run: (*parser).callonAnnotatedExpression3,
	expr: &seqExpr{
	pos: position{line: 408, col: 5, offset: 11876},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 408, col: 5, offset: 11876},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 408, col: 7, offset: 11878},
	name: "OperatorExpression",
},
},
&labeledExpr{
	pos: position{line: 408, col: 26, offset: 11897},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 408, col: 28, offset: 11899},
	expr: &seqExpr{
	pos: position{line: 408, col: 29, offset: 11900},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 408, col: 29, offset: 11900},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 408, col: 31, offset: 11902},
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
	pos: position{line: 413, col: 1, offset: 12027},
	expr: &actionExpr{
	pos: position{line: 413, col: 13, offset: 12041},
	run: (*parser).callonEmptyList1,
	expr: &seqExpr{
	pos: position{line: 413, col: 13, offset: 12041},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 413, col: 13, offset: 12041},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 413, col: 17, offset: 12045},
	name: "_",
},
&litMatcher{
	pos: position{line: 413, col: 19, offset: 12047},
	val: "]",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 413, col: 23, offset: 12051},
	name: "_",
},
&litMatcher{
	pos: position{line: 413, col: 25, offset: 12053},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 413, col: 29, offset: 12057},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 413, col: 32, offset: 12060},
	name: "List",
},
&ruleRefExpr{
	pos: position{line: 413, col: 37, offset: 12065},
	name: "_",
},
&labeledExpr{
	pos: position{line: 413, col: 39, offset: 12067},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 413, col: 41, offset: 12069},
	name: "ImportExpression",
},
},
	},
},
},
},
{
	name: "OperatorExpression",
	pos: position{line: 417, col: 1, offset: 12132},
	expr: &ruleRefExpr{
	pos: position{line: 417, col: 22, offset: 12155},
	name: "ImportAltExpression",
},
},
{
	name: "ImportAltExpression",
	pos: position{line: 419, col: 1, offset: 12176},
	expr: &ruleRefExpr{
	pos: position{line: 419, col: 23, offset: 12200},
	name: "PlusExpression",
},
},
{
	name: "MorePlus",
	pos: position{line: 421, col: 1, offset: 12216},
	expr: &actionExpr{
	pos: position{line: 421, col: 12, offset: 12229},
	run: (*parser).callonMorePlus1,
	expr: &seqExpr{
	pos: position{line: 421, col: 12, offset: 12229},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 421, col: 12, offset: 12229},
	name: "_",
},
&litMatcher{
	pos: position{line: 421, col: 14, offset: 12231},
	val: "+",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 421, col: 18, offset: 12235},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 421, col: 21, offset: 12238},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 421, col: 23, offset: 12240},
	name: "ListAppendExpression",
},
},
	},
},
},
},
{
	name: "PlusExpression",
	pos: position{line: 422, col: 1, offset: 12279},
	expr: &actionExpr{
	pos: position{line: 423, col: 7, offset: 12304},
	run: (*parser).callonPlusExpression1,
	expr: &seqExpr{
	pos: position{line: 423, col: 7, offset: 12304},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 423, col: 7, offset: 12304},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 423, col: 13, offset: 12310},
	name: "ListAppendExpression",
},
},
&labeledExpr{
	pos: position{line: 423, col: 34, offset: 12331},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 423, col: 39, offset: 12336},
	expr: &ruleRefExpr{
	pos: position{line: 423, col: 39, offset: 12336},
	name: "MorePlus",
},
},
},
	},
},
},
},
{
	name: "MoreListAppend",
	pos: position{line: 432, col: 1, offset: 12564},
	expr: &actionExpr{
	pos: position{line: 432, col: 18, offset: 12583},
	run: (*parser).callonMoreListAppend1,
	expr: &seqExpr{
	pos: position{line: 432, col: 18, offset: 12583},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 432, col: 18, offset: 12583},
	name: "_",
},
&litMatcher{
	pos: position{line: 432, col: 20, offset: 12585},
	val: "#",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 432, col: 24, offset: 12589},
	name: "_",
},
&labeledExpr{
	pos: position{line: 432, col: 26, offset: 12591},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 432, col: 28, offset: 12593},
	name: "TimesExpression",
},
},
	},
},
},
},
{
	name: "ListAppendExpression",
	pos: position{line: 433, col: 1, offset: 12627},
	expr: &actionExpr{
	pos: position{line: 433, col: 24, offset: 12652},
	run: (*parser).callonListAppendExpression1,
	expr: &seqExpr{
	pos: position{line: 433, col: 24, offset: 12652},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 433, col: 24, offset: 12652},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 433, col: 30, offset: 12658},
	name: "TimesExpression",
},
},
&labeledExpr{
	pos: position{line: 433, col: 46, offset: 12674},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 433, col: 51, offset: 12679},
	expr: &ruleRefExpr{
	pos: position{line: 433, col: 51, offset: 12679},
	name: "MoreListAppend",
},
},
},
	},
},
},
},
{
	name: "MoreTimes",
	pos: position{line: 442, col: 1, offset: 12870},
	expr: &actionExpr{
	pos: position{line: 442, col: 13, offset: 12884},
	run: (*parser).callonMoreTimes1,
	expr: &seqExpr{
	pos: position{line: 442, col: 13, offset: 12884},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 442, col: 13, offset: 12884},
	name: "_",
},
&litMatcher{
	pos: position{line: 442, col: 15, offset: 12886},
	val: "*",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 442, col: 19, offset: 12890},
	name: "_",
},
&labeledExpr{
	pos: position{line: 442, col: 21, offset: 12892},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 442, col: 23, offset: 12894},
	name: "ApplicationExpression",
},
},
	},
},
},
},
{
	name: "TimesExpression",
	pos: position{line: 443, col: 1, offset: 12934},
	expr: &actionExpr{
	pos: position{line: 444, col: 7, offset: 12960},
	run: (*parser).callonTimesExpression1,
	expr: &seqExpr{
	pos: position{line: 444, col: 7, offset: 12960},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 444, col: 7, offset: 12960},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 444, col: 13, offset: 12966},
	name: "ApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 444, col: 35, offset: 12988},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 444, col: 40, offset: 12993},
	expr: &ruleRefExpr{
	pos: position{line: 444, col: 40, offset: 12993},
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
	pos: position{line: 453, col: 1, offset: 13223},
	expr: &actionExpr{
	pos: position{line: 453, col: 25, offset: 13249},
	run: (*parser).callonApplicationExpression1,
	expr: &seqExpr{
	pos: position{line: 453, col: 25, offset: 13249},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 453, col: 25, offset: 13249},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 453, col: 27, offset: 13251},
	name: "FirstApplicationExpression",
},
},
&labeledExpr{
	pos: position{line: 453, col: 54, offset: 13278},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 453, col: 59, offset: 13283},
	expr: &seqExpr{
	pos: position{line: 453, col: 60, offset: 13284},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 453, col: 60, offset: 13284},
	name: "_1",
},
&ruleRefExpr{
	pos: position{line: 453, col: 63, offset: 13287},
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
	pos: position{line: 462, col: 1, offset: 13537},
	expr: &choiceExpr{
	pos: position{line: 463, col: 8, offset: 13575},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 463, col: 8, offset: 13575},
	run: (*parser).callonFirstApplicationExpression2,
	expr: &seqExpr{
	pos: position{line: 463, col: 8, offset: 13575},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 463, col: 8, offset: 13575},
	name: "Some",
},
&ruleRefExpr{
	pos: position{line: 463, col: 13, offset: 13580},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 463, col: 16, offset: 13583},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 463, col: 18, offset: 13585},
	name: "ImportExpression",
},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 464, col: 8, offset: 13640},
	name: "ImportExpression",
},
	},
},
},
{
	name: "ImportExpression",
	pos: position{line: 466, col: 1, offset: 13658},
	expr: &choiceExpr{
	pos: position{line: 466, col: 20, offset: 13679},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 466, col: 20, offset: 13679},
	name: "Import",
},
&ruleRefExpr{
	pos: position{line: 466, col: 29, offset: 13688},
	name: "SelectorExpression",
},
	},
},
},
{
	name: "SelectorExpression",
	pos: position{line: 468, col: 1, offset: 13708},
	expr: &actionExpr{
	pos: position{line: 468, col: 22, offset: 13731},
	run: (*parser).callonSelectorExpression1,
	expr: &seqExpr{
	pos: position{line: 468, col: 22, offset: 13731},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 468, col: 22, offset: 13731},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 468, col: 24, offset: 13733},
	name: "PrimitiveExpression",
},
},
&labeledExpr{
	pos: position{line: 468, col: 44, offset: 13753},
	label: "ls",
	expr: &zeroOrMoreExpr{
	pos: position{line: 468, col: 47, offset: 13756},
	expr: &seqExpr{
	pos: position{line: 468, col: 48, offset: 13757},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 468, col: 48, offset: 13757},
	name: "_",
},
&litMatcher{
	pos: position{line: 468, col: 50, offset: 13759},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 468, col: 54, offset: 13763},
	name: "_",
},
&ruleRefExpr{
	pos: position{line: 468, col: 56, offset: 13765},
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
	pos: position{line: 478, col: 1, offset: 13998},
	expr: &choiceExpr{
	pos: position{line: 479, col: 7, offset: 14028},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 479, col: 7, offset: 14028},
	name: "DoubleLiteral",
},
&ruleRefExpr{
	pos: position{line: 480, col: 7, offset: 14048},
	name: "NaturalLiteral",
},
&ruleRefExpr{
	pos: position{line: 481, col: 7, offset: 14069},
	name: "IntegerLiteral",
},
&ruleRefExpr{
	pos: position{line: 482, col: 7, offset: 14090},
	name: "TextLiteral",
},
&actionExpr{
	pos: position{line: 483, col: 7, offset: 14108},
	run: (*parser).callonPrimitiveExpression6,
	expr: &seqExpr{
	pos: position{line: 483, col: 7, offset: 14108},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 483, col: 7, offset: 14108},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 483, col: 11, offset: 14112},
	name: "_",
},
&labeledExpr{
	pos: position{line: 483, col: 13, offset: 14114},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 483, col: 15, offset: 14116},
	name: "RecordTypeOrLiteral",
},
},
&ruleRefExpr{
	pos: position{line: 483, col: 35, offset: 14136},
	name: "_",
},
&litMatcher{
	pos: position{line: 483, col: 37, offset: 14138},
	val: "}",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 484, col: 7, offset: 14166},
	name: "NonEmptyListLiteral",
},
&ruleRefExpr{
	pos: position{line: 485, col: 7, offset: 14192},
	name: "Identifier",
},
&actionExpr{
	pos: position{line: 486, col: 7, offset: 14209},
	run: (*parser).callonPrimitiveExpression16,
	expr: &seqExpr{
	pos: position{line: 486, col: 7, offset: 14209},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 486, col: 7, offset: 14209},
	val: "(",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 486, col: 11, offset: 14213},
	name: "_",
},
&labeledExpr{
	pos: position{line: 486, col: 14, offset: 14216},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 486, col: 16, offset: 14218},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 486, col: 27, offset: 14229},
	name: "_",
},
&litMatcher{
	pos: position{line: 486, col: 29, offset: 14231},
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
	pos: position{line: 488, col: 1, offset: 14254},
	expr: &choiceExpr{
	pos: position{line: 489, col: 7, offset: 14284},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 489, col: 7, offset: 14284},
	run: (*parser).callonRecordTypeOrLiteral2,
	expr: &litMatcher{
	pos: position{line: 489, col: 7, offset: 14284},
	val: "=",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 490, col: 7, offset: 14339},
	name: "NonEmptyRecordType",
},
&ruleRefExpr{
	pos: position{line: 491, col: 7, offset: 14364},
	name: "NonEmptyRecordLiteral",
},
&actionExpr{
	pos: position{line: 492, col: 7, offset: 14392},
	run: (*parser).callonRecordTypeOrLiteral6,
	expr: &litMatcher{
	pos: position{line: 492, col: 7, offset: 14392},
	val: "",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "RecordTypeField",
	pos: position{line: 494, col: 1, offset: 14438},
	expr: &actionExpr{
	pos: position{line: 494, col: 19, offset: 14458},
	run: (*parser).callonRecordTypeField1,
	expr: &seqExpr{
	pos: position{line: 494, col: 19, offset: 14458},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 494, col: 19, offset: 14458},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 494, col: 24, offset: 14463},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 494, col: 33, offset: 14472},
	name: "_",
},
&litMatcher{
	pos: position{line: 494, col: 35, offset: 14474},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 494, col: 39, offset: 14478},
	name: "_1",
},
&labeledExpr{
	pos: position{line: 494, col: 42, offset: 14481},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 494, col: 47, offset: 14486},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordType",
	pos: position{line: 497, col: 1, offset: 14543},
	expr: &actionExpr{
	pos: position{line: 497, col: 18, offset: 14562},
	run: (*parser).callonMoreRecordType1,
	expr: &seqExpr{
	pos: position{line: 497, col: 18, offset: 14562},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 497, col: 18, offset: 14562},
	name: "_",
},
&litMatcher{
	pos: position{line: 497, col: 20, offset: 14564},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 497, col: 24, offset: 14568},
	name: "_",
},
&labeledExpr{
	pos: position{line: 497, col: 26, offset: 14570},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 497, col: 28, offset: 14572},
	name: "RecordTypeField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordType",
	pos: position{line: 498, col: 1, offset: 14604},
	expr: &actionExpr{
	pos: position{line: 499, col: 7, offset: 14633},
	run: (*parser).callonNonEmptyRecordType1,
	expr: &seqExpr{
	pos: position{line: 499, col: 7, offset: 14633},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 499, col: 7, offset: 14633},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 499, col: 13, offset: 14639},
	name: "RecordTypeField",
},
},
&labeledExpr{
	pos: position{line: 499, col: 29, offset: 14655},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 499, col: 34, offset: 14660},
	expr: &ruleRefExpr{
	pos: position{line: 499, col: 34, offset: 14660},
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
	pos: position{line: 509, col: 1, offset: 15056},
	expr: &actionExpr{
	pos: position{line: 509, col: 22, offset: 15079},
	run: (*parser).callonRecordLiteralField1,
	expr: &seqExpr{
	pos: position{line: 509, col: 22, offset: 15079},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 509, col: 22, offset: 15079},
	label: "name",
	expr: &ruleRefExpr{
	pos: position{line: 509, col: 27, offset: 15084},
	name: "AnyLabel",
},
},
&ruleRefExpr{
	pos: position{line: 509, col: 36, offset: 15093},
	name: "_",
},
&litMatcher{
	pos: position{line: 509, col: 38, offset: 15095},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 509, col: 42, offset: 15099},
	name: "_",
},
&labeledExpr{
	pos: position{line: 509, col: 44, offset: 15101},
	label: "expr",
	expr: &ruleRefExpr{
	pos: position{line: 509, col: 49, offset: 15106},
	name: "Expression",
},
},
	},
},
},
},
{
	name: "MoreRecordLiteral",
	pos: position{line: 512, col: 1, offset: 15163},
	expr: &actionExpr{
	pos: position{line: 512, col: 21, offset: 15185},
	run: (*parser).callonMoreRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 512, col: 21, offset: 15185},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 512, col: 21, offset: 15185},
	name: "_",
},
&litMatcher{
	pos: position{line: 512, col: 23, offset: 15187},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 512, col: 27, offset: 15191},
	name: "_",
},
&labeledExpr{
	pos: position{line: 512, col: 29, offset: 15193},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 512, col: 31, offset: 15195},
	name: "RecordLiteralField",
},
},
	},
},
},
},
{
	name: "NonEmptyRecordLiteral",
	pos: position{line: 513, col: 1, offset: 15230},
	expr: &actionExpr{
	pos: position{line: 514, col: 7, offset: 15262},
	run: (*parser).callonNonEmptyRecordLiteral1,
	expr: &seqExpr{
	pos: position{line: 514, col: 7, offset: 15262},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 514, col: 7, offset: 15262},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 514, col: 13, offset: 15268},
	name: "RecordLiteralField",
},
},
&labeledExpr{
	pos: position{line: 514, col: 32, offset: 15287},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 514, col: 37, offset: 15292},
	expr: &ruleRefExpr{
	pos: position{line: 514, col: 37, offset: 15292},
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
	pos: position{line: 524, col: 1, offset: 15694},
	expr: &actionExpr{
	pos: position{line: 524, col: 12, offset: 15707},
	run: (*parser).callonMoreList1,
	expr: &seqExpr{
	pos: position{line: 524, col: 12, offset: 15707},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 524, col: 12, offset: 15707},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 524, col: 16, offset: 15711},
	name: "_",
},
&labeledExpr{
	pos: position{line: 524, col: 18, offset: 15713},
	label: "e",
	expr: &ruleRefExpr{
	pos: position{line: 524, col: 20, offset: 15715},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 524, col: 31, offset: 15726},
	name: "_",
},
	},
},
},
},
{
	name: "NonEmptyListLiteral",
	pos: position{line: 526, col: 1, offset: 15745},
	expr: &actionExpr{
	pos: position{line: 527, col: 7, offset: 15775},
	run: (*parser).callonNonEmptyListLiteral1,
	expr: &seqExpr{
	pos: position{line: 527, col: 7, offset: 15775},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 527, col: 7, offset: 15775},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 527, col: 11, offset: 15779},
	name: "_",
},
&labeledExpr{
	pos: position{line: 527, col: 13, offset: 15781},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 527, col: 19, offset: 15787},
	name: "Expression",
},
},
&ruleRefExpr{
	pos: position{line: 527, col: 30, offset: 15798},
	name: "_",
},
&labeledExpr{
	pos: position{line: 527, col: 32, offset: 15800},
	label: "rest",
	expr: &zeroOrMoreExpr{
	pos: position{line: 527, col: 37, offset: 15805},
	expr: &ruleRefExpr{
	pos: position{line: 527, col: 37, offset: 15805},
	name: "MoreList",
},
},
},
&litMatcher{
	pos: position{line: 527, col: 47, offset: 15815},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "EOF",
	pos: position{line: 537, col: 1, offset: 16091},
	expr: &notExpr{
	pos: position{line: 537, col: 7, offset: 16099},
	expr: &anyMatcher{
	line: 537, col: 8, offset: 16100,
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
 return url.ParseRequestURI(string(c.text)) 
}

func (p *parser) callonHttpRaw1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHttpRaw1()
}

func (c *current) onIPv6address1() (interface{}, error) {
    addr := net.ParseIP(string(c.text))
    if addr == nil { return nil, errors.New("Malformed IPv6 address") }
    return string(c.text), nil
}

func (p *parser) callonIPv6address1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIPv6address1()
}

func (c *current) onHttp1(u interface{}) (interface{}, error) {
 return MakeRemote(u.(*url.URL)) 
}

func (p *parser) callonHttp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHttp1(stack["u"])
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

func (c *current) onMoreListAppend1(e interface{}) (interface{}, error) {
 return e, nil 
}

func (p *parser) callonMoreListAppend1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMoreListAppend1(stack["e"])
}

func (c *current) onListAppendExpression1(first, rest interface{}) (interface{}, error) {
    a := first.(Expr)
    if rest == nil { return a, nil }
    for _, b := range rest.([]interface{}) {
        a = ListAppend{L: a, R: b.(Expr)}
    }
    return a, nil
}

func (p *parser) callonListAppendExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onListAppendExpression1(stack["first"], stack["rest"])
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

