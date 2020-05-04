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
		Pseudo: []PseudoObj{PseudoObj{
			Synonym: "foobar",
			Action:  VWalk,
		}},
	}
	Rooms         = Object{}
	NotHereObject = Object{
		Desc:   "souch thing",
		Action: NotHereObjectFcn,
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
		Action:   CretinFcn,
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
)

func NotHereObjectFcn(arg ActArg) bool {
	return false
}

func CretinFcn(arg ActArg) bool {
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
