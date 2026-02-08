package engine

import (
	"testing"
)

func TestGetObjectGetOneSelectsOne(t *testing.T) {
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

	a := &Object{Desc: "rock", Synonyms: []string{"rock"}, In: room}
	b := &Object{Desc: "rock", Synonyms: []string{"rock"}, In: room}
	room.AddChild(a)
	room.AddChild(b)

	G.Search.Syn.Norm = "rock"
	G.Search.Syn.Orig = "rock"
	G.Search.Syn.Types = WordTypes{WordObj}
	G.Search.LocFlags = LocSet(LocOnGrnd)
	G.Params.GetType = GetOne

	res := GetObject(true, true)
	if len(res) != 1 {
		t.Fatalf("expected GetOne to return a single object, got %d", len(res))
	}
	if res[0] != a && res[0] != b {
		t.Fatalf("expected GetOne to return one of the rocks")
	}
}
