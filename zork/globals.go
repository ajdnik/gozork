package zork

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
		Print("Those things aren't here!", Newline)
		return true
	}
	G.Params.Continue = NumUndef
	G.Params.InQuotes = false
	isDir := false
	if G.DirObj == &NotHereObject {
		isDir = true
	}
	if G.Winner == G.Player {
		Print("You can't see any ", NoNewline)
		NotHerePrint(isDir)
		Print(" here!", Newline)
		return true
	}
	Print("The ", NoNewline)
	PrintObject(G.Winner)
	Print(" seems confused. \"I don't see any ", NoNewline)
	NotHerePrint(isDir)
	Print(" here!\"", Newline)
	return true
}

func NotHerePrint(isDir bool) {
	if G.Params.ShldOrphan {
		if G.NotHere.Adj.IsSet() {
			Print(G.NotHere.Adj.Orig+" ", NoNewline)
		}
		if G.NotHere.Syn.IsSet() {
			Print(G.NotHere.Syn.Orig, NoNewline)
		}
		return
	}
	if isDir {
		for idx, wrd := range G.ParsedSyntx.ObjOrClause1 {
			if idx != 0 {
				Print(" ", NoNewline)
			}
			Print(wrd.Orig, NoNewline)
		}
		return
	}
	for idx, wrd := range G.ParsedSyntx.ObjOrClause2 {
		if idx != 0 {
			Print(" ", NoNewline)
		}
		Print(wrd.Orig, NoNewline)
	}
}

func SailorFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "tell" {
		G.Params.Continue = NumUndef
		G.Params.InQuotes = false
		Print("You can't talk to the sailor that way.", Newline)
		return true
	}
	if G.ActVerb.Norm == "examine" {
		Print("There is no sailor to be seen.", Newline)
		return true
	}
	if G.ActVerb.Norm == "hello" {
		G.HelloSailor++
		if G.HelloSailor%20 == 0 {
			Print("You seem to be repeating yourself.", Newline)
		} else if G.HelloSailor%10 == 0 {
			Print("I think that phrase is getting a bit worn out.", Newline)
		} else {
			Print("Nothing happens here.", Newline)
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
		Print("The ground is too hard for digging here.", Newline)
		return true
	}
	return false
}

func GrueFunction(arg ActArg) bool {
	if G.ActVerb.Norm == "examine" {
		Print("The grue is a sinister, lurking presence in the dark places of the earth. Its favorite diet is adventurers, but its insatiable appetite is tempered by its fear of light. No grue has ever been seen by the light of day, and few have survived its fearsome jaws to tell the tale.", Newline)
		return true
	}
	if G.ActVerb.Norm == "find" {
		Print("There is no grue here, but I'm sure there is at least one lurking in the darkness nearby. I wouldn't let my light go out if I were you!", Newline)
		return true
	}
	if G.ActVerb.Norm == "listen" {
		Print("It makes no sound but is always lurking in the darkness nearby.", Newline)
		return true
	}
	return false
}

func CretinFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "tell" {
		G.Params.Continue = NumUndef
		G.Params.InQuotes = false
		Print("Talking to yourself is said to be a sign of impending mental collapse.", Newline)
		return true
	}
	if G.ActVerb.Norm == "give" && G.IndirObj == &Me {
		Perform(ActionVerb{Norm: "take", Orig: "take"}, G.DirObj, nil)
		return true
	}
	if G.ActVerb.Norm == "make" {
		Print("Only you can do that.", Newline)
		return true
	}
	if G.ActVerb.Norm == "disembark" {
		Print("You'll have to do that on your own.", Newline)
		return true
	}
	if G.ActVerb.Norm == "eat" {
		Print("Auto-cannibalism is not the answer.", Newline)
		return true
	}
	if G.ActVerb.Norm == "attack" || G.ActVerb.Norm == "mung" {
		if G.IndirObj != nil && G.IndirObj.Has(FlgWeapon) {
			return JigsUp("If you insist.... Poof, you're dead!", false)
		}
		Print("Suicide is not the answer.", Newline)
		return true
	}
	if G.ActVerb.Norm == "throw" && G.DirObj == &Me {
		Print("Why don't you just walk like normal people?", Newline)
		return true
	}
	if G.ActVerb.Norm == "take" {
		Print("How romantic!", Newline)
		return true
	}
	if G.ActVerb.Norm == "examine" {
		if G.Here == &MirrorRoom1 || G.Here == &MirrorRoom2 {
			Print("Your image in the mirror looks tired.", Newline)
			return true
		}
		Print("That's difficult unless your eyes are prehensile.", Newline)
		return true
	}
	return false
}

func PathObject(arg ActArg) bool {
	if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "follow" {
		Print("You must specify a direction to go.", Newline)
		return true
	}
	if G.ActVerb.Norm == "find" {
		Print("I can't help you there....", Newline)
		return true
	}
	if G.ActVerb.Norm == "dig" {
		Print("Not a chance.", Newline)
		return true
	}
	return false
}

func StairsFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "through" {
		Print("You should say whether you want to go up or down.", Newline)
		return true
	}
	return false
}

func ZorkmidFunction(arg ActArg) bool {
	if G.ActVerb.Norm == "examine" {
		Print("The zorkmid is the unit of currency of the Great Underground Empire.", Newline)
		return true
	}
	if G.ActVerb.Norm == "find" {
		Print("The best way to find zorkmids is to go out and look for them.", Newline)
		return true
	}
	return false
}
