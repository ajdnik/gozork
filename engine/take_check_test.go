package engine

import (
	"bytes"
	"strings"
	"testing"
)

func TestITakeCheckRequiresHave(t *testing.T) {
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

	ok := ITakeCheck([]*Object{obj}, LocSet(LocHave))
	if ok {
		t.Fatalf("expected ITakeCheck to fail when object is not held")
	}
	if !strings.Contains(out.String(), "You don't have the lamp") {
		t.Fatalf("expected missing-have message, got %q", out.String())
	}
}

func TestITakeCheckNotHere(t *testing.T) {
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

	var out bytes.Buffer
	G.GameOutput = &out

	ok := ITakeCheck([]*Object{G.NotHereObj}, LocSet(LocHave))
	if ok {
		t.Fatalf("expected ITakeCheck to fail for NotHereObj with LocHave")
	}
	if !strings.Contains(out.String(), "You don't have that") {
		t.Fatalf("expected NotHereObj message, got %q", out.String())
	}
}

func TestITakeCheckLocHaveAndLocTakeNotHere(t *testing.T) {
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

	var out bytes.Buffer
	G.GameOutput = &out

	ok := ITakeCheck([]*Object{G.NotHereObj}, LocSet(LocHave, LocTake))
	if ok {
		t.Fatalf("expected ITakeCheck to fail for NotHereObj")
	}
	if !strings.Contains(out.String(), "You don't have that") {
		t.Fatalf("expected NotHereObj message, got %q", out.String())
	}
}
