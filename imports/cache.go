package imports

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/philandstuff/dhall-golang/binary"
	"github.com/philandstuff/dhall-golang/core"
)

func dhallCacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return path.Join(cacheDir, "dhall"), nil
}

func fetchFromCache(hash []byte) core.Term {
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
	expr, err := binary.DecodeAsCbor(reader)
	if err != nil {
		log.Println(err)
		return nil
	}
	return expr
}

func saveToCache(hash []byte, e core.Term) {
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
	err = binary.EncodeAsCbor(file, e)
}
