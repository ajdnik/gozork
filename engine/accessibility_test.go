package engine

import (
	"testing"
)

func TestIsAccessible(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	room := &Object{Desc: "room"}
	player := &Object{Desc: "player", In: room}
	G.Player = player
	G.Winner = player
	G.Here = room
	G.RoomsObj = &Object{Desc: "rooms"}
	room.In = G.RoomsObj

	obj := &Object{Desc: "obj", In: room}
	if !IsAccessible(obj) {
		t.Fatalf("expected object in room to be accessible")
	}

	obj.Flags = FlgInvis
	if IsAccessible(obj) {
		t.Fatalf("expected invisible object to be inaccessible")
	}
	obj.Flags = 0

	global := &Object{Desc: "global"}
	G.GlobalObj = global
	obj.In = global
	if !IsAccessible(obj) {
		t.Fatalf("expected global object to be accessible")
	}

	localGlobal := &Object{Desc: "localglobal"}
	G.LocalGlobalObj = localGlobal
	room.Global = []*Object{obj}
	obj.In = localGlobal
	if !IsAccessible(obj) {
		t.Fatalf("expected local-global object to be accessible in room")
	}
}
