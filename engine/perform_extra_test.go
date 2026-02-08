package engine

import (
	"testing"
)

func TestPerformNotHereAction(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	room := &Object{Desc: "room"}
	player := &Object{Desc: "player", In: room}
	it := &Object{Desc: "it", In: room}
	room.AddChild(it)
	G.RoomsObj = &Object{Desc: "rooms"}
	room.In = G.RoomsObj
	G.Player = player
	G.Winner = player
	G.Here = room
	G.ItPronounObj = &Object{Desc: "it-pronoun"}
	G.Params.ItObj = it

	called := false
	G.NotHereObj = &Object{Desc: "not here"}
	G.NotHereObj.Action = func(ActionArg) bool {
		called = true
		return true
	}

	res := Perform(ActionVerb{Norm: "test", Orig: "test"}, G.NotHereObj, nil)
	if res != PerfHndld {
		t.Fatalf("expected PerfHndld, got %v", res)
	}
	if !called {
		t.Fatalf("expected NotHere action to be called")
	}
}

func TestPerformPreActionAndActionNormFallback(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	room := &Object{Desc: "room"}
	player := &Object{Desc: "player", In: room}
	it := &Object{Desc: "it", In: room}
	room.AddChild(it)
	G.RoomsObj = &Object{Desc: "rooms"}
	room.In = G.RoomsObj
	G.Player = player
	G.Winner = player
	G.Here = room
	G.ItPronounObj = &Object{Desc: "it-pronoun"}
	G.Params.ItObj = it
	G.NotHereObj = &Object{Desc: "not here"}

	preCalled := false
	actCalled := false
	G.PreActions["look"] = func(ActionArg) bool {
		preCalled = true
		return false
	}
	G.Actions["look"] = func(ActionArg) bool {
		actCalled = true
		return true
	}

	res := Perform(ActionVerb{Norm: "look", Orig: "see"}, nil, nil)
	if res != PerfHndld {
		t.Fatalf("expected PerfHndld, got %v", res)
	}
	if !preCalled || !actCalled {
		t.Fatalf("expected both pre-action and action to run, pre=%v action=%v", preCalled, actCalled)
	}
}

func TestPerformContainerContFcn(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	room := &Object{Desc: "room"}
	player := &Object{Desc: "player", In: room}
	container := &Object{Desc: "box", In: room}
	obj := &Object{Desc: "coin", In: container}
	it := &Object{Desc: "it", In: room}
	room.AddChild(it)
	G.RoomsObj = &Object{Desc: "rooms"}
	room.In = G.RoomsObj
	G.Player = player
	G.Winner = player
	G.Here = room
	G.ItPronounObj = &Object{Desc: "it-pronoun"}
	G.Params.ItObj = it
	G.NotHereObj = &Object{Desc: "not here"}

	called := false
	container.ContFcn = func(ActionArg) bool {
		called = true
		return true
	}

	res := Perform(ActionVerb{Norm: "take", Orig: "take"}, obj, nil)
	if res != PerfHndld {
		t.Fatalf("expected PerfHndld, got %v", res)
	}
	if !called {
		t.Fatalf("expected container ContFcn to be called")
	}
}
