package binary_test

import (
	"fmt"
	"io"
	"os"
	"path"
	"testing"

	"github.com/philandstuff/dhall-golang/v6/binary"
	"github.com/philandstuff/dhall-golang/v6/core"
	"github.com/philandstuff/dhall-golang/v6/imports"
	"github.com/philandstuff/dhall-golang/v6/internal"
	"github.com/philandstuff/dhall-golang/v6/term"
)

func BenchmarkDecodeLargeExpression(b *testing.B) {
	file, err := os.Open("../dhall-lang/tests/parser/success/largeExpressionB.dhallb")
	if err != nil {
		b.Fatal(err)
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(fileinfo.Size())
	for n := 0; n < b.N; n++ {
		file.Seek(0, io.SeekStart)
		binary.DecodeAsCbor(file)
	}
}

func BenchmarkDecodePrelude(b *testing.B) {
	file := internal.NewLocalImport("../dhall-lang/Prelude/package.dhall", term.Code)
	t, err := imports.Load(file)
	if err != nil {
		b.Fatal(err)
	}
	v := core.Eval(t)
	cache, err := imports.StandardCache()
	if err != nil {
		b.Fatal(err)
	}
	hash, err := binary.SemanticHash(v)
	if err != nil {
		b.Fatal(err)
	}
	// ensure Prelude is saved into the cache
	cache.Save(hash, core.Quote(v))

	cacheDir, err := imports.DhallCacheDir()
	if err != nil {
		b.Fatal(err)
	}
	hash16 := fmt.Sprintf("%x", hash)
	cborfile, err := os.Open(path.Join(cacheDir, hash16))
	if err != nil {
		b.Fatal(err)
	}
	fileinfo, err := cborfile.Stat()
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(fileinfo.Size())
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		cborfile.Seek(0, io.SeekStart)
		binary.DecodeAsCbor(cborfile)
	}
}
