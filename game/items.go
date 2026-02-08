package game

import . "github.com/ajdnik/gozork/engine"

var (
	// ================================================================
	// GLOBAL OBJECTS (In: &localGlobals or &globalObjects)
	// ================================================================

	board = Object{
		In:       &localGlobals,
		Synonyms: []string{"boards", "board"},
		Desc:     "board",
		Flags:    FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	teeth = Object{
		In:       &globalObjects,
		Synonyms: []string{"overboard", "teeth"},
		Desc:     "set of teeth",
		Flags:    FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	wall = Object{
		In:         &globalObjects,
		Synonyms:   []string{"wall", "walls"},
		Adjectives: []string{"surrounding"},
		Desc:       "surrounding wall",
	}
	graniteWall = Object{
		In:         &globalObjects,
		Synonyms:   []string{"wall"},
		Adjectives: []string{"granite"},
		Desc:       "granite wall",
		// Action set in finalizeGameObjects to avoid init cycle
	}
	songbird = Object{
		In:         &localGlobals,
		Synonyms:   []string{"bird", "songbird"},
		Adjectives: []string{"song"},
		Desc:       "songbird",
		Flags:      FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	whiteHouse = Object{
		In:         &localGlobals,
		Synonyms:   []string{"house"},
		Adjectives: []string{"white", "beautiful", "colonial"},
		Desc:       "white house",
		Flags:      FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	forest = Object{
		In:       &localGlobals,
		Synonyms: []string{"forest", "trees", "pines", "hemlocks"},
		Desc:     "forest",
		Flags:    FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	tree = Object{
		In:         &localGlobals,
		Synonyms:   []string{"tree", "branch"},
		Adjectives: []string{"large", "storm"},
		Desc:       "tree",
		Flags:      FlgNoDesc | FlgClimb,
	}
	globalWater = Object{
		In:       &localGlobals,
		Synonyms: []string{"water", "quantity"},
		Desc:     "water",
		Flags:    FlgDrink,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	kitchenWindow = Object{
		In:         &localGlobals,
		Synonyms:   []string{"window"},
		Adjectives: []string{"kitchen", "small"},
		Desc:       "kitchen window",
		Flags:      FlgDoor | FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	chimney = Object{
		In:         &localGlobals,
		Synonyms:   []string{"chimney"},
		Adjectives: []string{"dark", "narrow"},
		Desc:       "chimney",
		Flags:      FlgClimb | FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	slide = Object{
		In:         &localGlobals,
		Synonyms:   []string{"chute", "ramp", "slide"},
		Adjectives: []string{"steep", "metal", "twisting"},
		Desc:       "chute",
		Flags:      FlgClimb,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	bodies = Object{
		In:         &localGlobals,
		Synonyms:   []string{"bodies", "body", "remains", "pile"},
		Adjectives: []string{"mangled"},
		Desc:       "pile of bodies",
		Flags:      FlgNoDesc | FlgTryTake,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	crack = Object{
		In:         &localGlobals,
		Synonyms:   []string{"crack"},
		Adjectives: []string{"narrow"},
		Desc:       "crack",
		Flags:      FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	grate = Object{
		In:       &localGlobals,
		Synonyms: []string{"grate", "grating"},
		Desc:     "grating",
		Flags:    FlgDoor | FlgNoDesc | FlgInvis,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	ladder = Object{
		In:         &localGlobals,
		Synonyms:   []string{"ladder"},
		Adjectives: []string{"wooden", "rickety", "narrow"},
		Desc:       "wooden ladder",
		Flags:      FlgNoDesc | FlgClimb,
	}
	climbableCliff = Object{
		In:         &localGlobals,
		Synonyms:   []string{"wall", "cliff", "walls", "ledge"},
		Adjectives: []string{"rocky", "sheer"},
		Desc:       "cliff",
		Flags:      FlgNoDesc | FlgClimb,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	whiteCliff = Object{
		In:         &localGlobals,
		Synonyms:   []string{"cliff", "cliffs"},
		Adjectives: []string{"white"},
		Desc:       "white cliffs",
		Flags:      FlgNoDesc | FlgClimb,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	rainbow = Object{
		In:       &localGlobals,
		Synonyms: []string{"rainbow"},
		Desc:     "rainbow",
		Flags:    FlgNoDesc | FlgClimb,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	river = Object{
		In:         &localGlobals,
		Synonyms:   []string{"river"},
		Adjectives: []string{"frigid"},
		Desc:       "river",
		Flags:      FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	boardedWindow = Object{
		In:         &localGlobals,
		Synonyms:   []string{"window"},
		Adjectives: []string{"boarded"},
		Desc:       "boarded window",
		Flags:      FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}

	// ================================================================
	// UNPLACED OBJECTS (created or swapped during gameplay)
	// ================================================================

	inflatedBoat = Object{
		Synonyms:   []string{"boat", "raft"},
		Adjectives: []string{"magic", "plastic", "seaworthy", "inflated", "inflatable"},
		Desc:       "magic boat",
		Flags:      FlgTake | FlgBurn | FlgVeh | FlgOpen | FlgSearch,
		// Action set in finalizeGameObjects to avoid init cycle
		Item:    &ItemData{Size: 20, Capacity: 100},
		Vehicle: &VehicleData{Type: FlgNonLand},
	}
	puncturedBoat = Object{
		Synonyms:   []string{"boat", "pile", "plastic"},
		Adjectives: []string{"plastic", "puncture", "large"},
		Desc:       "punctured boat",
		Flags:      FlgTake | FlgBurn,
		// Action set in finalizeGameObjects to avoid init cycle
		Item: &ItemData{Size: 20},
	}
	brokenLamp = Object{
		Synonyms:   []string{"lamp", "lantern"},
		Adjectives: []string{"broken"},
		Desc:       "broken lantern",
		Flags:      FlgTake,
	}
	gunk = Object{
		Synonyms:   []string{"gunk", "piece", "slag"},
		Adjectives: []string{"small", "vitreous"},
		Desc:       "small piece of vitreous slag",
		Flags:      FlgTake | FlgTryTake,
		// Action set in finalizeGameObjects to avoid init cycle
		Item: &ItemData{Size: 10},
	}
	hotBell = Object{
		Synonyms:   []string{"bell"},
		Adjectives: []string{"brass", "hot", "red", "small"},
		Desc:       "red hot brass bell",
		Flags:      FlgTryTake,
		LongDesc:   "On the ground is a red hot bell.",
		// Action set in finalizeGameObjects to avoid init cycle
	}
	brokenEgg = Object{
		Synonyms:   []string{"egg", "treasure"},
		Adjectives: []string{"broken", "birds", "encrusted", "jewel"},
		Desc:       "broken jewel-encrusted egg",
		Flags:      FlgTake | FlgCont | FlgOpen,
		LongDesc:   "There is a somewhat ruined egg here.",
		Item:       &ItemData{TValue: 2, Capacity: 6},
	}
	bauble = Object{
		Synonyms:   []string{"bauble", "treasure"},
		Adjectives: []string{"brass", "beautiful"},
		Desc:       "beautiful brass bauble",
		Flags:      FlgTake,
		Item:       &ItemData{Value: 1, TValue: 1},
	}
	diamond = Object{
		Synonyms:   []string{"diamond", "treasure"},
		Adjectives: []string{"huge", "enormous"},
		Desc:       "huge diamond",
		Flags:      FlgTake,
		LongDesc:   "There is an enormous diamond (perfectly cut) here.",
		Item:       &ItemData{Value: 10, TValue: 10},
	}

	// ================================================================
	// OBJECTS IN ROOMS
	// ================================================================

	// mountains
	mountainRange = Object{
		In:         &mountains,
		Synonyms:   []string{"mountain", "range"},
		Adjectives: []string{"impassable", "flathead"},
		Desc:       "mountain range",
		Flags:      FlgNoDesc | FlgClimb,
		// Action set in finalizeGameObjects to avoid init cycle
	}

	// West of House
	frontDoor = Object{
		In:         &westOfHouse,
		Synonyms:   []string{"door"},
		Adjectives: []string{"front", "boarded"},
		Desc:       "door",
		Flags:      FlgDoor | FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	mailbox = Object{
		In:         &westOfHouse,
		Synonyms:   []string{"mailbox", "box"},
		Adjectives: []string{"small"},
		Desc:       "small mailbox",
		Flags:      FlgCont | FlgTryTake,
		// Action set in finalizeGameObjects to avoid init cycle
		Item: &ItemData{Capacity: 10},
	}

	// Stone barrow
	barrowDoor = Object{
		In:         &stoneBarrow,
		Synonyms:   []string{"door"},
		Adjectives: []string{"huge", "stone"},
		Desc:       "stone door",
		Flags:      FlgDoor | FlgNoDesc | FlgOpen,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	barrow = Object{
		In:         &stoneBarrow,
		Synonyms:   []string{"barrow", "tomb"},
		Adjectives: []string{"massive", "stone"},
		Desc:       "stone barrow",
		Flags:      FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}

	// Living Room
	trophyCase = Object{
		In:         &livingRoom,
		Synonyms:   []string{"case"},
		Adjectives: []string{"trophy"},
		Desc:       "trophy case",
		Flags:      FlgTrans | FlgCont | FlgNoDesc | FlgTryTake | FlgSearch,
		// Action set in finalizeGameObjects to avoid init cycle
		Item: &ItemData{Capacity: 10000},
	}
	rug = Object{
		In:         &livingRoom,
		Synonyms:   []string{"rug", "carpet"},
		Adjectives: []string{"large", "oriental"},
		Desc:       "carpet",
		Flags:      FlgNoDesc | FlgTryTake,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	trapDoor = Object{
		In:         &livingRoom,
		Synonyms:   []string{"door", "trapdoor", "trap-door", "cover"},
		Adjectives: []string{"trap", "dusty"},
		Desc:       "trap door",
		Flags:      FlgDoor | FlgNoDesc | FlgInvis,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	woodenDoor = Object{
		In:         &livingRoom,
		Synonyms:   []string{"door", "lettering", "writing"},
		Adjectives: []string{"wooden", "gothic", "strange", "west"},
		Desc:       "wooden door",
		Flags:      FlgRead | FlgDoor | FlgNoDesc | FlgTrans,
		// Action set in finalizeGameObjects to avoid init cycle
		Text: "The engravings translate to \"This space intentionally left blank.\"",
	}
	sword = Object{
		In:         &livingRoom,
		Synonyms:   []string{"sword", "orcrist", "glamdring", "blade"},
		Adjectives: []string{"elvish", "old", "antique"},
		Desc:       "sword",
		Flags:      FlgTake | FlgWeapon | FlgTryTake,
		FirstDesc:  "Above the trophy case hangs an elvish sword of great antiquity.",
		// Action set in finalizeGameObjects to avoid init cycle
		Item: &ItemData{Size: 30, TValue: 0},
	}
	lamp = Object{
		In:         &livingRoom,
		Synonyms:   []string{"lamp", "lantern", "light"},
		Adjectives: []string{"brass"},
		Desc:       "brass lantern",
		Flags:      FlgTake | FlgLight,
		FirstDesc:  "A battery-powered brass lantern is on the trophy case.",
		// Action set in finalizeGameObjects to avoid init cycle
		LongDesc: "There is a brass lantern (battery-powered) here.",
		Item:     &ItemData{Size: 15},
	}

	// kitchen
	kitchenTable = Object{
		In:         &kitchen,
		Synonyms:   []string{"table"},
		Adjectives: []string{"kitchen"},
		Desc:       "kitchen table",
		Flags:      FlgNoDesc | FlgCont | FlgOpen | FlgSurf,
		Item:       &ItemData{Capacity: 50},
	}

	// attic
	atticTable = Object{
		In:       &attic,
		Synonyms: []string{"table"},
		Desc:     "table",
		Flags:    FlgNoDesc | FlgCont | FlgOpen | FlgSurf,
		Item:     &ItemData{Capacity: 40},
	}
	rope = Object{
		In:         &attic,
		Synonyms:   []string{"rope", "hemp", "coil"},
		Adjectives: []string{"large"},
		Desc:       "rope",
		Flags:      FlgTake | FlgSacred | FlgTryTake,
		// Action set in finalizeGameObjects to avoid init cycle
		FirstDesc: "A large coil of rope is lying in the corner.",
		Item:      &ItemData{Size: 10},
	}

	// Entrance to Hades
	ghosts = Object{
		In:         &enteranceToHades,
		Synonyms:   []string{"ghosts", "spirits", "fiends", "force"},
		Adjectives: []string{"invisible", "evil"},
		Desc:       "number of ghosts",
		Flags:      FlgPerson | FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}

	// Land of Living Dead
	skull = Object{
		In:         &landOfLivingDead,
		Synonyms:   []string{"skull", "head", "treasure"},
		Adjectives: []string{"crystal"},
		Desc:       "crystal skull",
		FirstDesc:  "Lying in one corner of the room is a beautifully carved crystal skull. it appears to be grinning at you rather nastily.",
		Flags:      FlgTake,
		Item:       &ItemData{Value: 10, TValue: 10},
	}

	// Shaft Room
	raisedBasket = Object{
		In:       &shaftRoom,
		Synonyms: []string{"cage", "dumbwaiter", "basket"},
		Desc:     "basket",
		Flags:    FlgTrans | FlgTryTake | FlgCont | FlgOpen,
		// Action set in finalizeGameObjects to avoid init cycle
		LongDesc: "At the end of the chain is a basket.",
		Item:     &ItemData{Capacity: 50},
	}

	// Lower Shaft
	loweredBasket = Object{
		In:         &lowerShaft,
		Synonyms:   []string{"cage", "dumbwaiter", "basket"},
		Adjectives: []string{"lowered"},
		Desc:       "basket",
		LongDesc:   "From the chain is suspended a basket.",
		Flags:      FlgTryTake,
		// Action set in finalizeGameObjects to avoid init cycle
	}

	// bat Room
	bat = Object{
		In:         &batRoom,
		Synonyms:   []string{"bat", "vampire"},
		Adjectives: []string{"vampire", "deranged"},
		Desc:       "bat",
		Flags:      FlgPerson | FlgTryTake,
		DescFcn:    batDescFcn,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	jade = Object{
		In:         &batRoom,
		Synonyms:   []string{"figurine", "treasure"},
		Adjectives: []string{"exquisite", "jade"},
		Desc:       "jade figurine",
		Flags:      FlgTake,
		LongDesc:   "There is an exquisite jade figurine here.",
		Item:       &ItemData{Size: 10, Value: 5, TValue: 5},
	}

	// North Temple
	bell = Object{
		In:         &northTemple,
		Synonyms:   []string{"bell"},
		Adjectives: []string{"small", "brass"},
		Desc:       "brass bell",
		Flags:      FlgTake,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	prayer = Object{
		In:         &northTemple,
		Synonyms:   []string{"prayer", "inscription"},
		Adjectives: []string{"ancient", "old"},
		Desc:       "prayer",
		Flags:      FlgRead | FlgSacred | FlgNoDesc,
		Text:       "The prayer is inscribed in an ancient script, rarely used today. it seems to be a philippic against small insects, absent-mindedness, and the picking up and dropping of small objects. The final verse consigns trespassers to the land of the dead. All evidence indicates that the beliefs of the ancient Zorkers were obscure.",
	}

	// South Temple
	altar = Object{
		In:       &southTemple,
		Synonyms: []string{"altar"},
		Desc:     "altar",
		Flags:    FlgNoDesc | FlgSurf | FlgCont | FlgOpen,
		Item:     &ItemData{Capacity: 50},
	}
	candles = Object{
		In:         &southTemple,
		Synonyms:   []string{"candles", "pair"},
		Adjectives: []string{"burning"},
		Desc:       "pair of candles",
		Flags:      FlgTake | FlgFlame | FlgOn | FlgLight,
		// Action set in finalizeGameObjects to avoid init cycle
		FirstDesc: "On the two ends of the altar are burning candles.",
		Item:      &ItemData{Size: 10},
	}

	// troll Room
	troll = Object{
		In:         &trollRoom,
		Synonyms:   []string{"troll"},
		Adjectives: []string{"nasty"},
		Desc:       "troll",
		Flags:      FlgPerson | FlgOpen | FlgTryTake,
		// Action set in finalizeGameObjects to avoid init cycle
		LongDesc: "A nasty-looking troll, brandishing a bloody axe, blocks all passages out of the room.",
		Combat:   &CombatData{Strength: 2},
	}

	// dam Room
	bolt = Object{
		In:         &damRoom,
		Synonyms:   []string{"bolt", "nut"},
		Adjectives: []string{"metal", "large"},
		Desc:       "bolt",
		Flags:      FlgNoDesc | FlgTurn | FlgTryTake,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	bubble = Object{
		In:         &damRoom,
		Synonyms:   []string{"bubble"},
		Adjectives: []string{"small", "green", "plastic"},
		Desc:       "green bubble",
		Flags:      FlgNoDesc | FlgTryTake,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	dam = Object{
		In:       &damRoom,
		Synonyms: []string{"dam", "gate", "gates", "fcd#3"},
		Desc:     "dam",
		Flags:    FlgNoDesc | FlgTryTake,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	controlPanel = Object{
		In:         &damRoom,
		Synonyms:   []string{"panel"},
		Adjectives: []string{"control"},
		Desc:       "control panel",
		Flags:      FlgNoDesc,
	}

	// dam Lobby
	match = Object{
		In:         &damLobby,
		Synonyms:   []string{"match", "matches", "matchbook"},
		Adjectives: []string{"match"},
		Desc:       "matchbook",
		Flags:      FlgRead | FlgTake,
		LongDesc:   "There is a matchbook whose cover says \"Visit Beautiful FCD#3\" here.",
		// Action set in finalizeGameObjects to avoid init cycle
		Text: "\n(Close cover before striking)\n\nYOU too can make BIG MONEY in the exciting field of PAPER SHUFFLING!\n\nMr. Anderson of Muddle, Mass. says: \"Before I took this course I was a lowly bit twiddler. Now with what I learned at GUE Tech I feel really important and can obfuscate and confuse with the best.\"\n\nDr. Blank had this to say: \"Ten short days ago all I could look forward to was a dead-end job as a doctor. Now I have a promising future and make really big Zorkmids.\"\n\nGUE Tech can't promise these fantastic results to everyone. But when you earn your degree from GUE Tech, your future will be brighter.",
		Item: &ItemData{Size: 2},
	}
	guide = Object{
		In:         &damLobby,
		Synonyms:   []string{"guide", "book", "books", "guidebooks"},
		Adjectives: []string{"tour", "guide"},
		Desc:       "tour guidebook",
		Flags:      FlgRead | FlgTake | FlgBurn,
		FirstDesc:  "Some guidebooks entitled \"Flood Control dam #3\" are on the reception desk.",
		Text:       "\"\tFlood Control dam #3\n\nFCD#3 was constructed in year 783 of the Great Underground Empire to harness the mighty Frigid river. This work was supported by a grant of 37 million zorkmids from your omnipotent local tyrant Lord Dimwit Flathead the Excessive. This impressive structure is composed of 370,000 cubic feet of concrete, is 256 feet tall at the center, and 193 feet wide at the top. The lake created behind the dam has a volume of 1.7 billion cubic feet, an area of 12 million square feet, and a shore line of 36 thousand feet.\n\nThe construction of FCD#3 took 112 days from ground breaking to the dedication. it required a work force of 384 slaves, 34 slave drivers, 12 engineers, 2 turtle doves, and a partridge in a pear tree. The work was managed by a command team composed of 2345 bureaucrats, 2347 secretaries (at least two of whom could type), 12,256 paper shufflers, 52,469 rubber stampers, 245,193 red tape processors, and nearly one million dead trees.\n\nWe will now point out some of the more interesting features of FCD#3 as we conduct you on a guided tour of the facilities:\n\n        1) You start your tour here in the dam Lobby. You will notice on your right that....\"",
	}

	// dam Base
	inflatableBoat = Object{
		In:         &damBase,
		Synonyms:   []string{"boat", "pile", "plastic", "valve"},
		Adjectives: []string{"plastic", "inflatable"},
		Desc:       "pile of plastic",
		Flags:      FlgTake | FlgBurn,
		// Action set in finalizeGameObjects to avoid init cycle
		LongDesc: "There is a folded pile of plastic here which has a small valve attached.",
		Item:     &ItemData{Size: 20},
	}

	// Maintenance Room
	toolChest = Object{
		In:         &maintenanceRoom,
		Synonyms:   []string{"chest", "chests", "group", "toolchests"},
		Adjectives: []string{"tool"},
		Desc:       "group of tool chests",
		Flags:      FlgCont | FlgOpen | FlgTryTake | FlgSacred,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	yellowButton = Object{
		In:         &maintenanceRoom,
		Synonyms:   []string{"button", "switch"},
		Adjectives: []string{"yellow"},
		Desc:       "yellow button",
		Flags:      FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	brownButton = Object{
		In:         &maintenanceRoom,
		Synonyms:   []string{"button", "switch"},
		Adjectives: []string{"brown"},
		Desc:       "brown button",
		Flags:      FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	redButton = Object{
		In:         &maintenanceRoom,
		Synonyms:   []string{"button", "switch"},
		Adjectives: []string{"red"},
		Desc:       "red button",
		Flags:      FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	blueButton = Object{
		In:         &maintenanceRoom,
		Synonyms:   []string{"button", "switch"},
		Adjectives: []string{"blue"},
		Desc:       "blue button",
		Flags:      FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	screwdriver = Object{
		In:         &maintenanceRoom,
		Synonyms:   []string{"screwdriver", "tool", "tools", "driver"},
		Adjectives: []string{"screw"},
		Desc:       "screwdriver",
		Flags:      FlgTake | FlgTool,
	}
	wrench = Object{
		In:       &maintenanceRoom,
		Synonyms: []string{"wrench", "tool", "tools"},
		Desc:     "wrench",
		Flags:    FlgTake | FlgTool,
		Item:     &ItemData{Size: 10},
	}
	tube = Object{
		In:       &maintenanceRoom,
		Synonyms: []string{"tube", "tooth", "paste"},
		Desc:     "tube",
		Flags:    FlgTake | FlgCont | FlgRead,
		// Action set in finalizeGameObjects to avoid init cycle
		LongDesc: "There is an object which looks like a tube of toothpaste here.",
		Text:     "---> Frobozz Magic gunk Company <---\n\tAll-Purpose gunk",
		Item:     &ItemData{Size: 5, Capacity: 7},
	}
	leak = Object{
		In:       &maintenanceRoom,
		Synonyms: []string{"leak", "drip", "pipe"},
		Desc:     "leak",
		Flags:    FlgNoDesc | FlgInvis,
		// Action set in finalizeGameObjects to avoid init cycle
	}

	// machine Room
	machine = Object{
		In:       &machineRoom,
		Synonyms: []string{"machine", "pdp10", "dryer", "lid"},
		Desc:     "machine",
		Flags:    FlgCont | FlgNoDesc | FlgTryTake,
		// Action set in finalizeGameObjects to avoid init cycle
		Item: &ItemData{Capacity: 50},
	}
	machineSwitch = Object{
		In:       &machineRoom,
		Synonyms: []string{"switch"},
		Desc:     "switch",
		Flags:    FlgNoDesc | FlgTurn,
		// Action set in finalizeGameObjects to avoid init cycle
	}

	// cyclops Room
	cyclops = Object{
		In:         &cyclopsRoom,
		Synonyms:   []string{"cyclops", "monster", "eye"},
		Adjectives: []string{"hungry", "giant"},
		Desc:       "cyclops",
		Flags:      FlgPerson | FlgNoDesc | FlgTryTake,
		// Action set in finalizeGameObjects to avoid init cycle
		Combat: &CombatData{Strength: 10000},
	}

	// Treasure Room
	chalice = Object{
		In:         &treasureRoom,
		Synonyms:   []string{"chalice", "cup", "silver", "treasure"},
		Adjectives: []string{"silver", "engravings"},
		Desc:       "chalice",
		Flags:      FlgTake | FlgTryTake | FlgCont,
		// Action set in finalizeGameObjects to avoid init cycle
		LongDesc: "There is a silver chalice, intricately engraved, here.",
		Item:     &ItemData{Size: 10, Value: 10, TValue: 5, Capacity: 5},
	}

	// gallery
	painting = Object{
		In:         &gallery,
		Synonyms:   []string{"painting", "art", "canvas", "treasure"},
		Adjectives: []string{"beautiful"},
		Desc:       "painting",
		Flags:      FlgTake | FlgBurn,
		// Action set in finalizeGameObjects to avoid init cycle
		FirstDesc: "Fortunately, there is still one chance for you to be a vandal, for on the far wall is a painting of unparalleled beauty.",
		LongDesc:  "A painting by a neglected genius is here.",
		Item:      &ItemData{Size: 15, Value: 4, TValue: 6},
	}

	// studio
	ownersManual = Object{
		In:         &studio,
		Synonyms:   []string{"manual", "piece", "paper"},
		Adjectives: []string{"zork", "owners", "small"},
		Desc:       "ZORK owner's manual",
		Flags:      FlgRead | FlgTake,
		FirstDesc:  "Loosely attached to a wall is a small piece of paper.",
		Text:       "Congratulations!\n\nYou are the privileged owner of ZORK I: The Great Underground Empire, a self-contained and self-maintaining universe. If used and maintained in accordance with normal operating practices for small universes, ZORK will provide many months of trouble-free operation.",
	}

	// Grating clearing
	leaves = Object{
		In:       &gratingClearing,
		Synonyms: []string{"leaves", "leaf", "pile"},
		Desc:     "pile of leaves",
		Flags:    FlgTake | FlgBurn | FlgTryTake,
		// Action set in finalizeGameObjects to avoid init cycle
		LongDesc: "On the ground is a pile of leaves.",
		Item:     &ItemData{Size: 25},
	}

	// Up a tree
	nest = Object{
		In:         &upATree,
		Synonyms:   []string{"nest"},
		Adjectives: []string{"birds"},
		Desc:       "bird's nest",
		Flags:      FlgTake | FlgBurn | FlgCont | FlgOpen | FlgSearch,
		FirstDesc:  "Beside you on the branch is a small bird's nest.",
		Item:       &ItemData{Capacity: 20},
	}

	// Sandy Cave
	sand = Object{
		In:       &sandyCave,
		Synonyms: []string{"sand"},
		Desc:     "sand",
		Flags:    FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	scarab = Object{
		In:         &sandyCave,
		Synonyms:   []string{"scarab", "bug", "beetle", "treasure"},
		Adjectives: []string{"beautiful", "carved", "jeweled"},
		Desc:       "beautiful jeweled scarab",
		Flags:      FlgTake | FlgInvis,
		Item:       &ItemData{Size: 8, Value: 5, TValue: 5},
	}

	// Sandy Beach
	shovel = Object{
		In:       &sandyBeach,
		Synonyms: []string{"shovel", "tool", "tools"},
		Desc:     "shovel",
		Flags:    FlgTake | FlgTool,
		Item:     &ItemData{Size: 15},
	}

	// Egypt Room
	coffin = Object{
		In:         &egyptRoom,
		Synonyms:   []string{"coffin", "casket", "treasure"},
		Adjectives: []string{"solid", "gold"},
		Desc:       "gold coffin",
		Flags:      FlgTake | FlgCont | FlgSacred | FlgSearch,
		LongDesc:   "The solid-gold coffin used for the burial of Ramses II is here.",
		Item:       &ItemData{Size: 55, Value: 10, TValue: 15, Capacity: 35},
	}

	// Round Room
	thief = Object{
		In:         &roundRoom,
		Synonyms:   []string{"thief", "robber", "man", "person"},
		Adjectives: []string{"shady", "suspicious", "seedy"},
		Desc:       "thief",
		Flags:      FlgPerson | FlgInvis | FlgCont | FlgOpen | FlgTryTake,
		// Action set in finalizeGameObjects to avoid init cycle
		LongDesc: "There is a suspicious-looking individual, holding a large bag, leaning against one wall. He is armed with a deadly stiletto.",
		Combat:   &CombatData{Strength: 5},
	}

	// reservoir
	trunk = Object{
		In:         &reservoir,
		Synonyms:   []string{"trunk", "chest", "jewels", "treasure"},
		Adjectives: []string{"old"},
		Desc:       "trunk of jewels",
		Flags:      FlgTake | FlgInvis,
		FirstDesc:  "Lying half buried in the mud is an old trunk, bulging with jewels.",
		LongDesc:   "There is an old trunk here, bulging with assorted jewels.",
		// Action set in finalizeGameObjects to avoid init cycle
		Item: &ItemData{Size: 35, Value: 15, TValue: 5},
	}

	// reservoir North
	pump = Object{
		In:         &reservoirNorth,
		Synonyms:   []string{"pump", "air-pump", "tool", "tools"},
		Adjectives: []string{"small", "hand-held"},
		Desc:       "hand-held air pump",
		Flags:      FlgTake | FlgTool,
	}

	// Atlantis Room
	trident = Object{
		In:         &atlantisRoom,
		Synonyms:   []string{"trident", "fork", "treasure"},
		Adjectives: []string{"poseidon", "own", "crystal"},
		Desc:       "crystal trident",
		Flags:      FlgTake,
		FirstDesc:  "On the shore lies Poseidon's own crystal trident.",
		Item:       &ItemData{Size: 20, Value: 4, TValue: 11},
	}

	// Mirror rooms
	mirror1 = Object{
		In:       &mirrorRoom1,
		Synonyms: []string{"reflection", "mirror", "enormous"},
		Desc:     "mirror",
		Flags:    FlgTryTake | FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	mirror2 = Object{
		In:       &mirrorRoom2,
		Synonyms: []string{"reflection", "mirror", "enormous"},
		Desc:     "mirror",
		Flags:    FlgTryTake | FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}

	// Dome Room
	railing = Object{
		In:         &domeRoom,
		Synonyms:   []string{"railing", "rail"},
		Adjectives: []string{"wooden"},
		Desc:       "wooden railing",
		Flags:      FlgNoDesc,
	}

	// torch Room
	pedestal = Object{
		In:         &torchRoom,
		Synonyms:   []string{"pedestal"},
		Adjectives: []string{"white", "marble"},
		Desc:       "pedestal",
		Flags:      FlgNoDesc | FlgCont | FlgOpen | FlgSurf,
		// Action set in finalizeGameObjects to avoid init cycle
		Item: &ItemData{Capacity: 30},
	}

	// engravings Cave
	engravings = Object{
		In:         &engravingsCave,
		Synonyms:   []string{"wall", "engravings", "inscription"},
		Adjectives: []string{"old", "ancient"},
		Desc:       "wall with engravings",
		Flags:      FlgRead | FlgSacred,
		LongDesc:   "There are old engravings on the walls here.",
		Text:       "The engravings were incised in the living rock of the cave wall by an unknown hand. They depict, in symbolic form, the beliefs of the ancient Zorkers. Skillfully interwoven with the bas reliefs are excerpts illustrating the major religious tenets of that time. Unfortunately, a later age seems to have considered them blasphemous and just as skillfully excised them.",
	}

	// Loud Room
	bar = Object{
		In:         &loudRoom,
		Synonyms:   []string{"bar", "platinum", "treasure"},
		Adjectives: []string{"platinum", "large"},
		Desc:       "platinum bar",
		Flags:      FlgTake | FlgSacred,
		LongDesc:   "On the ground is a large platinum bar.",
		Item:       &ItemData{Size: 20, Value: 10, TValue: 5},
	}

	// End of rainbow
	potOfGold = Object{
		In:         &endOfRainbow,
		Synonyms:   []string{"pot", "gold", "treasure"},
		Adjectives: []string{"gold"},
		Desc:       "pot of gold",
		Flags:      FlgTake | FlgInvis,
		FirstDesc:  "At the end of the rainbow is a pot of gold.",
		Item:       &ItemData{Size: 15, Value: 10, TValue: 10},
	}

	// river 4
	buoy = Object{
		In:         &river4,
		Synonyms:   []string{"buoy"},
		Adjectives: []string{"red"},
		Desc:       "red buoy",
		Flags:      FlgTake | FlgCont,
		FirstDesc:  "There is a red buoy here (probably a warning).",
		// Action set in finalizeGameObjects to avoid init cycle
		Item: &ItemData{Size: 10, Capacity: 20},
	}

	// Gas Room
	bracelet = Object{
		In:         &gasRoom,
		Synonyms:   []string{"bracelet", "jewel", "sapphire", "treasure"},
		Adjectives: []string{"sapphire"},
		Desc:       "sapphire-encrusted bracelet",
		Flags:      FlgTake,
		Item:       &ItemData{Size: 10, Value: 5, TValue: 5},
	}

	// Dead End 5
	coal = Object{
		In:         &deadEnd5,
		Synonyms:   []string{"coal", "pile", "heap"},
		Adjectives: []string{"small"},
		Desc:       "small pile of coal",
		Flags:      FlgTake | FlgBurn,
		Item:       &ItemData{Size: 20},
	}

	// Timber Room
	timbers = Object{
		In:         &timberRoom,
		Synonyms:   []string{"timbers", "pile"},
		Adjectives: []string{"wooden", "broken"},
		Desc:       "broken timber",
		Flags:      FlgTake,
		Item:       &ItemData{Size: 50},
	}

	// Maze 5
	bones = Object{
		In:       &maze5,
		Synonyms: []string{"bones", "skeleton", "body"},
		Desc:     "skeleton",
		Flags:    FlgTryTake | FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	burnedOutLantern = Object{
		In:         &maze5,
		Synonyms:   []string{"lantern", "lamp"},
		Adjectives: []string{"rusty", "burned", "dead", "useless"},
		Desc:       "burned-out lantern",
		Flags:      FlgTake,
		FirstDesc:  "The deceased adventurer's useless lantern is here.",
		Item:       &ItemData{Size: 20},
	}
	bagOfCoins = Object{
		In:         &maze5,
		Synonyms:   []string{"bag", "coins", "treasure"},
		Adjectives: []string{"old", "leather"},
		Desc:       "leather bag of coins",
		Flags:      FlgTake,
		LongDesc:   "An old leather bag, bulging with coins, is here.",
		// Action set in finalizeGameObjects to avoid init cycle
		Item: &ItemData{Size: 15, Value: 10, TValue: 5},
	}
	rustyKnife = Object{
		In:         &maze5,
		Synonyms:   []string{"knives", "knife"},
		Adjectives: []string{"rusty"},
		Desc:       "rusty knife",
		Flags:      FlgTake | FlgTryTake | FlgWeapon | FlgTool,
		// Action set in finalizeGameObjects to avoid init cycle
		FirstDesc: "Beside the skeleton is a rusty knife.",
		Item:      &ItemData{Size: 20},
	}
	keys = Object{
		In:         &maze5,
		Synonyms:   []string{"key"},
		Adjectives: []string{"skeleton"},
		Desc:       "skeleton key",
		Flags:      FlgTake | FlgTool,
		Item:       &ItemData{Size: 10},
	}

	// ================================================================
	// OBJECTS IN OBJECTS
	// ================================================================

	// In Trophy Case
	mapObj = Object{
		In:         &trophyCase,
		Synonyms:   []string{"parchment", "map"},
		Adjectives: []string{"antique", "old", "ancient"},
		Desc:       "ancient map",
		Flags:      FlgInvis | FlgRead | FlgTake,
		FirstDesc:  "In the trophy case is an ancient parchment which appears to be a map.",
		Text:       "The map shows a forest with three clearings. The largest clearing contains a house. Three paths leave the large clearing. One of these paths, leading southwest, is marked \"To Stone barrow\".",
		Item:       &ItemData{Size: 2},
	}

	// In mailbox
	advertisement = Object{
		In:         &mailbox,
		Synonyms:   []string{"advertisement", "leaflet", "booklet", "mail"},
		Adjectives: []string{"small"},
		Desc:       "leaflet",
		Flags:      FlgRead | FlgTake | FlgBurn,
		LongDesc:   "A small leaflet is on the ground.",
		Text:       "\"WELCOME TO ZORK!\n\nZORK is a game of adventure, danger, and low cunning. In it you will explore some of the most amazing territory ever seen by mortals. No computer should be without one!\"",
		Item:       &ItemData{Size: 2},
	}

	// In kitchen Table
	bottle = Object{
		In:         &kitchenTable,
		Synonyms:   []string{"bottle", "container"},
		Adjectives: []string{"clear", "glass"},
		Desc:       "glass bottle",
		Flags:      FlgTake | FlgTrans | FlgCont,
		// Action set in finalizeGameObjects to avoid init cycle
		FirstDesc: "A bottle is sitting on the table.",
		Item:      &ItemData{Capacity: 4},
	}
	sandwichBag = Object{
		In:         &kitchenTable,
		Synonyms:   []string{"bag", "sack"},
		Adjectives: []string{"brown", "elongated", "smelly"},
		Desc:       "brown sack",
		Flags:      FlgTake | FlgCont | FlgBurn,
		FirstDesc:  "On the table is an elongated brown sack, smelling of hot peppers.",
		// Action set in finalizeGameObjects to avoid init cycle
		Item: &ItemData{Size: 9, Capacity: 9},
	}

	// In attic Table
	knife = Object{
		In:         &atticTable,
		Synonyms:   []string{"knives", "knife", "blade"},
		Adjectives: []string{"nasty", "unrusty"},
		Desc:       "nasty knife",
		Flags:      FlgTake | FlgWeapon | FlgTryTake,
		FirstDesc:  "On a table is a nasty-looking knife.",
		// Action set in finalizeGameObjects to avoid init cycle
	}

	// In bottle
	water = Object{
		In:       &bottle,
		Synonyms: []string{"water", "quantity", "liquid", "h2o"},
		Desc:     "quantity of water",
		Flags:    FlgTryTake | FlgTake | FlgDrink,
		// Action set in finalizeGameObjects to avoid init cycle
		Item: &ItemData{Size: 4},
	}

	// In Sandwich Bag
	lunch = Object{
		In:         &sandwichBag,
		Synonyms:   []string{"food", "sandwich", "lunch", "dinner"},
		Adjectives: []string{"hot", "pepper"},
		Desc:       "lunch",
		Flags:      FlgTake | FlgFood,
		LongDesc:   "A hot pepper sandwich is here.",
	}
	garlic = Object{
		In:       &sandwichBag,
		Synonyms: []string{"garlic", "clove"},
		Desc:     "clove of garlic",
		Flags:    FlgTake | FlgFood,
		// Action set in finalizeGameObjects to avoid init cycle
		Item: &ItemData{Size: 4},
	}

	// In altar
	book = Object{
		In:         &altar,
		Synonyms:   []string{"book", "prayer", "page", "books"},
		Adjectives: []string{"large", "black"},
		Desc:       "black book",
		Flags:      FlgRead | FlgTake | FlgCont | FlgBurn | FlgTurn,
		// Action set in finalizeGameObjects to avoid init cycle
		FirstDesc: "On the altar is a large black book, open to page 569.",
		Text:      "Commandment #12592\n\nOh ye who go about saying unto each: \"Hello sailor\":\nDost thou know the magnitude of thy sin before the gods?\nYea, verily, thou shalt be ground between two stones.\nShall the angry gods cast thy body into the whirlpool?\nSurely, thy eye shall be put out with a sharp stick!\nEven unto the ends of the earth shalt thou wander and\nUnto the land of the dead shalt thou be sent at last.\nSurely thou shalt repent of thy cunning.",
		Item:      &ItemData{Size: 10},
	}

	// In coffin
	sceptre = Object{
		In:         &coffin,
		Synonyms:   []string{"sceptre", "scepter", "treasure"},
		Adjectives: []string{"sharp", "egyptian", "ancient", "enameled"},
		Desc:       "sceptre",
		Flags:      FlgTake | FlgWeapon,
		LongDesc:   "An ornamented sceptre, tapering to a sharp point, is here.",
		// Action set in finalizeGameObjects to avoid init cycle
		FirstDesc: "A sceptre, possibly that of ancient Egypt itself, is in the coffin. The sceptre is ornamented with colored enamel, and tapers to a sharp point.",
		Item:      &ItemData{Size: 3, Value: 4, TValue: 6},
	}

	// In nest
	egg = Object{
		In:         &nest,
		Synonyms:   []string{"egg", "treasure"},
		Adjectives: []string{"birds", "encrusted", "jeweled"},
		Desc:       "jewel-encrusted egg",
		Flags:      FlgTake | FlgCont | FlgSearch,
		// Action set in finalizeGameObjects to avoid init cycle
		FirstDesc: "In the bird's nest is a large egg encrusted with precious jewels, apparently scavenged by a childless songbird. The egg is covered with fine gold inlay, and ornamented in lapis lazuli and mother-of-pearl. Unlike most eggs, this one is hinged and closed with a delicate looking clasp. The egg appears extremely fragile.",
		Item:      &ItemData{Value: 5, TValue: 5, Capacity: 6},
	}

	// In egg
	canary = Object{
		In:         &egg,
		Synonyms:   []string{"canary", "treasure"},
		Adjectives: []string{"clockwork", "gold", "golden"},
		Desc:       "golden clockwork canary",
		Flags:      FlgTake | FlgSearch,
		// Action set in finalizeGameObjects to avoid init cycle
		FirstDesc: "There is a golden clockwork canary nestled in the egg. it has ruby eyes and a silver beak. through a crystal window below its left wing you can see intricate machinery inside. it appears to have wound down.",
		Item:      &ItemData{Value: 6, TValue: 4},
	}

	// In tube
	putty = Object{
		In:         &tube,
		Synonyms:   []string{"material", "gunk"},
		Adjectives: []string{"viscous"},
		Desc:       "viscous material",
		Flags:      FlgTake | FlgTool,
		// Action set in finalizeGameObjects to avoid init cycle
		Item: &ItemData{Size: 6},
	}

	// In troll
	axe = Object{
		In:         &troll,
		Synonyms:   []string{"axe", "ax"},
		Adjectives: []string{"bloody"},
		Desc:       "bloody axe",
		Flags:      FlgWeapon | FlgTryTake | FlgTake | FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
		Item: &ItemData{Size: 25},
	}

	// In thief
	largeBag = Object{
		In:         &thief,
		Synonyms:   []string{"bag"},
		Adjectives: []string{"large", "thiefs"},
		Desc:       "large bag",
		Flags:      FlgTryTake | FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	stiletto = Object{
		In:         &thief,
		Synonyms:   []string{"stiletto"},
		Adjectives: []string{"vicious"},
		Desc:       "stiletto",
		Flags:      FlgWeapon | FlgTryTake | FlgTake | FlgNoDesc,
		// Action set in finalizeGameObjects to avoid init cycle
		Item: &ItemData{Size: 10},
	}

	// In pedestal
	torch = Object{
		In:         &pedestal,
		Synonyms:   []string{"torch", "ivory", "treasure"},
		Adjectives: []string{"flaming", "ivory"},
		Desc:       "torch",
		Flags:      FlgTake | FlgFlame | FlgOn | FlgLight,
		// Action set in finalizeGameObjects to avoid init cycle
		FirstDesc: "Sitting on the pedestal is a flaming torch, made of ivory.",
		Item:      &ItemData{Size: 20, Value: 14, TValue: 6},
	}

	// In Inflated Boat
	boatLabel = Object{
		In:         &inflatedBoat,
		Synonyms:   []string{"label", "fineprint", "print"},
		Adjectives: []string{"tan", "fine"},
		Desc:       "tan label",
		Flags:      FlgRead | FlgTake | FlgBurn,
		Text:       "  !!!!FROBOZZ MAGIC BOAT COMPANY!!!!\n\nHello, sailor!\n\nInstructions for use:\n\n   To get into a body of water, say \"Launch\".\n   To get to shore, say \"Land\" or the direction in which you want to maneuver the boat.\n\nWarranty:\n\n  This boat is guaranteed against all defects for a period of 76 milliseconds from date of purchase or until first used, whichever comes first.\n\nWarning:\n   This boat is made of thin plastic.\n   Good Luck!",
		Item:       &ItemData{Size: 2},
	}

	// In buoy
	emerald = Object{
		In:         &buoy,
		Synonyms:   []string{"emerald", "treasure"},
		Adjectives: []string{"large"},
		Desc:       "large emerald",
		Flags:      FlgTake,
		Item:       &ItemData{Value: 5, TValue: 10},
	}

	// In Broken egg (unplaced)
	brokenCanary = Object{
		In:         &brokenEgg,
		Synonyms:   []string{"canary", "treasure"},
		Adjectives: []string{"broken", "clockwork", "gold", "golden"},
		Desc:       "broken clockwork canary",
		Flags:      FlgTake,
		// Action set in finalizeGameObjects to avoid init cycle
		FirstDesc: "There is a golden clockwork canary nestled in the egg. it seems to have recently had a bad experience. The mountings for its jewel-like eyes are empty, and its silver beak is crumpled. through a cracked crystal window below its left wing you can see the remains of intricate machinery. it is not clear what result winding it would have, as the mainspring seems sprung.",
		Item:      &ItemData{TValue: 1},
	}
)
