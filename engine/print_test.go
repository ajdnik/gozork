package engine

import (
	"bytes"
	"strings"
	"testing"
)

func TestThingPrintUsesTheAndOriginal(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	var out bytes.Buffer
	G.GameOutput = &out

	G.ParsedSyntx.ObjOrClause1 = []LexItem{
		{Norm: "lamp", Orig: "lamp", Types: WordTypes{WordObj}},
		{Norm: ".", Orig: ".", Types: nil},
	}
	G.Params.ShldOrphan = false
	G.Params.HasMerged = false

	ThingPrint(true, true)
	if got := strings.TrimSpace(out.String()); got != "the lamp" {
		t.Fatalf("expected ThingPrint to include article, got %q", got)
	}
}

func TestWhichPrintListsOptions(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	var out bytes.Buffer
	G.GameOutput = &out

	G.ParsedSyntx.ObjOrClause1 = []LexItem{
		{Norm: "lamp", Orig: "lamp", Types: WordTypes{WordObj}},
	}

	obj1 := &Object{Desc: "brass lamp"}
	obj2 := &Object{Desc: "silver lamp"}

	WhichPrint(true, []*Object{obj1, obj2})
	got := out.String()
	if !strings.Contains(got, "Which") || !strings.Contains(got, "brass lamp") || !strings.Contains(got, "silver lamp") {
		t.Fatalf("unexpected WhichPrint output: %q", got)
	}
}

func TestWhichPrintIndirectUsesClause2(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	var out bytes.Buffer
	G.GameOutput = &out

	G.ParsedSyntx.ObjOrClause2 = []LexItem{
		{Norm: "key", Orig: "key", Types: WordTypes{WordObj}},
	}

	obj1 := &Object{Desc: "small key"}
	obj2 := &Object{Desc: "large key"}

	WhichPrint(false, []*Object{obj1, obj2})
	got := out.String()
	if !strings.Contains(got, "small key") || !strings.Contains(got, "large key") {
		t.Fatalf("unexpected WhichPrint output: %q", got)
	}
}

func TestWhichPrintMultipleOptions(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	var out bytes.Buffer
	G.GameOutput = &out

	G.ParsedSyntx.ObjOrClause1 = []LexItem{
		{Norm: "key", Orig: "key", Types: WordTypes{WordObj}},
	}

	objs := []*Object{
		{Desc: "brass key"},
		{Desc: "silver key"},
		{Desc: "iron key"},
	}

	WhichPrint(true, objs)
	got := out.String()
	if !strings.Contains(got, "brass key") || !strings.Contains(got, "silver key") || !strings.Contains(got, "iron key") {
		t.Fatalf("expected all options to be listed, got %q", got)
	}
}

func TestThingPrintPronounsAndNumbers(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	var out bytes.Buffer
	G.GameOutput = &out
	G.MeObj = &Object{Desc: "myself"}
	G.Params.Number = 17

	G.ParsedSyntx.ObjOrClause1 = []LexItem{
		{Norm: "me", Orig: "me", Types: WordTypes{WordObj}},
		{Norm: "intnum", Orig: "17", Types: nil},
	}

	ThingPrint(true, false)
	got := strings.TrimSpace(out.String())
	if got != "myself 17" {
		t.Fatalf("expected pronoun and number output, got %q", got)
	}
}

func TestThingPrintHandlesItPronoun(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	var out bytes.Buffer
	G.GameOutput = &out
	room := &Object{Desc: "room"}
	G.RoomsObj = &Object{Desc: "rooms"}
	room.In = G.RoomsObj
	G.Here = room
	G.Winner = &Object{Desc: "player", In: room}
	G.Params.ItObj = &Object{Desc: "golden idol", In: room}

	G.ParsedSyntx.ObjOrClause1 = []LexItem{
		{Norm: "it", Orig: "it", Types: WordTypes{WordObj}},
	}

	ThingPrint(true, false)
	got := strings.TrimSpace(out.String())
	if got != "golden idol" {
		t.Fatalf("expected it pronoun to resolve, got %q", got)
	}
}

func TestThingPrintSpacingWithPunctuation(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	var out bytes.Buffer
	G.GameOutput = &out

	G.ParsedSyntx.ObjOrClause1 = []LexItem{
		{Norm: "lamp", Orig: "lamp", Types: WordTypes{WordObj}},
		{Norm: ",", Orig: ",", Types: nil},
		{Norm: "shiny", Orig: "shiny", Types: WordTypes{WordAdj}},
		{Norm: ".", Orig: ".", Types: nil},
	}

	ThingPrint(true, false)
	got := strings.TrimSpace(out.String())
	if got != "lamp, shiny" {
		t.Fatalf("expected punctuation spacing, got %q", got)
	}
}

func TestThingPrintWhenOrphanedUsesNorm(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	var out bytes.Buffer
	G.GameOutput = &out
	G.Params.ShldOrphan = true

	G.ParsedSyntx.ObjOrClause1 = []LexItem{
		{Norm: "lamp", Orig: "lamp", Types: WordTypes{WordObj}},
		{Norm: "shiny", Orig: "shiny", Types: WordTypes{WordAdj}},
	}

	ThingPrint(true, false)
	got := strings.TrimSpace(out.String())
	if got != "lamp shiny" {
		t.Fatalf("expected ThingPrint to use Norm when orphaned, got %q", got)
	}
}

func TestThingPrintIntnumWithPunctuation(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	var out bytes.Buffer
	G.GameOutput = &out
	G.Params.Number = 42

	G.ParsedSyntx.ObjOrClause1 = []LexItem{
		{Norm: "intnum", Orig: "42", Types: nil},
		{Norm: ".", Orig: ".", Types: nil},
	}

	ThingPrint(true, false)
	got := strings.TrimSpace(out.String())
	if got != "42" {
		t.Fatalf("expected intnum to print without punctuation, got %q", got)
	}
}
