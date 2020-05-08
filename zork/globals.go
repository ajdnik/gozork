package zork

var (
	GlobalObjects = Object{
		Flags: []Flag{FlgKludge, FlgInvis, FlgTouch, FlgSurf, FlgTryTake, FlgOpen, FlgSearch, FlgTrans, FlgOn, FlgLand, FlgFight, FlgStagg, FlgWear},
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
		Flags:    []Flag{FlgNoDesc, FlgTouch},
	}
	Hands = Object{
		In:         &GlobalObjects,
		Synonyms:   []string{"pair", "hands", "hand"},
		Adjectives: []string{"bare"},
		Desc:       "pair of hands",
		Flags:      []Flag{FlgNoDesc, FlgTool},
	}
	Me = Object{
		In:       &GlobalObjects,
		Synonyms: []string{"me", "myself", "self", "cretin"},
		Desc:     "you",
		Flags:    []Flag{FlgPerson},
		// Action:   CretinFcn,
	}
	Adventurer = Object{
		Synonyms: []string{"adventurer"},
		Desc:     "cretin",
		Flags:    []Flag{FlgNoDesc, FlgInvis, FlgSacred, FlgPerson},
	}
	Stairs = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"stairs", "steps", "staircase", "stairway"},
		Adjectives: []string{"stone", "dark", "marble", "forbidding", "steep"},
		Desc:       "stairs",
		Flags:      []Flag{FlgNoDesc, FlgClimb},
		Action:     StairsFcn,
	}
	Intnum = Object{
		In:       &GlobalObjects,
		Synonyms: []string{"intnum"},
		Desc:     "number",
		Flags:    []Flag{FlgTool},
	}
	Blessings = Object{
		In:       &GlobalObjects,
		Synonyms: []string{"blessings", "graces"},
		Desc:     "blessings",
		Flags:    []Flag{FlgNoDesc},
	}
	Sailor = Object{
		In:       &GlobalObjects,
		Synonyms: []string{"sailor", "footpad", "aviator"},
		Desc:     "sailor",
		Flags:    []Flag{FlgNoDesc},
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
		Flags:    []Flag{FlgNoDesc},
	}
	PathObj = Object{
		In:         &GlobalObjects,
		Synonyms:   []string{"trail", "path"},
		Adjectives: []string{"forest", "narrow", "long", "winding"},
		Desc:       "passage",
		Flags:      []Flag{FlgNoDesc},
		Action:     PathObject,
	}
	Zorkmid = Object{
		In:       &GlobalObjects,
		Synonyms: []string{"zorkmid"},
		Desc:     "zorkmid",
		Action:   ZorkmidFunction,
	}
	LoadAllowed = 100
	LoadMax     = 100
)

func NotHereObjectFcn(arg ActArg) bool {
	if DirObj == &NotHereObject && IndirObj == &NotHereObject {
		Print("Those things aren't here!", Newline)
		return true
	}
	Params.Continue = NumUndef
	Params.InQuotes = false
	isDir := false
	if DirObj == &NotHereObject {
		isDir = true
	}
	if Winner == Player {
		Print("You can't see any ", NoNewline)
		NotHerePrint(isDir)
		Print(" here!", Newline)
		return true
	}
	Print("The ", NoNewline)
	PrintObject(Winner)
	Print(" seems confused. \"I don't see any ", NoNewline)
	NotHerePrint(isDir)
	Print(" here!\"", Newline)
	return true
}

func NotHerePrint(isDir bool) {
	if Params.ShldOrphan {
		if NotHere.Adj.IsSet() {
			Print(NotHere.Adj.Orig+" ", NoNewline)
		}
		if NotHere.Syn.IsSet() {
			Print(NotHere.Syn.Orig, NoNewline)
		}
		return
	}
	if isDir {
		for idx, wrd := range ParsedSyntx.ObjOrClause1 {
			if idx != 0 {
				Print(" ", NoNewline)
			}
			Print(wrd.Orig, NoNewline)
		}
		return
	}
	for idx, wrd := range ParsedSyntx.ObjOrClause2 {
		if idx != 0 {
			Print(" ", NoNewline)
		}
		Print(wrd.Orig, NoNewline)
	}
}

func SailorFcn(arg ActArg) bool {
	if ActVerb == "tell" {
		Params.Continue = NumUndef
		Params.InQuotes = false
		Print("You can't talk to the sailor that way.", Newline)
		return true
	}
	if ActVerb == "examine" {
		Print("There is no sailor to be seen.", Newline)
		return true
	}
	if ActVerb == "hello" {
		HelloSailor++
		if HelloSailor%20 == 0 {
			Print("You seem to be repeating yourself.", Newline)
		} else if HelloSailor%10 == 0 {
			Print("I think that phrase is getting a bit worn out.", Newline)
		} else {
			Print("Nothing happens here.", Newline)
		}
		return true
	}
	return false
}

func GroundFunction(arg ActArg) bool {
	if (ActVerb == "put" || ActVerb == "put on") && IndirObj == &Ground {
		Perform("drop", DirObj, nil)
		return true
	}
	if Here == &SandyCave {
		return SandFunction(ActUnk)
	}
	if ActVerb == "dig" {
		Print("The ground is too hard for digging here.", Newline)
		return true
	}
	return false
}

func GrueFunction(arg ActArg) bool {
	if ActVerb == "examine" {
		Print("The grue is a sinister, lurking presence in the dark places of the earth. Its favorite diet is adventurers, but its insatiable appetite is tempered by its fear of light. No grue has ever been seen by the light of day, and few have survived its fearsome jaws to tell the tale.", Newline)
		return true
	}
	if ActVerb == "find" {
		Print("There is no grue here, but I'm sure there is at least one lurking in the darkness nearby. I wouldn't let my light go out if I were you!", Newline)
		return true
	}
	if ActVerb == "listen" {
		Print("It makes no sound but is always lurking in the darkness nearby.", Newline)
		return true
	}
	return false
}

func CretinFcn(arg ActArg) bool {
	if ActVerb == "tell" {
		Params.Continue = NumUndef
		Params.InQuotes = false
		Print("Talking to yourself is said to be a sign of impending mental collapse.", Newline)
		return true
	}
	if ActVerb == "give" && IndirObj == &Me {
		Perform("take", DirObj, nil)
		return true
	}
	if ActVerb == "make" {
		Print("Only you can do that.", Newline)
		return true
	}
	if ActVerb == "disembark" {
		Print("You'll have to do that on your own.", Newline)
		return true
	}
	if ActVerb == "eat" {
		Print("Auto-cannibalism is not the answer.", Newline)
		return true
	}
	if ActVerb == "attack" || ActVerb == "mung" {
		if IndirObj != nil && IndirObj.Has(FlgWeapon) {
			return JigsUp("If you insist.... Poof, you're dead!", false)
		}
		Print("Suicide is not the answer.", Newline)
		return true
	}
	if ActVerb == "throw" && DirObj == &Me {
		Print("Why don't you just walk like normal people?", Newline)
		return true
	}
	if ActVerb == "take" {
		Print("How romantic!", Newline)
		return true
	}
	if ActVerb == "examine" {
		if Here == &Mirror1 || Here == &Mirror2 {
			Print("Your image in the mirror looks tired.", Newline)
			return true
		}
		Print("That's difficult unless your eyes are prehensile.", Newline)
		return true
	}
	return false
}

func PathObject(arg ActArg) bool {
	if ActVerb == "take" || ActVerb == "follow" {
		Print("You must specify a direction to go.", Newline)
		return true
	}
	if ActVerb == "find" {
		Print("I can't help you there....", Newline)
		return true
	}
	if ActVerb == "dig" {
		Print("Not a chance.", Newline)
		return true
	}
	return false
}

func StairsFcn(arg ActArg) bool {
	if ActVerb == "through" {
		Print("You should say whether you want to go up or down.", Newline)
		return true
	}
	return false
}

func ZorkmidFunction(arg ActArg) bool {
	if ActVerb == "examine" {
		Print("The zorkmid is the unit of currency of the Great Underground Empire.", Newline)
		return true
	}
	if ActVerb == "find" {
		Print("The best way to find zorkmids is to go out and look for them.", Newline)
		return true
	}
	return false
}
