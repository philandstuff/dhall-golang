package main_test

import (
	"bytes"
	"io"
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
	"TestParserAccepts/recordA.dhall",
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
	// FIXME: accessEncodedTypeA parses, so why doesn't it typecheck?
	"TestTypechecks/accessEncodedTypeA.dhall",
	"TestTypechecks/accessTypeA.dhall",
	"TestTypechecks/prelude",
	"TestTypechecks/recordOfRecordOfTypesA.dhall",
	"TestTypechecks/recordOfTypesA.dhall",
	"TestTypechecks/simple/access",
	"TestTypechecks/simple/alternativesAreTypesA.dhall",
	// breaks because no alpha normalization yet
	"TestTypechecks/simple/anonymousFunctionsInTypesA.dhall",
	"TestTypechecks/simple/fieldsAreTypesA.dhall",
	// breaks because no alpha normalization yet
	"TestTypechecks/simple/kindParameterA.dhall",
	"TestTypechecks/simple/mergeEquivalenceA.dhall",
	"TestTypechecks/simple/mixedFieldAccessA.dhall",
	"TestTypechecks/simple/unionsOfTypesA.dhall",
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
	if reflect.DeepEqual(expected, actual) {
		pass(t)
	} else {
		failf(t, "Expected %+v to equal %+v", actual, expected)
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
	var cbor codec.CborHandle
	var json codec.JsonHandle
	runTestOnFilePairs(t, "dhall-lang/tests/parser/success/",
		"A.dhall", "B.json",
		func(t *testing.T, aReader, bReader io.Reader) {
			buf := new(bytes.Buffer)
			parsed, err := parser.ParseReader(t.Name(), aReader)
			expectNoError(t, err)
			aEnc := codec.NewEncoder(buf, &cbor)
			err = aEnc.Encode(parsed)
			expectNoError(t, err)
			aDec := codec.NewDecoder(buf, &cbor)
			var actual interface{}
			err = aDec.Decode(&actual)
			expectNoError(t, err)

			bDec := codec.NewDecoder(bReader, &json)
			var expected interface{}
			err = bDec.Decode(&expected)
			expectNoError(t, err)
			expectEqual(t, expected, actual)
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

			typeOfA, err := parsedA.(ast.Expr).TypeWith(ast.EmptyContext())
			expectNoError(t, err)
			expectEqual(t, typeOfA, parsedB)
		})
}
