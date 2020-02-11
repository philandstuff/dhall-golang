package binary

import (
	"bytes"
	"crypto/sha256"

	"github.com/philandstuff/dhall-golang/core"
	"github.com/philandstuff/dhall-golang/term"
)

// SemanticHash returns the semantic hash of an expression.
// The semantic hash is defined as the multihash-encoded sha256 sum of the CBOR
// representation of the fully alpha-beta-normalized expression.
func SemanticHash(e term.Term) ([]byte, error) {
	norm := core.AlphaBetaEval(e)
	var buf bytes.Buffer
	err := EncodeAsCbor(&buf, core.Quote(norm))
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(buf.Bytes())
	return append([]byte{0x12, 0x20}, hash[:]...), nil
}
