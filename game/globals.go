package game

import . "github.com/ajdnik/gozork/engine"

var (
	GlobalObjects = Object{
		Flags: FlgKludge | FlgInvis | FlgTouch | FlgSurf | FlgTryTake | FlgOpen | FlgSearch | FlgTrans | FlgOn | FlgLand | FlgFight | FlgStagg | FlgWear,
	}
	LocalGlobals = Object{
		In:        &GlobalObjects,
		Synonyms:  []string{"zzmgck"},
		Global:    []*Object{&GlobalObjects},
		DescFcn:   PathObject,
		FirstDesc: "F",
		LongDesc:  "F",
		/* Pseudo: []PseudoObj{PseudoObj{
			Synonym: "foobar",
			Action:  VWalk,
		}}, */
	}
	Rooms         = Object{}
	NotHereObject = Object{
		Desc: "souch thing",
		// Action: NotHereObjectFcn,
	}
	PseudoObject = Object{
		In:     &LocalGlobals,
		Desc:   "pseudo",
		Action: CretinFcn,
	}
	It = Object{
		In:       &GlobalObjects,
		Synonyms: []string{"it", "them", "her", "him"},
		Desc:     "random object",
		Flags:    FlgNoDesc | FlgTouch,
	}
	Hands = Object{
		In:         &GlobalObjects,
		Synonyms:   []string{"pair", "hands", "hand"},
		Adjectives: []string{"bare"},
		Desc:       "pair of hands",
		Flags:      FlgNoDesc | FlgTool,
	}
	Me = Object{
		In:       &GlobalObjects,
		Synonyms: []string{"me", "myself", "self", "cretin"},
		Desc:     "you",
		Flags:    FlgPerson,
		// Action:   CretinFcn,
	}
	Adventurer = Object{
		Synonyms: []string{"adventurer"},
		Desc:     "cretin",
		Flags:    FlgNoDesc | FlgInvis | FlgSacred | FlgPerson,
	}
	Stairs = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"stairs", "steps", "staircase", "stairway"},
		Adjectives: []string{"stone", "dark", "marble", "forbidding", "steep"},
		Desc:       "stairs",
		Flags:      FlgNoDesc | FlgClimb,
		Action:     StairsFcn,
	}
	Intnum = Object{
		In:       &GlobalObjects,
		Synonyms: []string{"intnum"},
		Desc:     "number",
		Flags:    FlgTool,
	}
	Blessings = Object{
		In:       &GlobalObjects,
		Synonyms: []string{"blessings", "graces"},
		Desc:     "blessings",
		Flags:    FlgNoDesc,
	}
	Sailor = Object{
		In:       &GlobalObjects,
		Synonyms: []string{"sailor", "footpad", "aviator"},
		Desc:     "sailor",
		Flags:    FlgNoDesc,
		Action:   SailorFcn,
	}
	Ground = Object{
		In:       &GlobalObjects,
		Synonyms: []string{"ground", "sand", "dirt", "floor"},
		Desc:     "ground",
		// Action:   GroundFunction,
	}
	Grue = Object{
		In:         &GlobalObjects,
		Synonyms:   []string{"grue"},
		Adjectives: []string{"lurking", "sinister", "hungry", "silent"},
		Desc:       "lurking grue",
		Action:     GrueFunction,
	}
	Lungs = Object{
		In:       &GlobalObjects,
		Synonyms: []string{"lungs", "air", "mouth", "breath"},
		Desc:     "blast of air",
		Flags:    FlgNoDesc,
	}
	PathObj = Object{
		In:         &GlobalObjects,
		Synonyms:   []string{"trail", "path"},
		Adjectives: []string{"forest", "narrow", "long", "winding"},
		Desc:       "passage",
		Flags:      FlgNoDesc,
		Action:     PathObject,
	}
	Zorkmid = Object{
		In:       &GlobalObjects,
		Synonyms: []string{"zorkmid"},
		Desc:     "zorkmid",
		Action:   ZorkmidFunction,
	}
)

func NotHereObjectFcn(arg ActArg) bool {
	if G.DirObj == &NotHereObject && G.IndirObj == &NotHereObject {
		Printf("Those things aren't here!\n")
		return true
	}
	G.Params.Continue = NumUndef
	G.Params.InQuotes = false
	isDir := false
	if G.DirObj == &NotHereObject {
		isDir = true
	}
	if G.Winner == G.Player {
		Printf("You can't see any ")
		NotHerePrint(isDir)
		Printf(" here!\n")
		return true
	}
	Printf("The %s seems confused. \"I don't see any ", G.Winner.Desc)
	NotHerePrint(isDir)
	Printf(" here!\"\n")
	return true
}

func NotHerePrint(isDir bool) {
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

func SailorFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "tell" {
		G.Params.Continue = NumUndef
		G.Params.InQuotes = false
		Printf("You can't talk to the sailor that way.\n")
		return true
	}
	if G.ActVerb.Norm == "examine" {
		Printf("There is no sailor to be seen.\n")
		return true
	}
	if G.ActVerb.Norm == "hello" {
		GD().HelloSailor++
		if GD().HelloSailor%20 == 0 {
			Printf("You seem to be repeating yourself.\n")
		} else if GD().HelloSailor%10 == 0 {
			Printf("I think that phrase is getting a bit worn out.\n")
		} else {
			Printf("Nothing happens here.\n")
		}
		return true
	}
	return false
}

func GroundFunction(arg ActArg) bool {
	if (G.ActVerb.Norm == "put" || G.ActVerb.Norm == "put on") && G.IndirObj == &Ground {
		Perform(ActionVerb{Norm: "drop", Orig: "drop"}, G.DirObj, nil)
		return true
	}
	if G.Here == &SandyCave {
		return SandFunction(ActUnk)
	}
	if G.ActVerb.Norm == "dig" {
		Printf("The ground is too hard for digging here.\n")
		return true
	}
	return false
}

func GrueFunction(arg ActArg) bool {
	if G.ActVerb.Norm == "examine" {
		Printf("The grue is a sinister, lurking presence in the dark places of the earth. Its favorite diet is adventurers, but its insatiable appetite is tempered by its fear of light. No grue has ever been seen by the light of day, and few have survived its fearsome jaws to tell the tale.\n")
		return true
	}
	if G.ActVerb.Norm == "find" {
		Printf("There is no grue here, but I'm sure there is at least one lurking in the darkness nearby. I wouldn't let my light go out if I were you!\n")
		return true
	}
	if G.ActVerb.Norm == "listen" {
		Printf("It makes no sound but is always lurking in the darkness nearby.\n")
		return true
	}
	return false
}

func CretinFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "tell" {
		G.Params.Continue = NumUndef
		G.Params.InQuotes = false
		Printf("Talking to yourself is said to be a sign of impending mental collapse.\n")
		return true
	}
	if G.ActVerb.Norm == "give" && G.IndirObj == &Me {
		Perform(ActionVerb{Norm: "take", Orig: "take"}, G.DirObj, nil)
		return true
	}
	if G.ActVerb.Norm == "make" {
		Printf("Only you can do that.\n")
		return true
	}
	if G.ActVerb.Norm == "disembark" {
		Printf("You'll have to do that on your own.\n")
		return true
	}
	if G.ActVerb.Norm == "eat" {
		Printf("Auto-cannibalism is not the answer.\n")
		return true
	}
	if G.ActVerb.Norm == "attack" || G.ActVerb.Norm == "mung" {
		if G.IndirObj != nil && G.IndirObj.Has(FlgWeapon) {
			return JigsUp("If you insist.... Poof, you're dead!", false)
		}
		Printf("Suicide is not the answer.\n")
		return true
	}
	if G.ActVerb.Norm == "throw" && G.DirObj == &Me {
		Printf("Why don't you just walk like normal people?\n")
		return true
	}
	if G.ActVerb.Norm == "take" {
		Printf("How romantic!\n")
		return true
	}
	if G.ActVerb.Norm == "examine" {
		if G.Here == &MirrorRoom1 || G.Here == &MirrorRoom2 {
			Printf("Your image in the mirror looks tired.\n")
			return true
		}
		Printf("That's difficult unless your eyes are prehensile.\n")
		return true
	}
	return false
}

func PathObject(arg ActArg) bool {
	if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "follow" {
		Printf("You must specify a direction to go.\n")
		return true
	}
	if G.ActVerb.Norm == "find" {
		Printf("I can't help you there....\n")
		return true
	}
	if G.ActVerb.Norm == "dig" {
		Printf("Not a chance.\n")
		return true
	}
	return false
}

func StairsFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "through" {
		Printf("You should say whether you want to go up or down.\n")
		return true
	}
	return false
}

func ZorkmidFunction(arg ActArg) bool {
	if G.ActVerb.Norm == "examine" {
		Printf("The zorkmid is the unit of currency of the Great Underground Empire.\n")
		return true
	}
	if G.ActVerb.Norm == "find" {
		Printf("The best way to find zorkmids is to go out and look for them.\n")
		return true
	}
	return false
}
