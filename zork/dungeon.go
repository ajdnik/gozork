package zork

var (
	InflatedBoat = Object{
		Synonyms:   []string{"boat", "raft"},
		Adjectives: []string{"magic", "plastic", "seaworthy", "inflated", "inflatable"},
		Desc:       "magic boat",
		Flags:      []Flag{FlgTake, FlgBurn, FlgVeh, FlgOpen, FlgSearch},
		Action:     RBoatFcn,
		Capacity:   100,
		Size:       20,
		VehType:    FlgNonLand,
	}
	WhiteHouse = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"house"},
		Adjectives: []string{"white", "beautiful", "colonial"},
		Desc:       "white house",
		Flags:      []Flag{FlgNoDesc},
		Action:     WhiteHouseFcn,
	}
	Board = Object{
		In:       &LocalGlobals,
		Synonyms: []string{"boards", "board"},
		Desc:     "board",
		Flags:    []Flag{FlgNoDesc},
		Action:   BoardFcn,
	}
	Forest = Object{
		In:       &LocalGlobals,
		Synonyms: []string{"forest", "trees", "pines", "hemlocks"},
		Desc:     "forest",
		Flags:    []Flag{FlgNoDesc},
		Action:   ForestFcn,
	}
	WestOfHouse = Object{
		In:     &Rooms,
		Desc:   "West of House",
		Action: WestHouseFcn,
		Flags:  []Flag{FlgLand, FlgOn, FlgSacred},
		Global: []*Object{&WhiteHouse, &Board, &Forest},
	}
	Mailbox = Object{
		In:         &WestOfHouse,
		Synonyms:   []string{"mailbox", "box"},
		Adjectives: []string{"small"},
		Desc:       "small mailbox",
		Flags:      []Flag{FlgCont, FlgTryTake},
		Capacity:   10,
		Action:     MailboxFcn,
	}
	Objects = []*Object{
		&Rooms,
		&InflatedBoat,
		&WestOfHouse,
		&Mailbox,
		&Hands,
		&Me,
		&GlobalObjects,
		&LocalGlobals,
		&NotHereObject,
		&PseudoObject,
		&Adventurer,
		&It,
		&WhiteHouse,
		&Board,
		&Forest,
	}
)
