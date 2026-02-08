package zork

// Conditional exit flags and game state globals
var (
	TrollFlag     bool
	CyclopsFlag   bool
	MagicFlag     bool
	LowTide       bool
	DomeFlag      bool
	EmptyHanded   bool
	LLDFlag       bool
	RainbowFlag   bool
	DeflateFlag   bool
	CoffinCure    bool
	GrateRevealed bool

	// Action-specific globals
	KitchenWindowFlag bool
	CageTop           = true
	RugMoved          bool
	GrUnlock          bool
	CycloWrath        int
	MirrorMung        bool
	GateFlag          bool
	GatesOpen         bool
	WaterLevel        int
	MatchCount        = 6
	EggSolve          bool
	ThiefHere         bool
	ThiefEngrossed    bool
	LoudFlag          bool
	SingSong          bool
	BuoyFlag          = true
	BeachDig          = -1
	LightShaft        = 13
	LampTableIdx      int
	CandleTableIdx    int
	XB                bool
	XC                bool
	Deflate           bool

	// RndSelect tables
	Dummy = RndSelect{
		Unselected: []string{
			"It is already closed.",
			"It is already open.",
			"It's already done.",
		},
	}
	SwimYuks = RndSelect{
		Unselected: []string{
			"I don't really feel like swimming.",
			"Swimming isn't usually allowed in dungeons.",
			"You'd need a submarine to go further.",
		},
	}
	BatDrops = []*Object{&Mine1, &Mine2, &Mine3, &Mine4, &LadderTop, &LadderBottom, &SqueekyRoom, &MineEntrance}

	Cyclomad = []string{
		"The cyclops seems somewhat agitated.",
		"The cyclops appears to be getting more agitated.",
		"The cyclops is moving about the room, looking for something.",
		"The cyclops was looking for salt and pepper. No doubt they are condiments for his upcoming snack.",
		"The cyclops is moving toward you in an unfriendly manner.",
		"You have two choices: 1. Leave  2. Become dinner.",
	}

	LoudRuns = []*Object{&DampCave, &RoundRoom, &DeepCanyon}

	BDigs = []string{
		"You seem to be digging a hole here.",
		"The hole is getting deeper, but that's about it.",
		"You are surrounded by a wall of sand on all sides.",
	}

	Drownings = []string{
		"up to your ankles.",
		"up to your shin.",
		"up to your knees.",
		"up to your hips.",
		"up to your waist.",
		"up to your chest.",
		"up to your neck.",
		"over your head.",
		"high in your lungs.",
	}

	RobberCDesc = "There is a suspicious-looking individual, holding a bag, leaning against one wall. He is armed with a vicious-looking stiletto."
	RobberUDesc = "There is a suspicious-looking individual lying unconscious on the ground."

	// River tables
	RiverSpeeds = []*Object{&River1, nil, &River2, nil, &River3, nil, &River4, nil, &River5, nil}
	RiverNext   = []*Object{&River1, &River2, &River2, &River3, &River3, &River4, &River4, &River5}
	RiverLaunch = []*Object{
		&DamBase, &River1,
		&WhiteCliffsNorth, &River3,
		&WhiteCliffsSouth, &River4,
		&Shore, &River5,
		&SandyBeach, &River4,
		&ReservoirSouth, &Reservoir,
		&ReservoirNorth, &Reservoir,
		&StreamView, &InStream,
	}

	// Lamp/Candle countdown tables
	LampTable = []interface{}{
		100, "The lamp appears a bit dimmer.",
		70, "The lamp is definitely dimmer now.",
		15, "The lamp is nearly out.",
		0,
	}
	CandleTable = []interface{}{
		20, "The candles grow shorter.",
		10, "The candles are becoming quite short.",
		5, "The candles won't last long now.",
		0,
	}

	ScoreMax = 350
)

var (
	// ================================================================
	// GLOBAL OBJECTS (In: &LocalGlobals or &GlobalObjects)
	// ================================================================

	Board = Object{
		In:       &LocalGlobals,
		Synonyms: []string{"boards", "board"},
		Desc:     "board",
		Flags:    []Flag{FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Teeth = Object{
		In:       &GlobalObjects,
		Synonyms: []string{"overboard", "teeth"},
		Desc:     "set of teeth",
		Flags:    []Flag{FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Wall = Object{
		In:         &GlobalObjects,
		Synonyms:   []string{"wall", "walls"},
		Adjectives: []string{"surrounding"},
		Desc:       "surrounding wall",
	}
	GraniteWall = Object{
		In:         &GlobalObjects,
		Synonyms:   []string{"wall"},
		Adjectives: []string{"granite"},
		Desc:       "granite wall",
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Songbird = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"bird", "songbird"},
		Adjectives: []string{"song"},
		Desc:       "songbird",
		Flags:      []Flag{FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	WhiteHouse = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"house"},
		Adjectives: []string{"white", "beautiful", "colonial"},
		Desc:       "white house",
		Flags:      []Flag{FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Forest = Object{
		In:       &LocalGlobals,
		Synonyms: []string{"forest", "trees", "pines", "hemlocks"},
		Desc:     "forest",
		Flags:    []Flag{FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Tree = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"tree", "branch"},
		Adjectives: []string{"large", "storm"},
		Desc:       "tree",
		Flags:      []Flag{FlgNoDesc, FlgClimb},
	}
	GlobalWater = Object{
		In:       &LocalGlobals,
		Synonyms: []string{"water", "quantity"},
		Desc:     "water",
		Flags:    []Flag{FlgDrink},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	KitchenWindow = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"window"},
		Adjectives: []string{"kitchen", "small"},
		Desc:       "kitchen window",
		Flags:      []Flag{FlgDoor, FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Chimney = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"chimney"},
		Adjectives: []string{"dark", "narrow"},
		Desc:       "chimney",
		Flags:      []Flag{FlgClimb, FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Slide = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"chute", "ramp", "slide"},
		Adjectives: []string{"steep", "metal", "twisting"},
		Desc:       "chute",
		Flags:      []Flag{FlgClimb},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Bodies = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"bodies", "body", "remains", "pile"},
		Adjectives: []string{"mangled"},
		Desc:       "pile of bodies",
		Flags:      []Flag{FlgNoDesc, FlgTryTake},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Crack = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"crack"},
		Adjectives: []string{"narrow"},
		Desc:       "crack",
		Flags:      []Flag{FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Grate = Object{
		In:       &LocalGlobals,
		Synonyms: []string{"grate", "grating"},
		Desc:     "grating",
		Flags:    []Flag{FlgDoor, FlgNoDesc, FlgInvis},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Ladder = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"ladder"},
		Adjectives: []string{"wooden", "rickety", "narrow"},
		Desc:       "wooden ladder",
		Flags:      []Flag{FlgNoDesc, FlgClimb},
	}
	ClimbableCliff = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"wall", "cliff", "walls", "ledge"},
		Adjectives: []string{"rocky", "sheer"},
		Desc:       "cliff",
		Flags:      []Flag{FlgNoDesc, FlgClimb},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	WhiteCliff = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"cliff", "cliffs"},
		Adjectives: []string{"white"},
		Desc:       "white cliffs",
		Flags:      []Flag{FlgNoDesc, FlgClimb},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Rainbow = Object{
		In:       &LocalGlobals,
		Synonyms: []string{"rainbow"},
		Desc:     "rainbow",
		Flags:    []Flag{FlgNoDesc, FlgClimb},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	River = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"river"},
		Adjectives: []string{"frigid"},
		Desc:       "river",
		Flags:      []Flag{FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	BoardedWindow = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"window"},
		Adjectives: []string{"boarded"},
		Desc:       "boarded window",
		Flags:      []Flag{FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// ================================================================
	// UNPLACED OBJECTS (created or swapped during gameplay)
	// ================================================================

	InflatedBoat = Object{
		Synonyms:   []string{"boat", "raft"},
		Adjectives: []string{"magic", "plastic", "seaworthy", "inflated", "inflatable"},
		Desc:       "magic boat",
		Flags:      []Flag{FlgTake, FlgBurn, FlgVeh, FlgOpen, FlgSearch},
		Capacity:   100,
		// Action set in FinalizeGameObjects to avoid init cycle
		Size:       20,
		VehType:    FlgNonLand,
	}
	PuncturedBoat = Object{
		Synonyms:   []string{"boat", "pile", "plastic"},
		Adjectives: []string{"plastic", "puncture", "large"},
		Desc:       "punctured boat",
		Flags:      []Flag{FlgTake, FlgBurn},
		// Action set in FinalizeGameObjects to avoid init cycle
		Size:       20,
	}
	BrokenLamp = Object{
		Synonyms:   []string{"lamp", "lantern"},
		Adjectives: []string{"broken"},
		Desc:       "broken lantern",
		Flags:      []Flag{FlgTake},
	}
	Gunk = Object{
		Synonyms:   []string{"gunk", "piece", "slag"},
		Adjectives: []string{"small", "vitreous"},
		Desc:       "small piece of vitreous slag",
		Flags:      []Flag{FlgTake, FlgTryTake},
		Size:       10,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	HotBell = Object{
		Synonyms:   []string{"bell"},
		Adjectives: []string{"brass", "hot", "red", "small"},
		Desc:       "red hot brass bell",
		Flags:      []Flag{FlgTryTake},
		LongDesc:   "On the ground is a red hot bell.",
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	BrokenEgg = Object{
		Synonyms:   []string{"egg", "treasure"},
		Adjectives: []string{"broken", "birds", "encrusted", "jewel"},
		Desc:       "broken jewel-encrusted egg",
		Flags:      []Flag{FlgTake, FlgCont, FlgOpen},
		Capacity:   6,
		TValue:     2,
		LongDesc:   "There is a somewhat ruined egg here.",
	}
	Bauble = Object{
		Synonyms:   []string{"bauble", "treasure"},
		Adjectives: []string{"brass", "beautiful"},
		Desc:       "beautiful brass bauble",
		Flags:      []Flag{FlgTake},
		Value:      1,
		TValue:     1,
	}
	Diamond = Object{
		Synonyms:   []string{"diamond", "treasure"},
		Adjectives: []string{"huge", "enormous"},
		Desc:       "huge diamond",
		Flags:      []Flag{FlgTake},
		LongDesc:   "There is an enormous diamond (perfectly cut) here.",
		Value:      10,
		TValue:     10,
	}

	// ================================================================
	// ROOMS - Forest and Outside
	// ================================================================

	WestOfHouse = Object{
		In:     &Rooms,
		Desc:   "West of House",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  []Flag{FlgLand, FlgOn, FlgSacred},
		Global: []*Object{&WhiteHouse, &Board, &Forest},
	}
	StoneBarrow = Object{
		In:       &Rooms,
		LongDesc: "You are standing in front of a massive barrow of stone. In the east face is a huge stone door which is open. You cannot see into the dark of the tomb.",
		Desc:     "Stone Barrow",
		Flags:    []Flag{FlgLand, FlgOn, FlgSacred},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	NorthOfHouse = Object{
		In:       &Rooms,
		LongDesc: "You are facing the north side of a white house. There is no door here, and all the windows are boarded up. To the north a narrow path winds through the trees.",
		Desc:     "North of House",
		Flags:    []Flag{FlgLand, FlgOn, FlgSacred},
		Global:   []*Object{&BoardedWindow, &Board, &WhiteHouse, &Forest},
	}
	SouthOfHouse = Object{
		In:       &Rooms,
		LongDesc: "You are facing the south side of a white house. There is no door here, and all the windows are boarded.",
		Desc:     "South of House",
		Flags:    []Flag{FlgLand, FlgOn, FlgSacred},
		Global:   []*Object{&BoardedWindow, &Board, &WhiteHouse, &Forest},
	}
	EastOfHouse = Object{
		In:     &Rooms,
		Desc:   "Behind House",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  []Flag{FlgLand, FlgOn, FlgSacred},
		Global: []*Object{&WhiteHouse, &KitchenWindow, &Forest},
	}
	Forest1 = Object{
		In:       &Rooms,
		LongDesc: "This is a forest, with trees in all directions. To the east, there appears to be sunlight.",
		Desc:     "Forest",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:    []Flag{FlgLand, FlgOn, FlgSacred},
		Global:   []*Object{&Tree, &Songbird, &WhiteHouse, &Forest},
	}
	Forest2 = Object{
		In:       &Rooms,
		LongDesc: "This is a dimly lit forest, with large trees all around.",
		Desc:     "Forest",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:    []Flag{FlgLand, FlgOn, FlgSacred},
		Global:   []*Object{&Tree, &Songbird, &WhiteHouse, &Forest},
	}
	Mountains = Object{
		In:       &Rooms,
		LongDesc: "The forest thins out, revealing impassable mountains.",
		Desc:     "Forest",
		Flags:    []Flag{FlgLand, FlgOn, FlgSacred},
		Global:   []*Object{&Tree, &WhiteHouse},
	}
	Forest3 = Object{
		In:       &Rooms,
		LongDesc: "This is a dimly lit forest, with large trees all around.",
		Desc:     "Forest",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:    []Flag{FlgLand, FlgOn, FlgSacred},
		Global:   []*Object{&Tree, &Songbird, &WhiteHouse, &Forest},
	}
	Path = Object{
		In:       &Rooms,
		LongDesc: "This is a path winding through a dimly lit forest. The path heads north-south here. One particularly large tree with some low branches stands at the edge of the path.",
		Desc:     "Forest Path",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:    []Flag{FlgLand, FlgOn, FlgSacred},
		Global:   []*Object{&Tree, &Songbird, &WhiteHouse, &Forest},
	}
	UpATree = Object{
		In:     &Rooms,
		Desc:   "Up a Tree",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  []Flag{FlgLand, FlgOn, FlgSacred},
		Global: []*Object{&Tree, &Forest, &Songbird, &WhiteHouse},
	}
	GratingClearing = Object{
		In:     &Rooms,
		Desc:   "Clearing",
		Flags:  []Flag{FlgLand, FlgOn, FlgSacred},
		Global: []*Object{&WhiteHouse, &Grate},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Clearing = Object{
		In:       &Rooms,
		LongDesc: "You are in a small clearing in a well marked forest path that extends to the east and west.",
		Desc:     "Clearing",
		Flags:    []Flag{FlgLand, FlgOn, FlgSacred},
		Global:   []*Object{&Tree, &Songbird, &WhiteHouse, &Forest},
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// ================================================================
	// ROOMS - House
	// ================================================================

	Kitchen = Object{
		In:     &Rooms,
		Desc:   "Kitchen",
		Flags:  []Flag{FlgLand, FlgOn, FlgSacred},
		Value:  10,
		Global: []*Object{&KitchenWindow, &Chimney, &Stairs},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Attic = Object{
		In:       &Rooms,
		LongDesc: "This is the attic. The only exit is a stairway leading down.",
		Desc:     "Attic",
		Flags:    []Flag{FlgLand, FlgSacred},
		Global:   []*Object{&Stairs},
	}
	LivingRoom = Object{
		In:     &Rooms,
		Desc:   "Living Room",
		Flags:  []Flag{FlgLand, FlgOn, FlgSacred},
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
		Flags:  []Flag{FlgLand},
		Value:  25,
		Global: []*Object{&Slide, &Stairs},
	}
	TrollRoom = Object{
		In:       &Rooms,
		LongDesc: "This is a small room with passages to the east and south and a forbidding hole leading west. Bloodstains and deep scratches (perhaps made by an axe) mar the walls.",
		Desc:     "The Troll Room",
		Flags:    []Flag{FlgLand},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	EastOfChasm = Object{
		In:       &Rooms,
		LongDesc: "You are on the east edge of a chasm, the bottom of which cannot be seen. A narrow passage goes north, and the path you are on continues to the east.",
		Desc:     "East of Chasm",
		Flags:    []Flag{FlgLand},
		Pseudo:   []PseudoObj{{Synonym: "chasm", Action: ChasmPseudo}},
	}
	Gallery = Object{
		In:       &Rooms,
		LongDesc: "This is an art gallery. Most of the paintings have been stolen by vandals with exceptional taste. The vandals left through either the north or west exits.",
		Desc:     "Gallery",
		Flags:    []Flag{FlgLand, FlgOn},
	}
	Studio = Object{
		In:       &Rooms,
		LongDesc: "This appears to have been an artist's studio. The walls and floors are splattered with paints of 69 different colors. Strangely enough, nothing of value is hanging here. At the south end of the room is an open door (also covered with paint). A dark and narrow chimney leads up from a fireplace; although you might be able to get up it, it seems unlikely you could get back down.",
		Desc:     "Studio",
		Flags:    []Flag{FlgLand},
		Global:   []*Object{&Chimney},
		Pseudo: []PseudoObj{
			{Synonym: "door", Action: DoorPseudo},
			{Synonym: "paint", Action: PaintPseudo},
		},
	}

	// ================================================================
	// ROOMS - Maze
	// ================================================================

	Maze1 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    []Flag{FlgLand, FlgMaze},
	}
	Maze2 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    []Flag{FlgLand, FlgMaze},
	}
	Maze3 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    []Flag{FlgLand, FlgMaze},
	}
	Maze4 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    []Flag{FlgLand, FlgMaze},
	}
	DeadEnd1 = Object{
		In:       &Rooms,
		LongDesc: "You have come to a dead end in the maze.",
		Desc:     "Dead End",
		Flags:    []Flag{FlgLand, FlgMaze},
	}
	Maze5 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike. A skeleton, probably the remains of a luckless adventurer, lies here.",
		Desc:     "Maze",
		Flags:    []Flag{FlgLand, FlgMaze},
	}
	DeadEnd2 = Object{
		In:       &Rooms,
		LongDesc: "You have come to a dead end in the maze.",
		Desc:     "Dead End",
		Flags:    []Flag{FlgLand, FlgMaze},
	}
	Maze6 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    []Flag{FlgLand, FlgMaze},
	}
	Maze7 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    []Flag{FlgLand, FlgMaze},
	}
	Maze8 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    []Flag{FlgLand, FlgMaze},
	}
	DeadEnd3 = Object{
		In:       &Rooms,
		LongDesc: "You have come to a dead end in the maze.",
		Desc:     "Dead End",
		Flags:    []Flag{FlgLand, FlgMaze},
	}
	Maze9 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    []Flag{FlgLand, FlgMaze},
	}
	Maze10 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    []Flag{FlgLand, FlgMaze},
	}
	Maze11 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    []Flag{FlgLand, FlgMaze},
	}
	GratingRoom = Object{
		In:     &Rooms,
		Desc:   "Grating Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  []Flag{FlgLand},
		Global: []*Object{&Grate},
	}
	Maze12 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    []Flag{FlgLand, FlgMaze},
	}
	DeadEnd4 = Object{
		In:       &Rooms,
		LongDesc: "You have come to a dead end in the maze.",
		Desc:     "Dead End",
		Flags:    []Flag{FlgLand, FlgMaze},
	}
	Maze13 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    []Flag{FlgLand, FlgMaze},
	}
	Maze14 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    []Flag{FlgLand, FlgMaze},
	}
	Maze15 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    []Flag{FlgLand, FlgMaze},
	}

	// ================================================================
	// ROOMS - Cyclops and Hideaway
	// ================================================================

	CyclopsRoom = Object{
		In:     &Rooms,
		Desc:   "Cyclops Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  []Flag{FlgLand},
		Global: []*Object{&Stairs},
	}
	StrangePassage = Object{
		In:       &Rooms,
		LongDesc: "This is a long passage. To the west is one entrance. On the east there is an old wooden door, with a large opening in it (about cyclops sized).",
		Desc:     "Strange Passage",
		Flags:    []Flag{FlgLand},
	}
	TreasureRoom = Object{
		In:       &Rooms,
		LongDesc: "This is a large room, whose east wall is solid granite. A number of discarded bags, which crumble at your touch, are scattered about on the floor. There is an exit down a staircase.",
		Desc:     "Treasure Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:    []Flag{FlgLand},
		Value:    25,
		Global:   []*Object{&Stairs},
	}

	// ================================================================
	// ROOMS - Reservoir Area
	// ================================================================

	ReservoirSouth = Object{
		In:     &Rooms,
		Desc:   "Reservoir South",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  []Flag{FlgLand},
		Global: []*Object{&GlobalWater},
		Pseudo: []PseudoObj{
			{Synonym: "lake", Action: LakePseudo},
			{Synonym: "chasm", Action: ChasmPseudo},
		},
	}
	Reservoir = Object{
		In:     &Rooms,
		Desc:   "Reservoir",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  []Flag{FlgNonLand},
		Global: []*Object{&GlobalWater},
		Pseudo: []PseudoObj{{Synonym: "stream", Action: StreamPseudo}},
	}
	ReservoirNorth = Object{
		In:     &Rooms,
		Desc:   "Reservoir North",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  []Flag{FlgLand},
		Global: []*Object{&GlobalWater, &Stairs},
		Pseudo: []PseudoObj{{Synonym: "lake", Action: LakePseudo}},
	}
	StreamView = Object{
		In:       &Rooms,
		LongDesc: "You are standing on a path beside a gently flowing stream. The path follows the stream, which flows from west to east.",
		Desc:     "Stream View",
		Flags:    []Flag{FlgLand},
		Global:   []*Object{&GlobalWater},
		Pseudo:   []PseudoObj{{Synonym: "stream", Action: StreamPseudo}},
	}
	InStream = Object{
		In:       &Rooms,
		LongDesc: "You are on the gently flowing stream. The upstream route is too narrow to navigate, and the downstream route is invisible due to twisting walls. There is a narrow beach to land on.",
		Desc:     "Stream",
		Flags:    []Flag{FlgNonLand},
		Global:   []*Object{&GlobalWater},
		Pseudo:   []PseudoObj{{Synonym: "stream", Action: StreamPseudo}},
	}

	// ================================================================
	// ROOMS - Mirror Rooms and Vicinity
	// ================================================================

	MirrorRoom1 = Object{
		In:     &Rooms,
		Desc:   "Mirror Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  []Flag{FlgLand},
	}
	MirrorRoom2 = Object{
		In:     &Rooms,
		Desc:   "Mirror Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  []Flag{FlgLand, FlgOn},
	}
	SmallCave = Object{
		In:       &Rooms,
		LongDesc: "This is a tiny cave with entrances west and north, and a staircase leading down.",
		Desc:     "Cave",
		Flags:    []Flag{FlgLand},
		Global:   []*Object{&Stairs},
	}
	TinyCave = Object{
		In:       &Rooms,
		LongDesc: "This is a tiny cave with entrances west and north, and a dark, forbidding staircase leading down.",
		Desc:     "Cave",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:    []Flag{FlgLand},
		Global:   []*Object{&Stairs},
	}
	ColdPassage = Object{
		In:       &Rooms,
		LongDesc: "This is a cold and damp corridor where a long east-west passageway turns into a southward path.",
		Desc:     "Cold Passage",
		Flags:    []Flag{FlgLand},
	}
	NarrowPassage = Object{
		In:       &Rooms,
		LongDesc: "This is a long and narrow corridor where a long north-south passageway briefly narrows even further.",
		Desc:     "Narrow Passage",
		Flags:    []Flag{FlgLand},
	}
	WindingPassage = Object{
		In:       &Rooms,
		LongDesc: "This is a winding passage. It seems that there are only exits on the east and north.",
		Desc:     "Winding Passage",
		Flags:    []Flag{FlgLand},
	}
	TwistingPassage = Object{
		In:       &Rooms,
		LongDesc: "This is a winding passage. It seems that there are only exits on the east and north.",
		Desc:     "Twisting Passage",
		Flags:    []Flag{FlgLand},
	}
	AtlantisRoom = Object{
		In:       &Rooms,
		LongDesc: "This is an ancient room, long under water. There is an exit to the south and a staircase leading up.",
		Desc:     "Atlantis Room",
		Flags:    []Flag{FlgLand},
		Global:   []*Object{&Stairs},
	}

	// ================================================================
	// ROOMS - Round Room and Vicinity
	// ================================================================

	EWPassage = Object{
		In:       &Rooms,
		LongDesc: "This is a narrow east-west passageway. There is a narrow stairway leading down at the north end of the room.",
		Desc:     "East-West Passage",
		Flags:    []Flag{FlgLand},
		Value:    5,
		Global:   []*Object{&Stairs},
	}
	RoundRoom = Object{
		In:       &Rooms,
		LongDesc: "This is a circular stone room with passages in all directions. Several of them have unfortunately been blocked by cave-ins.",
		Desc:     "Round Room",
		Flags:    []Flag{FlgLand},
	}
	DeepCanyon = Object{
		In:     &Rooms,
		Desc:   "Deep Canyon",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  []Flag{FlgLand},
		Global: []*Object{&Stairs},
	}
	DampCave = Object{
		In:       &Rooms,
		LongDesc: "This cave has exits to the west and east, and narrows to a crack toward the south. The earth is particularly damp here.",
		Desc:     "Damp Cave",
		Flags:    []Flag{FlgLand},
		Global:   []*Object{&Crack},
	}
	LoudRoom = Object{
		In:     &Rooms,
		Desc:   "Loud Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  []Flag{FlgLand},
		Global: []*Object{&Stairs},
	}
	NSPassage = Object{
		In:       &Rooms,
		LongDesc: "This is a high north-south passage, which forks to the northeast.",
		Desc:     "North-South Passage",
		Flags:    []Flag{FlgLand},
	}
	ChasmRoom = Object{
		In:       &Rooms,
		LongDesc: "A chasm runs southwest to northeast and the path follows it. You are on the south side of the chasm, where a crack opens into a passage.",
		Desc:     "Chasm",
		Flags:    []Flag{FlgLand},
		Global:   []*Object{&Crack, &Stairs},
		Pseudo:   []PseudoObj{{Synonym: "chasm", Action: ChasmPseudo}},
	}

	// ================================================================
	// ROOMS - Hades
	// ================================================================

	EnteranceToHades = Object{
		In:     &Rooms,
		Desc:   "Entrance to Hades",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  []Flag{FlgLand, FlgOn},
		Global: []*Object{&Bodies},
		Pseudo: []PseudoObj{
			{Synonym: "gate", Action: GatePseudo},
			{Synonym: "gates", Action: GatePseudo},
		},
	}
	LandOfLivingDead = Object{
		In:       &Rooms,
		LongDesc: "You have entered the Land of the Living Dead. Thousands of lost souls can be heard weeping and moaning. In the corner are stacked the remains of dozens of previous adventurers less fortunate than yourself. A passage exits to the north.",
		Desc:     "Land of the Dead",
		Flags:    []Flag{FlgLand, FlgOn},
		Global:   []*Object{&Bodies},
	}

	// ================================================================
	// ROOMS - Dome, Temple, Egypt
	// ================================================================

	EngravingsCave = Object{
		In:       &Rooms,
		LongDesc: "You have entered a low cave with passages leading northwest and east.",
		Desc:     "Engravings Cave",
		Flags:    []Flag{FlgLand},
	}
	EgyptRoom = Object{
		In:       &Rooms,
		LongDesc: "This is a room which looks like an Egyptian tomb. There is an ascending staircase to the west.",
		Desc:     "Egyptian Room",
		Flags:    []Flag{FlgLand},
		Global:   []*Object{&Stairs},
	}
	DomeRoom = Object{
		In:     &Rooms,
		Desc:   "Dome Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  []Flag{FlgLand},
		Pseudo: []PseudoObj{{Synonym: "dome", Action: DomePseudo}},
	}
	TorchRoom = Object{
		In:     &Rooms,
		Desc:   "Torch Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  []Flag{FlgLand},
		Global: []*Object{&Stairs},
		Pseudo: []PseudoObj{{Synonym: "dome", Action: DomePseudo}},
	}
	NorthTemple = Object{
		In:       &Rooms,
		LongDesc: "This is the north end of a large temple. On the east wall is an ancient inscription, probably a prayer in a long-forgotten language. Below the prayer is a staircase leading down. The west wall is solid granite. The exit to the north end of the room is through huge marble pillars.",
		Desc:     "Temple",
		Flags:    []Flag{FlgLand, FlgOn, FlgSacred},
		Global:   []*Object{&Stairs},
	}
	SouthTemple = Object{
		In:       &Rooms,
		LongDesc: "This is the south end of a large temple. In front of you is what appears to be an altar. In one corner is a small hole in the floor which leads into darkness. You probably could not get back up it.",
		Desc:     "Altar",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:    []Flag{FlgLand, FlgOn, FlgSacred},
	}

	// ================================================================
	// ROOMS - Flood Control Dam #3
	// ================================================================

	DamRoom = Object{
		In:     &Rooms,
		Desc:   "Dam",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  []Flag{FlgLand, FlgOn},
		Global: []*Object{&GlobalWater},
	}
	DamLobby = Object{
		In:       &Rooms,
		LongDesc: "This room appears to have been the waiting room for groups touring the dam. There are open doorways here to the north and east marked \"Private\", and there is a path leading south over the top of the dam.",
		Desc:     "Dam Lobby",
		Flags:    []Flag{FlgLand, FlgOn},
	}
	MaintenanceRoom = Object{
		In:       &Rooms,
		LongDesc: "This is what appears to have been the maintenance room for Flood Control Dam #3. Apparently, this room has been ransacked recently, for most of the valuable equipment is gone. On the wall in front of you is a group of buttons colored blue, yellow, brown, and red. There are doorways to the west and south.",
		Desc:     "Maintenance Room",
		Flags:    []Flag{FlgLand},
	}

	// ================================================================
	// ROOMS - River Area
	// ================================================================

	DamBase = Object{
		In:       &Rooms,
		LongDesc: "You are at the base of Flood Control Dam #3, which looms above you and to the north. The river Frigid is flowing by here. Along the river are the White Cliffs which seem to form giant walls stretching from north to south along the shores of the river as it winds its way downstream.",
		Desc:     "Dam Base",
		Flags:    []Flag{FlgLand, FlgOn, FlgSacred},
		Global:   []*Object{&GlobalWater, &River},
	}
	River1 = Object{
		In:       &Rooms,
		LongDesc: "You are on the Frigid River in the vicinity of the Dam. The river flows quietly here. There is a landing on the west shore.",
		Desc:     "Frigid River",
		Flags:    []Flag{FlgNonLand, FlgSacred, FlgOn},
		Global:   []*Object{&GlobalWater, &River},
	}
	River2 = Object{
		In:       &Rooms,
		LongDesc: "The river turns a corner here making it impossible to see the Dam. The White Cliffs loom on the east bank and large rocks prevent landing on the west.",
		Desc:     "Frigid River",
		Flags:    []Flag{FlgNonLand, FlgSacred},
		Global:   []*Object{&GlobalWater, &River},
	}
	River3 = Object{
		In:       &Rooms,
		LongDesc: "The river descends here into a valley. There is a narrow beach on the west shore below the cliffs. In the distance a faint rumbling can be heard.",
		Desc:     "Frigid River",
		Flags:    []Flag{FlgNonLand, FlgSacred},
		Global:   []*Object{&GlobalWater, &River},
	}
	WhiteCliffsNorth = Object{
		In:       &Rooms,
		LongDesc: "You are on a narrow strip of beach which runs along the base of the White Cliffs. There is a narrow path heading south along the Cliffs and a tight passage leading west into the cliffs themselves.",
		Desc:     "White Cliffs Beach",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:    []Flag{FlgLand, FlgSacred},
		Global:   []*Object{&GlobalWater, &WhiteCliff, &River},
	}
	WhiteCliffsSouth = Object{
		In:       &Rooms,
		LongDesc: "You are on a rocky, narrow strip of beach beside the Cliffs. A narrow path leads north along the shore.",
		Desc:     "White Cliffs Beach",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:    []Flag{FlgLand, FlgSacred},
		Global:   []*Object{&GlobalWater, &WhiteCliff, &River},
	}
	River4 = Object{
		In:       &Rooms,
		LongDesc: "The river is running faster here and the sound ahead appears to be that of rushing water. On the east shore is a sandy beach. A small area of beach can also be seen below the cliffs on the west shore.",
		Desc:     "Frigid River",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:    []Flag{FlgNonLand, FlgSacred},
		Global:   []*Object{&GlobalWater, &River},
	}
	River5 = Object{
		In:       &Rooms,
		LongDesc: "The sound of rushing water is nearly unbearable here. On the east shore is a large landing area.",
		Desc:     "Frigid River",
		Flags:    []Flag{FlgNonLand, FlgSacred, FlgOn},
		Global:   []*Object{&GlobalWater, &River},
	}
	Shore = Object{
		In:       &Rooms,
		LongDesc: "You are on the east shore of the river. The water here seems somewhat treacherous. A path travels from north to south here, the south end quickly turning around a sharp corner.",
		Desc:     "Shore",
		Flags:    []Flag{FlgLand, FlgSacred, FlgOn},
		Global:   []*Object{&GlobalWater, &River},
	}
	SandyBeach = Object{
		In:       &Rooms,
		LongDesc: "You are on a large sandy beach on the east shore of the river, which is flowing quickly by. A path runs beside the river to the south here, and a passage is partially buried in sand to the northeast.",
		Desc:     "Sandy Beach",
		Flags:    []Flag{FlgLand, FlgSacred},
		Global:   []*Object{&GlobalWater, &River},
	}
	SandyCave = Object{
		In:       &Rooms,
		LongDesc: "This is a sand-filled cave whose exit is to the southwest.",
		Desc:     "Sandy Cave",
		Flags:    []Flag{FlgLand},
	}
	AragainFalls = Object{
		In:     &Rooms,
		Desc:   "Aragain Falls",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  []Flag{FlgLand, FlgSacred, FlgOn},
		Global: []*Object{&GlobalWater, &River, &Rainbow},
	}
	OnRainbow = Object{
		In:       &Rooms,
		LongDesc: "You are on top of a rainbow (I bet you never thought you would walk on a rainbow), with a magnificent view of the Falls. The rainbow travels east-west here.",
		Desc:     "On the Rainbow",
		Flags:    []Flag{FlgLand, FlgOn, FlgSacred},
		Global:   []*Object{&Rainbow},
	}
	EndOfRainbow = Object{
		In:       &Rooms,
		LongDesc: "You are on a small, rocky beach on the continuation of the Frigid River past the Falls. The beach is narrow due to the presence of the White Cliffs. The river canyon opens here and sunlight shines in from above. A rainbow crosses over the falls to the east and a narrow path continues to the southwest.",
		Desc:     "End of Rainbow",
		Flags:    []Flag{FlgLand, FlgOn},
		Global:   []*Object{&GlobalWater, &Rainbow, &River},
	}
	CanyonBottom = Object{
		In:       &Rooms,
		LongDesc: "You are beneath the walls of the river canyon which may be climbable here. The lesser part of the runoff of Aragain Falls flows by below. To the north is a narrow path.",
		Desc:     "Canyon Bottom",
		Flags:    []Flag{FlgLand, FlgOn, FlgSacred},
		Global:   []*Object{&GlobalWater, &ClimbableCliff, &River},
	}
	CliffMiddle = Object{
		In:       &Rooms,
		LongDesc: "You are on a ledge about halfway up the wall of the river canyon. You can see from here that the main flow from Aragain Falls twists along a passage which it is impossible for you to enter. Below you is the canyon bottom. Above you is more cliff, which appears climbable.",
		Desc:     "Rocky Ledge",
		Flags:    []Flag{FlgLand, FlgOn, FlgSacred},
		Global:   []*Object{&ClimbableCliff, &River},
	}
	CanyonView = Object{
		In:       &Rooms,
		LongDesc: "You are at the top of the Great Canyon on its west wall. From here there is a marvelous view of the canyon and parts of the Frigid River upstream. Across the canyon, the walls of the White Cliffs join the mighty ramparts of the Flathead Mountains to the east. Following the Canyon upstream to the north, Aragain Falls may be seen, complete with rainbow. The mighty Frigid River flows out from a great dark cavern. To the west and south can be seen an immense forest, stretching for miles around. A path leads northwest. It is possible to climb down into the canyon from here.",
		Desc:     "Canyon View",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:    []Flag{FlgLand, FlgOn, FlgSacred},
		Global:   []*Object{&ClimbableCliff, &River, &Rainbow},
	}

	// ================================================================
	// ROOMS - Coal Mine Area
	// ================================================================

	MineEntrance = Object{
		In:       &Rooms,
		LongDesc: "You are standing at the entrance of what might have been a coal mine. The shaft enters the west wall, and there is another exit on the south end of the room.",
		Desc:     "Mine Entrance",
		Flags:    []Flag{FlgLand},
	}
	SqueekyRoom = Object{
		In:       &Rooms,
		LongDesc: "You are in a small room. Strange squeaky sounds may be heard coming from the passage at the north end. You may also escape to the east.",
		Desc:     "Squeaky Room",
		Flags:    []Flag{FlgLand},
	}
	BatRoom = Object{
		In:     &Rooms,
		Desc:   "Bat Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  []Flag{FlgLand, FlgSacred},
	}
	ShaftRoom = Object{
		In:       &Rooms,
		LongDesc: "This is a large room, in the middle of which is a small shaft descending through the floor into darkness below. To the west and the north are exits from this room. Constructed over the top of the shaft is a metal framework to which a heavy iron chain is attached.",
		Desc:     "Shaft Room",
		Flags:    []Flag{FlgLand},
		Pseudo:   []PseudoObj{{Synonym: "chain", Action: ChainPseudo}},
	}
	SmellyRoom = Object{
		In:       &Rooms,
		LongDesc: "This is a small nondescript room. However, from the direction of a small descending staircase a foul odor can be detected. To the south is a narrow tunnel.",
		Desc:     "Smelly Room",
		Flags:    []Flag{FlgLand},
		Global:   []*Object{&Stairs},
		Pseudo: []PseudoObj{
			{Synonym: "odor", Action: GasPseudo},
			{Synonym: "gas", Action: GasPseudo},
		},
	}
	GasRoom = Object{
		In:       &Rooms,
		LongDesc: "This is a small room which smells strongly of coal gas. There is a short climb up some stairs and a narrow tunnel leading east.",
		Desc:     "Gas Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:    []Flag{FlgLand, FlgSacred},
		Global:   []*Object{&Stairs},
		Pseudo: []PseudoObj{
			{Synonym: "gas", Action: GasPseudo},
			{Synonym: "odor", Action: GasPseudo},
		},
	}
	LadderTop = Object{
		In:       &Rooms,
		LongDesc: "This is a very small room. In the corner is a rickety wooden ladder, leading downward. It might be safe to descend. There is also a staircase leading upward.",
		Desc:     "Ladder Top",
		Flags:    []Flag{FlgLand},
		Global:   []*Object{&Ladder, &Stairs},
	}
	LadderBottom = Object{
		In:       &Rooms,
		LongDesc: "This is a rather wide room. On one side is the bottom of a narrow wooden ladder. To the west and the south are passages leaving the room.",
		Desc:     "Ladder Bottom",
		Flags:    []Flag{FlgLand},
		Global:   []*Object{&Ladder},
	}
	DeadEnd5 = Object{
		In:       &Rooms,
		LongDesc: "You have come to a dead end in the mine.",
		Desc:     "Dead End",
		Flags:    []Flag{FlgLand},
	}
	TimberRoom = Object{
		In:       &Rooms,
		LongDesc: "This is a long and narrow passage, which is cluttered with broken timbers. A wide passage comes from the east and turns at the west end of the room into a very narrow passageway. From the west comes a strong draft.",
		Desc:     "Timber Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:    []Flag{FlgLand, FlgSacred},
	}
	LowerShaft = Object{
		In:       &Rooms,
		LongDesc: "This is a small drafty room in which is the bottom of a long shaft. To the south is a passageway and to the east a very narrow passage. In the shaft can be seen a heavy iron chain.",
		Desc:     "Drafty Room",
		Flags:    []Flag{FlgLand, FlgSacred},
		Pseudo:   []PseudoObj{{Synonym: "chain", Action: ChainPseudo}},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	MachineRoom = Object{
		In:     &Rooms,
		Desc:   "Machine Room",
		Flags:  []Flag{FlgLand},
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// ================================================================
	// ROOMS - Coal Mine
	// ================================================================

	Mine1 = Object{
		In:       &Rooms,
		LongDesc: "This is a nondescript part of a coal mine.",
		Desc:     "Coal Mine",
		Flags:    []Flag{FlgLand},
	}
	Mine2 = Object{
		In:       &Rooms,
		LongDesc: "This is a nondescript part of a coal mine.",
		Desc:     "Coal Mine",
		Flags:    []Flag{FlgLand},
	}
	Mine3 = Object{
		In:       &Rooms,
		LongDesc: "This is a nondescript part of a coal mine.",
		Desc:     "Coal Mine",
		Flags:    []Flag{FlgLand},
	}
	Mine4 = Object{
		In:       &Rooms,
		LongDesc: "This is a nondescript part of a coal mine.",
		Desc:     "Coal Mine",
		Flags:    []Flag{FlgLand},
	}
	SlideRoom = Object{
		In:       &Rooms,
		LongDesc: "This is a small chamber, which appears to have been part of a coal mine. On the south wall of the chamber the letters \"Granite Wall\" are etched in the rock. To the east is a long passage, and there is a steep metal slide twisting downward. To the north is a small opening.",
		Desc:     "Slide Room",
		Flags:    []Flag{FlgLand},
		Global:   []*Object{&Slide},
	}

	// ================================================================
	// OBJECTS IN ROOMS
	// ================================================================

	// Mountains
	MountainRange = Object{
		In:         &Mountains,
		Synonyms:   []string{"mountain", "range"},
		Adjectives: []string{"impassable", "flathead"},
		Desc:       "mountain range",
		Flags:      []Flag{FlgNoDesc, FlgClimb},
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// West of House
	FrontDoor = Object{
		In:         &WestOfHouse,
		Synonyms:   []string{"door"},
		Adjectives: []string{"front", "boarded"},
		Desc:       "door",
		Flags:      []Flag{FlgDoor, FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Mailbox = Object{
		In:         &WestOfHouse,
		Synonyms:   []string{"mailbox", "box"},
		Adjectives: []string{"small"},
		Desc:       "small mailbox",
		Flags:      []Flag{FlgCont, FlgTryTake},
		Capacity:   10,
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// Stone Barrow
	BarrowDoor = Object{
		In:         &StoneBarrow,
		Synonyms:   []string{"door"},
		Adjectives: []string{"huge", "stone"},
		Desc:       "stone door",
		Flags:      []Flag{FlgDoor, FlgNoDesc, FlgOpen},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Barrow = Object{
		In:         &StoneBarrow,
		Synonyms:   []string{"barrow", "tomb"},
		Adjectives: []string{"massive", "stone"},
		Desc:       "stone barrow",
		Flags:      []Flag{FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// Living Room
	TrophyCase = Object{
		In:         &LivingRoom,
		Synonyms:   []string{"case"},
		Adjectives: []string{"trophy"},
		Desc:       "trophy case",
		Flags:      []Flag{FlgTrans, FlgCont, FlgNoDesc, FlgTryTake, FlgSearch},
		Capacity:   10000,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Rug = Object{
		In:         &LivingRoom,
		Synonyms:   []string{"rug", "carpet"},
		Adjectives: []string{"large", "oriental"},
		Desc:       "carpet",
		Flags:      []Flag{FlgNoDesc, FlgTryTake},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	TrapDoor = Object{
		In:         &LivingRoom,
		Synonyms:   []string{"door", "trapdoor", "trap-door", "cover"},
		Adjectives: []string{"trap", "dusty"},
		Desc:       "trap door",
		Flags:      []Flag{FlgDoor, FlgNoDesc, FlgInvis},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	WoodenDoor = Object{
		In:         &LivingRoom,
		Synonyms:   []string{"door", "lettering", "writing"},
		Adjectives: []string{"wooden", "gothic", "strange", "west"},
		Desc:       "wooden door",
		Flags:      []Flag{FlgRead, FlgDoor, FlgNoDesc, FlgTrans},
		// Action set in FinalizeGameObjects to avoid init cycle
		Text:       "The engravings translate to \"This space intentionally left blank.\"",
	}
	Sword = Object{
		In:         &LivingRoom,
		Synonyms:   []string{"sword", "orcrist", "glamdring", "blade"},
		Adjectives: []string{"elvish", "old", "antique"},
		Desc:       "sword",
		Flags:      []Flag{FlgTake, FlgWeapon, FlgTryTake},
		FirstDesc:  "Above the trophy case hangs an elvish sword of great antiquity.",
		// Action set in FinalizeGameObjects to avoid init cycle
		Size:       30,
		TValue:     0,
	}
	Lamp = Object{
		In:         &LivingRoom,
		Synonyms:   []string{"lamp", "lantern", "light"},
		Adjectives: []string{"brass"},
		Desc:       "brass lantern",
		Flags:      []Flag{FlgTake, FlgLight},
		FirstDesc:  "A battery-powered brass lantern is on the trophy case.",
		// Action set in FinalizeGameObjects to avoid init cycle
		LongDesc:   "There is a brass lantern (battery-powered) here.",
		Size:       15,
	}

	// Kitchen
	KitchenTable = Object{
		In:         &Kitchen,
		Synonyms:   []string{"table"},
		Adjectives: []string{"kitchen"},
		Desc:       "kitchen table",
		Flags:      []Flag{FlgNoDesc, FlgCont, FlgOpen, FlgSurf},
		Capacity:   50,
	}

	// Attic
	AtticTable = Object{
		In:       &Attic,
		Synonyms: []string{"table"},
		Desc:     "table",
		Flags:    []Flag{FlgNoDesc, FlgCont, FlgOpen, FlgSurf},
		Capacity: 40,
	}
	Rope = Object{
		In:         &Attic,
		Synonyms:   []string{"rope", "hemp", "coil"},
		Adjectives: []string{"large"},
		Desc:       "rope",
		Flags:      []Flag{FlgTake, FlgSacred, FlgTryTake},
		// Action set in FinalizeGameObjects to avoid init cycle
		FirstDesc:  "A large coil of rope is lying in the corner.",
		Size:       10,
	}

	// Entrance to Hades
	Ghosts = Object{
		In:         &EnteranceToHades,
		Synonyms:   []string{"ghosts", "spirits", "fiends", "force"},
		Adjectives: []string{"invisible", "evil"},
		Desc:       "number of ghosts",
		Flags:      []Flag{FlgPerson, FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// Land of Living Dead
	Skull = Object{
		In:         &LandOfLivingDead,
		Synonyms:   []string{"skull", "head", "treasure"},
		Adjectives: []string{"crystal"},
		Desc:       "crystal skull",
		FirstDesc:  "Lying in one corner of the room is a beautifully carved crystal skull. It appears to be grinning at you rather nastily.",
		Flags:      []Flag{FlgTake},
		Value:      10,
		TValue:     10,
	}

	// Shaft Room
	RaisedBasket = Object{
		In:         &ShaftRoom,
		Synonyms:   []string{"cage", "dumbwaiter", "basket"},
		Desc:       "basket",
		Flags:      []Flag{FlgTrans, FlgTryTake, FlgCont, FlgOpen},
		// Action set in FinalizeGameObjects to avoid init cycle
		LongDesc:   "At the end of the chain is a basket.",
		Capacity:   50,
	}

	// Lower Shaft
	LoweredBasket = Object{
		In:         &LowerShaft,
		Synonyms:   []string{"cage", "dumbwaiter", "basket"},
		Adjectives: []string{"lowered"},
		Desc:       "basket",
		LongDesc:   "From the chain is suspended a basket.",
		Flags:      []Flag{FlgTryTake},
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// Bat Room
	Bat = Object{
		In:         &BatRoom,
		Synonyms:   []string{"bat", "vampire"},
		Adjectives: []string{"vampire", "deranged"},
		Desc:       "bat",
		Flags:      []Flag{FlgPerson, FlgTryTake},
		DescFcn:    BatDescFcn,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Jade = Object{
		In:         &BatRoom,
		Synonyms:   []string{"figurine", "treasure"},
		Adjectives: []string{"exquisite", "jade"},
		Desc:       "jade figurine",
		Flags:      []Flag{FlgTake},
		LongDesc:   "There is an exquisite jade figurine here.",
		Size:       10,
		Value:      5,
		TValue:     5,
	}

	// North Temple
	Bell = Object{
		In:         &NorthTemple,
		Synonyms:   []string{"bell"},
		Adjectives: []string{"small", "brass"},
		Desc:       "brass bell",
		Flags:      []Flag{FlgTake},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Prayer = Object{
		In:         &NorthTemple,
		Synonyms:   []string{"prayer", "inscription"},
		Adjectives: []string{"ancient", "old"},
		Desc:       "prayer",
		Flags:      []Flag{FlgRead, FlgSacred, FlgNoDesc},
		Text:       "The prayer is inscribed in an ancient script, rarely used today. It seems to be a philippic against small insects, absent-mindedness, and the picking up and dropping of small objects. The final verse consigns trespassers to the land of the dead. All evidence indicates that the beliefs of the ancient Zorkers were obscure.",
	}

	// South Temple
	Altar = Object{
		In:       &SouthTemple,
		Synonyms: []string{"altar"},
		Desc:     "altar",
		Flags:    []Flag{FlgNoDesc, FlgSurf, FlgCont, FlgOpen},
		Capacity: 50,
	}
	Candles = Object{
		In:         &SouthTemple,
		Synonyms:   []string{"candles", "pair"},
		Adjectives: []string{"burning"},
		Desc:       "pair of candles",
		Flags:      []Flag{FlgTake, FlgFlame, FlgOn, FlgLight},
		// Action set in FinalizeGameObjects to avoid init cycle
		FirstDesc:  "On the two ends of the altar are burning candles.",
		Size:       10,
	}

	// Troll Room
	Troll = Object{
		In:         &TrollRoom,
		Synonyms:   []string{"troll"},
		Adjectives: []string{"nasty"},
		Desc:       "troll",
		Flags:      []Flag{FlgPerson, FlgOpen, FlgTryTake},
		// Action set in FinalizeGameObjects to avoid init cycle
		LongDesc:   "A nasty-looking troll, brandishing a bloody axe, blocks all passages out of the room.",
		Strength:   2,
	}

	// Dam Room
	Bolt = Object{
		In:         &DamRoom,
		Synonyms:   []string{"bolt", "nut"},
		Adjectives: []string{"metal", "large"},
		Desc:       "bolt",
		Flags:      []Flag{FlgNoDesc, FlgTurn, FlgTryTake},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Bubble = Object{
		In:         &DamRoom,
		Synonyms:   []string{"bubble"},
		Adjectives: []string{"small", "green", "plastic"},
		Desc:       "green bubble",
		Flags:      []Flag{FlgNoDesc, FlgTryTake},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Dam = Object{
		In:         &DamRoom,
		Synonyms:   []string{"dam", "gate", "gates", "fcd#3"},
		Desc:       "dam",
		Flags:      []Flag{FlgNoDesc, FlgTryTake},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	ControlPanel = Object{
		In:         &DamRoom,
		Synonyms:   []string{"panel"},
		Adjectives: []string{"control"},
		Desc:       "control panel",
		Flags:      []Flag{FlgNoDesc},
	}

	// Dam Lobby
	Match = Object{
		In:         &DamLobby,
		Synonyms:   []string{"match", "matches", "matchbook"},
		Adjectives: []string{"match"},
		Desc:       "matchbook",
		Flags:      []Flag{FlgRead, FlgTake},
		LongDesc:   "There is a matchbook whose cover says \"Visit Beautiful FCD#3\" here.",
		// Action set in FinalizeGameObjects to avoid init cycle
		Size:       2,
		Text:       "\n(Close cover before striking)\n\nYOU too can make BIG MONEY in the exciting field of PAPER SHUFFLING!\n\nMr. Anderson of Muddle, Mass. says: \"Before I took this course I was a lowly bit twiddler. Now with what I learned at GUE Tech I feel really important and can obfuscate and confuse with the best.\"\n\nDr. Blank had this to say: \"Ten short days ago all I could look forward to was a dead-end job as a doctor. Now I have a promising future and make really big Zorkmids.\"\n\nGUE Tech can't promise these fantastic results to everyone. But when you earn your degree from GUE Tech, your future will be brighter.",
	}
	Guide = Object{
		In:         &DamLobby,
		Synonyms:   []string{"guide", "book", "books", "guidebooks"},
		Adjectives: []string{"tour", "guide"},
		Desc:       "tour guidebook",
		Flags:      []Flag{FlgRead, FlgTake, FlgBurn},
		FirstDesc:  "Some guidebooks entitled \"Flood Control Dam #3\" are on the reception desk.",
		Text:       "\"\tFlood Control Dam #3\n\nFCD#3 was constructed in year 783 of the Great Underground Empire to harness the mighty Frigid River. This work was supported by a grant of 37 million zorkmids from your omnipotent local tyrant Lord Dimwit Flathead the Excessive. This impressive structure is composed of 370,000 cubic feet of concrete, is 256 feet tall at the center, and 193 feet wide at the top. The lake created behind the dam has a volume of 1.7 billion cubic feet, an area of 12 million square feet, and a shore line of 36 thousand feet.\n\nThe construction of FCD#3 took 112 days from ground breaking to the dedication. It required a work force of 384 slaves, 34 slave drivers, 12 engineers, 2 turtle doves, and a partridge in a pear tree. The work was managed by a command team composed of 2345 bureaucrats, 2347 secretaries (at least two of whom could type), 12,256 paper shufflers, 52,469 rubber stampers, 245,193 red tape processors, and nearly one million dead trees.\n\nWe will now point out some of the more interesting features of FCD#3 as we conduct you on a guided tour of the facilities:\n\n        1) You start your tour here in the Dam Lobby. You will notice on your right that....\"",
	}

	// Dam Base
	InflatableBoat = Object{
		In:         &DamBase,
		Synonyms:   []string{"boat", "pile", "plastic", "valve"},
		Adjectives: []string{"plastic", "inflatable"},
		Desc:       "pile of plastic",
		Flags:      []Flag{FlgTake, FlgBurn},
		// Action set in FinalizeGameObjects to avoid init cycle
		LongDesc:   "There is a folded pile of plastic here which has a small valve attached.",
		Size:       20,
	}

	// Maintenance Room
	ToolChest = Object{
		In:         &MaintenanceRoom,
		Synonyms:   []string{"chest", "chests", "group", "toolchests"},
		Adjectives: []string{"tool"},
		Desc:       "group of tool chests",
		Flags:      []Flag{FlgCont, FlgOpen, FlgTryTake, FlgSacred},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	YellowButton = Object{
		In:         &MaintenanceRoom,
		Synonyms:   []string{"button", "switch"},
		Adjectives: []string{"yellow"},
		Desc:       "yellow button",
		Flags:      []Flag{FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	BrownButton = Object{
		In:         &MaintenanceRoom,
		Synonyms:   []string{"button", "switch"},
		Adjectives: []string{"brown"},
		Desc:       "brown button",
		Flags:      []Flag{FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	RedButton = Object{
		In:         &MaintenanceRoom,
		Synonyms:   []string{"button", "switch"},
		Adjectives: []string{"red"},
		Desc:       "red button",
		Flags:      []Flag{FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	BlueButton = Object{
		In:         &MaintenanceRoom,
		Synonyms:   []string{"button", "switch"},
		Adjectives: []string{"blue"},
		Desc:       "blue button",
		Flags:      []Flag{FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Screwdriver = Object{
		In:         &MaintenanceRoom,
		Synonyms:   []string{"screwdriver", "tool", "tools", "driver"},
		Adjectives: []string{"screw"},
		Desc:       "screwdriver",
		Flags:      []Flag{FlgTake, FlgTool},
	}
	Wrench = Object{
		In:       &MaintenanceRoom,
		Synonyms: []string{"wrench", "tool", "tools"},
		Desc:     "wrench",
		Flags:    []Flag{FlgTake, FlgTool},
		Size:     10,
	}
	Tube = Object{
		In:       &MaintenanceRoom,
		Synonyms: []string{"tube", "tooth", "paste"},
		Desc:     "tube",
		Flags:    []Flag{FlgTake, FlgCont, FlgRead},
		// Action set in FinalizeGameObjects to avoid init cycle
		LongDesc: "There is an object which looks like a tube of toothpaste here.",
		Capacity: 7,
		Size:     5,
		Text:     "---> Frobozz Magic Gunk Company <---\n\tAll-Purpose Gunk",
	}
	Leak = Object{
		In:       &MaintenanceRoom,
		Synonyms: []string{"leak", "drip", "pipe"},
		Desc:     "leak",
		Flags:    []Flag{FlgNoDesc, FlgInvis},
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// Machine Room
	Machine = Object{
		In:       &MachineRoom,
		Synonyms: []string{"machine", "pdp10", "dryer", "lid"},
		Desc:     "machine",
		Flags:    []Flag{FlgCont, FlgNoDesc, FlgTryTake},
		Capacity: 50,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	MachineSwitch = Object{
		In:       &MachineRoom,
		Synonyms: []string{"switch"},
		Desc:     "switch",
		Flags:    []Flag{FlgNoDesc, FlgTurn},
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// Cyclops Room
	Cyclops = Object{
		In:         &CyclopsRoom,
		Synonyms:   []string{"cyclops", "monster", "eye"},
		Adjectives: []string{"hungry", "giant"},
		Desc:       "cyclops",
		Flags:      []Flag{FlgPerson, FlgNoDesc, FlgTryTake},
		// Action set in FinalizeGameObjects to avoid init cycle
		Strength:   10000,
	}

	// Treasure Room
	Chalice = Object{
		In:         &TreasureRoom,
		Synonyms:   []string{"chalice", "cup", "silver", "treasure"},
		Adjectives: []string{"silver", "engravings"},
		Desc:       "chalice",
		Flags:      []Flag{FlgTake, FlgTryTake, FlgCont},
		// Action set in FinalizeGameObjects to avoid init cycle
		LongDesc:   "There is a silver chalice, intricately engraved, here.",
		Capacity:   5,
		Size:       10,
		Value:      10,
		TValue:     5,
	}

	// Gallery
	Painting = Object{
		In:         &Gallery,
		Synonyms:   []string{"painting", "art", "canvas", "treasure"},
		Adjectives: []string{"beautiful"},
		Desc:       "painting",
		Flags:      []Flag{FlgTake, FlgBurn},
		// Action set in FinalizeGameObjects to avoid init cycle
		FirstDesc:  "Fortunately, there is still one chance for you to be a vandal, for on the far wall is a painting of unparalleled beauty.",
		LongDesc:   "A painting by a neglected genius is here.",
		Size:       15,
		Value:      4,
		TValue:     6,
	}

	// Studio
	OwnersManual = Object{
		In:         &Studio,
		Synonyms:   []string{"manual", "piece", "paper"},
		Adjectives: []string{"zork", "owners", "small"},
		Desc:       "ZORK owner's manual",
		Flags:      []Flag{FlgRead, FlgTake},
		FirstDesc:  "Loosely attached to a wall is a small piece of paper.",
		Text:       "Congratulations!\n\nYou are the privileged owner of ZORK I: The Great Underground Empire, a self-contained and self-maintaining universe. If used and maintained in accordance with normal operating practices for small universes, ZORK will provide many months of trouble-free operation.",
	}

	// Grating Clearing
	Leaves = Object{
		In:       &GratingClearing,
		Synonyms: []string{"leaves", "leaf", "pile"},
		Desc:     "pile of leaves",
		Flags:    []Flag{FlgTake, FlgBurn, FlgTryTake},
		// Action set in FinalizeGameObjects to avoid init cycle
		LongDesc: "On the ground is a pile of leaves.",
		Size:     25,
	}

	// Up a Tree
	Nest = Object{
		In:         &UpATree,
		Synonyms:   []string{"nest"},
		Adjectives: []string{"birds"},
		Desc:       "bird's nest",
		Flags:      []Flag{FlgTake, FlgBurn, FlgCont, FlgOpen, FlgSearch},
		FirstDesc:  "Beside you on the branch is a small bird's nest.",
		Capacity:   20,
	}

	// Sandy Cave
	Sand = Object{
		In:       &SandyCave,
		Synonyms: []string{"sand"},
		Desc:     "sand",
		Flags:    []Flag{FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Scarab = Object{
		In:         &SandyCave,
		Synonyms:   []string{"scarab", "bug", "beetle", "treasure"},
		Adjectives: []string{"beautiful", "carved", "jeweled"},
		Desc:       "beautiful jeweled scarab",
		Flags:      []Flag{FlgTake, FlgInvis},
		Size:       8,
		Value:      5,
		TValue:     5,
	}

	// Sandy Beach
	Shovel = Object{
		In:       &SandyBeach,
		Synonyms: []string{"shovel", "tool", "tools"},
		Desc:     "shovel",
		Flags:    []Flag{FlgTake, FlgTool},
		Size:     15,
	}

	// Egypt Room
	Coffin = Object{
		In:         &EgyptRoom,
		Synonyms:   []string{"coffin", "casket", "treasure"},
		Adjectives: []string{"solid", "gold"},
		Desc:       "gold coffin",
		Flags:      []Flag{FlgTake, FlgCont, FlgSacred, FlgSearch},
		LongDesc:   "The solid-gold coffin used for the burial of Ramses II is here.",
		Capacity:   35,
		Size:       55,
		Value:      10,
		TValue:     15,
	}

	// Round Room
	Thief = Object{
		In:         &RoundRoom,
		Synonyms:   []string{"thief", "robber", "man", "person"},
		Adjectives: []string{"shady", "suspicious", "seedy"},
		Desc:       "thief",
		Flags:      []Flag{FlgPerson, FlgInvis, FlgCont, FlgOpen, FlgTryTake},
		// Action set in FinalizeGameObjects to avoid init cycle
		LongDesc:   "There is a suspicious-looking individual, holding a large bag, leaning against one wall. He is armed with a deadly stiletto.",
		Strength:   5,
	}

	// Reservoir
	Trunk = Object{
		In:         &Reservoir,
		Synonyms:   []string{"trunk", "chest", "jewels", "treasure"},
		Adjectives: []string{"old"},
		Desc:       "trunk of jewels",
		Flags:      []Flag{FlgTake, FlgInvis},
		FirstDesc:  "Lying half buried in the mud is an old trunk, bulging with jewels.",
		LongDesc:   "There is an old trunk here, bulging with assorted jewels.",
		// Action set in FinalizeGameObjects to avoid init cycle
		Size:       35,
		Value:      15,
		TValue:     5,
	}

	// Reservoir North
	Pump = Object{
		In:         &ReservoirNorth,
		Synonyms:   []string{"pump", "air-pump", "tool", "tools"},
		Adjectives: []string{"small", "hand-held"},
		Desc:       "hand-held air pump",
		Flags:      []Flag{FlgTake, FlgTool},
	}

	// Atlantis Room
	Trident = Object{
		In:         &AtlantisRoom,
		Synonyms:   []string{"trident", "fork", "treasure"},
		Adjectives: []string{"poseidon", "own", "crystal"},
		Desc:       "crystal trident",
		Flags:      []Flag{FlgTake},
		FirstDesc:  "On the shore lies Poseidon's own crystal trident.",
		Size:       20,
		Value:      4,
		TValue:     11,
	}

	// Mirror rooms
	Mirror1 = Object{
		In:         &MirrorRoom1,
		Synonyms:   []string{"reflection", "mirror", "enormous"},
		Desc:       "mirror",
		Flags:      []Flag{FlgTryTake, FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Mirror2 = Object{
		In:         &MirrorRoom2,
		Synonyms:   []string{"reflection", "mirror", "enormous"},
		Desc:       "mirror",
		Flags:      []Flag{FlgTryTake, FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// Dome Room
	Railing = Object{
		In:         &DomeRoom,
		Synonyms:   []string{"railing", "rail"},
		Adjectives: []string{"wooden"},
		Desc:       "wooden railing",
		Flags:      []Flag{FlgNoDesc},
	}

	// Torch Room
	Pedestal = Object{
		In:         &TorchRoom,
		Synonyms:   []string{"pedestal"},
		Adjectives: []string{"white", "marble"},
		Desc:       "pedestal",
		Flags:      []Flag{FlgNoDesc, FlgCont, FlgOpen, FlgSurf},
		// Action set in FinalizeGameObjects to avoid init cycle
		Capacity:   30,
	}

	// Engravings Cave
	Engravings = Object{
		In:       &EngravingsCave,
		Synonyms: []string{"wall", "engravings", "inscription"},
		Adjectives: []string{"old", "ancient"},
		Desc:     "wall with engravings",
		Flags:    []Flag{FlgRead, FlgSacred},
		LongDesc: "There are old engravings on the walls here.",
		Text:     "The engravings were incised in the living rock of the cave wall by an unknown hand. They depict, in symbolic form, the beliefs of the ancient Zorkers. Skillfully interwoven with the bas reliefs are excerpts illustrating the major religious tenets of that time. Unfortunately, a later age seems to have considered them blasphemous and just as skillfully excised them.",
	}

	// Loud Room
	Bar = Object{
		In:         &LoudRoom,
		Synonyms:   []string{"bar", "platinum", "treasure"},
		Adjectives: []string{"platinum", "large"},
		Desc:       "platinum bar",
		Flags:      []Flag{FlgTake, FlgSacred},
		LongDesc:   "On the ground is a large platinum bar.",
		Size:       20,
		Value:      10,
		TValue:     5,
	}

	// End of Rainbow
	PotOfGold = Object{
		In:         &EndOfRainbow,
		Synonyms:   []string{"pot", "gold", "treasure"},
		Adjectives: []string{"gold"},
		Desc:       "pot of gold",
		Flags:      []Flag{FlgTake, FlgInvis},
		FirstDesc:  "At the end of the rainbow is a pot of gold.",
		Size:       15,
		Value:      10,
		TValue:     10,
	}

	// River 4
	Buoy = Object{
		In:         &River4,
		Synonyms:   []string{"buoy"},
		Adjectives: []string{"red"},
		Desc:       "red buoy",
		Flags:      []Flag{FlgTake, FlgCont},
		FirstDesc:  "There is a red buoy here (probably a warning).",
		Capacity:   20,
		Size:       10,
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// Gas Room
	Bracelet = Object{
		In:         &GasRoom,
		Synonyms:   []string{"bracelet", "jewel", "sapphire", "treasure"},
		Adjectives: []string{"sapphire"},
		Desc:       "sapphire-encrusted bracelet",
		Flags:      []Flag{FlgTake},
		Size:       10,
		Value:      5,
		TValue:     5,
	}

	// Dead End 5
	Coal = Object{
		In:         &DeadEnd5,
		Synonyms:   []string{"coal", "pile", "heap"},
		Adjectives: []string{"small"},
		Desc:       "small pile of coal",
		Flags:      []Flag{FlgTake, FlgBurn},
		Size:       20,
	}

	// Timber Room
	Timbers = Object{
		In:         &TimberRoom,
		Synonyms:   []string{"timbers", "pile"},
		Adjectives: []string{"wooden", "broken"},
		Desc:       "broken timber",
		Flags:      []Flag{FlgTake},
		Size:       50,
	}

	// Maze 5
	Bones = Object{
		In:         &Maze5,
		Synonyms:   []string{"bones", "skeleton", "body"},
		Desc:       "skeleton",
		Flags:      []Flag{FlgTryTake, FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	BurnedOutLantern = Object{
		In:         &Maze5,
		Synonyms:   []string{"lantern", "lamp"},
		Adjectives: []string{"rusty", "burned", "dead", "useless"},
		Desc:       "burned-out lantern",
		Flags:      []Flag{FlgTake},
		FirstDesc:  "The deceased adventurer's useless lantern is here.",
		Size:       20,
	}
	BagOfCoins = Object{
		In:         &Maze5,
		Synonyms:   []string{"bag", "coins", "treasure"},
		Adjectives: []string{"old", "leather"},
		Desc:       "leather bag of coins",
		Flags:      []Flag{FlgTake},
		LongDesc:   "An old leather bag, bulging with coins, is here.",
		// Action set in FinalizeGameObjects to avoid init cycle
		Size:       15,
		Value:      10,
		TValue:     5,
	}
	RustyKnife = Object{
		In:         &Maze5,
		Synonyms:   []string{"knives", "knife"},
		Adjectives: []string{"rusty"},
		Desc:       "rusty knife",
		Flags:      []Flag{FlgTake, FlgTryTake, FlgWeapon, FlgTool},
		// Action set in FinalizeGameObjects to avoid init cycle
		FirstDesc:  "Beside the skeleton is a rusty knife.",
		Size:       20,
	}
	Keys = Object{
		In:         &Maze5,
		Synonyms:   []string{"key"},
		Adjectives: []string{"skeleton"},
		Desc:       "skeleton key",
		Flags:      []Flag{FlgTake, FlgTool},
		Size:       10,
	}

	// ================================================================
	// OBJECTS IN OBJECTS
	// ================================================================

	// In Trophy Case
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

	// In Mailbox
	Advertisement = Object{
		In:         &Mailbox,
		Synonyms:   []string{"advertisement", "leaflet", "booklet", "mail"},
		Adjectives: []string{"small"},
		Desc:       "leaflet",
		Flags:      []Flag{FlgRead, FlgTake, FlgBurn},
		LongDesc:   "A small leaflet is on the ground.",
		Size:       2,
		Text:       "\"WELCOME TO ZORK!\n\nZORK is a game of adventure, danger, and low cunning. In it you will explore some of the most amazing territory ever seen by mortals. No computer should be without one!\"",
	}

	// In Kitchen Table
	Bottle = Object{
		In:         &KitchenTable,
		Synonyms:   []string{"bottle", "container"},
		Adjectives: []string{"clear", "glass"},
		Desc:       "glass bottle",
		Flags:      []Flag{FlgTake, FlgTrans, FlgCont},
		// Action set in FinalizeGameObjects to avoid init cycle
		FirstDesc:  "A bottle is sitting on the table.",
		Capacity:   4,
	}
	SandwichBag = Object{
		In:         &KitchenTable,
		Synonyms:   []string{"bag", "sack"},
		Adjectives: []string{"brown", "elongated", "smelly"},
		Desc:       "brown sack",
		Flags:      []Flag{FlgTake, FlgCont, FlgBurn},
		FirstDesc:  "On the table is an elongated brown sack, smelling of hot peppers.",
		Capacity:   9,
		Size:       9,
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// In Attic Table
	Knife = Object{
		In:         &AtticTable,
		Synonyms:   []string{"knives", "knife", "blade"},
		Adjectives: []string{"nasty", "unrusty"},
		Desc:       "nasty knife",
		Flags:      []Flag{FlgTake, FlgWeapon, FlgTryTake},
		FirstDesc:  "On a table is a nasty-looking knife.",
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// In Bottle
	Water = Object{
		In:         &Bottle,
		Synonyms:   []string{"water", "quantity", "liquid", "h2o"},
		Desc:       "quantity of water",
		Flags:      []Flag{FlgTryTake, FlgTake, FlgDrink},
		// Action set in FinalizeGameObjects to avoid init cycle
		Size:       4,
	}

	// In Sandwich Bag
	Lunch = Object{
		In:         &SandwichBag,
		Synonyms:   []string{"food", "sandwich", "lunch", "dinner"},
		Adjectives: []string{"hot", "pepper"},
		Desc:       "lunch",
		Flags:      []Flag{FlgTake, FlgFood},
		LongDesc:   "A hot pepper sandwich is here.",
	}
	Garlic = Object{
		In:       &SandwichBag,
		Synonyms: []string{"garlic", "clove"},
		Desc:     "clove of garlic",
		Flags:    []Flag{FlgTake, FlgFood},
		// Action set in FinalizeGameObjects to avoid init cycle
		Size:     4,
	}

	// In Altar
	Book = Object{
		In:         &Altar,
		Synonyms:   []string{"book", "prayer", "page", "books"},
		Adjectives: []string{"large", "black"},
		Desc:       "black book",
		Flags:      []Flag{FlgRead, FlgTake, FlgCont, FlgBurn, FlgTurn},
		// Action set in FinalizeGameObjects to avoid init cycle
		FirstDesc:  "On the altar is a large black book, open to page 569.",
		Size:       10,
		Text:       "Commandment #12592\n\nOh ye who go about saying unto each: \"Hello sailor\":\nDost thou know the magnitude of thy sin before the gods?\nYea, verily, thou shalt be ground between two stones.\nShall the angry gods cast thy body into the whirlpool?\nSurely, thy eye shall be put out with a sharp stick!\nEven unto the ends of the earth shalt thou wander and\nUnto the land of the dead shalt thou be sent at last.\nSurely thou shalt repent of thy cunning.",
	}

	// In Coffin
	Sceptre = Object{
		In:         &Coffin,
		Synonyms:   []string{"sceptre", "scepter", "treasure"},
		Adjectives: []string{"sharp", "egyptian", "ancient", "enameled"},
		Desc:       "sceptre",
		Flags:      []Flag{FlgTake, FlgWeapon},
		LongDesc:   "An ornamented sceptre, tapering to a sharp point, is here.",
		// Action set in FinalizeGameObjects to avoid init cycle
		FirstDesc:  "A sceptre, possibly that of ancient Egypt itself, is in the coffin. The sceptre is ornamented with colored enamel, and tapers to a sharp point.",
		Size:       3,
		Value:      4,
		TValue:     6,
	}

	// In Nest
	Egg = Object{
		In:         &Nest,
		Synonyms:   []string{"egg", "treasure"},
		Adjectives: []string{"birds", "encrusted", "jeweled"},
		Desc:       "jewel-encrusted egg",
		Flags:      []Flag{FlgTake, FlgCont, FlgSearch},
		// Action set in FinalizeGameObjects to avoid init cycle
		Value:      5,
		TValue:     5,
		Capacity:   6,
		FirstDesc:  "In the bird's nest is a large egg encrusted with precious jewels, apparently scavenged by a childless songbird. The egg is covered with fine gold inlay, and ornamented in lapis lazuli and mother-of-pearl. Unlike most eggs, this one is hinged and closed with a delicate looking clasp. The egg appears extremely fragile.",
	}

	// In Egg
	Canary = Object{
		In:         &Egg,
		Synonyms:   []string{"canary", "treasure"},
		Adjectives: []string{"clockwork", "gold", "golden"},
		Desc:       "golden clockwork canary",
		Flags:      []Flag{FlgTake, FlgSearch},
		// Action set in FinalizeGameObjects to avoid init cycle
		Value:      6,
		TValue:     4,
		FirstDesc:  "There is a golden clockwork canary nestled in the egg. It has ruby eyes and a silver beak. Through a crystal window below its left wing you can see intricate machinery inside. It appears to have wound down.",
	}

	// In Tube
	Putty = Object{
		In:         &Tube,
		Synonyms:   []string{"material", "gunk"},
		Adjectives: []string{"viscous"},
		Desc:       "viscous material",
		Flags:      []Flag{FlgTake, FlgTool},
		Size:       6,
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// In Troll
	Axe = Object{
		In:         &Troll,
		Synonyms:   []string{"axe", "ax"},
		Adjectives: []string{"bloody"},
		Desc:       "bloody axe",
		Flags:      []Flag{FlgWeapon, FlgTryTake, FlgTake, FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
		Size:       25,
	}

	// In Thief
	LargeBag = Object{
		In:         &Thief,
		Synonyms:   []string{"bag"},
		Adjectives: []string{"large", "thiefs"},
		Desc:       "large bag",
		Flags:      []Flag{FlgTryTake, FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Stiletto = Object{
		In:         &Thief,
		Synonyms:   []string{"stiletto"},
		Adjectives: []string{"vicious"},
		Desc:       "stiletto",
		Flags:      []Flag{FlgWeapon, FlgTryTake, FlgTake, FlgNoDesc},
		// Action set in FinalizeGameObjects to avoid init cycle
		Size:       10,
	}

	// In Pedestal
	Torch = Object{
		In:         &Pedestal,
		Synonyms:   []string{"torch", "ivory", "treasure"},
		Adjectives: []string{"flaming", "ivory"},
		Desc:       "torch",
		Flags:      []Flag{FlgTake, FlgFlame, FlgOn, FlgLight},
		// Action set in FinalizeGameObjects to avoid init cycle
		FirstDesc:  "Sitting on the pedestal is a flaming torch, made of ivory.",
		Size:       20,
		Value:      14,
		TValue:     6,
	}

	// In Inflated Boat
	BoatLabel = Object{
		In:         &InflatedBoat,
		Synonyms:   []string{"label", "fineprint", "print"},
		Adjectives: []string{"tan", "fine"},
		Desc:       "tan label",
		Flags:      []Flag{FlgRead, FlgTake, FlgBurn},
		Size:       2,
		Text:       "  !!!!FROBOZZ MAGIC BOAT COMPANY!!!!\n\nHello, Sailor!\n\nInstructions for use:\n\n   To get into a body of water, say \"Launch\".\n   To get to shore, say \"Land\" or the direction in which you want to maneuver the boat.\n\nWarranty:\n\n  This boat is guaranteed against all defects for a period of 76 milliseconds from date of purchase or until first used, whichever comes first.\n\nWarning:\n   This boat is made of thin plastic.\n   Good Luck!",
	}

	// In Buoy
	Emerald = Object{
		In:         &Buoy,
		Synonyms:   []string{"emerald", "treasure"},
		Adjectives: []string{"large"},
		Desc:       "large emerald",
		Flags:      []Flag{FlgTake},
		Value:      5,
		TValue:     10,
	}

	// In Broken Egg (unplaced)
	BrokenCanary = Object{
		In:         &BrokenEgg,
		Synonyms:   []string{"canary", "treasure"},
		Adjectives: []string{"broken", "clockwork", "gold", "golden"},
		Desc:       "broken clockwork canary",
		Flags:      []Flag{FlgTake},
		// Action set in FinalizeGameObjects to avoid init cycle
		TValue:     1,
		FirstDesc:  "There is a golden clockwork canary nestled in the egg. It seems to have recently had a bad experience. The mountings for its jewel-like eyes are empty, and its silver beak is crumpled. Through a cracked crystal window below its left wing you can see the remains of intricate machinery. It is not clear what result winding it would have, as the mainspring seems sprung.",
	}

	// ================================================================
	// TABLES
	// ================================================================

	InHouseAround = []*Object{
		&LivingRoom,
		&Kitchen,
		&Attic,
		&Kitchen,
	}
	HouseAround = []*Object{
		&WestOfHouse,
		&NorthOfHouse,
		&EastOfHouse,
		&SouthOfHouse,
		&WestOfHouse,
	}
	ForestAround = []*Object{
		&Forest1,
		&Forest2,
		&Forest3,
		&Path,
		&Clearing,
		&Forest1,
	}
	AboveGround = []*Object{
		&WestOfHouse,
		&NorthOfHouse,
		&EastOfHouse,
		&SouthOfHouse,
		&Forest1,
		&Forest2,
		&Forest3,
		&Path,
		&Clearing,
		&GratingClearing,
		&CanyonView,
	}

	// ================================================================
	// MASTER OBJECT LIST
	// ================================================================

	Objects = []*Object{
		// Core
		&Rooms, &GlobalObjects, &LocalGlobals,
		&NotHereObject, &PseudoObject,
		&Hands, &Me, &Adventurer, &It,
		&Stairs, &Intnum, &Blessings, &Sailor, &Ground, &Grue, &Lungs, &PathObj, &Zorkmid,
		// Global objects
		&Board, &Teeth, &Wall, &GraniteWall, &Songbird,
		&WhiteHouse, &Forest, &Tree, &GlobalWater,
		&KitchenWindow, &Chimney, &Slide, &Bodies, &Crack, &Grate,
		&Ladder, &ClimbableCliff, &WhiteCliff, &Rainbow, &River,
		&BoardedWindow,
		// Unplaced objects
		&InflatedBoat, &PuncturedBoat, &BrokenLamp, &Gunk, &HotBell,
		&BrokenEgg, &Bauble, &Diamond,
		// Forest and outside rooms
		&WestOfHouse, &StoneBarrow, &NorthOfHouse, &SouthOfHouse, &EastOfHouse,
		&Forest1, &Forest2, &Mountains, &Forest3, &Path,
		&UpATree, &GratingClearing, &Clearing,
		// House rooms
		&Kitchen, &Attic, &LivingRoom,
		// Cellar and vicinity
		&Cellar, &TrollRoom, &EastOfChasm, &Gallery, &Studio,
		// Maze
		&Maze1, &Maze2, &Maze3, &Maze4, &DeadEnd1,
		&Maze5, &DeadEnd2, &Maze6, &Maze7, &Maze8, &DeadEnd3,
		&Maze9, &Maze10, &Maze11, &GratingRoom,
		&Maze12, &DeadEnd4, &Maze13, &Maze14, &Maze15,
		// Cyclops and hideaway
		&CyclopsRoom, &StrangePassage, &TreasureRoom,
		// Reservoir area
		&ReservoirSouth, &Reservoir, &ReservoirNorth, &StreamView, &InStream,
		// Mirror rooms and vicinity
		&MirrorRoom1, &MirrorRoom2, &SmallCave, &TinyCave,
		&ColdPassage, &NarrowPassage, &WindingPassage, &TwistingPassage, &AtlantisRoom,
		// Round room and vicinity
		&EWPassage, &RoundRoom, &DeepCanyon, &DampCave, &LoudRoom, &NSPassage, &ChasmRoom,
		// Hades
		&EnteranceToHades, &LandOfLivingDead,
		// Dome, temple, egypt
		&EngravingsCave, &EgyptRoom, &DomeRoom, &TorchRoom, &NorthTemple, &SouthTemple,
		// Flood control dam
		&DamRoom, &DamLobby, &MaintenanceRoom,
		// River area
		&DamBase, &River1, &River2, &River3,
		&WhiteCliffsNorth, &WhiteCliffsSouth,
		&River4, &River5, &Shore, &SandyBeach, &SandyCave,
		&AragainFalls, &OnRainbow, &EndOfRainbow,
		&CanyonBottom, &CliffMiddle, &CanyonView,
		// Coal mine area
		&MineEntrance, &SqueekyRoom, &BatRoom, &ShaftRoom,
		&SmellyRoom, &GasRoom, &LadderTop, &LadderBottom, &DeadEnd5,
		&TimberRoom, &LowerShaft, &MachineRoom,
		&Mine1, &Mine2, &Mine3, &Mine4, &SlideRoom,
		// Objects in rooms
		&MountainRange, &FrontDoor, &Mailbox, &BarrowDoor, &Barrow,
		&TrophyCase, &Rug, &TrapDoor, &WoodenDoor, &Sword, &Lamp,
		&KitchenTable, &AtticTable, &Rope,
		&Ghosts, &Skull,
		&RaisedBasket, &LoweredBasket,
		&Bat, &Jade, &Bell, &Prayer, &Altar, &Candles,
		&Troll, &Bolt, &Bubble, &Dam, &ControlPanel,
		&Match, &Guide, &InflatableBoat,
		&ToolChest, &YellowButton, &BrownButton, &RedButton, &BlueButton,
		&Screwdriver, &Wrench, &Tube, &Leak,
		&Machine, &MachineSwitch,
		&Cyclops, &Chalice, &Painting, &OwnersManual,
		&Leaves, &Nest, &Sand, &Scarab, &Shovel,
		&Coffin, &Thief, &Trunk, &Pump, &Trident,
		&Mirror1, &Mirror2, &Railing, &Pedestal, &Engravings,
		&Bar, &PotOfGold, &Buoy, &Bracelet, &Coal, &Timbers,
		&Bones, &BurnedOutLantern, &BagOfCoins, &RustyKnife, &Keys,
		// Objects in objects
		&Map, &Advertisement, &Bottle, &SandwichBag, &Knife,
		&Water, &Lunch, &Garlic, &Book, &Sceptre,
		&Egg, &Canary, &Putty, &Axe,
		&LargeBag, &Stiletto, &Torch, &BoatLabel, &Emerald,
		&BrokenCanary,
	}
)

// InitRoomExits sets all room directional properties.
// This is done in a function to avoid circular reference issues
// between rooms that reference each other.
func InitRoomExits() {
	// West of House
	WestOfHouse.North = DirProps{UExit: true, RExit: &NorthOfHouse}
	WestOfHouse.South = DirProps{UExit: true, RExit: &SouthOfHouse}
	WestOfHouse.NorthEast = DirProps{UExit: true, RExit: &NorthOfHouse}
	WestOfHouse.SouthEast = DirProps{UExit: true, RExit: &SouthOfHouse}
	WestOfHouse.West = DirProps{UExit: true, RExit: &Forest1}
	WestOfHouse.East = DirProps{NExit: "The door is boarded and you can't remove the boards."}
	WestOfHouse.SouthWest = DirProps{CExit: func() bool { return WonGame }, RExit: &StoneBarrow}
	WestOfHouse.Into = DirProps{CExit: func() bool { return WonGame }, RExit: &StoneBarrow}

	// Stone Barrow
	StoneBarrow.NorthEast = DirProps{UExit: true, RExit: &WestOfHouse}

	// North of House
	NorthOfHouse.SouthWest = DirProps{UExit: true, RExit: &WestOfHouse}
	NorthOfHouse.SouthEast = DirProps{UExit: true, RExit: &EastOfHouse}
	NorthOfHouse.West = DirProps{UExit: true, RExit: &WestOfHouse}
	NorthOfHouse.East = DirProps{UExit: true, RExit: &EastOfHouse}
	NorthOfHouse.North = DirProps{UExit: true, RExit: &Path}
	NorthOfHouse.South = DirProps{NExit: "The windows are all boarded."}

	// South of House
	SouthOfHouse.West = DirProps{UExit: true, RExit: &WestOfHouse}
	SouthOfHouse.East = DirProps{UExit: true, RExit: &EastOfHouse}
	SouthOfHouse.NorthEast = DirProps{UExit: true, RExit: &EastOfHouse}
	SouthOfHouse.NorthWest = DirProps{UExit: true, RExit: &WestOfHouse}
	SouthOfHouse.South = DirProps{UExit: true, RExit: &Forest3}
	SouthOfHouse.North = DirProps{NExit: "The windows are all boarded."}

	// East of House
	EastOfHouse.North = DirProps{UExit: true, RExit: &NorthOfHouse}
	EastOfHouse.South = DirProps{UExit: true, RExit: &SouthOfHouse}
	EastOfHouse.SouthWest = DirProps{UExit: true, RExit: &SouthOfHouse}
	EastOfHouse.NorthWest = DirProps{UExit: true, RExit: &NorthOfHouse}
	EastOfHouse.East = DirProps{UExit: true, RExit: &Clearing}
	EastOfHouse.West = DirProps{DExit: &KitchenWindow, RExit: &Kitchen}
	EastOfHouse.Into = DirProps{DExit: &KitchenWindow, RExit: &Kitchen}

	// Forest 1
	Forest1.Up = DirProps{NExit: "There is no tree here suitable for climbing."}
	Forest1.North = DirProps{UExit: true, RExit: &GratingClearing}
	Forest1.East = DirProps{UExit: true, RExit: &Path}
	Forest1.South = DirProps{UExit: true, RExit: &Forest3}
	Forest1.West = DirProps{NExit: "You would need a machete to go further west."}

	// Forest 2
	Forest2.Up = DirProps{NExit: "There is no tree here suitable for climbing."}
	Forest2.North = DirProps{NExit: "The forest becomes impenetrable to the north."}
	Forest2.East = DirProps{UExit: true, RExit: &Mountains}
	Forest2.South = DirProps{UExit: true, RExit: &Clearing}
	Forest2.West = DirProps{UExit: true, RExit: &Path}

	// Mountains
	Mountains.Up = DirProps{NExit: "The mountains are impassable."}
	Mountains.North = DirProps{UExit: true, RExit: &Forest2}
	Mountains.East = DirProps{NExit: "The mountains are impassable."}
	Mountains.South = DirProps{UExit: true, RExit: &Forest2}
	Mountains.West = DirProps{UExit: true, RExit: &Forest2}

	// Forest 3
	Forest3.Up = DirProps{NExit: "There is no tree here suitable for climbing."}
	Forest3.North = DirProps{UExit: true, RExit: &Clearing}
	Forest3.East = DirProps{NExit: "The rank undergrowth prevents eastward movement."}
	Forest3.South = DirProps{NExit: "Storm-tossed trees block your way."}
	Forest3.West = DirProps{UExit: true, RExit: &Forest1}
	Forest3.NorthWest = DirProps{UExit: true, RExit: &SouthOfHouse}

	// Path
	Path.Up = DirProps{UExit: true, RExit: &UpATree}
	Path.North = DirProps{UExit: true, RExit: &GratingClearing}
	Path.East = DirProps{UExit: true, RExit: &Forest2}
	Path.South = DirProps{UExit: true, RExit: &NorthOfHouse}
	Path.West = DirProps{UExit: true, RExit: &Forest1}

	// Up a Tree
	UpATree.Down = DirProps{UExit: true, RExit: &Path}
	UpATree.Up = DirProps{NExit: "You cannot climb any higher."}

	// Grating Clearing
	GratingClearing.North = DirProps{NExit: "The forest becomes impenetrable to the north."}
	GratingClearing.East = DirProps{UExit: true, RExit: &Forest2}
	GratingClearing.West = DirProps{UExit: true, RExit: &Forest1}
	GratingClearing.South = DirProps{UExit: true, RExit: &Path}
	GratingClearing.Down = DirProps{FExit: GratingExitFcn}

	// Clearing
	Clearing.Up = DirProps{NExit: "There is no tree here suitable for climbing."}
	Clearing.East = DirProps{UExit: true, RExit: &CanyonView}
	Clearing.North = DirProps{UExit: true, RExit: &Forest2}
	Clearing.South = DirProps{UExit: true, RExit: &Forest3}
	Clearing.West = DirProps{UExit: true, RExit: &EastOfHouse}

	// Kitchen
	Kitchen.East = DirProps{DExit: &KitchenWindow, RExit: &EastOfHouse}
	Kitchen.West = DirProps{UExit: true, RExit: &LivingRoom}
	Kitchen.Out = DirProps{DExit: &KitchenWindow, RExit: &EastOfHouse}
	Kitchen.Up = DirProps{UExit: true, RExit: &Attic}
	Kitchen.Down = DirProps{CExit: func() bool { return false }, RExit: &Studio, CExitStr: "Only Santa Claus climbs down chimneys."}

	// Attic
	Attic.Down = DirProps{UExit: true, RExit: &Kitchen}

	// Living Room
	LivingRoom.East = DirProps{UExit: true, RExit: &Kitchen}
	LivingRoom.West = DirProps{CExit: func() bool { return MagicFlag }, RExit: &StrangePassage, CExitStr: "The door is nailed shut."}
	LivingRoom.Down = DirProps{FExit: TrapDoorExitFcn}

	// Cellar
	Cellar.North = DirProps{UExit: true, RExit: &TrollRoom}
	Cellar.South = DirProps{UExit: true, RExit: &EastOfChasm}
	Cellar.Up = DirProps{DExit: &TrapDoor, RExit: &LivingRoom}
	Cellar.West = DirProps{NExit: "You try to ascend the ramp, but it is impossible, and you slide back down."}

	// Troll Room
	TrollRoom.South = DirProps{UExit: true, RExit: &Cellar}
	TrollRoom.East = DirProps{CExit: func() bool { return TrollFlag }, RExit: &EWPassage, CExitStr: "The troll fends you off with a menacing gesture."}
	TrollRoom.West = DirProps{CExit: func() bool { return TrollFlag }, RExit: &Maze1, CExitStr: "The troll fends you off with a menacing gesture."}

	// East of Chasm
	EastOfChasm.North = DirProps{UExit: true, RExit: &Cellar}
	EastOfChasm.East = DirProps{UExit: true, RExit: &Gallery}
	EastOfChasm.Down = DirProps{NExit: "The chasm probably leads straight to the infernal regions."}

	// Gallery
	Gallery.West = DirProps{UExit: true, RExit: &EastOfChasm}
	Gallery.North = DirProps{UExit: true, RExit: &Studio}

	// Studio
	Studio.South = DirProps{UExit: true, RExit: &Gallery}
	Studio.Up = DirProps{FExit: UpChimneyFcn}

	// Maze 1
	Maze1.East = DirProps{UExit: true, RExit: &TrollRoom}
	Maze1.North = DirProps{UExit: true, RExit: &Maze1}
	Maze1.South = DirProps{UExit: true, RExit: &Maze2}
	Maze1.West = DirProps{UExit: true, RExit: &Maze4}

	// Maze 2
	Maze2.South = DirProps{UExit: true, RExit: &Maze1}
	Maze2.Down = DirProps{FExit: MazeDiodesFcn}
	Maze2.East = DirProps{UExit: true, RExit: &Maze3}

	// Maze 3
	Maze3.West = DirProps{UExit: true, RExit: &Maze2}
	Maze3.North = DirProps{UExit: true, RExit: &Maze4}
	Maze3.Up = DirProps{UExit: true, RExit: &Maze5}

	// Maze 4
	Maze4.West = DirProps{UExit: true, RExit: &Maze3}
	Maze4.North = DirProps{UExit: true, RExit: &Maze1}
	Maze4.East = DirProps{UExit: true, RExit: &DeadEnd1}

	// Dead End 1
	DeadEnd1.South = DirProps{UExit: true, RExit: &Maze4}

	// Maze 5
	Maze5.East = DirProps{UExit: true, RExit: &DeadEnd2}
	Maze5.North = DirProps{UExit: true, RExit: &Maze3}
	Maze5.SouthWest = DirProps{UExit: true, RExit: &Maze6}

	// Dead End 2
	DeadEnd2.West = DirProps{UExit: true, RExit: &Maze5}

	// Maze 6
	Maze6.Down = DirProps{UExit: true, RExit: &Maze5}
	Maze6.East = DirProps{UExit: true, RExit: &Maze7}
	Maze6.West = DirProps{UExit: true, RExit: &Maze6}
	Maze6.Up = DirProps{UExit: true, RExit: &Maze9}

	// Maze 7
	Maze7.Up = DirProps{UExit: true, RExit: &Maze14}
	Maze7.West = DirProps{UExit: true, RExit: &Maze6}
	Maze7.Down = DirProps{FExit: MazeDiodesFcn}
	Maze7.East = DirProps{UExit: true, RExit: &Maze8}
	Maze7.South = DirProps{UExit: true, RExit: &Maze15}

	// Maze 8
	Maze8.NorthEast = DirProps{UExit: true, RExit: &Maze7}
	Maze8.West = DirProps{UExit: true, RExit: &Maze8}
	Maze8.SouthEast = DirProps{UExit: true, RExit: &DeadEnd3}

	// Dead End 3
	DeadEnd3.North = DirProps{UExit: true, RExit: &Maze8}

	// Maze 9
	Maze9.North = DirProps{UExit: true, RExit: &Maze6}
	Maze9.Down = DirProps{FExit: MazeDiodesFcn}
	Maze9.East = DirProps{UExit: true, RExit: &Maze10}
	Maze9.South = DirProps{UExit: true, RExit: &Maze13}
	Maze9.West = DirProps{UExit: true, RExit: &Maze12}
	Maze9.NorthWest = DirProps{UExit: true, RExit: &Maze9}

	// Maze 10
	Maze10.East = DirProps{UExit: true, RExit: &Maze9}
	Maze10.West = DirProps{UExit: true, RExit: &Maze13}
	Maze10.Up = DirProps{UExit: true, RExit: &Maze11}

	// Maze 11
	Maze11.NorthEast = DirProps{UExit: true, RExit: &GratingRoom}
	Maze11.Down = DirProps{UExit: true, RExit: &Maze10}
	Maze11.NorthWest = DirProps{UExit: true, RExit: &Maze13}
	Maze11.SouthWest = DirProps{UExit: true, RExit: &Maze12}

	// Grating Room
	GratingRoom.SouthWest = DirProps{UExit: true, RExit: &Maze11}
	GratingRoom.Up = DirProps{DExit: &Grate, RExit: &GratingClearing, DExitStr: "The grating is closed."}

	// Maze 12
	Maze12.Down = DirProps{FExit: MazeDiodesFcn}
	Maze12.SouthWest = DirProps{UExit: true, RExit: &Maze11}
	Maze12.East = DirProps{UExit: true, RExit: &Maze13}
	Maze12.Up = DirProps{UExit: true, RExit: &Maze9}
	Maze12.North = DirProps{UExit: true, RExit: &DeadEnd4}

	// Dead End 4
	DeadEnd4.South = DirProps{UExit: true, RExit: &Maze12}

	// Maze 13
	Maze13.East = DirProps{UExit: true, RExit: &Maze9}
	Maze13.Down = DirProps{UExit: true, RExit: &Maze12}
	Maze13.South = DirProps{UExit: true, RExit: &Maze10}
	Maze13.West = DirProps{UExit: true, RExit: &Maze11}

	// Maze 14
	Maze14.West = DirProps{UExit: true, RExit: &Maze15}
	Maze14.NorthWest = DirProps{UExit: true, RExit: &Maze14}
	Maze14.NorthEast = DirProps{UExit: true, RExit: &Maze7}
	Maze14.South = DirProps{UExit: true, RExit: &Maze7}

	// Maze 15
	Maze15.West = DirProps{UExit: true, RExit: &Maze14}
	Maze15.South = DirProps{UExit: true, RExit: &Maze7}
	Maze15.SouthEast = DirProps{UExit: true, RExit: &CyclopsRoom}

	// Cyclops Room
	CyclopsRoom.NorthWest = DirProps{UExit: true, RExit: &Maze15}
	CyclopsRoom.East = DirProps{CExit: func() bool { return MagicFlag }, RExit: &StrangePassage, CExitStr: "The east wall is solid rock."}
	CyclopsRoom.Up = DirProps{CExit: func() bool { return CyclopsFlag }, RExit: &TreasureRoom, CExitStr: "The cyclops doesn't look like he'll let you past."}

	// Strange Passage
	StrangePassage.West = DirProps{UExit: true, RExit: &CyclopsRoom}
	StrangePassage.Into = DirProps{UExit: true, RExit: &CyclopsRoom}
	StrangePassage.East = DirProps{UExit: true, RExit: &LivingRoom}

	// Treasure Room
	TreasureRoom.Down = DirProps{UExit: true, RExit: &CyclopsRoom}

	// Reservoir South
	ReservoirSouth.SouthEast = DirProps{UExit: true, RExit: &DeepCanyon}
	ReservoirSouth.SouthWest = DirProps{UExit: true, RExit: &ChasmRoom}
	ReservoirSouth.East = DirProps{UExit: true, RExit: &DamRoom}
	ReservoirSouth.West = DirProps{UExit: true, RExit: &StreamView}
	ReservoirSouth.North = DirProps{CExit: func() bool { return LowTide }, RExit: &Reservoir, CExitStr: "You would drown."}

	// Reservoir
	Reservoir.North = DirProps{UExit: true, RExit: &ReservoirNorth}
	Reservoir.South = DirProps{UExit: true, RExit: &ReservoirSouth}
	Reservoir.Up = DirProps{UExit: true, RExit: &InStream}
	Reservoir.West = DirProps{UExit: true, RExit: &InStream}
	Reservoir.Down = DirProps{NExit: "The dam blocks your way."}

	// Reservoir North
	ReservoirNorth.North = DirProps{UExit: true, RExit: &AtlantisRoom}
	ReservoirNorth.South = DirProps{CExit: func() bool { return LowTide }, RExit: &Reservoir, CExitStr: "You would drown."}

	// Stream View
	StreamView.East = DirProps{UExit: true, RExit: &ReservoirSouth}
	StreamView.West = DirProps{NExit: "The stream emerges from a spot too small for you to enter."}

	// In Stream
	InStream.Up = DirProps{NExit: "The channel is too narrow."}
	InStream.West = DirProps{NExit: "The channel is too narrow."}
	InStream.Land = DirProps{UExit: true, RExit: &StreamView}
	InStream.Down = DirProps{UExit: true, RExit: &Reservoir}
	InStream.East = DirProps{UExit: true, RExit: &Reservoir}

	// Mirror Room 1
	MirrorRoom1.North = DirProps{UExit: true, RExit: &ColdPassage}
	MirrorRoom1.West = DirProps{UExit: true, RExit: &TwistingPassage}
	MirrorRoom1.East = DirProps{UExit: true, RExit: &SmallCave}

	// Mirror Room 2
	MirrorRoom2.West = DirProps{UExit: true, RExit: &WindingPassage}
	MirrorRoom2.North = DirProps{UExit: true, RExit: &NarrowPassage}
	MirrorRoom2.East = DirProps{UExit: true, RExit: &TinyCave}

	// Small Cave
	SmallCave.North = DirProps{UExit: true, RExit: &MirrorRoom1}
	SmallCave.Down = DirProps{UExit: true, RExit: &AtlantisRoom}
	SmallCave.South = DirProps{UExit: true, RExit: &AtlantisRoom}
	SmallCave.West = DirProps{UExit: true, RExit: &TwistingPassage}

	// Tiny Cave
	TinyCave.North = DirProps{UExit: true, RExit: &MirrorRoom2}
	TinyCave.West = DirProps{UExit: true, RExit: &WindingPassage}
	TinyCave.Down = DirProps{UExit: true, RExit: &EnteranceToHades}

	// Cold Passage
	ColdPassage.South = DirProps{UExit: true, RExit: &MirrorRoom1}
	ColdPassage.West = DirProps{UExit: true, RExit: &SlideRoom}

	// Narrow Passage
	NarrowPassage.North = DirProps{UExit: true, RExit: &RoundRoom}
	NarrowPassage.South = DirProps{UExit: true, RExit: &MirrorRoom2}

	// Winding Passage
	WindingPassage.North = DirProps{UExit: true, RExit: &MirrorRoom2}
	WindingPassage.East = DirProps{UExit: true, RExit: &TinyCave}

	// Twisting Passage
	TwistingPassage.North = DirProps{UExit: true, RExit: &MirrorRoom1}
	TwistingPassage.East = DirProps{UExit: true, RExit: &SmallCave}

	// Atlantis Room
	AtlantisRoom.Up = DirProps{UExit: true, RExit: &SmallCave}
	AtlantisRoom.South = DirProps{UExit: true, RExit: &ReservoirNorth}

	// EW Passage
	EWPassage.East = DirProps{UExit: true, RExit: &RoundRoom}
	EWPassage.West = DirProps{UExit: true, RExit: &TrollRoom}
	EWPassage.Down = DirProps{UExit: true, RExit: &ChasmRoom}
	EWPassage.North = DirProps{UExit: true, RExit: &ChasmRoom}

	// Round Room
	RoundRoom.East = DirProps{UExit: true, RExit: &LoudRoom}
	RoundRoom.West = DirProps{UExit: true, RExit: &EWPassage}
	RoundRoom.North = DirProps{UExit: true, RExit: &NSPassage}
	RoundRoom.South = DirProps{UExit: true, RExit: &NarrowPassage}
	RoundRoom.SouthEast = DirProps{UExit: true, RExit: &EngravingsCave}

	// Deep Canyon
	DeepCanyon.NorthWest = DirProps{UExit: true, RExit: &ReservoirSouth}
	DeepCanyon.East = DirProps{UExit: true, RExit: &DamRoom}
	DeepCanyon.SouthWest = DirProps{UExit: true, RExit: &NSPassage}
	DeepCanyon.Down = DirProps{UExit: true, RExit: &LoudRoom}

	// Damp Cave
	DampCave.West = DirProps{UExit: true, RExit: &LoudRoom}
	DampCave.East = DirProps{UExit: true, RExit: &WhiteCliffsNorth}
	DampCave.South = DirProps{NExit: "It is too narrow for most insects."}

	// Loud Room
	LoudRoom.East = DirProps{UExit: true, RExit: &DampCave}
	LoudRoom.West = DirProps{UExit: true, RExit: &RoundRoom}
	LoudRoom.Up = DirProps{UExit: true, RExit: &DeepCanyon}

	// NS Passage
	NSPassage.North = DirProps{UExit: true, RExit: &ChasmRoom}
	NSPassage.NorthEast = DirProps{UExit: true, RExit: &DeepCanyon}
	NSPassage.South = DirProps{UExit: true, RExit: &RoundRoom}

	// Chasm Room
	ChasmRoom.NorthEast = DirProps{UExit: true, RExit: &ReservoirSouth}
	ChasmRoom.SouthWest = DirProps{UExit: true, RExit: &EWPassage}
	ChasmRoom.Up = DirProps{UExit: true, RExit: &EWPassage}
	ChasmRoom.South = DirProps{UExit: true, RExit: &NSPassage}
	ChasmRoom.Down = DirProps{NExit: "Are you out of your mind?"}

	// Entrance to Hades
	EnteranceToHades.Up = DirProps{UExit: true, RExit: &TinyCave}
	EnteranceToHades.Into = DirProps{CExit: func() bool { return LLDFlag }, RExit: &LandOfLivingDead, CExitStr: "Some invisible force prevents you from passing through the gate."}
	EnteranceToHades.South = DirProps{CExit: func() bool { return LLDFlag }, RExit: &LandOfLivingDead, CExitStr: "Some invisible force prevents you from passing through the gate."}

	// Land of Living Dead
	LandOfLivingDead.Out = DirProps{UExit: true, RExit: &EnteranceToHades}
	LandOfLivingDead.North = DirProps{UExit: true, RExit: &EnteranceToHades}

	// Engravings Cave
	EngravingsCave.NorthWest = DirProps{UExit: true, RExit: &RoundRoom}
	EngravingsCave.East = DirProps{UExit: true, RExit: &DomeRoom}

	// Egypt Room
	EgyptRoom.West = DirProps{UExit: true, RExit: &NorthTemple}
	EgyptRoom.Up = DirProps{UExit: true, RExit: &NorthTemple}

	// Dome Room
	DomeRoom.West = DirProps{UExit: true, RExit: &EngravingsCave}
	DomeRoom.Down = DirProps{CExit: func() bool { return DomeFlag }, RExit: &TorchRoom, CExitStr: "You cannot go down without fracturing many bones."}

	// Torch Room
	TorchRoom.Up = DirProps{NExit: "You cannot reach the rope."}
	TorchRoom.South = DirProps{UExit: true, RExit: &NorthTemple}
	TorchRoom.Down = DirProps{UExit: true, RExit: &NorthTemple}

	// North Temple
	NorthTemple.Down = DirProps{UExit: true, RExit: &EgyptRoom}
	NorthTemple.East = DirProps{UExit: true, RExit: &EgyptRoom}
	NorthTemple.North = DirProps{UExit: true, RExit: &TorchRoom}
	NorthTemple.Out = DirProps{UExit: true, RExit: &TorchRoom}
	NorthTemple.Up = DirProps{UExit: true, RExit: &TorchRoom}
	NorthTemple.South = DirProps{UExit: true, RExit: &SouthTemple}

	// South Temple
	SouthTemple.North = DirProps{UExit: true, RExit: &NorthTemple}
	SouthTemple.Down = DirProps{CExit: func() bool { return CoffinCure }, RExit: &TinyCave, CExitStr: "You haven't a prayer of getting the coffin down there."}

	// Dam Room
	DamRoom.South = DirProps{UExit: true, RExit: &DeepCanyon}
	DamRoom.Down = DirProps{UExit: true, RExit: &DamBase}
	DamRoom.East = DirProps{UExit: true, RExit: &DamBase}
	DamRoom.North = DirProps{UExit: true, RExit: &DamLobby}
	DamRoom.West = DirProps{UExit: true, RExit: &ReservoirSouth}

	// Dam Lobby
	DamLobby.South = DirProps{UExit: true, RExit: &DamRoom}
	DamLobby.North = DirProps{UExit: true, RExit: &MaintenanceRoom}
	DamLobby.East = DirProps{UExit: true, RExit: &MaintenanceRoom}

	// Maintenance Room
	MaintenanceRoom.South = DirProps{UExit: true, RExit: &DamLobby}
	MaintenanceRoom.West = DirProps{UExit: true, RExit: &DamLobby}

	// Dam Base
	DamBase.North = DirProps{UExit: true, RExit: &DamRoom}
	DamBase.Up = DirProps{UExit: true, RExit: &DamRoom}

	// River 1
	River1.Up = DirProps{NExit: "You cannot go upstream due to strong currents."}
	River1.West = DirProps{UExit: true, RExit: &DamBase}
	River1.Land = DirProps{UExit: true, RExit: &DamBase}
	River1.Down = DirProps{UExit: true, RExit: &River2}
	River1.East = DirProps{NExit: "The White Cliffs prevent your landing here."}

	// River 2
	River2.Up = DirProps{NExit: "You cannot go upstream due to strong currents."}
	River2.Down = DirProps{UExit: true, RExit: &River3}
	River2.Land = DirProps{NExit: "There is no safe landing spot here."}
	River2.East = DirProps{NExit: "The White Cliffs prevent your landing here."}
	River2.West = DirProps{NExit: "Just in time you steer away from the rocks."}

	// River 3
	River3.Up = DirProps{NExit: "You cannot go upstream due to strong currents."}
	River3.Down = DirProps{UExit: true, RExit: &River4}
	River3.Land = DirProps{UExit: true, RExit: &WhiteCliffsNorth}
	River3.West = DirProps{UExit: true, RExit: &WhiteCliffsNorth}

	// White Cliffs North
	WhiteCliffsNorth.South = DirProps{CExit: func() bool { return DeflateFlag }, RExit: &WhiteCliffsSouth, CExitStr: "The path is too narrow."}
	WhiteCliffsNorth.West = DirProps{CExit: func() bool { return DeflateFlag }, RExit: &DampCave, CExitStr: "The path is too narrow."}

	// White Cliffs South
	WhiteCliffsSouth.North = DirProps{CExit: func() bool { return DeflateFlag }, RExit: &WhiteCliffsNorth, CExitStr: "The path is too narrow."}

	// River 4
	River4.Up = DirProps{NExit: "You cannot go upstream due to strong currents."}
	River4.Down = DirProps{UExit: true, RExit: &River5}
	River4.Land = DirProps{NExit: "You can land either to the east or the west."}
	River4.West = DirProps{UExit: true, RExit: &WhiteCliffsSouth}
	River4.East = DirProps{UExit: true, RExit: &SandyBeach}

	// River 5
	River5.Up = DirProps{NExit: "You cannot go upstream due to strong currents."}
	River5.East = DirProps{UExit: true, RExit: &Shore}
	River5.Land = DirProps{UExit: true, RExit: &Shore}

	// Shore
	Shore.North = DirProps{UExit: true, RExit: &SandyBeach}
	Shore.South = DirProps{UExit: true, RExit: &AragainFalls}

	// Sandy Beach
	SandyBeach.NorthEast = DirProps{UExit: true, RExit: &SandyCave}
	SandyBeach.South = DirProps{UExit: true, RExit: &Shore}

	// Sandy Cave
	SandyCave.SouthWest = DirProps{UExit: true, RExit: &SandyBeach}

	// Aragain Falls
	AragainFalls.West = DirProps{CExit: func() bool { return RainbowFlag }, RExit: &OnRainbow}
	AragainFalls.Down = DirProps{NExit: "It's a long way..."}
	AragainFalls.North = DirProps{UExit: true, RExit: &Shore}
	AragainFalls.Up = DirProps{CExit: func() bool { return RainbowFlag }, RExit: &OnRainbow}

	// On Rainbow
	OnRainbow.West = DirProps{UExit: true, RExit: &EndOfRainbow}
	OnRainbow.East = DirProps{UExit: true, RExit: &AragainFalls}

	// End of Rainbow
	EndOfRainbow.Up = DirProps{CExit: func() bool { return RainbowFlag }, RExit: &OnRainbow}
	EndOfRainbow.NorthEast = DirProps{CExit: func() bool { return RainbowFlag }, RExit: &OnRainbow}
	EndOfRainbow.East = DirProps{CExit: func() bool { return RainbowFlag }, RExit: &OnRainbow}
	EndOfRainbow.SouthWest = DirProps{UExit: true, RExit: &CanyonBottom}

	// Canyon Bottom
	CanyonBottom.Up = DirProps{UExit: true, RExit: &CliffMiddle}
	CanyonBottom.North = DirProps{UExit: true, RExit: &EndOfRainbow}

	// Cliff Middle
	CliffMiddle.Up = DirProps{UExit: true, RExit: &CanyonView}
	CliffMiddle.Down = DirProps{UExit: true, RExit: &CanyonBottom}

	// Canyon View
	CanyonView.East = DirProps{UExit: true, RExit: &CliffMiddle}
	CanyonView.Down = DirProps{UExit: true, RExit: &CliffMiddle}
	CanyonView.NorthWest = DirProps{UExit: true, RExit: &Clearing}
	CanyonView.West = DirProps{UExit: true, RExit: &Forest3}
	CanyonView.South = DirProps{NExit: "Storm-tossed trees block your way."}

	// Mine Entrance
	MineEntrance.South = DirProps{UExit: true, RExit: &SlideRoom}
	MineEntrance.Into = DirProps{UExit: true, RExit: &SqueekyRoom}
	MineEntrance.West = DirProps{UExit: true, RExit: &SqueekyRoom}

	// Squeaky Room
	SqueekyRoom.North = DirProps{UExit: true, RExit: &BatRoom}
	SqueekyRoom.East = DirProps{UExit: true, RExit: &MineEntrance}

	// Bat Room
	BatRoom.South = DirProps{UExit: true, RExit: &SqueekyRoom}
	BatRoom.East = DirProps{UExit: true, RExit: &ShaftRoom}

	// Shaft Room
	ShaftRoom.Down = DirProps{NExit: "You wouldn't fit and would die if you could."}
	ShaftRoom.West = DirProps{UExit: true, RExit: &BatRoom}
	ShaftRoom.North = DirProps{UExit: true, RExit: &SmellyRoom}

	// Smelly Room
	SmellyRoom.Down = DirProps{UExit: true, RExit: &GasRoom}
	SmellyRoom.South = DirProps{UExit: true, RExit: &ShaftRoom}

	// Gas Room
	GasRoom.Up = DirProps{UExit: true, RExit: &SmellyRoom}
	GasRoom.East = DirProps{UExit: true, RExit: &Mine1}

	// Ladder Top
	LadderTop.Down = DirProps{UExit: true, RExit: &LadderBottom}
	LadderTop.Up = DirProps{UExit: true, RExit: &Mine4}

	// Ladder Bottom
	LadderBottom.South = DirProps{UExit: true, RExit: &DeadEnd5}
	LadderBottom.West = DirProps{UExit: true, RExit: &TimberRoom}
	LadderBottom.Up = DirProps{UExit: true, RExit: &LadderTop}

	// Dead End 5
	DeadEnd5.North = DirProps{UExit: true, RExit: &LadderBottom}

	// Timber Room
	TimberRoom.East = DirProps{UExit: true, RExit: &LadderBottom}
	TimberRoom.West = DirProps{CExit: func() bool { return EmptyHanded }, RExit: &LowerShaft, CExitStr: "You cannot fit through this passage with that load."}

	// Lower Shaft
	LowerShaft.South = DirProps{UExit: true, RExit: &MachineRoom}
	LowerShaft.Out = DirProps{CExit: func() bool { return EmptyHanded }, RExit: &TimberRoom, CExitStr: "You cannot fit through this passage with that load."}
	LowerShaft.East = DirProps{CExit: func() bool { return EmptyHanded }, RExit: &TimberRoom, CExitStr: "You cannot fit through this passage with that load."}

	// Machine Room
	MachineRoom.North = DirProps{UExit: true, RExit: &LowerShaft}

	// Mine 1
	Mine1.North = DirProps{UExit: true, RExit: &GasRoom}
	Mine1.East = DirProps{UExit: true, RExit: &Mine1}
	Mine1.NorthEast = DirProps{UExit: true, RExit: &Mine2}

	// Mine 2
	Mine2.North = DirProps{UExit: true, RExit: &Mine2}
	Mine2.South = DirProps{UExit: true, RExit: &Mine1}
	Mine2.SouthEast = DirProps{UExit: true, RExit: &Mine3}

	// Mine 3
	Mine3.South = DirProps{UExit: true, RExit: &Mine3}
	Mine3.SouthWest = DirProps{UExit: true, RExit: &Mine4}
	Mine3.East = DirProps{UExit: true, RExit: &Mine2}

	// Mine 4
	Mine4.North = DirProps{UExit: true, RExit: &Mine3}
	Mine4.West = DirProps{UExit: true, RExit: &Mine4}
	Mine4.Down = DirProps{UExit: true, RExit: &LadderTop}

	// Slide Room
	SlideRoom.East = DirProps{UExit: true, RExit: &ColdPassage}
	SlideRoom.North = DirProps{UExit: true, RExit: &MineEntrance}
	SlideRoom.Down = DirProps{UExit: true, RExit: &Cellar}
}

func FinalizeGameObjects() {
	// Set up room exits
	InitRoomExits()

	// Add TrapDoor to Cellar's globals (done here to avoid potential init ordering issues)
	Cellar.Global = append(Cellar.Global, &TrapDoor)

	// Set action functions that would cause init cycles
	WhiteHouse.Action = WhiteHouseFcn
	Mailbox.Action = MailboxFcn
	Forest.Action = ForestFcn
	KitchenWindow.Action = KitchenWindowFcn
	Chimney.Action = ChimneyFcn
	Grate.Action = GrateFcn
	ClimbableCliff.Action = CliffObjectFcn
	Gunk.Action = GunkFcn
	GratingClearing.Action = ClearingFcn
	Clearing.Action = ForestRoomFcn
	Kitchen.Action = KitchenFcn
	LivingRoom.Action = LivingRoomFcn
	StoneBarrow.Action = StoneBarrowFcn
	Rainbow.Action = RainbowFcn
	InflatedBoat.Action = RBoatFcn
	HotBell.Action = HotBellFcn
	TrollRoom.Action = TrollRoomFcn
	LowerShaft.Action = NoObjsFcn
	MachineRoom.Action = MachineRoomFcn
	TrophyCase.Action = TrophyCaseFcn
	Sword.Action = SwordFcn
	Lamp.Action = LanternFcn
	Bell.Action = BellFcn
	Match.Action = MatchFcn
	ToolChest.Action = ToolChestFcn
	Machine.Action = MachineFcn
	MachineSwitch.Action = MachineSwitchFcn
	Sceptre.Action = SceptreFcn
	Board.Action = BoardFcn
	Teeth.Action = TeethFcn
	GraniteWall.Action = GraniteWallFcn
	Songbird.Action = SongbirdFcn
	GlobalWater.Action = WaterFcn
	Slide.Action = SlideFcn
	Bodies.Action = BodyFcn
	Crack.Action = CrackFcn
	WhiteCliff.Action = WhiteCliffFcn
	River.Action = RiverFcn
	BoardedWindow.Action = BoardedWindowFcn
	PuncturedBoat.Action = PuncturedBoatFcn
	WestOfHouse.Action = WestHouseFcn
	EastOfHouse.Action = EastHouseFcn
	Forest1.Action = ForestRoomFcn
	Forest2.Action = ForestRoomFcn
	Forest3.Action = ForestRoomFcn
	Path.Action = ForestRoomFcn
	UpATree.Action = TreeRoomFcn
	Cellar.Action = CellarFcn
	GratingRoom.Action = Maze11Fcn
	CyclopsRoom.Action = CyclopsRoomFcn
	TreasureRoom.Action = TreasureRoomFcn
	ReservoirSouth.Action = ReservoirSouthFcn
	Reservoir.Action = ReservoirFcn
	ReservoirNorth.Action = ReservoirNorthFcn
	MirrorRoom1.Action = MirrorRoomFcn
	MirrorRoom2.Action = MirrorRoomFcn
	TinyCave.Action = Cave2RoomFcn
	DeepCanyon.Action = DeepCanyonFcn
	LoudRoom.Action = LoudRoomFcn
	EnteranceToHades.Action = LLDRoomFcn
	DomeRoom.Action = DomeRoomFcn
	TorchRoom.Action = TorchRoomFcn
	SouthTemple.Action = SouthTempleFcn
	DamRoom.Action = DamRoomFcn
	WhiteCliffsNorth.Action = WhiteCliffsFcn
	WhiteCliffsSouth.Action = WhiteCliffsFcn
	River4.Action = Rivr4RoomFcn
	AragainFalls.Action = FallsRoomFcn
	CanyonView.Action = CanyonViewFcn
	BatRoom.Action = BatsRoomFcn
	GasRoom.Action = BoomRoomFcn
	TimberRoom.Action = NoObjsFcn
	MountainRange.Action = MountainRangeFcn
	FrontDoor.Action = FrontDoorFcn
	BarrowDoor.Action = BarrowDoorFcn
	Barrow.Action = BarrowFcn
	Rug.Action = RugFcn
	TrapDoor.Action = TrapDoorFcn
	WoodenDoor.Action = FrontDoorFcn
	Rope.Action = RopeFcn
	Ghosts.Action = GhostsFcn
	RaisedBasket.Action = BasketFcn
	LoweredBasket.Action = BasketFcn
	Bat.Action = BatFcn
	Candles.Action = CandlesFcn
	Troll.Action = TrollFcn
	Bolt.Action = BoltFcn
	Bubble.Action = BubbleFcn
	Dam.Action = DamFunction
	InflatableBoat.Action = InflatableBoatFcn
	YellowButton.Action = ButtonFcn
	BrownButton.Action = ButtonFcn
	RedButton.Action = ButtonFcn
	BlueButton.Action = ButtonFcn
	Tube.Action = TubeFcn
	Leak.Action = LeakFcn
	Cyclops.Action = CyclopsFcn
	Chalice.Action = ChaliceFcn
	Painting.Action = PaintingFcn
	Leaves.Action = LeafPileFcn
	Sand.Action = SandFunction
	Thief.Action = RobberFcn
	Trunk.Action = TrunkFcn
	Mirror1.Action = MirrorMirrorFcn
	Mirror2.Action = MirrorMirrorFcn
	Pedestal.Action = DumbContainerFcn
	Buoy.Action = TreasureInsideFcn
	Bones.Action = SkeletonFcn
	BagOfCoins.Action = BagOfCoinsFcn
	RustyKnife.Action = RustyKnifeFcn
	Bottle.Action = BottleFcn
	SandwichBag.Action = SandwichBagFcn
	Knife.Action = KnifeFcn
	Water.Action = WaterFcn
	Garlic.Action = GarlicFcn
	Book.Action = BlackBookFcn
	Egg.Action = EggObjectFcn
	Canary.Action = CanaryObjectFcn
	Putty.Action = PuttyFcn
	Axe.Action = AxeFcn
	LargeBag.Action = LargeBagFcn
	Stiletto.Action = StiletteFcn
	Torch.Action = TorchFcn
	BrokenCanary.Action = CanaryObjectFcn

	// Set up globals.go object actions
	LocalGlobals.Pseudo = []PseudoObj{{
		Synonym: "foobar",
		Action:  VWalk,
	}}
	NotHereObject.Action = NotHereObjectFcn
	Me.Action = CretinFcn
	Ground.Action = GroundFunction

	// Initialize villain table
	Villains = []*VillainEntry{
		{Villain: &Troll, Best: &Sword, BestAdv: 1, Prob: 0, Msgs: &TrollMelee},
		{Villain: &Thief, Best: &Knife, BestAdv: 1, Prob: 0, Msgs: &ThiefMelee},
		{Villain: &Cyclops, Best: nil, BestAdv: 0, Prob: 0, Msgs: &CyclopsMelee},
	}
}
