package engine

import (
	"bytes"
	"strings"
	"testing"
)

func TestWordTypesEquals(t *testing.T) {
	a := WordTypes{WordVerb, WordObj}
	b := WordTypes{WordObj, WordVerb}
	if !a.Equals(b) {
		t.Fatalf("expected Equals to ignore order")
	}
	c := WordTypes{WordVerb}
	if a.Equals(c) {
		t.Fatalf("expected Equals to differ on length")
	}
}

func TestLexItemSetAndMatches(t *testing.T) {
	itm := LexItem{Norm: "take", Orig: "take", Types: WordTypes{WordVerb}}
	var dst LexItem
	dst.Set(itm)
	if !dst.Matches(itm) {
		t.Fatalf("expected Set to copy values")
	}
	if !dst.IsSet() {
		t.Fatalf("expected IsSet to be true")
	}
	dst.Clear()
	if dst.IsSet() {
		t.Fatalf("expected Clear to reset item")
	}
}

func TestTokenize(t *testing.T) {
	toks := Tokenize("go2west, now!")
	want := []string{"go", "2", "west", ",", "now", "!"}
	if len(toks) != len(want) {
		t.Fatalf("unexpected token count: %v", toks)
	}
	for i := range want {
		if toks[i] != want[i] {
			t.Fatalf("token %d: expected %q, got %q", i, want[i], toks[i])
		}
	}
}

func TestTokenizeEdgeCases(t *testing.T) {
	toks := Tokenize("a1b..c")
	want := []string{"a", "1", "b", "..", "c"}
	if len(toks) != len(want) {
		t.Fatalf("unexpected token count: %v", toks)
	}
	for i := range want {
		if toks[i] != want[i] {
			t.Fatalf("token %d: expected %q, got %q", i, want[i], toks[i])
		}
	}
}

func TestLexAndRead(t *testing.T) {
	oldG := G
	oldVocab := Vocabulary
	G = NewGameState()
	t.Cleanup(func() {
		G = oldG
		Vocabulary = oldVocab
	})

	Vocabulary = map[string]WordItem{
		"take": {Norm: "take", Types: WordTypes{WordVerb}},
		"lamp": {Norm: "lamp", Types: WordTypes{WordObj}},
	}

	G.GameInput = bytes.NewBufferString("Take lamp\n")
	InitReader()
	txt, items := Read()
	if strings.ToLower(txt) != "take lamp" {
		t.Fatalf("expected normalized input, got %q", txt)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 lex items, got %d", len(items))
	}
	if items[0].Norm != "take" || !items[0].Types.Has(WordVerb) {
		t.Fatalf("expected verb token for take")
	}
	if items[1].Norm != "lamp" || !items[1].Types.Has(WordObj) {
		t.Fatalf("expected object token for lamp")
	}
}

func TestLexUnknownToken(t *testing.T) {
	oldVocab := Vocabulary
	Vocabulary = make(map[string]WordItem)
	t.Cleanup(func() { Vocabulary = oldVocab })

	items := Lex([]string{"??"})
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Norm != "??" || items[0].Orig != "??" {
		t.Fatalf("expected unknown token to preserve norm/orig")
	}
	if items[0].Types != nil {
		t.Fatalf("expected unknown token to have nil types")
	}
}

func TestReadEOF(t *testing.T) {
	oldG := G
	oldVocab := Vocabulary
	G = NewGameState()
	t.Cleanup(func() {
		G = oldG
		Vocabulary = oldVocab
	})

	Vocabulary = map[string]WordItem{
		"take": {Norm: "take", Types: WordTypes{WordVerb}},
	}
	G.GameInput = bytes.NewBufferString("take")
	G.Reader = nil

	txt, items := Read()
	if txt != "take" {
		t.Fatalf("expected txt to be returned on EOF, got %q", txt)
	}
	if !G.InputExhausted {
		t.Fatalf("expected InputExhausted to be true at EOF")
	}
	if len(items) != 1 || items[0].Norm != "take" {
		t.Fatalf("unexpected lex items at EOF: %+v", items)
	}
}
