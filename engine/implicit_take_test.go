package engine

import (
	"bytes"
	"strings"
	"testing"
)

func TestITakeCheckImplicitTake(t *testing.T) {
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
	G.NotHereObj = &Object{Desc: "not here"}

	obj := &Object{Desc: "lamp", Synonyms: []string{"lamp"}, In: room}
	room.AddChild(obj)

	var out bytes.Buffer
	G.GameOutput = &out

	called := false
	G.ITakeFunc = func(vb bool) bool {
		called = true
		return true
	}

	ok := ITakeCheck([]*Object{obj}, LocSet(LocTake))
	if !ok {
		t.Fatalf("expected ITakeCheck to succeed when implicit take occurs")
	}
	if !called {
		t.Fatalf("expected ITakeFunc to be called")
	}
	if !strings.Contains(out.String(), "(Taken)") {
		t.Fatalf("expected implicit take message, got %q", out.String())
	}
}
