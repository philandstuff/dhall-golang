package imports

import (
	"bytes"
	"fmt"

	"github.com/philandstuff/dhall-golang/v5/binary"
	"github.com/philandstuff/dhall-golang/v5/core"
	"github.com/philandstuff/dhall-golang/v5/parser"
	"github.com/philandstuff/dhall-golang/v5/term"
	. "github.com/philandstuff/dhall-golang/v5/term"
)

// Load takes a Term and resolves all imports
func Load(e Term, ancestors ...Fetchable) (Term, error) {
	cache, err := StandardCache()
	if err != nil {
		return nil, err
	}
	return LoadWith(cache, e, ancestors...)
}

// LoadWith takes a Term and resolves all imports, using cache for
// saving and fetching imports
func LoadWith(cache DhallCache, e Term, ancestors ...Fetchable) (Term, error) {
	switch e := e.(type) {
	case Import:
		here := e.Fetchable
		origin := term.NullOrigin
		if len(ancestors) >= 1 {
			origin = ancestors[len(ancestors)-1].Origin()

			var err error
			here, err = here.ChainOnto(ancestors[len(ancestors)-1])
			if err != nil {
				return nil, err
			}
		}
		if e.ImportMode == Location {
			return here.AsLocation(), nil
		}

		for _, ancestor := range ancestors {
			if ancestor == here {
				return nil, fmt.Errorf("Detected import cycle in %s", ancestor)
			}
		}
		if e.Hash != nil {
			// fetch from cache if available
			if expr := cache.Fetch(e.Hash); expr != nil {
				return expr, nil
			}
		}
		imports := append(ancestors, here)
		content, err := here.Fetch(origin)
		if err != nil {
			return nil, err
		}
		var expr Term
		if e.ImportMode == RawText {
			expr = PlainText(content)
		} else {
			// dynamicExpr may contain more imports
			dynamicExpr, err := parser.Parse(here.String(), []byte(content))
			if err != nil {
				return nil, err
			}

			// recursively load any more imports
			expr, err = LoadWith(cache, dynamicExpr, imports...)
			if err != nil {
				return nil, err
			}

			// ensure that expr typechecks in empty context
			_, err = core.TypeOf(expr)
			if err != nil {
				return nil, err
			}
		}

		// evaluate expression
		exprVal := core.Eval(expr)
		expr = core.Quote(exprVal)

		// check hash, if supplied
		if e.Hash != nil {
			actualHash, err := binary.SemanticHash(exprVal)
			if err != nil {
				return nil, err
			}
			if !bytes.Equal(e.Hash, actualHash[:]) {
				return nil, fmt.Errorf("Failed integrity check: expected %x but saw %x", e.Hash, actualHash)
			}
			// store in cache
			cache.Save(actualHash, core.QuoteAlphaNormal(exprVal))
		}
		return expr, nil
	case Op:
		if e.OpCode == ImportAltOp {
			resolvedL, err := LoadWith(cache, e.L, ancestors...)
			if err == nil {
				return resolvedL, nil
			}
			resolvedR, err := LoadWith(cache, e.R, ancestors...)
			if err != nil {
				return nil, err
			}
			return resolvedR, nil
		}
		resolvedL, err := LoadWith(cache, e.L, ancestors...)
		if err != nil {
			return nil, err
		}
		resolvedR, err := LoadWith(cache, e.R, ancestors...)
		if err != nil {
			return nil, err
		}
		return Op{OpCode: e.OpCode, L: resolvedL, R: resolvedR}, nil
	default:
		// Const, NaturalLit, etc
		return term.MaybeTransformSubexprs(e, func(t Term) (Term, error) {
			return LoadWith(cache, t, ancestors...)
		})
	}
}
