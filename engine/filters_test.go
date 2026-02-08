package engine

import (
	"testing"
)

func TestIsThisItFilters(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	obj := &Object{
		Desc:       "brass lamp",
		Synonyms:   []string{"lamp"},
		Adjectives: []string{"brass"},
		Flags:      FlgTake,
	}

	G.Search.Syn = LexItem{Norm: "lamp", Orig: "lamp", Types: WordTypes{WordObj}}
	G.Search.Adj = LexItem{Norm: "brass", Orig: "brass", Types: WordTypes{WordAdj}}
	G.Search.ObjFlags = FlgTake
	if !IsThisIt(obj) {
		t.Fatalf("expected IsThisIt to match syn/adj/flags")
	}

	G.Search.Adj = LexItem{Norm: "red", Orig: "red", Types: WordTypes{WordAdj}}
	if IsThisIt(obj) {
		t.Fatalf("expected IsThisIt to fail on adjective mismatch")
	}

	G.Search.Adj = LexItem{Norm: "brass", Orig: "brass", Types: WordTypes{WordAdj}}
	G.Search.ObjFlags = FlgOpen
	if IsThisIt(obj) {
		t.Fatalf("expected IsThisIt to fail on flag mismatch")
	}

	obj.Flags = FlgTake | FlgInvis
	G.Search.ObjFlags = FlgTake
	if IsThisIt(obj) {
		t.Fatalf("expected IsThisIt to fail for invisible objects")
	}
}
