package zork

type LocFlag int

const (
	LocHeld LocFlag = iota
	LocCarried
	LocInRoom
	LocOnGrnd
	LocTake
	LocMany
	LocHave
)

func (lf LocFlag) In(flgs LocFlags) bool {
	for _, flg := range flgs {
		if flg == lf {
			return true
		}
	}
	return false
}

type LocFlags []LocFlag

func (lfs *LocFlags) All() LocFlags {
	*lfs = LocFlags{LocHeld, LocCarried, LocInRoom, LocOnGrnd, LocTake, LocMany, LocHave}
	return *lfs
}

func (lfs LocFlags) HasAll() bool {
	return len(lfs) == 7
}

type VrbAction func(ActArg) bool

type ObjProp struct {
	ObjFlags []Flag
	LocFlags LocFlags
	HasObj   bool
}

type Syntx struct {
	NormVerb  string
	Verb      string
	VrbPrep   string
	Obj1      ObjProp
	ObjPrep   string
	Obj2      ObjProp
	Action    VrbAction
	PreAction VrbAction
}

func (s *Syntx) NumObjects() int {
	if !s.Obj1.HasObj {
		return 0
	}
	if !s.Obj2.HasObj {
		return 1
	}
	return 2
}

func (s *Syntx) IsVrbPrep(prep string) bool {
	return s.VrbPrep == prep
}

func (s *Syntx) IsObjPrep(prep string) bool {
	return s.ObjPrep == prep
}

func (s *Syntx) GetActionVerb() string {
	av := s.Verb
	if len(s.VrbPrep) > 0 {
		av += " " + s.VrbPrep
	}
	return av
}

func (s *Syntx) GetNormVerb() string {
	if vrb, ok := NormVerbs[s.GetActionVerb()]; ok {
		return vrb
	}
	return s.GetActionVerb()
}

var (
	BuzzWords = []string{
		"again", "g", "oops", "a", "an", "the",
		"is", "and", "of", "then", "all", "one",
		"but", "except", ".", ",", "\"", "yes",
		"no", "y", "here",
	}
	Synonyms = map[string]string{
		"using":      "with",
		"through":    "with",
		"thru":       "with",
		"into":       "in",
		"inside":     "in",
		"onto":       "on",
		"underneath": "under",
		"beneath":    "under",
		"below":      "under",
		"n":          "north",
		"s":          "south",
		"e":          "east",
		"w":          "west",
		"d":          "down",
		"u":          "up",
		"nw":         "northwest",
		"ne":         "northeast",
		"sw":         "southwest",
		"se":         "southeast",
		"superbrief": "super",
		"i":          "inventory",
		"q":          "quit",
		"l":          "look",
		"stare":      "look",
		"gaze":       "look",
		"describe":   "examine",
		"what":       "examine",
		"whats":      "examine",
		"go":         "walk",
		"run":        "walk",
		"step":       "walk",
		"proceed":    "walk",
		"xyzzy":      "plugh",
		"awake":      "wake",
		"suprise":    "wake",
		"startle":    "wake",
		"reply":      "answer",
		"fight":      "attack",
		"hurt":       "attack",
		"injure":     "attack",
		"hit":        "attack",
		"murder":     "kill",
		"slay":       "kill",
		"dispatch":   "kill",
		"rap":        "knock",
		"clean":      "brush",
		"sit":        "climb",
		"get":        "take",
		"hold":       "take",
		"carry":      "take",
		"remove":     "take",
		"grab":       "take",
		"catch":      "take",
		"incinerate": "burn",
		"ignite":     "burn",
		"lose":       "chomp",
		"barf":       "chomp",
		"ford":       "cross",
		"slice":      "cut",
		"pierce":     "cut",
		"shit":       "curse",
		"fuck":       "curse",
		"damn":       "curse",
		"damage":     "destroy",
		"break":      "destroy",
		"block":      "destroy",
		"smash":      "destroy",
		"imbibe":     "drink",
		"swallow":    "drink",
		"consume":    "eat",
		"taste":      "eat",
		"bite":       "eat",
		"banish":     "exorcise",
		"cast":       "exorcise",
		"drive":      "exorcise",
		"begone":     "exorcise",
		"douse":      "extinguish",
		"where":      "find",
		"seek":       "find",
		"see":        "find",
		"pursue":     "follow",
		"chase":      "follow",
		"come":       "follow",
		"donate":     "give",
		"offer":      "give",
		"feed":       "give",
		"hi":         "hello",
		"chant":      "incant",
		"leap":       "jump",
		"dive":       "jump",
		"taunt":      "kick",
		"oil":        "lubricate",
		"grease":     "lubricate",
		"liquify":    "melt",
		"sigh":       "mumble",
		"ulysses":    "odysseus",
		"glue":       "plug",
		"patch":      "plug",
		"repair":     "plug",
		"fix":        "plug",
		"spill":      "pour",
		"tug":        "pull",
		"yank":       "pull",
		"press":      "push",
		"stuff":      "put",
		"insert":     "put",
		"place":      "put",
		"hide":       "put",
		"lift":       "raise",
		"molest":     "rape",
		"skim":       "read",
		"peal":       "ring",
		"touch":      "rub",
		"feel":       "rub",
		"pat":        "rub",
		"pet":        "rub",
		"hop":        "skip",
		"sniff":      "smell",
		"bathe":      "swim",
		"wade":       "swim",
		"thrust":     "swing",
		"ask":        "tell",
		"hurl":       "throw",
		"chuck":      "throw",
		"toss":       "throw",
		"fasten":     "tie",
		"secure":     "tie",
		"attach":     "tie",
		"temple":     "treasure",
		"set":        "turn",
		"flip":       "turn",
		"shut":       "turn",
		"free":       "untie",
		"release":    "untie",
		"unfasten":   "untie",
		"unattach":   "untie",
		"unhook":     "untie",
		"z":          "wait",
		"brandish":   "wave",
		"winnage":    "win",
		"scream":     "yell",
		"shout":      "yell",
	}
	Directions = []string{
		"north", "east", "west", "south", "northeast",
		"northwest", "southeast", "southwest", "up",
		"down", "in", "out", "land",
	}
	Commands = []Syntx{
		{
			Verb:   "verbose",
			Action: VVerbose,
		},
		{
			Verb:   "brief",
			Action: VBrief,
		},
		{
			Verb:   "super",
			Action: VSuperBrief,
		},
		{
			Verb:   "inventory",
			Action: VInventory,
		},
		{
			Verb:   "quit",
			Action: VQuit,
		},
		{
			Verb:   "restart",
			Action: VRestart,
		},
		{
			Verb:   "restore",
			Action: VRestore,
		},
		{
			Verb:   "save",
			Action: VSave,
		},
		{
			Verb:   "score",
			Action: VScore,
		},
		{
			Verb:   "script",
			Action: VScript,
		},
		{
			Verb:   "unscript",
			Action: VUnscript,
		},
		{
			Verb:   "version",
			Action: VVersion,
		},
		{
			Verb:   "verify",
			Action: VVerify,
		},
		{
			Verb:   "answer",
			Action: VAnswer,
		},
		{
			Verb:     "activate",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgLight}, LocFlags: LocFlags{LocHeld, LocCarried, LocOnGrnd, LocInRoom}},
			Action:   VLampOn,
			NormVerb: "lamp on",
		},
		{
			Verb:    "attack",
			Obj1:    ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgPerson}, LocFlags: LocFlags{LocInRoom, LocOnGrnd}},
			ObjPrep: "with",
			Obj2:    ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgWeapon}, LocFlags: LocFlags{LocHeld, LocCarried, LocHave}},
			Action:  VAttack,
		},
		{
			Verb:   "back",
			Action: VBack,
		},
		{
			Verb:   "blast",
			Action: VBlast,
		},
		{
			Verb:     "blow",
			VrbPrep:  "out",
			Obj1:     ObjProp{HasObj: true},
			Action:   VLampOff,
			NormVerb: "lamp off",
		},
		{
			Verb:     "blow",
			VrbPrep:  "up",
			Obj1:     ObjProp{HasObj: true},
			Action:   VBlast,
			NormVerb: "blast",
		},
		{
			Verb:     "blow",
			VrbPrep:  "in",
			Obj1:     ObjProp{HasObj: true},
			Action:   VBreathe,
			NormVerb: "breathe",
		},
		{
			Verb:     "blow",
			VrbPrep:  "up",
			Obj1:     ObjProp{HasObj: true},
			ObjPrep:  "with",
			Obj2:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgTool}, LocFlags: LocFlags{LocOnGrnd, LocInRoom, LocHeld, LocCarried}},
			Action:   VInflate,
			NormVerb: "inflate",
		},
		{
			Verb:      "board",
			Obj1:      ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgVeh}, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action:    VBoard,
			PreAction: PreBoard,
		},
		{
			Verb:   "brush",
			Obj1:   ObjProp{HasObj: true, LocFlags: LocFlags{LocOnGrnd, LocInRoom, LocCarried, LocHeld}},
			Action: VBrush,
		},
		{
			Verb:    "brush",
			Obj1:    ObjProp{HasObj: true, LocFlags: LocFlags{LocOnGrnd, LocInRoom, LocCarried, LocHeld}},
			ObjPrep: "with",
			Obj2:    ObjProp{HasObj: true},
			Action:  VBrush,
		},
		{
			Verb:   "bug",
			Action: VBug,
		},
		{
			Verb:      "burn",
			Obj1:      ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgBurn}, LocFlags: LocFlags{LocInRoom, LocOnGrnd, LocHeld, LocCarried}},
			ObjPrep:   "with",
			Obj2:      ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgFlame}, LocFlags: LocFlags{LocHeld, LocCarried, LocHave, LocInRoom, LocOnGrnd}},
			Action:    VBurn,
			PreAction: PreBurn,
		},
		{
			Verb:      "burn",
			VrbPrep:   "down",
			Obj1:      ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgBurn}, LocFlags: LocFlags{LocInRoom, LocOnGrnd, LocHeld, LocCarried}},
			ObjPrep:   "with",
			Obj2:      ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgFlame}, LocFlags: LocFlags{LocHeld, LocCarried, LocHave, LocInRoom, LocOnGrnd}},
			Action:    VBurn,
			PreAction: PreBurn,
			NormVerb:  "burn",
		},
		{
			Verb:   "chomp",
			Action: VChomp,
		},
		{
			Verb:    "climb",
			VrbPrep: "up",
			Obj1:    ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgKludge}},
			Action:  VClimbUp,
		},
		{
			Verb:    "climb",
			VrbPrep: "up",
			Obj1:    ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgClimb}, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action:  VClimbUp,
		},
		{
			Verb:    "climb",
			VrbPrep: "down",
			Obj1:    ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgKludge}},
			Action:  VClimbDown,
		},
		{
			Verb:    "climb",
			VrbPrep: "down",
			Obj1:    ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgClimb}, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action:  VClimbDown,
		},
		{
			Verb:     "climb",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgClimb}, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action:   VClimbFoo,
			NormVerb: "climb foo",
		},
		{
			Verb:      "climb",
			VrbPrep:   "in",
			Obj1:      ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgVeh}, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action:    VBoard,
			PreAction: PreBoard,
			NormVerb:  "board",
		},
		{
			Verb:    "climb",
			VrbPrep: "on",
			Obj1:    ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgVeh}, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action:  VClimbOn,
		},
		{
			Verb:   "close",
			Obj1:   ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgDoor}, LocFlags: LocFlags{LocOnGrnd, LocInRoom, LocHeld, LocCarried}},
			Action: VClose,
		},
		{
			Verb:   "command",
			Obj1:   ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgPerson}},
			Action: VCommand,
		},
		{
			Verb:   "count",
			Obj1:   ObjProp{HasObj: true},
			Action: VCount,
		},
		{
			Verb:   "cross",
			Obj1:   ObjProp{HasObj: true},
			Action: VCross,
		},
		{
			Verb:   "curse",
			Action: VCurses,
		},
		{
			Verb:   "curse",
			Obj1:   ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgPerson}},
			Action: VCurses,
		},
		{
			Verb:    "cut",
			Obj1:    ObjProp{HasObj: true},
			ObjPrep: "with",
			Obj2:    ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgWeapon}, LocFlags: LocFlags{LocHeld, LocCarried}},
			Action:  VCut,
		},
		{
			Verb:   "deflate",
			Obj1:   ObjProp{HasObj: true},
			Action: VDeflate,
		},
		{
			Verb:      "destroy",
			Obj1:      ObjProp{HasObj: true, LocFlags: LocFlags{LocOnGrnd, LocInRoom, LocHeld, LocCarried}},
			ObjPrep:   "with",
			Obj2:      ObjProp{HasObj: true, LocFlags: LocFlags{LocHeld, LocCarried, LocTake}},
			Action:    VMung,
			PreAction: PreMung,
			NormVerb:  "mung",
		},
		{
			Verb:      "destroy",
			VrbPrep:   "down",
			Obj1:      ObjProp{HasObj: true, LocFlags: LocFlags{LocOnGrnd, LocInRoom, LocHeld, LocCarried}},
			ObjPrep:   "with",
			Obj2:      ObjProp{HasObj: true, LocFlags: LocFlags{LocHeld, LocCarried, LocTake}},
			Action:    VMung,
			PreAction: PreMung,
			NormVerb:  "mung",
		},
		{
			Verb:     "destroy",
			VrbPrep:  "in",
			Obj1:     ObjProp{HasObj: true, LocFlags: LocFlags{LocOnGrnd, LocInRoom, LocHeld, LocCarried}},
			Action:   VOpen,
			NormVerb: "open",
		},
		{
			Verb:     "dig",
			VrbPrep:  "in",
			Obj1:     ObjProp{HasObj: true, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action:   VDig,
			NormVerb: "dig",
		},
		{
			Verb:     "dig",
			VrbPrep:  "in",
			Obj1:     ObjProp{HasObj: true, LocFlags: LocFlags{LocInRoom, LocOnGrnd}},
			ObjPrep:  "with",
			Obj2:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgTool}, LocFlags: LocFlags{LocHeld, LocCarried, LocHave}},
			Action:   VDig,
			NormVerb: "dig",
		},
		{
			Verb:    "dig",
			Obj1:    ObjProp{HasObj: true, LocFlags: LocFlags{LocInRoom, LocOnGrnd}},
			ObjPrep: "with",
			Obj2:    ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgTool}, LocFlags: LocFlags{LocHeld, LocCarried, LocHave}},
			Action:  VDig,
		},
		{
			Verb:   "disembark",
			Obj1:   ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgVeh}, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action: VDisembark,
		},
		{
			Verb:   "disenchant",
			Obj1:   ObjProp{HasObj: true},
			Action: VDisenchant,
		},
		{
			Verb:   "drink",
			Obj1:   ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgDrink}, LocFlags: LocFlags{LocHeld, LocCarried, LocOnGrnd, LocInRoom}},
			Action: VDrink,
		},
		{
			Verb:    "drink",
			VrbPrep: "from",
			Obj1:    ObjProp{HasObj: true, LocFlags: LocFlags{LocHeld, LocCarried}},
			Action:  VDrinkFrom,
		},
		{
			Verb:      "drop",
			Obj1:      ObjProp{HasObj: true, LocFlags: LocFlags{LocHeld, LocHave, LocMany}},
			Action:    VDrop,
			PreAction: PreDrop,
		},
		{
			Verb:      "eat",
			Obj1:      ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgFood}, LocFlags: LocFlags{LocHeld, LocCarried, LocOnGrnd, LocInRoom, LocTake}},
			Action:    VDrop,
			PreAction: PreDrop,
		},
		{
			Verb:   "echo",
			Action: VEcho,
		},
		{
			Verb:   "enchant",
			Obj1:   ObjProp{HasObj: true, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action: VEnchant,
		},
		{
			Verb:   "enter",
			Action: VEnter,
		},
		{
			Verb:   "exit",
			Action: VExit,
		},
		{
			Verb:   "exit",
			Obj1:   ObjProp{HasObj: true},
			Action: VExit,
		},
		{
			Verb:   "examine",
			Obj1:   ObjProp{HasObj: true, LocFlags: LocFlags{LocMany}},
			Action: VExamine,
		},
		{
			Verb:     "examine",
			VrbPrep:  "in",
			Obj1:     ObjProp{HasObj: true, LocFlags: LocFlags{LocHeld, LocCarried, LocInRoom, LocOnGrnd, LocMany}},
			Action:   VLookInside,
			NormVerb: "look inside",
		},
		{
			Verb:     "examine",
			VrbPrep:  "on",
			Obj1:     ObjProp{HasObj: true, LocFlags: LocFlags{LocHeld, LocCarried, LocInRoom, LocOnGrnd, LocMany}},
			Action:   VLookInside,
			NormVerb: "look inside",
		},
		{
			Verb:     "extinguish",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgOn}, LocFlags: LocFlags{LocHeld, LocCarried, LocOnGrnd, LocInRoom, LocTake, LocHave}},
			Action:   VLampOff,
			NormVerb: "lamp off",
		},
		{
			Verb:      "fill",
			Obj1:      ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgCont}, LocFlags: LocFlags{LocHeld, LocCarried, LocOnGrnd, LocInRoom}},
			ObjPrep:   "with",
			Obj2:      ObjProp{HasObj: true},
			Action:    VFill,
			PreAction: PreFill,
		},
		{
			Verb:      "fill",
			Obj1:      ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgCont}, LocFlags: LocFlags{LocHeld, LocCarried, LocOnGrnd, LocInRoom}},
			Action:    VFill,
			PreAction: PreFill,
		},
		{
			Verb:   "find",
			Obj1:   ObjProp{HasObj: true},
			Action: VFind,
		},
		{
			Verb:   "follow",
			Action: VFollow,
		},
		{
			Verb:   "follow",
			Obj1:   ObjProp{HasObj: true},
			Action: VFollow,
		},
		{
			Verb:   "frobozz",
			Action: VFrobozz,
		},
		{
			Verb:      "give",
			Obj1:      ObjProp{HasObj: true, LocFlags: LocFlags{LocMany, LocHeld, LocHave}},
			ObjPrep:   "to",
			Obj2:      ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgPerson}, LocFlags: LocFlags{LocOnGrnd}},
			Action:    VGive,
			PreAction: PreGive,
		},
		{
			Verb:   "hatch",
			Obj1:   ObjProp{HasObj: true},
			Action: VHatch,
		},
		{
			Verb:   "hello",
			Action: VHello,
		},
		{
			Verb:   "hello",
			Obj1:   ObjProp{HasObj: true},
			Action: VHello,
		},
		{
			Verb:   "incant",
			Action: VIncant,
		},
		{
			Verb:    "inflate",
			Obj1:    ObjProp{HasObj: true},
			ObjPrep: "with",
			Obj2:    ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgTool}, LocFlags: LocFlags{LocOnGrnd, LocInRoom, LocHeld, LocCarried}},
			Action:  VInflate,
		},
		{
			Verb:   "jump",
			Action: VLeap,
		},
		{
			Verb:     "jump",
			VrbPrep:  "over",
			Obj1:     ObjProp{HasObj: true},
			Action:   VLeap,
			NormVerb: "leap",
		},
		{
			Verb:     "jump",
			VrbPrep:  "across",
			Obj1:     ObjProp{HasObj: true},
			Action:   VLeap,
			NormVerb: "leap",
		},
		{
			Verb:     "jump",
			VrbPrep:  "in",
			Obj1:     ObjProp{HasObj: true},
			Action:   VLeap,
			NormVerb: "leap",
		},
		{
			Verb:     "jump",
			VrbPrep:  "from",
			Obj1:     ObjProp{HasObj: true},
			Action:   VLeap,
			NormVerb: "leap",
		},
		{
			Verb:     "jump",
			VrbPrep:  "off",
			Obj1:     ObjProp{HasObj: true},
			Action:   VLeap,
			NormVerb: "leap",
		},
		{
			Verb:   "kick",
			Obj1:   ObjProp{HasObj: true},
			Action: VKick,
		},
		{
			Verb:     "kill",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgPerson}, LocFlags: LocFlags{LocInRoom, LocOnGrnd}},
			ObjPrep:  "with",
			Obj2:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgWeapon}, LocFlags: LocFlags{LocHeld, LocCarried, LocHave}},
			Action:   VAttack,
			NormVerb: "attack",
		},
		{
			Verb:   "leave",
			Action: VLeave,
		},
		{
			Verb:      "leave",
			Obj1:      ObjProp{HasObj: true},
			Action:    VDrop,
			PreAction: PreDrop,
			NormVerb:  "drop",
		},
		{
			Verb:     "stab",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgPerson}, LocFlags: LocFlags{LocInRoom, LocOnGrnd}},
			ObjPrep:  "with",
			Obj2:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgWeapon}, LocFlags: LocFlags{LocHeld, LocCarried, LocHave}},
			Action:   VAttack,
			NormVerb: "attack",
		},
		{
			Verb:   "kiss",
			Obj1:   ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgPerson}, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action: VKiss,
		},
		{
			Verb:    "knock",
			VrbPrep: "at",
			Obj1:    ObjProp{HasObj: true},
			Action:  VKnock,
		},
		{
			Verb:    "knock",
			VrbPrep: "on",
			Obj1:    ObjProp{HasObj: true},
			Action:  VKnock,
		},
		{
			Verb:     "knock",
			VrbPrep:  "down",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgPerson}, LocFlags: LocFlags{LocInRoom, LocOnGrnd}},
			Action:   VAttack,
			NormVerb: "attack",
		},
		{
			Verb:   "launch",
			Obj1:   ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgVeh}},
			Action: VLaunch,
		},
		{
			Verb:    "lean",
			VrbPrep: "on",
			Obj1:    ObjProp{HasObj: true, LocFlags: LocFlags{LocHeld, LocHave}},
			Action:  VLeanOn,
		},
		{
			Verb:     "light",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgLight}, LocFlags: LocFlags{LocHeld, LocCarried, LocOnGrnd, LocInRoom, LocTake, LocHave}},
			Action:   VLampOn,
			NormVerb: "lamp on",
		},
		{
			Verb:      "light",
			Obj1:      ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgBurn}, LocFlags: LocFlags{LocInRoom, LocOnGrnd, LocHeld, LocCarried}},
			ObjPrep:   "with",
			Obj2:      ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgFlame}, LocFlags: LocFlags{LocHeld, LocCarried, LocHave, LocInRoom, LocOnGrnd}},
			Action:    VBurn,
			PreAction: PreBurn,
			NormVerb:  "burn",
		},
		{
			Verb:     "listen",
			VrbPrep:  "to",
			Obj1:     ObjProp{HasObj: true},
			Action:   VListen,
			NormVerb: "listen",
		},
		{
			Verb:     "listen",
			VrbPrep:  "for",
			Obj1:     ObjProp{HasObj: true},
			Action:   VListen,
			NormVerb: "listen",
		},
		{
			Verb:    "lock",
			Obj1:    ObjProp{HasObj: true, LocFlags: LocFlags{LocInRoom, LocOnGrnd}},
			ObjPrep: "with",
			Obj2:    ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgTool}, LocFlags: LocFlags{LocHeld, LocCarried, LocInRoom, LocOnGrnd, LocTake}},
			Action:  VLock,
		},
		{
			Verb:   "look",
			Action: VLook,
		},
		{
			Verb:     "look",
			VrbPrep:  "around",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgKludge}},
			Action:   VLook,
			NormVerb: "look",
		},
		{
			Verb:     "look",
			VrbPrep:  "up",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgKludge}},
			Action:   VLook,
			NormVerb: "look",
		},
		{
			Verb:     "look",
			VrbPrep:  "down",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgKludge}},
			Action:   VLook,
			NormVerb: "look",
		},
		{
			Verb:     "look",
			VrbPrep:  "at",
			Obj1:     ObjProp{HasObj: true, LocFlags: LocFlags{LocHeld, LocCarried, LocInRoom, LocOnGrnd, LocMany}},
			Action:   VExamine,
			NormVerb: "examine",
		},
		{
			Verb:     "look",
			VrbPrep:  "for",
			Obj1:     ObjProp{HasObj: true},
			Action:   VFind,
			NormVerb: "find",
		},
		{
			Verb:    "look",
			VrbPrep: "on",
			Obj1:    ObjProp{HasObj: true},
			Action:  VLookOn,
		},
		{
			Verb:     "look",
			VrbPrep:  "with",
			Obj1:     ObjProp{HasObj: true, LocFlags: LocFlags{LocHeld, LocCarried, LocInRoom, LocOnGrnd, LocMany}},
			Action:   VLookInside,
			NormVerb: "look inside",
		},
		{
			Verb:    "look",
			VrbPrep: "under",
			Obj1:    ObjProp{HasObj: true},
			Action:  VLookUnder,
		},
		{
			Verb:    "look",
			VrbPrep: "behind",
			Obj1:    ObjProp{HasObj: true},
			Action:  VLookBehind,
		},
		{
			Verb:     "look",
			VrbPrep:  "in",
			Obj1:     ObjProp{HasObj: true, LocFlags: LocFlags{LocHeld, LocCarried, LocInRoom, LocOnGrnd, LocMany}},
			Action:   VLookInside,
			NormVerb: "look inside",
		},
		{
			Verb:   "lower",
			Obj1:   ObjProp{HasObj: true},
			Action: VLower,
		},
		{
			Verb:    "lubricate",
			Obj1:    ObjProp{HasObj: true},
			ObjPrep: "with",
			Obj2:    ObjProp{HasObj: true, LocFlags: LocFlags{LocHeld, LocCarried}},
			Action:  VOil,
		},
		{
			Verb:   "make",
			Obj1:   ObjProp{HasObj: true},
			Action: VMake,
		},
		{
			Verb:    "melt",
			Obj1:    ObjProp{HasObj: true},
			ObjPrep: "with",
			Obj2:    ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgFlame}, LocFlags: LocFlags{LocHeld, LocCarried, LocOnGrnd, LocInRoom}},
			Action:  VMelt,
		},
		{
			Verb:      "move",
			Obj1:      ObjProp{HasObj: true, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action:    VMove,
			PreAction: PreMove,
		},
		{
			Verb:      "roll",
			VrbPrep:   "up",
			Obj1:      ObjProp{HasObj: true, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action:    VMove,
			PreAction: PreMove,
			NormVerb:  "move",
		},
		{
			Verb:      "roll",
			Obj1:      ObjProp{HasObj: true, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action:    VMove,
			PreAction: PreMove,
			NormVerb:  "move",
		},
		{
			Verb:   "mumble",
			Action: VMumble,
		},
		{
			Verb:   "odysseus",
			Action: VOdysseus,
		},
		{
			Verb:   "open",
			Obj1:   ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgDoor}, LocFlags: LocFlags{LocOnGrnd, LocInRoom, LocHeld, LocCarried}},
			Action: VOpen,
		},
		{
			Verb:     "open",
			VrbPrep:  "up",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgDoor}, LocFlags: LocFlags{LocOnGrnd, LocInRoom, LocHeld, LocCarried}},
			Action:   VOpen,
			NormVerb: "open",
		},
		{
			Verb:    "open",
			Obj1:    ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgDoor}, LocFlags: LocFlags{LocOnGrnd, LocInRoom, LocHeld, LocCarried}},
			ObjPrep: "with",
			Obj2:    ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgTool}, LocFlags: LocFlags{LocOnGrnd, LocInRoom, LocHeld, LocCarried, LocHave}},
			Action:  VOpen,
		},
		{
			Verb:   "plugh",
			Action: VAdvent,
		},
		{
			Verb:      "poke",
			Obj1:      ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgPerson}, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			ObjPrep:   "with",
			Obj2:      ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgWeapon}, LocFlags: LocFlags{LocHeld, LocCarried, LocHave}},
			Action:    VMung,
			PreAction: PreMung,
			NormVerb:  "mung",
		},
		{
			Verb:      "puncture",
			Obj1:      ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgPerson}, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			ObjPrep:   "with",
			Obj2:      ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgWeapon}, LocFlags: LocFlags{LocHeld, LocCarried, LocHave}},
			Action:    VMung,
			PreAction: PreMung,
			NormVerb:  "mung",
		},
		{
			Verb:      "pour",
			Obj1:      ObjProp{HasObj: true, LocFlags: LocFlags{LocHeld, LocCarried}},
			Action:    VDrop,
			PreAction: PreDrop,
			NormVerb:  "drop",
		},
		{
			Verb:      "pour",
			Obj1:      ObjProp{HasObj: true, LocFlags: LocFlags{LocHeld, LocCarried}},
			ObjPrep:   "in",
			Obj2:      ObjProp{HasObj: true},
			Action:    VDrop,
			PreAction: PreDrop,
			NormVerb:  "drop",
		},
		{
			Verb:      "pour",
			Obj1:      ObjProp{HasObj: true, LocFlags: LocFlags{LocHeld, LocCarried}},
			ObjPrep:   "from",
			Obj2:      ObjProp{HasObj: true},
			Action:    VDrop,
			PreAction: PreDrop,
			NormVerb:  "drop",
		},
		{
			Verb:      "pull",
			Obj1:      ObjProp{HasObj: true, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action:    VMove,
			PreAction: PreMove,
			NormVerb:  "move",
		},
		{
			Verb:      "pull",
			VrbPrep:   "on",
			Obj1:      ObjProp{HasObj: true, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action:    VMove,
			PreAction: PreMove,
			NormVerb:  "move",
		},
		{
			Verb:      "pull",
			VrbPrep:   "up",
			Obj1:      ObjProp{HasObj: true, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action:    VMove,
			PreAction: PreMove,
			NormVerb:  "move",
		},
		{
			Verb:      "put",
			VrbPrep:   "down",
			Obj1:      ObjProp{HasObj: true, LocFlags: LocFlags{LocHeld, LocMany}},
			Action:    VDrop,
			PreAction: PreDrop,
			NormVerb:  "drop",
		},
		{
			Verb:     "put",
			VrbPrep:  "out",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgOn}, LocFlags: LocFlags{LocHeld, LocCarried, LocOnGrnd, LocInRoom, LocTake, LocHave}},
			Action:   VLampOff,
			NormVerb: "lamp off",
		},
		{
			Verb:     "search",
			VrbPrep:  "for",
			Obj1:     ObjProp{HasObj: true},
			Action:   VFind,
			NormVerb: "find",
		},
		{
			Verb:     "strike",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgPerson}, LocFlags: LocFlags{LocInRoom, LocOnGrnd}},
			ObjPrep:  "with",
			Obj2:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgWeapon}, LocFlags: LocFlags{LocHeld, LocCarried, LocHave, LocInRoom, LocOnGrnd}},
			Action:   VAttack,
			NormVerb: "attack",
		},
		{
			Verb:     "turn",
			VrbPrep:  "on",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgLight}, LocFlags: LocFlags{LocHeld, LocCarried, LocOnGrnd, LocInRoom}},
			Action:   VLampOn,
			NormVerb: "lamp on",
		},
		{
			Verb:     "turn",
			VrbPrep:  "on",
			Obj1:     ObjProp{HasObj: true},
			ObjPrep:  "with",
			Obj2:     ObjProp{HasObj: true, LocFlags: LocFlags{LocHave}},
			Action:   VLampOn,
			NormVerb: "lamp on",
		},
		{
			Verb:     "turn",
			VrbPrep:  "off",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgOn}, LocFlags: LocFlags{LocHeld, LocCarried, LocOnGrnd, LocInRoom, LocTake, LocHave}},
			Action:   VLampOff,
			NormVerb: "lamp off",
		},
		{
			Verb:      "take",
			VrbPrep:   "in",
			Obj1:      ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgVeh}, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action:    VBoard,
			PreAction: PreBoard,
			NormVerb:  "board",
		},
		{
			Verb:     "take",
			VrbPrep:  "out",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgKludge}, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action:   VDisembark,
			NormVerb: "disembark",
		},
		{
			Verb:     "take",
			VrbPrep:  "on",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgVeh}, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action:   VClimbOn,
			NormVerb: "climb on",
		},
		{
			Verb:     "wake",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgPerson}, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action:   VAlarm,
			NormVerb: "alarm",
		},
		{
			Verb:     "wake",
			VrbPrep:  "up",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgPerson}, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action:   VAlarm,
			NormVerb: "alarm",
		},
		{
			Verb:     "walk",
			Action:   VWalkAround,
			NormVerb: "walk around",
		},
		{
			Verb:   "walk",
			Obj1:   ObjProp{HasObj: true},
			Action: VWalk,
		},
		{
			Verb:     "walk",
			VrbPrep:  "away",
			Obj1:     ObjProp{HasObj: true},
			Action:   VWalk,
			NormVerb: "walk",
		},
		{
			Verb:     "walk",
			VrbPrep:  "over",
			Obj1:     ObjProp{HasObj: true},
			Action:   VLeap,
			NormVerb: "leap",
		},
		{
			Verb:    "walk",
			VrbPrep: "to",
			Obj1:    ObjProp{HasObj: true},
			Action:  VWalkTo,
		},
		{
			Verb:    "walk",
			VrbPrep: "around",
			Obj1:    ObjProp{HasObj: true},
			Action:  VWalkAround,
		},
		{
			Verb:     "walk",
			VrbPrep:  "up",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgClimb}, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action:   VClimbUp,
			NormVerb: "climb up",
		},
		{
			Verb:     "walk",
			VrbPrep:  "down",
			Obj1:     ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgClimb}, LocFlags: LocFlags{LocOnGrnd, LocInRoom}},
			Action:   VClimbDown,
			NormVerb: "climb down",
		},
	}
	Actions    = make(map[string]VrbAction)
	PreActions = make(map[string]VrbAction)
	NormVerbs  = make(map[string]string)
)

func addToVocab(wrd string, typ WordTyp) {
	v, ok := Vocabulary[wrd]
	if !ok {
		Vocabulary[wrd] = WordItm{
			Norm:  wrd,
			Types: WordTypes{typ},
		}
	} else {
		Vocabulary[wrd] = WordItm{
			Norm:  wrd,
			Types: append(v.Types, typ),
		}
	}
}

func BuildVocabulary() {
	// Add buzz words
	for _, bw := range BuzzWords {
		addToVocab(bw, WordBuzz)
	}
	// Add verbs
	for _, cmd := range Commands {
		addToVocab(cmd.Verb, WordVerb)
		if len(cmd.VrbPrep) > 0 {
			addToVocab(cmd.VrbPrep, WordPrep)
		}
		if len(cmd.ObjPrep) > 0 {
			addToVocab(cmd.ObjPrep, WordPrep)
		}
		Actions[cmd.GetActionVerb()] = cmd.Action
		if cmd.PreAction != nil {
			PreActions[cmd.GetActionVerb()] = cmd.PreAction
		}
		NormVerbs[cmd.GetActionVerb()] = cmd.GetActionVerb()
		if len(cmd.NormVerb) > 0 {
			NormVerbs[cmd.GetActionVerb()] = cmd.NormVerb
		}
	}
	// Add directions
	for _, dir := range Directions {
		addToVocab(dir, WordDir)
	}
	// Add objects
	for _, obj := range Objects {
		if obj.Synonyms != nil {
			for _, s := range obj.Synonyms {
				addToVocab(s, WordObj)
			}
		}
		if obj.Adjectives != nil {
			for _, a := range obj.Adjectives {
				addToVocab(a, WordAdj)
			}
		}
	}
	// Add synonyms
	for key, val := range Synonyms {
		if _, ok := Vocabulary[key]; !ok {
			Vocabulary[key] = WordItm{
				Norm:  val,
				Types: nil,
			}
			if el, ok := Vocabulary[val]; ok {
				Vocabulary[key] = WordItm{
					Norm:  val,
					Types: append(WordTypes{}, el.Types...),
				}
			}
		}
	}
}
