package dhall_test

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

	"github.com/philandstuff/dhall-golang/binary"
	"github.com/philandstuff/dhall-golang/core"
	"github.com/philandstuff/dhall-golang/imports"
	"github.com/philandstuff/dhall-golang/parser"
	"github.com/pkg/errors"
)

var slowTests = []string{
	"TestParserAccepts/largeExpressionA",
	"TestTypeInference/preludeA",
}

var expectedFailures = []string{
	// needs `using`
	"TestParserAccepts/unit/import/Headers",
	"TestParserAccepts/unit/import/inlineUsing",
	"TestParserAccepts/unit/import/parenthesizeUsing",
	"TestTypeInferenceFails/customHeadersUsingBoundVariable",
	"TestImport/customHeadersA.dhall",
	"TestImport/headerForwardingA.dhall",
	"TestImport/noHeaderForwardingA.dhall",
	"TestImportFails/customHeadersUsingBoundVariable",

	// needs bigint support
	"TestNormalization/simple/integerToDoubleA.dhall",
	"TestSemanticHash/simple/integerToDouble",

	// needs quoted paths in URLs
	"TestParserAccepts/unit/import/quotedPaths",     // needs.. quoted paths
	"TestParserAccepts/unit/import/urls/quotedPath", // needs quotedPaths

	// urgh IEEE floating point; NaN != NaN :(
	"TestTypeInference/unit/AssertNaNA",

	// other
	"TestParserAccepts/unit/import/urls/potPourri", // net/url doesn't parse authorities in the way the test expects

	// in dhall-golang, duplicate fields & alternatives are a parse error, not a
	// type error
	"TestTypeInferenceFails/unit/RecordTypeDuplicateFields.dhall",
	"TestTypeInferenceFails/unit/RecordLitDuplicateFields.dhall",
	"TestTypeInferenceFails/unit/RecordProjectionDuplicateFields.dhall",
	"TestTypeInferenceFails/unit/UnionTypeDuplicateVariants",
	"TestTypeInferenceFails/unit/README", // FIXME, shouldn't need excluding

	// since NbE, we don't deal with unbound variables in the
	// "standard" way, because we don't have a shift() function
	"TestAlphaNormalization/unit/FunctionNestedBindingXXFree",

	// this passes with an empty cache. the tests should disable the
	// cache for TestTypeInference but don't yet.
	"TestTypeInference/preludeA",
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
			t.Skip("Skipping expected failure")
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

func expectEqualTerms(t *testing.T, expected, actual core.Term) {
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
				if testing.Short() {
					for _, prefix := range slowTests {
						if strings.HasPrefix(t.Name(), prefix) {
							t.Skip("Skipping slow tests in short mode")
						}
					}
				}
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
					if testing.Short() {
						for _, prefix := range slowTests {
							if strings.HasPrefix(t.Name(), prefix) {
								t.Skip("Skipping slow tests in short mode")
							}
						}
					}

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
	t.Parallel()
	runTestOnEachFile(t, "dhall-lang/tests/parser/failure/", func(t *testing.T, testPath string) {
		_, err := parser.ParseFile(testPath)

		expectError(t, err)
	})
}

func TestParserAccepts(t *testing.T) {
	t.Parallel()
	runTestOnFilePairs(t, "dhall-lang/tests/parser/success/",
		"A.dhall", "B.dhallb",
		func(t *testing.T, aPath, bPath string) {
			actualBuf := new(bytes.Buffer)
			parsed, err := parser.ParseFile(aPath)
			expectNoError(t, err)

			err = binary.EncodeAsCbor(actualBuf, parsed.(core.Term))
			expectNoError(t, err)

			expected, err := ioutil.ReadFile(bPath)
			expectNoError(t, err)
			expectEqualCbor(t, expected, actualBuf.Bytes())
		})
}

func BenchmarkParserLargeExpression(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := parser.ParseFile("dhall-lang/tests/parser/success/largeExpressionA.dhall")
		if err != nil {
			b.Fatalf("Parse error: %v", err)
		}
	}
}

func isSimpleTest(testName string) bool {
	return strings.Contains(testName, "/unit/") ||
		strings.Contains(testName, "/simple/")
}

func TestTypeInference(t *testing.T) {
	t.Parallel()
	runTestOnFilePairs(t, "dhall-lang/tests/type-inference/success/",
		"A.dhall", "B.dhall",
		func(t *testing.T, aPath, bPath string) {
			parsedA, err := parser.ParseFile(aPath)
			expectNoError(t, err)

			parsedB, err := parser.ParseFile(bPath)
			expectNoError(t, err)

			var resolvedA core.Term
			if isSimpleTest(t.Name()) {
				resolvedA = parsedA.(core.Term)
			} else {
				// FIXME: should disable cache here
				resolvedA, err = imports.Load(parsedA.(core.Term), core.Local(aPath))
				expectNoError(t, err)
			}

			inferredType, err := core.TypeOf(resolvedA)
			expectNoError(t, err)

			expectEqualTerms(t, parsedB.(core.Term), inferredType)
		})
}

func TestTypeInferenceFails(t *testing.T) {
	t.Parallel()
	runTestOnEachFile(t, "dhall-lang/tests/type-inference/failure/", func(t *testing.T, testPath string) {
		parsed, err := parser.ParseFile(testPath)

		expectNoError(t, err)

		expr, ok := parsed.(core.Term)
		if !ok {
			failf(t, "Expected core.Term, got %+v\n", parsed)
		}

		_, err = core.TypeOf(expr)

		expectError(t, err)
	})
}

func TestAlphaNormalization(t *testing.T) {
	t.Parallel()
	runTestOnFilePairs(t, "dhall-lang/tests/alpha-normalization/success/",
		"A.dhall", "B.dhall",
		func(t *testing.T, aPath, bPath string) {
			parsedA, err := parser.ParseFile(aPath)
			expectNoError(t, err)

			parsedB, err := parser.ParseFile(bPath)
			expectNoError(t, err)

			normA := core.Quote(core.AlphaBetaEval(parsedA.(core.Term)))
			normB := core.Quote(core.AlphaBetaEval(parsedB.(core.Term)))

			expectEqualTerms(t, normB, normA)
		})
}

func TestNormalization(t *testing.T) {
	t.Parallel()
	runTestOnFilePairs(t, "dhall-lang/tests/normalization/success/",
		"A.dhall", "B.dhall",
		func(t *testing.T, aPath, bPath string) {
			parsedA, err := parser.ParseFile(aPath)
			expectNoError(t, err)

			parsedB, err := parser.ParseFile(bPath)
			expectNoError(t, err)

			var resolvedA, resolvedB core.Term
			if isSimpleTest(t.Name()) {
				resolvedA = parsedA.(core.Term)
				resolvedB = parsedB.(core.Term)
			} else {

				resolvedA, err = imports.Load(parsedA.(core.Term), core.Local(aPath))
				expectNoError(t, err)

				resolvedB, err = imports.Load(parsedB.(core.Term), core.Local(bPath))
				expectNoError(t, err)
			}

			normA := core.Eval(resolvedA)

			expectEqualTerms(t, resolvedB, core.Quote(normA))
		})
}

func TestImportFails(t *testing.T) {
	t.Parallel()
	runTestOnEachFile(t, "dhall-lang/tests/import/failure/", func(t *testing.T, testPath string) {
		parsed, err := parser.ParseFile(testPath)
		expectNoError(t, err)

		_, err = imports.Load(parsed.(core.Term), core.Local(testPath))
		expectError(t, err)
	})
}

func TestImport(t *testing.T) {
	t.Parallel()
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

			resolvedA, err := imports.Load(parsedA.(core.Term), core.Local(aPath))
			expectNoError(t, err)

			resolvedB, err := imports.Load(parsedB.(core.Term), core.Local(bPath))
			expectNoError(t, err)

			expectEqualTerms(t, resolvedB, resolvedA)
		})
}

func TestSemanticHash(t *testing.T) {
	t.Parallel()
	sha256re := regexp.MustCompile("^sha256:([0-9a-fA-F]{64})\n$")
	runTestOnFilePairs(t, "dhall-lang/tests/semantic-hash/success/",
		"A.dhall", "B.hash",
		func(t *testing.T, aPath, bPath string) {
			parsedA, err := parser.ParseFile(aPath)
			expectNoError(t, err)

			var resolvedA core.Term
			if isSimpleTest(t.Name()) {
				resolvedA = parsedA.(core.Term)
			} else {
				resolvedA, err = imports.Load(parsedA.(core.Term), core.Local(aPath))
				expectNoError(t, err)
			}

			actualHash, err := binary.SemanticHash(resolvedA)
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
	t.Parallel()
	runTestOnFilePairs(t, "dhall-lang/tests/binary-decode/success/",
		"A.dhallb", "B.dhall",
		func(t *testing.T, aPath, bPath string) {
			aReader, err := os.Open(aPath)
			expectNoError(t, err)
			defer aReader.Close()
			expr, err := binary.DecodeAsCbor(aReader)
			expectNoError(t, err)

			parsedB, err := parser.ParseFile(bPath)
			expectNoError(t, err)

			expectEqualTerms(t, parsedB.(core.Term), expr)
		})
}

func TestBinaryDecodeFails(t *testing.T) {
	t.Parallel()
	runTestOnEachFile(t, "dhall-lang/tests/binary-decode/failure/", func(t *testing.T, testPath string) {
		reader, err := os.Open(testPath)
		if err != nil {
			t.Fatal(err)
		}
		defer reader.Close()
		_, err = binary.DecodeAsCbor(reader)
		expectError(t, err)
	})
}
