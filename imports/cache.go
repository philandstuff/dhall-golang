package imports

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"path"

	"github.com/philandstuff/dhall-golang/binary"
	"github.com/philandstuff/dhall-golang/term"
)

// DhallCache is an interface for caching implementations.
type DhallCache interface {
	// Fetch fetches a Term from the cache
	Fetch(hash []byte) term.Term
	// Save saves a Term to the cache
	Save(hash []byte, term term.Term)
}

// A LocalCache is a cache for normalized Dhall expressions, stored in
// binary form.
//
// Note that a LocalCache has a shared hash.Hash instance, and as such
// is not safe for concurrent use across goroutines.
type LocalCache struct {
	path string
	hash hash.Hash
}

// FIXME: we might want to move the hash.Hash implementation to a
// context?

// NewLocalCache creates a new LocalCache, with the cache store at the
// given path.
func NewLocalCache(path string) LocalCache {
	return LocalCache{path, sha256.New()}
}

// Fetch searches the LocalCache for a term at the index given by
// hash.  If the hash isn't in the cache, returns nil.
func (l LocalCache) Fetch(hash []byte) term.Term {
	// FIXME: don't swallow these errors, maybe?
	hash16 := fmt.Sprintf("%x", hash)

	reader, err := os.Open(path.Join(l.path, hash16))
	if err != nil {
		return nil
	}

	defer l.hash.Reset()
	io.Copy(l.hash, reader)
	if bytes.Compare(hash[2:], l.hash.Sum(nil)) != 0 {
		log.Printf("warning: invalid cache entry for %x, ignoring\n", hash)
		return nil
	}

	reader.Seek(0, io.SeekStart)

	expr, err := binary.DecodeAsCbor(reader)
	if err != nil {
		log.Println(err)
		return nil
	}
	return expr
}

// Save saves the given Term to the LocalCache at the given hash.
func (l LocalCache) Save(hash []byte, e term.Term) {
	hash16 := fmt.Sprintf("%x", hash)
	file, err := os.Create(path.Join(l.path, hash16))
	if err != nil {
		return
	}
	defer file.Close()
	// ignores returned error
	binary.EncodeAsCbor(file, e)
}

// StandardCache is the standard DhallCache implementation.  It is a
// LocalCache in the standard Dhall cache directory.
func StandardCache() (DhallCache, error) {
	cacheDir, err := dhallCacheDir()
	if err != nil {
		return nil, err
	}
	return NewLocalCache(cacheDir), nil
}

func dhallCacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return path.Join(cacheDir, "dhall"), nil
}

// NoCache is a DhallCache which doesn't do any caching.  It might be
// useful for testing.
type NoCache struct{}

// Fetch always returns nil.
func (NoCache) Fetch([]byte) term.Term { return nil }

// Save does nothing.
func (NoCache) Save([]byte, term.Term) {}
