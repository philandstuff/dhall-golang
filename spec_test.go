// +build spec

package main_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/philandstuff/dhall-golang/parser"
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
	"TestParserAccepts/functionTypeA.dhall", // requires arrow-expression (ie Bool -> Bool instead of forall (_ : Bool) -> Bool)
	"TestParserAccepts/ifThenElseA.dhall",
	"TestParserAccepts/importAltA.dhall",
	"TestParserAccepts/interpolatedDoubleQuotedStringA.dhall",
	"TestParserAccepts/interpolatedSingleQuotedStringA.dhall",
	"TestParserAccepts/labelA.dhall", // requires let
	"TestParserAccepts/largeExpressionA.dhall",
	"TestParserAccepts/letA.dhall",
	"TestParserAccepts/mergeA.dhall",
	"TestParserAccepts/multiletA.dhall",
	"TestParserAccepts/operatorsA.dhall",
	"TestParserAccepts/parenthesizeUsingA.dhall",
	"TestParserAccepts/pathTerminationA.dhall",
	"TestParserAccepts/pathsA.dhall",
	"TestParserAccepts/quotedLabelA.dhall",
	"TestParserAccepts/quotedPathsA.dhall",
	"TestParserAccepts/recordA.dhall",
	"TestParserAccepts/reservedPrefixA.dhall", // requires let
	"TestParserAccepts/singleQuotedStringA.dhall",
	"TestParserAccepts/templateA.dhall",
	"TestParserAccepts/unicodeDoubleQuotedStringA.dhall",
	"TestParserAccepts/unionA.dhall",
	"TestParserAccepts/urlsA.dhall",
}

func pass(t *testing.T) {
	for _, name := range expectedFailures {
		if t.Name() == name {
			t.Error("Expected failure, but actually passed")
		}
	}
}

func failf(t *testing.T, format string, args ...interface{}) {
	for _, name := range expectedFailures {
		if t.Name() == name {
			t.Skipf(format, args...)
			return
		}
	}
	t.Errorf(format, args...)
}

func expectError(t *testing.T, err error) {
	if err == nil {
		failf(t, "Expected file to fail to parse, but it parsed successfully")
	} else {
		pass(t)
	}
}

func expectNoError(t *testing.T, err error) {
	if err != nil {
		failf(t, "Expected file to parse successfully, but got error %v", err)
	} else {
		pass(t)
	}
}

func TestParserRejects(t *testing.T) {
	failuresDir := "./dhall-lang/tests/parser/failure/"
	files, err := ioutil.ReadDir(failuresDir)
	if err != nil {
		t.Fatalf("Couldn't read dhall-lang tests: %v\n(Have you pulled submodules?)\n", err)
	}

	for _, f := range files {
		t.Run(f.Name(), func(t *testing.T) {
			reader, openerr := os.Open(failuresDir + f.Name())
			defer reader.Close()
			if openerr != nil {
				t.Fatal(openerr)
			}

			_, err := parser.ParseReader(f.Name(), reader)

			expectError(t, err)
		})
	}
}

func TestParserAccepts(t *testing.T) {
	successesDir := "./dhall-lang/tests/parser/success/"
	files, err := filepath.Glob(successesDir + "*A.dhall")
	if err != nil {
		t.Fatalf("Couldn't read dhall-lang tests: %v\n(Have you pulled submodules?)\n", err)
	}

	for _, f := range files {
		reader, openerr := os.Open(f)
		name := filepath.Base(f)
		defer reader.Close()
		if openerr != nil {
			t.Fatal(openerr)
		}
		t.Run(name, func(t *testing.T) {

			_, err := parser.ParseReader(name, reader)

			expectNoError(t, err)
		})
	}
}
