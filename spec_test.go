package main_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/philandstuff/dhall-golang/ast"
	"github.com/philandstuff/dhall-golang/parser"
	"github.com/ugorji/go/codec"
)

var expectedFailures = []string{
	"TestParserAccepts/annotationsA.dhall", // requires records, list append, optionals
	"TestParserAccepts/asTextA.dhall",
	"TestParserAccepts/builtinsA.dhall",
	"TestParserAccepts/collectionImportTypeA.dhall",
	"TestParserAccepts/constructorsA.dhall",
	// FIXME binary encoding doesn't match here
	"TestParserAccepts/doubleA.dhall",
	"TestParserAccepts/doubleQuotedStringA.dhall",
	"TestParserAccepts/environmentVariablesA.dhall",
	"TestParserAccepts/escapedDoubleQuotedStringA.dhall",
	"TestParserAccepts/escapedSingleQuotedStringA.dhall",
	"TestParserAccepts/fieldsA.dhall",
	"TestParserAccepts/importAltA.dhall",
	"TestParserAccepts/interpolatedDoubleQuotedStringA.dhall",
	"TestParserAccepts/interpolatedSingleQuotedStringA.dhall",
	"TestParserAccepts/largeExpressionA.dhall",
	"TestParserAccepts/mergeA.dhall",
	"TestParserAccepts/operatorsA.dhall",
	"TestParserAccepts/parenthesizeUsingA.dhall",
	"TestParserAccepts/pathTerminationA.dhall",
	"TestParserAccepts/pathsA.dhall",
	"TestParserAccepts/quotedLabelA.dhall",
	"TestParserAccepts/quotedPathsA.dhall",
	"TestParserAccepts/singleQuotedStringA.dhall",
	"TestParserAccepts/templateA.dhall",
	"TestParserAccepts/unicodePathsA.dhall",
	"TestParserAccepts/unicodeDoubleQuotedStringA.dhall",
	"TestParserAccepts/unionA.dhall",
	"TestParserAccepts/urlsA.dhall",
	"TestTypecheckFails/combineMixedRecords.dhall",
	"TestTypecheckFails/duplicateFields.dhall",
	"TestTypecheckFails/mixedUnions.dhall",
	"TestTypecheckFails/preferMixedRecords.dhall",
	"TestTypechecks/prelude",
	"TestTypechecks/recordOfRecordOfTypesA.dhall",
	"TestTypechecks/simple/access/1A.dhall",
	"TestTypechecks/simple/alternativesAreTypesA.dhall",
	"TestTypechecks/simple/mergeEquivalenceA.dhall",
	"TestTypechecks/simple/mixedFieldAccessA.dhall",
	"TestTypechecks/simple/unionsOfTypesA.dhall",
	"TestNormalization/haskell-tutorial/access/1A.dhall",
	"TestNormalization/haskell-tutorial/combineTypes",
	"TestNormalization/haskell-tutorial/prefer",
	"TestNormalization/haskell-tutorial/projection",
	"TestNormalization/multiline",
	"TestNormalization/prelude",
	"TestNormalization/remoteSystemsA.dhall",
	"TestNormalization/simple/doubleShowA.dhall",
	"TestNormalization/simple/integerShowA.dhall",
	"TestNormalization/simple/integerToDoubleA.dhall",
	"TestNormalization/simple/listBuildA.dhall",
	"TestNormalization/simple/multiLineA.dhall",
	"TestNormalization/simple/naturalBuildA.dhall",
	"TestNormalization/simple/naturalShowA.dhall",
	"TestNormalization/simple/naturalToIntegerA.dhall",
	"TestNormalization/simple/optional",
	"TestNormalization/simple/sortOperatorA.dhall",
	"TestNormalization/simplifications/andA.dhall",
	"TestNormalization/simplifications/eqA.dhall",
	"TestNormalization/simplifications/neA.dhall",
	"TestNormalization/simplifications/orA.dhall",
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
		failf(t, "Expected file to fail to parse, but it parsed successfully")
	}
}

func expectNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		failf(t, "Expected file to parse successfully, but got error %v", err)
	}
}

func expectEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		failf(t, "Expected %+v to equal %+v", actual, expected)
	}
}

func expectEqualExprs(t *testing.T, expected, actual ast.Expr) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		failf(t, "Expected %v to equal %v", actual, expected)
	}
}

func expectEqualBytes(t *testing.T, expected, actual []byte) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		failf(t, "Expected %x to equal %x", actual, expected)
	}
}

func runTestOnEachFile(
	t *testing.T,
	dir string,
	test func(*testing.T, io.Reader),
) {
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if err != nil {
				return err
			}
			reader, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			defer reader.Close()
			name := strings.Replace(path, dir, "", 1)
			t.Run(name, func(t *testing.T) {
				test(t, reader)
				pass(t)
			})
			return nil
		})
	if err != nil {
		t.Fatalf("Couldn't read dhall-lang tests: %v\n(Have you pulled submodules?)\n", err)
	}
}

func runTestOnFilePair(t *testing.T, name, pathA, pathB string, test func(*testing.T, io.Reader, io.Reader)) {
	aReader, err := os.Open(pathA)
	defer aReader.Close()
	if err != nil {
		t.Fatal(err)
	}
	bReader, err := os.Open(pathB)
	defer bReader.Close()
	if err != nil {
		t.Fatal(err)
	}
	t.Run(name, func(t *testing.T) {
		test(t, aReader, bReader)
		pass(t)
	})
}

func runTestOnFilePairs(
	t *testing.T,
	dir, suffixA, suffixB string,
	test func(*testing.T, io.Reader, io.Reader),
) {
	err := filepath.Walk(dir,
		func(aPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(aPath, suffixA) {
				bPath := strings.Replace(aPath, suffixA, suffixB, 1)
				testName := strings.Replace(aPath, dir, "", 1)

				runTestOnFilePair(t, testName, aPath, bPath, test)
			}
			return nil
		})
	if err != nil {
		t.Fatalf("Couldn't read dhall-lang tests: %v\n(Have you pulled submodules?)\n", err)
	}
}

func TestParserRejects(t *testing.T) {
	runTestOnEachFile(t, "dhall-lang/tests/parser/failure/", func(t *testing.T, reader io.Reader) {
		_, err := parser.ParseReader(t.Name(), reader)

		expectError(t, err)
	})
}

func TestParserAccepts(t *testing.T) {
	cbor := ast.NewCborHandle()
	runTestOnFilePairs(t, "dhall-lang/tests/parser/success/",
		"A.dhall", "B.dhallb",
		func(t *testing.T, aReader, bReader io.Reader) {
			actualBuf := new(bytes.Buffer)
			parsed, err := parser.ParseReader(t.Name(), aReader)
			expectNoError(t, err)

			aEnc := codec.NewEncoder(actualBuf, &cbor)
			err = aEnc.Encode(parsed)
			expectNoError(t, err)

			expected, err := ioutil.ReadAll(bReader)
			expectNoError(t, err)
			expectEqualBytes(t, expected, actualBuf.Bytes())
		})
}

func TestTypecheckFails(t *testing.T) {
	runTestOnEachFile(t, "dhall-lang/tests/typecheck/failure/", func(t *testing.T, reader io.Reader) {
		parsed, err := parser.ParseReader(t.Name(), reader)

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
		func(t *testing.T, aReader, bReader io.Reader) {
			parsedA, err := parser.ParseReader(t.Name(), aReader)
			expectNoError(t, err)

			parsedB, err := parser.ParseReader(t.Name(), bReader)
			expectNoError(t, err)

			annot := ast.Annot{
				Expr:       parsedA.(ast.Expr),
				Annotation: parsedB.(ast.Expr),
			}
			_, err = annot.TypeWith(ast.EmptyContext())
			expectNoError(t, err)
		})
}

func TestNormalization(t *testing.T) {
	runTestOnFilePairs(t, "dhall-lang/tests/normalization/success/",
		"A.dhall", "B.dhall",
		func(t *testing.T, aReader, bReader io.Reader) {
			parsedA, err := parser.ParseReader(t.Name(), aReader)
			expectNoError(t, err)

			parsedB, err := parser.ParseReader(t.Name(), bReader)
			expectNoError(t, err)

			_, err = parsedA.(ast.Expr).TypeWith(ast.EmptyContext())
			expectNoError(t, err)

			normA := parsedA.(ast.Expr).Normalize()
			normB := parsedB.(ast.Expr).Normalize()

			expectEqualExprs(t, normB, normA)
		})
}
