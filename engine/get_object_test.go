package engine

import (
	"bytes"
	"strings"
	"testing"
)

func TestGetObjectFindsHeldAndRoomObjects(t *testing.T) {
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

	held := &Object{Desc: "coin", Synonyms: []string{"coin"}, In: player}
	roomObj := &Object{Desc: "coin", Synonyms: []string{"coin"}, In: room}
	player.AddChild(held)
	room.AddChild(roomObj)

	G.Search.Syn.Norm = "coin"
	G.Search.Syn.Orig = "coin"
	G.Search.Syn.Types = WordTypes{WordObj}
	G.Search.LocFlags = LocSet(LocHeld, LocOnGrnd)
	G.Params.GetType = GetAll

	res := GetObject(true, true)
	if len(res) != 2 {
		t.Fatalf("expected both held and room objects, got %d", len(res))
	}
}

func TestGetObjectGetAllReturnsAllMatches(t *testing.T) {
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
	G.Params.GetType = GetAll

	res := GetObject(true, true)
	if len(res) != 2 {
		t.Fatalf("expected GetAll to return both rocks, got %d", len(res))
	}
}

// ---- get_object_inhibit_test.go ----

func TestGetObjectInhibitReturnsEmpty(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.Params.GetType = GetInhibit
	G.Search.Syn.Norm = "lamp"
	G.Search.Syn.Orig = "lamp"
	G.Search.Syn.Types = WordTypes{WordObj}

	res := GetObject(true, true)
	if res == nil {
		t.Fatalf("expected empty slice, got nil")
	}
	if len(res) != 0 {
		t.Fatalf("expected empty slice, got %d", len(res))
	}
}

// ---- get_object_locflags_test.go ----

func TestGetObjectRestoresLocFlagsAfterGetAll(t *testing.T) {
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

	obj := &Object{Desc: "lamp", Synonyms: []string{"lamp"}, In: room}
	room.AddChild(obj)

	G.Search.Syn = LexItem{Norm: "lamp", Orig: "lamp", Types: WordTypes{WordObj}}
	G.Search.LocFlags = LocSet(LocOnGrnd)
	G.Params.GetType = GetAll

	before := G.Search.LocFlags
	_ = GetObject(true, true)
	if G.Search.LocFlags != before {
		t.Fatalf("expected LocFlags to be restored after GetAll")
	}
}

// ---- get_object_orphan_test.go ----

func TestGetObjectTriggersWhichPrintAndOrphan(t *testing.T) {
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

	a := &Object{Desc: "brass key", Synonyms: []string{"key"}, In: room}
	b := &Object{Desc: "silver key", Synonyms: []string{"key"}, In: room}
	room.AddChild(a)
	room.AddChild(b)

	var out bytes.Buffer
	G.GameOutput = &out

	G.Search.Syn = LexItem{Norm: "key", Orig: "key", Types: WordTypes{WordObj}}
	G.Search.LocFlags = LocSet(LocOnGrnd)
	G.Params.GetType = GetUndef

	res := GetObject(true, true)
	if res != nil {
		t.Fatalf("expected GetObject to return nil when ambiguous")
	}
	if !G.Params.ShldOrphan {
		t.Fatalf("expected ShldOrphan to be set")
	}
	if !strings.Contains(out.String(), "Which") {
		t.Fatalf("expected WhichPrint output, got %q", out.String())
	}
}
