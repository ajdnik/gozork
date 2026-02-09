package engine

import (
	"bytes"
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

func TestResetGameStatePreservesRefs(t *testing.T) {
	oldG := G
	G = NewGameState()
	defer func() { G = oldG }()

	out := &bytes.Buffer{}
	in := strings.NewReader("x")
	rng := rand.New(rand.NewSource(1))
	clockFuncs := map[string]func() bool{"noop": func() bool { return true }}
	objs := []*Object{{Desc: "dummy"}}
	rooms := &Object{Desc: "rooms"}
	global := &Object{Desc: "global"}
	local := &Object{Desc: "local"}
	notHere := &Object{Desc: "notHere"}
	pseudo := &Object{Desc: "pseudo"}
	itObj := &Object{Desc: "it"}
	meObj := &Object{Desc: "me"}
	handsObj := &Object{Desc: "hands"}
	data := &struct{ X int }{X: 1}

	G.GameOutput = out
	G.GameInput = in
	G.Rand = rng
	G.ClockFuncs = clockFuncs
	G.AllObjects = objs
	G.RoomsObj = rooms
	G.GlobalObj = global
	G.LocalGlobalObj = local
	G.NotHereObj = notHere
	G.PseudoObj = pseudo
	G.ItPronounObj = itObj
	G.MeObj = meObj
	G.HandsObj = handsObj
	G.GameData = data

	ResetGameState()

	if G.GameOutput != out {
		t.Fatalf("expected GameOutput preserved")
	}
	if G.GameInput != in {
		t.Fatalf("expected GameInput preserved")
	}
	if G.Rand != rng {
		t.Fatalf("expected Rand preserved")
	}
	if reflect.ValueOf(G.ClockFuncs).Pointer() != reflect.ValueOf(clockFuncs).Pointer() {
		t.Fatalf("expected ClockFuncs preserved")
	}
	if len(G.AllObjects) != len(objs) || (len(objs) > 0 && G.AllObjects[0] != objs[0]) {
		t.Fatalf("expected AllObjects preserved")
	}
	if G.RoomsObj != rooms {
		t.Fatalf("expected RoomsObj preserved")
	}
	if G.GlobalObj != global {
		t.Fatalf("expected GlobalObj preserved")
	}
	if G.LocalGlobalObj != local {
		t.Fatalf("expected LocalGlobalObj preserved")
	}
	if G.NotHereObj != notHere {
		t.Fatalf("expected NotHereObj preserved")
	}
	if G.PseudoObj != pseudo {
		t.Fatalf("expected PseudoObj preserved")
	}
	if G.ItPronounObj != itObj {
		t.Fatalf("expected ItPronounObj preserved")
	}
	if G.MeObj != meObj {
		t.Fatalf("expected MeObj preserved")
	}
	if G.HandsObj != handsObj {
		t.Fatalf("expected HandsObj preserved")
	}
	if G.GameData != data {
		t.Fatalf("expected GameData preserved")
	}
}
