package core

type textValBuilder struct {
	Chunks chunks
	Suffix string
}

func (t *textValBuilder) appendStr(s string) {
	t.Suffix += s
}

func (t *textValBuilder) appendValue(v Value) {
	if plainText, ok := v.(PlainTextLit); ok {
		t.Suffix += string(plainText)
	} else if interpolatedText, ok := v.(interpolatedText); ok {
		t.Chunks = append(t.Chunks, chunk{
			Prefix: t.Suffix + interpolatedText.Chunks[0].Prefix,
			Expr:   interpolatedText.Chunks[0].Expr,
		})
		t.Chunks = append(t.Chunks, interpolatedText.Chunks[1:]...)
		t.Suffix = interpolatedText.Suffix
	} else {
		t.Chunks = append(t.Chunks, chunk{Prefix: t.Suffix, Expr: v})
		t.Suffix = ""
	}
}

func (t *textValBuilder) value() Value {
	if len(t.Chunks) == 0 {
		return PlainTextLit(t.Suffix)
	} else if len(t.Chunks) == 1 && t.Chunks[0].Prefix == "" && t.Suffix == "" {
		// special case: "${x}" -> x
		return t.Chunks[0].Expr
	}
	return interpolatedText(*t)
}
