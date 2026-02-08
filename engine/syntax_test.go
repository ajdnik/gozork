package engine

import (
	"testing"
)

func TestLocFlagsBitset(t *testing.T) {
	flags := LocSet(LocHeld, LocTake)
	if !LocHeld.In(flags) {
		t.Fatalf("expected LocHeld to be set")
	}
	if LocHave.In(flags) {
		t.Fatalf("did not expect LocHave to be set")
	}
	if flags.HasAll() {
		t.Fatalf("did not expect HasAll to be true")
	}

	var all LocFlags
	all.All()
	if !all.HasAll() {
		t.Fatalf("expected HasAll to be true after All()")
	}
	for _, lf := range []LocFlag{LocHeld, LocCarried, LocInRoom, LocOnGrnd, LocTake, LocMany, LocHave} {
		if !lf.In(all) {
			t.Fatalf("expected %v to be set in LocAll", lf)
		}
	}
}

// ---- syntax_build_test.go ----

func TestBuildVocabulary(t *testing.T) {
	oldG := G
	oldVocab := Vocabulary
	G = NewGameState()
	Vocabulary = make(map[string]WordItem)
	t.Cleanup(func() {
		G = oldG
		Vocabulary = oldVocab
	})

	room := &Object{
		Desc:       "room",
		Synonyms:   []string{"lantern"},
		Adjectives: []string{"brass"},
	}
	G.AllObjects = []*Object{room}

	act := func(ActionArg) bool { return true }
	pre := func(ActionArg) bool { return true }
	commands := []Syntax{
		{
			Verb:      "take",
			Action:    act,
			PreAction: pre,
		},
		{
			NormVerb: "look",
			Verb:     "examine",
			Action:   act,
		},
	}

	synonyms := map[string]string{
		"grab": "take",
	}

	BuildVocabulary(commands, []string{"the"}, synonyms)

	if len(Commands) != 2 {
		t.Fatalf("expected Commands to be set by BuildVocabulary")
	}

	if _, ok := G.Actions["take"]; !ok {
		t.Fatalf("expected actions for base verb")
	}
	if _, ok := G.PreActions["take"]; !ok {
		t.Fatalf("expected pre-actions for base verb")
	}
	if _, ok := G.Actions["look"]; !ok {
		t.Fatalf("expected actions for normalized verb")
	}
	if _, ok := G.PreActions["look"]; ok {
		t.Fatalf("did not expect pre-actions for look without PreAction")
	}

	if _, ok := Vocabulary["north"]; !ok {
		t.Fatalf("expected directions to be added to vocabulary")
	}
	if _, ok := Vocabulary["lantern"]; !ok {
		t.Fatalf("expected object synonyms to be added to vocabulary")
	}
	if _, ok := Vocabulary["brass"]; !ok {
		t.Fatalf("expected object adjectives to be added to vocabulary")
	}

	if _, ok := Vocabulary["the"]; !ok {
		t.Fatalf("expected buzz words to be added to vocabulary")
	}

	if syn, ok := Vocabulary["grab"]; !ok || syn.Norm != "take" {
		t.Fatalf("expected synonyms to map to normalized term")
	}
	if syn := Vocabulary["grab"]; !syn.Types.Has(WordVerb) {
		t.Fatalf("expected synonym to inherit types from target word")
	}
}

// ---- syntax_vocab_test.go ----

func TestAddToVocab(t *testing.T) {
	oldVocab := Vocabulary
	Vocabulary = make(map[string]WordItem)
	t.Cleanup(func() { Vocabulary = oldVocab })

	AddToVocab("lamp", WordObj)
	if v, ok := Vocabulary["lamp"]; !ok || v.Norm != "lamp" || !v.Types.Has(WordObj) {
		t.Fatalf("expected lamp to be added as WordObj")
	}

	AddToVocab("lamp", WordAdj)
	v := Vocabulary["lamp"]
	if !v.Types.Has(WordObj) || !v.Types.Has(WordAdj) {
		t.Fatalf("expected lamp to include both WordObj and WordAdj")
	}
}

func TestSyntaxNumObjects(t *testing.T) {
	s := Syntax{}
	if s.NumObjects() != 0 {
		t.Fatalf("expected zero objects when Obj1.HasObj is false")
	}
	s.Obj1.HasObj = true
	if s.NumObjects() != 1 {
		t.Fatalf("expected one object when only Obj1 is set")
	}
	s.Obj2.HasObj = true
	if s.NumObjects() != 2 {
		t.Fatalf("expected two objects when Obj1 and Obj2 are set")
	}
}

func TestWordTypesHas(t *testing.T) {
	types := WordTypes{WordAdj, WordVerb}
	if !types.Has(WordAdj) {
		t.Fatalf("expected WordTypes.Has to find WordAdj")
	}
	if types.Has(WordObj) {
		t.Fatalf("expected WordTypes.Has to be false for missing type")
	}
}
