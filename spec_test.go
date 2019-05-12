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
	"TestParserAccepts/annotationsA.dhall", // requires records, list append, optionals
	"TestParserAccepts/builtinsA.dhall",    // Haven't implemented all builtins
	"TestParserAccepts/collectionImportTypeA.dhall",
	"TestParserAccepts/constructorsA.dhall",
	// FIXME binary encoding doesn't match here
	"TestParserAccepts/doubleA.dhall",
	"TestParserAccepts/largeExpressionA.dhall",
	"TestParserAccepts/merge", // requires Natural/even
	"TestParserAccepts/operatorsA.dhall",
	"TestParserAccepts/quotedBoundVariableA.dhall",
	"TestParserAccepts/quotedLabelA.dhall",
	"TestParserAccepts/text/interpolatedDoubleQuotedStringA.dhall", // needs Natural/show
	"TestParserAccepts/text/interpolatedSingleQuotedStringA.dhall",
	"TestParserAccepts/text/interpolationA.dhall",
	"TestParserAccepts/text/templateA.dhall",
	"TestParserAccepts/unionA.dhall",
	"TestParserAccepts/unit/import/hash",
	"TestParserAccepts/unit/import/importAlt",
	"TestParserAccepts/unit/import/parenthesizeUsing",
	"TestParserAccepts/unit/import/pathTerminationUnion", // needs union literals, which I won't implement
	"TestParserAccepts/unit/import/quotedPaths",          // needs.. quoted paths
	"TestParserAccepts/unit/import/unicodePaths",         // needs quoted paths
	"TestParserAccepts/unit/import/urlsA",                // needs all the URL features
	"TestParserAccepts/unit/import/urls/potPourri",       // needs all the URL features
	"TestParserAccepts/unit/import/urls/quotedPath",      // needs quotedPaths
	"TestParserAccepts/unit/Quoted",
	"TestParserAccepts/unit/UnionLit",           // not going to implement union literals
	"TestParserAccepts/whitespaceBuffetA.dhall", // requires Natural/even
	"TestTypecheckFails/combineMixedRecords.dhall",
	"TestTypecheckFails/duplicateFields.dhall",
	"TestTypecheckFails/preferMixedRecords.dhall",
	"TestTypecheckFails/unit/MergeAlternative",
	"TestTypecheckFails/unit/MergeHandlerNotFunction",
	"TestTypecheckFails/unit/MergeHandlerNotMatchAlternativeType",
	"TestTypecheckFails/unit/MergeHandlersWithDifferentType",
	"TestTypecheckFails/unit/MergeLhsNotRecord",
	"TestTypecheckFails/unit/MergeWithWrongAnnotation",
	"TestTypecheckFails/unit/README", // FIXME, shouldn't need excluding
	"TestTypecheckFails/unit/RecordProjection",
	"TestTypecheckFails/unit/RecursiveRecordMerge",
	"TestTypecheckFails/unit/RecursiveRecordTypeMerge",
	"TestTypecheckFails/unit/RightBiasedRecordMerge",
	"TestTypecheckFails/unit/Some",
	"TestTypechecks/prelude",
	"TestTypechecks/simple/alternativesAreTypesA.dhall",
	"TestTypechecks/simple/mergeEquivalenceA.dhall",
	"TestTypechecks/simple/unionsOfTypesA.dhall",
	"TestTypeInference/simple/alternativesAreTypesA.dhall",
	"TestTypeInference/unit/DoubleShow",
	"TestTypeInference/unit/IntegerShow",
	"TestTypeInference/unit/IntegerToDouble",
	"TestTypeInference/unit/ListBuild",
	"TestTypeInference/unit/ListFold",
	"TestTypeInference/unit/ListHead",
	"TestTypeInference/unit/ListIndex",
	"TestTypeInference/unit/ListLast",
	"TestTypeInference/unit/ListLength",
	"TestTypeInference/unit/ListReverse",
	"TestTypeInference/unit/MergeOneA.dhall",               // uses union literals
	"TestTypeInference/unit/MergeOneWithAnnotationA.dhall", // uses union literals
	"TestTypeInference/unit/NaturalBuild",
	"TestTypeInference/unit/NaturalEven",
	"TestTypeInference/unit/NaturalFold",
	"TestTypeInference/unit/NaturalIsZero",
	"TestTypeInference/unit/NaturalOdd",
	"TestTypeInference/unit/NaturalShow",
	"TestTypeInference/unit/NaturalToInteger",
	"TestTypeInference/unit/OldOptional",
	"TestTypeInference/unit/OptionalBuild",
	"TestTypeInference/unit/OptionalFold",
	"TestTypeInference/unit/RecordNestedKind",
	"TestTypeInference/unit/RecordTypeKindLike",
	"TestTypeInference/unit/RecordTypeNestedKind",
	"TestTypeInference/unit/RecordProjection",
	"TestTypeInference/unit/RecursiveRecordMerge",
	"TestTypeInference/unit/RecursiveRecordTypeMerge",
	"TestTypeInference/unit/RightBiasedRecordMerge",
	"TestTypeInference/unit/TextShow",
	"TestTypeInference/unit/TypeAnnotationSort",
	"TestTypeInference/unit/UnionOne", // deprecated union literal syntax
	"TestNormalization/haskell-tutorial/combineTypes",
	"TestNormalization/haskell-tutorial/prefer",
	"TestNormalization/haskell-tutorial/projection",
	"TestNormalization/multiline",
	"TestNormalization/prelude/Bool/and",
	"TestNormalization/prelude/Bool/even",
	"TestNormalization/prelude/Bool/odd",
	"TestNormalization/prelude/Bool/or",
	"TestNormalization/prelude/Double",
	"TestNormalization/prelude/Integer",
	"TestNormalization/prelude/List",
	"TestNormalization/prelude/Natural",
	"TestNormalization/prelude/Optional",
	"TestNormalization/prelude/Text",
	"TestNormalization/remoteSystemsA.dhall",
	"TestNormalization/simple/doubleShowA.dhall",
	"TestNormalization/simple/integerShowA.dhall",
	"TestNormalization/simple/integerToDoubleA.dhall",
	"TestNormalization/simple/listBuildA.dhall",
	"TestNormalization/simple/naturalBuildA.dhall",
	"TestNormalization/simple/naturalShowA.dhall",
	"TestNormalization/simple/naturalToIntegerA.dhall",
	"TestNormalization/simple/optional",
	"TestNormalization/simple/sortOperatorA.dhall",
	"TestNormalization/unit/DoubleShow",
	"TestNormalization/unit/IntegerShow",
	"TestNormalization/unit/IntegerToDouble",
	"TestNormalization/unit/ListBuild",
	"TestNormalization/unit/ListFold",
	"TestNormalization/unit/ListHead",
	"TestNormalization/unit/ListIndexed",
	"TestNormalization/unit/ListLast",
	"TestNormalization/unit/ListLength",
	"TestNormalization/unit/ListReverse",
	"TestNormalization/unit/MergeWithTypeA", // uses union literals
	"TestNormalization/unit/NaturalBuild",
	"TestNormalization/unit/NaturalEven",
	"TestNormalization/unit/NaturalFold",
	"TestNormalization/unit/NaturalIsZero",
	"TestNormalization/unit/NaturalOdd",
	"TestNormalization/unit/NaturalShow",
	"TestNormalization/unit/NaturalToInteger",
	"TestNormalization/unit/NoneNaturalA", // I don't intend to implement this; it will disappear from the standard
	"TestNormalization/unit/OptionalBuild",
	"TestNormalization/unit/OptionalFold",
	"TestNormalization/unit/RecordProjection",
	"TestNormalization/unit/RecursiveRecordMerge",
	"TestNormalization/unit/RecursiveRecordTypeMerge",
	"TestNormalization/unit/RightBiasedRecordMerge",
	"TestNormalization/unit/TextShow",
	"TestNormalization/unit/UnionNormalize",
	"TestNormalization/unit/UnionSort",
	"TestImportFails/alternative",
	"TestImport/alternative", // needs alternative operator
	"TestImport/fieldOrder",  // needs import hashes
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

			annot := ast.Annot{
				Expr:       parsedA.(ast.Expr),
				Annotation: parsedB.(ast.Expr),
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
