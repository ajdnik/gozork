package engine

import (
	"testing"
)

func TestSnarfemAllAndExcept(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	room := &Object{Desc: "room"}
	player := &Object{Desc: "player", In: room}
	G.Player = player
	G.Winner = player
	G.Here = room
	G.Lit = true
	G.RoomsObj = &Object{Desc: "rooms"}
	room.In = G.RoomsObj

	a := &Object{Desc: "lamp", Synonyms: []string{"lamp"}, In: room}
	b := &Object{Desc: "key", Synonyms: []string{"key"}, In: room}
	room.AddChild(a)
	room.AddChild(b)

	G.Search.LocFlags = LocSet(LocOnGrnd)
	G.Params.GetType = GetUndef
	G.Params.Buts = nil

	wrds := []LexItem{
		{Norm: "all", Orig: "all", Types: WordTypes{WordObj}},
		{Norm: "except", Orig: "except", Types: nil},
		{Norm: "key", Orig: "key", Types: WordTypes{WordObj}},
	}

	res := Snarfem(true, wrds)
	if len(res) != 2 {
		t.Fatalf("expected Snarfem to return all objects, got %d", len(res))
	}
	if len(G.Params.Buts) != 1 || G.Params.Buts[0] != b {
		t.Fatalf("expected Snarfem to record key as exclusion")
	}
}

func TestSnarfemAndCombinesObjects(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	room := &Object{Desc: "room"}
	player := &Object{Desc: "player", In: room}
	G.Player = player
	G.Winner = player
	G.Here = room
	G.Lit = true
	G.RoomsObj = &Object{Desc: "rooms"}
	room.In = G.RoomsObj

	a := &Object{Desc: "lamp", Synonyms: []string{"lamp"}, In: room}
	b := &Object{Desc: "key", Synonyms: []string{"key"}, In: room}
	room.AddChild(a)
	room.AddChild(b)

	G.Search.LocFlags = LocSet(LocOnGrnd)
	G.Params.GetType = GetUndef

	wrds := []LexItem{
		{Norm: "lamp", Orig: "lamp", Types: WordTypes{WordObj}},
		{Norm: "and", Orig: "and", Types: nil},
		{Norm: "key", Orig: "key", Types: WordTypes{WordObj}},
	}

	res := Snarfem(true, wrds)
	if len(res) != 2 {
		t.Fatalf("expected Snarfem to return both objects, got %d", len(res))
	}
	if !G.Params.HasAnd {
		t.Fatalf("expected HasAnd to be set when using and")
	}
}

func TestSnarfemOfInhibitsGetType(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.Params.GetType = GetUndef
	G.Search.LocFlags = LocSet(LocOnGrnd)

	wrds := []LexItem{
		{Norm: "of", Orig: "of", Types: nil},
	}

	res := Snarfem(true, wrds)
	if res == nil {
		t.Fatalf("expected Snarfem to return empty slice, got nil")
	}
	if len(res) != 0 {
		t.Fatalf("expected Snarfem to return no objects, got %d", len(res))
	}
	if G.Params.GetType != GetInhibit {
		t.Fatalf("expected GetType to be set to GetInhibit")
	}
}

func TestSnarfemOneOfKeepsGetOne(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.Params.GetType = GetUndef
	G.Search.LocFlags = LocSet(LocOnGrnd)

	wrds := []LexItem{
		{Norm: "one", Orig: "one", Types: WordTypes{WordObj}},
		{Norm: "of", Orig: "of", Types: nil},
	}

	res := Snarfem(true, wrds)
	if res != nil {
		t.Fatalf("expected Snarfem to return nil when noun is missing")
	}
	if G.Params.GetType != GetOne {
		t.Fatalf("expected GetType to remain GetOne")
	}
}

func TestSnarfemAdjectiveOfObject(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	room := &Object{Desc: "room"}
	player := &Object{Desc: "player", In: room}
	G.Player = player
	G.Winner = player
	G.Here = room
	G.Lit = true
	G.RoomsObj = &Object{Desc: "rooms"}
	room.In = G.RoomsObj

	lamp := &Object{Desc: "lamp", Synonyms: []string{"lamp"}, Adjectives: []string{"brass"}, In: room}
	room.AddChild(lamp)

	G.Search.LocFlags = LocSet(LocOnGrnd)
	G.Params.GetType = GetUndef

	wrds := []LexItem{
		{Norm: "brass", Orig: "brass", Types: WordTypes{WordAdj}},
		{Norm: "lamp", Orig: "lamp", Types: WordTypes{WordObj}},
	}

	res := Snarfem(true, wrds)
	if len(res) != 1 || res[0] != lamp {
		t.Fatalf("expected Snarfem to find brass lamp, got %d", len(res))
	}
}

func TestSnarfemOneWithAdjective(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	room := &Object{Desc: "room"}
	player := &Object{Desc: "player", In: room}
	G.Player = player
	G.Winner = player
	G.Here = room
	G.Lit = true
	G.RoomsObj = &Object{Desc: "rooms"}
	room.In = G.RoomsObj

	lamp := &Object{Desc: "lamp", Synonyms: []string{"lamp"}, Adjectives: []string{"brass"}, In: room}
	room.AddChild(lamp)

	G.Search.LocFlags = LocSet(LocOnGrnd)
	G.Params.GetType = GetUndef

	wrds := []LexItem{
		{Norm: "one", Orig: "one", Types: WordTypes{WordObj}},
		{Norm: "brass", Orig: "brass", Types: WordTypes{WordAdj}},
		{Norm: "lamp", Orig: "lamp", Types: WordTypes{WordObj}},
	}

	res := Snarfem(true, wrds)
	if len(res) != 1 || res[0] != lamp {
		t.Fatalf("expected Snarfem to return brass lamp, got %d", len(res))
	}
}
