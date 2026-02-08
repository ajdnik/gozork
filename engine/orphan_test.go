package engine

import (
	"testing"
)

func TestOrphanCopiesParsedSyntax(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.ParsedSyntx.Verb = LexItem{Norm: "take", Orig: "take", Types: WordTypes{WordVerb}}
	G.ParsedSyntx.ObjOrClause1 = []LexItem{
		{Norm: "lamp", Orig: "lamp", Types: WordTypes{WordObj}},
	}
	G.Params.ObjOrClauseCnt = 1

	s1 := &Syntax{VrbPrep: "with"}
	s2 := &Syntax{ObjPrep: "using"}

	Orphan(s1, s2)
	if !G.OrphanedSyntx.Verb.IsSet() {
		t.Fatalf("expected Orphan to copy verb")
	}
	if len(G.OrphanedSyntx.ObjOrClause1) != 0 {
		t.Fatalf("expected ObjOrClause1 to be cleared when first syntax provided")
	}
	if G.OrphanedSyntx.Prep1.Norm != "with" {
		t.Fatalf("expected prep1 to be set from first syntax")
	}
}
