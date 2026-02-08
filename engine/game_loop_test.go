package engine

import (
	"bytes"
	"testing"
)

func TestMainLoopQuit(t *testing.T) {
	oldG := G
	oldVocab := Vocabulary
	oldCommands := Commands
	G = NewGameState()
	Vocabulary = make(map[string]WordItem)
	t.Cleanup(func() {
		G = oldG
		Vocabulary = oldVocab
		Commands = oldCommands
	})

	room := &Object{Desc: "room", Flags: FlgOn}
	player := &Object{Desc: "player", In: room}
	G.Player = player
	G.Winner = player
	G.Here = room
	G.ItPronounObj = &Object{Desc: "it"}
	G.NotHereObj = &Object{Desc: "not here"}
	G.AllObjects = []*Object{room, player}

	G.GameInput = bytes.NewBufferString("quit\n")
	G.GameOutput = &bytes.Buffer{}

	commands := []Syntax{
		{
			Verb: "quit",
			Action: func(ActionArg) bool {
				Quit()
				return true
			},
		},
	}
	BuildVocabulary(commands, nil, nil)

	MainLoop()
	if !G.QuitRequested {
		t.Fatalf("expected MainLoop to set QuitRequested")
	}
}

func TestMainLoopContinue(t *testing.T) {
	oldG := G
	oldVocab := Vocabulary
	oldCommands := Commands
	G = NewGameState()
	Vocabulary = make(map[string]WordItem)
	t.Cleanup(func() {
		G = oldG
		Vocabulary = oldVocab
		Commands = oldCommands
	})

	room := &Object{Desc: "room", Flags: FlgOn}
	player := &Object{Desc: "player", In: room}
	G.Player = player
	G.Winner = player
	G.Here = room
	G.ItPronounObj = &Object{Desc: "it"}
	G.NotHereObj = &Object{Desc: "not here"}
	G.AllObjects = []*Object{room, player}

	G.GameInput = bytes.NewBufferString("look\nquit\n")
	var out bytes.Buffer
	G.GameOutput = &out

	commands := []Syntax{
		{
			Verb: "look",
			Action: func(ActionArg) bool {
				return true
			},
		},
		{
			Verb: "quit",
			Action: func(ActionArg) bool {
				Quit()
				return true
			},
		},
	}
	BuildVocabulary(commands, nil, nil)

	MainLoop()
	if !G.QuitRequested {
		t.Fatalf("expected MainLoop to handle continue and quit, output=%q", out.String())
	}
}

func TestMainLoopDarkRoomNoObjects(t *testing.T) {
	oldG := G
	oldVocab := Vocabulary
	oldCommands := Commands
	G = NewGameState()
	Vocabulary = make(map[string]WordItem)
	t.Cleanup(func() {
		G = oldG
		Vocabulary = oldVocab
		Commands = oldCommands
	})

	room := &Object{Desc: "room", Flags: FlgLand}
	player := &Object{Desc: "player", In: room}
	G.Player = player
	G.Winner = player
	G.Here = room
	G.ItPronounObj = &Object{Desc: "it"}
	G.NotHereObj = &Object{Desc: "not here"}
	G.RoomsObj = &Object{Desc: "rooms"}
	room.In = G.RoomsObj
	G.AllObjects = []*Object{room, player}
	G.AlwaysLit = false

	G.GameInput = bytes.NewBufferString("take all\nquit\n")
	var out bytes.Buffer
	G.GameOutput = &out

	commands := []Syntax{
		{
			Verb: "take",
			Obj1: ObjProp{HasObj: true, ObjFlags: FlgTake, LocFlags: LocSet(LocOnGrnd, LocMany)},
			Action: func(ActionArg) bool {
				return true
			},
		},
		{
			Verb: "quit",
			Action: func(ActionArg) bool {
				Quit()
				return true
			},
		},
	}
	BuildVocabulary(commands, []string{"all"}, nil)

	MainLoop()
	if !G.QuitRequested {
		t.Fatalf("expected MainLoop to quit")
	}
	if !bytes.Contains(out.Bytes(), []byte("It's too dark to see.")) {
		t.Fatalf("expected dark-room message, output=%q", out.String())
	}
}

func TestMainLoopNotClearWhenNoObjects(t *testing.T) {
	oldG := G
	oldVocab := Vocabulary
	oldCommands := Commands
	G = NewGameState()
	Vocabulary = make(map[string]WordItem)
	t.Cleanup(func() {
		G = oldG
		Vocabulary = oldVocab
		Commands = oldCommands
	})

	room := &Object{Desc: "room", Flags: FlgLand}
	player := &Object{Desc: "player", In: room}
	G.Player = player
	G.Winner = player
	G.Here = room
	G.ItPronounObj = &Object{Desc: "it"}
	G.NotHereObj = &Object{Desc: "not here"}
	G.RoomsObj = &Object{Desc: "rooms"}
	room.In = G.RoomsObj
	G.AllObjects = []*Object{room, player}
	G.AlwaysLit = true

	G.GameInput = bytes.NewBufferString("take all\nquit\n")
	var out bytes.Buffer
	G.GameOutput = &out

	commands := []Syntax{
		{
			Verb: "take",
			Obj1: ObjProp{HasObj: true, ObjFlags: FlgTake, LocFlags: LocSet(LocOnGrnd)},
			Action: func(ActionArg) bool {
				return true
			},
		},
		{
			Verb: "quit",
			Action: func(ActionArg) bool {
				Quit()
				return true
			},
		},
	}
	BuildVocabulary(commands, []string{"all"}, nil)

	MainLoop()
	if !G.QuitRequested {
		t.Fatalf("expected MainLoop to quit")
	}
	if !bytes.Contains(out.Bytes(), []byte("It's not clear what you're referring to.")) {
		t.Fatalf("expected not-clear message, output=%q", out.String())
	}
}

func TestMainLoopMultipleObjects(t *testing.T) {
	oldG := G
	oldVocab := Vocabulary
	oldCommands := Commands
	G = NewGameState()
	Vocabulary = make(map[string]WordItem)
	t.Cleanup(func() {
		G = oldG
		Vocabulary = oldVocab
		Commands = oldCommands
	})

	room := &Object{Desc: "room", Flags: FlgLand}
	player := &Object{Desc: "player", In: room}
	lamp := &Object{Desc: "lamp", Synonyms: []string{"lamp"}, Flags: FlgTake, In: room}
	key := &Object{Desc: "key", Synonyms: []string{"key"}, Flags: FlgTake, In: room}
	room.AddChild(lamp)
	room.AddChild(key)

	G.Player = player
	G.Winner = player
	G.Here = room
	G.ItPronounObj = &Object{Desc: "it"}
	G.NotHereObj = &Object{Desc: "not here"}
	G.RoomsObj = &Object{Desc: "rooms"}
	room.In = G.RoomsObj
	G.AllObjects = []*Object{room, player, lamp, key}
	G.AlwaysLit = true

	G.GameInput = bytes.NewBufferString("take lamp and key\nquit\n")
	var out bytes.Buffer
	G.GameOutput = &out

	commands := []Syntax{
		{
			Verb: "take",
			Obj1: ObjProp{HasObj: true, ObjFlags: FlgTake, LocFlags: LocSet(LocOnGrnd, LocMany)},
			Action: func(ActionArg) bool {
				return true
			},
		},
		{
			Verb: "quit",
			Action: func(ActionArg) bool {
				Quit()
				return true
			},
		},
	}
	BuildVocabulary(commands, []string{"and"}, nil)

	MainLoop()
	if !G.QuitRequested {
		t.Fatalf("expected MainLoop to quit")
	}
	if !bytes.Contains(out.Bytes(), []byte("lamp:")) || !bytes.Contains(out.Bytes(), []byte("key:")) {
		t.Fatalf("expected object prefixes for multiple objects, output=%q", out.String())
	}
}

func TestMainLoopCallsRoomActionAtEnd(t *testing.T) {
	oldG := G
	oldVocab := Vocabulary
	oldCommands := Commands
	G = NewGameState()
	Vocabulary = make(map[string]WordItem)
	t.Cleanup(func() {
		G = oldG
		Vocabulary = oldVocab
		Commands = oldCommands
	})

	room := &Object{Desc: "room", Flags: FlgLand}
	player := &Object{Desc: "player", In: room}
	G.Player = player
	G.Winner = player
	G.Here = room
	G.ItPronounObj = &Object{Desc: "it"}
	G.NotHereObj = &Object{Desc: "not here"}
	G.RoomsObj = &Object{Desc: "rooms"}
	room.In = G.RoomsObj
	G.AllObjects = []*Object{room, player}
	G.AlwaysLit = true

	calls := []ActionArg{}
	room.Action = func(arg ActionArg) bool {
		calls = append(calls, arg)
		return arg == ActEnd
	}

	G.GameInput = bytes.NewBufferString("look\nquit\n")
	G.GameOutput = &bytes.Buffer{}

	commands := []Syntax{
		{
			Verb:   "look",
			Action: func(ActionArg) bool { return true },
		},
		{
			Verb: "quit",
			Action: func(ActionArg) bool {
				Quit()
				return true
			},
		},
	}
	BuildVocabulary(commands, nil, nil)

	MainLoop()
	if !G.QuitRequested {
		t.Fatalf("expected MainLoop to quit")
	}
	if len(calls) == 0 {
		t.Fatalf("expected room action to be called")
	}
	foundEnd := false
	for _, c := range calls {
		if c == ActEnd {
			foundEnd = true
			break
		}
	}
	if !foundEnd {
		t.Fatalf("expected ActEnd to be passed to room action, got %+v", calls)
	}
}

func TestMainLoopWalkDirection(t *testing.T) {
	oldG := G
	oldVocab := Vocabulary
	oldCommands := Commands
	G = NewGameState()
	Vocabulary = make(map[string]WordItem)
	t.Cleanup(func() {
		G = oldG
		Vocabulary = oldVocab
		Commands = oldCommands
	})

	room := &Object{Desc: "room", Flags: FlgLand}
	player := &Object{Desc: "player", In: room}
	G.Player = player
	G.Winner = player
	G.Here = room
	G.ItPronounObj = &Object{Desc: "it"}
	G.NotHereObj = &Object{Desc: "not here"}
	G.RoomsObj = &Object{Desc: "rooms"}
	room.In = G.RoomsObj
	G.AllObjects = []*Object{room, player}

	called := false
	commands := []Syntax{
		{
			Verb: "walk",
			Action: func(ActionArg) bool {
				called = true
				return true
			},
		},
		{
			Verb: "quit",
			Action: func(ActionArg) bool {
				Quit()
				return true
			},
		},
	}

	G.GameInput = bytes.NewBufferString("north\nquit\n")
	G.GameOutput = &bytes.Buffer{}
	BuildVocabulary(commands, nil, nil)

	MainLoop()
	if !G.QuitRequested {
		t.Fatalf("expected MainLoop to quit")
	}
	if !called {
		t.Fatalf("expected walk action to be called")
	}
}
