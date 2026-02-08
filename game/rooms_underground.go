package game

import . "github.com/ajdnik/gozork/engine"

var (
	// ================================================================
	// ROOMS - Cyclops and Hideaway
	// ================================================================

	CyclopsRoom = Object{
		In:   &Rooms,
		Desc: "Cyclops Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgLand,
		Global: []*Object{&Stairs},
	}
	StrangePassage = Object{
		In:       &Rooms,
		LongDesc: "This is a long passage. To the west is one entrance. On the east there is an old wooden door, with a large opening in it (about cyclops sized).",
		Desc:     "Strange Passage",
		Flags:    FlgLand,
	}
	TreasureRoom = Object{
		In:       &Rooms,
		LongDesc: "This is a large room, whose east wall is solid granite. A number of discarded bags, which crumble at your touch, are scattered about on the floor. There is an exit down a staircase.",
		Desc:     "Treasure Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgLand,
		Global: []*Object{&Stairs},
		Item:   &ItemData{Value: 25},
	}

	// ================================================================
	// ROOMS - Reservoir Area
	// ================================================================

	ReservoirSouth = Object{
		In:   &Rooms,
		Desc: "Reservoir South",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgLand,
		Global: []*Object{&GlobalWater},
		Pseudo: []PseudoObj{
			{Synonym: "lake", Action: LakePseudo},
			{Synonym: "chasm", Action: ChasmPseudo},
		},
	}
	Reservoir = Object{
		In:   &Rooms,
		Desc: "Reservoir",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgNonLand,
		Global: []*Object{&GlobalWater},
		Pseudo: []PseudoObj{{Synonym: "stream", Action: StreamPseudo}},
	}
	ReservoirNorth = Object{
		In:   &Rooms,
		Desc: "Reservoir North",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgLand,
		Global: []*Object{&GlobalWater, &Stairs},
		Pseudo: []PseudoObj{{Synonym: "lake", Action: LakePseudo}},
	}
	StreamView = Object{
		In:       &Rooms,
		LongDesc: "You are standing on a path beside a gently flowing stream. The path follows the stream, which flows from west to east.",
		Desc:     "Stream View",
		Flags:    FlgLand,
		Global:   []*Object{&GlobalWater},
		Pseudo:   []PseudoObj{{Synonym: "stream", Action: StreamPseudo}},
	}
	InStream = Object{
		In:       &Rooms,
		LongDesc: "You are on the gently flowing stream. The upstream route is too narrow to navigate, and the downstream route is invisible due to twisting walls. There is a narrow beach to land on.",
		Desc:     "Stream",
		Flags:    FlgNonLand,
		Global:   []*Object{&GlobalWater},
		Pseudo:   []PseudoObj{{Synonym: "stream", Action: StreamPseudo}},
	}

	// ================================================================
	// ROOMS - Mirror Rooms and Vicinity
	// ================================================================

	MirrorRoom1 = Object{
		In:   &Rooms,
		Desc: "Mirror Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags: FlgLand,
	}
	MirrorRoom2 = Object{
		In:   &Rooms,
		Desc: "Mirror Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags: FlgLand | FlgOn,
	}
	SmallCave = Object{
		In:       &Rooms,
		LongDesc: "This is a tiny cave with entrances west and north, and a staircase leading down.",
		Desc:     "Cave",
		Flags:    FlgLand,
		Global:   []*Object{&Stairs},
	}
	TinyCave = Object{
		In:       &Rooms,
		LongDesc: "This is a tiny cave with entrances west and north, and a dark, forbidding staircase leading down.",
		Desc:     "Cave",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgLand,
		Global: []*Object{&Stairs},
	}
	ColdPassage = Object{
		In:       &Rooms,
		LongDesc: "This is a cold and damp corridor where a long east-west passageway turns into a southward path.",
		Desc:     "Cold Passage",
		Flags:    FlgLand,
	}
	NarrowPassage = Object{
		In:       &Rooms,
		LongDesc: "This is a long and narrow corridor where a long north-south passageway briefly narrows even further.",
		Desc:     "Narrow Passage",
		Flags:    FlgLand,
	}
	WindingPassage = Object{
		In:       &Rooms,
		LongDesc: "This is a winding passage. It seems that there are only exits on the east and north.",
		Desc:     "Winding Passage",
		Flags:    FlgLand,
	}
	TwistingPassage = Object{
		In:       &Rooms,
		LongDesc: "This is a winding passage. It seems that there are only exits on the east and north.",
		Desc:     "Twisting Passage",
		Flags:    FlgLand,
	}
	AtlantisRoom = Object{
		In:       &Rooms,
		LongDesc: "This is an ancient room, long under water. There is an exit to the south and a staircase leading up.",
		Desc:     "Atlantis Room",
		Flags:    FlgLand,
		Global:   []*Object{&Stairs},
	}

	// ================================================================
	// ROOMS - Round Room and Vicinity
	// ================================================================

	EWPassage = Object{
		In:       &Rooms,
		LongDesc: "This is a narrow east-west passageway. There is a narrow stairway leading down at the north end of the room.",
		Desc:     "East-West Passage",
		Flags:    FlgLand,
		Global:   []*Object{&Stairs},
		Item:     &ItemData{Value: 5},
	}
	RoundRoom = Object{
		In:       &Rooms,
		LongDesc: "This is a circular stone room with passages in all directions. Several of them have unfortunately been blocked by cave-ins.",
		Desc:     "Round Room",
		Flags:    FlgLand,
	}
	DeepCanyon = Object{
		In:   &Rooms,
		Desc: "Deep Canyon",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgLand,
		Global: []*Object{&Stairs},
	}
	DampCave = Object{
		In:       &Rooms,
		LongDesc: "This cave has exits to the west and east, and narrows to a crack toward the south. The earth is particularly damp here.",
		Desc:     "Damp Cave",
		Flags:    FlgLand,
		Global:   []*Object{&Crack},
	}
	LoudRoom = Object{
		In:   &Rooms,
		Desc: "Loud Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgLand,
		Global: []*Object{&Stairs},
	}
	NSPassage = Object{
		In:       &Rooms,
		LongDesc: "This is a high north-south passage, which forks to the northeast.",
		Desc:     "North-South Passage",
		Flags:    FlgLand,
	}
	ChasmRoom = Object{
		In:       &Rooms,
		LongDesc: "A chasm runs southwest to northeast and the path follows it. You are on the south side of the chasm, where a crack opens into a passage.",
		Desc:     "Chasm",
		Flags:    FlgLand,
		Global:   []*Object{&Crack, &Stairs},
		Pseudo:   []PseudoObj{{Synonym: "chasm", Action: ChasmPseudo}},
	}

	// ================================================================
	// ROOMS - Hades
	// ================================================================

	EnteranceToHades = Object{
		In:   &Rooms,
		Desc: "Entrance to Hades",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgLand | FlgOn,
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
		Flags:    FlgLand | FlgOn,
		Global:   []*Object{&Bodies},
	}

	// ================================================================
	// ROOMS - Dome, Temple, Egypt
	// ================================================================

	EngravingsCave = Object{
		In:       &Rooms,
		LongDesc: "You have entered a low cave with passages leading northwest and east.",
		Desc:     "Engravings Cave",
		Flags:    FlgLand,
	}
	EgyptRoom = Object{
		In:       &Rooms,
		LongDesc: "This is a room which looks like an Egyptian tomb. There is an ascending staircase to the west.",
		Desc:     "Egyptian Room",
		Flags:    FlgLand,
		Global:   []*Object{&Stairs},
	}
	DomeRoom = Object{
		In:   &Rooms,
		Desc: "Dome Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgLand,
		Pseudo: []PseudoObj{{Synonym: "dome", Action: DomePseudo}},
	}
	TorchRoom = Object{
		In:   &Rooms,
		Desc: "Torch Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgLand,
		Global: []*Object{&Stairs},
		Pseudo: []PseudoObj{{Synonym: "dome", Action: DomePseudo}},
	}
	NorthTemple = Object{
		In:       &Rooms,
		LongDesc: "This is the north end of a large temple. On the east wall is an ancient inscription, probably a prayer in a long-forgotten language. Below the prayer is a staircase leading down. The west wall is solid granite. The exit to the north end of the room is through huge marble pillars.",
		Desc:     "Temple",
		Flags:    FlgLand | FlgOn | FlgSacred,
		Global:   []*Object{&Stairs},
	}
	SouthTemple = Object{
		In:       &Rooms,
		LongDesc: "This is the south end of a large temple. In front of you is what appears to be an altar. In one corner is a small hole in the floor which leads into darkness. You probably could not get back up it.",
		Desc:     "Altar",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags: FlgLand | FlgOn | FlgSacred,
	}

	// ================================================================
	// ROOMS - Flood Control Dam #3
	// ================================================================

	DamRoom = Object{
		In:   &Rooms,
		Desc: "Dam",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgLand | FlgOn,
		Global: []*Object{&GlobalWater},
	}
	DamLobby = Object{
		In:       &Rooms,
		LongDesc: "This room appears to have been the waiting room for groups touring the dam. There are open doorways here to the north and east marked \"Private\", and there is a path leading south over the top of the dam.",
		Desc:     "Dam Lobby",
		Flags:    FlgLand | FlgOn,
	}
	MaintenanceRoom = Object{
		In:       &Rooms,
		LongDesc: "This is what appears to have been the maintenance room for Flood Control Dam #3. Apparently, this room has been ransacked recently, for most of the valuable equipment is gone. On the wall in front of you is a group of buttons colored blue, yellow, brown, and red. There are doorways to the west and south.",
		Desc:     "Maintenance Room",
		Flags:    FlgLand,
	}

	// ================================================================
	// ROOMS - River Area
	// ================================================================

	DamBase = Object{
		In:       &Rooms,
		LongDesc: "You are at the base of Flood Control Dam #3, which looms above you and to the north. The river Frigid is flowing by here. Along the river are the White Cliffs which seem to form giant walls stretching from north to south along the shores of the river as it winds its way downstream.",
		Desc:     "Dam Base",
		Flags:    FlgLand | FlgOn | FlgSacred,
		Global:   []*Object{&GlobalWater, &River},
	}
	River1 = Object{
		In:       &Rooms,
		LongDesc: "You are on the Frigid River in the vicinity of the Dam. The river flows quietly here. There is a landing on the west shore.",
		Desc:     "Frigid River",
		Flags:    FlgNonLand | FlgSacred | FlgOn,
		Global:   []*Object{&GlobalWater, &River},
	}
	River2 = Object{
		In:       &Rooms,
		LongDesc: "The river turns a corner here making it impossible to see the Dam. The White Cliffs loom on the east bank and large rocks prevent landing on the west.",
		Desc:     "Frigid River",
		Flags:    FlgNonLand | FlgSacred,
		Global:   []*Object{&GlobalWater, &River},
	}
	River3 = Object{
		In:       &Rooms,
		LongDesc: "The river descends here into a valley. There is a narrow beach on the west shore below the cliffs. In the distance a faint rumbling can be heard.",
		Desc:     "Frigid River",
		Flags:    FlgNonLand | FlgSacred,
		Global:   []*Object{&GlobalWater, &River},
	}
	WhiteCliffsNorth = Object{
		In:       &Rooms,
		LongDesc: "You are on a narrow strip of beach which runs along the base of the White Cliffs. There is a narrow path heading south along the Cliffs and a tight passage leading west into the cliffs themselves.",
		Desc:     "White Cliffs Beach",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgLand | FlgSacred,
		Global: []*Object{&GlobalWater, &WhiteCliff, &River},
	}
	WhiteCliffsSouth = Object{
		In:       &Rooms,
		LongDesc: "You are on a rocky, narrow strip of beach beside the Cliffs. A narrow path leads north along the shore.",
		Desc:     "White Cliffs Beach",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgLand | FlgSacred,
		Global: []*Object{&GlobalWater, &WhiteCliff, &River},
	}
	River4 = Object{
		In:       &Rooms,
		LongDesc: "The river is running faster here and the sound ahead appears to be that of rushing water. On the east shore is a sandy beach. A small area of beach can also be seen below the cliffs on the west shore.",
		Desc:     "Frigid River",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgNonLand | FlgSacred,
		Global: []*Object{&GlobalWater, &River},
	}
	River5 = Object{
		In:       &Rooms,
		LongDesc: "The sound of rushing water is nearly unbearable here. On the east shore is a large landing area.",
		Desc:     "Frigid River",
		Flags:    FlgNonLand | FlgSacred | FlgOn,
		Global:   []*Object{&GlobalWater, &River},
	}
	Shore = Object{
		In:       &Rooms,
		LongDesc: "You are on the east shore of the river. The water here seems somewhat treacherous. A path travels from north to south here, the south end quickly turning around a sharp corner.",
		Desc:     "Shore",
		Flags:    FlgLand | FlgSacred | FlgOn,
		Global:   []*Object{&GlobalWater, &River},
	}
	SandyBeach = Object{
		In:       &Rooms,
		LongDesc: "You are on a large sandy beach on the east shore of the river, which is flowing quickly by. A path runs beside the river to the south here, and a passage is partially buried in sand to the northeast.",
		Desc:     "Sandy Beach",
		Flags:    FlgLand | FlgSacred,
		Global:   []*Object{&GlobalWater, &River},
	}
	SandyCave = Object{
		In:       &Rooms,
		LongDesc: "This is a sand-filled cave whose exit is to the southwest.",
		Desc:     "Sandy Cave",
		Flags:    FlgLand,
	}
	AragainFalls = Object{
		In:   &Rooms,
		Desc: "Aragain Falls",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgLand | FlgSacred | FlgOn,
		Global: []*Object{&GlobalWater, &River, &Rainbow},
	}
	OnRainbow = Object{
		In:       &Rooms,
		LongDesc: "You are on top of a rainbow (I bet you never thought you would walk on a rainbow), with a magnificent view of the Falls. The rainbow travels east-west here.",
		Desc:     "On the Rainbow",
		Flags:    FlgLand | FlgOn | FlgSacred,
		Global:   []*Object{&Rainbow},
	}
	EndOfRainbow = Object{
		In:       &Rooms,
		LongDesc: "You are on a small, rocky beach on the continuation of the Frigid River past the Falls. The beach is narrow due to the presence of the White Cliffs. The river canyon opens here and sunlight shines in from above. A rainbow crosses over the falls to the east and a narrow path continues to the southwest.",
		Desc:     "End of Rainbow",
		Flags:    FlgLand | FlgOn,
		Global:   []*Object{&GlobalWater, &Rainbow, &River},
	}
	CanyonBottom = Object{
		In:       &Rooms,
		LongDesc: "You are beneath the walls of the river canyon which may be climbable here. The lesser part of the runoff of Aragain Falls flows by below. To the north is a narrow path.",
		Desc:     "Canyon Bottom",
		Flags:    FlgLand | FlgOn | FlgSacred,
		Global:   []*Object{&GlobalWater, &ClimbableCliff, &River},
	}
	CliffMiddle = Object{
		In:       &Rooms,
		LongDesc: "You are on a ledge about halfway up the wall of the river canyon. You can see from here that the main flow from Aragain Falls twists along a passage which it is impossible for you to enter. Below you is the canyon bottom. Above you is more cliff, which appears climbable.",
		Desc:     "Rocky Ledge",
		Flags:    FlgLand | FlgOn | FlgSacred,
		Global:   []*Object{&ClimbableCliff, &River},
	}
	CanyonView = Object{
		In:       &Rooms,
		LongDesc: "You are at the top of the Great Canyon on its west wall. From here there is a marvelous view of the canyon and parts of the Frigid River upstream. Across the canyon, the walls of the White Cliffs join the mighty ramparts of the Flathead Mountains to the east. Following the Canyon upstream to the north, Aragain Falls may be seen, complete with rainbow. The mighty Frigid River flows out from a great dark cavern. To the west and south can be seen an immense forest, stretching for miles around. A path leads northwest. It is possible to climb down into the canyon from here.",
		Desc:     "Canyon View",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgLand | FlgOn | FlgSacred,
		Global: []*Object{&ClimbableCliff, &River, &Rainbow},
	}

	// ================================================================
	// ROOMS - Coal Mine Area
	// ================================================================

	MineEntrance = Object{
		In:       &Rooms,
		LongDesc: "You are standing at the entrance of what might have been a coal mine. The shaft enters the west wall, and there is another exit on the south end of the room.",
		Desc:     "Mine Entrance",
		Flags:    FlgLand,
	}
	SqueekyRoom = Object{
		In:       &Rooms,
		LongDesc: "You are in a small room. Strange squeaky sounds may be heard coming from the passage at the north end. You may also escape to the east.",
		Desc:     "Squeaky Room",
		Flags:    FlgLand,
	}
	BatRoom = Object{
		In:   &Rooms,
		Desc: "Bat Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags: FlgLand | FlgSacred,
	}
	ShaftRoom = Object{
		In:       &Rooms,
		LongDesc: "This is a large room, in the middle of which is a small shaft descending through the floor into darkness below. To the west and the north are exits from this room. Constructed over the top of the shaft is a metal framework to which a heavy iron chain is attached.",
		Desc:     "Shaft Room",
		Flags:    FlgLand,
		Pseudo:   []PseudoObj{{Synonym: "chain", Action: ChainPseudo}},
	}
	SmellyRoom = Object{
		In:       &Rooms,
		LongDesc: "This is a small nondescript room. However, from the direction of a small descending staircase a foul odor can be detected. To the south is a narrow tunnel.",
		Desc:     "Smelly Room",
		Flags:    FlgLand,
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
		Flags:  FlgLand | FlgSacred,
		Global: []*Object{&Stairs},
		Pseudo: []PseudoObj{
			{Synonym: "gas", Action: GasPseudo},
			{Synonym: "odor", Action: GasPseudo},
		},
	}
	LadderTop = Object{
		In:       &Rooms,
		LongDesc: "This is a very small room. In the corner is a rickety wooden ladder, leading downward. It might be safe to descend. There is also a staircase leading upward.",
		Desc:     "Ladder Top",
		Flags:    FlgLand,
		Global:   []*Object{&Ladder, &Stairs},
	}
	LadderBottom = Object{
		In:       &Rooms,
		LongDesc: "This is a rather wide room. On one side is the bottom of a narrow wooden ladder. To the west and the south are passages leaving the room.",
		Desc:     "Ladder Bottom",
		Flags:    FlgLand,
		Global:   []*Object{&Ladder},
	}
	DeadEnd5 = Object{
		In:       &Rooms,
		LongDesc: "You have come to a dead end in the mine.",
		Desc:     "Dead End",
		Flags:    FlgLand,
	}
	TimberRoom = Object{
		In:       &Rooms,
		LongDesc: "This is a long and narrow passage, which is cluttered with broken timbers. A wide passage comes from the east and turns at the west end of the room into a very narrow passageway. From the west comes a strong draft.",
		Desc:     "Timber Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags: FlgLand | FlgSacred,
	}
	LowerShaft = Object{
		In:       &Rooms,
		LongDesc: "This is a small drafty room in which is the bottom of a long shaft. To the south is a passageway and to the east a very narrow passage. In the shaft can be seen a heavy iron chain.",
		Desc:     "Drafty Room",
		Flags:    FlgLand | FlgSacred,
		Pseudo:   []PseudoObj{{Synonym: "chain", Action: ChainPseudo}},
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	MachineRoom = Object{
		In:    &Rooms,
		Desc:  "Machine Room",
		Flags: FlgLand,
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// ================================================================
	// ROOMS - Coal Mine
	// ================================================================

	Mine1 = Object{
		In:       &Rooms,
		LongDesc: "This is a nondescript part of a coal mine.",
		Desc:     "Coal Mine",
		Flags:    FlgLand,
	}
	Mine2 = Object{
		In:       &Rooms,
		LongDesc: "This is a nondescript part of a coal mine.",
		Desc:     "Coal Mine",
		Flags:    FlgLand,
	}
	Mine3 = Object{
		In:       &Rooms,
		LongDesc: "This is a nondescript part of a coal mine.",
		Desc:     "Coal Mine",
		Flags:    FlgLand,
	}
	Mine4 = Object{
		In:       &Rooms,
		LongDesc: "This is a nondescript part of a coal mine.",
		Desc:     "Coal Mine",
		Flags:    FlgLand,
	}
	SlideRoom = Object{
		In:       &Rooms,
		LongDesc: "This is a small chamber, which appears to have been part of a coal mine. On the south wall of the chamber the letters \"Granite Wall\" are etched in the rock. To the east is a long passage, and there is a steep metal slide twisting downward. To the north is a small opening.",
		Desc:     "Slide Room",
		Flags:    FlgLand,
		Global:   []*Object{&Slide},
	}
)
