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

	"github.com/philandstuff/dhall-golang/ast"
)

var g = &grammar{
	rules: []*rule{
		{
			name: "CompleteExpression",
			pos:  position{line: 7, col: 1, offset: 71},
			expr: &actionExpr{
				pos: position{line: 7, col: 22, offset: 94},
				run: (*parser).callonCompleteExpression1,
				expr: &seqExpr{
					pos: position{line: 7, col: 22, offset: 94},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 7, col: 22, offset: 94},
							name: "Whitespace",
						},
						&labeledExpr{
							pos:   position{line: 7, col: 33, offset: 105},
							label: "e",
							expr: &ruleRefExpr{
								pos:  position{line: 7, col: 35, offset: 107},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 7, col: 46, offset: 118},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 7, col: 57, offset: 129},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "EOL",
			pos:  position{line: 9, col: 1, offset: 152},
			expr: &choiceExpr{
				pos: position{line: 9, col: 7, offset: 160},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 9, col: 7, offset: 160},
						val:        "\n",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 9, col: 14, offset: 167},
						val:        "\r\n",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "NotEOL",
			pos:  position{line: 11, col: 1, offset: 175},
			expr: &charClassMatcher{
				pos:        position{line: 11, col: 10, offset: 186},
				val:        "[\\t\\u0020-\\U0010ffff]",
				chars:      []rune{'\t'},
				ranges:     []rune{' ', '\U0010ffff'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "LineComment",
			pos:  position{line: 13, col: 1, offset: 209},
			expr: &actionExpr{
				pos: position{line: 13, col: 15, offset: 225},
				run: (*parser).callonLineComment1,
				expr: &seqExpr{
					pos: position{line: 13, col: 15, offset: 225},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 13, col: 15, offset: 225},
							val:        "--",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 13, col: 20, offset: 230},
							label: "content",
							expr: &actionExpr{
								pos: position{line: 13, col: 29, offset: 239},
								run: (*parser).callonLineComment5,
								expr: &zeroOrMoreExpr{
									pos: position{line: 13, col: 29, offset: 239},
									expr: &ruleRefExpr{
										pos:  position{line: 13, col: 29, offset: 239},
										name: "NotEOL",
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 13, col: 68, offset: 278},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "WhitespaceChunk",
			pos:  position{line: 15, col: 1, offset: 307},
			expr: &choiceExpr{
				pos: position{line: 15, col: 19, offset: 327},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 15, col: 19, offset: 327},
						val:        " ",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 15, col: 25, offset: 333},
						val:        "\t",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 15, col: 32, offset: 340},
						name: "EOL",
					},
					&ruleRefExpr{
						pos:  position{line: 15, col: 38, offset: 346},
						name: "LineComment",
					},
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 17, col: 1, offset: 377},
			expr: &zeroOrMoreExpr{
				pos: position{line: 17, col: 14, offset: 392},
				expr: &ruleRefExpr{
					pos:  position{line: 17, col: 14, offset: 392},
					name: "WhitespaceChunk",
				},
			},
		},
		{
			name: "NonemptyWhitespace",
			pos:  position{line: 19, col: 1, offset: 410},
			expr: &oneOrMoreExpr{
				pos: position{line: 19, col: 22, offset: 433},
				expr: &ruleRefExpr{
					pos:  position{line: 19, col: 22, offset: 433},
					name: "WhitespaceChunk",
				},
			},
		},
		{
			name: "SimpleLabel",
			pos:  position{line: 21, col: 1, offset: 451},
			expr: &actionExpr{
				pos: position{line: 21, col: 15, offset: 467},
				run: (*parser).callonSimpleLabel1,
				expr: &seqExpr{
					pos: position{line: 21, col: 15, offset: 467},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 21, col: 15, offset: 467},
							val:        "[A-Za-z_]",
							chars:      []rune{'_'},
							ranges:     []rune{'A', 'Z', 'a', 'z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 21, col: 25, offset: 477},
							expr: &charClassMatcher{
								pos:        position{line: 21, col: 25, offset: 477},
								val:        "[A-Za-z0-9_/-]",
								chars:      []rune{'_', '/', '-'},
								ranges:     []rune{'A', 'Z', 'a', 'z', '0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "Label",
			pos:  position{line: 25, col: 1, offset: 541},
			expr: &actionExpr{
				pos: position{line: 25, col: 9, offset: 551},
				run: (*parser).callonLabel1,
				expr: &seqExpr{
					pos: position{line: 25, col: 9, offset: 551},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 25, col: 9, offset: 551},
							label: "label",
							expr: &ruleRefExpr{
								pos:  position{line: 25, col: 15, offset: 557},
								name: "SimpleLabel",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 25, col: 27, offset: 569},
							name: "Whitespace",
						},
					},
				},
			},
		},
		{
			name: "ReservedRaw",
			pos:  position{line: 27, col: 1, offset: 603},
			expr: &choiceExpr{
				pos: position{line: 27, col: 15, offset: 619},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 27, col: 15, offset: 619},
						val:        "Bool",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 28, col: 5, offset: 630},
						val:        "Optional",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 29, col: 5, offset: 645},
						val:        "None",
						ignoreCase: false,
					},
					&actionExpr{
						pos: position{line: 30, col: 5, offset: 656},
						run: (*parser).callonReservedRaw5,
						expr: &litMatcher{
							pos:        position{line: 30, col: 5, offset: 656},
							val:        "Natural",
							ignoreCase: false,
						},
					},
					&litMatcher{
						pos:        position{line: 31, col: 5, offset: 698},
						val:        "Integer",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 32, col: 5, offset: 712},
						val:        "Double",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 33, col: 5, offset: 725},
						val:        "Text",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 34, col: 5, offset: 736},
						val:        "List",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 35, col: 5, offset: 747},
						val:        "True",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 36, col: 5, offset: 758},
						val:        "False",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 37, col: 5, offset: 770},
						val:        "NaN",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 38, col: 5, offset: 780},
						val:        "Infinity",
						ignoreCase: false,
					},
					&actionExpr{
						pos: position{line: 39, col: 5, offset: 795},
						run: (*parser).callonReservedRaw15,
						expr: &litMatcher{
							pos:        position{line: 39, col: 5, offset: 795},
							val:        "Type",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 40, col: 5, offset: 831},
						run: (*parser).callonReservedRaw17,
						expr: &litMatcher{
							pos:        position{line: 40, col: 5, offset: 831},
							val:        "Kind",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 41, col: 5, offset: 867},
						run: (*parser).callonReservedRaw19,
						expr: &litMatcher{
							pos:        position{line: 41, col: 5, offset: 867},
							val:        "Sort",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Reserved",
			pos:  position{line: 43, col: 1, offset: 900},
			expr: &actionExpr{
				pos: position{line: 43, col: 12, offset: 913},
				run: (*parser).callonReserved1,
				expr: &seqExpr{
					pos: position{line: 43, col: 12, offset: 913},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 43, col: 12, offset: 913},
							label: "word",
							expr: &ruleRefExpr{
								pos:  position{line: 43, col: 17, offset: 918},
								name: "ReservedRaw",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 43, col: 29, offset: 930},
							name: "Whitespace",
						},
					},
				},
			},
		},
		{
			name: "OpenParens",
			pos:  position{line: 45, col: 1, offset: 963},
			expr: &seqExpr{
				pos: position{line: 45, col: 14, offset: 978},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 45, col: 14, offset: 978},
						val:        "(",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 45, col: 18, offset: 982},
						name: "Whitespace",
					},
				},
			},
		},
		{
			name: "CloseParens",
			pos:  position{line: 46, col: 1, offset: 993},
			expr: &seqExpr{
				pos: position{line: 46, col: 15, offset: 1009},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 46, col: 15, offset: 1009},
						val:        ")",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 46, col: 19, offset: 1013},
						name: "Whitespace",
					},
				},
			},
		},
		{
			name: "At",
			pos:  position{line: 47, col: 1, offset: 1024},
			expr: &seqExpr{
				pos: position{line: 47, col: 6, offset: 1031},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 47, col: 6, offset: 1031},
						val:        "@",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 47, col: 10, offset: 1035},
						name: "Whitespace",
					},
				},
			},
		},
		{
			name: "Colon",
			pos:  position{line: 48, col: 1, offset: 1046},
			expr: &seqExpr{
				pos: position{line: 48, col: 9, offset: 1056},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 48, col: 9, offset: 1056},
						val:        ":",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 48, col: 13, offset: 1060},
						name: "NonemptyWhitespace",
					},
				},
			},
		},
		{
			name: "Lambda",
			pos:  position{line: 50, col: 1, offset: 1080},
			expr: &seqExpr{
				pos: position{line: 50, col: 10, offset: 1091},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 50, col: 11, offset: 1092},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 50, col: 11, offset: 1092},
								val:        "\\",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 50, col: 18, offset: 1099},
								val:        "λ",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 50, col: 23, offset: 1105},
						name: "Whitespace",
					},
				},
			},
		},
		{
			name: "Forall",
			pos:  position{line: 51, col: 1, offset: 1116},
			expr: &seqExpr{
				pos: position{line: 51, col: 10, offset: 1127},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 51, col: 11, offset: 1128},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 51, col: 11, offset: 1128},
								val:        "forall",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 51, col: 22, offset: 1139},
								val:        "∀",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 51, col: 27, offset: 1146},
						name: "Whitespace",
					},
				},
			},
		},
		{
			name: "Arrow",
			pos:  position{line: 52, col: 1, offset: 1157},
			expr: &seqExpr{
				pos: position{line: 52, col: 9, offset: 1167},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 52, col: 10, offset: 1168},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 52, col: 10, offset: 1168},
								val:        "->",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 52, col: 17, offset: 1175},
								val:        "→",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 52, col: 22, offset: 1182},
						name: "Whitespace",
					},
				},
			},
		},
		{
			name: "NaturalLiteral",
			pos:  position{line: 56, col: 1, offset: 1226},
			expr: &actionExpr{
				pos: position{line: 56, col: 18, offset: 1245},
				run: (*parser).callonNaturalLiteral1,
				expr: &oneOrMoreExpr{
					pos: position{line: 56, col: 18, offset: 1245},
					expr: &charClassMatcher{
						pos:        position{line: 56, col: 18, offset: 1245},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "DeBruijn",
			pos:  position{line: 64, col: 1, offset: 1393},
			expr: &actionExpr{
				pos: position{line: 64, col: 12, offset: 1406},
				run: (*parser).callonDeBruijn1,
				expr: &seqExpr{
					pos: position{line: 64, col: 12, offset: 1406},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 64, col: 12, offset: 1406},
							val:        "@",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 64, col: 16, offset: 1410},
							label: "index",
							expr: &ruleRefExpr{
								pos:  position{line: 64, col: 22, offset: 1416},
								name: "NaturalLiteral",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 64, col: 37, offset: 1431},
							name: "Whitespace",
						},
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 66, col: 1, offset: 1471},
			expr: &actionExpr{
				pos: position{line: 66, col: 14, offset: 1486},
				run: (*parser).callonIdentifier1,
				expr: &seqExpr{
					pos: position{line: 66, col: 14, offset: 1486},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 66, col: 14, offset: 1486},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 66, col: 19, offset: 1491},
								name: "Label",
							},
						},
						&labeledExpr{
							pos:   position{line: 66, col: 25, offset: 1497},
							label: "index",
							expr: &zeroOrOneExpr{
								pos: position{line: 66, col: 31, offset: 1503},
								expr: &ruleRefExpr{
									pos:  position{line: 66, col: 31, offset: 1503},
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
			pos:  position{line: 74, col: 1, offset: 1674},
			expr: &actionExpr{
				pos: position{line: 74, col: 28, offset: 1703},
				run: (*parser).callonIdentifierReservedPrefix1,
				expr: &seqExpr{
					pos: position{line: 74, col: 28, offset: 1703},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 74, col: 28, offset: 1703},
							label: "reserved",
							expr: &ruleRefExpr{
								pos:  position{line: 74, col: 37, offset: 1712},
								name: "ReservedRaw",
							},
						},
						&labeledExpr{
							pos:   position{line: 74, col: 49, offset: 1724},
							label: "suffix",
							expr: &actionExpr{
								pos: position{line: 74, col: 57, offset: 1732},
								run: (*parser).callonIdentifierReservedPrefix6,
								expr: &oneOrMoreExpr{
									pos: position{line: 74, col: 57, offset: 1732},
									expr: &charClassMatcher{
										pos:        position{line: 74, col: 57, offset: 1732},
										val:        "[A-Za-z0-9/_-]",
										chars:      []rune{'/', '_', '-'},
										ranges:     []rune{'A', 'Z', 'a', 'z', '0', '9'},
										ignoreCase: false,
										inverted:   false,
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 74, col: 105, offset: 1780},
							name: "Whitespace",
						},
						&labeledExpr{
							pos:   position{line: 74, col: 116, offset: 1791},
							label: "index",
							expr: &zeroOrOneExpr{
								pos: position{line: 74, col: 122, offset: 1797},
								expr: &ruleRefExpr{
									pos:  position{line: 74, col: 122, offset: 1797},
									name: "DeBruijn",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Expression",
			pos:  position{line: 87, col: 1, offset: 2184},
			expr: &choiceExpr{
				pos: position{line: 88, col: 7, offset: 2205},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 88, col: 7, offset: 2205},
						run: (*parser).callonExpression2,
						expr: &seqExpr{
							pos: position{line: 88, col: 7, offset: 2205},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 88, col: 7, offset: 2205},
									name: "Lambda",
								},
								&ruleRefExpr{
									pos:  position{line: 88, col: 14, offset: 2212},
									name: "OpenParens",
								},
								&labeledExpr{
									pos:   position{line: 88, col: 25, offset: 2223},
									label: "label",
									expr: &ruleRefExpr{
										pos:  position{line: 88, col: 31, offset: 2229},
										name: "Label",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 88, col: 37, offset: 2235},
									name: "Colon",
								},
								&labeledExpr{
									pos:   position{line: 88, col: 43, offset: 2241},
									label: "t",
									expr: &ruleRefExpr{
										pos:  position{line: 88, col: 45, offset: 2243},
										name: "Expression",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 88, col: 56, offset: 2254},
									name: "CloseParens",
								},
								&ruleRefExpr{
									pos:  position{line: 88, col: 68, offset: 2266},
									name: "Arrow",
								},
								&labeledExpr{
									pos:   position{line: 88, col: 74, offset: 2272},
									label: "body",
									expr: &ruleRefExpr{
										pos:  position{line: 88, col: 79, offset: 2277},
										name: "Expression",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 91, col: 7, offset: 2406},
						run: (*parser).callonExpression15,
						expr: &seqExpr{
							pos: position{line: 91, col: 7, offset: 2406},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 91, col: 7, offset: 2406},
									name: "Forall",
								},
								&ruleRefExpr{
									pos:  position{line: 91, col: 14, offset: 2413},
									name: "OpenParens",
								},
								&labeledExpr{
									pos:   position{line: 91, col: 25, offset: 2424},
									label: "label",
									expr: &ruleRefExpr{
										pos:  position{line: 91, col: 31, offset: 2430},
										name: "Label",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 91, col: 37, offset: 2436},
									name: "Colon",
								},
								&labeledExpr{
									pos:   position{line: 91, col: 43, offset: 2442},
									label: "t",
									expr: &ruleRefExpr{
										pos:  position{line: 91, col: 45, offset: 2444},
										name: "Expression",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 91, col: 56, offset: 2455},
									name: "CloseParens",
								},
								&ruleRefExpr{
									pos:  position{line: 91, col: 68, offset: 2467},
									name: "Arrow",
								},
								&labeledExpr{
									pos:   position{line: 91, col: 74, offset: 2473},
									label: "body",
									expr: &ruleRefExpr{
										pos:  position{line: 91, col: 79, offset: 2478},
										name: "Expression",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 94, col: 7, offset: 2599},
						name: "AnnotatedExpression",
					},
				},
			},
		},
		{
			name: "AnnotatedExpression",
			pos:  position{line: 96, col: 1, offset: 2620},
			expr: &ruleRefExpr{
				pos:  position{line: 96, col: 23, offset: 2644},
				name: "OperatorExpression",
			},
		},
		{
			name: "OperatorExpression",
			pos:  position{line: 98, col: 1, offset: 2664},
			expr: &ruleRefExpr{
				pos:  position{line: 98, col: 22, offset: 2687},
				name: "ImportAltExpression",
			},
		},
		{
			name: "ImportAltExpression",
			pos:  position{line: 100, col: 1, offset: 2708},
			expr: &ruleRefExpr{
				pos:  position{line: 101, col: 7, offset: 2774},
				name: "ApplicationExpression",
			},
		},
		{
			name: "ApplicationExpression",
			pos:  position{line: 103, col: 1, offset: 2797},
			expr: &actionExpr{
				pos: position{line: 104, col: 7, offset: 2829},
				run: (*parser).callonApplicationExpression1,
				expr: &seqExpr{
					pos: position{line: 104, col: 7, offset: 2829},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 104, col: 7, offset: 2829},
							label: "e",
							expr: &ruleRefExpr{
								pos:  position{line: 104, col: 9, offset: 2831},
								name: "ImportExpression",
							},
						},
						&zeroOrMoreExpr{
							pos: position{line: 104, col: 26, offset: 2848},
							expr: &seqExpr{
								pos: position{line: 104, col: 27, offset: 2849},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 104, col: 27, offset: 2849},
										name: "WhitespaceChunk",
									},
									&ruleRefExpr{
										pos:  position{line: 104, col: 43, offset: 2865},
										name: "ImportExpression",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ImportExpression",
			pos:  position{line: 106, col: 1, offset: 2911},
			expr: &ruleRefExpr{
				pos:  position{line: 106, col: 20, offset: 2932},
				name: "SelectorExpression",
			},
		},
		{
			name: "SelectorExpression",
			pos:  position{line: 108, col: 1, offset: 2952},
			expr: &ruleRefExpr{
				pos:  position{line: 108, col: 22, offset: 2975},
				name: "PrimitiveExpression",
			},
		},
		{
			name: "PrimitiveExpression",
			pos:  position{line: 110, col: 1, offset: 2996},
			expr: &choiceExpr{
				pos: position{line: 111, col: 7, offset: 3026},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 111, col: 7, offset: 3026},
						name: "NaturalLiteral",
					},
					&ruleRefExpr{
						pos:  position{line: 112, col: 7, offset: 3047},
						name: "IdentifierReservedPrefix",
					},
					&ruleRefExpr{
						pos:  position{line: 113, col: 7, offset: 3078},
						name: "Reserved",
					},
					&ruleRefExpr{
						pos:  position{line: 114, col: 7, offset: 3093},
						name: "Identifier",
					},
					&actionExpr{
						pos: position{line: 115, col: 7, offset: 3110},
						run: (*parser).callonPrimitiveExpression6,
						expr: &seqExpr{
							pos: position{line: 115, col: 7, offset: 3110},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 115, col: 7, offset: 3110},
									name: "OpenParens",
								},
								&labeledExpr{
									pos:   position{line: 115, col: 18, offset: 3121},
									label: "e",
									expr: &ruleRefExpr{
										pos:  position{line: 115, col: 20, offset: 3123},
										name: "Expression",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 115, col: 31, offset: 3134},
									name: "CloseParens",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 117, col: 1, offset: 3165},
			expr: &notExpr{
				pos: position{line: 117, col: 7, offset: 3173},
				expr: &anyMatcher{
					line: 117, col: 8, offset: 3174,
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

func (c *current) onReservedRaw5() (interface{}, error) {
	return ast.Natural, nil
}

func (p *parser) callonReservedRaw5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw5()
}

func (c *current) onReservedRaw15() (interface{}, error) {
	return ast.Type, nil
}

func (p *parser) callonReservedRaw15() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw15()
}

func (c *current) onReservedRaw17() (interface{}, error) {
	return ast.Kind, nil
}

func (p *parser) callonReservedRaw17() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw17()
}

func (c *current) onReservedRaw19() (interface{}, error) {
	return ast.Sort, nil
}

func (p *parser) callonReservedRaw19() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReservedRaw19()
}

func (c *current) onReserved1(word interface{}) (interface{}, error) {
	return word, nil
}

func (p *parser) callonReserved1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReserved1(stack["word"])
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

func (c *current) onDeBruijn1(index interface{}) (interface{}, error) {
	return index.(int), nil
}

func (p *parser) callonDeBruijn1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDeBruijn1(stack["index"])
}

func (c *current) onIdentifier1(name, index interface{}) (interface{}, error) {
	if index != nil {
		return ast.Var{Name: name.(string), Index: index.(int)}, nil
	} else {
		return ast.Var{Name: name.(string)}, nil
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
		return ast.Var{Name: name, Index: index.(int)}, nil
	} else {
		return ast.Var{Name: name}, nil
	}
}

func (p *parser) callonIdentifierReservedPrefix1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifierReservedPrefix1(stack["reserved"], stack["suffix"], stack["index"])
}

func (c *current) onExpression2(label, t, body interface{}) (interface{}, error) {
	return &ast.LambdaExpr{Label: label.(string), Type: t.(ast.Expr), Body: body.(ast.Expr)}, nil

}

func (p *parser) callonExpression2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression2(stack["label"], stack["t"], stack["body"])
}

func (c *current) onExpression15(label, t, body interface{}) (interface{}, error) {
	return &ast.Pi{Label: label.(string), Type: t.(ast.Expr), Body: body.(ast.Expr)}, nil

}

func (p *parser) callonExpression15() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression15(stack["label"], stack["t"], stack["body"])
}

func (c *current) onApplicationExpression1(e interface{}) (interface{}, error) {
	return e, nil
}

func (p *parser) callonApplicationExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onApplicationExpression1(stack["e"])
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
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch = errors.New("no match found")
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
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
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
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
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
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug   bool
	depth   int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
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
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
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
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
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
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
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
