package zork

var (
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
	// ZIL: <LTABLE (PURE) RIVER-1 4 RIVER-2 4 RIVER-3 3 RIVER-4 2 RIVER-5 1>
	RiverSpeedMap = map[*Object]int{
		&River1: 4,
		&River2: 4,
		&River3: 3,
		&River4: 2,
		&River5: 1,
	}
	// RiverNext maps each river room to the next one downstream.
	RiverNext = map[*Object]*Object{
		&River1: &River2,
		&River2: &River3,
		&River3: &River4,
		&River4: &River5,
	}
	// RiverLaunch maps a room to the river room you enter when launching from it.
	RiverLaunch = map[*Object]*Object{
		&DamBase:          &River1,
		&WhiteCliffsNorth: &River3,
		&WhiteCliffsSouth: &River4,
		&Shore:            &River5,
		&SandyBeach:       &River4,
		&ReservoirSouth:   &Reservoir,
		&ReservoirNorth:   &Reservoir,
		&StreamView:       &InStream,
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
	// TABLES
	// ================================================================

	// InHouseAround maps each interior room to the next in a circular walk.
	InHouseAround = map[*Object]*Object{
		&LivingRoom: &Kitchen,
		&Kitchen:    &Attic,
		&Attic:      &Kitchen,
	}
	// HouseAround maps each exterior room to the next in a circular walk.
	HouseAround = map[*Object]*Object{
		&WestOfHouse:  &NorthOfHouse,
		&NorthOfHouse: &EastOfHouse,
		&EastOfHouse:  &SouthOfHouse,
		&SouthOfHouse: &WestOfHouse,
	}
	// ForestAround maps each forest room to the next in a circular walk.
	ForestAround = map[*Object]*Object{
		&Forest1:   &Forest2,
		&Forest2:   &Forest3,
		&Forest3:   &Path,
		&Path:      &Clearing,
		&Clearing:  &Forest1,
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
	WestOfHouse.SetExit(North, DirProps{UExit: true, RExit: &NorthOfHouse})
	WestOfHouse.SetExit(South, DirProps{UExit: true, RExit: &SouthOfHouse})
	WestOfHouse.SetExit(NorthEast, DirProps{UExit: true, RExit: &NorthOfHouse})
	WestOfHouse.SetExit(SouthEast, DirProps{UExit: true, RExit: &SouthOfHouse})
	WestOfHouse.SetExit(West, DirProps{UExit: true, RExit: &Forest1})
	WestOfHouse.SetExit(East, DirProps{NExit: "The door is boarded and you can't remove the boards."})
	WestOfHouse.SetExit(SouthWest, DirProps{CExit: func() bool { return G.WonGame }, RExit: &StoneBarrow})
	WestOfHouse.SetExit(In, DirProps{CExit: func() bool { return G.WonGame }, RExit: &StoneBarrow})

	// Stone Barrow
	StoneBarrow.SetExit(NorthEast, DirProps{UExit: true, RExit: &WestOfHouse})

	// North of House
	NorthOfHouse.SetExit(SouthWest, DirProps{UExit: true, RExit: &WestOfHouse})
	NorthOfHouse.SetExit(SouthEast, DirProps{UExit: true, RExit: &EastOfHouse})
	NorthOfHouse.SetExit(West, DirProps{UExit: true, RExit: &WestOfHouse})
	NorthOfHouse.SetExit(East, DirProps{UExit: true, RExit: &EastOfHouse})
	NorthOfHouse.SetExit(North, DirProps{UExit: true, RExit: &Path})
	NorthOfHouse.SetExit(South, DirProps{NExit: "The windows are all boarded."})

	// South of House
	SouthOfHouse.SetExit(West, DirProps{UExit: true, RExit: &WestOfHouse})
	SouthOfHouse.SetExit(East, DirProps{UExit: true, RExit: &EastOfHouse})
	SouthOfHouse.SetExit(NorthEast, DirProps{UExit: true, RExit: &EastOfHouse})
	SouthOfHouse.SetExit(NorthWest, DirProps{UExit: true, RExit: &WestOfHouse})
	SouthOfHouse.SetExit(South, DirProps{UExit: true, RExit: &Forest3})
	SouthOfHouse.SetExit(North, DirProps{NExit: "The windows are all boarded."})

	// East of House
	EastOfHouse.SetExit(North, DirProps{UExit: true, RExit: &NorthOfHouse})
	EastOfHouse.SetExit(South, DirProps{UExit: true, RExit: &SouthOfHouse})
	EastOfHouse.SetExit(SouthWest, DirProps{UExit: true, RExit: &SouthOfHouse})
	EastOfHouse.SetExit(NorthWest, DirProps{UExit: true, RExit: &NorthOfHouse})
	EastOfHouse.SetExit(East, DirProps{UExit: true, RExit: &Clearing})
	EastOfHouse.SetExit(West, DirProps{DExit: &KitchenWindow, RExit: &Kitchen})
	EastOfHouse.SetExit(In, DirProps{DExit: &KitchenWindow, RExit: &Kitchen})

	// Forest 1
	Forest1.SetExit(Up, DirProps{NExit: "There is no tree here suitable for climbing."})
	Forest1.SetExit(North, DirProps{UExit: true, RExit: &GratingClearing})
	Forest1.SetExit(East, DirProps{UExit: true, RExit: &Path})
	Forest1.SetExit(South, DirProps{UExit: true, RExit: &Forest3})
	Forest1.SetExit(West, DirProps{NExit: "You would need a machete to go further west."})

	// Forest 2
	Forest2.SetExit(Up, DirProps{NExit: "There is no tree here suitable for climbing."})
	Forest2.SetExit(North, DirProps{NExit: "The forest becomes impenetrable to the north."})
	Forest2.SetExit(East, DirProps{UExit: true, RExit: &Mountains})
	Forest2.SetExit(South, DirProps{UExit: true, RExit: &Clearing})
	Forest2.SetExit(West, DirProps{UExit: true, RExit: &Path})

	// Mountains
	Mountains.SetExit(Up, DirProps{NExit: "The mountains are impassable."})
	Mountains.SetExit(North, DirProps{UExit: true, RExit: &Forest2})
	Mountains.SetExit(East, DirProps{NExit: "The mountains are impassable."})
	Mountains.SetExit(South, DirProps{UExit: true, RExit: &Forest2})
	Mountains.SetExit(West, DirProps{UExit: true, RExit: &Forest2})

	// Forest 3
	Forest3.SetExit(Up, DirProps{NExit: "There is no tree here suitable for climbing."})
	Forest3.SetExit(North, DirProps{UExit: true, RExit: &Clearing})
	Forest3.SetExit(East, DirProps{NExit: "The rank undergrowth prevents eastward movement."})
	Forest3.SetExit(South, DirProps{NExit: "Storm-tossed trees block your way."})
	Forest3.SetExit(West, DirProps{UExit: true, RExit: &Forest1})
	Forest3.SetExit(NorthWest, DirProps{UExit: true, RExit: &SouthOfHouse})

	// Path
	Path.SetExit(Up, DirProps{UExit: true, RExit: &UpATree})
	Path.SetExit(North, DirProps{UExit: true, RExit: &GratingClearing})
	Path.SetExit(East, DirProps{UExit: true, RExit: &Forest2})
	Path.SetExit(South, DirProps{UExit: true, RExit: &NorthOfHouse})
	Path.SetExit(West, DirProps{UExit: true, RExit: &Forest1})

	// Up a Tree
	UpATree.SetExit(Down, DirProps{UExit: true, RExit: &Path})
	UpATree.SetExit(Up, DirProps{NExit: "You cannot climb any higher."})

	// Grating Clearing
	GratingClearing.SetExit(North, DirProps{NExit: "The forest becomes impenetrable to the north."})
	GratingClearing.SetExit(East, DirProps{UExit: true, RExit: &Forest2})
	GratingClearing.SetExit(West, DirProps{UExit: true, RExit: &Forest1})
	GratingClearing.SetExit(South, DirProps{UExit: true, RExit: &Path})
	GratingClearing.SetExit(Down, DirProps{FExit: GratingExitFcn})

	// Clearing
	Clearing.SetExit(Up, DirProps{NExit: "There is no tree here suitable for climbing."})
	Clearing.SetExit(East, DirProps{UExit: true, RExit: &CanyonView})
	Clearing.SetExit(North, DirProps{UExit: true, RExit: &Forest2})
	Clearing.SetExit(South, DirProps{UExit: true, RExit: &Forest3})
	Clearing.SetExit(West, DirProps{UExit: true, RExit: &EastOfHouse})

	// Kitchen
	Kitchen.SetExit(East, DirProps{DExit: &KitchenWindow, RExit: &EastOfHouse})
	Kitchen.SetExit(West, DirProps{UExit: true, RExit: &LivingRoom})
	Kitchen.SetExit(Out, DirProps{DExit: &KitchenWindow, RExit: &EastOfHouse})
	Kitchen.SetExit(Up, DirProps{UExit: true, RExit: &Attic})
	Kitchen.SetExit(Down, DirProps{CExit: func() bool { return false }, RExit: &Studio, CExitStr: "Only Santa Claus climbs down chimneys."})

	// Attic
	Attic.SetExit(Down, DirProps{UExit: true, RExit: &Kitchen})

	// Living Room
	LivingRoom.SetExit(East, DirProps{UExit: true, RExit: &Kitchen})
	LivingRoom.SetExit(West, DirProps{CExit: func() bool { return G.MagicFlag }, RExit: &StrangePassage, CExitStr: "The door is nailed shut."})
	LivingRoom.SetExit(Down, DirProps{FExit: TrapDoorExitFcn})

	// Cellar
	Cellar.SetExit(North, DirProps{UExit: true, RExit: &TrollRoom})
	Cellar.SetExit(South, DirProps{UExit: true, RExit: &EastOfChasm})
	Cellar.SetExit(Up, DirProps{DExit: &TrapDoor, RExit: &LivingRoom})
	Cellar.SetExit(West, DirProps{NExit: "You try to ascend the ramp, but it is impossible, and you slide back down."})

	// Troll Room
	TrollRoom.SetExit(South, DirProps{UExit: true, RExit: &Cellar})
	TrollRoom.SetExit(East, DirProps{CExit: func() bool { return G.TrollFlag }, RExit: &EWPassage, CExitStr: "The troll fends you off with a menacing gesture."})
	TrollRoom.SetExit(West, DirProps{CExit: func() bool { return G.TrollFlag }, RExit: &Maze1, CExitStr: "The troll fends you off with a menacing gesture."})

	// East of Chasm
	EastOfChasm.SetExit(North, DirProps{UExit: true, RExit: &Cellar})
	EastOfChasm.SetExit(East, DirProps{UExit: true, RExit: &Gallery})
	EastOfChasm.SetExit(Down, DirProps{NExit: "The chasm probably leads straight to the infernal regions."})

	// Gallery
	Gallery.SetExit(West, DirProps{UExit: true, RExit: &EastOfChasm})
	Gallery.SetExit(North, DirProps{UExit: true, RExit: &Studio})

	// Studio
	Studio.SetExit(South, DirProps{UExit: true, RExit: &Gallery})
	Studio.SetExit(Up, DirProps{FExit: UpChimneyFcn})

	// Maze 1
	Maze1.SetExit(East, DirProps{UExit: true, RExit: &TrollRoom})
	Maze1.SetExit(North, DirProps{UExit: true, RExit: &Maze1})
	Maze1.SetExit(South, DirProps{UExit: true, RExit: &Maze2})
	Maze1.SetExit(West, DirProps{UExit: true, RExit: &Maze4})

	// Maze 2
	Maze2.SetExit(South, DirProps{UExit: true, RExit: &Maze1})
	Maze2.SetExit(Down, DirProps{FExit: MazeDiodesFcn})
	Maze2.SetExit(East, DirProps{UExit: true, RExit: &Maze3})

	// Maze 3
	Maze3.SetExit(West, DirProps{UExit: true, RExit: &Maze2})
	Maze3.SetExit(North, DirProps{UExit: true, RExit: &Maze4})
	Maze3.SetExit(Up, DirProps{UExit: true, RExit: &Maze5})

	// Maze 4
	Maze4.SetExit(West, DirProps{UExit: true, RExit: &Maze3})
	Maze4.SetExit(North, DirProps{UExit: true, RExit: &Maze1})
	Maze4.SetExit(East, DirProps{UExit: true, RExit: &DeadEnd1})

	// Dead End 1
	DeadEnd1.SetExit(South, DirProps{UExit: true, RExit: &Maze4})

	// Maze 5
	Maze5.SetExit(East, DirProps{UExit: true, RExit: &DeadEnd2})
	Maze5.SetExit(North, DirProps{UExit: true, RExit: &Maze3})
	Maze5.SetExit(SouthWest, DirProps{UExit: true, RExit: &Maze6})

	// Dead End 2
	DeadEnd2.SetExit(West, DirProps{UExit: true, RExit: &Maze5})

	// Maze 6
	Maze6.SetExit(Down, DirProps{UExit: true, RExit: &Maze5})
	Maze6.SetExit(East, DirProps{UExit: true, RExit: &Maze7})
	Maze6.SetExit(West, DirProps{UExit: true, RExit: &Maze6})
	Maze6.SetExit(Up, DirProps{UExit: true, RExit: &Maze9})

	// Maze 7
	Maze7.SetExit(Up, DirProps{UExit: true, RExit: &Maze14})
	Maze7.SetExit(West, DirProps{UExit: true, RExit: &Maze6})
	Maze7.SetExit(Down, DirProps{FExit: MazeDiodesFcn})
	Maze7.SetExit(East, DirProps{UExit: true, RExit: &Maze8})
	Maze7.SetExit(South, DirProps{UExit: true, RExit: &Maze15})

	// Maze 8
	Maze8.SetExit(NorthEast, DirProps{UExit: true, RExit: &Maze7})
	Maze8.SetExit(West, DirProps{UExit: true, RExit: &Maze8})
	Maze8.SetExit(SouthEast, DirProps{UExit: true, RExit: &DeadEnd3})

	// Dead End 3
	DeadEnd3.SetExit(North, DirProps{UExit: true, RExit: &Maze8})

	// Maze 9
	Maze9.SetExit(North, DirProps{UExit: true, RExit: &Maze6})
	Maze9.SetExit(Down, DirProps{FExit: MazeDiodesFcn})
	Maze9.SetExit(East, DirProps{UExit: true, RExit: &Maze10})
	Maze9.SetExit(South, DirProps{UExit: true, RExit: &Maze13})
	Maze9.SetExit(West, DirProps{UExit: true, RExit: &Maze12})
	Maze9.SetExit(NorthWest, DirProps{UExit: true, RExit: &Maze9})

	// Maze 10
	Maze10.SetExit(East, DirProps{UExit: true, RExit: &Maze9})
	Maze10.SetExit(West, DirProps{UExit: true, RExit: &Maze13})
	Maze10.SetExit(Up, DirProps{UExit: true, RExit: &Maze11})

	// Maze 11
	Maze11.SetExit(NorthEast, DirProps{UExit: true, RExit: &GratingRoom})
	Maze11.SetExit(Down, DirProps{UExit: true, RExit: &Maze10})
	Maze11.SetExit(NorthWest, DirProps{UExit: true, RExit: &Maze13})
	Maze11.SetExit(SouthWest, DirProps{UExit: true, RExit: &Maze12})

	// Grating Room
	GratingRoom.SetExit(SouthWest, DirProps{UExit: true, RExit: &Maze11})
	GratingRoom.SetExit(Up, DirProps{DExit: &Grate, RExit: &GratingClearing, DExitStr: "The grating is closed."})

	// Maze 12
	Maze12.SetExit(Down, DirProps{FExit: MazeDiodesFcn})
	Maze12.SetExit(SouthWest, DirProps{UExit: true, RExit: &Maze11})
	Maze12.SetExit(East, DirProps{UExit: true, RExit: &Maze13})
	Maze12.SetExit(Up, DirProps{UExit: true, RExit: &Maze9})
	Maze12.SetExit(North, DirProps{UExit: true, RExit: &DeadEnd4})

	// Dead End 4
	DeadEnd4.SetExit(South, DirProps{UExit: true, RExit: &Maze12})

	// Maze 13
	Maze13.SetExit(East, DirProps{UExit: true, RExit: &Maze9})
	Maze13.SetExit(Down, DirProps{UExit: true, RExit: &Maze12})
	Maze13.SetExit(South, DirProps{UExit: true, RExit: &Maze10})
	Maze13.SetExit(West, DirProps{UExit: true, RExit: &Maze11})

	// Maze 14
	Maze14.SetExit(West, DirProps{UExit: true, RExit: &Maze15})
	Maze14.SetExit(NorthWest, DirProps{UExit: true, RExit: &Maze14})
	Maze14.SetExit(NorthEast, DirProps{UExit: true, RExit: &Maze7})
	Maze14.SetExit(South, DirProps{UExit: true, RExit: &Maze7})

	// Maze 15
	Maze15.SetExit(West, DirProps{UExit: true, RExit: &Maze14})
	Maze15.SetExit(South, DirProps{UExit: true, RExit: &Maze7})
	Maze15.SetExit(SouthEast, DirProps{UExit: true, RExit: &CyclopsRoom})

	// Cyclops Room
	CyclopsRoom.SetExit(NorthWest, DirProps{UExit: true, RExit: &Maze15})
	CyclopsRoom.SetExit(East, DirProps{CExit: func() bool { return G.MagicFlag }, RExit: &StrangePassage, CExitStr: "The east wall is solid rock."})
	CyclopsRoom.SetExit(Up, DirProps{CExit: func() bool { return G.CyclopsFlag }, RExit: &TreasureRoom, CExitStr: "The cyclops doesn't look like he'll let you past."})

	// Strange Passage
	StrangePassage.SetExit(West, DirProps{UExit: true, RExit: &CyclopsRoom})
	StrangePassage.SetExit(In, DirProps{UExit: true, RExit: &CyclopsRoom})
	StrangePassage.SetExit(East, DirProps{UExit: true, RExit: &LivingRoom})

	// Treasure Room
	TreasureRoom.SetExit(Down, DirProps{UExit: true, RExit: &CyclopsRoom})

	// Reservoir South
	ReservoirSouth.SetExit(SouthEast, DirProps{UExit: true, RExit: &DeepCanyon})
	ReservoirSouth.SetExit(SouthWest, DirProps{UExit: true, RExit: &ChasmRoom})
	ReservoirSouth.SetExit(East, DirProps{UExit: true, RExit: &DamRoom})
	ReservoirSouth.SetExit(West, DirProps{UExit: true, RExit: &StreamView})
	ReservoirSouth.SetExit(North, DirProps{CExit: func() bool { return G.LowTide }, RExit: &Reservoir, CExitStr: "You would drown."})

	// Reservoir
	Reservoir.SetExit(North, DirProps{UExit: true, RExit: &ReservoirNorth})
	Reservoir.SetExit(South, DirProps{UExit: true, RExit: &ReservoirSouth})
	Reservoir.SetExit(Up, DirProps{UExit: true, RExit: &InStream})
	Reservoir.SetExit(West, DirProps{UExit: true, RExit: &InStream})
	Reservoir.SetExit(Down, DirProps{NExit: "The dam blocks your way."})

	// Reservoir North
	ReservoirNorth.SetExit(North, DirProps{UExit: true, RExit: &AtlantisRoom})
	ReservoirNorth.SetExit(South, DirProps{CExit: func() bool { return G.LowTide }, RExit: &Reservoir, CExitStr: "You would drown."})

	// Stream View
	StreamView.SetExit(East, DirProps{UExit: true, RExit: &ReservoirSouth})
	StreamView.SetExit(West, DirProps{NExit: "The stream emerges from a spot too small for you to enter."})

	// In Stream
	InStream.SetExit(Up, DirProps{NExit: "The channel is too narrow."})
	InStream.SetExit(West, DirProps{NExit: "The channel is too narrow."})
	InStream.SetExit(Land, DirProps{UExit: true, RExit: &StreamView})
	InStream.SetExit(Down, DirProps{UExit: true, RExit: &Reservoir})
	InStream.SetExit(East, DirProps{UExit: true, RExit: &Reservoir})

	// Mirror Room 1
	MirrorRoom1.SetExit(North, DirProps{UExit: true, RExit: &ColdPassage})
	MirrorRoom1.SetExit(West, DirProps{UExit: true, RExit: &TwistingPassage})
	MirrorRoom1.SetExit(East, DirProps{UExit: true, RExit: &SmallCave})

	// Mirror Room 2
	MirrorRoom2.SetExit(West, DirProps{UExit: true, RExit: &WindingPassage})
	MirrorRoom2.SetExit(North, DirProps{UExit: true, RExit: &NarrowPassage})
	MirrorRoom2.SetExit(East, DirProps{UExit: true, RExit: &TinyCave})

	// Small Cave
	SmallCave.SetExit(North, DirProps{UExit: true, RExit: &MirrorRoom1})
	SmallCave.SetExit(Down, DirProps{UExit: true, RExit: &AtlantisRoom})
	SmallCave.SetExit(South, DirProps{UExit: true, RExit: &AtlantisRoom})
	SmallCave.SetExit(West, DirProps{UExit: true, RExit: &TwistingPassage})

	// Tiny Cave
	TinyCave.SetExit(North, DirProps{UExit: true, RExit: &MirrorRoom2})
	TinyCave.SetExit(West, DirProps{UExit: true, RExit: &WindingPassage})
	TinyCave.SetExit(Down, DirProps{UExit: true, RExit: &EnteranceToHades})

	// Cold Passage
	ColdPassage.SetExit(South, DirProps{UExit: true, RExit: &MirrorRoom1})
	ColdPassage.SetExit(West, DirProps{UExit: true, RExit: &SlideRoom})

	// Narrow Passage
	NarrowPassage.SetExit(North, DirProps{UExit: true, RExit: &RoundRoom})
	NarrowPassage.SetExit(South, DirProps{UExit: true, RExit: &MirrorRoom2})

	// Winding Passage
	WindingPassage.SetExit(North, DirProps{UExit: true, RExit: &MirrorRoom2})
	WindingPassage.SetExit(East, DirProps{UExit: true, RExit: &TinyCave})

	// Twisting Passage
	TwistingPassage.SetExit(North, DirProps{UExit: true, RExit: &MirrorRoom1})
	TwistingPassage.SetExit(East, DirProps{UExit: true, RExit: &SmallCave})

	// Atlantis Room
	AtlantisRoom.SetExit(Up, DirProps{UExit: true, RExit: &SmallCave})
	AtlantisRoom.SetExit(South, DirProps{UExit: true, RExit: &ReservoirNorth})

	// EW Passage
	EWPassage.SetExit(East, DirProps{UExit: true, RExit: &RoundRoom})
	EWPassage.SetExit(West, DirProps{UExit: true, RExit: &TrollRoom})
	EWPassage.SetExit(Down, DirProps{UExit: true, RExit: &ChasmRoom})
	EWPassage.SetExit(North, DirProps{UExit: true, RExit: &ChasmRoom})

	// Round Room
	RoundRoom.SetExit(East, DirProps{UExit: true, RExit: &LoudRoom})
	RoundRoom.SetExit(West, DirProps{UExit: true, RExit: &EWPassage})
	RoundRoom.SetExit(North, DirProps{UExit: true, RExit: &NSPassage})
	RoundRoom.SetExit(South, DirProps{UExit: true, RExit: &NarrowPassage})
	RoundRoom.SetExit(SouthEast, DirProps{UExit: true, RExit: &EngravingsCave})

	// Deep Canyon
	DeepCanyon.SetExit(NorthWest, DirProps{UExit: true, RExit: &ReservoirSouth})
	DeepCanyon.SetExit(East, DirProps{UExit: true, RExit: &DamRoom})
	DeepCanyon.SetExit(SouthWest, DirProps{UExit: true, RExit: &NSPassage})
	DeepCanyon.SetExit(Down, DirProps{UExit: true, RExit: &LoudRoom})

	// Damp Cave
	DampCave.SetExit(West, DirProps{UExit: true, RExit: &LoudRoom})
	DampCave.SetExit(East, DirProps{UExit: true, RExit: &WhiteCliffsNorth})
	DampCave.SetExit(South, DirProps{NExit: "It is too narrow for most insects."})

	// Loud Room
	LoudRoom.SetExit(East, DirProps{UExit: true, RExit: &DampCave})
	LoudRoom.SetExit(West, DirProps{UExit: true, RExit: &RoundRoom})
	LoudRoom.SetExit(Up, DirProps{UExit: true, RExit: &DeepCanyon})

	// NS Passage
	NSPassage.SetExit(North, DirProps{UExit: true, RExit: &ChasmRoom})
	NSPassage.SetExit(NorthEast, DirProps{UExit: true, RExit: &DeepCanyon})
	NSPassage.SetExit(South, DirProps{UExit: true, RExit: &RoundRoom})

	// Chasm Room
	ChasmRoom.SetExit(NorthEast, DirProps{UExit: true, RExit: &ReservoirSouth})
	ChasmRoom.SetExit(SouthWest, DirProps{UExit: true, RExit: &EWPassage})
	ChasmRoom.SetExit(Up, DirProps{UExit: true, RExit: &EWPassage})
	ChasmRoom.SetExit(South, DirProps{UExit: true, RExit: &NSPassage})
	ChasmRoom.SetExit(Down, DirProps{NExit: "Are you out of your mind?"})

	// Entrance to Hades
	EnteranceToHades.SetExit(Up, DirProps{UExit: true, RExit: &TinyCave})
	EnteranceToHades.SetExit(In, DirProps{CExit: func() bool { return G.LLDFlag }, RExit: &LandOfLivingDead, CExitStr: "Some invisible force prevents you from passing through the gate."})
	EnteranceToHades.SetExit(South, DirProps{CExit: func() bool { return G.LLDFlag }, RExit: &LandOfLivingDead, CExitStr: "Some invisible force prevents you from passing through the gate."})

	// Land of Living Dead
	LandOfLivingDead.SetExit(Out, DirProps{UExit: true, RExit: &EnteranceToHades})
	LandOfLivingDead.SetExit(North, DirProps{UExit: true, RExit: &EnteranceToHades})

	// Engravings Cave
	EngravingsCave.SetExit(NorthWest, DirProps{UExit: true, RExit: &RoundRoom})
	EngravingsCave.SetExit(East, DirProps{UExit: true, RExit: &DomeRoom})

	// Egypt Room
	EgyptRoom.SetExit(West, DirProps{UExit: true, RExit: &NorthTemple})
	EgyptRoom.SetExit(Up, DirProps{UExit: true, RExit: &NorthTemple})

	// Dome Room
	DomeRoom.SetExit(West, DirProps{UExit: true, RExit: &EngravingsCave})
	DomeRoom.SetExit(Down, DirProps{CExit: func() bool { return G.DomeFlag }, RExit: &TorchRoom, CExitStr: "You cannot go down without fracturing many bones."})

	// Torch Room
	TorchRoom.SetExit(Up, DirProps{NExit: "You cannot reach the rope."})
	TorchRoom.SetExit(South, DirProps{UExit: true, RExit: &NorthTemple})
	TorchRoom.SetExit(Down, DirProps{UExit: true, RExit: &NorthTemple})

	// North Temple
	NorthTemple.SetExit(Down, DirProps{UExit: true, RExit: &EgyptRoom})
	NorthTemple.SetExit(East, DirProps{UExit: true, RExit: &EgyptRoom})
	NorthTemple.SetExit(North, DirProps{UExit: true, RExit: &TorchRoom})
	NorthTemple.SetExit(Out, DirProps{UExit: true, RExit: &TorchRoom})
	NorthTemple.SetExit(Up, DirProps{UExit: true, RExit: &TorchRoom})
	NorthTemple.SetExit(South, DirProps{UExit: true, RExit: &SouthTemple})

	// South Temple
	SouthTemple.SetExit(North, DirProps{UExit: true, RExit: &NorthTemple})
	SouthTemple.SetExit(Down, DirProps{CExit: func() bool { return G.CoffinCure }, RExit: &TinyCave, CExitStr: "You haven't a prayer of getting the coffin down there."})

	// Dam Room
	DamRoom.SetExit(South, DirProps{UExit: true, RExit: &DeepCanyon})
	DamRoom.SetExit(Down, DirProps{UExit: true, RExit: &DamBase})
	DamRoom.SetExit(East, DirProps{UExit: true, RExit: &DamBase})
	DamRoom.SetExit(North, DirProps{UExit: true, RExit: &DamLobby})
	DamRoom.SetExit(West, DirProps{UExit: true, RExit: &ReservoirSouth})

	// Dam Lobby
	DamLobby.SetExit(South, DirProps{UExit: true, RExit: &DamRoom})
	DamLobby.SetExit(North, DirProps{UExit: true, RExit: &MaintenanceRoom})
	DamLobby.SetExit(East, DirProps{UExit: true, RExit: &MaintenanceRoom})

	// Maintenance Room
	MaintenanceRoom.SetExit(South, DirProps{UExit: true, RExit: &DamLobby})
	MaintenanceRoom.SetExit(West, DirProps{UExit: true, RExit: &DamLobby})

	// Dam Base
	DamBase.SetExit(North, DirProps{UExit: true, RExit: &DamRoom})
	DamBase.SetExit(Up, DirProps{UExit: true, RExit: &DamRoom})

	// River 1
	River1.SetExit(Up, DirProps{NExit: "You cannot go upstream due to strong currents."})
	River1.SetExit(West, DirProps{UExit: true, RExit: &DamBase})
	River1.SetExit(Land, DirProps{UExit: true, RExit: &DamBase})
	River1.SetExit(Down, DirProps{UExit: true, RExit: &River2})
	River1.SetExit(East, DirProps{NExit: "The White Cliffs prevent your landing here."})

	// River 2
	River2.SetExit(Up, DirProps{NExit: "You cannot go upstream due to strong currents."})
	River2.SetExit(Down, DirProps{UExit: true, RExit: &River3})
	River2.SetExit(Land, DirProps{NExit: "There is no safe landing spot here."})
	River2.SetExit(East, DirProps{NExit: "The White Cliffs prevent your landing here."})
	River2.SetExit(West, DirProps{NExit: "Just in time you steer away from the rocks."})

	// River 3
	River3.SetExit(Up, DirProps{NExit: "You cannot go upstream due to strong currents."})
	River3.SetExit(Down, DirProps{UExit: true, RExit: &River4})
	River3.SetExit(Land, DirProps{UExit: true, RExit: &WhiteCliffsNorth})
	River3.SetExit(West, DirProps{UExit: true, RExit: &WhiteCliffsNorth})

	// White Cliffs North
	WhiteCliffsNorth.SetExit(South, DirProps{CExit: func() bool { return G.DeflateFlag }, RExit: &WhiteCliffsSouth, CExitStr: "The path is too narrow."})
	WhiteCliffsNorth.SetExit(West, DirProps{CExit: func() bool { return G.DeflateFlag }, RExit: &DampCave, CExitStr: "The path is too narrow."})

	// White Cliffs South
	WhiteCliffsSouth.SetExit(North, DirProps{CExit: func() bool { return G.DeflateFlag }, RExit: &WhiteCliffsNorth, CExitStr: "The path is too narrow."})

	// River 4
	River4.SetExit(Up, DirProps{NExit: "You cannot go upstream due to strong currents."})
	River4.SetExit(Down, DirProps{UExit: true, RExit: &River5})
	River4.SetExit(Land, DirProps{NExit: "You can land either to the east or the west."})
	River4.SetExit(West, DirProps{UExit: true, RExit: &WhiteCliffsSouth})
	River4.SetExit(East, DirProps{UExit: true, RExit: &SandyBeach})

	// River 5
	River5.SetExit(Up, DirProps{NExit: "You cannot go upstream due to strong currents."})
	River5.SetExit(East, DirProps{UExit: true, RExit: &Shore})
	River5.SetExit(Land, DirProps{UExit: true, RExit: &Shore})

	// Shore
	Shore.SetExit(North, DirProps{UExit: true, RExit: &SandyBeach})
	Shore.SetExit(South, DirProps{UExit: true, RExit: &AragainFalls})

	// Sandy Beach
	SandyBeach.SetExit(NorthEast, DirProps{UExit: true, RExit: &SandyCave})
	SandyBeach.SetExit(South, DirProps{UExit: true, RExit: &Shore})

	// Sandy Cave
	SandyCave.SetExit(SouthWest, DirProps{UExit: true, RExit: &SandyBeach})

	// Aragain Falls
	AragainFalls.SetExit(West, DirProps{CExit: func() bool { return G.RainbowFlag }, RExit: &OnRainbow})
	AragainFalls.SetExit(Down, DirProps{NExit: "It's a long way..."})
	AragainFalls.SetExit(North, DirProps{UExit: true, RExit: &Shore})
	AragainFalls.SetExit(Up, DirProps{CExit: func() bool { return G.RainbowFlag }, RExit: &OnRainbow})

	// On Rainbow
	OnRainbow.SetExit(West, DirProps{UExit: true, RExit: &EndOfRainbow})
	OnRainbow.SetExit(East, DirProps{UExit: true, RExit: &AragainFalls})

	// End of Rainbow
	EndOfRainbow.SetExit(Up, DirProps{CExit: func() bool { return G.RainbowFlag }, RExit: &OnRainbow})
	EndOfRainbow.SetExit(NorthEast, DirProps{CExit: func() bool { return G.RainbowFlag }, RExit: &OnRainbow})
	EndOfRainbow.SetExit(East, DirProps{CExit: func() bool { return G.RainbowFlag }, RExit: &OnRainbow})
	EndOfRainbow.SetExit(SouthWest, DirProps{UExit: true, RExit: &CanyonBottom})

	// Canyon Bottom
	CanyonBottom.SetExit(Up, DirProps{UExit: true, RExit: &CliffMiddle})
	CanyonBottom.SetExit(North, DirProps{UExit: true, RExit: &EndOfRainbow})

	// Cliff Middle
	CliffMiddle.SetExit(Up, DirProps{UExit: true, RExit: &CanyonView})
	CliffMiddle.SetExit(Down, DirProps{UExit: true, RExit: &CanyonBottom})

	// Canyon View
	CanyonView.SetExit(East, DirProps{UExit: true, RExit: &CliffMiddle})
	CanyonView.SetExit(Down, DirProps{UExit: true, RExit: &CliffMiddle})
	CanyonView.SetExit(NorthWest, DirProps{UExit: true, RExit: &Clearing})
	CanyonView.SetExit(West, DirProps{UExit: true, RExit: &Forest3})
	CanyonView.SetExit(South, DirProps{NExit: "Storm-tossed trees block your way."})

	// Mine Entrance
	MineEntrance.SetExit(South, DirProps{UExit: true, RExit: &SlideRoom})
	MineEntrance.SetExit(In, DirProps{UExit: true, RExit: &SqueekyRoom})
	MineEntrance.SetExit(West, DirProps{UExit: true, RExit: &SqueekyRoom})

	// Squeaky Room
	SqueekyRoom.SetExit(North, DirProps{UExit: true, RExit: &BatRoom})
	SqueekyRoom.SetExit(East, DirProps{UExit: true, RExit: &MineEntrance})

	// Bat Room
	BatRoom.SetExit(South, DirProps{UExit: true, RExit: &SqueekyRoom})
	BatRoom.SetExit(East, DirProps{UExit: true, RExit: &ShaftRoom})

	// Shaft Room
	ShaftRoom.SetExit(Down, DirProps{NExit: "You wouldn't fit and would die if you could."})
	ShaftRoom.SetExit(West, DirProps{UExit: true, RExit: &BatRoom})
	ShaftRoom.SetExit(North, DirProps{UExit: true, RExit: &SmellyRoom})

	// Smelly Room
	SmellyRoom.SetExit(Down, DirProps{UExit: true, RExit: &GasRoom})
	SmellyRoom.SetExit(South, DirProps{UExit: true, RExit: &ShaftRoom})

	// Gas Room
	GasRoom.SetExit(Up, DirProps{UExit: true, RExit: &SmellyRoom})
	GasRoom.SetExit(East, DirProps{UExit: true, RExit: &Mine1})

	// Ladder Top
	LadderTop.SetExit(Down, DirProps{UExit: true, RExit: &LadderBottom})
	LadderTop.SetExit(Up, DirProps{UExit: true, RExit: &Mine4})

	// Ladder Bottom
	LadderBottom.SetExit(South, DirProps{UExit: true, RExit: &DeadEnd5})
	LadderBottom.SetExit(West, DirProps{UExit: true, RExit: &TimberRoom})
	LadderBottom.SetExit(Up, DirProps{UExit: true, RExit: &LadderTop})

	// Dead End 5
	DeadEnd5.SetExit(North, DirProps{UExit: true, RExit: &LadderBottom})

	// Timber Room
	TimberRoom.SetExit(East, DirProps{UExit: true, RExit: &LadderBottom})
	TimberRoom.SetExit(West, DirProps{CExit: func() bool { return G.EmptyHanded }, RExit: &LowerShaft, CExitStr: "You cannot fit through this passage with that load."})

	// Lower Shaft
	LowerShaft.SetExit(South, DirProps{UExit: true, RExit: &MachineRoom})
	LowerShaft.SetExit(Out, DirProps{CExit: func() bool { return G.EmptyHanded }, RExit: &TimberRoom, CExitStr: "You cannot fit through this passage with that load."})
	LowerShaft.SetExit(East, DirProps{CExit: func() bool { return G.EmptyHanded }, RExit: &TimberRoom, CExitStr: "You cannot fit through this passage with that load."})

	// Machine Room
	MachineRoom.SetExit(North, DirProps{UExit: true, RExit: &LowerShaft})

	// Mine 1
	Mine1.SetExit(North, DirProps{UExit: true, RExit: &GasRoom})
	Mine1.SetExit(East, DirProps{UExit: true, RExit: &Mine1})
	Mine1.SetExit(NorthEast, DirProps{UExit: true, RExit: &Mine2})

	// Mine 2
	Mine2.SetExit(North, DirProps{UExit: true, RExit: &Mine2})
	Mine2.SetExit(South, DirProps{UExit: true, RExit: &Mine1})
	Mine2.SetExit(SouthEast, DirProps{UExit: true, RExit: &Mine3})

	// Mine 3
	Mine3.SetExit(South, DirProps{UExit: true, RExit: &Mine3})
	Mine3.SetExit(SouthWest, DirProps{UExit: true, RExit: &Mine4})
	Mine3.SetExit(East, DirProps{UExit: true, RExit: &Mine2})

	// Mine 4
	Mine4.SetExit(North, DirProps{UExit: true, RExit: &Mine3})
	Mine4.SetExit(West, DirProps{UExit: true, RExit: &Mine4})
	Mine4.SetExit(Down, DirProps{UExit: true, RExit: &LadderTop})

	// Slide Room
	SlideRoom.SetExit(East, DirProps{UExit: true, RExit: &ColdPassage})
	SlideRoom.SetExit(North, DirProps{UExit: true, RExit: &MineEntrance})
	SlideRoom.SetExit(Down, DirProps{UExit: true, RExit: &Cellar})
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
	G.Villains = []*VillainEntry{
		{Villain: &Troll, Best: &Sword, BestAdv: 1, Prob: 0, Msgs: &TrollMelee},
		{Villain: &Thief, Best: &Knife, BestAdv: 1, Prob: 0, Msgs: &ThiefMelee},
		{Villain: &Cyclops, Best: nil, BestAdv: 0, Prob: 0, Msgs: &CyclopsMelee},
	}
}
