package zork

var (
	// ================================================================
	// ROOMS - Forest and Outside
	// ================================================================

	WestOfHouse = Object{
		In:     &Rooms,
		Desc:   "West of House",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgLand | FlgOn | FlgSacred,
		Global: []*Object{&WhiteHouse, &Board, &Forest},
	}
	StoneBarrow = Object{
		In:       &Rooms,
		LongDesc: "You are standing in front of a massive barrow of stone. In the east face is a huge stone door which is open. You cannot see into the dark of the tomb.",
		Desc:     "Stone Barrow",
		Flags:    FlgLand | FlgOn | FlgSacred,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	NorthOfHouse = Object{
		In:       &Rooms,
		LongDesc: "You are facing the north side of a white house. There is no door here, and all the windows are boarded up. To the north a narrow path winds through the trees.",
		Desc:     "North of House",
		Flags:    FlgLand | FlgOn | FlgSacred,
		Global:   []*Object{&BoardedWindow, &Board, &WhiteHouse, &Forest},
	}
	SouthOfHouse = Object{
		In:       &Rooms,
		LongDesc: "You are facing the south side of a white house. There is no door here, and all the windows are boarded.",
		Desc:     "South of House",
		Flags:    FlgLand | FlgOn | FlgSacred,
		Global:   []*Object{&BoardedWindow, &Board, &WhiteHouse, &Forest},
	}
	EastOfHouse = Object{
		In:     &Rooms,
		Desc:   "Behind House",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgLand | FlgOn | FlgSacred,
		Global: []*Object{&WhiteHouse, &KitchenWindow, &Forest},
	}
	Forest1 = Object{
		In:       &Rooms,
		LongDesc: "This is a forest, with trees in all directions. To the east, there appears to be sunlight.",
		Desc:     "Forest",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:    FlgLand | FlgOn | FlgSacred,
		Global:   []*Object{&Tree, &Songbird, &WhiteHouse, &Forest},
	}
	Forest2 = Object{
		In:       &Rooms,
		LongDesc: "This is a dimly lit forest, with large trees all around.",
		Desc:     "Forest",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:    FlgLand | FlgOn | FlgSacred,
		Global:   []*Object{&Tree, &Songbird, &WhiteHouse, &Forest},
	}
	Mountains = Object{
		In:       &Rooms,
		LongDesc: "The forest thins out, revealing impassable mountains.",
		Desc:     "Forest",
		Flags:    FlgLand | FlgOn | FlgSacred,
		Global:   []*Object{&Tree, &WhiteHouse},
	}
	Forest3 = Object{
		In:       &Rooms,
		LongDesc: "This is a dimly lit forest, with large trees all around.",
		Desc:     "Forest",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:    FlgLand | FlgOn | FlgSacred,
		Global:   []*Object{&Tree, &Songbird, &WhiteHouse, &Forest},
	}
	Path = Object{
		In:       &Rooms,
		LongDesc: "This is a path winding through a dimly lit forest. The path heads north-south here. One particularly large tree with some low branches stands at the edge of the path.",
		Desc:     "Forest Path",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:    FlgLand | FlgOn | FlgSacred,
		Global:   []*Object{&Tree, &Songbird, &WhiteHouse, &Forest},
	}
	UpATree = Object{
		In:     &Rooms,
		Desc:   "Up a Tree",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgLand | FlgOn | FlgSacred,
		Global: []*Object{&Tree, &Forest, &Songbird, &WhiteHouse},
	}
	GratingClearing = Object{
		In:     &Rooms,
		Desc:   "Clearing",
		Flags:  FlgLand | FlgOn | FlgSacred,
		Global: []*Object{&WhiteHouse, &Grate},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Clearing = Object{
		In:       &Rooms,
		LongDesc: "You are in a small clearing in a well marked forest path that extends to the east and west.",
		Desc:     "Clearing",
		Flags:    FlgLand | FlgOn | FlgSacred,
		Global:   []*Object{&Tree, &Songbird, &WhiteHouse, &Forest},
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// ================================================================
	// ROOMS - House
	// ================================================================

	Kitchen = Object{
		In:     &Rooms,
		Desc:   "Kitchen",
		Flags:  FlgLand | FlgOn | FlgSacred,
		Value:  10,
		Global: []*Object{&KitchenWindow, &Chimney, &Stairs},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Attic = Object{
		In:       &Rooms,
		LongDesc: "This is the attic. The only exit is a stairway leading down.",
		Desc:     "Attic",
		Flags:    FlgLand | FlgSacred,
		Global:   []*Object{&Stairs},
	}
	LivingRoom = Object{
		In:     &Rooms,
		Desc:   "Living Room",
		Flags:  FlgLand | FlgOn | FlgSacred,
		Global: []*Object{&Stairs},
		Pseudo: []PseudoObj{
			{Synonym: "nails", Action: NailsPseudo},
			{Synonym: "nail", Action: NailsPseudo},
		},
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// ================================================================
	// ROOMS - Cellar and Vicinity
	// ================================================================

	Cellar = Object{
		In:     &Rooms,
		Desc:   "Cellar",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgLand,
		Value:  25,
		Global: []*Object{&Slide, &Stairs},
	}
	TrollRoom = Object{
		In:       &Rooms,
		LongDesc: "This is a small room with passages to the east and south and a forbidding hole leading west. Bloodstains and deep scratches (perhaps made by an axe) mar the walls.",
		Desc:     "The Troll Room",
		Flags:    FlgLand,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	EastOfChasm = Object{
		In:       &Rooms,
		LongDesc: "You are on the east edge of a chasm, the bottom of which cannot be seen. A narrow passage goes north, and the path you are on continues to the east.",
		Desc:     "East of Chasm",
		Flags:    FlgLand,
		Pseudo:   []PseudoObj{{Synonym: "chasm", Action: ChasmPseudo}},
	}
	Gallery = Object{
		In:       &Rooms,
		LongDesc: "This is an art gallery. Most of the paintings have been stolen by vandals with exceptional taste. The vandals left through either the north or west exits.",
		Desc:     "Gallery",
		Flags:    FlgLand | FlgOn,
	}
	Studio = Object{
		In:       &Rooms,
		LongDesc: "This appears to have been an artist's studio. The walls and floors are splattered with paints of 69 different colors. Strangely enough, nothing of value is hanging here. At the south end of the room is an open door (also covered with paint). A dark and narrow chimney leads up from a fireplace; although you might be able to get up it, it seems unlikely you could get back down.",
		Desc:     "Studio",
		Flags:    FlgLand,
		Global:   []*Object{&Chimney},
		Pseudo: []PseudoObj{
			{Synonym: "door", Action: DoorPseudo},
			{Synonym: "paint", Action: PaintPseudo},
		},
	}
)
