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
	LivingRoom = Object{
		In:     &Rooms,
		Desc:   "Living Room",
		Action: LivingRoomFcn,
		Flags:  []Flag{FlgLand, FlgOn, FlgSacred},
		Global: []*Object{&Stairs},
		Pseudo: []PseudoObj{
			{Synonym: "nails", Action: NailsPseudo},
			{Synonym: "nail", Action: NailsPseudo},
		},
	}
	TrophyCase = Object{
		In:         &LivingRoom,
		Synonyms:   []string{"case"},
		Adjectives: []string{"trophy"},
		Desc:       "trophy case",
		Flags:      []Flag{FlgTrans, FlgCont, FlgNoDesc, FlgTryTake, FlgSearch},
		Action:     TrophyCaseFcn,
		Capacity:   10000,
	}
	Map = Object{
		In:         &TrophyCase,
		Synonyms:   []string{"parchment", "map"},
		Adjectives: []string{"antique", "old", "ancient"},
		Desc:       "ancient map",
		Flags:      []Flag{FlgInvis, FlgRead, FlgTake},
		FirstDesc:  "In the trophy case is an ancient parchment which appears to be a map.",
		Size:       2,
		Text:       "The map shows a forest with three clearings. The largest clearing contains a house. Three paths leave the large clearing. One of these paths, leading southwest, is marked \"To Stone Barrow\".",
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
		// Action:     MailboxFcn,
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
		&TrophyCase,
		&LivingRoom,
		&Stairs,
		&Map,
	}
)

func FinalizeGameObjects() {
	Mailbox.Action = MailboxFcn
}
