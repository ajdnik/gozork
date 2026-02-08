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
	// ZIL: <LTABLE (PURE) RIVER-1 RIVER-2 RIVER-3 RIVER-4 RIVER-5>
	// Lkp finds the item and returns the next element — must be a flat sequential list.
	RiverNext = []*Object{&River1, &River2, &River3, &River4, &River5}
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
	WestOfHouse.SouthWest = DirProps{CExit: func() bool { return G.WonGame }, RExit: &StoneBarrow}
	WestOfHouse.Into = DirProps{CExit: func() bool { return G.WonGame }, RExit: &StoneBarrow}

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
	LivingRoom.West = DirProps{CExit: func() bool { return G.MagicFlag }, RExit: &StrangePassage, CExitStr: "The door is nailed shut."}
	LivingRoom.Down = DirProps{FExit: TrapDoorExitFcn}

	// Cellar
	Cellar.North = DirProps{UExit: true, RExit: &TrollRoom}
	Cellar.South = DirProps{UExit: true, RExit: &EastOfChasm}
	Cellar.Up = DirProps{DExit: &TrapDoor, RExit: &LivingRoom}
	Cellar.West = DirProps{NExit: "You try to ascend the ramp, but it is impossible, and you slide back down."}

	// Troll Room
	TrollRoom.South = DirProps{UExit: true, RExit: &Cellar}
	TrollRoom.East = DirProps{CExit: func() bool { return G.TrollFlag }, RExit: &EWPassage, CExitStr: "The troll fends you off with a menacing gesture."}
	TrollRoom.West = DirProps{CExit: func() bool { return G.TrollFlag }, RExit: &Maze1, CExitStr: "The troll fends you off with a menacing gesture."}

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
	CyclopsRoom.East = DirProps{CExit: func() bool { return G.MagicFlag }, RExit: &StrangePassage, CExitStr: "The east wall is solid rock."}
	CyclopsRoom.Up = DirProps{CExit: func() bool { return G.CyclopsFlag }, RExit: &TreasureRoom, CExitStr: "The cyclops doesn't look like he'll let you past."}

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
	ReservoirSouth.North = DirProps{CExit: func() bool { return G.LowTide }, RExit: &Reservoir, CExitStr: "You would drown."}

	// Reservoir
	Reservoir.North = DirProps{UExit: true, RExit: &ReservoirNorth}
	Reservoir.South = DirProps{UExit: true, RExit: &ReservoirSouth}
	Reservoir.Up = DirProps{UExit: true, RExit: &InStream}
	Reservoir.West = DirProps{UExit: true, RExit: &InStream}
	Reservoir.Down = DirProps{NExit: "The dam blocks your way."}

	// Reservoir North
	ReservoirNorth.North = DirProps{UExit: true, RExit: &AtlantisRoom}
	ReservoirNorth.South = DirProps{CExit: func() bool { return G.LowTide }, RExit: &Reservoir, CExitStr: "You would drown."}

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
	EnteranceToHades.Into = DirProps{CExit: func() bool { return G.LLDFlag }, RExit: &LandOfLivingDead, CExitStr: "Some invisible force prevents you from passing through the gate."}
	EnteranceToHades.South = DirProps{CExit: func() bool { return G.LLDFlag }, RExit: &LandOfLivingDead, CExitStr: "Some invisible force prevents you from passing through the gate."}

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
	DomeRoom.Down = DirProps{CExit: func() bool { return G.DomeFlag }, RExit: &TorchRoom, CExitStr: "You cannot go down without fracturing many bones."}

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
	SouthTemple.Down = DirProps{CExit: func() bool { return G.CoffinCure }, RExit: &TinyCave, CExitStr: "You haven't a prayer of getting the coffin down there."}

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
	WhiteCliffsNorth.South = DirProps{CExit: func() bool { return G.DeflateFlag }, RExit: &WhiteCliffsSouth, CExitStr: "The path is too narrow."}
	WhiteCliffsNorth.West = DirProps{CExit: func() bool { return G.DeflateFlag }, RExit: &DampCave, CExitStr: "The path is too narrow."}

	// White Cliffs South
	WhiteCliffsSouth.North = DirProps{CExit: func() bool { return G.DeflateFlag }, RExit: &WhiteCliffsNorth, CExitStr: "The path is too narrow."}

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
	AragainFalls.West = DirProps{CExit: func() bool { return G.RainbowFlag }, RExit: &OnRainbow}
	AragainFalls.Down = DirProps{NExit: "It's a long way..."}
	AragainFalls.North = DirProps{UExit: true, RExit: &Shore}
	AragainFalls.Up = DirProps{CExit: func() bool { return G.RainbowFlag }, RExit: &OnRainbow}

	// On Rainbow
	OnRainbow.West = DirProps{UExit: true, RExit: &EndOfRainbow}
	OnRainbow.East = DirProps{UExit: true, RExit: &AragainFalls}

	// End of Rainbow
	EndOfRainbow.Up = DirProps{CExit: func() bool { return G.RainbowFlag }, RExit: &OnRainbow}
	EndOfRainbow.NorthEast = DirProps{CExit: func() bool { return G.RainbowFlag }, RExit: &OnRainbow}
	EndOfRainbow.East = DirProps{CExit: func() bool { return G.RainbowFlag }, RExit: &OnRainbow}
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
	TimberRoom.West = DirProps{CExit: func() bool { return G.EmptyHanded }, RExit: &LowerShaft, CExitStr: "You cannot fit through this passage with that load."}

	// Lower Shaft
	LowerShaft.South = DirProps{UExit: true, RExit: &MachineRoom}
	LowerShaft.Out = DirProps{CExit: func() bool { return G.EmptyHanded }, RExit: &TimberRoom, CExitStr: "You cannot fit through this passage with that load."}
	LowerShaft.East = DirProps{CExit: func() bool { return G.EmptyHanded }, RExit: &TimberRoom, CExitStr: "You cannot fit through this passage with that load."}

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
	G.Villains = []*VillainEntry{
		{Villain: &Troll, Best: &Sword, BestAdv: 1, Prob: 0, Msgs: &TrollMelee},
		{Villain: &Thief, Best: &Knife, BestAdv: 1, Prob: 0, Msgs: &ThiefMelee},
		{Villain: &Cyclops, Best: nil, BestAdv: 0, Prob: 0, Msgs: &CyclopsMelee},
	}
}
