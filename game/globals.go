package game

import . "github.com/ajdnik/gozork/engine"

var (
	globalObjects = Object{
		Flags: FlgKludge | FlgInvis | FlgTouch | FlgSurf | FlgTryTake | FlgOpen | FlgSearch | FlgTrans | FlgOn | FlgLand | FlgFight | FlgStaggered | FlgWear,
	}
	localGlobals = Object{
		In:        &globalObjects,
		Synonyms:  []string{"zzmgck"},
		Global:    []*Object{&globalObjects},
		DescFcn:   pathObject,
		FirstDesc: "F",
		LongDesc:  "F",
		/* Pseudo: []PseudoObj{PseudoObj{
			Synonym: "foobar",
			Action:  vWalk,
		}}, */
	}
	rooms         = Object{}
	notHereObject = Object{
		Desc: "souch thing",
		// Action: notHereObjectFcn,
	}
	pseudoObject = Object{
		In:     &localGlobals,
		Desc:   "pseudo",
		Action: cretinFcn,
	}
	it = Object{
		In:       &globalObjects,
		Synonyms: []string{"it", "them", "her", "him"},
		Desc:     "random object",
		Flags:    FlgNoDesc | FlgTouch,
	}
	hands = Object{
		In:         &globalObjects,
		Synonyms:   []string{"pair", "hands", "hand"},
		Adjectives: []string{"bare"},
		Desc:       "pair of hands",
		Flags:      FlgNoDesc | FlgTool,
	}
	me = Object{
		In:       &globalObjects,
		Synonyms: []string{"me", "myself", "self", "cretin"},
		Desc:     "you",
		Flags:    FlgPerson,
		// Action:   cretinFcn,
	}
	adventurer = Object{
		Synonyms: []string{"adventurer"},
		Desc:     "cretin",
		Flags:    FlgNoDesc | FlgInvis | FlgSacred | FlgPerson,
	}
	stairs = Object{
		In:         &localGlobals,
		Synonyms:   []string{"stairs", "steps", "staircase", "stairway"},
		Adjectives: []string{"stone", "dark", "marble", "forbidding", "steep"},
		Desc:       "stairs",
		Flags:      FlgNoDesc | FlgClimb,
		Action:     stairsFcn,
	}
	intnum = Object{
		In:       &globalObjects,
		Synonyms: []string{"intnum"},
		Desc:     "number",
		Flags:    FlgTool,
	}
	blessings = Object{
		In:       &globalObjects,
		Synonyms: []string{"blessings", "graces"},
		Desc:     "blessings",
		Flags:    FlgNoDesc,
	}
	sailor = Object{
		In:       &globalObjects,
		Synonyms: []string{"sailor", "footpad", "aviator"},
		Desc:     "sailor",
		Flags:    FlgNoDesc,
		Action:   sailorFcn,
	}
	ground = Object{
		In:       &globalObjects,
		Synonyms: []string{"ground", "sand", "dirt", "floor"},
		Desc:     "ground",
		// Action:   groundFunction,
	}
	grue = Object{
		In:         &globalObjects,
		Synonyms:   []string{"grue"},
		Adjectives: []string{"lurking", "sinister", "hungry", "silent"},
		Desc:       "lurking grue",
		Action:     grueFunction,
	}
	lungs = Object{
		In:       &globalObjects,
		Synonyms: []string{"lungs", "air", "mouth", "breath"},
		Desc:     "blast of air",
		Flags:    FlgNoDesc,
	}
	pathObj = Object{
		In:         &globalObjects,
		Synonyms:   []string{"trail", "path"},
		Adjectives: []string{"forest", "narrow", "long", "winding"},
		Desc:       "passage",
		Flags:      FlgNoDesc,
		Action:     pathObject,
	}
	zorkmid = Object{
		In:       &globalObjects,
		Synonyms: []string{"zorkmid"},
		Desc:     "zorkmid",
		Action:   zorkmidFunction,
	}
)

func notHereObjectFcn(arg ActionArg) bool {
	if G.DirObj == &notHereObject && G.IndirObj == &notHereObject {
		Printf("Those things aren't here!\n")
		return true
	}
	G.Params.Continue = NumUndef
	G.Params.InQuotes = false
	isDir := G.DirObj == &notHereObject
	if G.Winner == G.Player {
		Printf("You can't see any ")
		notHerePrint(isDir)
		Printf(" here!\n")
		return true
	}
	Printf("The %s seems confused. \"I don't see any ", G.Winner.Desc)
	notHerePrint(isDir)
	Printf(" here!\"\n")
	return true
}

func notHerePrint(isDir bool) {
	if G.Params.ShldOrphan {
		if G.NotHere.Adj.IsSet() {
			Printf("%s", G.NotHere.Adj.Orig+" ")
		}
		if G.NotHere.Syn.IsSet() {
			Printf("%s", G.NotHere.Syn.Orig)
		}
		return
	}
	if isDir {
		for idx, wrd := range G.ParsedSyntx.ObjOrClause1 {
			if idx != 0 {
				Printf(" ")
			}
			Printf("%s", wrd.Orig)
		}
		return
	}
	for idx, wrd := range G.ParsedSyntx.ObjOrClause2 {
		if idx != 0 {
			Printf(" ")
		}
		Printf("%s", wrd.Orig)
	}
}

func sailorFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "tell":
		G.Params.Continue = NumUndef
		G.Params.InQuotes = false
		Printf("You can't talk to the sailor that way.\n")
		return true
	case "examine":
		Printf("There is no sailor to be seen.\n")
		return true
	case "hello":
		gD().HelloSailor++
		if gD().HelloSailor%20 == 0 {
			Printf("You seem to be repeating yourself.\n")
		} else if gD().HelloSailor%10 == 0 {
			Printf("I think that phrase is getting a bit worn out.\n")
		} else {
			Printf("Nothing happens here.\n")
		}
		return true
	}
	return false
}

func groundFunction(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "put", "put on":
		if G.IndirObj == &ground {
			Perform(ActionVerb{Norm: "drop", Orig: "drop"}, G.DirObj, nil)
			return true
		}
	}
	if G.Here == &sandyCave {
		return sandFunction(ActUnk)
	}
	switch G.ActVerb.Norm {
	case "dig":
		Printf("The ground is too hard for digging here.\n")
		return true
	}
	return false
}

func grueFunction(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "examine":
		Printf("The grue is a sinister, lurking presence in the dark places of the earth. Its favorite diet is adventurers, but its insatiable appetite is tempered by its fear of light. No grue has ever been seen by the light of day, and few have survived its fearsome jaws to tell the tale.\n")
		return true
	case "find":
		Printf("There is no grue here, but I'm sure there is at least one lurking in the darkness nearby. I wouldn't let my light go out if I were you!\n")
		return true
	case "listen":
		Printf("it makes no sound but is always lurking in the darkness nearby.\n")
		return true
	}
	return false
}

func cretinFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "tell":
		G.Params.Continue = NumUndef
		G.Params.InQuotes = false
		Printf("Talking to yourself is said to be a sign of impending mental collapse.\n")
		return true
	case "give":
		if G.IndirObj == &me {
			Perform(ActionVerb{Norm: "take", Orig: "take"}, G.DirObj, nil)
			return true
		}
	case "make":
		Printf("Only you can do that.\n")
		return true
	case "disembark":
		Printf("You'll have to do that on your own.\n")
		return true
	case "eat":
		Printf("Auto-cannibalism is not the answer.\n")
		return true
	case "attack", "mung":
		if G.IndirObj != nil && G.IndirObj.Has(FlgWeapon) {
			return jigsUp("If you insist.... Poof, you're dead!", false)
		}
		Printf("Suicide is not the answer.\n")
		return true
	case "throw":
		if G.DirObj == &me {
			Printf("Why don't you just walk like normal people?\n")
			return true
		}
	case "take":
		Printf("How romantic!\n")
		return true
	case "examine":
		if G.Here == &mirrorRoom1 || G.Here == &mirrorRoom2 {
			Printf("Your image in the mirror looks tired.\n")
			return true
		}
		Printf("That's difficult unless your eyes are prehensile.\n")
		return true
	}
	return false
}

func pathObject(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "take", "follow":
		Printf("You must specify a direction to go.\n")
		return true
	case "find":
		Printf("I can't help you there....\n")
		return true
	case "dig":
		Printf("Not a chance.\n")
		return true
	}
	return false
}

func stairsFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "through":
		Printf("You should say whether you want to go up or down.\n")
		return true
	}
	return false
}

func zorkmidFunction(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "examine":
		Printf("The zorkmid is the unit of currency of the Great Underground Empire.\n")
		return true
	case "find":
		Printf("The best way to find zorkmids is to go out and look for them.\n")
		return true
	}
	return false
}
