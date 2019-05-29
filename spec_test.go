package main_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/philandstuff/dhall-golang/ast"
	"github.com/philandstuff/dhall-golang/imports"
	"github.com/philandstuff/dhall-golang/parser"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

var expectedFailures = []string{
	"TestParserAccepts/annotationsA.dhall", // requires old optional syntax
	// FIXME binary encoding doesn't match here
	"TestParserAccepts/doubleA.dhall",
	"TestParserAccepts/unionA.dhall", // deprecated syntax
	"TestParserAccepts/unit/import/hash",
	"TestParserAccepts/unit/import/parenthesizeUsing",
	"TestParserAccepts/unit/import/pathTerminationUnion", // needs union literals, which I won't implement
	"TestParserAccepts/unit/import/quotedPaths",          // needs.. quoted paths
	"TestParserAccepts/unit/import/unicodePaths",         // needs quoted paths
	"TestParserAccepts/unit/import/urlsA",                // needs all the URL features
	"TestParserAccepts/unit/import/urls/potPourri",       // needs all the URL features
	"TestParserAccepts/unit/import/urls/quotedPath",      // needs quotedPaths
	"TestParserAccepts/unit/UnionLit",                    // not going to implement union literals
	"TestTypecheckFails/duplicateFields.dhall",           // in dhall-golang, duplicate fields a parse error, not a type error
	"TestTypecheckFails/unit/README",                     // FIXME, shouldn't need excluding
	"TestTypecheckFails/unit/RecordProjection",
	"TestTypeInference/simple/alternativesAreTypesA.dhall", // old union literals
	"TestTypeInference/unit/OldOptional",
	"TestTypeInference/unit/RecordProjection",
	"TestTypeInference/unit/UnionLiteral", // deprecated union literal syntax
	"TestNormalization/haskell-tutorial/projection",
	"TestNormalization/prelude/Double",
	"TestNormalization/simple/doubleShowA.dhall",      // getting precision right is tricky
	"TestNormalization/simple/integerToDoubleA.dhall", // requires bigint representation, which the standard itself does not require
	"TestNormalization/unit/DoubleShowValue",          // getting precision right is tricky
	"TestNormalization/unit/NoneNaturalA",             // I don't intend to implement this; it will disappear from the standard
	"TestNormalization/unit/RecordProjection",
	"TestNormalization/unit/UnionLiteral",
	"TestImport/alternativeHashMismatchA.dhall", // needs alternative operator
	"TestImport/fieldOrder",                     // needs import hashes
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
	test func(*testing.T, io.Reader),
) {
	err := filepath.Walk(dir,
		func(testPath string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if err != nil {
				return err
			}
			reader, err := os.Open(testPath)
			if err != nil {
				t.Fatal(err)
			}
			defer reader.Close()
			name := strings.Replace(testPath, dir, "", 1)
			wd, _ := os.Getwd()
			os.Chdir(path.Dir(testPath))
			t.Run(name, func(t *testing.T) {
				test(t, reader)
				pass(t)
			})
			os.Chdir(wd)
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

	wd, _ := os.Getwd()
	os.Chdir(path.Dir(pathA))
	t.Run(name, func(t *testing.T) {
		test(t, aReader, bReader)
		pass(t)
	})
	os.Chdir(wd)
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
			expectEqualCbor(t, expected, actualBuf.Bytes())
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

			resolvedA, err := imports.Load(parsedA.(ast.Expr))
			expectNoError(t, err)

			resolvedB, err := imports.Load(parsedB.(ast.Expr))
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
		func(t *testing.T, aReader, bReader io.Reader) {
			parsedA, err := parser.ParseReader(t.Name(), aReader)
			expectNoError(t, err)

			parsedB, err := parser.ParseReader(t.Name(), bReader)
			expectNoError(t, err)

			inferredType, err := parsedA.(ast.Expr).TypeWith(ast.EmptyContext())
			expectNoError(t, err)

			expectEqualExprs(t, parsedB.(ast.Expr), inferredType)
		})
}

func TestAlphaNormalization(t *testing.T) {
	runTestOnFilePairs(t, "dhall-lang/tests/alpha-normalization/success/",
		"A.dhall", "B.dhall",
		func(t *testing.T, aReader, bReader io.Reader) {
			parsedA, err := parser.ParseReader(t.Name(), aReader)
			expectNoError(t, err)

			parsedB, err := parser.ParseReader(t.Name(), bReader)
			expectNoError(t, err)

			normA := parsedA.(ast.Expr).AlphaNormalize()
			normB := parsedB.(ast.Expr).AlphaNormalize()

			expectEqualExprs(t, normB, normA)
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

			resolvedA, err := imports.Load(parsedA.(ast.Expr))
			expectNoError(t, err)

			resolvedB, err := imports.Load(parsedB.(ast.Expr))
			expectNoError(t, err)

			normA := resolvedA.(ast.Expr).Normalize()
			normB := resolvedB.(ast.Expr).Normalize()

			expectEqualExprs(t, normB, normA)
		})
}

func TestImportFails(t *testing.T) {
	runTestOnEachFile(t, "dhall-lang/tests/import/failure/", func(t *testing.T, reader io.Reader) {
		parsed, err := parser.ParseReader(t.Name(), reader)
		expectNoError(t, err)

		_, err = imports.Load(parsed.(ast.Expr))
		expectError(t, err)
	})
}

func TestImport(t *testing.T) {
	runTestOnFilePairs(t, "dhall-lang/tests/import/success/",
		"A.dhall", "B.dhall",
		func(t *testing.T, aReader, bReader io.Reader) {
			parsedA, err := parser.ParseReader(t.Name(), aReader)
			expectNoError(t, err)

			parsedB, err := parser.ParseReader(t.Name(), bReader)
			expectNoError(t, err)

			resolvedA, err := imports.Load(parsedA.(ast.Expr))
			expectNoError(t, err)

			resolvedB, err := imports.Load(parsedB.(ast.Expr))
			expectNoError(t, err)

			expectEqualExprs(t, resolvedB, resolvedA)
		})
}
