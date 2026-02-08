package engine

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseTblSetClear(t *testing.T) {
	src := ParseTbl{
		Verb:         LexItem{Norm: "take", Orig: "take", Types: WordTypes{WordVerb}},
		Prep1:        LexItem{Norm: "with", Orig: "with", Types: WordTypes{WordPrep}},
		ObjOrClause1: []LexItem{{Norm: "lamp", Orig: "lamp", Types: WordTypes{WordObj}}},
		Obj1Start:    1,
		Obj1End:      2,
		Prep2:        LexItem{Norm: "using", Orig: "using", Types: WordTypes{WordPrep}},
		ObjOrClause2: []LexItem{{Norm: "sword", Orig: "sword", Types: WordTypes{WordObj}}},
	}

	var dst ParseTbl
	dst.Set(src)

	src.Verb.Norm = "drop"
	src.ObjOrClause1[0].Norm = "box"
	if dst.Verb.Norm != "take" {
		t.Fatalf("expected Set to copy verb data")
	}
	if dst.ObjOrClause1[0].Norm != "lamp" {
		t.Fatalf("expected Set to clone object clause")
	}
	if dst.Obj1Start != 1 || dst.Obj1End != 2 {
		t.Fatalf("expected Set to copy range indices")
	}
	if dst.Prep2.Norm != "using" {
		t.Fatalf("expected Set to copy second prep")
	}

	dst.Clear()
	if dst.Verb.IsSet() || dst.Prep1.IsSet() || dst.Prep2.IsSet() {
		t.Fatalf("expected Clear to reset lexical items")
	}
	if dst.ObjOrClause1 != nil || dst.ObjOrClause2 != nil {
		t.Fatalf("expected Clear to reset object clauses")
	}
	if dst.Obj1Start != NumUndef || dst.Obj1End != NumUndef {
		t.Fatalf("expected Clear to reset range indices")
	}
}

func TestClausePropsIsSet(t *testing.T) {
	var cp ClauseProps
	if cp.IsSet() {
		t.Fatalf("expected zero ClauseProps to be unset")
	}
	cp.Type = Clause1
	if !cp.IsSet() {
		t.Fatalf("expected ClauseProps to be set when type is Clause1")
	}
}

// ---- parser_clause_test.go ----

func TestMkBuzz(t *testing.T) {
	bz := mkBuzz("then")
	if bz.Norm != "then" || bz.Orig != "then" {
		t.Fatalf("expected mkBuzz to set Norm and Orig, got %+v", bz)
	}
	if !bz.Types.Has(WordBuzz) {
		t.Fatalf("expected mkBuzz to tag WordBuzz, got %+v", bz.Types)
	}
}

func TestClauseWithPrepAndObject(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	prep := LexItem{Norm: "with", Orig: "with", Types: WordTypes{WordPrep}}
	obj := LexItem{Norm: "lamp", Orig: "lamp", Types: WordTypes{WordObj}}
	G.LexRes = []LexItem{prep, obj}
	G.Params.ObjOrClauseCnt = 1
	G.Params.BufLen = 1

	ok, end := Clause(0, prep)
	if !ok || end != 1 {
		t.Fatalf("expected Clause to succeed and end at 1, ok=%v end=%d", ok, end)
	}
	if !G.ParsedSyntx.Prep1.Is("with") {
		t.Fatalf("expected Prep1 to be set to with")
	}
	if len(G.ParsedSyntx.ObjOrClause1) != 1 || !G.ParsedSyntx.ObjOrClause1[0].Is("lamp") {
		t.Fatalf("expected ObjOrClause1 to contain lamp, got %+v", G.ParsedSyntx.ObjOrClause1)
	}
}

func TestClausePrepOnlyEnds(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	prep := LexItem{Norm: "with", Orig: "with", Types: WordTypes{WordPrep}}
	G.LexRes = []LexItem{prep}
	G.Params.ObjOrClauseCnt = 1
	G.Params.BufLen = 0

	ok, end := Clause(0, prep)
	if !ok || end != -1 {
		t.Fatalf("expected Clause to succeed and end at -1, ok=%v end=%d", ok, end)
	}
	if G.Params.ObjOrClauseCnt != 0 {
		t.Fatalf("expected ObjOrClauseCnt to decrement, got %d", G.Params.ObjOrClauseCnt)
	}
	if !G.ParsedSyntx.Prep1.Is("with") {
		t.Fatalf("expected Prep1 to be set to with")
	}
}

func TestClauseStopsOnThen(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	lamp := LexItem{Norm: "lamp", Orig: "lamp", Types: WordTypes{WordObj}}
	and := LexItem{Norm: "and", Orig: "and", Types: WordTypes{WordBuzz}}
	key := LexItem{Norm: "key", Orig: "key", Types: WordTypes{WordObj}}
	then := LexItem{Norm: "then", Orig: "then", Types: WordTypes{WordBuzz}}
	G.LexRes = []LexItem{lamp, and, key, then}
	G.Params.ObjOrClauseCnt = 1
	G.Params.BufLen = 3

	ok, end := Clause(0, lamp)
	if !ok || end != 2 {
		t.Fatalf("expected Clause to stop before then, ok=%v end=%d", ok, end)
	}
	if len(G.ParsedSyntx.ObjOrClause1) != 3 {
		t.Fatalf("expected ObjOrClause1 to include lamp and key, got %+v", G.ParsedSyntx.ObjOrClause1)
	}
}

// ---- parser_errors_test.go ----

func TestParseUnknownWord(t *testing.T) {
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
	G.Lit = true
	G.AlwaysLit = true
	G.AllObjects = []*Object{room, player}
	G.Params.Continue = NumUndef

	G.GameInput = bytes.NewBufferString("blorf\n")
	var out bytes.Buffer
	G.GameOutput = &out

	BuildVocabulary(nil, nil, nil)

	if ok := Parse(); ok {
		t.Fatalf("expected Parse to fail for unknown word")
	}
	if !strings.Contains(out.String(), "I don't know the word") {
		t.Fatalf("expected unknown word message, got %q", out.String())
	}
}

func TestParseMissingNoun(t *testing.T) {
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
	G.AllObjects = []*Object{room, player}
	G.Params.Continue = NumUndef

	G.GameInput = bytes.NewBufferString("take\n")
	var out bytes.Buffer
	G.GameOutput = &out

	commands := []Syntax{
		{
			Verb: "take",
			Obj1: ObjProp{HasObj: true},
		},
	}
	BuildVocabulary(commands, nil, nil)

	if ok := Parse(); ok {
		t.Fatalf("expected Parse to fail for missing noun")
	}
	if !strings.Contains(out.String(), "What do you want to take") {
		t.Fatalf("expected orphan prompt for missing noun, got %q", out.String())
	}
}

func TestParseCantUseWord(t *testing.T) {
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
	G.AllObjects = []*Object{room, player}
	G.Params.Continue = NumUndef

	G.GameInput = bytes.NewBufferString("take north\n")
	var out bytes.Buffer
	G.GameOutput = &out

	commands := []Syntax{
		{
			Verb: "take",
			Obj1: ObjProp{HasObj: true},
		},
	}
	BuildVocabulary(commands, nil, nil)

	if ok := Parse(); ok {
		t.Fatalf("expected Parse to fail for invalid word usage")
	}
	if !strings.Contains(out.String(), "You used the word") {
		t.Fatalf("expected cant-use message, got %q", out.String())
	}
	if !strings.Contains(out.String(), "north") {
		t.Fatalf("expected cant-use message to mention word, got %q", out.String())
	}
}

func TestParseOopsReplacesWord(t *testing.T) {
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
	G.AllObjects = []*Object{room, player}
	G.Params.Continue = NumUndef

	G.GameInput = bytes.NewBufferString("lok\noops look\n")
	var out bytes.Buffer
	G.GameOutput = &out

	commands := []Syntax{
		{
			Verb: "look",
		},
	}
	BuildVocabulary(commands, nil, nil)

	if ok := Parse(); ok {
		t.Fatalf("expected Parse to fail for misspelled verb")
	}
	if !strings.Contains(out.String(), "I don't know the word") {
		t.Fatalf("expected unknown word message, got %q", out.String())
	}

	out.Reset()
	if ok := Parse(); !ok {
		t.Fatalf("expected Parse to succeed after OOPS take, output=%q", out.String())
	}
	if G.ParsedSyntx.Verb.Norm != "look" {
		t.Fatalf("expected verb to be corrected to look, got %q", G.ParsedSyntx.Verb.Norm)
	}
}

// ---- parser_integration_test.go ----

func TestParseSimpleVerb(t *testing.T) {
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
	G.Lit = true
	G.AllObjects = []*Object{room, player}
	G.Params.Continue = NumUndef

	G.GameInput = bytes.NewBufferString("look\n")
	var out bytes.Buffer
	G.GameOutput = &out

	commands := []Syntax{
		{
			Verb:   "look",
			Action: func(ActionArg) bool { return true },
		},
	}
	BuildVocabulary(commands, nil, nil)

	if ok := Parse(); !ok {
		t.Fatalf("expected Parse to succeed for simple verb, output=%q", out.String())
	}
	if G.ActVerb.Norm != "look" {
		t.Fatalf("expected ActVerb.Norm to be look, got %q", G.ActVerb.Norm)
	}
	if G.DetectedSyntx == nil || G.DetectedSyntx.NumObjects() != 0 {
		t.Fatalf("expected syntax with zero objects to be detected")
	}
}

func TestParseDirection(t *testing.T) {
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
	G.Lit = true
	G.AllObjects = []*Object{room, player}
	G.Params.Continue = NumUndef

	G.GameInput = bytes.NewBufferString("north\n")
	var out bytes.Buffer
	G.GameOutput = &out

	BuildVocabulary(nil, nil, nil)

	if ok := Parse(); !ok {
		t.Fatalf("expected Parse to succeed for direction, output=%q", out.String())
	}
	if G.ActVerb.Norm != "walk" {
		t.Fatalf("expected ActVerb.Norm to be walk, got %q", G.ActVerb.Norm)
	}
	if !G.Params.HasWalkDir || G.Params.WalkDir != North {
		t.Fatalf("expected walk direction to be North")
	}
}

// ---- parser_parse_more_test.go ----

func TestParseAgainWithoutHistory(t *testing.T) {
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
	G.AllObjects = []*Object{room, player}
	G.Params.Continue = NumUndef

	G.GameInput = bytes.NewBufferString("again\n")
	G.GameOutput = &bytes.Buffer{}

	BuildVocabulary(nil, nil, nil)

	if ok := Parse(); ok {
		t.Fatalf("expected Parse to fail for again without history")
	}
}

func TestParseOopsWithoutUnknown(t *testing.T) {
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
	G.AllObjects = []*Object{room, player}
	G.Params.Continue = NumUndef

	G.GameInput = bytes.NewBufferString("oops look\n")
	G.GameOutput = &bytes.Buffer{}

	BuildVocabulary(nil, nil, nil)

	if ok := Parse(); ok {
		t.Fatalf("expected Parse to fail for oops without unknown")
	}
}

func TestParseThenSetsContinue(t *testing.T) {
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
	G.AllObjects = []*Object{room, player}
	G.Params.Continue = NumUndef

	G.GameInput = bytes.NewBufferString("look then look\n")
	G.GameOutput = &bytes.Buffer{}

	commands := []Syntax{
		{
			Verb:   "look",
			Action: func(ActionArg) bool { return true },
		},
	}
	BuildVocabulary(commands, []string{"then"}, nil)

	if ok := Parse(); !ok {
		t.Fatalf("expected Parse to succeed for then sequence")
	}
	if G.Params.Continue == NumUndef {
		t.Fatalf("expected Continue to be set for then")
	}
}

func TestParseLeadingThenTreatedAsBuzz(t *testing.T) {
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
	G.AllObjects = []*Object{room, player}
	G.Params.Continue = NumUndef

	G.GameInput = bytes.NewBufferString("then look\n")
	G.GameOutput = &bytes.Buffer{}

	commands := []Syntax{
		{
			Verb:   "look",
			Action: func(ActionArg) bool { return true },
		},
	}
	BuildVocabulary(commands, []string{"then", "the"}, nil)

	if ok := Parse(); !ok {
		t.Fatalf("expected Parse to treat leading then as buzz")
	}
	if G.ActVerb.Norm != "look" {
		t.Fatalf("expected ActVerb.Norm to be look, got %q", G.ActVerb.Norm)
	}
}

func TestParseAgainRepeatsMistake(t *testing.T) {
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
	G.AllObjects = []*Object{room, player}
	G.Params.Continue = NumUndef

	G.Again.Buf = []LexItem{{Norm: "look", Orig: "look", Types: WordTypes{WordVerb}}}
	G.ParserOk = false

	G.GameInput = bytes.NewBufferString("again\n")
	G.GameOutput = &bytes.Buffer{}

	BuildVocabulary(nil, nil, nil)

	if ok := Parse(); ok {
		t.Fatalf("expected Parse to fail for repeating mistake")
	}
	if out := G.GameOutput.(*bytes.Buffer).String(); !bytes.Contains([]byte(out), []byte("repeat a mistake")) {
		t.Fatalf("expected repeat-mistake message, got %q", out)
	}
}

func TestParseAgainUnexpectedWord(t *testing.T) {
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
	G.AllObjects = []*Object{room, player}
	G.Params.Continue = NumUndef

	G.Again.Buf = []LexItem{{Norm: "look", Orig: "look", Types: WordTypes{WordVerb}}}
	G.ParserOk = true

	G.GameInput = bytes.NewBufferString("again look\n")
	G.GameOutput = &bytes.Buffer{}

	BuildVocabulary(nil, nil, nil)

	if ok := Parse(); ok {
		t.Fatalf("expected Parse to fail for unexpected word after again")
	}
	if out := G.GameOutput.(*bytes.Buffer).String(); !bytes.Contains([]byte(out), []byte("couldn't understand")) {
		t.Fatalf("expected parse error message, got %q", out)
	}
}

func TestParseDirectionWithAnd(t *testing.T) {
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
	G.AllObjects = []*Object{room, player}
	G.Params.Continue = NumUndef

	G.GameInput = bytes.NewBufferString("north and south\n")
	G.GameOutput = &bytes.Buffer{}

	BuildVocabulary(nil, []string{"and"}, nil)

	if ok := Parse(); !ok {
		t.Fatalf("expected Parse to succeed for direction with and")
	}
	if !G.Params.HasWalkDir || G.Params.WalkDir != North {
		t.Fatalf("expected walk direction to be North")
	}
	if G.Params.Continue == NumUndef {
		t.Fatalf("expected Continue to be set for chained direction")
	}
}

func TestParseOopsWithPunctuation(t *testing.T) {
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
	G.AllObjects = []*Object{room, player}
	G.Params.Continue = NumUndef

	G.GameInput = bytes.NewBufferString("oops , look\n")
	G.GameOutput = &bytes.Buffer{}

	BuildVocabulary(nil, nil, nil)

	if ok := Parse(); ok {
		t.Fatalf("expected Parse to fail for oops with punctuation")
	}
	if out := G.GameOutput.(*bytes.Buffer).String(); !bytes.Contains([]byte(out), []byte("There was no word")) {
		t.Fatalf("expected oops failure message, got %q", out)
	}
}

// ---- parser_snarf_objects_test.go ----

func TestSnarfObjectsWithButMerge(t *testing.T) {
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

	lamp := &Object{Desc: "lamp", Synonyms: []string{"lamp"}, Flags: FlgTake, In: room}
	key := &Object{Desc: "key", Synonyms: []string{"key"}, Flags: FlgTake, In: room}
	room.AddChild(lamp)
	room.AddChild(key)

	G.Params.GetType = GetUndef
	G.DetectedSyntx = &Syntax{
		Obj1: ObjProp{LocFlags: LocSet(LocOnGrnd)},
	}
	G.ParsedSyntx.ObjOrClause1 = []LexItem{
		{Norm: "all", Orig: "all", Types: WordTypes{WordObj}},
		{Norm: "except", Orig: "except", Types: nil},
		{Norm: "key", Orig: "key", Types: WordTypes{WordObj}},
	}

	if ok := SnarfObjects(); !ok {
		t.Fatalf("expected SnarfObjects to succeed")
	}
	if len(G.Params.Buts) != 1 || G.Params.Buts[0] != key {
		t.Fatalf("expected key to be recorded as exclusion, got %+v", G.Params.Buts)
	}
	if len(G.DirObjPossibles) != 1 || G.DirObjPossibles[0] != key {
		t.Fatalf("expected DirObjPossibles to be filtered to key, got %+v", G.DirObjPossibles)
	}
}

func TestButMergeKeepsOnlyButs(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	a := &Object{Desc: "a"}
	b := &Object{Desc: "b"}
	c := &Object{Desc: "c"}
	G.Params.Buts = []*Object{b, a}

	res := ButMerge([]*Object{a, b, c})
	if len(res) != 2 || res[0] != b || res[1] != a {
		t.Fatalf("expected ButMerge to keep buts in order, got %+v", res)
	}
}

// ---- parser_syntax_test.go ----

func TestFindWhatIMeanKludgeReturnsRooms(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	rooms := &Object{Desc: "rooms"}
	G.RoomsObj = rooms

	if got := FindWhatIMean(FlgKludge, LocFlags(0), ""); got != rooms {
		t.Fatalf("expected FindWhatIMean to return RoomsObj for FlgKludge")
	}
}

func TestFindWhatIMeanPrepOutHands(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	room := &Object{Desc: "room"}
	player := &Object{Desc: "player", In: room}
	hands := &Object{Desc: "hands", Synonyms: []string{"hands"}, Flags: FlgTake, In: player}
	player.AddChild(hands)
	G.Player = player
	G.Winner = player
	G.Here = room
	G.HandsObj = hands
	G.Lit = true
	G.RoomsObj = &Object{Desc: "rooms"}
	room.In = G.RoomsObj
	G.AllObjects = []*Object{room, player, hands}
	G.GameOutput = &bytes.Buffer{}

	got := FindWhatIMean(FlgTake, LocSet(LocHeld), "out")
	if got != hands {
		t.Fatalf("expected FindWhatIMean to return hands")
	}
	if out := G.GameOutput.(*bytes.Buffer).String(); !strings.Contains(out, "(out of your hands)") {
		t.Fatalf("expected out-of-hands prompt, got %q", out)
	}
}

func TestFindWhatIMeanPrintsObject(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	room := &Object{Desc: "room"}
	player := &Object{Desc: "player", In: room}
	lamp := &Object{Desc: "lamp", Synonyms: []string{"lamp"}, Flags: FlgTake, In: room}
	room.AddChild(lamp)
	G.Player = player
	G.Winner = player
	G.Here = room
	G.Lit = true
	G.RoomsObj = &Object{Desc: "rooms"}
	room.In = G.RoomsObj
	G.AllObjects = []*Object{room, player, lamp}
	G.GameOutput = &bytes.Buffer{}

	got := FindWhatIMean(FlgTake, LocSet(LocOnGrnd), "")
	if got != lamp {
		t.Fatalf("expected FindWhatIMean to return lamp")
	}
	if out := G.GameOutput.(*bytes.Buffer).String(); !strings.Contains(out, "(lamp)") {
		t.Fatalf("expected object prompt, got %q", out)
	}
}

func TestSyntaxCheckUnknownSentence(t *testing.T) {
	oldG := G
	oldCommands := Commands
	G = NewGameState()
	t.Cleanup(func() {
		G = oldG
		Commands = oldCommands
	})

	G.ParsedSyntx.Verb = LexItem{Norm: "take", Orig: "take", Types: WordTypes{WordVerb}}
	G.Params.ObjOrClauseCnt = 0
	Commands = []Syntax{{Verb: "look"}}
	G.GameOutput = &bytes.Buffer{}

	if ok := SyntaxCheck(); ok {
		t.Fatalf("expected SyntaxCheck to fail for unknown sentence")
	}
	if out := G.GameOutput.(*bytes.Buffer).String(); !strings.Contains(out, "That sentence isn't one I recognize") {
		t.Fatalf("expected unrecognized sentence output, got %q", out)
	}
}

func TestSyntaxCheckFindWhatIMean(t *testing.T) {
	oldG := G
	oldCommands := Commands
	G = NewGameState()
	t.Cleanup(func() {
		G = oldG
		Commands = oldCommands
	})

	room := &Object{Desc: "room"}
	player := &Object{Desc: "player", In: room}
	lamp := &Object{Desc: "lamp", Synonyms: []string{"lamp"}, Flags: FlgTake, In: room}
	room.AddChild(lamp)
	G.Player = player
	G.Winner = player
	G.Here = room
	G.Lit = true
	G.RoomsObj = &Object{Desc: "rooms"}
	room.In = G.RoomsObj
	G.AllObjects = []*Object{room, player, lamp}
	G.GameOutput = &bytes.Buffer{}

	G.ParsedSyntx.Verb = LexItem{Norm: "take", Orig: "take", Types: WordTypes{WordVerb}}
	G.Params.ObjOrClauseCnt = 0
	Commands = []Syntax{
		{
			Verb: "take",
			Obj1: ObjProp{HasObj: true, ObjFlags: FlgTake, LocFlags: LocSet(LocOnGrnd)},
		},
	}

	if ok := SyntaxCheck(); !ok {
		t.Fatalf("expected SyntaxCheck to succeed via FindWhatIMean")
	}
	if G.DetectedSyntx == nil || len(G.DirObjPossibles) != 1 || G.DirObjPossibles[0] != lamp {
		t.Fatalf("expected DirObjPossibles to be set to lamp")
	}
	if out := G.GameOutput.(*bytes.Buffer).String(); !strings.Contains(out, "(lamp)") {
		t.Fatalf("expected FindWhatIMean prompt, got %q", out)
	}
}

func TestSyntaxCheckCanNotOrphanForNPC(t *testing.T) {
	oldG := G
	oldCommands := Commands
	G = NewGameState()
	t.Cleanup(func() {
		G = oldG
		Commands = oldCommands
	})

	room := &Object{Desc: "room"}
	player := &Object{Desc: "player", In: room}
	npc := &Object{Desc: "npc", In: room}
	G.Player = player
	G.Winner = npc
	G.Here = room
	G.Lit = true
	G.GlobalObj = &Object{Desc: "global"}
	G.RoomsObj = &Object{Desc: "rooms"}
	room.In = G.RoomsObj
	G.AllObjects = []*Object{room, player, npc}
	G.GameOutput = &bytes.Buffer{}

	G.ParsedSyntx.Verb = LexItem{Norm: "take", Orig: "take", Types: WordTypes{WordVerb}}
	G.Params.ObjOrClauseCnt = 0
	Commands = []Syntax{
		{
			Verb: "take",
			Obj1: ObjProp{HasObj: true, ObjFlags: FlgTake, LocFlags: LocSet(LocOnGrnd)},
		},
	}

	if ok := SyntaxCheck(); ok {
		t.Fatalf("expected SyntaxCheck to fail for NPC orphan case")
	}
	if out := G.GameOutput.(*bytes.Buffer).String(); !strings.Contains(out, "I don't understand") {
		t.Fatalf("expected CanNotOrphan output, got %q", out)
	}
}

func TestCanNotOrphan(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.GameOutput = &bytes.Buffer{}
	CanNotOrphan()
	if out := G.GameOutput.(*bytes.Buffer).String(); !strings.Contains(out, "I don't understand") {
		t.Fatalf("expected CanNotOrphan output, got %q", out)
	}
}
