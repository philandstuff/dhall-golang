package imports

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/philandstuff/dhall-golang/binary"
	"github.com/philandstuff/dhall-golang/core"
)

// DhallCache is an interface for caching implementations.
type DhallCache interface {
	// Fetch fetches a Term from the cache
	Fetch(hash []byte) core.Term
	// Save saves a Term to the cache
	Save(hash []byte, term core.Term)
}

// StandardCache is the standard DhallCache implementation.
type StandardCache struct{}

func dhallCacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return path.Join(cacheDir, "dhall"), nil
}

// Fetch implements DhallCache
func (StandardCache) Fetch(hash []byte) core.Term {
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

// Save implements DhallCache
func (StandardCache) Save(hash []byte, e core.Term) {
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

// NoCache is a DhallCache which doesn't do any caching.  It might be
// useful for testing.
type NoCache struct{}

// Fetch implements DhallCache
func (NoCache) Fetch([]byte) core.Term { return nil }

// Save implements DhallCache
func (NoCache) Save([]byte, core.Term) {}
