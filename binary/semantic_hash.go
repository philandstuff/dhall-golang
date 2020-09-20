package binary

import (
	"bytes"
	"crypto/sha256"

	"github.com/philandstuff/dhall-golang/v5/core"
)

// SemanticHash returns the semantic hash of an evaluated expression.
// The semantic hash is defined as the multihash-encoded sha256 sum of the CBOR
// representation of the fully alpha-beta-normalized expression.
func SemanticHash(e core.Value) ([]byte, error) {
	norm := core.QuoteAlphaNormal(e)
	var buf bytes.Buffer
	err := EncodeAsCbor(&buf, norm)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(buf.Bytes())
	return append([]byte{0x12, 0x20}, hash[:]...), nil
}
