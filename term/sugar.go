package term

// Desugar implements `desugar-with` from the Dhall standard
func (w With) Desugar() Term {
	if len(w.Path) == 1 {
		return Op{
			OpCode: RightBiasedRecordMergeOp,
			L:      w.Record,
			R:      RecordLit{w.Path[0]: w.Value},
		}
	}
	// TODO: add a let
	return Op{
		OpCode: RightBiasedRecordMergeOp,
		L:      w.Record,
		R: RecordLit{w.Path[0]: With{
			Field{w.Record, w.Path[0]},
			w.Path[1:],
			w.Value,
		}.Desugar()},
	}
}
