package main_test

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/philandstuff/dhall-golang/ast"
	"github.com/philandstuff/dhall-golang/imports"
	"github.com/philandstuff/dhall-golang/parser"
	"github.com/pkg/errors"
)

var expectedFailures = []string{
	"TestParserAccepts/unionA.dhall", // deprecated syntax
	"TestParserAccepts/unit/import/inlineUsing",
	"TestParserAccepts/unit/import/parenthesizeUsing",
	"TestParserAccepts/unit/import/pathTerminationUnion", // needs union literals, which I won't implement
	"TestParserAccepts/unit/import/quotedPaths",          // needs.. quoted paths
	"TestParserAccepts/unit/import/unicodePaths",         // needs quoted paths
	"TestParserAccepts/unit/import/urls/potPourri",       // net/url doesn't parse authorities in the way the test expects
	"TestParserAccepts/unit/import/urls/quotedPath",      // needs quotedPaths
	"TestParserAccepts/unit/UnionLit",                    // not going to implement union literals
	"TestTypecheckFails/duplicateFields.dhall",           // in dhall-golang, duplicate fields a parse error, not a type error
	"TestTypecheckFails/unit/README",                     // FIXME, shouldn't need excluding
	"TestTypecheckFails/customHeadersUsingBoundVariable",
	"TestTypeInference/simple/alternativesAreTypesA.dhall", // old union literals
	"TestTypeInference/unit/UnionLiteral",                  // deprecated union literal syntax
	"TestNormalization/prelude/JSON/Type",                  // test bug, fixed in dhall-lang/dhall-lang#599
	"TestNormalization/simple/integerToDoubleA.dhall",      // requires bigint representation, which the standard itself does not require
	"TestNormalization/unit/UnionLiteral",
	"TestImport/customHeadersA.dhall",
	"TestImport/headerForwardingA.dhall",
	"TestImport/noHeaderForwardingA.dhall",
	"TestSemanticHash/simple/integerToDouble", // oversized int (we don't support bigints)
}

func pass(t *testing.T) {
	t.Helper()
	for _, prefix := range expectedFailures {
		if strings.HasPrefix(t.Name(), prefix) {
			t.Fatal("Expected failure, but actually passed")
		}
	}
}

func failf(t *testing.T, format string, args ...interface{}) {
	t.Helper()
	for _, prefix := range expectedFailures {
		if strings.HasPrefix(t.Name(), prefix) {
			t.Skipf(format, args...)
			return
		}
	}
	t.Fatalf(format, args...)
}

func expectError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		failf(t, "Expected error, but saw none")
	}
}

func expectNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		failf(t, "Got unexpected error %v", err)
	}
}

func expectEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		failf(t, "Expected %+v to equal %+v", actual, expected)
	}
}

func expectEqualBytes(t *testing.T, expected, actual []byte) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		failf(t, "Expected %x to equal %x", actual, expected)
	}
}

func expectEqualExprs(t *testing.T, expected, actual ast.Expr) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		failf(t, "Expected %v to equal %v", actual, expected)
	}
}

func cbor2Pretty(cbor []byte) (string, error) {
	// cbor2pretty.rb comes from the cbor-diag gem
	cmd := exec.Command("cbor2pretty.rb")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	go func() {
		defer stdin.Close()
		stdin.Write(cbor)
	}()
	out, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return "", errors.Wrap(err, "error calling cbor2pretty.rb. stderr was: "+string(exitError.Stderr))
		}
		return "", errors.Wrap(err, "error calling cbor2pretty.rb")
	}
	return string(out), nil
}

func expectEqualCbor(t *testing.T, expected, actual []byte) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		actualPretty, err := cbor2Pretty(actual)
		if err != nil {
			failf(t, "Couldn't decode actual CBOR: %v", err)
		}
		expectedPretty, err := cbor2Pretty(expected)
		if err != nil {
			failf(t, "Couldn't decode expected CBOR: %v", err)
		}
		failf(t, "Expected\n%s\nto equal\n%s\n", actualPretty, expectedPretty)
	}
}

func runTestOnEachFile(
	t *testing.T,
	dir string,
	test func(*testing.T, string),
) {
	err := filepath.Walk(dir,
		func(testPath string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if err != nil {
				return err
			}
			name := strings.Replace(testPath, dir, "", 1)
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				test(t, testPath)
				pass(t)
			})
			return nil
		})
	if err != nil {
		t.Fatalf("Couldn't read dhall-lang tests: %v\n(Have you pulled submodules?)\n", err)
	}
}

func runTestOnFilePairs(
	t *testing.T,
	dir, suffixA, suffixB string,
	test func(*testing.T, string, string),
) {
	err := filepath.Walk(dir,
		func(aPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(aPath, suffixA) {
				bPath := strings.Replace(aPath, suffixA, suffixB, 1)
				testName := strings.Replace(aPath, dir, "", 1)

				t.Run(testName, func(t *testing.T) {
					t.Parallel()
					test(t, aPath, bPath)
					pass(t)
				})
			}
			return nil
		})
	if err != nil {
		t.Fatalf("Couldn't read dhall-lang tests: %v\n(Have you pulled submodules?)\n", err)
	}
}

func TestParserRejects(t *testing.T) {
	runTestOnEachFile(t, "dhall-lang/tests/parser/failure/", func(t *testing.T, testPath string) {
		_, err := parser.ParseFile(testPath)

		expectError(t, err)
	})
}

func TestParserAccepts(t *testing.T) {
	runTestOnFilePairs(t, "dhall-lang/tests/parser/success/",
		"A.dhall", "B.dhallb",
		func(t *testing.T, aPath, bPath string) {
			actualBuf := new(bytes.Buffer)
			parsed, err := parser.ParseFile(aPath)
			expectNoError(t, err)

			err = ast.EncodeAsCbor(actualBuf, parsed.(ast.Expr))
			expectNoError(t, err)

			expected, err := ioutil.ReadFile(bPath)
			expectNoError(t, err)
			expectEqualCbor(t, expected, actualBuf.Bytes())
		})
}

func TestTypecheckFails(t *testing.T) {
	runTestOnEachFile(t, "dhall-lang/tests/typecheck/failure/", func(t *testing.T, testPath string) {
		parsed, err := parser.ParseFile(testPath)

		expectNoError(t, err)

		expr, ok := parsed.(ast.Expr)
		if !ok {
			failf(t, "Expected ast.Expr, got %+v\n", parsed)
		}

		_, err = expr.TypeWith(ast.EmptyContext())

		expectError(t, err)
	})
}

func TestTypechecks(t *testing.T) {
	runTestOnFilePairs(t, "dhall-lang/tests/typecheck/success/",
		"A.dhall", "B.dhall",
		func(t *testing.T, aPath, bPath string) {
			parsedA, err := parser.ParseFile(aPath)
			expectNoError(t, err)

			parsedB, err := parser.ParseFile(bPath)
			expectNoError(t, err)

			resolvedA, err := imports.Load(parsedA.(ast.Expr), ast.Local(aPath))
			expectNoError(t, err)

			resolvedB, err := imports.Load(parsedB.(ast.Expr), ast.Local(bPath))
			expectNoError(t, err)

			annot := ast.Annot{
				Expr:       resolvedA.(ast.Expr),
				Annotation: resolvedB.(ast.Expr),
			}
			_, err = annot.TypeWith(ast.EmptyContext())
			expectNoError(t, err)
		})
}

func TestTypeInference(t *testing.T) {
	runTestOnFilePairs(t, "dhall-lang/tests/type-inference/success/",
		"A.dhall", "B.dhall",
		func(t *testing.T, aPath, bPath string) {
			parsedA, err := parser.ParseFile(aPath)
			expectNoError(t, err)

			parsedB, err := parser.ParseFile(bPath)
			expectNoError(t, err)

			inferredType, err := parsedA.(ast.Expr).TypeWith(ast.EmptyContext())
			expectNoError(t, err)

			expectEqualExprs(t, parsedB.(ast.Expr), inferredType)
		})
}

func TestAlphaNormalization(t *testing.T) {
	runTestOnFilePairs(t, "dhall-lang/tests/alpha-normalization/success/",
		"A.dhall", "B.dhall",
		func(t *testing.T, aPath, bPath string) {
			parsedA, err := parser.ParseFile(aPath)
			expectNoError(t, err)

			parsedB, err := parser.ParseFile(bPath)
			expectNoError(t, err)

			normA := parsedA.(ast.Expr).AlphaNormalize()
			normB := parsedB.(ast.Expr).AlphaNormalize()

			expectEqualExprs(t, normB, normA)
		})
}

func TestNormalization(t *testing.T) {
	runTestOnFilePairs(t, "dhall-lang/tests/normalization/success/",
		"A.dhall", "B.dhall",
		func(t *testing.T, aPath, bPath string) {
			parsedA, err := parser.ParseFile(aPath)
			expectNoError(t, err)

			parsedB, err := parser.ParseFile(bPath)
			expectNoError(t, err)

			resolvedA, err := imports.Load(parsedA.(ast.Expr), ast.Local(aPath))
			expectNoError(t, err)

			resolvedB, err := imports.Load(parsedB.(ast.Expr), ast.Local(bPath))
			expectNoError(t, err)

			normA := resolvedA.(ast.Expr).Normalize()
			normB := resolvedB.(ast.Expr).Normalize()

			expectEqualExprs(t, normB, normA)
		})
}

func TestImportFails(t *testing.T) {
	runTestOnEachFile(t, "dhall-lang/tests/import/failure/", func(t *testing.T, testPath string) {
		parsed, err := parser.ParseFile(testPath)
		expectNoError(t, err)

		_, err = imports.Load(parsed.(ast.Expr), ast.Local(testPath))
		expectError(t, err)
	})
}

func TestImport(t *testing.T) {
	cwd, err := os.Getwd()
	expectNoError(t, err)
	os.Setenv("XDG_CACHE_HOME", cwd+"/dhall-lang/tests/import/cache")
	runTestOnFilePairs(t, "dhall-lang/tests/import/success/",
		"A.dhall", "B.dhall",
		func(t *testing.T, aPath, bPath string) {
			parsedA, err := parser.ParseFile(aPath)
			expectNoError(t, err)

			parsedB, err := parser.ParseFile(bPath)
			expectNoError(t, err)

			resolvedA, err := imports.Load(parsedA.(ast.Expr), ast.Local(aPath))
			expectNoError(t, err)

			resolvedB, err := imports.Load(parsedB.(ast.Expr), ast.Local(bPath))
			expectNoError(t, err)

			expectEqualExprs(t, resolvedB, resolvedA)
		})
}

func TestSemanticHash(t *testing.T) {
	sha256re := regexp.MustCompile("^sha256:([0-9a-fA-F]{64})\n$")
	runTestOnFilePairs(t, "dhall-lang/tests/semantic-hash/success/",
		"A.dhall", "B.hash",
		func(t *testing.T, aPath, bPath string) {
			parsedA, err := parser.ParseFile(aPath)
			expectNoError(t, err)
			resolvedA, err := imports.Load(parsedA.(ast.Expr), ast.Local(aPath))
			expectNoError(t, err)

			actualHash, err := ast.SemanticHash(resolvedA)
			expectNoError(t, err)

			expectedHashStr, err := ioutil.ReadFile(bPath)
			expectNoError(t, err)

			groups := sha256re.FindSubmatch(expectedHashStr)
			if groups == nil {
				t.Fatalf("Couldn't parse expected hash string %s", expectedHashStr)
			}

			expectedRawHash := make([]byte, 32)
			_, err = hex.Decode(expectedRawHash, groups[1])
			expectNoError(t, err)

			expectEqualBytes(t, append([]byte{0x12, 0x20}, expectedRawHash...), actualHash)
		})
}

func TestBinaryDecode(t *testing.T) {
	runTestOnFilePairs(t, "dhall-lang/tests/binary-decode/success/",
		"A.dhallb", "B.dhall",
		func(t *testing.T, aPath, bPath string) {
			aReader, err := os.Open(aPath)
			expectNoError(t, err)
			defer aReader.Close()
			expr, err := ast.DecodeAsCbor(aReader)
			expectNoError(t, err)

			parsedB, err := parser.ParseFile(bPath)
			expectNoError(t, err)

			expectEqualExprs(t, parsedB.(ast.Expr), expr)
		})
}

func TestBinaryDecodeFails(t *testing.T) {
	runTestOnEachFile(t, "dhall-lang/tests/binary-decode/failure/", func(t *testing.T, testPath string) {
		reader, err := os.Open(testPath)
		if err != nil {
			t.Fatal(err)
		}
		defer reader.Close()
		_, err = ast.DecodeAsCbor(reader)
		expectError(t, err)
	})
}
