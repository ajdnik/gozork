package engine

import (
	"testing"
)

func TestGlobalCheck(t *testing.T) {
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

	G.GlobalObj = &Object{Desc: "global"}
	G.LocalGlobalObj = &Object{Desc: "localglobal"}

	globalItem := &Object{Desc: "rope", Synonyms: []string{"rope"}, In: G.GlobalObj}
	localItem := &Object{Desc: "rope", Synonyms: []string{"rope"}, In: G.LocalGlobalObj}
	G.GlobalObj.AddChild(globalItem)
	G.LocalGlobalObj.AddChild(localItem)
	room.Global = []*Object{localItem}

	G.Search.Syn = LexItem{Norm: "rope", Orig: "rope", Types: WordTypes{WordObj}}
	G.Search.ObjFlags = 0
	G.Search.LocFlags = LocSet(LocInRoom)

	res := GlobalCheck()
	if len(res) != 1 || res[0] != localItem {
		t.Fatalf("expected GlobalCheck to return room global item, got %d", len(res))
	}

	room.Global = nil
	res = GlobalCheck()
	if len(res) != 1 || res[0] != globalItem {
		t.Fatalf("expected GlobalCheck to fall back to global list, got %d", len(res))
	}
}

func TestGlobalCheckPseudoAndRoomsFallback(t *testing.T) {
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
	G.GlobalObj = &Object{Desc: "global"}
	G.PseudoObj = &Object{Desc: "pseudo"}

	room.Pseudo = []PseudoObj{
		{
			Synonym: "spirit",
			Action:  func(ActionArg) bool { return true },
		},
	}

	G.Search.Syn = LexItem{Norm: "spirit", Orig: "spirit", Types: WordTypes{WordObj}}
	G.Search.ObjFlags = 0
	G.Search.LocFlags = LocSet(LocInRoom)

	res := GlobalCheck()
	if len(res) != 1 || res[0] != G.PseudoObj {
		t.Fatalf("expected GlobalCheck to return pseudo object, got %+v", res)
	}
	if G.PseudoObj.Action == nil {
		t.Fatalf("expected pseudo action to be attached")
	}

	room.Pseudo = nil
	G.ActVerb = ActionVerb{Norm: "search", Orig: "search"}
	G.Search.LocFlags = LocSet(LocHave)
	G.Search.ObjFlags = FlgTake
	G.Search.Syn.Clear()

	roomsChild := &Object{Desc: "roomchild", Synonyms: []string{"roomchild"}, Flags: FlgTake, In: G.RoomsObj}
	G.RoomsObj.AddChild(roomsChild)

	res = GlobalCheck()
	if len(res) != 1 || res[0] != roomsChild {
		t.Fatalf("expected GlobalCheck to search rooms list, got %+v", res)
	}
}
