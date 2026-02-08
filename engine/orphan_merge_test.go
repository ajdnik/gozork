package engine

import (
	"testing"
)

func TestOrphanMergeCopiesObjectClause(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.OrphanedSyntx.Verb = LexItem{Norm: "take", Orig: "take", Types: WordTypes{WordVerb}}
	G.OrphanedSyntx.ObjOrClause1 = []LexItem{}

	G.ParsedSyntx.Verb = LexItem{Norm: "take", Orig: "take", Types: WordTypes{WordVerb}}
	G.ParsedSyntx.ObjOrClause1 = []LexItem{{Norm: "lamp", Orig: "lamp", Types: WordTypes{WordObj}}}
	G.ParsedSyntx.Obj1Start = 0
	G.ParsedSyntx.Obj1End = 1
	G.Params.ObjOrClauseCnt = 1
	G.LexRes = []LexItem{{Norm: "lamp", Orig: "lamp", Types: WordTypes{WordObj}}}

	OrphanMerge()

	if !G.Params.HasMerged {
		t.Fatalf("expected HasMerged to be true")
	}
	if len(G.OrphanedSyntx.ObjOrClause1) != 1 || G.OrphanedSyntx.ObjOrClause1[0].Norm != "lamp" {
		t.Fatalf("expected orphaned clause to contain lamp")
	}
	if len(G.ParsedSyntx.ObjOrClause1) != 1 || G.ParsedSyntx.ObjOrClause1[0].Norm != "lamp" {
		t.Fatalf("expected merged clause to contain lamp")
	}
}

// ---- orphan_merge_adj_test.go ----

func TestOrphanMergeAdjectiveClause(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.OrphanedSyntx.Verb = LexItem{Norm: "take", Orig: "take", Types: WordTypes{WordVerb}}
	G.OrphanedSyntx.ObjOrClause1 = []LexItem{
		{Norm: "brass", Orig: "brass", Types: WordTypes{WordAdj}},
		{Norm: "lamp", Orig: "lamp", Types: WordTypes{WordObj}},
	}

	G.Params.AdjClause.Type = Clause1
	G.Params.AdjClause.Syn = LexItem{Norm: "lamp", Orig: "lamp", Types: WordTypes{WordObj}}
	G.Params.AdjClause.Adj = LexItem{Norm: "brass", Orig: "brass", Types: WordTypes{WordAdj}}
	G.Params.ObjOrClauseCnt = 1

	G.ParsedSyntx.Verb = LexItem{Norm: "brass", Orig: "brass", Types: WordTypes{WordAdj}}
	G.ParsedSyntx.Obj1Start = 0
	G.ParsedSyntx.Obj1End = 1
	G.LexRes = []LexItem{{Norm: "brass", Orig: "brass", Types: WordTypes{WordAdj}}}

	OrphanMerge()

	if !G.Params.HasMerged {
		t.Fatalf("expected HasMerged to be true")
	}
	if len(G.OrphanedSyntx.ObjOrClause1) != 3 || G.OrphanedSyntx.ObjOrClause1[0].Norm != "brass" {
		t.Fatalf("expected adjective to be inserted into orphaned clause")
	}
}

// ---- orphan_merge_nclause_test.go ----

func TestOrphanMergeNounClause(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.OrphanedSyntx.Verb = LexItem{Norm: "take", Orig: "take", Types: WordTypes{WordVerb}}
	G.OrphanedSyntx.ObjOrClause1 = []LexItem{
		{Norm: "lamp", Orig: "lamp", Types: WordTypes{WordObj}},
	}

	G.Params.AdjClause.Type = Clause1
	G.Params.AdjClause.Syn = LexItem{Norm: "lamp", Orig: "lamp", Types: WordTypes{WordObj}}
	G.Params.AdjClause.Adj = LexItem{Norm: "brass", Orig: "brass", Types: WordTypes{WordAdj}}
	G.Params.ObjOrClauseCnt = 1

	G.ParsedSyntx.Verb = LexItem{Norm: "take", Orig: "take", Types: WordTypes{WordVerb}}
	G.ParsedSyntx.Obj1Start = 0
	G.ParsedSyntx.Obj1End = 1
	G.ParsedSyntx.ObjOrClause1 = []LexItem{
		{Norm: "sword", Orig: "sword", Types: WordTypes{WordObj}},
	}
	G.LexRes = []LexItem{{Norm: "sword", Orig: "sword", Types: WordTypes{WordObj}}}

	OrphanMerge()

	if !G.Params.HasMerged {
		t.Fatalf("expected HasMerged to be true")
	}
	if len(G.OrphanedSyntx.ObjOrClause1) != 1 || G.OrphanedSyntx.ObjOrClause1[0].Norm != "sword" {
		t.Fatalf("expected orphaned clause to be replaced with sword")
	}
}
