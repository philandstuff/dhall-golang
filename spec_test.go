package main_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/philandstuff/dhall-golang/parser"
)

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

			if err == nil {
				t.Errorf("Expected file %s to fail to parse, but it parsed successfully", f.Name())
			}
		})
	}
}
