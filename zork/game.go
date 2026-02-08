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

// PerformFatal is set by action handlers to signal the equivalent of ZIL's
// RFATAL return. When set, the Perform function returns PerfFatal instead
// of PerfHndld, which causes the main loop to stop the multi-object loop
// and clear the continuation flag.
var PerformFatal bool

// RFatal signals a fatal action result (equivalent to ZIL's <RFATAL>).
// Call this in an action handler and return its result to stop the
// Perform chain and the multi-object loop in the main loop.
func RFatal() bool {
	PerformFatal = true
	return true
}

// Save, Restore, and Restart are function variables that start as stubs.
// They are replaced with real implementations in InitGame (from save.go)
// to break the init cycle between package-level variable initializers and
// the Objects slice.
var (
	Save    = func() bool { return false }
	Restore = func() bool { return false }
	Restart = func() bool { return false }
)

func Quit() {
	panic(ErrQuit)
}

// Verify checks the integrity of the game file.
// In the original Z-machine, VERIFY computes a checksum of the story file.
// This is not applicable to the Go port; we always return true.
func Verify() bool {
	return true
}

// ResetGameState resets all mutable global state so a fresh game can be started.
// This is essential for tests that run multiple games in one process.
func ResetGameState() {
	// Reset clock
	QueueItms = [30]QueueItm{}
	QueueInts = 30
	QueueDmns = 30
	ClockWait = false

	// Reset perform state
	PerformFatal = false

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
	initSaveSystem()

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
				Print("It's too dark to see.", Newline)
				res = PerfNotHndld
			} else {
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

// callHandler invokes an action handler and checks for a fatal signal.
// Returns the appropriate PerfRet and whether the chain should stop.
func callHandler(fn func(ActArg) bool, arg ActArg) (PerfRet, bool) {
	result := fn(arg)
	if PerformFatal {
		PerformFatal = false
		return PerfFatal, true
	}
	if result {
		return PerfHndld, true
	}
	return PerfNotHndld, false
}

func Perform(a ActionVerb, o, i *Object) PerfRet {
	// Save old globals so nested Perform calls don't corrupt the outer
	// call's state. This mirrors ZIL's save/restore of PRSA, PRSO, PRSI.
	oldActVerb := ActVerb
	oldDirObj := DirObj
	oldIndirObj := IndirObj
	defer func() {
		ActVerb = oldActVerb
		DirObj = oldDirObj
		IndirObj = oldIndirObj
	}()

	// Clear any stale fatal flag (e.g. from parser's ITake call).
	PerformFatal = false
	if (o == &It || i == &It) && !IsAccessible(Params.ItObj) {
		Print("I don't see what you are referring to.", Newline)
		return PerfFatal
	}
	if o == &It {
		o = Params.ItObj
	}
	if i == &It {
		i = Params.ItObj
	}
	// Set globals from parameters (like ZIL's SETG PRSA/PRSO/PRSI).
	ActVerb = a
	DirObj = o
	// Track "it" for non-walk commands (ZIL: <NOT <VERB? WALK>>).
	if o != nil && a.Norm != "walk" {
		Params.ItObj = o
	}
	IndirObj = i

	if o == &NotHereObject || i == &NotHereObject {
		if ret, done := callHandler(NotHereObjectFcn, ActUnk); done {
			return ret
		}
	}
	if Winner != nil && Winner.Action != nil {
		if ret, done := callHandler(Winner.Action, ActUnk); done {
			return ret
		}
	}
	if l := Winner.Location(); l != nil && l.Action != nil {
		if ret, done := callHandler(l.Action, ActBegin); done {
			return ret
		}
	}
	if act, ok := PreActions[a.Orig]; ok && act != nil {
		if ret, done := callHandler(act, ActUnk); done {
			return ret
		}
	} else if a.Norm != a.Orig {
		if act, ok := PreActions[a.Norm]; ok && act != nil {
			if ret, done := callHandler(act, ActUnk); done {
				return ret
			}
		}
	}
	if i != nil && i.Action != nil {
		if ret, done := callHandler(i.Action, ActUnk); done {
			return ret
		}
	}
	if o != nil && a.Norm != "walk" && o.Location() != nil && o.Location().ContFcn != nil {
		if ret, done := callHandler(o.Location().ContFcn, ActUnk); done {
			return ret
		}
	}
	if o != nil && a.Norm != "walk" && o.Action != nil {
		if ret, done := callHandler(o.Action, ActUnk); done {
			return ret
		}
	}
	// Try the action by Orig first, then by Norm (for normalized verbs
	// like "turn on" -> "lamp on")
	if act, ok := Actions[a.Orig]; ok && act != nil {
		if ret, done := callHandler(act, ActUnk); done {
			return ret
		}
	} else if a.Norm != a.Orig {
		if act, ok := Actions[a.Norm]; ok && act != nil {
			if ret, done := callHandler(act, ActUnk); done {
				return ret
			}
		}
	}
	return PerfNotHndld
}
