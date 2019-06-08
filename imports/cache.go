package imports

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/philandstuff/dhall-golang/ast"
)

func dhallCacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return path.Join(cacheDir, "dhall"), nil
}

func fetchFromCache(hash []byte) ast.Expr {
	// FIXME: don't swallow these errors, maybe?
	hash16 := fmt.Sprintf("%x", hash)
	dir, err := dhallCacheDir()
	if err != nil {
		return nil
	}
	reader, err := os.Open(path.Join(dir, hash16))
	if err != nil {
		return nil
	}
	expr, err := ast.DecodeAsCbor(reader)
	if err != nil {
		log.Println(err)
		return nil
	}
	return expr
}

func saveToCache(hash []byte, e ast.Expr) {
	hash16 := fmt.Sprintf("%x", hash)
	dir, err := dhallCacheDir()
	if err != nil {
		return
	}
	file, err := os.Create(path.Join(dir, hash16))
	if err != nil {
		return
	}
	defer file.Close()
	err = ast.EncodeAsCbor(file, e)
}
