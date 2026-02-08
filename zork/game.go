package zork

import "os"

type PerfRet int

const (
	PerfNotHndld PerfRet = iota
	PerfHndld
	PerfFatal
)

// ErrQuit is a sentinel value used to signal game exit via panic/recover.
// This allows tests to catch quit without os.Exit terminating the process.
var ErrQuit = "game-quit"

// RFatal signals a fatal action result (equivalent to ZIL's <RFATAL>).
// It sets G.PerformFatal, which makes Perform return PerfFatal.
// Call this in an action handler and return its result to stop the
// Perform chain and the multi-object loop in the main loop.
func RFatal() bool {
	G.PerformFatal = true
	return true
}

func Quit() {
	panic(ErrQuit)
}

// Verify checks the integrity of the game file.
// In the original Z-machine, VERIFY computes a checksum of the story file.
// This is not applicable to the Go port; we always return true.
func Verify() bool {
	return true
}

// ResetGameState creates a fresh GameState with all defaults and restores the
// object tree. Tests and restart use this to get a clean game.
func ResetGameState() {
	// Preserve I/O handles if they were set (e.g. by tests)
	var out, in_ = G.GameOutput, G.GameInput
	G = NewGameState()
	if out != nil {
		G.GameOutput = out
	}
	if in_ != nil {
		G.GameInput = in_
	}
	// Restore object tree to initial state
	ResetObjectTree()
}

// InitGame sets up all game state for a fresh game. Call once before MainLoop.
func InitGame() {
	if G == nil {
		G = NewGameState()
	}
	ResetGameState()
	initClockFuncs()
	FinalizeGameObjects()
	BuildObjectTree()
	BuildVocabulary()
	InitReader()
	initSaveSystem()

	Queue("IFight", -1).Run = true
	Queue("ISword", -1)
	Queue("IThief", -1).Run = true
	Queue("ICandles", 40)
	Queue("ILantern", 200)
	InflatedBoat.VehType = FlgNonLand
	Def1Res[1] = Def1[2]
	Def1Res[2] = Def1[4]
	Def2Res[2] = Def2B[2]
	Def2Res[3] = Def2B[4]
	Def3Res[1] = Def3A[2]
	Def3Res[3] = Def3B[2]
	G.Here = &WestOfHouse
	ThisIsIt(&Mailbox)
	G.Lit = true
	G.Winner = &Adventurer
	G.Player = G.Winner
	G.Winner.MoveTo(G.Here)
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
	if !G.Here.Has(FlgTouch) {
		VVersion(ActUnk)
		NewLine()
	}
	VLook(ActUnk)
	MainLoop()
}

func MainLoop() {
	G.Params.Continue = NumUndef
	for {
		if G.InputExhausted {
			return
		}
		G.ParserOk = Parse()
		if G.InputExhausted {
			return
		}
		if !G.ParserOk {
			G.Params.Continue = NumUndef
			continue
		}
		if G.Params.ItObj != nil && IsAccessible(G.Params.ItObj) {
			found := false
			for idx, obj := range G.IndirObjPossibles {
				if obj == &It {
					G.IndirObjPossibles[idx] = G.Params.ItObj
					found = true
				}
			}
			if !found {
				for idx, obj := range G.DirObjPossibles {
					if obj == &It {
						G.DirObjPossibles[idx] = G.Params.ItObj
					}
				}
			}
		}
		numObj := 1
		var obj *Object
		isDir := true
		if len(G.DirObjPossibles) == 0 {
			numObj = 0
		} else if len(G.DirObjPossibles) > 1 {
			if len(G.IndirObjPossibles) == 0 {
				obj = nil
			} else {
				obj = G.IndirObjPossibles[0]
			}
			numObj = len(G.DirObjPossibles)
		} else if len(G.IndirObjPossibles) > 1 {
			isDir = false
			obj = G.DirObjPossibles[0]
			numObj = len(G.IndirObjPossibles)
		}
		if obj == nil && len(G.IndirObjPossibles) == 1 {
			obj = G.IndirObjPossibles[0]
		}
		var res PerfRet
		if G.ActVerb.Norm == "walk" && len(G.Params.WalkDir) != 0 {
			res = Perform(G.ActVerb, G.DirObj, nil)
		} else if numObj == 0 {
			if G.DetectedSyntx.NumObjects() == 0 {
				res = Perform(G.ActVerb, nil, nil)
				G.DirObj = nil
			} else if !G.Lit {
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
					obj1 = G.DirObjPossibles[i]
					dir, indir = obj1, obj
				} else {
					obj1 = G.IndirObjPossibles[i]
					dir, indir = obj, obj1
				}
				if numObj > 1 || (len(G.ParsedSyntx.ObjOrClause1) > 0 && G.ParsedSyntx.ObjOrClause1[0].Is("all")) {
					if dir == &NotHereObject {
						notHere++
						continue
					}
					if G.ActVerb.Norm == "take" &&
						indir != nil &&
						len(G.ParsedSyntx.ObjOrClause1) > 0 &&
						G.ParsedSyntx.ObjOrClause1[0].Is("all") &&
						dir != nil &&
						!dir.IsIn(indir) {
						continue
					}
					if l := dir.Location(); G.Params.GetType == GetAll &&
						G.ActVerb.Norm == "take" &&
						((l != G.Winner &&
							l != G.Here &&
							l != G.Winner.Location() &&
							l != indir &&
							!l.Has(FlgSurf)) ||
							!(dir.Has(FlgTake) || dir.Has(FlgTryTake))) {
						continue
					}
					if obj1 == &It {
						PrintObject(G.Params.ItObj)
					} else {
						PrintObject(obj1)
					}
					Print(": ", NoNewline)
				}
				G.DirObj, G.IndirObj = dir, indir
				tmp = true
				res = Perform(G.ActVerb, G.DirObj, G.IndirObj)
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
		if l := G.Winner.Location(); res != PerfFatal && l != nil && l.Action != nil {
			if l.Action(ActEnd) {
				res = PerfHndld
			} else {
				res = PerfNotHndld
			}
		}
		if res == PerfFatal {
			G.Params.Continue = -1
		}
		if G.ActVerb.Norm == "tell" || G.ActVerb.Norm == "brief" || G.ActVerb.Norm == "superbrief" || G.ActVerb.Norm == "verbose" || G.ActVerb.Norm == "save" || G.ActVerb.Norm == "version" || G.ActVerb.Norm == "quit" || G.ActVerb.Norm == "restart" || G.ActVerb.Norm == "score" || G.ActVerb.Norm == "script" || G.ActVerb.Norm == "unscript" || G.ActVerb.Norm == "restore" {
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
	if G.PerformFatal {
		G.PerformFatal = false
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
	oldActVerb := G.ActVerb
	oldDirObj := G.DirObj
	oldIndirObj := G.IndirObj
	defer func() {
		G.ActVerb = oldActVerb
		G.DirObj = oldDirObj
		G.IndirObj = oldIndirObj
	}()

	// Clear any stale fatal flag (e.g. from parser's ITake call).
	G.PerformFatal = false
	if (o == &It || i == &It) && !IsAccessible(G.Params.ItObj) {
		Print("I don't see what you are referring to.", Newline)
		return PerfFatal
	}
	if o == &It {
		o = G.Params.ItObj
	}
	if i == &It {
		i = G.Params.ItObj
	}
	// Set globals from parameters (like ZIL's SETG PRSA/PRSO/PRSI).
	G.ActVerb = a
	G.DirObj = o
	// Track "it" for non-walk commands (ZIL: <NOT <VERB? WALK>>).
	if o != nil && a.Norm != "walk" {
		G.Params.ItObj = o
	}
	G.IndirObj = i

	if o == &NotHereObject || i == &NotHereObject {
		if ret, done := callHandler(NotHereObjectFcn, ActUnk); done {
			return ret
		}
	}
	if G.Winner != nil && G.Winner.Action != nil {
		if ret, done := callHandler(G.Winner.Action, ActUnk); done {
			return ret
		}
	}
	if l := G.Winner.Location(); l != nil && l.Action != nil {
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
