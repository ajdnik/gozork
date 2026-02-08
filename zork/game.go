package zork

import "os"

var (
	ParserOk bool
	Script   bool
	Player   *Object
)

type PerfRet int

const (
	PerfNotHndld PerfRet = iota
	PerfHndld
	PerfFatal
)

// ErrQuit is a sentinel value used to signal game exit via panic/recover.
// This allows tests to catch quit without os.Exit terminating the process.
var ErrQuit = "game-quit"

func Restart() bool {
	// TODO: Implement restart
	return false
}

func Restore() bool {
	// TODO: Implement restore
	return false
}

func Quit() {
	panic(ErrQuit)
}

func Save() bool {
	// TODO: Implement save
	return false
}

func Verify() bool {
	// TODO: Implement disk verify
	return false
}

// ResetGameState resets all mutable global state so a fresh game can be started.
// This is essential for tests that run multiple games in one process.
func ResetGameState() {
	// Reset clock
	QueueItms = [30]QueueItm{}
	QueueInts = 30
	QueueDmns = 30
	ClockWait = false

	// Reset parser/game state
	ParserOk = false
	Dead = false
	Deaths = 0
	Score = 0
	BaseScore = 0
	Moves = 0
	Lit = false
	SuperBrief = false
	Verbose = false
	WonGame = false
	HelloSailor = 0
	IsSprayed = false
	Lucky = true
	FumbleNumber = 7
	FumbleProb = 8
	DescObj = nil

	// Reset dungeon flags
	TrollFlag = false
	CyclopsFlag = false
	MagicFlag = false
	LowTide = false
	DomeFlag = false
	EmptyHanded = false
	LLDFlag = false
	RainbowFlag = false
	DeflateFlag = false
	CoffinCure = false
	GrateRevealed = false
	KitchenWindowFlag = false
	CageTop = true
	RugMoved = false
	GrUnlock = false
	CycloWrath = 0
	MirrorMung = false
	GateFlag = false
	GatesOpen = false
	WaterLevel = 0
	MatchCount = 6
	EggSolve = false
	ThiefHere = false
	ThiefEngrossed = false
	LoudFlag = false
	SingSong = false
	BuoyFlag = true
	BeachDig = -1
	LightShaft = 13
	XB = false
	XC = false
	Deflate = false
	LampTableIdx = 0
	CandleTableIdx = 0

	// Reset carry limits
	LoadAllowed = 100
	LoadMax = 100

	// Reset input state
	InputExhausted = false

	// Reset parser internals
	DirObjPossibles = nil
	IndirObjPossibles = nil
	LexRes = nil
	Params = ParseProps{}
	Reserv = ReserveProps{}
	Again = AgainProps{}
	Oops = OopsProps{}
	ActVerb = ActionVerb{}
	DirObj = nil
	IndirObj = nil
	Winner = nil
	Here = nil
	DetectedSyntx = nil
	ParsedSyntx = ParseTbl{}
	OrphanedSyntx = ParseTbl{}

	// Restore object tree to initial state
	ResetObjectTree()
}

// InitGame sets up all game state for a fresh game. Call once before MainLoop.
func InitGame() {
	ResetGameState()
	FinalizeGameObjects()
	BuildObjectTree()
	BuildVocabulary()
	InitReader()

	Queue(IFight, -1).Run = true
	Queue(ISword, -1)
	Queue(IThief, -1).Run = true
	Queue(ICandles, 40)
	Queue(ILantern, 200)
	InflatedBoat.VehType = FlgNonLand
	Def1Res[1] = Def1[2]
	Def1Res[2] = Def1[4]
	Def2Res[2] = Def2B[2]
	Def2Res[3] = Def2B[4]
	Def3Res[1] = Def3A[2]
	Def3Res[3] = Def3B[2]
	Here = &WestOfHouse
	ThisIsIt(&Mailbox)
	Lit = true
	Winner = &Adventurer
	Player = Winner
	Winner.MoveTo(Here)
}

func Run() {
	defer func() {
		if r := recover(); r != nil {
			if r == ErrQuit {
				os.Exit(0)
			}
			panic(r) // re-panic for unexpected panics
		}
	}()

	InitGame()
	if !Here.Has(FlgTouch) {
		VVersion(ActUnk)
		NewLine()
	}
	VLook(ActUnk)
	MainLoop()
}

func MainLoop() {
	Params.Continue = NumUndef
	for {
		if InputExhausted {
			return
		}
		ParserOk = Parse()
		if InputExhausted {
			return
		}
		if !ParserOk {
			Params.Continue = NumUndef
			continue
		}
		if Params.ItObj != nil && IsAccessible(Params.ItObj) {
			found := false
			for idx, obj := range IndirObjPossibles {
				if obj == &It {
					IndirObjPossibles[idx] = Params.ItObj
					found = true
				}
			}
			if !found {
				for idx, obj := range DirObjPossibles {
					if obj == &It {
						DirObjPossibles[idx] = Params.ItObj
					}
				}
			}
		}
		numObj := 1
		var obj *Object
		isDir := true
		if len(DirObjPossibles) == 0 {
			numObj = 0
		} else if len(DirObjPossibles) > 1 {
			if len(IndirObjPossibles) == 0 {
				obj = nil
			} else {
				obj = IndirObjPossibles[0]
			}
			numObj = len(DirObjPossibles)
		} else if len(IndirObjPossibles) > 1 {
			isDir = false
			obj = DirObjPossibles[0]
			numObj = len(IndirObjPossibles)
		}
		if obj == nil && len(IndirObjPossibles) == 1 {
			obj = IndirObjPossibles[0]
		}
		var res PerfRet
		if ActVerb.Norm == "walk" && len(Params.WalkDir) != 0 {
			res = Perform(ActVerb, DirObj, nil)
		} else if numObj == 0 {
			if DetectedSyntx.NumObjects() == 0 {
				res = Perform(ActVerb, nil, nil)
				DirObj = nil
			} else if !Lit {
				Print("It's not clear what you're referring to.", Newline)
				res = PerfNotHndld
			}
		} else {
			tmp := false
			notHere := 0
			var obj1 *Object
			var indir, dir *Object
			for i := 0; i < numObj; i++ {
				if isDir {
					obj1 = DirObjPossibles[i]
					dir, indir = obj1, obj
				} else {
					obj1 = IndirObjPossibles[i]
					dir, indir = obj, obj1
				}
				if numObj > 1 || (len(ParsedSyntx.ObjOrClause1) > 0 && ParsedSyntx.ObjOrClause1[0].Is("all")) {
					if dir == &NotHereObject {
						notHere++
						continue
					}
					if ActVerb.Norm == "take" &&
						indir != nil &&
						len(ParsedSyntx.ObjOrClause1) > 0 &&
						ParsedSyntx.ObjOrClause1[0].Is("all") &&
						dir != nil &&
						!dir.IsIn(indir) {
						continue
					}
					if l := dir.Location(); Params.GetType == GetAll &&
						ActVerb.Norm == "take" &&
						((l != Winner &&
							l != Here &&
							l != Winner.Location() &&
							l != indir &&
							!l.Has(FlgSurf)) ||
							!(dir.Has(FlgTake) || dir.Has(FlgTryTake))) {
						continue
					}
					if obj1 == &It {
						PrintObject(Params.ItObj)
					} else {
						PrintObject(obj1)
					}
					Print(": ", NoNewline)
				}
				DirObj, IndirObj = dir, indir
				tmp = true
				res = Perform(ActVerb, DirObj, IndirObj)
				if res == PerfFatal {
					break
				}
			}
			if notHere > 0 {
				Print("The ", NoNewline)
				if notHere != numObj {
					Print("other ", NoNewline)
				}
				Print("object", NoNewline)
				if notHere != 1 {
					Print("s", NoNewline)
				}
				Print(" that you mentioned ", NoNewline)
				if notHere != 1 {
					Print("are", NoNewline)
				} else {
					Print("is", NoNewline)
				}
				Print("n't here.", Newline)
			} else if !tmp {
				Print("There's nothing here you can take.", Newline)
			}
		}
		if l := Winner.Location(); res != PerfFatal && l != nil && l.Action != nil {
			if l.Action(ActEnd) {
				res = PerfHndld
			} else {
				res = PerfNotHndld
			}
		}
		if res == PerfFatal {
			Params.Continue = -1
		}
		if ActVerb.Norm == "tell" || ActVerb.Norm == "brief" || ActVerb.Norm == "superbrief" || ActVerb.Norm == "verbose" || ActVerb.Norm == "save" || ActVerb.Norm == "version" || ActVerb.Norm == "quit" || ActVerb.Norm == "restart" || ActVerb.Norm == "score" || ActVerb.Norm == "script" || ActVerb.Norm == "unscript" || ActVerb.Norm == "restore" {
			continue
		} else {
			Clocker()
		}
	}
}

func Perform(a ActionVerb, o, i *Object) PerfRet {
	if (o == &It || i == &It) && IsAccessible(Params.ItObj) {
		Print("I don't see what you are referring to.", Newline)
		return PerfFatal
	}
	if o == &It {
		o = Params.ItObj
	}
	if i == &It {
		i = Params.ItObj
	}
	if o != nil && IndirObj != &It && a.Norm == "walk" {
		Params.ItObj = o
	}
	if o == &NotHereObject || i == &NotHereObject {
		if NotHereObjectFcn(ActUnk) {
			return PerfHndld
		}
	}
	if Winner != nil && Winner.Action != nil {
		if Winner.Action(ActUnk) {
			return PerfHndld
		}
	}
	if l := Winner.Location(); l != nil && l.Action != nil {
		if l.Action(ActBegin) {
			return PerfHndld
		}
	}
	if act, ok := PreActions[a.Orig]; ok && act != nil {
		if act(ActUnk) {
			return PerfHndld
		}
	} else if a.Norm != a.Orig {
		if act, ok := PreActions[a.Norm]; ok && act != nil {
			if act(ActUnk) {
				return PerfHndld
			}
		}
	}
	if i != nil && i.Action != nil {
		if i.Action(ActUnk) {
			return PerfHndld
		}
	}
	if o != nil && a.Norm != "walk" && o.Location() != nil && o.Location().ContFcn != nil {
		if o.Location().ContFcn(ActUnk) {
			return PerfHndld
		}
	}
	if o != nil && a.Norm != "walk" && o.Action != nil {
		if o.Action(ActUnk) {
			return PerfHndld
		}
	}
	// Try the action by Orig first, then by Norm (for normalized verbs
	// like "turn on" -> "lamp on")
	if act, ok := Actions[a.Orig]; ok && act != nil {
		if act(ActUnk) {
			return PerfHndld
		}
	} else if a.Norm != a.Orig {
		if act, ok := Actions[a.Norm]; ok && act != nil {
			if act(ActUnk) {
				return PerfHndld
			}
		}
	}
	return PerfNotHndld
}
