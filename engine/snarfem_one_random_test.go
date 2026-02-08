package engine

import (
	"testing"
)

func TestSnarfemOneWithMultipleMatches(t *testing.T) {
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

	a := &Object{Desc: "lamp", Synonyms: []string{"lamp"}, Adjectives: []string{"brass"}, In: room}
	b := &Object{Desc: "lamp", Synonyms: []string{"lamp"}, Adjectives: []string{"brass"}, In: room}
	room.AddChild(a)
	room.AddChild(b)

	G.Search.LocFlags = LocSet(LocOnGrnd)
	G.Params.GetType = GetUndef

	wrds := []LexItem{
		{Norm: "one", Orig: "one", Types: WordTypes{WordObj}},
		{Norm: "brass", Orig: "brass", Types: WordTypes{WordAdj}},
		{Norm: "lamp", Orig: "lamp", Types: WordTypes{WordObj}},
	}

	res := Snarfem(true, wrds)
	if len(res) != 1 {
		t.Fatalf("expected Snarfem to return one object, got %d", len(res))
	}
	if res[0] != a && res[0] != b {
		t.Fatalf("expected Snarfem to return one of the lamps")
	}
}
