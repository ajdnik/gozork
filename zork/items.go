package zork

var (
	// ================================================================
	// GLOBAL OBJECTS (In: &LocalGlobals or &GlobalObjects)
	// ================================================================

	Board = Object{
		In:       &LocalGlobals,
		Synonyms: []string{"boards", "board"},
		Desc:     "board",
		Flags:    FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Teeth = Object{
		In:       &GlobalObjects,
		Synonyms: []string{"overboard", "teeth"},
		Desc:     "set of teeth",
		Flags:    FlgNoDesc,
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
		Flags:      FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	WhiteHouse = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"house"},
		Adjectives: []string{"white", "beautiful", "colonial"},
		Desc:       "white house",
		Flags:      FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Forest = Object{
		In:       &LocalGlobals,
		Synonyms: []string{"forest", "trees", "pines", "hemlocks"},
		Desc:     "forest",
		Flags:    FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Tree = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"tree", "branch"},
		Adjectives: []string{"large", "storm"},
		Desc:       "tree",
		Flags:      FlgNoDesc | FlgClimb,
	}
	GlobalWater = Object{
		In:       &LocalGlobals,
		Synonyms: []string{"water", "quantity"},
		Desc:     "water",
		Flags:    FlgDrink,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	KitchenWindow = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"window"},
		Adjectives: []string{"kitchen", "small"},
		Desc:       "kitchen window",
		Flags:      FlgDoor | FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Chimney = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"chimney"},
		Adjectives: []string{"dark", "narrow"},
		Desc:       "chimney",
		Flags:      FlgClimb | FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Slide = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"chute", "ramp", "slide"},
		Adjectives: []string{"steep", "metal", "twisting"},
		Desc:       "chute",
		Flags:      FlgClimb,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Bodies = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"bodies", "body", "remains", "pile"},
		Adjectives: []string{"mangled"},
		Desc:       "pile of bodies",
		Flags:      FlgNoDesc | FlgTryTake,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Crack = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"crack"},
		Adjectives: []string{"narrow"},
		Desc:       "crack",
		Flags:      FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Grate = Object{
		In:       &LocalGlobals,
		Synonyms: []string{"grate", "grating"},
		Desc:     "grating",
		Flags:    FlgDoor | FlgNoDesc | FlgInvis,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Ladder = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"ladder"},
		Adjectives: []string{"wooden", "rickety", "narrow"},
		Desc:       "wooden ladder",
		Flags:      FlgNoDesc | FlgClimb,
	}
	ClimbableCliff = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"wall", "cliff", "walls", "ledge"},
		Adjectives: []string{"rocky", "sheer"},
		Desc:       "cliff",
		Flags:      FlgNoDesc | FlgClimb,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	WhiteCliff = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"cliff", "cliffs"},
		Adjectives: []string{"white"},
		Desc:       "white cliffs",
		Flags:      FlgNoDesc | FlgClimb,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Rainbow = Object{
		In:       &LocalGlobals,
		Synonyms: []string{"rainbow"},
		Desc:     "rainbow",
		Flags:    FlgNoDesc | FlgClimb,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	River = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"river"},
		Adjectives: []string{"frigid"},
		Desc:       "river",
		Flags:      FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	BoardedWindow = Object{
		In:         &LocalGlobals,
		Synonyms:   []string{"window"},
		Adjectives: []string{"boarded"},
		Desc:       "boarded window",
		Flags:      FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// ================================================================
	// UNPLACED OBJECTS (created or swapped during gameplay)
	// ================================================================

	InflatedBoat = Object{
		Synonyms:   []string{"boat", "raft"},
		Adjectives: []string{"magic", "plastic", "seaworthy", "inflated", "inflatable"},
		Desc:       "magic boat",
		Flags:      FlgTake | FlgBurn | FlgVeh | FlgOpen | FlgSearch,
		Capacity:   100,
		// Action set in FinalizeGameObjects to avoid init cycle
		Size:       20,
		VehType:    FlgNonLand,
	}
	PuncturedBoat = Object{
		Synonyms:   []string{"boat", "pile", "plastic"},
		Adjectives: []string{"plastic", "puncture", "large"},
		Desc:       "punctured boat",
		Flags:      FlgTake | FlgBurn,
		// Action set in FinalizeGameObjects to avoid init cycle
		Size:       20,
	}
	BrokenLamp = Object{
		Synonyms:   []string{"lamp", "lantern"},
		Adjectives: []string{"broken"},
		Desc:       "broken lantern",
		Flags:      FlgTake,
	}
	Gunk = Object{
		Synonyms:   []string{"gunk", "piece", "slag"},
		Adjectives: []string{"small", "vitreous"},
		Desc:       "small piece of vitreous slag",
		Flags:      FlgTake | FlgTryTake,
		Size:       10,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	HotBell = Object{
		Synonyms:   []string{"bell"},
		Adjectives: []string{"brass", "hot", "red", "small"},
		Desc:       "red hot brass bell",
		Flags:      FlgTryTake,
		LongDesc:   "On the ground is a red hot bell.",
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	BrokenEgg = Object{
		Synonyms:   []string{"egg", "treasure"},
		Adjectives: []string{"broken", "birds", "encrusted", "jewel"},
		Desc:       "broken jewel-encrusted egg",
		Flags:      FlgTake | FlgCont | FlgOpen,
		Capacity:   6,
		TValue:     2,
		LongDesc:   "There is a somewhat ruined egg here.",
	}
	Bauble = Object{
		Synonyms:   []string{"bauble", "treasure"},
		Adjectives: []string{"brass", "beautiful"},
		Desc:       "beautiful brass bauble",
		Flags:      FlgTake,
		Value:      1,
		TValue:     1,
	}
	Diamond = Object{
		Synonyms:   []string{"diamond", "treasure"},
		Adjectives: []string{"huge", "enormous"},
		Desc:       "huge diamond",
		Flags:      FlgTake,
		LongDesc:   "There is an enormous diamond (perfectly cut) here.",
		Value:      10,
		TValue:     10,
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
		Flags:      FlgNoDesc | FlgClimb,
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// West of House
	FrontDoor = Object{
		In:         &WestOfHouse,
		Synonyms:   []string{"door"},
		Adjectives: []string{"front", "boarded"},
		Desc:       "door",
		Flags:      FlgDoor | FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Mailbox = Object{
		In:         &WestOfHouse,
		Synonyms:   []string{"mailbox", "box"},
		Adjectives: []string{"small"},
		Desc:       "small mailbox",
		Flags:      FlgCont | FlgTryTake,
		Capacity:   10,
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// Stone Barrow
	BarrowDoor = Object{
		In:         &StoneBarrow,
		Synonyms:   []string{"door"},
		Adjectives: []string{"huge", "stone"},
		Desc:       "stone door",
		Flags:      FlgDoor | FlgNoDesc | FlgOpen,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Barrow = Object{
		In:         &StoneBarrow,
		Synonyms:   []string{"barrow", "tomb"},
		Adjectives: []string{"massive", "stone"},
		Desc:       "stone barrow",
		Flags:      FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// Living Room
	TrophyCase = Object{
		In:         &LivingRoom,
		Synonyms:   []string{"case"},
		Adjectives: []string{"trophy"},
		Desc:       "trophy case",
		Flags:      FlgTrans | FlgCont | FlgNoDesc | FlgTryTake | FlgSearch,
		Capacity:   10000,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Rug = Object{
		In:         &LivingRoom,
		Synonyms:   []string{"rug", "carpet"},
		Adjectives: []string{"large", "oriental"},
		Desc:       "carpet",
		Flags:      FlgNoDesc | FlgTryTake,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	TrapDoor = Object{
		In:         &LivingRoom,
		Synonyms:   []string{"door", "trapdoor", "trap-door", "cover"},
		Adjectives: []string{"trap", "dusty"},
		Desc:       "trap door",
		Flags:      FlgDoor | FlgNoDesc | FlgInvis,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	WoodenDoor = Object{
		In:         &LivingRoom,
		Synonyms:   []string{"door", "lettering", "writing"},
		Adjectives: []string{"wooden", "gothic", "strange", "west"},
		Desc:       "wooden door",
		Flags:      FlgRead | FlgDoor | FlgNoDesc | FlgTrans,
		// Action set in FinalizeGameObjects to avoid init cycle
		Text:       "The engravings translate to \"This space intentionally left blank.\"",
	}
	Sword = Object{
		In:         &LivingRoom,
		Synonyms:   []string{"sword", "orcrist", "glamdring", "blade"},
		Adjectives: []string{"elvish", "old", "antique"},
		Desc:       "sword",
		Flags:      FlgTake | FlgWeapon | FlgTryTake,
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
		Flags:      FlgTake | FlgLight,
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
		Flags:      FlgNoDesc | FlgCont | FlgOpen | FlgSurf,
		Capacity:   50,
	}

	// Attic
	AtticTable = Object{
		In:       &Attic,
		Synonyms: []string{"table"},
		Desc:     "table",
		Flags:    FlgNoDesc | FlgCont | FlgOpen | FlgSurf,
		Capacity: 40,
	}
	Rope = Object{
		In:         &Attic,
		Synonyms:   []string{"rope", "hemp", "coil"},
		Adjectives: []string{"large"},
		Desc:       "rope",
		Flags:      FlgTake | FlgSacred | FlgTryTake,
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
		Flags:      FlgPerson | FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// Land of Living Dead
	Skull = Object{
		In:         &LandOfLivingDead,
		Synonyms:   []string{"skull", "head", "treasure"},
		Adjectives: []string{"crystal"},
		Desc:       "crystal skull",
		FirstDesc:  "Lying in one corner of the room is a beautifully carved crystal skull. It appears to be grinning at you rather nastily.",
		Flags:      FlgTake,
		Value:      10,
		TValue:     10,
	}

	// Shaft Room
	RaisedBasket = Object{
		In:         &ShaftRoom,
		Synonyms:   []string{"cage", "dumbwaiter", "basket"},
		Desc:       "basket",
		Flags:      FlgTrans | FlgTryTake | FlgCont | FlgOpen,
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
		Flags:      FlgTryTake,
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// Bat Room
	Bat = Object{
		In:         &BatRoom,
		Synonyms:   []string{"bat", "vampire"},
		Adjectives: []string{"vampire", "deranged"},
		Desc:       "bat",
		Flags:      FlgPerson | FlgTryTake,
		DescFcn:    BatDescFcn,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Jade = Object{
		In:         &BatRoom,
		Synonyms:   []string{"figurine", "treasure"},
		Adjectives: []string{"exquisite", "jade"},
		Desc:       "jade figurine",
		Flags:      FlgTake,
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
		Flags:      FlgTake,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Prayer = Object{
		In:         &NorthTemple,
		Synonyms:   []string{"prayer", "inscription"},
		Adjectives: []string{"ancient", "old"},
		Desc:       "prayer",
		Flags:      FlgRead | FlgSacred | FlgNoDesc,
		Text:       "The prayer is inscribed in an ancient script, rarely used today. It seems to be a philippic against small insects, absent-mindedness, and the picking up and dropping of small objects. The final verse consigns trespassers to the land of the dead. All evidence indicates that the beliefs of the ancient Zorkers were obscure.",
	}

	// South Temple
	Altar = Object{
		In:       &SouthTemple,
		Synonyms: []string{"altar"},
		Desc:     "altar",
		Flags:    FlgNoDesc | FlgSurf | FlgCont | FlgOpen,
		Capacity: 50,
	}
	Candles = Object{
		In:         &SouthTemple,
		Synonyms:   []string{"candles", "pair"},
		Adjectives: []string{"burning"},
		Desc:       "pair of candles",
		Flags:      FlgTake | FlgFlame | FlgOn | FlgLight,
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
		Flags:      FlgPerson | FlgOpen | FlgTryTake,
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
		Flags:      FlgNoDesc | FlgTurn | FlgTryTake,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Bubble = Object{
		In:         &DamRoom,
		Synonyms:   []string{"bubble"},
		Adjectives: []string{"small", "green", "plastic"},
		Desc:       "green bubble",
		Flags:      FlgNoDesc | FlgTryTake,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Dam = Object{
		In:         &DamRoom,
		Synonyms:   []string{"dam", "gate", "gates", "fcd#3"},
		Desc:       "dam",
		Flags:      FlgNoDesc | FlgTryTake,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	ControlPanel = Object{
		In:         &DamRoom,
		Synonyms:   []string{"panel"},
		Adjectives: []string{"control"},
		Desc:       "control panel",
		Flags:      FlgNoDesc,
	}

	// Dam Lobby
	Match = Object{
		In:         &DamLobby,
		Synonyms:   []string{"match", "matches", "matchbook"},
		Adjectives: []string{"match"},
		Desc:       "matchbook",
		Flags:      FlgRead | FlgTake,
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
		Flags:      FlgRead | FlgTake | FlgBurn,
		FirstDesc:  "Some guidebooks entitled \"Flood Control Dam #3\" are on the reception desk.",
		Text:       "\"\tFlood Control Dam #3\n\nFCD#3 was constructed in year 783 of the Great Underground Empire to harness the mighty Frigid River. This work was supported by a grant of 37 million zorkmids from your omnipotent local tyrant Lord Dimwit Flathead the Excessive. This impressive structure is composed of 370,000 cubic feet of concrete, is 256 feet tall at the center, and 193 feet wide at the top. The lake created behind the dam has a volume of 1.7 billion cubic feet, an area of 12 million square feet, and a shore line of 36 thousand feet.\n\nThe construction of FCD#3 took 112 days from ground breaking to the dedication. It required a work force of 384 slaves, 34 slave drivers, 12 engineers, 2 turtle doves, and a partridge in a pear tree. The work was managed by a command team composed of 2345 bureaucrats, 2347 secretaries (at least two of whom could type), 12,256 paper shufflers, 52,469 rubber stampers, 245,193 red tape processors, and nearly one million dead trees.\n\nWe will now point out some of the more interesting features of FCD#3 as we conduct you on a guided tour of the facilities:\n\n        1) You start your tour here in the Dam Lobby. You will notice on your right that....\"",
	}

	// Dam Base
	InflatableBoat = Object{
		In:         &DamBase,
		Synonyms:   []string{"boat", "pile", "plastic", "valve"},
		Adjectives: []string{"plastic", "inflatable"},
		Desc:       "pile of plastic",
		Flags:      FlgTake | FlgBurn,
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
		Flags:      FlgCont | FlgOpen | FlgTryTake | FlgSacred,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	YellowButton = Object{
		In:         &MaintenanceRoom,
		Synonyms:   []string{"button", "switch"},
		Adjectives: []string{"yellow"},
		Desc:       "yellow button",
		Flags:      FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	BrownButton = Object{
		In:         &MaintenanceRoom,
		Synonyms:   []string{"button", "switch"},
		Adjectives: []string{"brown"},
		Desc:       "brown button",
		Flags:      FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	RedButton = Object{
		In:         &MaintenanceRoom,
		Synonyms:   []string{"button", "switch"},
		Adjectives: []string{"red"},
		Desc:       "red button",
		Flags:      FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	BlueButton = Object{
		In:         &MaintenanceRoom,
		Synonyms:   []string{"button", "switch"},
		Adjectives: []string{"blue"},
		Desc:       "blue button",
		Flags:      FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Screwdriver = Object{
		In:         &MaintenanceRoom,
		Synonyms:   []string{"screwdriver", "tool", "tools", "driver"},
		Adjectives: []string{"screw"},
		Desc:       "screwdriver",
		Flags:      FlgTake | FlgTool,
	}
	Wrench = Object{
		In:       &MaintenanceRoom,
		Synonyms: []string{"wrench", "tool", "tools"},
		Desc:     "wrench",
		Flags:    FlgTake | FlgTool,
		Size:     10,
	}
	Tube = Object{
		In:       &MaintenanceRoom,
		Synonyms: []string{"tube", "tooth", "paste"},
		Desc:     "tube",
		Flags:    FlgTake | FlgCont | FlgRead,
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
		Flags:    FlgNoDesc | FlgInvis,
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// Machine Room
	Machine = Object{
		In:       &MachineRoom,
		Synonyms: []string{"machine", "pdp10", "dryer", "lid"},
		Desc:     "machine",
		Flags:    FlgCont | FlgNoDesc | FlgTryTake,
		Capacity: 50,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	MachineSwitch = Object{
		In:       &MachineRoom,
		Synonyms: []string{"switch"},
		Desc:     "switch",
		Flags:    FlgNoDesc | FlgTurn,
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// Cyclops Room
	Cyclops = Object{
		In:         &CyclopsRoom,
		Synonyms:   []string{"cyclops", "monster", "eye"},
		Adjectives: []string{"hungry", "giant"},
		Desc:       "cyclops",
		Flags:      FlgPerson | FlgNoDesc | FlgTryTake,
		// Action set in FinalizeGameObjects to avoid init cycle
		Strength:   10000,
	}

	// Treasure Room
	Chalice = Object{
		In:         &TreasureRoom,
		Synonyms:   []string{"chalice", "cup", "silver", "treasure"},
		Adjectives: []string{"silver", "engravings"},
		Desc:       "chalice",
		Flags:      FlgTake | FlgTryTake | FlgCont,
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
		Flags:      FlgTake | FlgBurn,
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
		Flags:      FlgRead | FlgTake,
		FirstDesc:  "Loosely attached to a wall is a small piece of paper.",
		Text:       "Congratulations!\n\nYou are the privileged owner of ZORK I: The Great Underground Empire, a self-contained and self-maintaining universe. If used and maintained in accordance with normal operating practices for small universes, ZORK will provide many months of trouble-free operation.",
	}

	// Grating Clearing
	Leaves = Object{
		In:       &GratingClearing,
		Synonyms: []string{"leaves", "leaf", "pile"},
		Desc:     "pile of leaves",
		Flags:    FlgTake | FlgBurn | FlgTryTake,
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
		Flags:      FlgTake | FlgBurn | FlgCont | FlgOpen | FlgSearch,
		FirstDesc:  "Beside you on the branch is a small bird's nest.",
		Capacity:   20,
	}

	// Sandy Cave
	Sand = Object{
		In:       &SandyCave,
		Synonyms: []string{"sand"},
		Desc:     "sand",
		Flags:    FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Scarab = Object{
		In:         &SandyCave,
		Synonyms:   []string{"scarab", "bug", "beetle", "treasure"},
		Adjectives: []string{"beautiful", "carved", "jeweled"},
		Desc:       "beautiful jeweled scarab",
		Flags:      FlgTake | FlgInvis,
		Size:       8,
		Value:      5,
		TValue:     5,
	}

	// Sandy Beach
	Shovel = Object{
		In:       &SandyBeach,
		Synonyms: []string{"shovel", "tool", "tools"},
		Desc:     "shovel",
		Flags:    FlgTake | FlgTool,
		Size:     15,
	}

	// Egypt Room
	Coffin = Object{
		In:         &EgyptRoom,
		Synonyms:   []string{"coffin", "casket", "treasure"},
		Adjectives: []string{"solid", "gold"},
		Desc:       "gold coffin",
		Flags:      FlgTake | FlgCont | FlgSacred | FlgSearch,
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
		Flags:      FlgPerson | FlgInvis | FlgCont | FlgOpen | FlgTryTake,
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
		Flags:      FlgTake | FlgInvis,
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
		Flags:      FlgTake | FlgTool,
	}

	// Atlantis Room
	Trident = Object{
		In:         &AtlantisRoom,
		Synonyms:   []string{"trident", "fork", "treasure"},
		Adjectives: []string{"poseidon", "own", "crystal"},
		Desc:       "crystal trident",
		Flags:      FlgTake,
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
		Flags:      FlgTryTake | FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Mirror2 = Object{
		In:         &MirrorRoom2,
		Synonyms:   []string{"reflection", "mirror", "enormous"},
		Desc:       "mirror",
		Flags:      FlgTryTake | FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// Dome Room
	Railing = Object{
		In:         &DomeRoom,
		Synonyms:   []string{"railing", "rail"},
		Adjectives: []string{"wooden"},
		Desc:       "wooden railing",
		Flags:      FlgNoDesc,
	}

	// Torch Room
	Pedestal = Object{
		In:         &TorchRoom,
		Synonyms:   []string{"pedestal"},
		Adjectives: []string{"white", "marble"},
		Desc:       "pedestal",
		Flags:      FlgNoDesc | FlgCont | FlgOpen | FlgSurf,
		// Action set in FinalizeGameObjects to avoid init cycle
		Capacity:   30,
	}

	// Engravings Cave
	Engravings = Object{
		In:       &EngravingsCave,
		Synonyms: []string{"wall", "engravings", "inscription"},
		Adjectives: []string{"old", "ancient"},
		Desc:     "wall with engravings",
		Flags:    FlgRead | FlgSacred,
		LongDesc: "There are old engravings on the walls here.",
		Text:     "The engravings were incised in the living rock of the cave wall by an unknown hand. They depict, in symbolic form, the beliefs of the ancient Zorkers. Skillfully interwoven with the bas reliefs are excerpts illustrating the major religious tenets of that time. Unfortunately, a later age seems to have considered them blasphemous and just as skillfully excised them.",
	}

	// Loud Room
	Bar = Object{
		In:         &LoudRoom,
		Synonyms:   []string{"bar", "platinum", "treasure"},
		Adjectives: []string{"platinum", "large"},
		Desc:       "platinum bar",
		Flags:      FlgTake | FlgSacred,
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
		Flags:      FlgTake | FlgInvis,
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
		Flags:      FlgTake | FlgCont,
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
		Flags:      FlgTake,
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
		Flags:      FlgTake | FlgBurn,
		Size:       20,
	}

	// Timber Room
	Timbers = Object{
		In:         &TimberRoom,
		Synonyms:   []string{"timbers", "pile"},
		Adjectives: []string{"wooden", "broken"},
		Desc:       "broken timber",
		Flags:      FlgTake,
		Size:       50,
	}

	// Maze 5
	Bones = Object{
		In:         &Maze5,
		Synonyms:   []string{"bones", "skeleton", "body"},
		Desc:       "skeleton",
		Flags:      FlgTryTake | FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	BurnedOutLantern = Object{
		In:         &Maze5,
		Synonyms:   []string{"lantern", "lamp"},
		Adjectives: []string{"rusty", "burned", "dead", "useless"},
		Desc:       "burned-out lantern",
		Flags:      FlgTake,
		FirstDesc:  "The deceased adventurer's useless lantern is here.",
		Size:       20,
	}
	BagOfCoins = Object{
		In:         &Maze5,
		Synonyms:   []string{"bag", "coins", "treasure"},
		Adjectives: []string{"old", "leather"},
		Desc:       "leather bag of coins",
		Flags:      FlgTake,
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
		Flags:      FlgTake | FlgTryTake | FlgWeapon | FlgTool,
		// Action set in FinalizeGameObjects to avoid init cycle
		FirstDesc:  "Beside the skeleton is a rusty knife.",
		Size:       20,
	}
	Keys = Object{
		In:         &Maze5,
		Synonyms:   []string{"key"},
		Adjectives: []string{"skeleton"},
		Desc:       "skeleton key",
		Flags:      FlgTake | FlgTool,
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
		Flags:      FlgInvis | FlgRead | FlgTake,
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
		Flags:      FlgRead | FlgTake | FlgBurn,
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
		Flags:      FlgTake | FlgTrans | FlgCont,
		// Action set in FinalizeGameObjects to avoid init cycle
		FirstDesc:  "A bottle is sitting on the table.",
		Capacity:   4,
	}
	SandwichBag = Object{
		In:         &KitchenTable,
		Synonyms:   []string{"bag", "sack"},
		Adjectives: []string{"brown", "elongated", "smelly"},
		Desc:       "brown sack",
		Flags:      FlgTake | FlgCont | FlgBurn,
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
		Flags:      FlgTake | FlgWeapon | FlgTryTake,
		FirstDesc:  "On a table is a nasty-looking knife.",
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// In Bottle
	Water = Object{
		In:         &Bottle,
		Synonyms:   []string{"water", "quantity", "liquid", "h2o"},
		Desc:       "quantity of water",
		Flags:      FlgTryTake | FlgTake | FlgDrink,
		// Action set in FinalizeGameObjects to avoid init cycle
		Size:       4,
	}

	// In Sandwich Bag
	Lunch = Object{
		In:         &SandwichBag,
		Synonyms:   []string{"food", "sandwich", "lunch", "dinner"},
		Adjectives: []string{"hot", "pepper"},
		Desc:       "lunch",
		Flags:      FlgTake | FlgFood,
		LongDesc:   "A hot pepper sandwich is here.",
	}
	Garlic = Object{
		In:       &SandwichBag,
		Synonyms: []string{"garlic", "clove"},
		Desc:     "clove of garlic",
		Flags:    FlgTake | FlgFood,
		// Action set in FinalizeGameObjects to avoid init cycle
		Size:     4,
	}

	// In Altar
	Book = Object{
		In:         &Altar,
		Synonyms:   []string{"book", "prayer", "page", "books"},
		Adjectives: []string{"large", "black"},
		Desc:       "black book",
		Flags:      FlgRead | FlgTake | FlgCont | FlgBurn | FlgTurn,
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
		Flags:      FlgTake | FlgWeapon,
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
		Flags:      FlgTake | FlgCont | FlgSearch,
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
		Flags:      FlgTake | FlgSearch,
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
		Flags:      FlgTake | FlgTool,
		Size:       6,
		// Action set in FinalizeGameObjects to avoid init cycle
	}

	// In Troll
	Axe = Object{
		In:         &Troll,
		Synonyms:   []string{"axe", "ax"},
		Adjectives: []string{"bloody"},
		Desc:       "bloody axe",
		Flags:      FlgWeapon | FlgTryTake | FlgTake | FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
		Size:       25,
	}

	// In Thief
	LargeBag = Object{
		In:         &Thief,
		Synonyms:   []string{"bag"},
		Adjectives: []string{"large", "thiefs"},
		Desc:       "large bag",
		Flags:      FlgTryTake | FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
	}
	Stiletto = Object{
		In:         &Thief,
		Synonyms:   []string{"stiletto"},
		Adjectives: []string{"vicious"},
		Desc:       "stiletto",
		Flags:      FlgWeapon | FlgTryTake | FlgTake | FlgNoDesc,
		// Action set in FinalizeGameObjects to avoid init cycle
		Size:       10,
	}

	// In Pedestal
	Torch = Object{
		In:         &Pedestal,
		Synonyms:   []string{"torch", "ivory", "treasure"},
		Adjectives: []string{"flaming", "ivory"},
		Desc:       "torch",
		Flags:      FlgTake | FlgFlame | FlgOn | FlgLight,
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
		Flags:      FlgRead | FlgTake | FlgBurn,
		Size:       2,
		Text:       "  !!!!FROBOZZ MAGIC BOAT COMPANY!!!!\n\nHello, Sailor!\n\nInstructions for use:\n\n   To get into a body of water, say \"Launch\".\n   To get to shore, say \"Land\" or the direction in which you want to maneuver the boat.\n\nWarranty:\n\n  This boat is guaranteed against all defects for a period of 76 milliseconds from date of purchase or until first used, whichever comes first.\n\nWarning:\n   This boat is made of thin plastic.\n   Good Luck!",
	}

	// In Buoy
	Emerald = Object{
		In:         &Buoy,
		Synonyms:   []string{"emerald", "treasure"},
		Adjectives: []string{"large"},
		Desc:       "large emerald",
		Flags:      FlgTake,
		Value:      5,
		TValue:     10,
	}

	// In Broken Egg (unplaced)
	BrokenCanary = Object{
		In:         &BrokenEgg,
		Synonyms:   []string{"canary", "treasure"},
		Adjectives: []string{"broken", "clockwork", "gold", "golden"},
		Desc:       "broken clockwork canary",
		Flags:      FlgTake,
		// Action set in FinalizeGameObjects to avoid init cycle
		TValue:     1,
		FirstDesc:  "There is a golden clockwork canary nestled in the egg. It seems to have recently had a bad experience. The mountings for its jewel-like eyes are empty, and its silver beak is crumpled. Through a cracked crystal window below its left wing you can see the remains of intricate machinery. It is not clear what result winding it would have, as the mainspring seems sprung.",
	}

)
