package game

import . "github.com/ajdnik/gozork/engine"

var (
	// RndSelect tables
	dummy = RndSelect{
		Unselected: []string{
			"it is already closed.",
			"it is already open.",
			"it's already done.",
		},
	}
	swimYuks = RndSelect{
		Unselected: []string{
			"I don't really feel like swimming.",
			"Swimming isn't usually allowed in dungeons.",
			"You'd need a submarine to go further.",
		},
	}
	batDrops = []*Object{&mine1, &mine2, &mine3, &mine4, &ladderTop, &ladderBottom, &squeekyRoom, &mineEntrance}

	cyclomad = []string{
		"The cyclops seems somewhat agitated.",
		"The cyclops appears to be getting more agitated.",
		"The cyclops is moving about the room, looking for something.",
		"The cyclops was looking for salt and pepper. No doubt they are condiments for his upcoming snack.",
		"The cyclops is moving toward you in an unfriendly manner.",
		"You have two choices: 1. Leave  2. Become dinner.",
	}

	loudRuns = []*Object{&dampCave, &roundRoom, &deepCanyon}

	bDigs = []string{
		"You seem to be digging a hole here.",
		"The hole is getting deeper, but that's about it.",
		"You are surrounded by a wall of sand on all sides.",
	}

	drownings = []string{
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

	robberCDesc = "There is a suspicious-looking individual, holding a bag, leaning against one wall. He is armed with a vicious-looking stiletto."
	robberUDesc = "There is a suspicious-looking individual lying unconscious on the ground."

	// river tables
	// ZIL: <LTABLE (PURE) RIVER-1 4 RIVER-2 4 RIVER-3 3 RIVER-4 2 RIVER-5 1>
	riverSpeedMap = map[*Object]int{
		&river1: 4,
		&river2: 4,
		&river3: 3,
		&river4: 2,
		&river5: 1,
	}
	// riverNext maps each river room to the next one downstream.
	riverNext = map[*Object]*Object{
		&river1: &river2,
		&river2: &river3,
		&river3: &river4,
		&river4: &river5,
	}
	// riverLaunch maps a room to the river room you enter when launching from it.
	riverLaunch = map[*Object]*Object{
		&damBase:          &river1,
		&whiteCliffsNorth: &river3,
		&whiteCliffsSouth: &river4,
		&shore:            &river5,
		&sandyBeach:       &river4,
		&reservoirSouth:   &reservoir,
		&reservoirNorth:   &reservoir,
		&streamView:       &inStream,
	}

	// lamp/Candle countdown tables
	lampTable = []interface{}{
		100, "The lamp appears a bit dimmer.",
		70, "The lamp is definitely dimmer now.",
		15, "The lamp is nearly out.",
		0,
	}
	candleTable = []interface{}{
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

	// inHouseAround maps each interior room to the next in a circular walk.
	inHouseAround = map[*Object]*Object{
		&livingRoom: &kitchen,
		&kitchen:    &attic,
		&attic:      &kitchen,
	}
	// houseAround maps each exterior room to the next in a circular walk.
	houseAround = map[*Object]*Object{
		&westOfHouse:  &northOfHouse,
		&northOfHouse: &eastOfHouse,
		&eastOfHouse:  &southOfHouse,
		&southOfHouse: &westOfHouse,
	}
	// forestAround maps each forest room to the next in a circular walk.
	forestAround = map[*Object]*Object{
		&forest1:  &forest2,
		&forest2:  &forest3,
		&forest3:  &path,
		&path:     &clearing,
		&clearing: &forest1,
	}
	aboveGround = []*Object{
		&westOfHouse,
		&northOfHouse,
		&eastOfHouse,
		&southOfHouse,
		&forest1,
		&forest2,
		&forest3,
		&path,
		&clearing,
		&gratingClearing,
		&canyonView,
	}

	// ================================================================
	// MASTER OBJECT LIST
	// ================================================================

	objects = []*Object{
		// Core
		&rooms, &globalObjects, &localGlobals,
		&notHereObject, &pseudoObject,
		&hands, &me, &adventurer, &it,
		&stairs, &intnum, &blessings, &sailor, &ground, &grue, &lungs, &pathObj, &zorkmid,
		// Global objects
		&board, &teeth, &wall, &graniteWall, &songbird,
		&whiteHouse, &forest, &tree, &globalWater,
		&kitchenWindow, &chimney, &slide, &bodies, &crack, &grate,
		&ladder, &climbableCliff, &whiteCliff, &rainbow, &river,
		&boardedWindow,
		// Unplaced objects
		&inflatedBoat, &puncturedBoat, &brokenLamp, &gunk, &hotBell,
		&brokenEgg, &bauble, &diamond,
		// forest and outside rooms
		&westOfHouse, &stoneBarrow, &northOfHouse, &southOfHouse, &eastOfHouse,
		&forest1, &forest2, &mountains, &forest3, &path,
		&upATree, &gratingClearing, &clearing,
		// House rooms
		&kitchen, &attic, &livingRoom,
		// cellar and vicinity
		&cellar, &trollRoom, &eastOfChasm, &gallery, &studio,
		// Maze
		&maze1, &maze2, &maze3, &maze4, &deadEnd1,
		&maze5, &deadEnd2, &maze6, &maze7, &maze8, &deadEnd3,
		&maze9, &maze10, &maze11, &gratingRoom,
		&maze12, &deadEnd4, &maze13, &maze14, &maze15,
		// cyclops and hideaway
		&cyclopsRoom, &strangePassage, &treasureRoom,
		// reservoir area
		&reservoirSouth, &reservoir, &reservoirNorth, &streamView, &inStream,
		// Mirror rooms and vicinity
		&mirrorRoom1, &mirrorRoom2, &smallCave, &tinyCave,
		&coldPassage, &narrowPassage, &windingPassage, &twistingPassage, &atlantisRoom,
		// Round room and vicinity
		&eWPassage, &roundRoom, &deepCanyon, &dampCave, &loudRoom, &nSPassage, &chasmRoom,
		// Hades
		&enteranceToHades, &landOfLivingDead,
		// Dome, temple, egypt
		&engravingsCave, &egyptRoom, &domeRoom, &torchRoom, &northTemple, &southTemple,
		// Flood control dam
		&damRoom, &damLobby, &maintenanceRoom,
		// river area
		&damBase, &river1, &river2, &river3,
		&whiteCliffsNorth, &whiteCliffsSouth,
		&river4, &river5, &shore, &sandyBeach, &sandyCave,
		&aragainFalls, &onRainbow, &endOfRainbow,
		&canyonBottom, &cliffMiddle, &canyonView,
		// coal mine area
		&mineEntrance, &squeekyRoom, &batRoom, &shaftRoom,
		&smellyRoom, &gasRoom, &ladderTop, &ladderBottom, &deadEnd5,
		&timberRoom, &lowerShaft, &machineRoom,
		&mine1, &mine2, &mine3, &mine4, &slideRoom,
		// G.AllObjects in rooms
		&mountainRange, &frontDoor, &mailbox, &barrowDoor, &barrow,
		&trophyCase, &rug, &trapDoor, &woodenDoor, &sword, &lamp,
		&kitchenTable, &atticTable, &rope,
		&ghosts, &skull,
		&raisedBasket, &loweredBasket,
		&bat, &jade, &bell, &prayer, &altar, &candles,
		&troll, &bolt, &bubble, &dam, &controlPanel,
		&match, &guide, &inflatableBoat,
		&toolChest, &yellowButton, &brownButton, &redButton, &blueButton,
		&screwdriver, &wrench, &tube, &leak,
		&machine, &machineSwitch,
		&cyclops, &chalice, &painting, &ownersManual,
		&leaves, &nest, &sand, &scarab, &shovel,
		&coffin, &thief, &trunk, &pump, &trident,
		&mirror1, &mirror2, &railing, &pedestal, &engravings,
		&bar, &potOfGold, &buoy, &bracelet, &coal, &timbers,
		&bones, &burnedOutLantern, &bagOfCoins, &rustyKnife, &keys,
		// G.AllObjects in objects
		&mapObj, &advertisement, &bottle, &sandwichBag, &knife,
		&water, &lunch, &garlic, &book, &sceptre,
		&egg, &canary, &putty, &axe,
		&largeBag, &stiletto, &torch, &boatLabel, &emerald,
		&brokenCanary,
	}
)

// initRoomExits sets all room directional properties.
// This is done in a function to avoid circular reference issues
// between rooms that reference each other.
func initRoomExits() {
	// West of House
	westOfHouse.SetExit(North, ExitProps{UExit: true, RExit: &northOfHouse})
	westOfHouse.SetExit(South, ExitProps{UExit: true, RExit: &southOfHouse})
	westOfHouse.SetExit(NorthEast, ExitProps{UExit: true, RExit: &northOfHouse})
	westOfHouse.SetExit(SouthEast, ExitProps{UExit: true, RExit: &southOfHouse})
	westOfHouse.SetExit(West, ExitProps{UExit: true, RExit: &forest1})
	westOfHouse.SetExit(East, ExitProps{NExit: "The door is boarded and you can't remove the boards."})
	westOfHouse.SetExit(SouthWest, ExitProps{CExit: func() bool { return gD().WonGame }, RExit: &stoneBarrow})
	westOfHouse.SetExit(In, ExitProps{CExit: func() bool { return gD().WonGame }, RExit: &stoneBarrow})

	// Stone barrow
	stoneBarrow.SetExit(NorthEast, ExitProps{UExit: true, RExit: &westOfHouse})

	// North of House
	northOfHouse.SetExit(SouthWest, ExitProps{UExit: true, RExit: &westOfHouse})
	northOfHouse.SetExit(SouthEast, ExitProps{UExit: true, RExit: &eastOfHouse})
	northOfHouse.SetExit(West, ExitProps{UExit: true, RExit: &westOfHouse})
	northOfHouse.SetExit(East, ExitProps{UExit: true, RExit: &eastOfHouse})
	northOfHouse.SetExit(North, ExitProps{UExit: true, RExit: &path})
	northOfHouse.SetExit(South, ExitProps{NExit: "The windows are all boarded."})

	// South of House
	southOfHouse.SetExit(West, ExitProps{UExit: true, RExit: &westOfHouse})
	southOfHouse.SetExit(East, ExitProps{UExit: true, RExit: &eastOfHouse})
	southOfHouse.SetExit(NorthEast, ExitProps{UExit: true, RExit: &eastOfHouse})
	southOfHouse.SetExit(NorthWest, ExitProps{UExit: true, RExit: &westOfHouse})
	southOfHouse.SetExit(South, ExitProps{UExit: true, RExit: &forest3})
	southOfHouse.SetExit(North, ExitProps{NExit: "The windows are all boarded."})

	// East of House
	eastOfHouse.SetExit(North, ExitProps{UExit: true, RExit: &northOfHouse})
	eastOfHouse.SetExit(South, ExitProps{UExit: true, RExit: &southOfHouse})
	eastOfHouse.SetExit(SouthWest, ExitProps{UExit: true, RExit: &southOfHouse})
	eastOfHouse.SetExit(NorthWest, ExitProps{UExit: true, RExit: &northOfHouse})
	eastOfHouse.SetExit(East, ExitProps{UExit: true, RExit: &clearing})
	eastOfHouse.SetExit(West, ExitProps{DExit: &kitchenWindow, RExit: &kitchen})
	eastOfHouse.SetExit(In, ExitProps{DExit: &kitchenWindow, RExit: &kitchen})

	// forest 1
	forest1.SetExit(Up, ExitProps{NExit: "There is no tree here suitable for climbing."})
	forest1.SetExit(North, ExitProps{UExit: true, RExit: &gratingClearing})
	forest1.SetExit(East, ExitProps{UExit: true, RExit: &path})
	forest1.SetExit(South, ExitProps{UExit: true, RExit: &forest3})
	forest1.SetExit(West, ExitProps{NExit: "You would need a machete to go further west."})

	// forest 2
	forest2.SetExit(Up, ExitProps{NExit: "There is no tree here suitable for climbing."})
	forest2.SetExit(North, ExitProps{NExit: "The forest becomes impenetrable to the north."})
	forest2.SetExit(East, ExitProps{UExit: true, RExit: &mountains})
	forest2.SetExit(South, ExitProps{UExit: true, RExit: &clearing})
	forest2.SetExit(West, ExitProps{UExit: true, RExit: &path})

	// mountains
	mountains.SetExit(Up, ExitProps{NExit: "The mountains are impassable."})
	mountains.SetExit(North, ExitProps{UExit: true, RExit: &forest2})
	mountains.SetExit(East, ExitProps{NExit: "The mountains are impassable."})
	mountains.SetExit(South, ExitProps{UExit: true, RExit: &forest2})
	mountains.SetExit(West, ExitProps{UExit: true, RExit: &forest2})

	// forest 3
	forest3.SetExit(Up, ExitProps{NExit: "There is no tree here suitable for climbing."})
	forest3.SetExit(North, ExitProps{UExit: true, RExit: &clearing})
	forest3.SetExit(East, ExitProps{NExit: "The rank undergrowth prevents eastward movement."})
	forest3.SetExit(South, ExitProps{NExit: "Storm-tossed trees block your way."})
	forest3.SetExit(West, ExitProps{UExit: true, RExit: &forest1})
	forest3.SetExit(NorthWest, ExitProps{UExit: true, RExit: &southOfHouse})

	// path
	path.SetExit(Up, ExitProps{UExit: true, RExit: &upATree})
	path.SetExit(North, ExitProps{UExit: true, RExit: &gratingClearing})
	path.SetExit(East, ExitProps{UExit: true, RExit: &forest2})
	path.SetExit(South, ExitProps{UExit: true, RExit: &northOfHouse})
	path.SetExit(West, ExitProps{UExit: true, RExit: &forest1})

	// Up a tree
	upATree.SetExit(Down, ExitProps{UExit: true, RExit: &path})
	upATree.SetExit(Up, ExitProps{NExit: "You cannot climb any higher."})

	// Grating clearing
	gratingClearing.SetExit(North, ExitProps{NExit: "The forest becomes impenetrable to the north."})
	gratingClearing.SetExit(East, ExitProps{UExit: true, RExit: &forest2})
	gratingClearing.SetExit(West, ExitProps{UExit: true, RExit: &forest1})
	gratingClearing.SetExit(South, ExitProps{UExit: true, RExit: &path})
	gratingClearing.SetExit(Down, ExitProps{FExit: gratingExitFcn})

	// clearing
	clearing.SetExit(Up, ExitProps{NExit: "There is no tree here suitable for climbing."})
	clearing.SetExit(East, ExitProps{UExit: true, RExit: &canyonView})
	clearing.SetExit(North, ExitProps{UExit: true, RExit: &forest2})
	clearing.SetExit(South, ExitProps{UExit: true, RExit: &forest3})
	clearing.SetExit(West, ExitProps{UExit: true, RExit: &eastOfHouse})

	// kitchen
	kitchen.SetExit(East, ExitProps{DExit: &kitchenWindow, RExit: &eastOfHouse})
	kitchen.SetExit(West, ExitProps{UExit: true, RExit: &livingRoom})
	kitchen.SetExit(Out, ExitProps{DExit: &kitchenWindow, RExit: &eastOfHouse})
	kitchen.SetExit(Up, ExitProps{UExit: true, RExit: &attic})
	kitchen.SetExit(Down, ExitProps{CExit: func() bool { return false }, RExit: &studio, CExitStr: "Only Santa Claus climbs down chimneys."})

	// attic
	attic.SetExit(Down, ExitProps{UExit: true, RExit: &kitchen})

	// Living Room
	livingRoom.SetExit(East, ExitProps{UExit: true, RExit: &kitchen})
	livingRoom.SetExit(West, ExitProps{CExit: func() bool { return gD().MagicFlag }, RExit: &strangePassage, CExitStr: "The door is nailed shut."})
	livingRoom.SetExit(Down, ExitProps{FExit: trapDoorExitFcn})

	// cellar
	cellar.SetExit(North, ExitProps{UExit: true, RExit: &trollRoom})
	cellar.SetExit(South, ExitProps{UExit: true, RExit: &eastOfChasm})
	cellar.SetExit(Up, ExitProps{DExit: &trapDoor, RExit: &livingRoom})
	cellar.SetExit(West, ExitProps{NExit: "You try to ascend the ramp, but it is impossible, and you slide back down."})

	// troll Room
	trollRoom.SetExit(South, ExitProps{UExit: true, RExit: &cellar})
	trollRoom.SetExit(East, ExitProps{CExit: func() bool { return gD().TrollFlag }, RExit: &eWPassage, CExitStr: "The troll fends you off with a menacing gesture."})
	trollRoom.SetExit(West, ExitProps{CExit: func() bool { return gD().TrollFlag }, RExit: &maze1, CExitStr: "The troll fends you off with a menacing gesture."})

	// East of Chasm
	eastOfChasm.SetExit(North, ExitProps{UExit: true, RExit: &cellar})
	eastOfChasm.SetExit(East, ExitProps{UExit: true, RExit: &gallery})
	eastOfChasm.SetExit(Down, ExitProps{NExit: "The chasm probably leads straight to the infernal regions."})

	// gallery
	gallery.SetExit(West, ExitProps{UExit: true, RExit: &eastOfChasm})
	gallery.SetExit(North, ExitProps{UExit: true, RExit: &studio})

	// studio
	studio.SetExit(South, ExitProps{UExit: true, RExit: &gallery})
	studio.SetExit(Up, ExitProps{FExit: upChimneyFcn})

	// Maze 1
	maze1.SetExit(East, ExitProps{UExit: true, RExit: &trollRoom})
	maze1.SetExit(North, ExitProps{UExit: true, RExit: &maze1})
	maze1.SetExit(South, ExitProps{UExit: true, RExit: &maze2})
	maze1.SetExit(West, ExitProps{UExit: true, RExit: &maze4})

	// Maze 2
	maze2.SetExit(South, ExitProps{UExit: true, RExit: &maze1})
	maze2.SetExit(Down, ExitProps{FExit: mazeDiodesFcn})
	maze2.SetExit(East, ExitProps{UExit: true, RExit: &maze3})

	// Maze 3
	maze3.SetExit(West, ExitProps{UExit: true, RExit: &maze2})
	maze3.SetExit(North, ExitProps{UExit: true, RExit: &maze4})
	maze3.SetExit(Up, ExitProps{UExit: true, RExit: &maze5})

	// Maze 4
	maze4.SetExit(West, ExitProps{UExit: true, RExit: &maze3})
	maze4.SetExit(North, ExitProps{UExit: true, RExit: &maze1})
	maze4.SetExit(East, ExitProps{UExit: true, RExit: &deadEnd1})

	// Dead End 1
	deadEnd1.SetExit(South, ExitProps{UExit: true, RExit: &maze4})

	// Maze 5
	maze5.SetExit(East, ExitProps{UExit: true, RExit: &deadEnd2})
	maze5.SetExit(North, ExitProps{UExit: true, RExit: &maze3})
	maze5.SetExit(SouthWest, ExitProps{UExit: true, RExit: &maze6})

	// Dead End 2
	deadEnd2.SetExit(West, ExitProps{UExit: true, RExit: &maze5})

	// Maze 6
	maze6.SetExit(Down, ExitProps{UExit: true, RExit: &maze5})
	maze6.SetExit(East, ExitProps{UExit: true, RExit: &maze7})
	maze6.SetExit(West, ExitProps{UExit: true, RExit: &maze6})
	maze6.SetExit(Up, ExitProps{UExit: true, RExit: &maze9})

	// Maze 7
	maze7.SetExit(Up, ExitProps{UExit: true, RExit: &maze14})
	maze7.SetExit(West, ExitProps{UExit: true, RExit: &maze6})
	maze7.SetExit(Down, ExitProps{FExit: mazeDiodesFcn})
	maze7.SetExit(East, ExitProps{UExit: true, RExit: &maze8})
	maze7.SetExit(South, ExitProps{UExit: true, RExit: &maze15})

	// Maze 8
	maze8.SetExit(NorthEast, ExitProps{UExit: true, RExit: &maze7})
	maze8.SetExit(West, ExitProps{UExit: true, RExit: &maze8})
	maze8.SetExit(SouthEast, ExitProps{UExit: true, RExit: &deadEnd3})

	// Dead End 3
	deadEnd3.SetExit(North, ExitProps{UExit: true, RExit: &maze8})

	// Maze 9
	maze9.SetExit(North, ExitProps{UExit: true, RExit: &maze6})
	maze9.SetExit(Down, ExitProps{FExit: mazeDiodesFcn})
	maze9.SetExit(East, ExitProps{UExit: true, RExit: &maze10})
	maze9.SetExit(South, ExitProps{UExit: true, RExit: &maze13})
	maze9.SetExit(West, ExitProps{UExit: true, RExit: &maze12})
	maze9.SetExit(NorthWest, ExitProps{UExit: true, RExit: &maze9})

	// Maze 10
	maze10.SetExit(East, ExitProps{UExit: true, RExit: &maze9})
	maze10.SetExit(West, ExitProps{UExit: true, RExit: &maze13})
	maze10.SetExit(Up, ExitProps{UExit: true, RExit: &maze11})

	// Maze 11
	maze11.SetExit(NorthEast, ExitProps{UExit: true, RExit: &gratingRoom})
	maze11.SetExit(Down, ExitProps{UExit: true, RExit: &maze10})
	maze11.SetExit(NorthWest, ExitProps{UExit: true, RExit: &maze13})
	maze11.SetExit(SouthWest, ExitProps{UExit: true, RExit: &maze12})

	// Grating Room
	gratingRoom.SetExit(SouthWest, ExitProps{UExit: true, RExit: &maze11})
	gratingRoom.SetExit(Up, ExitProps{DExit: &grate, RExit: &gratingClearing, DExitStr: "The grating is closed."})

	// Maze 12
	maze12.SetExit(Down, ExitProps{FExit: mazeDiodesFcn})
	maze12.SetExit(SouthWest, ExitProps{UExit: true, RExit: &maze11})
	maze12.SetExit(East, ExitProps{UExit: true, RExit: &maze13})
	maze12.SetExit(Up, ExitProps{UExit: true, RExit: &maze9})
	maze12.SetExit(North, ExitProps{UExit: true, RExit: &deadEnd4})

	// Dead End 4
	deadEnd4.SetExit(South, ExitProps{UExit: true, RExit: &maze12})

	// Maze 13
	maze13.SetExit(East, ExitProps{UExit: true, RExit: &maze9})
	maze13.SetExit(Down, ExitProps{UExit: true, RExit: &maze12})
	maze13.SetExit(South, ExitProps{UExit: true, RExit: &maze10})
	maze13.SetExit(West, ExitProps{UExit: true, RExit: &maze11})

	// Maze 14
	maze14.SetExit(West, ExitProps{UExit: true, RExit: &maze15})
	maze14.SetExit(NorthWest, ExitProps{UExit: true, RExit: &maze14})
	maze14.SetExit(NorthEast, ExitProps{UExit: true, RExit: &maze7})
	maze14.SetExit(South, ExitProps{UExit: true, RExit: &maze7})

	// Maze 15
	maze15.SetExit(West, ExitProps{UExit: true, RExit: &maze14})
	maze15.SetExit(South, ExitProps{UExit: true, RExit: &maze7})
	maze15.SetExit(SouthEast, ExitProps{UExit: true, RExit: &cyclopsRoom})

	// cyclops Room
	cyclopsRoom.SetExit(NorthWest, ExitProps{UExit: true, RExit: &maze15})
	cyclopsRoom.SetExit(East, ExitProps{CExit: func() bool { return gD().MagicFlag }, RExit: &strangePassage, CExitStr: "The east wall is solid rock."})
	cyclopsRoom.SetExit(Up, ExitProps{CExit: func() bool { return gD().CyclopsFlag }, RExit: &treasureRoom, CExitStr: "The cyclops doesn't look like he'll let you past."})

	// Strange Passage
	strangePassage.SetExit(West, ExitProps{UExit: true, RExit: &cyclopsRoom})
	strangePassage.SetExit(In, ExitProps{UExit: true, RExit: &cyclopsRoom})
	strangePassage.SetExit(East, ExitProps{UExit: true, RExit: &livingRoom})

	// Treasure Room
	treasureRoom.SetExit(Down, ExitProps{UExit: true, RExit: &cyclopsRoom})

	// reservoir South
	reservoirSouth.SetExit(SouthEast, ExitProps{UExit: true, RExit: &deepCanyon})
	reservoirSouth.SetExit(SouthWest, ExitProps{UExit: true, RExit: &chasmRoom})
	reservoirSouth.SetExit(East, ExitProps{UExit: true, RExit: &damRoom})
	reservoirSouth.SetExit(West, ExitProps{UExit: true, RExit: &streamView})
	reservoirSouth.SetExit(North, ExitProps{CExit: func() bool { return gD().LowTide }, RExit: &reservoir, CExitStr: "You would drown."})

	// reservoir
	reservoir.SetExit(North, ExitProps{UExit: true, RExit: &reservoirNorth})
	reservoir.SetExit(South, ExitProps{UExit: true, RExit: &reservoirSouth})
	reservoir.SetExit(Up, ExitProps{UExit: true, RExit: &inStream})
	reservoir.SetExit(West, ExitProps{UExit: true, RExit: &inStream})
	reservoir.SetExit(Down, ExitProps{NExit: "The dam blocks your way."})

	// reservoir North
	reservoirNorth.SetExit(North, ExitProps{UExit: true, RExit: &atlantisRoom})
	reservoirNorth.SetExit(South, ExitProps{CExit: func() bool { return gD().LowTide }, RExit: &reservoir, CExitStr: "You would drown."})

	// Stream View
	streamView.SetExit(East, ExitProps{UExit: true, RExit: &reservoirSouth})
	streamView.SetExit(West, ExitProps{NExit: "The stream emerges from a spot too small for you to enter."})

	// In Stream
	inStream.SetExit(Up, ExitProps{NExit: "The channel is too narrow."})
	inStream.SetExit(West, ExitProps{NExit: "The channel is too narrow."})
	inStream.SetExit(Land, ExitProps{UExit: true, RExit: &streamView})
	inStream.SetExit(Down, ExitProps{UExit: true, RExit: &reservoir})
	inStream.SetExit(East, ExitProps{UExit: true, RExit: &reservoir})

	// Mirror Room 1
	mirrorRoom1.SetExit(North, ExitProps{UExit: true, RExit: &coldPassage})
	mirrorRoom1.SetExit(West, ExitProps{UExit: true, RExit: &twistingPassage})
	mirrorRoom1.SetExit(East, ExitProps{UExit: true, RExit: &smallCave})

	// Mirror Room 2
	mirrorRoom2.SetExit(West, ExitProps{UExit: true, RExit: &windingPassage})
	mirrorRoom2.SetExit(North, ExitProps{UExit: true, RExit: &narrowPassage})
	mirrorRoom2.SetExit(East, ExitProps{UExit: true, RExit: &tinyCave})

	// Small Cave
	smallCave.SetExit(North, ExitProps{UExit: true, RExit: &mirrorRoom1})
	smallCave.SetExit(Down, ExitProps{UExit: true, RExit: &atlantisRoom})
	smallCave.SetExit(South, ExitProps{UExit: true, RExit: &atlantisRoom})
	smallCave.SetExit(West, ExitProps{UExit: true, RExit: &twistingPassage})

	// Tiny Cave
	tinyCave.SetExit(North, ExitProps{UExit: true, RExit: &mirrorRoom2})
	tinyCave.SetExit(West, ExitProps{UExit: true, RExit: &windingPassage})
	tinyCave.SetExit(Down, ExitProps{UExit: true, RExit: &enteranceToHades})

	// Cold Passage
	coldPassage.SetExit(South, ExitProps{UExit: true, RExit: &mirrorRoom1})
	coldPassage.SetExit(West, ExitProps{UExit: true, RExit: &slideRoom})

	// Narrow Passage
	narrowPassage.SetExit(North, ExitProps{UExit: true, RExit: &roundRoom})
	narrowPassage.SetExit(South, ExitProps{UExit: true, RExit: &mirrorRoom2})

	// Winding Passage
	windingPassage.SetExit(North, ExitProps{UExit: true, RExit: &mirrorRoom2})
	windingPassage.SetExit(East, ExitProps{UExit: true, RExit: &tinyCave})

	// Twisting Passage
	twistingPassage.SetExit(North, ExitProps{UExit: true, RExit: &mirrorRoom1})
	twistingPassage.SetExit(East, ExitProps{UExit: true, RExit: &smallCave})

	// Atlantis Room
	atlantisRoom.SetExit(Up, ExitProps{UExit: true, RExit: &smallCave})
	atlantisRoom.SetExit(South, ExitProps{UExit: true, RExit: &reservoirNorth})

	// EW Passage
	eWPassage.SetExit(East, ExitProps{UExit: true, RExit: &roundRoom})
	eWPassage.SetExit(West, ExitProps{UExit: true, RExit: &trollRoom})
	eWPassage.SetExit(Down, ExitProps{UExit: true, RExit: &chasmRoom})
	eWPassage.SetExit(North, ExitProps{UExit: true, RExit: &chasmRoom})

	// Round Room
	roundRoom.SetExit(East, ExitProps{UExit: true, RExit: &loudRoom})
	roundRoom.SetExit(West, ExitProps{UExit: true, RExit: &eWPassage})
	roundRoom.SetExit(North, ExitProps{UExit: true, RExit: &nSPassage})
	roundRoom.SetExit(South, ExitProps{UExit: true, RExit: &narrowPassage})
	roundRoom.SetExit(SouthEast, ExitProps{UExit: true, RExit: &engravingsCave})

	// Deep Canyon
	deepCanyon.SetExit(NorthWest, ExitProps{UExit: true, RExit: &reservoirSouth})
	deepCanyon.SetExit(East, ExitProps{UExit: true, RExit: &damRoom})
	deepCanyon.SetExit(SouthWest, ExitProps{UExit: true, RExit: &nSPassage})
	deepCanyon.SetExit(Down, ExitProps{UExit: true, RExit: &loudRoom})

	// Damp Cave
	dampCave.SetExit(West, ExitProps{UExit: true, RExit: &loudRoom})
	dampCave.SetExit(East, ExitProps{UExit: true, RExit: &whiteCliffsNorth})
	dampCave.SetExit(South, ExitProps{NExit: "it is too narrow for most insects."})

	// Loud Room
	loudRoom.SetExit(East, ExitProps{UExit: true, RExit: &dampCave})
	loudRoom.SetExit(West, ExitProps{UExit: true, RExit: &roundRoom})
	loudRoom.SetExit(Up, ExitProps{UExit: true, RExit: &deepCanyon})

	// NS Passage
	nSPassage.SetExit(North, ExitProps{UExit: true, RExit: &chasmRoom})
	nSPassage.SetExit(NorthEast, ExitProps{UExit: true, RExit: &deepCanyon})
	nSPassage.SetExit(South, ExitProps{UExit: true, RExit: &roundRoom})

	// Chasm Room
	chasmRoom.SetExit(NorthEast, ExitProps{UExit: true, RExit: &reservoirSouth})
	chasmRoom.SetExit(SouthWest, ExitProps{UExit: true, RExit: &eWPassage})
	chasmRoom.SetExit(Up, ExitProps{UExit: true, RExit: &eWPassage})
	chasmRoom.SetExit(South, ExitProps{UExit: true, RExit: &nSPassage})
	chasmRoom.SetExit(Down, ExitProps{NExit: "Are you out of your mind?"})

	// Entrance to Hades
	enteranceToHades.SetExit(Up, ExitProps{UExit: true, RExit: &tinyCave})
	enteranceToHades.SetExit(In, ExitProps{CExit: func() bool { return gD().LLDFlag }, RExit: &landOfLivingDead, CExitStr: "Some invisible force prevents you from passing through the gate."})
	enteranceToHades.SetExit(South, ExitProps{CExit: func() bool { return gD().LLDFlag }, RExit: &landOfLivingDead, CExitStr: "Some invisible force prevents you from passing through the gate."})

	// Land of Living Dead
	landOfLivingDead.SetExit(Out, ExitProps{UExit: true, RExit: &enteranceToHades})
	landOfLivingDead.SetExit(North, ExitProps{UExit: true, RExit: &enteranceToHades})

	// engravings Cave
	engravingsCave.SetExit(NorthWest, ExitProps{UExit: true, RExit: &roundRoom})
	engravingsCave.SetExit(East, ExitProps{UExit: true, RExit: &domeRoom})

	// Egypt Room
	egyptRoom.SetExit(West, ExitProps{UExit: true, RExit: &northTemple})
	egyptRoom.SetExit(Up, ExitProps{UExit: true, RExit: &northTemple})

	// Dome Room
	domeRoom.SetExit(West, ExitProps{UExit: true, RExit: &engravingsCave})
	domeRoom.SetExit(Down, ExitProps{CExit: func() bool { return gD().DomeFlag }, RExit: &torchRoom, CExitStr: "You cannot go down without fracturing many bones."})

	// torch Room
	torchRoom.SetExit(Up, ExitProps{NExit: "You cannot reach the rope."})
	torchRoom.SetExit(South, ExitProps{UExit: true, RExit: &northTemple})
	torchRoom.SetExit(Down, ExitProps{UExit: true, RExit: &northTemple})

	// North Temple
	northTemple.SetExit(Down, ExitProps{UExit: true, RExit: &egyptRoom})
	northTemple.SetExit(East, ExitProps{UExit: true, RExit: &egyptRoom})
	northTemple.SetExit(North, ExitProps{UExit: true, RExit: &torchRoom})
	northTemple.SetExit(Out, ExitProps{UExit: true, RExit: &torchRoom})
	northTemple.SetExit(Up, ExitProps{UExit: true, RExit: &torchRoom})
	northTemple.SetExit(South, ExitProps{UExit: true, RExit: &southTemple})

	// South Temple
	southTemple.SetExit(North, ExitProps{UExit: true, RExit: &northTemple})
	southTemple.SetExit(Down, ExitProps{CExit: func() bool { return gD().CoffinCure }, RExit: &tinyCave, CExitStr: "You haven't a prayer of getting the coffin down there."})

	// dam Room
	damRoom.SetExit(South, ExitProps{UExit: true, RExit: &deepCanyon})
	damRoom.SetExit(Down, ExitProps{UExit: true, RExit: &damBase})
	damRoom.SetExit(East, ExitProps{UExit: true, RExit: &damBase})
	damRoom.SetExit(North, ExitProps{UExit: true, RExit: &damLobby})
	damRoom.SetExit(West, ExitProps{UExit: true, RExit: &reservoirSouth})

	// dam Lobby
	damLobby.SetExit(South, ExitProps{UExit: true, RExit: &damRoom})
	damLobby.SetExit(North, ExitProps{UExit: true, RExit: &maintenanceRoom})
	damLobby.SetExit(East, ExitProps{UExit: true, RExit: &maintenanceRoom})

	// Maintenance Room
	maintenanceRoom.SetExit(South, ExitProps{UExit: true, RExit: &damLobby})
	maintenanceRoom.SetExit(West, ExitProps{UExit: true, RExit: &damLobby})

	// dam Base
	damBase.SetExit(North, ExitProps{UExit: true, RExit: &damRoom})
	damBase.SetExit(Up, ExitProps{UExit: true, RExit: &damRoom})

	// river 1
	river1.SetExit(Up, ExitProps{NExit: "You cannot go upstream due to strong currents."})
	river1.SetExit(West, ExitProps{UExit: true, RExit: &damBase})
	river1.SetExit(Land, ExitProps{UExit: true, RExit: &damBase})
	river1.SetExit(Down, ExitProps{UExit: true, RExit: &river2})
	river1.SetExit(East, ExitProps{NExit: "The White Cliffs prevent your landing here."})

	// river 2
	river2.SetExit(Up, ExitProps{NExit: "You cannot go upstream due to strong currents."})
	river2.SetExit(Down, ExitProps{UExit: true, RExit: &river3})
	river2.SetExit(Land, ExitProps{NExit: "There is no safe landing spot here."})
	river2.SetExit(East, ExitProps{NExit: "The White Cliffs prevent your landing here."})
	river2.SetExit(West, ExitProps{NExit: "Just in time you steer away from the rocks."})

	// river 3
	river3.SetExit(Up, ExitProps{NExit: "You cannot go upstream due to strong currents."})
	river3.SetExit(Down, ExitProps{UExit: true, RExit: &river4})
	river3.SetExit(Land, ExitProps{UExit: true, RExit: &whiteCliffsNorth})
	river3.SetExit(West, ExitProps{UExit: true, RExit: &whiteCliffsNorth})

	// White Cliffs North
	whiteCliffsNorth.SetExit(South, ExitProps{CExit: func() bool { return gD().DeflateFlag }, RExit: &whiteCliffsSouth, CExitStr: "The path is too narrow."})
	whiteCliffsNorth.SetExit(West, ExitProps{CExit: func() bool { return gD().DeflateFlag }, RExit: &dampCave, CExitStr: "The path is too narrow."})

	// White Cliffs South
	whiteCliffsSouth.SetExit(North, ExitProps{CExit: func() bool { return gD().DeflateFlag }, RExit: &whiteCliffsNorth, CExitStr: "The path is too narrow."})

	// river 4
	river4.SetExit(Up, ExitProps{NExit: "You cannot go upstream due to strong currents."})
	river4.SetExit(Down, ExitProps{UExit: true, RExit: &river5})
	river4.SetExit(Land, ExitProps{NExit: "You can land either to the east or the west."})
	river4.SetExit(West, ExitProps{UExit: true, RExit: &whiteCliffsSouth})
	river4.SetExit(East, ExitProps{UExit: true, RExit: &sandyBeach})

	// river 5
	river5.SetExit(Up, ExitProps{NExit: "You cannot go upstream due to strong currents."})
	river5.SetExit(East, ExitProps{UExit: true, RExit: &shore})
	river5.SetExit(Land, ExitProps{UExit: true, RExit: &shore})

	// shore
	shore.SetExit(North, ExitProps{UExit: true, RExit: &sandyBeach})
	shore.SetExit(South, ExitProps{UExit: true, RExit: &aragainFalls})

	// Sandy Beach
	sandyBeach.SetExit(NorthEast, ExitProps{UExit: true, RExit: &sandyCave})
	sandyBeach.SetExit(South, ExitProps{UExit: true, RExit: &shore})

	// Sandy Cave
	sandyCave.SetExit(SouthWest, ExitProps{UExit: true, RExit: &sandyBeach})

	// Aragain Falls
	aragainFalls.SetExit(West, ExitProps{CExit: func() bool { return gD().RainbowFlag }, RExit: &onRainbow})
	aragainFalls.SetExit(Down, ExitProps{NExit: "it's a long way..."})
	aragainFalls.SetExit(North, ExitProps{UExit: true, RExit: &shore})
	aragainFalls.SetExit(Up, ExitProps{CExit: func() bool { return gD().RainbowFlag }, RExit: &onRainbow})

	// On rainbow
	onRainbow.SetExit(West, ExitProps{UExit: true, RExit: &endOfRainbow})
	onRainbow.SetExit(East, ExitProps{UExit: true, RExit: &aragainFalls})

	// End of rainbow
	endOfRainbow.SetExit(Up, ExitProps{CExit: func() bool { return gD().RainbowFlag }, RExit: &onRainbow})
	endOfRainbow.SetExit(NorthEast, ExitProps{CExit: func() bool { return gD().RainbowFlag }, RExit: &onRainbow})
	endOfRainbow.SetExit(East, ExitProps{CExit: func() bool { return gD().RainbowFlag }, RExit: &onRainbow})
	endOfRainbow.SetExit(SouthWest, ExitProps{UExit: true, RExit: &canyonBottom})

	// Canyon Bottom
	canyonBottom.SetExit(Up, ExitProps{UExit: true, RExit: &cliffMiddle})
	canyonBottom.SetExit(North, ExitProps{UExit: true, RExit: &endOfRainbow})

	// Cliff Middle
	cliffMiddle.SetExit(Up, ExitProps{UExit: true, RExit: &canyonView})
	cliffMiddle.SetExit(Down, ExitProps{UExit: true, RExit: &canyonBottom})

	// Canyon View
	canyonView.SetExit(East, ExitProps{UExit: true, RExit: &cliffMiddle})
	canyonView.SetExit(Down, ExitProps{UExit: true, RExit: &cliffMiddle})
	canyonView.SetExit(NorthWest, ExitProps{UExit: true, RExit: &clearing})
	canyonView.SetExit(West, ExitProps{UExit: true, RExit: &forest3})
	canyonView.SetExit(South, ExitProps{NExit: "Storm-tossed trees block your way."})

	// Mine Entrance
	mineEntrance.SetExit(South, ExitProps{UExit: true, RExit: &slideRoom})
	mineEntrance.SetExit(In, ExitProps{UExit: true, RExit: &squeekyRoom})
	mineEntrance.SetExit(West, ExitProps{UExit: true, RExit: &squeekyRoom})

	// Squeaky Room
	squeekyRoom.SetExit(North, ExitProps{UExit: true, RExit: &batRoom})
	squeekyRoom.SetExit(East, ExitProps{UExit: true, RExit: &mineEntrance})

	// bat Room
	batRoom.SetExit(South, ExitProps{UExit: true, RExit: &squeekyRoom})
	batRoom.SetExit(East, ExitProps{UExit: true, RExit: &shaftRoom})

	// Shaft Room
	shaftRoom.SetExit(Down, ExitProps{NExit: "You wouldn't fit and would die if you could."})
	shaftRoom.SetExit(West, ExitProps{UExit: true, RExit: &batRoom})
	shaftRoom.SetExit(North, ExitProps{UExit: true, RExit: &smellyRoom})

	// Smelly Room
	smellyRoom.SetExit(Down, ExitProps{UExit: true, RExit: &gasRoom})
	smellyRoom.SetExit(South, ExitProps{UExit: true, RExit: &shaftRoom})

	// Gas Room
	gasRoom.SetExit(Up, ExitProps{UExit: true, RExit: &smellyRoom})
	gasRoom.SetExit(East, ExitProps{UExit: true, RExit: &mine1})

	// ladder Top
	ladderTop.SetExit(Down, ExitProps{UExit: true, RExit: &ladderBottom})
	ladderTop.SetExit(Up, ExitProps{UExit: true, RExit: &mine4})

	// ladder Bottom
	ladderBottom.SetExit(South, ExitProps{UExit: true, RExit: &deadEnd5})
	ladderBottom.SetExit(West, ExitProps{UExit: true, RExit: &timberRoom})
	ladderBottom.SetExit(Up, ExitProps{UExit: true, RExit: &ladderTop})

	// Dead End 5
	deadEnd5.SetExit(North, ExitProps{UExit: true, RExit: &ladderBottom})

	// Timber Room
	timberRoom.SetExit(East, ExitProps{UExit: true, RExit: &ladderBottom})
	timberRoom.SetExit(West, ExitProps{CExit: func() bool { return gD().EmptyHanded }, RExit: &lowerShaft, CExitStr: "You cannot fit through this passage with that load."})

	// Lower Shaft
	lowerShaft.SetExit(South, ExitProps{UExit: true, RExit: &machineRoom})
	lowerShaft.SetExit(Out, ExitProps{CExit: func() bool { return gD().EmptyHanded }, RExit: &timberRoom, CExitStr: "You cannot fit through this passage with that load."})
	lowerShaft.SetExit(East, ExitProps{CExit: func() bool { return gD().EmptyHanded }, RExit: &timberRoom, CExitStr: "You cannot fit through this passage with that load."})

	// machine Room
	machineRoom.SetExit(North, ExitProps{UExit: true, RExit: &lowerShaft})

	// Mine 1
	mine1.SetExit(North, ExitProps{UExit: true, RExit: &gasRoom})
	mine1.SetExit(East, ExitProps{UExit: true, RExit: &mine1})
	mine1.SetExit(NorthEast, ExitProps{UExit: true, RExit: &mine2})

	// Mine 2
	mine2.SetExit(North, ExitProps{UExit: true, RExit: &mine2})
	mine2.SetExit(South, ExitProps{UExit: true, RExit: &mine1})
	mine2.SetExit(SouthEast, ExitProps{UExit: true, RExit: &mine3})

	// Mine 3
	mine3.SetExit(South, ExitProps{UExit: true, RExit: &mine3})
	mine3.SetExit(SouthWest, ExitProps{UExit: true, RExit: &mine4})
	mine3.SetExit(East, ExitProps{UExit: true, RExit: &mine2})

	// Mine 4
	mine4.SetExit(North, ExitProps{UExit: true, RExit: &mine3})
	mine4.SetExit(West, ExitProps{UExit: true, RExit: &mine4})
	mine4.SetExit(Down, ExitProps{UExit: true, RExit: &ladderTop})

	// slide Room
	slideRoom.SetExit(East, ExitProps{UExit: true, RExit: &coldPassage})
	slideRoom.SetExit(North, ExitProps{UExit: true, RExit: &mineEntrance})
	slideRoom.SetExit(Down, ExitProps{UExit: true, RExit: &cellar})
}

func finalizeGameObjects() {
	// Set up room exits
	initRoomExits()

	// Add trapDoor to cellar's globals (done here to avoid potential init ordering issues)
	cellar.Global = append(cellar.Global, &trapDoor)

	// Set action functions that would cause init cycles
	whiteHouse.Action = whiteHouseFcn
	mailbox.Action = mailboxFcn
	forest.Action = forestFcn
	kitchenWindow.Action = kitchenWindowFcn
	chimney.Action = chimneyFcn
	grate.Action = grateFcn
	climbableCliff.Action = cliffObjectFcn
	gunk.Action = gunkFcn
	gratingClearing.Action = clearingFcn
	clearing.Action = forestRoomFcn
	kitchen.Action = kitchenFcn
	livingRoom.Action = livingRoomFcn
	stoneBarrow.Action = stoneBarrowFcn
	rainbow.Action = rainbowFcn
	inflatedBoat.Action = rBoatFcn
	hotBell.Action = hotBellFcn
	trollRoom.Action = trollRoomFcn
	lowerShaft.Action = noObjsFcn
	machineRoom.Action = machineRoomFcn
	trophyCase.Action = trophyCaseFcn
	sword.Action = swordFcn
	lamp.Action = lanternFcn
	bell.Action = bellFcn
	match.Action = matchFcn
	toolChest.Action = toolChestFcn
	machine.Action = machineFcn
	machineSwitch.Action = machineSwitchFcn
	sceptre.Action = sceptreFcn
	board.Action = boardFcn
	teeth.Action = teethFcn
	graniteWall.Action = graniteWallFcn
	songbird.Action = songbirdFcn
	globalWater.Action = waterFcn
	slide.Action = slideFcn
	bodies.Action = bodyFcn
	crack.Action = crackFcn
	whiteCliff.Action = whiteCliffFcn
	river.Action = riverFcn
	boardedWindow.Action = boardedWindowFcn
	puncturedBoat.Action = puncturedBoatFcn
	westOfHouse.Action = westHouseFcn
	eastOfHouse.Action = eastHouseFcn
	forest1.Action = forestRoomFcn
	forest2.Action = forestRoomFcn
	forest3.Action = forestRoomFcn
	path.Action = forestRoomFcn
	upATree.Action = treeRoomFcn
	cellar.Action = cellarFcn
	gratingRoom.Action = maze11Fcn
	cyclopsRoom.Action = cyclopsRoomFcn
	treasureRoom.Action = treasureRoomFcn
	reservoirSouth.Action = reservoirSouthFcn
	reservoir.Action = reservoirFcn
	reservoirNorth.Action = reservoirNorthFcn
	mirrorRoom1.Action = mirrorRoomFcn
	mirrorRoom2.Action = mirrorRoomFcn
	tinyCave.Action = cave2RoomFcn
	deepCanyon.Action = deepCanyonFcn
	loudRoom.Action = loudRoomFcn
	enteranceToHades.Action = lLDRoomFcn
	domeRoom.Action = domeRoomFcn
	torchRoom.Action = torchRoomFcn
	southTemple.Action = southTempleFcn
	damRoom.Action = damRoomFcn
	whiteCliffsNorth.Action = whiteCliffsFcn
	whiteCliffsSouth.Action = whiteCliffsFcn
	river4.Action = rivr4RoomFcn
	aragainFalls.Action = fallsRoomFcn
	canyonView.Action = canyonViewFcn
	batRoom.Action = batsRoomFcn
	gasRoom.Action = boomRoomFcn
	timberRoom.Action = noObjsFcn
	mountainRange.Action = mountainRangeFcn
	frontDoor.Action = frontDoorFcn
	barrowDoor.Action = barrowDoorFcn
	barrow.Action = barrowFcn
	rug.Action = rugFcn
	trapDoor.Action = trapDoorFcn
	woodenDoor.Action = frontDoorFcn
	rope.Action = ropeFcn
	ghosts.Action = ghostsFcn
	raisedBasket.Action = basketFcn
	loweredBasket.Action = basketFcn
	bat.Action = batFcn
	candles.Action = candlesFcn
	troll.Action = trollFcn
	bolt.Action = boltFcn
	bubble.Action = bubbleFcn
	dam.Action = damFunction
	inflatableBoat.Action = inflatableBoatFcn
	yellowButton.Action = buttonFcn
	brownButton.Action = buttonFcn
	redButton.Action = buttonFcn
	blueButton.Action = buttonFcn
	tube.Action = tubeFcn
	leak.Action = leakFcn
	cyclops.Action = cyclopsFcn
	chalice.Action = chaliceFcn
	painting.Action = paintingFcn
	leaves.Action = leafPileFcn
	sand.Action = sandFunction
	thief.Action = robberFcn
	trunk.Action = trunkFcn
	mirror1.Action = mirrorMirrorFcn
	mirror2.Action = mirrorMirrorFcn
	pedestal.Action = dumbContainerFcn
	buoy.Action = treasureInsideFcn
	bones.Action = skeletonFcn
	bagOfCoins.Action = bagOfCoinsFcn
	rustyKnife.Action = rustyKnifeFcn
	bottle.Action = bottleFcn
	sandwichBag.Action = sandwichBagFcn
	knife.Action = knifeFcn
	water.Action = waterFcn
	garlic.Action = garlicFcn
	book.Action = blackBookFcn
	egg.Action = eggObjectFcn
	canary.Action = canaryObjectFcn
	putty.Action = puttyFcn
	axe.Action = axeFcn
	largeBag.Action = largeBagFcn
	stiletto.Action = stiletteFcn
	torch.Action = torchFcn
	brokenCanary.Action = canaryObjectFcn

	// Set up globals.go object actions
	localGlobals.Pseudo = []PseudoObj{{
		Synonym: "foobar",
		Action:  vWalk,
	}}
	notHereObject.Action = notHereObjectFcn
	me.Action = cretinFcn
	ground.Action = groundFunction

	// Initialize villain table
	gD().Villains = []*VillainEntry{
		{Villain: &troll, Best: &sword, BestAdv: 1, Prob: 0, Msgs: &trollMelee},
		{Villain: &thief, Best: &knife, BestAdv: 1, Prob: 0, Msgs: &thiefMelee},
		{Villain: &cyclops, Best: nil, BestAdv: 0, Prob: 0, Msgs: &cyclopsMelee},
	}
}
