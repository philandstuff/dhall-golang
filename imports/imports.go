package imports

import (
	"bytes"
	"fmt"

	"github.com/philandstuff/dhall-golang/binary"
	"github.com/philandstuff/dhall-golang/core"
	. "github.com/philandstuff/dhall-golang/core"
	"github.com/philandstuff/dhall-golang/parser"
)

func ResolveStringAsExpr(name, content string) (Term, error) {
	expr, err := parser.Parse(name, []byte(content))
	if err != nil {
		return nil, err
	}
	return expr.(Term), nil
}

func Load(e Term, ancestors ...Fetchable) (Term, error) {
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
			if expr := fetchFromCache(e.Hash); expr != nil {
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
			dynamicExpr, err := ResolveStringAsExpr(here.Name(), content)
			if err != nil {
				return nil, err
			}

			// recursively load any more imports
			expr, err = Load(dynamicExpr, imports...)
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
			saveToCache(actualHash, expr)
		}
		return expr, nil
	case LambdaTerm:
		resolvedType, err := Load(e.Type, ancestors...)
		if err != nil {
			return nil, err
		}
		resolvedBody, err := Load(e.Body, ancestors...)
		if err != nil {
			return nil, err
		}
		return LambdaTerm{
			Label: e.Label,
			Type:  resolvedType,
			Body:  resolvedBody,
		}, nil
	case PiTerm:
		resolvedType, err := Load(e.Type, ancestors...)
		if err != nil {
			return nil, err
		}
		resolvedBody, err := Load(e.Body, ancestors...)
		if err != nil {
			return nil, err
		}
		return PiTerm{
			Label: e.Label,
			Type:  resolvedType,
			Body:  resolvedBody,
		}, nil
	case AppTerm:
		resolvedFn, err := Load(e.Fn, ancestors...)
		if err != nil {
			return nil, err
		}
		resolvedArg, err := Load(e.Arg, ancestors...)
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
				newBindings[i].Annotation, err = Load(binding.Annotation, ancestors...)
				if err != nil {
					return nil, err
				}
			}
			newBindings[i].Value, err = Load(binding.Value, ancestors...)
			if err != nil {
				return nil, err
			}
		}
		resolvedBody, err := Load(e.Body, ancestors...)
		if err != nil {
			return nil, err
		}
		return Let{Bindings: newBindings, Body: resolvedBody}, nil
	case Annot:
		resolvedExpr, err := Load(e.Expr, ancestors...)
		if err != nil {
			return nil, err
		}
		resolvedAnnotation, err := Load(e.Annotation, ancestors...)
		if err != nil {
			return nil, err
		}
		return Annot{Expr: resolvedExpr, Annotation: resolvedAnnotation}, nil
	case TextLitTerm:
		var newChunks Chunks
		for _, chunk := range e.Chunks {
			resolvedExpr, err := Load(chunk.Expr, ancestors...)
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
		resolvedCond, err := Load(e.Cond, ancestors...)
		if err != nil {
			return nil, err
		}
		resolvedT, err := Load(e.T, ancestors...)
		if err != nil {
			return nil, err
		}
		resolvedF, err := Load(e.F, ancestors...)
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
			resolvedL, err := Load(e.L, ancestors...)
			if err == nil {
				return resolvedL, nil
			}
			resolvedR, err := Load(e.R, ancestors...)
			if err != nil {
				return nil, err
			}
			return resolvedR, nil
		}
		resolvedL, err := Load(e.L, ancestors...)
		if err != nil {
			return nil, err
		}
		resolvedR, err := Load(e.R, ancestors...)
		if err != nil {
			return nil, err
		}
		return OpTerm{OpCode: e.OpCode, L: resolvedL, R: resolvedR}, nil
	case EmptyList:
		resolvedType, err := Load(e.Type, ancestors...)
		if err != nil {
			return nil, err
		}
		return EmptyList{Type: resolvedType}, nil
	case NonEmptyList:
		newList := make(NonEmptyList, len(e))
		for i, item := range e {
			var err error
			newList[i], err = Load(item, ancestors...)
			if err != nil {
				return nil, err
			}
		}
		return newList, nil
	case RecordType:
		newRecord := make(RecordType, len(e))
		for k, v := range e {
			var err error
			newRecord[k], err = Load(v, ancestors...)
			if err != nil {
				return nil, err
			}
		}
		return newRecord, nil
	case RecordLit:
		newRecord := make(RecordLit, len(e))
		for k, v := range e {
			var err error
			newRecord[k], err = Load(v, ancestors...)
			if err != nil {
				return nil, err
			}
		}
		return newRecord, nil
	case Field:
		newRecord, err := Load(e.Record, ancestors...)
		if err != nil {
			return nil, err
		}
		return Field{Record: newRecord, FieldName: e.FieldName}, nil
		// TODO: other new terms (Project, ProjectType etc)
	default:
		// Const, NaturalLit, etc
		return e, nil
	}
}
