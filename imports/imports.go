package imports

import (
	"bytes"
	"fmt"

	"github.com/philandstuff/dhall-golang/binary"
	"github.com/philandstuff/dhall-golang/core"
	. "github.com/philandstuff/dhall-golang/core"
	"github.com/philandstuff/dhall-golang/parser"
)

func resolveStringAsExpr(name, content string) (Term, error) {
	expr, err := parser.Parse(name, []byte(content))
	if err != nil {
		return nil, err
	}
	return expr.(Term), nil
}

// Load takes a Term and resolves all imports
func Load(e Term, ancestors ...Fetchable) (Term, error) {
	return LoadWith(StandardCache{}, e, ancestors...)
}

// LoadWith takes a Term and resolves all imports, using cache for
// saving and fetching imports
func LoadWith(cache DhallCache, e Term, ancestors ...Fetchable) (Term, error) {
	switch e := e.(type) {
	case Import:
		here := e.Fetchable
		origin := core.NullOrigin
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
			expr = TextLitTerm{Suffix: content}
		} else {
			// dynamicExpr may contain more imports
			dynamicExpr, err := resolveStringAsExpr(here.Name(), content)
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
		// check hash, if supplied
		if e.Hash != nil {
			actualHash, err := binary.SemanticHash(expr)
			if err != nil {
				return nil, err
			}
			if !bytes.Equal(e.Hash, actualHash[:]) {
				return nil, fmt.Errorf("Failed integrity check: expected %x but saw %x", e.Hash, actualHash)
			}
			// store in cache
			cache.Save(actualHash, expr)
		}
		return expr, nil
	case LambdaTerm:
		resolvedType, err := LoadWith(cache, e.Type, ancestors...)
		if err != nil {
			return nil, err
		}
		resolvedBody, err := LoadWith(cache, e.Body, ancestors...)
		if err != nil {
			return nil, err
		}
		return LambdaTerm{
			Label: e.Label,
			Type:  resolvedType,
			Body:  resolvedBody,
		}, nil
	case PiTerm:
		resolvedType, err := LoadWith(cache, e.Type, ancestors...)
		if err != nil {
			return nil, err
		}
		resolvedBody, err := LoadWith(cache, e.Body, ancestors...)
		if err != nil {
			return nil, err
		}
		return PiTerm{
			Label: e.Label,
			Type:  resolvedType,
			Body:  resolvedBody,
		}, nil
	case AppTerm:
		resolvedFn, err := LoadWith(cache, e.Fn, ancestors...)
		if err != nil {
			return nil, err
		}
		resolvedArg, err := LoadWith(cache, e.Arg, ancestors...)
		if err != nil {
			return nil, err
		}
		return Apply(
			resolvedFn,
			resolvedArg,
		), nil
	case Let:
		newBindings := make([]Binding, len(e.Bindings))
		for i, binding := range e.Bindings {
			var err error
			newBindings[i].Variable = binding.Variable
			if binding.Annotation != nil {
				newBindings[i].Annotation, err = LoadWith(cache, binding.Annotation, ancestors...)
				if err != nil {
					return nil, err
				}
			}
			newBindings[i].Value, err = LoadWith(cache, binding.Value, ancestors...)
			if err != nil {
				return nil, err
			}
		}
		resolvedBody, err := LoadWith(cache, e.Body, ancestors...)
		if err != nil {
			return nil, err
		}
		return Let{Bindings: newBindings, Body: resolvedBody}, nil
	case Annot:
		resolvedExpr, err := LoadWith(cache, e.Expr, ancestors...)
		if err != nil {
			return nil, err
		}
		resolvedAnnotation, err := LoadWith(cache, e.Annotation, ancestors...)
		if err != nil {
			return nil, err
		}
		return Annot{Expr: resolvedExpr, Annotation: resolvedAnnotation}, nil
	case TextLitTerm:
		var newChunks Chunks
		for _, chunk := range e.Chunks {
			resolvedExpr, err := LoadWith(cache, chunk.Expr, ancestors...)
			if err != nil {
				return nil, err
			}
			newChunks = append(newChunks, Chunk{
				Prefix: chunk.Prefix,
				Expr:   resolvedExpr,
			})
		}
		return TextLitTerm{newChunks, e.Suffix}, nil
	case IfTerm:
		resolvedCond, err := LoadWith(cache, e.Cond, ancestors...)
		if err != nil {
			return nil, err
		}
		resolvedT, err := LoadWith(cache, e.T, ancestors...)
		if err != nil {
			return nil, err
		}
		resolvedF, err := LoadWith(cache, e.F, ancestors...)
		if err != nil {
			return nil, err
		}
		return IfTerm{
			Cond: resolvedCond,
			T:    resolvedT,
			F:    resolvedF,
		}, nil
	case OpTerm:
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
		return OpTerm{OpCode: e.OpCode, L: resolvedL, R: resolvedR}, nil
	case EmptyList:
		resolvedType, err := LoadWith(cache, e.Type, ancestors...)
		if err != nil {
			return nil, err
		}
		return EmptyList{Type: resolvedType}, nil
	case NonEmptyList:
		newList := make(NonEmptyList, len(e))
		for i, item := range e {
			var err error
			newList[i], err = LoadWith(cache, item, ancestors...)
			if err != nil {
				return nil, err
			}
		}
		return newList, nil
	case Some:
		val, err := LoadWith(cache, e.Val, ancestors...)
		if err != nil {
			return nil, err
		}
		return Some{val}, nil
	case RecordType:
		newRecord := make(RecordType, len(e))
		for k, v := range e {
			var err error
			newRecord[k], err = LoadWith(cache, v, ancestors...)
			if err != nil {
				return nil, err
			}
		}
		return newRecord, nil
	case RecordLit:
		newRecord := make(RecordLit, len(e))
		for k, v := range e {
			var err error
			newRecord[k], err = LoadWith(cache, v, ancestors...)
			if err != nil {
				return nil, err
			}
		}
		return newRecord, nil
	case ToMap:
		record, err := LoadWith(cache, e.Record, ancestors...)
		if err != nil {
			return nil, err
		}
		typ, err := LoadWith(cache, e.Type, ancestors...)
		if err != nil {
			return nil, err
		}
		return ToMap{Record: record, Type: typ}, nil
	case Field:
		newRecord, err := LoadWith(cache, e.Record, ancestors...)
		if err != nil {
			return nil, err
		}
		return Field{Record: newRecord, FieldName: e.FieldName}, nil
	case Project:
		newRecord, err := LoadWith(cache, e.Record, ancestors...)
		if err != nil {
			return nil, err
		}
		return Project{Record: newRecord, FieldNames: e.FieldNames}, nil
	case ProjectType:
		record, err := LoadWith(cache, e.Record, ancestors...)
		if err != nil {
			return nil, err
		}
		typ, err := LoadWith(cache, e.Selector, ancestors...)
		if err != nil {
			return nil, err
		}
		return ProjectType{Record: record, Selector: typ}, nil
	case UnionType:
		result := make(UnionType, len(e))
		for k, v := range e {
			var err error
			if v == nil {
				result[k] = nil
				continue
			}
			result[k], err = LoadWith(cache, v, ancestors...)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	case Merge:
		handler, err := LoadWith(cache, e.Handler, ancestors...)
		if err != nil {
			return nil, err
		}
		union, err := LoadWith(cache, e.Union, ancestors...)
		if err != nil {
			return nil, err
		}
		return Merge{Handler: handler, Union: union}, nil
	case Assert:
		annot, err := LoadWith(cache, e.Annotation, ancestors...)
		if err != nil {
			return nil, err
		}
		return Assert{Annotation: annot}, nil
	default:
		// Const, NaturalLit, etc
		return e, nil
	}
}
